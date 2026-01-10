package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ButtonMiddleX = (float64(WindowWidth) - ButtonWidth) / 2
	ButtonSpacing = 100

	// Create room button position
	RoomsMenuCreateRoomBtnX = ButtonMiddleX - ButtonWidth - ButtonSpacing
	RoomsMenuCreateRoomBtnY = 50.0

	// Join game button position
	RoomsMenuJoinGameBtnX = ButtonMiddleX
	RoomsMenuJoinGameBtnY = 50.0

	// Against bot button position
	RoomsMenuPlayBotBtnX = ButtonMiddleX + ButtonWidth + ButtonSpacing
	RoomsMenuPlayBotBtnY = 50.0

	// Back button position
	RoomsMenuBackBtnX = ButtonMiddleX - ButtonWidth - ButtonSpacing
	RoomsMenuBackBtnY = 390.0

	// Text field
	RoomsMenuTextFieldX    = (float64(WindowWidth) - RoomsMenuTextFieldW) / 2
	RoomsMenuTextFieldY    = (float64(WindowHeight)-RoomsMenuTextFieldH)/2 - 100
	RoomsMenuTextFieldW    = 300
	RoomsMenuTextFieldH    = 50
	RoomsMenuTextFieldFont = 14
)

// RoomsMenu represents the rooms menu UI.
type RoomsMenu struct {
	Rooms         []*Room
	RoomIndex     int
	BtnPlayBot    *Button
	BtnCreateRoom *Button
	BtnJoinGame   *Button
	BtnBack       *Button
	RoomField     *TextField
}

// NewRoomsMenu creates a new RoomsMenu instance.
func NewRoomsMenu() *RoomsMenu {
	menu := &RoomsMenu{}

	// Create buttons
	menu.BtnCreateRoom = NewButton(RoomsMenuCreateRoomBtnX, RoomsMenuCreateRoomBtnY, ButtonWidth, ButtonHeight, "Create Room", BigFontFace)
	menu.BtnPlayBot = NewButton(RoomsMenuPlayBotBtnX, RoomsMenuPlayBotBtnY, ButtonWidth, ButtonHeight, "Against Bot", BigFontFace)
	menu.BtnJoinGame = NewButton(RoomsMenuJoinGameBtnX, RoomsMenuJoinGameBtnY, ButtonWidth, ButtonHeight, "Join Game", BigFontFace)
	menu.BtnBack = NewButton(RoomsMenuBackBtnX, RoomsMenuBackBtnY, ButtonWidth, ButtonHeight, "Back", BigFontFace)

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
