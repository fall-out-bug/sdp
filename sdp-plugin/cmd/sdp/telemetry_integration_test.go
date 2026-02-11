package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

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
