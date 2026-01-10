package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	// Spinning wheel parameters
	WheelTintRed   = 0.8
	WheelTintGreen = 0.8
	WheelTintBlue  = 1.0

	// Room ID text height
	WaitingMenuRoomTextY = (float64(WindowHeight) / 2.0) - 100
)

// WaitingMenu represents the waiting screen UI.
type WaitingMenu struct {
	RotationAngle float64
	RoomId        string
}

// Draw draws the waiting menu to the screen, including the spinning wheel and room ID.
func (waitingMenu *WaitingMenu) Draw(screen *ebiten.Image) {
	screen.DrawImage(WaitingMenuImage, nil)
	screenCenterX := float64(WindowWidth) / 2.0
	screenCenterY := float64(WindowHeight) / 2.0

	// Draw the spinning wheel
	w := WheelImage.Bounds().Dx()
	h := WheelImage.Bounds().Dy()
	halfW := float64(w) / 2.0
	halfH := float64(h) / 2.0

	wheelOpt := &ebiten.DrawImageOptions{}

	wheelOpt.GeoM.Translate(-halfW, -halfH)
	wheelOpt.GeoM.Rotate(waitingMenu.RotationAngle)
	wheelOpt.GeoM.Translate(screenCenterX, screenCenterY)
	wheelOpt.ColorScale.Scale(WheelTintRed, WheelTintGreen, WheelTintBlue, 1)

	screen.DrawImage(WheelImage, wheelOpt)

	// Draw the text
	waitingRoomText := fmt.Sprintf("Room ID : %s", waitingMenu.RoomId)
	textOpt := &text.DrawOptions{}

	textWidth, _ := text.Measure(waitingRoomText, SmallGameFont, textOpt.LineSpacing)
	x := (screenCenterX - (textWidth / 2))

	textOpt.GeoM.Translate(x, WaitingMenuRoomTextY)

	textOpt.ColorScale.ScaleWithColor(color.Black)

	text.Draw(screen, waitingRoomText, SmallGameFont, textOpt)
}
