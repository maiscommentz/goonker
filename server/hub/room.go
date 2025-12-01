package hub

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"Goonker/common"
	"Goonker/server/logic"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Player struct {
	Conn *websocket.Conn
	ID   common.PlayerID
}

type Room struct {
	ID        string
	Players   map[common.PlayerID]*Player
	Logic     *logic.GameLogic
	
	mutex     sync.Mutex
	IsBotGame bool
}

func NewRoom(id string, isBot bool) *Room {
	return &Room{
		ID:        id,
		Players:   make(map[common.PlayerID]*Player),
		Logic:     logic.NewGameLogic(),
		IsBotGame: isBot,
	}
}

// AddPlayer add a new player to the room and starts listening to their messages
func (r *Room) AddPlayer(conn *websocket.Conn) common.PlayerID {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Determine the player ID (1 or 2)
	var pid common.PlayerID
	if _, ok := r.Players[common.P1]; !ok {
		pid = common.P1
	} else if _, ok := r.Players[common.P2]; !ok {
		pid = common.P2
	} else {
		return common.Empty // Room full
	}

	r.Players[pid] = &Player{Conn: conn, ID: pid}
	
	// Launch listener in a separate goroutine
	go r.listenPlayer(pid, conn)

	// If the room is full or if it's a bot game, start
	if r.IsFull() || (r.IsBotGame && pid == common.P1) {
		go r.startGame()
	}

	return pid
}

func (r *Room) IsFull() bool {
	if r.IsBotGame {
		return len(r.Players) >= 1
	}
	return len(r.Players) == 2
}

func (r *Room) startGame() {
	log.Printf("Room %s: Starting game", r.ID)
	r.broadcastGameStart()
	r.broadcastUpdate()
}

// listenPlayer listens to incoming messages from a specific client
func (r *Room) listenPlayer(pid common.PlayerID, conn *websocket.Conn) {
	ctx := context.Background()
	defer func() {
		// Cleanup on disconnection
		r.mutex.Lock()
		delete(r.Players, pid)
		r.mutex.Unlock()
		conn.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		var packet common.Packet
		err := wsjson.Read(ctx, conn, &packet)
		if err != nil {
			log.Printf("Room %s: Player %d disconnected", r.ID, pid)
			return
		}

		if packet.Type == common.MsgClick {
			var payload common.ClickPayload
			if err := json.Unmarshal(packet.Data, &payload); err == nil {
				r.handleMove(pid, payload.X, payload.Y)
			}
		}
	}
}

func (r *Room) handleMove(pid common.PlayerID, x, y int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Apply the move via pure game logic
	err := r.Logic.ApplyMove(pid, x, y)
	if err != nil {
		log.Printf("Invalid move from %d: %v", pid, err)
		return
	}

	// Send the update to everyone
	r.broadcastUpdate_Locked()

	// If it's a Bot Game and the game is not over, the bot plays
	if r.IsBotGame && !r.Logic.GameOver && r.Logic.Turn == common.P2 {
		go func() {
			// Launch the bot in a goroutine to avoid blocking the mutex for too long
			bx, by := logic.GetBotMove(r.Logic)
			if bx != -1 {
				// Call handleMove for the bot
				// Note: handleMove takes a Lock, so we must call it outside the current lock.
				// That's why we're in a `go func` here.
				r.handleMove(common.P2, bx, by)
			}
		}()
	}
}

func (r *Room) broadcastGameStart() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for pid, p := range r.Players {
		payload := common.GameStartPayload{
			YouAre: pid,
			OpponentID: "unknown",
		}
		data, _ := json.Marshal(payload)
		packet := common.Packet{Type: common.MsgGameStart, Data: data}
		
		go wsjson.Write(context.Background(), p.Conn, packet)
	}
}

func (r *Room) broadcastUpdate() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.broadcastUpdate_Locked()
}

func (r *Room) broadcastUpdate_Locked() {
	log.Printf("Room %s: Broadcasting update", r.ID)
	r.Logic.PrintConsoleBoard()
	payload := common.UpdatePayload{
		Board: r.Logic.Board,
		Turn:  r.Logic.Turn,
	}
	data, _ := json.Marshal(payload)
	packet := common.Packet{Type: common.MsgUpdate, Data: data}

	for _, p := range r.Players {
		// Timeout of 5 seconds for sending
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		wsjson.Write(ctx, p.Conn, packet)
		cancel()
	}
}