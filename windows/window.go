package windows

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/jibaru/gominiwin"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	gdi32            = syscall.NewLazyDLL("gdi32.dll")   // Asegúrate de cargar gdi32.dll
	createSolidBrush = gdi32.NewProc("CreateSolidBrush") // Busca CreateSolidBrush en gdi32.dll
	deleteObject     = gdi32.NewProc("DeleteObject")     // También está en gdi32.dll
	createWindowEx   = user32.NewProc("CreateWindowExW")
	defWindowProc    = user32.NewProc("DefWindowProcW")
	// destroyWindow    = user32.NewProc("DestroyWindow")
	dispatchMessage = user32.NewProc("DispatchMessageW")
	getMessage      = user32.NewProc("GetMessageW")
	loadCursor      = user32.NewProc("LoadCursorW")
	postQuitMessage = user32.NewProc("PostQuitMessage")
	registerClassEx = user32.NewProc("RegisterClassExW")
	showWindow      = user32.NewProc("ShowWindow")
	updateWindow    = user32.NewProc("UpdateWindow")
	translateMsg    = user32.NewProc("TranslateMessage")
	beginPaint      = user32.NewProc("BeginPaint")
	invalidateRect  = user32.NewProc("InvalidateRect")
	endPaint        = user32.NewProc("EndPaint")
	// getDC            = user32.NewProc("GetDC")
	// releaseDC        = user32.NewProc("ReleaseDC")
	rectangle  = gdi32.NewProc("Rectangle")
	setPixel   = gdi32.NewProc("SetPixel")     // Dibujar un píxel
	moveToEx   = gdi32.NewProc("MoveToEx")     // Establecer el punto inicial de una línea
	lineTo     = gdi32.NewProc("LineTo")       // Dibujar una línea
	textOut    = gdi32.NewProc("TextOutW")     // Escribir texto en el dispositivo de contexto (HDC)
	messageBox = user32.NewProc("MessageBoxW") // Función para mostrar cuadros de diálogo estándar
	ellipse    = gdi32.NewProc("Ellipse")      // Función para dibujar círculos/elipses
)

const (
	WS_OVERLAPPEDWINDOW = 0x00CF0000
	WS_VISIBLE          = 0x10000000
	CW_USEDEFAULT       = 0x80000000
	SW_SHOW             = 5
	CS_HREDRAW          = 0x0002
	CS_VREDRAW          = 0x0001
	IDC_ARROW           = 32512
)

