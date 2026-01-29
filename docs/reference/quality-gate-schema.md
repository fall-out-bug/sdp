# Quality Gate Configuration Guide

**Version:** 1.0.0
**Schema File:** `quality-gate.toml`

## Overview

The `quality-gate.toml` file defines quality thresholds for code validation in SDP projects. It provides a centralized, declarative way to enforce coding standards across your team.

## Quick Start

1. **Create** `quality-gate.toml` in your project root
2. **Configure** quality thresholds (see sections below)
3. **Run** validation:
   ```python
   from sdp.quality import QualityGateValidator

   validator = QualityGateValidator()
   violations = validator.validate_directory("src/")
   validator.print_report()
   ```

## Configuration Sections

### [coverage] - Test Coverage

Controls test coverage requirements.

```toml
[coverage]
enabled = true              # Enable/disable coverage checks
minimum = 80                # Target coverage percentage
fail_under = 80             # Minimum coverage to pass CI/CD
exclude_patterns = [        # Patterns to exclude from coverage
    "*/tests/*",
    "*/test_*.py",
    "*/__pycache__/*",
    "*/migrations/*"
]
```

**Best Practices:**
- Set `minimum` to your target goal
- Set `fail_under` to the minimum acceptable threshold
- Exclude test files, generated code, and migrations
- Use 80% for most projects, 90%+ for critical systems

### [complexity] - Cyclomatic Complexity

Controls function complexity limits.

```toml
[complexity]
enabled = true              # Enable/disable complexity checks
max_cc = 10                 # Maximum complexity per function
max_average_cc = 5          # Maximum average complexity per file
```

**What is Cyclomatic Complexity?**
- Measures the number of independent paths through code
- Higher complexity = harder to test and maintain
- Each `if`, `for`, `while`, `except` adds 1

**Thresholds:**
- 1-10: Simple, low risk
- 11-20: Moderate complexity
- 21+: High complexity, should refactor

**Best Practices:**
- Keep functions under 10 complexity
- Extract helper methods for complex logic
- Use early returns to reduce nesting

### [file_size] - File Size Limits

Controls maximum file size in lines of code.

```toml
[file_size]
enabled = true              # Enable/disable file size checks
max_lines = 200             # Maximum lines per file
max_imports = 20            # Maximum import statements per file
max_functions = 15          # Maximum functions/classes per file
```

**Rationale:**
- Smaller files are easier to understand
- Forces single responsibility principle
- Improves testability

**Best Practices:**
- Split large files by responsibility
- Group related functions
- Consider package structure for 20+ files

### [type_hints] - Type Hinting

Controls Python type hint requirements.

```toml
[type_hints]
enabled = true              # Enable/disable type hint checks
require_return_types = true # Require -> return annotations
require_param_types = true  # Require parameter type annotations
strict_mode = true          # Enable strict mypy checks
allow_implicit_any = false  # Forbid implicit Any types
```

**Why Type Hints?**
- Catch bugs before runtime
- Improve IDE autocomplete
- Self-documenting code
- Required for modern Python (3.10+)

**Best Practices:**
- Always type function signatures
- Use `from typing import ...` for complex types
- Run `mypy --strict` in CI/CD
- Avoid `Any` unless necessary

### [error_handling] - Error Handling Patterns

Controls error handling quality.

```toml
[error_handling]
enabled = true                  # Enable/disable error handling checks
forbid_bare_except = true       # Forbid `except:` without exception type
forbid_bare_raise = true        # Forbid `raise` without exception object
forbid_pass_with_except = true  # Forbid empty except blocks
require_explicit_errors = true  # Require custom exception types
```

**Anti-Patterns to Avoid:**
```python
# BAD: Bare except
try:
    ...
except:  # Catches everything, including SystemExit
    pass

# BAD: Empty except
try:
    ...
except Exception:
    pass  # Silently ignores errors

# GOOD: Specific exception
try:
    ...
except ValueError as e:
    logger.error(f"Invalid value: {e}")
    raise
```

**Best Practices:**
- Catch specific exceptions
- Log errors before re-raising
- Use context managers (`with`) for resources
- Define domain-specific exceptions

### [architecture] - Clean Architecture

Enforces layer separation and dependency rules.

```toml
[architecture]
enabled = true                  # Enable/disable architecture checks
enforce_layer_separation = true # Enforce clean architecture layers
allowed_layer_imports = [       # Permitted import patterns
    "domain <- application",
    "domain <- infrastructure",
    "domain <- presentation",
    "application <- infrastructure",
    "application <- presentation",
    "infrastructure <- presentation"
]
forbid_violations = [           # Forbidden import patterns
    "infrastructure -> domain",
    "presentation -> domain",
    "presentation -> application"
]
```

