package windows

import (
	"errors"
	"unsafe"

	gominiwin "github.com/jibaru/gominiwin"
)

var _ gominiwin.Input = &Input{}

type Input struct {
	hwnd uintptr
}

func (wi *Input) Key() int {
	// Llamamos a GetAsyncKeyState para cada tecla
	for key := 0; key < 256; key++ {
		state, _, _ := user32.NewProc("GetAsyncKeyState").Call(uintptr(key))
		if state&0x8000 != 0 { // Bit alto indica que la tecla está presionada
			return key
		}
	}
	return 0 // Ninguna tecla está presionada
}

func (wi *Input) MouseState() (gominiwin.MouseState, error) {
	var pt POINT

	// Obtener la posición global del cursor
	ret, _, _ := user32.NewProc("GetCursorPos").Call(uintptr(unsafe.Pointer(&pt)))
	if ret == 0 {
		return gominiwin.MouseState{}, errors.New("failed to get mouse position")
	}

	// Convertir coordenadas a relativas a la ventana
	user32.NewProc("ScreenToClient").Call(wi.hwnd, uintptr(unsafe.Pointer(&pt)))

	// Determinar el estado de los botones del ratón
	leftState, _, _ := user32.NewProc("GetAsyncKeyState").Call(0x01)  // VK_LBUTTON
	rightState, _, _ := user32.NewProc("GetAsyncKeyState").Call(0x02) // VK_RBUTTON

	return gominiwin.MouseState{
		X:            float64(pt.X),
		Y:            float64(pt.Y),
		LeftClicked:  leftState&0x8000 != 0,
		RightClicked: rightState&0x8000 != 0,
	}, nil
}

func (wi *Input) IsMouseInsideWindow() (bool, error) {
	var pt POINT
	var rect RECT

	// Obtener la posición global del cursor
	ret, _, _ := user32.NewProc("GetCursorPos").Call(uintptr(unsafe.Pointer(&pt)))
	if ret == 0 {
		return false, errors.New("failed to get mouse position")
	}

	// Obtener las dimensiones del área del cliente
	user32.NewProc("GetClientRect").Call(wi.hwnd, uintptr(unsafe.Pointer(&rect)))

	// Convertir coordenadas del cursor a relativas a la ventana
	user32.NewProc("ScreenToClient").Call(wi.hwnd, uintptr(unsafe.Pointer(&pt)))

	// Verificar si está dentro del área del cliente
	return pt.X >= rect.Left && pt.X <= rect.Right && pt.Y >= rect.Top && pt.Y <= rect.Bottom, nil
}
