package quality

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (c *Checker) CheckFileSize() (*FileSizeResult, error) {
	result := &FileSizeResult{
		Threshold: 200,
	}

	var totalLOC int
	var totalFiles int
	var sumLOC int

	// Walk the project directory
	err := filepath.Walk(c.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		// Skip certain directories
		if strings.Contains(path, "vendor/") ||
		   strings.Contains(path, "node_modules/") ||
		   strings.Contains(path, ".git/") ||
		   strings.Contains(path, "target/") ||
		   strings.Contains(path, "__pycache__/") ||
		   strings.Contains(path, ".venv/") ||
		   strings.Contains(path, "venv/") {
			return nil
		}

		// Check file extensions based on project type
		var shouldCheck bool
		switch c.projectType {
		case Python:
			shouldCheck = strings.HasSuffix(path, ".py")
		case Go:
			shouldCheck = strings.HasSuffix(path, ".go")
		case Java:
			shouldCheck = strings.HasSuffix(path, ".java")
		}

		if !shouldCheck {
			return nil
		}

		// Skip test files for size check (they can be longer)
		if strings.Contains(filepath.Base(path), "test") ||
		   strings.Contains(filepath.Base(path), "_test") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		lines := strings.Split(string(content), "\n")
		loc := len(lines)

		totalFiles++
		totalLOC += loc
		sumLOC += loc

		if loc > result.Threshold {
			// Make path relative to project path
			relPath, _ := filepath.Rel(c.projectPath, path)
			result.Violators = append(result.Violators, FileViolation{
				File: relPath,
				LOC:  loc,
			})
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking project: %w", err)
	}

	if totalFiles > 0 {
		result.AverageLOC = totalLOC / totalFiles
	}

	result.TotalFiles = totalFiles
	result.Passed = len(result.Violators) == 0

	return result, nil
}
