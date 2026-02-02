# Python Quick Start

SDP workflow for Python projects with pytest, mypy, and ruff.

## Prerequisites

```bash
# Python 3.10+
python --version

# Install SDP dependencies
pip install pytest pytest-cov mypy ruff
```

## Project Structure

```
my-project/
├── src/
│   └── myapp/
│       ├── __init__.py
│       └── service.py
├── tests/
│   ├── __init__.py
│   └── test_service.py
├── pyproject.toml
└── .claude/
    └── skills/
```

## Workflow

### 1. Initialize Project

```bash
# Create feature
@feature "Add user authentication"
```

Claude will interview you about:
- Mission and users
- Technical approach
- Success criteria
- Tradeoffs

### 2. Plan Workstreams

```bash
@design feature-auth
```

Claude will:
- Explore codebase
- Design workstream decomposition
- Define dependencies
- Request approval

### 3. Execute Workstream

```bash
@build 00-001-01
```

SDP will:
1. Detect Python project (pyproject.toml found)
2. Run tests: `pytest tests/ -v`
3. Run coverage: `pytest --cov=src/ --cov-fail-under=80`
4. Run type checking: `mypy src/ --strict`
5. Run linting: `ruff check src/`
6. AI validators check architecture, errors, complexity

### 4. Review Quality

```bash
@review feature-auth
```

SDP will run AI validators:
- `/coverage-validator` - Analyzes test coverage
- `/architecture-validator` - Checks Clean Architecture
- `/error-validator` - Finds bare except clauses
- `/complexity-validator` - Identifies complex code

## Quality Gates

### Test Coverage ≥80%

```bash
pytest tests/unit/ --cov=src/ --cov-fail-under=80
```

**Output:**
```
---------- coverage: platform linux, python 3.10 ----------
Name                          Stmts   Miss  Cover   Missing
---------------------------------------------------------
src/service.py                   50      5    90%    23-27
src/models.py                    30      2    93%    45, 67
---------------------------------------------------------
TOTAL                            80      7    91%    ✅ PASS
```

### Type Checking

```bash
mypy src/ --strict
```

**Output:**
```
Success: no issues found in 2 source files ✅ PASS
```

### Linting

```bash
ruff check src/
```

**Output:**
```
All checks passed! ✅ PASS
```

### File Size

```bash
find src/ -name "*.py" -exec wc -l {} + | awk '$1 > 200'
```

**Output:**
```
(no output = all files <200 LOC) ✅ PASS
```

## Example Workflow

```bash
# Start feature
@feature "Add user login"

# Plan workstreams
@design feature-login

# Execute first workstream
@build 00-001-01

# Expected output:
# ✓ Project type detected: Python (pyproject.toml)
# ✓ Running tests: pytest tests/ -v
# ✓ Coverage: 85% (≥80% required)
# ✓ Type checking: mypy src/ --strict
# ✓ Linting: ruff check src/
# ✓ AI validators: PASS
#
# Workstream 00-001-01 complete!

# Execute next workstream
@build 00-001-02

# Review all workstreams
@review feature-login

# Deploy
@deploy feature-login
```

## Common Issues

### Issue: Coverage <80%

**Solution:** Add more tests
```python
# tests/test_service.py
def test_user_creation():
    user = User("Alice")
    assert user.name == "Alice"

def test_user_invalid_name():
    with pytest.raises(ValueError):
        User("")  # Add this test
```

### Issue: mypy errors

**Solution:** Add type hints
```python
# Before (FAIL)
def process(data, count):
    return data * count

# After (PASS)
def process(data: str, count: int) -> str:
    return data * count
```

### Issue: ruff errors

**Solution:** Fix linting issues
```bash
ruff check src/ --fix
```

## Tips

1. **Use virtual environment:**
   ```bash
   python -m venv .venv
   source .venv/bin/activate
   pip install -r requirements.txt
   ```

2. **Run tests in watch mode:**
   ```bash
   pytest -f tests/
   ```

3. **Check coverage interactively:**
   ```bash
   pytest --cov=src/ --cov-report=html
   open htmlcov/index.html
   ```

4. **Pre-commit hooks:**
   ```bash
   # SDP hooks run automatically on git commit
   git commit -m "feat: add feature"
   # Hooks run pytest, mypy, ruff automatically
   ```

## Next Steps

- [Java Quick Start](../java/QUICKSTART.md)
- [Go Quick Start](../go/QUICKSTART.md)
- [Full Tutorial](../../TUTORIAL.md)
