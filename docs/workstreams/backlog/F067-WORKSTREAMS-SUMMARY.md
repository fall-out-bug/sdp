# F067: Repository Hardening and Alignment — Workstream Summary

> Feature: F067 | Priority: P0 | 14 workstreams (8 original + 6 supplementary)
> Source: Two independent audits — PR #32 + cloud agent analysis

## Goal

Align SDP repository behavior across agent configs, distribution pipelines, development quality gates, and codebase structure so policy equals implementation.

## Dependency Graph

```text
00-067-01 Baseline Reality Snapshot
   │
   ├─► 00-067-02 Prompt Source-of-Truth and Drift Guard
   │      ├─► 00-067-03 Cursor/OpenCode Adapter Consistency
   │      └─► 00-067-13 Skill/Agent Optimization (size, versioning)
   │
   ├─► 00-067-04 Go Toolchain Alignment
   │      ├─► 00-067-05 Quality Gate Enforcement
   │      │      └─► 00-067-06 Release Contract Hardening
   │      ├─► 00-067-10 Go Module Consolidation (go.work)
   │      └─► 00-067-12 Developer Onboarding + Dependabot
   │
   ├─► 00-067-07 Repository Hygiene (binaries/evidence)
   │
   ├─► 00-067-09 LOC Violation Remediation
   │      └─► 00-067-14 Linter Exclusion Reduction
   │
   └─► 00-067-11 Git Structural Hygiene (submodule/Python)

00-067-08 Documentation + Migration + Rollout
    depends on: 02, 03, 05, 06, 07, 13
```

## Workstream Registry

| WS ID | Title | Priority | Size | Depends On | Source |
|-------|-------|----------|------|------------|--------|
| 00-067-01 | Baseline Reality Snapshot and Target Metrics | P0 | SMALL | — | PR #32 |
| 00-067-02 | Prompt Source-of-Truth and Drift Guard | P0 | MEDIUM | 01 | PR #32 |
| 00-067-03 | Cursor/OpenCode Adapter Consistency | P0 | SMALL | 02 | PR #32 |
| 00-067-04 | Go Toolchain Alignment (root/plugin/CI/Docker) | P0 | SMALL | 01 | PR #32 |
| 00-067-05 | Quality Gate Enforcement Alignment | P0 | MEDIUM | 04 | PR #32 |
| 00-067-06 | Release and Verify Action Contract Hardening | P0 | MEDIUM | 04, 05 | PR #32 |
| 00-067-07 | Repository Hygiene (binaries/evidence policy) | P1 | SMALL | 01 | PR #32 |
| 00-067-08 | Documentation, Migration Notes, and Rollout | P1 | MEDIUM | 02,03,05,06,07,13 | PR #32 |
| 00-067-09 | LOC Violation Remediation (split oversized files) | P0 | LARGE | 01 | Agent |
| 00-067-10 | Go Module Consolidation (dual source + go.work) | P1 | MEDIUM | 04 | Agent |
| 00-067-11 | Git Structural Hygiene (submodule/Python/paths) | P1 | SMALL | 01 | Agent |
| 00-067-12 | Developer Onboarding and Dependabot Setup | P1 | SMALL | 04, 05 | Agent |
| 00-067-13 | Skill and Agent Optimization (size, versioning) | P2 | MEDIUM | 02 | Agent |
| 00-067-14 | Linter Exclusion Reduction (errcheck, gocognit) | P1 | LARGE | 09 | Agent |

## Execution Phases

### Phase 1: Baseline (WS-01)
Freeze metrics before any changes.

### Phase 2: Core Alignment (WS-02, 03, 04, 05 — parallel tracks)
- **Track A** (agents): 02 → 03
- **Track B** (toolchain): 04 → 05

### Phase 3: Structural Remediation (WS-09, 10, 11 — parallel with Phase 2)
- Split oversized files (09)
- Consolidate modules (10)
- Clean Git structure (11)

### Phase 4: Distribution + Hardening (WS-06, 14 — after Phase 2+3)
- Harden release contract (06)
- Reduce linter exclusions (14, depends on file splits from 09)

### Phase 5: Rollout (WS-07, 08, 12, 13 — final)
- Repository hygiene (07)
- Documentation and migration (08)
- Developer onboarding (12)
- Skill optimization (13)

## Coverage: PR #32 vs Agent Analysis

| Finding | PR #32 WS | Agent WS | Status |
|---------|-----------|----------|--------|
| Duplicate agent definitions | 02 | — | Covered |
| OpenCode broken paths | 03 | — | Covered |
| Go version chaos | 04 | — | Covered |
| Coverage threshold mismatch | 05 | — | Covered |
| Release gap (v0.7.1 → v0.10.0) | 06 | — | Covered |
| Tracked binaries | 07 | — | Covered |
| Migration docs | 08 | — | Covered |
| **LOC violations (20+ files)** | — | **09** | **NEW** |
| **Dual source tree / go.work** | — | **10** | **NEW** |
| **Self-referential submodule** | — | **11** | **NEW** |
| **Python artifacts** | — | **11** | **NEW** |
| **Missing DEVELOPMENT.md** | — | **12** | **NEW** |
| **Missing Go Dependabot** | — | **12** | **NEW** |
| **Large skill files (700 LOC)** | — | **13** | **NEW** |
| **Agent versioning gaps** | — | **13** | **NEW** |
| **Linter exclusion sprawl** | — | **14** | **NEW** |

## Feature-Level Deliverables

- Deterministic and validated agent adapter paths
- Consistent Go toolchain and quality thresholds
- Enforced guard/coverage behavior in CI
- Reliable release/install artifact contract
- Explicit repository policy for generated binaries and evidence logs
- Own code passes own quality gates (LOC <=200, coverage >=80%)
- Linter exclusions reduced by >=50%
- go.work for multi-module development
- DEVELOPMENT.md for <15 minute contributor onboarding
- All agent definitions versioned and indexed
- No structural Git debt (submodule, Python artifacts)
