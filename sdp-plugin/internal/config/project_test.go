package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_LogPathValidationWhenProjectRootEmpty(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".sdp")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(cfgDir, "config.yml")
	// Malicious log_path that would escape when resolved from cwd
	cfgContent := "version: 1\nevidence:\n  enabled: true\n  log_path: \"../../etc/passwd\"\n"
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0644); err != nil {
		t.Fatal(err)
	}
	origWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	// Load with empty projectRoot — should still validate log_path (using cwd)
	_, err := Load("")
	if err == nil {
		t.Error("expected error for log_path outside root when projectRoot empty, got nil")
	}
}
