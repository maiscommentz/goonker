package logic

import (
	"fmt"
	"strings"

	"Goonker/common"
	"errors"
)

// Game constants
const (
	// Maximum number of moves in a Tic-Tac-Toe game
	MaxMoves = common.BoardSize * common.BoardSize

	// Display Symbols (for console debug)
	SymbolP1    = "X"
	SymbolP2    = "O"
	SymbolEmpty = "."
	SeparatorV  = " | "
	SeparatorH  = "---------"
)

// Error messages
var (
	ErrGameOver     = errors.New("game is over")
	ErrNotYourTurn  = errors.New("not your turn")
	ErrOutOfBounds  = errors.New("out of bounds")
	ErrCellOccupied = errors.New("cell already occupied")
)

// GameLogic manages the state of a Tic-Tac-Toe game
type GameLogic struct {
	Board       [common.BoardSize][common.BoardSize]common.PlayerID
	Turn        common.PlayerID
	Winner      common.PlayerID
	GameOver    bool
	SymbolCount int
}

// NewGameLogic initializes a new game state.
func NewGameLogic() *GameLogic {
	return &GameLogic{
		Turn: common.P1, // X always starts
	}
}

func (g *GameLogic) ShouldTriggerChallenge(player common.PlayerID, x, y int) bool {
	return g.Board[x][y] != player && g.Board[x][y] != common.Empty
}

// ApplyMove attempts to play a move. Returns an error if invalid. Or true if a minigame must start
func (g *GameLogic) ApplyMove(player common.PlayerID, x, y int) error {
	// Validate move
	if g.GameOver {
		return ErrGameOver
	}
	if player != g.Turn {
		return ErrNotYourTurn
	}
	if x < 0 || x > common.BoardSize-1 || y < 0 || y > common.BoardSize-1 {
		return ErrOutOfBounds
	}
	if g.Board[x][y] == player {
		return ErrCellOccupied
	}

	if g.Board[x][y] == common.Empty {
		// Place the player symbol
		g.Board[x][y] = player
		g.SymbolCount++
	}

	// Check for win or draw
	if g.checkWin(player) {
		g.Winner = player
		g.GameOver = true
	} else if g.SymbolCount >= MaxMoves {
		g.GameOver = true // Draw
	} else {
		// Toggle turn
		if g.Turn == common.P1 {
			g.Turn = common.P2
		} else {
			g.Turn = common.P1
		}
	}

	return nil
}

// DeleteMove empties the given board cell
func (g *GameLogic) DeleteMove(x, y int) {
	g.Board[x][y] = common.Empty
	g.SymbolCount--
}

// checkWin scans rows, columns, and diagonals for a complete line.
func (g *GameLogic) checkWin(p common.PlayerID) bool {
	board := g.Board
	boardSize := common.BoardSize

	// Check Rows and Columns
	for i := range boardSize {
		rowWin, colWin := true, true
		for j := range boardSize {
			// Check Column i (varying rows j)
			if board[i][j] != p {
				colWin = false
			}
			// Check Row i (varying cols j)
			if board[j][i] != p {
				rowWin = false
			}
		}
		if rowWin || colWin {
			return true
		}
	}

	// Check Diagonals
	diag1, diag2 := true, true
	for i := range boardSize {
		// Top-left to bottom-right (0,0 -> 1,1 -> 2,2)
		if board[i][i] != p {
			diag1 = false
		}
		// Bottom-left to top-right (0,2 -> 1,1 -> 2,0)
		if board[i][boardSize-1-i] != p {
			diag2 = false
		}
	}

	return diag1 || diag2
}

// PrintConsoleBoard renders the board state to the console for debugging.
func (g *GameLogic) PrintConsoleBoard() {
	for y := range common.BoardSize {
		var line []string
		for x := range common.BoardSize {
			cell := g.Board[x][y]
			switch cell {
			case common.P1:
				line = append(line, SymbolP1)
			case common.P2:
				line = append(line, SymbolP2)
			default:
				line = append(line, SymbolEmpty)
			}
		}
		fmt.Println(strings.Join(line, SeparatorV))

		// Print horizontal separator only between rows
		if y < common.BoardSize-1 {
			fmt.Println(SeparatorH)
		}
	}
}
