# F067: Repository Hardening and Alignment — Workstream Summary

> Feature: F067 | Priority: P0 | 8 workstreams

## Goal

Align SDP repository behavior across agent configs, distribution pipelines, and development quality gates so policy equals implementation.

## Dependency Graph

```text
00-067-01 Baseline Reality Snapshot
   ├─► 00-067-02 Prompt Source-of-Truth and Drift Guard
   │      └─► 00-067-03 Cursor/OpenCode Adapter Consistency
   │
   ├─► 00-067-04 Go Toolchain Alignment
   │      └─► 00-067-05 Quality Gate Enforcement Alignment
   │              └─► 00-067-06 Release + Verify Contract Hardening
   │
   └─► 00-067-07 Repository Hygiene

00-067-08 Documentation + Migration + Rollout
    depends on: 02,03,05,06,07
```

## Workstream Registry

| WS ID | Title | Est. | Depends On |
|-------|-------|------|------------|
| 00-067-01 | Baseline Reality Snapshot and Target Metrics | 1.5h | - |
| 00-067-02 | Prompt Source-of-Truth and Drift Guard | 2h | 00-067-01 |
| 00-067-03 | Cursor/OpenCode Adapter Consistency | 1.5h | 00-067-02 |
| 00-067-04 | Go Toolchain Alignment (root/plugin/CI/Docker) | 1.5h | 00-067-01 |
| 00-067-05 | Quality Gate Enforcement Alignment | 2h | 00-067-04 |
| 00-067-06 | Release and Verify Action Contract Hardening | 2.5h | 00-067-04,00-067-05 |
| 00-067-07 | Repository Hygiene (tracked binaries/evidence policy) | 1.5h | 00-067-01 |
| 00-067-08 | Documentation, Migration Notes, and Rollout | 1.5h | 00-067-02,00-067-03,00-067-05,00-067-06,00-067-07 |

## Feature-Level Deliverables

- Deterministic and validated agent adapter paths
- Consistent Go toolchain and quality thresholds
- Enforced guard/coverage behavior in CI
- Reliable release/install artifact contract
- Explicit repository policy for generated binaries and evidence logs
