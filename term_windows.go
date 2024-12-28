//go:build windows

package progress

import (
	"os"
	"syscall"

	"golang.org/x/sys/windows"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	getConsoleOutputCP = kernel32.NewProc("GetConsoleOutputCP")
)

func getTermCols() (int, error) {
	stdOut := windows.Handle(os.Stdout.Fd())

	var csbi windows.ConsoleScreenBufferInfo

	err := windows.GetConsoleScreenBufferInfo(stdOut, &csbi)
	if err != nil {
		return 0, err
	}

	width := csbi.Window.Right - csbi.Window.Left + 1

	return int(width), nil
}

func supportsUnicode() bool {
	stdout := windows.Handle(windows.Stdout)
	var mode uint32
	err := windows.GetConsoleMode(stdout, &mode)
	if err != nil {
		return false
	}

	// Check if ENABLE_VIRTUAL_TERMINAL_PROCESSING is supported
	// This is available in Windows 10 and later
	const ENABLE_VIRTUAL_TERMINAL_PROCESSING uint32 = 0x0004
	if mode&ENABLE_VIRTUAL_TERMINAL_PROCESSING != 0 {
		return true
	}

	// Check code page
	ret, _, _ := getConsoleOutputCP.Call()
	cp := uint32(ret)

	// UTF-8 code page is 65001
	return cp == 65001
}
