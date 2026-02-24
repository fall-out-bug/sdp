package ciloop

import (
	"context"
	"time"
)

// LoopResult is the outcome of RunLoop.
type LoopResult int

const (
	ResultGreen     LoopResult = iota // all checks passed
	ResultEscalated                   // escalation triggered
	ResultMaxIter                     // max iterations exceeded
)

// DefaultMaxPendingRetries is the default cap on PENDING-only polling rounds.
// A round is a poll that returns only PENDING/IN_PROGRESS checks.
// Zero means unlimited (use for short-lived tests only).
const DefaultMaxPendingRetries = 60

// Fixer attempts to fix a set of auto-fixable failing checks.
// Returns an error if the fix cannot be applied.
type Fixer interface {
	Fix(checks []CheckResult) error
}

// LoopOptions configures RunLoop behaviour.
type LoopOptions struct {
	// Context allows cancellation (e.g. SIGINT/SIGTERM). When cancelled, RunLoop returns ResultEscalated.
	Context context.Context
	PRNumber int
	MaxIter  int
	// MaxPendingRetries caps how many consecutive PENDING-only rounds before escalation.
	// Zero disables the cap (tests only).
	MaxPendingRetries int
	PollDelay         time.Duration
	RetryDelay        time.Duration
	Poller            *Poller
	// OnEscalate is called when a non-auto-fixable failure is detected or Fixer is nil.
	OnEscalate func(checks []CheckResult) error
	// OnPollError is called when GetChecks fails (before returning). Use to save checkpoint defensively.
	OnPollError func(err error)
	// Fixer handles auto-fixable failures.
	// When nil, auto-fixable failures escalate immediately (same as non-auto-fixable).
	Fixer Fixer
}

// RunLoop polls CI checks until green, escalation, or max iterations.
//
// PENDING/IN_PROGRESS checks trigger a RetryDelay wait without consuming an iteration.
// Up to MaxPendingRetries consecutive pending-only rounds are allowed; after that, escalate.
// FAILURE checks are classified: non-auto-fixable (or auto-fixable with nil Fixer) → escalate.
// Auto-fixable failures with a Fixer: call Fixer.Fix, increment iter, re-poll.
//
// Exit criteria:
//   - ResultGreen     when IsAllGreen
//   - ResultEscalated when OnEscalate is called or on error
//   - ResultMaxIter   when iter >= MaxIter
func RunLoop(opts LoopOptions) (LoopResult, error) {
	iter := 0
	pendingRounds := 0
	for {
		if opts.Context != nil {
			select {
			case <-opts.Context.Done():
				return ResultEscalated, opts.Context.Err()
			default:
			}
		}
		if opts.PollDelay > 0 {
			if opts.Context != nil {
				select {
				case <-opts.Context.Done():
					return ResultEscalated, opts.Context.Err()
				case <-time.After(opts.PollDelay):
				}
			} else {
				time.Sleep(opts.PollDelay)
			}
		}

		checks, err := opts.Poller.GetChecks(opts.PRNumber)
		if err != nil {
			if opts.OnPollError != nil {
				opts.OnPollError(err)
			}
			return ResultEscalated, err
		}

		if IsAllGreen(checks) {
			return ResultGreen, nil
		}

		pending := FilterByState(checks, StatePending)
		inProgress := FilterByState(checks, StateInProgress)
		if len(pending)+len(inProgress) > 0 {
			pendingRounds++
			if opts.MaxPendingRetries > 0 && pendingRounds >= opts.MaxPendingRetries {
				if opts.OnEscalate != nil {
					if err := opts.OnEscalate(checks); err != nil {
						return ResultEscalated, err
					}
				}
				return ResultEscalated, nil
			}
			if opts.RetryDelay > 0 {
				if opts.Context != nil {
					select {
					case <-opts.Context.Done():
						return ResultEscalated, opts.Context.Err()
					case <-time.After(opts.RetryDelay):
					}
				} else {
					time.Sleep(opts.RetryDelay)
				}
			}
			continue
		}
		pendingRounds = 0

		failing := append(FilterByState(checks, StateFailure), FilterByState(checks, StateError)...)
		if len(failing) == 0 {
			return ResultGreen, nil
		}

		escalateChecks := make([]CheckResult, 0)
		autoFixChecks := make([]CheckResult, 0)
		for _, c := range failing {
			if Classify(c.Name) == ClassAutoFixable && opts.Fixer != nil {
				autoFixChecks = append(autoFixChecks, c)
			} else {
				escalateChecks = append(escalateChecks, c)
			}
		}

		if len(escalateChecks) > 0 {
			if opts.OnEscalate != nil {
				if err := opts.OnEscalate(escalateChecks); err != nil {
					return ResultEscalated, err
				}
			}
			return ResultEscalated, nil
		}

		// Auto-fixable failures with Fixer: count iteration and attempt fix.
		iter++
		if iter >= opts.MaxIter {
			return ResultMaxIter, nil
		}

		if err := opts.Fixer.Fix(autoFixChecks); err != nil {
			if opts.OnEscalate != nil {
				if escErr := opts.OnEscalate(autoFixChecks); escErr != nil {
					return ResultEscalated, escErr
				}
			}
			return ResultEscalated, err
		}
	}
}
