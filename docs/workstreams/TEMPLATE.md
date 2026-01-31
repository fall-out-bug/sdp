---
ws_id: PP-FFF-WW
feature: FFF
status: backlog|active|completed|blocked
size: SMALL|MEDIUM|LARGE
project_id: PP
github_issue: null
assignee: null
depends_on:
  - PP-FFF-WW  # Optional: list of dependent WS IDs
---

## WS-PP-FFF-WW: Title

### 🎯 Goal

**What must WORK after completing this WS:**
- [First specific outcome]
- [Second specific outcome]

**Acceptance Criteria:**
- [ ] AC1: [First criterion - specific, measurable]
- [ ] AC2: [Second criterion - specific, measurable]
- [ ] AC3: [Third criterion - specific, measurable]
- [ ] Coverage ≥ 80%
- [ ] mypy --strict passes

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

[Background information about the workstream]

### Dependencies

[List dependencies or write "None" for no dependencies]

### Input Files

[List input files or sections to read]

### Steps

1. **[Step 1 title]**

   [Detailed instructions for step 1]

2. **[Step 2 title]**

   [Detailed instructions for step 2]

### Code

```python
# Ready-to-use code for copy-paste
# Full type hints
```

### Expected Outcome

[Description of expected outcome]

### Scope Estimate

- Files: ~[number]
- Lines: ~[number] ([SMALL|MEDIUM|LARGE])
- Tokens: ~[number]

### Completion Criteria

```bash
# Verification commands
pytest tests/unit/test_module.py -v
pytest --cov=module --cov-fail-under=80
mypy module/ --strict
```

### Constraints

- DO NOT [constraint 1]
- DO NOT [constraint 2]

---

## Execution Report

**Executed by:** ______
**Date:** ______
**Duration:** ______ minutes

### Goal Status
- [ ] AC1-AC3 — ✅

**Goal Achieved:** ______

### Files Changed
| File | Action | LOC |
|------|--------|-----|
|      |        |     |

### Statistics
- **Files Changed:** ______
- **Lines Added:** ______
- **Lines Removed:** ______
- **Test Coverage:** ______ %
- **Tests Passed:** ______
- **Tests Failed:** ______

### Deviations from Plan
- ______

### Commit
______
