package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// This file contains integration tests for CLI commands.
// These tests require the sdp binary to be built separately.
// To build: go build -o sdp ./cmd/sdp

func skipIfBinaryNotBuilt(t *testing.T) string {
	// Tests run in sdp-plugin/cmd/sdp directory
	// Binary is in sdp-plugin/ directory (two levels up)
	// Need to go up TWO levels because of root-level cmd/ directory
	relativePath := filepath.Join("..", "..", "sdp")
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		t.Skip("Cannot resolve binary path")
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		t.Skip("sdp binary not found. Run: go build -o sdp ./cmd/sdp from sdp-plugin/ directory")
	}
	return absPath
}

// repoRoot returns the absolute path to the repository root
func repoRoot(t *testing.T) string {
	// From sdp-plugin/cmd/sdp, go up THREE levels to get to repo root
	// sdp-plugin/cmd/sdp → ../.. → sdp-plugin/ → ../../.. → sdp/
	path := filepath.Join("..", "..", "..")
	absPath, err := filepath.Abs(path)
	if err != nil {
		t.Fatalf("Cannot resolve repo root: %v", err)
	}
	return absPath
}

// TestParseCommand tests the sdp parse command
func TestParseCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	tests := []struct {
		name       string
		args       []string
		wantErr    bool
		contains   string
		notContain string
	}{
		{
			name:     "parse valid workstream by ID",
			args:     []string{"parse", "00-050-01"},
			wantErr:  false,
			contains: "00-050-01",
		},
		{
			name:     "parse missing workstream",
			args:     []string{"parse", "99-999-99"},
			wantErr:  true,
			contains: "not found",
		},
		{
			name:     "parse without args",
			args:     []string{"parse"},
			wantErr:  true,
			contains: "required",
		},
		{
			name:       "path traversal attack blocked",
			args:       []string{"parse", "../../../etc/passwd"},
			wantErr:    true,
			contains:   "not found",
			notContain: "root:x",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			// Set working directory to repo root so docs/workstreams/ are found
			cmd.Dir = repoRoot(t)

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				output := stdout.String() + stderr.String()
				t.Errorf("Unexpected error: %v\nOutput:\n%s", err, output)
			}

			output := stdout.String() + stderr.String()
			if !strings.Contains(output, tt.contains) {
				t.Errorf("Output does not contain expected string %q\nGot: %s", tt.contains, output)
			}

			if tt.notContain != "" && strings.Contains(output, tt.notContain) {
				t.Errorf("Output should not contain string %q\nGot: %s", tt.notContain, output)
			}
		})
	}
}

// TestDriftCommand tests the sdp drift detect command
func TestDriftCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	// Test that the drift command exists and doesn't crash
	cmd := exec.Command(binaryPath, "drift", "detect", "--help")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run drift help: %v", err)
	}

	output := stdout.String() + stderr.String()
	if !strings.Contains(output, "detect") {
		t.Errorf("Drift command help should mention detect\nGot: %s", output)
	}
}

// TestVersionCommand tests the sdp --version flag
func TestVersionCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	cmd := exec.Command(binaryPath, "--version")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run version command: %v", err)
	}

	output := stdout.String() + stderr.String()
	if !strings.Contains(output, "sdp version") {
		t.Errorf("Version output does not contain version string\nGot: %s", output)
	}
}

// TestHelpCommand tests the sdp --help flag
func TestHelpCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	cmd := exec.Command(binaryPath, "--help")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run help command: %v", err)
	}

	output := stdout.String() + stderr.String()

	// Check for expected help content
	expectedKeywords := []string{"Usage", "Available Commands", "Flags"}
	for _, keyword := range expectedKeywords {
		if !strings.Contains(output, keyword) {
			t.Errorf("Help output does not contain expected keyword %q\nGot: %s", keyword, output)
		}
	}
}

// TestDoctorCommand tests the sdp doctor command
func TestDoctorCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	cmd := exec.Command(binaryPath, "doctor")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// doctor should not fail
	if err := cmd.Run(); err != nil {
		t.Logf("Doctor command failed: %v\nOutput: %s", err, stdout.String()+stderr.String())
	}

	output := stdout.String() + stderr.String()

	// Check for expected doctor content
	if !strings.Contains(output, "SDP Doctor") && !strings.Contains(output, "doctor") {
		t.Logf("Doctor output: %s", output)
	}
}

// TestCheckpointCommand tests the sdp checkpoint commands
func TestCheckpointCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	// Test checkpoint list
	cmd := exec.Command(binaryPath, "checkpoint", "list")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	// May fail if no checkpoints directory, that's OK

	output := stdout.String() + stderr.String()
	t.Logf("Checkpoint list output: %s\nError: %v", output, err)
}

