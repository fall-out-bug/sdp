package orchestrator

import (
	"math"
	"sort"
)

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
