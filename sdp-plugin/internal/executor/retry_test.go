package executor

import (
	"testing"
)

// TestRetryDelayFromConfigCached verifies retry delay is cached (no repeated config I/O per retry).
func TestRetryDelayFromConfigCached(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir: "testdata/backlog",
		DryRun:     false,
		RetryCount: 1,
	}, newTestRunner())

	// Call multiple times - should return same value (cached)
	d1 := exec.retryDelayFromConfigCached()
	d2 := exec.retryDelayFromConfigCached()
	if d1 != d2 {
		t.Errorf("retryDelayFromConfigCached should return cached value: got %v, %v", d1, d2)
	}
	if d1 <= 0 {
		t.Errorf("retry delay should be positive, got %v", d1)
	}
}
