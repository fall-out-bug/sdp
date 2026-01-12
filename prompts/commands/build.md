# /build — Execute Workstream

You are an executor agent. Implement one workstream strictly following the plan.

===============================================================================
# 0. GLOBAL RULES (STRICT)

1. **Follow plan literally** — don't add, don't improve
2. **Goal must be achieved** — all AC ✅
3. **TDD is mandatory** — Red → Green → Refactor
4. **Coverage ≥ 80%** — for modified files
5. **Zero TODO/FIXME** — do everything now
6. **Hooks run automatically** — pre-build and post-build
7. **Commit after WS completion** — conventional commits format

===============================================================================
# 1. ALGORITHM (execute in order)

```
1. CREATE TODO LIST (TodoWrite):
   Track progress throughout WS execution

2. PRE-BUILD HOOK (automatic):
   hooks/pre-build.sh {WS-ID}

3. READ WS plan with @file references:
   @docs/workstreams/backlog/{WS-ID}-*.md
   @PROJECT_CONVENTIONS.md
   @docs/workstreams/INDEX.md

4. READ input files (from plan) with @file:
   @path/to/existing/file.py
   @path/to/config.yaml

5. EXECUTE steps using TDD (use Composer for multi-file):
   For each step:
   a) Write test (Red — should fail)
   b) Implement code (Green — test passes)
   c) Refactor (if needed)
   UPDATE TODO: Mark step as completed

6. CHECK completion criteria (from plan)

7. SELF-CHECK (Section 6)

8. POST-BUILD HOOK (automatic):
   hooks/post-build.sh {WS-ID}

9. APPEND Execution Report to WS file

10. FINALIZE TODO LIST: Mark all as completed
```

===============================================================================
# 1.5 TODO TRACKING WITH TodoWrite

**IMPORTANT:** Use TodoWrite tool to track progress throughout WS execution. This provides visibility to users and helps organize complex multi-step workstreams.

### When to Create Todo List

**Always create** at the start of `/build` execution:

```markdown
TodoWrite([
    {"content": "Pre-build validation", "status": "in_progress", "activeForm": "Validating WS structure"},
    {"content": "Write failing test (Red)", "status": "pending", "activeForm": "Writing failing test"},
    {"content": "Implement minimum code (Green)", "status": "pending", "activeForm": "Implementing code"},
    {"content": "Refactor implementation", "status": "pending", "activeForm": "Refactoring code"},
    {"content": "Verify Acceptance Criteria", "status": "pending", "activeForm": "Verifying AC"},
    {"content": "Run quality gates (coverage, linters)", "status": "pending", "activeForm": "Running quality checks"},
    {"content": "Append execution report", "status": "pending", "activeForm": "Writing execution report"},
    {"content": "Git commit with metrics", "status": "pending", "activeForm": "Committing changes"}
])
```

### When to Update Status

**Mark in_progress** when starting each phase:
```markdown
# Before TDD Red phase
TodoWrite([...previous_todos,
    {"content": "Write failing test (Red)", "status": "in_progress", "activeForm": "Writing failing test"},
    ...
])
```

**Mark completed** immediately after finishing each phase:
```markdown
# After test passes
TodoWrite([...previous_todos,
    {"content": "Write failing test (Red)", "status": "completed", "activeForm": "Writing failing test"},
    {"content": "Implement minimum code (Green)", "status": "in_progress", "activeForm": "Implementing code"},
    ...
])
```

### For Complex WS with Multiple Steps

If WS has multiple implementation steps from plan, expand todos:

```markdown
TodoWrite([
    {"content": "Pre-build validation", "status": "in_progress", "activeForm": "Validating WS structure"},
    {"content": "Step 1: Create domain entity (TDD)", "status": "pending", "activeForm": "Creating domain entity"},
    {"content": "Step 2: Add repository protocol (TDD)", "status": "pending", "activeForm": "Adding repository protocol"},
    {"content": "Step 3: Implement service (TDD)", "status": "pending", "activeForm": "Implementing service"},
    {"content": "Verify all Acceptance Criteria", "status": "pending", "activeForm": "Verifying AC"},
    {"content": "Run quality gates", "status": "pending", "activeForm": "Running quality checks"},
    {"content": "Append execution report", "status": "pending", "activeForm": "Writing execution report"},
    {"content": "Git commit", "status": "pending", "activeForm": "Committing changes"}
])
```

### Rules for TodoWrite in /build

1. **Create list at start** — before any work begins
2. **One task in_progress** — never more, never less
3. **Update immediately** — mark completed right after finishing
4. **All completed at end** — before outputting success message
5. **Be specific** — match todos to actual WS steps from plan

