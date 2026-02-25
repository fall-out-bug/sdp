package verify

import (
	"context"
	"os/exec"
)

// CoverageResult is the minimal result needed for verification.
// Production implementations (e.g. quality.Checker) can adapt to this.
type CoverageResult struct {
	Coverage  float64
	Threshold float64
	Report    string
}

// CoverageChecker runs coverage analysis. Injectable for tests.
type CoverageChecker interface {
	CheckCoverage(ctx context.Context) (*CoverageResult, error)
}

// PathValidator validates that a path is within a base directory.
// Injectable for tests.
type PathValidator interface {
	ValidatePathInDirectory(baseDir, targetPath string) error
}

// CommandRunner creates a validated command for execution.
// Injectable for tests.
type CommandRunner interface {
	SafeCommand(ctx context.Context, name string, args ...string) (*exec.Cmd, error)
}
