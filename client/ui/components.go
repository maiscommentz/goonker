package ui

import (
	"Goonker/common"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Button struct {
	X, Y, Width, Height float64
	Image               *ebiten.Image
	Text                string
}

type Menu struct {
	MenuImage *ebiten.Image
	BtnPlay   *Button
	BtnQuit   *Button
}

type Grid struct {
	Col        int
	BoardImage *ebiten.Image
	BoardData  [3][3]common.PlayerID
}

type Cell struct {
	Btn    Button
	Symbol common.PlayerID
}

type Drawable interface {
	Draw(screen *ebiten.Image)
}

func NewButton(x, y, w, h float64, text string) *Button {
	b := &Button{
		X: x, Y: y, Width: w, Height: h,
		Text: text,
	}

	dc := gg.NewContext(int(w), int(h))

	dc.DrawRoundedRectangle(0, 0, w, h, 10)
	dc.SetHexColor("#2C3E50")
	dc.Fill()

	// dc.LoadFontFace("arial.ttf", 24)
	dc.SetHexColor("#FFFFFF")
	dc.DrawStringAnchored(text, w/2, h/2, 0.5, 0.35)

	b.Image = ebiten.NewImageFromImage(dc.Image())

	return b
}

func NewGrid(col, row int) *Grid {
	g := &Grid{
		Col: col,
	}

	gridImage := DrawGrid(col)
	g.BoardImage = ebiten.NewImageFromImage(gridImage)

	return g
}

func (b *Button) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(b.X, b.Y)
	screen.DrawImage(b.Image, opts)
}

func (b *Button) IsClicked() bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		fmx, fmy := float64(mx), float64(my)

		return fmx >= b.X && fmx <= b.X+b.Width &&
			fmy >= b.Y && fmy <= b.Y+b.Height
	}
	return false
}
