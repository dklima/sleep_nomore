//go:build windows && integration
// +build windows,integration

package main

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestSignalHandling(t *testing.T) {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	go func() {
		time.Sleep(100 * time.Millisecond)
		sigChan <- os.Interrupt
	}()

	select {
	case sig := <-sigChan:
		if sig != os.Interrupt {
			t.Errorf("Expected os.Interrupt, got %v", sig)
		}
	case <-time.After(1 * time.Second):
		t.Error("Signal handling timeout")
	}
}

func TestAppToggleState(t *testing.T) {
	app := &App{active: false}

	iconGreen := []byte{0x01, 0x02, 0x03}
	iconRed := []byte{0x04, 0x05, 0x06}

	if app.active {
		t.Error("App should start inactive")
	}

	app.active = !app.active
	if !app.active {
		t.Error("App should be active after toggle")
	}

	if app.active {
		err := preventSleep()
		if err != nil {
			t.Logf("Warning: Could not prevent sleep: %v", err)
		}
	}

	app.active = !app.active
	if app.active {
		t.Error("App should be inactive after second toggle")
	}

	if !app.active {
		err := allowSleep()
		if err != nil {
			t.Logf("Warning: Could not allow sleep: %v", err)
		}
	}

	app.toggleActive(iconGreen, iconRed)
	if !app.active {
		t.Error("App should be active after toggleActive call")
	}

	app.toggleActive(iconGreen, iconRed)
	if app.active {
		t.Error("App should be inactive after second toggleActive call")
	}
}

func TestMultipleTogglesStability(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stability test in short mode")
	}

	app := &App{active: false}
	iconGreen := generateDefaultIcon(true)
	iconRed := generateDefaultIcon(false)

	for i := 0; i < 10; i++ {
		app.toggleActive(iconGreen, iconRed)
		time.Sleep(10 * time.Millisecond)
	}

	expectedActive := (10 % 2) != 0
	if app.active != expectedActive {
		t.Errorf("After 10 toggles, active should be %v, got %v", expectedActive, app.active)
	}

	if app.active {
		allowSleep()
	}
}

func TestConcurrentAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent access test in short mode")
	}

	app := &App{active: false}
	iconGreen := generateDefaultIcon(true)
	iconRed := generateDefaultIcon(false)

	done := make(chan bool)

	for i := 0; i < 3; i++ {
		go func() {
			for j := 0; j < 5; j++ {
				app.toggleActive(iconGreen, iconRed)
				time.Sleep(20 * time.Millisecond)
			}
			done <- true
		}()
	}

	for i := 0; i < 3; i++ {
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Error("Concurrent access test timeout")
		}
	}

	if app.active {
		allowSleep()
	}
}
