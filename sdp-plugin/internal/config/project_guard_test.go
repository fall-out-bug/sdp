package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultGuardRules(t *testing.T) {
	rules := DefaultGuardRules()
	if rules == nil {
		t.Fatal("DefaultGuardRules returned nil")
	}
	if rules.Version != 1 {
		t.Errorf("expected version 1, got %d", rules.Version)
	}
	if len(rules.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rules.Rules))
	}
	if rules.Rules[0].ID != "max-file-loc" {
		t.Errorf("expected first rule id max-file-loc, got %s", rules.Rules[0].ID)
	}
	if rules.Rules[1].ID != "coverage-threshold" {
		t.Errorf("expected second rule id coverage-threshold, got %s", rules.Rules[1].ID)
	}
}

func TestLoadGuardRules_FileNotExist_ReturnsDefault(t *testing.T) {
	dir := t.TempDir()
	rules, err := LoadGuardRules(dir)
	if err != nil {
		t.Fatalf("LoadGuardRules failed: %v", err)
	}
	if rules == nil {
		t.Fatal("expected non-nil rules")
	}
	if rules.Version != 1 {
		t.Errorf("expected default version 1, got %d", rules.Version)
	}
}

func TestLoadGuardRules_ValidFile(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".sdp")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatal(err)
	}
	rulesPath := filepath.Join(cfgDir, "guard-rules.yml")
	content := `version: 1
rules:
  - id: custom-rule
    enabled: true
    severity: warning
`
	if err := os.WriteFile(rulesPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	rules, err := LoadGuardRules(dir)
	if err != nil {
		t.Fatalf("LoadGuardRules failed: %v", err)
	}
	if len(rules.Rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(rules.Rules))
	}
	if rules.Rules[0].ID != "custom-rule" {
		t.Errorf("expected rule id custom-rule, got %s", rules.Rules[0].ID)
	}
	if rules.Rules[0].Severity != "warning" {
		t.Errorf("expected severity warning, got %s", rules.Rules[0].Severity)
	}
}

func TestLoadGuardRules_InvalidYAML_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".sdp")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatal(err)
	}
	rulesPath := filepath.Join(cfgDir, "guard-rules.yml")
	if err := os.WriteFile(rulesPath, []byte("invalid: ["), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadGuardRules(dir)
	if err == nil {
		t.Error("expected parse error for invalid YAML")
	}
}

func TestLoadGuardRules_InvalidVersion_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".sdp")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatal(err)
	}
	rulesPath := filepath.Join(cfgDir, "guard-rules.yml")
	content := `version: 0
rules: []
`
	if err := os.WriteFile(rulesPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadGuardRules(dir)
	if err == nil {
		t.Error("expected error for version < 1")
	}
}

func TestLoadGuardRules_MissingRuleID_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".sdp")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatal(err)
	}
	rulesPath := filepath.Join(cfgDir, "guard-rules.yml")
	content := `version: 1
rules:
  - enabled: true
    severity: error
`
	if err := os.WriteFile(rulesPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadGuardRules(dir)
	if err == nil {
		t.Error("expected error for missing rule id")
	}
}

func TestLoadGuardRules_MissingSeverity_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".sdp")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatal(err)
	}
	rulesPath := filepath.Join(cfgDir, "guard-rules.yml")
	content := `version: 1
rules:
  - id: my-rule
    enabled: true
`
	if err := os.WriteFile(rulesPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadGuardRules(dir)
	if err == nil {
		t.Error("expected error for missing severity")
	}
}

func TestLoadGuardRules_InvalidSeverity_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".sdp")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatal(err)
	}
	rulesPath := filepath.Join(cfgDir, "guard-rules.yml")
	content := `version: 1
rules:
  - id: my-rule
    enabled: true
    severity: critical
`
	if err := os.WriteFile(rulesPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadGuardRules(dir)
	if err == nil {
		t.Error("expected error for invalid severity")
	}
}
