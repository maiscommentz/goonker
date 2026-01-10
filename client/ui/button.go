package ui

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
)

const (
	// Button configuration
	ButtonWidth        = 200.0
	ButtonHeight       = 60.0
	ButtonCornerRadius = 10.0
	ButtonTextYAnchor  = 0.35
)

// Button represents a clickable UI button.
type Button struct {
	X, Y, Width, Height float64
	Image               *ebiten.Image
	Text                string
}

// NewButton creates a new Button instance.
func NewButton(x, y, w, h float64, text string, fontFace font.Face) *Button {
	b := &Button{
		X: x, Y: y, Width: w, Height: h,
		Text: text,
	}

	dc := gg.NewContext(int(w), int(h))

	dc.DrawRoundedRectangle(0, 0, w, h, ButtonCornerRadius)
	dc.SetHexColor(gridBorderColor)
	dc.Fill()

	dc.SetFontFace(fontFace)
	dc.SetHexColor(gridBackgroundColor)
	dc.DrawStringAnchored(text, w/2, h/2, 0.5, ButtonTextYAnchor)

	b.Image = ebiten.NewImageFromImage(dc.Image())

	return b
}

// Draw the button to the screen.
func (b *Button) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(b.X, b.Y)
	screen.DrawImage(b.Image, opts)
}

// Check if a button is clicked.
func (b *Button) IsClicked() bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		return b.Contains(float64(mx), float64(my))
	}
	return false
}

// Contains checks if the given coordinates are within the button's bounds.
func (b *Button) Contains(x, y float64) bool {
	return x >= b.X && x <= b.X+b.Width &&
		y >= b.Y && y <= b.Y+b.Height
}
