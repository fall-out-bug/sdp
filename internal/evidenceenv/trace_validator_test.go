package evidenceenv

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestValidateTraceChain_Complete(t *testing.T) {
	events := []TraceEvent{
		{At: "2025-01-01T10:00:00Z", Phase: "claimed"},
		{At: "2025-01-01T10:01:00Z", Phase: "execute"},
		{At: "2025-01-01T10:02:00Z", Phase: "verify"},
		{At: "2025-01-01T10:03:00Z", Phase: "publish"},
	}
	res := ValidateTraceChain(events)
	if !res.OK {
		t.Errorf("expected OK=true for complete chain, got OK=false, missing=%v", res.Missing)
	}
	if len(res.Warnings) > 0 {
		t.Errorf("expected no warnings for complete chain, got %v", res.Warnings)
	}
}

func TestValidateTraceChain_Incomplete(t *testing.T) {
	events := []TraceEvent{
		{At: "2025-01-01T10:00:00Z", Phase: "claimed"},
		{At: "2025-01-01T10:01:00Z", Phase: "execute"},
		// missing verify and publish/review
	}
	res := ValidateTraceChain(events)
	if res.OK {
		t.Errorf("expected OK=false for incomplete chain")
	}
	if len(res.Missing) == 0 {
		t.Errorf("expected missing phases, got %v", res.Missing)
	}
	if len(res.Warnings) == 0 {
		t.Errorf("expected warnings for incomplete chain, got %v", res.Warnings)
	}
}

func TestValidateTraceChain_ReviewInsteadOfPublish(t *testing.T) {
	events := []TraceEvent{
		{At: "2025-01-01T10:00:00Z", Phase: "claimed"},
		{At: "2025-01-01T10:01:00Z", Phase: "execute"},
		{At: "2025-01-01T10:02:00Z", Phase: "verify"},
		{At: "2025-01-01T10:03:00Z", Phase: "review"},
	}
	res := ValidateTraceChain(events)
	if !res.OK {
		t.Errorf("expected OK=true when review present, got OK=false, missing=%v", res.Missing)
	}
}

func TestValidateTraceChain_IgnoresHeartbeat(t *testing.T) {
	events := []TraceEvent{
		{At: "2025-01-01T10:00:00Z", Phase: "claimed"},
		{At: "2025-01-01T10:01:00Z", Phase: "heartbeat"},
		{At: "2025-01-01T10:02:00Z", Phase: "execute"},
		{At: "2025-01-01T10:03:00Z", Phase: "verify"},
		{At: "2025-01-01T10:04:00Z", Phase: "publish"},
	}
	res := ValidateTraceChain(events)
	if !res.OK {
		t.Errorf("expected OK=true, heartbeat should be ignored, got OK=false, missing=%v", res.Missing)
	}
}

func TestDetectTraceGaps(t *testing.T) {
	// gap > 5 min between execute and verify
	events := []TraceEvent{
		{At: "2025-01-01T10:00:00Z", Phase: "claimed"},
		{At: "2025-01-01T10:01:00Z", Phase: "execute"},
		{At: "2025-01-01T10:10:00Z", Phase: "verify"}, // 9 min gap
		{At: "2025-01-01T10:11:00Z", Phase: "publish"},
	}
	res := ValidateTraceChain(events)
	if !res.OK {
		t.Fatalf("chain should be complete: %v", res.Missing)
	}
	if len(res.Gaps) == 0 {
		t.Errorf("expected trace gap to be detected")
	}
}

func TestLoadTraceEventsFromRunFile(t *testing.T) {
	dir := t.TempDir()
	runsDir := filepath.Join(dir, ".sdp", "runs")
	if err := os.MkdirAll(runsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(runsDir, "run1.json")
	if err := os.WriteFile(path, []byte(`{"run_id":"run1","events":[{"at":"2025-01-01T10:00:00Z","phase":"execute"}]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	evts := LoadTraceEventsFromRunFile(dir, "run1")
	if len(evts) != 1 || evts[0].Phase != "execute" {
		t.Errorf("expected 1 event with phase execute, got %v", evts)
	}
	evts = LoadTraceEventsFromRunFile(dir, "nonexistent")
	if evts != nil {
		t.Errorf("expected nil for missing file, got %v", evts)
	}
}

func TestAddTraceValidationToEvidence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ev.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"intent":{"issue_id":"test"}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	tvRes := TraceValidationResult{OK: false, Missing: []string{"verify"}, Warnings: []string{"trace incomplete"}}
	if err := AddTraceValidationToEvidence(path, tvRes); err != nil {
		t.Fatalf("AddTraceValidationToEvidence: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := payload["trace_validation"]; !ok {
		t.Error("expected trace_validation in evidence")
	}
}
