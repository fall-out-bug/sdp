package executor

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// testRunner is a mock WorkstreamRunner for tests.
// Replicates the original executeWorkstreamMock behavior:
// 00-054-02 fails on first attempt, succeeds on retry.
// All other workstreams always succeed.
type testRunner struct {
	mu    sync.Mutex
	calls map[string]int
}

func newTestRunner() *testRunner {
	return &testRunner{calls: make(map[string]int)}
}

func (r *testRunner) Run(_ context.Context, wsID string) error {
	r.mu.Lock()
	r.calls[wsID]++
	attempt := r.calls[wsID]
	r.mu.Unlock()

	if wsID == "00-054-02" && attempt == 1 {
		return fmt.Errorf("mock execution failure for %s", wsID)
	}
	return nil
}

// blockingRunner is a WorkstreamRunner that blocks until ctx is cancelled.
// Used to test mid-execution context cancellation.
type blockingRunner struct{}

func newBlockingRunner() *blockingRunner { return &blockingRunner{} }

func (r *blockingRunner) Run(ctx context.Context, _ string) error {
	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// Keep running, check ctx again
		}
	}
}
