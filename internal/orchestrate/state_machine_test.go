package orchestrate_test

import (
	"testing"

	"github.com/fall-out-bug/sdp/internal/orchestrate"
)

func TestComputeNextAction(t *testing.T) {
	workstreams := []string{"00-004-01", "00-004-02"}
	projectRoot := "."

	tests := []struct {
		name    string
		cp      *orchestrate.Checkpoint
		wantAct string
		wantWS  string
		wantPR  int
		wantErr bool
	}{
		{
			name: "init returns init action",
			cp: &orchestrate.Checkpoint{
				FeatureID: "F004", Phase: orchestrate.PhaseInit,
				Workstreams: []orchestrate.WSStatus{},
			},
			wantAct: "init",
		},
		{
			name: "build with pending WS returns build",
			cp: &orchestrate.Checkpoint{
				FeatureID: "F004", Phase: orchestrate.PhaseBuild,
				Workstreams: []orchestrate.WSStatus{
					{ID: "00-004-01", Status: "pending"},
					{ID: "00-004-02", Status: "pending"},
				},
			},
			wantAct: "build",
			wantWS:  "00-004-01",
		},
		{
			name: "build with in_progress WS returns build",
			cp: &orchestrate.Checkpoint{
				FeatureID: "F004", Phase: orchestrate.PhaseBuild,
				Workstreams: []orchestrate.WSStatus{
					{ID: "00-004-01", Status: "in_progress"},
					{ID: "00-004-02", Status: "pending"},
				},
			},
			wantAct: "build",
			wantWS:  "00-004-01",
		},
		{
			name: "build all done returns review",
			cp: &orchestrate.Checkpoint{
				FeatureID: "F004", Phase: orchestrate.PhaseBuild,
				Workstreams: []orchestrate.WSStatus{
					{ID: "00-004-01", Status: "done"},
					{ID: "00-004-02", Status: "done"},
				},
			},
			wantAct: "review",
		},
		{
			name: "review returns review",
			cp: &orchestrate.Checkpoint{
				FeatureID: "F004", Phase: orchestrate.PhaseReview,
			},
			wantAct: "review",
		},
		{
			name: "pr returns pr",
			cp: &orchestrate.Checkpoint{
				FeatureID: "F004", Phase: orchestrate.PhasePR,
			},
			wantAct: "pr",
		},
		{
			name: "ci with PRNumber returns ci-loop",
			cp: &orchestrate.Checkpoint{
				FeatureID: "F004", Phase: orchestrate.PhaseCI,
				PRNumber: intPtr(42),
			},
			wantAct: "ci-loop",
			wantPR:  42,
		},
		{
			name: "ci without PRNumber returns ci-loop with 0",
			cp: &orchestrate.Checkpoint{
				FeatureID: "F004", Phase: orchestrate.PhaseCI,
			},
			wantAct: "ci-loop",
			wantPR:  0,
		},
		{
			name: "done returns done",
			cp: &orchestrate.Checkpoint{
				FeatureID: "F004", Phase: orchestrate.PhaseDone,
			},
			wantAct: "done",
		},
		{
			name: "unknown phase returns error",
			cp: &orchestrate.Checkpoint{
				FeatureID: "F004", Phase: "unknown",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orchestrate.ComputeNextAction(tt.cp, workstreams, projectRoot)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Action != tt.wantAct {
				t.Errorf("action = %q, want %q", got.Action, tt.wantAct)
			}
			if tt.wantWS != "" && got.WSID != tt.wantWS {
				t.Errorf("ws_id = %q, want %q", got.WSID, tt.wantWS)
			}
			if tt.wantPR != 0 && got.PR != tt.wantPR {
				t.Errorf("pr = %d, want %d", got.PR, tt.wantPR)
			}
		})
	}
}

func intPtr(n int) *int { return &n }

