package verify

import (
	"context"
	"os/exec"

	"github.com/fall-out-bug/sdp/internal/quality"
	"github.com/fall-out-bug/sdp/internal/security"
)

// qualityCoverageChecker adapts quality.Checker to CoverageChecker.
type qualityCoverageChecker struct {
	checker *quality.Checker
}

func (q *qualityCoverageChecker) CheckCoverage(ctx context.Context) (*CoverageResult, error) {
	result, err := q.checker.CheckCoverage(ctx)
	if err != nil {
		return nil, err
	}
	return &CoverageResult{
		Coverage:  result.Coverage,
		Threshold: result.Threshold,
		Report:    result.Report,
	}, nil
}

// securityPathValidator adapts security.ValidatePathInDirectory to PathValidator.
type securityPathValidator struct{}

func (securityPathValidator) ValidatePathInDirectory(baseDir, targetPath string) error {
	return security.ValidatePathInDirectory(baseDir, targetPath)
}

// securityCommandRunner adapts security.SafeCommand to CommandRunner.
type securityCommandRunner struct{}

func (securityCommandRunner) SafeCommand(ctx context.Context, name string, args ...string) (*exec.Cmd, error) {
	return security.SafeCommand(ctx, name, args...)
}

// defaultCoverageChecker returns a CoverageChecker for the given project root.
func defaultCoverageChecker(projectRoot string) (CoverageChecker, error) {
	checker, err := quality.NewChecker(projectRoot)
	if err != nil {
		return nil, err
	}
	return &qualityCoverageChecker{checker: checker}, nil
}

// defaultPathValidator returns the production PathValidator.
func defaultPathValidator() PathValidator {
	return securityPathValidator{}
}

// defaultCommandRunner returns the production CommandRunner.
func defaultCommandRunner() CommandRunner {
	return securityCommandRunner{}
}
