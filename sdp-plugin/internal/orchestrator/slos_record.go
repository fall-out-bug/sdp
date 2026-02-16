package orchestrator

import (
	"fmt"
	"time"
)

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
