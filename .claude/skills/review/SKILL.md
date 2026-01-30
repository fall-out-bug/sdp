---
name: review
description: Quality review with traceability check
tools: Read, Shell, Grep
---

# @review - Quality Review

Review feature by validating workstreams against quality gates and traceability.

## Invocation (BEADS-001)

Accepts **both** formats:

- `@review F01` — Feature ID (markdown workflow)
- `@review sdp-xxx` — Beads task ID (parent feature)

## Quick Reference

| Step | Action | Gate |
|------|--------|------|
| 0 | Resolve workstreams | beads_id → bd list, or markdown ls |
| 1 | List WS | All WS found |
| 2 | Traceability | All ACs have tests |
| 3 | Quality gates | All checks pass |
| 4 | Goal check | All ACs achieved |
| 5 | Verdict | APPROVED or CHANGES_REQUESTED |
| 6 | Post-review (if CHANGES_REQUESTED) | Report + @issue for bugs, WS under same feature |

## Workflow

### Step 0: Resolve Workstreams (when Beads enabled)

**Beads workflow** (bd installed, `.beads/` exists):
```bash
# Get sub-tasks (workstreams) under feature
bd list --parent {feature-id} --json
# Resolve beads_id → ws_id via .beads-sdp-mapping.jsonl for trace check
```

**Markdown workflow:**
```bash
ls docs/workstreams/completed/{feature-id}-*.md
```

### Step 1: List Workstreams

```bash
# Beads: bd list --parent {beads_id}
bd list --parent {feature-id}
```

Or for markdown workflow:

```bash
ls docs/workstreams/completed/{feature-id}-*.md
```

### Step 2: Check Traceability

For each workstream, verify all ACs have mapped tests using the traceability CLI:

```bash
# ws_id from mapping (beads_id → sdp_id) or from markdown filename
sdp trace check {WS-ID}
```

The command will:
- Extract all ACs from the workstream
- Check for test mappings
- Display traceability table
- Exit 1 if any AC is unmapped

Example output:

```
Traceability Report: 00-032-01
==================================================
| AC | Description | Test | Status |
|----|-------------|------|--------|
| AC1 | User can login | `test_user_login` | ✅ |
| AC2 | User can logout | - | ❌ |

Coverage: 50% (1/2 ACs mapped)
Status: ❌ INCOMPLETE (1 unmapped)
```

**Gate:** All ACs must have mapped tests (100% coverage).

If traceability check fails (exit code 1) → **CHANGES_REQUESTED**

**Auto-detection:** If mappings are missing, try auto-detection first:

```bash
sdp trace auto {WS-ID} --apply
```

This will automatically detect mappings from:
- Test docstrings (e.g., `"""Tests AC1"""`)
- Test function names (e.g., `test_ac1_user_login`)
- Keyword matching between AC descriptions and test names

### Step 3: Quality Gates

```bash
# All tests pass
pytest tests/ -v

# Coverage ≥80%
pytest --cov=src --cov-fail-under=80

# Type checking
mypy src/ --strict

# Linting
ruff check src/

# No except:pass
grep -r "except:" src/ | grep "pass"

# Files <200 LOC
find src/ -name "*.py" -exec sh -c 'lines=$(wc -l < "$1"); [ $lines -gt 200 ] && echo "$1: $lines lines"' _ {} \;
```

### Step 4: Goal Achievement

For each WS, verify:
- [ ] All ACs have passing tests
- [ ] Implementation matches description
- [ ] No TODO/FIXME in code

### Step 5: Verdict

**APPROVED** if:
- All ACs traceable to tests
- All tests pass
- All quality gates pass

**CHANGES_REQUESTED** if any fails.

No middle ground. No "approved with notes."

### Step 6: Post-Review Actions (when CHANGES_REQUESTED)

**6.1 Record verdict**
- Save report to `docs/reports/{YYYY-MM-DD}-{reviewed-id}-review.md`
- Include: verdict, AC status, quality gates, required actions

**6.2 Update reviewed item**
- Add to frontmatter: `review_verdict: CHANGES_REQUESTED`, `review_report: ../../reports/{date}-{id}-review.md`
- Add link to report in body

**6.3 Route findings — do NOT create new feature**

| Finding type | Action | Output |
|--------------|--------|--------|
| **Bugs** (failing tests, mypy/ruff, runtime errors) | @issue | `docs/issues/{ID}-{slug}.md` → route to /bugfix |
| **Planned work** (missing AC, new tests) | Add WS to **same feature** | `docs/workstreams/backlog/` with existing feature ID |
| **Pre-existing tech debt** | @issue for triage | docs/issues/ or backlog |

**6.4 Feature ID rule**
- **Never create new feature** for review follow-up
- Use `feature:` from reviewed workstreams or epic's parent
- Epic (e.g. BEADS-001) → use parent feature (e.g. F032), not F033

**6.5 Issue vs Workstream**
- Failing tests, errors → **Bug** → @issue → /bugfix
- Missing tests (AC), new capability → **Planned** → WS under same feature

## Quality Gates

See [Quality Gates Reference](../../docs/reference/quality-gates.md)

## Errors

| Error | Cause | Fix |
|-------|-------|-----|
| Missing trace | AC has no test | Add test for AC |
| Coverage <80% | Insufficient tests | Add more tests |
| Goal not met | AC not working | Fix implementation |

## See Also

- [Post-Review Fix Plan](../../docs/plans/2026-01-30-review-skill-post-review-fix.md) — Issue vs WS, feature ID rule
- [@issue skill](../issue/SKILL.md) — Bugs → docs/issues/ → /bugfix
- [Full Review Spec](../../docs/reference/review-spec.md)
- [Traceability Guide](../../docs/reference/traceability.md)
