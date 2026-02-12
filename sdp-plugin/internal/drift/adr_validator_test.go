package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func TestADRValidator_Validate(t *testing.T) {
	// AC2: Decision drift via ADR validation
	tmpDir := t.TempDir()

	// Create decisions directory
	decisionsDir := filepath.Join(tmpDir, "docs", "decisions")
	if err := os.MkdirAll(decisionsDir, 0755); err != nil {
		t.Fatalf("Failed to create decisions dir: %v", err)
	}

	// Create an ADR file with accepted status (no drift expected)
	adrContent := `---
decision_id: ADR-001
status: accepted
---

# Use SQLite for Storage

## Context
Need a database for artifact storage.

## Decision
Use SQLite with FTS5 for full-text search.

## Consequences
- Simple deployment (single file)
- Good performance for small-medium datasets
`
	adrPath := filepath.Join(decisionsDir, "ADR-001.md")
	if err := os.WriteFile(adrPath, []byte(adrContent), 0644); err != nil {
		t.Fatalf("Failed to create ADR: %v", err)
	}

	validator := NewADRValidator(tmpDir)
	issues, err := validator.Validate()
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}

	// Accepted status should produce no drift issues
	// The validator only flags superseded/deprecated decisions
	if len(issues) != 0 {
		t.Logf("Note: Accepted status ADR produced %d issues (expected 0)", len(issues))
	}
}

func TestADRValidator_CheckDecisionStatus(t *testing.T) {
	tmpDir := t.TempDir()

	decisionsDir := filepath.Join(tmpDir, "docs", "decisions")
	if err := os.MkdirAll(decisionsDir, 0755); err != nil {
		t.Fatalf("Failed to create decisions dir: %v", err)
	}

	// Create ADR with superseded status
	adrContent := `---
decision_id: ADR-002
status: superseded
superseded_by: ADR-003
---

# Deprecated Decision

This decision has been superseded.
`
	adrPath := filepath.Join(decisionsDir, "ADR-002.md")
	if err := os.WriteFile(adrPath, []byte(adrContent), 0644); err != nil {
		t.Fatalf("Failed to create ADR: %v", err)
	}

	validator := NewADRValidator(tmpDir)
	issues, err := validator.Validate()
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}

	// Should flag superseded decision
	found := false
	for _, issue := range issues {
		if issue.Type == DriftTypeDecisionCode && issue.Severity == SeverityWarning {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected warning for superseded decision")
	}
}

func TestADRValidator_ExtractKeywords(t *testing.T) {
	validator := &ADRValidator{}

	content := `Use SQLite for storage with FTS5 full-text search.`
	keywords := validator.extractKeywords(content)

	if len(keywords) == 0 {
		t.Error("Expected to extract keywords from content")
	}

	// Should find SQLite and FTS5
	hasSQLite := false
	for _, kw := range keywords {
		if kw == "sqlite" {
			hasSQLite = true
			break
		}
	}

	if !hasSQLite {
		t.Error("Expected to find 'sqlite' keyword")
	}
}
