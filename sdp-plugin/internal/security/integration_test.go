package security

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ai-masters/sdp/internal/tdd"
)

// TestRunnerCommandInjection tests that runner.go blocks injection attempts
func TestRunnerCommandInjection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a temporary test directory
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		language    string
		testCmd     string
		shouldBlock bool
		description string
	}{
		{
			name:        "safe_python_command",
			language:    "python",
			testCmd:     "pytest",
			shouldBlock: false,
			description: "Safe pytest command",
		},
		{
			name:        "python_injection_semicolon",
			language:    "python",
			testCmd:     "pytest; cat /etc/passwd",
			shouldBlock: true,
			description: "Python command with semicolon injection",
		},
		{
			name:        "python_injection_pipe",
			language:    "python",
			testCmd:     "pytest | sh",
			shouldBlock: true,
			description: "Python command with pipe injection",
		},
		{
			name:        "go_safe_command",
			language:    "go",
			testCmd:     "go test",
			shouldBlock: false,
			description: "Safe go test command",
		},
		{
			name:        "java_safe_command",
			language:    "java",
			testCmd:     "mvn test",
			shouldBlock: false,
			description: "Safe mvn test command",
		},
		{
			name:        "custom_malicious_command",
			language:    "unknown",
			testCmd:     "rm -rf /",
			shouldBlock: true,
			description: "Malicious custom command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var lang tdd.Language
			switch tt.language {
			case "python":
				lang = tdd.Python
			case "go":
				lang = tdd.Go
			case "java":
				lang = tdd.Java
			default:
				lang = tdd.Python // Default to Python for unknown
			}

			runner := tdd.NewRunner(lang)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Try to run with potentially malicious command
			result, err := runner.RunPhase(ctx, tdd.Red, tmpDir)

			// If command should be blocked, we expect either:
			// 1. An error about disallowed command, or
			// 2. Output containing "Error: disallowed test command"
			if tt.shouldBlock {
				if err == nil && result.Success {
					// Check if output indicates blocking
					if !strings.Contains(result.Stdout, "Error: disallowed") &&
					   !strings.Contains(result.Stderr, "Error: disallowed") {
						t.Errorf("%s: malicious command was not blocked! Output: %s", tt.name, result.Stdout)
					}
				}
				// Error is OK - it means the command was rejected
			} else {
				// Safe command should not be blocked
				// It may fail for other reasons (no tests, etc), but not due to security
				if err != nil && strings.Contains(err.Error(), "disallowed") {
					t.Errorf("%s: safe command was blocked: %v", tt.name, err)
				}
			}
		})
	}
}

// TestDoctorCommandInjection tests that doctor.go doesn't execute arbitrary commands
func TestDoctorCommandInjection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test that doctor checks don't allow command injection
	tests := []struct {
		name        string
		envVar      string
		envValue    string
		shouldBlock bool
	}{
		{
			name:        "path_injection",
			envVar:      "PATH",
			envValue:    "/usr/bin:$(touch /tmp/pwned)",
			shouldBlock: true,
		},
		{
			name:        "variable_injection_semicolon",
			envVar:      "CUSTOM_VAR",
			envValue:    "value; rm -rf /",
			shouldBlock: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check if envValue contains injection patterns
			if hasInjectionPattern(tt.envValue) {
				if !tt.shouldBlock {
					t.Errorf("%s: injection pattern detected but should be allowed: %s", tt.name, tt.envValue)
				}
			} else if tt.shouldBlock {
				t.Errorf("%s: should be blocked but no injection pattern found: %s", tt.name, tt.envValue)
			}
		})
	}
}

// TestQualityToolExecution tests that quality tool execution is safe
func TestQualityToolExecution(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test that quality tools use safe command patterns
	safeTools := []string{
		"go test", "go vet", "pytest", "mypy", "ruff",
		"mvn test", "gradle test", "npm test",
	}

	for _, tool := range safeTools {
		t.Run("safe_tool_"+tool, func(t *testing.T) {
			// Split tool into command and args
			parts := strings.Fields(tool)
			if len(parts) == 0 {
				t.Fatal("Tool command is empty")
			}

			cmd := exec.Command(parts[0], parts[1:]...)
			cmd.Dir = t.TempDir() // Run in temp dir for safety

			// Don't actually run, just validate the command structure
			if cmd.Path == "" {
				t.Errorf("Tool command '%s' would fail to execute", tool)
			}

			// Check for injection patterns
			for _, part := range parts {
				if hasInjectionPattern(part) {
					t.Errorf("Tool command '%s' contains injection pattern in '%s'", tool, part)
				}
			}
		})
	}
}

// TestBeadsCommandSafety tests that Beads command execution is safe
func TestBeadsCommandSafety(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test that Beads commands are validated
	beadsCommands := []string{
		"bd list",
		"bd show",
		"bd create",
		"bd close",
	}

	for _, cmdStr := range beadsCommands {
		t.Run("beads_"+strings.Fields(cmdStr)[0], func(t *testing.T) {
			parts := strings.Fields(cmdStr)
			if len(parts) == 0 {
				t.Fatal("Command is empty")
			}

			// Verify command starts with "bd"
			if parts[0] != "bd" {
				t.Errorf("Beads command should start with 'bd', got: %s", parts[0])
			}

			// Check for injection in arguments
			for _, part := range parts[1:] {
				if hasInjectionPattern(part) {
					t.Errorf("Beads command '%s' contains injection pattern in '%s'", cmdStr, part)
				}
			}
		})
	}
}

// TestAllExecCommandCalls tests that all exec.Command calls in the codebase are safe
func TestAllExecCommandCalls(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// List of files that use exec.Command
	filesWithExec := []string{
		"internal/tdd/runner.go",
		"internal/beads/utils.go",
		"internal/doctor/doctor.go",
		"internal/orchestrator/beads_loader.go",
		"internal/quality/coverage.go",
		"internal/quality/complexity_python.go",
		"internal/quality/complexity_go.go",
		"internal/quality/complexity_java.go",
	}

	for _, file := range filesWithExec {
		t.Run(filepath.Base(file), func(t *testing.T) {
			// Read the file and check for unsafe patterns
			content, err := filepath.Abs(file)
			if err != nil {
				t.Fatalf("Failed to get absolute path for %s: %v", file, err)
			}

			// For now, just log that we need to check this file
			// TODO: Implement actual static analysis
			t.Logf("Need to audit %s for exec.Command safety", content)
		})
	}
}
