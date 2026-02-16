package orchestrator

import (
	"bytes"
	"log/slog"
	"math"
	"strings"
	"testing"
	"time"
)

// TestNewSLOTracker tests SLO tracker creation
func TestNewSLOTracker(t *testing.T) {
	tracker := NewSLOTracker(nil)

	if tracker == nil {
		t.Fatal("NewSLOTracker() returned nil")
	}

	if tracker.checkpointSaveMetrics == nil {
		t.Error("checkpointSaveMetrics not initialized")
	}
	if tracker.wsExecutionMetrics == nil {
		t.Error("wsExecutionMetrics not initialized")
	}
	if tracker.graphBuildMetrics == nil {
		t.Error("graphBuildMetrics not initialized")
	}
	if tracker.recoveryMetrics == nil {
		t.Error("recoveryMetrics not initialized")
	}
}

// TestRecordCheckpointSave tests checkpoint save latency recording
func TestRecordCheckpointSave(t *testing.T) {
	tracker := NewSLOTracker(nil)

	// Record some measurements
	tracker.RecordCheckpointSave(50 * time.Millisecond)
	tracker.RecordCheckpointSave(100 * time.Millisecond)
	tracker.RecordCheckpointSave(150 * time.Millisecond)

	count := tracker.GetCheckpointSaveMetrics()
	if count != 3 {
		t.Errorf("Expected 3 measurements, got %d", count)
	}
}

// TestRecordWSExecution tests workstream execution time recording
func TestRecordWSExecution(t *testing.T) {
	logger := &OrchestratorLogger{}
	tracker := NewSLOTracker(logger)

	// Record normal execution
	tracker.RecordWSExecution("00-001-01", 20*time.Minute)

	count := tracker.GetWSExecutionMetrics()
	if count != 1 {
		t.Errorf("Expected 1 measurement, got %d", count)
	}
}

// TestRecordWSExecutionSLOBreach tests SLO breach detection
func TestRecordWSExecutionSLOBreach(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	ol := &OrchestratorLogger{logger: logger}
	tracker := NewSLOTracker(ol)

	// Record execution that exceeds target
	tracker.RecordWSExecution("00-001-01", 35*time.Minute)

	// Should log SLO breach
	logOutput := buf.String()
	if !strings.Contains(logOutput, "SLO breach") {
		t.Error("Expected SLO breach log message")
	}
}

// TestRecordGraphBuild tests graph build time recording
func TestRecordGraphBuild(t *testing.T) {
	tracker := NewSLOTracker(nil)

	tracker.RecordGraphBuild(100, 2*time.Second)

	count := tracker.GetGraphBuildMetrics()
	if count != 1 {
		t.Errorf("Expected 1 measurement, got %d", count)
	}
}

// TestRecordRecovery tests recovery success/failure recording
func TestRecordRecovery(t *testing.T) {
	tracker := NewSLOTracker(nil)

	tracker.RecordRecovery(true)
	tracker.RecordRecovery(true)
	tracker.RecordRecovery(false)

	count, successCount := tracker.GetRecoveryMetrics()
	if count != 3 {
		t.Errorf("Expected 3 attempts, got %d", count)
	}
	if successCount != 2 {
		t.Errorf("Expected 2 successes, got %d", successCount)
	}
}

// TestCalculatePercentile tests percentile calculation
func TestCalculatePercentile(t *testing.T) {
	tracker := NewSLOTracker(nil)

	// Record measurements: 50, 100, 150, 200, 250 ms
	for _, d := range []time.Duration{
		50 * time.Millisecond,
		100 * time.Millisecond,
		150 * time.Millisecond,
		200 * time.Millisecond,
		250 * time.Millisecond,
	} {
		tracker.RecordCheckpointSave(d)
	}

	status := tracker.GetSLOStatus()

	// p95 should be 250ms (highest value for 5 measurements)
	expected95th := 0.250 // seconds
	if math.Abs(status.CheckpointSaveLatency-expected95th) > 0.001 {
		t.Errorf("p95 = %v, want %v", status.CheckpointSaveLatency, expected95th)
	}
}

