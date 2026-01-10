package ui

import (
	"Goonker/common"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Grid represents the game board grid.
type Grid struct {
	Col       int
	BoardData [GridCol][GridCol]common.PlayerID
}

// OnClick checks for mouse input and returns the clicked cell coordinates.
func (g *Grid) OnClick() (int, int, bool) {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return -1, -1, false
	}

	mx, my := ebiten.CursorPosition()
	return g.GetCellAt(mx, my)
}

// GetCellAt calculates which cell contains the given screen coordinates.
func (g *Grid) GetCellAt(mx, my int) (int, int, bool) {
	// Calculate grid width and height
	gridW, gridH := GridImage.Bounds().Dx(), GridImage.Bounds().Dy()

	// Calculate offsets to center the grid
	offsetX := (WindowWidth - gridW) / 2
	offsetY := (WindowHeight - gridH) / 2

	// Calculate local coordinates inside the grid
	localX := mx - offsetX
	localY := my - offsetY

	// Check if the click is inside the grid
	if localX < 0 || localY < 0 || localX >= gridW || localY >= gridH {
		return -1, -1, false
	}

	cellSize := gridW / GridCol

	// Determine cell coordinates
	cellX := localX / cellSize
	cellY := localY / cellSize

	// Clamp values just to be safe (though equality check above handles most)
	if cellX >= g.Col || cellY >= g.Col {
		return -1, -1, false
	}

	return cellX, cellY, true
}
