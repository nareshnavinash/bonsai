package ui

import (
	"fmt"
	"os"
	"time"
)

var frames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type Spinner struct {
	message string
	stop    chan struct{}
	done    chan struct{}
}

func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}
}

func (s *Spinner) Start() {
	go func() {
		defer close(s.done)
		i := 0
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-s.stop:
				fmt.Fprint(os.Stderr, "\r\033[K")
				return
			case <-ticker.C:
				fmt.Fprintf(os.Stderr, "\r%s %s", frames[i%len(frames)], s.message)
				i++
			}
		}
	}()
}

func (s *Spinner) Stop(finalMessage string) {
	close(s.stop)
	<-s.done
	if finalMessage != "" {
		fmt.Fprintln(os.Stderr, finalMessage)
	}
}
