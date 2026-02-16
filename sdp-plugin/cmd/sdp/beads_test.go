package main

import (
	"testing"
)

// TestBeadsCmd tests the beads command structure
func TestBeadsCmd(t *testing.T) {
	cmd := beadsCmd()

	// Test command structure
	if cmd.Use != "beads" {
		t.Errorf("beadsCmd() has wrong use: %s", cmd.Use)
	}

	// Test subcommands
	subcommands := []string{"ready", "show", "update", "sync"}
	for _, sub := range subcommands {
		found := false
		for _, c := range cmd.Commands() {
			if c.Name() == sub {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("beadsCmd() missing subcommand: %s", sub)
		}
	}
}

// TestBeadsShowCmdArgs tests the beads show command argument validation
func TestBeadsShowCmdArgs(t *testing.T) {
	cmd := beadsShowCmd()

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "no args",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "valid task id",
			args:        []string{"sdp-abc"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cmd.Args(cmd, tt.args)

			if tt.expectError && err == nil {
				t.Errorf("beadsShowCmd() args validation expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("beadsShowCmd() args validation unexpected error: %v", err)
			}
		})
	}
}