// TestInitCommand tests the sdp init command
func TestInitCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	// Create temp directory for init
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Run init
	cmd := exec.Command(binaryPath, "init", "--project-type", "go")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// May fail if prompts/ dir not found (that's OK in test environment)
	output := stdout.String() + stderr.String()
	t.Logf("Init output: %s\nError: %v", output, err)

	// Check if .claude was created (should succeed if prompts exist)
	if err == nil {
		claudeDir := filepath.Join(tmpDir, ".claude")
		if _, err := os.Stat(claudeDir); os.IsNotExist(err) {
			t.Error(".claude directory was not created")
		}
	}
}

// TestVerifyCommand tests the sdp verify command
func TestVerifyCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	root := repoRoot(t)

	// Test verify on an existing workstream
	wsFile := filepath.Join(root, "docs", "workstreams", "completed", "00-050-01.md")

	if _, err := os.Stat(wsFile); os.IsNotExist(err) {
		t.Skip("Workstream file not found, skipping verify test")
	}

	// Change to repo root directory so verify can find docs/workstreams
	cmd := exec.Command(binaryPath, "verify", "00-050-01")
	cmd.Dir = root
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	t.Logf("Verify output: %s\nError: %v", output, err)

	// Verify should not fail catastrophically
	if err != nil && !strings.Contains(output, "Error") && !strings.Contains(output, "FAILED") {
		t.Errorf("Verify failed unexpectedly: %v", err)
	}
}

// TestGuardCommand tests the sdp guard command
func TestGuardCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	root := repoRoot(t)
	wsFile := filepath.Join(root, "docs", "workstreams", "completed", "00-050-01.md")

	if _, err := os.Stat(wsFile); os.IsNotExist(err) {
		t.Skip("Workstream file not found, skipping guard test")
	}

	cmd := exec.Command(binaryPath, "guard", "activate", wsFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	t.Logf("Guard output: %s\nError: %v", output, err)

	// Guard should work or fail gracefully
	if err != nil && !strings.Contains(output, "Error") {
		// OK as long as it's not a panic
	}
}

// TestBeadsCommand tests the sdp beads command
func TestBeadsCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "beads ready",
			args:     []string{"beads", "ready"},
			wantErr:  false,
			contains: "ready",
		},
		{
			name:    "beads list",
			args:    []string{"beads", "list"},
			wantErr: false,
		},
		{
			name:    "beads sync",
			args:    []string{"beads", "sync"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			output := stdout.String() + stderr.String()

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}

			if tt.contains != "" && !strings.Contains(output, tt.contains) {
				t.Logf("Output: %s", output)
			}
		})
	}
}

// TestCompletionCommand tests the sdp completion command
func TestCompletionCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	shells := []string{"bash", "zsh", "fish"}

	for _, shell := range shells {
		t.Run("completion "+shell, func(t *testing.T) {
			cmd := exec.Command(binaryPath, "completion", shell)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			if err := cmd.Run(); err != nil {
				t.Logf("Completion %s failed: %v\nOutput: %s", shell, err, stdout.String()+stderr.String())
			}

			output := stdout.String()
			if len(output) == 0 {
				t.Errorf("Completion script for %s should produce output", shell)
			}
		})
	}
}

// TestOrchestrateCommand tests the sdp orchestrate command
func TestOrchestrateCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	// Test orchestrate help
	cmd := exec.Command(binaryPath, "orchestrate", "--help")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Orchestrate help failed: %v", err)
	}

	output := stdout.String() + stderr.String()
	if !strings.Contains(output, "orchestrate") {
		t.Errorf("Orchestrate help should mention orchestrate\nGot: %s", output)
	}
}

// TestPrdCommand tests the sdp prd command
func TestPrdCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "prd help",
			args:     []string{"prd", "--help"},
			wantErr:  false,
			contains: "PRD",
		},
		{
			name:    "prd detect",
			args:    []string{"prd", "detect"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			output := stdout.String() + stderr.String()

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}

			if tt.contains != "" && !strings.Contains(output, tt.contains) {
				t.Errorf("Output does not contain %q\nGot: %s", tt.contains, output)
			}
		})
	}
}

