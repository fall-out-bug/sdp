package resolver

import (
	"regexp"
	"strings"
)

// Pattern definitions for ID detection (fast, no I/O)
var (
	// Standard workstream: 00-064-01
	workstreamPattern = regexp.MustCompile(`^\d{2}-\d{3}-\d{2}$`)

	// Fix workstream: 99-F064-0001
	fixWorkstreamPattern = regexp.MustCompile(`^\d{2}-[A-Z]\d{3}-\d{4}$`)

	// Beads ID: sdp-ushh, abc-123
	beadsPattern = regexp.MustCompile(`^[a-z]{3}-[a-z0-9]+$`)

	// Issue ID: ISSUE-0001
	issuePattern = regexp.MustCompile(`^ISSUE-\d+$`)
)

// DetectIDType determines the type of identifier from its pattern
// This is a pure function with no I/O - fast pattern matching only
func DetectIDType(id string) IDType {
	id = strings.TrimSpace(id)

	if id == "" {
		return TypeUnknown
	}

	// Check patterns in order of specificity
	if workstreamPattern.MatchString(id) || fixWorkstreamPattern.MatchString(id) {
		return TypeWorkstream
	}

	if beadsPattern.MatchString(id) {
		return TypeBeads
	}

	if issuePattern.MatchString(id) {
		return TypeIssue
	}

	return TypeUnknown
}
