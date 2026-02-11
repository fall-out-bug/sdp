//go:build windows

package acceptance

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

// Run executes the command with timeout (AC2, AC3, AC4).
// Windows version doesn't use process groups (not available).
func (r *Runner) Run(ctx context.Context) (*Result, error) {
	start := time.Now()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	// Create command with context
	// Note: Windows doesn't have /bin/sh, use cmd.exe
	cmd := exec.CommandContext(ctx, "cmd", "/c", r.Command)
	if r.Dir != "" {
		cmd.Dir = r.Dir
	}

	// Capture output
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
	if r.Expected != "" && result.Output != "" && !strings.Contains(result.Output, r.Expected) {
		result.Passed = false
		result.Error = "expected output not found: " + r.Expected
	}
	return result, nil
}
