package runtime

import (
	"context"
	"testing"
	"time"
)

func TestDegradedMode_EnterExit(t *testing.T) {
	dm := NewDegradedMode()

	if dm.IsActive() {
		t.Error("Should not be active initially")
	}

	dm.Enter("test reason")

	if !dm.IsActive() {
		t.Error("Should be active after Enter")
	}
	if dm.Reason() != "test reason" {
		t.Errorf("Expected 'test reason', got %s", dm.Reason())
	}
	if dm.Since().IsZero() {
		t.Error("Since should be set")
	}

	dm.Exit()

	if dm.IsActive() {
		t.Error("Should not be active after Exit")
	}
	if dm.Duration() != 0 {
		t.Errorf("Duration should be 0 when inactive, got %v", dm.Duration())
	}
}

func TestDegradedMode_Duration(t *testing.T) {
	dm := NewDegradedMode()

	dm.Enter("test")
	time.Sleep(10 * time.Millisecond)

	dur := dm.Duration()
	if dur < 10*time.Millisecond {
		t.Errorf("Duration should be at least 10ms, got %v", dur)
	}
}

func TestDegradedMode_Cache(t *testing.T) {
	dm := NewDegradedMode()

	// Set and get
	dm.SetCache("key1", "value1")

	val, ok := dm.GetCache("key1")
	if !ok {
		t.Error("Expected to find key1")
	}
	if val != "value1" {
		t.Errorf("Expected 'value1', got %v", val)
	}

	// Non-existent key
	_, ok = dm.GetCache("nonexistent")
	if ok {
		t.Error("Should not find nonexistent key")
	}

	// Keys
	keys := dm.CachedKeys()
	if len(keys) != 1 || keys[0] != "key1" {
		t.Errorf("Expected [key1], got %v", keys)
	}

	// Clear
	dm.ClearCache()
	_, ok = dm.GetCache("key1")
	if ok {
		t.Error("Should not find key1 after clear")
	}
}

func TestDegradedMode_EnterIdempotent(t *testing.T) {
	dm := NewDegradedMode()

	dm.Enter("reason1")
	since1 := dm.Since()

	dm.Enter("reason2") // Should not update

	if dm.Reason() != "reason1" {
		t.Errorf("Reason should remain 'reason1', got %s", dm.Reason())
	}
	if dm.Since() != since1 {
		t.Error("Since should not change on repeated Enter")
	}
}

func TestServiceChecker(t *testing.T) {
	checker := NewServiceChecker("http://test", 1*time.Second, 100*time.Millisecond)

	// First check
	available := checker.IsAvailable(context.Background())
	if !available {
		t.Error("Should be available")
	}

	// Check period should prevent re-check (returns cached value)
	// The cached value is true from first check
	available = checker.IsAvailable(context.Background())
	if !available {
		t.Error("Should use cached status (true) within check period")
	}
}

func TestCheckServiceHealth(t *testing.T) {
	dm := NewDegradedMode()
	checker := NewServiceChecker("http://test", 1*time.Second, 100*time.Millisecond)

	err := CheckServiceHealth(context.Background(), checker, dm, "test-service")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestWithDegradedMode_Success(t *testing.T) {
	dm := NewDegradedMode()
	fallbackCalled := false

	err := WithDegradedMode(context.Background(), dm,
		func() error { return nil },
		func() error { fallbackCalled = true; return nil },
	)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if fallbackCalled {
		t.Error("Fallback should not be called on success")
	}
}

func TestWithDegradedMode_Active(t *testing.T) {
	dm := NewDegradedMode()
	dm.Enter("test")
	fallbackCalled := false

	err := WithDegradedMode(context.Background(), dm,
		func() error { panic("should not be called") },
		func() error { fallbackCalled = true; return nil },
	)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !fallbackCalled {
		t.Error("Fallback should be called when degraded")
	}
}

func TestGlobalDegradedMode(t *testing.T) {
	// Reset global state
	GlobalDegradedMode.Exit()

	if GlobalDegradedMode.IsActive() {
		t.Error("Global should not be active initially")
	}

	GlobalDegradedMode.Enter("test")
	if !GlobalDegradedMode.IsActive() {
		t.Error("Global should be active")
	}

	GlobalDegradedMode.Exit()
}
