package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall-out-bug/sdp/internal/evidence"
)

// TestDeployCmd_ApprovalEvent tests that deploy command emits approval event (F056-01 AC3, AC4)
func TestDeployCmd_ApprovalEvent(t *testing.T) {
	evidence.ResetGlobalWriter()
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()

	// Create .sdp/config.yml to enable evidence
	cfgDir := filepath.Join(tmpDir, ".sdp")
	logDir := filepath.Join(cfgDir, "log")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cfgPath := filepath.Join(cfgDir, "config.yml")
	cfgContent := "version: 1\nevidence:\n  enabled: true\n  log_path: \".sdp/log/events.jsonl\"\n"
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Test deploy with explicit flags
	cmd := deployCmd()
	if err := cmd.Flags().Set("target", "main"); err != nil {
		t.Fatalf("set target: %v", err)
	}
	if err := cmd.Flags().Set("sha", "abc123def456"); err != nil {
		t.Fatalf("set sha: %v", err)
	}
	if err := cmd.Flags().Set("who", "CI"); err != nil {
		t.Fatalf("set who: %v", err)
	}

	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("deployCmd() failed: %v", err)
	}

	// Verify event was written
	logPath := filepath.Join(tmpDir, ".sdp", "log", "events.jsonl")
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read events.jsonl: %v", err)
	}
	content := string(data)
	if content == "" {
		t.Error("events.jsonl is empty")
	}

	// Verify approval event contains expected fields (AC4)
	if !containsAll(content, `"type":"approval"`, `"target_branch":"main"`, `"commit_sha":"abc123def456"`, `"approved_by":"CI"`) {
		t.Errorf("approval event missing required fields: %s", content)
	}
}

// TestDeployCmd_DefaultValues tests deploy with default values
func TestDeployCmd_DefaultValues(t *testing.T) {
	evidence.ResetGlobalWriter()
	originalResolveSHA := deployResolveSHA
	originalResolveApprover := deployResolveApprover
	deployResolveSHA = func() (string, error) {
		return "deadbeefcafebabe", nil
	}
	deployResolveApprover = func() (string, error) {
		return "Test Runner", nil
	}
	t.Cleanup(func() {
		deployResolveSHA = originalResolveSHA
		deployResolveApprover = originalResolveApprover
	})

	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()

	cfgDir := filepath.Join(tmpDir, ".sdp")
	logDir := filepath.Join(cfgDir, "log")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cfgPath := filepath.Join(cfgDir, "config.yml")
	cfgContent := "version: 1\nevidence:\n  enabled: false\n"
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	cmd := deployCmd()
	output, err := captureDeployOutput(t, func() error {
		return cmd.RunE(cmd, []string{})
	})
	if err != nil {
		t.Fatalf("deployCmd() with defaults failed: %v", err)
	}
	if !strings.Contains(output, "Approval recorded: deadbee -> main (Test Runner)") {
		t.Fatalf("deploy output missing resolved defaults:\n%s", output)
	}
}

// TestDeployCmd_EvidenceDisabled tests deploy when evidence is disabled
func TestDeployCmd_EvidenceDisabled(t *testing.T) {
	evidence.ResetGlobalWriter()
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()

	cfgDir := filepath.Join(tmpDir, ".sdp")
	logDir := filepath.Join(cfgDir, "log")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cfgPath := filepath.Join(cfgDir, "config.yml")
	cfgContent := "version: 1\nevidence:\n  enabled: false\n"
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	cmd := deployCmd()
	if err := cmd.Flags().Set("sha", "test123"); err != nil {
		t.Fatalf("set sha: %v", err)
	}

	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("deployCmd() with evidence disabled failed: %v", err)
	}

	// Verify no event was written
	logPath := filepath.Join(tmpDir, ".sdp", "log", "events.jsonl")
	if _, err := os.Stat(logPath); !os.IsNotExist(err) {
		data, _ := os.ReadFile(logPath)
		t.Errorf("events.jsonl should not exist or be empty when evidence disabled: %s", string(data))
	}
}

// containsAll checks if s contains all substrings
func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}

func captureDeployOutput(t *testing.T, run func() error) (string, error) {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w

	runErr := run()

	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("read output: %v", err)
	}
	return buf.String(), runErr
}
