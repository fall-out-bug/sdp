---
name: tdd
description: Enforce Test-Driven Development discipline: Red -> Green -> Refactor (INTERNAL - used by @build)
tools: Read, Write, Edit, Bash
---

# /tdd - Test-Driven Development (INTERNAL)

**INTERNAL SKILL** — Automatically called by `/build`, not invoked directly by users.

Enforce TDD discipline with Red-Green-Refactor cycle.

**Language-Agnostic:** Works with Python, Java, Go, and other languages.

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
- Type checking passes (language-specific)

## Language-Specific Examples

### Python Example

```
# Phase 1: RED - Write Failing Test

File: tests/unit/test_validators.py
```python
def test_email_validation():
    validator = EmailValidator()
    assert validator.is_valid("user@example.com") is True
    assert validator.is_valid("invalid-email") is False
```

Running test...
$ pytest tests/unit/test_validators.py
FAILED - ModuleNotFoundError: No module named 'validators'
✓ Test fails as expected (RED phase complete)

# Phase 2: GREEN - Minimal Implementation

File: src/validators.py
```python
import re

class EmailValidator:
    def is_valid(self, email: str) -> bool:
        pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
        return bool(re.match(pattern, email))
```

Running test...
$ pytest tests/unit/test_validators.py
PASSED
✓ Test passes (GREEN phase complete)

# Phase 3: REFACTOR - Improve Code

File: src/validators.py
```python
import re
from typing import Final

EMAIL_PATTERN: Final[str] = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'

class EmailValidator:
    def is_valid(self, email: str) -> bool:
        return bool(re.match(EMAIL_PATTERN, email))
```

Running test...
$ pytest tests/unit/test_validators.py
PASSED
✓ Tests still pass after refactor

Running quality checks...
$ pytest --cov=src/ tests/unit/test_validators.py
Coverage: 100% ✓

$ mypy src/validators.py --strict
Success: no issues found ✓
```

### Java Example

```
# Phase 1: RED - Write Failing Test

File: src/test/java/com/example/EmailValidatorTest.java
```java
import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;

class EmailValidatorTest {
    @Test
    void testEmailValidation() {
        EmailValidator validator = new EmailValidator();
        assertTrue(validator.isValid("user@example.com"));
        assertFalse(validator.isValid("invalid-email"));
    }
}
```

Running test...
$ mvn test -Dtest=EmailValidatorTest
FAILED - EmailValidator class not found
✓ Test fails as expected (RED phase complete)

# Phase 2: GREEN - Minimal Implementation

File: src/main/java/com/example/EmailValidator.java
```java
package com.example;

import java.util.regex.Pattern;

public class EmailValidator {
    private static final Pattern EMAIL_PATTERN = Pattern.compile(
        "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
    );

    public boolean isValid(String email) {
        return EMAIL_PATTERN.matcher(email).matches();
    }
}
```

Running test...
$ mvn test -Dtest=EmailValidatorTest
PASSED
✓ Test passes (GREEN phase complete)

# Phase 3: REFACTOR - Improve Code

Running test...
$ mvn test
PASSED
✓ Tests still pass after refactor

Running quality checks...
$ mvn verify
JaCoCo coverage: 100% ✓
```

### Go Example

```
# Phase 1: RED - Write Failing Test

File: service/validator_test.go
```go
package service

import "testing"

func TestEmailValidation(t *testing.T) {
    validator := NewEmailValidator()
    if !validator.IsValid("user@example.com") {
        t.Error("expected valid email to pass")
    }
    if validator.IsValid("invalid-email") {
        t.Error("expected invalid email to fail")
    }
}
```

Running test...
$ go test -v -run TestEmailValidation
--- FAIL: TestEmailValidation (0.00s)
    validator_test.go:10: undefined: NewEmailValidator
✓ Test fails as expected (RED phase complete)

# Phase 2: GREEN - Minimal Implementation

File: service/validator.go
```go
package service

import "regexp"

var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type EmailValidator struct{}

func NewEmailValidator() *EmailValidator {
    return &EmailValidator{}
}

func (e *EmailValidator) IsValid(email string) bool {
    return emailPattern.MatchString(email)
}
```

Running test...
$ go test -v -run TestEmailValidation
--- PASS: TestEmailValidation (0.00s)
✓ Test passes (GREEN phase complete)

# Phase 3: REFACTOR - Improve Code

Running test...
$ go test -v
--- PASS: TestEmailValidation (0.00s)
✓ Tests still pass after refactor

Running quality checks...
$ go test -coverprofile=coverage.out ./...
$ go tool cover -func=coverage.out | grep total
github.com/user/service/    100.0% ✓
```

## Quality Gates Summary

| Gate | Python | Java | Go |
|------|--------|------|-----|
| Tests | pytest | mvn test | go test |
| Coverage | pytest-cov ≥80% | JaCoCo ≥80% | go tool cover ≥80% |
| Type Check | mypy --strict | javac -Xlint:all | go vet |
| File Size | wc -l <200 | wc -l <200 | wc -l <200 |

**All gates must PASS for TDD cycle completion.**

## Common Mistakes

❌ **Writing implementation before test** - Violates Red phase
❌ **Writing too much code in Green** - Should be minimal
❌ **Skipping Refactor** - Code quality degrades
❌ **Committing without tests** - Violates TDD discipline

✅ **Test FIRST** - Always write test before code
✅ **Minimal Green** - Just enough to pass
✅ **Refactor thoroughly** - Clean up before commit
✅ **Commit each cycle** - Small, atomic changes

## See Also

- [@build Skill](../build/SKILL.md) - Calls /tdd automatically
- [Python Quick Start](../../docs/examples/python/QUICKSTART.md)
- [Java Quick Start](../../docs/examples/java/QUICKSTART.md)
- [Go Quick Start](../../docs/examples/go/QUICKSTART.md)
