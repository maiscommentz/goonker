package ui

import (
	"Goonker/common"

	"github.com/hajimehoshi/ebiten/v2"
)

func RenderMenu(screen *ebiten.Image, menu *Menu) {
	screen.DrawImage(menu.MenuImage, nil)
	menu.BtnPlay.Draw(screen)
	menu.BtnQuit.Draw(screen)
}

func RenderGame(screen *ebiten.Image, grid *Grid) {
	screen.DrawImage(grid.BoardImage, nil)

	for x := 0; x < grid.Col; x++ {
		for y := 0; y < grid.Col; y++ {
			switch grid.BoardData[x][y] {
			case common.P1:
				crossImage := ebiten.NewImageFromImage(DrawCross(x, y))
				screen.DrawImage(crossImage, nil)
			case common.P2:
				circleImage := ebiten.NewImageFromImage(DrawCircle(x, y))
				screen.DrawImage(circleImage, nil)
			}
		}
	}
}
