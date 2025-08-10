package main

import (
	"embed"
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

const (
	ES_CONTINUOUS       = 0x80000000
	ES_SYSTEM_REQUIRED  = 0x00000001
	ES_DISPLAY_REQUIRED = 0x00000002
)

type App struct {
	active     bool
	menuToggle *systray.MenuItem
	menuQuit   *systray.MenuItem
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
	ret, _, err := setThreadExecState.Call(
		uintptr(ES_CONTINUOUS),
	)
	if ret == 0 {
		return err
	}
	return nil
}

func (app *App) onExit() {
	if app.active {
		allowSleep()
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