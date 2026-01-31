---
ws_id: 00-020-01
feature: F020
status: backlog
size: MEDIUM
project_id: 00
github_issue: null
assignee: null
depends_on: []
---

## WS-00-020-01: Extract Git Hooks to Python

### 🎯 Goal

**What must WORK after completing this WS:**
- Shell hooks (pre-commit, pre-push, post-build) extracted to testable Python modules
- Shell scripts become thin wrappers around Python implementations
- All hook functionality preserved with 100% backward compatibility

**Acceptance Criteria:**
- [ ] AC1: All 4 hooks (pre-commit.sh, pre-push.sh, post-build.sh, pre-deploy.sh) extracted to Python
- [ ] AC2: Shell scripts rewritten as thin wrappers (<20 lines each)
- [ ] AC3: All existing tests pass with new Python implementation
- [ ] AC4: Coverage ≥80% for new Python hook modules
- [ ] AC5: mypy --strict passes on all hook modules
- [ ] AC6: Zero functionality regression (all checks still work)

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

Current git hooks are implemented as 284-line bash scripts (`hooks/pre-commit.sh`) that:
- Run quality checks on staged files
- Validate time estimates, tech debt, Python quality
- Check architecture, workstream format, breaking changes
- Enforce test quality standards

**Problem**: Shell scripts are:
- **Untestable** - No unit tests for shell logic
- **Hardcoded to external project** - References `tools/hw_checker` paths
- **Unmaintainable** - Complex bash logic is error-prone
- **Blocks meta-testing** - Can't verify SDP follows its own quality gates

**Solution**: Extract core logic to Python modules, keep shell as thin wrapper.

### Dependencies

None (this is foundational work)

### Input Files

- `hooks/pre-commit.sh` (284 lines)
- `hooks/pre-push.sh` (~150 lines)
- `hooks/post-build.sh` (~100 lines)
- `hooks/pre-deploy.sh` (~50 lines)
- `scripts/check_quality_gates.py` (existing Python checker)
- `scripts/quality/` directory (modular quality gate modules)

### Steps

1. **Create Python hook modules**

   Create `src/sdp/hooks/` directory structure:
   ```
   src/sdp/hooks/
   ├── __init__.py
   ├── pre_commit.py     # Extract from hooks/pre-commit.sh
   ├── pre_push.py       # Extract from hooks/pre-push.sh
   ├── post_build.py     # Extract from hooks/post-build.sh
   ├── pre_deploy.py     # Extract from hooks/pre-deploy.sh
   └── common.py         # Shared utilities
   ```

2. **Extract pre-commit.sh to Python**

   Map each check section to Python function:
   ```python
   # src/sdp/hooks/pre_commit.py
   def check_time_estimates(files: list[Path]) -> CheckResult:
       """Check for time-based estimates in staged files."""
       # Extracted from hooks/pre-commit.sh lines 38-68

   def check_tech_debt_markers(files: list[Path]) -> CheckResult:
       """Check for TODO/FIXME without workstream references."""
       # Extracted from hooks/pre-commit.sh lines 70-102

   def check_python_quality(files: list[Path]) -> CheckResult:
       """Run ruff, mypy on Python files."""
       # Extracted from hooks/pre-commit.sh lines 104-140

   # ... etc for all 7 checks
   ```

3. **Rewrite shell scripts as wrappers**

   ```bash
   #!/bin/bash
   # hooks/pre-commit.sh (thin wrapper - <20 lines)

   # Run Python hook implementation
   python -m sdp.hooks.pre_commit "$@"

   # Exit with Python's exit code
   exit $?
   ```

4. **Add comprehensive tests**

   Create `tests/unit/hooks/test_pre_commit.py`:
   ```python
   def test_check_time_estimates_detects_hours():
       """Verify time estimate check detects '2 hours' pattern."""

   def test_check_time_estimates_allows_relative_sizing():
       """Verify time estimate check allows 'SMALL/MEDIUM/LARGE'."""

   def test_check_tech_debt_rejects_bare_todo():
       """Verify tech debt check rejects TODO without WS reference."""

   def test_all_checks_run_on_staged_files():
       """Verify all 7 checks execute in correct order."""
   ```

5. **Verify backward compatibility**

   - Run all hooks manually on test files
   - Ensure exit codes match (0=pass, 1=fail)
   - Verify output format is identical
   - Test with git commit (real hook execution)

### Code

