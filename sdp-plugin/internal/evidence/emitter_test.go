package evidence

import (
	"errors"
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
	ResetGlobalWriter()
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
	ResetGlobalWriter()
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
	if err := EmitSync(ev); err != nil {
		t.Fatalf("EmitSync: %v", err)
	}
	logPath := filepath.Join(dir, ".sdp", "log", "events.jsonl")
	// Retry Stat to reduce flakiness on slow filesystems (sdp-yout, CI)
	for i := 0; i < 25; i++ {
		if _, err := os.Stat(logPath); err == nil {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Errorf("events.jsonl not created after 25 retries (625ms)")
}

func TestValidateEvent(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		ev := PlanEvent("00-054-01", nil)
		fillDefaults(ev)
		if err := ValidateEvent(ev); err != nil {
			t.Errorf("ValidateEvent(valid): %v", err)
		}
	})
	t.Run("nil", func(t *testing.T) {
		if err := ValidateEvent(nil); err == nil {
			t.Error("ValidateEvent(nil): want error")
		} else if !errors.Is(err, ErrEventInvalid) {
			t.Errorf("ValidateEvent(nil): want ErrEventInvalid, got %v", err)
		}
	})
	t.Run("missing type", func(t *testing.T) {
		ev := &Event{ID: "x", Timestamp: "2026-01-01T00:00:00Z", Type: ""}
		if err := ValidateEvent(ev); err == nil {
			t.Error("ValidateEvent(missing type): want error")
		}
	})
}

func TestEmit_ReturnsValidationError(t *testing.T) {
	if err := Emit(nil); err != nil {
		t.Errorf("Emit(nil): want nil, got %v", err)
	}
	ev := &Event{Type: "plan"} // missing ID, Timestamp before fillDefaults - but fillDefaults adds them
	fillDefaults(ev)
	ev.Type = "" // now invalid
	if err := Emit(ev); err == nil {
		t.Error("Emit(invalid): want error")
	}
}
