package model

import (
	"sort"
)

// Profile represents a model configuration
type Profile struct {
	Name      string  // Profile name (e.g., "fast", "quality")
	ModelID   string  // Actual model identifier
	Speed     string  // "fast", "balanced", "quality"
	MaxTokens int     // Maximum context length
	CostPer1K float64 // Cost per 1K tokens
}

// RoutingRule defines when to use a specific profile
type RoutingRule struct {
	TaskType    string // Type of task (e.g., "planning", "code", "review")
	MinComplex  int    // Minimum complexity (0-10)
	MaxComplex  int    // Maximum complexity (0-10)
	Profile     string // Profile to use
	Priority    int    // Rule priority (higher = more important)
}

// Router handles model selection
type Router struct {
	profiles map[string]Profile
	rules    []RoutingRule
	default_ string // default profile
}

// NewRouter creates a new router with default configuration
func NewRouter() *Router {
	r := &Router{
		profiles: make(map[string]Profile),
		rules:    []RoutingRule{},
		default_: "balanced",
	}

	// Add default profiles
	r.AddProfile(Profile{
		Name: "fast", ModelID: "claude-sonnet", Speed: "fast",
		MaxTokens: 200000, CostPer1K: 0.003,
	})
	r.AddProfile(Profile{
		Name: "balanced", ModelID: "claude-sonnet", Speed: "balanced",
		MaxTokens: 200000, CostPer1K: 0.003,
	})
	r.AddProfile(Profile{
		Name: "quality", ModelID: "claude-opus", Speed: "quality",
		MaxTokens: 200000, CostPer1K: 0.015,
	})

	// Add default rules
	r.AddRule(RoutingRule{TaskType: "planning", MinComplex: 7, MaxComplex: 10, Profile: "quality", Priority: 10})
	r.AddRule(RoutingRule{TaskType: "planning", MinComplex: 0, MaxComplex: 6, Profile: "balanced", Priority: 5})
	r.AddRule(RoutingRule{TaskType: "code", MinComplex: 0, MaxComplex: 10, Profile: "balanced", Priority: 5})
	r.AddRule(RoutingRule{TaskType: "review", MinComplex: 0, MaxComplex: 10, Profile: "quality", Priority: 10})
	r.AddRule(RoutingRule{TaskType: "debug", MinComplex: 0, MaxComplex: 10, Profile: "fast", Priority: 10})
	r.AddRule(RoutingRule{TaskType: "quick", MinComplex: 0, MaxComplex: 10, Profile: "fast", Priority: 5})

	return r
}

// AddProfile adds a model profile
func (r *Router) AddProfile(p Profile) {
	r.profiles[p.Name] = p
}

// AddRule adds a routing rule
func (r *Router) AddRule(rule RoutingRule) {
	r.rules = append(r.rules, rule)
	// Sort by priority descending
	sort.Slice(r.rules, func(i, j int) bool {
		return r.rules[i].Priority > r.rules[j].Priority
	})
}

// SetDefault sets the default profile
func (r *Router) SetDefault(profile string) {
	r.default_ = profile
}

// SelectModel selects the best model for a task
func (r *Router) SelectModel(taskType string, complexity int) Profile {
	// Find matching rules
	for _, rule := range r.rules {
		if (rule.TaskType == taskType || rule.TaskType == "*") &&
			complexity >= rule.MinComplex && complexity <= rule.MaxComplex {
			if p, ok := r.profiles[rule.Profile]; ok {
				return p
			}
		}
	}

	// Fall back to default
	if p, ok := r.profiles[r.default_]; ok {
		return p
	}

	// Ultimate fallback
	return Profile{Name: "default", ModelID: "claude-sonnet", Speed: "balanced"}
}

// SelectBySpeed selects a model based on desired speed
func (r *Router) SelectBySpeed(speed string) Profile {
	for _, p := range r.profiles {
		if p.Speed == speed {
			return p
		}
	}
	return r.SelectModel("", 5)
}

// SelectByCost selects the cheapest model that meets requirements
func (r *Router) SelectByCost(maxCostPer1K float64) Profile {
	var best Profile
	var found bool

	for _, p := range r.profiles {
		if p.CostPer1K <= maxCostPer1K {
			if !found || p.CostPer1K < best.CostPer1K {
				best = p
				found = true
			}
		}
	}

	if !found {
		return r.SelectModel("", 5)
	}
	return best
}

// GetProfile returns a profile by name
func (r *Router) GetProfile(name string) (Profile, bool) {
	p, ok := r.profiles[name]
	return p, ok
}

// ListProfiles returns all available profiles
func (r *Router) ListProfiles() []Profile {
	profiles := make([]Profile, 0, len(r.profiles))
	for _, p := range r.profiles {
		profiles = append(profiles, p)
	}
	return profiles
}

// EstimateCost estimates the cost for a given token count and profile
func EstimateCost(tokens int, profile Profile) float64 {
	return float64(tokens) / 1000 * profile.CostPer1K
}
