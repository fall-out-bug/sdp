# F070: Failure and Recovery UX — Workstream Summary

> Feature: F070 | Priority: P0 | 5 workstreams
> Theme: fast diagnosis, safe recovery, low frustration

## Goal

Cut mean-time-to-recovery for failed workstreams by making failures explicit, actionable, and resumable.

## Dependency Graph

```text
00-070-01 Failure Taxonomy + Error Contract
   ├─► 00-070-02 Recovery Playbook Engine
   ├─► 00-070-03 Resume/Checkpoint UX Improvements
   └─► 00-070-04 Post-Failure Diagnostics Report

00-070-05 Recovery Runbook + Adoption
    depends on: 02, 03, 04
```

## Workstream Registry

| WS ID | Title | Priority | Size | Depends On |
|-------|-------|----------|------|------------|
| 00-070-01 | Failure Taxonomy and Error Codes | P0 | SMALL | — |
| 00-070-02 | Recovery Playbook Engine | P0 | MEDIUM | 01 |
| 00-070-03 | Resume/Checkpoint UX Hardening | P0 | MEDIUM | 01 |
| 00-070-04 | Post-Failure Diagnostics and Triage Report | P1 | MEDIUM | 01 |
| 00-070-05 | Recovery Documentation and Rollout Gate | P0 | SMALL | 02,03,04 |

## Execution Phases

### Phase 1: Taxonomy
- WS-01 standardizes failure semantics.

### Phase 2: Recovery System
- WS-02 playbooks
- WS-03 resume flow
- WS-04 diagnostics report

### Phase 3: Operationalization
- WS-05 docs, runbooks, release gate.

## Feature-Level Deliverables

- Standard failure classes and error-code contract
- Recovery playbooks surfaced directly in failures
- Reliable checkpoint resume flow with clear state
- Diagnostics report usable by humans and agents
