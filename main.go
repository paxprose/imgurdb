package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 1 {
		readFromImage(os.Args[1])
		return
	}
	generateImage(dir)
}

func generateImage(dir string) {
	f, err := os.ReadFile(dir + "/main.go")
	if err != nil {
		log.Fatal(err)
	}
	encoded := hex.EncodeToString([]byte(f))
	var (
		width  = 25
		height = 27
		x, y   int
	)
	img := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	hexbuf := bytes.Buffer{}
	for i, c := range encoded {
		if i%6 == 0 && i > 0 {
			setColor(hexbuf.String(), x, y, img)
			hexbuf.Reset()
			x++
			if x >= width {
				y++
				x = 0
			}
		}
		hexbuf.WriteRune(c)
	}
	if hexbuf.Len() > 0 {
		setColor(hexbuf.String(), x, y, img)
	}
	file, _ := os.Create("./main.png")
	png.Encode(file, img)
}

func setColor(hexVal string, x, y int, img *image.NRGBA) {
	hexColor := hexVal + strings.Repeat("0", 6-len(hexVal))
	cb, _ := hex.DecodeString(hexColor)
	img.Set(x, y, color.RGBA{R: cb[0], G: cb[1], B: cb[2], A: 255})
}

func readFromImage(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	img, err := png.Decode(f)
	max := img.Bounds().Max
	if err != nil {
		log.Fatal(err)
	}
	outbuf := bytes.Buffer{}
	for y := 0; y < max.Y; y++ {
		for x := 0; x < max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			outbuf.WriteString(hex.EncodeToString([]byte{byte(r), byte(g), byte(b)}))
		}
	}
	out, err := hex.DecodeString(outbuf.String())
	fmt.Printf(string(out))
}
