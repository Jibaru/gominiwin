//go:build windows

package main

import (
	"time"

	"github.com/jibaru/gominiwin/colors"
	"github.com/jibaru/gominiwin/windows"
)

type Point struct {
	x, y float32
}

type Polyline []Point

type Icon struct {
	Point
	size float32
}

func main() {
	w, err := windows.New("Paint", 1024, 768)
	if err != nil {
		panic(err)
	}

	polylines := []Polyline{}
	var currPolyline Polyline
	icon := Icon{Point: Point{x: 0, y: 0}, size: 10}

	go func() {
		for {
			inside, x, y := w.MouseState()

			if inside && w.MouseLeftClicked() {
				pos := Point{x, y}
				lenght := len(currPolyline)
				if lenght == 0 || currPolyline[lenght-1] != pos {
					currPolyline = append(currPolyline, pos)
				}
			} else if len(currPolyline) > 0 {
				polylines = append(polylines, currPolyline)
				currPolyline = make(Polyline, 0)
			}

			w.Clear()

			w.SetColor(colors.White)
			w.FilledRectangle(0, 0, float32(w.Width()), float32(w.Height()))

			icon.x = x
			icon.y = y

			// Draw icon
			w.SetColor(colors.Black)
			w.Line(icon.x-icon.size, icon.y, icon.x+icon.size, icon.y)
			w.Line(icon.x, icon.y-icon.size, icon.x, icon.y+icon.size)

			// Draw past polylines
			for _, line := range polylines {
				for i := 1; i < len(line); i++ {
					prev := line[i-1]
					curr := line[i]
					w.Line(prev.x, prev.y, curr.x, curr.y)
				}
			}

			// Draw current polyline only if we are currently painting it
			if len(currPolyline) > 1 {
				for i := 1; i < len(currPolyline); i++ {
					prev := currPolyline[i-1]
					curr := currPolyline[i]
					w.Line(prev.x, prev.y, curr.x, curr.y)
				}
			}

			w.Refresh()
			time.Sleep(15 * time.Millisecond)
		}
	}()

	w.Start()
}
