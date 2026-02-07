# Orchestrator Service Level Objectives (SLOs)

## Overview

This document defines the Service Level Objectives (SLOs) and Service Level Indicators (SLIs) for the SDP Orchestrator component. These objectives ensure reliable, performant execution of workstream orchestration in production environments.

**Target Audience:** SRE team, Orchestrator developers, Monitoring team

**Last Updated:** 2026-02-06

## SLIs (Service Level Indicators)

Service Level Indicators are quantitative measures of service level.

### SLI 1: Checkpoint Save Latency

**Definition:** Time to serialize checkpoint state to disk

**Measurement:** 95th percentile latency over 7-day rolling window

**Metric Source:** `SLOTracker.RecordCheckpointSave()`

**Unit:** Seconds

**Calculation:**
```
p95_checkpoint_save = percentile(checkpoint_save_durations, 95)
```

---

### SLI 2: Workstream Execution Time

**Definition:** Time from workstream start to completion (including retries)

**Measurement:** 95th percentile duration over 7-day rolling window

**Metric Source:** `SLOTracker.RecordWSExecution()`

**Unit:** Seconds

**Calculation:**
```
p95_ws_execution = percentile(workstream_execution_durations, 95)
```

---

### SLI 3: Dependency Graph Build Time

**Definition:** Time to parse workstream dependencies and build execution graph

**Measurement:** 95th percentile duration for graphs with ≥100 nodes

**Metric Source:** `SLOTracker.RecordGraphBuild()`

**Unit:** Seconds

**Calculation:**
```
p95_graph_build = percentile(graph_build_durations, 95)
```

---

### SLI 4: Checkpoint Recovery Success Rate

**Definition:** Ratio of successful checkpoint recoveries to total recovery attempts

**Measurement:** Success percentage over 7-day rolling window

**Metric Source:** `SLOTracker.RecordRecovery()`

**Unit:** Percentage (0-1)

**Calculation:**
```
recovery_success_rate = successful_recoveries / total_recovery_attempts
```

---

## SLO Targets

Service Level Objectives are the specific goals for each SLI.

| SLI | Target | Measurement Window |
|-----|--------|-------------------|
| **Checkpoint Save Latency** | p95 < 100ms | 7-day rolling |
| **Workstream Execution Time** | p95 < 30min | 7-day rolling |
| **Dependency Graph Build** | p95 < 5s (100 nodes) | 7-day rolling |
| **Checkpoint Recovery Success** | 99.9% success | 7-day rolling |

## Error Budgets

Error budgets represent the allowable amount of "bad" service within a measurement window.

### Checkpoint Save Latency

**Target:** p95 < 100ms
**Budget:** 5% of checkpoint saves can exceed 100ms
**Budget Burn Rate:** 100ms breaches per 1000 saves

### Workstream Execution Time

**Target:** p95 < 30min
**Budget:** 5% of workstreams can exceed 30min
**Budget Burn Rate:** 30min breaches per 1000 workstreams

### Dependency Graph Build

**Target:** p95 < 5s
**Budget:** 5% of builds can exceed 5s
**Budget Burn Rate:** 5s breaches per 1000 builds

### Checkpoint Recovery Success

**Target:** 99.9% success
**Budget:** 0.1% failure rate allowed
**Budget Burn Rate:** 1 failure per 1000 attempts

## Alerting Thresholds

Alerts fire when error budget is at risk or SLO is breached.

| SLO | Warning Threshold | Critical Threshold |
|-----|------------------|-------------------|
| **Checkpoint Save Latency** | p95 > 120ms | p95 > 150ms |
| **Workstream Execution Time** | p95 > 40min | p95 > 45min |
| **Dependency Graph Build** | p95 > 8s | p95 > 10s |
| **Checkpoint Recovery Success** | < 99.7% success | < 99.5% success |

### Alert Actions

**Warning Alerts:**
- Notify SRE team
- Investigate performance degradation
- Check system resources (CPU, memory, disk I/O)

**Critical Alerts:**
- Page on-call SRE
- Pause new orchestration executions
- Initiate incident response

## Implementation

### SLO Tracker

The orchestrator includes an SLO tracker that records metrics and evaluates compliance:

```go
import "github.com/fall-out-bug/sdp/internal/orchestrator"

// Create tracker with logger
tracker := orchestrator.NewSLOTracker(logger)

// Record metrics
tracker.RecordCheckpointSave(duration)
tracker.RecordWSExecution(wsID, duration)
tracker.RecordGraphBuild(nodeCount, duration)
tracker.RecordRecovery(success)

// Check compliance
status := tracker.GetSLOStatus()
if !status.OverallSLOCompliance {
    // Handle SLO breach
}
```

### SLOStatus Structure

```go
type SLOStatus struct {
    CheckpointSaveLatency       float64  // p95 in seconds
    CheckpointSaveLatencyOK     bool
    WSExecutionTime             float64  // p95 in seconds
    WSExecutionTimeOK           bool
    GraphBuildTime              float64  // p95 in seconds
    GraphBuildTimeOK            bool
    RecoverySuccessRate         float64  // 0-1
    RecoverySuccessRateOK       bool
    OverallSLOCompliance        bool
}
```

## Integration with Logging

SLO measurements are automatically logged through the `OrchestratorLogger`:

- **SLO Breach Logged:** When measurement exceeds target
- **SLO Status Included:** In `orchestration_complete` event
- **Structured Fields:** All SLO metrics available for monitoring systems

### Example Log Entry

```json
{
  "level": "warn",
  "msg": "SLO breach: graph build time exceeds target",
  "duration_seconds": 6.5,
  "target_seconds": 5.0,
  "node_count": 150,
  "timestamp": "2026-02-06T20:00:00Z",
  "correlation_id": "feature-123-1754544000"
}
```

## Monitoring Integration

### Prometheus Metrics (Future)

SLO metrics will be exported to Prometheus for dashboarding and alerting:

```promql
# Checkpoint save latency p95
histogram_quantile(0.95, rate(orchestrator_checkpoint_save_duration_seconds_bucket[7d]))

# Workstream execution time p95
histogram_quantile(0.95, rate(orchestrator_ws_execution_duration_seconds_bucket[7d]))

# Recovery success rate
sum(rate(orchestrator_recovery_success_total[7d])) / sum(rate(orchestrator_recovery_attempts_total[7d]))
```

### Grafana Dashboards

Create dashboards showing:
- SLO compliance status (pass/fail)
- p95 latencies over time
- Error budget consumption
- Success rate trends

## Performance Testing

### Load Testing

Validate SLO compliance under load:

```bash
# Test checkpoint save latency (target: p95 < 100ms)
for i in {1..1000}; do
  time orchestrator checkpoint save
done

# Test graph build time (target: p95 < 5s for 100 nodes)
orchestrator graph build --nodes=100 --iterations=100

# Test recovery success (target: 99.9%)
orchestrator recovery test --attempts=1000
```

### SLO Validation

Run SLO validation suite:

```bash
go test ./internal/orchestrator/... -run TestSLO -v
```

## Change Management

### SLO Modification Process

1. **Propose Change:** Create ADR documenting SLO change rationale
2. **Validate:** Test against production traffic shadow
3. **Review:** SRE team approval required
4. **Implement:** Update SLO constants in `slos.go`
5. **Communicate:** Notify stakeholders of SLO change
6. **Monitor:** Track compliance for 30 days post-change

### Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-02-06 | Initial SLO definitions |

## References

- [Google SRE Book - Service Level Objectives](https://sre.google/sre-book/service-level-objectives/)
- [Monitoring SOP](../sops/monitoring.md) (TODO)
- [Incident Response SOP](../sops/incident-response.md) (TODO)

## Appendix: SLO Calculations

### Percentile Calculation

SLOs use the 95th percentile (p95), which means 95% of measurements must meet the target.

**Example:** For checkpoint save latency:
- 1000 checkpoint saves recorded
- p95 = 950th value when sorted
- If 950th value = 90ms, SLO is met (target: < 100ms)
- If 950th value = 110ms, SLO is breached

### Error Budget Calculation

**Weekly Budget Calculation (168 hours):**

```
# For checkpoint save latency (p95 < 100ms)
max_breaches_per_week = total_saves * 0.05

# Example: 10,000 saves per week
max_breaches = 10,000 * 0.05 = 500 breaches allowed
```

**Budget Burn Rate:**

```
burn_rate = breaches_per_hour / max_breaches_per_week

# Example: 50 breaches in first hour
burn_rate = 50 / 500 = 0.1 (10% of budget burned in 1 hour)
# At this rate, entire budget exhausted in 10 hours
```

---

**Document Owner:** SRE Team
**Review Frequency:** Quarterly
**Next Review:** 2026-05-06