**Clean Architecture Layers:**
1. **Domain** - Business logic, entities (innermost)
2. **Application** - Use cases, orchestration
3. **Infrastructure** - Database, API, external services
4. **Presentation** - HTTP, CLI, UI (outermost)

**Dependency Rule:**
- Dependencies point **inward** only
- Outer layers can use inner layers
- Inner layers never use outer layers

**Example Violation:**
```python
# BAD: Domain importing Infrastructure
# src/domain/entities.py
from infrastructure.database import Base  # VIOLATION

class User(Base):  # Domain depends on Infra
    pass

# GOOD: Infrastructure depends on Domain
# src/infrastructure/database.py
from domain.entities import User  # OK

class Base:
    pass
```

### [documentation] - Documentation Requirements

Controls code documentation standards.

```toml
[documentation]
enabled = true                      # Enable/disable documentation checks
require_docstrings = false          # Require docstrings on all functions
min_docstring_coverage = 0.5        # Minimum ratio of documented functions
require_module_docstrings = true    # Require module-level docstrings
require_class_docstrings = false    # Require class docstrings
require_function_docstrings = false # Require function docstrings
```

**Best Practices:**
- Document **why**, not **what**
- Use Google or NumPy docstring style
- Include examples for complex functions
- Keep docstrings up-to-date

**Example:**
```python
def calculate_user_score(user: User) -> float:
    """Calculate the user's reputation score based on activity.

    The score is computed as a weighted sum of:
    - Number of posts (weight: 2.0)
    - Number of comments (weight: 1.0)
    - Account age in days (weight: 0.1)

    Args:
        user: The user object to calculate score for.

    Returns:
        The user's reputation score as a float.

    Raises:
        ValueError: If user has negative activity counts.

    Example:
        >>> user = User(posts=10, comments=20, age_days=100)
        >>> calculate_user_score(user)
        50.0
    """
    ...
```

### [testing] - Test Quality

Controls test quality requirements.

```toml
[testing]
enabled = true                  # Enable/disable test quality checks
require_test_for_new_code = true # Require tests for new code
min_test_to_code_ratio = 0.8    # Minimum test LOC / code LOC ratio
require_fast_marker = true      # Require pytest 'fast' marker on unit tests
forbid_print_statements = true  # Forbid print() in tests (use assertions)
```

**Best Practices:**
- Write tests before code (TDD)
- Use descriptive test names
- One assertion per test
- Use fixtures for setup

**Example:**
```python
# tests/test_user.py
import pytest

def test_user_score_calculation():
    """Test that user score is calculated correctly."""
    user = User(posts=10, comments=20)
    assert calculate_user_score(user) == 40.0

def test_user_score_with_negative_posts_raises_error():
    """Test that negative post counts raise ValueError."""
    user = User(posts=-1, comments=0)
    with pytest.raises(ValueError):
        calculate_user_score(user)
```

### [naming] - Naming Conventions

Controls variable/function naming standards.

```toml
[naming]
enabled = true                    # Enable/disable naming checks
enforce_pep8 = true               # Enforce PEP 8 naming conventions
allow_single_letter = false       # Forbid single-letter variable names
min_variable_name_length = 3      # Minimum variable name length
max_variable_name_length = 50     # Maximum variable name length
```

**PEP 8 Conventions:**
- `snake_case` for functions and variables
- `PascalCase` for classes
- `UPPER_CASE` for constants
- `_leading_underscore` for private/internal

**Examples:**
```python
# GOOD
user_name = "Alice"
class UserAccount: pass
MAX_CONNECTIONS = 100
def calculate_score(): pass

# BAD
userName = "Bob"  # camelCase
class useraccount: pass  # lowercase
max_connections = 100  # should be UPPER
def CalculateScore(): pass  # PascalCase
```

### [security] - Security Checks

Controls security-related validations.

```toml
[security]
enabled = true                       # Enable/disable security checks
forbid_hardcoded_secrets = true      # Detect hardcoded passwords/API keys
forbid_sql_injection_patterns = true # Detect SQL injection risks
forbid_eval_usage = true             # Forbid eval() usage
require_https_urls = true            # Require HTTPS in URLs
```

