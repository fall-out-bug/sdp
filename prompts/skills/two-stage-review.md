# Two-Stage Code Review Protocol

**Purpose:** Catch "well-written but wrong" bugs by separating spec compliance from code quality.

**Key Insight:** Don't waste time perfecting wrong code. First verify correctness, then polish quality.

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

**Rule:** Stage 2 only runs if Stage 1 passes. If Stage 1 fails → fix → re-review Stage 1.

---

## Stage 1: Spec Compliance (BLOCKING)

**Question:** Does the code match the specification exactly?

**Goal:** Catch "well-written but wrong" bugs before polishing.

### Checklist

#### 1. Goal Achievement (CRITICAL)

```bash
# Read Goal from WS file
grep -A20 "### 🎯 Goal" WS-XXX-YY.md

# Check each Acceptance Criterion
# - AC1: ... → verify it works (✅/❌)
# - AC2: ... → verify it works (✅/❌)
```

**Metrics:**
- Target: 100% AC passed
- Actual: {X}/{Y} AC passed ({percentage}%)
- Status: ✅ / 🔴 BLOCKING

**If ANY AC ❌ → Stage 1 FAILED → CHANGES REQUESTED**

---

#### 2. Specification Alignment

**Check:** Does implementation match the spec exactly?

```bash
# Compare WS spec with implementation
# - Are all required features present?
# - Are any features missing?
# - Are any extra features added (over-engineering)?
```

**Questions:**
- [ ] All required features from spec are implemented?
- [ ] No missing functionality?
- [ ] No over-engineering (extra features not in spec)?
- [ ] No under-engineering (simplified beyond spec)?

**Status:** ✅ / 🔴 BLOCKING

---

#### 3. Acceptance Criteria Coverage

**Check:** Each AC has corresponding implementation and verification.

```bash
# For each AC in WS file:
# 1. Find corresponding code
# 2. Verify it works
# 3. Check if tests cover it
```

**Example:**
- AC1: "Feature X works" → Code exists? → Test exists? → Test passes?
- AC2: "Error handling for Y" → Code exists? → Test exists? → Test passes?

**Status:** ✅ / 🔴 BLOCKING

---

#### 4. No Over-Engineering

**Check:** Implementation doesn't add unnecessary complexity.

**Red Flags:**
- [ ] Extra features not in spec
- [ ] Overly complex patterns for simple requirements
- [ ] Premature optimization
- [ ] Unnecessary abstractions

**Status:** ✅ / ⚠️ WARNING / 🔴 BLOCKING

---

#### 5. No Under-Engineering

**Check:** Implementation doesn't skip required functionality.

**Red Flags:**
- [ ] Missing required features
- [ ] Simplified beyond spec requirements
- [ ] Missing error handling specified in spec
- [ ] Missing edge cases from spec

**Status:** ✅ / 🔴 BLOCKING

---

### Stage 1 Verdict

**PASS:** All checks ✅ → Proceed to Stage 2

**FAIL:** Any check 🔴 → CHANGES REQUESTED → Fix → Re-review Stage 1

**Output Format:**

```markdown
## Stage 1: Spec Compliance

**Date:** {YYYY-MM-DD}
**Reviewer:** {agent}

| Check | Status | Notes |
|-------|--------|-------|
| Goal Achievement | ✅ / 🔴 | {X}/{Y} AC passed |
| Specification Alignment | ✅ / 🔴 | {notes} |
| AC Coverage | ✅ / 🔴 | {coverage details} |
| No Over-Engineering | ✅ / ⚠️ / 🔴 | {notes} |
| No Under-Engineering | ✅ / 🔴 | {notes} |

**Verdict:** ✅ PASS / 🔴 FAIL

**Issues (if FAIL):**
1. {Issue 1}
2. {Issue 2}
```

---

## Stage 2: Code Quality (Only if Stage 1 Passes)

**Question:** Is the code well-written?

**Goal:** Ensure maintainability, security, and best practices.

