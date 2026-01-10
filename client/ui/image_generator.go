package ui

import (
	"Goonker/client/assets"
	"Goonker/common"
	"image/color"
	"log"
	"math"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

const (
	gridSize            = WindowHeight - 20
	lineWidth           = 12.0
	gridBackgroundColor = "#F4F6F7"
	gridBorderColor     = "#2C3E50"
	cellSize            = (gridSize / common.BoardSize)
	symbolLength        = cellSize/2 - 2*lineWidth

	// Wheel
	WheelSize         = 64
	WheelDots         = 12
	WheelRadius       = 22.0
	WheelDotMinRadius = 2.0
	WheelDotSizeRange = 3.0
	WheelMinAlpha     = 50
	WheelAlphaRange   = 205

	// Menu
	TitleYRatio      = 5
	TitleYRatioRooms = 8
)

var (
	BigFontFace      font.Face
	SmallFontFace    font.Face
	GridImage        *ebiten.Image
	CircleImage      *ebiten.Image
	CrossImage       *ebiten.Image
	WheelImage       *ebiten.Image
	MainMenuImage    *ebiten.Image
	WaitingMenuImage *ebiten.Image
	GameMenuImage    *ebiten.Image
	WinMenuImage     *ebiten.Image
	LoseMenuImage    *ebiten.Image
	DrawMenuImage    *ebiten.Image
	RoomsMenuImage   *ebiten.Image
	NoRoomsImage     *ebiten.Image
)

// InitImages initializes all the game images.
func InitImages() {
	initFont()

	DrawGrid(GridCol)
	DrawCircle()
	DrawCross()
	DrawWaitingWheel()
	DrawMainMenu(WindowWidth, WindowHeight, GameTitle)
	DrawWaitingMenu(WindowWidth, WindowHeight)
	DrawGameMenu(WindowWidth, WindowHeight)
	DrawWinMenu(WindowWidth, WindowHeight)
	DrawLoseMenu(WindowWidth, WindowHeight)
	DrawDrawMenu(WindowWidth, WindowHeight)
	DrawRoomsMenu(WindowWidth, WindowHeight)
}

func initFont() {
	// Read the font file from the embed filesystem
	fontBytes, err := assets.AssetsFS.ReadFile("font.ttf")
	if err != nil {
		log.Fatal(err)
	}

	// Parse the font
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	// Create the font face
	BigFontFace = truetype.NewFace(font, &truetype.Options{
		Size: 24,
	})
	SmallFontFace = truetype.NewFace(font, &truetype.Options{
		Size: 16,
	})
}

// DrawGrid draws the grid image.
func DrawGrid(col int) {
	dc := gg.NewContext(gridSize, gridSize)

	// Draw the background
	dc.SetHexColor(gridBackgroundColor)
	dc.Clear()

	// Draw the lines
	dc.SetHexColor(gridBorderColor)
	dc.SetLineWidth(lineWidth)
	dc.SetLineCap(gg.LineCapRound)

	cellSize := float64(gridSize) / float64(col)

	for i := range col {
		pos := float64(i) * cellSize

		// Vertical line
		dc.DrawLine(pos, 0, pos, float64(gridSize))

		// Horizontal line
		dc.DrawLine(0, pos, float64(gridSize), pos)
	}
	dc.Stroke()

	// Outer border
	offset := lineWidth / 2
	dc.DrawRectangle(offset, offset, float64(gridSize)-lineWidth, float64(gridSize)-lineWidth)
	dc.Stroke()

	GridImage = ebiten.NewImageFromImage(dc.Image())
}

// DrawCircle draws the circle symbol image.
func DrawCircle() {
	dc := gg.NewContext(gridSize, gridSize)

	centerX := float64(cellSize + cellSize/2)
	centerY := float64(cellSize + cellSize/2)

	dc.SetHexColor(gridBorderColor)
	dc.SetLineWidth(lineWidth)
	dc.DrawCircle(centerX, centerY, symbolLength)

	dc.Stroke()

	CircleImage = ebiten.NewImageFromImage(dc.Image())
}

// DrawCross draws the cross symbol image.
func DrawCross() {
	dc := gg.NewContext(gridSize, gridSize)

	dc.SetHexColor(gridBorderColor)
	dc.SetLineWidth(lineWidth)
	dc.SetLineCap(gg.LineCapRound)

	centerX := float64(cellSize + cellSize/2)
	centerY := float64(cellSize + cellSize/2)

	// Diagonal from top-left to bottom-right \
	dc.DrawLine(centerX-symbolLength, centerY-symbolLength, centerX+symbolLength, centerY+symbolLength)

	// Diagonal from bottom-left to top-right /
	dc.DrawLine(centerX-symbolLength, centerY+symbolLength, centerX+symbolLength, centerY-symbolLength)

	dc.Stroke()

	CrossImage = ebiten.NewImageFromImage(dc.Image())
}

// DrawWaitingWheel draws the waiting wheel animation frames.
func DrawWaitingWheel() {
	dc := gg.NewContext(WheelSize, WheelSize)

	cx, cy := float64(WheelSize)/2, float64(WheelSize)/2
	radius := WheelRadius

	for i := range WheelDots {
		angle := float64(i) * (2 * math.Pi) / float64(WheelDots)
		x := cx + math.Cos(angle)*radius
		y := cy + math.Sin(angle)*radius

		progress := float64(i) / float64(WheelDots)

		r := WheelDotMinRadius + (WheelDotSizeRange * progress)

		alpha := uint8(WheelMinAlpha + (WheelAlphaRange * progress))

		col := color.RGBA{R: 0, G: 0, B: 0, A: alpha}

		dc.SetColor(col)
		dc.DrawCircle(x, y, r)
		dc.Fill()
	}

	WheelImage = ebiten.NewImageFromImage(dc.Image())
}

// Draw the image for the main menu.
func DrawMainMenu(width, height int, title string) {
	dc := gg.NewContext(width, height)

	dc.SetHexColor(gridBackgroundColor)
	dc.Clear()

	dc.SetFontFace(BigFontFace)

	// Game title
	dc.SetHexColor(gridBorderColor)
	dc.DrawStringAnchored(title, float64(width/2), float64(height)/TitleYRatio, 0.5, 0.5)

	MainMenuImage = ebiten.NewImageFromImage(dc.Image())
}

// Draw the image for the waiting menu.
func DrawWaitingMenu(width, height int) {
	dc := gg.NewContext(width, height)

	dc.SetHexColor(gridBackgroundColor)
	dc.Clear()

	dc.SetFontFace(BigFontFace)

	dc.SetHexColor(gridBorderColor)
	dc.DrawStringAnchored("Waiting for another player...", float64(width/2), float64(height)/TitleYRatio, 0.5, 0.5)

	WaitingMenuImage = ebiten.NewImageFromImage(dc.Image())
}

// Draw the image for the game menu.
func DrawGameMenu(width, height int) {
	dc := gg.NewContext(width, height)

	dc.SetHexColor(gridBackgroundColor)
	dc.Clear()

	dc.SetFontFace(SmallFontFace)

	dc.SetHexColor(gridBorderColor)
	dc.DrawStringAnchored("Playing Goonker", (float64(width/2)-(gridSize/2))/2, float64(height)/TitleYRatio, 0.5, 0.5)

	GameMenuImage = ebiten.NewImageFromImage(dc.Image())
}

// Draw the image for the win menu.
func DrawWinMenu(width, height int) {
	dc := gg.NewContext(width, height)

	dc.SetHexColor(gridBackgroundColor)
	dc.Clear()

	dc.SetFontFace(BigFontFace)

	dc.SetHexColor(gridBorderColor)
	dc.DrawStringAnchored("You won !", float64(width/2), float64(height)/TitleYRatio, 0.5, 0.5)

	WinMenuImage = ebiten.NewImageFromImage(dc.Image())
}

// Draw the image for the lose menu.
func DrawLoseMenu(width, height int) {
	dc := gg.NewContext(width, height)

	dc.SetHexColor(gridBackgroundColor)
	dc.Clear()

	dc.SetFontFace(BigFontFace)

	dc.SetHexColor(gridBorderColor)
	dc.DrawStringAnchored("You lost :(", float64(width/2), float64(height)/TitleYRatio, 0.5, 0.5)

	LoseMenuImage = ebiten.NewImageFromImage(dc.Image())
}

// Draw the image for the draw menu.
func DrawDrawMenu(width, height int) {
	dc := gg.NewContext(width, height)

	dc.SetHexColor(gridBackgroundColor)
	dc.Clear()

	dc.SetFontFace(BigFontFace)

	dc.SetHexColor(gridBorderColor)
	dc.DrawStringAnchored("It's a draw...", float64(width/2), float64(height)/TitleYRatio, 0.5, 0.5)

	DrawMenuImage = ebiten.NewImageFromImage(dc.Image())
}

// Draw the image for the rooms menu.
func DrawRoomsMenu(width, height int) {
	dc := gg.NewContext(width, height)

	dc.SetHexColor(gridBackgroundColor)
	dc.Clear()

	dc.SetFontFace(BigFontFace)

	dc.SetHexColor(gridBorderColor)
	dc.DrawStringAnchored("Enter room ID", float64(width/2)-RoomsMenuTextFieldW, RoomsMenuTextFieldY, 0.5, 1.5)

	RoomsMenuImage = ebiten.NewImageFromImage(dc.Image())
}
