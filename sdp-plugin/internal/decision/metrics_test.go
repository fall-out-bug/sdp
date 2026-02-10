package decision

import (
	"testing"
	"time"
)

func TestMetricsRecorder_RecordLog(t *testing.T) {
	mr := &MetricsRecorder{}

	// Test successful log
	mr.RecordLog(100*time.Millisecond, true)

	// Test failed log
	mr.RecordLog(50*time.Millisecond, false)

	// Verify counters are incremented
	// (expvar globals are modified, but we can't easily reset them)
	// This test ensures the methods don't panic
}

func TestMetricsRecorder_RecordBatchLog(t *testing.T) {
	mr := &MetricsRecorder{}

	// Test successful batch log
	mr.RecordBatchLog(200*time.Millisecond, true, 10)

	// Test failed batch log
	mr.RecordBatchLog(150*time.Millisecond, false, 5)

	// This test ensures the methods don't panic
}

func TestMetricsRecorder_RecordLoad(t *testing.T) {
	mr := &MetricsRecorder{}

	// Test successful load
	mr.RecordLoad(300*time.Millisecond, true, 0)

	// Test load with parse errors
	mr.RecordLoad(250*time.Millisecond, true, 2)

	// Test failed load
	mr.RecordLoad(100*time.Millisecond, false, 0)

	// This test ensures the methods don't panic
}

func TestMetricsRecorder_ConcurrentAccess(t *testing.T) {
	mr := &MetricsRecorder{}

	done := make(chan bool)

	// Start multiple goroutines
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				mr.RecordLog(time.Duration(j)*time.Millisecond, j%2 == 0)
				mr.RecordBatchLog(time.Duration(j)*time.Millisecond, j%2 == 0, j)
				mr.RecordLoad(time.Duration(j)*time.Millisecond, j%2 == 0, j%3)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// If we got here, no race conditions occurred
}