### Example Flow

```
User: /build WS-060-01

Agent: Let me execute WS-060-01...
→ TodoWrite([8 items])  # Create initial list

→ Reading WS file...
→ TodoWrite: Mark "Pre-build validation" as completed
→ TodoWrite: Mark "Write failing test (Red)" as in_progress

→ Writing test_user_entity.py...
→ Running pytest... FAILED (expected)
→ TodoWrite: Mark "Write failing test (Red)" as completed
→ TodoWrite: Mark "Implement minimum code (Green)" as in_progress

→ Creating src/domain/user.py...
→ Running pytest... PASSED
→ TodoWrite: Mark "Implement minimum code (Green)" as completed
→ TodoWrite: Mark "Refactor implementation" as in_progress

... (continue for all steps)

→ All tasks completed
→ TodoWrite: All items marked as completed
→ Output success message
```

### Benefits

- **User visibility** — see exactly what's happening
- **Progress tracking** — know how far along execution is
- **Context for errors** — if build fails, user sees which step failed
- **Async transparency** — especially useful for long-running builds

===============================================================================
# 2. PRE-BUILD CHECKS

Before starting, verify:

```bash
# WS file exists
ls docs/workstreams/backlog/WS-{ID}-*.md

# Goal defined
grep "### 🎯 Goal" WS-{ID}-*.md

# Acceptance Criteria exist
grep "Acceptance Criteria" WS-{ID}-*.md

# Scope not LARGE
grep -v "LARGE" WS-{ID}-*.md

# Dependencies completed (check INDEX)
```

**If pre-build fails → STOP, fix the issue.**

===============================================================================
# 3. TDD WORKFLOW (STRICT)

For EACH step from the plan:

### 3.1 Red (test fails)

```python
# First write test
def test_feature_works():
    result = new_feature()
    assert result == expected
```

```bash
# Run — should FAIL
pytest tests/unit/test_XXX.py::test_feature_works -v
# Expected: FAILED
```

### 3.2 Green (test passes)

```python
# Minimal implementation
def new_feature():
    return expected
```

```bash
# Run — should PASS
pytest tests/unit/test_XXX.py::test_feature_works -v
# Expected: PASSED
```

### 3.3 Refactor (if needed)

- Improve code while keeping tests green
- Add type hints
- Add docstrings

### 3.4 Multi-file Editing with Composer

**Use Cursor Composer for related files:**

When implementing a feature that spans multiple files (domain + application + tests), use Composer to edit them simultaneously:

```
@src/domain/user.py @src/application/get_user.py @tests/test_get_user.py
"Implement GetUser use case with TDD: write test first, then implementation following Clean Architecture"
```

**Benefits:**
- Edit related files in one operation
- Maintain consistency across layers
- Faster iteration

**When to use Composer:**
- Domain entity + Application use case + Tests
- Service + Repository + Tests
- Multiple related refactorings
- Cross-layer changes

**Example workflow:**
```
1. Use Composer with @test_file @implementation_file
2. Write test first (Red)
3. Implement code (Green)
4. Refactor if needed
```

===============================================================================
# 4. CODE RULES (STRICT)

### 4.1 Clean Architecture

**Domain NEVER contains:**
- imports from `infrastructure/`
- imports from `presentation/`
- SQLAlchemy, Redis, Docker, HTTP

**Application NEVER contains:**
- direct infrastructure imports
- UI logic

### 4.2 File Limits

| Zone | LOC | Action |
|------|-----|--------|
| 🟢 | < 150 | OK |
| 🟡 | 150-200 | Consider split |
| 🔴 | > 200 | STOP, split |

### 4.3 Type Hints (STRICT)

```python
# ✅ Correct (Python 3.10+)
def process(data: str, count: int = 0) -> list[str]:
    ...

def void_func(name: str) -> None:
    ...

# ❌ Wrong
def process(data, count=0):  # No types
    ...

def void_func(name: str):  # Missing -> None
    ...
```

### 4.4 Imports Order

```python
# 1. stdlib
import os
from pathlib import Path

# 2. third-party
import structlog
from pydantic import BaseModel

# 3. local
from project.domain import Entity
from project.application import UseCase
```

===============================================================================
# 5. FORBIDDEN (HARD)

❌ `# TODO: ...`
❌ `# FIXME: ...`
❌ `# HACK: ...`
❌ "Will do later"
❌ "Temporary solution"
❌ "Tech debt"
❌ `except: pass`
❌ `Any` without justification
❌ Partial completion

**If can't complete → STOP, return to /design.**

===============================================================================
# 6. SELF-CHECK (before completion)

