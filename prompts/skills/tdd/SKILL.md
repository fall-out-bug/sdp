---
name: tdd
description: "Enforce Test-Driven Development discipline: Red -> Green -> Refactor (INTERNAL - used by @build)"
tools:
  - Read
  - Write
  - Edit
  - Bash
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
- `go test ./...` passes
- `go vet ./...` passes

## Full Cycle Example

```
# Phase 1: RED - Write Failing Test

→ Writing test for email validation...

File: internal/validators/validators_test.go
```go
func TestEmailValidation(t *testing.T) {
    v := NewEmailValidator()
    if !v.IsValid("user@example.com") {
        t.Error("expected valid email to pass")
    }
    if v.IsValid("invalid-email") {
        t.Error("expected invalid email to fail")
    }
}
```

→ Running test...
$ go test ./internal/validators/...
FAIL: undefined: NewEmailValidator
✓ Test fails as expected (RED phase complete)

# Phase 2: GREEN - Minimal Implementation

→ Writing minimal implementation...

File: internal/validators/validators.go
```go
package validators

import "regexp"

func NewEmailValidator() *EmailValidator { return &EmailValidator{} }

type EmailValidator struct{}

func (e *EmailValidator) IsValid(email string) bool {
    re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return re.MatchString(email)
}
```

→ Running test...
$ go test ./internal/validators/...
PASS
✓ Test passes (GREEN phase complete)

# Phase 3: REFACTOR - Improve Code

→ Refactoring: Extract pattern as package variable...

→ Running test...
$ go test ./internal/validators/...
PASS
✓ Tests still pass after refactor

→ Running quality gates...
$ go vet ./...
$ go build ./...
✓ Quality gates passed

# Phase 4: COMMIT - Save State

→ Creating commit...
$ git add internal/validators/
$ git commit -m "feat(validators): add email validation"

✓ TDD cycle complete! Ready for next AC.
```
