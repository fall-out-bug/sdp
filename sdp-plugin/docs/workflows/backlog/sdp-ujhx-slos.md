# sdp-ujhx: Define Orchestrator-Specific SLOs

> **Issue ID:** sdp-ujhx
> **Type:** Bug (Critical - P0)
> **Priority:** 0
> **Source:** SRE Review finding

## Goal

Define Service Level Objectives (SLOs) and Service Level Indicators (SLIs) for orchestrator operations to enable production monitoring and alerting.

## Problem

SRE review identified that orchestrator lacks specific SLOs for:
- Checkpoint save latency
- Workstream execution time
- Dependency graph build time
- Checkpoint recovery success rate

Without SLOs, we cannot:
- Detect performance degradation
- Set up meaningful alerts
- Measure operational health
- Validate performance requirements

## Acceptance Criteria

### AC1: SLO Documentation
- [ ] Create SLO document in `docs/slos/`
- [ ] Define SLIs (Service Level Indicators) for each orchestrator operation
- [ ] Define SLO targets (e.g., 95th percentile latency)
- [ ] Document error budget calculations

### AC2: Checkpoint Save Latency
- [ ] Target: p95 < 100ms
- [ ] SLI: Time to serialize checkpoint to disk
- [ ] Alert if p95 > 150ms (error budget burn)

### AC3: Workstream Execution Time
- [ ] Target: p95 < 30min per workstream
- [ ] SLI: Time from WS start to completion
- [ ] Alert if p95 > 45min

### AC4: Dependency Graph Build
- [ ] Target: p95 < 5s for 100 nodes
- [ ] SLI: Time to parse dependencies and build graph
- [ ] Alert if p95 > 10s

### AC5: Checkpoint Recovery Success
- [ ] Target: 99.9% success rate
- [ ] SLI: Successful recovery / Total recovery attempts
- [ ] Alert if success rate < 99.5%

### AC6: Integration with Logging
- [ ] Log SLO measurements
- [ ] Include SLO status in orchestration_complete event
- [ ] Export metrics for monitoring system

## Scope Files

**NEW:**
- `sdp-plugin/docs/slos/orchestrator.md` (~200 LOC) - SLO definitions
- `sdp-plugin/internal/orchestrator/slos.go` (~150 LOC) - SLO tracking
- `sdp-plugin/internal/orchestrator/slos_test.go` (~200 LOC) - SLO tests

**MODIFY:**
- `sdp-plugin/internal/orchestrator/feature_coordinator.go` (add SLO tracking)
- `sdp-plugin/internal/orchestrator/orchestrator.go` (add SLO tracking)

## Implementation Steps

1. Create SLO documentation:
   ```markdown
   # Orchestrator SLOs

   ## Overview
   ## SLIs (Service Level Indicators)
   ## SLO Targets
   ## Error Budgets
   ## Alerting Thresholds
   ```

2. Implement SLO tracker:
   ```go
   type SLOTracker struct {
       mu              sync.RWMutex
       measurements    map[string]*Metric
       logger          *OrchestratorLogger
   }

   func (st *SLOTracker) RecordCheckpointSave(duration time.Duration)
   func (st *SLOTracker) RecordWSExecution(wsID string, duration time.Duration)
   func (st *SLOTracker) RecordGraphBuild(nodeCount int, duration time.Duration)
   func (st *SLOTracker) RecordRecovery(success bool)
   func (st *SLOTracker) GetSLOStatus() SLOStatus
   ```

3. Integrate into orchestrator:
   - Track checkpoint save latency
   - Track workstream execution time
   - Track graph build time
   - Track recovery success rate

4. Add tests:
   - Test SLO recording
   - Test SLO calculation (percentiles)
   - Test SLO status evaluation
   - Test error budget tracking

## SLO Definitions

### Checkpoint Save Latency
**SLI:** 95th percentile time to serialize checkpoint to disk
**Target:** < 100ms
**Error Budget:** 5% of requests can exceed 100ms
**Alert:** p95 > 150ms

### Workstream Execution Time
**SLI:** 95th percentile time from WS start to completion
**Target:** < 30min
**Error Budget:** 5% of workstreams can exceed 30min
**Alert:** p95 > 45min

### Dependency Graph Build
**SLI:** 95th percentile time to build dependency graph
**Target:** < 5s for 100 nodes
**Error Budget:** 5% of builds can exceed 5s
**Alert:** p95 > 10s

### Checkpoint Recovery Success
**SLI:** Successful recoveries / Total recovery attempts
**Target:** 99.9%
**Error Budget:** 0.1% failure rate
**Alert:** < 99.5% success

## Quality Gates

- Test coverage ≥80%
- go vet clean
- SLO calculations verified (percentile computation)
- No hardcoded thresholds (use constants)

## Dependencies

- OrchestratorLogger (for SLO event logging)
- Existing checkpoint system

## Estimated Scope

- ~200 LOC documentation
- ~150 LOC implementation
- ~200 LOC tests
- Duration: 2 hours
