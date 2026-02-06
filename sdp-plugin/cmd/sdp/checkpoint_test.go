package main

import (
	"testing"
)

// TestCheckpointCreateCmd tests the checkpoint create command structure
func TestCheckpointCreateCmd(t *testing.T) {
	cmd := checkpointCreateCmd

	if cmd == nil {
		t.Fatal("checkpointCreateCmd is nil")
	}

	// Test command structure
	if cmd.Use != "create <id> <feature-id>" {
		t.Errorf("checkpointCreateCmd has wrong use: %s", cmd.Use)
	}

	// Test args validation
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
			name:        "valid checkpoint",
			args:        []string{"test-checkpoint", "F042"},
			expectError: false,
		},
		{
			name:        "only id",
			args:        []string{"test-checkpoint"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cmd.Args(cmd, tt.args)

			if tt.expectError && err == nil {
				t.Errorf("checkpointCreateCmd() args validation expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("checkpointCreateCmd() args validation unexpected error: %v", err)
			}
		})
	}
}

// TestCheckpointResumeCmd tests the checkpoint resume command structure
func TestCheckpointResumeCmd(t *testing.T) {
	cmd := checkpointResumeCmd

	if cmd == nil {
		t.Fatal("checkpointResumeCmd is nil")
	}

	if cmd.Use != "resume <id>" {
		t.Errorf("checkpointResumeCmd has wrong use: %s", cmd.Use)
	}
}

// TestCheckpointListCmd tests the checkpoint list command structure
func TestCheckpointListCmd(t *testing.T) {
	cmd := checkpointListCmd

	if cmd == nil {
		t.Fatal("checkpointListCmd is nil")
	}

	if cmd.Use != "list" {
		t.Errorf("checkpointListCmd has wrong use: %s", cmd.Use)
	}
}

// TestCheckpointCleanCmd tests the checkpoint clean command structure
func TestCheckpointCleanCmd(t *testing.T) {
	cmd := checkpointCleanCmd

	if cmd == nil {
		t.Fatal("checkpointCleanCmd is nil")
	}

	if cmd.Use != "clean" {
		t.Errorf("checkpointCleanCmd has wrong use: %s", cmd.Use)
	}

	// Check that age flag exists
	if cmd.Flags().Lookup("age") == nil {
		t.Error("checkpointCleanCmd missing --age flag")
	}
}

// TestCheckpointCmdStructure tests the checkpoint command structure
func TestCheckpointCmdStructure(t *testing.T) {
	cmd := checkpointCmd

	// Test command structure
	if cmd.Use != "checkpoint" {
		t.Errorf("checkpointCmd has wrong use: %s", cmd.Use)
	}

	// Test subcommands
	expectedSubcommands := []string{"create", "resume", "list", "clean"}
	for _, expected := range expectedSubcommands {
		found := false
		for _, c := range cmd.Commands() {
			if c.Name() == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("checkpointCmd missing subcommand: %s", expected)
		}
	}

	// Test persistent flags
	if cmd.PersistentFlags().Lookup("dir") == nil {
		t.Error("checkpointCmd missing --dir persistent flag")
	}
}
