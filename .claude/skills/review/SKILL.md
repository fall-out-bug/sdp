---
name: review
description: Review feature quality using Beads task hierarchy. Validates all workstreams against quality gates.
tools: Read, Write, Edit, Bash, Glob, Grep
version: 2.0.0-beads
---

# @review - Quality Review (Beads Integration)

Review feature by validating all workstreams against quality gates using Beads task hierarchy.

## When to Use

- After all feature workstreams are complete
- Before merging to main branch
- To ensure quality standards are met
- As part of code review process

## Beads vs Markdown Workflow

**This skill reads workstreams from Beads task hierarchy.**

For traditional markdown workflow, use existing review process.

## Invocation

```bash
@review bd-0001
```

**Environment Variables:**
- `BEADS_USE_MOCK=true` - Use mock Beads (default for dev)

## Workflow

### Step 1: Initialize Review

```python
from sdp.beads import create_beads_client
import os

use_mock = os.getenv("BEADS_USE_MOCK", "true").lower() == "true"
client = create_beads_client(use_mock=use_mock)
```

### Step 2: Read Feature and Workstreams

```python
# Get feature task
feature = client.get_task(feature_id)

if not feature:
    print(f"❌ Feature not found: {feature_id}")
    return

# Get all sub-tasks (workstreams)
workstreams = client.list_tasks(parent_id=feature_id)

if not workstreams:
    print(f"❌ No workstreams found for feature: {feature_id}")
    return

print(f"📋 Reviewing feature: {feature.title}")
print(f"   Workstreams: {len(workstreams)}")
```

### Step 3: Validate Quality Gates

For each workstream, check:

```python
from sdp.beads import BeadsStatus

passed = 0
failed = 0
issues = []

for ws in workstreams:
    # Check if completed
    if ws.status != BeadsStatus.CLOSED:
        issues.append(f"{ws.id}: Not complete (status={ws.status.value})")
        failed += 1
        continue
    
    # Check quality gates
    ws_issues = validate_quality_gates(ws)
    
    if ws_issues:
        issues.append(f"{ws.id}: {len(ws_issues)} issues")
        failed += 1
    else:
        passed += 1
```

### Quality Gates Checklist

```markdown
**Coverage:**
- [ ] Test coverage ≥ 80%
- [ ] All critical paths covered

**Code Quality:**
- [ ] No `except: pass` statements
- [ ] No TODO/FIXME in production code
- [ ] Files < 200 LOC (AI-readability)
- [ ] Type hints present
- [ ] MyPy strict mode passes

**Clean Architecture:**
- [ ] No layer violations
- [ ] Dependencies point inward
- [ ] Domain logic pure

**Documentation:**
- [ ] Docstrings on public APIs
- [ ] README updated if needed
- [ ] Changelog entry added

**Tests:**
- [ ] Unit tests present
- [ ] Integration tests if needed
- [ ] All tests passing
```

### Step 4: Generate Review Report

```python
# Print summary
print(f"\n{'='*60}")
print(f"Review Report for {feature.title}")
print(f"{'='*60}")
print(f"Workstreams: {len(workstreams)}")
print(f"Passed: {passed} ✅")
print(f"Failed: {failed} ❌")

if issues:
    print(f"\nIssues:")
    for issue in issues:
        print(f"  ❌ {issue}")
else:
    print(f"\n✅ All workstreams passed review!")

print(f"{'='*60}")

# Update feature status based on review
if failed == 0:
    client.update_task_status(feature_id, BeadsStatus.CLOSED)
    print(f"\n✅ Feature {feature_id} approved!")
else:
    client.update_task_status(feature_id, BeadsStatus.BLOCKED)
    print(f"\n❌ Feature {feature_id} blocked - fix issues and re-run @review")
```

## Output

**Success:**
```
============================================================
Review Report for User Authentication
============================================================
Workstreams: 3
Passed: 3 ✅
Failed: 0 ❌

✅ All workstreams passed review!

============================================================
✅ Feature bd-0001 approved!
```

**Failure:**
```
============================================================
Review Report for User Authentication
============================================================
Workstreams: 3
Passed: 2 ✅
Failed: 1 ❌

Issues:
  ❌ bd-0001.3: 3 issues
     - Coverage: 65% (need ≥80%)
     - File src/auth/service.py: 250 LOC (need <200 LOC)
     - Missing type hints on authenticate()

============================================================
❌ Feature bd-0001 blocked - fix issues and re-run @review
```

## Example Session

```bash
# After @oneshot completes
@review bd-0001

# Output:
📋 Reviewing feature: User Authentication
   Workstreams: 3

Checking quality gates...
  bd-0001.1: Domain entities ✅
  bd-0001.2: Repository layer ✅
  bd-0001.3: Service layer ❌

Issues found in bd-0001.3:
  - Coverage: 72% (need ≥80%)
  - Missing docstrings

# Fix issues...
@build bd-0001.3  # Add tests, docstrings

# Re-run review
@review bd-0001

# Output:
✅ All workstreams passed review!
✅ Feature bd-0001 approved!
```

## Quality Gates Reference

| Gate | Requirement | Check |
|------|------------|-------|
| **Coverage** | ≥80% | `pytest --cov` |
| **File size** | <200 LOC | `wc -l` |
| **Type hints** | Present | `mypy --strict` |
| **Clean arch** | No violations | Manual review |
| **No TODOs** | None in code | `grep -r "TODO"` |

## Integration with Existing Review

This skill extends existing review process:
- Reads from Beads instead of file system
- Validates same quality gates
- Updates Beads task status
- Compatible with existing review checklist

---

**Version:** 2.0.0-beads
**Status:** Beads Integration
**See Also:** `@idea`, `@design`, `@build`, `@oneshot`
