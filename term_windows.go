//go:build windows

package progress

import (
	"os"

	"golang.org/x/sys/windows"
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
