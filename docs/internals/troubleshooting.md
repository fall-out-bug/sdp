# SDP Troubleshooting Guide

This guide helps you diagnose and resolve common issues when working with SDP (Spec-Driven Protocol).

## Table of Contents

- [Error Categories](#error-categories)
- [Common Errors](#common-errors)
- [Quality Gate Failures](#quality-gate-failures)
- [Test Failures](#test-failures)
- [Build Issues](#build-issues)
- [Hook Failures](#hook-failures)
- [Dependency Problems](#dependency-problems)
- [Configuration Errors](#configuration-errors)

---

## Error Categories

SDP errors are organized into categories for easier troubleshooting:

| Category | Description | Common Examples |
|----------|-------------|-----------------|
| **validation** | File format and structure validation | Missing WS sections, invalid frontmatter |
| **build** | Pre/post-build validation failures | Missing Goal, failed acceptance criteria |
| **test** | Test execution failures | Failing unit tests, low coverage |
| **configuration** | Config file issues | Invalid TOML, missing keys |
| **dependency** | Workstream dependency issues | Dependency not completed |
| **hook** | Git/build hook failures | Pre-commit, post-build failures |
| **artifact** | Output validation failures | Code quality, documentation |
| **beads** | Beads task management issues | Task not found, server down |
| **coverage** | Test coverage issues | Coverage below threshold |

---

## Common Errors

### BeadsNotFoundError

**Symptom:** `Beads task 'TASK-XXX' not found`

**Causes:**
- Beads server not running
- Incorrect task ID
- BEADS_HOME not configured

**Solutions:**

1. Check Beads status:
   ```bash
   beads status
   ```

2. List available tasks:
   ```bash
   beads tasks
   ```

3. Verify BEADS_HOME:
   ```bash
   echo $BEADS_HOME
   ```

4. Start Beads server if needed:
   ```bash
   beads server start
   ```

**Documentation:** [Beads Integration](./beads.md)

---

### CoverageTooLowError

**Symptom:** `Coverage 65.5% is below required 80.0%`

**Causes:**
- Insufficient test coverage
- Missing edge case tests
- Untested utility functions

**Solutions:**

1. Run coverage with detailed report:
   ```bash
   pytest --cov=module.name --cov-report=term-missing
   ```

2. Identify uncovered lines:
   ```bash
   pytest --cov=module.name --cov-report=html
   open htmlcov/index.html
   ```

3. Add tests for missing coverage:
   ```python
   def test_uncovered_edge_case():
       """Test the edge case that's currently uncovered."""
       result = module.function(edge_case_input)
       assert result == expected_output
   ```

4. Verify coverage threshold:
   ```bash
   pytest --cov=module.name --cov-fail-under=80
   ```

**Documentation:** [Quality Gates](./quality-gates.md#coverage)

---

### WorkstreamValidationError

**Symptom:** `Workstream 'WS-001-01' validation failed with 3 error(s)`

**Causes:**
- Missing required sections
- Invalid frontmatter format
- Incorrect workstream ID format

**Solutions:**

1. Review validation errors:
   ```bash
   sdp core validate-ws docs/workstreams/backlog/WS-001-01.md
   ```

2. Check required sections:
   - Goal section (### 🎯 Goal or ### 🎯 Цель)
   - Acceptance Criteria
   - Dependencies (if any)

3. Verify frontmatter format:
   ```yaml
   ---
   ws_id: WS-001-01
   feature: F01
   title: "Workstream Title"
   status: backlog
   size: MEDIUM
   ---
   ```

4. Compare with template:
   ```bash
   cat docs/workstreams/template.md
   ```

**Documentation:** [Workstream Format](./workstreams.md#format)

---

### DependencyNotFoundError

**Symptom:** `Dependency 'WS-001-01' not found for workstream 'WS-001-02'`

**Causes:**
- Dependency workstream doesn't exist
- Dependency not yet completed
- Incorrect dependency ID

**Solutions:**

1. Check INDEX.md for dependency status:
   ```bash
   grep "WS-001-01" docs/workstreams/INDEX.md
   ```

2. Verify dependency exists:
   ```bash
   ls docs/workstorms/backlog/WS-001-01*.md
   ls docs/workstreams/completed/WS-001-01*.md
   ```

3. Update frontmatter if ID is wrong:
   ```yaml
   ---
   dependency: "WS-001-01"  # Correct ID
   ---
   ```

4. Mark dependency as completed if done:
   ```bash
   mv docs/workstreams/active/WS-001-01.md docs/workstreams/completed/
   ```

**Documentation:** [Dependencies](./workstreams.md#dependencies)

---

## Quality Gate Failures

### File Size Violations

**Symptom:** `File too large: module.py (250 lines)`

**Solution:**

1. Split large files into smaller modules:
   ```python
   # Before: module.py (250 lines)
   # After:
   # - module.py (main logic, 150 lines)
   # - module_utils.py (helpers, 50 lines)
   # - module_types.py (types, 50 lines)
   ```

2. Extract related functionality:
   ```python
   # module.py
   from module_utils import helper_function
   from module_types import CustomType
   ```

3. Verify file sizes:
   ```bash
   find src -name "*.py" -exec wc -l {} \; | awk '$1 > 200'
   ```

**Documentation:** [Quality Gates](./quality-gates.md#file-size)

---

### Complexity Violations

**Symptom:** `Cyclomatic complexity too high: function_name (CC=15)`

**Solution:**

1. Refactor complex functions:
   ```python
   # Before: CC=15
   def complex_function(data):
       # ... 15 decision points

   # After: CC=5 each
   def complex_function(data):
       preprocessed = preprocess(data)
       validated = validate(preprocessed)
       return process(validated)
   ```

2. Extract helper functions:
   ```python
   def preprocess(data):
       """Extract preprocessing logic."""
       ...

   def validate(data):
       """Extract validation logic."""
       ...
   ```

3. Check complexity:
   ```bash
   radon cc src/sdp/module.py -a
   ```

**Documentation:** [Quality Gates](./quality-gates.md#complexity)

---

## Test Failures

### Test Execution Failed

**Symptom:** `Tests failed: 8/10 passed, 2 failed`

**Solution:**

1. Run tests with verbose output:
   ```bash
   pytest tests/unit/test_module.py -v
   ```

2. Run only failing tests:
   ```bash
   pytest tests/unit/test_module.py::test_failing_test -v
   ```

3. Check for regressions:
   ```bash
   git diff HEAD~1 | grep "def test_"
   ```

4. Fix failing tests:
   ```python
   def test_failing_test():
       """Test that was failing."""
       # Fix implementation or test assertion
       result = module.function(input_data)
       assert result == expected_output
   ```

5. Re-run until all pass:
   ```bash
   pytest tests/unit/test_module.py --tb=short
   ```

**Documentation:** [Testing](./testing.md)

---

## Build Issues

### Pre-build Validation Failed

**Symptom:** `Build validation failed: Goal section (pre-build)`

**Solution:**

1. Check what's missing:
   ```bash
   ./hooks/pre-build.sh WS-001-01
   ```

2. Add missing Goal section:
   ```markdown
   ### 🎯 Goal (Цель)
   Implement feature X to solve problem Y.
   ```

3. Verify Acceptance Criteria:
   ```markdown
   ## Acceptance Criteria
   - [ ] Criterion 1
   - [ ] Criterion 2
   ```

4. Re-run pre-build check:
   ```bash
   ./hooks/pre-build.sh WS-001-01
   ```

**Documentation:** [Build Process](./building.md)

---

### Post-build Validation Failed

**Symptom:** `Build validation failed: Coverage check (post-build)`

**Solution:**

1. Run coverage manually:
   ```bash
   pytest tests/unit/test_module.py --cov=module --cov-report=term-missing
   ```

2. Add missing tests:
   ```python
   def test_uncovered_function():
       """Test previously uncovered function."""
       result = module.uncovered_function()
       assert result is not None
   ```

3. Verify threshold:
   ```bash
   pytest tests/unit/test_module.py --cov=module --cov-fail-under=80
   ```

4. Re-run post-build check:
   ```bash
   ./hooks/post-build.sh WS-001-01 module.name
   ```

**Documentation:** [Build Process](./building.md#post-build)

---

## Hook Failures

### Pre-commit Hook Failed

**Symptom:** `Hook 'pre-commit' failed during pre-commit (exit code: 1)`

**Solution:**

1. Review hook output:
   ```bash
   git commit -m "message"
   # Hook will run and show output
   ```

2. Identify failing check:
   - Time estimates in WS files?
   - Tech debt markers?
   - Bare except clauses?
   - Large files?

3. Fix the issue:
   ```bash
   # Example: Remove time estimate
   git diff --cached
   # Edit file to remove "2 hours", "1 day", etc.
   ```

4. Stage fixed files:
   ```bash
   git add file.py
   ```

5. Commit again:
   ```bash
   git commit -m "message"
   ```

**Bypass (not recommended):**
```bash
SKIP_CHECK=1 git commit -m "message"
```

**Documentation:** [Git Hooks](./hooks.md)

---

### Post-build Hook Failed

**Symptom:** `Hook 'post-build' failed during post-build (exit code: 1)`

**Solution:**

1. Run post-build manually:
   ```bash
   ./hooks/post-build.sh WS-001-01
   ```

2. Fix failing checks:
   - Regression tests passing?
   - TODO/FIXME markers removed?
   - Execution Report appended?

3. Re-run post-build:
   ```bash
   ./hooks/post-build.sh WS-001-01
   ```

**Documentation:** [Build Hooks](./building.md#hooks)

---

## Dependency Problems

### Circular Dependency

**Symptom:** `CircularDependencyError: WS-001-01 → WS-001-02 → WS-001-01`

**Solution:**

1. Visualize dependency graph:
   ```bash
   sdp design graph F01
   ```

2. Break the cycle:
   - Extract common functionality to new workstream
   - Reorder workstreams
   - Make one workstream independent

3. Update dependencies:
   ```yaml
   # WS-001-01
   dependency: "Independent"  # Remove circular dep

   # WS-001-02
   dependency: "WS-001-01"  # Keep this one
   ```

**Documentation:** [Dependencies](./workstreams.md#circular-dependencies)

---

### Missing Dependency

**Symptom:** `MissingDependencyError: WS-001-02 depends on WS-001-01 which doesn't exist`

**Solution:**

1. Create missing workstream:
   ```bash
   cp docs/workstreams/template.md docs/workstreams/backlog/WS-001-01.md
   ```

2. Or remove dependency if not needed:
   ```yaml
   # WS-001-02
   dependency: "Independent"  # Was: "WS-001-01"
   ```

3. Re-run validation:
   ```bash
   ./hooks/pre-build.sh WS-001-02
   ```

**Documentation:** [Dependencies](./workstreams.md#creating)

---

## Configuration Errors

### Invalid Configuration

**Symptom:** `ConfigValidationError: quality-gate.toml validation failed`

**Solution:**

1. Validate configuration:
   ```bash
   python -m sdp.quality.config validate quality-gate.toml
   ```

2. Check TOML syntax:
   ```bash
   python -c "import tomllib; tomllib.load(open('quality-gate.toml'))"
   ```

3. Fix configuration errors:
   ```toml
   # quality-gate.toml

   [file_size]
   max_lines = 200  # Valid integer

   [coverage]
   min_percentage = 80.0  # Valid float
   ```

4. Reload configuration:
   ```bash
   python -m sdp.quality.config load quality-gate.toml
   ```

**Documentation:** [Configuration](./configuration.md)

---

### Missing Configuration Keys

**Symptom:** `ConfigurationError: Missing required keys: ['max_file_size', 'min_coverage']`

**Solution:**

1. Identify missing keys:
   ```bash
   python -m sdp.quality.config check quality-gate.toml
   ```

2. Add missing keys:
   ```toml
   # quality-gate.toml

   [file_size]
   max_lines = 200  # Required

   [coverage]
   min_percentage = 80.0  # Required
   ```

3. Verify configuration:
   ```bash
   python -m sdp.quality.config validate quality-gate.toml
   ```

**Documentation:** [Configuration Schema](./configuration.md#schema)

---

## Getting Help

If you're still stuck after trying these solutions:

1. **Check Documentation:**
   - [README](../README.md)
   - [PROTOCOL](../PROTOCOL.md)
   - [Quality Gates](./quality-gates.md)
   - [Workstreams](./workstreams.md)

2. **Search Issues:**
   - [GitHub Issues](https://github.com/your-org/sdp/issues)

3. **Ask for Help:**
   - Create a GitHub issue with:
     - Error message (full output)
     - Steps to reproduce
     - Your environment (OS, Python version)
     - What you've already tried

4. **Community:**
   - Join discussion forums
   - Check existing troubleshooting threads

---

## Quick Reference

### Common Commands

```bash
# Validate workstream
sdp core parse-ws docs/workstreams/backlog/WS-001-01.md

# Run tests
pytest tests/unit/ -v

# Check coverage
pytest --cov=module.name --cov-report=term-missing

# Run quality gates
python -m sdp.quality.validator

# Check dependencies
sdp design graph F01

# Pre-build check
./hooks/pre-build.sh WS-001-01

# Post-build check
./hooks/post-build.sh WS-001-01
```

### Error Categories

```python
from sdp.errors import (
    SDPError,
    BeadsNotFoundError,
    CoverageTooLowError,
    WorkstreamValidationError,
    # ... etc
)

# Format error for terminal
from sdp.errors import format_error_for_terminal
print(format_error_for_terminal(error))
```

---

**Last Updated:** 2025-01-29
**SDP Version:** 0.3.0
