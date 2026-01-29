---
name: build
description: Execute single workstream with TDD cycle enforcement and progress tracking
tools: Read, Write, Edit, Bash, Skill
---

# /build - Execute Workstream

Execute a single workstream following TDD discipline with progress reporting.

## When to Use

- Implementing a single workstream
- After @design creates workstreams
- For focused, tracked development

## Workflow

### Step 1: Read Workstream

Read the workstream file to understand requirements:

```bash
@docs/workstreams/in_progress/WS-XXX-YY.md
# or
@docs/workstreams/backlog/WS-XXX-YY.md
```

### Step 2: Verify Prerequisites

Check that all prerequisites are met:

```bash
# From WS Prerequisites section
pytest tests/unit/XXX/ -v
# or other check commands
```

**If prerequisites fail:** Stop and report dependency issue.

### Step 3: Execute Steps with TDD

For each step in the workstream:

1. **Mark progress:** "Step N/M: [Step description]"
2. **Call /tdd** for Red-Green-Refactor cycle
3. **Verify output** matches expected result
4. **Self-review** after each step

```python
# TDD execution pattern
Skill("tdd")
# Follows: Red (write failing test) -> Green (minimal code) -> Refactor
```

### Step 4: Verify Acceptance Criteria

After all steps complete:

```bash
# Run all tests
pytest tests/unit/ -v

# Check coverage
pytest --cov=src --cov-fail-under=80

# Type check
mypy src/ --strict
```

### Step 5: Append Execution Report

Generate and append execution report to the workstream file:

```python
from sdp.report.generator import ReportGenerator

generator = ReportGenerator(ws_id="WS-XXX-YY")
generator.start_timer()
# ... execute workstream ...
generator.stop_timer()

# Collect statistics
stats = generator.collect_stats(
    files_changed=[("src/module.py", "modified", 100)],
    coverage_pct=85.0,
    tests_passed=12,
    tests_failed=0,
    deviations=["Added extra validation for edge case"]
)

# Get current commit
import subprocess
commit_hash = subprocess.run(
    ["git", "rev-parse", "HEAD"],
    capture_output=True,
    text=True
).stdout.strip()

# Append report
generator.append_report(
    stats,
    executed_by="developer-name",
    commit_hash=commit_hash
)
```

### Step 6: Move to Completed

```bash
mv docs/workstreams/in_progress/WS-XXX-YY.md docs/workstreams/completed/
```

## Progress Reporting

Report progress after each step:

```markdown
✅ Step 1/5: Create module skeleton
✅ Step 2/5: Implement core class
✅ Step 3/5: Add error handling
...
```

## Quality Gates

After each step:
- [ ] Test written BEFORE implementation
- [ ] Test verified FAILING in Red phase
- [ ] Only minimal code in Green phase
- [ ] All tests passing after Refactor

## Output

- Working code with tests
- Updated workstream status
- Progress tracking

## Next Step

`/build WS-XXX-YY+1` for next workstream
