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
