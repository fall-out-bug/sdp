package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const configDir = ".sdp"
const configFile = "config.yml"

// Config holds project-level SDP settings.
type Config struct {
	Version    int               `yaml:"version"`
	Acceptance AcceptanceSection `yaml:"acceptance"`
	Evidence   EvidenceSection   `yaml:"evidence"`
	Quality    QualitySection    `yaml:"quality"`
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