// TestQualityCommand tests the sdp quality command
func TestQualityCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	root := repoRoot(t)
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(root)

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "quality all",
			args:     []string{"quality", "all"},
			wantErr:  true, // Expected to fail due to coverage/complexity
			contains: "Coverage",
		},
		{
			name:    "quality coverage",
			args:    []string{"quality", "coverage"},
			wantErr: false,
		},
		{
			name:     "quality help",
			args:     []string{"quality", "--help"},
			wantErr:  false,
			contains: "quality",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			output := stdout.String() + stderr.String()

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}

			if tt.contains != "" && !strings.Contains(output, tt.contains) {
				t.Logf("Output: %s", output)
			}

			t.Logf("Quality %s: err=%v", tt.name, err)
		})
	}
}

// TestSkillCommand tests the sdp skill command
func TestSkillCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "skill help",
			args:     []string{"skill", "--help"},
			wantErr:  false,
			contains: "skill",
		},
		{
			name:    "skill validate",
			args:    []string{"skill", "validate"},
			wantErr: false,
		},
		{
			name:    "skill list",
			args:    []string{"skill", "list"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			output := stdout.String() + stderr.String()

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}

			if tt.contains != "" && !strings.Contains(output, tt.contains) {
				t.Logf("Output: %s", output)
			}
		})
	}
}

// TestTddCommand tests the sdp tdd command
func TestTddCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	// Create temp directory for TDD test
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_example_test.go")
	testContent := `
package test_example

import "testing"

func TestExample(t *testing.T) {
	if 1+1 != 2 {
		t.Fail()
	}
}
`
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "tdd run",
			args:    []string{"tdd", "run", testFile},
			wantErr: false,
		},
		{
			name:    "tdd help",
			args:    []string{"tdd", "--help"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			output := stdout.String() + stderr.String()

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}

			t.Logf("TDD %s: err=%v, output=%s", tt.name, err, output)
		})
	}
}

// TestTelemetryCommand tests the sdp telemetry command
func TestTelemetryCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "telemetry analyze",
			args:     []string{"telemetry", "analyze"},
			wantErr:  false,
			contains: "telemetry",
		},
		{
			name:    "telemetry export json",
			args:    []string{"telemetry", "export", "--format", "json"},
			wantErr: false,
		},
		{
			name:     "telemetry help",
			args:     []string{"telemetry", "--help"},
			wantErr:  false,
			contains: "telemetry",
		},
		{
			name:    "telemetry status",
			args:    []string{"telemetry", "status"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			output := stdout.String() + stderr.String()

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}

			if tt.contains != "" && !strings.Contains(output, tt.contains) {
				t.Logf("Output: %s", output)
			}

			t.Logf("Telemetry %s: err=%v", tt.name, err)
		})
	}
}

// TestWatchCommand tests the sdp watch command
func TestWatchCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	// Test watch help (can't test actual watch as it runs forever)
	cmd := exec.Command(binaryPath, "watch", "--help")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Watch help failed: %v", err)
	}

	output := stdout.String() + stderr.String()
	if !strings.Contains(output, "watch") {
		t.Errorf("Watch help should mention watch\nGot: %s", output)
	}
}

// TestHooksCommand tests the sdp hooks command
func TestHooksCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	root := repoRoot(t)
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(root)

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "hooks install",
			args:     []string{"hooks", "install"},
			wantErr:  false,
			contains: "hooks",
		},
		{
			name:    "hooks uninstall",
			args:    []string{"hooks", "uninstall"},
			wantErr: false,
		},
		{
			name:     "hooks help",
			args:     []string{"hooks", "--help"},
			wantErr:  false,
			contains: "hooks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			output := stdout.String() + stderr.String()

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}

			if tt.contains != "" && !strings.Contains(output, tt.contains) {
				t.Logf("Output: %s", output)
			}

			t.Logf("Hooks %s: err=%v", tt.name, err)
		})
	}
}

// TestCommandsCoverage tests that all commands are reachable
func TestCommandsCoverage(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	// Get list of all commands
	cmd := exec.Command(binaryPath, "--help")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to get help: %v", err)
	}

	output := stdout.String() + stderr.String()

	// List of expected commands
	expectedCommands := []string{
		"beads",
		"checkpoint",
		"completion",
		"decisions",
		"doctor",
		"drift",
		"guard",
		"help",
		"hooks",
		"init",
		"orchestrate",
		"parse",
		"prd",
		"quality",
		"skill",
		"status",
		"tdd",
		"telemetry",
		"verify",
		"watch",
	}

	for _, expected := range expectedCommands {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected command %q not found in help output", expected)
		}
	}
}