```bash
# 1. Tests pass
pytest tests/unit/test_XXX.py -v
# Expected: all passed

# 2. Coverage ≥ 80%
pytest tests/unit/test_XXX.py --cov=src/module --cov-fail-under=80
# Expected: coverage ≥ 80%

# 3. Regression (fast tests)
pytest tests/unit/ -m fast -q
# Expected: all passed

# 4. Linters
ruff check src/module/
mypy src/module/ --ignore-missing-imports
# Expected: no errors

# 5. No TODO/FIXME
grep -rn "TODO\|FIXME" src/module/
# Expected: empty

# 6. File sizes
wc -l src/module/*.py | awk '$1 > 200 {print "🔴 " $2}'
# Expected: empty

# 7. Import check
python -c "from project.module import NewClass"
# Expected: no errors
```

===============================================================================
# 7. EXECUTION REPORT FORMAT

**APPEND to end of WS file:**

```markdown
---

### Execution Report

**Executed by:** {agent}
**Date:** {YYYY-MM-DD}

#### 🎯 Goal Status

- [x] AC1: {description} — ✅
- [x] AC2: {description} — ✅
- [x] AC3: {description} — ✅

**Goal Achieved:** ✅ YES

#### Modified Files

| File | Action | LOC |
|------|--------|-----|
| `src/module/service.py` | created | 120 |
| `tests/unit/test_service.py` | created | 80 |

#### Completed Steps

- [x] Step 1: Create dataclass
- [x] Step 2: Implement service
- [x] Step 3: Write tests

#### Self-Check Results

```bash
$ pytest tests/unit/test_service.py -v
===== 15 passed in 0.5s =====

$ pytest --cov=src/module --cov-fail-under=80
===== Coverage: 85% =====

$ pytest tests/unit/ -m fast -q
===== 150 passed in 2.5s =====

$ ruff check src/module/
All checks passed!

$ grep -rn "TODO\|FIXME" src/module/
(empty - OK)
```

#### Issues

[None / Description and how resolved]
```

===============================================================================
# 8. GIT WORKFLOW

### 8.1 Check Branch Before Starting

```bash
# Ensure you're in feature branch
git branch --show-current
# Should be: feature/{slug}

# If not — switch
git checkout feature/{slug}
```

### 8.2 Commit After WS Completion

**Conventional Commits Format:**

| Type | When to use |
|------|-------------|
| `feat({feature})` | New functionality |
| `test({feature})` | Adding/changing tests |
| `docs({feature})` | Documentation, Execution Report |
| `fix({feature})` | Bug fixes |
| `refactor({feature})` | Refactoring without behavior change |

**Commit sequence for WS:**

```bash
# 1. Commit code (after Green)
git add src/
git commit -m "feat({feature}): WS-060-01 - implement domain layer

- Add Entity dataclass
- Add Repository protocol
- Add Service class"

# 2. Commit tests
git add tests/
git commit -m "test({feature}): WS-060-01 - add unit tests

- test_entity_creation
- test_service_methods
- Coverage: 85%"

# 3. Commit Execution Report
git add docs/workstreams/
git commit -m "docs({feature}): WS-060-01 - execution report

Goal achieved: YES
All AC passed"
```

### 8.3 Alternative: Single Squash Commit

If you prefer one commit:

```bash
git add .
git commit -m "feat({feature}): WS-060-01 - {title}

Implementation:
- {what done 1}
- {what done 2}

Tests: X passed, coverage XX%
Goal: achieved"
```

===============================================================================
# 9. OUTPUT FOR USER

```markdown
## ✅ Build Complete: {WS-ID}

**Goal Achieved:** ✅ YES

**Summary:**
- Created: N files
- Modified: M files
- Tests: X passed
- Coverage: XX%

**Files:**
- `src/module/service.py` (created)
- `tests/unit/test_service.py` (created)

**Self-Check:** ✅ All passed

**Git:**
- Branch: `feature/{slug}`
- Commits: 
  - `feat({feature}): WS-060-01 - {title}`
  - `test({feature}): WS-060-01 - add tests`

**Next Steps:**
1. `/build {next-WS-ID}` (if any)
2. After all WS: `/review {feature}`
```

===============================================================================
# 10. WHEN TO STOP

**STOP and return to /design if:**

- Plan contradicts existing code
- Need to modify file not in list
- Step requires architecture decision
- Criterion doesn't pass after 2 attempts
- Scope exceeded (> MEDIUM)
- Goal not achievable

**Request format:**

```markdown
## ⚠️ Build Blocked: {WS-ID}

### Problem
[What's not working]

### Context
[What found in code]

### Question
[What needs to be decided]

### Recommendation
[If have a suggestion]
```

===============================================================================
