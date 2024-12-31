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

	tickrate   time.Duration
	theme      Theme
	delimiters bool
	counter    bool

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

// NewProgressBar returns a new progress bar with the given label, total, tickrate, theme, delimiters and counter
func NewProgressBar(label string, total int64, tickrate time.Duration, theme generator, delimiters, counter bool) *Bar {
	bar := &Bar{
		total:   total,
		label:   label,
		current: 0,

		length: len([]rune(label)),
		digits: int(math.Log10(float64(total))) + 1,

		tickrate:   tickrate,
		theme:      theme(),
		delimiters: delimiters,
		counter:    counter,

		running: false,
		stopped: true,

		stop:  make(chan struct{}),
		abort: make(chan struct{}),
		wg:    sync.WaitGroup{},
	}

	return bar
}

// NewProgressBarWithTheme returns a new progress bar with the given label, total and theme
func NewProgressBarWithTheme(label string, total int64, theme Theme) *Bar {
	return NewProgressBar(label, total, Fps20, func() Theme { return theme }, false, false)
}

// NewDefaultProgressBar returns a new progress bar with the given label and total
func NewDefaultProgressBar(label string, total int64) *Bar {
	return NewProgressBar(label, total, Fps20, ThemeBlocks, false, false)
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

			ticker = time.NewTicker(p.tickrate)
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
	columns := TerminalWidth()

	current := atomic.LoadInt64(&p.current)

	var (
		buffer bytes.Buffer
		suffix bytes.Buffer
	)

	buffer.WriteRune('\r')

	// Add the label
	if p.label != "" {
		buffer.WriteString(p.label)
		buffer.WriteString(" ")

		columns -= p.length + 1
	}

	// Build the suffix
	suffix.WriteString(" ")

	// Add the count
	if p.counter {
		suffix.WriteString(fmt.Sprintf("%*d/%*d ", p.digits, current, p.digits, p.total))
	}

	// Add the percentage
	percentage := float64(current) / float64(p.total)

	suffix.WriteString(fmt.Sprintf("%5.1f%%", percentage*100))

	// Update remaining space
	columns -= suffix.Len()

	// Add the left delimiter
	addDelimiter(&buffer, &columns, p.delimiters)

	// Add the bar
	p.theme(&buffer, percentage, columns)

	// Add the right delimiter
	addDelimiter(&buffer, &columns, p.delimiters)

	// Add the suffix
	buffer.Write(suffix.Bytes())

	// Print the progress bar
	fmt.Fprint(os.Stdout, buffer.String())

	return nil
}

func addDelimiter(buffer *bytes.Buffer, columns *int, delimiters bool) {
	if !delimiters {
		return
	}

	*columns--

	if SupportsUnicode() {
		buffer.WriteRune('│')
	} else {
		buffer.WriteRune('|')
	}
}
