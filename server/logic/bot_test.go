package logic

import (
	"Goonker/common"
	"testing"
)

func TestGetBotMove(t *testing.T) {
	// Defend against immediate loss
	/*
	   X X .
	   . O .
	   . . .
	   Turn: O (Bot)
	   Expected: 0,2 (block X)
	*/
	logic := NewGameLogic()
	logic.Board[0][0] = common.P1
	logic.Board[0][1] = common.P1
	logic.Board[1][1] = common.P2
	logic.Turn = common.P2

	x, y := GetBotMove(logic)
	if x != 0 || y != 2 {
		t.Errorf("Expected defensive move at 0,2, got %d,%d", x, y)
	}

	// Take winning move
	/*
	   O O .
	   . X .
	   X . .
	   Turn: O (Bot)
	   Expected: 0,2 (win)
	*/
	logic = NewGameLogic()
	logic.Board[0][0] = common.P2
	logic.Board[0][1] = common.P2
	logic.Board[1][1] = common.P1
	logic.Board[2][0] = common.P1
	logic.Turn = common.P2

	x, y = GetBotMove(logic)
	if x != 0 || y != 2 {
		t.Errorf("Expected winning move at 0,2, got %d,%d", x, y)
	}
}

func TestIsBoardFull(t *testing.T) {
	board := [common.BoardSize][common.BoardSize]common.PlayerID{}
	if isBoardFull(board) {
		t.Error("Expected empty board to not be full")
	}

	for x := range common.BoardSize {
		for y := range common.BoardSize {
			board[x][y] = common.P1
		}
	}

	if !isBoardFull(board) {
		t.Error("Expected filled board to be full")
	}
}
