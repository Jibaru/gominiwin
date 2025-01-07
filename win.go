package gominiwin

type Canva interface {
	Point(x, y float32)
	Line(x1, y1, x2, y2 float32)
	Rectangle(left, top, right, bottom float32)
	FilledRectangle(left, top, right, bottom float32)
	Circle(centerX, centerY, radius float32)
	FilledCircle(centerX, centerY, radius float32)
	SetColor(c int)
	SetColorRGB(r, g, b int)
}

type Input interface {
	KeyPressed() int
	MouseState() (bool, float32, float32)
	IsMouseInside() bool
	MouseX() float32
	MouseY() float32
	MouseButtons() (bool, bool)
	MouseLeftClicked() bool
	MouseRightClicked() bool
}

type Win interface {
	Input
	Canva
	Start()
	Clear()
	Refresh()
	Width() int
	Height() int
	Resize(newWidth, newHeight int)
	Close()
}
