package orchestrate_test

import (
	"errors"
	"testing"

	"github.com/fall-out-bug/sdp/internal/orchestrate"
)

func TestValidateTransition_Valid(t *testing.T) {
	tests := []struct {
		name string
		from string
		to   string
		cp   *orchestrate.Checkpoint
	}{
		{
			name: "init to build",
			from: orchestrate.PhaseInit,
			to:   orchestrate.PhaseBuild,
			cp:   &orchestrate.Checkpoint{Phase: orchestrate.PhaseInit},
		},
		{
			name: "build to build (more workstreams)",
			from: orchestrate.PhaseBuild,
			to:   orchestrate.PhaseBuild,
			cp:   &orchestrate.Checkpoint{Phase: orchestrate.PhaseBuild},
		},
		{
			name: "build to review (all done)",
			from: orchestrate.PhaseBuild,
			to:   orchestrate.PhaseReview,
			cp: &orchestrate.Checkpoint{
				Phase: orchestrate.PhaseBuild,
				Workstreams: []orchestrate.WSStatus{
					{ID: "00-028-01", Status: "done"},
				},
			},
		},
		{
			name: "review to pr (approved)",
			from: orchestrate.PhaseReview,
			to:   orchestrate.PhasePR,
			cp: &orchestrate.Checkpoint{
				Phase:  orchestrate.PhaseReview,
				Review: &orchestrate.ReviewStatus{Status: "approved"},
			},
		},
		{
			name: "pr to ci",
			from: orchestrate.PhasePR,
			to:   orchestrate.PhaseCI,
			cp:   &orchestrate.Checkpoint{Phase: orchestrate.PhasePR},
		},
		{
			name: "ci to done",
			from: orchestrate.PhaseCI,
			to:   orchestrate.PhaseDone,
			cp:   &orchestrate.Checkpoint{Phase: orchestrate.PhaseCI},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := orchestrate.ValidateTransition(tt.from, tt.to, tt.cp, nil)
			if err != nil {
				t.Errorf("expected valid transition, got error: %v", err)
			}
		})
	}
}

func TestValidateTransition_Invalid(t *testing.T) {
	tests := []struct {
		name string
		from string
		to   string
		cp   *orchestrate.Checkpoint
	}{
		{
			name: "init to done (skip phases)",
			from: orchestrate.PhaseInit,
			to:   orchestrate.PhaseDone,
			cp:   &orchestrate.Checkpoint{Phase: orchestrate.PhaseInit},
		},
		{
			name: "build to done (skip review+pr+ci)",
			from: orchestrate.PhaseBuild,
			to:   orchestrate.PhaseDone,
			cp:   &orchestrate.Checkpoint{Phase: orchestrate.PhaseBuild},
		},
		{
			name: "review to done (skip pr+ci)",
			from: orchestrate.PhaseReview,
			to:   orchestrate.PhaseDone,
			cp:   &orchestrate.Checkpoint{Phase: orchestrate.PhaseReview},
		},
		{
			name: "build to review but workstreams not done",
			from: orchestrate.PhaseBuild,
			to:   orchestrate.PhaseReview,
			cp: &orchestrate.Checkpoint{
				Phase: orchestrate.PhaseBuild,
				Workstreams: []orchestrate.WSStatus{
					{ID: "00-028-01", Status: "pending"},
				},
			},
		},
		{
			name: "review to pr but review not approved",
			from: orchestrate.PhaseReview,
			to:   orchestrate.PhasePR,
			cp: &orchestrate.Checkpoint{
				Phase:  orchestrate.PhaseReview,
				Review: &orchestrate.ReviewStatus{Status: "pending"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := orchestrate.ValidateTransition(tt.from, tt.to, tt.cp, nil)
			if err == nil {
				t.Errorf("expected error for invalid transition %s→%s, got nil", tt.from, tt.to)
			}
			var fsmErr *orchestrate.FSMViolationError
			if !errors.As(err, &fsmErr) {
				t.Errorf("expected FSMViolationError, got %T: %v", err, err)
			}
		})
	}
}

func TestValidateAdvance_PreCheck(t *testing.T) {
	cp := &orchestrate.Checkpoint{
		Phase: orchestrate.PhaseInit,
	}
	workstreams := []string{"00-028-01"}
	// init → build should be valid
	if err := orchestrate.ValidateAdvance(cp, workstreams); err != nil {
		t.Errorf("ValidateAdvance from init: unexpected error: %v", err)
	}
}
