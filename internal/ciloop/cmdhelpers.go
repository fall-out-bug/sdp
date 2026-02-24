package ciloop

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const execRunnerTimeout = 30 * time.Second
const gitOperationTimeout = 60 * time.Second

// ExecRunner implements CommandRunner with process context and timeout.
// When ctx is cancelled (e.g. SIGTERM), Run returns promptly.
type ExecRunner struct {
	Ctx context.Context
}

// Run runs the command with ExecRunnerTimeout; respects Ctx cancellation.
func (e *ExecRunner) Run(name string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(e.Ctx, execRunnerTimeout)
	defer cancel()
	return exec.CommandContext(ctx, name, args...).Output()
}

// SanitizeLabel returns a label-safe string (alphanumeric and hyphen only).
func SanitizeLabel(s string) string {
	var b strings.Builder
	for _, r := range s {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	out := b.String()
	if out == "" {
		return "F000"
	}
	return out
}

// GitCommitter implements Committer via git CLI.
type GitCommitter struct{}

// AllFilesCommitter commits all changes (for deterministic fixers: goimports, go mod tidy).
type AllFilesCommitter struct{}

// Commit stages tracked files and commits (used by deterministic auto-fixers).
func (g *AllFilesCommitter) Commit(ctx context.Context, msg string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	runCtx, cancel := context.WithTimeout(ctx, gitOperationTimeout)
	defer cancel()
	add := exec.CommandContext(runCtx, "git", "add", "-u")
	add.Stdout = os.Stdout
	add.Stderr = os.Stderr
	if err := add.Run(); err != nil {
		return err
	}
	cmd := exec.CommandContext(runCtx, "git", "commit", "-m", msg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Push pushes the current branch.
func (g *AllFilesCommitter) Push(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	runCtx, cancel := context.WithTimeout(ctx, gitOperationTimeout)
	defer cancel()
	cmd := exec.CommandContext(runCtx, "git", "push")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Commit adds .sdp/ci-fixes/ and commits with the given message.
func (g *GitCommitter) Commit(ctx context.Context, msg string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	runCtx, cancel := context.WithTimeout(ctx, gitOperationTimeout)
	defer cancel()
	add := exec.CommandContext(runCtx, "git", "add", ".sdp/ci-fixes/")
	add.Stdout = os.Stdout
	add.Stderr = os.Stderr
	if err := add.Run(); err != nil {
		return err
	}
	cmd := exec.CommandContext(runCtx, "git", "commit", "-m", msg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Push pushes the current branch.
func (g *GitCommitter) Push(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	runCtx, cancel := context.WithTimeout(ctx, gitOperationTimeout)
	defer cancel()
	cmd := exec.CommandContext(runCtx, "git", "push")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GhLogFetcher implements LogFetcher via gh CLI.
type GhLogFetcher struct {
	Runner CommandRunner
}

// FailedLogs returns the log output of the most recent failed run for the current branch.
func (g *GhLogFetcher) FailedLogs(prNumber int) (string, error) {
	// Use Runner for git branch (respects Runner's context/timeout)
	out, err := g.Runner.Run("git", "branch", "--show-current")
	if err != nil {
		return "", fmt.Errorf("current branch: %w", err)
	}
	branch := strings.TrimSpace(string(out))
	runID, err := g.Runner.Run("gh", "run", "list",
		"--branch", branch,
		"--json", "databaseId,conclusion",
		"--jq", `.[] | select(.conclusion == "failure") | .databaseId`,
	)
	if err != nil {
		return "", fmt.Errorf("list failed runs: %w", err)
	}
	id := strings.TrimSpace(string(runID))
	if id == "" {
		return "", fmt.Errorf("no failed run found for PR #%d", prNumber)
	}
	if nl := strings.Index(id, "\n"); nl > 0 {
		id = id[:nl]
	}
	logOut, err := g.Runner.Run("gh", "run", "view", id, "--log-failed")
	if err != nil {
		return "", fmt.Errorf("fetch run logs: %w", err)
	}
	return string(logOut), nil
}
