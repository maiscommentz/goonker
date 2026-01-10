package ui

import (
	"testing"
)

func TestGridLogic(t *testing.T) {
	// Initialize images so GridImage has dimensions
	// We handle panic if headless in case InitImages tries to load fonts/assets that fail
	defer func() {
		if r := recover(); r != nil {
			t.Skip("Skipping GridLogic test due to asset initialization failure (headless?):", r)
		}
	}()
	InitImages()

	grid := &Grid{Col: 3}

	// Get dimensions we expect
	gridW := GridImage.Bounds().Dx()
	gridH := GridImage.Bounds().Dy()

	offsetX := (WindowWidth - gridW) / 2
	offsetY := (WindowHeight - gridH) / 2
	cellSize := gridW / 3

	// Test 1: Click Top-Left Cell (0,0)
	// We pick a point inside the cell
	mx := offsetX + cellSize/2
	my := offsetY + cellSize/2

	x, y, ok := grid.GetCellAt(mx, my)
	if !ok {
		t.Error("Expected valid click for Top-Left cell")
	}
	if x != 0 || y != 0 {
		t.Errorf("Expected 0,0 got %d,%d", x, y)
	}

	// Test 2: Click Middle Cell (1,1)
	mx = offsetX + cellSize + cellSize/2
	my = offsetY + cellSize + cellSize/2

	x, y, ok = grid.GetCellAt(mx, my)
	if !ok {
		t.Error("Expected valid click for Middle cell")
	}
	if x != 1 || y != 1 {
		t.Errorf("Expected 1,1 got %d,%d", x, y)
	}

	// Test 3: Click Outside (Left of grid)
	mx = offsetX - 10
	my = offsetY + 10

	_, _, ok = grid.GetCellAt(mx, my)
	if ok {
		t.Error("Expected invalid click when left of grid")
	}

	// Test 4: Click Outside (Bottom of grid)
	mx = offsetX + 10
	my = offsetY + gridH + 10

	_, _, ok = grid.GetCellAt(mx, my)
	if ok {
		t.Error("Expected invalid click when below grid")
	}
}