**Anti-Patterns:**
```python
# BAD: Hardcoded secrets
password = "SuperSecret123"
api_key = "sk-1234567890"

# BAD: SQL injection risk
query = f"SELECT * FROM users WHERE name='{user_input}'"

# BAD: eval usage
result = eval(user_input)  # Arbitrary code execution

# GOOD: Use environment variables
import os
password = os.getenv("DB_PASSWORD")
assert password, "DB_PASSWORD not set"

# GOOD: Parameterized queries
cursor.execute("SELECT * FROM users WHERE name=?", (user_input,))

# GOOD: Use literal_eval for safe parsing
from ast import literal_eval
result = literal_eval(user_input)
```

### [performance] - Performance Checks

Controls performance-related validations.

```toml
[performance]
enabled = true                           # Enable/disable performance checks
forbid_sql_queries_in_loops = true       # Detect N+1 query patterns
max_nesting_depth = 5                    # Maximum nesting depth
warn_large_string_concatenation = true   # Warn about string += in loops
```

**Anti-Patterns:**
```python
# BAD: SQL query in loop
for user in users:
    posts = db.execute(f"SELECT * FROM posts WHERE user_id={user.id}")

# GOOD: Batch query
user_ids = [u.id for u in users]
posts = db.execute("SELECT * FROM posts WHERE user_id IN (?)", (user_ids,))

# BAD: String concatenation in loop
result = ""
for item in items:
    result += str(item)  # O(n²) performance

# GOOD: Use join
result = "".join(str(item) for item in items)
```

## Usage Examples

### Example 1: Basic Setup

```toml
# quality-gate.toml
[coverage]
minimum = 80
fail_under = 70

[complexity]
max_cc = 10

[file_size]
max_lines = 200
```

```python
# validate.py
from sdp.quality import QualityGateValidator

validator = QualityGateValidator()
violations = validator.validate_directory("src/")

if violations:
    print(f"Found {len(violations)} violations!")
    for v in violations:
        print(f"  {v}")
    exit(1)
else:
    print("All quality gates passed!")
```

### Example 2: Custom Configuration

```python
from sdp.quality.config import QualityGateConfig

# Load custom config
validator = QualityGateValidator(config_path="config/custom-gate.toml")
violations = validator.validate_file("src/module.py")

# Print report
validator.print_report()
```

### Example 3: CI/CD Integration

```yaml
# .github/workflows/quality-gate.yml
name: Quality Gate

on: [pull_request]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-python@v4
        with:
          python-version: '3.10'

      - name: Install dependencies
        run: |
          pip install -r requirements.txt

      - name: Run quality gate
        run: |
          python -c "
          from sdp.quality import QualityGateValidator
          validator = QualityGateValidator()
          violations = validator.validate_directory('src/')
          validator.print_report()
          exit(len([v for v in violations if v.severity == 'error']))
          "
```

## Integration with SDP Tools

### Pre-commit Hook

```bash
# .git/hooks/pre-commit
#!/bin/bash
python -c "
from sdp.quality import QualityGateValidator
validator = QualityGateValidator()
violations = validator.validate_directory('src/')
if violations:
    print('Quality gate violations found!')
    validator.print_report()
    exit(1)
"
```

### VS Code Extension

```json
// .vscode/settings.json
{
    "python.linting.enabled": true,
    "python.linting.sdpQualityEnabled": true,
    "python.linting.sdpQualityArgs": ["--config", "quality-gate.toml"]
}
```

## FAQ

**Q: Should I commit quality-gate.toml?**
A: Yes! Quality standards should be shared across the team.

**Q: Can I have multiple quality-gate.toml files?**
A: Yes, place them in subdirectories for project-specific rules.

**Q: How do I temporarily disable a check?**
A: Set `enabled = false` for that section, or use inline comments:
```python
# quality-gate: disable-next-line=complexity
def complex_function():
    ...
```

**Q: What if my legacy code fails quality gates?**
A: Use `exclude_patterns` or introduce standards gradually:
```toml
[complexity]
enabled = true
max_cc = 20  # Start lenient, tighten over time
```

**Q: Can I auto-fix violations?**
A: Some violations can be auto-fixed with tools like:
- `ruff check --fix` for formatting
- `mypy` for type hints
- Manual refactoring for complexity

## References

- [SDP Protocol](../PROTOCOL.md)
- [Clean Architecture](../PRINCIPLES.md)
- [Python Type Hints](https://docs.python.org/3/library/typing.html)
- [PEP 8 Style Guide](https://pep8.org/)
- [Cyclomatic Complexity](https://en.wikipedia.org/wiki/Cyclomatic_complexity)

---

**Version:** 1.0.0
**Last Updated:** 2025-01-29
**Maintainer:** SDP Team
