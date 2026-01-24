---
ws_id: 00-191-01
project_id: 00
feature: F003
status: completed
size: MEDIUM
github_issue: null
assignee: null
started: 2026-01-11
completed: 2026-01-11
blocked_reason: null
---

## 02-191-01: Two-Stage Code Review

### 🎯 Goal

**What must WORK after this WS is complete:**
- `/codereview` skill has two stages: Spec Compliance → Code Quality
- Stage 1 (Spec) catches "well-written but wrong" bugs
- Stage 2 (Quality) only runs after Stage 1 passes
- Review loop: issue found → fix → re-review same stage

**Acceptance Criteria:**
- [x] AC1: `sdp/prompts/skills/two-stage-review.md` created
- [x] AC2: Stage 1 checklist (spec compliance, AC coverage, no over/under-engineering)
- [x] AC3: Stage 2 checklist (clean code, tests, security, maintainability)
- [x] AC4: `/codereview` skill updated to use two stages
- [x] AC5: Documentation with examples

---

### Context

From Superpowers: Two-stage review catches bugs that single-pass review misses.
- Stage 1 (Spec): Does code match specification exactly?
- Stage 2 (Quality): Is code well-written?

Key insight: Don't waste time perfecting wrong code.

---

### Dependencies

00--04 (Core package ready)

---

### Steps

1. Create two-stage-review.md skill prompt
2. Define Stage 1 (Spec Compliance) checklist
3. Define Stage 2 (Code Quality) checklist
4. Update /codereview to orchestrate stages
5. Add review loop logic (fail → fix → re-review)
6. Write documentation and examples

---

### Scope Estimate

- **Files:** 3 created/modified
- **Lines:** ~400 (prompts: ~300, integration: ~100)
- **Size:** MEDIUM

---

## Execution Report

**Date:** 2026-01-11  
**Status:** ✅ COMPLETED

### Summary

Implemented Two-Stage Code Review Protocol that separates spec compliance from code quality to catch "well-written but wrong" bugs.

### Files Created/Modified

1. **Created:** `sdp/prompts/skills/two-stage-review.md`
   - Two-stage review protocol with Stage 1 (Spec Compliance) and Stage 2 (Code Quality)
   - Review loop logic (fail → fix → re-review same stage)
   - Complete checklists for both stages
   - ~476 lines

2. **Modified:** `sdp/prompts/commands/codereview.md`
   - Updated to use Two-Stage Review Protocol
   - Integrated Stage 1 and Stage 2 checklists
   - Added review loop instructions
   - Updated output format to show both stages

3. **Modified:** `.claude/skills/codereview/SKILL.md`
   - Updated workflow description to reference Two-Stage Review
   - Added Stage 1 and Stage 2 execution instructions
   - Documented review loop logic

4. **Created:** `sdp/docs/two-stage-review.md`
   - Comprehensive documentation with examples
   - Three example scenarios (Stage 1 failure, Stage 2 failure, both pass)
   - Integration guide with `/codereview` command
   - ~400 lines

### Acceptance Criteria Verification

- ✅ **AC1:** `sdp/prompts/skills/two-stage-review.md` created
- ✅ **AC2:** Stage 1 checklist includes:
  - Goal Achievement (CRITICAL)
  - Specification Alignment
  - AC Coverage
  - No Over-Engineering
  - No Under-Engineering
- ✅ **AC3:** Stage 2 checklist includes:
  - Tests & Coverage
  - Regression
  - AI-Readiness
  - Clean Architecture
  - Type Hints
  - Error Handling
  - Security
  - No Tech Debt
  - Documentation
  - Git History
- ✅ **AC4:** `/codereview` skill updated:
  - `sdp/prompts/commands/codereview.md` uses two-stage protocol
  - `.claude/skills/codereview/SKILL.md` references two-stage review
- ✅ **AC5:** Documentation created:
  - `sdp/docs/two-stage-review.md` with 3 examples and integration guide

### Implementation Details

**Stage 1: Spec Compliance**
- Focuses on correctness: Does code match spec exactly?
- Checks Goal Achievement, Specification Alignment, AC Coverage
- Prevents over-engineering and under-engineering
- BLOCKING: If Stage 1 fails, Stage 2 doesn't run

**Stage 2: Code Quality**
- Focuses on quality: Is code well-written?
- Only runs if Stage 1 passes
- Checks tests, coverage, clean architecture, security, etc.

**Review Loop**
- If Stage 1 fails → Fix → Re-review Stage 1 only
- If Stage 2 fails → Fix → Re-review Stage 2 only (Stage 1 already passed)
- Prevents wasted effort polishing incorrect code

### Goal Achievement

✅ **All Goal requirements met:**
- `/codereview` skill has two stages: Spec Compliance → Code Quality ✅
- Stage 1 (Spec) catches "well-written but wrong" bugs ✅
- Stage 2 (Quality) only runs after Stage 1 passes ✅
- Review loop: issue found → fix → re-review same stage ✅

### Next Steps

1. Test the two-stage review with an actual feature review
2. Gather feedback on the protocol effectiveness
3. Consider adding telemetry/metrics for review stages

---

### Human Verification (UAT)

**Quick Smoke Test (30 sec):**
```bash
# Verify files exist
ls sdp/prompts/skills/two-stage-review.md
ls sdp/docs/two-stage-review.md

# Verify codereview references two-stage
grep -i "two-stage" sdp/prompts/commands/codereview.md | head -3
```

**Expected:** All files exist, codereview references two-stage review.

**Detailed Scenarios:**
1. Run `/codereview` on a feature and verify it uses two stages
2. Verify Stage 1 runs before Stage 2
3. Verify review loop works (fail → fix → re-review)

**Red Flags:**
- ❌ Stage 2 runs even if Stage 1 fails
- ❌ Review loop re-runs both stages instead of failed stage only
- ❌ Missing checklists in two-stage-review.md
