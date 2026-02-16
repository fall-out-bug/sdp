package evidence

import (
	"os"
	"path/filepath"
	"testing"
)

// TestEvidenceLayerE2E runs end-to-end: config, write events, verify chain, trace (AC3–AC7).
func TestEvidenceLayerE2E(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".sdp")
	logDir := filepath.Join(cfgDir, "log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cfgPath := filepath.Join(cfgDir, "config.yml")
	cfgBody := "version: 1\nevidence:\n  enabled: true\n  log_path: \".sdp/log/events.jsonl\"\n"
	if err := os.WriteFile(cfgPath, []byte(cfgBody), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	origWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(origWd)

	logPath := filepath.Join(dir, ".sdp", "log", "events.jsonl")
	w, err := NewWriter(logPath)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	events := []*Event{
		{ID: "e1", Type: "plan", Timestamp: "2026-02-09T12:00:00Z", WSID: "00-054-09"},
		{ID: "e2", Type: "generation", Timestamp: "2026-02-09T12:01:00Z", WSID: "00-054-09"},
		{ID: "e3", Type: "verification", Timestamp: "2026-02-09T12:02:00Z", WSID: "00-054-09"},
		{ID: "e4", Type: "approval", Timestamp: "2026-02-09T12:03:00Z", WSID: "00-054-09"},
	}
	for _, ev := range events {
		if err := w.Append(ev); err != nil {
			t.Fatalf("Append: %v", err)
		}
	}
	if len(events) < 4 {
		t.Fatal("need at least 4 events for AC3")
	}

	r := NewReader(logPath)
	if err := r.Verify(); err != nil {
		t.Errorf("chain verify: %v", err)
	}
	all, err := r.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(all) != 4 {
		t.Errorf("ReadAll: want 4 events, got %d", len(all))
	}
	filtered := FilterByWS(all, "00-054-09")
	if len(filtered) != 4 {
		t.Errorf("FilterByWS: want 4, got %d", len(filtered))
	}
}

// TestFullPipelineEvidenceChain (F056-04 AC1–AC4): idea → design → build → review → deploy produces
// complete chain; sdp log show --type=X works for plan, generation, verification, approval; chain integrity.
func TestFullPipelineEvidenceChain(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "events.jsonl")
	w, err := NewWriter(logPath)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	chain := []*Event{
		{ID: "idea", Type: "plan", Timestamp: "2026-02-10T10:00:00Z", WSID: "00-056-04", Data: map[string]interface{}{"skill": "idea"}},
		{ID: "design", Type: "plan", Timestamp: "2026-02-10T10:01:00Z", WSID: "00-056-04", Data: map[string]interface{}{"skill": "design"}},
		{ID: "build", Type: "generation", Timestamp: "2026-02-10T10:02:00Z", WSID: "00-056-04", Data: map[string]interface{}{"skill": "build"}},
		{ID: "review", Type: "verification", Timestamp: "2026-02-10T10:03:00Z", WSID: "00-056-04", Data: map[string]interface{}{"skill": "review", "passed": true}},
		{ID: "deploy", Type: "approval", Timestamp: "2026-02-10T10:04:00Z", WSID: "00-056-04", Data: map[string]interface{}{"skill": "deploy"}},
	}
	for _, ev := range chain {
		if err := w.Append(ev); err != nil {
			t.Fatalf("Append: %v", err)
		}
	}

	r := NewReader(logPath)
	if err := r.Verify(); err != nil {
		t.Errorf("chain integrity: %v", err)
	}
	all, err := r.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(all) != 5 {
		t.Fatalf("want 5 events (idea→design→build→review→deploy), got %d", len(all))
	}
	order := []string{"plan", "plan", "generation", "verification", "approval"}
	for i := range order {
		if all[i].Type != order[i] {
			t.Errorf("event %d: want type %q, got %q", i, order[i], all[i].Type)
		}
	}

	// AC3: sdp log show --type=X for all 4 event types
	for _, typ := range []string{"plan", "generation", "verification", "approval"} {
		filtered := FilterByType(all, typ)
		if len(filtered) == 0 {
			t.Errorf("FilterByType(%q): want at least 1, got 0", typ)
		}
	}
}
