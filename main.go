package main

import (
	"time"

	"github.com/jibaru/gominiwin/windows"
)

func main() {
	var a, b, c, d float32 = 0.0, 0.0, 30.0, 30.0
	w, h := 800, 600
	go func() {
		for {
			windows.Resize(w, h)

			k := windows.KeyPressed()
			if k == 'A' {
				windows.Close()
			}

			windows.Clear()
			windows.SetColor(windows.BLUE)
			windows.FilledRectangle(a, b, c, d)

			windows.Refresh()

			a += 10
			b += 10
			c += 10
			d += 10
			w += 1
			h += 1

			time.Sleep(60 * time.Millisecond)
		}
	}()

	windows.Start()
}
