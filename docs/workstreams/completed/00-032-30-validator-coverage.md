---
ws_id: 00-032-30
title: Add Tests for F032 Validators (Coverage ≥80%)
status: completed
feature: F032
project_id: 00
size: L
estimated_loc: 500
dependencies:
  - 00-032-29
completed: "2026-01-30"
review_source: docs/reports/2026-01-30-F032-review.md
---

# 00-032-30: Add Tests for F032 Validators (Coverage ≥80%)

## Goal

Increase test coverage for F032 validator modules from current levels to ≥80%.

## Context

F032 review found overall coverage 52.68%, with these F032-scope modules at critical levels:

| Module | Current | Target |
|--------|---------|--------|
| `validators/supersede_checker.py` | 0% | ≥80% |
| `validators/time_estimate_checker.py` | 0% | ≥80% |
| `validators/ws_completion.py` | 0% | ≥80% |
| `validators/ws_template_checker.py` | 0% | ≥80% |
| `validators/capability_tier_checks_interface.py` | 15% | ≥80% |
| `validators/capability_tier_checks_contract.py` | 57% | ≥80% |
| `validators/capability_tier.py` | 61% | ≥80% |

## Acceptance Criteria

- [x] AC1: `supersede_checker.py` coverage ≥80%
- [x] AC2: `time_estimate_checker.py` coverage ≥80%
- [x] AC3: `ws_completion.py` coverage ≥80%
- [x] AC4: `ws_template_checker.py` coverage ≥80%
- [x] AC5: `capability_tier_checks_*.py` coverage ≥80%
- [x] AC6: All new tests pass
- [x] AC7: Overall F032 module coverage ≥80%

## Technical Approach

1. **Identify uncovered branches:**
   ```bash
   pytest --cov=src/sdp/validators --cov-report=term-missing tests/
   ```

2. **Write tests per module** using TDD:
   - Happy path
   - Edge cases (empty input, invalid format)
   - Error conditions

3. **Test file structure:**
   ```
   tests/unit/validators/
   ├── test_supersede_checker.py
   ├── test_time_estimate_checker.py
   ├── test_ws_completion.py
   └── test_ws_template_checker.py
   ```

## Dependencies

- **00-032-29** — AC mappings must exist before writing tests that reference them

## Out of Scope

- Refactoring validators for complexity (see Issue 003)
- Integration tests (unit tests only)

## Execution Report

**Date:** 2026-01-30

**New test files:**
- `tests/unit/validators/test_supersede_checker.py` — 9 tests, 83% coverage
- `tests/unit/validators/test_time_estimate_checker.py` — 16 tests, 97% coverage
- `tests/unit/validators/test_ws_template_checker.py` — 6 tests, 81% coverage
- `tests/unit/validators/test_ws_completion.py` — 14 tests, 83% coverage
- `tests/unit/validators/test_capability_tier_checks_interface.py` — 9 tests, 100% coverage
- `tests/unit/validators/test_capability_tier_checks_contract.py` — 9 tests, 100% coverage
- `tests/unit/validators/test_capability_tier_checks_scope.py` — 4 tests, 100% coverage

**Result:** Validator module coverage 84% (≥80% gate). 115 tests pass.
