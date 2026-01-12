# Workstream: [Short Description]

**ID:** WS-XXX-YY  
**Feature:** F-XXX  
**Status:** READY  
**Owner:** [Agent/Developer]  
**Complexity:** [SMALL/MEDIUM/LARGE]

---

## Goal

[Clear, one-sentence goal of this workstream]

---

## Context

[Background information: why this workstream exists, how it fits into the feature]

---

## Scope

### In Scope
- ✅ [Item 1]
- ✅ [Item 2]
- ✅ [Item 3]

### Out of Scope
- ❌ [Item 1]
- ❌ [Item 2]

---

## Dependencies

**Depends On:**
- [ ] WS-XXX-YY: [Description]

**Blocks:**
- [ ] WS-XXX-YY: [Description]

---

## Acceptance Criteria

- [ ] Criterion 1: [Specific, measurable]
- [ ] Criterion 2: [Specific, measurable]
- [ ] Criterion 3: [Specific, measurable]
- [ ] Test coverage ≥ 80%
- [ ] Type hints complete (mypy --strict passes)
- [ ] No TODO/FIXME in code
- [ ] Clean Architecture followed

---

## Implementation Plan

### Tasks

1. **[Task 1 Title]**
   - Subtask 1.1
   - Subtask 1.2

2. **[Task 2 Title]**
   - Subtask 2.1
   - Subtask 2.2

3. **[Task 3 Title]**
   - Subtask 3.1

---

## DO / DON'T

### Architecture

**✅ DO:**
- Keep domain layer pure (no external dependencies)
- Define ports in Application layer
- Implement adapters in Infrastructure layer
- Follow Clean Architecture boundaries

**❌ DON'T:**
- Import Infrastructure in Domain
- Put business logic in controllers
- Skip port definitions
- Mix layers

### Code Quality

**✅ DO:**
- Write tests first (TDD: Red → Green → Refactor)
- Use explicit type hints everywhere
- Handle errors explicitly with clear messages
- Keep functions small (< 30 lines)
- Use meaningful names (no abbreviations)

**❌ DON'T:**
- Skip tests
- Use `except: pass` or bare exceptions
- Create God Objects (> 200 LOC)
- Use magic numbers/strings
- Leave TODO/FIXME without followup WS

### Testing

**✅ DO:**
- Test happy path and error cases
- Use test doubles for external services
- Assert specific exceptions with messages
- Aim for 100% coverage on critical paths

**❌ DON'T:**
- Test implementation details
- Mock domain entities
- Skip edge cases
- Write tests without assertions

### Documentation

**✅ DO:**
- Write docstrings for public APIs
- Include examples in docstrings
- Update README when adding features
- Document WHY, not WHAT

**❌ DON'T:**
- Write obvious comments ("increment i")
- Leave outdated comments
- Skip docstrings for complex logic
- Over-comment simple code

---

## Project-Specific DO/DON'T

> **Fill this section based on your `PROJECT_CONVENTIONS.md`**

### [Your Domain Area]

**✅ DO:**
- [Add project-specific best practice]
- [Add another best practice]

**❌ DON'T:**
- [Add project-specific anti-pattern]
- [Add another anti-pattern]

---

## Files to Modify/Create

**Create:**
- [ ] `path/to/new_file.py`
- [ ] `path/to/test_new_file.py`

**Modify:**
- [ ] `path/to/existing_file.py`
- [ ] `path/to/test_existing_file.py`

**Delete:**
- [ ] `path/to/deprecated_file.py`

---

## Test Plan

### Unit Tests
- [ ] Test case 1: [Description]
- [ ] Test case 2: [Description]
- [ ] Test case 3: [Description]

### Integration Tests
- [ ] Integration test 1: [Description]
- [ ] Integration test 2: [Description]

### Edge Cases
- [ ] Edge case 1: [Description]
- [ ] Edge case 2: [Description]

---

## Execution Report

### Completed Tasks
- [x] Task 1: [Summary]
- [x] Task 2: [Summary]
- [x] Task 3: [Summary]

### Test Results
```
pytest tests/ --cov=src
```

**Coverage:** [X]%  
**Tests Passed:** [X]/[Y]

### Decisions Made

**Decision 1:** [What was decided and why]

**Decision 2:** [What was decided and why]

### Blockers Encountered

**Blocker 1:** [Description and resolution]

**Blocker 2:** [Description and resolution]

---

## Review Checklist

**Code Quality:**
- [ ] All acceptance criteria met
- [ ] DO/DON'T rules followed
- [ ] Test coverage ≥ 80%
- [ ] Type hints complete (mypy --strict passes)
- [ ] No `except: pass` or silent failures
- [ ] Files < 200 LOC
- [ ] Functions < 30 LOC
- [ ] Cyclomatic complexity < 10

**Architecture:**
- [ ] Clean Architecture boundaries respected
- [ ] Domain layer pure (no external deps)
- [ ] Ports defined in Application layer
- [ ] Adapters in Infrastructure layer

**Testing:**
- [ ] Unit tests for business logic
- [ ] Integration tests for adapters
- [ ] Edge cases covered
- [ ] No flaky tests

**Documentation:**
- [ ] Docstrings for public APIs
- [ ] README updated if needed
- [ ] No TODO/FIXME in code
- [ ] Execution report filled

---

## Sign-off

**Completed By:** [Name/Agent]  
**Reviewed By:** [Name/Agent]  
**Date:** YYYY-MM-DD  
**Status:** [READY → EXECUTING → DONE]

---

**Version:** 1.0  
**Last Updated:** YYYY-MM-DD
