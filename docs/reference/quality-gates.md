# SDP Quality Gates

**Binary pass/fail criteria** for workstream completion and quality assurance.

---

## Overview

Quality gates are **MUST PASS** criteria. Each gate has:
- ✅ **PASS** - Criteria met, workstream approved
- ❌ **FAIL** - Criteria not met, must fix before completion
- 🔧 **Measurement** - How to check the criterion
- 📋 **Example** - What pass/fail looks like

---

## 1. AI-Readiness Gate

### Criteria

| # | Rule | Limit | Measurement |
|---|------|-------|-------------|
| 1.1 | Max lines per file | < 200 LOC | `find src/ -name "*.py" -exec wc -l {} + | awk '$1 > 200'` |
| 1.2 | Max cyclomatic complexity | < 10 | `ruff check src/ --select=C901` |
| 1.3 | Type hints coverage | 100% | `mypy src/ --strict` (no errors) |
| 1.4 | Function length | < 50 lines | Manual review or `radon cc src/ -a -s` |

### PASS Criteria

```bash
# All checks must return empty (no violations)
find src/ -name "*.py" -exec wc -l {} + | awk '$1 > 200'  # No output
ruff check src/ --select=C901  # No C901 errors
mypy src/ --strict  # Success: 0 errors
```

### FAIL Examples

❌ **File too long:**
```python
# service.py (250 lines) - FAIL
def complex_function():
    # 150 lines of logic
    pass
```

✅ **Correct:**
```python
# service.py (180 lines) - PASS
def core_function():
    # 30 lines
    pass

def helper_function():
    # 20 lines
    pass
```

❌ **Missing type hints:**
```python
def process(data):  # FAIL - no type hints
    return data.upper()
```

✅ **Correct:**
```python
def process(data: str) -> str:  # PASS - full type hints
    return data.upper()
```

---

## 2. Clean Architecture Gate

### Criteria

| # | Rule | Description | Measurement |
|---|------|-------------|-------------|
| 2.1 | Domain layer isolation | No dependencies on other layers | `grep -r "from.*infrastructure" src/domain/` (must be empty) |
| 2.2 | Application layer | Can only depend on Domain | `grep -r "from.*presentation" src/application/` (must be empty) |
| 2.3 | Infrastructure layer | Can depend on Domain + Application | Manual review of imports |
| 2.4 | No circular dependencies | Across any layers | `radon cc src/ -s` (no cycles) |

### PASS Criteria

```bash
# Domain layer must have NO external dependencies
grep -r "from.*infrastructure" src/domain/  # No output
grep -r "from.*application" src/domain/  # No output

# Application layer must NOT import from presentation
grep -r "from.*presentation" src/application/  # No output
```

### FAIL Examples

❌ **Domain layer violation:**
```python
# src/domain/entities/user.py
from src.infrastructure.persistence import Database  # FAIL - domain imports infrastructure

class UserEntity:
    def save(self):
        db = Database()
```

✅ **Correct:**
```python
# src/domain/entities/user.py
class UserEntity:
    def __init__(self, name: str, email: str):
        self.name = name
        self.email = email
    # PASS - domain has no infrastructure dependencies
```

❌ **Circular dependency:**
```python
# src/domain/entity.py
from src.application.service import Service  # FAIL - domain imports application

# src/application/service.py
from src.domain.entity import Entity  # Creates circular dependency
```

✅ **Correct:**
```python
# src/domain/entity.py
class Entity:
    pass  # PASS - no imports to application

# src/application/service.py
from src.domain.entity import Entity  # OK - application can import domain
```

---

## 3. Error Handling Gate

### Criteria

| # | Rule | Pattern | Measurement |
|---|------|---------|-------------|
| 3.1 | No bare except | Forbidden | `grep -rn "except:" src/` (must be empty) |
| 3.2 | No broad exceptions | `except Exception` only with logging | `grep -rn "except Exception" src/ \| grep -v "exc_info"` |
| 3.3 | Explicit error types | Must catch specific exceptions | Manual review |
| 3.4 | All errors logged | With context | Manual review |

### PASS Criteria

```bash
# No bare except clauses
grep -rn "except:" src/  # No output

# No bare Exception catches
grep -rn "except Exception" src/ | grep -v "exc_info"  # No output
```

