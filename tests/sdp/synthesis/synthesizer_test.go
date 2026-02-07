package synthesis

import (
	"errors"
	"testing"

	"github.com/fall-out-bug/sdp/src/sdp/synthesis"
)

// TestSynthesizer_New_CreatesEmptySynthesizer tests AC3: Synthesizer interface
func TestSynthesizer_New_CreatesEmptySynthesizer(t *testing.T) {
	// Act
	synth := synthesis.NewSynthesizer()

	// Assert
	if synth == nil {
		t.Fatal("Expected synthesizer to be created, got nil")
	}

	proposals := synth.GetProposals()
	if len(proposals) != 0 {
		t.Errorf("Expected 0 proposals, got %d", len(proposals))
	}
}

// TestSynthesizer_AddProposal_AddsToSynthesizer tests AC3: AddProposal
func TestSynthesizer_AddProposal_AddsToSynthesizer(t *testing.T) {
	// Arrange
	synth := synthesis.NewSynthesizer()
	proposal := synthesis.NewProposal("agent-1", map[string]string{"approach": "TDD"}, 0.9, "reasoning")

	// Act
	synth.AddProposal(proposal)

	// Assert
	proposals := synth.GetProposals()
	if len(proposals) != 1 {
		t.Fatalf("Expected 1 proposal, got %d", len(proposals))
	}

	if proposals[0].AgentID != "agent-1" {
		t.Errorf("Expected agent ID agent-1, got %s", proposals[0].AgentID)
	}
}

// TestSynthesizer_Clear_RemovesAllProposals tests AC3: Clear
func TestSynthesizer_Clear_RemovesAllProposals(t *testing.T) {
	// Arrange
	synth := synthesis.NewSynthesizer()
	synth.AddProposal(synthesis.NewProposal("agent-1", nil, 0.9, ""))
	synth.AddProposal(synthesis.NewProposal("agent-2", nil, 0.8, ""))

	// Act
	synth.Clear()

	// Assert
	proposals := synth.GetProposals()
	if len(proposals) != 0 {
		t.Errorf("Expected 0 proposals after clear, got %d", len(proposals))
	}
}

// TestSynthesizer_Synthesize_UnanimousAgreement tests AC2: Rule 1 - Unanimous
func TestSynthesizer_Synthesize_UnanimousAgreement(t *testing.T) {
	// Arrange
	synth := synthesis.NewSynthesizer()
	solution := map[string]string{"approach": "TDD"}

	// All agents agree on same solution
	synth.AddProposal(synthesis.NewProposal("agent-1", solution, 0.9, "Use TDD"))
	synth.AddProposal(synthesis.NewProposal("agent-2", solution, 0.85, "TDD is best"))
	synth.AddProposal(synthesis.NewProposal("agent-3", solution, 0.95, "Quality first"))

	// Act
	result, err := synth.Synthesize()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Rule != "unanimous" {
		t.Errorf("Expected rule 'unanimous', got %s", result.Rule)
	}

	if result.Solution == nil {
		t.Error("Expected solution to be set")
	}
}

// TestSynthesizer_Synthesize_DomainExpertise tests AC2: Rule 2 - Domain expertise
func TestSynthesizer_Synthesize_DomainExpertise(t *testing.T) {
	// Arrange
	synth := synthesis.NewSynthesizer()

	// Different solutions, agent-3 has highest confidence
	synth.AddProposal(synthesis.NewProposal("agent-1", map[string]string{"approach": "TDD"}, 0.7, "Low confidence"))
	synth.AddProposal(synthesis.NewProposal("agent-2", map[string]string{"approach": "spike"}, 0.6, "Not sure"))
	synth.AddProposal(synthesis.NewProposal("agent-3", map[string]string{"approach": "TDD"}, 0.95, "Expert in this"))

	// Act
	result, err := synth.Synthesize()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Rule != "domain_expertise" {
		t.Errorf("Expected rule 'domain_expertise', got %s", result.Rule)
	}

	if result.WinningAgent != "agent-3" {
		t.Errorf("Expected winning agent 'agent-3', got %s", result.WinningAgent)
	}
}

// TestSynthesizer_Synthesize_EscalatesToHuman tests AC2: Rule 5 - Escalate
func TestSynthesizer_Synthesize_EscalatesToHuman(t *testing.T) {
	// Arrange
	synth := synthesis.NewSynthesizer()

	// Conflicting solutions with equal confidence - cannot resolve
	synth.AddProposal(synthesis.NewProposal("agent-1", map[string]string{"approach": "TDD"}, 0.8, "Use TDD"))
	synth.AddProposal(synthesis.NewProposal("agent-2", map[string]string{"approach": "spike"}, 0.8, "Use spike"))

	// Act
	result, err := synth.Synthesize()

	// Assert
	if err == nil {
		t.Error("Expected error for unresolvable conflict, got nil")
	}

	if !errors.Is(err, synthesis.ErrCannotSynthesize) {
		t.Errorf("Expected ErrCannotSynthesize, got %v", err)
	}

	if result != nil {
		t.Error("Expected nil result for escalation, got non-nil")
	}
}

// TestSynthesizer_DetectConflict_NoConflict tests AC4: No conflict detection
func TestSynthesizer_DetectConflict_NoConflict(t *testing.T) {
	// Arrange
	synth := synthesis.NewSynthesizer()
	solution := map[string]string{"approach": "TDD"}
	synth.AddProposal(synthesis.NewProposal("agent-1", solution, 0.9, ""))
	synth.AddProposal(synthesis.NewProposal("agent-2", solution, 0.85, ""))

	// Act
	conflict := synth.DetectConflict()

	// Assert
	if conflict != synthesis.NoConflict {
		t.Errorf("Expected NoConflict, got %v", conflict)
	}
}

// TestSynthesizer_DetectConflict_MajorConflict tests AC4: Major conflict detection
func TestSynthesizer_DetectConflict_MajorConflict(t *testing.T) {
	// Arrange
	synth := synthesis.NewSynthesizer()
	synth.AddProposal(synthesis.NewProposal("agent-1", map[string]string{"approach": "TDD"}, 0.9, ""))
	synth.AddProposal(synthesis.NewProposal("agent-2", map[string]string{"approach": "spike"}, 0.8, ""))

	// Act
	conflict := synth.DetectConflict()

	// Assert
	if conflict != synthesis.MajorConflict {
		t.Errorf("Expected MajorConflict, got %v", conflict)
	}
}
