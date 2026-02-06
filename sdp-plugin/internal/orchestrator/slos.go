package orchestrator

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

// SLO Constants
const (
	// CheckpointSaveLatencyTarget is the p95 target for checkpoint save latency
	CheckpointSaveLatencyTarget = 100 * time.Millisecond
	// CheckpointSaveLatencyAlert is the alert threshold for checkpoint save latency
	CheckpointSaveLatencyAlert = 150 * time.Millisecond

	// WSExecutionTimeTarget is the p95 target for workstream execution time
	WSExecutionTimeTarget = 30 * time.Minute
	// WSExecutionTimeAlert is the alert threshold for workstream execution time
	WSExecutionTimeAlert = 45 * time.Minute

	// GraphBuildTimeTarget is the p95 target for dependency graph build time
	GraphBuildTimeTarget = 5 * time.Second
	// GraphBuildTimeAlert is the alert threshold for dependency graph build time
	GraphBuildTimeAlert = 10 * time.Second

	// RecoverySuccessTarget is the target success rate for checkpoint recovery
	RecoverySuccessTarget = 0.999
	// RecoverySuccessAlert is the alert threshold for checkpoint recovery success rate
	RecoverySuccessAlert = 0.995
)

// Metric tracks measurements for an SLI
type Metric struct {
	mu           sync.RWMutex
	values       []float64 // Duration in seconds or boolean (0/1 for success rate)
	count        int
	successCount int // For success rate metrics
}

// SLOStatus represents the current SLO status
type SLOStatus struct {
	CheckpointSaveLatency       float64 // p95 in seconds
	CheckpointSaveLatencyOK     bool
	WSExecutionTime             float64 // p95 in seconds
	WSExecutionTimeOK           bool
	GraphBuildTime              float64 // p95 in seconds
	GraphBuildTimeOK            bool
	RecoverySuccessRate         float64 // 0-1
	RecoverySuccessRateOK       bool
	OverallSLOCompliance        bool
}

// SLOTracker tracks SLO measurements for orchestrator operations
type SLOTracker struct {
	mu                    sync.RWMutex
	checkpointSaveMetrics *Metric
	wsExecutionMetrics    *Metric
	graphBuildMetrics     *Metric
	recoveryMetrics       *Metric
	logger                *OrchestratorLogger
}

// NewSLOTracker creates a new SLO tracker
func NewSLOTracker(logger *OrchestratorLogger) *SLOTracker {
	return &SLOTracker{
		checkpointSaveMetrics: &Metric{},
		wsExecutionMetrics:    &Metric{},
		graphBuildMetrics:     &Metric{},
		recoveryMetrics:       &Metric{},
		logger:                logger,
	}
}

// RecordCheckpointSave records a checkpoint save latency measurement
func (st *SLOTracker) RecordCheckpointSave(duration time.Duration) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.checkpointSaveMetrics.mu.Lock()
	st.checkpointSaveMetrics.values = append(st.checkpointSaveMetrics.values, duration.Seconds())
	st.checkpointSaveMetrics.mu.Unlock()
}

// RecordWSExecution records a workstream execution time measurement
func (st *SLOTracker) RecordWSExecution(wsID string, duration time.Duration) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.wsExecutionMetrics.mu.Lock()
	st.wsExecutionMetrics.values = append(st.wsExecutionMetrics.values, duration.Seconds())
	st.wsExecutionMetrics.mu.Unlock()

	if st.logger != nil {
		if duration.Seconds() > WSExecutionTimeTarget.Seconds() {
			st.logger.LogWSError(wsID, 0, fmt.Errorf("SLO breach: execution time %v exceeds target %v",
				duration.Round(time.Second), WSExecutionTimeTarget))
		}
	}
}

// RecordGraphBuild records a dependency graph build time measurement
func (st *SLOTracker) RecordGraphBuild(nodeCount int, duration time.Duration) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.graphBuildMetrics.mu.Lock()
	st.graphBuildMetrics.values = append(st.graphBuildMetrics.values, duration.Seconds())
	st.graphBuildMetrics.mu.Unlock()

	if st.logger != nil {
		if duration.Seconds() > GraphBuildTimeTarget.Seconds() {
			st.logger.logger.Warn("SLO breach: graph build time exceeds target",
				"duration_seconds", duration.Seconds(),
				"target_seconds", GraphBuildTimeTarget.Seconds(),
				"node_count", nodeCount)
		}
	}
}

