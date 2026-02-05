package tdd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Runner executes TDD phases for different programming languages
type Runner struct {
	language    Language
	testCmd     string
	projectRoot string
}

// PhaseResult represents the result of running a TDD phase
type PhaseResult struct {
	Phase    Phase
	Success  bool
	Duration time.Duration
	Stdout   string
	Stderr   string
	Error    error
}

// RunPhase executes a single TDD phase
func (r *Runner) RunPhase(ctx context.Context, phase Phase, wsPath string) (*PhaseResult, error) {
	start := time.Now()

	// Build command based on language
	cmd := r.buildTestCommand(wsPath)

	// Set working directory to project root if set
	if r.projectRoot != "" {
		cmd.Dir = r.projectRoot
	}

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Start the command
	if err := cmd.Start(); err != nil {
		return &PhaseResult{
			Phase:    phase,
			Success:  false,
			Duration: time.Since(start),
			Error:    fmt.Errorf("failed to start command: %w", err),
		}, err
	}

	// Wait for command to complete or context cancellation
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		// Context cancelled - kill the process
		cmd.Process.Kill()
		return &PhaseResult{
			Phase:    phase,
			Success:  false,
			Duration: time.Since(start),
			Error:    ctx.Err(),
		}, ctx.Err()
	case err := <-done:
		// Command completed
		result := &PhaseResult{
			Phase:    phase,
			Success:  err == nil,
			Duration: time.Since(start),
			Stdout:   stdout.String(),
			Stderr:   stderr.String(),
			Error:    err,
		}

		// Validate result based on phase expectations
		if phase == Red {
			// Red phase expects failure
			if err == nil {
				return result, fmt.Errorf("Red phase expected failure but tests passed")
			}
			return result, nil
		}

		// Green and Refactor phases expect success
		if err != nil {
			return result, fmt.Errorf("phase %s failed: %w", phase, err)
		}

		return result, nil
	}
}

// RunAllPhases executes all TDD phases in sequence
func (r *Runner) RunAllPhases(ctx context.Context, wsPath string) ([]*PhaseResult, error) {
	phases := []Phase{Red, Green, Refactor}
	results := make([]*PhaseResult, 0, len(phases))

	for _, phase := range phases {
		result, err := r.RunPhase(ctx, phase, wsPath)
		if err != nil && phase != Red {
			// Green and Refactor must succeed
			return results, fmt.Errorf("phase %s failed: %w", phase, err)
		}
		results = append(results, result)
	}

	return results, nil
}

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
		return exec.Command(r.testCmd, wsPath)
	}
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
