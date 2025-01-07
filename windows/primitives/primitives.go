package primitives

import (
	"syscall"
	"unsafe"
)

var (
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	user32   = syscall.NewLazyDLL("user32.dll")

	pArc                    = gdi32.NewProc("Arc")
	pBeginPath              = gdi32.NewProc("BeginPath")
	pBitBlt                 = gdi32.NewProc("BitBlt")
	pCreateCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	pCreateCompatibleDC     = gdi32.NewProc("CreateCompatibleDC")
	pCreatePen              = gdi32.NewProc("CreatePen")
	pCreateSolidBrush       = gdi32.NewProc("CreateSolidBrush")
	pDeleteDC               = gdi32.NewProc("DeleteDC")
	pDeleteObject           = gdi32.NewProc("DeleteObject")
	pEllipse                = gdi32.NewProc("Ellipse")
	pEndPath                = gdi32.NewProc("EndPath")
	pFillPath               = gdi32.NewProc("FillPath")
	pGetStockObject         = gdi32.NewProc("GetStockObject")
	pGetTextExtentPoint32   = gdi32.NewProc("GetTextExtentPoint32W")
	pLineTo                 = gdi32.NewProc("LineTo")
	pMoveToEx               = gdi32.NewProc("MoveToEx")
	pSelectObject           = gdi32.NewProc("SelectObject")
	pSetBkMode              = gdi32.NewProc("SetBkMode")
	pSetPixel               = gdi32.NewProc("SetPixel")
	pSetTextColor           = gdi32.NewProc("SetTextColor")
	pStrokePath             = gdi32.NewProc("StrokePath")
	pTextOut                = gdi32.NewProc("TextOutW")
	pSetDCPenColor          = gdi32.NewProc("SetDCPenColor")
	pRectangle              = gdi32.NewProc("Rectangle")

	pGetModuleHandle = kernel32.NewProc("GetModuleHandleW")

	pAdjustWindowRect = user32.NewProc("AdjustWindowRect")
	pBeginPaint       = user32.NewProc("BeginPaint")
	pCreateWindowExW  = user32.NewProc("CreateWindowExW")
	pDefWindowProcW   = user32.NewProc("DefWindowProcW")
	pDispatchMessageW = user32.NewProc("DispatchMessageW")
	pEndPaint         = user32.NewProc("EndPaint")
	pFillRect         = user32.NewProc("FillRect")
	pGetClientRect    = user32.NewProc("GetClientRect")
	pGetDC            = user32.NewProc("GetDC")
	pGetMessageW      = user32.NewProc("GetMessageW")
	pInvalidateRect   = user32.NewProc("InvalidateRect")
	pLoadCursorW      = user32.NewProc("LoadCursorW")
	pMessageBoxW      = user32.NewProc("MessageBoxW")
	pPostMessage      = user32.NewProc("PostMessageW")
	pPostQuitMessage  = user32.NewProc("PostQuitMessage")
	pRegisterClassExW = user32.NewProc("RegisterClassExW")
	pReleaseDC        = user32.NewProc("ReleaseDC")
	pSetRect          = user32.NewProc("SetRect")
	pSetTimer         = user32.NewProc("SetTimer")
	pSetWindowPos     = user32.NewProc("SetWindowPos")
	pShowWindow       = user32.NewProc("ShowWindow")
	pTranslateMessage = user32.NewProc("TranslateMessage")
	pUpdateWindow     = user32.NewProc("UpdateWindow")
	pPeekMessage      = user32.NewProc("PeekMessage")
)

