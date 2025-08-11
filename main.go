// Package main implements a Windows system tray application that prevents
// the computer from entering sleep mode.
package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/getlantern/systray"
	"golang.org/x/sys/windows"
)

//go:embed assets/icon_green.ico
//go:embed assets/icon_red.ico
var assets embed.FS

var (
	kernel32           = windows.NewLazySystemDLL("kernel32.dll")
	setThreadExecState = kernel32.NewProc("SetThreadExecutionState")
)

// Windows API constants for SetThreadExecutionState
const (
	// ES_CONTINUOUS - Informs the system that the state being set should remain in effect until
	// the next call that uses ES_CONTINUOUS and one of the other state flags is cleared.
	ES_CONTINUOUS = 0x80000000
	// ES_SYSTEM_REQUIRED - Forces the system to be in the working state by resetting the system idle timer.
	ES_SYSTEM_REQUIRED = 0x00000001
	// ES_DISPLAY_REQUIRED - Forces the display to be on by resetting the display idle timer.
	ES_DISPLAY_REQUIRED = 0x00000002
)

// App represents the main application state and UI elements.
type App struct {
	active     bool              // Current state of sleep prevention
	menuToggle *systray.MenuItem // Menu item for toggling sleep prevention
	menuQuit   *systray.MenuItem // Menu item for quitting the application
}

func main() {
	app := &App{
		active: false,
	}
	systray.Run(app.onReady, app.onExit)
}

func (app *App) onReady() {
	iconRed, err := assets.ReadFile("assets/icon_red.ico")
	if err != nil {
		log.Printf("Error loading red icon: %v", err)
		iconRed = generateDefaultIcon(false)
	}

	iconGreen, err := assets.ReadFile("assets/icon_green.ico")
	if err != nil {
		log.Printf("Error loading green icon: %v", err)
		iconGreen = generateDefaultIcon(true)
	}

	systray.SetIcon(iconRed)
	systray.SetTitle("Sleep No More")
	systray.SetTooltip("Click to toggle awake mode")

	app.menuToggle = systray.AddMenuItem("Enable", "Enable awake mode")
	systray.AddSeparator()
	app.menuQuit = systray.AddMenuItem("Quit", "Quit the application")

	go func() {
		for {
			select {
			case <-app.menuToggle.ClickedCh:
				app.toggleActive(iconGreen, iconRed)
			case <-app.menuQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()

	go app.handleSignals()
}

func (app *App) toggleActive(iconGreen, iconRed []byte) {
	app.active = !app.active

	if app.active {
		if err := preventSleep(); err != nil {
			log.Printf("Error preventing sleep: %v", err)
			app.active = false
			return
		}
		systray.SetIcon(iconGreen)
		app.menuToggle.SetTitle("Disable")
		app.menuToggle.SetTooltip("Disable awake mode")
		systray.SetTooltip("Sleep No More - ACTIVE")
	} else {
		if err := allowSleep(); err != nil {
			log.Printf("Error allowing sleep: %v", err)
		}
		systray.SetIcon(iconRed)
		app.menuToggle.SetTitle("Enable")
		app.menuToggle.SetTooltip("Enable awake mode")
		systray.SetTooltip("Sleep No More - INACTIVE")
	}
}

func preventSleep() error {
	ret, _, err := setThreadExecState.Call(
		uintptr(ES_CONTINUOUS | ES_SYSTEM_REQUIRED | ES_DISPLAY_REQUIRED),
	)
	if ret == 0 {
		return err
	}
	return nil
}

func allowSleep() error {
	// Clear the previous state by calling ES_CONTINUOUS without any other flags
	// This removes the ES_SYSTEM_REQUIRED and ES_DISPLAY_REQUIRED flags
	ret, _, err := setThreadExecState.Call(
		uintptr(ES_CONTINUOUS),
	)
	if ret == 0 {
		return fmt.Errorf("failed to clear execution state: %v", err)
	}
	
	// Additionally, reset to default state (allow system to sleep)
	ret, _, err = setThreadExecState.Call(0)
	if ret == 0 {
		return fmt.Errorf("failed to reset execution state: %v", err)
	}
	
	return nil
}

func (app *App) onExit() {
	if app.active {
		_ = allowSleep()
	}
}

func (app *App) handleSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	systray.Quit()
}

func generateDefaultIcon(green bool) []byte {
	icoHeader := []byte{
		0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x10, 0x10,
		0x00, 0x00, 0x01, 0x00, 0x20, 0x00, 0x68, 0x04,
		0x00, 0x00, 0x16, 0x00, 0x00, 0x00,
	}

	bmpHeader := []byte{
		0x28, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00,
		0x20, 0x00, 0x00, 0x00, 0x01, 0x00, 0x20, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	var color [4]byte
	if green {
		color = [4]byte{0x00, 0xFF, 0x00, 0xFF}
	} else {
		color = [4]byte{0x00, 0x00, 0xFF, 0xFF}
	}

	pixels := make([]byte, 1024)
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			dx := x - 8
			dy := y - 8
			if dx*dx+dy*dy <= 49 {
				idx := (y*16 + x) * 4
				copy(pixels[idx:idx+4], color[:])
			}
		}
	}

	result := make([]byte, 0, len(icoHeader)+len(bmpHeader)+len(pixels))
	result = append(result, icoHeader...)
	result = append(result, bmpHeader...)
	result = append(result, pixels...)

	return result
}

func init() {
	if runtime := os.Getenv("GOOS"); runtime != "" && runtime != "windows" {
		log.Fatal("This program only works on Windows")
	}
}
