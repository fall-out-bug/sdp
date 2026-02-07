package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall-out-bug/sdp/internal/drift"
)

// TestGuardActivateCmd tests the guard activate command
func TestGuardActivateCmd(t *testing.T) {
	// Create temp config directory
	tmpDir := t.TempDir()

	// Create fake command
	cmd := guardActivate()

	tests := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "no args",
			args:        []string{},
			expectError: true,
			errorMsg:    "requires exactly 1 arg",
		},
		{
			name:        "valid ws id",
			args:        []string{"00-001-01"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Recover from panic caused by cobra.ExactArgs validation
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectError {
						t.Errorf("guardActivate() panicked with: %v", r)
					}
				}
			}()

			// Set config dir env
			oldConfigDir := os.Getenv("XDG_CONFIG_HOME")
			os.Setenv("XDG_CONFIG_HOME", tmpDir)
			defer os.Setenv("XDG_CONFIG_HOME", oldConfigDir)

			err := cmd.RunE(cmd, tt.args)

			if tt.expectError && err == nil {
				t.Errorf("guardActivate() expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("guardActivate() unexpected error: %v", err)
			}
			if tt.expectError && tt.errorMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("guardActivate() error = %q, should contain %q", err.Error(), tt.errorMsg)
				}
			}
		})
	}
}

// TestGuardCheckCmd tests the guard check command
func TestGuardCheckCmd(t *testing.T) {
	// Create temp config directory
	tmpDir := t.TempDir()

	// Create fake command
	cmd := guardCheck()

	tests := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "no args",
			args:        []string{},
			expectError: true,
			errorMsg:    "requires exactly 1 arg",
		},
		{
			name:        "valid file path - blocked when no active WS",
			args:        []string{"internal/file.go"},
			expectError: true, // No active WS means editing is blocked
			errorMsg:    "No active WS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Recover from panic caused by cobra.ExactArgs validation
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectError {
						t.Errorf("guardCheck() panicked with: %v", r)
					}
				}
			}()

			// Set config dir env
			oldConfigDir := os.Getenv("XDG_CONFIG_HOME")
			os.Setenv("XDG_CONFIG_HOME", tmpDir)
			defer os.Setenv("XDG_CONFIG_HOME", oldConfigDir)

			err := cmd.RunE(cmd, tt.args)

			if tt.expectError && err == nil {
				t.Errorf("guardCheck() expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("guardCheck() unexpected error: %v", err)
			}
		})
	}
}

// TestGuardStatusCmd tests the guard status command
func TestGuardStatusCmd(t *testing.T) {
	// Create temp config directory
	tmpDir := t.TempDir()

	// Create fake command
	cmd := guardStatus()

	// Set config dir env
	oldConfigDir := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", oldConfigDir)

	// Run command
	err := cmd.RunE(cmd, []string{})

	// Should succeed (no active WS is valid state)
	if err != nil {
		t.Errorf("guardStatus() failed: %v", err)
	}
}

// TestGuardDeactivateCmd tests the guard deactivate command
func TestGuardDeactivateCmd(t *testing.T) {
	// Create temp config directory
	tmpDir := t.TempDir()

	// Create fake command
	cmd := guardDeactivate()

	// Set config dir env
	oldConfigDir := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", oldConfigDir)

	// Run command
	err := cmd.RunE(cmd, []string{})

	// Should succeed (deactivating when inactive is ok)
	if err != nil {
		t.Errorf("guardDeactivate() failed: %v", err)
	}
}

// TestFindDriftProjectRoot tests project root detection
func TestFindDriftProjectRoot(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() string
		cleanup     func(string)
		expectError bool
	}{
		{
			name: "project root with docs",
			setup: func() string {
				tmpDir := t.TempDir()
				os.MkdirAll(filepath.Join(tmpDir, "docs"), 0755)
				originalWd, _ := os.Getwd()
				os.Chdir(tmpDir)
				return originalWd
			},
			cleanup: func(originalWd string) {
				os.Chdir(originalWd)
			},
			expectError: false,
		},
		{
			name: "project root with .beads",
			setup: func() string {
				tmpDir := t.TempDir()
				os.MkdirAll(filepath.Join(tmpDir, ".beads"), 0755)
				originalWd, _ := os.Getwd()
				os.Chdir(tmpDir)
				return originalWd
			},
			cleanup: func(originalWd string) {
				os.Chdir(originalWd)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalWd := tt.setup()
			defer tt.cleanup(originalWd)

			root, err := findDriftProjectRoot()

			if tt.expectError && err == nil {
				t.Errorf("findDriftProjectRoot() expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("findDriftProjectRoot() unexpected error: %v", err)
			}
			if !tt.expectError && root == "" {
				t.Errorf("findDriftProjectRoot() returned empty root")
			}
		})
	}
}

// TestFindDriftWorkstreamFile tests workstream file finding
func TestFindDriftWorkstreamFile(t *testing.T) {
	// Create temp directory with workstream file
	tmpDir := t.TempDir()
	backlogDir := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(backlogDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	wsID := "00-050-01"
	wsPath := filepath.Join(backlogDir, wsID+".md")
	wsContent := `---
ws_id: 00-050-01
feature: F050
status: completed
---
# Test Workstream
`
	if err := os.WriteFile(wsPath, []byte(wsContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		wsID        string
		expectError bool
	}{
		{
			name:        "existing workstream in backlog",
			wsID:        wsID,
			expectError: false,
		},
		{
			name:        "non-existent workstream",
			wsID:        "99-999-99",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foundPath, err := findDriftWorkstreamFile(tmpDir, tt.wsID)

			if tt.expectError && err == nil {
				t.Errorf("findDriftWorkstreamFile(%q) expected error but got none", tt.wsID)
			}
			if !tt.expectError && err != nil {
				t.Errorf("findDriftWorkstreamFile(%q) unexpected error: %v", tt.wsID, err)
			}
			if !tt.expectError && !strings.Contains(foundPath, tt.wsID+".md") {
				t.Errorf("findDriftWorkstreamFile(%q) = %q, should contain %q", tt.wsID, foundPath, tt.wsID+".md")
			}
		})
	}
}

// TestCountDriftErrors tests error counting in drift reports
func TestCountDriftErrors(t *testing.T) {
	report := &drift.DriftReport{
		Issues: []drift.DriftIssue{
			{Status: drift.StatusError, File: "file1.go"},
			{Status: drift.StatusWarning, File: "file2.go"},
			{Status: drift.StatusError, File: "file3.go"},
		},
	}

	count := countDriftErrors(report)
	if count != 2 {
		t.Errorf("countDriftErrors() = %d, want %d", count, 2)
	}
}

// TestCountDriftWarnings tests warning counting in drift reports
func TestCountDriftWarnings(t *testing.T) {
	report := &drift.DriftReport{
		Issues: []drift.DriftIssue{
			{Status: drift.StatusError, File: "file1.go"},
			{Status: drift.StatusWarning, File: "file2.go"},
			{Status: drift.StatusWarning, File: "file3.go"},
		},
	}

	count := countDriftWarnings(report)
	if count != 2 {
		t.Errorf("countDriftWarnings() = %d, want %d", count, 2)
	}
}