const (
	WS_EX_APPWINDOW     = 0x00040000
	WS_EX_CLIENTEDGE    = 0x00000200
	WS_EX_TOPMOST       = 0x00000008
	WS_BORDER           = 0x00800000
	WS_CAPTION          = 0x00C00000
	WS_CHILD            = 0x40000000
	WS_CLIPSIBLINGS     = 0x04000000
	WS_DISABLED         = 0x08000000
	WS_DLGFRAME         = 0x00400000
	WS_GROUP            = 0x00020000
	WS_HSCROLL          = 0x00100000
	WS_MINIMIZE         = 0x20000000
	WS_MAXIMIZEBOX      = 0x00010000
	WS_SYSMENU          = 0x00080000
	WS_TABSTOP          = 0x00010000
	WS_THICKFRAME       = 0x00040000
	WS_VSCROLL          = 0x00200000
	WS_OVERLAPPEDWINDOW = 0x00CF0000

	WM_CREATE      = 0x0001
	WM_DESTROY     = 0x0002
	WM_PAINT       = 0x000F
	WM_SIZING      = 0x0214
	WS_VISIBLE     = 0x10000000
	WM_CLOSE       = 0x0010
	WM_QUIT        = 0x0012
	WM_COMMAND     = 0x0111
	WM_LBUTTONDOWN = 0x0201
	WM_RBUTTONDOWN = 0x0204
	WM_MOUSEMOVE   = 0x0200
	WM_KEYDOWN     = 0x0100
	WM_KEYUP       = 0x0101
	WM_CHAR        = 0x0102
	WM_SIZE        = 0x0005
	WM_ERASEBKGND  = 0x0014
	WM_NCCREATE    = 0x0081
	WM_MOUSELEAVE  = 0x02A3
	WM_LBUTTONUP   = 0x0202
	WM_RBUTTONUP   = 0x0205
	WM_TIMER       = 0x0113

	WMSZ_LEFT        = 1
	WMSZ_RIGHT       = 2
	WMSZ_TOP         = 3
	WMSZ_TOPLEFT     = 4
	WMSZ_TOPRIGHT    = 5
	WMSZ_BOTTOM      = 6
	WMSZ_BOTTOMLEFT  = 7
	WMSZ_BOTTOMRIGHT = 8

	VK_NUMPAD0   = 0x60
	VK_NUMPAD1   = 0x61
	VK_NUMPAD2   = 0x62
	VK_NUMPAD3   = 0x63
	VK_NUMPAD4   = 0x64
	VK_NUMPAD5   = 0x65
	VK_NUMPAD6   = 0x66
	VK_NUMPAD7   = 0x67
	VK_NUMPAD8   = 0x68
	VK_NUMPAD9   = 0x69
	VK_MULTIPLY  = 0x6A
	VK_ADD       = 0x6B
	VK_SEPARATOR = 0x6C
	VK_SUBTRACT  = 0x6D
	VK_DECIMAL   = 0x6E
	VK_DIVIDE    = 0x6F
	VK_F1        = 0x70
	VK_F2        = 0x71
	VK_F3        = 0x72
	VK_F4        = 0x73
	VK_F5        = 0x74
	VK_F6        = 0x75
	VK_F7        = 0x76
	VK_F8        = 0x77
	VK_F9        = 0x78
	VK_F10       = 0x79
	VK_F11       = 0x7A
	VK_F12       = 0x7B
	VK_F13       = 0x7C
	VK_F14       = 0x7D
	VK_F15       = 0x7E
	VK_F16       = 0x7F
	VK_F17       = 0x80
	VK_F18       = 0x81
	VK_F19       = 0x82
	VK_F20       = 0x83
	VK_F21       = 0x84
	VK_F22       = 0x85
	VK_F23       = 0x86
	VK_F24       = 0x87

	VK_ESCAPE     = 0x1B
	VK_CLEAR      = 0x0C
	VK_RETURN     = 0x0D
	VK_CONVERT    = 0x1C
	VK_NONCONVERT = 0x1D
	VK_ACCEPT     = 0x1E
	VK_MODECHANGE = 0x1F
	VK_SPACE      = 0x20
	VK_PRIOR      = 0x21
	VK_NEXT       = 0x22
	VK_END        = 0x23
	VK_HOME       = 0x24
	VK_LEFT       = 0x25
	VK_UP         = 0x26
	VK_RIGHT      = 0x27
	VK_DOWN       = 0x28
	VK_SELECT     = 0x29
	VK_PRINT      = 0x2A
	VK_EXECUTE    = 0x2B
	VK_SNAPSHOT   = 0x2C
	VK_INSERT     = 0x2D
	VK_DELETE     = 0x2E
	VK_HELP       = 0x2F

	IDOK             = 1
	IDCANCEL         = 2
	IDABORT          = 3
	IDRETRY          = 4
	IDIGNORE         = 5
	IDYES            = 6
	IDNO             = 7
	IDC_ARROW uint16 = 32512

	MK_LBUTTON = 0x0001
	MK_RBUTTON = 0x0002
	MK_SHIFT   = 0x0004
	MK_CONTROL = 0x0008

	MB_OK       = 0x00000000
	MB_OKCANCEL = 0x00000001

	PS_SOLID   = 0x00000000
	CS_DBLCLKS = 0x0008

	COLOR_WINDOW = 5
	SW_SHOW      = 5

	CW_USEDEFAULT = 0x8000
	HWND_DESKTOP  = 0
	SRCCOPY       = 0x00CC0020
	SWP_NOMOVE    = 0x0002

	BLACK_BRUSH = 4

	FALSE = 0
	TRUE  = 1
	NULL  = 0
)

