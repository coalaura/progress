package progress

import (
	"fmt"
	"testing"
	"time"
)

func TestProgressBar(t *testing.T) {
	fmt.Println("- - - Progress Bars - - -")

	themes := map[string]generator{
		"Blocks ": ThemeBlocks,
		"Braille": ThemeBraille,
		"Dots   ": ThemeDots,
		"Pixels ": ThemePixels,
		"Shades ": ThemeShades,
	}

	for name, theme := range themes {
		testBarTheme(name, theme)
	}

	fmt.Println()
}

func TestLoadingSpinner(t *testing.T) {
	fmt.Println("- - - Loading Spinner - - -")

	spin := NewLoadingSpinner()

	spin.Start()

	time.Sleep(time.Duration(3000) * time.Millisecond)

	spin.Stop()

	fmt.Println()
}

func testBarTheme(name string, theme generator) {
	bar := NewProgressBar(name, 500, theme, true, true)

	bar.Start()
	bar.Start()
	bar.Start()

	for !bar.Finished() {
		time.Sleep(10 * time.Millisecond)

		bar.Increment()
	}

	bar.Stop()
	bar.Abort()
	bar.Stop()
	bar.Abort()
}
