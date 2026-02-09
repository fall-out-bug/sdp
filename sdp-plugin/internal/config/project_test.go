package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}
	if cfg.Version != 1 {
		t.Errorf("Version: want 1, got %d", cfg.Version)
	}
	if cfg.Acceptance.Command == "" {
		t.Error("Acceptance.Command should be set")
	}
	if cfg.Acceptance.Timeout == "" {
		t.Error("Acceptance.Timeout should be set")
	}
	if cfg.Acceptance.Expected == "" {
		t.Error("Acceptance.Expected should be set")
	}
	if cfg.Quality.CoverageThreshold != 80 {
		t.Errorf("Quality.CoverageThreshold: want 80, got %d", cfg.Quality.CoverageThreshold)
	}
	if cfg.Quality.MaxFileLOC != 200 {
		t.Errorf("Quality.MaxFileLOC: want 200, got %d", cfg.Quality.MaxFileLOC)
	}
}

func TestLoadMissingFileReturnsDefaults(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg == nil {
		t.Fatal("Load returned nil config")
	}
	if cfg.Version != 1 {
		t.Errorf("Version: want 1, got %d", cfg.Version)
	}
}

func TestLoadValidYAML(t *testing.T) {
	dir := t.TempDir()
	sdpDir := filepath.Join(dir, configDir)
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path := filepath.Join(sdpDir, configFile)
	content := "version: 1\nacceptance:\n  command: \"custom test\"\n  timeout: 60s\n  expected: OK\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Acceptance.Command != "custom test" {
		t.Errorf("Acceptance.Command: want custom test, got %s", cfg.Acceptance.Command)
	}
	if cfg.Acceptance.Timeout != "60s" {
		t.Errorf("Acceptance.Timeout: want 60s, got %s", cfg.Acceptance.Timeout)
	}
	if cfg.Acceptance.Expected != "OK" {
		t.Errorf("Acceptance.Expected: want OK, got %s", cfg.Acceptance.Expected)
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	sdpDir := filepath.Join(dir, configDir)
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path := filepath.Join(sdpDir, configFile)
	if err := os.WriteFile(path, []byte("invalid: yaml: ["), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	_, err := Load(dir)
	if err == nil {
		t.Error("Load should fail on invalid YAML")
	}
}

func TestValidate(t *testing.T) {
	cfg := DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("DefaultConfig should validate: %v", err)
	}
	cfg.Version = 0
	if err := cfg.Validate(); err == nil {
		t.Error("Version 0 should fail validation")
	}
}

func TestFindProjectRoot(t *testing.T) {
	dir := t.TempDir()
	sdpDir := filepath.Join(dir, configDir)
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	subDir := filepath.Join(dir, "sub", "deep")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("mkdir sub: %v", err)
	}
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	if err := os.Chdir(subDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	root, err := FindProjectRoot()
	if err != nil {
		t.Fatalf("FindProjectRoot: %v", err)
	}
	normDir, _ := filepath.EvalSymlinks(dir)
	normRoot, _ := filepath.EvalSymlinks(root)
	if normRoot != normDir {
		t.Errorf("FindProjectRoot: want %s, got %s", normDir, normRoot)
	}
}