type ATOM = uint16
type BOOL = int32
type COLORREF = uint32
type DWORD = uint32
type HBITMAP = uintptr
type HBRUSH = uintptr
type HCURSOR = uintptr
type HDC = uintptr
type HDCMEM = uintptr
type HGDIOBJ = uintptr
type HICON = uintptr
type HINSTANCE = uintptr
type HMENU = uintptr
type HPEN = uintptr
type HWND = uintptr
type LONG = int32
type LPARAM = uintptr
type LPARAM_PTR = uintptr
type LPCWSTR = *uint16
type LPVOID = uintptr
type LRESULT = uintptr
type UINT = uint32
type UINT_PTR = uintptr
type WPARAM = uintptr

type MSG struct {
	HWnd    HWND
	Message uint32
	WParam  WPARAM
	LParam  LPARAM
	Time    uint32
	Pt      POINT
}

type PAINTSTRUCT struct {
	Hdc         HDC
	FErase      int32
	Rect        RECT
	Restore     int32
	IncUpdate   int32
	RGBReserved [32]byte
}

type POINT struct {
	X, Y int32
}

type RECT struct {
	Left, Top, Right, Bottom int32
}

type SIZE struct {
	CX int32
	CY int32
}

type WNDCLASSEX struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     HINSTANCE
	HIcon         uintptr
	HCursor       uintptr
	HbrBackground uintptr
	LpszMenuName  uintptr
	LpszClassName uintptr
	HIconSm       uintptr
}

func BeginPaint(hWnd HWND) (HDC, *PAINTSTRUCT, error) {
	var ps PAINTSTRUCT = PAINTSTRUCT{}
	ret, _, err := pBeginPaint.Call(hWnd, uintptr(unsafe.Pointer(&ps)))
	if ret == 0 {
		return 0, nil, err
	}

	return ret, &ps, nil
}

func BitBlt(hdcDest HDC, xDest, yDest, width, height int32, hdcSrc HDC, xSrc, ySrc int32, rop DWORD) bool {
	ret, _, _ := pBitBlt.Call(
		uintptr(hdcDest),
		uintptr(xDest),
		uintptr(yDest),
		uintptr(width),
		uintptr(height),
		uintptr(hdcSrc),
		uintptr(xSrc),
		uintptr(ySrc),
		uintptr(rop),
	)
	return ret != 0
}

func CreateWindowEx(
	exStyle DWORD,
	className uintptr,
	windowName uintptr,
	style DWORD,
	x, y, width, height int32,
	parent HWND,
	menu HMENU,
	instance HINSTANCE,
	param LPVOID,
) (HWND, error) {
	ret, _, err := pCreateWindowExW.Call(
		uintptr(exStyle),
		uintptr(className),
		uintptr(windowName),
		uintptr(style),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(parent),
		uintptr(menu),
		uintptr(instance),
		param,
	)
	if ret == 0 {
		return 0, err
	}
	return HWND(ret), nil
}

