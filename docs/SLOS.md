# SLOs/SLIs for SDP CLI Tool

> **Version:** 1.0
> **Last Updated:** 2026-02-06
> **Scope:** SDP Developer CLI Tool

## Context

SDP is a **developer CLI tool**, not a web service. Traditional service SLOs (uptime, latency, p95/p99 response times) do not apply.

Instead, we define **CLI-specific SLOs** that measure:
- Tool availability (can developers get and use it?)
- Installation success (does setup work?)
- Command performance (is it fast enough?)
- Reliability (does it crash?)
- Telemetry (are we monitoring usage?)

---

## SLO 1: Binary Availability

### **Objective**
Developers can successfully download the SDP binary within 30 seconds.

### **SLI: Download Success Rate**
- **Metric:** Percentage of successful binary downloads
- **Measurement:** GitHub release downloads / unique visitors to download page
- **Target:** ≥ 99%
- **Measurement Window:** Rolling 7 days

### **SLI: Download Latency**
- **Metric:** Time to complete binary download (start → finish)
- **Target:** p95 ≤ 10 seconds
- **Measurement:** Telemetry timestamp (download_start) → (download_complete)

### **Error Budget**
- 1% failure rate = ~43 minutes downtime per month acceptable

### **Alerting**
- **Critical:** Download success rate < 95% for 1 hour
- **Warning:** Download success rate < 98% for 1 hour

---

## SLO 2: Installation Success Rate

### **Objective**
Developers can install SDP and run `sdp --help` successfully on first try.

### **SLI: Installation Success Rate**
- **Metric:** Percentage of successful installations (binary executes without error)
- **Measurement:** Telemetry: `sdp init` executions / total installations
- **Target:** ≥ 98%
- **Measurement Window:** Rolling 7 days

### **SLI: Platform Coverage**
- **Metric:** Success rate per platform (darwin-arm64, darwin-amd64, linux-amd64, windows-amd64)
- **Target:** ≥ 95% per platform
- **Measurement Window:** Rolling 30 days

### **Error Budget**
- 2% failure rate = ~1.4 hours downtime per month acceptable

### **Alerting**
- **Critical:** Installation success < 90% for any platform
- **Warning:** Installation success < 95% for any platform

---

## SLO 3: Command Execution Latency

### **Objective**
SDP commands execute fast enough to not interrupt developer workflow.

### **SLI: Command Latency (by command)**
- **Metric:** p95 command execution time (start → exit)
- **Target:**
  - `sdp --help`: ≤ 100ms
  - `sdp init`: ≤ 2 seconds (interactive)
  - `sdp doctor`: ≤ 500ms
  - `sdp quality check`: ≤ 5 seconds
  - `sdp guard activate`: ≤ 100ms
  - `sdp guard complete`: ≤ 100ms
  - `sdp verify`: ≤ 3 seconds

- **Measurement Window:** Rolling 24 hours
- **Measurement Tool:** Telemetry (command_start, command_complete timestamps)

### **SLI: Command Timeouts**
- **Metric:** Percentage of commands timing out (>30 seconds)
- **Target:** ≤ 0.1%
- **Measurement Window:** Rolling 7 days

### **Alerting**
- **Critical:** Any command p95 > 2x target for 1 hour
- **Warning:** Any command p95 > 1.5x target for 1 hour

---

## SLO 4: Error Rate by Command

### **Objective**
SDP commands execute reliably without crashes or errors.

### **SLI: Command Error Rate**
- **Metric:** Percentage of commands exiting with error (non-zero exit code)
- **Target:** ≤ 5% overall
- **Target by Command:**
  - `sdp init`: ≤ 2% (interactive, sensitive to environment)
  - `sdp doctor`: ≤ 1% (readiness checks)
  - `sdp quality check`: ≤ 10% (expected failures for low quality)
  - `sdp guard activate/complete`: ≤ 1%
  - `sdp verify`: ≤ 5%

- **Measurement Window:** Rolling 7 days
- **Measurement Tool:** Telemetry (command_exit_code)

### **SLI: Crash Rate**
- **Metric:** Percentage of commands crashing (signal: SIGSEGV, SIGABRT, panic)
- **Target:** ≤ 0.01%
- **Measurement Window:** Rolling 30 days

### **Error Budget**
- 5% error rate = ~3.6 hours downtime per month acceptable

### **Alerting**
- **Critical:** Error rate > 10% overall
- **Critical:** Crash rate > 0.1%
- **Warning:** Error rate > 7.5% overall

---

## SLO 5: Telemetry Opt-In Rate

### **Objective**
Understand how developers use SDP without being intrusive.

### **SLI: Opt-In Rate**
- **Metric:** Percentage of installations with telemetry enabled
- **Target:** ≥ 60%
- **Measurement Window:** Rolling 30 days

### **SLI: Data Quality**
- **Metric:** Percentage of telemetry events with valid schema
- **Target:** ≥ 99.9%
- **Measurement Window:** Rolling 7 days

### **Privacy Requirements**
- **No PII** (personally identifiable information) in telemetry
- **No code snippets** from user projects
- **Opt-out mechanism:** `sdp telemetry disable`

**Note:** This is a **directional metric** (higher is better), not a hard SLO.

---

## SLO 6: User Satisfaction

