# SDP Verify Action - SLOs and SLIs

## Service Level Objectives (SLOs)

### Reliability

| SLO | Target | Measurement Window |
|-----|--------|-------------------|
| Availability (uptime) | 99.5% | 30 days |
| Success Rate | 99% | 7 days |
| Data Integrity (evidence chain) | 100% | Per execution |

### Performance

| SLO | Target | Measurement Window |
|-----|--------|-------------------|
| Execution Time (p50) | < 2 minutes | Per execution |
| Execution Time (p95) | < 3 minutes | Per execution |
| Execution Time (p99) | < 5 minutes | Per execution |
| Binary Download Time | < 30 seconds | Per download |
| PR Comment Post Time | < 10 seconds | Per comment |

### Correctness

| SLO | Target | Measurement Window |
|-----|--------|-------------------|
| Gate Accuracy | 100% | Per execution |
| False Positive Rate | < 1% | 30 days |
| False Negative Rate | < 0.5% | 30 days |

## Service Level Indicators (SLIs)

### 1. Execution Success Rate

**Definition:** Percentage of action runs that complete successfully (exit code 0).

**Formula:**
```
Success Rate = (Successful Runs / Total Runs) × 100%
```

**Measurement:**
- GitHub Actions workflow runs
- Filter: `sdp-verify-dogfood.yml` and external usage
- Success: `conclusion == "success"`

**Target:** ≥ 99%

**Alert Threshold:** < 98% (investigate immediately)

### 2. Execution Duration (Latency)

**Definition:** Time from action start to completion.

**Formula:**
```
Duration = (Completion Time - Start Time)
```

**Measurement:**
- GitHub Actions API: `jobs[].steps[].duration_ms`
- Track: Install SDP CLI, Run Gates, Post Comment

**Target (p50):** < 120 seconds (2 minutes)
**Target (p95):** < 180 seconds (3 minutes)
**Target (p99):** < 300 seconds (5 minutes)

**Alert Threshold:** p95 > 240 seconds (4 minutes)

### 3. Binary Download Success Rate

**Definition:** Percentage of successful SDP CLI binary downloads.

**Formula:**
```
Download Success Rate = (Successful Downloads / Total Download Attempts) × 100%
```

**Measurement:**
- Action logs: "✅ Download successful" vs "❌ Failed to download"
- Count retries vs failures

**Target:** ≥ 99.5% (including retries)

**Alert Threshold:** < 99%

### 4. PR Comment Success Rate

**Definition:** Percentage of successful PR comment posts/updates.

**Formula:**
```
Comment Success Rate = (Successful Comments / Total Comment Attempts) × 100%
```

**Measurement:**
- Action logs: "✅ Comment posted/updated" vs "⚠️ Warning: Failed to post"
- Track both create and update operations

**Target:** ≥ 98% (non-blocking failures allowed)

**Alert Threshold:** < 95%

### 5. Evidence Chain Integrity

**Definition:** Percentage of runs with intact evidence chain.

**Formula:**
```
Chain Integrity = (Valid Chains / Total Evidence Checks) × 100%
```

**Measurement:**
- Action output: `sdp log trace --verify`
- Parse: "✅ Evidence chain integrity verified"

**Target:** 100% (critical for compliance)

**Alert Threshold:** < 100% (investigate immediately)

### 6. Gate Accuracy

**Definition:** Correctness of gate results (true positives/negatives).

**Formula:**
```
Gate Accuracy = (Correct Results / Total Gate Executions) × 100%
```

**Measurement:**
- Manual sampling of failed gates
- Compare with local `sdp verify` results

**Target:** 100% (critical for trust)

**Alert Threshold:** < 100%

## Error Budget

### Monthly Error Budget Calculation

**Monthly Runs:** ~1,000 (estimated, 50 PRs × 20 runs/day)

**Target Success Rate:** 99%

**Allowed Failures:**
```
Error Budget = Total Runs × (1 - Target Success Rate)
Error Budget = 1,000 × (1 - 0.99) = 10 failures/month
```

### Error Budget Consumption

**Per Failure Impact:**
```
Budget Consumed = (1 / Error Budget) × 100%
Budget Consumed = (1 / 10) × 100% = 10% per failure
```

**Burn Rate:**
- **Fast Burn:** > 10 failures in < 10 days (pause releases)
- **Normal Burn:** 10 failures in 30 days (acceptable)
- **Slow Burn:** < 10 failures in 30 days (healthy)

## Monitoring

### Metrics Collection

**GitHub Actions Logs:**
```yaml
- Enable workflow run logs retention: 90 days
- Export metrics to:
  - GitHub Actions API (built-in)
  - External monitoring (optional: Datadog, Prometheus)
```

**Key Metrics to Track:**
1. `workflow_run.conclusion` (success/failure)
2. `workflow_run.duration_ms` (execution time)
3. `job.steps[].name` contains "✅" or "❌"
4. Retry counts (from action logs)

### Alerting

**Critical Alerts (Page Immediately):**
- Success rate < 95%
- Evidence chain integrity failure
- p99 execution time > 10 minutes

**Warning Alerts (Investigate within 24h):**
- Success rate < 98%
- p95 execution time > 4 minutes
- Binary download failure rate > 2%
- PR comment failure rate > 5%

### Dashboards

**Recommended Dashboard Widgets:**
1. Success Rate (last 7/30 days)
2. Execution Duration (p50, p95, p99)
3. Failure Reasons (breakdown)
4. Evidence Chain Integrity Status
5. Binary Download Success Rate
6. PR Comment Success Rate

## Reporting

### Weekly Report

**Metrics:**
- Total runs, success rate
- p50/p95/p99 latency
- Top 3 failure reasons
- Error budget remaining

### Monthly Report

**Metrics:**
- All weekly metrics
- SLO attainment (%)
- Error budget consumption
- Incidents and postmortems
- Improvement roadmap

## Incident Response

### Severity Levels

**SEV1 (Critical):**
- Complete action failure (> 50% failure rate)
- Evidence chain integrity compromised
- Action: Investigate immediately, pause releases

**SEV2 (High):**
- Degraded performance (p95 > 5 minutes)
- Intermittent failures (10-50% failure rate)
- Action: Investigate within 1 hour

**SEV3 (Medium):**
- Minor performance degradation
- Low failure rate (< 10%)
- Action: Investigate within 24 hours

## Continuous Improvement

### Quarterly Review

**Topics:**
- SLO attainment review
- Error budget analysis
- Performance optimization opportunities
- Reliability improvements
- User feedback integration

---

**Last Updated:** 2026-02-10
**Next Review:** 2026-05-10
**Owner:** SDP Maintainers
