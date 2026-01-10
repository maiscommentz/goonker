package ui

import "github.com/hajimehoshi/ebiten/v2"

// Button positions
const (
	MainMenuPlayBtnY = 200.0
	MainMenuQuitBtnY = 280.0
)

// MainMenu represents the main menu UI.
type MainMenu struct {
	BtnPlay *Button
	BtnQuit *Button
}

// NewMainMenu creates a new MainMenu instance.
func NewMainMenu() *MainMenu {
	menu := &MainMenu{}

	// Center buttons
	centerX := (float64(WindowWidth) - ButtonWidth) / 2

	// Create buttons
	menu.BtnPlay = NewButton(centerX, MainMenuPlayBtnY, ButtonWidth, ButtonHeight, "Play", BigFontFace)
	menu.BtnQuit = NewButton(centerX, MainMenuQuitBtnY, ButtonWidth, ButtonHeight, "Quit", BigFontFace)

	return menu
}

// Draw the main menu to the screen.
func (m *MainMenu) Draw(screen *ebiten.Image) {
	screen.DrawImage(MainMenuImage, nil)
	m.BtnPlay.Draw(screen)
	m.BtnQuit.Draw(screen)
}
