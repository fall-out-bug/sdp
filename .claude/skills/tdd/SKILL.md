---
name: tdd
description: "Enforce Test-Driven Development discipline: Red -> Green -> Refactor (INTERNAL - used by @build)"
tools: Read, Write, Edit, Bash
---

# /tdd - Test-Driven Development (INTERNAL)

**INTERNAL SKILL** — Automatically called by `/build`, not invoked directly by users.

Enforce TDD discipline with Red-Green-Refactor cycle.

## Purpose

Called automatically by `@build` to ensure:
- Tests written BEFORE implementation
- Minimal code in Green phase
- Refactoring doesn't break tests

## The TDD Cycle

### Phase 1: RED - Write Failing Test

1. **Write test FIRST** - before any implementation code
2. **Run test** - verify it FAILS with expected error
3. **NO implementation yet** - if you wrote code, you cheated

### Phase 2: GREEN - Minimal Implementation

1. **Write minimal code** - just enough to make test pass
2. **Run test** - verify it PASSES
3. **NO refactoring yet** - that comes next

### Phase 3: REFACTOR - Improve Code

1. **Improve code** - clean up, extract, rename
2. **Run test** - verify it STILL PASSES
3. **Add more tests** if new edge cases discovered

### Phase 4: COMMIT - Save Working State

1. **Commit** - each cycle ends in a commit
2. **Message** - describes what was built

## Self-Review Checklist

After each cycle:
- [ ] Test written BEFORE implementation
- [ ] Test verified FAILING in Red phase
- [ ] Only minimal code in Green phase
- [ ] All tests passing after Refactor
- [ ] Commit created with conventional message

## Exit When

- All acceptance criteria met
- Coverage >= 80%
- mypy --strict passes

## Full Cycle Example

```
# Phase 1: RED - Write Failing Test

→ Writing test for email validation...

File: tests/unit/test_validators.py
```python
def test_email_validation():
    validator = EmailValidator()
    assert validator.is_valid("user@example.com") is True
    assert validator.is_valid("invalid-email") is False
```

→ Running test...
$ pytest tests/unit/test_validators.py
FAILED - ModuleNotFoundError: No module named 'validators'
✓ Test fails as expected (RED phase complete)

# Phase 2: GREEN - Minimal Implementation

→ Writing minimal implementation...

File: src/validators.py
```python
import re

class EmailValidator:
    def is_valid(self, email: str) -> bool:
        pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
        return bool(re.match(pattern, email))
```

→ Running test...
$ pytest tests/unit/test_validators.py
PASSED
✓ Test passes (GREEN phase complete)

# Phase 3: REFACTOR - Improve Code

→ Refactoring: Extract pattern as constant...

File: src/validators.py
```python
import re
from typing import Final

EMAIL_PATTERN: Final[str] = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'

class EmailValidator:
    def is_valid(self, email: str) -> bool:
        return bool(re.match(EMAIL_PATTERN, email))
```

→ Running test...
$ pytest tests/unit/test_validators.py
PASSED
✓ Tests still pass after refactor

→ Running coverage...
$ pytest --cov=src/validators tests/unit/test_validators.py
Coverage: 100%
✓ Coverage >= 80%

→ Running type check...
$ mypy src/validators.py --strict
Success: no issues found
✓ Type checking passed

# Phase 4: COMMIT - Save State

→ Creating commit...
$ git add src/validators.py tests/unit/test_validators.py
$ git commit -m "feat(validators): add email validation

- Add EmailValidator class with regex pattern
- Cover basic valid/invalid cases
- 100% test coverage"

✓ TDD cycle complete! Ready for next AC.
```