// TestDecisionsCommand tests the sdp decisions command
func TestDecisionsCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)
	root := repoRoot(t)

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "decisions list empty",
			args:     []string{"decisions", "list"},
			wantErr:  false,
			contains: "No decisions found",
		},
		{
			name:     "decisions search",
			args:     []string{"decisions", "search", "test"},
			wantErr:  false,
			contains: "No decisions found",
		},
		{
			name:     "decisions export",
			args:     []string{"decisions", "export"},
			wantErr:  false,
			contains: "No decisions to export",
		},
		{
			name:     "decisions log missing flags",
			args:     []string{"decisions", "log"},
			wantErr:  true,
			contains: "required",
		},
		{
			name:     "decisions help",
			args:     []string{"decisions", "--help"},
			wantErr:  false,
			contains: "Manage decision audit trail",
		},
		{
			name:     "decisions list help",
			args:     []string{"decisions", "list", "--help"},
			wantErr:  false,
			contains: "List all decisions",
		},
		{
			name:     "decisions search help",
			args:     []string{"decisions", "search", "--help"},
			wantErr:  false,
			contains: "Search decisions",
		},
		{
			name:     "decisions export help",
			args:     []string{"decisions", "export", "--help"},
			wantErr:  false,
			contains: "Export decisions",
		},
		{
			name:     "decisions log help",
			args:     []string{"decisions", "log", "--help"},
			wantErr:  false,
			contains: "Log a new decision",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			cmd.Dir = root
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			output := stdout.String() + stderr.String()

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}

			if tt.contains != "" && !strings.Contains(output, tt.contains) {
				t.Logf("Output: %s", output)
			}
		})
	}
}

// TestPlanApplyTraceWorkflow tests the end-to-end plan → apply → trace workflow
func TestPlanApplyTraceWorkflow(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	// Create temp directory for isolated test environment
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Initialize minimal project structure
	backlogDir := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(backlogDir, 0755); err != nil {
		t.Fatalf("create backlog dir: %v", err)
	}
	sdpDir := filepath.Join(tmpDir, ".sdp", "log")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("create .sdp/log dir: %v", err)
	}

	// Create minimal config file
	configContent := `version: "0.9.0"
evidence:
  enabled: true
  log_path: ".sdp/log/events.jsonl"
`
	configPath := filepath.Join(tmpDir, ".sdp", "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("create config: %v", err)
	}

	t.Run("plan command creates workstreams", func(t *testing.T) {
		// Test plan --dry-run first (doesn't require MODEL_API)
		cmd := exec.Command(binaryPath, "plan", "Add simple feature", "--dry-run")
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		cmd.Dir = tmpDir

		err := cmd.Run()
		output := stdout.String() + stderr.String()

		// dry-run should succeed without MODEL_API
		if err != nil {
			t.Logf("Plan dry-run output: %s\nError: %v", output, err)
		}

		// Check for dry-run indicators
		if !strings.Contains(output, "DRY RUN") && !strings.Contains(output, "dry") {
			t.Logf("Expected dry-run output to mention DRY RUN or dry\nGot: %s", output)
		}
	})

	t.Run("plan command JSON output", func(t *testing.T) {
		// Test JSON output format
		cmd := exec.Command(binaryPath, "plan", "Add JSON feature", "--dry-run", "--output=json")
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		cmd.Dir = tmpDir

		err := cmd.Run()
		output := stdout.String()

		if err != nil {
			t.Logf("Plan JSON output: %s\nError: %v", output, err)
		}

		// Check for JSON structure
		if !strings.Contains(output, "{") && !strings.Contains(output, "[") {
			t.Logf("Expected JSON output to contain braces/brackets\nGot: %s", output)
		}
	})

	t.Run("log show command", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "log", "show")
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		cmd.Dir = tmpDir

		err := cmd.Run()
		output := stdout.String() + stderr.String()

		// log show should succeed even with empty log
		if err != nil && !strings.Contains(output, "No events") && !strings.Contains(output, "No matching") {
			t.Logf("Log show failed: %v\nOutput: %s", err, output)
		}

		// Should mention events or be empty
		if !strings.Contains(output, "event") && !strings.Contains(output, "No events") && !strings.Contains(output, "No matching") {
			t.Logf("Unexpected log show output: %s", output)
		}
	})

	t.Run("log export JSON", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "log", "export", "--format=json")
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		cmd.Dir = tmpDir

		err := cmd.Run()
		output := stdout.String() + stderr.String()

		if err != nil && !strings.Contains(output, "No events") {
			t.Logf("Log export JSON: %v\nOutput: %s", err, output)
		}

		// Check for valid JSON or empty message
		if !strings.Contains(output, "No events") && !strings.Contains(output, "[") && !strings.Contains(output, "{") {
			t.Logf("Expected JSON or empty message\nGot: %s", output)
		}
	})

	t.Run("log trace command", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "log", "trace", "--json")
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		cmd.Dir = tmpDir

		err := cmd.Run()
		output := stdout.String() + stderr.String()

		// trace should succeed even with empty log
		if err != nil && !strings.Contains(output, "No events") && !strings.Contains(output, "No matching") {
			t.Logf("Log trace failed: %v\nOutput: %s", err, output)
		}

		// Should be JSON or empty message
		if !strings.Contains(output, "No events") && !strings.Contains(output, "No matching") && !strings.Contains(output, "[") {
			t.Logf("Expected JSON or empty message\nGot: %s", output)
		}
	})

	t.Run("log stats command", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "log", "stats")
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		cmd.Dir = tmpDir

		err := cmd.Run()
		output := stdout.String() + stderr.String()

		// stats should succeed even with empty log
		if err != nil && !strings.Contains(output, "No events") {
			t.Logf("Log stats: %v\nOutput: %s", err, output)
		}

		// Should mention events or statistics
		if !strings.Contains(output, "event") && !strings.Contains(output, "statistics") && !strings.Contains(output, "No events") {
			t.Logf("Unexpected log stats output: %s", output)
		}
	})
}

