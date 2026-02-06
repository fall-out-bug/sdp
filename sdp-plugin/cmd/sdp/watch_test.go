package main

import (
	"os"
	"testing"
)

// TestWatchCmd tests the watch command initialization
func TestWatchCmd(t *testing.T) {
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	// Create a simple go file to watch
	if err := os.WriteFile("test.go", []byte("package main\n\nfunc main() {}"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd := watchCmd()

	// Set quiet mode to avoid output
	if err := cmd.Flags().Set("quiet", "true"); err != nil {
		t.Fatalf("Failed to set quiet flag: %v", err)
	}

	// Test that command can be created and flags work
	if cmd == nil {
		t.Error("watchCmd() returned nil")
	}

	// We can't actually run the watch command because it blocks waiting for files
	// Just verify the command structure is correct
	if cmd.Use != "watch" {
		t.Errorf("watchCmd() has wrong use: %s", cmd.Use)
	}

	// Test flag parsing
	include, _ := cmd.Flags().GetStringSlice("include")
	if len(include) == 0 {
		t.Error("watchCmd() default include patterns not set")
	}

	exclude, _ := cmd.Flags().GetStringSlice("exclude")
	if len(exclude) == 0 {
		t.Error("watchCmd() default exclude patterns not set")
	}
}
