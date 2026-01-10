package ui

import (
	"testing"
	"time"
)

func TestTimerLogic(t *testing.T) {
	// Create a timer of 1 second
	timer := NewTimer(time.Second)

	if !timer.IsRunning {
		t.Error("New timer should be running")
	}

	// Ratio should be 1.0 initially
	if r := timer.Ratio(); r != 1.0 {
		t.Errorf("Expected ratio 1.0, got %f", r)
	}

	// Update (simulate 1 tick)
	timer.Update()

	if timer.CurrentDuration >= timer.TotalDuration {
		t.Error("Timer did not decrease duration")
	}

	// Force expiration
	timer.CurrentDuration = 5 * time.Millisecond // very small
	timer.Update()                               // should expire approx

	// Force it to 0
	timer.CurrentDuration = 0

	// Reset
	callbackCalled := false
	timer = NewTimer(time.Second)
	timer.OnEnd = func() {
		callbackCalled = true
	}

	// Fast forward
	// 60 ticks should empty it
	for i := 0; i < 62; i++ {
		timer.Update()
	}

	if !callbackCalled {
		t.Error("OnEnd callback was not called")
	}
	if timer.IsRunning {
		t.Error("Timer should not be running after end")
	}
	if timer.CurrentDuration != 0 {
		t.Errorf("Duration should be 0, got %v", timer.CurrentDuration)
	}
}
