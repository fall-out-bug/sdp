# Review Command Full Specification

This document contains the complete specification for `@review`.
For quick reference, see [SKILL.md](../../.claude/skills/review/SKILL.md).

## Overview

The review command validates a feature by checking all workstreams against quality gates and ensuring traceability from Acceptance Criteria to tests.

## Quality Gates Checklist

### Coverage

- [ ] Test coverage ≥ 80%
- [ ] All critical paths covered
- [ ] Branch coverage for conditionals
- [ ] Edge cases tested

### Code Quality

- [ ] No `except: pass` statements
- [ ] No TODO/FIXME in production code
- [ ] Files < 200 LOC (AI-readability)
- [ ] Type hints present on all functions
- [ ] MyPy strict mode passes
- [ ] Ruff linter passes with no errors

### Clean Architecture

- [ ] No layer violations
- [ ] Dependencies point inward
- [ ] Domain logic pure (no infrastructure imports)
- [ ] Application layer doesn't import presentation
- [ ] Infrastructure implements ports from application

### Documentation

- [ ] Docstrings on public APIs
- [ ] README updated if needed
- [ ] Changelog entry added
- [ ] ADR created for architectural decisions

### Tests

- [ ] Unit tests present for all functions
- [ ] Integration tests if needed
- [ ] All tests passing
- [ ] Test names descriptive
- [ ] Assertions meaningful

### Traceability

- [ ] All ACs have corresponding tests
- [ ] Test names reference AC numbers
- [ ] Execution report shows AC completion
- [ ] No orphaned tests (tests without AC)

## Detailed Review Process

### 1. Feature Scope Verification

```bash
# List all workstreams for feature
bd list --parent {feature-id}

# Or for markdown
ls docs/workstreams/completed/{feature-id}-*.md
```

Verify:
- All planned workstreams are complete
- No workstreams left in backlog or in_progress
- Execution reports present for all WS

### 2. Traceability Matrix

For each workstream, verify AC→Test traceability using the traceability CLI:

```bash
sdp trace check {WS-ID}
```

This will automatically:
- Extract all ACs from the workstream
- Check for test mappings in metadata
- Display traceability table
- Exit 1 if any AC is unmapped

Example output:

```
Traceability Report: 00-032-01
==================================================
| AC | Description | Test | Status |
|----|-------------|------|--------|
| AC1 | User can login | `test_user_login` | ✅ |
| AC2 | Invalid creds rejected | `test_login_invalid` | ✅ |
| AC3 | Session created | - | ❌ |

Coverage: 67% (2/3 ACs mapped)
Status: ❌ INCOMPLETE (1 unmapped)
```

**Auto-detection:** If mappings are missing, try auto-detection:

```bash
sdp trace auto {WS-ID} --apply
```

This automatically detects mappings from:
- Test docstrings: `"""Tests AC1: User can login"""`
- Test names: `test_ac1_user_login`, `test_acceptance_criterion_1`
- Keyword matching between AC descriptions and test content

**Manual mapping:** If auto-detection fails, add manually:

```bash
sdp trace add {WS-ID} --ac AC1 --test test_user_login --file tests/test_auth.py
```

**Requirements:**
- All ACs must have mapped tests (100% coverage)
- Tests must be passing (use `pytest` to verify)
- Mappings stored in Beads task metadata

**Verification:**
```bash
# Check all WS for a feature
for ws in $(bd list --parent {feature-id}); do
  sdp trace check "$ws" || echo "FAILED: $ws"
done
```

### 3. Coverage Analysis

```bash
# Generate coverage report
pytest --cov=src --cov-report=html --cov-report=term-missing

# Open HTML report
open htmlcov/index.html
```

Review:
- Overall coverage percentage
- Files with low coverage
- Uncovered lines
- Branch coverage

### 4. Type Checking

```bash
# Run mypy in strict mode
mypy src/ --strict --show-error-codes

# Check for type ignore comments
grep -r "type: ignore" src/
```

All type errors must be fixed. No `type: ignore` without explanation.

### 5. Linting

```bash
# Run ruff
ruff check src/ tests/

# Check for complexity
ruff check src/ --select C90 --max-complexity 10
```