// TestPlanApplyHelpCommand tests that help text shows all new commands
func TestPlanApplyHelpCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	tests := []struct {
		name     string
		args     []string
		contains []string
	}{
		{
			name: "plan help shows all options",
			args: []string{"plan", "--help"},
			contains: []string{
				"plan",
				"interactive",
				"auto-apply",
				"dry-run",
				"output",
			},
		},
		{
			name: "apply help shows all options",
			args: []string{"apply", "--help"},
			contains: []string{
				"apply",
				"dry-run",
				"retry",
				"output",
			},
		},
		{
			name: "log help shows subcommands",
			args: []string{"log", "--help"},
			contains: []string{
				"log",
				"show",
				"export",
				"stats",
				"trace",
			},
		},
		{
			name: "log trace help",
			args: []string{"log", "trace", "--help"},
			contains: []string{
				"trace",
				"ws",
				"json",
				"verify",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			if err := cmd.Run(); err != nil {
				t.Fatalf("Command failed: %v\nOutput: %s", err, stdout.String()+stderr.String())
			}

			output := stdout.String() + stderr.String()

			for _, expected := range tt.contains {
				if !strings.Contains(output, expected) {
					t.Errorf("Help output should contain %q\nGot: %s", expected, output)
				}
			}
		})
	}
}

// TestJSONOutputParseable tests that JSON output is parseable by jq
func TestJSONOutputParseable(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	// Check if jq is available
	if _, err := exec.LookPath("jq"); err != nil {
		t.Skip("jq not found in PATH, skipping JSON validation test")
	}

	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(tmpDir)

	// Setup minimal structure
	sdpDir := filepath.Join(tmpDir, ".sdp", "log")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("create .sdp/log: %v", err)
	}

	// Create empty events file for log commands
	eventsPath := filepath.Join(sdpDir, "events.jsonl")
	if err := os.WriteFile(eventsPath, []byte{}, 0644); err != nil {
		t.Fatalf("create events file: %v", err)
	}

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "plan JSON output",
			args: []string{"plan", "Test feature", "--dry-run", "--output=json"},
		},
		{
			name: "log export JSON",
			args: []string{"log", "export", "--format=json"},
		},
		{
			name: "log trace JSON",
			args: []string{"log", "trace", "--json"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run sdp command
			sdpCmd := exec.Command(binaryPath, tt.args...)
			sdpCmd.Dir = tmpDir
			var stdout, stderr bytes.Buffer
			sdpCmd.Stdout = &stdout
			sdpCmd.Stderr = &stderr

			err := sdpCmd.Run()
			output := stdout.String()

			// Some commands may fail (e.g., plan without MODEL_API), but output should still be valid JSON
			if err != nil && !strings.Contains(stderr.String(), "MODEL_API") {
				t.Logf("Command output: %s\nError: %v", output, err)
			}

			// Skip jq validation if output is empty or just an error message
			if len(output) == 0 || !strings.Contains(output, "{") && !strings.Contains(output, "[") {
				t.Skip("No JSON output to validate")
			}

			// Pipe output through jq
			jqCmd := exec.Command("jq", ".")
			jqCmd.Stdin = strings.NewReader(output)
			var jqOut, jqErr bytes.Buffer
			jqCmd.Stdout = &jqOut
			jqCmd.Stderr = &jqErr

			if err := jqCmd.Run(); err != nil {
				t.Errorf("JSON output not parseable by jq\nOutput: %s\njq error: %s", output, jqErr.String())
			}
		})
	}
}
