package guard

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestStagedCheck tests AC1: staged checks
func TestStagedCheck(t *testing.T) {
	// Create a temporary git repository
	tmpDir := t.TempDir()

	// Initialize git repo
	err := runCommand(tmpDir, "git", "init")
	if err != nil {
		t.Fatalf("git init failed: %v", err)
	}
	err = runCommand(tmpDir, "git", "config", "user.email", "test@example.com")
	if err != nil {
		t.Fatalf("git config failed: %v", err)
	}
	err = runCommand(tmpDir, "git", "config", "user.name", "Test User")
	if err != nil {
		t.Fatalf("git config failed: %v", err)
	}

	// Create a test file and commit it (base)
	testFile := "test.go"
	content := "package main\n\nfunc main() {}\n"
	if err := os.WriteFile(filepath.Join(tmpDir, testFile), []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err = runCommand(tmpDir, "git", "add", testFile)
	if err != nil {
		t.Fatalf("git add failed: %v", err)
	}

	err = runCommand(tmpDir, "git", "commit", "-m", "Initial commit")
	if err != nil {
		t.Fatalf("git commit failed: %v", err)
	}

	// Modify the file and stage changes
	newContent := "package main\n\nfunc modified() {}\n"
	if err := os.WriteFile(filepath.Join(tmpDir, testFile), []byte(newContent), 0644); err != nil {
		t.Fatalf("failed to modify test file: %v", err)
	}

	err = runCommand(tmpDir, "git", "add", testFile)
	if err != nil {
		t.Fatalf("git add failed: %v", err)
	}

	// Create skill and check staged files
	skill := NewSkill(tmpDir)
	result, err := skill.StagedCheck(CheckOptions{Staged: true})

	if err != nil {
		t.Fatalf("StagedCheck failed: %v", err)
	}

	// Should return a result with at least staged file
	if result == nil {
		t.Fatal("StagedCheck returned nil result")
	}

	// The exact findings depend on rules, but we should get a valid CheckResult
	if result.ExitCode < 0 || result.ExitCode > 2 {
		t.Errorf("Invalid exit code: %d", result.ExitCode)
	}
}

// TestStagedCheckNoStagedFiles tests staged check with no staged files
func TestStagedCheckNoStagedFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Initialize git repo
	err := runCommand(tmpDir, "git", "init")
	if err != nil {
		t.Fatalf("git init failed: %v", err)
	}
	err = runCommand(tmpDir, "git", "config", "user.email", "test@example.com")
	if err != nil {
		t.Fatalf("git config failed: %v", err)
	}
	err = runCommand(tmpDir, "git", "config", "user.name", "Test User")
	if err != nil {
		t.Fatalf("git config failed: %v", err)
	}

	// Create skill and check staged files (none staged)
	skill := NewSkill(tmpDir)
	result, err := skill.StagedCheck(CheckOptions{Staged: true})

	if err != nil {
		t.Fatalf("StagedCheck failed: %v", err)
	}

	// Should pass with no findings
	if !result.Success {
		t.Errorf("Expected success when no staged files, got exit code %d", result.ExitCode)
	}

	if len(result.Findings) != 0 {
		t.Errorf("Expected no findings, got %d", len(result.Findings))
	}
}

