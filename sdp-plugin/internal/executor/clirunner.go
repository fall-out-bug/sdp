package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

// CLIRunner runs a command for each workstream. Implements WorkstreamRunner.
// Use with NewExecutor to delegate workstream execution to a CLI (e.g. sdp build).
type CLIRunner struct {
	Command string   // Executable (e.g. "sdp")
	Args    []string // Base args (e.g. ["build"]). wsID is appended.
}

// NewCLIRunner creates a CLIRunner. Args are base args; wsID is appended per Run.
// Example: NewCLIRunner("sdp", "build") runs "sdp build <wsID>".
func NewCLIRunner(command string, args ...string) *CLIRunner {
	return &CLIRunner{Command: command, Args: args}
}

// Run executes the command with wsID appended to args.
func (r *CLIRunner) Run(ctx context.Context, wsID string) error {
	args := append(r.Args, wsID)
	cmd := exec.CommandContext(ctx, r.Command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %s: %w", r.Command, wsID, err)
	}
	return nil
}