All linting errors must be fixed.

### 6. Architecture Validation

Manual review for layer violations:

```bash
# Check domain doesn't import infrastructure
grep -r "from.*infrastructure" src/domain/
grep -r "import.*infrastructure" src/domain/

# Check application doesn't import presentation
grep -r "from.*presentation" src/application/
grep -r "import.*presentation" src/application/
```

### 7. File Size Check

```bash
# Find large files
find src/ -name "*.py" -exec sh -c 'lines=$(wc -l < "$1"); [ $lines -gt 200 ] && echo "$1: $lines lines"' _ {} \;
```

Files over 200 LOC must be split.

### 8. Code Smell Check

```bash
# Check for except: pass
grep -r "except:" src/ | grep "pass"

# Check for TODOs
grep -r "TODO\|FIXME" src/

# Check for print statements (should use logging)
grep -r "print(" src/

# Check for hardcoded values
grep -r "localhost\|127.0.0.1" src/
```

## Verdict Decision Tree

```
All ACs have tests? 
  No → CHANGES_REQUESTED
  Yes ↓
  
All tests pass?
  No → CHANGES_REQUESTED  
  Yes ↓
  
Coverage ≥80%?
  No → CHANGES_REQUESTED
  Yes ↓
  
MyPy passes?
  No → CHANGES_REQUESTED
  Yes ↓
  
Ruff passes?
  No → CHANGES_REQUESTED
  Yes ↓
  
Files <200 LOC?
  No → CHANGES_REQUESTED
  Yes ↓
  
No code smells?
  No → CHANGES_REQUESTED
  Yes ↓
  
APPROVED ✅
```

## Review Report Template

```markdown
# Review Report: {Feature Title}

**Feature ID:** {feature-id}
**Reviewer:** {name}
**Date:** {date}

## Summary

- Workstreams reviewed: {count}
- Tests executed: {count}
- Coverage: {percentage}%
- Verdict: {APPROVED | CHANGES_REQUESTED}

## Traceability Matrix

| WS-ID | AC | Test | Status |
|-------|----|----- |--------|
| ... | ... | ... | ... |

## Quality Gates

- [ ] Coverage ≥80%
- [ ] MyPy strict passes
- [ ] Ruff passes
- [ ] Files <200 LOC
- [ ] No code smells
- [ ] Clean architecture

## Issues Found

{List any issues}

## Recommendation

{APPROVED | CHANGES_REQUESTED}

{Detailed explanation}
```

## Common Review Issues

### Issue: Missing Tests for AC

**Problem:** AC defined but no corresponding test

**Fix:** Add test that validates the AC

**Example:**
```python
def test_ac2_session_expires():
    """AC2: Session expires after timeout."""
    session = create_session()
    time.sleep(SESSION_TIMEOUT + 1)
    assert session.is_expired() is True
```

### Issue: Coverage Below 80%

**Problem:** Overall coverage is 65%

**Fix:** Add tests for uncovered branches

**Commands:**
```bash
# Find uncovered lines
pytest --cov=src --cov-report=term-missing

# Focus on specific module
pytest --cov=src/auth --cov-report=html
open htmlcov/index.html
```

### Issue: Layer Violation

**Problem:** Domain imports from infrastructure

**Fix:** Invert dependency using ports

**Before:**
```python
# src/domain/user.py
from src.infrastructure.database import UserRepository  # ❌
```

**After:**
```python
# src/application/ports.py
from abc import ABC, abstractmethod

class UserRepository(ABC):
    @abstractmethod
    def save(self, user: User) -> None: ...

# src/infrastructure/database.py
from src.application.ports import UserRepository

class SQLUserRepository(UserRepository):
    def save(self, user: User) -> None:
        # Implementation
```

### Issue: File Too Large

**Problem:** `src/service.py` is 350 lines

**Fix:** Split into multiple files

**Refactor:**
```
src/auth/
  service.py (150 lines)
  validators.py (100 lines)
  models.py (100 lines)
```

## Examples

See [docs/examples/review/](../examples/review/) for:
- Sample review reports
- Traceability matrices
- Common issue resolutions
