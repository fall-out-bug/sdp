package synthesis

import (
	"testing"
)

// TestUnanimousRule_Name verifies rule name
func TestUnanimousRule_Name(t *testing.T) {
	rule := NewUnanimousRule()
	if rule.Name() != "unanimous" {
		t.Errorf("expected name 'unanimous', got %s", rule.Name())
	}
}

// TestUnanimousRule_Priority verifies rule priority
func TestUnanimousRule_Priority(t *testing.T) {
	rule := NewUnanimousRule()
	if rule.Priority() != 1 {
		t.Errorf("expected priority 1, got %d", rule.Priority())
	}
}

// TestUnanimousRule_CanApply_True verifies unanimous case
func TestUnanimousRule_CanApply_True(t *testing.T) {
	rule := NewUnanimousRule()
	proposals := []*Proposal{
		NewProposal("agent-1", "same", 0.9, ""),
		NewProposal("agent-2", "same", 0.8, ""),
	}

	if !rule.CanApply(proposals) {
		t.Error("should apply for unanimous proposals")
	}
}

// TestUnanimousRule_CanApply_False verifies non-unanimous case
func TestUnanimousRule_CanApply_False(t *testing.T) {
	rule := NewUnanimousRule()
	proposals := []*Proposal{
		NewProposal("agent-1", "solution-1", 0.9, ""),
		NewProposal("agent-2", "solution-2", 0.8, ""),
	}

	if rule.CanApply(proposals) {
		t.Error("should not apply for different proposals")
	}
}

// TestUnanimousRule_CanApply_SingleProposal verifies single proposal case
func TestUnanimousRule_CanApply_SingleProposal(t *testing.T) {
	rule := NewUnanimousRule()
	proposals := []*Proposal{
		NewProposal("agent-1", "solution", 0.9, ""),
	}

	if rule.CanApply(proposals) {
		t.Error("should not apply for single proposal")
	}
}

// TestUnanimousRule_Apply verifies rule application
func TestUnanimousRule_Apply(t *testing.T) {
	rule := NewUnanimousRule()
	proposals := []*Proposal{
		NewProposal("agent-1", "solution", 0.9, ""),
		NewProposal("agent-2", "solution", 0.8, ""),
	}

	result, err := rule.Apply(proposals)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	if result.Solution != "solution" {
		t.Errorf("expected solution 'solution', got %v", result.Solution)
	}
	if result.Rule != "unanimous" {
		t.Errorf("expected rule 'unanimous', got %s", result.Rule)
	}
}

// TestDomainExpertiseRule_Name verifies rule name
func TestDomainExpertiseRule_Name(t *testing.T) {
	rule := NewDomainExpertiseRule()
	if rule.Name() != "domain_expertise" {
		t.Errorf("expected name 'domain_expertise', got %s", rule.Name())
	}
}

// TestDomainExpertiseRule_Priority verifies rule priority
func TestDomainExpertiseRule_Priority(t *testing.T) {
	rule := NewDomainExpertiseRule()
	if rule.Priority() != 2 {
		t.Errorf("expected priority 2, got %d", rule.Priority())
	}
}

// TestDomainExpertiseRule_CanApply_True verifies unique highest confidence
func TestDomainExpertiseRule_CanApply_True(t *testing.T) {
	rule := NewDomainExpertiseRule()
	proposals := []*Proposal{
		NewProposal("agent-1", "s1", 0.9, ""),
		NewProposal("agent-2", "s2", 0.7, ""),
	}

	if !rule.CanApply(proposals) {
		t.Error("should apply when unique highest confidence")
	}
}

// TestDomainExpertiseRule_CanApply_False_Tie verifies tie case
func TestDomainExpertiseRule_CanApply_False_Tie(t *testing.T) {
	rule := NewDomainExpertiseRule()
	proposals := []*Proposal{
		NewProposal("agent-1", "s1", 0.9, ""),
		NewProposal("agent-2", "s2", 0.9, ""),
	}

	if rule.CanApply(proposals) {
		t.Error("should not apply when tied confidence")
	}
}

// TestDomainExpertiseRule_CanApply_Empty verifies empty case
func TestDomainExpertiseRule_CanApply_Empty(t *testing.T) {
	rule := NewDomainExpertiseRule()
	proposals := []*Proposal{}

	if rule.CanApply(proposals) {
		t.Error("should not apply for empty proposals")
	}
}

// TestDomainExpertiseRule_Apply verifies rule application
func TestDomainExpertiseRule_Apply(t *testing.T) {
	rule := NewDomainExpertiseRule()
	proposals := []*Proposal{
		NewProposal("agent-1", "best-solution", 0.95, ""),
		NewProposal("agent-2", "other-solution", 0.7, ""),
	}

	result, err := rule.Apply(proposals)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	if result.Solution != "best-solution" {
		t.Errorf("expected best solution, got %v", result.Solution)
	}
	if result.WinningAgent != "agent-1" {
		t.Errorf("expected agent-1 to win, got %s", result.WinningAgent)
	}
}

// TestSynthesisRule_Interface verifies interface compliance
func TestSynthesisRule_Interface(t *testing.T) {
	var _ SynthesisRule = NewUnanimousRule()
	var _ SynthesisRule = NewDomainExpertiseRule()
}
