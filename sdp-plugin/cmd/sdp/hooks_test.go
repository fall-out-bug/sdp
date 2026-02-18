package main

import (
	"testing"

	"github.com/spf13/cobra"
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
	var installCmdFound bool
	var installCmd *cobra.Command
	for _, sub := range subcommands {
		found := false
		for _, c := range cmd.Commands() {
			if c.Name() == sub {
				found = true
				if sub == "install" {
					installCmdFound = true
					installCmd = c
				}
				break
			}
		}
		if !found {
			t.Errorf("hooksCmd() missing subcommand: %s", sub)
		}
	}

	if !installCmdFound || installCmd == nil {
		t.Fatal("install subcommand not found")
	}

	if installCmd.Flags().Lookup("with-provenance") == nil {
		t.Error("hooks install missing --with-provenance flag")
	}
}
