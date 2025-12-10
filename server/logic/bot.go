package logic

import (
	"Goonker/common"
	"math"
	"time"
)

// Constants for bot behavior
const (
	BotThinkDelay = 500 * time.Millisecond
	InvalidCoord  = -1
	MaxDepth      = (common.BoardSize * common.BoardSize) + 1
)

// GetBotMove implements the minimax algorithm to find the best move for the bot.
func GetBotMove(logic *GameLogic) (int, int) {
	// Simulate "thinking" time for natural gameplay flow
	time.Sleep(BotThinkDelay)

	// Create a copy of the board to evaluate moves
	currentBoard := logic.Board

	// Initialize variables to track the best move
	bestScore := math.Inf(-1)
	moveX, moveY := InvalidCoord, InvalidCoord

	// Iterate through all possible moves
	for x := range common.BoardSize {
		for y := range common.BoardSize {
			if currentBoard[x][y] == common.Empty {
				// Simulate the move
				currentBoard[x][y] = common.P2

				// Evaluate the move
				score := minimax(currentBoard, 0, false)

				// Revert the move
				currentBoard[x][y] = common.Empty

				// Update the best move if this move is better
				if float64(score) > bestScore {
					bestScore = float64(score)
					moveX, moveY = x, y
				}
			}
		}
	}

	// Return the best move
	return moveX, moveY
}

// Minimax algorithm to evaluate the board
func minimax(board [common.BoardSize][common.BoardSize]common.PlayerID, depth int, isMaximizing bool) int {
	simulatedGame := &GameLogic{Board: board}

	if simulatedGame.checkWin(common.P2) {
		return MaxDepth - depth
	}
	if simulatedGame.checkWin(common.P1) {
		return depth - MaxDepth
	}
	if isBoardFull(board) {
		return 0
	}

	if isMaximizing {
		maxEval := math.Inf(-1)
		for x := range common.BoardSize {
			for y := range common.BoardSize {
				if board[x][y] == common.Empty {
					board[x][y] = common.P2
					eval := float64(minimax(board, depth+1, false))
					board[x][y] = common.Empty
					maxEval = math.Max(maxEval, eval)
				}
			}
		}
		return int(maxEval)
	} else {
		minEval := math.Inf(1)
		for x := range common.BoardSize {
			for y := range common.BoardSize {
				if board[x][y] == common.Empty {
					board[x][y] = common.P1
					eval := float64(minimax(board, depth+1, true))
					board[x][y] = common.Empty
					minEval = math.Min(minEval, eval)
				}
			}
		}
		return int(minEval)
	}
}

// Aide pour savoir si le plateau est rempli
func isBoardFull(board [common.BoardSize][common.BoardSize]common.PlayerID) bool {
	for x := range common.BoardSize {
		for y := range common.BoardSize {
			if board[x][y] == common.Empty {
				return false
			}
		}
	}
	return true
}
