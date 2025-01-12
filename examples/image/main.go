//go:build windows

package main

import (
	"image"
	_ "image/jpeg"
	"os"
	"time"

	"github.com/jibaru/gominiwin"
)

type RGB struct {
	R, G, B int
}

func main() {
	file, err := os.Open("image.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	pixels := make([][]RGB, height)
	for y := 0; y < height; y++ {
		pixels[y] = make([]RGB, width)
		for x := 0; x < width; x++ {
			color := img.At(x, y)
			r, g, b, _ := color.RGBA()

			pixels[y][x] = RGB{
				R: int(r >> 8),
				G: int(g >> 8),
				B: int(b >> 8),
			}
		}
	}

	w, err := gominiwin.NewWindowsWin("Image", width, height)
	if err != nil {
		panic(err)
	}

	go func() {
		time.Sleep(1 * time.Second)

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				rgb := pixels[y][x]

				w.SetColorRGB(rgb.R, rgb.G, rgb.B)
				w.Point(float32(x), float32(y))
			}
		}

		w.Refresh()
	}()

	w.Start()
}
