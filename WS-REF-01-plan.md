# WS-REF-01: Refactoring Plan

## Current State Analysis

### Critical Issues
1. **parser_python.py**: BROKEN - syntax error at line 194 ("return" outside function)
2. **28 files** exceed 200 LOC limit
3. **Overall coverage**: 44% (need 80%)
4. **mypy errors**: 5 errors in parser_python.py

### File Size Violations (>200 LOC)
- validators/capability_tier.py: 592 LOC
- cli.py: 459 LOC
- beads/client.py: 463 LOC
- And 25 other files

### Priority Order (P0 → P3)

**P0 - CRITICAL (Broken code):**
1. Fix parser_python.py syntax error
2. Fix mypy errors in parser_python.py

**P1 - HIGH (User-facing):**
3. Split cli.py (459 LOC) → multiple modules
4. Increase test coverage to 80%

**P2 - MEDIUM (Developer experience):**
5. Split other large files
6. Reduce complexity

**P3 - LOW (Nice to have):**
7. Documentation improvements
8. Performance optimizations

## Execution Strategy

Given the scope, I'll focus on P0-P1 to meet the core quality gates.

### Phase 1: Fix parser_python.py (P0)
- [ ] Fix syntax error (line 194)
- [ ] Fix all mypy errors
- [ ] Add type hints
- [ ] Write tests

### Phase 2: Split cli.py (P1)
- [ ] Create src/sdp/cli/workstream.py
- [ ] Create src/sdp/cli/quality.py
- [ ] Create src/sdp/cli/tier.py
- [ ] Update cli/__init__.py

### Phase 3: Increase Coverage (P1)
- [ ] Write tests for prd/ module
- [ ] Write tests for cli/ commands
- [ ] Write tests for beads/ module
- [ ] Target: 80% overall

## Quality Gate Checklist

- [ ] All files < 200 LOC (CLI split)
- [ ] mypy --strict passes (parser_python.py fixed)
- [ ] Test coverage ≥ 80%
- [ ] All tests passing
- [ ] No bare exceptions
- [ ] Type hints everywhere
