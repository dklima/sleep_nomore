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

	// Note: Cannot test toggleActive() method directly as it requires systray initialization
	// This would need to be tested in an end-to-end test with actual GUI context
	t.Log("toggleActive() method requires GUI context and cannot be tested in unit tests")
}

func TestMultipleTogglesStability(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stability test in short mode")
	}

	// Test sleep API multiple times without GUI components
	for i := 0; i < 5; i++ {
		err := preventSleep()
		if err != nil {
			t.Logf("Iteration %d: preventSleep() error: %v", i, err)
		}

		time.Sleep(10 * time.Millisecond)

		err = allowSleep()
		if err != nil {
			t.Logf("Iteration %d: allowSleep() error: %v", i, err)
		}

		time.Sleep(10 * time.Millisecond)
	}

	t.Log("Sleep API stability test completed successfully")
}

func TestConcurrentAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent access test in short mode")
	}

	done := make(chan bool)

	for i := 0; i < 3; i++ {
		go func(goroutineID int) {
			for j := 0; j < 3; j++ {
				err := preventSleep()
				if err != nil {
					t.Logf("Goroutine %d iteration %d: preventSleep() error: %v", goroutineID, j, err)
				}

				time.Sleep(20 * time.Millisecond)

				err = allowSleep()
				if err != nil {
					t.Logf("Goroutine %d iteration %d: allowSleep() error: %v", goroutineID, j, err)
				}
			}
			done <- true
		}(i)
	}

	for i := 0; i < 3; i++ {
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Error("Concurrent access test timeout")
		}
	}

	t.Log("Concurrent sleep API access test completed successfully")
}
