package main

import (
	"time"

	"github.com/jibaru/gominiwin"
	"github.com/jibaru/gominiwin/colors"
)

func main() {
	w, err := gominiwin.New("Figures", 800, 600)
	if err != nil {
		panic(err)
	}

	go func() {
		w.SetColor(colors.Blue)
		w.FilledRectangle(0, 0, 800, 600)
		for {
			w.SetColor(colors.Cyan)
			w.FilledCircle(100, 100, 50)

			w.SetColor(colors.Red)
			w.Circle(123, 188, 77)

			w.SetColor(colors.Green)
			w.Line(0, 0, 800, 600)

			w.SetColorRGB(40, 11, 123)
			w.Rectangle(90, 110, 635, 441)

			w.SetColor(colors.White)
			w.SetText(50, 50, "Hello world!")

			w.SetColor(colors.Magenta)
			w.FilledRectangle(500, 400, 700, 600)

			w.SetColor(colors.Yellow)
			for i := 0.0; i < 600; i += 10 {
				w.Point(400, float32(i))
			}

			w.Refresh()

			time.Sleep(36 * time.Millisecond)
		}
	}()

	w.Start()
}
