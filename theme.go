package progress

import (
	"bytes"
)

type Theme func(*bytes.Buffer, float64, int)

/* Unicode Themes */

func ThemeBlocks() Theme {
	return asciiOrUnicode([]rune{' ', '▏', '▎', '▍', '▌', '▋', '▊', '▉', '█'}, ThemeBlocksAscii)
}

func ThemeBraille() Theme {
	return asciiOrUnicode([]rune{' ', '⡀', '⡄', '⡆', '⡇', '⣇', '⣧', '⣷', '⣿'}, ThemeBrailleAscii)
}

func ThemeDots() Theme {
	return asciiOrUnicode([]rune{' ', '⠁', '⠃', '⠇', '⠏', '⠟', '⠿', '⡿', '⣿'}, ThemeDotsAscii)
}

func ThemePixels() Theme {
	return asciiOrUnicode([]rune{' ', '⣀', '⣤', '⣶', '⣿', '⣿', '⣿', '⣿', '⣿'}, ThemePixelsAscii)
}

func ThemeShades() Theme {
	return asciiOrUnicode([]rune{'░', '░', '▒', '▒', '▓', '▓', '█', '█', '█'}, ThemeShadesAscii)
}

/* ASCII Themes */

func ThemeBlocksAscii() Theme {
	return NewThemeFromBlocks([]rune{' ', '.', '-', '=', '#'})
}

func ThemeBrailleAscii() Theme {
	return NewThemeFromBlocks([]rune{' ', '.', ',', '*', 'o', 'O', '@'})
}

func ThemeDotsAscii() Theme {
	return NewThemeFromBlocks([]rune{' ', '.', ':', ';', '8'})
}

func ThemePixelsAscii() Theme {
	return NewThemeFromBlocks([]rune{' ', '.', ':', '|', '#'})
}

func ThemeShadesAscii() Theme {
	return NewThemeFromBlocks([]rune{'.', '-', '=', '+', '#'})
}

// NewProgressTheme creates a new theme with the given blocks
func NewThemeFromBlocks(blocks []rune) Theme {
	var (
		maxI = len(blocks) - 1
		maxF = float64(maxI)
	)

	return func(buf *bytes.Buffer, percent float64, width int) {
		// Calculate how many full columns we need
		fullWidth := int(percent * float64(width))

		// Calculate the fractional part for the tip block
		fraction := (percent * float64(width)) - float64(fullWidth)
		tipBlock := int(fraction * maxF)

		// Add full blocks
		for i := 0; i < fullWidth; i++ {
			buf.WriteRune(blocks[maxI]) // The full block
		}

		// Add the tip block if we're not at 100% and we have space
		if fullWidth < width {
			buf.WriteRune(blocks[tipBlock])

			// Fill the rest with empty blocks
			for i := fullWidth + 1; i < width; i++ {
				buf.WriteRune(blocks[0]) // The empty block
			}
		}
	}
}

type generator func() Theme

func asciiOrUnicode(unicode []rune, ascii generator) Theme {
	if supportsUnicode() {
		return NewThemeFromBlocks(unicode)
	}

	return ascii()
}
