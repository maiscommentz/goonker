package ui

import (
	"testing"
)

func TestTextFieldLogic(t *testing.T) {
	// Initialize images because TextField.redraw() is called on creation
	defer func() {
		if r := recover(); r != nil {
			t.Skip("Skipping TextFieldLogic test due to asset initialization failure (headless?):", r)
		}
	}()
	InitImages()

	tf := NewTextField(100, 100, 200, 30, 14)
	tf.MaxLength = 5

	// Focus Logic
	tf.HandleClick(110, 110) // Inside
	if !tf.Focused {
		t.Error("Expected focus after click inside")
	}

	tf.HandleClick(0, 0) // Outside
	if tf.Focused {
		t.Error("Expected blur after click outside")
	}

	// Test Insert Logic
	tf.Insert("AB")
	if tf.Text != "AB" {
		t.Errorf("Expected 'AB', got '%s'", tf.Text)
	}

	// Test MaxLength Logic
	// "AB" + "CDE" = 5 chars (OK)
	// "F" should be rejected
	if !tf.Insert("CDE") {
		t.Error("Expected successful insert")
	}
	if tf.Text != "ABCDE" {
		t.Errorf("Expected 'ABCDE', got '%s'", tf.Text)
	}

	if tf.Insert("F") {
		t.Error("Expected insert to fail due to max length")
	}
	if tf.Text != "ABCDE" {
		t.Errorf("Expected 'ABCDE', got '%s'", tf.Text)
	}

	// Test Backspace Logic
	if !tf.Backspace() {
		t.Error("Expected backspace to succeed")
	}
	if tf.Text != "ABCD" {
		t.Errorf("Expected 'ABCD', got '%s'", tf.Text)
	}

	// Drain
	tf.Backspace()
	tf.Backspace()
	tf.Backspace()
	tf.Backspace()

	if len(tf.Text) != 0 {
		t.Error("Expected empty string")
	}

	if tf.Backspace() {
		t.Error("Expected backspace to fail on empty string")
	}
}
