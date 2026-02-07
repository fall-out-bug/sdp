package synthesis

import (
	"sort"
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
	sort.Slice(e.rules, func(i, j int) bool {
		return e.rules[i].Priority() < e.rules[j].Priority()
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
