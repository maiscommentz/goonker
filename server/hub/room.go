package hub

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"Goonker/common"
	"Goonker/server/logic"
	"Goonker/server/utils"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Hub configuration constants
const (
	WriteTimeout      = 5 * time.Second
	CloseMessage      = "Goodbye"
	MaxPlayers        = 2
	MaxPlayersWithBot = 1
)

// Player represents a connected player in the room
type Player struct {
	Conn *websocket.Conn
	ID   common.PlayerID
}

// Room represents a game room with players and game logic
type Room struct {
	ID      string
	Players map[common.PlayerID]*Player
	Logic   *logic.GameLogic

	mutex     sync.Mutex
	IsBotGame bool

	// Challenge
	challengedMove     common.ClickPayload
	challengeAnswerKey int
	challengedPlayer   common.PlayerID
	challengeTimer     *time.Timer
}

// NewRoom creates a new Room instance.
func NewRoom(id string, isBot bool) *Room {
	return &Room{
		ID:        id,
		Players:   make(map[common.PlayerID]*Player),
		Logic:     logic.NewGameLogic(),
		IsBotGame: isBot,
	}
}

// AddPlayer assigns an ID (P1/P2) to the connecting player and starts listening.
func (r *Room) AddPlayer(conn *websocket.Conn) common.PlayerID {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Determine the player ID (1 or 2)
	var pid common.PlayerID
	if _, ok := r.Players[common.P1]; !ok {
		pid = common.P1
	} else if r.IsBotGame {
		return common.Empty // Bot game already has player 1
	} else if _, ok := r.Players[common.P2]; !ok {
		pid = common.P2
	} else {
		return common.Empty // Room full
	}

	r.Players[pid] = &Player{Conn: conn, ID: pid}

	// Start listening to this client on a separate goroutine
	go r.listenPlayer(pid, conn)

	// Check if the game is ready to start
	if r.IsFull() {
		go r.startGame()
	}

	return pid
}

// IsFull checks if the room has enough players to start the game.
// In Bot games, only 1 player is needed, otherwise 2 players are required.
func (r *Room) IsFull() bool {
	if r.IsBotGame {
		return len(r.Players) >= MaxPlayersWithBot
	}
	return len(r.Players) == MaxPlayers
}

// startGame initializes the game and notifies players.
func (r *Room) startGame() {
	log.Printf("Room %s: Starting game", r.ID)
	r.broadcastGameStart()
	r.broadcastUpdate()
}

// listenPlayer listens to incoming messages from a specific client.
// It manages the connection lifecycle and handles disconnections.
func (r *Room) listenPlayer(pid common.PlayerID, conn *websocket.Conn) {
	ctx := context.Background()

	// Cleanup triggers on function exit (connection closed or error)
	defer func() {
		r.mutex.Lock()
		delete(r.Players, pid)
		r.mutex.Unlock()

		err := conn.Close(websocket.StatusNormalClosure, CloseMessage)
		if err != nil {
			log.Println(err)
		}

		// Auto-remove room if empty
		if len(r.Players) == 0 {
			GlobalHub.RemoveRoom(r.ID)
			log.Printf("Room %s: All players disconnected, room removed", r.ID)
		} else {
			log.Printf("Room %s: Player %d disconnected, waiting for new player", r.ID, pid)
		}
	}()

	// Listen for incoming messages
	for {
		var packet common.Packet
		err := wsjson.Read(ctx, conn, &packet)
		if err != nil {
			return
		}

		// Handle Click messages
		switch packet.Type {
		case common.MsgClick:
			var payload common.ClickPayload
			if err := json.Unmarshal(packet.Data, &payload); err == nil {
				if r.Logic.ShouldTriggerChallenge(pid, payload.X, payload.Y) {
					r.challengedMove = payload
					r.challengedPlayer = pid
					r.startChallenge(conn)
				} else {
					r.handleMove(pid, payload.X, payload.Y)
				}
			}
		case common.MsgGetRooms:
			r.sendRooms(conn)
		case common.MsgAnswer:
			var payload common.AnswerPayload
			if err := json.Unmarshal(packet.Data, &payload); err == nil {
				if r.challengeTimer != nil {
					r.challengeTimer.Stop()
				}
				if payload.Answer == r.challengeAnswerKey {
					log.Println("Challenge completed successfully")
					r.Logic.DeleteMove(r.challengedMove.X, r.challengedMove.Y)
					// Play the move
					r.handleMove(pid, r.challengedMove.X, r.challengedMove.Y)
				} else {
					// Play the move, but no challenge this time
					r.handleMove(pid, r.challengedMove.X, r.challengedMove.Y)
					log.Println("Challenge failed")
				}
			}
		default:
			log.Printf("Room %s: Unknown message type: %s", r.ID, packet.Type)
		}
	}
}

