package synthesis

import (
	"slices"
)

// RuleEngine executes synthesis rules in priority order
type RuleEngine struct {
	rules []SynthesisRule
}

// NewRuleEngine creates a new rule engine
func NewRuleEngine() *RuleEngine {
	return &RuleEngine{
		rules: make([]SynthesisRule, 0),
	}
}

// AddRule adds a rule to the engine
func (e *RuleEngine) AddRule(rule SynthesisRule) {
	e.rules = append(e.rules, rule)

	// Sort by priority (1 = highest)
	slices.SortFunc(e.rules, func(a, b SynthesisRule) int {
		switch {
		case a.Priority() < b.Priority():
			return -1
		case a.Priority() > b.Priority():
			return 1
		default:
			return 0
		}
	})
}

// Execute executes rules in priority order
// First rule that can apply wins
func (e *RuleEngine) Execute(proposals []*Proposal) (*SynthesisResult, error) {
	for _, rule := range e.rules {
		if rule.CanApply(proposals) {
			return rule.Apply(proposals)
		}
	}

	return nil, ErrCannotSynthesize
}

// GetRules returns all rules (for testing)
func (e *RuleEngine) GetRules() []SynthesisRule {
	return e.rules
}

// DefaultRuleEngine creates a rule engine with default rules
func DefaultRuleEngine() *RuleEngine {
	engine := NewRuleEngine()
	engine.AddRule(NewUnanimousRule())
	engine.AddRule(NewDomainExpertiseRule())
	return engine
}
