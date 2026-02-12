package drift

import (
	"testing"
	"time"
)

func TestEnhancedDriftReport(t *testing.T) {
	// AC4: Drift report generated on demand with severity levels
	report := &EnhancedDriftReport{
		WorkstreamID: "00-051-01",
		Timestamp:    time.Now(),
		DriftTypes: []DriftTypeReport{
			{
				Type:     DriftTypeCodeDocs,
				Severity: SeverityError,
				Issues: []EnhancedDriftIssue{
					{File: "test.go", Line: 10, Message: "Missing function", Severity: SeverityError},
				},
			},
			{
				Type:     DriftTypeDecisionCode,
				Severity: SeverityWarning,
				Issues: []EnhancedDriftIssue{
					{File: "ADR-001.md", Message: "Superseded decision", Severity: SeverityWarning},
				},
			},
		},
	}

	if report.WorkstreamID != "00-051-01" {
		t.Errorf("Expected workstream ID 00-051-01, got %s", report.WorkstreamID)
	}

	if len(report.DriftTypes) != 2 {
		t.Errorf("Expected 2 drift types, got %d", len(report.DriftTypes))
	}

	// Generate verdict
	report.Verdict = report.GenerateVerdict()
	if report.Verdict != VerdictFail {
		t.Errorf("Expected verdict FAIL for error severity, got %s", report.Verdict)
	}
}

func TestEnhancedDriftReport_SeverityLevels(t *testing.T) {
	tests := []struct {
		name     string
		report   *EnhancedDriftReport
		expected string
	}{
		{
			name: "error severity",
			report: &EnhancedDriftReport{
				DriftTypes: []DriftTypeReport{
					{Severity: SeverityError},
				},
			},
			expected: VerdictFail,
		},
		{
			name: "warning severity",
			report: &EnhancedDriftReport{
				DriftTypes: []DriftTypeReport{
					{Severity: SeverityWarning},
				},
			},
			expected: VerdictWarning,
		},
		{
			name: "info severity",
			report: &EnhancedDriftReport{
				DriftTypes: []DriftTypeReport{
					{Severity: SeverityInfo},
				},
			},
			expected: VerdictPass,
		},
		{
			name: "no issues",
			report: &EnhancedDriftReport{
				DriftTypes: []DriftTypeReport{},
			},
			expected: VerdictPass,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.report.Verdict = tt.report.GenerateVerdict()
			if tt.report.Verdict != tt.expected {
				t.Errorf("Expected verdict %s, got %s", tt.expected, tt.report.Verdict)
			}
		})
	}
}

func TestEnhancedDriftReport_String(t *testing.T) {
	report := &EnhancedDriftReport{
		WorkstreamID: "00-051-01",
		Timestamp:    time.Now(),
		DriftTypes: []DriftTypeReport{
			{
				Type:     DriftTypeCodeDocs,
				Severity: SeverityWarning,
				Issues: []EnhancedDriftIssue{
					{File: "test.go", Message: "Test issue"},
				},
				Suggestions: []string{"Fix the issue"},
			},
		},
	}

	output := report.String()

	if output == "" {
		t.Error("Expected non-empty string output")
	}
}

func TestEnhancedDriftIssue(t *testing.T) {
	issue := EnhancedDriftIssue{
		File:     "test.go",
		Line:     42,
		Message:  "Test message",
		Severity: SeverityError,
	}

	if issue.File != "test.go" {
		t.Errorf("Expected file test.go, got %s", issue.File)
	}
	if issue.Line != 42 {
		t.Errorf("Expected line 42, got %d", issue.Line)
	}
	if issue.Severity != SeverityError {
		t.Errorf("Expected severity ERROR, got %s", issue.Severity)
	}
}
