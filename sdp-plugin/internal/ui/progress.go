package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// ProgressBar displays a progress bar for long-running operations
type ProgressBar struct {
	total    int64
	current  int64
	width    int
	output   io.Writer
	lastTime time.Time
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int64) *ProgressBar {
	return &ProgressBar{
		total:    total,
		width:    50,
		output:   os.Stderr,
		lastTime: time.Now(),
	}
}

// SetWidth sets the width of the progress bar
func (p *ProgressBar) SetWidth(width int) {
	p.width = width
}

// SetOutput sets the output writer
func (p *ProgressBar) SetOutput(w io.Writer) {
	p.output = w
}

// Add increments the progress by n
func (p *ProgressBar) Add(n int64) {
	p.current += n
	p.render()
}

// Set sets the current progress
func (p *ProgressBar) Set(n int64) {
	p.current = n
	p.render()
}

// render updates the progress bar display
func (p *ProgressBar) render() {
	// Throttle updates to avoid flickering (max 10 updates per second)
	now := time.Now()
	if now.Sub(p.lastTime) < 100*time.Millisecond && p.current < p.total {
		return
	}
	p.lastTime = now

	var percent float64
	var filled int
	var bar string

	if p.total == 0 {
		// Handle zero total - show empty bar with indeterminate status
		percent = 0.0
		filled = 0
		bar = strings.Repeat("░", p.width)
	} else {
		percent = float64(p.current) / float64(p.total)
		if percent > 1.0 {
			percent = 1.0
		}
		// Calculate filled width
		filled = int(percent * float64(p.width))
		bar = strings.Repeat("█", filled) + strings.Repeat("░", p.width-filled)
	}

	// Format percentage
	percentStr := fmt.Sprintf("%.1f%%", percent*100)

	// Format progress
	currentStr := formatBytes(p.current)
	totalStr := formatBytes(p.total)

	// Render line with carriage return
	line := fmt.Sprintf("\r%s [%s] %s/%s", InfoSymbol(), bar, currentStr, totalStr)

	// If terminal supports color, add percentage at end
	if !NoColor && isTerminal() {
		line += " " + BoldText(percentStr)
	} else {
		line += " " + percentStr
	}

	fmt.Fprint(p.output, line)

	// Complete with newline when done
	if p.current >= p.total {
		fmt.Fprintln(p.output)
	}
}

// Complete marks the progress bar as complete
func (p *ProgressBar) Complete() {
	p.current = p.total
	p.render()
}

// formatBytes formats a byte count for display
func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

// Spinner displays a spinning animation for indeterminate progress
type Spinner struct {
	frames   []string
	current  int
	message  string
	output   io.Writer
	lastTime time.Time
	active   bool
}

// NewSpinner creates a new spinner
func NewSpinner(message string) *Spinner {
	return &Spinner{
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		message: message,
		output:  os.Stderr,
	}
}

// SetFrames sets custom spinner frames
func (s *Spinner) SetFrames(frames []string) {
	s.frames = frames
}

// SetOutput sets the output writer
func (s *Spinner) SetOutput(w io.Writer) {
	s.output = w
}

// Start starts the spinner animation
func (s *Spinner) Start() {
	s.active = true
	s.render()
}

// Stop stops the spinner with a success message
func (s *Spinner) Stop() {
	s.active = false
	fmt.Fprint(s.output, "\r") // Clear line
}

// StopWithError stops the spinner with an error message
func (s *Spinner) StopWithError() {
	s.active = false
	fmt.Fprint(s.output, "\r") // Clear line
}

// render updates the spinner display
func (s *Spinner) render() {
	if !s.active {
		return
	}

	// Throttle updates
	now := time.Now()
	if now.Sub(s.lastTime) < 100*time.Millisecond {
		return
	}
	s.lastTime = now

	frame := s.frames[s.current%len(s.frames)]
	s.current++

	line := fmt.Sprintf("\r%s %s %s", frame, InfoSymbol(), s.message)
	fmt.Fprint(s.output, line)
}
