package ui

import (
	"Goonker/common"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Grid struct {
	Col       int
	BoardData [GridCol][GridCol]common.PlayerID
}

// Get the clicked cell.
func (g *Grid) OnClick() (int, int, bool) {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return -1, -1, false
	}

	mx, my := ebiten.CursorPosition()

	gridW, gridH := GridImage.Bounds().Dx(), GridImage.Bounds().Dy()

	offsetX := (WindowWidth - gridW) / 2
	offsetY := (WindowHeight - gridH) / 2

	localX := mx - offsetX
	localY := my - offsetY

	if localX < 0 || localY < 0 || localX >= gridW || localY >= gridH {
		return -1, -1, false
	}

	cellSize := gridW / GridCol

	cellX := localX / cellSize
	cellY := localY / cellSize

	return cellX, cellY, true
}
