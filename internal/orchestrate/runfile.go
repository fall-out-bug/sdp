package orchestrate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fall-out-bug/sdp/internal/sdputil"
)

type runFileJSON struct {
	RunID        string            `json:"run_id"`
	FeatureID    string            `json:"feature_id"`
	Orchestrator string            `json:"orchestrator"`
	Branch       string            `json:"branch"`
	StartedAt    string            `json:"started_at"`
	Events       []runFileEventJSON `json:"events"`
	LastPhase    string            `json:"last_phase"`
	LastState    string            `json:"last_state"`
}

type runFileEventJSON struct {
	At    string `json:"at"`
	Phase string `json:"phase"`
	State string `json:"state"`
}

// EnsureRunFile creates the initial run file for a feature (atomic write).
func EnsureRunFile(dir, featureID, branch string) error {
	if err := sdputil.ValidateFeatureID(featureID); err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	runID := fmt.Sprintf("oneshot-%s-%s", featureID, time.Now().UTC().Format("20060102T150405Z"))
	path := filepath.Join(dir, runID+".json")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("mkdir runs dir: %w", err)
	}
	rf := runFileJSON{
		RunID:        runID,
		FeatureID:    featureID,
		Orchestrator: "sdp-orchestrate",
		Branch:       branch,
		StartedAt:    now,
		Events:       []runFileEventJSON{{At: now, Phase: "init", State: "ok"}},
		LastPhase:    "init",
		LastState:    "ok",
	}
	body, err := json.MarshalIndent(rf, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal run file: %w", err)
	}
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, body, 0o644); err != nil {
		return fmt.Errorf("write run file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("rename run file: %w", err)
	}
	return nil
}
