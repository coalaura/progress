//go:build !windows

package progress

import (
	"syscall"
	"unsafe"
)

func getTermCols() (int, error) {
	var ws struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}

	// TIOCGWINSZ is the terminal IO control request for getting window size
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdout),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)),
	)

	if err != 0 {
		return 0, err
	}

	return int(ws.Col), nil
}
