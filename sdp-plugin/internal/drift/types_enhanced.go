package drift

import (
	"time"
)

// DriftType represents the type of drift
type DriftType string

const (
	DriftTypeCodeDocs     DriftType = "code_docs"
	DriftTypeDecisionCode DriftType = "decision_code"
	DriftTypeDocsDocs     DriftType = "docs_docs"
)

// Severity represents the severity level
type Severity string

const (
	SeverityError   Severity = "ERROR"
	SeverityWarning Severity = "WARNING"
	SeverityInfo    Severity = "INFO"
)

// Verdict constants
const (
	VerdictPass    = "PASS"
	VerdictWarning = "WARNING"
	VerdictFail    = "FAIL"
)

// EnhancedDriftIssue represents a single drift issue with enhanced details
type EnhancedDriftIssue struct {
	File     string
	Line     int
	Message  string
	Severity Severity
}

// DriftTypeReport represents a report for a specific drift type
type DriftTypeReport struct {
	Type        DriftType
	Severity    Severity
	Issues      []EnhancedDriftIssue
	Suggestions []string
}

// EnhancedDriftReport represents a complete drift detection report
type EnhancedDriftReport struct {
	WorkstreamID string
	Timestamp    time.Time
	DriftTypes   []DriftTypeReport
	Verdict      string
}

// GenerateVerdict determines the overall verdict based on drift types
func (r *EnhancedDriftReport) GenerateVerdict() string {
	for _, dt := range r.DriftTypes {
		if dt.Severity == SeverityError {
			return VerdictFail
		}
	}
	for _, dt := range r.DriftTypes {
		if dt.Severity == SeverityWarning {
			return VerdictWarning
		}
	}
	return VerdictPass
}

// String returns a formatted report
func (r *EnhancedDriftReport) String() string {
	result := "## Enhanced Drift Report\n\n"
	for _, dt := range r.DriftTypes {
		result += "### " + string(dt.Type) + " (" + string(dt.Severity) + ")\n"
		for _, issue := range dt.Issues {
			result += "- " + issue.File
			if issue.Line > 0 {
				result += ":" + string(rune(issue.Line))
			}
			result += ": " + issue.Message + "\n"
		}
	}
	result += "\n**Verdict:** " + r.Verdict + "\n"
	return result
}