### **Objective**
Developers find SDP useful and recommend it.

### **SLI: Retention Rate**
- **Metric:** Percentage of users who run SDP commands >10 times in 30 days
- **Target:** ≥ 70%
- **Measurement Window:** Rolling 30 days

### **SLI: GitHub Stars**
- **Metric:** Growth rate of GitHub repository stars
- **Target:** ≥ 10% month-over-month
- **Measurement Window:** Monthly

**Note:** These are **lagging indicators**, not real-time SLOs.

---

## Non-Goals (What We DON'T Measure)

1. **Code Coverage**
   - Coverage is a **quality gate**, not an SLO
   - Measured during `sdp quality check`, not telemetry

2. **Test Pass Rate**
   - Tests failing is **expected** during TDD Red phase
   - Not a reliability indicator

3. **Documentation Completeness**
   - Qualitative measure, hard to quantify
   - Validated via user feedback

4. **Feature Usage**
   - Which commands are used is **directional**
   - Low usage ≠ bad (may be power user features)

---

## Measurement Architecture

### **Telemetry Events**

```json
{
  "event": "command_execution",
  "timestamp": "2026-02-06T12:00:00Z",
  "session_id": "uuid",
  "command": "sdp init",
  "args": ["--project-type", "python"],
  "exit_code": 0,
  "duration_ms": 1234,
  "platform": "darwin-arm64",
  "version": "1.0.0",
  "telemetry_enabled": true
}
```

### **Privacy**

**Collected:**
- Command name
- Exit code
- Duration
- Platform (OS, architecture)
- Version

**NOT Collected:**
- File paths
- Code snippets
- User input (answers to questions)
- Project names
- User IDs or IP addresses

### **Storage**

- **Backend:** GitHub Releases (download count)
- **Telemetry:** OpenTelemetry + Honeycomb/New Relic (or self-hosted)
- **Retention:** 90 days for telemetry events

---

## SLO Review Process

### **Frequency**
- **Quarterly** SLO reviews (every 3 months)
- **Annual** SLO redefinition (adjust targets based on maturity)

### **Review Checklist**
- [ ] Are targets realistic given current maturity?
- [ ] Are measurements still accurate?
- [ ] Any new SLOs needed? Any obsolete SLOs to remove?
- [ ] Are error budgets being consumed appropriately?

### **SLO Changes**
SLO changes require:
1. Proposal document (draft-SLOS-v1.1.md)
2. 2-week review period
3. User feedback (GitHub issue)
4. Approval by maintainer
5. Version bump (SLOS.md v1.0 → v1.1)

---

## Dashboard Configuration

### **Recommended Metrics (Grafana/Loki)**

**Panel 1: Availability**
- Download success rate (7-day trend)
- Installation success rate by platform (30-day)

**Panel 2: Performance**
- Command latency p95 by command (heatmap)
- Timeout rate (7-day trend)

**Panel 3: Reliability**
- Error rate by command (7-day trend)
- Crash rate (30-day trend)

**Panel 4: Usage**
- Opt-in rate (30-day trend)
- Retention rate (30-day)
- GitHub stars growth (monthly)

---

## Examples

### **Example 1: Download Success Rate**

**Measurement:**
- Week 1: 500 downloads, 490 successful = 98% ✅
- Week 2: 600 downloads, 570 successful = 95% ⚠️
- Week 3: 700 downloads, 686 successful = 98% ✅

**Action Required for Week 2:**
- Investigate download failures (CDN issue? corrupted binary?)
- Check GitHub Actions release upload logs

---

### **Example 2: Command Latency**

**Measurement:**
```
sdp doctor p95 latency:
- Week 1: 450ms ✅ (target: ≤500ms)
- Week 2: 650ms ❌ (target: ≤500ms)
- Week 3: 420ms ✅ (target: ≤500ms)
```

**Action Required for Week 2:**
- Profile `sdp doctor` command
- Identify slow subcommand (Git check? Python version check?)
- Optimize or parallelize checks

---

### **Example 3: Error Rate**

**Measurement:**
```
sdp quality check error rate:
- Week 1: 8% ❌ (target: ≤5%)
- Week 2: 6% ⚠️ (target: ≤5%)
- Week 3: 4% ✅ (target: ≤5%)
```

**Action Required for Week 1:**
- Investigate error types (coverage failures? type checking?)
- Check if common errors are false positives
- Adjust quality gates if too strict

---

## References

- **Google SRE Book:** https://sre.google/sre-book/table-of-contents/
- **SLI/SLO Framework:** https://sre.google/sre-book/service-level-objectives/
- **CLI Telemetry Best Practices:** https://clig.dev/#cli-telemetry

---

## Appendix: Error Budget Calculation

### **Example: SLO 2 - Installation Success Rate (98%)**

**Error Budget:** 2% failure rate per month

**Monthly Minutes:** 43,200 minutes (30 days × 24 hours × 60 minutes)

**Error Budget:** 43,200 × 0.02 = **864 minutes/month** = **14.4 hours/month**

**Usage:**
- If installation is broken for 14 hours, we can still meet SLO
- If broken for >14 hours, we must:
  1. Declare incident
  2. Fix installation
  3. Postmortem

---

**Version:** 1.0
**Next Review:** 2026-05-06
