package ui

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TextField represents a text input field.
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

// NewTextField creates a new TextField instance.
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

// redraw redraws the text field image.
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
	dc.SetFontFace(BigFontFace)
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

// Update handles user input for the text field.
func (tf *TextField) Update() {
	// Handle click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		tf.HandleClick(float64(mx), float64(my))
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
			if tf.Insert(string(r)) {
				needsRedraw = true
			}
		}
	}

	// Backspace
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if tf.Backspace() {
			needsRedraw = true
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		tf.Focused = false
		needsRedraw = true
	}

	if needsRedraw {
		tf.redraw()
	}
}

// HandleClick updates focus state based on click position
func (tf *TextField) HandleClick(x, y float64) {
	wasFocused := tf.Focused
	tf.Focused = x >= tf.X && x <= tf.X+tf.Width &&
		y >= tf.Y && y <= tf.Y+tf.Height

	if wasFocused != tf.Focused {
		tf.redraw()
	}
}

// Insert adds a string to the text field if max length permits. Returns true if changed.
func (tf *TextField) Insert(s string) bool {
	if len(tf.Text)+len(s) <= tf.MaxLength {
		tf.Text += s
		tf.cursorVisible = true
		tf.cursorTimer = 0
		return true
	}
	return false
}

// Backspace removes the last character. Returns true if changed.
func (tf *TextField) Backspace() bool {
	if len(tf.Text) > 0 {
		tf.Text = tf.Text[:len(tf.Text)-1]
		tf.cursorVisible = true
		tf.cursorTimer = 0
		return true
	}
	return false
}

// Draw draws the text field to the screen.
func (tf *TextField) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(tf.X, tf.Y)
	screen.DrawImage(tf.Image, opts)
}
