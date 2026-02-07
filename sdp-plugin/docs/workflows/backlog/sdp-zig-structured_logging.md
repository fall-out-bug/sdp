# sdp-zig: Add Structured Logging to Orchestrator

> **Issue ID:** sdp-zig
> **Type:** Bug (Critical - P0)
> **Priority:** 0
> **Source:** SRE Review finding

## Goal

Add structured logging to orchestrator to enable production debugging and monitoring.

## Problem

SRE review identified that orchestrator lacks structured logging, making it impossible to:
- Debug production issues
- Trace execution flow
- Monitor orchestrator performance
- Correlate events across workstreams

## Acceptance Criteria

### AC1: slog Integration
- [ ] Import `log/slog` package in orchestrator components
- [ ] Configure slog with JSON handler for production
- [ ] Add logger to FeatureCoordinator and Orchestrator structs

### AC2: Correlation IDs
- [ ] Generate unique correlation ID per feature execution
- [ ] Include correlation ID in all log messages
- [ ] Pass correlation ID to spawned agents

### AC3: Key Events Logged
- [ ] Orchestration started: feature_id, correlation_id, timestamp
- [ ] Workstream started: ws_id, correlation_id
- [ ] Workstream completed: ws_id, duration, coverage, commit
- [ ] Workstream failed: ws_id, error, attempt
- [ ] Checkpoint saved: checkpoint_file, ws_count
- [ ] Orchestration completed: feature_id, duration, total_ws, success/failure

### AC4: Structured Fields
- [ ] All logs include: level, msg, time, correlation_id
- [ ] Context-specific fields: feature_id, ws_id, error, duration
- [ ] No unstructured fmt.Printf/Terminal output (except user-facing)

### AC5: Log Levels
- [ ] INFO: Normal operations (start, complete, checkpoint)
- [ ] WARN: Retry attempts, degraded performance
- [ ] ERROR: Workstream failures, blocking errors
- [ ] DEBUG: Detailed diagnostics (disabled in production)

### AC6: Testing
- [ ] Unit tests for log output capture
- [ ] Verify correlation ID propagation
- [ ] Test log level filtering

## Scope Files

**MODIFY:**
- `internal/orchestrator/feature_coordinator.go` (add logger)
- `internal/orchestrator/orchestrator.go` (add logger)

**NEW:**
- `internal/orchestrator/logging.go` (~100 LOC) - logging utilities
- `internal/orchestrator/logging_test.go` (~150 LOC)

## Implementation Steps

1. Create logging utilities:
   ```go
   type OrchestratorLogger struct {
       logger *slog.Logger
       correlationID string
   }

   func NewOrchestratorLogger(featureID string) *OrchestratorLogger
   func (ol *OrchestratorLogger) WithCorrelationID() *OrchestratorLogger
   func (ol *OrchestratorLogger) LogStart(featureID string)
   func (ol *OrchestratorLogger) LogWSStart(wsID string)
   func (ol *OrchestratorLogger) LogWSComplete(wsID string, duration time.Duration, coverage float64)
   func (ol *OrchestratorLogger) LogWSCheckpoint(wsID string, checkpointFile string)
   ```

2. Integrate into FeatureCoordinator:
   - Add logger field
   - Log orchestration start/completion
   - Log workstream events
   - Include correlation ID in all logs

3. Integrate into Orchestrator:
   - Add logger field
   - Log dependency graph building
   - Log execution order
   - Log retry attempts

4. Add tests:
   - Test correlation ID generation and propagation
   - Test log output capture
   - Test structured field presence

## Quality Gates

- Test coverage ≥80%
- go vet clean
- No fmt.Printf for logging (use slog instead)
- JSON handler for production logs

## Dependencies

None (uses Go stdlib log/slog)

## Estimated Scope

- ~100 LOC implementation
- ~150 LOC tests
- Duration: 2 hours
