package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
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

	// Test init with python project type and auto flag to avoid interactive prompts
	cmd := initCmd()
	if err := cmd.Flags().Set("project-type", "python"); err != nil {
		t.Fatalf("Failed to set project-type flag: %v", err)
	}
	if err := cmd.Flags().Set("auto", "true"); err != nil {
		t.Fatalf("Failed to set auto flag: %v", err)
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
	if err := cmd.Flags().Set("auto", "true"); err != nil {
		t.Fatalf("Failed to set auto flag: %v", err)
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

// TestInitCmdWithHeadless tests init with --headless flag
func TestInitCmdWithHeadless(t *testing.T) {
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

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := initCmd()
	if err := cmd.Flags().Set("headless", "true"); err != nil {
		w.Close()
		os.Stdout = oldStdout
		t.Fatalf("Failed to set headless flag: %v", err)
	}
	if err := cmd.Flags().Set("project-type", "go"); err != nil {
		w.Close()
		os.Stdout = oldStdout
		t.Fatalf("Failed to set project-type flag: %v", err)
	}

	err := cmd.RunE(cmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// The command should not return an error for successful headless init
	if err != nil {
		t.Errorf("initCmd() with headless failed: %v", err)
	}

	// Verify JSON output
	if !strings.Contains(output, `"success"`) {
		t.Errorf("Headless output should be JSON with success field, got: %s", output)
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("Output is not valid JSON: %v\nOutput was: %s", err, output)
	}
}

// TestInitCmdWithDryRun tests init with --dry-run flag
func TestInitCmdWithDryRun(t *testing.T) {
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
	if err := cmd.Flags().Set("dry-run", "true"); err != nil {
		t.Fatalf("Failed to set dry-run flag: %v", err)
	}

	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("initCmd() with dry-run failed: %v", err)
	}

	// Should NOT create .claude in dry-run mode
	if _, err := os.Stat(".claude"); !os.IsNotExist(err) {
		t.Error("Dry-run should not create .claude directory")
	}
}

// TestInitCmdWithForce tests init with --force flag
func TestInitCmdWithForce(t *testing.T) {
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()

	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	// Create existing .claude/settings.json
	if err := os.MkdirAll(".claude", 0755); err != nil {
		t.Fatalf("Failed to create .claude dir: %v", err)
	}
	if err := os.WriteFile(".claude/settings.json", []byte(`{"old": true}`), 0644); err != nil {
		t.Fatalf("Failed to create existing settings: %v", err)
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
	if err := cmd.Flags().Set("force", "true"); err != nil {
		t.Fatalf("Failed to set force flag: %v", err)
	}

	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("initCmd() with force failed: %v", err)
	}
}

// TestInitCmdWithNoEvidence tests init with --no-evidence flag
func TestInitCmdWithNoEvidence(t *testing.T) {
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
	if err := cmd.Flags().Set("no-evidence", "true"); err != nil {
		t.Fatalf("Failed to set no-evidence flag: %v", err)
	}

	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("initCmd() with no-evidence failed: %v", err)
	}

	// Verify settings.json has evidence disabled
	content, err := os.ReadFile(".claude/settings.json")
	if err != nil {
		t.Fatalf("Failed to read settings: %v", err)
	}

	if !strings.Contains(string(content), `"enabled": false`) {
		t.Error("Settings should have evidence disabled")
	}
}

// TestInitCmdWithSkills tests init with --skills flag
func TestInitCmdWithSkills(t *testing.T) {
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
	if err := cmd.Flags().Set("skills", "feature,build,review"); err != nil {
		t.Fatalf("Failed to set skills flag: %v", err)
	}

	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("initCmd() with skills failed: %v", err)
	}

	// Verify settings.json has custom skills
	content, err := os.ReadFile(".claude/settings.json")
	if err != nil {
		t.Fatalf("Failed to read settings: %v", err)
	}

	if !strings.Contains(string(content), "feature") {
		t.Error("Settings should contain custom skills")
	}
}

// TestInitCmdWithGuidedAlias tests backward-compatible --guided alias.
func TestInitCmdWithGuidedAlias(t *testing.T) {
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()

	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	if err := os.MkdirAll("prompts/skills", 0755); err != nil {
		t.Fatalf("Failed to create prompts dir: %v", err)
	}
	if err := os.WriteFile("prompts/skills/test.md", []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create test prompt: %v", err)
	}

	cmd := initCmd()
	if err := cmd.Flags().Set("guided", "true"); err != nil {
		t.Fatalf("Failed to set guided flag: %v", err)
	}
	if err := cmd.Flags().Set("auto", "true"); err != nil {
		t.Fatalf("Failed to set auto flag: %v", err)
	}

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("initCmd() with guided alias failed: %v", err)
	}

	if _, err := os.Stat(".claude"); os.IsNotExist(err) {
		t.Error("initCmd() with guided alias did not create .claude directory")
	}
}

// TestInitCmdFlags tests that all flags are properly registered
func TestInitCmdFlags(t *testing.T) {
	cmd := initCmd()

	// Check all flags exist
	expectedFlags := []string{
		"project-type", "name", "skip-beads", "skills",
		"auto", "headless", "guided", "interactive", "output",
		"force", "dry-run", "no-evidence",
	}

	for _, flag := range expectedFlags {
		if cmd.Flags().Lookup(flag) == nil {
			t.Errorf("Flag %s not found", flag)
		}
	}
}
