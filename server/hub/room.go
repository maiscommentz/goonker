package hub

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"Goonker/common" // Adjust to your module path

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Player wraps the connection
type Player struct {
	Conn   *websocket.Conn
	Symbol common.PlayerID // 1 (X) or 2 (O)
}

// Room represents one match
type Room struct {
	ID string

	// Game State
	mu    sync.Mutex
	board [3][3]common.PlayerID
	turn  common.PlayerID

	// Players
	p1 *Player
	p2 *Player
}

func NewRoom(id string) *Room {
	return &Room{
		ID:    id,
		board: [3][3]common.PlayerID{}, // Empty board
		turn:  common.P1,               // X starts
	}
}

// Join handles the player connection lifecycle
func (r *Room) Join(conn *websocket.Conn, ctx context.Context) {
	r.mu.Lock()

	// Assign Player Slot
	var player *Player
	var symbol common.PlayerID

	if r.p1 == nil {
		symbol = common.P1
		player = &Player{Conn: conn, Symbol: symbol}
		r.p1 = player
		log.Printf("[%s] Player 1 joined", r.ID)
	} else if r.p2 == nil {
		symbol = common.P2
		player = &Player{Conn: conn, Symbol: symbol}
		r.p2 = player
		log.Printf("[%s] Player 2 joined", r.ID)
	} else {
		r.mu.Unlock()
		conn.Close(websocket.StatusPolicyViolation, "Room full")
		return
	}
	r.mu.Unlock()

	// Check if game can start
	r.checkStart()

	// Listen Loop (Blocks)
	r.listen(player, ctx)

	// Cleanup on disconnect
	r.mu.Lock()
	if r.p1 == player {
		r.p1 = nil
	} else if r.p2 == player {
		r.p2 = nil
	}
	r.mu.Unlock()
	log.Printf("[%s] Player %d disconnected", r.ID, symbol)
}

func (r *Room) checkStart() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.p1 != nil && r.p2 != nil {
		log.Printf("[%s] Both players ready. Starting game.", r.ID)
		
		// Send "Game Start" to P1
		r.sendPacket(r.p1, common.Packet{
			Type: common.MsgGameStart,
			Data: mustMarshal(common.GameStartPayload{YouAre: common.P1}),
		})

		// Send "Game Start" to P2
		r.sendPacket(r.p2, common.Packet{
			Type: common.MsgGameStart,
			Data: mustMarshal(common.GameStartPayload{YouAre: common.P2}),
		})

		// Broadcast initial board
		r.broadcastUpdate()
	}
}

func (r *Room) listen(player *Player, ctx context.Context) {
	for {
		var packet common.Packet
		err := wsjson.Read(ctx, player.Conn, &packet)
		if err != nil {
			return // Connection closed
		}

		// Handle Incoming Packet
		r.handlePacket(player, packet)
	}
}

func (r *Room) handlePacket(player *Player, packet common.Packet) {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch packet.Type {
	case common.MsgClick:
		// 1. Validate Turn
		if r.turn != player.Symbol {
			log.Printf("[%s] Ignored move from %d (not their turn)", r.ID, player.Symbol)
			return
		}

		// 2. Parse Payload
		var payload common.ClickPayload
		if err := json.Unmarshal(packet.Data, &payload); err != nil {
			return
		}

		// 3. Validate Move (Bounds & Empty)
		if payload.X < 0 || payload.X > 2 || payload.Y < 0 || payload.Y > 2 {
			return
		}
		if r.board[payload.X][payload.Y] != common.Empty {
			// CONQUER MECHANIC WOULD GO HERE (Check if enemy occupied)
			log.Printf("Cell occupied. Trigger minigame logic here later.")
			return
		}

		// 4. Update Board
		r.board[payload.X][payload.Y] = player.Symbol
		
		// 5. Toggle Turn
		if r.turn == common.P1 {
			r.turn = common.P2
		} else {
			r.turn = common.P1
		}

		// 6. Broadcast new state
		r.broadcastUpdate()
	}
}

func (r *Room) broadcastUpdate() {
	// Prepare Payload
	update := common.UpdatePayload{
		Board: r.board,
		Turn:  r.turn,
	}
	data := mustMarshal(update)

	pkt := common.Packet{
		Type: common.MsgUpdate,
		Data: data,
	}

	if r.p1 != nil {
		r.sendPacket(r.p1, pkt)
	}
	if r.p2 != nil {
		r.sendPacket(r.p2, pkt)
	}
}

func (r *Room) sendPacket(p *Player, pkt common.Packet) {
	// In a real app, use a context with timeout
	go func() {
		wsjson.Write(context.Background(), p.Conn, pkt)
	}()
}

// Helper to marshal JSON without error checking (for internal structs)
func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

// HasPlayers checks if there are any connected players
func (r *Room) HasPlayers() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.p1 != nil || r.p2 != nil
}