package progress

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type ProgressBar struct {
	Total   int
	Current int
	Label   string
	Theme   ProgressTheme

	stop  chan struct{}
	abort chan struct{}
	wg    sync.WaitGroup
}

// NewProgressBar returns a new progress bar with the given label, total and the default theme
func NewProgressBar(label string, total int) *ProgressBar {
	return NewProgressBarWithTheme(label, total, ThemeDefault)
}

// NewProgressBarWithTheme returns a new progress bar with the given label, total and theme
func NewProgressBarWithTheme(label string, total int, theme ProgressTheme) *ProgressBar {
	return &ProgressBar{
		Total:   total,
		Current: 0,
		Label:   label,
		Theme:   theme,

		stop:  make(chan struct{}),
		abort: make(chan struct{}),
		wg:    sync.WaitGroup{},
	}
}

// Increment increments the progress bar by 1
func (p *ProgressBar) Increment() {
	p.IncrementBy(1)
}

// IncrementBy increments the progress bar by the given amount
func (p *ProgressBar) IncrementBy(amount int) {
	p.Current += amount

	if p.Current > p.Total {
		p.Current = p.Total
	}
}

// Reset resets the progress bar's current to 0
func (p *ProgressBar) Reset() {
	p.Current = 0
}

// Finished returns true if the progress bar has finished (reached its total)
func (p *ProgressBar) Finished() bool {
	return p.Current == p.Total
}

// Start starts the progress bar draw-goroutine
func (p *ProgressBar) Start() {
	p.wg.Add(1)

	go func() {
		var aborted bool

		defer func() {
			if !aborted {
				p.Current = p.Total
				p.draw()
			}

			fmt.Println()

			p.wg.Done()
		}()

		ticker := time.NewTicker(150 * time.Millisecond)

		lastCurrent := -1

		for {
			select {
			case <-ticker.C:
				if p.Current != lastCurrent {
					p.draw()

					lastCurrent = p.Current
				}
			case <-p.stop:
				return
			case <-p.abort:
				aborted = true

				return
			}
		}
	}()
}

// Stop stops the progress bar, prints the final progress and waits for the draw goroutine to finish
func (p *ProgressBar) Stop() {
	close(p.stop)

	p.wg.Wait()
}

// Abort stops the progress bar, does not print the final progress and waits for the draw goroutine to finish
func (p *ProgressBar) Abort() {
	close(p.abort)

	p.wg.Wait()
}

func (p *ProgressBar) draw() error {
	cols, err := getTermCols()
	if err != nil {
		return err
	}

	percentage := float64(p.Current) / float64(p.Total) * 100

	left := fmt.Sprintf("%s [", p.Label)
	right := fmt.Sprintf("] %5.1f%%", percentage)

	width := cols - len([]rune(left)) - len([]rune(right)) - 1

	filled := int(percentage / 100 * float64(width))
	empty := width - filled

	bar := strings.Repeat(p.Theme.Filled, filled)

	if filled < width {
		bar += p.Theme.Tip

		empty -= len([]rune(p.Theme.Tip))
	}

	bar += strings.Repeat(p.Theme.Empty, empty)

	fmt.Printf("%s%s%s\r", left, bar, right)

	return nil
}