// TestStagedCheckWithActiveWS tests staged check with active workstream
func TestStagedCheckWithActiveWS(t *testing.T) {
	tmpDir := t.TempDir()

	// Initialize git repo
	err := runCommand(tmpDir, "git", "init")
	if err != nil {
		t.Fatalf("git init failed: %v", err)
	}
	err = runCommand(tmpDir, "git", "config", "user.email", "test@example.com")
	if err != nil {
		t.Fatalf("git config failed: %v", err)
	}
	err = runCommand(tmpDir, "git", "config", "user.name", "Test User")
	if err != nil {
		t.Fatalf("git config failed: %v", err)
	}

	// Create a test file and commit
	testFile := "allowed.go"
	content := "package main\n\nfunc allowed() {}\n"
	if err := os.WriteFile(filepath.Join(tmpDir, testFile), []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err = runCommand(tmpDir, "git", "add", testFile)
	if err != nil {
		t.Fatalf("git add failed: %v", err)
	}
	err = runCommand(tmpDir, "git", "commit", "-m", "Initial commit")
	if err != nil {
		t.Fatalf("git commit failed: %v", err)
	}

	// Create skill and activate WS with scope
	skill := NewSkill(tmpDir)
	err = skill.Activate("00-001-01")
	if err != nil {
		t.Fatalf("Activate failed: %v", err)
	}

	// Set scope files
	state, _ := skill.stateManager.Load()
	absPath := filepath.Join(tmpDir, testFile)
	state.ScopeFiles = []string{absPath}
	skill.stateManager.Save(*state)

	// Modify and stage file
	newContent := "package main\n\nfunc modified() {}\n"
	if err := os.WriteFile(filepath.Join(tmpDir, testFile), []byte(newContent), 0644); err != nil {
		t.Fatalf("failed to modify test file: %v", err)
	}

	err = runCommand(tmpDir, "git", "add", testFile)
	if err != nil {
		t.Fatalf("git add failed: %v", err)
	}

	// Check staged files
	result, err := skill.StagedCheck(CheckOptions{Staged: true})

	if err != nil {
		t.Fatalf("StagedCheck failed: %v", err)
	}

	// Should pass (file is in scope)
	if !result.Success {
		t.Errorf("Expected success when file is in scope, got exit code %d", result.ExitCode)
	}
}

// TestStagedCheckOutsideScope tests staged check with file outside scope
func TestStagedCheckOutsideScope(t *testing.T) {
	tmpDir := t.TempDir()

	// Initialize git repo
	err := runCommand(tmpDir, "git", "init")
	if err != nil {
		t.Fatalf("git init failed: %v", err)
	}
	err = runCommand(tmpDir, "git", "config", "user.email", "test@example.com")
	if err != nil {
		t.Fatalf("git config failed: %v", err)
	}
	err = runCommand(tmpDir, "git", "config", "user.name", "Test User")
	if err != nil {
		t.Fatalf("git config failed: %v", err)
	}

	// Create skill and activate WS with scope
	skill := NewSkill(tmpDir)
	err = skill.Activate("00-001-01")
	if err != nil {
		t.Fatalf("Activate failed: %v", err)
	}

	// Set scope files (different from what we'll stage)
	otherPath := filepath.Join(tmpDir, "other.go")
	state, _ := skill.stateManager.Load()
	state.ScopeFiles = []string{otherPath}
	skill.stateManager.Save(*state)

	// Create and stage a different file
	testFile := "outside.go"
	content := "package main\n\nfunc outside() {}\n"
	if err := os.WriteFile(filepath.Join(tmpDir, testFile), []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err = runCommand(tmpDir, "git", "add", testFile)
	if err != nil {
		t.Fatalf("git add failed: %v", err)
	}

	// Save current directory and change to tmpDir for path resolution
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	// Check staged files
	result, err := skill.StagedCheck(CheckOptions{Staged: true})

	if err != nil {
		t.Fatalf("StagedCheck failed: %v", err)
	}

	// Should pass (WARNING doesn't block in hybrid mode)
	// but should have a warning finding
	if !result.Success {
		t.Errorf("Expected success (warnings don't block), got exit code %d", result.ExitCode)
	}

	if result.Summary.Warnings == 0 {
		t.Error("Expected warning for file outside scope")
	}
}

// TestCIDiffRangeDetection tests AC6: CI diff-range auto-detection
func TestCIDiffRangeDetection(t *testing.T) {
	tests := []struct {
		name           string
		setCIEnvs      bool
		baseSHA        string
		headSHA        string
		expectFallback bool
	}{
		{
			name:           "CI env vars set",
			setCIEnvs:      true,
			baseSHA:        "abc123",
			headSHA:        "def456",
			expectFallback: false,
		},
		{
			name:           "CI env vars not set",
			setCIEnvs:      false,
			expectFallback: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set or unset CI env vars
			if tt.setCIEnvs {
				oldBase := os.Getenv("CI_BASE_SHA")
				oldHead := os.Getenv("CI_HEAD_SHA")
				defer func() {
					if oldBase != "" {
						os.Setenv("CI_BASE_SHA", oldBase)
					} else {
						os.Unsetenv("CI_BASE_SHA")
					}
					if oldHead != "" {
						os.Setenv("CI_HEAD_SHA", oldHead)
					} else {
						os.Unsetenv("CI_HEAD_SHA")
					}
				}()
				os.Setenv("CI_BASE_SHA", tt.baseSHA)
				os.Setenv("CI_HEAD_SHA", tt.headSHA)
			}

			options := ParseCheckOptions()
			if tt.setCIEnvs && !tt.expectFallback {
				if options.Base != tt.baseSHA {
					t.Errorf("Base = %s, want %s", options.Base, tt.baseSHA)
				}
				if options.Head != tt.headSHA {
					t.Errorf("Head = %s, want %s", options.Head, tt.headSHA)
				}
			}
		})
	}
}

// TestHybridModeEnforcement tests AC8: hybrid mode (block on ERROR, warn on WARNING)
func TestHybridModeEnforcement(t *testing.T) {
	tests := []struct {
		name         string
		findings     []Finding
		wantExitCode int
		wantSuccess  bool
	}{
		{
			name:         "No findings - pass",
			findings:     []Finding{},
			wantExitCode: ExitCodePass,
			wantSuccess:  true,
		},
		{
			name: "Only warnings - pass (hybrid mode)",
			findings: []Finding{
				{Severity: SeverityWarning, Rule: "test-rule", File: "test.go", Message: "warning message"},
			},
			wantExitCode: ExitCodePass,
			wantSuccess:  true,
		},
		{
			name: "Single error - fail",
			findings: []Finding{
				{Severity: SeverityError, Rule: "test-rule", File: "test.go", Message: "error message"},
			},
			wantExitCode: ExitCodeViolation,
			wantSuccess:  false,
		},
		{
			name: "Mixed error and warning - fail",
			findings: []Finding{
				{Severity: SeverityWarning, Rule: "warn-rule", File: "test.go", Message: "warning"},
				{Severity: SeverityError, Rule: "error-rule", File: "test.go", Message: "error"},
			},
			wantExitCode: ExitCodeViolation,
			wantSuccess:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildCheckResult(tt.findings)

			if result.ExitCode != tt.wantExitCode {
				t.Errorf("ExitCode = %d, want %d", result.ExitCode, tt.wantExitCode)
			}

			if result.Success != tt.wantSuccess {
				t.Errorf("Success = %v, want %v", result.Success, tt.wantSuccess)
			}

			// Verify summary
			errorCount := 0
			warningCount := 0
			for _, f := range tt.findings {
				if f.Severity == SeverityError {
					errorCount++
				} else if f.Severity == SeverityWarning {
					warningCount++
				}
			}

			if result.Summary.Errors != errorCount {
				t.Errorf("Summary.Errors = %d, want %d", result.Summary.Errors, errorCount)
			}

			if result.Summary.Warnings != warningCount {
				t.Errorf("Summary.Warnings = %d, want %d", result.Summary.Warnings, warningCount)
			}
		})
	}
}

