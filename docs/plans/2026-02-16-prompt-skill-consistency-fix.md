# Analysis: Enforcement of Existing Protocol

> **Status:** Revised after user feedback
> **Date:** 2026-02-16
> **Goal:** Minimal changes for enforcement of existing protocol

---

## Key Insight

**The protocol ALREADY EXISTS.** The problem is not lack of functionality, but **enforcement**.

### Existing Protocol

```
@oneshot F067  ->  @review F067  ->  @deploy F067
    |                 |                |
    v                 v                v
 Execute WS      APPROVED?         Merge PR
                  |
                  +- YES -> proceed
                  +-- NO -> fix loop
```

### What Was Violated

| Step | Protocol | What Happened |
|------|----------|---------------|
| 1 | @oneshot loads WS from backlog/ | OK |
| 2 | @review checks feature | **SKIPPED** |
| 3 | @deploy only after APPROVED | **SKIPPED** - merged PR directly |
| 4 | WS status update | Done manually, not per protocol |

---

## Root Cause

### Why the Protocol Did Not Work

1. **Context Loss** - After compaction, forgot about roadmap
2. **No enforcement gate** - Nothing blocked merge without @review
3. **Skill invocation optional** - Could skip @review and nothing broke
4. **Done = PR merged** - Habit of considering "done" by PR, not by verdict

### Why Proposed Solutions Were Wrong

| Proposal | Why Wrong |
|----------|-----------|
| New @milestone skill | Protocol already exists |
| .sdp/milestones.json | Overkill |
| sdp status reconcile | WS manage should be in protocol |
| sdp guard feature-complete | @review already does this |

---

## Minimal Fixes

### 1. Context Preservation (CLAUDE.md)

**Problem:** After compaction, milestone context is lost.

**Solution:** Add section to CLAUDE.md (already loads every session):

```markdown
## Milestone Context

Current milestone: **M1 "T-shirt"**

M1 Features: F054, F063, F064, F067, F068, F070, F075, F076
M2 Features: F060, F071, F073, F077, F078
M3 Features: F057, F058, F069, F072, F074, F079
M4 Features: F055, F056, F059, F061

**Warning:** Only work on current milestone features unless explicitly requested.
```

**Implementation:** 5 minutes, just add section.

### 2. Session-Start Pattern

**Problem:** New session does not check roadmap.

**Solution:** Create `.claude/patterns/session-start.md`:

```markdown
# Session Start Protocol

Before any work:
1. Read current milestone from CLAUDE.md
2. Check CHANGELOG for recent changes
3. Verify: Does the work belong to current milestone?
```

**Implementation:** 5 minutes, create file.

### 3. Review Gate (enforcement)

**Problem:** Nothing blocks PR merge without @review APPROVED.

**Solution A (CI):** Add GitHub Action:

```yaml
# .github/workflows/review-gate.yml
name: Review Gate
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  check-review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Check review verdict
        run: |
          if [ -f .sdp/review_verdict.json ]; then
            verdict=$(jq -r '.verdict' .sdp/review_verdict.json)
            if [ "$verdict" != "APPROVED" ]; then
              echo "ERROR: Feature not reviewed. Run @review first."
              exit 1
            fi
          fi
```

**Solution B (Skill update):** Update @deploy skill:

```markdown
## Pre-merge Check

Before merging:
1. Check if .sdp/review_verdict.json exists
2. If not: "Run @review first"
3. If exists but not APPROVED: "Fix issues first"
4. Only proceed if APPROVED
```

**Implementation:** 15-30 minutes.

### 4. @review Output (Evidence)

**Problem:** @review verdict is not persisted.

**Solution:** @review should create `.sdp/review_verdict.json`:

```json
{
  "feature": "F067",
  "verdict": "APPROVED",
  "timestamp": "2026-02-16T10:00:00Z",
  "reviewers": ["qa", "security", "devops", "sre", "techlead", "docs"],
  "summary": "All checks passed"
}
```

**Implementation:** Update @review skill, add file saving.

---

## Implementation Plan

### Phase 1: Quick Wins (30 minutes)

- [x] ~~Analysis document~~
- [x] Add Milestone Context to CLAUDE.md
- [x] Create session-start.md pattern

### Phase 2: Enforcement (1 hour)

- [x] Update @review skill: save verdict to .sdp/review_verdict.json
- [x] Update @deploy skill: check verdict before merge
- [ ] Add review-gate.yml GitHub Action

### Phase 3: Verification

- [ ] Test full flow: @oneshot -> @review -> @deploy
- [ ] Verify @deploy is blocked without APPROVED

---

## What We Are NOT Doing

| Rejected | Why |
|----------|-----|
| @milestone skill | Overkill, CLAUDE.md is sufficient |
| .sdp/milestones.json | Duplicates ROADMAP |
| sdp status reconcile | WS manage in protocol |
| sdp guard feature-complete | @review already does this |
| milestone field in Beads | Not needed for enforcement |

---

## Success Criteria

| Metric | Before | After |
|--------|--------|-------|
| PR without @review | Possible | Blocked |
| Milestone check | Manual | In CLAUDE.md |
| Session context | Lost | Restored |
| Review evidence | None | .sdp/review_verdict.json |

---

## Summary

**Problem:** Protocol existed, but was not enforced.

**Solution:** Minimal changes:
1. Context in CLAUDE.md
2. Session-start pattern
3. Review verdict gate

**Without:** New skills, new commands, new code.

---

*Ready to implement minimal fixes.*