// TestCalculatePercentileEmpty tests percentile calculation with no data
func TestCalculatePercentileEmpty(t *testing.T) {
	tracker := NewSLOTracker(nil)

	status := tracker.GetSLOStatus()

	if status.CheckpointSaveLatency != 0.0 {
		t.Errorf("Expected 0.0 for empty metric, got %v", status.CheckpointSaveLatency)
	}
}

// TestCalculateSuccessRate tests success rate calculation
func TestCalculateSuccessRate(t *testing.T) {
	tracker := NewSLOTracker(nil)

	// 9 successes, 1 failure = 90% success rate
	for i := 0; i < 9; i++ {
		tracker.RecordRecovery(true)
	}
	tracker.RecordRecovery(false)

	status := tracker.GetSLOStatus()

	expectedRate := 0.9
	if math.Abs(status.RecoverySuccessRate-expectedRate) > 0.001 {
		t.Errorf("Success rate = %v, want %v", status.RecoverySuccessRate, expectedRate)
	}
}

// TestCalculateSuccessRateNoAttempts tests success rate with no attempts
func TestCalculateSuccessRateNoAttempts(t *testing.T) {
	tracker := NewSLOTracker(nil)

	status := tracker.GetSLOStatus()

	if status.RecoverySuccessRate != 1.0 {
		t.Errorf("Expected 1.0 (100%%) for no attempts, got %v", status.RecoverySuccessRate)
	}
}

// TestSLOStatusCheckpointSave tests checkpoint save SLO status
func TestSLOStatusCheckpointSave(t *testing.T) {
	tracker := NewSLOTracker(nil)

	// All measurements under target (100ms)
	for _, d := range []time.Duration{
		50 * time.Millisecond,
		75 * time.Millisecond,
		90 * time.Millisecond,
	} {
		tracker.RecordCheckpointSave(d)
	}

	status := tracker.GetSLOStatus()

	if !status.CheckpointSaveLatencyOK {
		t.Error("Expected checkpoint save SLO to be OK")
	}
}

// TestSLOStatusCheckpointSaveBreached tests checkpoint save SLO breach
func TestSLOStatusCheckpointSaveBreached(t *testing.T) {
	tracker := NewSLOTracker(nil)

	// Some measurements exceed target
	tracker.RecordCheckpointSave(50 * time.Millisecond)
	tracker.RecordCheckpointSave(75 * time.Millisecond)
	tracker.RecordCheckpointSave(150 * time.Millisecond) // Exceeds 100ms target

	status := tracker.GetSLOStatus()

	if status.CheckpointSaveLatencyOK {
		t.Error("Expected checkpoint save SLO to be breached")
	}
}

// TestSLOStatusWSExecution tests workstream execution SLO status
func TestSLOStatusWSExecution(t *testing.T) {
	tracker := NewSLOTracker(nil)

	// All measurements under target (30min)
	for _, d := range []time.Duration{
		20 * time.Minute,
		25 * time.Minute,
		28 * time.Minute,
	} {
		tracker.RecordWSExecution("00-001-01", d)
	}

	status := tracker.GetSLOStatus()

	if !status.WSExecutionTimeOK {
		t.Error("Expected WS execution SLO to be OK")
	}
}

// TestSLOStatusWSExecutionBreached tests WS execution SLO breach
func TestSLOStatusWSExecutionBreached(t *testing.T) {
	tracker := NewSLOTracker(nil)

	// Some measurements exceed target
	tracker.RecordWSExecution("00-001-01", 20*time.Minute)
	tracker.RecordWSExecution("00-001-02", 25*time.Minute)
	tracker.RecordWSExecution("00-001-03", 35*time.Minute) // Exceeds 30min target

	status := tracker.GetSLOStatus()

	if status.WSExecutionTimeOK {
		t.Error("Expected WS execution SLO to be breached")
	}
}

// TestSLOStatusGraphBuild tests graph build SLO status
func TestSLOStatusGraphBuild(t *testing.T) {
	tracker := NewSLOTracker(nil)

	// All measurements under target (5s)
	for _, d := range []time.Duration{
		2 * time.Second,
		3 * time.Second,
		4 * time.Second,
	} {
		tracker.RecordGraphBuild(100, d)
	}

	status := tracker.GetSLOStatus()

	if !status.GraphBuildTimeOK {
		t.Error("Expected graph build SLO to be OK")
	}
}

