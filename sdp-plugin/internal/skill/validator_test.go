package skill

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewValidator(t *testing.T) {
	validator := NewValidator()

	if validator == nil {
		t.Fatal("NewValidator returned nil")
	}

	if validator.maxLines != 150 {
		t.Errorf("maxLines = %d, want 150", validator.maxLines)
	}

	if validator.warningLines != 100 {
		t.Errorf("warningLines = %d, want 100", validator.warningLines)
	}
}

func TestValidateFileValid(t *testing.T) {
	tmpDir := t.TempDir()
	validator := NewValidator()

	// Create valid skill file
	content := `---
name: test
---

## Quick Reference
Test content

## Workflow
Test workflow

## See Also
Test refs
`

	skillPath := filepath.Join(tmpDir, "SKILL.md")
	os.WriteFile(skillPath, []byte(content), 0644)

	result, err := validator.ValidateFile(skillPath)
	if err != nil {
		t.Fatalf("ValidateFile failed: %v", err)
	}

	if !result.IsValid {
		t.Errorf("IsValid = false, want true (errors: %v)", result.Errors)
	}

	if len(result.Errors) > 0 {
		t.Errorf("Errors = %v, want empty", result.Errors)
	}
}

func TestValidateFileTooLong(t *testing.T) {
	tmpDir := t.TempDir()
	validator := NewValidator()

	// Create skill file with too many lines
	content := "---\n"
	for i := 0; i < 160; i++ {
		content += "line\n"
	}

	skillPath := filepath.Join(tmpDir, "SKILL.md")
	os.WriteFile(skillPath, []byte(content), 0644)

	result, err := validator.ValidateFile(skillPath)
	if err != nil {
		t.Fatalf("ValidateFile failed: %v", err)
	}

	if result.IsValid {
		t.Error("IsValid should be false for file with >150 lines")
	}

	if len(result.Errors) == 0 {
		t.Error("Should have error for too long file")
	}
}

func TestValidateFileWarningLength(t *testing.T) {
	tmpDir := t.TempDir()
	validator := NewValidator()

	// Create skill file with warning length
	content := "---\n## Quick Reference\n## Workflow\n## See Also\n"
	for i := 0; i < 110; i++ {
		content += "line\n"
	}

	skillPath := filepath.Join(tmpDir, "SKILL.md")
	os.WriteFile(skillPath, []byte(content), 0644)

	result, err := validator.ValidateFile(skillPath)
	if err != nil {
		t.Fatalf("ValidateFile failed: %v", err)
	}

	if len(result.Warnings) == 0 {
		t.Error("Should have warning for file with >100 lines")
	}
}

func TestValidateFileMissingSection(t *testing.T) {
	tmpDir := t.TempDir()
	validator := NewValidator()

	// Create skill file missing required sections
	content := `---
name: test
---

## Quick Reference
Test
`

	skillPath := filepath.Join(tmpDir, "SKILL.md")
	os.WriteFile(skillPath, []byte(content), 0644)

	result, err := validator.ValidateFile(skillPath)
	if err != nil {
		t.Fatalf("ValidateFile failed: %v", err)
	}

	if result.IsValid {
		t.Error("IsValid should be false when sections missing")
	}

	if len(result.Errors) < 2 {
		t.Errorf("Should have at least 2 errors (missing sections), got %d", len(result.Errors))
	}
}

func TestValidateFileNoFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	validator := NewValidator()

	// Create skill file without frontmatter
	content := `## Quick Reference
## Workflow
## See Also
`

	skillPath := filepath.Join(tmpDir, "SKILL.md")
	os.WriteFile(skillPath, []byte(content), 0644)

	result, err := validator.ValidateFile(skillPath)
	if err != nil {
		t.Fatalf("ValidateFile failed: %v", err)
	}

	if result.IsValid {
		t.Error("IsValid should be false without frontmatter")
	}

	found := false
	for _, err := range result.Errors {
		if err == "Missing frontmatter (must start with ---)" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Should have error for missing frontmatter")
	}
}

func TestCheckReferences(t *testing.T) {
	tmpDir := t.TempDir()
	validator := NewValidator()

	// Create skill file with broken reference
	content := `---
name: test
---

## Quick Reference
See [missing](./missing.md)

## Workflow
Test

## See Also
Test
`

	skillPath := filepath.Join(tmpDir, "SKILL.md")
	os.WriteFile(skillPath, []byte(content), 0644)

	result, err := validator.ValidateFile(skillPath)
	if err != nil {
		t.Fatalf("ValidateFile failed: %v", err)
	}

	if len(result.Warnings) == 0 {
		t.Error("Should have warning for broken reference")
	}
}

