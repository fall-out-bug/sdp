# F074: Layered Public OSS Adoption and Packaging Profiles — Workstream Summary

> Feature: F074 | Priority: P0 | 6 workstreams
> Theme: progressive adoption without full-stack lock-in

## Goal

Operationalize SDP as a progressive public OSS stack (`L0` -> `L1` -> `L2`) with explicit packaging boundaries.

## Dependency Graph

```text
00-074-01 Layer Contracts + Profile Spec
   ├─► 00-074-02 L0 Claude Plugin Distribution (MIT)
   ├─► 00-074-03 L1 Safety Bundle Distribution via Brew (MIT)
   ├─► 00-074-04 L2 Core Orchestrator Distribution via Brew (MIT)
   └─► 00-074-05 Public OSS Boundary and Decoupling Guardrails

00-074-06 Adoption Docs + Migration + Release Gates
    depends on: 02, 03, 04, 05
```

## Workstream Registry

| WS ID | Title | Priority | Size | Depends On |
|-------|-------|----------|------|------------|
| 00-074-01 | Layer Contracts and Install Profile Specification | P0 | SMALL | 00-073-05 |
| 00-074-02 | L0 Protocol Distribution via Claude Plugin (MIT) | P0 | MEDIUM | 01 |
| 00-074-03 | L1 Safety Bundle Packaging via Homebrew (MIT) | P0 | MEDIUM | 01 |
| 00-074-04 | L2 Core Orchestrator Packaging via Homebrew (MIT) | P0 | MEDIUM | 01,03 |
| 00-074-05 | Public OSS Boundary and Decoupling Guardrails | P0 | SMALL | 01 |
| 00-074-06 | Adoption Documentation, Migration Paths, and Release Gates | P1 | SMALL | 02,03,04,05 |

## Execution Phases

### Phase 1: Contracts
- WS-01 defines profiles and non-breaking boundaries.

### Phase 2: Distribution Tracks
- WS-02 L0 plugin track
- WS-03 L1 brew track
- WS-04 L2 brew track
- WS-05 OSS boundary guardrails

### Phase 3: Rollout
- WS-06 docs, migration guidance, and release criteria.

## Feature-Level Deliverables

- Explicit install profiles for protocol, safety, and core usage
- MIT-compliant distribution path for `L0-L2`
- Enforced decoupling guardrails for public OSS scope
- Migration docs that let teams upgrade progressively
