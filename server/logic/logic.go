package logic

import (
	"Goonker/common"
	"errors"
)

// GameLogic manages the state of a Tic-Tac-Toe game
type GameLogic struct {
	Board [3][3]common.PlayerID
	Turn  common.PlayerID
	Winner common.PlayerID
	GameOver bool
	Moves int
}

func NewGameLogic() *GameLogic {
	return &GameLogic{
		Turn: common.P1, // X always starts
	}
}

// ApplyMove attempts to play a move. Returns an error if invalid.
func (g *GameLogic) ApplyMove(player common.PlayerID, x, y int) error {
	if g.GameOver {
		return errors.New("game is over")
	}
	if player != g.Turn {
		return errors.New("not your turn")
	}
	if x < 0 || x > 2 || y < 0 || y > 2 {
		return errors.New("out of bounds")
	}
	if g.Board[x][y] != common.Empty {
		return errors.New("cell already occupied")
	}

	// Apply the move
	g.Board[x][y] = player
	g.Moves++

	// Check for win or draw
	if g.checkWin(player) {
		g.Winner = player
		g.GameOver = true
	} else if g.Moves >= 9 {
		g.GameOver = true // Draw
	} else {
		// Change turn
		if g.Turn == common.P1 {
			g.Turn = common.P2
		} else {
			g.Turn = common.P1
		}
	}

	return nil
}

func (g *GameLogic) checkWin(p common.PlayerID) bool {
	b := g.Board
	// Rows and Columns
	for i := 0; i < 3; i++ {
		if b[i][0] == p && b[i][1] == p && b[i][2] == p { return true }
		if b[0][i] == p && b[1][i] == p && b[2][i] == p { return true }
	}
	// Diagonals
	if b[0][0] == p && b[1][1] == p && b[2][2] == p { return true }
	if b[0][2] == p && b[1][1] == p && b[2][0] == p { return true }
	
	return false
}

func (g *GameLogic) PrintConsoleBoard() {
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			cell := g.Board[x][y]
			var symbol string
			if cell == common.P1 {
				symbol = "X"
			} else if cell == common.P2 {
				symbol = "O"
			} else {
				symbol = "."
			}
			if x < 2 {
				symbol += " | "
			}
			print(symbol)
		}
		println()
		if y < 2 {
			println("---------")
		}
	}
}
