package vision

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/src/sdp/vision"
)

func TestExtractFeaturesFromPRD(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Create test PRD
	prdPath := filepath.Join(tmpDir, "PRD.md")
	prdContent := `# Product Requirements Document

## Features

### P0 (Must Have)
- Feature 1: User authentication
- Feature 2: Task creation

### P1 (Should Have)
- Feature 3: Calendar integration
- Feature 4: Notifications

### P2 (Nice to Have)
- Feature 5: Analytics
`
	err := os.WriteFile(prdPath, []byte(prdContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test PRD: %v", err)
	}

	// Extract features
	features, err := vision.ExtractFeaturesFromPRD(prdPath)
	if err != nil {
		t.Fatalf("Failed to extract features: %v", err)
	}

	// Should extract only P0 and P1 (not P2)
	t.Logf("Found %d features:", len(features))
	for i, f := range features {
		t.Logf("  [%d] %s (priority: %s)", i, f.Title, f.Priority)
	}
	if len(features) != 4 {
		t.Errorf("Expected 4 features (P0 + P1), got %d", len(features))
	}

	// Check first feature
	if features[0].Title != "User authentication" {
		t.Errorf("Expected 'User authentication', got '%s'", features[0].Title)
	}
	if features[0].Priority != "P0" {
		t.Errorf("Expected priority 'P0', got '%s'", features[0].Priority)
	}

	// Check third feature (first P1)
	if features[2].Priority != "P1" {
		t.Errorf("Expected priority 'P1', got '%s'", features[2].Priority)
	}
}

func TestFeatureDraftSlug(t *testing.T) {
	draft := vision.FeatureDraft{
		Title:    "Calendar Integration",
		Priority: "P1",
	}

	slug := draft.Slug()
	if slug != "calendar-integration" {
		t.Errorf("Expected 'calendar-integration', got '%s'", slug)
	}
}
