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

// Room represents a single room item in the rooms list.
type Room struct {
	JoinBtn *Button
	Id      string
	Image   *ebiten.Image
}

// NewRoom creates a new Room UI component.
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
	room.JoinBtn = NewButton(0, 0, ButtonWidth/3, ButtonHeight/2, "Join", BigFontFace)

	return room
}
