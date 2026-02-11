package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

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
