//go:build windows

package windows

import (
	"github.com/jibaru/gominiwin"
	"github.com/jibaru/gominiwin/colors"
)

type window struct{}

var _ gominiwin.Win = &window{}
var windowInstance *window = &window{}

func (w *window) Start() {
	Start()
}

func (w *window) Point(x, y float32) {
	Point(x, y)
}
func (w *window) Line(x1, y1, x2, y2 float32) {
	Line(x1, y1, x2, y2)
}
func (w *window) Rectangle(left, top, right, bottom float32) {
	Rectangle(left, top, right, bottom)
}

func (w *window) FilledRectangle(left, top, right, bottom float32) {
	FilledRectangle(left, top, right, bottom)
}
func (w *window) Circle(centerX, centerY, radius float32) {
	Circle(centerX, centerY, radius)
}
func (w *window) FilledCircle(centerX, centerY, radius float32) {
	FilledCircle(centerX, centerY, radius)
}
func (w *window) SetColor(c colors.Color) {
	SetColor(c)
}
func (w *window) SetColorRGB(r, g, b int) {
	SetColorRGB(r, g, b)
}

func (w *window) SetText(x, y float32, content string) {
	SetText(x, y, content)
}

func (w *window) KeyPressed() int {
	return KeyPressed()
}
func (w *window) MouseState() (bool, float32, float32) {
	return MouseState()
}
func (w *window) IsMouseInside() bool {
	return IsMouseInside()
}
func (w *window) MouseX() float32            { return MouseX() }
func (w *window) MouseY() float32            { return MouseY() }
func (w *window) MouseButtons() (bool, bool) { return MouseButtons() }
func (w *window) MouseLeftClicked() bool     { return MouseLeftClicked() }
func (w *window) MouseRightClicked() bool    { return MouseRightClicked() }

func (w *window) Clear()                         { Clear() }
func (w *window) Refresh()                       { Refresh() }
func (w *window) Width() int                     { return Width() }
func (w *window) Height() int                    { return Height() }
func (w *window) Resize(newWidth, newHeight int) { Resize(newWidth, newHeight) }
func (w *window) Close()                         { Close() }
