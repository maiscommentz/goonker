package ui

import (
	"image"
	"log"

	"github.com/fogleman/gg"
)

const (
	gridSize            = 600
	lineWidth           = 12.0
	gridBackgroundColor = "#F4F6F7"
	gridBorderColor     = "#2C3E50"

	cellSize = (gridSize / 3)

	symbolLength = cellSize/2 - 2*lineWidth
)

func DrawGrid(col int) image.Image {
	dc := gg.NewContext(gridSize, gridSize)

	// Draw the background
	dc.SetHexColor(gridBackgroundColor)
	dc.Clear()

	// Draw the lines
	dc.SetHexColor(gridBorderColor)
	dc.SetLineWidth(lineWidth)
	dc.SetLineCap(gg.LineCapRound)

	cellSize := float64(gridSize) / float64(col)

	for i := 1; i < col; i++ {
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

	return dc.Image()
}

func DrawCircle(col, row int) image.Image {
	dc := gg.NewContext(gridSize, gridSize)

	centerX := float64(col*cellSize + cellSize/2)
	centerY := float64(row*cellSize + cellSize/2)

	dc.SetHexColor(gridBorderColor)
	dc.SetLineWidth(lineWidth)
	dc.DrawCircle(centerX, centerY, symbolLength)

	dc.Stroke()

	return dc.Image()
}

func DrawCross(col, row int) image.Image {
	dc := gg.NewContext(gridSize, gridSize)

	dc.SetHexColor(gridBorderColor)
	dc.SetLineWidth(lineWidth)
	dc.SetLineCap(gg.LineCapRound)

	centerX := float64(col*cellSize + cellSize/2)
	centerY := float64(row*cellSize + cellSize/2)

	// Diagonal from top-left to bottom-right \
	dc.DrawLine(centerX-symbolLength, centerY-symbolLength, centerX+symbolLength, centerY+symbolLength)

	// Diagonal from bottom-left to top-right /
	dc.DrawLine(centerX-symbolLength, centerY+symbolLength, centerX+symbolLength, centerY-symbolLength)

	dc.Stroke()

	return dc.Image()
}

func DrawMenu(width, height int, title string) image.Image {
	dc := gg.NewContext(width, height)

	dc.SetHexColor(gridBackgroundColor)
	dc.Clear()

	// Load the font
	if err := dc.LoadFontFace("arial.ttf", 48); err != nil {
		log.Println("warning, couldn't load the font")
	}

	// Game title
	dc.SetHexColor("#2C3E50")
	dc.DrawStringAnchored(title, float64(width/2), float64(height)/5, 0.5, 0.5)

	return dc.Image()
}