func TestAdvanceFullLifecycle(t *testing.T) {
	workstreams := []string{"00-004-01", "00-004-02"}

	t.Run("init to build", func(t *testing.T) {
		cp := orchestrate.CreateInitialCheckpoint("F004", "feature/F004-x", workstreams)
		if err := orchestrate.Advance(cp, workstreams, ""); err != nil {
			t.Fatal(err)
		}
		if cp.Phase != orchestrate.PhaseBuild {
			t.Errorf("phase = %q, want build", cp.Phase)
		}
		if len(cp.Workstreams) != 2 {
			t.Errorf("workstreams = %d, want 2", len(cp.Workstreams))
		}
		for i, ws := range cp.Workstreams {
			if ws.Status != "pending" {
				t.Errorf("workstream[%d].status = %q, want pending", i, ws.Status)
			}
		}
	})

	t.Run("build first WS to build second WS", func(t *testing.T) {
		cp := &orchestrate.Checkpoint{
			FeatureID: "F004", Phase: orchestrate.PhaseBuild,
			Workstreams: []orchestrate.WSStatus{
				{ID: "00-004-01", Status: "pending"},
				{ID: "00-004-02", Status: "pending"},
			},
		}
		if err := orchestrate.Advance(cp, workstreams, "abc123"); err != nil {
			t.Fatal(err)
		}
		if cp.Phase != orchestrate.PhaseBuild {
			t.Errorf("phase = %q, want build (second WS)", cp.Phase)
		}
		if cp.Workstreams[0].Status != "done" || cp.Workstreams[0].Commit != "abc123" {
			t.Errorf("first WS should be done with commit abc123, got %+v", cp.Workstreams[0])
		}
		if cp.Workstreams[1].Status != "pending" {
			t.Errorf("second WS should still be pending, got %q", cp.Workstreams[1].Status)
		}
	})

	t.Run("build all done to review", func(t *testing.T) {
		cp := &orchestrate.Checkpoint{
			FeatureID: "F004", Phase: orchestrate.PhaseBuild,
			Workstreams: []orchestrate.WSStatus{
				{ID: "00-004-01", Status: "done"},
				{ID: "00-004-02", Status: "done"},
			},
		}
		if err := orchestrate.Advance(cp, workstreams, ""); err != nil {
			t.Fatal(err)
		}
		if cp.Phase != orchestrate.PhaseReview {
			t.Errorf("phase = %q, want review", cp.Phase)
		}
	})

	t.Run("review to pr", func(t *testing.T) {
		cp := &orchestrate.Checkpoint{
			FeatureID: "F004", Phase: orchestrate.PhaseReview,
			Review: &orchestrate.ReviewStatus{Status: "pending"},
		}
		if err := orchestrate.Advance(cp, workstreams, ""); err != nil {
			t.Fatal(err)
		}
		if cp.Phase != orchestrate.PhasePR {
			t.Errorf("phase = %q, want pr", cp.Phase)
		}
		if cp.Review != nil && cp.Review.Status != "approved" {
			t.Errorf("review status = %q, want approved", cp.Review.Status)
		}
	})

	t.Run("pr to ci", func(t *testing.T) {
		cp := &orchestrate.Checkpoint{
			FeatureID: "F004", Phase: orchestrate.PhasePR,
		}
		if err := orchestrate.Advance(cp, workstreams, ""); err != nil {
			t.Fatal(err)
		}
		if cp.Phase != orchestrate.PhaseCI {
			t.Errorf("phase = %q, want ci", cp.Phase)
		}
	})

	t.Run("ci to done", func(t *testing.T) {
		cp := &orchestrate.Checkpoint{
			FeatureID: "F004", Phase: orchestrate.PhaseCI,
		}
		if err := orchestrate.Advance(cp, workstreams, ""); err != nil {
			t.Fatal(err)
		}
		if cp.Phase != orchestrate.PhaseDone {
			t.Errorf("phase = %q, want done", cp.Phase)
		}
	})

	t.Run("done to done no-op", func(t *testing.T) {
		cp := &orchestrate.Checkpoint{
			FeatureID: "F004", Phase: orchestrate.PhaseDone,
		}
		if err := orchestrate.Advance(cp, workstreams, ""); err != nil {
			t.Fatal(err)
		}
		if cp.Phase != orchestrate.PhaseDone {
			t.Errorf("phase = %q, want done (no-op)", cp.Phase)
		}
	})
}

func TestAdvanceInitToBuild(t *testing.T) {
	cp := orchestrate.CreateInitialCheckpoint("F004", "feature/F004-x", []string{"00-004-01", "00-004-02"})
	if cp.Phase != orchestrate.PhaseInit {
		t.Errorf("expected init phase, got %s", cp.Phase)
	}
	err := orchestrate.Advance(cp, []string{"00-004-01", "00-004-02"}, "")
	if err != nil {
		t.Fatal(err)
	}
	if cp.Phase != orchestrate.PhaseBuild {
		t.Errorf("expected build phase, got %s", cp.Phase)
	}
	if len(cp.Workstreams) != 2 {
		t.Errorf("expected 2 workstreams, got %d", len(cp.Workstreams))
	}
}

func TestAdvanceBuildToReview(t *testing.T) {
	cp := &orchestrate.Checkpoint{
		FeatureID:   "F004",
		Phase:       orchestrate.PhaseBuild,
		Workstreams: []orchestrate.WSStatus{
			{ID: "00-004-01", Status: "done"},
			{ID: "00-004-02", Status: "done"},
		},
	}
	err := orchestrate.Advance(cp, []string{"00-004-01", "00-004-02"}, "")
	if err != nil {
		t.Fatal(err)
	}
	if cp.Phase != orchestrate.PhaseReview {
		t.Errorf("expected review phase, got %s", cp.Phase)
	}
}
