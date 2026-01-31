---
ws_id: 00-034-03
feature: F034
status: completed
complexity: LARGE
project_id: "00"
depends_on:
  - 00-034-02
---

# Workstream: Increase Test Coverage to 80%+

**ID:** 00-034-03  
**Feature:** F034 (A+ Quality Initiative)  
**Status:** READY  
**Owner:** AI Agent  
**Complexity:** LARGE (~2000 LOC new tests)

---

## Goal

Поднять test coverage с 68% до 80%+ по всем модулям `src/sdp/`.

---

## Context

Текущее состояние (из `pytest --cov`):

| Module | Current Coverage | Target | Gap |
|--------|-----------------|--------|-----|
| `validators/capability_tier.py` | 61% | 80% | +19% |
| `validators/capability_tier_t2_t3.py` | 64% | 80% | +16% |
| `beads/skills_oneshot.py` | 67% | 80% | +13% |
| `unified/checkpoint/` | 70% | 80% | +10% |
| `github/issue_sync.py` | 72% | 80% | +8% |
| `beads/sync.py` | 74% | 80% | +6% |
| Overall | **68%** | **80%** | **+12%** |

**Uncovered lines:** ~2,688 (need to cover ~1,500 to reach 80%)

---

## Scope

### In Scope
- ✅ Add unit tests for uncovered lines
- ✅ Add integration tests for critical paths
- ✅ Add edge case tests
- ✅ Fix any bugs discovered during testing

### Out of Scope
- ❌ E2E tests
- ❌ Performance tests
- ❌ UI tests

---

## Dependencies

**Depends On:**
- [x] 00-034-01: Split Large Files (Phase 1)
- [x] 00-034-02: Split Large Files (Phase 2)

**Blocks:**
- None (final quality gate)

---

## Acceptance Criteria

- [ ] `pytest tests/ --cov=src/sdp --cov-fail-under=80` passes
- [ ] No module below 75% coverage
- [ ] All critical paths have integration tests
- [ ] No flaky tests

---

## Implementation Plan

### Task 1: Cover `validators/capability_tier*.py` (Priority: HIGH)

**Target:** 61% → 80%

**Missing coverage areas:**
- T2/T3 edge cases (invalid configs)
- Tier extraction from malformed WS
- Error handling paths

**Tests to add:**
```python
# tests/unit/validators/test_capability_tier.py
def test_tier_extraction_invalid_ws():
    """Test tier extraction with missing fields."""
    
def test_tier_validation_t2_constraints():
    """Test T2 tier constraints."""
    
def test_tier_validation_t3_constraints():
    """Test T3 tier constraints."""
    
def test_tier_config_malformed():
    """Test handling of malformed tier config."""
```

### Task 2: Cover `beads/skills_oneshot.py` (Priority: HIGH)

**Target:** 67% → 80%

**Missing coverage areas:**
- Checkpoint save/restore
- Error recovery
- Parallel execution paths

**Tests to add:**
```python
# tests/unit/beads/test_skills_oneshot.py
def test_oneshot_checkpoint_save():
    """Test checkpoint is saved after each WS."""
    
def test_oneshot_checkpoint_restore():
    """Test execution resumes from checkpoint."""
    
def test_oneshot_error_recovery():
    """Test graceful error handling."""
```

### Task 3: Cover `unified/checkpoint/` (Priority: MEDIUM)

**Target:** 70% → 80%

**Missing coverage:**
- Serialization edge cases
- Corrupted checkpoint handling
- Storage errors

### Task 4: Cover `github/issue_sync.py` (Priority: MEDIUM)

**Target:** 72% → 80%

**Missing coverage:**
- Rate limiting
- Network errors
- Pagination

### Task 5: Cover `beads/sync.py` (Priority: LOW)

**Target:** 74% → 80%

**Missing coverage:**
- Conflict resolution
- Partial sync
- Dry-run mode

### Task 6: Integration Tests for Critical Paths

**Add integration tests:**
- [ ] Full @build workflow (Red → Green → Refactor)
- [ ] Full @oneshot workflow (multi-WS execution)
- [ ] Beads sync with real files
- [ ] GitHub sync with mock API

---

## DO / DON'T

### Testing

**✅ DO:**
- Test behavior, not implementation
- Use fixtures for common setup
- Test error paths, not just happy paths
- Use parametrized tests for variations

**❌ DON'T:**
- Mock domain entities
- Test private methods directly
- Create brittle tests tied to implementation
- Skip edge cases

---

## Files to Create

**Create:**
- [ ] `tests/unit/validators/test_capability_tier_coverage.py`
- [ ] `tests/unit/beads/test_skills_oneshot_coverage.py`
- [ ] `tests/unit/unified/test_checkpoint_coverage.py`
- [ ] `tests/unit/github/test_issue_sync_coverage.py`
- [ ] `tests/unit/beads/test_sync_coverage.py`
- [ ] `tests/integration/test_build_workflow.py`
- [ ] `tests/integration/test_oneshot_workflow.py`

**Modify:**
- [ ] Existing test files (add missing test cases)

---

## Test Plan

### Verification
```bash
# Run full coverage report
pytest tests/ --cov=src/sdp --cov-report=term-missing

# Verify 80% threshold
pytest tests/ --cov=src/sdp --cov-fail-under=80

# Check no module below 75%
pytest tests/ --cov=src/sdp --cov-report=html
# Review htmlcov/index.html
```

### CI Integration
- [ ] Add `--cov-fail-under=80` to CI pipeline
- [ ] Block PRs below threshold

---

## Estimated Effort

| Task | New Tests | LOC |
|------|-----------|-----|
| capability_tier | ~15 tests | ~400 LOC |
| skills_oneshot | ~12 tests | ~350 LOC |
| checkpoint | ~8 tests | ~250 LOC |
| issue_sync | ~8 tests | ~250 LOC |
| sync | ~6 tests | ~200 LOC |
| Integration | ~10 tests | ~500 LOC |
| **Total** | **~59 tests** | **~1950 LOC** |

---

**Version:** 1.0  
**Created:** 2026-01-31
