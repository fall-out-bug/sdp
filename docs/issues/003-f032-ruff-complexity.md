# Issue 003: F032 Review — Ruff Complexity Violations

> **Source:** [F032 Review 2026-01-30](../reports/2026-01-30-F032-review.md)
> **Severity:** P2 (MEDIUM)
> **Route:** @bugfix
> **Feature:** F032
> **Deferred from:** Issue 002 (AC: "C901/E501 in other modules deferred")

**Agent command:** `@bugfix "F032 ruff complexity violations" --feature=F032 --issue-id=003`

---

## Symptom

8 ruff violations remain after Issue 002 bugfix:

| File | Rule | Issue |
|------|------|-------|
| `cli/skill.py` | C901 | `validate_skill` complexity 14 > 10 |
| `cli/tier_promote.py` | C901 | `tier_promote_check` complexity 11 > 10 |
| `hooks/ws_complete.py` | C901 | `run` complexity 12 > 10 |
| `validators/ws_completion.py` | C901 | `_parse_ws_file` complexity 17 > 10 |
| `validators/ws_template_checker.py` | C901 | complexity 14 > 10 |
| `validators/ws_template_checker.py` | E501 | 2× line length |
| `validators/supersede_checker.py` | E501 | line length |

---

## Acceptance Criteria

- [x] All C901 violations fixed (complexity ≤ 10)
- [x] All E501 violations fixed (line length ≤ 88)
- [x] `ruff check src/sdp/` passes with 0 errors
- [x] No regressions (existing tests pass)

---

## Technical Approach

C901 fixes require extracting helper functions:

```python
# Before: one 17-branch function
def _parse_ws_file(path: Path) -> dict:
    # 17 branches...

# After: split into focused helpers
def _parse_ws_file(path: Path) -> dict:
    frontmatter = _extract_frontmatter(path)
    acs = _extract_acceptance_criteria(path)
    return _build_ws_dict(frontmatter, acs)
```

E501 fixes: split long strings, use variables.

---

## Resolution (bugfix/003-f032-ruff-complexity)

**Date:** 2026-01-30

**Fixes:**
1. **cli/skill.py** — Extracted `_check_line_count`, `_check_sections`, `_check_refs`, `_output_validation`
2. **cli/tier_promote.py** — Extracted `_get_next_tier`, `_show_validation_result`
3. **hooks/ws_complete.py** — Extracted `_handle_bypass`, `_handle_verification_passed`, `_handle_verification_failed`
4. **validators/ws_completion.py** — Extracted `_parse_frontmatter_scope`, `_parse_verification_commands`
5. **validators/ws_template_checker.py** — Extracted `_extract_ws_id`, `_check_required_sections`, `_check_code_blocks`; split E501 long lines
6. **validators/supersede_checker.py** — Split E501 long line (SupersedeResult)
