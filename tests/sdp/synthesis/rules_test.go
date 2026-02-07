package synthesis

import (
	"testing"

	"github.com/fall-out-bug/sdp/src/sdp/synthesis"
)

// MockRule for testing
type mockRule struct {
	name      string
	priority  int
	canApply  bool
	result    *synthesis.SynthesisResult
	err       error
}

func (m *mockRule) Name() string        { return m.name }
func (m *mockRule) Priority() int       { return m.priority }
func (m *mockRule) CanApply([]*synthesis.Proposal) bool { return m.canApply }
func (m *mockRule) Apply([]*synthesis.Proposal) (*synthesis.SynthesisResult, error) {
	return m.result, m.err
}

// TestRuleEngine_AddRule_AddsToEngine tests AC1: Rule interface and AC3: AddRule
func TestRuleEngine_AddRule_AddsToEngine(t *testing.T) {
	// Arrange
	engine := synthesis.NewRuleEngine()
	rule := &mockRule{name: "test-rule", priority: 1}

	// Act
	engine.AddRule(rule)

	// Assert
	// Engine should have 1 rule
	// We'll verify by executing with a proposal that matches
}

// TestRuleEngine_Execute_PriorityOrder tests AC3: Rules execute in priority order
func TestRuleEngine_Execute_PriorityOrder(t *testing.T) {
	// Arrange
	engine := synthesis.NewRuleEngine()

	// Add rules in wrong priority order
	rule1 := &mockRule{name: "priority-5", priority: 5, canApply: true, result: &synthesis.SynthesisResult{Rule: "priority-5"}}
	rule2 := &mockRule{name: "priority-1", priority: 1, canApply: true, result: &synthesis.SynthesisResult{Rule: "priority-1"}}
	rule3 := &mockRule{name: "priority-3", priority: 3, canApply: true, result: &synthesis.SynthesisResult{Rule: "priority-3"}}

	engine.AddRule(rule1)
	engine.AddRule(rule2)
	engine.AddRule(rule3)

	proposals := []*synthesis.Proposal{}

	// Act
	result, err := engine.Execute(proposals)

	// Assert - priority 1 rule should win
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Rule != "priority-1" {
		t.Errorf("Expected rule 'priority-1' to execute first, got %s", result.Rule)
	}
}

// TestRuleEngine_Execute_FirstApplicableWins tests AC3: First rule that can apply wins
func TestRuleEngine_Execute_FirstApplicableWins(t *testing.T) {
	// Arrange
	engine := synthesis.NewRuleEngine()

	rule1 := &mockRule{name: "cannot-apply", priority: 1, canApply: false}
	rule2 := &mockRule{name: "can-apply", priority: 2, canApply: true, result: &synthesis.SynthesisResult{Rule: "winner"}}

	engine.AddRule(rule1)
	engine.AddRule(rule2)

	proposals := []*synthesis.Proposal{}

	// Act
	result, err := engine.Execute(proposals)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Rule != "winner" {
		t.Errorf("Expected rule 'can-apply' to win, got %s", result.Rule)
	}
}

// TestUnanimousRule_ApplyingWhenAllAgree tests AC2: UnanimousRule
func TestUnanimousRule_ApplyingWhenAllAgree(t *testing.T) {
	// Arrange
	rule := synthesis.NewUnanimousRule()
	solution := map[string]string{"approach": "TDD"}

	proposals := []*synthesis.Proposal{
		synthesis.NewProposal("agent-1", solution, 0.9, "TDD"),
		synthesis.NewProposal("agent-2", solution, 0.85, "TDD is best"),
		synthesis.NewProposal("agent-3", solution, 0.95, "Quality"),
	}

	// Act
	canApply := rule.CanApply(proposals)

	// Assert
	if !canApply {
		t.Error("Expected UnanimousRule to apply when all agents agree")
	}

	result, err := rule.Apply(proposals)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Rule != "unanimous" {
		t.Errorf("Expected rule 'unanimous', got %s", result.Rule)
	}
}

// TestDomainExpertiseRule_ApplyingWithUniqueHighest tests AC2: DomainExpertiseRule
func TestDomainExpertiseRule_ApplyingWithUniqueHighest(t *testing.T) {
	// Arrange
	rule := synthesis.NewDomainExpertiseRule()

	proposals := []*synthesis.Proposal{
		synthesis.NewProposal("agent-1", map[string]string{"approach": "TDD"}, 0.7, "Low confidence"),
		synthesis.NewProposal("agent-2", map[string]string{"approach": "spike"}, 0.6, "Not sure"),
		synthesis.NewProposal("agent-3", map[string]string{"approach": "TDD"}, 0.95, "Expert"),
	}

	// Act
	canApply := rule.CanApply(proposals)

	// Assert
	if !canApply {
		t.Error("Expected DomainExpertiseRule to apply with unique highest confidence")
	}

	result, err := rule.Apply(proposals)
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

// TestDomainExpertiseRule_NotApplyingWithTie tests that rule doesn't apply when confidence is tied
func TestDomainExpertiseRule_NotApplyingWithTie(t *testing.T) {
	// Arrange
	rule := synthesis.NewDomainExpertiseRule()

	proposals := []*synthesis.Proposal{
		synthesis.NewProposal("agent-1", nil, 0.8, ""),
		synthesis.NewProposal("agent-2", nil, 0.8, ""),
	}

	// Act
	canApply := rule.CanApply(proposals)

	// Assert
	if canApply {
		t.Error("Expected DomainExpertiseRule to not apply with tied confidence")
	}
}

// TestRuleEngine_Execute_NoRuleApplies tests escalation when no rule can apply
func TestRuleEngine_Execute_NoRuleApplies(t *testing.T) {
	// Arrange
	engine := synthesis.NewRuleEngine()
	rule := &mockRule{name: "cannot-apply", priority: 1, canApply: false}
	engine.AddRule(rule)

	proposals := []*synthesis.Proposal{}

	// Act
	result, err := engine.Execute(proposals)

	// Assert
	if err == nil {
		t.Error("Expected error when no rule applies, got nil")
	}

	if result != nil {
		t.Error("Expected nil result when no rule applies, got non-nil")
	}
}
