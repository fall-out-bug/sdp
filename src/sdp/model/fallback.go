package model

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// ErrAllModelsFailed is returned when all models in chain fail
var ErrAllModelsFailed = errors.New("all models in fallback chain failed")

// FallbackStats tracks fallback statistics
type FallbackStats struct {
	TotalCalls   int64
	Fallbacks    int64
	Successes    int64
	ModelErrors  map[string]int64
	mu           sync.Mutex
}

// NewFallbackStats creates new stats
func NewFallbackStats() *FallbackStats {
	return &FallbackStats{
		ModelErrors: make(map[string]int64),
	}
}

// RecordCall records a call attempt
func (s *FallbackStats) RecordCall() {
	atomic.AddInt64(&s.TotalCalls, 1)
}

// RecordFallback records a fallback
func (s *FallbackStats) RecordFallback() {
	atomic.AddInt64(&s.Fallbacks, 1)
}

// RecordSuccess records a success
func (s *FallbackStats) RecordSuccess() {
	atomic.AddInt64(&s.Successes, 1)
}

// RecordError records an error for a model
func (s *FallbackStats) RecordError(model string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ModelErrors[model]++
}

// GetStats returns current statistics
func (s *FallbackStats) GetStats() (total, fallbacks, successes int64) {
	return atomic.LoadInt64(&s.TotalCalls),
		atomic.LoadInt64(&s.Fallbacks),
		atomic.LoadInt64(&s.Successes)
}

// FallbackChain defines a chain of models to try
type FallbackChain struct {
	primary   string
	fallbacks []string
	stats     *FallbackStats
	router    *Router
	costAware bool
}

// NewFallbackChain creates a new fallback chain
func NewFallbackChain(primary string, fallbacks []string, router *Router) *FallbackChain {
	return &FallbackChain{
		primary:   primary,
		fallbacks: fallbacks,
		stats:     NewFallbackStats(),
		router:    router,
		costAware: false,
	}
}

// SetCostAware enables or disables cost-aware fallback
func (f *FallbackChain) SetCostAware(enabled bool) {
	f.costAware = enabled
}

// Execute runs the function with fallback support
func (f *FallbackChain) Execute(ctx context.Context, fn func(model string) error) error {
	f.stats.RecordCall()

	// Build model list
	models := append([]string{f.primary}, f.fallbacks...)

	// Sort by cost if cost-aware
	if f.costAware && f.router != nil {
		models = f.sortByCost(models)
	}

	var lastErr error
	for i, model := range models {
		if i > 0 {
			f.stats.RecordFallback()
		}

		err := fn(model)
		if err == nil {
			f.stats.RecordSuccess()
			return nil
		}

		f.stats.RecordError(model)
		lastErr = err

		// Check if context cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	return errors.Join(ErrAllModelsFailed, lastErr)
}

// sortByCost sorts models by cost (cheapest first)
func (f *FallbackChain) sortByCost(models []string) []string {
	type modelCost struct {
		name string
		cost float64
	}

	costs := make([]modelCost, len(models))
	for i, m := range models {
		costs[i] = modelCost{name: m, cost: f.getModelCost(m)}
	}

	// Sort by cost ascending
	for i := 0; i < len(costs)-1; i++ {
		for j := i + 1; j < len(costs); j++ {
			if costs[i].cost > costs[j].cost {
				costs[i], costs[j] = costs[j], costs[i]
			}
		}
	}

	result := make([]string, len(costs))
	for i, c := range costs {
		result[i] = c.name
	}
	return result
}

// getModelCost returns the cost for a model
func (f *FallbackChain) getModelCost(model string) float64 {
	if f.router == nil {
		return 0
	}
	if p, ok := f.router.GetProfile(model); ok {
		return p.CostPer1K
	}
	return 0
}

// GetStats returns the fallback statistics
func (f *FallbackChain) GetStats() *FallbackStats {
	return f.stats
}

// GetPrimary returns the primary model
func (f *FallbackChain) GetPrimary() string {
	return f.primary
}

// GetFallbacks returns the fallback models
func (f *FallbackChain) GetFallbacks() []string {
	return f.fallbacks
}

// DefaultChains returns default fallback chains
func DefaultChains(router *Router) map[string]*FallbackChain {
	return map[string]*FallbackChain{
		"planning": NewFallbackChain("quality", []string{"balanced", "fast"}, router),
		"code":     NewFallbackChain("balanced", []string{"fast", "quality"}, router),
		"review":   NewFallbackChain("quality", []string{"balanced"}, router),
		"debug":    NewFallbackChain("fast", []string{"balanced"}, router),
		"quick":    NewFallbackChain("fast", []string{}, router),
	}
}
