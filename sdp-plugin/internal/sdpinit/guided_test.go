package sdpinit

import (
	"os"
	"testing"
)

func TestRunGuided_AllStepsPass(t *testing.T) {
	// Create temp directory and initialize git repo
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Initialize git repo for the test
	if err := os.WriteFile("go.mod", []byte("module test"), 0644); err != nil {
		t.Fatalf("create go.mod: %v", err)
	}

	// Create .git directory to simulate git repo
	if err := os.MkdirAll(".git", 0755); err != nil {
		t.Fatalf("create .git: %v", err)
	}

	cfg := GuidedConfig{
		AutoFix: false,
	}

	result, err := RunGuided(cfg)
	if err != nil {
		t.Fatalf("RunGuided() error: %v", err)
	}

	if result == nil {
		t.Fatal("RunGuided() returned nil result")
	}

	if len(result.Steps) == 0 {
		t.Fatal("RunGuided() returned no steps")
	}

	// Git step should pass (git is typically installed)
	gitStep := result.Steps[0]
	if !gitStep.Passed {
		t.Logf("Warning: Git step failed, but this might be expected in some environments")
	}

	// Git repo step should pass (we created .git)
	repoStep := result.Steps[1]
	if !repoStep.Passed {
		t.Errorf("Git repo step should pass, got: %s", repoStep.Message)
	}

	// Project type step should pass (we created go.mod)
	projectStep := result.Steps[2]
	if !projectStep.Passed {
		t.Errorf("Project type step should pass, got: %s", projectStep.Message)
	}
}

func TestRunGuided_MissingGitRepo(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Don't create .git directory
	cfg := GuidedConfig{
		AutoFix: false,
	}

	result, err := RunGuided(cfg)
	if err != nil {
		t.Fatalf("RunGuided() error: %v", err)
	}

	// Git repo step should fail
	repoStep := result.Steps[1]
	if repoStep.Passed {
		t.Error("Git repo step should fail without .git directory")
	}
}

func TestRunGuided_UnknownProjectType(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Create .git but no project files
	if err := os.MkdirAll(".git", 0755); err != nil {
		t.Fatalf("create .git: %v", err)
	}

	cfg := GuidedConfig{
		AutoFix: false,
	}

	result, err := RunGuided(cfg)
	if err != nil {
		t.Fatalf("RunGuided() error: %v", err)
	}

	// Project type step should fail
	projectStep := result.Steps[2]
	if projectStep.Passed {
		t.Error("Project type step should fail with unknown project type")
	}
}

func TestRunGuided_AutoFix(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Create go.mod for project detection
	if err := os.WriteFile("go.mod", []byte("module test"), 0644); err != nil {
		t.Fatalf("create go.mod: %v", err)
	}

	// Don't create .git - let auto-fix handle it
	cfg := GuidedConfig{
		AutoFix: true,
	}

	result, err := RunGuided(cfg)
	if err != nil {
		t.Fatalf("RunGuided() error: %v", err)
	}

	// After auto-fix, .git should exist
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		t.Error("Auto-fix should have created .git directory")
	}

	// Check that result reflects successful fix
	if len(result.NextSteps) == 0 {
		t.Error("Result should have next steps")
	}
}

func TestCheckGitStep(t *testing.T) {
	result := checkGitStep()

	// Git should be installed in most environments
	if result.Message == "" {
		t.Error("checkGitStep should return a message")
	}
}

func TestCheckGitRepoStep(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })

	// Test without .git
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	result := checkGitRepoStep()
	if result.Passed {
		t.Error("checkGitRepoStep should fail without .git")
	}

	// Test with .git
	if err := os.MkdirAll(".git", 0755); err != nil {
		t.Fatalf("create .git: %v", err)
	}

	result = checkGitRepoStep()
	if !result.Passed {
		t.Error("checkGitRepoStep should pass with .git")
	}
}

func TestFixGitRepo(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Fix should create .git
	err := fixGitRepo()
	if err != nil {
		t.Fatalf("fixGitRepo() error: %v", err)
	}

	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		t.Error("fixGitRepo should create .git directory")
	}
}

