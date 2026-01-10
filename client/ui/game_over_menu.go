package ui

import "github.com/hajimehoshi/ebiten/v2"

// Button positions
const (
	GameOverMenuBackBtnY = 350.0
)

type GameOverMenu struct {
	BtnBack *Button
}

// Constructor for the game over menu.
func NewGameOverMenu() *GameOverMenu {
	menu := &GameOverMenu{}

	// Center buttons
	centerX := (float64(WindowWidth) - ButtonWidth) / 2

	// Create buttons
	menu.BtnBack = NewButton(centerX, GameOverMenuBackBtnY, ButtonWidth, ButtonHeight, "Back", BigFontFace)

	return menu
}

// Draw the game over menu to the screen.
func (m *GameOverMenu) Draw(screen *ebiten.Image) {
	m.BtnBack.Draw(screen)
}
