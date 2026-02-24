package ciloop

import (
	"context"
	"fmt"
	"time"
)

// DeterministicFirstFixer wraps an inner Fixer: tries deterministic fixers first,
// only invokes inner Fixer if they don't produce changes.
type DeterministicFirstFixer struct {
	ProjectRoot   string
	Registry      *AutofixerRegistry
	Runner        CommandRunner
	Committer     Committer
	LogFetcher    LogFetcher
	DecisionLog   func(decision, rationale string) error
	RunFileLogger func(fixerNames []string, duration time.Duration)
	Inner         Fixer
	PRNumber      int
	Ctx           context.Context // for cancellation (e.g. SIGTERM)
}

// Fix implements Fixer: tries deterministic fixers first, then inner Fixer.
func (d *DeterministicFirstFixer) Fix(checks []CheckResult) error {
	log, err := d.LogFetcher.FailedLogs(d.PRNumber)
	if err != nil {
		return fmt.Errorf("fetch CI logs: %w", err)
	}

	ctx := d.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	changed, err := RunDeterministicFixers(ctx, d.ProjectRoot, log, d.Registry, d.Committer, d.DecisionLog, d.RunFileLogger)
	if err != nil {
		return err
	}
	if changed {
		return nil
	}

	return d.Inner.Fix(checks)
}
