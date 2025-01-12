//go:build linux

package gominiwin

import "github.com/jibaru/gominiwin/linux"

func NewLinuxWin(title string, width, height int) (Win, error) {
	return linux.New(title, width, height)
}
