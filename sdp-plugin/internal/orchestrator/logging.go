package orchestrator

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	correlationIDKey contextKey = "correlation_id"
)

// OrchestratorLogger wraps slog.Logger with correlation ID support
type OrchestratorLogger struct {
	logger        *slog.Logger
	correlationID string
	featureID     string
}

// NewOrchestratorLogger creates a new logger with correlation ID
func NewOrchestratorLogger(featureID string) *OrchestratorLogger {
	correlationID := fmt.Sprintf("%s-%d", featureID, time.Now().UnixNano())

	logger := slog.Default()
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}

	return &OrchestratorLogger{
		logger:        logger,
		correlationID: correlationID,
		featureID:     featureID,
	}
}

// NewOrchestratorLoggerWithHandler creates a logger with custom handler (for testing)
func NewOrchestratorLoggerWithHandler(featureID string, handler slog.Handler) *OrchestratorLogger {
	correlationID := fmt.Sprintf("%s-%d", featureID, time.Now().UnixNano())

	return &OrchestratorLogger{
		logger:        slog.New(handler),
		correlationID: correlationID,
		featureID:     featureID,
	}
}

// WithContext adds correlation ID to context
func (ol *OrchestratorLogger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, correlationIDKey, ol.correlationID)
}

// LogStart logs orchestration start event
func (ol *OrchestratorLogger) LogStart(featureID string) {
	ol.logger.Info("orchestration_start",
		"feature_id", featureID,
		"correlation_id", ol.correlationID,
		"timestamp", time.Now().UTC().Format(time.RFC3339),
	)
}

// LogWSStart logs workstream start event
func (ol *OrchestratorLogger) LogWSStart(wsID string) {
	ol.logger.Info("workstream_start",
		"ws_id", wsID,
		"feature_id", ol.featureID,
		"correlation_id", ol.correlationID,
		"timestamp", time.Now().UTC().Format(time.RFC3339),
	)
}

// LogWSComplete logs workstream completion event
func (ol *OrchestratorLogger) LogWSComplete(wsID string, duration time.Duration, coverage float64) {
	ol.logger.Info("workstream_complete",
		"ws_id", wsID,
		"feature_id", ol.featureID,
		"correlation_id", ol.correlationID,
		"duration_seconds", duration.Seconds(),
		"coverage_percent", coverage,
		"timestamp", time.Now().UTC().Format(time.RFC3339),
	)
}

// LogWSCheckpoint logs checkpoint save event
func (ol *OrchestratorLogger) LogWSCheckpoint(wsID, checkpointFile string) {
	ol.logger.Info("checkpoint_saved",
		"ws_id", wsID,
		"checkpoint_file", checkpointFile,
		"feature_id", ol.featureID,
		"correlation_id", ol.correlationID,
		"timestamp", time.Now().UTC().Format(time.RFC3339),
	)
}

// LogWSError logs workstream error event
func (ol *OrchestratorLogger) LogWSError(wsID string, attempt int, err error) {
	ol.logger.Error("workstream_error",
		"ws_id", wsID,
		"feature_id", ol.featureID,
		"correlation_id", ol.correlationID,
		"attempt", attempt,
		"error", err.Error(),
		"timestamp", time.Now().UTC().Format(time.RFC3339),
	)
}

// LogOrchestrationComplete logs orchestration completion event
func (ol *OrchestratorLogger) LogOrchestrationComplete(featureID string, duration time.Duration, totalWS int, success bool) {
	ol.logger.Info("orchestration_complete",
		"feature_id", featureID,
		"correlation_id", ol.correlationID,
		"duration_seconds", duration.Seconds(),
		"total_workstreams", totalWS,
		"success", success,
		"timestamp", time.Now().UTC().Format(time.RFC3339),
	)
}

// LogRetry logs retry attempt
func (ol *OrchestratorLogger) LogRetry(wsID string, attempt int, maxAttempts int) {
	ol.logger.Warn("retry_attempt",
		"ws_id", wsID,
		"feature_id", ol.featureID,
		"correlation_id", ol.correlationID,
		"attempt", attempt,
		"max_attempts", maxAttempts,
		"timestamp", time.Now().UTC().Format(time.RFC3339),
	)
}

// LogDependencyGraph logs dependency graph building
func (ol *OrchestratorLogger) LogDependencyGraph(nodeCount, edgeCount int) {
	ol.logger.Info("dependency_graph_built",
		"feature_id", ol.featureID,
		"correlation_id", ol.correlationID,
		"nodes", nodeCount,
		"edges", edgeCount,
		"timestamp", time.Now().UTC().Format(time.RFC3339),
	)
}

// GetCorrelationID returns the correlation ID
func (ol *OrchestratorLogger) GetCorrelationID() string {
	return ol.correlationID
}

// GetFeatureID returns the feature ID
func (ol *OrchestratorLogger) GetFeatureID() string {
	return ol.featureID
}
