package evidence

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fall-out-bug/sdp/internal/verify"
)

func TestExtractLesson_NilResult(t *testing.T) {
	l := ExtractLesson("00-054-01", nil)
	if l.WSID != "00-054-01" {
		t.Errorf("WSID: want 00-054-01, got %s", l.WSID)
	}
	if l.Outcome != "unknown" {
		t.Errorf("Outcome: want unknown, got %s", l.Outcome)
	}
}

func TestExtractLesson_Passed(t *testing.T) {
	result := &verify.VerificationResult{
		WSID:   "00-054-01",
		Passed: true,
		Checks: []verify.CheckResult{
			{Name: "File", Passed: true, Message: "schema/index.json"},
		},
	}
	l := ExtractLesson("00-054-01", result)
	if l.Outcome != "passed" {
		t.Errorf("Outcome: want passed, got %s", l.Outcome)
	}
	if len(l.WhatWorked) != 1 {
		t.Errorf("WhatWorked: want 1, got %d", len(l.WhatWorked))
	}
}

func TestExtractLesson_Failed(t *testing.T) {
	result := &verify.VerificationResult{
		WSID:   "00-054-01",
		Passed: false,
		Checks: []verify.CheckResult{
			{Name: "File", Passed: false, Message: "Missing: x.go"},
		},
		MissingFiles: []string{"x.go"},
		Duration:     time.Second,
	}
	l := ExtractLesson("00-054-01", result)
	if l.Outcome != "failed" {
		t.Errorf("Outcome: want failed, got %s", l.Outcome)
	}
	if len(l.WhatFailed) == 0 {
		t.Error("WhatFailed should not be empty")
	}
	if l.Category != "verification" {
		t.Errorf("Category: want verification, got %s", l.Category)
	}
}

func TestLesson_MatchesOutcome(t *testing.T) {
	l := Lesson{Outcome: "failed"}
	if !l.MatchesOutcome("failed") {
		t.Error("MatchesOutcome(failed) should be true")
	}
	if !l.MatchesOutcome("") {
		t.Error("MatchesOutcome('') should be true")
	}
	if l.MatchesOutcome("passed") {
		t.Error("MatchesOutcome(passed) should be false")
	}
}

func TestEmitLesson_WritesWhenEnabled(t *testing.T) {
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
	lesson := Lesson{
		WSID:       "00-054-13",
		Outcome:    "passed",
		WhatWorked: []string{"Check: ok"},
		Category:   "verification",
	}
	EmitLesson(lesson)
	logPath := filepath.Join(dir, ".sdp", "log", "events.jsonl")
	r := NewReader(logPath)
	events, err := r.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	var found bool
	for _, e := range events {
		if e.Type == "lesson" && e.WSID == "00-054-13" {
			found = true
			break
		}
	}
	if !found {
		t.Error("EmitLesson: no lesson event found in evidence log")
	}
}
