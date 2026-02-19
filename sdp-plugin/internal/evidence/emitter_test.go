package evidence

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestModelID(t *testing.T) {
	os.Unsetenv("SDP_MODEL_ID")
	os.Unsetenv("OPENCODE_MODEL")
	os.Unsetenv("ANTHROPIC_MODEL")
	os.Unsetenv("OPENAI_MODEL")
	os.Unsetenv("MODEL")
	if got := ModelID(); got != "unknown" {
		t.Errorf("ModelID(): want unknown, got %s", got)
	}
	os.Setenv("OPENCODE_MODEL", "openai/gpt-5.3-codex")
	defer os.Unsetenv("OPENCODE_MODEL")
	if got := ModelID(); got != "openai/gpt-5.3-codex" {
		t.Errorf("ModelID(): want openai/gpt-5.3-codex, got %s", got)
	}
	os.Setenv("SDP_MODEL_ID", "claude-sonnet")
	defer os.Unsetenv("SDP_MODEL_ID")
	if got := ModelID(); got != "claude-sonnet" {
		t.Errorf("ModelID(): want claude-sonnet, got %s", got)
	}
}

func TestEnabled(t *testing.T) {
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
