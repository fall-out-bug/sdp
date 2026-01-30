# Issue 002: F032 Review — Quality Gate Failures

> **Source:** [F032 Review 2026-01-30](../reports/2026-01-30-F032-review.md)
> **Severity:** P2 (MEDIUM)
> **Route:** @bugfix
> **Feature:** F032

**Agent command:** `@bugfix "F032 review quality gate failures" --feature=F032 --issue-id=002`

---

## Symptom

F032 review verdict: **CHANGES_REQUESTED**

1. **Traceability blocked** — `sdp trace check {WS-ID}` fails:
   - Real Beads: `KeyError: 'task_id'` in `sdp.beads.models` (BeadsTask.from_dict)
   - Mock Beads: `WS not found: {ws_id}`

2. **ValidationCheck undefined** — F821 in:
   - `capability_tier_checks_contract.py` (3×)
   - `capability_tier_checks_interface.py` (2×)

3. **Ruff violations** — F401 (unused imports), E501 (line length), C901 (complexity > 10)

---

## Routing

| Step | Skill | Action |
|------|-------|--------|
| 1 | @bugfix | Fix with TDD, branch from dev |
| 2 | @review | Re-review F032 after fix |

---

## Acceptance Criteria

- [x] `sdp trace check {WS-ID}` works (real or mock)
- [x] ValidationCheck imported/defined in capability_tier checks
- [x] F821/F401 fixed; C901/E501 in other modules deferred
- [x] No regressions (80+ tests pass for changed modules)

---

## Resolution (bugfix/002-f032-review-quality-gates)

**Date:** 2026-01-30

**Fixes:**
1. **Traceability KeyError** — BeadsTask.from_dict: handle `task_id` or `id` for dependencies
2. **Traceability WS not found** — Markdown fallback: read from docs/workstreams/ when Beads has no task
3. **Trace CLI** — Register `sdp trace` in main CLI
4. **ValidationCheck F821** — Add module-level imports in capability_tier_checks_contract, _interface, _scope, t0_t1, t2_t3
5. **E501** — Split long lines in health_checks/checks.py, beads.py