// sendRooms sends the available rooms to the client.
func (r *Room) sendRooms(conn *websocket.Conn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	rooms := GlobalHub.GetAvailableRooms()
	payload := common.RoomsPayload{Rooms: rooms}
	r.sendJson(conn, common.MsgRooms, payload)
}

func (r *Room) startChallenge(conn *websocket.Conn) {
	log.Println("Start challenge")
	r.mutex.Lock()
	defer r.mutex.Unlock()

	challenge, err := utils.PickChallenge()
	if err != nil {
		log.Fatal(err)
		return
	}

	payload := common.ChallengePayload{Question: challenge.Question, Answers: challenge.Answers}
	r.challengeAnswerKey = challenge.AnswerKey
	r.sendJson(conn, common.MsgChallenge, payload)

	r.challengeTimer = time.AfterFunc(common.ChallengeTime*time.Second, func() {
		r.handleChallengeTimeout()
	})
}

func (r *Room) handleChallengeTimeout() {
	log.Println("Challenge time ran out")
	r.handleMove(r.challengedPlayer, r.challengedMove.X, r.challengedMove.Y)
}

// handleMove coordinates game logic updates and notifications. Returns true if a challenge must start.
func (r *Room) handleMove(pid common.PlayerID, x, y int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Apply the move via pure game logic
	err := r.Logic.ApplyMove(pid, x, y)
	if err != nil {
		log.Printf("Invalid move from %d: %v", pid, err)
	}

	// Send the updated board state to all players
	r.broadcastUpdate_Locked()

	// Check if game is over, broadcast the result
	if r.Logic.GameOver {
		r.broadcastGameOver()
	}

	// If it's a Bot Game and the game is not over, the bot plays.
	// Launch the bot in a goroutine to avoid blocking the mutex for too long.
	if r.IsBotGame && !r.Logic.GameOver && r.Logic.Turn == common.P2 {
		// Take a snapshot of the current game logic
		logicSnapshot := r.Logic
		go func(snapshot *logic.GameLogic) {
			botX, botY := logic.GetBotMove(snapshot)
			if botX != logic.InvalidCoord {
				// Valid move returned
				r.handleMove(common.P2, botX, botY)
			}
		}(logicSnapshot) // Pass a snapshot to avoid race conditions
	}
}

// broadcastGameStart notifies all players that the game is starting.
func (r *Room) broadcastGameStart() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for pid, p := range r.Players {
		payload := common.GameStartPayload{
			YouAre: pid,
		}
		r.sendJson(p.Conn, common.MsgGameStart, payload)
	}
}

// broadcastUpdate sends the current game state to all players.
func (r *Room) broadcastUpdate() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.broadcastUpdate_Locked()
}

// broadcastUpdate_Locked sends the current game state to all players.
func (r *Room) broadcastUpdate_Locked() {
	log.Printf("Room %s: Broadcasting update", r.ID)
	r.Logic.PrintConsoleBoard()
	payload := common.UpdatePayload{
		Board: r.Logic.Board,
		Turn:  r.Logic.Turn,
	}

	for _, p := range r.Players {
		r.sendJson(p.Conn, common.MsgUpdate, payload)
	}
}

// broadcastGameOver notifies all players that the game has ended.
func (r *Room) broadcastGameOver() {
	log.Printf("Room %s: Broadcasting game over", r.ID)
	payload := common.GameOverPayload{
		Winner: r.Logic.Winner,
	}

	for _, p := range r.Players {
		r.sendJson(p.Conn, common.MsgGameOver, payload)
	}
}

// sendJson helps to reduce boilerplate and enforce timeouts
func (r *Room) sendJson(c *websocket.Conn, msgType string, payload interface{}) {
	data, _ := json.Marshal(payload)
	packet := common.Packet{Type: msgType, Data: data}

	ctx, cancel := context.WithTimeout(context.Background(), WriteTimeout)
	defer cancel()

	if err := wsjson.Write(ctx, c, packet); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}
