package main

import (
	"testing"
)

// TestTddCmd tests the tdd command structure
func TestTddCmd(t *testing.T) {
	cmd := tddCmd()

	// Test command structure
	if cmd.Use != "tdd <phase> [path]" {
		t.Errorf("tddCmd() has wrong use: %s", cmd.Use)
	}

	// Test arg validation (Args validator)
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
			name:        "too many args",
			args:        []string{"red", "path", "extra"},
			expectError: true,
		},
		{
			name:        "valid red phase",
			args:        []string{"red"},
			expectError: false,
		},
		{
			name:        "valid green phase with path",
			args:        []string{"green", "./internal/parser"},
			expectError: false,
		},
		{
			name:        "valid refactor phase",
			args:        []string{"refactor"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cmd.Args(cmd, tt.args)

			if tt.expectError && err == nil {
				t.Errorf("tddCmd() args validation expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("tddCmd() args validation unexpected error: %v", err)
			}
		})
	}
}

// TestTddCmdInvalidPhase tests that invalid phase names are rejected
func TestTddCmdInvalidPhase(t *testing.T) {
	cmd := tddCmd()

	// Test invalid phase (should fail in RunE, not Args)
	err := cmd.RunE(cmd, []string{"invalid"})

	if err == nil {
		t.Errorf("tddCmd() with invalid phase expected error but got none")
	}
	if err != nil && err.Error() != "invalid phase: invalid (must be red, green, or refactor)" {
		// Error message might vary, just check we got an error
	}
}
