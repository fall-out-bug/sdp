# F072: Interop and Migration Path — Workstream Summary

> Feature: F072 | Priority: P1 | 5 workstreams
> Theme: import/export and low-friction migration from adjacent tools

## Goal

Reduce migration cost by supporting structured interop with major open-source spec/task ecosystems.

## Dependency Graph

```text
00-072-01 Interop Target Matrix + Contract
   ├─► 00-072-02 Import Pipeline (OpenSpec/Taskmaster/CCPM)
   ├─► 00-072-03 Export Pipeline (SDP -> external)
   └─► 00-072-04 Migration Wizard + Dry-Run Validator

00-072-05 Compatibility Docs + Lifecycle Policy
    depends on: 02, 03, 04
```

## Workstream Registry

| WS ID | Title | Priority | Size | Depends On |
|-------|-------|----------|------|------------|
| 00-072-01 | Interop Matrix and Schema Contracts | P1 | SMALL | — |
| 00-072-02 | Import Pipeline for External Specs/Tasks | P1 | LARGE | 01 |
| 00-072-03 | Export Pipeline from SDP Protocol | P2 | MEDIUM | 01 |
| 00-072-04 | Migration Wizard and Dry-Run Validation | P1 | MEDIUM | 02,03 |
| 00-072-05 | Compatibility Documentation and Version Policy | P1 | SMALL | 02,03,04 |

## Execution Phases

### Phase 1: Contract
- WS-01 defines supported formats, fidelity levels, and compatibility promises.

### Phase 2: Conversion Engine
- WS-02 import
- WS-03 export
- WS-04 guided migration

### Phase 3: Productization
- WS-05 policy and lifecycle commitments.

## Feature-Level Deliverables

- Explicit compatibility targets and support boundaries
- Import/export path with fidelity guarantees
- Migration dry-run that surfaces losses before writes
- Versioned compatibility policy for operators
