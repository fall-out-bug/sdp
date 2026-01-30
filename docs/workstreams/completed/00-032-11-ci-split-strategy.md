---
ws_id: 00-032-11
feature: F032
status: completed
completed: "2026-01-30"
size: SMALL
project_id: 00
github_issue: null
assignee: Claude
depends_on: []
---

## 00-032-11: CI Split Strategy

### 🎯 Goal

**What must WORK after completing this WS:**
- ADR документ определяет какие checks блокируют, какие предупреждают
- Config файл `ci-gates.toml` с настройками

**Acceptance Criteria:**
- [x] AC1: ADR `docs/adr/008-ci-split-strategy.md` created
- [x] AC2: Critical checks defined (coverage, mypy, tests)
- [x] AC3: Warning checks defined (file size, complexity)
- [x] AC4: `ci-gates.toml` created with configuration

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: Текущий CI использует `continue-on-error: true` везде. Ничего не блокирует.

**Solution**: Разделить checks на critical (block) и warning (comment).

### Dependencies

None (independent)

### Steps

1. **Create ADR**

   ```markdown
   # docs/adr/008-ci-split-strategy.md
   # ADR-008: CI Split Strategy
   
   ## Status
   Accepted
   
   ## Context
   
   Current CI workflow uses `continue-on-error: true` for all checks.
   This means:
   - PRs can merge with failing tests
   - Agents learn they can ignore CI failures
   - Quality degrades over time
   
   We need to balance:
   - Strictness (blocking bad code)
   - Developer experience (not blocking on minor issues)
   
   ## Decision
   
   Split checks into two categories:
   
   ### Critical (Block PR)
   
   | Check | Threshold | Reason |
   |-------|-----------|--------|
   | Tests | All pass | Core functionality |
   | Coverage | ≥80% | Prevent untested code |
   | mypy strict | No errors | Type safety |
   | ruff errors | No errors | Code quality |
   
   These run without `continue-on-error`. Failure blocks merge.
   
   ### Warning (Comment Only)
   
   | Check | Threshold | Reason |
   |-------|-----------|--------|
   | File size | <200 LOC | AI readability |
   | Complexity | CC <10 | Maintainability |
   | ruff warnings | Report | Style suggestions |
   
   These run with `continue-on-error: true`. Post comment but don't block.
   
   ## Implementation
   
   Two separate workflows:
   - `ci-critical.yml` — Required status check
   - `ci-warnings.yml` — Informational
   
   ## Consequences
   
   ### Positive
   - Bad code can't merge
   - Clear distinction critical vs nice-to-have
   - Agents learn which rules are hard requirements
   
   ### Negative
   - More CI configuration
   - Some PRs will be blocked (intended)
   - Need branch protection setup
   ```

2. **Create config file**

   ```toml
   # ci-gates.toml
   # CI Quality Gates Configuration
   
   [critical]
   # These BLOCK PR merge
   
   tests_pass = true
   coverage_threshold = 80
   mypy_strict = true
   ruff_errors = true
   
   [warning]
   # These COMMENT but don't block
   
   file_size_loc = 200
   complexity_cc = 10
   ruff_warnings = true
   
   [thresholds]
   # Numeric thresholds
   
   coverage_min = 80
   coverage_warn = 85
   complexity_max = 10
   complexity_warn = 7
   file_size_max = 200
   file_size_warn = 150
   ```

### Output Files

- `docs/adr/008-ci-split-strategy.md`
- `ci-gates.toml`

### Completion Criteria

```bash
# ADR exists
test -f docs/adr/008-ci-split-strategy.md

# Config exists
test -f ci-gates.toml

# Config is valid TOML
python -c "import tomllib; tomllib.load(open('ci-gates.toml', 'rb'))"
```

---

## Execution Report

**Executed by:** Claude (AI Agent)  
**Date:** 2026-01-30

### Goal Status
- [x] AC1-AC4 — ✅

**Goal Achieved:** YES

### Implementation Details

Created two key files:

1. **ADR-008** (`docs/adr/008-ci-split-strategy.md`)
   - Documents decision to split CI into critical/warning
   - Lists critical checks: tests, coverage ≥80%, mypy strict, ruff errors
   - Lists warning checks: file size <200 LOC, complexity CC <10, ruff warnings
   - Explains consequences and rationale

2. **Config** (`ci-gates.toml`)
   - Structured TOML configuration
   - Three sections: critical, warning, thresholds
   - All numeric thresholds parameterized
   - Validated with tomllib (valid TOML)

### Verification

```bash
$ python3 -c "import tomllib; tomllib.load(open('ci-gates.toml', 'rb'))"
✅ Valid TOML
Critical checks: ['tests_pass', 'coverage_threshold', 'mypy_strict', 'ruff_errors']
Warning checks: ['file_size_loc', 'complexity_cc', 'ruff_warnings']
Thresholds: ['coverage_min', 'coverage_warn', 'complexity_max', 'complexity_warn', 'file_size_max', 'file_size_warn']
```

### Notes

Foundation is ready for workflow implementation in WS 00-032-12 and 00-032-13.
