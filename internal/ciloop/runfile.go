package ciloop

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/sdputil"
)

// RunEvent is a single event appended to a run file.
type RunEvent struct {
	At    string `json:"at"`
	Phase string `json:"phase"`
	State string `json:"state"`
	Notes string `json:"notes,omitempty"`
}

// RunFile mirrors the .sdp/runs/{run-id}.json schema.
type RunFile struct {
	RunID        string     `json:"run_id"`
	FeatureID    string     `json:"feature_id"`
	Orchestrator string     `json:"orchestrator"`
	Branch       string     `json:"branch"`
	StartedAt    string     `json:"started_at"`
	Events       []RunEvent `json:"events"`
	LastPhase    string     `json:"last_phase"`
	LastState    string     `json:"last_state"`
}

// maxRunEventFieldBytes caps phase/state/notes length to avoid disk DoS.
const maxRunEventFieldBytes = 1024

func truncateField(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}

// AppendRunEvent finds the latest run file for featureID in dir and appends an event.
func AppendRunEvent(dir, featureID, phase, state, notes string) error {
	if err := sdputil.ValidateFeatureID(featureID); err != nil {
		return err
	}
	phase = truncateField(phase, maxRunEventFieldBytes)
	state = truncateField(state, maxRunEventFieldBytes)
	notes = truncateField(notes, maxRunEventFieldBytes)
	path, err := findRunFile(dir, featureID)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read run file: %w", err)
	}
	var rf RunFile
	if err := json.NewDecoder(io.LimitReader(bytes.NewReader(data), sdputil.MaxJSONDecodeBytes)).Decode(&rf); err != nil {
		return fmt.Errorf("parse run file: %w", err)
	}
	rf.Events = append(rf.Events, RunEvent{
		At:    time.Now().UTC().Format(time.RFC3339),
		Phase: phase,
		State: state,
		Notes: notes,
	})
	rf.LastPhase = phase
	rf.LastState = state
	out, err := json.MarshalIndent(rf, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal run file: %w", err)
	}
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, out, 0o644); err != nil {
		return fmt.Errorf("write run file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("rename run file: %w", err)
	}
	return nil
}

func findRunFile(dir, featureID string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("read runs dir %s: %w", dir, err)
	}
	prefix := "oneshot-" + featureID + "-"
	var matches []string
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), prefix) && strings.HasSuffix(e.Name(), ".json") {
			matches = append(matches, e.Name())
		}
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("no run file found for feature %s in %s", featureID, dir)
	}
	sort.Slice(matches, func(i, j int) bool {
		si := strings.TrimSuffix(matches[i], ".json")
		sj := strings.TrimSuffix(matches[j], ".json")
		ni := strings.TrimPrefix(si, prefix)
		nj := strings.TrimPrefix(sj, prefix)
		vi, ei := strconv.Atoi(ni)
		vj, ej := strconv.Atoi(nj)
		if ei == nil && ej == nil {
			return vi < vj // ascending: last in slice = latest
		}
		return si < sj // fallback: string sort (e.g. timestamps)
	})
	return filepath.Join(dir, matches[len(matches)-1]), nil
}
