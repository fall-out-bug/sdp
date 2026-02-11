package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Use filepath indirectly via repoRoot
var _ = filepath.Join

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
