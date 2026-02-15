package main

import (
	"os"
	"testing"

	"github.com/fall-out-bug/sdp/internal/sdpinit"
)

// TestDetectProjectType tests automatic project type detection via sdpinit
func TestDetectProjectType(t *testing.T) {
	tests := []struct {
		name         string
		createFile   string
		expectedType string
	}{
		{
			name:         "python project",
			createFile:   "pyproject.toml",
			expectedType: "python",
		},
		{
			name:         "go project",
			createFile:   "go.mod",
			expectedType: "go",
		},
		{
			name:         "node project",
			createFile:   "package.json",
			expectedType: "node",
		},
		{
			name:         "unknown project",
			createFile:   "",
			expectedType: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create test file if specified
			if tt.createFile != "" {
				filePath := tmpDir + "/" + tt.createFile
				if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			// Change to temp directory
			originalWd, _ := os.Getwd()
			t.Cleanup(func() { os.Chdir(originalWd) })
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("Failed to chdir: %v", err)
			}

			result := sdpinit.DetectProjectType()
			if result != tt.expectedType {
				t.Errorf("DetectProjectType() = %s, want %s", result, tt.expectedType)
			}
		})
	}
}

// TestInitCmd tests the init command
func TestInitCmd(t *testing.T) {
	// Get original working directory (repo root)
	originalWd, _ := os.Getwd()

	// Create temp directory
	tmpDir := t.TempDir()

	// Change to temp directory
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	// Create prompts directory (init command requires it)
	if err := os.MkdirAll("prompts/skills", 0755); err != nil {
		t.Fatalf("Failed to create prompts dir: %v", err)
	}
	if err := os.WriteFile("prompts/skills/test.md", []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create test prompt: %v", err)
	}

	// Test init with python project type
	cmd := initCmd()
	if err := cmd.Flags().Set("project-type", "python"); err != nil {
		t.Fatalf("Failed to set project-type flag: %v", err)
	}

	// Run init - this will create .claude directory
	err := cmd.RunE(cmd, []string{})

	// Should succeed
	if err != nil {
		t.Errorf("initCmd() failed: %v", err)
	}

	// Check that .claude directory was created
	if _, err := os.Stat(".claude"); os.IsNotExist(err) {
		t.Error("initCmd() did not create .claude directory")
	}
}

// TestInitCmdWithSkipBeads tests init with skip-beads flag
func TestInitCmdWithSkipBeads(t *testing.T) {
	// Get original working directory (repo root)
	originalWd, _ := os.Getwd()

	// Create temp directory
	tmpDir := t.TempDir()

	// Change to temp directory
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	// Create prompts directory (init command requires it)
	if err := os.MkdirAll("prompts/skills", 0755); err != nil {
		t.Fatalf("Failed to create prompts dir: %v", err)
	}
	if err := os.WriteFile("prompts/skills/test.md", []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create test prompt: %v", err)
	}

	cmd := initCmd()
	if err := cmd.Flags().Set("project-type", "go"); err != nil {
		t.Fatalf("Failed to set project-type flag: %v", err)
	}
	if err := cmd.Flags().Set("skip-beads", "true"); err != nil {
		t.Fatalf("Failed to set skip-beads flag: %v", err)
	}

	err := cmd.RunE(cmd, []string{})

	// Should succeed
	if err != nil {
		t.Errorf("initCmd() with skip-beads failed: %v", err)
	}

	// Check that .claude directory was created
	if _, err := os.Stat(".claude"); os.IsNotExist(err) {
		t.Error("initCmd() did not create .claude directory")
	}
}

// TestInitCmdWithAuto tests init with --auto flag
func TestInitCmdWithAuto(t *testing.T) {
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()

	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	// Create prompts directory
	if err := os.MkdirAll("prompts/skills", 0755); err != nil {
		t.Fatalf("Failed to create prompts dir: %v", err)
	}
	if err := os.WriteFile("prompts/skills/test.md", []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create test prompt: %v", err)
	}

	cmd := initCmd()
	if err := cmd.Flags().Set("auto", "true"); err != nil {
		t.Fatalf("Failed to set auto flag: %v", err)
	}

	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("initCmd() with auto failed: %v", err)
	}

	if _, err := os.Stat(".claude"); os.IsNotExist(err) {
		t.Error("initCmd() did not create .claude directory")
	}
}
