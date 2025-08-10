//go:build windows
// +build windows

package main

import (
	"testing"
	"time"
)

func TestPreventSleep(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping sleep prevention test in short mode")
	}

	err := preventSleep()
	if err != nil {
		t.Logf("Warning: preventSleep() returned error: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	err = allowSleep()
	if err != nil {
		t.Logf("Warning: allowSleep() returned error: %v", err)
	}
}

func TestAllowSleep(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping sleep allowance test in short mode")
	}

	_ = preventSleep()
	time.Sleep(100 * time.Millisecond)

	err := allowSleep()
	if err != nil {
		t.Logf("Warning: allowSleep() returned error: %v", err)
	}
}

func TestSleepToggleSequence(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping sleep toggle sequence test in short mode")
	}

	for i := 0; i < 3; i++ {
		err := preventSleep()
		if err != nil {
			t.Logf("Iteration %d: preventSleep() error: %v", i, err)
		}

		time.Sleep(50 * time.Millisecond)

		err = allowSleep()
		if err != nil {
			t.Logf("Iteration %d: allowSleep() error: %v", i, err)
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func TestAppOnExit(t *testing.T) {
	app := &App{active: false}

	app.onExit()

	app.active = true
	app.onExit()
}
