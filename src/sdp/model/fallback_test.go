package model

import (
	"context"
	"errors"
	"testing"
)

func TestFallbackChain_Execute_PrimarySuccess(t *testing.T) {
	r := NewRouter()
	fc := NewFallbackChain("balanced", []string{"fast", "quality"}, r)

	calls := []string{}
	err := fc.Execute(context.Background(), func(model string) error {
		calls = append(calls, model)
		return nil
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(calls) != 1 || calls[0] != "balanced" {
		t.Errorf("Expected only primary call, got %v", calls)
	}
}

func TestFallbackChain_Execute_Fallback(t *testing.T) {
	r := NewRouter()
	fc := NewFallbackChain("balanced", []string{"fast", "quality"}, r)

	calls := []string{}
	callCount := 0

	err := fc.Execute(context.Background(), func(model string) error {
		calls = append(calls, model)
		callCount++
		if callCount < 3 {
			return errors.New("fail")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(calls) != 3 {
		t.Errorf("Expected 3 calls (primary + 2 fallbacks), got %v", calls)
	}
}

func TestFallbackChain_Execute_AllFail(t *testing.T) {
	r := NewRouter()
	fc := NewFallbackChain("balanced", []string{"fast"}, r)

	err := fc.Execute(context.Background(), func(model string) error {
		return errors.New("always fail")
	})

	if err == nil {
		t.Error("Expected error when all models fail")
	}
	if !errors.Is(err, ErrAllModelsFailed) {
		t.Errorf("Expected ErrAllModelsFailed, got %v", err)
	}
}

func TestFallbackChain_Execute_ContextCancel(t *testing.T) {
	r := NewRouter()
	fc := NewFallbackChain("balanced", []string{"fast"}, r)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := fc.Execute(ctx, func(model string) error {
		return errors.New("fail")
	})

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestFallbackChain_CostAware(t *testing.T) {
	r := NewRouter()
	fc := NewFallbackChain("quality", []string{"fast", "balanced"}, r)
	fc.SetCostAware(true)

	calls := []string{}
	err := fc.Execute(context.Background(), func(model string) error {
		calls = append(calls, model)
		return errors.New("fail")
	})

	if err == nil {
		t.Error("Expected error")
	}

	// In cost-aware mode, should try cheapest first
	// fast (0.003) < balanced (0.003) < quality (0.015)
	// So order should be: fast, balanced, quality (or fast, quality, balanced)
	if len(calls) < 2 {
		t.Errorf("Expected at least 2 calls, got %d", len(calls))
	}
	// First call should be the cheapest (fast)
	if calls[0] != "fast" {
		t.Errorf("Expected fast first in cost-aware mode, got %v", calls)
	}
}

func TestFallbackChain_Stats(t *testing.T) {
	r := NewRouter()
	fc := NewFallbackChain("balanced", []string{"fast"}, r)

	// Success on primary
	fc.Execute(context.Background(), func(model string) error { return nil })

	total, fallbacks, successes := fc.GetStats().GetStats()
	if total != 1 || successes != 1 || fallbacks != 0 {
		t.Errorf("Stats mismatch: total=%d, successes=%d, fallbacks=%d", total, successes, fallbacks)
	}

	// Fallback case
	fc.Execute(context.Background(), func(model string) error {
		if model == "balanced" {
			return errors.New("fail")
		}
		return nil
	})

	total, fallbacks, successes = fc.GetStats().GetStats()
	if total != 2 || successes != 2 || fallbacks != 1 {
		t.Errorf("Stats mismatch after fallback: total=%d, successes=%d, fallbacks=%d", total, successes, fallbacks)
	}
}

func TestFallbackChain_Getters(t *testing.T) {
	r := NewRouter()
	fc := NewFallbackChain("balanced", []string{"fast", "quality"}, r)

	if fc.GetPrimary() != "balanced" {
		t.Errorf("Expected primary 'balanced', got %s", fc.GetPrimary())
	}

	fallbacks := fc.GetFallbacks()
	if len(fallbacks) != 2 {
		t.Errorf("Expected 2 fallbacks, got %d", len(fallbacks))
	}
}

func TestDefaultChains(t *testing.T) {
	r := NewRouter()
	chains := DefaultChains(r)

	if len(chains) == 0 {
		t.Error("DefaultChains returned empty")
	}

	if _, ok := chains["planning"]; !ok {
		t.Error("Missing planning chain")
	}
	if _, ok := chains["code"]; !ok {
		t.Error("Missing code chain")
	}
}

func TestFallbackStats(t *testing.T) {
	stats := NewFallbackStats()

	stats.RecordCall()
	stats.RecordCall()
	stats.RecordFallback()
	stats.RecordSuccess()
	stats.RecordError("model1")
	stats.RecordError("model1")

	total, fallbacks, successes := stats.GetStats()
	if total != 2 {
		t.Errorf("Expected 2 total, got %d", total)
	}
	if fallbacks != 1 {
		t.Errorf("Expected 1 fallback, got %d", fallbacks)
	}
	if successes != 1 {
		t.Errorf("Expected 1 success, got %d", successes)
	}
	if stats.ModelErrors["model1"] != 2 {
		t.Errorf("Expected 2 errors for model1, got %d", stats.ModelErrors["model1"])
	}
}