func TestCheckProjectTypeStep(t *testing.T) {
	tests := []struct {
		name         string
		createFile   string
		expectPass   bool
		expectDetail string
	}{
		{
			name:         "go project",
			createFile:   "go.mod",
			expectPass:   true,
			expectDetail: "go",
		},
		{
			name:         "python project",
			createFile:   "pyproject.toml",
			expectPass:   true,
			expectDetail: "python",
		},
		{
			name:         "node project",
			createFile:   "package.json",
			expectPass:   true,
			expectDetail: "node",
		},
		{
			name:       "unknown project",
			createFile: "",
			expectPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			originalWd, _ := os.Getwd()
			t.Cleanup(func() { os.Chdir(originalWd) })

			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("chdir: %v", err)
			}

			if tt.createFile != "" {
				if err := os.WriteFile(tt.createFile, []byte("test"), 0644); err != nil {
					t.Fatalf("create file: %v", err)
				}
			}

			result := checkProjectTypeStep()

			if result.Passed != tt.expectPass {
				t.Errorf("checkProjectTypeStep() passed = %v, want %v", result.Passed, tt.expectPass)
			}

			if tt.expectPass && result.Details != tt.expectDetail {
				t.Errorf("checkProjectTypeStep() details = %v, want %v", result.Details, tt.expectDetail)
			}
		})
	}
}

func TestCheckProjectTypeWithOverride(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// No project files, but explicit type provided
	result := CheckProjectTypeWithOverride("go")
	if !result.Passed {
		t.Error("CheckProjectTypeWithOverride should pass with explicit type")
	}
	if result.Details != "go" {
		t.Errorf("Details = %s, want 'go'", result.Details)
	}

	// Empty explicit type should fall back to detection
	result = CheckProjectTypeWithOverride("")
	if result.Passed {
		t.Error("CheckProjectTypeWithOverride with empty type should fail on unknown project")
	}
}

func TestCheckClaudeCodeStep(t *testing.T) {
	result := checkClaudeCodeStep()

	// Claude Code is optional, so should always pass
	if !result.Passed {
		t.Error("checkClaudeCodeStep should always pass (optional)")
	}

	if result.Message == "" {
		t.Error("checkClaudeCodeStep should return a message")
	}
}

func TestCheckBeadsStep(t *testing.T) {
	result := checkBeadsStep()

	// Beads is optional, so should always pass
	if !result.Passed {
		t.Error("checkBeadsStep should always pass (optional)")
	}

	if result.Message == "" {
		t.Error("checkBeadsStep should return a message")
	}
}

func TestGetGuidedSteps(t *testing.T) {
	steps := getGuidedSteps()

	if len(steps) == 0 {
		t.Fatal("getGuidedSteps should return at least one step")
	}

	// Check each step has required fields
	for _, step := range steps {
		if step.ID == "" {
			t.Error("Step missing ID")
		}
		if step.Name == "" {
			t.Error("Step missing Name")
		}
		if step.Description == "" {
			t.Error("Step missing Description")
		}
		if step.Check == nil {
			t.Error("Step missing Check function")
		}
	}

	// Verify expected step order
	expectedIDs := []string{"git", "git-repo", "project-type", "claude-code", "beads"}
	for i, expectedID := range expectedIDs {
		if i >= len(steps) {
			break
		}
		if steps[i].ID != expectedID {
			t.Errorf("Step %d ID = %s, want %s", i, steps[i].ID, expectedID)
		}
	}
}

func TestPrintGuidedResult(t *testing.T) {
	// This test just verifies PrintGuidedResult doesn't panic
	result := &GuidedResult{
		Steps: []GuidedStepResult{
			{Passed: true, Message: "Test pass"},
			{Passed: false, Message: "Test fail"},
		},
		AllPassed: false,
		NextSteps: []string{"Fix things", "Try again"},
	}

	// Should not panic
	PrintGuidedResult(result)

	// Test with all passed
	result.AllPassed = true
	PrintGuidedResult(result)
}

func TestGuidedStepResult_Fields(t *testing.T) {
	result := GuidedStepResult{
		Passed:  true,
		Message: "Test message",
		Details: "Test details",
	}

	if !result.Passed {
		t.Error("Passed should be true")
	}
	if result.Message != "Test message" {
		t.Errorf("Message = %s, want 'Test message'", result.Message)
	}
	if result.Details != "Test details" {
		t.Errorf("Details = %s, want 'Test details'", result.Details)
	}
}

func TestGuidedResult_Fields(t *testing.T) {
	result := &GuidedResult{
		Steps: []GuidedStepResult{
			{Passed: true},
		},
		AllPassed: true,
		NextSteps: []string{"step1"},
	}

	if len(result.Steps) != 1 {
		t.Errorf("Steps length = %d, want 1", len(result.Steps))
	}
	if !result.AllPassed {
		t.Error("AllPassed should be true")
	}
	if len(result.NextSteps) != 1 {
		t.Errorf("NextSteps length = %d, want 1", len(result.NextSteps))
	}
}
