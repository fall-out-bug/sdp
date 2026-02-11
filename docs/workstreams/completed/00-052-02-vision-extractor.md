# 00-052-02: Vision Extractor Implementation

> **Beads ID:** sdp-wbyc
> **Feature:** F052 - Multi-Agent SDP + @vision + @reality
> **Phase:** 1A - Strategic Skills (@vision)
> **Size:** MEDIUM
> **Duration:** 2-3 days
> **Dependencies:**
> - 00-052-01 (@vision Skill Structure)

## Goal

Implement feature extractor that parses PRD and generates feature drafts.

## Acceptance Criteria

- **AC1:** `src/sdp/vision/extractor.go` created with FeatureExtractor struct
- **AC2:** `ExtractFeaturesFromPRD()` function parses PRD and extracts P0/P1 features
- **AC3:** Each feature gets slug, title, priority, description
- **AC4:** Feature drafts written to `docs/drafts/feature-{slug}.md`
- **AC5:** `tests/sdp/vision/extractor_test.go` with ≥80% coverage

## Files

**Create:**
- `src/sdp/vision/extractor.go` - Feature extraction logic
- `tests/sdp/vision/extractor_test.go` - Test suite

**Modify:**
- None

## Steps

### Step 1: Implement Feature Extractor

Create `src/sdp/vision/extractor.go`:

```go
package vision

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Feature represents a product feature extracted from PRD
type Feature struct {
	Slug        string
	Title       string
	Priority    string // P0, P1, P2, P3
	Description string
}

// FeatureExtractor extracts features from PRD
type FeatureExtractor struct {
	prdPath string
	draftsDir string
}

// NewFeatureExtractor creates a new extractor
func NewFeatureExtractor(prdPath, draftsDir string) *FeatureExtractor {
	return &FeatureExtractor{
		prdPath: prdPath,
		draftsDir: draftsDir,
	}
}

// ExtractFeaturesFromPRD parses PRD and extracts features
func (e *FeatureExtractor) ExtractFeaturesFromPRD() ([]*Feature, error) {
	// Read PRD
	content, err := os.ReadFile(e.prdPath)
	if err != nil {
		return nil, err
	}

	// Find P0/P1 sections
	p0Features := e.extractFeaturesByPriority(string(content), "P0")
	p1Features := e.extractFeaturesByPriority(string(content), "P1")

	features := append(p0Features, p1Features...)
	return features, nil
}

// extractFeaturesByPriority finds features under priority heading
func (e *FeatureExtractor) extractFeaturesByPriority(content, priority string) []*Feature {
	features := []*Feature{}

	// Find priority section
	re := regexp.MustCompile(`### ` + priority + `.*?\n((?:-.*?\n)*)`)
	matches := re.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			lines := strings.Split(match[1], "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "- Feature") || strings.HasPrefix(line, "-") {
					feature := e.parseFeatureLine(line, priority)
					if feature != nil {
						features = append(features, feature)
					}
				}
			}
		}
	}

	return features
}

// parseFeatureLine parses "Feature N: Title - Description" format
func (e *FeatureExtractor) parseFeatureLine(line, priority string) *Feature {
	// Remove "- " prefix
	line = strings.TrimPrefix(line, "- ")

	// Extract title and description
	parts := strings.SplitN(line, ":", 2)
	if len(parts) < 2 {
		return nil
	}

	title := strings.TrimSpace(parts[0])
	title = strings.TrimPrefix(title, "Feature ")
	title = strings.TrimSpace(title)

	description := strings.TrimSpace(parts[1])
	description = strings.TrimPrefix(description, "-")
	description = strings.TrimSpace(description)

	// Generate slug
	slug := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
	slug = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(slug, "")

	return &Feature{
		Slug:        slug,
		Title:       title,
		Priority:    priority,
		Description: description,
	}
}

// WriteFeatureDrafts creates markdown files for each feature
func (e *FeatureExtractor) WriteFeatureDrafts(features []*Feature) error {
	// Ensure drafts directory exists
	if err := os.MkdirAll(e.draftsDir, 0755); err != nil {
		return err
	}

	for _, feature := range features {
		draft := e.generateFeatureDraft(feature)
		path := filepath.Join(e.draftsDir, "feature-"+feature.Slug+".md")

		if err := os.WriteFile(path, []byte(draft), 0644); err != nil {
			return err
		}
	}

	return nil
}

// generateFeatureDraft creates markdown content for feature
func (e *FeatureExtractor) generateFeatureDraft(feature *Feature) string {
	return `# Feature: ` + feature.Title + `

> **Priority:** ` + feature.Priority + `
> **Slug:** ` + feature.Slug + `

## Description

` + feature.Description + `

## Requirements

TBD - Extracted from PRD

## Success Criteria

TBD - Define in @idea phase

## Dependencies

TBD - Identify in @design phase
`
}
```

### Step 2: Write Tests

Create `tests/sdp/vision/extractor_test.go`:

```go
package vision

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/src/sdp/vision"
)

func TestExtractFeaturesFromPRD(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	prdPath := filepath.Join(tmpDir, "PRD.md")
	prdContent := `
# Product Requirements

## P0 (Must Have)
- Feature 1: User authentication - Login/register with email
- Feature 2: Task creation - Users can create tasks

## P1 (Should Have)
- Feature 3: Calendar integration - Sync with external calendars
`
	os.WriteFile(prdPath, []byte(prdContent), 0644)

	extractor := vision.NewFeatureExtractor(prdPath, tmpDir)

	// Act
	features, err := extractor.ExtractFeaturesFromPRD()
	if err != nil {
		t.Fatalf("ExtractFeaturesFromPRD failed: %v", err)
	}

	// Assert
	if len(features) != 3 {
		t.Errorf("Expected 3 features, got %d", len(features))
	}

	// Check P0 features
	p0Count := 0
	for _, f := range features {
		if f.Priority == "P0" {
			p0Count++
		}
	}
	if p0Count != 2 {
		t.Errorf("Expected 2 P0 features, got %d", p0Count)
	}
}

func TestFeatureSlug(t *testing.T) {
	// Test slug generation
}
```

### Step 3: Verify Coverage

```bash
go test -v -coverprofile=coverage.out ./tests/sdp/vision
go tool cover -func=coverage.out | grep extractor.go
```

Target: ≥80% coverage

## Quality Gates

- Files < 200 LOC each
- Coverage ≥80%
- No `go vet` warnings
- Tests use t.TempDir() for cleanup

## Success Metrics

- Can parse PRD and extract P0/P1 features
- Generates valid markdown drafts
- Slug generation produces clean URLs
- Tests cover all extraction paths