func TestValidateAll(t *testing.T) {
	tmpDir := t.TempDir()
	validator := NewValidator()

	// Create skills directory
	skillsDir := filepath.Join(tmpDir, "skills")
	os.MkdirAll(skillsDir, 0755)

	// Create valid skill
	validSkillDir := filepath.Join(skillsDir, "valid")
	os.MkdirAll(validSkillDir, 0755)
	validContent := `---
name: valid
---

## Quick Reference
Test

## Workflow
Test

## See Also
Test
`
	os.WriteFile(filepath.Join(validSkillDir, "SKILL.md"), []byte(validContent), 0644)

	// Create invalid skill
	invalidSkillDir := filepath.Join(skillsDir, "invalid")
	os.MkdirAll(invalidSkillDir, 0755)
	os.WriteFile(filepath.Join(invalidSkillDir, "SKILL.md"), []byte("invalid"), 0644)

	results, err := validator.ValidateAll(skillsDir)
	if err != nil {
		t.Fatalf("ValidateAll failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Results count = %d, want 2", len(results))
	}

	if !results["valid"].IsValid {
		t.Error("valid skill should be valid")
	}

	if results["invalid"].IsValid {
		t.Error("invalid skill should be invalid")
	}
}

func TestValidateAllNoDirectory(t *testing.T) {
	validator := NewValidator()

	_, err := validator.ValidateAll("/nonexistent/directory")
	if err == nil {
		t.Error("ValidateAll should return error for non-existent directory")
	}
}

func TestListSkills(t *testing.T) {
	tmpDir := t.TempDir()

	// Create skills directory
	skillsDir := filepath.Join(tmpDir, "skills")
	os.MkdirAll(skillsDir, 0755)

	// Create skill directories
	os.MkdirAll(filepath.Join(skillsDir, "skill1"), 0755)
	os.MkdirAll(filepath.Join(skillsDir, "skill2"), 0755)

	// Create a file (not a directory)
	os.WriteFile(filepath.Join(skillsDir, "file.txt"), []byte("test"), 0644)

	skills, err := ListSkills(skillsDir)
	if err != nil {
		t.Fatalf("ListSkills failed: %v", err)
	}

	if len(skills) != 2 {
		t.Errorf("Skills count = %d, want 2", len(skills))
	}
}

func TestListSkillsNoDirectory(t *testing.T) {
	_, err := ListSkills("/nonexistent/directory")
	if err == nil {
		t.Error("ListSkills should return error for non-existent directory")
	}
}

func TestReadSkillContent(t *testing.T) {
	tmpDir := t.TempDir()

	// Create skill directory and file
	skillsDir := filepath.Join(tmpDir, "skills")
	skillDir := filepath.Join(skillsDir, "test")
	os.MkdirAll(skillDir, 0755)

	expectedContent := "Test content\nLine 2"
	skillPath := filepath.Join(skillDir, "SKILL.md")
	os.WriteFile(skillPath, []byte(expectedContent), 0644)

	content, err := ReadSkillContent(skillsDir, "test")
	if err != nil {
		t.Fatalf("ReadSkillContent failed: %v", err)
	}

	if content != expectedContent {
		t.Errorf("Content = %s, want %s", content, expectedContent)
	}
}

func TestReadSkillContentNotExists(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := ReadSkillContent(tmpDir, "nonexistent")
	if err == nil {
		t.Error("ReadSkillContent should return error for non-existent skill")
	}
}

func TestCountLines(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test file
	filePath := filepath.Join(tmpDir, "test.txt")
	content := "line1\nline2\nline3\n"
	os.WriteFile(filePath, []byte(content), 0644)

	count, err := CountLines(filePath)
	if err != nil {
		t.Fatalf("CountLines failed: %v", err)
	}

	if count != 3 {
		t.Errorf("Count = %d, want 3", count)
	}
}

func TestCountLinesNotExists(t *testing.T) {
	_, err := CountLines("/nonexistent/file.txt")
	if err == nil {
		t.Error("CountLines should return error for non-existent file")
	}
}
