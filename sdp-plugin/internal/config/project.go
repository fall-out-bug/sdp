package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const configDir = ".sdp"
const configFile = "config.yml"
const guardRulesFile = "guard-rules.yml"

// Config holds project-level SDP settings.
type Config struct {
	Version    int               `yaml:"version"`
	Acceptance AcceptanceSection `yaml:"acceptance"`
	Evidence   EvidenceSection   `yaml:"evidence"`
	Quality    QualitySection    `yaml:"quality"`
	Guard      GuardSection      `yaml:"guard"`
}

// AcceptanceSection holds acceptance test gate settings.
type AcceptanceSection struct {
	Command  string `yaml:"command"`
	Timeout  string `yaml:"timeout"`
	Expected string `yaml:"expected"`
}

// EvidenceSection holds evidence log settings (stub for WS-04).
type EvidenceSection struct {
	Enabled bool   `yaml:"enabled"`
	LogPath string `yaml:"log_path"`
}

// QualitySection holds quality gate settings (stub for WS-06).
type QualitySection struct {
	CoverageThreshold int `yaml:"coverage_threshold"`
	MaxFileLOC        int `yaml:"max_file_loc"`
}

// GuardSection holds guard policy settings (WS-063-03).
type GuardSection struct {
	Mode            string            `yaml:"mode"`
	RulesFile       string            `yaml:"rules_file"`
	SeverityMapping map[string]string `yaml:"severity_mapping"`
}

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

// DefaultConfig returns config with sensible defaults (AC4).
func DefaultConfig() *Config {
	return &Config{
		Version: 1,
		Acceptance: AcceptanceSection{
			Command:  "go test ./... -run TestSmoke",
			Timeout:  "30s",
			Expected: "PASS",
		},
		Evidence: EvidenceSection{
			Enabled: true,
			LogPath: ".sdp/log/events.jsonl",
		},
		Quality: QualitySection{
			CoverageThreshold: 80,
			MaxFileLOC:        200,
		},
		Guard: GuardSection{
			Mode:      "standard",
			RulesFile: ".sdp/guard-rules.yml",
			SeverityMapping: map[string]string{
				"error":   "block",
				"warning": "warn",
				"info":    "log",
			},
		},
	}
}

// Load reads .sdp/config.yml from projectRoot and merges with defaults (AC2, AC3, AC4).
func Load(projectRoot string) (*Config, error) {
	path := filepath.Join(projectRoot, configDir, configFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}
	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

// Validate returns an error if config is invalid.
func (c *Config) Validate() error {
	if c.Version < 1 {
		return fmt.Errorf("version must be >= 1, got %d", c.Version)
	}
	return nil
}

// FindProjectRoot walks up from cwd to find a directory containing .sdp or .git.
func FindProjectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	current := cwd
	for {
		if _, err := os.Stat(filepath.Join(current, configDir)); err == nil {
			return current, nil
		}
		if _, err := os.Stat(filepath.Join(current, ".git")); err == nil {
			return current, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			return cwd, nil
		}
		current = parent
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
