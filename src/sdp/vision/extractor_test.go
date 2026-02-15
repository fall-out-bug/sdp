package vision

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFeatureDraft_Slug(t *testing.T) {
	tests := []struct {
		title    string
		expected string
	}{
		{"User Authentication", "user-authentication"},
		{"API Rate Limiting", "api-rate-limiting"},
		{"Feature 123: Test", "feature-123-test"},
		{"Test@#$%Special!Chars", "testspecialchars"},
		{"Multiple Spaces", "multiple-spaces"},
		{"UPPERCASE", "uppercase"},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			f := &FeatureDraft{Title: tt.title}
			slug := f.Slug()
			if slug != tt.expected {
				t.Errorf("Slug() = %q, want %q", slug, tt.expected)
			}
		})
	}
}

func TestExtractFeaturesFromPRD(t *testing.T) {
	prdContent := `# Product Requirements Document

## Overview
This is the overview.

## Features

### P0 (Critical)
- Feature 1: User Authentication
- Feature 2: Database Connection

### P1 (High)
- 3: Email Notifications
- 4: API Rate Limiting

### P2 (Medium)
- Feature 5: Dark Mode

## Other Section
Some content.
`
	tempDir := t.TempDir()
	prdPath := filepath.Join(tempDir, "prd.md")
	os.WriteFile(prdPath, []byte(prdContent), 0644)

	features, err := ExtractFeaturesFromPRD(prdPath)
	if err != nil {
		t.Fatalf("ExtractFeaturesFromPRD failed: %v", err)
	}

	// Should extract P0 and P1 features, skip P2
	if len(features) != 4 {
		t.Errorf("expected 4 features, got %d", len(features))
	}

	// Check first feature
	if len(features) > 0 && features[0].Title != "User Authentication" {
		t.Errorf("expected first feature 'User Authentication', got %q", features[0].Title)
	}

	// Check priority
	if len(features) > 0 && features[0].Priority != "P0" {
		t.Errorf("expected P0 priority, got %s", features[0].Priority)
	}
}

func TestExtractFeaturesFromPRD_NoFile(t *testing.T) {
	_, err := ExtractFeaturesFromPRD("/nonexistent/path/prd.md")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestExtractFeaturesFromPRD_NoFeaturesSection(t *testing.T) {
	prdContent := `# Document
## Overview
No features section here.
`
	tempDir := t.TempDir()
	prdPath := filepath.Join(tempDir, "prd.md")
	os.WriteFile(prdPath, []byte(prdContent), 0644)

	features, err := ExtractFeaturesFromPRD(prdPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(features) != 0 {
		t.Errorf("expected 0 features, got %d", len(features))
	}
}

func TestExtractFeaturesFromPRD_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	prdPath := filepath.Join(tempDir, "empty.md")
	os.WriteFile(prdPath, []byte(""), 0644)

	features, err := ExtractFeaturesFromPRD(prdPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(features) != 0 {
		t.Errorf("expected 0 features from empty file, got %d", len(features))
	}
}

func TestFeatureDraft_Fields(t *testing.T) {
	f := FeatureDraft{
		Title:       "Test Feature",
		Description: "A test feature",
		Priority:    "P0",
	}

	if f.Title != "Test Feature" {
		t.Error("title not set")
	}
	if f.Description != "A test feature" {
		t.Error("description not set")
	}
	if f.Priority != "P0" {
		t.Error("priority not set")
	}
}
