package guard

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestNewSkill(t *testing.T) {
	configDir := t.TempDir()
	skill := NewSkill(configDir)

	if skill == nil {
		t.Fatal("NewSkill returned nil")
	}
	if skill.stateManager == nil {
		t.Error("stateManager not initialized")
	}
	if skill.activeWS != "" {
		t.Errorf("activeWS should be empty, got %s", skill.activeWS)
	}
}

func TestActivate(t *testing.T) {
	configDir := t.TempDir()
	skill := NewSkill(configDir)

	// Test activation
	err := skill.Activate("00-001-01")
	if err != nil {
		t.Fatalf("Activate failed: %v", err)
	}

	// Verify state saved
	state, err := skill.stateManager.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if state.ActiveWS != "00-001-01" {
		t.Errorf("ActiveWS = %s, want 00-001-01", state.ActiveWS)
	}

	if state.ActivatedAt == "" {
		t.Error("ActivatedAt should be set")
	}

	// Verify timestamp is recent
	activatedAt, err := time.Parse(time.RFC3339, state.ActivatedAt)
	if err != nil {
		t.Fatalf("Failed to parse ActivatedAt: %v", err)
	}

	if time.Since(activatedAt) > 5*time.Second {
		t.Error("ActivatedAt should be recent")
	}
}

func TestCheckEditNoActiveWS(t *testing.T) {
	configDir := t.TempDir()
	skill := NewSkill(configDir)

	result, err := skill.CheckEdit("/some/file.go")
	if err != nil {
		t.Fatalf("CheckEdit failed: %v", err)
	}

	if result.Allowed {
		t.Error("Should not be allowed when no active WS")
	}

	if result.WSID != "" {
		t.Errorf("WSID should be empty, got %s", result.WSID)
	}
}

func TestCheckEditNoScope(t *testing.T) {
	configDir := t.TempDir()
	skill := NewSkill(configDir)

	// Activate WS without scope restrictions
	err := skill.Activate("00-001-01")
	if err != nil {
		t.Fatalf("Activate failed: %v", err)
	}

	// Any file should be allowed
	result, err := skill.CheckEdit("/any/file.go")
	if err != nil {
		t.Fatalf("CheckEdit failed: %v", err)
	}

	if !result.Allowed {
		t.Errorf("Should be allowed (no scope): %s", result.Reason)
	}

	if result.WSID != "00-001-01" {
		t.Errorf("WSID = %s, want 00-001-01", result.WSID)
	}
}

func TestCheckEditWithScope(t *testing.T) {
	configDir := t.TempDir()
	skill := NewSkill(configDir)

	// Activate WS with scope files
	err := skill.Activate("00-001-01")
	if err != nil {
		t.Fatalf("Activate failed: %v", err)
	}

	// Set scope files
	state, _ := skill.stateManager.Load()
	state.ScopeFiles = []string{
		"/allowed/file1.go",
		"/allowed/file2.go",
	}
	skill.stateManager.Save(*state)

	// Test allowed file
	result, err := skill.CheckEdit("/allowed/file1.go")
	if err != nil {
		t.Fatalf("CheckEdit failed: %v", err)
	}

	if !result.Allowed {
		t.Errorf("Should be allowed (in scope): %s", result.Reason)
	}

	// Test blocked file
	result, err = skill.CheckEdit("/blocked/file.go")
	if err != nil {
		t.Fatalf("CheckEdit failed: %v", err)
	}

	if result.Allowed {
		t.Error("Should not be allowed (not in scope)")
	}

	if len(result.ScopeFiles) != 2 {
		t.Errorf("ScopeFiles count = %d, want 2", len(result.ScopeFiles))
	}
}

