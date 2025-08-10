//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	generateIcon("assets/icon_green.ico", color.RGBA{0, 255, 0, 255})
	generateIcon("assets/icon_red.ico", color.RGBA{255, 0, 0, 255})
}

func generateIcon(filename string, c color.Color) {
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			dx := float64(x) - 7.5
			dy := float64(y) - 7.5
			if dx*dx+dy*dy <= 49 {
				img.Set(x, y, c)
			}
		}
	}

	buf := new(bytes.Buffer)
	png.Encode(buf, img)
	pngData := buf.Bytes()

	icoFile, _ := os.Create(filename)
	defer icoFile.Close()

	binary.Write(icoFile, binary.LittleEndian, uint16(0))
	binary.Write(icoFile, binary.LittleEndian, uint16(1))
	binary.Write(icoFile, binary.LittleEndian, uint16(1))

	binary.Write(icoFile, binary.LittleEndian, uint8(16))
	binary.Write(icoFile, binary.LittleEndian, uint8(16))
	binary.Write(icoFile, binary.LittleEndian, uint8(0))
	binary.Write(icoFile, binary.LittleEndian, uint8(0))
	binary.Write(icoFile, binary.LittleEndian, uint16(0))
	binary.Write(icoFile, binary.LittleEndian, uint16(32))
	binary.Write(icoFile, binary.LittleEndian, uint32(len(pngData)))
	binary.Write(icoFile, binary.LittleEndian, uint32(22))

	icoFile.Write(pngData)
}