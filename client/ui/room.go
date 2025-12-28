package ui

import (
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

// Rooms list
const (
	RoomsLineWidth  = 2
	RoomsRowPadding = 10
	RoomsRowGap     = 70.0
)

type Room struct {
	JoinBtn *Button
	Id      string
	Image   *ebiten.Image
}

// Constructor for the room.
func NewRoom(id string) *Room {
	// Room row dimensions
	width := WindowWidth / 2
	height := TitleFontSize

	// Initialize room
	room := &Room{
		Id: id,
	}

	dc := gg.NewContext(width, height)

	// Draw bottom separator line
	dc.SetHexColor(gridBorderColor)
	dc.SetLineWidth(RoomsLineWidth)
	dc.DrawLine(0, float64(height), float64(width), float64(height))
	dc.Stroke()

	// Load font for room name
	if err := dc.LoadFontFace(FontPath, TextFontSize); err != nil {
		log.Printf("Error loading font: %v", err)
	}

	// Draw Room Name (Left aligned)
	dc.SetHexColor(gridBorderColor)
	dc.DrawStringAnchored(id, RoomsRowPadding, float64(height)/3, 0.0, 0.5)

	room.Image = ebiten.NewImageFromImage(dc.Image())

	// Initialize Join Button
	room.JoinBtn = NewButton(0, 0, ButtonWidth/3, ButtonHeight/2, "Join", TextFontSize)

	return room
}

// Draw the room at the specified index
func (r *Room) Draw(screen *ebiten.Image, index int) {
	// Constants for layout
	listX := float64(RoomsMenuBtnX + ButtonWidth + RoomsMenuBtnX)
	listY := float64(RoomsMenuBackBtnY)
	yPos := listY + float64(index)*RoomsRowGap

	// Draw the Row Image (Name + Line)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(listX, yPos)
	screen.DrawImage(r.Image, op)

	// Update Button Position
	btnX := listX + listX + ButtonWidth/2
	r.JoinBtn.X = btnX
	r.JoinBtn.Y = yPos

	// Draw Button
	r.JoinBtn.Draw(screen)
}
