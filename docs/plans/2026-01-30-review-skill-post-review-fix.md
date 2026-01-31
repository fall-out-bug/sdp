# Review Skill: Post-Review Actions Fix

**Date:** 2026-01-30  
**Problem:** Agent creates new feature (F033), workstreams instead of issues, doesn't record review properly.

---

## Root Cause Analysis

### What the agent did wrong (BEADS-001 review)

1. **Created new feature F033** — Follow-up work should stay under reviewed feature/epic
2. **Created workstreams** — Failing tests = bugs → should be @issue first, then @bugfix
3. **No issues created** — Bugs go to `docs/issues/` via @issue
4. **No clear recording** — Review skill has no Step 6: Post-Review Actions

### What the review skill lacks

| Gap | Current | Required |
|-----|---------|----------|
| Post-verdict actions | None | Step 6: Record + Route |
| Issue vs WS decision | Not documented | Bugs → @issue, Planned → WS |
| Feature ID rule | Not documented | Use reviewed feature, never create new |
| Report location | Not specified | docs/reports/{date}-{id}-review.md |
| Epic/feature update | Not specified | Add verdict + report link to frontmatter |

---

## Proposed Fix: Add Step 6 to Review Skill

### Step 6: Post-Review Actions (when CHANGES_REQUESTED)

**6.1 Record verdict**
- Save report to `docs/reports/{YYYY-MM-DD}-{reviewed-id}-review.md`
- Use reviewed ID: feature (F01), epic (BEADS-001), or beads_id

**6.2 Update reviewed item**
- Add to frontmatter: `review_verdict: CHANGES_REQUESTED`, `review_report: path`
- Add to body: link to report

**6.3 Route findings — CRITICAL**

| Finding type | Action | Output | Do NOT |
|--------------|--------|--------|--------|
| **Bugs** (failing tests, mypy/ruff in reviewed code, runtime errors) | @issue | docs/issues/{ID}-{slug}.md | Create WS |
| **Planned work** (missing AC, new tests for epic) | Add WS to **same feature** | docs/workstreams/backlog/ | Create new feature |
| **Pre-existing tech debt** (project-wide) | @issue for triage | docs/issues/ or backlog | Create new feature |

**6.4 Feature ID rule**
- **Never create new feature** for review follow-up
- Use `feature:` from reviewed workstreams or epic's parent feature
- If Epic (BEADS-001): use parent feature (e.g. F032) or epic's existing feature ref

**6.5 Issue vs Workstream**

```
Failing tests?        → Bug → @issue → /bugfix
Runtime error?        → Bug → @issue → /hotfix or /bugfix
Missing tests (AC)?   → Planned work → WS under same feature
New capability?       → Planned work → WS under same feature
Pre-existing debt?    → @issue → triage → backlog or WS in existing feature
```

---

## Implementation Plan

### WS-1: Update review skill (SKILL.md)

Add Step 6 section with:
- 6.1 Record verdict (report path)
- 6.2 Update reviewed item (frontmatter)
- 6.3 Route findings (table: Bug→@issue, Planned→WS)
- 6.4 Feature ID rule (never create new)
- 6.5 Issue vs WS decision tree

### WS-2: Update review-spec.md

Add "Post-Review Actions" section with same content for reference.

### WS-3: Create docs/issues/ + template

Ensure docs/issues/ exists, add issue template if needed.

---

## Rollback for F033

- Move 00-033-01, 00-033-02 to use feature: F032 (or BEADS-001 parent)
- Or: mark F033 as "review follow-up convention" — but user said no new feature
- Recommended: Change 00-033-* feature to F032, delete 00-033-00-feature-overview or merge into F032
