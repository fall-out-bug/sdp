package acceptance

import (
	"time"
)

// Runner runs the acceptance test command (AC2, AC3).
type Runner struct {
	Command  string
	Timeout  time.Duration
	Expected string // substring match (AC5)
	Dir      string // optional: run command in this directory (project root)
}

// Result is the outcome of Run (AC4, AC8).
type Result struct {
	Passed   bool
	Duration time.Duration
	Output   string
	Error    string
}

// ParseTimeout converts a string like "30s" to time.Duration.
func ParseTimeout(s string) (time.Duration, error) {
	if s == "" {
		return 30 * time.Second, nil
	}
	return time.ParseDuration(s)
}
