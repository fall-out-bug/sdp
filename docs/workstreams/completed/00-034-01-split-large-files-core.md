---
ws_id: 00-034-01
feature: F034
status: completed
complexity: MEDIUM
project_id: "00"
---

# Workstream: Split Large Files (Phase 1: Core)

**ID:** 00-034-01  
**Feature:** F034 (A+ Quality Initiative)  
**Status:** READY  
**Owner:** AI Agent  
**Complexity:** MEDIUM (~800 LOC refactoring)

---

## Goal

Разбить 10 файлов >200 LOC в `src/sdp/core/`, `src/sdp/cli/`, `src/sdp/validators/` на модули <200 LOC каждый.

---

## Context

Анализ показал 25 файлов, нарушающих правило <200 LOC. Phase 1 фокусируется на core-модулях:

**Топ нарушители (Phase 1):**

| File | LOC | Action |
|------|-----|--------|
| `core/workstream.py` | 402 | Split → parser, validator, models |
| `validators/ws_completion.py` | 358 | Split → checkers, reporter |
| `validators/capability_tier.py` | 329 | Split → extractors, validators |
| `validators/capability_tier_t2_t3.py` | 290 | Split → t2_checks, t3_checks |
| `cli/main.py` | 285 | Split → commands, app |
| `cli/workstream.py` | 267 | Split → commands, handlers |
| `core/decomposition.py` | 248 | Split → analyzer, builder |
| `validators/capability_tier_checks.py` | 235 | Split → contract, scope |
| `cli/beads.py` | 228 | Split → commands, formatters |
| `core/feature.py` | 215 | Split → loader, models |

---

## Scope

### In Scope
- ✅ Split 10 files listed above
- ✅ Update all imports across codebase
- ✅ Ensure all existing tests pass
- ✅ Add `__all__` exports for backward compatibility

### Out of Scope
- ❌ Files in `beads/`, `unified/` (Phase 2)
- ❌ Adding new tests (Phase 3)
- ❌ Changing public APIs
- ❌ Performance optimization

---

## Dependencies

**Depends On:**
- None (can start immediately)

**Blocks:**
- 00-034-02: Split Large Files (Phase 2)
- 00-034-03: Increase Test Coverage
- 00-034-05: Extract Domain Layer

---

## Acceptance Criteria

- [ ] All 10 files split to <200 LOC each
- [ ] `find src/sdp/core src/sdp/cli src/sdp/validators -name "*.py" -exec wc -l {} + | awk '$1 > 200'` returns empty
- [ ] All existing tests pass: `pytest tests/ -x`
- [ ] mypy --strict passes
- [ ] No import errors: `python -c "import sdp"`
- [ ] Backward compatibility: existing code using old imports still works via re-exports

---

## Implementation Plan

### Task 1: Split `core/workstream.py` (402 LOC)

**Target structure:**
```
core/
├── workstream/
│   ├── __init__.py      # Re-exports (backward compat)
│   ├── models.py        # Workstream, WorkstreamStatus dataclasses (~80 LOC)
│   ├── parser.py        # parse_workstream(), load_workstream() (~120 LOC)
│   └── validator.py     # validate_workstream(), WorkstreamValidator (~100 LOC)
└── workstream.py        # DEPRECATED: re-exports from workstream/ (~20 LOC)
```

**Steps:**
1. Create `core/workstream/` directory
2. Extract dataclasses → `models.py`
3. Extract parsing logic → `parser.py`
4. Extract validation → `validator.py`
5. Update `__init__.py` with re-exports
6. Update old `workstream.py` to re-export (deprecation warning)
7. Run tests, fix imports

### Task 2: Split `validators/ws_completion.py` (358 LOC)

**Target structure:**
```
validators/
├── ws_completion/
│   ├── __init__.py      # Re-exports
│   ├── checkers.py      # Individual check functions (~150 LOC)
│   ├── reporter.py      # Report generation (~100 LOC)
│   └── models.py        # CompletionResult, CheckResult (~50 LOC)
└── ws_completion.py     # DEPRECATED: re-exports
```

### Task 3: Split `validators/capability_tier*.py` (3 files)

**Target structure:**
```
validators/
├── capability_tier/
│   ├── __init__.py
│   ├── models.py        # TierResult, TierConfig
│   ├── extractors.py    # Extract tier info from WS
│   ├── validators.py    # Validate tier constraints
│   ├── t0_t1.py         # T0/T1 specific checks
│   └── t2_t3.py         # T2/T3 specific checks
```

### Task 4: Split `cli/main.py` (285 LOC)

**Target structure:**
```
cli/
├── main.py              # App entry, group registration (~50 LOC)
├── commands/
│   ├── __init__.py
│   ├── doctor.py        # Already exists
│   ├── guard.py         # Already exists
│   └── ...
```

### Task 5: Split `cli/workstream.py`, `cli/beads.py`

Similar pattern: extract command handlers into separate modules.

### Task 6: Split `core/decomposition.py`, `core/feature.py`

Extract analyzer/builder logic into submodules.

---

## DO / DON'T

### Refactoring

**✅ DO:**
- Keep public API unchanged (re-export from old locations)
- Add deprecation warnings for direct imports from old files
- Use `__all__` to control exports
- Run tests after each file split

**❌ DON'T:**
- Change function signatures
- Remove existing exports
- Break import paths without re-exports
- Combine unrelated functionality

---

## Files to Modify/Create

**Create:**
- [ ] `src/sdp/core/workstream/__init__.py`
- [ ] `src/sdp/core/workstream/models.py`
- [ ] `src/sdp/core/workstream/parser.py`
- [ ] `src/sdp/core/workstream/validator.py`
- [ ] `src/sdp/validators/ws_completion/__init__.py`
- [ ] `src/sdp/validators/ws_completion/checkers.py`
- [ ] `src/sdp/validators/ws_completion/reporter.py`
- [ ] `src/sdp/validators/capability_tier/__init__.py`
- [ ] `src/sdp/validators/capability_tier/models.py`
- [ ] `src/sdp/validators/capability_tier/extractors.py`
- [ ] `src/sdp/validators/capability_tier/validators.py`
- [ ] ... (similar for other splits)

**Modify:**
- [ ] `src/sdp/core/workstream.py` → re-exports only
- [ ] `src/sdp/validators/ws_completion.py` → re-exports only
- [ ] All files importing from split modules

---

## Test Plan

### Unit Tests
- [ ] All existing tests in `tests/unit/test_workstream*.py` pass
- [ ] All existing tests in `tests/unit/test_validators*.py` pass
- [ ] All existing tests in `tests/unit/test_cli*.py` pass

### Integration Tests
- [ ] `sdp doctor` works
- [ ] `sdp guard check` works
- [ ] `python -c "from sdp.core.workstream import Workstream"` works

### Regression
- [ ] `pytest tests/ --tb=short` — all green
- [ ] `mypy src/sdp --strict` — no errors

---

## Review Checklist

**Code Quality:**
- [ ] All split files <200 LOC
- [ ] Re-exports maintain backward compatibility
- [ ] Deprecation warnings added for old imports
- [ ] No circular imports

**Testing:**
- [ ] All existing tests pass
- [ ] No new tests needed (Phase 3)

---

**Version:** 1.0  
**Created:** 2026-01-31