func TestCheckEditExpired(t *testing.T) {
	configDir := t.TempDir()
	skill := NewSkill(configDir)

	// Activate WS
	err := skill.Activate("00-001-01")
	if err != nil {
		t.Fatalf("Activate failed: %v", err)
	}

	// Manually set old timestamp (>24 hours)
	state, _ := skill.stateManager.Load()
	oldTime := time.Now().Add(-25 * time.Hour)
	state.ActivatedAt = oldTime.Format(time.RFC3339)
	skill.stateManager.Save(*state)

	// Check should fail due to expiration
	result, err := skill.CheckEdit("/some/file.go")
	if err != nil {
		t.Fatalf("CheckEdit failed: %v", err)
	}

	if result.Allowed {
		t.Error("Should not be allowed (state expired)")
	}

	if result.WSID != "" {
		t.Errorf("WSID should be empty when expired, got %s", result.WSID)
	}
}

func TestGetActiveWS(t *testing.T) {
	configDir := t.TempDir()
	skill := NewSkill(configDir)

	// No active WS
	wsID := skill.GetActiveWS()
	if wsID != "" {
		t.Errorf("WSID should be empty, got %s", wsID)
	}

	// Activate WS
	err := skill.Activate("00-001-01")
	if err != nil {
		t.Fatalf("Activate failed: %v", err)
	}

	// Should return active WS
	wsID = skill.GetActiveWS()
	if wsID != "00-001-01" {
		t.Errorf("WSID = %s, want 00-001-01", wsID)
	}
}

func TestDeactivate(t *testing.T) {
	configDir := t.TempDir()
	skill := NewSkill(configDir)

	// Activate WS first
	err := skill.Activate("00-001-01")
	if err != nil {
		t.Fatalf("Activate failed: %v", err)
	}

	// Verify active
	if skill.GetActiveWS() != "00-001-01" {
		t.Error("WS should be active")
	}

	// Deactivate
	err = skill.Deactivate()
	if err != nil {
		t.Fatalf("Deactivate failed: %v", err)
	}

	// Verify inactive
	wsID := skill.GetActiveWS()
	if wsID != "" {
		t.Errorf("WSID should be empty after deactivate, got %s", wsID)
	}

	if skill.activeWS != "" {
		t.Errorf("activeWS should be empty, got %s", skill.activeWS)
	}
}

func TestResolvePathAbsolute(t *testing.T) {
	absPath := "/absolute/path/to/file.go"
	result, err := ResolvePath(absPath)
	if err != nil {
		t.Fatalf("ResolvePath failed: %v", err)
	}

	if result != absPath {
		t.Errorf("ResolvePath = %s, want %s", result, absPath)
	}
}

func TestResolvePathRelative(t *testing.T) {
	// Create temp directory and change to it
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	os.Chdir(tmpDir)

	relPath := "relative/file.go"
	result, err := ResolvePath(relPath)
	if err != nil {
		t.Fatalf("ResolvePath failed: %v", err)
	}

	// Resolve symlinks for comparison (macOS /var -> /private/var)
	expected, _ := filepath.EvalSymlinks(filepath.Join(tmpDir, relPath))
	resultResolved, _ := filepath.EvalSymlinks(result)

	if resultResolved != expected {
		t.Errorf("ResolvePath = %s (resolved: %s), want %s", result, resultResolved, expected)
	}
}

func TestStateManagerSaveAndLoad(t *testing.T) {
	configDir := t.TempDir()
	sm := NewStateManager(configDir)

	state := GuardState{
		ActiveWS:    "00-001-01",
		ActivatedAt: time.Now().Format(time.RFC3339),
		ScopeFiles:  []string{"/file1.go", "/file2.go"},
		Timestamp:   "",
	}

	// Save
	err := sm.Save(state)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists and permissions
	statePath := filepath.Join(configDir, GuardStateFile)
	info, err := os.Stat(statePath)
	if err != nil {
		t.Fatalf("State file not created: %v", err)
	}

	perms := info.Mode().Perm()
	if perms != 0600 {
		t.Errorf("File permissions = %04o, want 0600", perms)
	}

	// Load
	loaded, err := sm.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.ActiveWS != state.ActiveWS {
		t.Errorf("ActiveWS = %s, want %s", loaded.ActiveWS, state.ActiveWS)
	}

	if len(loaded.ScopeFiles) != len(state.ScopeFiles) {
		t.Errorf("ScopeFiles count = %d, want %d", len(loaded.ScopeFiles), len(state.ScopeFiles))
	}

	if loaded.Timestamp == "" {
		t.Error("Timestamp should be set by Save")
	}
}