### FAIL Examples

❌ **Bare except (SWALLOWS ALL ERRORS):**
```python
try:
    risky_operation()
except:  # FAIL - catches everything, including KeyboardInterrupt
    pass  # Silent failure - impossible to debug
```

✅ **Correct:**
```python
try:
    risky_operation()
except SpecificError as e:
    logger.error(f"Failed: {e}")  # PASS - explicit error, logged
    raise  # Re-raise or handle appropriately
```

❌ **Broad Exception without logging:**
```python
try:
    operation()
except Exception:  # FAIL - too broad, no logging
    pass
```

✅ **Correct:**
```python
try:
    operation()
except ValueError as e:
    logger.warning(f"Invalid value: {e}")  # PASS - specific error
    raise
except DatabaseError as e:
    logger.error(f"Database error: {e}")  # PASS - specific error
    raise
```

---

## 4. Test Coverage Gate

### Criteria

| # | Rule | Minimum | Measurement |
|---|------|---------|-------------|
| 4.1 | Line coverage | ≥ 80% | `pytest --cov=src/ --cov-fail-under=80` |
| 4.2 | Branch coverage | ≥ 70% | `pytest --cov=src/ --cov-branch` |
| 4.3 | All public APIs tested | 100% | Manual review |
| 4.4 | No skipped tests | 0 skipped | `pytest -v` (no "SKIPPED") |

### PASS Criteria

```bash
# Coverage must be ≥ 80%
pytest --cov=src/ --cov-fail-under=80 --cov-report=term-missing

# Output must show:
# TOTAL 80% (or higher)
# FAIL < 80% would exit with error code 1
```

### FAIL Examples

❌ **Insufficient coverage:**
```
Name                      Stmts   Miss  Cover   Missing
-------------------------------------------------------
src/module.py                50     15    70%   23-37
FAIL - Required coverage 80% but achieved 70%
```

✅ **Correct:**
```
Name                      Stmts   Miss  Cover   Missing
-------------------------------------------------------
src/module.py                50      8    84%   45-47, 89
TOTAL                       50      8    84%
PASS ✓
```

❌ **Skipped test (flaky test avoidance):**
```python
def test_feature():
    # pytest.skip("Skipping for now")  # FAIL - test skipped
    assert True
```

✅ **Correct:**
```python
@pytest.mark.parametrize("input,expected", [
    ("valid", "expected"),
    ("invalid", "error"),
])
def test_feature(input, expected):  # PASS - no skipping
    assert process(input) == expected
```

---

## 5. Code Quality Gate

### Criteria

| # | Rule | Standard | Measurement |
|---|------|----------|-------------|
| 5.1 | No linting errors | All ruff rules pass | `ruff check src/` (exit 0) |
| 5.2 | No formatting issues | All files formatted | `ruff format --check src/` |
| 5.3 | No TODOs without followup | Each TODO has WS ID | `grep -rn "TODO" src/ \| grep -v "WS-"` |
| 5.4 | No commented code | Clean codebase | Manual review |

### PASS Criteria

```bash
# All ruff checks pass
ruff check src/  # Exit 0, no output

# All files formatted
ruff format --check src/  # Exit 0, no "would reformat" messages

# No orphaned TODOs
grep -rn "TODO" src/ | grep -v "WS-"  # No output
```

### FAIL Examples

❌ **Linting error:**
```python
import os, sys  # FAIL - multiple imports on one line (E401)

def function( ):  # FAIL - extra whitespace (E211)
    pass
```

✅ **Correct:**
```python
import os
import sys  # PASS - separate imports

def function():  # PASS - no extra whitespace
    pass
```

❌ **TODO without followup:**
```python
# TODO: Add error handling  # FAIL - no WS ID
def process():
    pass
```

✅ **Correct:**
```python
# TODO: Add error handling (WS-123-45)  # PASS - tracked in workstream
def process():
    pass
```

---

## 6. Documentation Gate

### Criteria

| # | Rule | Requirement | Measurement |
|---|------|-------------|-------------|
| 6.1 | All public functions documented | Docstrings | `interrogate src/ -vvv` |
| 6.2 | Module docstrings | All modules | Manual review |
| 6.3 | README examples | Run without errors | Manual execution |
| 6.4 | Type hints in docs | sphinx-autodoc | `make docs` (no warnings) |

