package windows

import (
	"runtime"
	"syscall"
	"unsafe"

	"github.com/jibaru/gominiwin/windows/primitives"
)

const (
	ESCAPE = iota
	LEFT
	RIGHT
	UP
	DOWN
	F1
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	SPACE
	RETURN
	NONE
)

const (
	BLACK = iota
	RED
	GREEN
	BLUE
	YELLOW
	MAGENTA
	CYAN
	WHITE
)

var colorRefs = []primitives.COLORREF{
	primitives.RGB(0, 0, 0),       // BLACK
	primitives.RGB(255, 0, 0),     // RED
	primitives.RGB(0, 255, 0),     // GREEN
	primitives.RGB(0, 0, 255),     // BLUE
	primitives.RGB(255, 255, 0),   // YELLOW
	primitives.RGB(255, 0, 255),   // MAGENTA
	primitives.RGB(0, 255, 255),   // CYAN
	primitives.RGB(255, 255, 255), // WHITE
}

const szClassName = "MiniWin"

var (
	// hWnd is the main window
	hWnd primitives.HWND
	// Bitmap to draw off-screen
	hBitmap primitives.HBITMAP
	// iWidth is the width of the main window
	iWidth int32 = 400
	// iHeight is the height of the main window
	iHeight int32 = 300
	// hDCMem is a in-memory device context
	hDCMem primitives.HDCMEM
	// keysPressed is a queue of keys
	keysPressed = make([]int, 0)
	// mouseInside indicated is mouse is in the 'client area'
	mouseInside bool
	// mouseXPos and mouseYPos indicate the mouse position
	mouseXPos, mouseYPos int32
	// mouseLeftClicked and mouseRightClicked indicate if mouse buttons are pressed
	mouseLeftClicked, mouseRightClicked bool
	// selectedColor is the color to paint anything
	selectedColor primitives.COLORREF = primitives.RGB(255, 255, 255)
)

func init() {
	// needed to prevent window freeze
	runtime.LockOSThread()
}

func realFrame(w, h int32) (rw, rh int32) {
	frame := primitives.RECT{Left: 0, Top: 0, Right: w, Bottom: h}
	primitives.AdjustWindowRect(&frame, primitives.WS_OVERLAPPEDWINDOW, false)
	rw = frame.Right - frame.Left
	rh = frame.Bottom - frame.Top
	return
}

func newMemDC(w, h int32) {
	if hDCMem != primitives.NULL {
		primitives.DeleteObject(hBitmap)
		primitives.DeleteDC(hDCMem)
	}

	hDC := primitives.GetDC(hWnd)
	hDCMem = primitives.CreateCompatibleDC(hDC)
	hBitmap = primitives.CreateCompatibleBitmap(hDC, w, h)
	primitives.SelectObject(hDCMem, hBitmap)
	primitives.SetBkMode(hDCMem, primitives.PS_SOLID)
}

func Start() error {
	classNamePtr, _ := syscall.UTF16PtrFromString(szClassName)
	hInstance, err := primitives.GetModuleHandle(nil)
	if err != nil {
		return err
	}

	var wincl primitives.WNDCLASSEX
	wincl.HInstance = hInstance
	wincl.LpszClassName = uintptr(unsafe.Pointer(classNamePtr))
	wincl.LpfnWndProc = syscall.NewCallback(windowProcedure)
	wincl.Style = primitives.CS_DBLCLKS
	wincl.CbSize = uint32(unsafe.Sizeof(wincl))
	wincl.LpszMenuName = 0
	wincl.CbClsExtra = 0
	wincl.CbWndExtra = 0
	background, err := primitives.GetStockObject(primitives.BLACK_BRUSH)
	if err != nil {
		return err
	}
	wincl.HbrBackground = background

	if _, err := primitives.RegisterClassEx(&wincl); err != nil {
		return err
	}

	var w, h int32
	w, h = realFrame(iWidth, iHeight)

	hWnd, err = primitives.CreateWindowEx(
		0, // Extended options
		uintptr(unsafe.Pointer(primitives.StringToUTF16Ptr(szClassName))),
		uintptr(unsafe.Pointer(primitives.StringToUTF16Ptr("Miniwin"))),
		primitives.WS_OVERLAPPEDWINDOW,
		primitives.CW_USEDEFAULT, // X pos
		primitives.CW_USEDEFAULT, // Y pos
		w,                        // widht
		h,                        // height
		primitives.HWND_DESKTOP,  // main window
		0,                        // no menu
		wincl.HInstance,          // app instance
		0,                        // no additional data
	)
	if err != nil {
		return err
	}

	hBitmap = 0 // starts bitmap
	primitives.ShowWindow(hWnd, 1)

	// messages loop
	var msg primitives.MSG
	for primitives.GetMessage(&msg, 0, 0, 0) {
		primitives.TranslateMessage(&msg)
		primitives.DispatchMessage(&msg)
	}

	return nil
}

