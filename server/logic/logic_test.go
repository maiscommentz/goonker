package logic

import (
	"Goonker/common"
	"testing"
)

func TestNewGameLogic(t *testing.T) {
	game := NewGameLogic()
	if game.Turn != common.P1 {
		t.Errorf("Expected Turn to be P1, got %d", game.Turn)
	}
	if game.GameOver {
		t.Error("Expected GameOver to be false")
	}
	if game.SymbolCount != 0 {
		t.Errorf("Expected SymbolCount to be 0, got %d", game.SymbolCount)
	}
}

func TestApplyMove(t *testing.T) {
	game := NewGameLogic()

	// Valid move
	err := game.ApplyMove(common.P1, 0, 0)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if game.Board[0][0] != common.P1 {
		t.Errorf("Expected board[0][0] to be P1, got %d", game.Board[0][0])
	}
	if game.Turn != common.P2 {
		t.Errorf("Expected Turn to be P2, got %d", game.Turn)
	}

	// Invalid move: wrong player
	err = game.ApplyMove(common.P1, 0, 1)
	if err != ErrNotYourTurn {
		t.Errorf("Expected ErrNotYourTurn, got %v", err)
	}

	// Invalid move: out of bounds
	err = game.ApplyMove(common.P2, -1, 0)
	if err != ErrOutOfBounds {
		t.Errorf("Expected ErrOutOfBounds, got %v", err)
	}
}

func mustMove(t *testing.T, game *GameLogic, p common.PlayerID, x, y int) {
	if err := game.ApplyMove(p, x, y); err != nil {
		t.Fatalf("ApplyMove(%v, %d, %d) failed: %v", p, x, y, err)
	}
}

func TestCheckWin(t *testing.T) {
	// Row win
	game := NewGameLogic()
	mustMove(t, game, common.P1, 0, 0) // X
	mustMove(t, game, common.P2, 1, 0) // O
	mustMove(t, game, common.P1, 0, 1) // X
	mustMove(t, game, common.P2, 1, 1) // O
	mustMove(t, game, common.P1, 0, 2) // X wins

	if !game.GameOver {
		t.Error("Expected game over")
	}
	if game.Winner != common.P1 {
		t.Errorf("Expected winner P1, got %d", game.Winner)
	}

	// Diagonal win
	game = NewGameLogic()
	mustMove(t, game, common.P1, 0, 0)
	mustMove(t, game, common.P2, 0, 1)
	mustMove(t, game, common.P1, 1, 1)
	mustMove(t, game, common.P2, 0, 2)
	mustMove(t, game, common.P1, 2, 2)

	if !game.GameOver {
		t.Error("Expected game over")
	}
	if game.Winner != common.P1 {
		t.Errorf("Expected winner P1, got %d", game.Winner)
	}
}

func TestDraw(t *testing.T) {
	// Fill board without win
	game := NewGameLogic()

	/*
	   X O X
	   X O O
	   O X X
	*/
	// P1: 0,0 (X)
	// P2: 0,1 (O)
	// P1: 0,2 (X)
	// P2: 1,1 (O)
	// P1: 1,0 (X)
	// P2: 1,2 (O)
	// P1: 2,1 (X)
	// P2: 2,0 (O)
	// P1: 2,2 (X) - Draw

	// Correction on moves to ensure no win and valid sequence
	// X O X
	// X O X
	// O X O

	// Real sequence:
	// P1(X): 0,0
	// P2(O): 1,1
	// P1(X): 0,2
	// P2(O): 0,1 .. wait logic...

	// Force the board state to test logic or simulate careful moves

	// X O X
	// X O O
	// O X X
	// Sequence:
	// 1. X -> 0,0
	// 2. O -> 0,1
	// 3. X -> 0,2
	// 4. O -> 1,1
	// 5. X -> 1,0
	// 6. O -> 1,2
	// 7. X -> 2,1
	// 8. O -> 2,0
	// 9. X -> 2,2

	mustMove(t, game, common.P1, 0, 0)
	mustMove(t, game, common.P2, 0, 1)
	mustMove(t, game, common.P1, 0, 2)
	mustMove(t, game, common.P2, 1, 1)
	mustMove(t, game, common.P1, 1, 0)
	mustMove(t, game, common.P2, 1, 2)
	mustMove(t, game, common.P1, 2, 1)
	mustMove(t, game, common.P2, 2, 0)
	mustMove(t, game, common.P1, 2, 2)

	if !game.GameOver {
		t.Error("Expected game over (draw)")
	}
	if game.Winner != common.Empty {
		t.Errorf("Expected no winner, got %d", game.Winner)
	}
}

func TestShouldTriggerChallenge(t *testing.T) {
	game := NewGameLogic()
	game.Board[0][0] = common.P2

	// P1 wants to take P2's cell
	if !game.ShouldTriggerChallenge(common.P1, 0, 0) {
		t.Error("Expected trigger challenge")
	}

	// P1 wants to take empty cell
	if game.ShouldTriggerChallenge(common.P1, 0, 1) {
		t.Error("Expected no trigger challenge for empty cell")
	}
}
