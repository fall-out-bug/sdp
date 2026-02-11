package parser

import (
	"fmt"
	"regexp"
)

// Workstream represents a parsed workstream markdown file
type Workstream struct {
	ID         string
	Parent     string
	Feature    string
	Status     string
	Size       string
	ProjectID  string
	Goal       string
	Acceptance []string
	Scope      Scope
}

// Scope represents the scope section of a workstream
type Scope struct {
	Implementation []string
	Tests          []string
}

// ValidationIssue represents a validation error or warning
type ValidationIssue struct {
	Field    string
	Message  string
	Severity string // ERROR, WARNING
}

// wsIDRegex validates the PP-FFF-SS format
var wsIDRegex = regexp.MustCompile(`^\d{2}-\d{3}-\d{2}$`)

// Validate validates the workstream ID format
func (ws *Workstream) Validate() error {
	if ws.ID == "" {
		return fmt.Errorf("workstream ID is required")
	}

	if !wsIDRegex.MatchString(ws.ID) {
		return fmt.Errorf("invalid workstream ID format: %s (expected PP-FFF-SS)", ws.ID)
	}

	return nil
}

// frontmatter represents the YAML frontmatter structure
type frontmatter struct {
	WSID      string `yaml:"ws_id"`
	Parent    string `yaml:"parent"`
	Feature   string `yaml:"feature"`              // Legacy field name (backward compat)
	FeatureID string `yaml:"feature_id,omitempty"` // Preferred field name
	Status    string `yaml:"status"`
	Size      string `yaml:"size"`
	ProjectID string `yaml:"project_id"`
}

// getFeature returns the feature ID, preferring feature_id over feature
func (fm *frontmatter) getFeature() string {
	if fm.FeatureID != "" {
		return fm.FeatureID
	}
	return fm.Feature
}
