package tdd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// buildTestCommand constructs the test command for the current language
func (r *Runner) buildTestCommand(wsPath string) *exec.Cmd {
	switch r.language {
	case Python:
		// Python: pytest
		return exec.Command("pytest", wsPath, "-v")
	case Go:
		// Go: go test - use package path or ./ for current dir
		cmd := exec.Command("go", "test", wsPath)
		return cmd
	case Java:
		// Java: mvn test
		return exec.Command("mvn", "test", "-f", wsPath)
	default:
		// Validate testCmd against whitelist before using
		if !isAllowedTestCommand(r.testCmd) {
			return exec.Command("echo", fmt.Sprintf("Error: disallowed test command '%s'", r.testCmd))
		}
		return exec.Command(r.testCmd, wsPath)
	}
}

// isAllowedTestCommand validates testCmd against a whitelist of safe commands
func isAllowedTestCommand(testCmd string) bool {
	// Whitelist of allowed test commands
	// Only allow specific, known-safe test runners
	allowedCommands := []string{
		"pytest",
		"pytest-3",
		"python -m pytest",
		"go test",
		"mvn test",
		"mvnw test",
		"gradle test",
		"./gradlew test",
		"gradlew test",
		"npm test",
		"yarn test",
		"pnpm test",
		"jest",
		"jasmine",
		"mocha",
		"cargo test",
		"dart test",
		"flutter test",
	}

	for _, allowed := range allowedCommands {
		if testCmd == allowed {
			return true
		}
	}
	return false
}

// NewRunner creates a new Runner for the specified language
func NewRunner(language Language) *Runner {
	testCmd := ""
	switch language {
	case Python:
		testCmd = "pytest"
	case Go:
		testCmd = "go test"
	case Java:
		testCmd = "mvn test"
	}

	// Find project root by looking for go.mod, pyproject.toml, or pom.xml
	// Start from current directory and go up until we find a project file
	projectRoot, err := findProjectRoot(".")
	if err != nil {
		// If not found, use current directory
		projectRoot = "."
	}

	return &Runner{
		language:    language,
		testCmd:     testCmd,
		projectRoot: projectRoot,
	}
}

// findProjectRoot finds the project root by searching for project files
func findProjectRoot(startPath string) (string, error) {
	// Get absolute path
	absPath, err := filepath.Abs(startPath)
	if err != nil {
		return "", err
	}

	currentPath := absPath

	for {
		// Check if any project file exists in current directory
		for _, file := range []string{"go.mod", "pyproject.toml", "pom.xml"} {
			if _, err := os.Stat(filepath.Join(currentPath, file)); err == nil {
				return currentPath, nil
			}
		}

		// Move to parent directory
		parent := filepath.Dir(currentPath)
		if parent == currentPath {
			// Reached root without finding project file
			return "", fmt.Errorf("no project file found")
		}
		currentPath = parent
	}
}
