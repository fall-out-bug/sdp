package tdd

import (
	"fmt"
	"os"
	"path/filepath"
)

// Language represents a programming language
type Language int

const (
	// Python language with pytest
	Python Language = iota
	// Go language with go test
	Go
	// Java language with mvn test
	Java
)

// String returns the string representation of the language
func (l Language) String() string {
	switch l {
	case Python:
		return "Python"
	case Go:
		return "Go"
	case Java:
		return "Java"
	default:
		return "Unknown"
	}
}

// DetectLanguage determines the project language by checking for project files
func (r *Runner) DetectLanguage(projectPath string) (Language, error) {
	// Check for Go project (go.mod)
	if _, err := os.Stat(filepath.Join(projectPath, "go.mod")); err == nil {
		return Go, nil
	}

	// Check for Python project (pyproject.toml)
	if _, err := os.Stat(filepath.Join(projectPath, "pyproject.toml")); err == nil {
		return Python, nil
	}

	// Check for Java project (pom.xml)
	if _, err := os.Stat(filepath.Join(projectPath, "pom.xml")); err == nil {
		return Java, nil
	}

	return -1, fmt.Errorf("unable to detect language: no project files found (go.mod, pyproject.toml, or pom.xml)")
}

// getTestCommand returns the appropriate test command for the language
func (r *Runner) getTestCommand() string {
	switch r.language {
	case Python:
		return "pytest"
	case Go:
		return "go test"
	case Java:
		return "mvn test"
	default:
		return ""
	}
}
