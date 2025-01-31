//go:build linux

package linux

import (
	"fmt"
	"math"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"

	"github.com/jibaru/gominiwin/colors"
	"github.com/jibaru/gominiwin/errors"
	"github.com/jibaru/gominiwin/keys"
)

var colorsRefs = []uint32{
	0x000000, // black
	0xFF0000, // red
	0x00FF00, // green
	0x0000FF, // blue
	0xFFFF00, // yellow
	0xFF00FF, // magenta
	0x00FFFF, // cyan
	0xFFFFFF, // white
}

type window struct {
	conn           *xgb.Conn
	screen         *xproto.ScreenInfo
	win            xproto.Window
	gc             xproto.Gcontext
	width, height  int
	color          uint32
	mouseX, mouseY float32
	mouseInside    bool
	mouseLeft      bool
	mouseRight     bool
	keysPressed    []int
	buffer         xproto.Pixmap
}

func New(title string, width, height int) (*window, error) {
	conn, err := xgb.NewConn()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to connect to X server: %w", errors.ErrCreateWinFailed, err)
	}

	screen := xproto.Setup(conn).DefaultScreen(conn)

	win, err := xproto.NewWindowId(conn)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("%w: failed to create window ID: %w", errors.ErrCreateWinFailed, err)
	}

	// Create the graphics context
	gc, err := xproto.NewGcontextId(conn)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("%w: failed to create graphics context: %w", errors.ErrCreateWinFailed, err)
	}

	buffer, err := xproto.NewPixmapId(conn)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("%w: failed to create pixmap: %w", errors.ErrCreateWinFailed, err)
	}

	xproto.CreateGC(conn, gc, xproto.Drawable(screen.Root), xproto.GcForeground, []uint32{screen.BlackPixel})
	xproto.CreatePixmap(conn, screen.RootDepth, buffer, xproto.Drawable(screen.Root), uint16(width), uint16(height))

	// Create the window
	xproto.CreateWindow(
		conn,
		screen.RootDepth,
		win,
		screen.Root,
		0, 0,
		uint16(width), uint16(height),
		10,
		xproto.WindowClassInputOutput,
		screen.RootVisual,
		xproto.CwBackPixel|xproto.CwEventMask,
		[]uint32{screen.WhitePixel, xproto.EventMaskExposure | xproto.EventMaskKeyPress | xproto.EventMaskPointerMotion | xproto.EventMaskButtonPress | xproto.EventMaskButtonRelease},
	)

	// Set the window title
	xproto.ChangeProperty(
		conn,
		xproto.PropModeReplace,
		win,
		xproto.AtomWmName,
		xproto.AtomString,
		8,
		uint32(len(title)),
		[]byte(title),
	)

	return &window{
		conn:   conn,
		screen: screen,
		win:    win,
		gc:     gc,
		width:  width,
		height: height,
		color:  screen.BlackPixel,
		buffer: buffer,
	}, nil
}

func (w *window) Start() {
	// Map the window (make it visible)
	xproto.MapWindow(w.conn, w.win)

	for {
		event, err := w.conn.WaitForEvent()
		if err != nil {
			fmt.Printf("Error waiting for event: %v\n", err)
			break
		}

		switch e := event.(type) {
		case xproto.KeyPressEvent:
			fmt.Println("key", e.Detail)
			ascii := int(e.Detail)
			w.keysPressed = append(w.keysPressed, ascii)
		case xproto.MotionNotifyEvent:
			w.mouseX, w.mouseY = float32(e.EventX), float32(e.EventY)
			w.mouseInside = e.EventX >= 0 && e.EventY >= 0 && int(e.EventX) < w.width && int(e.EventY) < w.height
		case xproto.ButtonPressEvent:
			if e.Detail == 1 {
				w.mouseLeft = true
			} else if e.Detail == 3 {
				w.mouseRight = true
			}
		case xproto.ButtonReleaseEvent:
			if e.Detail == 1 {
				w.mouseLeft = false
			} else if e.Detail == 3 {
				w.mouseRight = false
			}
		}
	}
}

func (w *window) Clear() {
	xproto.ClearArea(w.conn, true, w.win, 0, 0, uint16(w.width), uint16(w.height))
}

func (w *window) Refresh() {
	xproto.CopyArea(w.conn, xproto.Drawable(w.buffer), xproto.Drawable(w.win), w.gc, 0, 0, 0, 0, uint16(w.width), uint16(w.height))

}

func (w *window) Width() int {
	return w.width
}

func (w *window) Height() int {
	return w.height
}

func (w *window) Resize(newWidth, newHeight int) {
	w.width = newWidth
	w.height = newHeight

	// Recreate the pixmap buffer with new dimensions
	xproto.FreePixmap(w.conn, w.buffer)
	newBuffer, err := xproto.NewPixmapId(w.conn)
	if err != nil {
		fmt.Printf("Error creating new pixmap: %v\n", err)
		return
	}
	w.buffer = newBuffer
	xproto.CreatePixmap(w.conn, w.screen.RootDepth, w.buffer, xproto.Drawable(w.win), uint16(newWidth), uint16(newHeight))

	xproto.ConfigureWindow(w.conn, w.win, xproto.ConfigWindowWidth|xproto.ConfigWindowHeight, []uint32{uint32(newWidth), uint32(newHeight)})
}

func (w *window) Close() {
	xproto.FreePixmap(w.conn, w.buffer)
	w.conn.Close()
}

