package sdpinit

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestNewHeadlessRunner(t *testing.T) {
	cfg := Config{
		ProjectType: "go",
	}

	runner := NewHeadlessRunner(cfg)

	if runner == nil {
		t.Fatal("NewHeadlessRunner returned nil")
	}

	if runner.output == nil {
		t.Error("Runner output should be initialized")
	}
}

func TestHeadlessRunner_Validate(t *testing.T) {
	tests := []struct {
		name        string
		projectType string
		conflicts   []string
		force       bool
		expectError bool
	}{
		{
			name:        "valid go project",
			projectType: "go",
			conflicts:   []string{},
			force:       false,
			expectError: false,
		},
		{
			name:        "valid node project",
			projectType: "node",
			conflicts:   []string{},
			force:       false,
			expectError: false,
		},
		{
			name:        "invalid project type",
			projectType: "invalid",
			conflicts:   []string{},
			force:       false,
			expectError: true,
		},
		{
			name:        "conflict without force",
			projectType: "go",
			conflicts:   []string{".claude/settings.json"},
			force:       false,
			expectError: true,
		},
		{
			name:        "conflict with force",
			projectType: "go",
			conflicts:   []string{".claude/settings.json"},
			force:       true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := NewHeadlessRunner(Config{
				ProjectType: tt.projectType,
				Force:       tt.force,
			})
			runner.preflight = &PreflightResult{
				ProjectType: tt.projectType,
				Conflicts:   tt.conflicts,
			}

			err := runner.validate()

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestHeadlessRunner_TrackCreatedFiles(t *testing.T) {
	runner := NewHeadlessRunner(Config{ProjectType: "go"})
	runner.trackCreatedFiles()

	if len(runner.output.Created) == 0 {
		t.Error("Should track created files")
	}

	// Check expected files
	expectedFiles := []string{".claude/", ".claude/settings.json"}
	for _, expected := range expectedFiles {
		found := false
		for _, created := range runner.output.Created {
			if created == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected created file: %s", expected)
		}
	}
}

func TestHeadlessRunner_Run_DryRun(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Create prompts directory
	if err := os.MkdirAll("prompts/skills", 0755); err != nil {
		t.Fatalf("mkdir prompts: %v", err)
	}

	cfg := Config{
		ProjectType: "go",
		DryRun:      true,
	}

	runner := NewHeadlessRunner(cfg)
	output, err := runner.Run()

	if err != nil {
		t.Fatalf("Run error: %v", err)
	}

	if !output.Success {
		t.Error("Dry run should succeed")
	}

	// In dry-run, files should not be created
	if _, err := os.Stat(".claude"); !os.IsNotExist(err) {
		t.Error("Dry run should not create .claude directory")
	}
}

func TestHeadlessRunner_Run_Actual(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Create prompts directory
	if err := os.MkdirAll("prompts/skills", 0755); err != nil {
		t.Fatalf("mkdir prompts: %v", err)
	}

	cfg := Config{
		ProjectType: "go",
		DryRun:      false,
	}

	runner := NewHeadlessRunner(cfg)
	output, err := runner.Run()

	if err != nil {
		t.Fatalf("Run error: %v", err)
	}

	if !output.Success {
		t.Error("Run should succeed")
	}

	// Files should be created
	if _, err := os.Stat(".claude"); os.IsNotExist(err) {
		t.Error("Run should create .claude directory")
	}
}

func TestHeadlessOutput_GetExitCode(t *testing.T) {
	tests := []struct {
		name     string
		output   *HeadlessOutput
		expected int
	}{
		{
			name: "success",
			output: &HeadlessOutput{
				Success: true,
			},
			expected: ExitSuccess,
		},
		{
			name: "error with code",
			output: &HeadlessOutput{
				Success:   false,
				ErrorCode: ExitValidationFailed,
			},
			expected: ExitValidationFailed,
		},
		{
			name: "error without code",
			output: &HeadlessOutput{
				Success: false,
			},
			expected: ExitError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.output.GetExitCode()
			if result != tt.expected {
				t.Errorf("GetExitCode() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestHeadlessOutput_OutputJSON(t *testing.T) {
	output := &HeadlessOutput{
		Success:     true,
		ProjectType: "go",
		Created:     []string{".claude/", ".claude/settings.json"},
		Config: &ConfigSummary{
			Skills:          []string{"feature", "build"},
			EvidenceEnabled: true,
			BeadsEnabled:    true,
		},
	}

	// Capture stdout (not easily done, so just verify it doesn't error)
	// In a real test we'd capture stdout
	err := output.OutputJSON()
	if err != nil {
		t.Errorf("OutputJSON error: %v", err)
	}
}

func TestHeadlessOutput_JSONMarshaling(t *testing.T) {
	output := &HeadlessOutput{
		Success:     true,
		ProjectType: "go",
		Created:     []string{".claude/"},
		Preflight: &PreflightResult{
			ProjectType: "go",
			HasGit:      true,
		},
	}

	data, err := json.Marshal(output)
	if err != nil {
		t.Fatalf("JSON marshal error: %v", err)
	}

	// Verify key fields
	jsonStr := string(data)
	if !strings.Contains(jsonStr, `"success":true`) {
		t.Error("JSON should contain success field")
	}
	if !strings.Contains(jsonStr, `"project_type":"go"`) {
		t.Error("JSON should contain project_type field")
	}
}

func TestRunHeadless(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Create prompts directory
	if err := os.MkdirAll("prompts/skills", 0755); err != nil {
		t.Fatalf("mkdir prompts: %v", err)
	}

	cfg := Config{
		ProjectType: "python",
		DryRun:      true,
	}

	output, err := RunHeadless(cfg)
	if err != nil {
		t.Fatalf("RunHeadless error: %v", err)
	}

	if output.ProjectType != "python" {
		t.Errorf("ProjectType = %s, want python", output.ProjectType)
	}
}

func TestHeadlessRunner_WithConflict(t *testing.T) {
	// Create temp directory with existing settings
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Create existing .claude/settings.json
	if err := os.MkdirAll(".claude", 0755); err != nil {
		t.Fatalf("mkdir .claude: %v", err)
	}
	if err := os.WriteFile(".claude/settings.json", []byte("{}"), 0644); err != nil {
		t.Fatalf("write settings: %v", err)
	}

	// Create prompts directory
	if err := os.MkdirAll("prompts/skills", 0755); err != nil {
		t.Fatalf("mkdir prompts: %v", err)
	}

	cfg := Config{
		ProjectType: "go",
		Force:       false,
	}

	runner := NewHeadlessRunner(cfg)
	_, err := runner.Run()

	// Should fail due to conflict
	if err == nil {
		t.Error("Expected error with conflict")
	}
}

func TestHeadlessRunner_ForceWithConflict(t *testing.T) {
	// Create temp directory with existing settings
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Create existing .claude/settings.json
	if err := os.MkdirAll(".claude", 0755); err != nil {
		t.Fatalf("mkdir .claude: %v", err)
	}
	if err := os.WriteFile(".claude/settings.json", []byte("{}"), 0644); err != nil {
		t.Fatalf("write settings: %v", err)
	}

	// Create prompts directory
	if err := os.MkdirAll("prompts/skills", 0755); err != nil {
		t.Fatalf("mkdir prompts: %v", err)
	}

	cfg := Config{
		ProjectType: "go",
		Force:       true,
		DryRun:      true,
	}

	runner := NewHeadlessRunner(cfg)
	output, err := runner.Run()

	// Should succeed with force
	if err != nil {
		t.Fatalf("Unexpected error with force: %v", err)
	}

	if !output.Success {
		t.Error("Should succeed with force flag")
	}
}

func TestHeadlessRunner_DetectProjectType(t *testing.T) {
	// Create temp directory with go.mod
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Create go.mod to detect go project
	if err := os.WriteFile("go.mod", []byte("module test"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	// Create prompts directory
	if err := os.MkdirAll("prompts/skills", 0755); err != nil {
		t.Fatalf("mkdir prompts: %v", err)
	}

	cfg := Config{
		// ProjectType not set - should be detected
		DryRun: true,
	}

	runner := NewHeadlessRunner(cfg)
	output, err := runner.Run()

	if err != nil {
		t.Fatalf("Run error: %v", err)
	}

	// Project type should be detected
	if output.ProjectType != "go" {
		t.Errorf("Detected ProjectType = %s, want go", output.ProjectType)
	}
}

func TestConfigSummary(t *testing.T) {
	summary := &ConfigSummary{
		Skills:          []string{"feature", "build"},
		EvidenceEnabled: true,
		BeadsEnabled:    false,
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("JSON marshal error: %v", err)
	}

	jsonStr := string(data)
	if !strings.Contains(jsonStr, "feature") {
		t.Error("JSON should contain skills")
	}
	if !strings.Contains(jsonStr, `"evidence_enabled":true`) {
		t.Error("JSON should contain evidence_enabled")
	}
}

func TestExitCodes(t *testing.T) {
	// Verify exit codes are distinct
	if ExitSuccess == ExitError {
		t.Error("ExitSuccess and ExitError should be different")
	}
	if ExitSuccess == ExitValidationFailed {
		t.Error("ExitSuccess and ExitValidationFailed should be different")
	}
	if ExitError == ExitValidationFailed {
		t.Error("ExitError and ExitValidationFailed should be different")
	}

	// Verify expected values
	if ExitSuccess != 0 {
		t.Errorf("ExitSuccess = %d, want 0", ExitSuccess)
	}
	if ExitError != 1 {
		t.Errorf("ExitError = %d, want 1", ExitError)
	}
	if ExitValidationFailed != 2 {
		t.Errorf("ExitValidationFailed = %d, want 2", ExitValidationFailed)
	}
}
