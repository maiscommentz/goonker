package ui

import (
	"Goonker/common"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	GameTitle    = "Goonker"
	WindowWidth  = 960
	WindowHeight = 540
	GridCol      = 3
)

func RenderMenu(screen *ebiten.Image, menu *Menu) {
	screen.DrawImage(menu.MenuImage, nil)
	menu.Draw(screen)
	menu.BtnPlay.Draw(screen)
	menu.BtnQuit.Draw(screen)
}

func RenderGame(screen *ebiten.Image, grid *Grid, myTurn bool) {
	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()
	gridWidth, gridHeight := grid.BoardImage.Bounds().Dx(), grid.BoardImage.Bounds().Dy()

	offsetX := float64(screenWidth-gridWidth) / 2
	offsetY := float64(screenHeight-gridHeight) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(offsetX, offsetY)
	screen.DrawImage(grid.BoardImage, op)

	for x := 0; x < grid.Col; x++ {
		for y := 0; y < grid.Col; y++ {
			var img *ebiten.Image

			switch grid.BoardData[x][y] {
			case common.P1:
				img = ebiten.NewImageFromImage(DrawCross(x, y))
			case common.P2:
				img = ebiten.NewImageFromImage(DrawCircle(x, y))
			}

			if img != nil {
				opSym := &ebiten.DrawImageOptions{}
				opSym.GeoM.Translate(offsetX, offsetY)
				screen.DrawImage(img, opSym)
			}
		}
	}
}