### Checklist

#### 1. Tests & Coverage

```bash
pytest tests/unit/test_XXX.py --cov=hw_checker/module --cov-report=term-missing
```

**Metrics:**
- Target: ≥80% coverage
- Actual: {coverage}%
- Status: ✅ (≥80%) / ⚠️ (70-79%) / 🔴 BLOCKING (<70%)

---

#### 2. Regression Tests

```bash
pytest tests/unit/ -m fast -q --tb=short
```

**Status:** ✅ All pass / 🔴 Failures

---

#### 3. AI-Readiness

```bash
# File size
wc -l src/hw_checker/module/*.py

# Complexity
ruff check src/hw_checker/module/ --select=C901
```

**Metrics:**
- File Size Target: <200 LOC
- Actual: max {max_loc} LOC in {filename}
- Status: ✅ (all <200) / ⚠️ (200-250) / 🔴 BLOCKING (>250)

- Complexity Target: CC <10
- Actual: avg CC {avg_cc}, max CC {max_cc}
- Status: ✅ (<10) / ⚠️ (10-15) / 🔴 BLOCKING (>15)

---

#### 4. Clean Architecture

```bash
# Domain doesn't import infrastructure
grep -r "from hw_checker.infrastructure" src/hw_checker/domain/
# Empty? ✅/❌

# Domain doesn't import presentation
grep -r "from hw_checker.presentation" src/hw_checker/domain/
# Empty? ✅/❌
```

**Status:** ✅ / 🔴 BLOCKING

---

#### 5. Type Hints

```bash
mypy src/hw_checker/module/ --strict --no-implicit-optional
# No errors? ✅/❌

# Check -> None for void
grep -rn "def.*:" src/hw_checker/module/*.py | grep -v "-> "
# Should be empty ✅
```

**Status:** ✅ / 🔴 BLOCKING

---

#### 6. Error Handling

```bash
# No except: pass
grep -rn "except.*:" src/hw_checker/module/ -A1 | grep "pass"
# Empty? ✅/❌

# No bare except
grep -rn "except:" src/hw_checker/module/
# Empty? ✅/❌
```

**Status:** ✅ / 🔴 BLOCKING

---

#### 7. Security

```bash
# No SQL injection
grep -rn "execute.*%" src/hw_checker/module/
# Empty? ✅/❌

# No shell injection
grep -rn "subprocess.*shell=True" src/hw_checker/module/
# Empty? ✅/❌

bandit -r src/hw_checker/module/ -ll
# No issues? ✅/❌
```

**Status:** ✅ / 🔴 BLOCKING

---

#### 8. No Tech Debt

```bash
grep -rn "TODO\|FIXME\|HACK\|XXX" src/hw_checker/module/
# Empty? ✅/❌
```

**Status:** ✅ / 🔴 BLOCKING

---

#### 9. Documentation

- [ ] Docstrings for public functions
- [ ] Type hints everywhere
- [ ] README updated (if needed)

**Status:** ✅ / ⚠️ WARNING

---

#### 10. Git History

```bash
# Check commits exist for WS
git log --oneline main..HEAD | grep "WS-XXX-YY"
# Should have commits ✅/❌

# Check commit format (conventional commits)
git log --oneline main..HEAD
# Should be: feat(), test(), docs(), fix()
```

**Status:** ✅ / ⚠️ WARNING

---

### Stage 2 Verdict

**PASS:** All checks ✅ → APPROVED

**FAIL:** Any check 🔴 → CHANGES REQUESTED → Fix → Re-review Stage 2

**Output Format:**