// RecordRecovery records a checkpoint recovery attempt
func (st *SLOTracker) RecordRecovery(success bool) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.recoveryMetrics.mu.Lock()
	st.recoveryMetrics.count++
	if success {
		st.recoveryMetrics.successCount++
	}
	// Store as 1.0 for success, 0.0 for failure
	value := 0.0
	if success {
		value = 1.0
	}
	st.recoveryMetrics.values = append(st.recoveryMetrics.values, value)
	st.recoveryMetrics.mu.Unlock()
}

// GetSLOStatus returns the current SLO status
func (st *SLOTracker) GetSLOStatus() SLOStatus {
	st.mu.RLock()
	defer st.mu.RUnlock()

	status := SLOStatus{}

	// Checkpoint save latency
	status.CheckpointSaveLatency = st.calculatePercentile(st.checkpointSaveMetrics, 95)
	status.CheckpointSaveLatencyOK = status.CheckpointSaveLatency <= CheckpointSaveLatencyTarget.Seconds()

	// Workstream execution time
	status.WSExecutionTime = st.calculatePercentile(st.wsExecutionMetrics, 95)
	status.WSExecutionTimeOK = status.WSExecutionTime <= WSExecutionTimeTarget.Seconds()

	// Graph build time
	status.GraphBuildTime = st.calculatePercentile(st.graphBuildMetrics, 95)
	status.GraphBuildTimeOK = status.GraphBuildTime <= GraphBuildTimeTarget.Seconds()

	// Recovery success rate
	status.RecoverySuccessRate = st.calculateSuccessRate(st.recoveryMetrics)
	status.RecoverySuccessRateOK = status.RecoverySuccessRate >= RecoverySuccessTarget

	// Overall compliance
	status.OverallSLOCompliance = status.CheckpointSaveLatencyOK &&
		status.WSExecutionTimeOK &&
		status.GraphBuildTimeOK &&
		status.RecoverySuccessRateOK

	return status
}

// calculatePercentile calculates the p-th percentile of measurements
func (st *SLOTracker) calculatePercentile(metric *Metric, p float64) float64 {
	metric.mu.RLock()
	defer metric.mu.RUnlock()

	if len(metric.values) == 0 {
		return 0.0
	}

	// Copy values to avoid modifying original
	values := make([]float64, len(metric.values))
	copy(values, metric.values)

	// Sort for percentile calculation
	sort.Float64s(values)

	// Calculate percentile index
	index := int(math.Ceil((p/100.0)*float64(len(values)))) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(values) {
		index = len(values) - 1
	}

	return values[index]
}

// calculateSuccessRate calculates the success rate from binary measurements
func (st *SLOTracker) calculateSuccessRate(metric *Metric) float64 {
	metric.mu.RLock()
	defer metric.mu.RUnlock()

	if metric.count == 0 {
		return 1.0 // No failures if no attempts
	}

	return float64(metric.successCount) / float64(metric.count)
}

// GetCheckpointSaveMetrics returns the count of checkpoint save measurements
func (st *SLOTracker) GetCheckpointSaveMetrics() int {
	st.checkpointSaveMetrics.mu.RLock()
	defer st.checkpointSaveMetrics.mu.RUnlock()
	return len(st.checkpointSaveMetrics.values)
}

// GetWSExecutionMetrics returns the count of workstream execution measurements
func (st *SLOTracker) GetWSExecutionMetrics() int {
	st.wsExecutionMetrics.mu.RLock()
	defer st.wsExecutionMetrics.mu.RUnlock()
	return len(st.wsExecutionMetrics.values)
}

// GetGraphBuildMetrics returns the count of graph build measurements
func (st *SLOTracker) GetGraphBuildMetrics() int {
	st.graphBuildMetrics.mu.RLock()
	defer st.graphBuildMetrics.mu.RUnlock()
	return len(st.graphBuildMetrics.values)
}

// GetRecoveryMetrics returns the count of recovery attempts and successes
func (st *SLOTracker) GetRecoveryMetrics() (int, int) {
	st.recoveryMetrics.mu.RLock()
	defer st.recoveryMetrics.mu.RUnlock()
	return st.recoveryMetrics.count, st.recoveryMetrics.successCount
}
