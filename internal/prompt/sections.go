package prompt

import (
	"strings"
)

// WorkstreamSpec holds task and boundary data for prompt section rendering.
// Callers construct from workstream markdown, IssueInput, or beads.Issue.
type WorkstreamSpec struct {
	ID                 string
	Title              string
	Description        string
	AcceptanceCriteria []string
	ScopeFiles         []string
	OutOfScope         []string
	SpecID             string
}

// BoundaryInput holds path/scope constraints for BoundarySection.
// Use AllowedPathPrefixes/ForbiddenPathPrefixes for path-based boundaries (llm.BoundarySpec).
// Use ScopeFiles/OutOfScope from WorkstreamSpec for workstream-based boundaries.
type BoundaryInput struct {
	AllowedPathPrefixes   []string
	ForbiddenPathPrefixes []string
	ControlPathPrefixes   []string
	ScopeFiles            []string
	OutOfScope           []string
}

// EvidenceInput holds checkpoint/evidence context for EvidenceSection.
// Callers populate from orchestrate.Checkpoint or evidence file content.
type EvidenceInput struct {
	Content      string   // raw evidence content (e.g. from .sdp/evidence/*.json)
	CompletedWS  []string // e.g. "00-025-01 (abc123)"
	ReviewStatus string
}

// TaskSectionForReview renders task in compact format for review prompts.
// Pure function: no side effects, no file I/O.
func TaskSectionForReview(ws WorkstreamSpec) string {
	var b strings.Builder
	b.WriteString("## Task\n")
	b.WriteString("ID: " + ws.ID + "\n")
	b.WriteString("Title: " + ws.Title + "\n")
	if ws.Description != "" {
		b.WriteString("Description: " + ws.Description + "\n")
	}
	return b.String()
}

// TaskSection renders task description and acceptance criteria.
// Pure function: no side effects, no file I/O.
func TaskSection(ws WorkstreamSpec) string {
	var b strings.Builder
	b.WriteString("## Task\n\n")
	b.WriteString("**ID:** " + ws.ID + "\n\n")
	b.WriteString("**Title:** " + ws.Title + "\n\n")
	if ws.Description != "" {
		b.WriteString("**Description:**\n")
		b.WriteString(ws.Description)
		b.WriteString("\n\n")
	}
	if len(ws.AcceptanceCriteria) > 0 {
		b.WriteString("**Acceptance Criteria:**\n")
		for _, ac := range ws.AcceptanceCriteria {
			b.WriteString("- ")
			b.WriteString(ac)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
	if ws.SpecID != "" {
		b.WriteString("**Spec ID:** " + ws.SpecID + "\n\n")
	}
	return b.String()
}

// BoundarySection renders scope files and out-of-scope constraints.
// Supports both path-prefix style (llm.BoundarySpec) and scope-files style (WorkstreamSpec).
// Pure function: no side effects, no file I/O.
func BoundarySection(in BoundaryInput) string {
	var b strings.Builder
	b.WriteString("## Constraints\n\n")
	if len(in.AllowedPathPrefixes) > 0 {
		b.WriteString("You may ONLY modify files under these path prefixes:\n")
		for _, p := range in.AllowedPathPrefixes {
			b.WriteString("- " + p + "\n")
		}
		b.WriteString("\n")
	}
	if len(in.ScopeFiles) > 0 {
		b.WriteString("Scope files (you may modify):\n")
		for _, f := range in.ScopeFiles {
			b.WriteString("- `" + f + "`\n")
		}
		b.WriteString("\n")
	}
	if len(in.ForbiddenPathPrefixes) > 0 || len(in.ControlPathPrefixes) > 0 || len(in.OutOfScope) > 0 {
		b.WriteString("You must NOT modify:\n")
		for _, p := range in.ForbiddenPathPrefixes {
			b.WriteString("- " + p + "\n")
		}
		for _, p := range in.ControlPathPrefixes {
			b.WriteString("- " + p + "\n")
		}
		for _, f := range in.OutOfScope {
			b.WriteString("- " + f + "\n")
		}
		b.WriteString("\n")
	}
	b.WriteString("Produce working, testable code. Run `go test ./...` to verify.\n")
	return b.String()
}

// AcceptanceCriteriaSection renders acceptance criteria for context packet.
// Pure function: no side effects, no file I/O.
func AcceptanceCriteriaSection(items []string) string {
	var b strings.Builder
	b.WriteString("### Acceptance Criteria\n\n")
	for _, ac := range items {
		b.WriteString("- ")
		b.WriteString(ac)
		b.WriteString("\n")
	}
	b.WriteString("\n")
	return b.String()
}

// ScopeFilesSection renders scope files list for context packet.
// Pure function: no side effects, no file I/O.
func ScopeFilesSection(files []string) string {
	var b strings.Builder
	b.WriteString("### Scope Files\n\n")
	for _, f := range files {
		b.WriteString("- ")
		b.WriteString(f)
		b.WriteString("\n")
	}
	b.WriteString("\n")
	return b.String()
}

// EvidenceSection renders evidence context for review prompts.
// Pure function: no side effects, no file I/O.
func EvidenceSection(in EvidenceInput) string {
	var b strings.Builder
	b.WriteString("\n## Evidence\n")
	if in.Content != "" {
		b.WriteString(in.Content)
	} else {
		b.WriteString("(no evidence file found)\n")
	}
	if len(in.CompletedWS) > 0 {
		b.WriteString("\n\n### Completed Workstreams\n")
		for _, ws := range in.CompletedWS {
			b.WriteString("- ")
			b.WriteString(ws)
			b.WriteString("\n")
		}
	}
	if in.ReviewStatus != "" {
		b.WriteString("\n### Review Status\n")
		b.WriteString(in.ReviewStatus)
		b.WriteString("\n")
	}
	return b.String()
}
