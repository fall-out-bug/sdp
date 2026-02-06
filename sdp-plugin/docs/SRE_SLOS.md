# CLI-Specific SLOs/SLIs for SDP

## Overview

SDP is a developer CLI tool distributed via GitHub Releases, not a hosted service. Traditional service SLOs (uptime, latency, availability) do not apply. This document defines CLI-specific metrics.

## SLOs (Service Level Objectives)

### Binary Availability
**Metric:** Download success rate from GitHub Releases
**Target:** 99.9%
**Measurement:** Failed downloads / Total download attempts
**Rationale:** Users must be able to install the tool

**Calculation:**
```
Binary Availability = (Successful Downloads / Total Download Attempts) × 100%
```

**Data Source:** GitHub Release analytics (download counts)

### Installation Success Rate
**Metric:** Successful installations / download attempts
**Target:** 95%
**Measurement:** Post-install smoke test (`sdp --version`)
**Rationale:** Installation should work flawlessly

**Calculation:**
```
Installation Success = (Successful `sdp --version` executions / Total installations) × 100%
```

**Data Source:** Telemetry opt-in data (command start events)

### Command Execution Latency
**Metric:** Median command duration
**Target:**
- `sdp init`: <5s
- `sdp doctor`: <1s
- `sdp parse`: <500ms
- Other commands: <2s

**Measurement:** Time from command invocation to completion
**Rationale:** Commands should feel responsive

**Data Source:** Telemetry (duration field in events)

### Error Rate
**Metric:** Failed commands / total commands
**Target:** <5% overall
**Measurement:** Commands with non-zero exit status
**Rationale:** Most commands should succeed

**Calculation:**
```
Error Rate = (Failed Commands / Total Commands) × 100%
```

**Data Source:** Telemetry (success/failure in events)

### Telemetry Opt-In Rate
**Metric:** Users enabling telemetry
**Target:** >20% (stretch goal)
**Measurement:** `~/.sdp/telemetry.jsonl` file exists
**Rationale:** Enough data to detect issues

**Data Source:** Cannot measure without privacy violation. Estimate based on support requests.

## SLIs (Service Level Indicators)

### Binary Size
**Indicator:** Compressed binary size by platform
**Target:**
- macOS ARM64: <8MB
- macOS AMD64: <8MB
- Linux AMD64: <8MB
- Windows AMD64: <10MB

**Measurement:** `ls -lh bin/sdp-*`

### Startup Time
**Indicator:** Time to first output
**Target:** <50ms (imperceptible delay)
**Measurement:** `time sdp --version`

### Test Coverage
**Indicator:** Code coverage percentage
**Target:** ≥75% average
**Measurement:** `go test -coverprofile=coverage.out`

### Documentation Coverage
**Indicator:** Commands with documentation
**Target:** 100% (all commands documented)
**Measurement:** Compare `sdp --help` output to docs/

## Error Budget

### Monthly Error Budget
**Target:** 99% uptime (allows ~7.2 hours downtime/month)
**Calculation:** 30 days × 24 hours × 1% = 7.2 hours

**Error Budget Consumption:**
- Critical bugs requiring immediate hotfix: 100% of budget
- Minor bugs: 10% of budget per incident
- Documentation issues: 0% of budget

## Monitoring

### Metrics Collection
1. **Telemetry Events:** Automatically collected on opt-in
   - Command start/complete
   - Success/failure
   - Duration
   - Error messages

2. **GitHub Analytics:** Passive collection
   - Release downloads
   - Clone counts
   - Issue/PR engagement

3. **Support Requests:** Manual tracking
   - GitHub issues (bug reports, feature requests)
   - GitHub Discussions (questions)
   - Community chat (if applicable)

### Alerting
**Alert on:**
- Binary Availability drops below 99%
- Error Rate exceeds 5% for >24 hours
- Critical security vulnerability (CVSS ≥ 8.0)
- Regression in test coverage (>5% drop)

**No Alert on:**
- Documentation updates (normal)
- Feature additions (normal)
- Performance degradation within SLO (normal)

## Reporting

### Monthly Report
**Contents:**
1. Binary Availability (with 3-month trend)
2. Installation Success Rate (with 3-month trend)
3. Top 5 Commands by Usage
4. Error Rate (with breakdown by command)
5. Median Latency by Command
6. Test Coverage Trend
7. Critical Incidents

**Distribution:**
- Posted to GitHub Discussions monthly
- Included in Release Notes
- Shared with team internally

## Measurement Methods

### Binary Availability
**How to measure:**
```bash
# GitHub Release analytics (manual check)
gh release view --repo fall-out-bug/sdp
```

**Alternative:**
- Download counters in README badges
- CDN logs (if using CDN)

### Installation Success Rate
**How to measure:**
```bash
# From telemetry data (opt-in only)
cd ~/.sdp
grep -c '"EventType":"CommandComplete"' telemetry.jsonl
```

**Limitation:** Cannot measure without telemetry opt-in.

### Command Execution Latency
**How to measure:**
```bash
# From telemetry data
cd ~/.sdp
jq '.duration' telemetry.jsonl | awk '{sum+=$1; count++} END {print sum/count}'
```

### Error Rate
**How to measure:**
```bash
# From telemetry data
cd ~/.sdp
jq -r 'select(.success == false) | .type' telemetry.jsonl | sort | uniq -c
```

## Targets Summary

| Metric | Target | Measurement Method |
|--------|--------|-------------------|
| Binary Availability | 99.9% | GitHub analytics |
| Installation Success | 95% | Telemetry (estimated) |
| Command Latency | <2s (p95) | Telemetry duration |
| Error Rate | <5% | Telemetry failures |
| Test Coverage | ≥75% | `go test -cover` |
| Binary Size | <10MB | `ls -lh bin/sdp-*` |
| Startup Time | <50ms | `time sdp --version` |

## Notes

1. **Passive Measurement:** Most metrics are passively collected via telemetry or GitHub analytics. No active monitoring infrastructure required.

2. **Privacy-First:** Telemetry is opt-in only. No PII collected. Users must explicitly enable.

3. **Best Effort:** Some metrics (Installation Success Rate, Telemetry Opt-In Rate) cannot be accurately measured without privacy violations.

4. **Focus on User Experience:** SLOs prioritize what users experience: installation, command speed, reliability.

5. **Continuous Improvement:** Targets will be adjusted based on real-world usage data.

## References

- [OWASP SLOs](https://owasp.org/www-project-slo/)
- [Google SRE Workbook](https://sre.google/sre-book/table-of-contents)
- [GitHub Release Analytics](https://docs.github.com/en/repositories/viewing-activity-and-data-for-your-repository/getting-insights-about-your-repository)