func DefWindowProc(hWnd HWND, msg UINT, wParam WPARAM, lParam LPARAM) uintptr {
	ret, _, _ := pDefWindowProcW.Call(uintptr(hWnd), uintptr(msg), uintptr(wParam), uintptr(lParam))
	return ret
}

func DeleteObject(hObject HGDIOBJ) bool {
	ret, _, _ := pDeleteObject.Call(hObject)
	return ret != 0
}

func DispatchMessage(msg *MSG) uintptr {
	ret, _, _ := pDispatchMessageW.Call(uintptr(unsafe.Pointer(msg)))
	return ret
}

func EndPaint(hWnd HWND, ps *PAINTSTRUCT) bool {
	ret, _, _ := pEndPaint.Call(hWnd, uintptr(unsafe.Pointer(ps)))
	return ret != 0
}

func FillRect(hdc HDC, rect *RECT, hBrush HBRUSH) error {
	ret, _, err := pFillRect.Call(uintptr(hdc), uintptr(unsafe.Pointer(rect)), hBrush)
	if ret == 0 {
		return err
	}
	return nil
}

func GetClientRect(hWnd HWND) (*RECT, error) {
	var rect RECT
	ret, _, err := pGetClientRect.Call(uintptr(hWnd), uintptr(unsafe.Pointer(&rect)))
	if ret == 0 {
		return nil, err
	}
	return &rect, nil
}

func GetDC(hWnd HWND) HDC {
	ret, _, _ := pGetDC.Call(uintptr(hWnd))
	return HDC(ret)
}

func GetMessage(msg *MSG, hWnd HWND, msgFilterMin, msgFilterMax UINT) bool {
	ret, _, _ := pGetMessageW.Call(uintptr(unsafe.Pointer(msg)), uintptr(hWnd), uintptr(msgFilterMin), uintptr(msgFilterMax))
	return ret > 0
}

func GetStockObject(object int32) (uintptr, error) {
	ret, _, err := pGetStockObject.Call(uintptr(object))
	if ret == 0 {
		return 0, err
	}
	return ret, nil
}

func InvalidateRect(hWnd HWND, rect *RECT, erase bool) {
	pInvalidateRect.Call(uintptr(hWnd), uintptr(unsafe.Pointer(rect)), uintptr(unsafe.Pointer(&erase)))
}

func LineTo(hdc HDC, x, y int32) bool {
	ret, _, _ := pLineTo.Call(uintptr(hdc), uintptr(x), uintptr(y))
	return ret != 0
}

func LoadCursor(hInstance HINSTANCE, lpCursorName LPCWSTR) (HCURSOR, error) {
	ret, _, err := pLoadCursorW.Call(hInstance, uintptr(unsafe.Pointer(lpCursorName)))
	if ret == 0 {
		return 0, err
	}
	return ret, nil
}

func MoveToEx(hdc HDC, x, y int32, lpPoint *POINT) bool {
	ret, _, _ := pMoveToEx.Call(uintptr(hdc), uintptr(x), uintptr(y), uintptr(unsafe.Pointer(lpPoint)))
	return ret != 0
}

func PostQuitMessage(exitCode int32) {
	pPostQuitMessage.Call(uintptr(exitCode))
}

func RegisterClassEx(wndClass *WNDCLASSEX) (uint16, error) {
	ret, _, err := pRegisterClassExW.Call(uintptr(unsafe.Pointer(wndClass)))
	if ret == 0 {
		return 0, err
	}
	return uint16(ret), nil
}

func ReleaseDC(hWnd HWND, hdc HDC) bool {
	ret, _, _ := pReleaseDC.Call(uintptr(hWnd), uintptr(hdc))
	return ret != 0
}

func ShowWindow(hWnd HWND, cmdShow int32) bool {
	ret, _, _ := pShowWindow.Call(uintptr(hWnd), uintptr(cmdShow))
	return ret != 0
}

func TranslateMessage(msg *MSG) {
	pTranslateMessage.Call(uintptr(unsafe.Pointer(msg)))
}

func PeekMessage(msg *MSG) {}

