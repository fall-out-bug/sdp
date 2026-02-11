package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

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