func (w *window) SetColor(c colors.Color) {
	if c >= 0 && c < len(colorsRefs) {
		w.color = colorsRefs[c]
		xproto.ChangeGC(w.conn, w.gc, xproto.GcForeground, []uint32{w.color})
	}
}

func (w *window) SetColorRGB(r, g, b int) {
	w.SetColor((r&0xFF)<<16 | (g&0xFF)<<8 | (b & 0xFF))
}

func (w *window) SetText(x, y float32, content string) {
	// Convert float coordinates to integers
	intX := int16(x)
	intY := int16(y)

	// Calculate the length of the string in bytes
	length := uint32(len(content))

	// Create an X11 font ID
	font, _ := xproto.NewFontId(w.conn)

	// Load a font
	_ = xproto.OpenFontChecked(w.conn, font, uint16(len("fixed")), "fixed").Check()

	// Set font in graphics context
	xproto.ChangeGC(w.conn, w.gc, xproto.GcFont, []uint32{uint32(font)})

	// Draw text
	xproto.ImageText8(w.conn, byte(length), xproto.Drawable(w.buffer), w.gc, intX, intY, content)
	xproto.CloseFont(w.conn, font)
}

func (w *window) applyColor() {
	xproto.ChangeGC(w.conn, w.gc, xproto.GcForeground, []uint32{w.color})
}

func (w *window) Point(x, y float32) {
	w.applyColor()
	xproto.PolyPoint(w.conn, xproto.CoordModeOrigin, xproto.Drawable(w.buffer), w.gc, []xproto.Point{{X: int16(x), Y: int16(y)}})
}

func (w *window) Line(x1, y1, x2, y2 float32) {
	w.applyColor()
	xproto.PolyLine(w.conn, xproto.CoordModeOrigin, xproto.Drawable(w.buffer), w.gc, []xproto.Point{
		{X: int16(x1), Y: int16(y1)},
		{X: int16(x2), Y: int16(y2)},
	})
}

func (w *window) Rectangle(left, top, right, bottom float32) {
	w.applyColor()
	xproto.PolyRectangle(w.conn, xproto.Drawable(w.buffer), w.gc, []xproto.Rectangle{
		{X: int16(left), Y: int16(top), Width: uint16(right - left), Height: uint16(bottom - top)},
	})
}

func (w *window) FilledRectangle(left, top, right, bottom float32) {
	w.applyColor()
	xproto.PolyFillRectangle(w.conn, xproto.Drawable(w.buffer), w.gc, []xproto.Rectangle{
		{X: int16(left), Y: int16(top), Width: uint16(right - left), Height: uint16(bottom - top)},
	})
}

func (w *window) Circle(centerX, centerY, radius float32) {
	w.applyColor()
	segments := []xproto.Segment{}
	for angle := 0.0; angle < 360; angle += 5 {
		x1 := centerX + radius*float32(math.Cos(angle*math.Pi/180))
		y1 := centerY + radius*float32(math.Sin(angle*math.Pi/180))
		x2 := centerX + radius*float32(math.Cos((angle+5)*math.Pi/180))
		y2 := centerY + radius*float32(math.Sin((angle+5)*math.Pi/180))
		segments = append(segments, xproto.Segment{X1: int16(x1), Y1: int16(y1), X2: int16(x2), Y2: int16(y2)})
	}
	xproto.PolySegment(w.conn, xproto.Drawable(w.buffer), w.gc, segments)
}

func (w *window) FilledCircle(centerX, centerY, radius float32) {
	w.applyColor()
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				w.Point(centerX+x, centerY+y)
			}
		}
	}
}

func (w *window) KeyPressed() int {
	if len(w.keysPressed) == 0 {
		return keys.None
	}

	var ret int
	key := w.keysPressed[0]

	switch key {
	case 9: // Escape key
		ret = keys.Escape
	case 36: // Return key
		ret = keys.Return
	case 65: // Space key
		ret = keys.Space
	case 111: // Arrow Up
		ret = keys.Up
	case 116: // Arrow Down
		ret = keys.Down
	case 113: // Arrow Left
		ret = keys.Left
	case 114: // Arrow Right
		ret = keys.Right
	case 67: // F1
		ret = keys.F1
	case 68: // F2
		ret = keys.F2
	case 69: // F3
		ret = keys.F3
	case 70: // F4
		ret = keys.F4
	case 71: // F5
		ret = keys.F5
	case 72: // F6
		ret = keys.F6
	case 73: // F7
		ret = keys.F7
	case 74: // F8
		ret = keys.F8
	case 75: // F9
		ret = keys.F9
	case 76: // F10
		ret = keys.F10
	default:
		key = ret
	}

	w.keysPressed = w.keysPressed[1:]

	return ret
}

func (w *window) MouseState() (bool, float32, float32) {
	return w.mouseInside, w.mouseX, w.mouseY
}

func (w *window) IsMouseInside() bool {
	return w.mouseInside
}

func (w *window) MouseX() float32 {
	return w.mouseX
}

func (w *window) MouseY() float32 {
	return w.mouseY
}

func (w *window) MouseButtons() (bool, bool) {
	return w.mouseLeft, w.mouseRight
}

func (w *window) MouseLeftClicked() bool {
	return w.mouseLeft
}

func (w *window) MouseRightClicked() bool {
	return w.mouseRight
}
