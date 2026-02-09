package evidence

import (
	"os"
	"path/filepath"
	"testing"
)

func TestModelID(t *testing.T) {
	os.Unsetenv("SDP_MODEL_ID")
	os.Unsetenv("ANTHROPIC_MODEL")
	if got := ModelID(); got != "unknown" {
		t.Errorf("ModelID(): want unknown, got %s", got)
	}
	os.Setenv("SDP_MODEL_ID", "claude-sonnet")
	defer os.Unsetenv("SDP_MODEL_ID")
	if got := ModelID(); got != "claude-sonnet" {
		t.Errorf("ModelID(): want claude-sonnet, got %s", got)
	}
}

func TestPlanEvent(t *testing.T) {
	ev := PlanEvent("00-054-01", []string{"schema/index.json"})
	if ev.Type != "plan" || ev.WSID != "00-054-01" {
		t.Errorf("PlanEvent: got %+v", ev)
	}
	if ev.Data == nil {
		t.Fatal("PlanEvent Data is nil")
	}
}

func TestVerificationEvent(t *testing.T) {
	ev := VerificationEvent("00-054-01", true, "coverage", 85.0)
	if ev.Type != "verification" || ev.WSID != "00-054-01" {
		t.Errorf("VerificationEvent: got %+v", ev)
	}
}

func TestGenerationEvent(t *testing.T) {
	ev := GenerationEvent("00-054-03", []string{"internal/evidence/types.go"})
	if ev.Type != "generation" || ev.WSID != "00-054-03" {
		t.Errorf("GenerationEvent: got %+v", ev)
	}
}

func TestDecisionEvent(t *testing.T) {
	ev := DecisionEvent("00-054-10", "How to store decisions?", "Evidence log", "Single source of truth", []string{"Separate file"}, 0.9, []string{"architecture"}, nil)
	if ev.Type != "decision" || ev.WSID != "00-054-10" {
		t.Errorf("DecisionEvent: got %+v", ev)
	}
	if ev.Data == nil {
		t.Fatal("DecisionEvent Data is nil")
	}
	rev := "evt-123"
	ev2 := DecisionEvent("00-054-10", "Q", "Revert", "Rationale", nil, 0, nil, &rev)
	if ev2.Data == nil {
		t.Fatal("DecisionEvent with reverses: Data is nil")
	}
}

func TestEnabled(t *testing.T) {
	// Enabled() depends on FindProjectRoot and config; just ensure it doesn't panic
	_ = Enabled()
}

func TestEmitSync_Enabled(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".sdp")
	logDir := filepath.Join(cfgDir, "log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cfgPath := filepath.Join(cfgDir, "config.yml")
	cfgContent := "version: 1\nevidence:\n  enabled: true\n  log_path: \".sdp/log/events.jsonl\"\n"
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	origWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(origWd)
	ev := PlanEvent("00-054-05", []string{"internal/evidence/emitter.go"})
	if err := emitSync(ev); err != nil {
		t.Fatalf("emitSync: %v", err)
	}
	logPath := filepath.Join(dir, ".sdp", "log", "events.jsonl")
	if _, err := os.Stat(logPath); err != nil {
		t.Errorf("events.jsonl not created: %v", err)
	}
}
