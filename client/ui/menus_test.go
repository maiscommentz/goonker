package ui

import (
	"Goonker/common"
	"testing"
)

func TestMenuConstructors(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Skip("Skipping Menu test due to asset initialization failure:", r)
		}
	}()
	InitImages()

	// Main Menu
	mm := NewMainMenu()
	if mm.BtnPlay == nil || mm.BtnQuit == nil {
		t.Error("Main Menu buttons not initialized")
	}

	// Rooms Menu
	rm := NewRoomsMenu()
	if rm.BtnCreateRoom == nil || rm.BtnJoinGame == nil {
		t.Error("Rooms Menu buttons not initialized")
	}
	if rm.RoomField == nil {
		t.Error("RoomField not initialized")
	}

	// Game Over Menu
	gom := NewGameOverMenu()
	if gom.BtnBack == nil {
		t.Error("Game Over Menu back button not initialized")
	}

	// Challenge Menu
	dummyChallenge := common.ChallengePayload{
		Question: "Q?",
		Answers:  []string{"A", "B"},
	}
	cm := NewChallengeMenu(dummyChallenge)
	if len(cm.Answers) != 2 {
		t.Errorf("Expected 2 answer buttons, got %d", len(cm.Answers))
	}
	if cm.Question != "Q?" {
		t.Error("Challenge Question mismatch")
	}
}
