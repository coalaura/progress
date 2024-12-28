//go:build !windows

package progress

import (
	"os"
	"strings"
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

func supportsUnicode() bool {
	// Check LANG environment variable
	lang := os.Getenv("LANG")
	if strings.Contains(strings.ToLower(lang), "utf-8") || strings.Contains(strings.ToLower(lang), "utf8") {
		return true
	}

	// Check LC_ALL and LC_CTYPE
	for _, env := range []string{"LC_ALL", "LC_CTYPE"} {
		val := os.Getenv(env)
		if strings.Contains(strings.ToLower(val), "utf-8") || strings.Contains(strings.ToLower(val), "utf8") {
			return true
		}
	}

	// Check TERM environment variable
	term := os.Getenv("TERM")
	if term == "xterm" || term == "xterm-256color" || strings.HasPrefix(term, "screen") {
		return true
	}

	return false
}
