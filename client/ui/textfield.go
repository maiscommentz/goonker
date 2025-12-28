package ui

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type TextField struct {
	X, Y          float64
	Width, Height float64
	Text          string
	Focused       bool
	MaxLength     int
	Image         *ebiten.Image

	cursorVisible bool
	cursorTimer   int
	fontSize      float64
}

func NewTextField(x, y, w, h float64, fontSize float64) *TextField {
	tf := &TextField{
		X:         x,
		Y:         y,
		Width:     w,
		Height:    h,
		Text:      "",
		MaxLength: 50,
		fontSize:  fontSize,
	}
	tf.redraw()
	return tf
}

func (tf *TextField) redraw() {
	dc := gg.NewContext(int(tf.Width), int(tf.Height))

	// Background
	dc.DrawRoundedRectangle(0, 0, tf.Width, tf.Height, 5)
	if tf.Focused {
		dc.SetColor(color.RGBA{240, 240, 255, 255})
	} else {
		dc.SetColor(color.RGBA{255, 255, 255, 255})
	}
	dc.Fill()

	// Border
	dc.DrawRoundedRectangle(0, 0, tf.Width, tf.Height, 5)
	dc.SetLineWidth(2)
	if tf.Focused {
		dc.SetColor(color.RGBA{0, 100, 255, 255})
	} else {
		dc.SetColor(color.RGBA{150, 150, 150, 255})
	}
	dc.Stroke()

	// Text
	dc.SetFontFace(FontFace)
	dc.SetColor(color.Black)
	dc.DrawString(tf.Text, 10, tf.Height/2+tf.fontSize/3)

	// Blinking cursor
	if tf.Focused && tf.cursorVisible {
		textWidth, _ := dc.MeasureString(tf.Text)
		cursorX := 10 + textWidth
		cursorY1 := tf.Height/2 - tf.fontSize/2
		cursorY2 := tf.Height/2 + tf.fontSize/2

		dc.SetLineWidth(2)
		dc.SetColor(color.Black)
		dc.DrawLine(cursorX, cursorY1, cursorX, cursorY2)
		dc.Stroke()
	}

	tf.Image = ebiten.NewImageFromImage(dc.Image())
}

func (tf *TextField) Update() {
	// Handle click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		wasFocused := tf.Focused
		tf.Focused = float64(mx) >= tf.X && float64(mx) <= tf.X+tf.Width &&
			float64(my) >= tf.Y && float64(my) <= tf.Y+tf.Height

		if wasFocused != tf.Focused {
			tf.redraw()
		}
	}

	if !tf.Focused {
		return
	}

	needsRedraw := false

	// Blinking cursor
	tf.cursorTimer++
	if tf.cursorTimer >= 30 {
		tf.cursorVisible = !tf.cursorVisible
		tf.cursorTimer = 0
		needsRedraw = true
	}

	// Handle input
	runes := ebiten.AppendInputChars(nil)
	if len(runes) > 0 {
		for _, r := range runes {
			if len(tf.Text) < tf.MaxLength {
				tf.Text += string(r)
			}
		}
		needsRedraw = true
		tf.cursorVisible = true
		tf.cursorTimer = 0
	}

	// Backspace
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(tf.Text) > 0 {
		tf.Text = tf.Text[:len(tf.Text)-1]
		needsRedraw = true
		tf.cursorVisible = true
		tf.cursorTimer = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		tf.Focused = false
		needsRedraw = true
	}

	if needsRedraw {
		tf.redraw()
	}
}

func (tf *TextField) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(tf.X, tf.Y)
	screen.DrawImage(tf.Image, opts)
}