func UpdateWindow(hWnd HWND) bool {
	ret, _, _ := pUpdateWindow.Call(uintptr(hWnd))
	return ret != 0
}

func DeleteDC(hdc HDC) bool {
	ret, _, _ := pDeleteDC.Call(uintptr(hdc))
	return ret != 0
}

func CreateCompatibleBitmap(hdc HDC, width, height int32) HBITMAP {
	ret, _, _ := pCreateCompatibleBitmap.Call(uintptr(hdc), uintptr(width), uintptr(height))
	return HBITMAP(ret)
}

func CreateCompatibleDC(hdc HDC) HDC {
	ret, _, _ := pCreateCompatibleDC.Call(uintptr(hdc))
	return HDC(ret)
}

func SelectObject(hdc HDC, hObject HGDIOBJ) uintptr {
	ret, _, _ := pSelectObject.Call(uintptr(hdc), hObject)
	return ret
}

func SetBkMode(hdc HDC, mode int32) int32 {
	ret, _, _ := pSetBkMode.Call(uintptr(hdc), uintptr(mode))
	return int32(ret)
}

func MessageBoxW(hWnd HWND, text, caption LPCWSTR, type_ UINT) int32 {
	ret, _, _ := pMessageBoxW.Call(
		uintptr(hWnd),
		uintptr(unsafe.Pointer(text)),
		uintptr(unsafe.Pointer(caption)),
		uintptr(type_),
	)
	return int32(ret)
}

func AdjustWindowRect(rect *RECT, style DWORD, bMenu bool) bool {
	ret, _, _ := pAdjustWindowRect.Call(
		uintptr(unsafe.Pointer(rect)),
		uintptr(style),
		uintptr(BoolToUintptr(bMenu)),
	)
	return ret != 0
}

func CreateSolidBrush(color COLORREF) HBRUSH {
	// Llama a la API de Windows para crear un pincel sólido
	hBrush, _, _ := pCreateSolidBrush.Call(uintptr(color))
	return HBRUSH(hBrush)
}

func SetRect(rect *RECT, left, top, right, bottom int32) bool {
	// Llama a la API de Windows para crear un pincel sólido
	res, _, _ := pSetRect.Call(uintptr(unsafe.Pointer(rect)), uintptr(left), uintptr(top), uintptr(right), uintptr(bottom))
	return res != 0
}

func SetPixel(hdc HDC, x, y int32, color COLORREF) (COLORREF, error) {
	ret, _, err := pSetPixel.Call(
		uintptr(hdc),
		uintptr(x),
		uintptr(y),
		uintptr(color),
	)
	if ret == 0 {
		return 0, err
	}
	return COLORREF(ret), nil
}

func BeginPath(hdc HDC) error {
	ret, _, err := pBeginPath.Call(uintptr(hdc))
	if ret == 0 {
		return err
	}
	return nil
}

func EndPath(hdc HDC) error {
	ret, _, err := pEndPath.Call(uintptr(hdc))
	if ret == 0 {
		return err
	}
	return nil
}

func CreatePen(style, width int32, color COLORREF) (HPEN, error) {
	ret, _, err := pCreatePen.Call(
		uintptr(style),
		uintptr(width),
		uintptr(color),
	)
	if ret == 0 {
		return 0, err
	}
	return HPEN(ret), nil
}

func StrokePath(hdc HDC) error {
	ret, _, err := pStrokePath.Call(uintptr(hdc))
	if ret == 0 {
		return err
	}
	return nil
}

func FillPath(hdc HDC) error {
	ret, _, err := pFillPath.Call(uintptr(hdc))
	if ret == 0 {
		return err
	}
	return nil
}

func Arc(hdc HDC, left, top, right, bottom int32, startX, startY, endX, endY int32) error {
	ret, _, err := pArc.Call(
		uintptr(hdc),
		uintptr(left),
		uintptr(top),
		uintptr(right),
		uintptr(bottom),
		uintptr(startX),
		uintptr(startY),
		uintptr(endX),
		uintptr(endY),
	)
	if ret == 0 {
		return err
	}
	return nil
}

