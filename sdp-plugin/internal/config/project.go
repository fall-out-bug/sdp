package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	Guard      GuardSection      `yaml:"guard"`
	Timeouts   TimeoutsSection   `yaml:"timeouts"`
}

// TimeoutsSection holds configurable timeouts (override via SDP_TIMEOUT_* env).
type TimeoutsSection struct {
	VerificationCommand string `yaml:"verification_command"`
	RetryDelay          string `yaml:"retry_delay"`
	BuildPhase          string `yaml:"build_phase"`
	ReviewPhase         string `yaml:"review_phase"`
	CoveragePython      string `yaml:"coverage_python"`
	CoverageGo          string `yaml:"coverage_go"`
	CoverageList        string `yaml:"coverage_list"`
	CoverageJava        string `yaml:"coverage_java"`
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
	CoverageThreshold   int      `yaml:"coverage_threshold"`
	MaxFileLOC          int      `yaml:"max_file_loc"`
	CoverageExclude     []string `yaml:"coverage_exclude"`
	ComplexityThreshold int      `yaml:"complexity_threshold"`
	ComplexityExclude   []string `yaml:"complexity_exclude"`
	SizeExclude         []string `yaml:"size_exclude"`
}

// GuardSection holds guard policy settings (WS-063-03).
type GuardSection struct {
	Mode            string            `yaml:"mode"`
	RulesFile       string            `yaml:"rules_file"`
	SeverityMapping map[string]string `yaml:"severity_mapping"`
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
			CoverageThreshold:   80,
			MaxFileLOC:          200,
			ComplexityThreshold: 40,
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
		Timeouts: TimeoutsSection{
			VerificationCommand: "60s",
			RetryDelay:          "100ms",
			BuildPhase:          "30m",
			ReviewPhase:         "15m",
			CoveragePython:      "30s",
			CoverageGo:          "60s",
			CoverageList:        "10s",
			CoverageJava:        "30s",
		},
	}
}

// Load reads .sdp/config.yml from projectRoot and merges with defaults (AC2, AC3, AC4).
// Validates evidence.log_path is within projectRoot to prevent path traversal.
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
	// Validate log_path is within project root (path traversal safety)
	if cfg.Evidence.LogPath != "" && projectRoot != "" {
		resolvedLog := cfg.Evidence.LogPath
		if !filepath.IsAbs(resolvedLog) {
			resolvedLog = filepath.Join(projectRoot, resolvedLog)
		}
		if err := validatePathWithinRoot(projectRoot, resolvedLog); err != nil {
			return nil, fmt.Errorf("evidence.log_path: %w", err)
		}
	}
	return cfg, nil
}

// validatePathWithinRoot returns nil if path is within root (no traversal outside root).
func validatePathWithinRoot(root, path string) error {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("resolve root: %w", err)
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}
	rel, err := filepath.Rel(absRoot, absPath)
	if err != nil {
		return fmt.Errorf("relative path: %w", err)
	}
	if strings.HasPrefix(rel, "..") {
		return fmt.Errorf("path outside project root")
	}
	return nil
}

// Validate returns an error if config is invalid.
func (c *Config) Validate() error {
	if c.Version < 1 {
		return fmt.Errorf("version must be >= 1, got %d", c.Version)
	}
	timeoutFields := []struct {
		name, val string
	}{
		{"timeouts.verification_command", c.Timeouts.VerificationCommand},
		{"timeouts.retry_delay", c.Timeouts.RetryDelay},
		{"timeouts.build_phase", c.Timeouts.BuildPhase},
		{"timeouts.review_phase", c.Timeouts.ReviewPhase},
		{"timeouts.coverage_python", c.Timeouts.CoveragePython},
		{"timeouts.coverage_go", c.Timeouts.CoverageGo},
		{"timeouts.coverage_list", c.Timeouts.CoverageList},
		{"timeouts.coverage_java", c.Timeouts.CoverageJava},
	}
	for _, f := range timeoutFields {
		if f.val != "" {
			if _, err := time.ParseDuration(f.val); err != nil {
				return fmt.Errorf("%s: invalid duration %q: %w", f.name, f.val, err)
			}
		}
	}
	if c.Acceptance.Timeout != "" {
		if _, err := time.ParseDuration(c.Acceptance.Timeout); err != nil {
			return fmt.Errorf("acceptance.timeout: invalid duration %q: %w", c.Acceptance.Timeout, err)
		}
	}
	return nil
}

// TimeoutFromEnv returns duration from env key (e.g. SDP_TIMEOUT_VERIFICATION) or fallback.
func TimeoutFromEnv(envKey string, fallback time.Duration) time.Duration {
	if v := os.Getenv(envKey); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

// TimeoutFromConfigOrEnv returns duration from config string, then env, then fallback.
func TimeoutFromConfigOrEnv(configVal, envKey string, fallback time.Duration) time.Duration {
	if configVal != "" {
		if d, err := time.ParseDuration(configVal); err == nil {
			return d
		}
	}
	return TimeoutFromEnv(envKey, fallback)
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
