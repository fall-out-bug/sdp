# SDP Development Setup

Guide to setting up development environment for SDP.

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Testing](#testing)
- [Debugging](#debugging)
- [Common Issues](#common-issues)

---

## Prerequisites

### Required Software

- **Python 3.14+** - Primary language
- **Git** - Version control
- **Make** (optional) - Build automation

### Required Python Packages

```bash
# Core dependencies
pip install click pydantic pyyaml requests

# Development dependencies
pip install pytest pytest-cov mypy ruff
```

### Optional Software

- **Beads CLI** - Task tracking
- **Claude Code** - AI-IDE (recommended)
- **Telegram Desktop** - Notifications
- **GitHub CLI** - GitHub operations

---

## Installation

### 1. Clone Repository

```bash
git clone https://github.com/fall-out-bug/sdp.git
cd sdp
```

### 2. Create Virtual Environment

```bash
# Using venv
python3 -m venv .venv
source .venv/bin/activate  # On Windows: .venv\Scripts\activate

# Or using conda
conda create -n sdp python=3.14
conda activate sdp
```

### 3. Install SDP

```bash
# Development installation
pip install -e .

# With dev dependencies
pip install -e ".[dev]"
```

### 4. Verify Installation

```bash
# Check version
sdp --version

# Run health checks
sdp doctor
```

---

## Configuration

### Environment Variables

**Create `.env` file:**
```bash
# Python
PYTHONPATH=src

# Beads (optional)
BEADS_HOME=~/.beads
BEADS_SERVER=http://localhost:8080

# GitHub (optional)
GITHUB_TOKEN=ghp_your_token_here
GITHUB_REPO=owner/repo

# Telegram (optional)
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_CHAT_ID=your_chat_id
```

### Quality Gate Configuration

**File:** `quality-gate.toml`

```toml
[coverage]
minimum = 80

[complexity]
max_cc = 10

[file_size]
max_lines = 200

[type_hints]
strict_mode = true
```

### Claude Code Settings

**File:** `.claude/settings.json`

```json
{
  "skills": [
    "feature",
    "design",
    "build",
    "review",
    "deploy",
    "oneshot",
    "debug",
    "issue",
    "hotfix",
    "bugfix"
  ]
}
```

---

## Testing

### Run All Tests

```bash
# Full test suite
pytest

# With coverage
pytest --cov=src/sdp --cov-report=term-missing

# Verbose output
pytest -v
```

### Run Specific Tests

```bash
# Unit tests only
pytest tests/unit/

# Integration tests only
pytest tests/integration/

# Specific test file
pytest tests/unit/test_workstream.py

# Specific test
pytest tests/unit/test_workstream.py::test_parse_workstream
```

### Test Markers

```bash
# Fast tests (no external deps)
pytest -m fast

# Integration tests
pytest -m integration

# Exclude slow tests
pytest -m "not slow"
```

### Coverage Reports

```bash
# Terminal coverage
pytest --cov=src/sdp --cov-report=term-missing

# HTML coverage
pytest --cov=src/sdp --cov-report=html
open htmlcov/index.html

# JSON coverage (for CI)
pytest --cov=src/sdp --cov-report=json
```

---

## Debugging

### Type Checking

```bash
# Check all code
mypy src/sdp/ --strict

# Check specific file
mypy src/sdp/workstream.py --strict

# Check with error summary
mypy src/sdp/ --strict --error-summary
```

### Linting

```bash
# Check all code
ruff check src/sdp/

# Fix auto-fixable issues
ruff check src/sdp/ --fix

# Check specific file
ruff check src/sdp/workstream.py
```

### Debugging Tests

```bash
# Run with debugger
pytest --pdb

# Run specific test with debugger
pytest tests/unit/test_workstream.py::test_parse_workstream --pdb

# Drop into PDB on failure
pytest --pdb -x
```

### Debugging Skills

```bash
# Enable verbose output
export DEBUG=1
export VERBOSE=1

# Run skill with debug
@build WS-001-01 --verbose
```

---

## Development Workflow

### 1. Make Changes

```bash
# Create feature branch
git checkout -b feature/my-change

# Edit files
vim src/sdp/module.py
```

### 2. Run Tests

```bash
# Run affected tests
pytest tests/unit/test_module.py -v

# Run full suite
pytest -m fast
```

### 3. Type Check

```bash
mypy src/sdp/module.py --strict
```

### 4. Check Coverage

```bash
pytest tests/unit/test_module.py --cov=src/sdp.module --cov-report=term-missing
```

### 5. Commit

```bash
git add src/sdp/module.py tests/unit/test_module.py
git commit -m "feat(module): Add new feature"
```

---

## Common Issues

### Import Errors

**Problem:** `ModuleNotFoundError: No module named 'sdp'`

**Solution:**
```bash
export PYTHONPATH=src
# Or
pip install -e .
```

---

### Mypy Errors

**Problem:** Mypy reports missing imports

**Solution:**
```bash
# Install mypy dependencies
mypy --install-types

# Or use stubs
pip install types-requests types-PyYAML
```

---

### Test Failures

**Problem:** Tests fail with "ImportError"

**Solution:**
```bash
# Ensure virtual environment is activated
source .venv/bin/activate

# Reinstall in development mode
pip install -e .
```

---

### Git Hooks Issues

**Problem:** Pre-commit hook fails

**Solution:**
```bash
# Skip hook (not recommended)
SKIP_CHECK=1 git commit

# Or fix hook issue
./hooks/pre-commit.sh
```

---

## IDE Setup

### VS Code

**Install extensions:**
- Python
- Pylance
- pytest

**Settings:**
```json
{
  "python.defaultInterpreterPath": "${workspaceFolder}/.venv/bin/python",
  "python.linting.enabled": true,
  "python.linting.mypyEnabled": true,
  "python.linting.ruffEnabled": true,
  "python.testing.pytestEnabled": true,
  "python.testing.pytestArgs": ["tests/"]
}
```

### PyCharm

**Configure:**
1. Settings → Project → Python Interpreter
2. Select `.venv` virtual environment
3. Settings → Tools → Python Integrated Tools
4. Set default test runner to pytest

---

## Continuous Integration

### GitHub Actions

**Workflow:** `.github/workflows/test.yml`

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-python@v4
        with:
          python-version: '3.14'
      - run: pip install -e ".[dev]"
      - run: pytest --cov=src/sdp --cov-report=xml
      - run: mypy src/sdp/ --strict
      - run: ruff check src/sdp/
```

---

## Performance Profiling

### Profile Tests

```bash
# Profile test execution
pytest --profile

# Profile with cProfile
python -m cProfile -o output.stats -m pytest
```

### Profile Coverage

```bash
# Find slow tests
pytest --durations=10

# Coverage profiling
pytest --cov=src/sdp --cov-profile
```

---

## Documentation

### Build Docs

```bash
# If using Sphinx
cd docs
make html

# Or serve with live reload
mkdocs serve
```

### Check Links

```bash
# Install markdown-link-check
npm install -g markdown-link-check

# Check all links
markdown-link-check docs/**/*.md
```

---

## Release Process

### 1. Update Version

**File:** `src/sdp/__init__.py`

```python
__version__ = "0.5.0"  # Bump this
```

### 2. Update Changelog

**File:** `CHANGELOG.md`

```markdown
## [0.5.0] - 2026-01-29

### Added
- New feature X

### Changed
- Improved Y
```

### 3. Create Tag

```bash
git tag -a v0.5.0 -m "Release v0.5.0"
git push origin v0.5.0
```

### 4. Build Distribution

```bash
pip install build
python -m build
```

---

## See Also

- [architecture.md](architecture.md) - System architecture
- [extending.md](extending.md) - Extension guide
- [contributing.md](contributing.md) - Contribution guide
- [reference/](../reference/) - User documentation

---

**Version:** SDP v0.5.0
**Updated:** 2026-01-29
