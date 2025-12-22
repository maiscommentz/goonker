package common

import "encoding/json"

// Represents a player in the game
type PlayerID byte

// PlayerID constants
const (
	Empty PlayerID = 0
	P1    PlayerID = 1 // X
	P2    PlayerID = 2 // O
)

// Game constants
const (
	BoardSize = 3
)

// Message types
const (
	MsgJoin      = "join"       // Client -> Server: "I want to join room X"
	MsgGetRooms  = "get_rooms"  // Client -> Server: "Get available rooms"
	MsgRooms     = "rooms"      // Server -> Client: "Available rooms"
	MsgGameStart = "game_start" // Server -> Client: "Match found, you are X"
	MsgClick     = "click"      // Client -> Server: "I clicked cell 4"
	MsgUpdate    = "update"     // Server -> Client: "New board state"
	MsgGameOver  = "game_over"  // Server -> Client: "Game over, result is X"
)

// Packet is the generic message structure for communication.
type Packet struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

// GameStartPayload is sent by server to notify game start.
type GameStartPayload struct {
	YouAre PlayerID `json:"you_are"` // 1 or 2
}

// ClickPayload is sent by client with (x,y) of clicked cell.
type ClickPayload struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// UpdatePayload is sent by server to sync the board.
type UpdatePayload struct {
	Board [BoardSize][BoardSize]PlayerID `json:"board"`
	Turn  PlayerID                       `json:"turn"` // Whose turn is it?
}

// JoinPayload is sent by client to join a room.
type JoinPayload struct {
	RoomID string `json:"room_id"`
	IsBot  bool   `json:"is_bot"` // Whether to play against a bot
}

// GameOverPayload is sent by server when game ends.
type GameOverPayload struct {
	Winner PlayerID `json:"winner"` // Who won? 0 for draw, 1 or 2 for players
}

// RoomsPayload is sent by server to notify available rooms.
type RoomsPayload struct {
	Rooms []string `json:"rooms"`
}
