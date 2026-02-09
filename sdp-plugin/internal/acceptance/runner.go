package acceptance

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

// Runner runs the acceptance test command (AC2, AC3).
type Runner struct {
	Command  string
	Timeout  time.Duration
	Expected string // substring match (AC5)
}

// Result is the outcome of Run (AC4, AC8).
type Result struct {
	Passed   bool
	Duration time.Duration
	Output   string
	Error    string
}

// Run executes the command with timeout (AC2, AC3, AC4).
func (r *Runner) Run(ctx context.Context) (*Result, error) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", r.Command)
	out, err := cmd.CombinedOutput()
	duration := time.Since(start)
	result := &Result{
		Duration: duration,
		Output:   string(out),
	}
	if err != nil {
		result.Passed = false
		result.Error = err.Error()
		return result, nil
	}
	result.Passed = cmd.ProcessState.ExitCode() == 0
	if r.Expected != "" && !strings.Contains(result.Output, r.Expected) {
		result.Passed = false
		result.Error = "expected output not found: " + r.Expected
	}
	return result, nil
}

// ParseTimeout converts a string like "30s" to time.Duration.
func ParseTimeout(s string) (time.Duration, error) {
	if s == "" {
		return 30 * time.Second, nil
	}
	return time.ParseDuration(s)
}
