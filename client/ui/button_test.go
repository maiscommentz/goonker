package ui

import (
	"testing"
)

func TestButtonLogic(t *testing.T) {
	btn := &Button{
		X: 10, Y: 10,
		Width: 100, Height: 50,
	}

	// Hit
	if !btn.Contains(15, 15) {
		t.Error("Expected hit at 15,15")
	}

	// Hit Edge
	if !btn.Contains(10, 10) {
		t.Error("Expected hit at top-left edge")
	}

	// Miss
	if btn.Contains(0, 0) {
		t.Error("Expected miss at 0,0")
	}

	// Miss Y
	if btn.Contains(50, 61) {
		t.Error("Expected miss at 61 Y")
	}
}
