package orchestrator

import (
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
	CheckpointSaveLatency   float64 // p95 in seconds
	CheckpointSaveLatencyOK bool
	WSExecutionTime         float64 // p95 in seconds
	WSExecutionTimeOK       bool
	GraphBuildTime          float64 // p95 in seconds
	GraphBuildTimeOK        bool
	RecoverySuccessRate     float64 // 0-1
	RecoverySuccessRateOK   bool
	OverallSLOCompliance    bool
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
