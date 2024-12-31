//go:build !windows

package progress

import (
	"os"
	"strings"
)

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
