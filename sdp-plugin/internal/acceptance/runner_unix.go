//go:build unix

package acceptance

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/fall-out-bug/sdp/internal/security"
)

// Run executes the command with timeout (AC2, AC3, AC4).
func (r *Runner) Run(ctx context.Context) (*Result, error) {
	start := time.Now()

	// Validate command for security (prevents shell injection)
	if err := security.ValidateTestCommand(r.Command); err != nil {
		return &Result{
			Passed:   false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("command validation failed: %v", err),
		}, nil
	}

	// Parse command into parts
	parts := strings.Fields(r.Command)
	if len(parts) == 0 {
		return &Result{
			Passed:   false,
			Duration: time.Since(start),
			Error:    "empty command",
		}, nil
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	// Create command directly (no shell)
	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
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
