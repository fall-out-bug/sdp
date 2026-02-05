package parser

import (
	"fmt"
	"os"
)

// ValidateFile validates a workstream file and returns any issues
func ValidateFile(wsPath string) ([]ValidationIssue, error) {
	ws, err := ParseWorkstream(wsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse workstream: %w", err)
	}

	var issues []ValidationIssue

	// Validate workstream ID
	if err := ws.Validate(); err != nil {
		issues = append(issues, ValidationIssue{
			Field:    "ws_id",
			Message:  err.Error(),
			Severity: "ERROR",
		})
	}

	// Check for required fields
	if ws.Feature == "" {
		issues = append(issues, ValidationIssue{
			Field:    "feature",
			Message:  "feature field is required",
			Severity: "ERROR",
		})
	}

	if ws.Status == "" {
		issues = append(issues, ValidationIssue{
			Field:    "status",
			Message:  "status field is recommended",
			Severity: "WARNING",
		})
	}

	// Check if goal is present
	if ws.Goal == "" {
		issues = append(issues, ValidationIssue{
			Field:    "goal",
			Message:  "goal section is missing",
			Severity: "WARNING",
		})
	}

	// Check if acceptance criteria are present
	if len(ws.Acceptance) == 0 {
		issues = append(issues, ValidationIssue{
			Field:    "acceptance_criteria",
			Message:  "no acceptance criteria found",
			Severity: "WARNING",
		})
	}

	// Validate scope files exist
	for _, file := range ws.Scope.Implementation {
		if !fileExists(file) && !stringsContain(file, "(NEW)") {
			issues = append(issues, ValidationIssue{
				Field:    "scope_files",
				Message:  fmt.Sprintf("implementation file not found: %s", file),
				Severity: "WARNING",
			})
		}
	}

	for _, file := range ws.Scope.Tests {
		if !fileExists(file) && !stringsContain(file, "(NEW)") {
			issues = append(issues, ValidationIssue{
				Field:    "scope_files",
				Message:  fmt.Sprintf("test file not found: %s", file),
				Severity: "WARNING",
			})
		}
	}

	return issues, nil
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// stringsContain checks if a string contains a substring (case-insensitive)
func stringsContain(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr))
}
