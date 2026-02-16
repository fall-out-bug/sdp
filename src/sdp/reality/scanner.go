package reality

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ScanResult represents the analysis of a project
type ScanResult struct {
	Language      string
	Framework     string
	TotalFiles    int
	LinesOfCode   int
	Directories   []string
	TestFiles     int
}

// ProjectScanner analyzes project structure
type ProjectScanner struct {
	projectPath string
}

// NewProjectScanner creates a new scanner
func NewProjectScanner(projectPath string) *ProjectScanner {
	return &ProjectScanner{
		projectPath: projectPath,
	}
}

// Scan analyzes the project
func (s *ProjectScanner) Scan() (*ScanResult, error) {
	result := &ScanResult{
		Directories: []string{},
	}

	// Detect language and framework
	lang, fw := DetectLanguage(s.projectPath)
	result.Language = lang
	result.Framework = fw

	// Walk directory tree
	err := filepath.Walk(s.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and directories
		if strings.HasPrefix(filepath.Base(path), ".") && path != s.projectPath {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip .git directory
		if filepath.Base(path) == ".git" {
			return filepath.SkipDir
		}

		// Skip node_modules, vendor, __pycache__
		if info.IsDir() {
			basename := filepath.Base(path)
			if basename == "node_modules" || basename == "vendor" || basename == "__pycache__" {
				return filepath.SkipDir
			}
			result.Directories = append(result.Directories, path)
			return nil
		}

		// Count files and lines
		result.TotalFiles++

		// Check if it's a test file
		if strings.Contains(filepath.Base(path), "_test.go") ||
		   strings.Contains(filepath.Base(path), ".test.") ||
		   strings.HasPrefix(filepath.Base(path), "test_") {
			result.TestFiles++
		}

		// Count lines of code
		lines, err := countLines(path)
		if err == nil {
			result.LinesOfCode += lines
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan project: %w", err)
	}

	return result, nil
}

// countLines counts non-empty lines in a file
func countLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "#") {
			count++
		}
	}

	return count, scanner.Err()
}
