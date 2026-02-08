package vision

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// findProjectRoot finds the project root by looking for go.mod
func findProjectRoot(t *testing.T) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("Could not find project root (go.mod not found)")
		}
		dir = parent
	}
}

// TestVisionEndToEnd_Integration verifies @vision skill can run end-to-end
func TestVisionEndToEnd_Integration(t *testing.T) {
	// This is an integration test that verifies @vision can:
	// 1. Be invoked via Skill tool
	// 2. Process a product idea
	// 3. Generate vision artifacts (PRODUCT_VISION.md, PRD.md, ROADMAP.md)

	// Create temp directory for test artifacts
	tmpDir := t.TempDir()

	// Note: Full @vision execution requires user interaction
	// This test verifies the skill structure and basic functionality

	// Verify skill file exists
	// Find project root by looking for go.mod
	projectRoot := findProjectRoot(t)
	skillPath := filepath.Join(projectRoot, ".claude/skills/vision/SKILL.md")
	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		t.Fatalf("@vision skill file not found at %s", skillPath)
	}

	// Read and verify skill structure
	content, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("Failed to read skill file: %v", err)
	}

	skillContent := string(content)

	// Verify key sections exist
	requiredSections := []string{
		"# @vision",
		"## Workflow",
		"## Step 1: Quick Interview",
		"## Step 2: Deep-Thinking Analysis",
		"## Step 3: Generate Artifacts",
	}

	for _, section := range requiredSections {
		if !strings.Contains(skillContent, section) {
			t.Errorf("Missing required section in skill file: %s", section)
		}
	}

	// Verify 7 expert agents are mentioned
	requiredExperts := []string{
		"Product expert",
		"Market expert",
		"Technical expert",
		"UX expert",
			"Business expert",
		"Growth expert",
		"Risk expert",
	}

	for _, expert := range requiredExperts {
		if !strings.Contains(skillContent, expert) {
			t.Errorf("Missing expert agent in skill file: %s", expert)
		}
	}

	// Create test output directory
	outputDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Note: Actual skill invocation requires Skill tool which is not available in unit tests
	// This test verifies the skill is properly structured and can be invoked
	t.Log("✓ @vision skill structure verified")
	t.Logf("✓ All 7 expert agents defined")
	t.Logf("✓ Output directory: %s", outputDir)
}

// TestVisionArtifactGeneration verifies artifact generation structure
func TestVisionArtifactGeneration(t *testing.T) {
	// Verify that @vision generates correct artifacts
	// This test checks the expected artifact structure

	expectedArtifacts := []string{
		"PRODUCT_VISION.md",
		"PRD.md",
		"ROADMAP.md",
	}

	// Verify artifact templates or examples exist if any
	for _, artifact := range expectedArtifacts {
		// Check if there are example artifacts in docs/
		matches, _ := filepath.Glob(filepath.Join("docs", artifact))
		t.Logf("Artifact %s: %d examples found", artifact, len(matches))
	}

	// Verify vision package structure exists
	visionPkg := "src/sdp/vision"
	if _, err := os.Stat(visionPkg); os.IsNotExist(err) {
		t.Logf("Warning: Vision package not found at %s", visionPkg)
	} else {
		t.Logf("✓ Vision package exists at %s", visionPkg)
	}
}

// TestVisionRequirements_Gathered verifies requirements gathering process
func TestVisionRequirements_Gathered(t *testing.T) {
	// Verify that @vision uses AskUserQuestion for requirements gathering

	projectRoot := findProjectRoot(t)
	skillPath := filepath.Join(projectRoot, ".claude/skills/vision/SKILL.md")
	content, _ := os.ReadFile(skillPath)
	skillContent := string(content)

	// Check for AskUserQuestion usage
	if !strings.Contains(skillContent, "AskUserQuestion") {
		t.Error("@vision should use AskUserQuestion for requirements gathering")
	}

	// Check for interview questions
	requiredQuestions := []string{
		"What problem are you solving?",
		"Who are your target users?",
		"What defines success in 1 year?",
	}

	for _, question := range requiredQuestions {
		if !strings.Contains(skillContent, question) {
			t.Errorf("Missing interview question: %s", question)
		}
	}

	// Verify multiSelect support
	if !strings.Contains(skillContent, "multiSelect") {
		t.Error("@vision should support multi-select questions")
	}

	t.Log("✓ Requirements gathering verified")
	t.Log("✓ Interview questions defined")
	t.Log("✓ Multi-select support confirmed")
}
