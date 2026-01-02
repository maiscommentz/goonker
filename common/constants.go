package common

// Represents a player in the game
type PlayerID byte

const (
	// Game constants
	BoardSize     = 3
	ChallengeTime = 8

	// PlayerID constants
	Empty PlayerID = 0
	P1    PlayerID = 1 // X
	P2    PlayerID = 2 // O
)
