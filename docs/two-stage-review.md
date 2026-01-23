# Two-Stage Code Review Protocol

**Version:** 1.0.0  
**Feature:** F191  
**Workstream:** WS-191-01

---

## Overview

The Two-Stage Code Review Protocol separates **spec compliance** from **code quality** to catch "well-written but wrong" bugs before polishing incorrect code.

**Key Insight:** Don't waste time perfecting wrong code. First verify correctness, then polish quality.

---

## Why Two Stages?

### Problem: Single-Pass Review Misses Bugs

Traditional single-pass review can miss critical issues:

1. **Well-written but wrong code** — Code is clean, tested, and follows best practices, but doesn't match the specification.
2. **Premature optimization** — Time spent polishing code that doesn't solve the right problem.
3. **Mixed concerns** — Spec issues and quality issues get conflated, making fixes harder.

### Solution: Sequential Stages

**Stage 1 (Spec Compliance):** Does the code match the specification exactly?  
**Stage 2 (Code Quality):** Is the code well-written?

Only Stage 2 runs if Stage 1 passes. This ensures we don't polish incorrect code.

---

## Review Flow

```
┌─────────────────┐
│  Stage 1: Spec  │
│  Compliance     │
└────────┬────────┘
         │
    ✅ Pass?
         │
    ┌────┴────┐
    │         │
   YES       NO
    │         │
    ▼         ▼
┌─────────┐  ┌──────────┐
│ Stage 2 │  │ Fix &    │
│ Quality │  │ Re-review│
└─────────┘  └──────────┘
```

---

## Stage 1: Spec Compliance

**Question:** Does the code match the specification exactly?

### Checklist

1. **Goal Achievement (CRITICAL)**
   - All Acceptance Criteria (AC) must pass
   - Target: 100% AC passed
   - If ANY AC fails → Stage 1 FAILED

2. **Specification Alignment**
   - All required features from spec are implemented
   - No missing functionality
   - No over-engineering (extra features not in spec)
   - No under-engineering (simplified beyond spec)

3. **AC Coverage**
   - Each AC has corresponding implementation
   - Each AC has verification (tests)
   - Tests pass for all AC

4. **No Over-Engineering**
   - No extra features not in spec
   - No overly complex patterns for simple requirements
   - No premature optimization
   - No unnecessary abstractions

5. **No Under-Engineering**
   - No missing required features
   - No simplified implementation beyond spec
   - No missing error handling from spec
   - No missing edge cases from spec

### Verdict

- **PASS:** All checks ✅ → Proceed to Stage 2
- **FAIL:** Any check 🔴 → CHANGES REQUESTED → Fix → Re-review Stage 1

---

## Stage 2: Code Quality

**Question:** Is the code well-written?

**Prerequisite:** Stage 1 must pass before Stage 2 runs.

### Checklist

1. **Tests & Coverage**
   - Target: ≥80% coverage
   - All tests pass

2. **Regression**
   - All existing tests still pass
   - No breaking changes

3. **AI-Readiness**
   - File size: <200 LOC
   - Cyclomatic complexity: <10

4. **Clean Architecture**
   - Domain doesn't import infrastructure
   - Domain doesn't import presentation
   - Dependencies point inward

5. **Type Hints**
   - 100% type hints
   - `mypy --strict` passes

6. **Error Handling**
   - No `except: pass`
   - No bare `except:`
   - Explicit error handling

7. **Security**
   - No SQL injection
   - No shell injection
   - `bandit` passes

8. **No Tech Debt**
   - No TODO/FIXME
   - No HACK/XXX markers

9. **Documentation**
   - Docstrings for public functions
   - README updated (if needed)

10. **Git History**
    - Commits exist for WS
    - Conventional commit format

### Verdict

- **PASS:** All checks ✅ → APPROVED
- **FAIL:** Any check 🔴 → CHANGES REQUESTED → Fix → Re-review Stage 2

---

## Review Loop Logic

### Flow

```
1. Run Stage 1
   ├─ PASS → Run Stage 2
   │         ├─ PASS → APPROVED
   │         └─ FAIL → Fix → Re-review Stage 2
   └─ FAIL → Fix → Re-review Stage 1
```

### Re-review Rules

1. **After fix:** Re-run the **failed stage only** (not both stages)
2. **Stage 1 fix:** Re-run Stage 1 → if PASS → proceed to Stage 2
3. **Stage 2 fix:** Re-run Stage 2 only (Stage 1 already passed)

### Example Workflow

```
Initial Review:
  Stage 1: FAIL (AC2 not working)
  → Fix AC2
  → Re-review Stage 1: PASS
  → Run Stage 2: FAIL (coverage 75%)
  → Fix coverage
  → Re-review Stage 2: PASS
  → APPROVED
```

---

## Examples

### Example 1: Stage 1 Failure (Spec Mismatch)

**Scenario:** Code is well-written but doesn't match spec.

**WS Spec:**
- AC1: Feature X returns list of items
- AC2: Feature X handles empty input gracefully

**Implementation:**
- Feature X returns dict (not list) ❌
- Feature X handles empty input ✅
- Code is clean, tested, 90% coverage ✅

