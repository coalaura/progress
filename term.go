package progress

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