// TestSLOStatusGraphBuildBreached tests graph build SLO breach
func TestSLOStatusGraphBuildBreached(t *testing.T) {
	tracker := NewSLOTracker(nil)

	// Some measurements exceed target
	tracker.RecordGraphBuild(100, 2*time.Second)
	tracker.RecordGraphBuild(100, 3*time.Second)
	tracker.RecordGraphBuild(100, 6*time.Second) // Exceeds 5s target

	status := tracker.GetSLOStatus()

	if status.GraphBuildTimeOK {
		t.Error("Expected graph build SLO to be breached")
	}
}

// TestSLOStatusRecoverySuccess tests recovery success SLO status
func TestSLOStatusRecoverySuccess(t *testing.T) {
	tracker := NewSLOTracker(nil)

	// 999 successes, 1 failure = 99.9% success rate (meets target)
	for i := 0; i < 999; i++ {
		tracker.RecordRecovery(true)
	}
	tracker.RecordRecovery(false)

	status := tracker.GetSLOStatus()

	if !status.RecoverySuccessRateOK {
		t.Error("Expected recovery success SLO to be OK")
	}
}

// TestSLOStatusRecoverySuccessBreached tests recovery success SLO breach
func TestSLOStatusRecoverySuccessBreached(t *testing.T) {
	tracker := NewSLOTracker(nil)

	// 95 successes, 5 failures = 95% success rate (below 99.9% target)
	for i := 0; i < 95; i++ {
		tracker.RecordRecovery(true)
	}
	for i := 0; i < 5; i++ {
		tracker.RecordRecovery(false)
	}

	status := tracker.GetSLOStatus()

	if status.RecoverySuccessRateOK {
		t.Error("Expected recovery success SLO to be breached")
	}
}

// TestOverallSLOCompliance tests overall SLO compliance
func TestOverallSLOCompliance(t *testing.T) {
	tracker := NewSLOTracker(nil)

	// Add measurements that meet all SLOs
	tracker.RecordCheckpointSave(50 * time.Millisecond)
	tracker.RecordWSExecution("00-001-01", 20*time.Minute)
	tracker.RecordGraphBuild(100, 2*time.Second)
	tracker.RecordRecovery(true)

	status := tracker.GetSLOStatus()

	if !status.OverallSLOCompliance {
		t.Error("Expected overall SLO compliance to be true")
	}
}

// TestOverallSLOComplianceBreached tests overall SLO compliance breach
func TestOverallSLOComplianceBreached(t *testing.T) {
	tracker := NewSLOTracker(nil)

	// Add measurements that breach one SLO
	tracker.RecordCheckpointSave(150 * time.Millisecond) // Breaches checkpoint save SLO

	status := tracker.GetSLOStatus()

	if status.OverallSLOCompliance {
		t.Error("Expected overall SLO compliance to be false")
	}
}

// TestConcurrentSLOAccess tests concurrent SLO tracking
func TestConcurrentSLOAccess(t *testing.T) {
	tracker := NewSLOTracker(nil)

	done := make(chan bool)

	// Start multiple goroutines recording metrics
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				tracker.RecordCheckpointSave(50 * time.Millisecond)
				tracker.RecordWSExecution("00-001-01", 20*time.Minute)
				tracker.RecordGraphBuild(100, 2*time.Second)
				tracker.RecordRecovery(true)
				tracker.GetSLOStatus()
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify counts
	expectedCount := 1000 // 10 goroutines * 100 iterations
	if tracker.GetCheckpointSaveMetrics() != expectedCount {
		t.Errorf("Expected %d checkpoint measurements, got %d", expectedCount, tracker.GetCheckpointSaveMetrics())
	}
	if tracker.GetWSExecutionMetrics() != expectedCount {
		t.Errorf("Expected %d WS execution measurements, got %d", expectedCount, tracker.GetWSExecutionMetrics())
	}
	if tracker.GetGraphBuildMetrics() != expectedCount {
		t.Errorf("Expected %d graph build measurements, got %d", expectedCount, tracker.GetGraphBuildMetrics())
	}
}
