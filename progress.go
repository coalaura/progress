package progress

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Bar struct {
	total   int64
	label   string
	current int64

	length int
	digits int
	theme  Theme

	running bool
	stopped bool

	stop  chan struct{}
	abort chan struct{}
	wg    sync.WaitGroup
}

const (
	Fps10 = 100 * time.Millisecond
	Fps20 = 50 * time.Millisecond
	Fps30 = 33 * time.Millisecond
	Fps60 = 16 * time.Millisecond
)

// default options
var (
	_tickrate       = Fps20
	_showDelimiters = false
	_showCounter    = false
)

// SetDefaults sets the default options (used when no options are passed)
func SetDefaults(tickrate time.Duration, showDelimiters, showCounter bool) {
	_tickrate = tickrate
	_showDelimiters = showDelimiters
	_showCounter = showCounter
}

// NewProgressBar returns a new progress bar with the given label, total and theme (or the default theme)
func NewProgressBar(label string, total int64, theme generator) *Bar {
	if theme == nil {
		theme = ThemeBlocksAscii
	}

	bar := &Bar{
		total:   total,
		label:   label,
		current: 0,

		length: len([]rune(label)),
		digits: int(math.Log10(float64(total))) + 1,
		theme:  theme(),

		running: false,
		stopped: true,

		stop:  make(chan struct{}),
		abort: make(chan struct{}),
		wg:    sync.WaitGroup{},
	}

	return bar
}

// Increment increments the progress bar by 1
func (p *Bar) Increment() {
	atomic.AddInt64(&p.current, 1)
}

// IncrementBy increments the progress bar by the given amount
func (p *Bar) IncrementBy(amount int64) {
	atomic.AddInt64(&p.current, amount)
}

// Reset resets the progress bar's current to 0
func (p *Bar) Reset() {
	atomic.StoreInt64(&p.current, 0)
}

// Finished returns true if the progress bar has finished (reached its total)
func (p *Bar) Finished() bool {
	return p.current == p.total
}

// Start starts the progress bar draw-goroutine
func (p *Bar) Start() {
	if p.running {
		return
	}

	p.running = true
	p.stopped = false

	p.wg.Add(1)

	go func() {
		var (
			last    int64
			current int64
			aborted bool

			ticker = time.NewTicker(_tickrate)
		)

		defer func() {
			if !aborted {
				atomic.StoreInt64(&p.current, p.total)

				p.draw()
			}

			fmt.Println()

			p.wg.Done()

			p.running = false
		}()

		p.draw()

		for {
			select {
			case <-ticker.C:
				current = atomic.LoadInt64(&p.current)

				if current != last {
					p.draw()

					last = current
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
func (p *Bar) Stop() {
	if p.stopped || !p.running {
		return
	}

	p.stopped = true

	close(p.stop)

	p.wg.Wait()
}

// Abort stops the progress bar, does not print the final progress and waits for the draw goroutine to finish
func (p *Bar) Abort() {
	if p.stopped || !p.running {
		return
	}

	p.stopped = true

	close(p.abort)

	p.wg.Wait()
}

func (p *Bar) draw() error {
	cols, err := getTermCols()
	if err != nil {
		return err
	}

	current := atomic.LoadInt64(&p.current)

	// Padding to avoid overflow
	cols -= 1

	var (
		buffer bytes.Buffer
		suffix bytes.Buffer
	)

	buffer.WriteRune('\r')

	// Add the label
	if p.label != "" {
		buffer.WriteString(p.label)
		buffer.WriteString(" ")

		cols -= p.length + 1
	}

	// Build the suffix
	suffix.WriteString(" ")

	// Add the count
	if _showCounter {
		suffix.WriteString(fmt.Sprintf("%*d/%*d ", p.digits, current, p.digits, p.total))
	}

	// Add the percentage
	percentage := float64(current) / float64(p.total)

	suffix.WriteString(fmt.Sprintf("%5.1f%%", percentage*100))

	// Update remaining space
	cols -= suffix.Len()

	// Add the left delimiter
	addDelimiter(&buffer)

	// Add the bar
	p.theme(&buffer, percentage, cols)

	// Add the right delimiter
	addDelimiter(&buffer)

	// Add the suffix
	buffer.Write(suffix.Bytes())

	// Print the progress bar
	fmt.Fprint(os.Stdout, buffer.String())

	return nil
}

func addDelimiter(buffer *bytes.Buffer) {
	if !_showDelimiters {
		return
	}

	if SupportsUnicode() {
		buffer.WriteRune('│')
	} else {
		buffer.WriteRune('|')
	}
}
