---
name: tdd
description: Enforce Test-Driven Development discipline: Red -> Green -> Refactor (INTERNAL - used by @build)
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
