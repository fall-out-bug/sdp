---
ws_id: 00-020-04
feature: F020
status: completed
size: SMALL
project_id: 00
github_issue: null
assignee: null
depends_on:
  - 00-020-01
---

## WS-00-020-04: Fix except Exception: pass in common.py (Issue 006)

### 🎯 Goal

**What must WORK after completing this WS:**
- No `except Exception: pass` without logging in `src/sdp/hooks/common.py`
- Quality gate "except Exception only with logging" passes for hooks module

**Acceptance Criteria:**
- [ ] AC1: Replace `except Exception: pass` at lines 41-42 (find_project_root) with specific exception or logging
- [ ] AC2: Replace `except Exception: pass` at lines 79-80 (find_workstream_dir) with specific exception or logging
- [ ] AC3: All tests pass, mypy --strict passes, ruff clean
- [ ] AC4: Hooks coverage remains ≥71% (no regression)

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Source:** Issue 006, F020 Review (2026-01-31)

**Problem:** Two instances of `except Exception: pass` without logging violate quality gate.

**Fix options:**
1. Add logging: `except Exception as e: logger.debug("...", e)`
2. Catch specific: `except (tomllib.TOMLDecodeError, OSError): pass`

### Dependency

- 00-020-01: Extract Git Hooks to Python (provides common.py)

### Input Files

- `src/sdp/hooks/common.py`
- `tests/unit/hooks/test_common.py`

### Steps

1. Replace except at lines 41-42 with `except (tomllib.TOMLDecodeError, OSError): pass`
2. Replace except at lines 79-80 with `except (tomllib.TOMLDecodeError, OSError): pass`
3. Add tests that verify fallback when TOML parse fails (optional, for coverage)
4. Run quality check: pytest, mypy, ruff

---

## Execution Report

**Completed:** 2026-01-31

### Completed Tasks

- [x] AC1: Replaced `except Exception: pass` at lines 41-42 with `except (tomllib.TOMLDecodeError, OSError): pass`
- [x] AC2: Replaced `except Exception: pass` at lines 79-80 with `except (tomllib.TOMLDecodeError, OSError): pass`
- [x] AC3: Added tests `test_find_project_root_skips_invalid_toml` and `test_find_workstream_dir_falls_through_on_invalid_toml`
- [x] AC4: All tests pass, mypy --strict passes, ruff clean

### Test Results

```
uv run pytest tests/unit/hooks/test_common.py -v
# 13 passed

uv run pytest --cov=src/sdp/hooks --cov-report=term-missing
# common.py: 93% coverage, hooks total: 92%
```

### Decisions

- **Specific exceptions over logging:** Used `except (tomllib.TOMLDecodeError, OSError)` instead of `except Exception` with logging. Avoids logger dependency in utility module; explicit about expected failure modes.

---

## Review Results (2026-01-31)

**Verdict:** APPROVED (F020 feature review)  
**Report:** [2026-01-31-F020-review.md](../../reports/2026-01-31-F020-review.md)
