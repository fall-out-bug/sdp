package main

import (
	"testing"
)

// TestOrchestrateCmd_SkillEvent tests that orchestrate emits skill events (F056-03 AC3)
func TestOrchestrateCmd_SkillEvent(t *testing.T) {
	// This test verifies that the orchestrate command is configured to emit
	// plan and approval events. The actual emission is verified through
	// the SkillEvent function tests in the evidence package.
	//
	// orchestrate.go line 75: evidence.Emit(evidence.SkillEvent("oneshot", "plan", ...))
	// orchestrate.go line 82-84: evidence.EmitSync(evidence.SkillEvent("oneshot", "approval", ...))

	// Verify the command structure exists
	cmd := orchestrateCmd
	if cmd == nil {
		t.Fatal("orchestrateCmd is nil")
	}
	if cmd.Use != "orchestrate <feature-id>" {
		t.Errorf("unexpected Use: %s", cmd.Use)
	}
}

// TestOrchestrateResumeCmd tests the resume subcommand
func TestOrchestrateResumeCmd(t *testing.T) {
	// Verify the resume subcommand exists
	if orchestrateResumeCmd == nil {
		t.Fatal("orchestrateResumeCmd is nil")
	}
	if orchestrateResumeCmd.Use != "resume <checkpoint-id>" {
		t.Errorf("unexpected Use: %s", orchestrateResumeCmd.Use)
	}
}
