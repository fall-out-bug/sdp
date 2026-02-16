# F069: Next-Step Engine and Guided Flow — Workstream Summary

> Feature: F069 | Priority: P0 | 5 workstreams
> Theme: reduce ambiguity after each command and error

## Goal

Give users deterministic "what to do next" guidance after every major SDP action.

## Dependency Graph

```text
00-069-01 Next-Step Engine Contract
   ├─► 00-069-02 Recommendation Engine (deterministic)
   ├─► 00-069-03 CLI Surfaces (status/help/error next-steps)
   └─► 00-069-04 Interactive Drive Loop Improvements

00-069-05 Metrics + Tuning + Rollout
    depends on: 02, 03, 04
```

## Workstream Registry

| WS ID | Title | Priority | Size | Depends On |
|-------|-------|----------|------|------------|
| 00-069-01 | Next-Step Decision Contract and State Model | P0 | SMALL | — |
| 00-069-02 | Deterministic Recommendation Engine | P0 | MEDIUM | 01 |
| 00-069-03 | Next-Step Surfaces in `status/help/error` | P0 | MEDIUM | 01 |
| 00-069-04 | Drive-Mode Guided Loop UX | P1 | MEDIUM | 02,03 |
| 00-069-05 | Recommendation Quality Metrics and Rollout | P0 | SMALL | 02,03,04 |

## Execution Phases

### Phase 1: Contract
- WS-01 defines state, outputs, and confidence levels.

### Phase 2: Engine + Surfaces
- WS-02 recommendation logic
- WS-03 CLI integration
- WS-04 interactive flow

### Phase 3: Stabilization
- WS-05 quality scoring and release criteria.

## Feature-Level Deliverables

- Shared "next-step" model for CLI and agents
- Deterministic recommendation output for common states
- Inline guidance in command output and failure cases
- Tunable quality metrics for recommendation usefulness
