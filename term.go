package progress

import (
	"os"

	"golang.org/x/term"
)

var (
	unicode *bool
)

// SupportsUnicode returns true if the current terminal supports unicode
func SupportsUnicode() bool {
	if unicode == nil {
		unicode = new(bool)

		*unicode = supportsUnicode()
	}

	return *unicode
}

func TerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80 // fallback
	}

	// Padding to avoid overflow
	return width - 1
}
