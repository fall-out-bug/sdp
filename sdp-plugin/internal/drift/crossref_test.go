package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCrossRefValidator_Validate(t *testing.T) {
	// AC3: Docs↔Docs drift via cross-reference validation
	tmpDir := t.TempDir()

	// Create docs directory
	docsDir := filepath.Join(tmpDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatalf("Failed to create docs dir: %v", err)
	}

	// Create doc with broken link
	docContent := `# API Documentation

See [Implementation Guide](./guide.md) for details.

Also check [Missing Doc](./missing.md) which doesn't exist.
`
	docPath := filepath.Join(docsDir, "api.md")
	if err := os.WriteFile(docPath, []byte(docContent), 0644); err != nil {
		t.Fatalf("Failed to create doc: %v", err)
	}

	// Create the referenced guide
	guideContent := `# Implementation Guide

Guide content here.
`
	guidePath := filepath.Join(docsDir, "guide.md")
	if err := os.WriteFile(guidePath, []byte(guideContent), 0644); err != nil {
		t.Fatalf("Failed to create guide: %v", err)
	}

	validator := NewCrossRefValidator(tmpDir)
	issues, err := validator.Validate()
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}

	// Should find broken link to missing.md
	found := false
	for _, issue := range issues {
		if issue.Type == DriftTypeDocsDocs && issue.Severity == SeverityWarning {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected warning for broken cross-reference")
	}
}

func TestCrossRefValidator_ExtractLinks(t *testing.T) {
	validator := &CrossRefValidator{}

	content := `# Document

See [Link 1](./file1.md) and [Link 2](../file2.md).

Also [External](https://example.com) should be ignored.
`
	links := validator.extractLinks(content)

	// Should find 2 internal links
	if len(links) != 2 {
		t.Errorf("Expected 2 internal links, got %d", len(links))
	}
}

func TestCrossRefValidator_CheckLink(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file
	filePath := filepath.Join(tmpDir, "exists.md")
	if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	validator := NewCrossRefValidator(tmpDir)

	// Check existing link
	if !validator.checkLink(tmpDir, "exists.md") {
		t.Error("Expected checkLink to return true for existing file")
	}

	// Check missing link
	if validator.checkLink(tmpDir, "missing.md") {
		t.Error("Expected checkLink to return false for missing file")
	}
}