func SetTextColor(hdc HDC, color COLORREF) (COLORREF, error) {
	ret, _, err := pSetTextColor.Call(uintptr(hdc), uintptr(color))
	if ret == 0xFFFFFFFF {
		return 0, err
	}
	return COLORREF(ret), nil
}

func TextOut(hdc HDC, x, y int32, text string) error {
	textPtr, err := syscall.UTF16PtrFromString(text)
	if err != nil {
		return err
	}
	ret, _, callErr := pTextOut.Call(
		uintptr(hdc),
		uintptr(x),
		uintptr(y),
		uintptr(unsafe.Pointer(textPtr)),
		uintptr(len(text)),
	)
	if ret == 0 {
		return callErr
	}
	return nil
}

func SetWindowPos(hWnd HWND, hWndInsertAfter HWND, x, y, cx, cy int32, flags uint32) error {
	ret, _, err := pSetWindowPos.Call(
		uintptr(hWnd),
		uintptr(hWndInsertAfter),
		uintptr(x),
		uintptr(y),
		uintptr(cx),
		uintptr(cy),
		uintptr(flags),
	)
	if ret == 0 {
		return err
	}
	return nil
}

func PostMessage(hWnd HWND, msg UINT, wParam WPARAM, lParam LPARAM) error {
	ret, _, err := pPostMessage.Call(
		uintptr(hWnd),
		uintptr(msg),
		uintptr(wParam),
		uintptr(lParam),
	)
	if ret == 0 {
		return err
	}
	return nil
}

func GetModuleHandle(moduleName *uint16) (HINSTANCE, error) {
	ret, _, err := pGetModuleHandle.Call(uintptr(unsafe.Pointer(moduleName)))
	if ret == 0 {
		return 0, err
	}
	return HINSTANCE(ret), nil
}

func Ellipse(hdc HDC, left, top, right, bottom int32) error {
	ret, _, err := pEllipse.Call(
		uintptr(hdc),
		uintptr(left),
		uintptr(top),
		uintptr(right),
		uintptr(bottom),
	)
	if ret == 0 {
		return err
	}
	return nil
}

func GetTextExtentPoint32(hdc HDC, text string) (SIZE, error) {
	textPtr, err := syscall.UTF16PtrFromString(text)
	if err != nil {
		return SIZE{}, err
	}
	var size SIZE
	ret, _, callErr := pGetTextExtentPoint32.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(textPtr)),
		uintptr(len(text)),
		uintptr(unsafe.Pointer(&size)),
	)
	if ret == 0 {
		return SIZE{}, callErr
	}
	return size, nil
}

func SetTimer(hWnd HWND, nIDEvent uintptr, uElapse uint32, lpTimerFunc uintptr) (uintptr, error) {
	ret, _, err := pSetTimer.Call(
		uintptr(hWnd),
		nIDEvent,
		uintptr(uElapse),
		lpTimerFunc,
	)
	if ret == 0 {
		return 0, err
	}
	return ret, nil
}

func SetDCPenColor(hdc HDC, color COLORREF) HGDIOBJ {
	ret, _, _ := pSetDCPenColor.Call(hdc, uintptr(unsafe.Pointer(&color)))
	return ret
}

func Rectangle(hdc HDC, left, top, right, bottom int32) error {
	ret, _, err := pRectangle.Call(
		hdc,
		uintptr(left),
		uintptr(top),
		uintptr(right),
		uintptr(bottom),
	)
	if ret == 0 {
		return err
	}
	return nil
}

func BoolToUintptr(b bool) uintptr {
	if b {
		return 1
	}
	return 0
}

func RGB(r, g, b byte) uint32 {
	return uint32(r) | (uint32(g) << 8) | (uint32(b) << 16)
}

func GET_X_LPARAM(lp uintptr) int32 {
	return int32(uint16(lp))
}

func GET_Y_LPARAM(lp uintptr) int32 {
	return int32(uint16(lp >> 16))
}

func StringToUTF16Ptr(v string) *uint16 {
	utf16, _ := syscall.UTF16PtrFromString(v)
	return utf16
}
