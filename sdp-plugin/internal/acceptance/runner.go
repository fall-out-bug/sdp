package acceptance

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
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

// Run executes the command with timeout (AC2, AC3, AC4).
func (r *Runner) Run(ctx context.Context) (*Result, error) {
	start := time.Now()

	cmd := exec.Command("sh", "-c", r.Command)
	if r.Dir != "" {
		cmd.Dir = r.Dir
	}
	// Set process group to ensure we can kill all child processes on timeout
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Use channels to coordinate between goroutines
	done := make(chan error, 1)
	var out []byte

	// Start command in goroutine
	go func() {
		var err error
		out, err = cmd.CombinedOutput()
		done <- err
	}()

	// Wait for either command to complete or timeout
	select {
	case err := <-done:
		// Command completed (either success or failure)
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

	case <-time.After(r.Timeout):
		// Timeout - kill the entire process group
		if cmd.Process != nil {
			// Kill negative PID kills the entire process group
			// nolint:errcheck // Error ignored - we're in timeout scenario
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		duration := time.Since(start)
		return &Result{
			Passed:   false,
			Duration: duration,
			Output:   string(out),
			Error:    fmt.Sprintf("command timed out after %v", r.Timeout),
		}, nil
	}
}

// ParseTimeout converts a string like "30s" to time.Duration.
func ParseTimeout(s string) (time.Duration, error) {
	if s == "" {
		return 30 * time.Second, nil
	}
	return time.ParseDuration(s)
}
