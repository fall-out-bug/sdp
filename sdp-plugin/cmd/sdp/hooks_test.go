package main

import (
	"testing"
)

// TestHooksCmd tests the hooks command structure
func TestHooksCmd(t *testing.T) {
	cmd := hooksCmd()

	// Test command structure
	if cmd.Use != "hooks" {
		t.Errorf("hooksCmd() has wrong use: %s", cmd.Use)
	}

	// Test subcommands
	subcommands := []string{"install", "uninstall"}
	for _, sub := range subcommands {
		found := false
		for _, c := range cmd.Commands() {
			if c.Name() == sub {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("hooksCmd() missing subcommand: %s", sub)
		}
	}
}