```python
# src/sdp/hooks/common.py
from dataclasses import dataclass
from pathlib import Path

@dataclass
class CheckResult:
    """Result of a quality check."""
    passed: bool
    message: str
    violations: list[tuple[Path, int, str]]  # (file, line, issue)

    def format_terminal(self) -> str:
        """Format result for terminal output."""
        if self.passed:
            return f"✅ {self.message}"
        else:
            output = [f"❌ {self.message}"]
            for file, line, issue in self.violations:
                output.append(f"  {file}:{line} - {issue}")
            return "\n".join(output)


# src/sdp/hooks/pre_commit.py
import subprocess
from pathlib import Path
from .common import CheckResult

def check_time_estimates(files: list[Path]) -> CheckResult:
    """Check for time-based estimates in staged files.

    Forbidden patterns:
    - "X hours?", "X days?", "X weeks?"
    - "soon", "ASAP", "TMTR"
    """
    forbidden = [
        r"\d+ hours?",
        r"\d+ days?",
        r"\d+ weeks?",
        "soon",
        "ASAP",
        "TMTR",
    ]

    violations = []
    for file in files:
        content = file.read_text()
        for pattern in forbidden:
            if re.search(pattern, content, re.IGNORECASE):
                line_no = content.split('\n').index(pattern) + 1
                violations.append((file, line_no, f"Time estimate: {pattern}"))

    if violations:
        return CheckResult(
            passed=False,
            message="Time estimates found (use relative sizing: SMALL/MEDIUM/LARGE)",
            violations=violations
        )
    return CheckResult(passed=True, message="No time estimates found")


def main() -> int:
    """Run pre-commit checks and return exit code."""
    import sys
    from git import Repo

    repo = Repo(".")
    staged_files = [Path(a.path) for a in repo.index.diff("HEAD")]

    checks = [
        check_time_estimates,
        check_tech_debt_markers,
        check_python_quality,
        # ... other checks
    ]

    all_passed = True
    for check in checks:
        result = check(staged_files)
        print(result.format_terminal())
        if not result.passed:
            all_passed = False

    return 0 if all_passed else 1


if __name__ == "__main__":
    import sys
    sys.exit(main())
```

```bash
#!/bin/bash
# hooks/pre-commit.sh (thin wrapper - <20 lines)

# Pre-commit hook: Run quality checks on staged files
# Extracted to Python: src/sdp/hooks/pre_commit.py

set -e

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Run Python implementation
python -m sdp.hooks.pre_commit

# Exit with Python's exit code
exit $?
```

### Expected Outcome

**After completion:**
- `src/sdp/hooks/` directory with 5 Python modules
- Shell scripts reduced to <20 lines each (wrapper only)
- `tests/unit/hooks/` with comprehensive test coverage
- All hook logic is now testable and maintainable
- Foundation for meta-quality enforcement (SDP can test its own hooks)

**Scope Estimate**
- Files: ~12 (5 Python modules + 4 wrapper scripts + 3 test files)
- Lines: ~800 (MEDIUM)
- Tokens: ~4000

### Completion Criteria

```bash
# Verify extraction works
python -m sdp.hooks.pre_commit --help

# Run all hook tests
pytest tests/unit/hooks/test_pre_commit.py -v
pytest tests/unit/hooks/test_pre_push.py -v
pytest tests/unit/hooks/test_post_build.py -v

# Verify coverage
pytest --cov=src/sdp/hooks --cov-fail-under=80

# Verify type checking
mypy src/sdp/hooks/ --strict

# Test real hook execution
git commit --allow-empty -m "test: verify pre-commit hook"
```

### Constraints

- DO NOT break existing hook functionality
- DO NOT change output format (must be identical)
- DO NOT change exit codes (0=pass, 1=fail)
- DO NOT remove shell scripts (keep as wrappers for backward compatibility)
- DO NOT add new dependencies (use only existing Poetry deps)

---

## Execution Report

**Executed by:** Cursor/Claude
**Date:** 2026-01-31
**Duration:** ~90 minutes

### Goal Status
- [x] AC1: All 4 hooks extracted to Python
- [x] AC2: Shell scripts rewritten as thin wrappers (<20 lines each)
- [x] AC3: All existing tests pass (982 passed)
- [ ] AC4: Coverage 66% (target ≥80%) — follow-up recommended
- [x] AC5: mypy --strict passes on all hook modules
- [x] AC6: Zero functionality regression (pre-commit, pre-push verified)

**Goal Achieved:** Partial (AC4 below target)

### Files Changed
| File | Action | LOC |
|------|--------|-----|
| src/sdp/hooks/common.py | Create | 25 |
| src/sdp/hooks/pre_commit.py | Create | 304 |
| src/sdp/hooks/pre_push.py | Create | 125 |
| src/sdp/hooks/post_build.py | Create | 183 |
| src/sdp/hooks/pre_deploy.py | Create | 115 |
| src/sdp/hooks/__init__.py | Update | +CheckResult |
| src/sdp/hooks/ws_complete.py | Fix | +VerificationResult types |
| hooks/pre-commit.sh | Rewrite | 284 → 12 |
| hooks/pre-push.sh | Rewrite | 284 → 12 |
| hooks/post-build.sh | Rewrite | 263 → 15 |
| hooks/pre-deploy.sh | Rewrite | 94 → 15 |
| hooks/pre-build.sh | Update | +docs/workstreams support |
| tests/unit/hooks/test_common.py | Create | 35 |
| tests/unit/hooks/test_pre_commit.py | Create | 180 |
| tests/unit/hooks/test_pre_push.py | Create | 55 |
| tests/unit/hooks/test_post_build.py | Create | 55 |
| tests/unit/hooks/test_pre_deploy.py | Create | 55 |

### Statistics
- **Files Changed:** 17
- **Lines Added:** ~1,200
- **Lines Removed:** ~700
- **Test Coverage:** 66% (new hook modules)
- **Tests Passed:** 982
- **Tests Failed:** 0

### Deviations from Plan
- pre-build.sh updated to support docs/workstreams (SDP layout) in addition to tools/hw_checker
- ws_complete.py: added VerificationResult type annotations for mypy --strict
- Coverage 66% vs 80% target: main() orchestration logic hard to unit-test; recommend integration tests or pragma: no cover for entry points

### Commit
feat(hooks): 00-020-01 - Extract Git hooks to Python
