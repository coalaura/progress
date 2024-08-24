package progress

import (
	"math/rand"
	"testing"
	"time"
)

func TestProgressBar(t *testing.T) {
	testBarTheme("Default", ThemeDefault)
	testBarTheme("Dots", ThemeDots)
	testBarTheme("Hash", ThemeHash)

	testBarTheme("BlocksUnicode", ThemeBlocksUnicode)
	testBarTheme("GradientUnicode", ThemeGradientUnicode)
}

func testBarTheme(label string, theme ProgressTheme) {
	bar := NewProgressBarWithTheme(label, 80, theme)

	bar.Start()
	bar.Start()
	bar.Start()

	for !bar.Finished() {
		time.Sleep(time.Duration(50+rand.Intn(200)) * time.Millisecond)

		bar.IncrementBy(1 + rand.Intn(3))
	}

	bar.Stop()
	bar.Abort()
	bar.Stop()
	bar.Abort()
}
