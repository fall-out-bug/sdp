package orchestrate

import (
	"fmt"
	"path/filepath"
)

// NextAction describes what the agent should do next.
type NextAction struct {
	Action  string `json:"action"`  // build, review, pr, ci-loop, done
	WSID    string `json:"ws_id,omitempty"`
	Feature string `json:"feature,omitempty"`
	PR      int    `json:"pr,omitempty"`
}

// ComputeNextAction returns the next action based on checkpoint state.
func ComputeNextAction(cp *Checkpoint, workstreams []string, projectRoot string) (*NextAction, error) {
	switch cp.Phase {
	case PhaseInit:
		return &NextAction{Action: "init"}, nil
	case PhaseBuild:
		for i, ws := range cp.Workstreams {
			if ws.Status != "done" {
				if ws.Status == "pending" {
					return &NextAction{Action: "build", WSID: workstreams[i], Feature: cp.FeatureID}, nil
				}
				return &NextAction{Action: "build", WSID: ws.ID, Feature: cp.FeatureID}, nil
			}
		}
		return &NextAction{Action: "review", Feature: cp.FeatureID}, nil
	case PhaseReview:
		return &NextAction{Action: "review", Feature: cp.FeatureID}, nil
	case PhasePR:
		return &NextAction{Action: "pr", Feature: cp.FeatureID}, nil
	case PhaseCI:
		pr := 0
		if cp.PRNumber != nil {
			pr = *cp.PRNumber
		}
		return &NextAction{Action: "ci-loop", Feature: cp.FeatureID, PR: pr}, nil
	case PhaseDone:
		return &NextAction{Action: "done"}, nil
	default:
		return nil, fmt.Errorf("unknown phase %q", cp.Phase)
	}
}

// CurrentBuildWS returns the workstream ID being built (first non-done) when in build phase.
func CurrentBuildWS(cp *Checkpoint) string {
	if cp.Phase != PhaseBuild {
		return ""
	}
	for _, ws := range cp.Workstreams {
		if ws.Status != "done" {
			return ws.ID
		}
	}
	return ""
}

// Advance transitions the checkpoint to the next phase.
// For build phase, result is the commit hash of the completed workstream.
func Advance(cp *Checkpoint, workstreams []string, result string) error {
	switch cp.Phase {
	case PhaseInit:
		cp.Phase = PhaseBuild
		cp.Workstreams = make([]WSStatus, len(workstreams))
		for i, ws := range workstreams {
			cp.Workstreams[i] = WSStatus{ID: ws, Status: "pending"}
		}
		return nil
	case PhaseBuild:
		for i := range cp.Workstreams {
			if cp.Workstreams[i].Status != "done" {
				cp.Workstreams[i].Status = "done"
				if result != "" {
					cp.Workstreams[i].Commit = result
				}
				cp.Workstreams[i].Attempts++
				break
			}
		}
		allDone := true
		for _, ws := range cp.Workstreams {
			if ws.Status != "done" {
				allDone = false
				break
			}
		}
		if allDone {
			cp.Phase = PhaseReview
			if cp.Review == nil {
				cp.Review = &ReviewStatus{Iteration: 0, Status: "pending"}
			}
		}
		return nil
	case PhaseReview:
		cp.Phase = PhasePR
		if cp.Review != nil {
			cp.Review.Status = "approved"
		}
		return nil
	case PhasePR:
		cp.Phase = PhaseCI
		return nil
	case PhaseCI:
		cp.Phase = PhaseDone
		return nil
	case PhaseDone:
		return nil
	default:
		return fmt.Errorf("unknown phase %q", cp.Phase)
	}
}

// CreateInitialCheckpoint builds a new checkpoint for a feature.
func CreateInitialCheckpoint(featureID, branch string, workstreams []string) *Checkpoint {
	ws := make([]WSStatus, len(workstreams))
	for i, id := range workstreams {
		ws[i] = WSStatus{ID: id, Status: "pending"}
	}
	return &Checkpoint{
		Schema:      "1.0",
		FeatureID:   featureID,
		Branch:      branch,
		Phase:       PhaseInit,
		Workstreams: ws,
		Review:      &ReviewStatus{Iteration: 0, Status: "pending"},
	}
}

// FindProjectRoot walks up from dir to find a directory containing docs/workstreams.
func FindProjectRoot(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	for d := abs; d != "" && d != "/"; d = filepath.Dir(d) {
		check := filepath.Join(d, "docs", "workstreams", "backlog")
		if _, err := filepath.Glob(filepath.Join(check, "*.md")); err == nil {
			if ents, _ := filepath.Glob(filepath.Join(check, "*.md")); len(ents) > 0 {
				return d, nil
			}
		}
	}
	return "", fmt.Errorf("project root not found (no docs/workstreams/backlog)")
}
