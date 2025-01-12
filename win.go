// Package gominiwin provides a minimal set of functions to create a window,
// draw on it, and handle basic input events. It is designed for creating simple
// games and graphical applications without requiring knowledge of underlying
// system APIs.
package gominiwin

import "github.com/jibaru/gominiwin/colors"

// Canva defines drawing operations for rendering shapes, text, and other graphical elements.
type Canva interface {
	// Point draws a point at the specified coordinates using the current color.
	//
	// Parameters:
	//   x: The X-coordinate of the point.
	//   y: The Y-coordinate of the point.
	Point(x, y float32)

	// Line draws a line from (x1, y1) to (x2, y2) using the current color.
	//
	// Parameters:
	//   x1: X-coordinate of the starting point.
	//   y1: Y-coordinate of the starting point.
	//   x2: X-coordinate of the ending point.
	//   y2: Y-coordinate of the ending point.
	Line(x1, y1, x2, y2 float32)

	// Rectangle draws a rectangle with the specified boundaries using the current color.
	// The rectangle is defined by the top-left and bottom-right corners.
	//
	// Parameters:
	//   left:   X-coordinate of the left edge.
	//   top:    Y-coordinate of the top edge.
	//   right:  X-coordinate of the right edge.
	//   bottom: Y-coordinate of the bottom edge.
	Rectangle(left, top, right, bottom float32)

	// FilledRectangle draws a filled rectangle with the specified boundaries using the current color.
	//
	// Parameters:
	//   left:   X-coordinate of the left edge.
	//   top:    Y-coordinate of the top edge.
	//   right:  X-coordinate of the right edge.
	//   bottom: Y-coordinate of the bottom edge.
	FilledRectangle(left, top, right, bottom float32)

	// Circle draws a circle with the specified center and radius using the current color.
	//
	// Parameters:
	//   centerX: X-coordinate of the circle's center.
	//   centerY: Y-coordinate of the circle's center.
	//   radius:  Radius of the circle.
	Circle(centerX, centerY, radius float32)

	// FilledCircle draws a filled circle with the specified center and radius using the current color.
	//
	// Parameters:
	//   centerX: X-coordinate of the circle's center.
	//   centerY: Y-coordinate of the circle's center.
	//   radius:  Radius of the circle.
	FilledCircle(centerX, centerY, radius float32)

	// SetColor sets the current drawing color.
	//
	// Parameters:
	//   c: A color in the `colors.Color` format.
	SetColor(c colors.Color)

	// SetColorRGB sets the current drawing color using RGB components.
	//
	// Parameters:
	//   r: Red component (0-255).
	//   g: Green component (0-255).
	//   b: Blue component (0-255).
	SetColorRGB(r, g, b int)

	// SetText draws text at the specified coordinates.
	//
	// Parameters:
	//   x:       X-coordinate of the text's starting position.
	//   y:       Y-coordinate of the text's starting position.
	//   content: The string to be drawn.
	SetText(x, y float32, content string)
}

// Input defines methods to handle keyboard and mouse input.
type Input interface {
	// KeyPressed returns the code of the last pressed key, or 0 if no key is pressed.
	KeyPressed() int

	// MouseState returns the current state of the mouse.
	//
	// Returns:
	//   inside:   Whether the mouse cursor is inside the window.
	//   mouseX:   The X-coordinate of the mouse cursor.
	//   mouseY:   The Y-coordinate of the mouse cursor.
	MouseState() (inside bool, mouseX, mouseY float32)

	// IsMouseInside checks if the mouse cursor is inside the window.
	IsMouseInside() bool

	// MouseX returns the current X-coordinate of the mouse cursor.
	MouseX() float32

	// MouseY returns the current Y-coordinate of the mouse cursor.
	MouseY() float32

	// MouseButtons returns the state of the mouse buttons.
	//
	// Returns:
	//   left:  Whether the left mouse button is pressed.
	//   right: Whether the right mouse button is pressed.
	MouseButtons() (left, right bool)

	// MouseLeftClicked checks if the left mouse button was clicked.
	MouseLeftClicked() bool

	// MouseRightClicked checks if the right mouse button was clicked.
	MouseRightClicked() bool
}

// Win combines the Canva and Input interfaces and provides window lifecycle methods.
type Win interface {
	Input
	Canva

	// Start begins the main event loop of the window. This method blocks the calling goroutine.
	Start()

	// Clear clears the contents of the window.
	Clear()

	// Refresh updates the window to display the latest drawing operations.
	Refresh()

	// Width returns the current width of the window in pixels.
	Width() int

	// Height returns the current height of the window in pixels.
	Height() int

	// Resize changes the dimensions of the window.
	//
	// Parameters:
	//   newWidth:  New width of the window in pixels.
	//   newHeight: New height of the window in pixels.
	Resize(newWidth, newHeight int)

	// Close closes the window and releases resources.
	Close()
}