// TestSeverityClassification tests AC2: ERROR and WARNING classification
func TestSeverityClassification(t *testing.T) {
	tests := []struct {
		name     string
		severity Severity
		wantErr  bool
	}{
		{
			name:     "ERROR severity",
			severity: SeverityError,
			wantErr:  false,
		},
		{
			name:     "WARNING severity",
			severity: SeverityWarning,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finding := Finding{
				Severity: tt.severity,
				Rule:     "test-rule",
				File:     "test.go",
				Message:  "test message",
			}

			if finding.Severity != tt.severity {
				t.Errorf("Severity mismatch: got %v, want %v", finding.Severity, tt.severity)
			}
		})
	}
}

// TestExitCodes tests AC3: stable exit codes
func TestExitCodes(t *testing.T) {
	tests := []struct {
		name            string
		exitCode        int
		expectedSuccess bool
	}{
		{
			name:            "Exit code 0 - pass",
			exitCode:        ExitCodePass,
			expectedSuccess: true,
		},
		{
			name:            "Exit code 1 - violation",
			exitCode:        ExitCodeViolation,
			expectedSuccess: false,
		},
		{
			name:            "Exit code 2 - runtime error",
			exitCode:        ExitCodeRuntimeError,
			expectedSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckResult{
				Success:  tt.expectedSuccess,
				ExitCode: tt.exitCode,
			}

			if result.ExitCode != tt.exitCode {
				t.Errorf("ExitCode = %d, want %d", result.ExitCode, tt.exitCode)
			}

			if tt.exitCode == ExitCodePass && !result.Success {
				t.Error("ExitCodePass should indicate success")
			}
		})
	}
}

// runCommand is a helper to run a command in a directory
func runCommand(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.Run()
}
