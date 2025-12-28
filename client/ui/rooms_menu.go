package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	RoomsMenuBackBtnY       = 390.0
	RoomsMenuPlayBotBtnX    = 600.0
	RoomsMenuPlayBotBtnY    = 50.0
	RoomsMenuJoinGameBtnX   = 350.0
	RoomsMenuJoinGameBtnY   = 50.0
	RoomsMenuCreateRoomBtnY = 50.0
	RoomsMenuBtnX           = 50.0

	// Text field
	RoomsMenuTextFieldX    = (float64(WindowWidth)-RoomsMenuTextFieldW)/2 + 300
	RoomsMenuTextFieldY    = (float64(WindowHeight)-RoomsMenuTextFieldH)/2 - 100
	RoomsMenuTextFieldW    = 300
	RoomsMenuTextFieldH    = 50
	RoomsMenuTextFieldFont = 16
)

type RoomsMenu struct {
	Rooms         []*Room
	RoomIndex     int
	BtnPlayBot    *Button
	BtnCreateRoom *Button
	BtnJoinGame   *Button
	BtnBack       *Button
	RoomField     *TextField
}

// Constructor for the play menu.
func NewRoomsMenu() *RoomsMenu {
	menu := &RoomsMenu{}

	// Create buttons
	menu.BtnBack = NewButton(RoomsMenuBtnX, RoomsMenuBackBtnY, ButtonWidth, ButtonHeight, "Back", SubtitleFontSize)
	menu.BtnPlayBot = NewButton(RoomsMenuPlayBotBtnX, RoomsMenuPlayBotBtnY, ButtonWidth, ButtonHeight, "Against Bot", SubtitleFontSize)
	menu.BtnCreateRoom = NewButton(RoomsMenuBtnX, RoomsMenuCreateRoomBtnY, ButtonWidth, ButtonHeight, "Create Room", SubtitleFontSize)
	menu.BtnJoinGame = NewButton(RoomsMenuJoinGameBtnX, RoomsMenuJoinGameBtnY, ButtonWidth, ButtonHeight, "Join Game", SubtitleFontSize)

	// Create textfield
	menu.RoomField = NewTextField(RoomsMenuTextFieldX, RoomsMenuTextFieldY, RoomsMenuTextFieldW, RoomsMenuTextFieldH, RoomsMenuTextFieldFont)

	return menu
}

// Draw the rooms menu to the screen.
func (m *RoomsMenu) Draw(screen *ebiten.Image) {
	screen.DrawImage(RoomsMenuImage, nil)
	m.BtnPlayBot.Draw(screen)
	m.BtnCreateRoom.Draw(screen)
	m.BtnJoinGame.Draw(screen)
	m.BtnBack.Draw(screen)
	m.RoomField.Draw(screen)
}
