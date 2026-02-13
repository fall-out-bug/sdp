# F067: Repository Hardening and Alignment (Agents + Distribution + Development)

> **Status:** Draft (ready for @design)
> **Created:** 2026-02-13
> **Input:** reality-style repository audit + CI/release telemetry

---

## Problem Statement

The repository has strong foundations (multi-agent workflow, rich tests, quality intent), but execution is inconsistent across three critical planes:

1. **Agent/runtime configuration drift**
   - Duplicate prompt trees diverge (`prompts/agents` vs `sdp-plugin/prompts/agents`)
   - Cursor/OpenCode command references point to non-canonical or missing paths
2. **Distribution/release instability**
   - Release assets are missing or naming is inconsistent with installer expectations
   - Verify action fallback logic and release naming are not reliably aligned
3. **Development quality gate mismatch**
   - Declared quality bars (80% coverage, strict guard enforcement, <200 LOC) are not uniformly enforced in CI
   - Toolchain versions are inconsistent across root module, plugin module, CI, and Docker image

This creates avoidable failure loops: CI churn, onboarding friction, and non-deterministic agent behavior.

---

## Goals

1. **Single source of truth for agent behavior**
   - Canonical prompt definitions are unambiguous and verifiable in CI
2. **Reliable distribution path**
   - Release artifacts, naming, and install scripts are deterministic and testable
3. **Enforced quality gates**
   - Policy in docs/config/CI matches actual pass-fail behavior
4. **Stable developer ergonomics**
   - Local and CI toolchains are aligned and reproducible
5. **Zero-ambiguity skill execution**
   - Cursor/OpenCode adapters resolve to valid, tested skill paths

---

## Non-Goals

- Redesigning SDP protocol semantics
- Introducing new feature-level product capabilities
- Rewriting all historical docs/workstreams
- Enforcing large-scale code refactors in this feature (only target highest-risk mismatches)

---

## Scope Tracks

### Track A: Agents and Skill Configuration
- Canonicalization of prompt source tree
- Drift detection in CI
- Adapter path consistency for Claude/Cursor/OpenCode

### Track B: Distribution and Release
- Release artifact naming contract
- Verify action installer compatibility
- Release workflow health and artifact validation

### Track C: Development and Quality Enforcement
- Go toolchain version alignment
- Coverage/guard enforcement consistency
- Repository hygiene (tracked binaries/log policy)

---

## Skill-Aligned Execution Plan

| SDP Skill | Usage in F067 | Expected Output |
|-----------|----------------|-----------------|
| `@reality` | Baseline metrics and mismatch inventory | Audit snapshot + measurable baseline |
| `@feature` | Confirm final constraints and boundaries | Locked feature-level ACs |
| `@design` | Create dependency-ordered WS graph | WS files with scope + AC |
| `@build` | Execute each WS with TDD and gates | Committed incremental hardening |
| `@review` | Validate quality and policy alignment | PASS/CHANGES_REQUESTED with findings |
| `@deploy` | Merge after all gates and compatibility checks | Safe rollout to integration branch |

---

## Proposed Architecture (Operational)

```text
                    Canonical Source Layer
        prompts/*  +  docs/reference/*  +  .sdp/*
                              |
                              v
                     Validation and Enforcement
       CI checks (link/paths, drift, coverage, guard, release contract)
                              |
                              v
                    Delivery and Consumption Layer
   .claude adapters | .cursor commands | .opencode config | verify action
```

Core principle: **policy, docs, and CI behavior must be isomorphic** (no contradictory thresholds or fallback paths).

---

## Workstreams

| WS ID | Title | Priority | Depends On |
|-------|-------|----------|------------|
| 00-067-01 | Baseline Reality Snapshot and Target Metrics | P0 | - |
| 00-067-02 | Prompt Source-of-Truth and Drift Guard | P0 | 00-067-01 |
| 00-067-03 | Cursor/OpenCode Adapter Consistency | P0 | 00-067-02 |
| 00-067-04 | Go Toolchain Alignment (root/plugin/CI/Docker) | P0 | 00-067-01 |
| 00-067-05 | Quality Gate Enforcement Alignment | P0 | 00-067-04 |
| 00-067-06 | Release and Verify Action Contract Hardening | P0 | 00-067-04, 00-067-05 |
| 00-067-07 | Repository Hygiene (tracked binaries/evidence policy) | P1 | 00-067-01 |
| 00-067-08 | Documentation, Migration Notes, and Rollout | P1 | 00-067-02, 00-067-03, 00-067-05, 00-067-06, 00-067-07 |

---

## Feature-Level Acceptance Criteria

- [ ] AC1: Prompt source-of-truth is unique and enforced via CI
- [ ] AC2: No broken skill/agent references in Cursor/OpenCode adapters
- [ ] AC3: CI toolchain versions are aligned and documented
- [ ] AC4: Coverage threshold policy is consistent across docs/config/workflow
- [ ] AC5: Guard checks fail builds when violation severity requires blocking
- [ ] AC6: Release workflow publishes expected assets with deterministic names
- [ ] AC7: Verify action installer resolves and validates release assets without fallback ambiguity
- [ ] AC8: Tracked generated binaries are removed from source control policy
- [ ] AC9: Evidence log storage policy is explicit (tracked vs generated) and documented
- [ ] AC10: Go CI stability improves over baseline window
- [ ] AC11: `@review` for F067 returns no P0/P1 process-quality findings
- [ ] AC12: All new/updated docs reference canonical paths only

---

## Baseline Metrics (from audit snapshot)

| Metric | Baseline | Target |
|--------|----------|--------|
| Go CI success rate (last 20) | 2/20 | >= 16/20 |
| Non-test Go files >200 LOC | 38 | <= 20 (phase target) |
| Duplicate verify workflow files | 2 active | 1 canonical |
| Broken adapter references | Present | 0 |
| Tracked generated binaries | Present | 0 |

---

## Risks and Mitigations

1. **Risk:** Hardening changes break existing user automation  
   **Mitigation:** Add compatibility mapping + migration notes (WS-08)
2. **Risk:** Tightening gates increases short-term CI failures  
   **Mitigation:** Stage rollout (warn -> block) with explicit timeline
3. **Risk:** Release naming changes break consumers  
   **Mitigation:** Support transition aliases for one release cycle

---

## Rollout Strategy

1. Stabilize configuration and enforcement (WS-01..05)
2. Harden distribution contract (WS-06)
3. Clean repository artifacts and policy debt (WS-07)
4. Publish migration/runbook updates (WS-08)

---

## Success Definition

F067 is complete when:
- CI behavior matches declared SDP quality policy
- Agent adapters are deterministic and non-drifting
- Release/install path is reproducible end-to-end
- Repository no longer accumulates generated binaries as source artifacts

**Next step:** `@design idea-F067-repo-hardening-alignment`
