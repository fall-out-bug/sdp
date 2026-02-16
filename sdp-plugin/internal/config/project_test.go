package config

import (
	"os"
	"path/filepath"
	"strings"
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

// AC1: Test loading guard rules from .sdp/guard-rules.yml
func TestLoadGuardRules(t *testing.T) {
	dir := t.TempDir()
	sdpDir := filepath.Join(dir, configDir)
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create guard rules file
	rulesPath := filepath.Join(sdpDir, "guard-rules.yml")
	rulesContent := `version: 1
rules:
  - id: "max-file-loc"
    enabled: true
    severity: "error"
    config:
      max_lines: 200
  - id: "coverage-threshold"
    enabled: true
    severity: "error"
    config:
      minimum: 80
`
	if err := os.WriteFile(rulesPath, []byte(rulesContent), 0644); err != nil {
		t.Fatalf("write guard rules: %v", err)
	}

	rules, err := LoadGuardRules(dir)
	if err != nil {
		t.Fatalf("LoadGuardRules: %v", err)
	}
	if rules == nil {
		t.Fatal("LoadGuardRules returned nil rules")
	}
	if rules.Version != 1 {
		t.Errorf("Version: want 1, got %d", rules.Version)
	}
	if len(rules.Rules) != 2 {
		t.Errorf("Rules count: want 2, got %d", len(rules.Rules))
	}
}

// AC1: Test missing guard rules returns defaults
func TestLoadGuardRulesMissingReturnsDefaults(t *testing.T) {
	dir := t.TempDir()
	sdpDir := filepath.Join(dir, configDir)
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	rules, err := LoadGuardRules(dir)
	if err != nil {
		t.Fatalf("LoadGuardRules: %v", err)
	}
	if rules == nil {
		t.Fatal("LoadGuardRules returned nil rules")
	}
	// Should have default rules
	if len(rules.Rules) == 0 {
		t.Error("Default rules should not be empty")
	}
}

// AC2: Test guard rules validation provides explicit errors
func TestGuardRulesValidation(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid rules",
			content: `version: 1
rules:
  - id: "max-file-loc"
    enabled: true
    severity: "error"
    config:
      max_lines: 200
`,
			wantErr: false,
		},
		{
			name: "missing version",
			content: `rules:
  - id: "max-file-loc"
    enabled: true
    severity: "error"
`,
			wantErr:     true,
			errContains: "version",
		},
		{
			name: "invalid severity",
			content: `version: 1
rules:
  - id: "max-file-loc"
    enabled: true
    severity: "invalid"
`,
			wantErr:     true,
			errContains: "severity",
		},
		{
			name: "missing rule id",
			content: `version: 1
rules:
  - enabled: true
    severity: "error"
`,
			wantErr:     true,
			errContains: "missing required field 'id'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			sdpDir := filepath.Join(dir, configDir)
			if err := os.MkdirAll(sdpDir, 0755); err != nil {
				t.Fatalf("mkdir: %v", err)
			}

			rulesPath := filepath.Join(sdpDir, "guard-rules.yml")
			if err := os.WriteFile(rulesPath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("write guard rules: %v", err)
			}

			rules, err := LoadGuardRules(dir)
			if tt.wantErr {
				if err == nil {
					t.Error("LoadGuardRules should return error")
				}
				if tt.errContains != "" && err != nil {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("Error should contain %q, got %q", tt.errContains, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("LoadGuardRules should not return error: %v", err)
				}
				if rules == nil {
					t.Error("LoadGuardRules returned nil rules")
				}
			}
		})
	}
}

// AC3: Test config.yml guard policy settings
func TestConfigGuardPolicySettings(t *testing.T) {
	dir := t.TempDir()
	sdpDir := filepath.Join(dir, configDir)
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	configContent := `version: 1
guard:
  mode: "strict"
  severity_mapping:
    error: "block"
    warning: "warn"
    info: "log"
  rules_file: ".sdp/guard-rules.yml"
`
	configPath := filepath.Join(sdpDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	// Check guard policy settings
	if cfg.Guard.Mode != "strict" {
		t.Errorf("Guard.Mode: want strict, got %s", cfg.Guard.Mode)
	}
	if cfg.Guard.RulesFile != ".sdp/guard-rules.yml" {
		t.Errorf("Guard.RulesFile: want .sdp/guard-rules.yml, got %s", cfg.Guard.RulesFile)
	}
	if len(cfg.Guard.SeverityMapping) == 0 {
		t.Error("Guard.SeverityMapping should not be empty")
	}
}

// AC3: Test default guard policy settings
func TestDefaultConfigGuardPolicy(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Guard.Mode == "" {
		t.Error("Guard.Mode should have default value")
	}
	if cfg.Guard.RulesFile == "" {
		t.Error("Guard.RulesFile should have default value")
	}
	if len(cfg.Guard.SeverityMapping) == 0 {
		t.Error("Guard.SeverityMapping should have default values")
	}
}
