package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSkillValidateCmd tests the skill validate command
func TestSkillValidateCmd(t *testing.T) {
	// Create temp directory with skill file
	tmpDir := t.TempDir()
	skillsDir := filepath.Join(tmpDir, ".claude", "skills")
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		t.Fatalf("Failed to create skills dir: %v", err)
	}

	// Create a valid skill file
	validSkill := filepath.Join(skillsDir, "test.md")
	skillContent := `---
name: test
description: Test skill file
---

# Test Skill

This is a test skill.

## Usage

Use this skill for testing.

## Quick Reference

- Test command
- Test option

## Workflow

1. Step 1
2. Step 2

## See Also

- Other skill
`
	if err := os.WriteFile(validSkill, []byte(skillContent), 0644); err != nil {
		t.Fatalf("Failed to create skill file: %v", err)
	}

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	cmd := skillValidate()

	tests := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "no args",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "valid skill file",
			args:        []string{validSkill},
			expectError: false,
		},
		{
			name:        "non-existent file",
			args:        []string{"nonexistent.md"},
			expectError: true,
			errorMsg:    "validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cmd.RunE(cmd, tt.args)

			if tt.expectError && err == nil {
				t.Errorf("skillValidate() expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("skillValidate() unexpected error: %v", err)
			}
			if tt.expectError && tt.errorMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("skillValidate() error = %q, should contain %q", err.Error(), tt.errorMsg)
				}
			}
		})
	}
}

// TestSkillValidateWithStrictFlag tests strict mode
func TestSkillValidateWithStrictFlag(t *testing.T) {
	// Create temp directory with skill file that has warnings
	tmpDir := t.TempDir()
	skillsDir := filepath.Join(tmpDir, ".claude", "skills")
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		t.Fatalf("Failed to create skills dir: %v", err)
	}

	// Create a skill file with long content (>100 lines to trigger warning)
	longSkill := filepath.Join(skillsDir, "long.md")
	var lines []string
	lines = append(lines, "# Long Skill")
	lines = append(lines, "")
	for i := 0; i < 110; i++ {
		lines = append(lines, "Line "+string(rune('0'+i%10)))
	}
	skillContent := strings.Join(lines, "\n")
	if err := os.WriteFile(longSkill, []byte(skillContent), 0644); err != nil {
		t.Fatalf("Failed to create skill file: %v", err)
	}

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	cmd := skillValidate()
	if err := cmd.Flags().Set("strict", "true"); err != nil {
		t.Fatalf("Failed to set strict flag: %v", err)
	}

	err := cmd.RunE(cmd, []string{longSkill})
	// Should error in strict mode with warnings
	if err == nil {
		t.Error("skillValidate() with strict flag expected error for warnings")
	}
}

// TestSkillCheckAllCmd tests the check-all command
func TestSkillCheckAllCmd(t *testing.T) {
	// Create temp directory with skills
	tmpDir := t.TempDir()
	skillsDir := filepath.Join(tmpDir, ".claude", "skills")
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		t.Fatalf("Failed to create skills dir: %v", err)
	}

	// Create valid skill files
	for _, name := range []string{"test1.md", "test2.md"} {
		skillPath := filepath.Join(skillsDir, name)
		content := "# Test Skill\n\nThis is a test skill."
		if err := os.WriteFile(skillPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create skill file: %v", err)
		}
	}

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	cmd := skillCheckAll()
	if err := cmd.Flags().Set("skills-dir", skillsDir); err != nil {
		t.Fatalf("Failed to set skills-dir flag: %v", err)
	}

	err := cmd.RunE(cmd, []string{})
	// Should succeed with all valid skills
	if err != nil {
		t.Errorf("skillCheckAll() failed: %v", err)
	}
}

// TestSkillListCmd tests the skill list command
func TestSkillListCmd(t *testing.T) {
	cmd := skillList()

	// Test command structure
	if cmd.Use != "list" {
		t.Errorf("skillList() has wrong use: %s", cmd.Use)
	}
}

// TestSkillShowCmd tests the skill show command
func TestSkillShowCmd(t *testing.T) {
	cmd := skillShow()

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "no args",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "non-existent skill",
			args:        []string{"nonexistent"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cmd.RunE(cmd, tt.args)

			if tt.expectError && err == nil {
				t.Errorf("skillShow() expected error but got none")
			}
		})
	}
}
