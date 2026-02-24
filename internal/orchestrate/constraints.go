package orchestrate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Constraint defines a single rule for agent behavior in a phase.
type Constraint struct {
	ID          string `yaml:"id"`
	Description string `yaml:"description"`
	Severity    string `yaml:"severity"` // warn, block, halt, escalate
	Check       string `yaml:"check"`    // scope-diff, command-pattern, file-pattern, file-exists
	Pattern     string `yaml:"pattern,omitempty"`
	Path        string `yaml:"path,omitempty"`
	Message     string `yaml:"message"`
}

// PhaseConstraints holds constraints for a specific phase.
type PhaseConstraints struct {
	Description string       `yaml:"description"`
	Constraints []Constraint `yaml:"constraints"`
}

// Containment thresholds.
type ContainmentThresholds struct {
	Warn      int `yaml:"warn"`
	Block     int `yaml:"block"`
	Halt      int `yaml:"halt"`
	Escalate  int `yaml:"escalate"`
}

// AgentConstraintConfig is the full config from .sdp/agent-constraints.yaml.
type AgentConstraintConfig struct {
	Version  string                      `yaml:"version"`
	Updated  string                      `yaml:"updated"`
	Phases   map[string]PhaseConstraints `yaml:"phases"`
	Containment struct {
		Thresholds ContainmentThresholds `yaml:"thresholds"`
	} `yaml:"containment"`
}

// ConstraintViolation records a rule that was triggered.
type ConstraintViolation struct {
	ConstraintID string
	Severity     string
	Message      string
}

// LoadConstraintConfig reads .sdp/agent-constraints.yaml.
// Returns empty config if file doesn't exist.
func LoadConstraintConfig(projectRoot string) (*AgentConstraintConfig, error) {
	path := filepath.Join(projectRoot, ".sdp", "agent-constraints.yaml")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &AgentConstraintConfig{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read constraints: %w", err)
	}
	var cfg AgentConstraintConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse constraints: %w", err)
	}
	return &cfg, nil
}

// CheckCommand evaluates agent-constraints for a shell command about to be executed.
// Returns violations (if any). Caller decides whether to block/halt.
func CheckCommand(cfg *AgentConstraintConfig, phase, command string) []ConstraintViolation {
	if cfg == nil {
		return nil
	}
	pc, ok := cfg.Phases[phase]
	if !ok {
		return nil
	}

	var violations []ConstraintViolation
	for _, c := range pc.Constraints {
		if c.Check != "command-pattern" {
			continue
		}
		if matchesPattern(command, c.Pattern) {
			violations = append(violations, ConstraintViolation{
				ConstraintID: c.ID,
				Severity:     c.Severity,
				Message:      c.Message,
			})
		}
	}
	return violations
}

// CheckFileAccess evaluates agent-constraints for a file about to be read or written.
func CheckFileAccess(cfg *AgentConstraintConfig, phase, filePath string) []ConstraintViolation {
	if cfg == nil {
		return nil
	}
	pc, ok := cfg.Phases[phase]
	if !ok {
		return nil
	}

	var violations []ConstraintViolation
	for _, c := range pc.Constraints {
		if c.Check != "file-pattern" {
			continue
		}
		if matchesPattern(filePath, c.Pattern) {
			violations = append(violations, ConstraintViolation{
				ConstraintID: c.ID,
				Severity:     c.Severity,
				Message:      c.Message,
			})
		}
	}
	return violations
}

// CheckRequiredFiles evaluates file-exists constraints.
func CheckRequiredFiles(cfg *AgentConstraintConfig, phase, projectRoot, featureID string) []ConstraintViolation {
	if cfg == nil {
		return nil
	}
	pc, ok := cfg.Phases[phase]
	if !ok {
		return nil
	}

	var violations []ConstraintViolation
	for _, c := range pc.Constraints {
		if c.Check != "file-exists" {
			continue
		}
		path := strings.ReplaceAll(c.Path, "{feature_id}", featureID)
		fullPath := filepath.Join(projectRoot, path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			violations = append(violations, ConstraintViolation{
				ConstraintID: c.ID,
				Severity:     c.Severity,
				Message:      c.Message,
			})
		}
	}
	return violations
}

// DetermineContainmentLevel returns the effective severity for a given violation count.
func DetermineContainmentLevel(cfg *AgentConstraintConfig, violationCount int) string {
	if cfg == nil {
		return "warn"
	}
	t := cfg.Containment.Thresholds
	switch {
	case violationCount >= t.Escalate && t.Escalate > 0:
		return "escalate"
	case violationCount >= t.Halt && t.Halt > 0:
		return "halt"
	case violationCount >= t.Block && t.Block > 0:
		return "block"
	default:
		return "warn"
	}
}

func matchesPattern(s, pattern string) bool {
	if pattern == "" {
		return false
	}
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return strings.Contains(s, pattern)
	}
	return matched
}
