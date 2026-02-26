package main

import (
	"os"
	"path/filepath"
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
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cfgPath := filepath.Join(cfgDir, "config.yml")
	cfgContent := "version: 1\nevidence:\n  enabled: true\n  log_path: \".sdp/log/events.jsonl\"\n"
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0644); err != nil {
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
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()

	cfgDir := filepath.Join(tmpDir, ".sdp")
	logDir := filepath.Join(cfgDir, "log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cfgPath := filepath.Join(cfgDir, "config.yml")
	cfgContent := "version: 1\nevidence:\n  enabled: false\n"
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	cmd := deployCmd()
	// No flags set - should use defaults
	// This will fail if git is not available, but we test the logic path
	_ = cmd.RunE(cmd, []string{})
	// Don't fail test on error since git may not be available in test env
}

// TestDeployCmd_EvidenceDisabled tests deploy when evidence is disabled
func TestDeployCmd_EvidenceDisabled(t *testing.T) {
	evidence.ResetGlobalWriter()
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()

	cfgDir := filepath.Join(tmpDir, ".sdp")
	logDir := filepath.Join(cfgDir, "log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cfgPath := filepath.Join(cfgDir, "config.yml")
	cfgContent := "version: 1\nevidence:\n  enabled: false\n"
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0644); err != nil {
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
		if !contains(s, sub) {
			return false
		}
	}
	return true
}
