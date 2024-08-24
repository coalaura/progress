package progress

type ProgressTheme struct {
	Empty  string
	Filled string
	Tip    string
}

var (
	// Default themes
	ThemeDefault = NewProgressTheme(" ", "=", ">")
	ThemeDots    = NewProgressTheme(".", ":", "o")
	ThemeHash    = NewProgressTheme("-", "#", ">")

	// Unicode themes
	ThemeBlocksUnicode   = NewProgressTheme(" ", "█", "▌")
	ThemeGradientUnicode = NewProgressTheme("░", "▓", "█")
)

// NewProgressTheme creates a new theme with empty being the empty character, filled being the filled character and tip being the tip character
func NewProgressTheme(empty, filled, tip string) ProgressTheme {
	return ProgressTheme{
		Empty:  empty,
		Filled: filled,
		Tip:    tip,
	}
}
