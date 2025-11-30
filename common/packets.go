package common

import "encoding/json"

type PlayerID int

const (
	Empty PlayerID = 0
	P1    PlayerID = 1 // X
	P2    PlayerID = 2 // O
)

// Message Types
const (
	MsgJoin      = "join"       // Client -> Server: "I want to play"
	MsgGameStart = "game_start" // Server -> Client: "Match found, you are X"
	MsgClick     = "click"      // Client -> Server: "I clicked cell 4"
	MsgUpdate    = "update"     // Server -> Client: "New board state"
)

// Packet is the generic container sent over the network
type Packet struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"` // The payload specific to the Type
}

// GameStartPayload: Sent by server when match starts
type GameStartPayload struct {
	YouAre     PlayerID `json:"you_are"`     // 1 or 2
	OpponentID string   `json:"opponent_id"` // For future use
}

// ClickPayload: Sent by client when clicking the board
type ClickPayload struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// UpdatePayload: Sent by server to sync the board
type UpdatePayload struct {
	Board [3][3]PlayerID `json:"board"`
	Turn  PlayerID       `json:"turn"` // Whose turn is it?
}