package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const guardRulesFile = "guard-rules.yml"

// GuardRules holds guard rule definitions.
type GuardRules struct {
	Version int         `yaml:"version"`
	Rules   []GuardRule `yaml:"rules"`
}

// GuardRule represents a single guard rule.
type GuardRule struct {
	ID       string                 `yaml:"id"`
	Enabled  bool                   `yaml:"enabled"`
	Severity string                 `yaml:"severity"`
	Config   map[string]interface{} `yaml:"config"`
}

// DefaultGuardRules returns default guard rules when no rules file exists.
func DefaultGuardRules() *GuardRules {
	return &GuardRules{
		Version: 1,
		Rules: []GuardRule{
			{
				ID:       "max-file-loc",
				Enabled:  true,
				Severity: "error",
				Config:   map[string]interface{}{"max_lines": 200},
			},
			{
				ID:       "coverage-threshold",
				Enabled:  true,
				Severity: "error",
				Config:   map[string]interface{}{"minimum": 80},
			},
		},
	}
}

// LoadGuardRules reads .sdp/guard-rules.yml from projectRoot and returns rules (AC1, AC2).
func LoadGuardRules(projectRoot string) (*GuardRules, error) {
	path := filepath.Join(projectRoot, configDir, guardRulesFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultGuardRules(), nil
		}
		return nil, fmt.Errorf("read guard rules: %w", err)
	}

	var rules GuardRules
	if err := yaml.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("parse guard rules: %w", err)
	}

	// Validate rules (AC2)
	if err := validateGuardRules(&rules); err != nil {
		return nil, err
	}

	return &rules, nil
}

// validateGuardRules validates guard rules structure and values (AC2).
func validateGuardRules(rules *GuardRules) error {
	if rules.Version < 1 {
		return fmt.Errorf("invalid version: must be >= 1, got %d", rules.Version)
	}

	validSeverities := map[string]bool{
		"error":   true,
		"warning": true,
		"info":    true,
	}

	for i, rule := range rules.Rules {
		if rule.ID == "" {
			return fmt.Errorf("rule at index %d: missing required field 'id'", i)
		}
		if rule.Severity == "" {
			return fmt.Errorf("rule %s: missing required field 'severity'", rule.ID)
		}
		if !validSeverities[rule.Severity] {
			return fmt.Errorf("rule %s: invalid severity %q, must be one of: error, warning, info",
				rule.ID, rule.Severity)
		}
	}

	return nil
}
