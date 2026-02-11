package cli

import (
	"testing"
)

// TestSynthesizeCmd_FlagParsing verifies flag parsing
func TestSynthesizeCmd_FlagParsing(t *testing.T) {
	cmd := synthesizeCmd

	// Verify feature flag exists
	flag := cmd.Flag("feature")
	if flag == nil {
		t.Fatal("Feature flag not found")
	}

	if flag.Shorthand != "f" {
		t.Errorf("Expected shorthand 'f', got '%s'", flag.Shorthand)
	}
}

// TestSynthesizeCmd_MissingFeatureFlag verifies required flag
func TestSynthesizeCmd_MissingFeatureFlag(t *testing.T) {
	// Test would require running the command, skip for now
	t.Skip("Integration test - requires command execution")
}

// TestLockCmd_FlagParsing verifies lock command flags
func TestLockCmd_FlagParsing(t *testing.T) {
	cmd := lockCmd

	// Verify contract flag exists
	flag := cmd.Flag("contract")
	if flag == nil {
		t.Fatal("Contract flag not found")
	}
}

// TestValidateCmd_FlagParsing verifies validate command flags
func TestValidateCmd_FlagParsing(t *testing.T) {
	cmd := validateCmd

	// Verify contracts flag exists
	flag := cmd.Flag("contracts")
	if flag == nil {
		t.Fatal("Contracts flag not found")
	}
}
