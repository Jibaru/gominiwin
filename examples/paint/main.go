package main

import (
	"time"

	"github.com/jibaru/gominiwin"
	"github.com/jibaru/gominiwin/colors"
	"github.com/jibaru/gominiwin/keys"
)

type Brush struct {
	x, y float32
	size float32
}

func (b Brush) Draw(w gominiwin.Win) {
	w.FilledRectangle(b.x, b.y, b.x+b.size, b.y+b.size)
}

func main() {
	w, err := gominiwin.New("Paint", 800, 600)
	if err != nil {
		panic(err)
	}

	brush := Brush{x: 0, y: 0, size: 40}
	brushColors := []colors.Color{
		colors.Black,
		colors.Red,
		colors.Blue,
		colors.Cyan,
		colors.Yellow,
		colors.Green,
	}
	colorIdx := 0
	var moveSpeed float32 = 20

	go func() {
		time.Sleep(1 * time.Second)

		w.SetColor(colors.White)
		w.FilledRectangle(0, 0, 800, 600)
		w.SetColor(brushColors[colorIdx])
		brush.Draw(w)

		w.Refresh()

		for {
			key := w.KeyPressed()
			for key == keys.None {
				key = w.KeyPressed()
			}

			switch key {
			case keys.Space:
				colorIdx = (colorIdx + 1) % len(brushColors)
				w.SetColor(brushColors[colorIdx])
			case keys.Up:
				brush.y -= moveSpeed
			case keys.Down:
				brush.y += moveSpeed
			case keys.Left:
				brush.x -= moveSpeed
			case keys.Right:
				brush.x += moveSpeed
			}

			brush.Draw(w)

			w.Refresh()
			time.Sleep(36 * time.Millisecond)
		}
	}()

	w.Start()
}
