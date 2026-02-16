# F068: UX Foundation and First-Run Experience — Workstream Summary

> Feature: F068 | Priority: P0 | 5 workstreams
> Theme: learning curve, discoverability, and time-to-first-value

## Goal

Make SDP usable for a new team in under 15 minutes from install to first successful workstream execution.

## Dependency Graph

```text
00-068-01 UX Baseline + Metrics
   ├─► 00-068-02 First-Run Setup Wizard
   ├─► 00-068-03 Command Discoverability + Help IA
   └─► 00-068-04 Quickstart Templates + Demo Flow

00-068-05 UX Rollout + Learning-Curve Gate
    depends on: 02, 03, 04
```

## Workstream Registry

| WS ID | Title | Priority | Size | Depends On |
|-------|-------|----------|------|------------|
| 00-068-01 | UX Baseline: Learning Curve and Time-to-First-Value | P0 | SMALL | — |
| 00-068-02 | Guided First-Run Setup (`sdp init --guided`) | P0 | MEDIUM | 01 |
| 00-068-03 | Help/Status Information Architecture | P0 | MEDIUM | 01 |
| 00-068-04 | Quickstart Templates and Demo Mode | P1 | MEDIUM | 01 |
| 00-068-05 | UX Rollout and Learning-Curve Release Gate | P0 | SMALL | 02,03,04 |

## Execution Phases

### Phase 1: Baseline
- WS-01 to freeze metrics and current pain points.

### Phase 2: UX Core (parallel)
- WS-02 setup guidance
- WS-03 discoverability
- WS-04 quickstart and templates

### Phase 3: Production Rollout
- WS-05 release gate and docs updates.

## Feature-Level Deliverables

- Standardized onboarding metrics and targets
- Guided setup path for first-time users
- Clear command discovery surface (`help`, `status`, `next`)
- Reproducible quickstart flow with known-good templates
- Release gate tied to UX KPIs (not only code quality KPIs)
