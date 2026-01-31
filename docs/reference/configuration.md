# SDP Configuration Reference

Complete reference for SDP configuration files and options.

---

## Table of Contents

- [Configuration Files](#configuration-files)
- [Quality Gate Configuration](#quality-gate-configuration)
- [Claude Code Settings](#claude-code-settings)
- [Environment Variables](#environment-variables)
- [Git Hooks](#git-hooks)

---

## Configuration Files

### Primary Configuration Files

| File | Location | Purpose |
|------|----------|---------|
| **quality-gate.toml** | `/` | Quality gate thresholds |
| **settings.json** | `.claude/` | Claude Code settings |
| **.env** | `/` | Environment variables |
| **pre-commit** | `.git/hooks/` | Git hooks |

---

## Quality Gate Configuration

### quality-gate.toml

**Location:** Repository root

**Purpose:** Define quality thresholds for validation

**Example:**
```toml
# SDP Quality Gate Configuration
# Version: 1.0.0

[coverage]
enabled = true
minimum = 80
fail_under = 80
exclude_patterns = [
    "*/tests/*",
    "*/test_*.py",
]

[complexity]
enabled = true
max_cc = 10
max_average_cc = 5

[file_size]
enabled = true
max_lines = 200
max_imports = 20
max_functions = 15

[type_hints]
enabled = true
require_return_types = true
require_param_types = true
strict_mode = true
allow_implicit_any = false

[error_handling]
enabled = true
forbid_bare_except = true
forbid_bare_raise = true
forbid_pass_with_except = true
require_explicit_errors = true

[architecture]
enabled = true
enforce_layer_separation = true

[documentation]
enabled = true
require_docstrings = false
min_docstring_coverage = 0.5
require_module_docstrings = true

[testing]
enabled = true
require_test_for_new_code = true
min_test_to_code_ratio = 0.8
require_fast_marker = true
forbid_print_statements = true

[naming]
enabled = true
enforce_pep8 = true
allow_single_letter = false
min_variable_name_length = 3
max_variable_name_length = 50

[security]
enabled = true
forbid_hardcoded_secrets = true
forbid_sql_injection_patterns = true
forbid_eval_usage = true
require_https_urls = true

[performance]
enabled = true
forbid_sql_queries_in_loops = true
max_nesting_depth = 5
warn_large_string_concatenation = true
```

### Configuration Sections

#### Coverage

Test coverage requirements.

```toml
[coverage]
enabled = true              # Enable coverage checks
minimum = 80                # Minimum percentage
fail_under = 80             # Fail threshold
exclude_patterns = [        # Files to exclude
    "*/tests/*",
    "*/test_*.py",
]
```

#### Complexity

Cyclomatic complexity limits.

```toml
[complexity]
enabled = true              # Enable complexity checks
max_cc = 10                 # Max complexity per function
max_average_cc = 5          # Max average complexity
```

#### File Size

File size limits (lines of code).

```toml
[file_size]
enabled = true              # Enable file size checks
max_lines = 200             # Max lines per file
max_imports = 20            # Max imports per file
max_functions = 15          # Max functions per file
```

#### Type Hints

Type hinting requirements.

```toml
[type_hints]
enabled = true                     # Enable type hint checks
require_return_types = true        # Require return type hints
require_param_types = true         # Require parameter type hints
strict_mode = true                 # Use mypy --strict
allow_implicit_any = false         # Disallow implicit Any
```

---

## Claude Code Settings

### .claude/settings.json

**Location:** `.claude/settings.json`

**Purpose:** Claude Code configuration

**Example:**
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
  ],
  "hooks": {
    "pre-commit": "hooks/pre-commit.sh",
    "post-build": "hooks/post-build.sh"
  },
  "qualityGates": {
    "enabled": true,
    "config": "quality-gate.toml"
  }
}
```

---

## Environment Variables

### Required Variables

```bash
# Python
PYTHONPATH=src
VIRTUAL_ENV=.venv

# Beads (optional)
BEADS_HOME=~/.beads
BEADS_SERVER=http://localhost:8080

# GitHub (optional)
GITHUB_TOKEN=ghp_*
GITHUB_REPO=owner/repo

# Telegram (optional)
TELEGRAM_BOT_TOKEN=*
TELEGRAM_CHAT_ID=*
```

### Optional Variables

```bash
# Skip quality gates (not recommended)
SKIP_QUALITY_GATES=1

# Skip commit checks
SKIP_COMMIT_CHECK=1

# Debug mode
DEBUG=1
VERBOSE=1
```

---

## Git Hooks

### Pre-commit Hook

**Location:** `.git/hooks/pre-commit` → `hooks/pre-commit.sh`

**Checks:**
- No time estimates in WS files
- No tech debt markers
- Python code quality
- Clean architecture violations
- WS file format

**Enable:**
```bash
ln -sf ../../hooks/pre-commit.sh .git/hooks/pre-commit
```

---

### Pre-build Hook

**Location:** `hooks/pre-build.sh`

**Checks:**
- Goal section exists
- Acceptance criteria present
- Scope not LARGE
- Dependencies completed

**Usage:**
```bash
./hooks/pre-build.sh WS-001-01
```

---

### Post-build Hook

**Location:** `hooks/post-build.sh`

**Checks:**
- Regression tests pass
- Linters pass
- No TODO/FIXME markers
- Coverage adequate
- Execution report appended

**Usage:**
```bash
./hooks/post-build.sh WS-001-01
```

---

## Validation

### Check Configuration

```bash
# Validate quality-gate.toml
python -m sdp.quality.config validate quality-gate.toml

# Check all configs
python -m sdp.doctor
```

---

### Test Configuration

```bash
# Run quality gates
python -m sdp.quality.validator check

# Check specific file
python -m sdp.quality.validator check src/sdp/module.py
```

---

## Quick Reference

### Common Config Tasks

**Change coverage threshold:**
```toml
[coverage]
minimum = 90  # Change from 80 to 90
```

**Disable specific check:**
```toml
[documentation]
enabled = false  # Disable doc checks
```

**Exclude files from coverage:**
```toml
[coverage]
exclude_patterns = [
    "*/tests/*",
    "*/legacy/*",  # Add legacy folder
]
```

---

## Best Practices

### DO ✅

1. **Keep config in sync** - Update quality-gate.toml when standards change
2. **Commit configs** - Track all configuration files
3. **Document changes** - Add comments for non-standard settings
4. **Test configs** - Validate after changes

### DON'T ❌

1. **Don't skip quality gates** - Avoid `SKIP_QUALITY_GATES=1`
2. **Don't lower standards** - Keep thresholds high
3. **Don't ignore mypy** - Use `--strict` mode
4. **Don't exclude too much** - Minimize exclude patterns

---

## See Also

- [quality-gates.md](quality-gates.md) - Quality gate details
- [commands.md](commands.md) - Command reference
- [error-handling.md](error-handling.md) - Error patterns

---

**Version:** SDP v0.5.0
**Updated:** 2026-01-29
