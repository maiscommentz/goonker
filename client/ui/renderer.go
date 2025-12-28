package ui

import (
	"Goonker/client/assets"
	"Goonker/common"
	"bytes"
	"image/color"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	GameTitle    = "Goonker"
	WindowWidth  = 960
	WindowHeight = 540
	GridCol      = 3

	// Font sizes
	TitleFontSize    = 48
	SubtitleFontSize = 20
	TextFontSize     = 12

	// Positions
	PlayerTurnTextYPos = 150

	// Assets
	FontPath = "font.ttf"
)

var (
	gameFaceSource *text.GoTextFaceSource
	GameFont       *text.GoTextFace
)

// Init rendering components, like the images, the fonts...
func Init() {
	InitImages()

	fontData, err := fs.ReadFile(assets.AssetsFS, FontPath)
	if err != nil {
		log.Fatal(err)
	}

	src, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		log.Fatal(err)
	}
	gameFaceSource = src

	GameFont = &text.GoTextFace{
		Source: gameFaceSource,
		Size:   TextFontSize,
	}
}

// Render the main menu.
func RenderMenu(screen *ebiten.Image, menu *MainMenu) {
	screen.DrawImage(MainMenuImage, nil)
	menu.Draw(screen)
}

// Render the rooms menu.
func RenderRoomsMenu(screen *ebiten.Image, menu *RoomsMenu) {
	screen.DrawImage(RoomsMenuImage, nil)
	menu.Draw(screen)
}

// Render the waiting game menu.
func RenderWaitingGame(screen *ebiten.Image, waitingMenu *WaitingMenu) {
	waitingMenu.Draw(screen)
}

// Render the game.
func RenderGame(screen *ebiten.Image, grid *Grid, myTurn bool) {
	screen.DrawImage(GameMenuImage, nil)

	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()
	gridWidth, gridHeight := GridImage.Bounds().Dx(), GridImage.Bounds().Dy()

	offsetX := float64(screenWidth-gridWidth) / 2
	offsetY := float64(screenHeight-gridHeight) / 2

	cellSize := float64(gridWidth) / float64(grid.Col)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(offsetX, offsetY)
	screen.DrawImage(GridImage, op)

	for x := 0; x < grid.Col; x++ {
		for y := 0; y < grid.Col; y++ {
			var img *ebiten.Image

			switch grid.BoardData[x][y] {
			case common.P1:
				img = CrossImage
			case common.P2:
				img = CircleImage
			}

			if img != nil {
				opSym := &ebiten.DrawImageOptions{}
				// Subtract 1 from x and y to align symbols with the grid cells,
				// because the grid drawing starts at (cellSize, cellSize) rather than (0,0).
				cellX := (float64(x) - 1) * cellSize
				cellY := (float64(y) - 1) * cellSize

				opSym.GeoM.Translate(cellX, cellY)
				opSym.GeoM.Translate(offsetX, offsetY)

				screen.DrawImage(img, opSym)
			}
		}
	}

	if myTurn {
		msg := "It's goonkin' time"

		op := &text.DrawOptions{}

		w, _ := text.Measure(msg, GameFont, op.LineSpacing)

		x := (float64((WindowWidth)/2) - (gridSize / 2) - w) / 2

		op.GeoM.Translate(x, PlayerTurnTextYPos)

		op.ColorScale.ScaleWithColor(color.Black)

		text.Draw(screen, msg, GameFont, op)
	}
}

// Render win screen.
func RenderWin(screen *ebiten.Image) {
	screen.DrawImage(WinMenuImage, nil)
}

// Render lose screen.
func RenderLose(screen *ebiten.Image) {
	screen.DrawImage(LoseMenuImage, nil)
}

// Render draw screen.
func RenderDraw(screen *ebiten.Image) {
	screen.DrawImage(DrawMenuImage, nil)
}