**Review Result:**

```
Stage 1: Spec Compliance
  ✅ Goal Achievement: ❌ FAIL (AC1 fails - returns dict not list)
  ✅ Specification Alignment: ❌ FAIL (returns wrong type)
  ✅ AC Coverage: ⚠️ WARNING (AC1 implementation wrong)
  ✅ No Over-Engineering: ✅ PASS
  ✅ No Under-Engineering: ✅ PASS

Verdict: 🔴 FAIL

Stage 2: Code Quality
  (Not run - Stage 1 failed)

Overall Verdict: CHANGES REQUESTED
```

**Action:** Fix AC1 to return list → Re-review Stage 1 → if PASS → run Stage 2.

---

### Example 2: Stage 2 Failure (Quality Issue)

**Scenario:** Code matches spec but has quality issues.

**WS Spec:**
- AC1: Feature X returns list of items ✅
- AC2: Feature X handles empty input ✅

**Implementation:**
- Feature X returns list ✅
- Feature X handles empty input ✅
- Code has 75% coverage ❌
- Code has TODO markers ❌

**Review Result:**

```
Stage 1: Spec Compliance
  ✅ Goal Achievement: ✅ PASS (2/2 AC passed)
  ✅ Specification Alignment: ✅ PASS
  ✅ AC Coverage: ✅ PASS
  ✅ No Over-Engineering: ✅ PASS
  ✅ No Under-Engineering: ✅ PASS

Verdict: ✅ PASS

Stage 2: Code Quality
  ✅ Tests & Coverage: ❌ FAIL (75% < 80%)
  ✅ No Tech Debt: ❌ FAIL (TODO markers found)
  ✅ All other checks: ✅ PASS

Verdict: 🔴 FAIL

Overall Verdict: CHANGES REQUESTED
```

**Action:** Fix coverage and remove TODOs → Re-review Stage 2 only → if PASS → APPROVED.

---

### Example 3: Both Stages Pass

**Scenario:** Code matches spec and is well-written.

**Review Result:**

```
Stage 1: Spec Compliance
  ✅ Goal Achievement: ✅ PASS (2/2 AC passed)
  ✅ Specification Alignment: ✅ PASS
  ✅ AC Coverage: ✅ PASS
  ✅ No Over-Engineering: ✅ PASS
  ✅ No Under-Engineering: ✅ PASS

Verdict: ✅ PASS

Stage 2: Code Quality
  ✅ Tests & Coverage: ✅ PASS (85%)
  ✅ Regression: ✅ PASS
  ✅ AI-Readiness: ✅ PASS
  ✅ Clean Architecture: ✅ PASS
  ✅ Type Hints: ✅ PASS
  ✅ Error Handling: ✅ PASS
  ✅ Security: ✅ PASS
  ✅ No Tech Debt: ✅ PASS
  ✅ Documentation: ✅ PASS
  ✅ Git History: ✅ PASS

Verdict: ✅ PASS

Overall Verdict: ✅ APPROVED
```

---

## Integration with /codereview

The `/codereview` command uses the Two-Stage Review Protocol:

```bash
/codereview F60
```

**Process:**

1. Find all WS in feature F60
2. For each WS:
   - Run Stage 1 (Spec Compliance)
   - If PASS → Run Stage 2 (Code Quality)
   - If FAIL → Report issues → Fix → Re-review failed stage
3. Cross-WS checks
4. Generate UAT guide (if all WS approved)
5. Output verdict

**Output:**

Each WS file gets review results appended:

```markdown
---

### Review Results

**Date:** 2026-01-11
**Reviewer:** agent
**Verdict:** APPROVED / CHANGES REQUESTED

#### Stage 1: Spec Compliance
[Stage 1 results]

#### Stage 2: Code Quality
[Stage 2 results]

#### Issues (if CHANGES REQUESTED)
[Issue list]
```

---

## Benefits

1. **Catches "well-written but wrong" bugs** — Spec issues found before polishing
2. **Saves time** — Don't polish incorrect code
3. **Clear separation** — Spec issues vs. quality issues are distinct
4. **Focused fixes** — Fix spec issues first, then quality
5. **Review loop efficiency** — Re-review only failed stage

---

## Key Principles

1. **Stage 1 First:** Always check spec compliance before code quality
2. **No Wasted Effort:** Don't polish code that doesn't match spec
3. **Clear Separation:** Spec issues vs. quality issues are different
4. **Review Loop:** Fix → re-review same stage (not both)
5. **Zero Tolerance:** No "minor issues" — all blockers must be fixed

---

## Related Documents

- **Protocol:** `sdp/prompts/skills/two-stage-review.md`
- **Command Prompt:** `sdp/prompts/commands/codereview.md`
- **Skill:** `.claude/skills/codereview/SKILL.md`
- **Feature:** F191 (Two-Stage Review)
- **Workstream:** WS-191-01

---

## Changelog

### v1.0.0 (2026-01-11)

- Initial implementation
- Two-stage review protocol
- Review loop logic
- Integration with `/codereview`
