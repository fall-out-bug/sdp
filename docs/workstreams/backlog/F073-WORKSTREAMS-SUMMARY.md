# F073: Trust and Explainability UX — Workstream Summary

> Feature: F073 | Priority: P0 | 5 workstreams
> Theme: make decisions explainable, traceable, and auditable

## Goal

Increase user trust by exposing why SDP made decisions, which gate failed, and how evidence supports outcomes.

## Dependency Graph

```text
00-073-01 Explainability Model + Levels
   ├─► 00-073-02 Recommendation Rationale UX
   ├─► 00-073-03 Evidence Trace UX Unification
   └─► 00-073-04 Policy Transparency + Gate Reasoning

00-073-05 Public Trust Pack + GA Criteria
    depends on: 02, 03, 04
```

## Workstream Registry

| WS ID | Title | Priority | Size | Depends On |
|-------|-------|----------|------|------------|
| 00-073-01 | Explainability Model and Disclosure Levels | P0 | SMALL | — |
| 00-073-02 | Recommendation Rationale in CLI/Reports | P0 | MEDIUM | 01 |
| 00-073-03 | Evidence Trace UX Unification | P0 | MEDIUM | 01 |
| 00-073-04 | Policy Transparency and Gate Reasoning | P0 | MEDIUM | 01 |
| 00-073-05 | Public Trust Pack and Release Criteria | P1 | SMALL | 02,03,04 |

## Execution Phases

### Phase 1: Model
- WS-01 defines explainability levels and mandatory rationale fields.

### Phase 2: UX Surfaces
- WS-02 rationale display
- WS-03 evidence trace improvements
- WS-04 gate transparency

### Phase 3: Release
- WS-05 public OSS trust packaging and GA readiness.

## Feature-Level Deliverables

- Unified explainability contract across CLI and docs
- Human-readable rationale for recommendations
- End-to-end evidence trace usability improvements
- Transparent policy and gate failure explanations