func windowProcedure(hWnd primitives.HWND, message primitives.UINT, wParam primitives.WPARAM, lParam primitives.LPARAM) primitives.LRESULT {
	switch message {
	case primitives.WM_CREATE:
		primitives.SetTimer(hWnd, 1, 36, primitives.NULL)
	case primitives.WM_SIZE:
		var rect *primitives.RECT
		rect, _ = primitives.GetClientRect(hWnd)
		w := rect.Right - rect.Left
		h := rect.Bottom - rect.Top

		if w == 0 && h == 0 {
			break // on minimize send WM_SIZE(0,0)
		}

		if hDCMem == 0 || w != iWidth || h != iHeight {
			newMemDC(w, h)
		}
	case primitives.WM_SIZING:
		rect := (*primitives.RECT)(unsafe.Pointer(lParam))
		var w, h int32
		w, h = realFrame(iWidth, iHeight)

		switch wParam {
		case primitives.WMSZ_BOTTOM:
			rect.Bottom = rect.Top + h
		case primitives.WMSZ_TOP:
			rect.Top = rect.Bottom - h
		case primitives.WMSZ_RIGHT:
			rect.Right = rect.Left + w
		case primitives.WMSZ_LEFT:
			rect.Left = rect.Right - w
		case primitives.WMSZ_TOPLEFT:
			rect.Top = rect.Bottom - h
			rect.Left = rect.Right - w
		case primitives.WMSZ_TOPRIGHT:
			rect.Top = rect.Bottom - h
			rect.Right = rect.Left + w
		case primitives.WMSZ_BOTTOMLEFT:
			rect.Bottom = rect.Top + h
			rect.Left = rect.Right - w
		case primitives.WMSZ_BOTTOMRIGHT:
			rect.Bottom = rect.Top + h
			rect.Right = rect.Left + w
		}

		return 1
	case primitives.WM_PAINT:
		hdc, ps, _ := primitives.BeginPaint(hWnd)
		primitives.SelectObject(hDCMem, primitives.HGDIOBJ(hBitmap))
		if hBitmap != 0 {
			primitives.BitBlt(hdc, 0, 0, iWidth, iHeight, hDCMem, 0, 0, primitives.SRCCOPY)
		}
		primitives.EndPaint(hWnd, ps)
	case primitives.WM_MOUSEMOVE:
		mouseInside = true
		mouseXPos = primitives.GET_X_LPARAM(lParam)
		mouseYPos = primitives.GET_Y_LPARAM(lParam)
		mouseLeftClicked = (wParam & primitives.MK_LBUTTON) != 0
		mouseRightClicked = (wParam & primitives.MK_RBUTTON) != 0
	case primitives.WM_MOUSELEAVE:
		mouseInside = false
	case primitives.WM_LBUTTONDOWN:
		mouseLeftClicked = true
	case primitives.WM_LBUTTONUP:
		mouseLeftClicked = false
	case primitives.WM_RBUTTONDOWN:
		mouseRightClicked = true
	case primitives.WM_RBUTTONUP:
		mouseRightClicked = false
	case primitives.WM_KEYDOWN:
		pushIt := false

		// Escape
		pushIt = pushIt || (wParam == primitives.VK_ESCAPE)

		// Flechas
		pushIt = pushIt || (wParam == primitives.VK_LEFT || wParam == primitives.VK_RIGHT || wParam == primitives.VK_UP || wParam == primitives.VK_DOWN)

		// Barra espaciadora
		pushIt = pushIt || (wParam == primitives.VK_SPACE)

		// Enter
		pushIt = pushIt || (wParam == primitives.VK_RETURN)

		// Números 0-9
		pushIt = pushIt || (wParam >= 48 && wParam <= 57)

		// Letras A-Z
		pushIt = pushIt || (wParam >= 65 && wParam <= 90)

		// Teclas de función
		for i := 0; i < 10; i++ {
			pushIt = pushIt || (wParam == (primitives.VK_F1 + uintptr(i)))
		}

		if pushIt {
			keysPressed = append(keysPressed, int(wParam))
		}

	case primitives.WM_DESTROY:
		primitives.DeleteObject(primitives.HGDIOBJ(hBitmap))
		primitives.DeleteDC(hDCMem)
		primitives.PostQuitMessage(0)
	default:
		return primitives.DefWindowProc(hWnd, message, wParam, lParam)
	}

	return 0
}

