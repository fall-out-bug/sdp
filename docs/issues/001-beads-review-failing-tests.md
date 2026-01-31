# Issue 001: BEADS-001 Review — 6 Failing Tests

> **Source:** [BEADS-001 Review 2026-01-30](../reports/2026-01-30-BEADS-001-review.md)
> **Severity:** P2 (MEDIUM)
> **Route:** @bugfix
> **Feature:** F032

**Agent command:** `@bugfix "6 failing tests from BEADS-001 review" --feature=F032 --issue-id=001`

---

## Symptom

6 tests fail when running `pytest tests/ -v`:

1. `test_beads_check_installed` — health_checks message format
2. `test_handles_beads_unavailable` — BeadsClientError not raised
3. `test_invalid_tier_case` — capability_tier ValueError
4. `test_nonexistent_file` — capability_tier file handling
5. `test_empty_body` — Workstream init missing status/size
6. `test_section_case_sensitivity` — Workstream init missing status/size

---

## Routing

| Step | Skill | Action |
|------|-------|--------|
| 1 | @bugfix | Fix with TDD, branch from develop |
| 2 | @review | Re-review after fix |

---

## Acceptance Criteria

- [x] All 6 tests pass
- [x] No regressions (6 tests fixed)
- [ ] Coverage ≥80% for changed modules
- [ ] mypy --strict passes

---

## Resolution (bugfix/001-beads-review-failing-tests)

**Date:** 2026-01-30

**Fixes:**
1. **test_invalid_tier_case** — `validate_workstream_tier`: require uppercase tier (reject "t0")
2. **test_nonexistent_file** — `validate_workstream_tier`: raise `WorkstreamParseError` when file not found
3. **test_empty_body** / **test_section_case_sensitivity** — Add `status=WorkstreamStatus.BACKLOG`, `size=WorkstreamSize.SMALL` to `Workstream()` in tests
4. **test_beads_check_installed** / **test_handles_beads_unavailable** — Fixed in WS 00-033-02 (already on dev)
5. **test_cli_beads_client** — Remove unused `BeadsNotInstalledError` import
