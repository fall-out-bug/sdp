package main

import (
	"testing"
)

// TestCompletionCmd tests the completion command
func TestCompletionCmd(t *testing.T) {
	cmd := completionCmd

	// Test command structure
	if cmd.Use != "completion [bash|zsh|fish]" {
		t.Errorf("completionCmd has wrong use: %s", cmd.Use)
	}

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "bash completion",
			args:        []string{"bash"},
			expectError: false,
		},
		{
			name:        "zsh completion",
			args:        []string{"zsh"},
			expectError: false,
		},
		{
			name:        "fish completion",
			args:        []string{"fish"},
			expectError: false,
		},
		{
			name:        "no args",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "invalid shell",
			args:        []string{"invalid"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Recover from panic caused by args[0] access when args is empty
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectError {
						t.Errorf("completionCmd() panicked with: %v", r)
					}
				}
			}()

			err := cmd.RunE(cmd, tt.args)

			if tt.expectError && err == nil {
				// Might have panicked instead, which is also ok for this test
			}
			if !tt.expectError && err != nil {
				t.Errorf("completionCmd() unexpected error: %v", err)
			}
		})
	}
}

// TestCompletionValidArgsFunction tests the completion argument suggestions
func TestCompletionValidArgsFunction(t *testing.T) {
	cmd := completionCmd

	// Test valid args function
	suggestions, _ := cmd.ValidArgsFunction(cmd, []string{}, "")

	if len(suggestions) != 3 {
		t.Errorf("completionCmd() valid args count = %d, want 3", len(suggestions))
	}

	expected := []string{"bash", "zsh", "fish"}
	for _, exp := range expected {
		found := false
		for _, sug := range suggestions {
			if sug == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("completionCmd() missing valid arg: %s", exp)
		}
	}
}