func KeyPressed() int {
	if len(keysPressed) == 0 {
		return NONE
	}

	var ret int
	key := keysPressed[0]

	switch key {
	case primitives.VK_LEFT:
		ret = LEFT
	case primitives.VK_RIGHT:
		ret = RIGHT
	case primitives.VK_UP:
		ret = UP
	case primitives.VK_DOWN:
		ret = DOWN
	case primitives.VK_ESCAPE:
		ret = ESCAPE
	case primitives.VK_SPACE:
		ret = SPACE
	case primitives.VK_RETURN:
		ret = RETURN
	case primitives.VK_F1:
		ret = F1
	case primitives.VK_F2:
		ret = F2
	case primitives.VK_F3:
		ret = F3
	case primitives.VK_F4:
		ret = F4
	case primitives.VK_F5:
		ret = F5
	case primitives.VK_F6:
		ret = F6
	case primitives.VK_F7:
		ret = F7
	case primitives.VK_F8:
		ret = F8
	case primitives.VK_F9:
		ret = F9
	case primitives.VK_F10:
		ret = F10
	default:
		ret = key
	}

	keysPressed = keysPressed[1:]
	return ret
}

func MouseState() (bool, float32, float32) {
	if !mouseInside {
		return false, 0, 0
	}

	return true, float32(mouseXPos), float32(mouseYPos)
}

func IsMouseInside() bool {
	return mouseInside
}

func MouseX() float32 {
	return float32(mouseXPos)
}

func MouseY() float32 {
	return float32(mouseYPos)
}

func MouseButtons() (bool, bool) {
	return mouseLeftClicked, mouseRightClicked
}

func MouseLeftClicked() bool {
	return mouseLeftClicked
}
func MouseRightClicked() bool {
	return mouseRightClicked
}

func Clear() {
	var rect primitives.RECT
	primitives.SetRect(&rect, 0, 0, iWidth, iHeight)
	hBrush := primitives.CreateSolidBrush(primitives.RGB(0, 0, 0))
	defer primitives.DeleteObject(hBrush)
	primitives.FillRect(hDCMem, &rect, hBrush)
}

func Refresh() {
	primitives.InvalidateRect(hWnd, nil, false)
}

func Point(x, y float32) {
	primitives.SetPixel(hDCMem, int32(x), int32(y), selectedColor)
}

