package logic

import (
	"Goonker/common"
	"errors"
)

// GameLogic gère l'état pur du jeu sans se soucier du réseau
type GameLogic struct {
	Board [3][3]common.PlayerID
	Turn  common.PlayerID
	Winner common.PlayerID
	GameOver bool
	Moves int
}

func NewGameLogic() *GameLogic {
	return &GameLogic{
		Turn: common.P1, // X commence toujours
	}
}

// ApplyMove tente de jouer un coup. Retourne une erreur si invalide.
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

	// Appliquer le coup
	g.Board[x][y] = player
	g.Moves++

	// Vérifier la victoire ou match nul
	if g.checkWin(player) {
		g.Winner = player
		g.GameOver = true
	} else if g.Moves >= 9 {
		g.GameOver = true // Match nul
	} else {
		// Changer le tour
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
	// Lignes et Colonnes
	for i := 0; i < 3; i++ {
		if b[i][0] == p && b[i][1] == p && b[i][2] == p { return true }
		if b[0][i] == p && b[1][i] == p && b[2][i] == p { return true }
	}
	// Diagonales
	if b[0][0] == p && b[1][1] == p && b[2][2] == p { return true }
	if b[0][2] == p && b[1][1] == p && b[2][0] == p { return true }
	
	return false
}