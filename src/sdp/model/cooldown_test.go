package model

import (
	"context"
	"testing"
	"time"
)

func TestCooldownManager_WaitForCooldown(t *testing.T) {
	config := CooldownConfig{
		DefaultCooldown: 50 * time.Millisecond,
		GlobalRate:      10 * time.Millisecond,
	}
	cm := NewCooldownManager(config)

	// First call should succeed immediately
	start := time.Now()
	err := cm.WaitForCooldown(context.Background(), "test-model")
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if elapsed > 20*time.Millisecond {
		t.Errorf("First call took too long: %v", elapsed)
	}

	cm.RecordCall("test-model")

	// Second call should wait for cooldown
	start = time.Now()
	err = cm.WaitForCooldown(context.Background(), "test-model")
	elapsed = time.Since(start)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if elapsed < 40*time.Millisecond {
		t.Errorf("Second call should have waited, but took only %v", elapsed)
	}
}

func TestCooldownManager_ContextCancellation(t *testing.T) {
	config := CooldownConfig{
		DefaultCooldown: 1 * time.Second,
		GlobalRate:      10 * time.Millisecond,
	}
	cm := NewCooldownManager(config)
	cm.RecordCall("test-model")

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := cm.WaitForCooldown(ctx, "test-model")
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestCooldownManager_SetCooldown(t *testing.T) {
	cm := NewCooldownManager(DefaultCooldownConfig)
	cm.SetCooldown("custom-model", 200*time.Millisecond)

	if cd := cm.GetCooldown("custom-model"); cd != 200*time.Millisecond {
		t.Errorf("Expected 200ms cooldown, got %v", cd)
	}
}

func TestCooldownManager_Backoff(t *testing.T) {
	config := CooldownConfig{
		DefaultCooldown: 100 * time.Millisecond,
		MaxBackoff:      1 * time.Second,
	}
	cm := NewCooldownManager(config)

	// Initial cooldown
	initial := cm.GetCooldown("test-model")
	if initial != 100*time.Millisecond {
		t.Errorf("Expected 100ms initial, got %v", initial)
	}

	// After backoff
	cm.Backoff("test-model")
	afterBackoff := cm.GetCooldown("test-model")
	if afterBackoff != 200*time.Millisecond {
		t.Errorf("Expected 200ms after backoff, got %v", afterBackoff)
	}

	// Multiple backoffs
	cm.Backoff("test-model")
	cm.Backoff("test-model")
	cm.Backoff("test-model")
	cm.Backoff("test-model")
	// Should cap at MaxBackoff
	if cd := cm.GetCooldown("test-model"); cd > config.MaxBackoff {
		t.Errorf("Cooldown %v exceeds max %v", cd, config.MaxBackoff)
	}
}

func TestCooldownManager_Reset(t *testing.T) {
	cm := NewCooldownManager(DefaultCooldownConfig)
	cm.SetCooldown("test-model", 500*time.Millisecond)
	cm.Backoff("test-model")

	cm.Reset("test-model")

	if cd := cm.GetCooldown("test-model"); cd != DefaultCooldownConfig.DefaultCooldown {
		t.Errorf("Reset should restore default, got %v", cd)
	}
}

func TestCooldownManager_ResetAll(t *testing.T) {
	cm := NewCooldownManager(DefaultCooldownConfig)
	cm.SetCooldown("model1", 500*time.Millisecond)
	cm.SetCooldown("model2", 600*time.Millisecond)

	cm.ResetAll()

	if cd := cm.GetCooldown("model1"); cd != DefaultCooldownConfig.DefaultCooldown {
		t.Errorf("ResetAll should clear model1, got %v", cd)
	}
}

func TestCooldownManager_TimeUntilNext(t *testing.T) {
	config := CooldownConfig{
		DefaultCooldown: 100 * time.Millisecond,
		GlobalRate:      10 * time.Millisecond,
	}
	cm := NewCooldownManager(config)

	// No calls yet - should be 0
	if wait := cm.TimeUntilNext("test-model"); wait != 0 {
		t.Errorf("Expected 0 wait initially, got %v", wait)
	}

	cm.RecordCall("test-model")

	// After call - should have wait time
	if wait := cm.TimeUntilNext("test-model"); wait < 90*time.Millisecond {
		t.Errorf("Expected ~100ms wait, got %v", wait)
	}
}

func TestCooldownManager_Concurrent(t *testing.T) {
	config := CooldownConfig{
		DefaultCooldown: 10 * time.Millisecond,
		GlobalRate:      5 * time.Millisecond,
	}
	cm := NewCooldownManager(config)

	done := make(chan bool)

	// Launch concurrent operations
	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				cm.WaitForCooldown(context.Background(), "test-model")
				cm.RecordCall("test-model")
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}
}
