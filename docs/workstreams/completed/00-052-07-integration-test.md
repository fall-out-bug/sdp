# 00-052-07: Test @vision + @reality Integration

> **Beads ID:** sdp-riqq
> **Feature:** F052 - Multi-Agent SDP + @vision + @reality
> **Phase:** 1B - Analysis Skills (@reality)
> **Size:** SMALL
> **Duration:** 1-2 days
> **Dependencies:**
> - 00-052-02 (Vision Extractor)
> - 00-052-05 (Project Scanner)

## Goal

End-to-end test of @vision + @reality on sample projects.

## Acceptance Criteria

- **AC1:** Test @vision on sample project (generates artifacts)
- **AC2:** Test @reality --quick on sample project (generates report)
- **AC3:** Test @reality --deep on sample project (spawns 8 experts)
- **AC4:** Verify vision vs reality gap analysis works
- **AC5:** Integration test suite in `tests/integration/vision_reality_test.go`

## Files

**Create:**
- `tests/integration/vision_reality_test.go` - Integration test suite
- `tests/fixtures/sample-project/` - Sample project for testing

**Modify:**
- None

## Steps

### Step 1: Create Sample Project

Create `tests/fixtures/sample-project/`:

```bash
mkdir -p tests/fixtures/sample-project
cd tests/fixtures/sample-project

# Create minimal Go project
cat > go.mod <<EOF
module sample

go 1.21
EOF

# Create main package
mkdir -p cmd/service
cat > cmd/service/main.go <<'EOF'
package main

func main() {
	// TODO: Add startup logic
}
EOF

# Create service package
mkdir -p internal/service
cat > internal/service/service.go <<'EOF'
package service

// Service handles business logic
type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Process() error {
	// TODO: Implement
	return nil
}
EOF
```

### Step 2: Write Integration Tests

Create `tests/integration/vision_reality_test.go`:

```go
package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/src/sdp/reality"
	"github.com/fall-out-bug/sdp/src/sdp/vision"
)

// TestVisionRealityIntegration tests full workflow
func TestVisionRealityIntegration(t *testing.T) {
	// Setup
	projectDir := filepath.Join("..", "fixtures", "sample-project")

	t.Run("@reality quick analysis", func(t *testing.T) {
		scanner := reality.NewProjectScanner(projectDir)

		// Detect language
		lang, err := scanner.DetectLanguage()
		if err != nil {
			t.Fatalf("DetectLanguage failed: %v", err)
		}

		if lang != "go" {
			t.Errorf("Expected 'go', got '%s'", lang)
		}

		// Scan project
		projectType, stats, err := scanner.ScanProject()
		if err != nil {
			t.Fatalf("ScanProject failed: %v", err)
		}

		if projectType.Language != "go" {
			t.Errorf("Expected language 'go', got '%s'", projectType.Language)
		}

		if stats.TotalFiles < 1 {
			t.Errorf("Expected at least 1 file, got %d", stats.TotalFiles)
		}
	})

	t.Run("@vision feature extraction", func(t *testing.T) {
		// Create temporary PRD
		tmpDir := t.TempDir()
		prdPath := filepath.Join(tmpDir, "PRD.md")
		prdContent := `
# Product Requirements

## P0 (Must Have)
- User authentication - Login/register with email
- Task creation - Users can create tasks

## P1 (Should Have)
- Calendar integration - Sync with external calendars
`
		os.WriteFile(prdPath, []byte(prdContent), 0644)

		extractor := vision.NewFeatureExtractor(prdPath, tmpDir)
		features, err := extractor.ExtractFeaturesFromPRD()
		if err != nil {
			t.Fatalf("ExtractFeaturesFromPRD failed: %v", err)
		}

		if len(features) != 3 {
			t.Errorf("Expected 3 features, got %d", len(features))
		}

		// Write feature drafts
		err = extractor.WriteFeatureDrafts(features)
		if err != nil {
			t.Fatalf("WriteFeatureDrafts failed: %v", err)
		}

		// Verify drafts created
		for _, feature := range features {
			draftPath := filepath.Join(tmpDir, "feature-"+feature.Slug+".md")
			if _, err := os.Stat(draftPath); os.IsNotExist(err) {
				t.Errorf("Feature draft not created: %s", draftPath)
			}
		}
	})
}

// TestVisionRealityGapAnalysis tests gap detection
func TestVisionRealityGapAnalysis(t *testing.T) {
	t.Run("Vision: Clean architecture", func(t *testing.T) {
		// Vision claims: "Clean architecture, domain-driven design"

		// Reality check: Layer violations?
		// TODO: Implement layer violation detection
	})

	t.Run("Vision: 90% test coverage", func(t *testing.T) {
		// Vision claims: "High test coverage"

		// Reality check: Actual coverage?
		// TODO: Implement coverage check
	})
}
```

### Step 3: Run Integration Tests

```bash
go test -v ./tests/integration/
```

Expected: All tests pass

### Step 4: Manual Verification

```bash
# Test @reality --quick
cd tests/fixtures/sample-project
@reality --quick

# Verify:
# - Detects Go project
# - Counts files correctly
# - Identifies structure
```

## Quality Gates

- Integration tests pass
- Manual verification successful
- Sample project is realistic
- Tests are isolated (use t.TempDir())

## Success Metrics

- @reality can analyze any project
- @vision can generate artifacts
- Integration workflow is smooth
- Gap analysis works end-to-end