func Line(x1, y1, x2, y2 float32) {
	primitives.BeginPath(hDCMem)
	primitives.MoveToEx(hDCMem, int32(x1), int32(y1), nil)
	primitives.LineTo(hDCMem, int32(x2), int32(y2))
	primitives.EndPath(hDCMem)

	hPen, _ := primitives.CreatePen(primitives.PS_SOLID, 1, selectedColor)
	defer primitives.DeleteObject(hPen)

	primitives.SelectObject(hDCMem, hPen)
	primitives.StrokePath(hDCMem)
}

func makeRectangle(left, top, right, bottom float32) {
	primitives.BeginPath(hDCMem)
	primitives.MoveToEx(hDCMem, int32(left), int32(top), nil)
	primitives.LineTo(hDCMem, int32(left), int32(bottom))
	primitives.LineTo(hDCMem, int32(right), int32(bottom))
	primitives.LineTo(hDCMem, int32(right), int32(top))
	primitives.LineTo(hDCMem, int32(left), int32(top))
	primitives.EndPath(hDCMem)
}

func Rectangle(left, top, right, bottom float32) {
	hPen, _ := primitives.CreatePen(primitives.PS_SOLID, 1, selectedColor)
	defer primitives.DeleteObject(hPen)

	orig := primitives.SelectObject(hDCMem, hPen)
	makeRectangle(left, top, right, bottom)
	primitives.StrokePath(hDCMem)
	primitives.SelectObject(hDCMem, orig)
}

func FilledRectangle(left, top, right, bottom float32) {
	hBrush := primitives.CreateSolidBrush(selectedColor)
	defer primitives.DeleteObject(hBrush)

	orig := primitives.SelectObject(hDCMem, hBrush)
	makeRectangle(left, top, right, bottom)
	primitives.FillPath(hDCMem)
	primitives.SelectObject(hDCMem, orig)
}

func makeCircle(centerX, centerY, radius float32) {
	primitives.BeginPath(hDCMem)
	primitives.Arc(hDCMem,
		int32(centerX-radius), int32(centerY-radius),
		int32(centerX+radius), int32(centerY+radius),
		int32(centerX-radius), int32(centerY-radius),
		int32(centerX-radius), int32(centerY-radius),
	)
	primitives.EndPath(hDCMem)
}

func Circle(centerX, centerY, radius float32) {
	hPen, _ := primitives.CreatePen(primitives.PS_SOLID, 1, selectedColor)
	defer primitives.DeleteObject(hPen)

	orig := primitives.SelectObject(hDCMem, hPen)
	makeCircle(centerX, centerY, radius)
	primitives.StrokePath(hDCMem)
	primitives.SelectObject(hDCMem, orig)
}

func FilledCircle(centerX, centerY, radius float32) {
	hBrush := primitives.CreateSolidBrush(selectedColor)
	defer primitives.DeleteObject(hBrush)

	orig := primitives.SelectObject(hDCMem, hBrush)
	makeCircle(centerX, centerY, radius)
	primitives.FillPath(hDCMem)
	primitives.SelectObject(hDCMem, orig)
}

func SetText(x, y float32, content string) {
	primitives.SetTextColor(hDCMem, selectedColor)
	primitives.TextOut(hDCMem, int32(x), int32(y), content)
}

func SetColor(c int) {
	if c >= 0 && c < len(colorRefs) {
		selectedColor = colorRefs[c]
	}
}

func SetColorRGB(r, g, b int) {
	selectedColor = primitives.RGB(byte(r), byte(g), byte(b))
}

func Width() int {
	return int(iWidth)
}

func Height() int {
	return int(iHeight)
}

func Resize(newWidth, newHeight int) {
	iWidth = int32(newWidth)
	iHeight = int32(newHeight)

	var w, h int32
	w, h = realFrame(iWidth, iHeight)

	primitives.SetWindowPos(hWnd, 0, 0, 0, w, h, primitives.SWP_NOMOVE)
	newMemDC(w, h)
}

func Close() {
	primitives.PostMessage(hWnd, primitives.WM_CLOSE, 0, 0)
}
