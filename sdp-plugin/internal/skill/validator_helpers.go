package skill

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

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
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close file %s: %v\n", filePath, cerr)
		}
	}()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
	}

	return count, scanner.Err()
}