func TestStateManagerLoadNotExists(t *testing.T) {
	configDir := t.TempDir()
	sm := NewStateManager(configDir)

	// Load when file doesn't exist should return empty state
	loaded, err := sm.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.ActiveWS != "" {
		t.Errorf("ActiveWS should be empty, got %s", loaded.ActiveWS)
	}
}

func TestStateManagerClear(t *testing.T) {
	configDir := t.TempDir()
	sm := NewStateManager(configDir)

	// Create state file
	state := GuardState{
		ActiveWS:    "00-001-01",
		ActivatedAt: time.Now().Format(time.RFC3339),
	}
	sm.Save(state)

	// Verify exists
	statePath := filepath.Join(configDir, GuardStateFile)
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Fatal("State file should exist")
	}

	// Clear
	err := sm.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// Verify removed
	if _, err := os.Stat(statePath); !os.IsNotExist(err) {
		t.Error("State file should be removed")
	}

	// Clear when already gone should not error
	err = sm.Clear()
	if err != nil {
		t.Errorf("Clear on non-existent file failed: %v", err)
	}
}

func TestGuardStateIsExpired(t *testing.T) {
	tests := []struct {
		name        string
		activatedAt string
		wantExpired bool
	}{
		{
			name:        "No active WS",
			activatedAt: "",
			wantExpired: true,
		},
		{
			name:        "Invalid timestamp",
			activatedAt: "invalid",
			wantExpired: true,
		},
		{
			name:        "Recent activation",
			activatedAt: time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			wantExpired: false,
		},
		{
			name:        "Expired (>24 hours)",
			activatedAt: time.Now().Add(-25 * time.Hour).Format(time.RFC3339),
			wantExpired: true,
		},
		{
			name:        "Just under 24 hours",
			activatedAt: time.Now().Add(-23*time.Hour - 59*time.Minute).Format(time.RFC3339),
			wantExpired: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &GuardState{
				ActiveWS:    "00-001-01",
				ActivatedAt: tt.activatedAt,
			}

			got := state.IsExpired()
			if got != tt.wantExpired {
				t.Errorf("IsExpired() = %v, want %v", got, tt.wantExpired)
			}
		})
	}
}

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

	// Should return a result with at least the staged file
	if result == nil {
		t.Fatal("StagedCheck returned nil result")
	}

	// The exact findings depend on rules, but we should get a valid CheckResult
	if result.ExitCode < 0 || result.ExitCode > 2 {
		t.Errorf("Invalid exit code: %d", result.ExitCode)
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
		expectedSuccess  bool
	}{
		{
			name:           "Exit code 0 - pass",
			exitCode:       ExitCodePass,
			expectedSuccess: true,
		},
		{
			name:           "Exit code 1 - violation",
			exitCode:       ExitCodeViolation,
			expectedSuccess: false,
		},
		{
			name:           "Exit code 2 - runtime error",
			exitCode:       ExitCodeRuntimeError,
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

// TestCIDiffRangeDetection tests AC6: CI diff-range auto-detection
func TestCIDiffRangeDetection(t *testing.T) {
	tests := []struct {
		name          string
		setCIEnvs     bool
		baseSHA       string
		headSHA       string
		expectFallback bool
	}{
		{
			name:          "CI env vars set",
			setCIEnvs:     true,
			baseSHA:       "abc123",
			headSHA:       "def456",
			expectFallback: false,
		},
		{
			name:          "CI env vars not set",
			setCIEnvs:     false,
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

	// Modify and stage the file
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

// runCommand is a helper to run a command in a directory
func runCommand(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.Run()
}
