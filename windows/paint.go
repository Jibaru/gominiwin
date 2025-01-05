package windows

import (
	"errors"
	"unsafe"

	"github.com/jibaru/gominiwin"
)

var _ gominiwin.Paint = &Paint{}

type Paint struct {
	hdc uintptr
	win *Win
}

func (w *Paint) SetBackground(color gominiwin.Color) error {
	return fillRect(w.hdc, RECT{0, 0, int32(w.win.width), int32(w.win.height)}, color)
}

func (w *Paint) DrawFilledRectangle(left, top, right, bottom int32, color gominiwin.Color) error {
	return fillRect(w.hdc, RECT{left, top, right, bottom}, color)
}

func (w *Paint) DrawUnfilledRectangle(left, top, right, bottom int32, color gominiwin.Color) error {
	// Crear un pincel (brush) o usar NULL_BRUSH para rectángulos vacíos
	var brush uintptr
	brush, _, _ = gdi32.NewProc("GetStockObject").Call(5) // NULL_BRUSH para no rellenar

	// Seleccionar el pincel actual
	oldBrush, _, _ := gdi32.NewProc("SelectObject").Call(w.hdc, brush)
	defer gdi32.NewProc("SelectObject").Call(w.hdc, oldBrush) // Restaurar el pincel original

	// Dibujar el rectángulo
	ret, _, _ := rectangle.Call(w.hdc, uintptr(left), uintptr(top), uintptr(right), uintptr(bottom))
	if ret == 0 {
		return errors.New("failed to draw rectangle")
	}
	return nil
}

// DrawPoint dibuja un punto en las coordenadas (x, y) con un color específico.
func (w *Paint) DrawPoint(x, y int32, color gominiwin.Color) error {
	ret, _, _ := setPixel.Call(w.hdc, uintptr(x), uintptr(y), uintptr(color))
	if ret == 0 {
		return errors.New("failed to draw point")
	}
	return nil
}

// DrawLine dibuja una línea desde el punto (x1, y1) al punto (x2, y2) con un color específico.
func (w *Paint) DrawLine(x1, y1, x2, y2 int32, color gominiwin.Color) error {
	// Establecer el punto inicial
	moveToEx.Call(w.hdc, uintptr(x1), uintptr(y1), 0)

	// Dibujar la línea al punto final
	lineTo.Call(w.hdc, uintptr(x2), uintptr(y2))

	// Cambiar el color del lápiz (opcional, ya que el color es determinado por el contexto)
	ret, _, _ := setPixel.Call(w.hdc, uintptr(x2), uintptr(y2), uintptr(color))
	if ret == 0 {
		return errors.New("failed to draw line")
	}
	return nil
}

func (w *Paint) DrawText(x, y int32, text string) error {
	textPtr := toUTF16Ptr(text)
	ret, _, _ := textOut.Call(w.hdc, uintptr(x), uintptr(y), uintptr(unsafe.Pointer(textPtr)), uintptr(len(text)))
	if ret == 0 {
		return errors.New("failed to draw text")
	}
	return nil
}

func (w *Paint) DrawFilledCircle(left, top, right, bottom int32, color gominiwin.Color) error {
	// Crear un pincel (brush) o usar NULL_BRUSH para círculos vacíos
	var brush uintptr
	brush, _, _ = createSolidBrush.Call(uintptr(color))
	if brush == 0 {
		return errors.New("failed to create solid brush")
	}
	defer deleteObject.Call(brush)

	// Seleccionar el pincel y dibujar
	oldBrush, _, _ := gdi32.NewProc("SelectObject").Call(w.hdc, brush)
	defer gdi32.NewProc("SelectObject").Call(w.hdc, oldBrush)

	// Dibujar el círculo/elipse
	ret, _, _ := ellipse.Call(w.hdc, uintptr(left), uintptr(top), uintptr(right), uintptr(bottom))
	if ret == 0 {
		return errors.New("failed to draw filled circle")
	}
	return nil
}

func (w *Paint) DrawUnfilledCircle(left, top, right, bottom int32, color gominiwin.Color) error {
	// Crear un pincel (brush) o usar NULL_BRUSH para círculos vacíos
	var brush uintptr
	brush, _, _ = gdi32.NewProc("GetStockObject").Call(5) // NULL_BRUSH

	// Seleccionar el pincel y dibujar
	oldBrush, _, _ := gdi32.NewProc("SelectObject").Call(w.hdc, brush)
	defer gdi32.NewProc("SelectObject").Call(w.hdc, oldBrush)

	// Dibujar el círculo/elipse
	ret, _, _ := ellipse.Call(w.hdc, uintptr(left), uintptr(top), uintptr(right), uintptr(bottom))
	if ret == 0 {
		return errors.New("failed to draw unfilled circle")
	}
	return nil
}
