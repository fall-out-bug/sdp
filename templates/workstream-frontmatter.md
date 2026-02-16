---
ws_id: PP-FFF-SS
feature: FFFF
status: backlog|active|completed|blocked
size: SMALL|MEDIUM|LARGE
project_id: PP
github_issue: null
assignee: null
depends_on:
  - PP-FFF-SS  # Optional: list of dependent WS IDs
---

## WS-PP-FFF-SS: Title

### Goal

**What should WORK after completing this WS:**
- [First specific outcome]
- [Second specific outcome]

**Acceptance Criteria:**
- [ ] AC1: [First criterion - specific, measurable]
- [ ] AC2: [Second criterion - specific, measurable]
- [ ] AC3: [Third criterion - specific, measurable]

**WARNING: WS is NOT complete until the Goal is achieved (all ACs checked).**

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

### Expected Result

[Description of expected outcome]

### Scope Estimate

- Files: ~[number]
- Lines: ~[number] ([SMALL|MEDIUM|LARGE])
- Tokens: ~[number]

### Completion Criteria

```bash
# Verification commands
test -f path/to/file
grep "expected content" path/to/file
echo "Verification passed"
```

### Constraints

- DO NOT [constraint 1]
- DO NOT [constraint 2]

---

## Execution Report

**Executed by:** [Name/Agent]
**Date:** YYYY-MM-DD

### Goal Status
- [x] AC1: [description] - DONE
- [x] AC2: [description] - DONE
- [x] AC3: [description] - DONE

**Goal Achieved:** DONE YES

### Files Changed
| File | Action | LOC |
|------|--------|-----|
| `path/to/file.py` | created | 120 |
| `path/to/test.py` | created | 80 |

### Self-Check Results
```bash
$ pytest tests/unit/test_module.py -v
===== 15 passed in 0.5s =====

$ pytest --cov=module --cov-fail-under=80
===== Coverage: 85% =====
```

### Commit
{commit_hash} - {commit_message}
