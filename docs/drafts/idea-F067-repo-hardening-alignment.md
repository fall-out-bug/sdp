# F067: Repository Hardening and Alignment

> **Status:** Draft (ready for @design)
> **Created:** 2026-02-13
> **Input:** Two independent repository audits (PR #32 + cloud agent analysis)
> **Workstreams:** 14 (8 original + 6 supplementary)

---

## Problem Statement

The repository has strong foundations (multi-agent workflow, rich tests, quality intent), but execution is inconsistent across five critical planes:

1. **Agent/runtime configuration drift**
   - Duplicate prompt trees diverge (`prompts/agents` vs `sdp-plugin/prompts/agents`)
   - Cursor/OpenCode command references point to non-canonical or missing paths
2. **Distribution/release instability**
   - Release assets are missing or naming is inconsistent with installer expectations
   - Latest GitHub release is v0.7.1 while CHANGELOG documents v0.10.0
   - `go install` path in docs points to root module instead of `sdp-plugin/cmd/sdp/`
3. **Development quality gate mismatch**
   - CI enforces 60% coverage but documentation requires 80%
   - `.golangci.yml` has 30+ exclusion rules disabling `errcheck` for most packages
   - 20+ source files exceed the declared 200 LOC limit (up to 572 LOC)
4. **Toolchain fragmentation**
   - Five different Go versions across root `go.mod` (1.25.6), `sdp-plugin/go.mod` (1.26), Dockerfile (1.21), README badge (1.24+), CI (1.26)
   - Dual Go module trees (`src/sdp/` vs `sdp-plugin/internal/`) with no `go.work`
5. **Structural debt**
   - Self-referential Git submodule (`.sdp/.sdp` → same repo)
   - Python-era configuration artifacts remain (`poetry.lock`, `quality-gate.toml` with pytest/mypy refs)
   - No contributor development setup documentation
   - Large skill files (700+ lines) consume excessive context window

---

## Goals

1. **Single source of truth for agent behavior** — canonical prompts are unambiguous and CI-enforced
2. **Reliable distribution path** — release artifacts, naming, and install scripts are deterministic
3. **Enforced quality gates** — policy in docs/config/CI matches actual behavior
4. **Stable developer ergonomics** — local and CI toolchains aligned, contributor onboarding smooth
5. **Zero-ambiguity skill execution** — Cursor/OpenCode adapters resolve to valid paths
6. **Practice what we preach** — own code passes own quality gates (LOC, coverage, linting)

---

## Non-Goals

- Redesigning SDP protocol semantics
- Introducing new feature-level product capabilities
- Rewriting all historical docs/workstreams
- Full elimination of all LOC violations (phased reduction target)

---

## Scope Tracks

### Track A: Agents and Skill Configuration (WS-02, WS-03, WS-13)
- Canonicalization of prompt source tree
- Drift detection in CI
- Adapter path consistency for Claude/Cursor/OpenCode
- Skill/agent versioning and size optimization

### Track B: Distribution and Release (WS-06)
- Release artifact naming contract
- Verify action installer compatibility
- Release workflow health and artifact validation

### Track C: Development and Quality Enforcement (WS-04, WS-05, WS-12, WS-14)
- Go toolchain version alignment
- Coverage/guard enforcement consistency
- Linter exclusion reduction
- Developer onboarding documentation

### Track D: Structural Remediation (WS-09, WS-10, WS-11)
- LOC violation remediation for highest-risk files
- Go module consolidation with go.work
- Git submodule cleanup and install path fixes

### Track E: Baseline and Rollout (WS-01, WS-07, WS-08)
- Baseline reality snapshot
- Repository hygiene
- Migration documentation

---

## Skill-Aligned Execution Plan

| SDP Skill | Usage in F067 | Expected Output |
|-----------|---------------|-----------------|
| `@reality` | Baseline metrics and mismatch inventory | Audit snapshot + measurable baseline |
| `@feature` | Confirm final constraints and boundaries | Locked feature-level ACs |
| `@design` | Create dependency-ordered WS graph | WS files with scope + AC |
| `@build` | Execute each WS with TDD and gates | Committed incremental hardening |
| `@review` | Validate quality and policy alignment | PASS/CHANGES_REQUESTED with findings |
| `@deploy` | Merge after all gates and checks | Safe rollout to integration branch |

---

## Workstreams

### Original (PR #32 scope)

| WS ID | Title | Priority | Depends On |
|-------|-------|----------|------------|
| 00-067-01 | Baseline Reality Snapshot and Target Metrics | P0 | — |
| 00-067-02 | Prompt Source-of-Truth and Drift Guard | P0 | 01 |
| 00-067-03 | Cursor/OpenCode Adapter Consistency | P0 | 02 |
| 00-067-04 | Go Toolchain Alignment (root/plugin/CI/Docker) | P0 | 01 |
| 00-067-05 | Quality Gate Enforcement Alignment | P0 | 04 |
| 00-067-06 | Release and Verify Action Contract Hardening | P0 | 04, 05 |
| 00-067-07 | Repository Hygiene (tracked binaries/evidence policy) | P1 | 01 |
| 00-067-08 | Documentation, Migration Notes, and Rollout | P1 | 02, 03, 05, 06, 07, 13 |

### Supplementary (cloud agent analysis)

| WS ID | Title | Priority | Depends On |
|-------|-------|----------|------------|
| 00-067-09 | LOC Violation Remediation (split oversized files) | P0 | 01 |
| 00-067-10 | Go Module Consolidation (dual source tree + go.work) | P1 | 04 |
| 00-067-11 | Git Structural Hygiene (submodule, install path, Python artifacts) | P1 | 01 |
| 00-067-12 | Developer Onboarding and Dependabot Setup | P1 | 04, 05 |
| 00-067-13 | Skill and Agent Optimization (size, versioning) | P2 | 02 |
| 00-067-14 | Linter Exclusion Reduction (errcheck, gocognit) | P1 | 09 |

---

## Feature-Level Acceptance Criteria

- [ ] AC1: Prompt source-of-truth is unique and enforced via CI
- [ ] AC2: No broken skill/agent references in Cursor/OpenCode adapters
- [ ] AC3: CI toolchain versions are aligned and documented
- [ ] AC4: Coverage threshold policy is consistent across docs/config/workflow
- [ ] AC5: Guard checks fail builds when violation severity requires blocking
- [ ] AC6: Release workflow publishes expected assets with deterministic names
- [ ] AC7: Verify action installer resolves artifacts without fallback ambiguity
- [ ] AC8: Tracked generated binaries are removed from source control
- [ ] AC9: Evidence log storage policy is documented
- [ ] AC10: Go CI stability improves over baseline (target >=16/20)
- [ ] AC11: `@review` for F067 returns no P0/P1 process-quality findings
- [ ] AC12: All docs reference canonical paths only
- [ ] AC13: No non-test Go source file exceeds 200 LOC (phased: <=20 violations)
- [ ] AC14: Go modules use go.work or are consolidated into single tree
- [ ] AC15: Self-referential submodule removed
- [ ] AC16: DEVELOPMENT.md exists with working local setup instructions
- [ ] AC17: All agent definitions have version headers
- [ ] AC18: Linter exclusions reduced by >=50%

---

## Baseline Metrics (from combined audit)

| Metric | Baseline | Phase 1 Target | Final Target |
|--------|----------|----------------|--------------|
| Go CI success rate (last 20) | 2/20 | >= 12/20 | >= 16/20 |
| Non-test Go files >200 LOC | 20+ | <= 10 | 0 |
| Duplicate agent definition sets | 2 | 1 | 1 |
| Broken adapter references | Present | 0 | 0 |
| Tracked generated binaries | Present | 0 | 0 |
| CI coverage threshold | 60% | 80% | 80% |
| Linter exclusion rules | 30+ | <= 15 | <= 5 |
| Go version declarations | 5 different | 1 | 1 |

---

## Risks and Mitigations

1. **Risk:** Hardening changes break existing user automation
   **Mitigation:** Migration notes + compatibility mapping (WS-08)
2. **Risk:** Tightening gates increases short-term CI failures
   **Mitigation:** Stage rollout (warn → block) with explicit timeline
3. **Risk:** LOC splitting introduces regressions
   **Mitigation:** TDD approach — tests must pass after each split (WS-09)
4. **Risk:** Module consolidation breaks import paths
   **Mitigation:** go.work as intermediate step before full merge (WS-10)

---

## Dependency Graph

```text
00-067-01 Baseline Reality Snapshot
   ├─► 00-067-02 Prompt Source-of-Truth
   │      ├─► 00-067-03 Cursor/OpenCode Adapters
   │      └─► 00-067-13 Skill/Agent Optimization
   │
   ├─► 00-067-04 Go Toolchain Alignment
   │      ├─► 00-067-05 Quality Gate Enforcement
   │      │      └─► 00-067-06 Release Contract Hardening
   │      ├─► 00-067-10 Go Module Consolidation
   │      └─► 00-067-12 Developer Onboarding + Dependabot
   │
   ├─► 00-067-07 Repository Hygiene
   │
   ├─► 00-067-09 LOC Violation Remediation
   │      └─► 00-067-14 Linter Exclusion Reduction
   │
   └─► 00-067-11 Git Structural Hygiene

00-067-08 Documentation + Migration + Rollout
    depends on: 02, 03, 05, 06, 07, 13
```

---

## Success Definition

F067 is complete when:
- CI behavior matches declared SDP quality policy
- Agent adapters are deterministic and non-drifting
- Release/install path is reproducible end-to-end
- Own code passes own quality gates (LOC, coverage, linting)
- Repository is clean of structural debt
- New contributors can set up locally in <15 minutes

**Next step:** `@design idea-F067-repo-hardening-alignment`
