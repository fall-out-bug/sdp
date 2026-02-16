# UAT Guide: F{XX} - {Feature Name}

**Created:** {YYYY-MM-DD}
**Feature:** F{XX}
**Workstreams:** WS-{XX}-01, WS-{XX}-02, ...

---

## Overview

{What the feature does in 2-3 sentences for a human}

---

## Prerequisites

Before testing, make sure:

- [ ] Docker is running (`docker ps`)
- [ ] `poetry install` executed in `tools/hw_checker/`
- [ ] `.env` or `hw_checker.yaml` is configured
- [ ] Database is accessible (if needed)
- [ ] Redis is running (if needed)

```bash
# Quick prerequisite check
cd tools/hw_checker
poetry run python -c "from hw_checker import __version__; print(f'Version: {__version__}')"
```

---

## Quick Verification (5 minutes)

### 1. Smoke Test

```bash
cd tools/hw_checker

# Main check
poetry run hwc {main_command}

# Expected result:
# {description of what should happen}
```

### 2. Visual Inspection

- [ ] Open {what to open: logs/UI/API}
- [ ] Verify that {what should be displayed}
- [ ] Make sure there are {no errors/warnings}

---

## Detailed Test Scenarios

### Scenario 1: Happy Path

**Description:** {main use case}

**Steps:**
1. {step 1}
2. {step 2}
3. {step 3}

**Expected:**
- {expectation 1}
- {expectation 2}

**Actual:** ____________________

**Status:** ⬜ Pass / ⬜ Fail

---

### Scenario 2: Error Handling

**Description:** {how the system handles errors}

**Steps:**
1. {trigger error condition}
2. {observe response}

**Expected:**
- Graceful error message (not a stack trace)
- Error logging
- System continues to work

**Actual:** ____________________

**Status:** ⬜ Pass / ⬜ Fail

---

### Scenario 3: Edge Cases

**Description:** {boundary cases}

**Steps:**
1. {edge case input}
2. {observe behavior}

**Expected:**
- {expected handling}

**Actual:** ____________________

**Status:** ⬜ Pass / ⬜ Fail

---

## Red Flags Checklist

**If you see any of these signs - the agent made a mistake:**

| # | Red Flag | What to Check | Severity |
|---|----------|---------------|----------|
| 1 | Stack trace in output | Logs, stderr | RED HIGH |
| 2 | Empty response | API response body | RED HIGH |
| 3 | Timeout (>30s) | Network, DB connection | YELLOW MEDIUM |
| 4 | Warnings in logs | Log files | YELLOW MEDIUM |
| 5 | Unexpected data format | Response structure | YELLOW MEDIUM |
| 6 | Deprecated warnings | Console output | GREEN LOW |

**What to do if you find a Red Flag:**
1. Copy the error message / screenshot
2. Check the corresponding WS Execution Report
3. Create an issue or go back to `/codereview`

---

## Code Sanity Checks

Quick check that the code is in order:

```bash
cd tools/hw_checker

# 1. No TODO/FIXME
grep -rn "TODO\|FIXME" src/hw_checker/{feature_module}/
# Expectation: empty

# 2. Reasonable file sizes
wc -l src/hw_checker/{feature_module}/*.py
# Expectation: all < 200 lines

# 3. Clean Architecture followed
grep -r "from hw_checker.infrastructure" src/hw_checker/domain/
# Expectation: empty

# 4. Tests pass
poetry run pytest tests/unit/test_{feature}*.py -v
# Expectation: all passed

# 5. Sufficient coverage
poetry run pytest tests/unit/test_{feature}*.py --cov=hw_checker/{feature_module} --cov-report=term-missing
# Expectation: >= 80%
```

---

## Performance Baseline (if applicable)

| Operation | Expected | Acceptable | Measured |
|----------|----------|------------|----------|
| {operation 1} | < 100ms | < 500ms | ___ms |
| {operation 2} | < 1s | < 5s | ___s |
| {operation 3} | < 5s | < 30s | ___s |

---

## Sign-off

### Pre-Sign-off Checklist

- [ ] All scenarios passed
- [ ] No red flags
- [ ] Code sanity checks passed
- [ ] Performance within baseline

### Approval

| Role | Name | Date | Signature |
|------|------|------|-----------|
| Developer (agent) | {agent} | {date} | CHECK |
| Reviewer | {reviewer} | {date} | BOX |
| **Human Tester** | ____________ | ____________ | BOX |

### Final Verdict

BOX **APPROVED** - ready for deploy
BOX **NEEDS WORK** - fixes required (see comments below)

### Comments

```
{comments from human tester}
```

---

## Related

- Feature Spec: `docs/specs/feature_{XX}/feature.md`
- Workstreams: `docs/workstreams/backlog/WS-{XX}-*.md`
- Review Results: see each WS file
