package progress

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"sync"
	"time"
)

type Bar struct {
	total   int
	label   string
	current int

	length int
	digits int

	tick      time.Duration
	counter   bool
	delimiter *rune
	theme     Theme

	running bool
	stopped bool

	stop  chan struct{}
	abort chan struct{}
	wg    sync.WaitGroup
}

type Config struct {
	// If the current and total should be displayed left of the percentage
	Counter bool

	// If the delimiters should be displayed
	Delimiters bool

	// The progress bar theme
	Theme Theme

	// The tick duration between each update (draw call)
	Tick time.Duration
}

// NewProgressBar returns a new progress bar with the given label, total and the default theme
func NewProgressBar(label string, total int) *Bar {
	return &Bar{
		total:   total,
		label:   label,
		current: 0,

		length: len([]rune(label)),
		digits: int(math.Log10(float64(total))) + 1,

		tick:      50 * time.Millisecond,
		counter:   false,
		delimiter: nil,
		theme:     nil,

		running: false,
		stopped: true,

		stop:  make(chan struct{}),
		abort: make(chan struct{}),
		wg:    sync.WaitGroup{},
	}
}

// WithConfig sets the progress bar to use the given config
func (p *Bar) WithConfig(config Config) *Bar {
	p.counter = config.Counter
	p.theme = config.Theme
	p.tick = config.Tick

	if config.Delimiters {
		var delimiter rune

		if supportsUnicode() {
			delimiter = '│'
		} else {
			delimiter = '|'
		}

		p.delimiter = &delimiter
	}

	return p
}

// Increment increments the progress bar by 1
func (p *Bar) Increment() {
	p.IncrementBy(1)
}

// IncrementBy increments the progress bar by the given amount
func (p *Bar) IncrementBy(amount int) {
	p.current += amount

	if p.current > p.total {
		p.current = p.total
	}
}

// Reset resets the progress bar's current to 0
func (p *Bar) Reset() {
	p.current = 0
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
			current int
			aborted bool

			ticker = time.NewTicker(p.tick)
		)

		defer func() {
			if !aborted {
				p.current = p.total
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
				if p.current != current {
					p.draw()

					current = p.current
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

	// Fallback to default theme if no theme is set
	if p.theme == nil {
		p.theme = ThemeBlocksAscii()
	}

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
	if p.counter {
		suffix.WriteString(fmt.Sprintf("%*d/%*d ", p.digits, p.current, p.digits, p.total))
	}

	// Add the percentage
	percentage := float64(p.current) / float64(p.total)

	suffix.WriteString(fmt.Sprintf("%5.1f%%", percentage*100))

	// Update remaining space
	cols -= suffix.Len()

	// Add the left delimiter
	if p.delimiter != nil {
		buffer.WriteRune(*p.delimiter)

		cols -= 2
	}

	// Add the bar
	p.theme(&buffer, percentage, cols)

	// Add the right delimiter
	if p.delimiter != nil {
		buffer.WriteRune(*p.delimiter)
	}

	// Add the suffix
	buffer.Write(suffix.Bytes())

	// Print the progress bar
	fmt.Fprint(os.Stdout, buffer.String())

	return nil
}