### PASS Criteria

```bash
# Check docstring coverage
interrogate src/ -vvv -f 90  # Requires 90% coverage

# All public modules have docstrings
for f in src/**/*.py; do
    head -10 "$f" | grep -q '"""'  # PASS
done
```

### FAIL Examples

❌ **Missing docstring:**
```python
def process(data: str) -> str:  # FAIL - no docstring
    return data.upper()
```

✅ **Correct:**
```python
def process(data: str) -> str:
    """Process input data by converting to uppercase.

    Args:
        data: Input string to process

    Returns:
        Uppercase version of input

    Examples:
        >>> process("hello")
        "HELLO"
    """
    return data.upper()  # PASS - full docstring
```

---

## 7. Build Gate

### Criteria

| # | Rule | Requirement | Measurement |
|---|------|-------------|-------------|
| 7.1 | Clean install | No setup errors | `pip install -e .` (exit 0) |
| 7.2 | All tests pass | pytest exit 0 | `pytest -x` |
| 7.3 | No import errors | All modules importable | `python -c "import src"` |

### PASS Criteria

```bash
# Installation succeeds
pip install -e .  # Exit 0

# All tests pass
pytest -x -v  # Exit 0, no failures

# All modules importable
python -c "from src import module"  # Exit 0
```

---

## 8. Security Gate

### Criteria

| # | Rule | Pattern | Measurement |
|---|------|---------|-------------|
| 8.1 | No hardcoded secrets | No passwords/tokens | `grep -rn "password\|token\|api_key" src/ \| grep -v "example\|test"` |
| 8.2 | No SQL injection | Parameterized queries | Manual review |
| 8.3 | No eval/exec | Forbidden functions | `grep -rn "\beval\b\|bexec\b" src/` |
| 8.4 | Input validation | All user inputs validated | Manual review |

### PASS Criteria

```bash
# No hardcoded secrets
grep -rn "password\|token\|api_key" src/ | grep -v "example\|test"  # No output

# No eval/exec
grep -rn "\beval\b\|bexec\b" src/  # No output
```

### FAIL Examples

❌ **Hardcoded secret:**
```python
API_KEY = "sk-1234567890"  # FAIL - hardcoded secret
```

✅ **Correct:**
```python
import os
API_KEY = os.getenv("API_KEY")  # PASS - from environment
```

❌ **SQL injection risk:**
```python
query = f"SELECT * FROM users WHERE name = '{user_input}'"  # FAIL
cursor.execute(query)  # Vulnerable to injection
```

✅ **Correct:**
```python
query = "SELECT * FROM users WHERE name = %s"
cursor.execute(query, (user_input,))  # PASS - parameterized
```

---

## Running All Gates

### Quick Check

```bash
#!/bin/bash
# Run all quality gates

set -e  # Fail on any error

echo "🔍 Running quality gates..."

# AI-Readiness
echo "📏 AI-Readiness..."
find src/ -name "*.py" -exec wc -l {} + | awk '$1 > 200' && exit 1 || true
ruff check src/ --select=C901
mypy src/ --strict

# Clean Architecture
echo "🏗️  Clean Architecture..."
grep -r "from.*infrastructure" src/domain/ && exit 1 || true
grep -r "from.*presentation" src/application/ && exit 1 || true

# Error Handling
echo "⚠️  Error Handling..."
grep -rn "except:" src/ && exit 1 || true

# Test Coverage
echo "🧪 Test Coverage..."
pytest --cov=src/ --cov-fail-under=80 --cov-report=term-missing

# Code Quality
echo "✨ Code Quality..."
ruff check src/
ruff format --check src/

echo "✅ All gates passed!"
```

### Detailed Report

```bash
# Run with coverage report
pytest --cov=src/ --cov-report=html
open htmlcov/index.html  # View detailed coverage

# Run with complexity report
radon cc src/ -a -s
radon mi src/  # Maintainability index

# Full linting report
ruff check src/ --output-format=github
```

---

## Troubleshooting

### Gate Fails: What Now?

1. **Check the measurement command output**
   ```bash
   # If coverage gate fails:
   pytest --cov=src/ --cov-report=term-missing
   # See which lines are missing coverage
   ```

