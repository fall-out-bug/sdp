# TDD Pattern

> Test-Driven Development discipline for all code changes

## The Cycle

```
┌─────────────────────────────────────┐
│                                     │
│   RED → Write failing test first    │
│    ↓                                │
│   GREEN → Write minimal code        │
│    ↓                                │
│   REFACTOR → Clean up code          │
│    ↓                                │
│   (repeat)                          │
│                                     │
└─────────────────────────────────────┘
```

## Rules

1. **RED**: You MUST write a failing test before any implementation code
2. **GREEN**: Write ONLY enough code to make the test pass
3. **REFACTOR**: Clean up code while keeping tests green
4. **NO EXCEPTIONS**: Even for "quick fixes" or "obvious" code

## Verification Commands

```bash
# Go
go test -v ./path/to/package -run TestFunctionName
go test -cover ./...

# Check coverage threshold (>=80%)
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total
```

## Forbidden Patterns

- Writing implementation before tests
- Skipping tests because "it's simple"
- Tests that always pass (no assertions)
- Commenting out failing tests

## Required Patterns

- Test file colocation (`*_test.go` next to source)
- Table-driven tests for multiple cases
- Descriptive test names (`TestFunction_scenario_expectedResult`)
- Explicit error handling in tests

## When To Apply

- ALL new functions/methods
- ALL bug fixes (write test that reproduces bug first)
- ALL refactoring (tests must exist and pass before/after)

## See Also

- `@build` - Executes workstream with TDD enforcement
- `@bugfix` - Bug fixes follow TDD
- `.claude/skills/tdd.md` - Full skill definition
