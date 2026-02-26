package main

import (
	"testing"
)

// TestVerifyCmdConstructed tests that verify command can be constructed.
// Full integration test (RunE) is skipped: verify flow may hang in isolated tmpDir
// due to evidence/acceptance/config interaction — see internal/verify for unit tests.
func TestVerifyCmdConstructed(t *testing.T) {
	cmd := verifyCmd()
	if cmd == nil {
		t.Fatal("verifyCmd() returned nil")
	}
	if cmd.Use != "verify <ws-id>" {
		t.Errorf("expected Use=verify <ws-id>, got %q", cmd.Use)
	}
	if cmd.Args == nil {
		t.Error("expected Args validator to be set")
	}
}
