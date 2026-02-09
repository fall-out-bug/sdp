package evidence

import (
	"os"
	"path/filepath"
	"testing"
	"time"
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

func TestPlanEventWithFeature(t *testing.T) {
	ev := PlanEventWithFeature("00-056-02", "F056", []string{"parse.go", "emitter.go"})
	if ev.Type != "plan" || ev.WSID != "00-056-02" {
		t.Errorf("PlanEventWithFeature: got %+v", ev)
	}
	m, _ := ev.Data.(map[string]interface{})
	if m["feature_id"] != "F056" {
		t.Errorf("feature_id: got %v", m["feature_id"])
	}
}

func TestPlanEventForDesign(t *testing.T) {
	ev := PlanEventForDesign("00-056-00", "F056", 4, []string{"a.md", "b.md"}, map[string]interface{}{"deps": "00-054-05"})
	if ev.Type != "plan" || ev.WSID != "00-056-00" {
		t.Errorf("PlanEventForDesign: got %+v", ev)
	}
	m, _ := ev.Data.(map[string]interface{})
	if m["feature_id"] != "F056" || m["ws_count"].(int) != 4 || m["deps"] != "00-054-05" {
		t.Errorf("PlanEventForDesign data: got %v", m)
	}
}

func TestPlanEventForIdea(t *testing.T) {
	qa := []QAPair{{Q: "Who?", A: "User"}, {Q: "What?", A: "Auth"}}
	ev := PlanEventForIdea("00-056-00", "F056", 2, "Auth requirements", qa)
	if ev.Type != "plan" || ev.WSID != "00-056-00" {
		t.Errorf("PlanEventForIdea: got %+v", ev)
	}
	m, _ := ev.Data.(map[string]interface{})
	if m["feature_id"] != "F056" || m["question_count"].(int) != 2 || m["summary"] != "Auth requirements" {
		t.Errorf("PlanEventForIdea data: got %v", m)
	}
	if m["qa_pairs"] == nil {
		t.Error("qa_pairs missing")
	}
}

func TestVerificationEvent(t *testing.T) {
	ev := VerificationEvent("00-054-01", true, "coverage", 85.0)
	if ev.Type != "verification" || ev.WSID != "00-054-01" {
		t.Errorf("VerificationEvent: got %+v", ev)
	}
}

func TestVerificationEventWithFindings(t *testing.T) {
	ev := VerificationEventWithFindings("00-056-01", false, "QA", 82.0, "Coverage below threshold")
	if ev.Type != "verification" || ev.WSID != "00-056-01" {
		t.Errorf("VerificationEventWithFindings: got %+v", ev)
	}
	if ev.Data == nil {
		t.Fatal("Data is nil")
	}
	m, _ := ev.Data.(map[string]interface{})
	if m["findings"] != "Coverage below threshold" {
		t.Errorf("findings: got %v", m["findings"])
	}
}

func TestApprovalEvent(t *testing.T) {
	ev := ApprovalEvent("00-000-00", "main", "abc123def", "CI")
	if ev.Type != "approval" || ev.WSID != "00-000-00" {
		t.Errorf("ApprovalEvent: got %+v", ev)
	}
	if ev.Data == nil {
		t.Fatal("Data is nil")
	}
	m, _ := ev.Data.(map[string]interface{})
	if m["target_branch"] != "main" || m["commit_sha"] != "abc123def" || m["approved_by"] != "CI" {
		t.Errorf("ApprovalEvent data: got %v", m)
	}
}

func TestGenerationEvent(t *testing.T) {
	ev := GenerationEvent("00-054-03", []string{"internal/evidence/types.go"})
	if ev.Type != "generation" || ev.WSID != "00-054-03" {
		t.Errorf("GenerationEvent: got %+v", ev)
	}
}

func TestLessonEvent(t *testing.T) {
	l := Lesson{WSID: "00-054-11", Outcome: "passed", WhatWorked: []string{"A: ok"}, Category: "verification"}
	ev := LessonEvent(l)
	if ev.Type != "lesson" || ev.WSID != "00-054-11" {
		t.Errorf("LessonEvent: got %+v", ev)
	}
	if ev.Data == nil {
		t.Fatal("LessonEvent Data is nil")
	}
	l2 := Lesson{WSID: "00-054-13", Outcome: "failed", WhatFailed: []string{"B: fail"}}
	ev2 := LessonEvent(l2)
	if ev2.Type != "lesson" {
		t.Errorf("LessonEvent failed: got %+v", ev2)
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
	if err := EmitSync(ev); err != nil {
		t.Fatalf("EmitSync: %v", err)
	}
	logPath := filepath.Join(dir, ".sdp", "log", "events.jsonl")
	if _, err := os.Stat(logPath); err != nil {
		t.Errorf("events.jsonl not created: %v", err)
	}
}

func TestEmit_EventuallyWrites(t *testing.T) {
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
	ev := VerificationEvent("00-054-12", true, "coverage", 82.0)
	Emit(ev)
	// Allow async goroutine to run
	waitForFile(t, filepath.Join(dir, ".sdp", "log", "events.jsonl"), 2*time.Second)
}

func waitForFile(t *testing.T, path string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Errorf("file %s not created within %v", path, timeout)
}
