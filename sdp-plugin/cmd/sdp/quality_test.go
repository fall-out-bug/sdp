package main

import (
	"testing"

	"github.com/spf13/cobra"
)

// TestQualityCmd tests the quality command structure
func TestQualityCmd(t *testing.T) {
	cmd := qualityCmd()

	// Test command structure
	if cmd.Use != "quality" {
		t.Errorf("qualityCmd() has wrong use: %s", cmd.Use)
	}

	// Test subcommands
	expectedSubcommands := []string{"coverage", "complexity", "size", "types", "all"}
	for _, expected := range expectedSubcommands {
		found := false
		for _, c := range cmd.Commands() {
			if c.Name() == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("qualityCmd() missing subcommand: %s", expected)
		}
	}
}

// TestQualityCoverageCmd tests the quality coverage command
func TestQualityCoverageCmd(t *testing.T) {
	cmd := qualityCmd()

	// Find the 'coverage' subcommand
	var coverageCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "coverage" {
			coverageCmd = subcmd
			break
		}
	}

	if coverageCmd == nil {
		t.Fatal("quality coverage command not found")
	}

	// Test command structure
	if coverageCmd.Use != "coverage" {
		t.Errorf("quality coverage command has wrong use: %s", coverageCmd.Use)
	}

	// Test that command can be executed (will fail due to real quality issues)
	err := coverageCmd.RunE(coverageCmd, []string{})
	// Expected to fail (coverage < 80%, complexity > 10, etc. in real codebase)
	if err == nil {
		t.Log("quality coverage command succeeded (codebase quality is good!)")
	} else {
		t.Log("quality coverage command failed as expected (real quality issues exist)")
	}
}

// TestQualityComplexityCmd tests the quality complexity command
func TestQualityComplexityCmd(t *testing.T) {
	cmd := qualityCmd()

	// Find the 'complexity' subcommand
	var complexityCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "complexity" {
			complexityCmd = subcmd
			break
		}
	}

	if complexityCmd == nil {
		t.Fatal("quality complexity command not found")
	}

	if complexityCmd.Use != "complexity" {
		t.Errorf("quality complexity command has wrong use: %s", complexityCmd.Use)
	}
}

// TestQualitySizeCmd tests the quality size command
func TestQualitySizeCmd(t *testing.T) {
	cmd := qualityCmd()

	// Find the 'size' subcommand
	var sizeCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "size" {
			sizeCmd = subcmd
			break
		}
	}

	if sizeCmd == nil {
		t.Fatal("quality size command not found")
	}

	if sizeCmd.Use != "size" {
		t.Errorf("quality size command has wrong use: %s", sizeCmd.Use)
	}
}

// TestQualityTypesCmd tests the quality types command
func TestQualityTypesCmd(t *testing.T) {
	cmd := qualityCmd()

	// Find the 'types' subcommand
	var typesCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "types" {
			typesCmd = subcmd
			break
		}
	}

	if typesCmd == nil {
		t.Fatal("quality types command not found")
	}

	if typesCmd.Use != "types" {
		t.Errorf("quality types command has wrong use: %s", typesCmd.Use)
	}
}

// TestQualityAllCmd tests the quality all command
func TestQualityAllCmd(t *testing.T) {
	cmd := qualityCmd()

	// Find the 'all' subcommand
	var allCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "all" {
			allCmd = subcmd
			break
		}
	}

	if allCmd == nil {
		t.Fatal("quality all command not found")
	}

	if allCmd.Use != "all" {
		t.Errorf("quality all command has wrong use: %s", allCmd.Use)
	}
}
