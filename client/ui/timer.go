package ui

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ClockWidth  = 300
	ClockHeight = 30
	ClockPosX   = (float64(WindowWidth) - ClockWidth) / 2
	ClockPosY   = 90

	TicksPerSeconds = 60

	RGBDefaultVal = 100
	MaxRGBAVal    = 255
)

var whiteImg *ebiten.Image

// TimerInit initializes the timer by creating a 1x1 white image for drawing rectangles.
func TimerInit() {
	whiteImg = ebiten.NewImage(1, 1)
	whiteImg.Fill(color.White)
}

// Timer handles the logic for a countdown
type Timer struct {
	TotalDuration   time.Duration
	CurrentDuration time.Duration
	IsRunning       bool
	OnEnd           func()
}

// NewTimer creates a new timer with a set duration.
func NewTimer(d time.Duration) *Timer {
	return &Timer{
		TotalDuration:   d,
		CurrentDuration: d,
		IsRunning:       true,
	}
}

// Update calculates the remaining time
func (t *Timer) Update() {
	if !t.IsRunning {
		return
	}

	// time.Second / 60 is approximately 16.66ms
	t.CurrentDuration -= time.Second / TicksPerSeconds

	if t.CurrentDuration <= 0 {
		t.CurrentDuration = 0
		t.IsRunning = false
		if t.OnEnd != nil {
			t.OnEnd()
		}
	}
}

// Ratio returns a value between 0.0 and 1.0 representing progress
func (t *Timer) Ratio() float32 {
	if t.TotalDuration == 0 {
		return 0
	}
	return float32(t.CurrentDuration) / float32(t.TotalDuration)
}

// drawRect draws a filled rectangle with the specified color at the given position and size.
func drawRect(screen *ebiten.Image, x, y, width, height float64, clr color.Color) {
	op := &ebiten.DrawImageOptions{}

	// Scale the 1x1 pixel to the target size
	op.GeoM.Scale(width, height)

	// Move it to the target position
	op.GeoM.Translate(x, y)

	// Color it
	op.ColorScale.ScaleWithColor(clr)

	screen.DrawImage(whiteImg, op)
}

// Draw renders the timer bar to the screen.
func (t *Timer) Draw(screen *ebiten.Image) {
	// Draw Background (Gray)
	bgColor := color.RGBA{RGBDefaultVal, RGBDefaultVal, RGBDefaultVal, MaxRGBAVal}
	drawRect(screen, ClockPosX, ClockPosY, ClockWidth, ClockHeight, bgColor)

	// Draw Foreground (Green -> Red)
	ratio := t.Ratio()
	currentWidth := float64(ClockWidth * ratio)

	// Dynamic color calculation
	c := color.RGBA{
		R: uint8(MaxRGBAVal * (1 - ratio)), // Red increases as time runs out
		G: uint8(MaxRGBAVal * ratio),       // Green decreases as time runs out
		B: 0,
		A: MaxRGBAVal,
	}

	drawRect(screen, ClockPosX, ClockPosY, currentWidth, ClockHeight, c)
}
