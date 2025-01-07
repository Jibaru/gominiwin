package gominiwin

type Color = uint32

const (
	BlackColor   Color = 0x000000
	RedColor     Color = 0xFF0000
	GreenColor   Color = 0x0000FF
	YellowColor  Color = 0xFFFF00
	MagentaColor Color = 0x00FFFF
	WhiteColor   Color = 0xFFFFFF
)

type MouseState struct {
	X, Y                      float64
	LeftClicked, RightClicked bool
}

type Paint interface {
	SetBackground(color Color) error
	DrawFilledRectangle(left, top, right, bottom int32, color Color) error
	DrawUnfilledRectangle(left, top, right, bottom int32, color Color) error
	DrawPoint(x, y int32, color Color) error
	DrawLine(x1, y1, x2, y2 int32, color uint32) error
	DrawText(x, y int32, text string) error
	DrawFilledCircle(left, top, right, bottom int32, color Color) error
	DrawUnfilledCircle(left, top, right, bottom int32, color Color) error
}

type Input interface {
	Key() int
	MouseState() (MouseState, error)
	IsMouseInsideWindow() (bool, error)
}

type Window interface {
	Run(onPaint func(Paint), onInput func(Input))
	Refresh()
	ShowMessage(title, message string)
}
