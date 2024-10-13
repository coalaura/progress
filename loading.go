package progress

import (
	"fmt"
	"sync"
	"time"
)

const LoadingFrames = "|/-\\"

type LoadingSpinner struct {
	frame   int
	running bool
	paused  bool

	stop chan struct{}
	wg   sync.WaitGroup
}

// NewLoadingSpinner returns a new loading spinner
func NewLoadingSpinner() *LoadingSpinner {
	return &LoadingSpinner{
		frame:   0,
		running: false,
		stop:    make(chan struct{}),
		wg:      sync.WaitGroup{},
	}
}

// Start starts the loading spinner draw-goroutine
func (l *LoadingSpinner) Start() {
	if l.running {
		return
	}

	l.running = true

	l.wg.Add(1)

	go func() {
		defer func() {
			l.wg.Done()

			l.running = false
		}()

		ticker := time.NewTicker(150 * time.Millisecond)

		for {
			select {
			case <-ticker.C:
				if !l.paused {
					l.step()
				}
			case <-l.stop:
				return
			}
		}
	}()
}

// Stop stops the loading spinner and waits for the draw goroutine to finish
func (l *LoadingSpinner) Stop() {
	if !l.running {
		return
	}

	close(l.stop)

	l.wg.Wait()
}

// Pause pauses the loading spinner
func (l *LoadingSpinner) Pause() {
	l.paused = true
}

// Resume resumes the loading spinner
func (l *LoadingSpinner) Resume() {
	l.paused = false
}

func (l *LoadingSpinner) step() {
	l.frame++
	l.frame %= len(LoadingFrames)

	fmt.Printf("%s\r", string(LoadingFrames[l.frame]))
}