```markdown
## Stage 2: Code Quality

**Date:** {YYYY-MM-DD}
**Reviewer:** {agent}

| Check | Status | Notes |
|-------|--------|-------|
| Tests & Coverage | ✅ / ⚠️ / 🔴 | {coverage}% |
| Regression | ✅ / 🔴 | {test_count} tests |
| AI-Readiness | ✅ / ⚠️ / 🔴 | max {loc} LOC, CC {cc} |
| Clean Architecture | ✅ / 🔴 | {notes} |
| Type Hints | ✅ / 🔴 | {notes} |
| Error Handling | ✅ / 🔴 | {notes} |
| Security | ✅ / 🔴 | {notes} |
| No Tech Debt | ✅ / 🔴 | {notes} |
| Documentation | ✅ / ⚠️ | {notes} |
| Git History | ✅ / ⚠️ | {notes} |

**Verdict:** ✅ PASS / 🔴 FAIL

**Issues (if FAIL):**
1. {Issue 1}
2. {Issue 2}
```

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

1. **After fix:** Re-run the failed stage only (not both stages)
2. **Stage 1 fix:** Re-run Stage 1 → if PASS → proceed to Stage 2
3. **Stage 2 fix:** Re-run Stage 2 only (Stage 1 already passed)

### Example

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

## Final Verdict

### APPROVED

**Conditions:**
- ✅ Stage 1: PASS
- ✅ Stage 2: PASS
- ✅ All blocking checks passed

### CHANGES REQUESTED

**Conditions:**
- ❌ Stage 1: FAIL (any blocking check)
- ❌ Stage 2: FAIL (any blocking check)

**Only two verdicts allowed:** APPROVED (all checks pass) or CHANGES REQUESTED (any check fails).

**❌ Anti-pattern:** Do NOT use partial approvals like "Approved with minor notes" or similar variants. Binary verdicts enforce clarity.

---

## Output Format (Per WS)

```markdown
---

### Review Results

**Date:** {YYYY-MM-DD}
**Reviewer:** {agent}
**Verdict:** APPROVED / CHANGES REQUESTED

#### Stage 1: Spec Compliance

| Check | Status | Notes |
|-------|--------|-------|
| Goal Achievement | ✅ | {X}/{Y} AC passed |
| Specification Alignment | ✅ | {notes} |
| AC Coverage | ✅ | {coverage details} |
| No Over-Engineering | ✅ | {notes} |
| No Under-Engineering | ✅ | {notes} |

**Stage 1 Verdict:** ✅ PASS / 🔴 FAIL

#### Stage 2: Code Quality

| Check | Status | Notes |
|-------|--------|-------|
| Tests & Coverage | ✅ | {coverage}% |
| Regression | ✅ | {test_count} tests |
| AI-Readiness | ✅ | max {loc} LOC, CC {cc} |
| Clean Architecture | ✅ | {notes} |
| Type Hints | ✅ | {notes} |
| Error Handling | ✅ | {notes} |
| Security | ✅ | {notes} |
| No Tech Debt | ✅ | {notes} |
| Documentation | ✅ | {notes} |
| Git History | ✅ | {notes} |

**Stage 2 Verdict:** ✅ PASS / 🔴 FAIL

#### Issues (if CHANGES REQUESTED)

| # | Stage | Severity | Description | How to Fix |
|---|-------|----------|-------------|------------|
| 1 | 1 | CRITICAL | AC2 not working | Fix X in Y |
| 2 | 2 | HIGH | Coverage 75% | Add tests for Z |
```

---

## Key Principles

1. **Stage 1 First:** Always check spec compliance before code quality
2. **No Wasted Effort:** Don't polish code that doesn't match spec
3. **Clear Separation:** Spec issues vs. quality issues are different
4. **Review Loop:** Fix → re-review same stage (not both)
5. **Zero Tolerance:** No "minor issues" — all blockers must be fixed

---

## Integration with /codereview

This two-stage protocol is used by `/codereview` command:

1. `/codereview F60` → For each WS:
   - Run Stage 1 (Spec Compliance)
   - If PASS → Run Stage 2 (Code Quality)
   - If FAIL → Report issues → Fix → Re-review

2. Review loop continues until both stages pass

3. Final verdict: APPROVED (both stages pass) or CHANGES REQUESTED (any stage fails)
