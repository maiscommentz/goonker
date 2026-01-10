package ui

import "github.com/hajimehoshi/ebiten/v2"

// Button positions
const (
	GameOverMenuBackBtnY = 350.0
)

// GameOverMenu represents the game over screen UI.
type GameOverMenu struct {
	BtnBack *Button
}

// NewGameOverMenu creates a new GameOverMenu instance.
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