type (
	WndProc func(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr

	WNDCLASSEX struct {
		Size       uint32
		Style      uint32
		WndProc    uintptr
		ClsExtra   int32
		WndExtra   int32
		Instance   uintptr
		Icon       uintptr
		Cursor     uintptr
		Background uintptr
		MenuName   *uint16
		ClassName  *uint16
		IconSm     uintptr
	}

	MSG struct {
		Hwnd    uintptr
		Message uint32
		WParam  uintptr
		LParam  uintptr
		Time    uint32
		Pt      POINT
	}

	POINT struct {
		X, Y int32
	}

	PAINTSTRUCT struct {
		Hdc        uintptr
		FErase     int32
		RcPaint    RECT
		FRestore   int32
		FIncUpdate int32
		Reserved   [32]byte
	}

	RECT struct {
		Left, Top, Right, Bottom int32
	}
)

var _ gominiwin.Window = &Win{}

type Win struct {
	title  string
	width  int
	height int
	hwnd   uintptr
}

func NewWin(title string, width int, height int) (*Win, error) {
	if runtime.GOOS != "windows" {
		return nil, errors.New("os not supported, it should be windows")
	}

	return &Win{
		title:  title,
		width:  width,
		height: height,
	}, nil
}

func (ww *Win) ShowMessage(title, message string) {
	titlePtr := toUTF16Ptr(title)
	messagePtr := toUTF16Ptr(message)

	// Ventana modal simplificada
	_, _, _ = messageBox.Call(0, uintptr(unsafe.Pointer(messagePtr)), uintptr(unsafe.Pointer(titlePtr)), 0)
}

func (ww *Win) Refresh() {
	if ww.hwnd == 0 {
		fmt.Println("Window handle is not initialized")
		return
	}

	// Invalidate the entire client area
	ret, _, lastErr := invalidateRect.Call(ww.hwnd, 0, 0)
	if ret == 0 {
		fmt.Println("Failed to refresh the window")
	}

	slog.Info("refresh last error", "error", lastErr)
}

func (ww *Win) Run(onPaint func(w gominiwin.Paint), onInput func(w gominiwin.Input)) {
	hInstance, _, _ := kernel32.NewProc("GetModuleHandleW").Call(0)

	windowProc := func(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
		slog.Info("message received", "msg", msg)
		switch msg {
		case 0x0002: // WM_DESTROY
			postQuitMessage.Call(0)
			return 0
		case 0x000F: // WM_PAINT
			var ps PAINTSTRUCT
			hdc, _, lastErr := beginPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
			if lastErr != nil {
				slog.Info("error begin paint", "err", lastErr)
			}

			if onPaint != nil {
				onPaint(&Paint{hdc: hdc, win: ww})
			}
			endPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
			return 0
		case 0x0100, 0x0200, 0x0201, 0x0202: // WM_KEYDOWN, WM_MOUSEMOVE, WM_LBUTTONDOWN, WM_LBUTTONUP
			if onInput != nil {
				w := &Input{hwnd: hwnd}
				onInput(w)
			}
			return 0
		default:
			ret, _, _ := defWindowProc.Call(hwnd, uintptr(msg), wParam, lParam)
			return ret
		}
	}

	hwnd := createMainWindow(ww, hInstance, windowProc)
	showWindow.Call(hwnd, SW_SHOW)
	updateWindow.Call(hwnd)

	ww.hwnd = hwnd

	// Bucle principal de mensajes
	var msg MSG
	for {
		ret, _, _ := getMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if ret == 0 {
			break
		}

		translateMsg.Call(uintptr(unsafe.Pointer(&msg)))
		dispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

func createMainWindow(w *Win, hInstance uintptr, windowProc WndProc) uintptr {
	className := toUTF16Ptr("MainWndClass")

	cursor, _, _ := loadCursor.Call(0, uintptr(IDC_ARROW))
	if cursor == 0 {
		panic("Failed to load cursor")
	}

	wndClass := WNDCLASSEX{
		Size:       uint32(unsafe.Sizeof(WNDCLASSEX{})),
		Style:      CS_HREDRAW | CS_VREDRAW,
		WndProc:    syscall.NewCallback(windowProc),
		Instance:   hInstance,
		Cursor:     cursor,
		Background: 0,
		ClassName:  className,
	}

	registerClassEx.Call(uintptr(unsafe.Pointer(&wndClass)))

	hwnd, _, _ := createWindowEx.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(toUTF16Ptr(w.title))),
		WS_OVERLAPPEDWINDOW|WS_VISIBLE,
		CW_USEDEFAULT, CW_USEDEFAULT, uintptr(w.width), uintptr(w.height),
		0, 0, hInstance, 0,
	)

	return hwnd
}

func toUTF16Ptr(s string) *uint16 {
	ptr, _ := syscall.UTF16PtrFromString(s)
	return ptr
}

func fillRect(hdc uintptr, rect RECT, color uint32) error {
	// Crear brocha con color sólido
	brush, _, _ := createSolidBrush.Call(uintptr(color))
	if brush == 0 {
		return errors.New("failed to create solid brush")
	}
	defer deleteObject.Call(brush) // Liberar recursos después de usar la brocha

	// Llenar el rectángulo
	user32.NewProc("FillRect").Call(hdc, uintptr(unsafe.Pointer(&rect)), brush)
	return nil
}