2. **Fix specific violations**
   ```bash
   # If ruff fails:
   ruff check src/ --fix  # Auto-fix where possible
   ```

3. **Re-run the gate**
   ```bash
   # Re-test after fixes
   pytest --cov=src/ --cov-fail-under=80
   ```

### Common Issues

| Issue | Cause | Solution |
|-------|-------|----------|
| Coverage < 80% | Missing test cases | Add tests for uncovered lines |
| File > 200 LOC | Too much logic | Split into smaller modules/functions |
| `mypy --strict` fails | Missing type hints | Add type hints to all functions |
| Import error | Circular dependency | Refactor to remove cycle |
| Test failures | Broken logic | Fix implementation or test |

---

## Local vs CI Behavior (AC4)

### Local Development

**Guard Mode:** `standard` (default)

**Behavior:**
- Violations displayed to user
- `error` severity: Blocks command
- `warning` severity: Displayed, continues
- `info` severity: Logged only

**Exit Codes:**
- `0`: All checks passed
- `1`: Error-level violations found
- `2`: Configuration or validation error

**Example:**
```bash
$ sdp guard check src/module.go
✅ max-file-loc: PASS (145 lines)
⚠️  max-function-length: WARN (process_func: 52 lines)
❌ coverage-threshold: ERROR (65% < 80%)
# Exit code: 1
```

### CI Environment

**Guard Mode:** `strict` (recommended for CI)

**Behavior:**
- All violations treated as errors
- Warnings cause build failure
- JSON output for parsing

**Exit Codes:**
- `0`: All checks passed
- `1`: Any violation found (error or warning)
- `2`: Configuration or validation error
- `3`: Diff-range detection failed

**Example:**
```bash
# In CI workflow
- name: Run guard checks
  run: sdp guard check --mode=strict --output=json > results.json

- name: Fail on violations
  run: |
    ERROR_COUNT=$(jq '[.[] | select(.severity == "error")] | length' results.json)
    if [ "$ERROR_COUNT" -gt 0 ]; then
      echo "❌ Found $ERROR_COUNT error(s)"
      exit 1
    fi
```

### CI Integration Pattern

**GitHub Actions workflow using same rules source as local (AC5):**

```yaml
name: Quality Gates

on: [push, pull_request]

jobs:
  guard-checks:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25.6'

      - name: Build SDP CLI
        run: |
          cd sdp-plugin
          go build -o sdp ./cmd/sdp

      - name: Run guard checks
        run: |
          # Uses .sdp/guard-rules.yml from repo
          # CI mode: strict, JSON output
          ./sdp guard check --mode=strict --output=json > guard-results.json

      - name: Display results
        if: always()
        run: |
          # Human-readable summary
          jq -r '.[] | "\(.severity): \(.rule_id) - \(.message)"' guard-results.json

      - name: Fail on errors
        run: |
          ERROR_COUNT=$(jq '[.[] | select(.severity == "error")] | length' guard-results.json)
          if [ "$ERROR_COUNT" -gt 0 ]; then
            echo "❌ Guard checks failed with $ERROR_COUNT error(s)"
            exit 1
          fi
```

### Exit Code Reference

| Code | Meaning | Local Behavior | CI Behavior |
|------|---------|----------------|-------------|
| 0 | Success | Proceed | Proceed |
| 1 | Violations | Error-level violations only | Any violation |
| 2 | Config Error | Invalid configuration | Invalid configuration |
| 3 | System Error | Runtime error | Runtime error |

---

## Continuous Integration

```yaml
name: Quality Gates

on: [push, pull_request]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-python@v4
        with:
          python-version: '3.11'
      - name: Install dependencies
        run: |
          pip install -e .
          pip install pytest pytest-cov ruff mypy
      - name: Run quality gates
        run: |
          pytest --cov=src/ --cov-fail-under=80
          ruff check src/
          mypy src/ --strict
```

---

## Version

**SDP Version:** 0.9.0
**Updated:** 2026-01-29
**Status:** Active (enforced in CI/CD)

---

## See Also

- [PROTOCOL.md](../PROTOCOL.md) - Full SDP specification
- [CODE_PATTERNS.md](../CODE_PATTERNS.md) - Code style examples
- [tests/](../tests/) - Test examples
