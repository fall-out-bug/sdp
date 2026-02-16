package synthesis

import (
	"errors"
	"testing"
)

// TestNewRuleEngine verifies engine creation
func TestNewRuleEngine(t *testing.T) {
	engine := NewRuleEngine()
	if engine == nil {
		t.Fatal("NewRuleEngine returned nil")
	}
	if engine.rules == nil {
		t.Error("rules slice not initialized")
	}
}

// TestRuleEngine_AddRule verifies adding rules
func TestRuleEngine_AddRule(t *testing.T) {
	engine := NewRuleEngine()
	engine.AddRule(NewDomainExpertiseRule())
	engine.AddRule(NewUnanimousRule())

	if len(engine.rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(engine.rules))
	}
}

// TestRuleEngine_AddRule_SortedByPriority verifies sorting by priority
func TestRuleEngine_AddRule_SortedByPriority(t *testing.T) {
	engine := NewRuleEngine()
	// Add in wrong order
	engine.AddRule(NewDomainExpertiseRule()) // priority 2
	engine.AddRule(NewUnanimousRule())       // priority 1

	rules := engine.GetRules()
	if rules[0].Priority() != 1 {
		t.Error("unanimous rule should be first (priority 1)")
	}
	if rules[1].Priority() != 2 {
		t.Error("domain expertise should be second (priority 2)")
	}
}

// TestRuleEngine_Execute_UnanimousRule verifies unanimous rule execution
func TestRuleEngine_Execute_UnanimousRule(t *testing.T) {
	engine := NewRuleEngine()
	engine.AddRule(NewUnanimousRule())
	engine.AddRule(NewDomainExpertiseRule())

	proposals := []*Proposal{
		NewProposal("agent-1", "same", 0.9, ""),
		NewProposal("agent-2", "same", 0.8, ""),
	}

	result, err := engine.Execute(proposals)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result.Rule != "unanimous" {
		t.Errorf("expected unanimous rule, got %s", result.Rule)
	}
}

// TestRuleEngine_Execute_DomainExpertiseRule verifies domain expertise execution
func TestRuleEngine_Execute_DomainExpertiseRule(t *testing.T) {
	engine := NewRuleEngine()
	engine.AddRule(NewUnanimousRule())
	engine.AddRule(NewDomainExpertiseRule())

	proposals := []*Proposal{
		NewProposal("agent-1", "solution-1", 0.9, ""),
		NewProposal("agent-2", "solution-2", 0.7, ""),
	}

	result, err := engine.Execute(proposals)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result.Rule != "domain_expertise" {
		t.Errorf("expected domain_expertise rule, got %s", result.Rule)
	}
}

// TestRuleEngine_Execute_NoRules verifies error when no rules apply
func TestRuleEngine_Execute_NoRules(t *testing.T) {
	engine := NewRuleEngine()
	// Empty engine

	proposals := []*Proposal{
		NewProposal("agent-1", "solution", 0.9, ""),
	}

	_, err := engine.Execute(proposals)
	if err == nil {
		t.Error("expected error for no rules")
	}
	if !errors.Is(err, ErrCannotSynthesize) {
		t.Errorf("expected ErrCannotSynthesize, got %v", err)
	}
}

// TestRuleEngine_Execute_NoMatchingRule verifies error when no rule can apply
func TestRuleEngine_Execute_NoMatchingRule(t *testing.T) {
	engine := NewRuleEngine()
	// Only unanimous rule - won't apply to single proposal
	engine.AddRule(NewUnanimousRule())

	proposals := []*Proposal{
		NewProposal("agent-1", "solution", 0.9, ""),
	}

	_, err := engine.Execute(proposals)
	if err == nil {
		t.Error("expected error when no rule can apply")
	}
}

// TestRuleEngine_GetRules verifies getting rules
func TestRuleEngine_GetRules(t *testing.T) {
	engine := NewRuleEngine()
	engine.AddRule(NewUnanimousRule())

	rules := engine.GetRules()
	if len(rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(rules))
	}
}

// TestDefaultRuleEngine verifies default engine creation
func TestDefaultRuleEngine(t *testing.T) {
	engine := DefaultRuleEngine()
	if engine == nil {
		t.Fatal("DefaultRuleEngine returned nil")
	}

	rules := engine.GetRules()
	if len(rules) != 2 {
		t.Errorf("expected 2 default rules, got %d", len(rules))
	}
}
