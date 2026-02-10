//go:build unix

package acceptance

import (
	"context"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// Run executes the command with timeout (AC2, AC3, AC4).
func (r *Runner) Run(ctx context.Context) (*Result, error) {
	start := time.Now()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	// Create command with context
	cmd := exec.CommandContext(ctx, "sh", "-c", r.Command)
	if r.Dir != "" {
		cmd.Dir = r.Dir
	}
	// Set process group to ensure we can kill all child processes
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

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
