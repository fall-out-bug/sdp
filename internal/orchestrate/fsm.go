package orchestrate

import (
	"fmt"
	"strings"
)

// TransitionKey identifies a state transition in the FSM.
type TransitionKey struct {
	From string
	To   string
}

// TransitionCondition describes when a transition is valid.
type TransitionCondition struct {
	// AllWorkstreamsDone is true when the transition requires all workstreams to be complete.
	AllWorkstreamsDone bool
	// ReviewApproved is true when the transition requires an approved review.
	ReviewApproved bool
	// Description explains the transition.
	Description string
}

// validTransitions is the declared FSM for the orchestrate state machine.
// Any transition not listed here is invalid and will be rejected.
var validTransitions = map[TransitionKey]TransitionCondition{
	{PhaseInit, PhaseBuild}: {
		Description: "init → build: begin workstream execution",
	},
	{PhaseBuild, PhaseBuild}: {
		Description: "build → build: complete one workstream, continue to next",
	},
	{PhaseBuild, PhaseReview}: {
		AllWorkstreamsDone: true,
		Description:        "build → review: all workstreams done, proceed to review",
	},
	{PhaseReview, PhasePR}: {
		ReviewApproved: true,
		Description:    "review → pr: review approved, create PR",
	},
	{PhasePR, PhaseCI}: {
		Description: "pr → ci: PR created, monitor CI",
	},
	{PhaseCI, PhaseDone}: {
		Description: "ci → done: CI passed, feature complete",
	},
	{PhaseDone, PhaseDone}: {
		Description: "done → done: idempotent (already complete)",
	},
}

// FSMViolationError is returned when a transition violates the FSM.
type FSMViolationError struct {
	From string
	To   string
	Why  string
}

func (e *FSMViolationError) Error() string {
	return fmt.Sprintf("FSM violation: %s → %s: %s", e.From, e.To, e.Why)
}

// ValidateTransition checks that a transition from `from` to `to` is declared
// in the FSM and that any conditions are met.
func ValidateTransition(from string, to string, cp *Checkpoint, workstreams []string) error {
	key := TransitionKey{From: from, To: to}
	cond, ok := validTransitions[key]
	if !ok {
		// Build error message with allowed transitions from current state
		var allowed []string
		for k := range validTransitions {
			if k.From == from {
				allowed = append(allowed, k.To)
			}
		}
		return &FSMViolationError{
			From: from,
			To:   to,
			Why:  fmt.Sprintf("not a valid transition (allowed from %s: [%s])", from, strings.Join(allowed, ", ")),
		}
	}

	if cond.AllWorkstreamsDone {
		allDone := true
		for _, ws := range cp.Workstreams {
			if ws.Status != "done" {
				allDone = false
				break
			}
		}
		if !allDone {
			return &FSMViolationError{
				From: from,
				To:   to,
				Why:  "condition not met: not all workstreams are done",
			}
		}
	}

	if cond.ReviewApproved {
		if cp.Review == nil || cp.Review.Status != "approved" {
			return &FSMViolationError{
				From: from,
				To:   to,
				Why:  "condition not met: review not approved",
			}
		}
	}

	return nil
}

// computeNextPhase determines what phase `Advance` will transition to.
// Used to pre-validate transitions before calling `Advance`.
func computeNextPhase(cp *Checkpoint, workstreams []string) string {
	switch cp.Phase {
	case PhaseInit:
		return PhaseBuild
	case PhaseBuild:
		// Count done workstreams (assume one more will be done after this advance)
		donePlus1 := 0
		for _, ws := range cp.Workstreams {
			if ws.Status == "done" {
				donePlus1++
			}
		}
		donePlus1++ // the current one being advanced
		if donePlus1 >= len(cp.Workstreams) {
			return PhaseReview
		}
		return PhaseBuild
	case PhaseReview:
		return PhasePR
	case PhasePR:
		return PhaseCI
	case PhaseCI:
		return PhaseDone
	default:
		return cp.Phase
	}
}

// ValidateAdvance pre-validates the transition that `Advance` will perform.
// Call this before `Advance` to enforce FSM conformance.
func ValidateAdvance(cp *Checkpoint, workstreams []string) error {
	to := computeNextPhase(cp, workstreams)
	return ValidateTransition(cp.Phase, to, cp, workstreams)
}

// FSMLog describes a recorded state transition for audit purposes.
type FSMLog struct {
	FeatureID string `json:"feature_id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Timestamp string `json:"timestamp"`
	WSID      string `json:"ws_id,omitempty"`
}
