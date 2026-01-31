# Build Command Full Specification

This document contains the complete specification for `@build`.
For quick reference, see [SKILL.md](../../.claude/skills/build/SKILL.md).

## Overview

The build command executes a single workstream following TDD.

## Prerequisites

- Workstream file exists in `docs/workstreams/backlog/`
- WS has Goal and Acceptance Criteria
- Dependencies are satisfied
- Scope is SMALL or MEDIUM

## Detailed Steps

### Pre-Build Validation

Run automatically or manually:

```bash
hooks/pre-build.sh {WS-ID}
```

Checks:
- WS file exists
- Goal section present
- AC defined
- Dependencies complete

### Guard Activation

The guard system ensures you only edit files within the workstream scope:

```bash
sdp guard activate {WS-ID}
```

This:
1. Reads `scope_files` from WS frontmatter
2. Creates `.sdp/active_ws.json` with allowed files
3. Enables pre-edit validation

### TDD Execution

For each Acceptance Criterion:

**Red Phase:**
1. Write a failing test that validates the AC
2. Run test to confirm it fails
3. Commit the failing test

**Green Phase:**
1. Write minimal code to make test pass
2. Run test to confirm it passes
3. Commit the passing code

**Refactor Phase:**
1. Improve code quality without changing behavior
2. Run all tests to confirm still passing
3. Commit the refactored code

### Quality Checks

Run after all ACs complete:

```bash
# Test coverage
pytest --cov=src --cov-report=term-missing --cov-fail-under=80

# Type checking
mypy src/ --strict

# Linting
ruff check src/

# File size check
find src/ -name "*.py" -exec wc -l {} \; | awk '$1 > 200 {print}'
```

### Completion

After all quality gates pass:

```bash
# Complete WS
sdp guard complete {WS-ID}

# Move to completed
mv docs/workstreams/backlog/{WS-ID}-*.md docs/workstreams/completed/

# Commit
git add .
git commit -m "feat({scope}): {WS-ID} - {title}"
```

## Progress Reporting

Report progress after each step:

```markdown
✅ Step 1/5: Create module skeleton
✅ Step 2/5: Implement core class
✅ Step 3/5: Add error handling
...
```

## Execution Report

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

## Common Issues

### Dependencies Not Met

**Error:** Pre-build validation fails with dependency error

**Fix:** Complete dependent workstreams first

### Files Outside Scope

**Error:** Guard rejects file edit

**Fix:** Update `scope_files` in WS frontmatter or use correct file

### Coverage Too Low

**Error:** Coverage check fails below 80%

**Fix:** Add more test cases to cover branches

### Type Errors

**Error:** mypy reports type errors

**Fix:** Add type hints or fix type mismatches

## Examples

See [docs/examples/build/](../examples/build/) for:
- TDD cycle examples
- Execution report samples
- Common patterns
