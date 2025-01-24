//go:build linux

package gominiwin

import "github.com/jibaru/gominiwin/linux"

func New(title string, width, height int) (Win, error) {
	return linux.New(title, width, height)
}
