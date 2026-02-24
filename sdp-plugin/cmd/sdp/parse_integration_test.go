package main

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Use filepath indirectly via repoRoot
var _ = filepath.Join

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
			args:     []string{"parse", "00-016-01"},
			wantErr:  false,
			contains: "00-016-01",
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
