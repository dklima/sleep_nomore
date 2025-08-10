package main

import (
	"bytes"
	"testing"
)

func TestGenerateDefaultIcon(t *testing.T) {
	tests := []struct {
		name     string
		isGreen  bool
		wantSize int
	}{
		{
			name:     "Generate green icon",
			isGreen:  true,
			wantSize: 1086, // ICO header (22) + BMP header (40) + pixel data (1024)
		},
		{
			name:     "Generate red icon",
			isGreen:  false,
			wantSize: 1086,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			icon := generateDefaultIcon(tt.isGreen)

			if len(icon) != tt.wantSize {
				t.Errorf("generateDefaultIcon() returned icon of size %d, want %d", len(icon), tt.wantSize)
			}

			if icon[0] != 0x00 || icon[1] != 0x00 || icon[2] != 0x01 || icon[3] != 0x00 {
				t.Error("Invalid ICO header signature")
			}

			if icon[4] != 0x01 || icon[5] != 0x00 {
				t.Error("Invalid icon count in header")
			}

			if icon[6] != 0x10 || icon[7] != 0x10 {
				t.Error("Invalid icon dimensions")
			}
		})
	}
}

func TestGenerateDefaultIconColors(t *testing.T) {
	greenIcon := generateDefaultIcon(true)
	redIcon := generateDefaultIcon(false)

	if bytes.Equal(greenIcon, redIcon) {
		t.Error("Green and red icons should be different")
	}

	headerSize := 22 + 40 // ICO header + BMP header
	greenPixels := greenIcon[headerSize:]
	redPixels := redIcon[headerSize:]

	if len(greenPixels) != 1024 || len(redPixels) != 1024 {
		t.Error("Invalid pixel data size")
	}

	centerOffset := (8*16 + 8) * 4 // Center pixel at (8,8)

	if greenPixels[centerOffset+1] != 0xFF {
		t.Error("Green icon center pixel should have green channel set")
	}

	if redPixels[centerOffset+2] != 0xFF {
		t.Error("Red icon center pixel should have red channel set")
	}
}

func TestAppStructInitialization(t *testing.T) {
	app := &App{
		active: false,
	}

	if app.active {
		t.Error("App should initialize with active = false")
	}

	if app.menuToggle != nil {
		t.Error("menuToggle should be nil before initialization")
	}

	if app.menuQuit != nil {
		t.Error("menuQuit should be nil before initialization")
	}
}

func BenchmarkGenerateDefaultIcon(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = generateDefaultIcon(i%2 == 0)
	}
}

func TestIconPixelGeneration(t *testing.T) {
	icon := generateDefaultIcon(true)
	headerSize := 22 + 40
	pixels := icon[headerSize:]

	coloredPixels := 0
	for i := 0; i < len(pixels); i += 4 {
		if pixels[i+3] == 0xFF { // Check alpha channel
			coloredPixels++
		}
	}

	if coloredPixels < 100 || coloredPixels > 200 {
		t.Errorf("Unexpected number of colored pixels: %d", coloredPixels)
	}
}
