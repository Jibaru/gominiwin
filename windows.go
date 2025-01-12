//go:build windows

package gominiwin

import "github.com/jibaru/gominiwin/windows"

func NewWindowsWin(title string, width, height int) (Win, error) {
	return windows.New(title, width, height)
}
