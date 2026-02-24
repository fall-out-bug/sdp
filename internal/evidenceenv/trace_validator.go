package evidenceenv

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TraceEvent is a minimal event for trace validation (phase only).
type TraceEvent struct {
	At    string
	Phase string
}

// TraceValidationResult holds the outcome of trace chain validation.
type TraceValidationResult struct {
	OK       bool     `json:"ok"`
	Missing  []string `json:"missing"`
	Warnings []string `json:"warnings"`
	Gaps     []string `json:"gaps,omitempty"`
}

// RequiredPhasesForSuccess are phases that must appear in a complete run trace.
// At least one of review/publish is required.
var RequiredPhasesForSuccess = []string{"execute", "verify"}

// OptionalTerminalPhases - at least one must be present for a complete chain.
var OptionalTerminalPhases = []string{"review", "publish"}

// ValidateTraceChain checks that the trace events contain all required phases.
// Missing phases produce warnings only; terminal transition is not blocked.
func ValidateTraceChain(events []TraceEvent) TraceValidationResult {
	phases := make(map[string]bool)
	var ordered []string
	for _, e := range events {
		p := strings.TrimSpace(e.Phase)
		if p == "" || p == "heartbeat" {
			continue
		}
		if !phases[p] {
			phases[p] = true
			ordered = append(ordered, p)
		}
	}

	var missing []string
	for _, req := range RequiredPhasesForSuccess {
		if !phases[req] {
			missing = append(missing, req)
		}
	}

	hasTerminal := false
	for _, opt := range OptionalTerminalPhases {
		if phases[opt] {
			hasTerminal = true
			break
		}
	}
	if !hasTerminal {
		missing = append(missing, "review|publish")
	}

	var warnings []string
	if len(missing) > 0 {
		warnings = append(warnings, "trace incomplete: missing phases "+strings.Join(missing, ", "))
	}

	gaps := detectTraceGaps(events)
	if len(gaps) > 0 {
		warnings = append(warnings, "trace gaps: "+strings.Join(gaps, "; "))
	}

	ok := len(missing) == 0
	return TraceValidationResult{
		OK:       ok,
		Missing:  missing,
		Warnings: warnings,
		Gaps:     gaps,
	}
}

// detectTraceGaps finds time gaps > 5 minutes between consecutive non-heartbeat events.
func detectTraceGaps(events []TraceEvent) []string {
	const gapThreshold = 5 * time.Minute
	var gaps []string
	var lastAt time.Time
	for _, e := range events {
		if e.Phase == "heartbeat" {
			continue
		}
		t, err := time.Parse(time.RFC3339Nano, e.At)
		if err != nil {
			t, err = time.Parse(time.RFC3339, e.At)
		}
		if err != nil {
			continue
		}
		if !lastAt.IsZero() && t.Sub(lastAt) > gapThreshold {
			gaps = append(gaps, lastAt.Format("15:04")+"-"+t.Format("15:04")+" ("+e.Phase+")")
		}
		lastAt = t
	}
	return gaps
}

// LoadTraceEventsFromRunFile reads events from a run file at workDir/.sdp/runs/{runID}.json.
// Returns nil if the file does not exist or cannot be parsed.
func LoadTraceEventsFromRunFile(workDir, runID string) []TraceEvent {
	path := filepath.Join(workDir, ".sdp", "runs", runID+".json")
	b, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var doc struct {
		Events []struct {
			At    string `json:"at"`
			Phase string `json:"phase"`
		} `json:"events"`
	}
	if err := json.Unmarshal(b, &doc); err != nil {
		return nil
	}
	out := make([]TraceEvent, len(doc.Events))
	for i, e := range doc.Events {
		out[i] = TraceEvent{At: e.At, Phase: e.Phase}
	}
	return out
}

// AddTraceValidationToEvidence reads an evidence file, adds trace_validation, and writes back.
// Used to report trace gaps and missing phases in the evidence payload.
func AddTraceValidationToEvidence(path string, res TraceValidationResult) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var payload map[string]any
	if err := json.Unmarshal(b, &payload); err != nil {
		return err
	}
	tv := map[string]any{
		"ok":       res.OK,
		"missing":  res.Missing,
		"warnings": res.Warnings,
		"gaps":     res.Gaps,
	}
	payload["trace_validation"] = tv
	out, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(out, '\n'), 0o644)
}
