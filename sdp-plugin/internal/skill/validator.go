package skill

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Required sections in skill files
var requiredSections = []string{
	"## Quick Reference",
	"## Workflow",
	"## See Also",
}

// ValidationResult represents validation results
type ValidationResult struct {
	IsValid  bool     `json:"is_valid"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
	LineCount int     `json:"line_count"`
}

// Validator checks skill files against standards
type Validator struct {
	maxLines      int
	warningLines  int
}

// NewValidator creates a new skill validator
func NewValidator() *Validator {
	return &Validator{
		maxLines:     150,
		warningLines: 100,
	}
}

// ValidateFile checks a skill file against standards
func (v *Validator) ValidateFile(skillPath string) (*ValidationResult, error) {
	// Read file
	content, err := os.ReadFile(skillPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	lineCount := len(lines)

	result := &ValidationResult{
		IsValid:   true,
		Errors:    []string{},
		Warnings:  []string{},
		LineCount: lineCount,
	}

	// Check line count
	if lineCount > v.maxLines {
		result.Errors = append(result.Errors, fmt.Sprintf("Too long: %d lines (max %d)", lineCount, v.maxLines))
		result.IsValid = false
	} else if lineCount > v.warningLines {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Consider shortening: %d lines (target %d)", lineCount, v.warningLines))
	}

	// Check required sections
	for _, section := range requiredSections {
		if !strings.Contains(contentStr, section) {
			result.Errors = append(result.Errors, fmt.Sprintf("Missing section: %s", section))
			result.IsValid = false
		}
	}

	// Check frontmatter
	if !strings.HasPrefix(contentStr, "---") {
		result.Errors = append(result.Errors, "Missing frontmatter (must start with ---)")
		result.IsValid = false
	}

	// Check references
	refWarnings := v.checkReferences(contentStr, skillPath)
	result.Warnings = append(result.Warnings, refWarnings...)

	return result, nil
}

// checkReferences validates that file references exist
func (v *Validator) checkReferences(content, skillPath string) []string {
	warnings := []string{}

	// Extract markdown references: [text](path)
	re := regexp.MustCompile(`\[.*?\]\((\.\.?/[^)]+)\)`)
	matches := re.FindAllStringSubmatch(content, -1)

	skillDir := filepath.Dir(skillPath)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		ref := match[1]

		// Skip absolute URLs and docs references
		if strings.HasPrefix(ref, "http") || strings.HasPrefix(ref, "../docs/") {
			continue
		}

		// Check if reference exists
		refPath := filepath.Join(skillDir, ref)
		if _, err := os.Stat(refPath); os.IsNotExist(err) {
			warnings = append(warnings, fmt.Sprintf("Reference may not exist: %s", ref))
		}
	}

	return warnings
}

// ValidateAll validates all skills in .claude/skills/
func (v *Validator) ValidateAll(skillsDir string) (map[string]*ValidationResult, error) {
	results := make(map[string]*ValidationResult)

	// Check if directory exists
	if _, err := os.Stat(skillsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("skills directory not found: %s", skillsDir)
	}

	// Walk through skill directories
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read skills directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillFile := filepath.Join(skillsDir, entry.Name(), "SKILL.md")
		if _, err := os.Stat(skillFile); err != nil {
			continue // Skip if SKILL.md doesn't exist
		}

		result, err := v.ValidateFile(skillFile)
		if err != nil {
			return nil, fmt.Errorf("failed to validate %s: %w", entry.Name(), err)
		}

		results[entry.Name()] = result
	}

	return results, nil
}

// ListSkills returns all skill directories
func ListSkills(skillsDir string) ([]string, error) {
	var skills []string

	// Check if directory exists
	if _, err := os.Stat(skillsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("skills directory not found: %s", skillsDir)
	}

	// Read skill directories
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read skills directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			skills = append(skills, entry.Name())
		}
	}

	return skills, nil
}

// ReadSkillContent returns the content of a skill file
func ReadSkillContent(skillsDir, skillName string) (string, error) {
	skillFile := filepath.Join(skillsDir, skillName, "SKILL.md")

	content, err := os.ReadFile(skillFile)
	if err != nil {
		return "", fmt.Errorf("failed to read skill file: %w", err)
	}

	return string(content), nil
}

// CountLines counts lines in a file
func CountLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
	}

	return count, scanner.Err()
}
