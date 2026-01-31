# Test Watch Mode

## Overview

The SDP project uses **pytest-watcher** for automated test re-running during development. This provides fast feedback (< 3 seconds) when you change code, which is essential for effective Test-Driven Development (TDD).

**Why pytest-watcher?**
- Actively maintained (last update: Aug 2024)
- Modern replacement for the unmaintained pytest-watch
- Cross-platform (Linux, macOS, Windows, BSD)
- Native file system monitoring (no polling)
- Smart delay handling for IDE formatters
- Pyproject.toml configuration support

## Installation

```bash
# Already in pyproject.toml
poetry install
```

## Usage

### Basic Watch Mode

Watch all unit tests (excludes slow integration tests):

```bash
bash scripts/watch-tests.sh
```

### Advanced Usage

Include integration tests:

```bash
bash scripts/watch-tests.sh --with-integration
```

Run with coverage report:

```bash
bash scripts/watch-tests.sh --cov
```

Filter tests by keyword:

```bash
bash scripts/watch-tests.sh -k test_specific_function
```

Run specific test file:

```bash
bash scripts/watch-tests.sh tests/unit/schema/test_validator.py
```

Combine options:

```bash
bash scripts/watch-tests.sh --with-integration --cov -k test_user_auth
```

### Direct pytest-watcher Usage

You can also use pytest-watcher directly without the wrapper script:

```bash
# Watch current directory
poetry run pytest-watcher .

# Watch with custom pytest args
poetry run pytest-watcher . -- -x --lf --nf

# Watch specific directory
poetry run pytest-watcher tests/unit/

# Watch with different patterns
poetry run pytest-watcher . --patterns '*.py,pyproject.toml'

# See all options
poetry run pytest-watcher --help
```

**Note:** Use `pytest-watcher` (not `ptw` which is from the abandoned pytest-watch package).

## Configuration

Default settings in `pyproject.toml`:

```toml
[tool.pytest-watcher]
now = false          # Don't run tests on startup
clear = true         # Clear screen before each run
delay = 0.2          # Wait 200ms for file saves/formatters
runner = "pytest"
runner_args = ["--ignore=tests/integration", "-v"]
patterns = ["*.py"]
ignore_patterns = []
```

### Key Options Explained

- **`now = false`**: Waits for file change before first run
  - Set to `true` to run tests immediately on start
- **`delay = 0.2`**: Prevents running tests while files are being saved
  - Useful for IDEs with auto-formatting on save
  - Set to `0` to disable delay
- **`clear = true`**: Cleans screen between runs for clarity

## TDD Workflow

1. **Start watcher** in terminal:
   ```bash
   bash scripts/watch-tests.sh
   ```

2. **Write failing test** (Red):
   - Edit or create test file
   - Watcher detects change and runs test
   - See test fail (expected)

3. **Write minimal code** (Green):
   - Edit source code
   - Watcher detects change and re-runs test
   - See test pass

4. **Refactor**:
   - Clean up code
   - Watcher continuously validates changes
   - All tests stay green

## Performance

- **First run**: ~1-2 seconds (test collection + execution)
- **Subsequent runs**: < 3 seconds (incremental)
- **Unit tests only**: ~50-100 tests in < 2 seconds
- **With integration**: 5-10 seconds (depends on external services)

## Keyboard Shortcuts (Interactive Mode)

When running in a terminal with TTY support:

- **Ctrl+C**: Stop the watcher
- Future enhancements: restart, clear, etc. (see pytest-watcher docs)

## Troubleshooting

### Watcher not detecting changes

Check file patterns:
```bash
poetry run pytest-watcher . --patterns '*.py' --ignore-patterns '*/migrations/*'
```

### Tests run too slowly

Exclude slow tests:
```bash
# Edit pyproject.toml runner_args
runner_args = ["--ignore=tests/integration", "--ignore=tests/slow", "-v"]
```

Or use pytest markers:
```python
@pytest.mark.slow
def test_slow_operation():
    # ...
```

Then:
```bash
poetry run pytest-watcher . -- -m "not slow"
```

### Terminal errors

If you see "Unable to initialize terminal state", the watcher will still work but interactive mode is disabled. This is normal in CI/CD or non-TTY environments.

## Comparison with Alternatives

| Tool | Status | Pros | Cons |
|------|--------|------|------|
| **pytest-watcher** | ✅ Active | Modern, configurable, cross-platform | Newer ecosystem |
| pytest-watch | ❌ Unmaintained | Widely known | Doesn't work for many users |
| pytest-xwatch | ⚠️ Unclear | Simple | Last update 2020 |
| entr | ✅ Active | Generic, powerful | Not pytest-specific, complex setup |

**Recommendation**: Use pytest-watcher for SDP projects.

## References

- [pytest-watcher Documentation](https://github.com/olzhasar/pytest-watcher)
- [pytest-watch Issue #121: Maintenance Status](https://github.com/joeyespo/pytest-watch/issues/121)
- [SDP TDD Protocol](../../PROTOCOL.md#test-driven-development-tdd)
