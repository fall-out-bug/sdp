package synthesis

import (
	"encoding/json"
	"testing"

	"github.com/fall-out-bug/sdp/src/sdp/synthesis"
)

// TestProposal_New_CreatesValidProposal tests AC1: Proposal data structure
func TestProposal_New_CreatesValidProposal(t *testing.T) {
	// Arrange
	agentID := "implementer"
	solution := map[string]interface{}{
		"approach": "TDD",
		"files":    []string{"main.go"},
	}
	confidence := 0.95
	reasoning := "Test-driven development ensures quality"

	// Act
	proposal := synthesis.NewProposal(agentID, solution, confidence, reasoning)

	// Assert
	if proposal == nil {
		t.Fatal("Expected proposal to be created, got nil")
	}

	if proposal.AgentID != agentID {
		t.Errorf("Expected agent ID %s, got %s", agentID, proposal.AgentID)
	}

	if proposal.Confidence != confidence {
		t.Errorf("Expected confidence %f, got %f", confidence, proposal.Confidence)
	}

	if proposal.Reasoning != reasoning {
		t.Errorf("Expected reasoning %s, got %s", reasoning, proposal.Reasoning)
	}

	// Verify timestamp is set
	if proposal.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set, got zero time")
	}

	// Verify solution is stored correctly
	if proposal.Solution == nil {
		t.Error("Expected solution to be stored, got nil")
	}
}

// TestProposal_JSONSerialization tests AC1: JSON serialization
func TestProposal_JSONSerialization(t *testing.T) {
	// Arrange
	proposal := synthesis.NewProposal(
		"agent-1",
		map[string]interface{}{"approach": "TDD"},
		0.9,
		"Test-driven approach",
	)

	// Act - serialize
	data, err := json.Marshal(proposal)
	if err != nil {
		t.Fatalf("Failed to marshal proposal: %v", err)
	}

	// Assert - verify JSON contains required fields
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal proposal: %v", err)
	}

	requiredFields := []string{"agent_id", "solution", "confidence", "reasoning", "timestamp"}
	for _, field := range requiredFields {
		if _, exists := raw[field]; !exists {
			t.Errorf("Missing required field in JSON: %s", field)
		}
	}

	// Act - deserialize
	restored := &synthesis.Proposal{}
	if err := json.Unmarshal(data, restored); err != nil {
		t.Fatalf("Failed to unmarshal proposal: %v", err)
	}

	// Assert - verify restored proposal matches original
	if restored.AgentID != proposal.AgentID {
		t.Errorf("Expected agent ID %s, got %s", proposal.AgentID, restored.AgentID)
	}

	if restored.Confidence != proposal.Confidence {
		t.Errorf("Expected confidence %f, got %f", proposal.Confidence, restored.Confidence)
	}

	if restored.Reasoning != proposal.Reasoning {
		t.Errorf("Expected reasoning %s, got %s", proposal.Reasoning, restored.Reasoning)
	}
}

// TestProposal_Equals_ComparesProposals tests AC1: Proposal comparison
func TestProposal_Equals_ComparesProposals(t *testing.T) {
	// Arrange
	p1 := synthesis.NewProposal("agent-1", map[string]interface{}{"x": 1}, 0.9, "reasoning")
	p2 := synthesis.NewProposal("agent-1", map[string]interface{}{"x": 1}, 0.9, "reasoning")
	p3 := synthesis.NewProposal("agent-2", map[string]interface{}{"x": 1}, 0.9, "reasoning")

	// Act & Assert
	if !p1.Equals(p2) {
		t.Error("Expected p1 to equal p2 (same content)")
	}

	if p1.Equals(p3) {
		t.Error("Expected p1 to not equal p3 (different agent ID)")
	}
}

// TestProposal_IsHigherConfidence tests confidence comparison
func TestProposal_IsHigherConfidence(t *testing.T) {
	// Arrange
	p1 := synthesis.NewProposal("agent-1", nil, 0.9, "")
	p2 := synthesis.NewProposal("agent-2", nil, 0.95, "")

	// Act & Assert
	if !p2.IsHigherConfidenceThan(p1) {
		t.Error("Expected p2 to have higher confidence than p1")
	}

	if p1.IsHigherConfidenceThan(p2) {
		t.Error("Expected p1 to not have higher confidence than p2")
	}
}
