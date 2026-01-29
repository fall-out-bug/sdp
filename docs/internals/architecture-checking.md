# Architecture Checking Guide

Portable clean architecture enforcement using configurable Python module.

## Overview

The architecture checker replaces hardcoded shell script patterns with a configurable Python module that:
- Reads layer patterns from `quality-gate.toml`
- Supports any project structure (hexagonal, onion, layered, custom)
- Validates imports across architectural boundaries
- Provides clear violation messages

## Quick Start

### 1. Configure Your Layers

Create or update `quality-gate.toml`:

```toml
[architecture]
enabled = true

[architecture.layers.domain]
path_regex = "(^|/)domain/"
module_regex = "(^|\\.)domain(\\.|$)"

[architecture.layers.application]
path_regex = "(^|/)application/"
module_regex = "(^|\\.)application(\\.|$)"

[architecture.layers.infrastructure]
path_regex = "(^|/)infrastructure/"
module_regex = "(^|\\.)infrastructure(\\.|$)"

# Forbidden dependencies
forbid_violations = [
    "domain -> application",
    "domain -> infrastructure",
]
```

### 2. Run Checks

```bash
# Check specific files
python scripts/check_architecture.py src/domain/entities/user.py

# Check staged files (from git)
python scripts/check_architecture.py --staged

# Use custom config
python scripts/check_architecture.py --config custom-quality.toml src/
```

## Configuration

### Layer Definition

Each layer requires:
- `name`: Layer identifier (used in rules)
- `path_regex`: Regex to match file paths
- `module_regex`: (Optional) Regex to match import statements

```toml
[architecture.layers.domain]
path_regex = "(^|/)domain/"
module_regex = "(^|\\.)domain(\\.|$)"
```

### Dependency Rules

Define forbidden layer dependencies:

```toml
forbid_violations = [
    "source_layer -> target_layer",
]
```

Example:
```toml
forbid_violations = [
    # Domain is innermost - cannot import from anyone
    "domain -> application",
    "domain -> infrastructure",
    "domain -> presentation",

    # Application can use domain, but not infrastructure/presentation
    "application -> infrastructure",
    "application -> presentation",
]
```

## Architecture Patterns

### Hexagonal Architecture

`quality-gate.toml`:

```toml
[architecture]
enabled = true

[architecture.layers.domain]
path_regex = "(^|/)domain/"
module_regex = "(^|\\.)domain(\\.|$)"

[architecture.layers.application]
path_regex = "(^|/)application/"
module_regex = "(^|\\.)application(\\.|$)"

[architecture.layers.infrastructure]
path_regex = "(^|/)infrastructure/"
module_regex = "(^|\\.)infrastructure(\\.|$)"

forbid_violations = [
    "domain -> application",
    "domain -> infrastructure",
    "application -> infrastructure",
]
```

See: `docs/examples/quality-gate-hexagonal.toml`

### Onion Architecture

`quality-gate.toml`:

```toml
[architecture]
enabled = true

[architecture.layers.domain]
path_regex = "(^|/)(domain|core)/"
module_regex = "(^|\\.)(domain|core)(\\.|$)"

[architecture.layers.application]
path_regex = "(^|/)application/"

[architecture.layers.infrastructure]
path_regex = "(^|/)infrastructure/"

forbid_violations = [
    "domain -> application",
    "domain -> infrastructure",
    "application -> infrastructure",
]
```

See: `docs/examples/quality-gate-onion.toml`

### Traditional Layered Architecture

`quality-gate.toml`:

```toml
[architecture]
enabled = true

[architecture.layers.presentation]
path_regex = "(^|/)(presentation|ui)/"

[architecture.layers.business]
path_regex = "(^|/)(business|service)/"

[architecture.layers.persistence]
path_regex = "(^|/)(persistence|dao)/"

forbid_violations = [
    "business -> presentation",
    "persistence -> presentation",
    "persistence -> business",
]
```

See: `docs/examples/quality-gate-layered.toml`

### Custom Architecture

Define your own layers:

```toml
[architecture]
enabled = true

[architecture.layers.core]
path_regex = "(^|/)src/core/"

[architecture.layers.plugins]
path_regex = "(^|/)plugins/"

[architecture.layers.ui]
path_regex = "(^|/)ui/"

forbid_violations = [
    "core -> plugins",
    "core -> ui",
    "plugins -> ui",
]
```

## Pre-Commit Integration

The hook `hooks/pre-commit.sh` automatically uses the architecture checker:

```bash
# When you commit Python files, architecture is checked
git add src/domain/entities.py
git commit -m "Add entity"

# If violation exists:
# ❌ Architecture violations found (1)
#
#   src/domain/entities.py:5
#     Architecture violation: domain cannot import from infrastructure
```

## Examples

### Violating Code

```python
# src/domain/entities/user.py
from infrastructure.database import Database  # ❌ VIOLATION

class UserEntity:
    def __init__(self, name: str):
        self.name = name
```

**Error:**
```
❌ Architecture violations found (1)

  src/domain/entities/user.py:1
    Architecture violation: domain cannot import from infrastructure
```

### Correct Code

```python
# src/domain/entities/user.py
from dataclasses import dataclass  # ✅ OK (std lib)

@dataclass
class UserEntity:
    name: str
    email: str
```

### Dependency Inversion

```python
# src/domain/entities/user.py
from abc import ABC, abstractmethod  # ✅ OK (define interface)

class UserRepository(ABC):
    @abstractmethod
    def save(self, user: UserEntity) -> None:
        pass

# src/infrastructure/persistence/sql_user_repo.py
from src.domain.entities import UserRepository  # ✅ OK (infra can import domain)

class SqlUserRepository(UserRepository):
    def save(self, user: UserEntity) -> None:
        # Implementation
        pass
```

## Troubleshooting

### Violations Not Detected

1. **Check configuration is loaded:**
   ```bash
   python -c "from sdp.quality.config import QualityGateConfigLoader; print(QualityGateConfigLoader().config.architecture.enabled)"
   ```

2. **Verify layer patterns match your files:**
   ```bash
   python -c "
   import re
   pattern = r'(^|/)domain/'
   print(bool(re.search(pattern, 'src/domain/entity.py')))
   "

   # Should print: True
   ```

3. **Check import detection:**
   ```bash
   python scripts/check_architecture.py --verbose src/domain/entity.py
   ```

### False Positives

If valid imports are flagged:

1. **Adjust layer patterns:**
   ```toml
   [architecture.layers.domain]
   path_regex = "(^|/)domain/"
   module_regex = "(^|\\.)myapp\\.domain(\\.|$)"  # More specific
   ```

2. **Whitelist specific imports:**
   ```toml
   [architecture]
   # Add to allowed_layer_imports if needed
   allowed_layer_imports = [
       "domain -> typing",
   ]
   ```

### Performance Issues

For large codebases:

1. **Exclude test files:**
   ```toml
   [architecture]
   enabled = true
   exclude_patterns = ["*/tests/*", "*/test_*.py"]
   ```

2. **Check only staged files:**
   ```bash
   # In pre-commit hook (default)
   python scripts/check_architecture.py --staged
   ```

## API Reference

### ArchitectureChecker Class

```python
from sdp.quality.architecture import ArchitectureChecker
from sdp.quality.config import QualityGateConfigLoader

config = QualityGateConfigLoader()
violations = []

checker = ArchitectureChecker(
    config=config.config.architecture,
    violations=violations,
)

# Check a file
import ast
with open("src/domain/entity.py") as f:
    tree = ast.parse(f.read())

checker.check_architecture(Path("src/domain/entity.py"), tree)

# Violations are appended to the list
for v in violations:
    print(f"{v.file_path}:{v.line_number}: {v.message}")
```

### Exit Codes

- `0`: All checks passed
- `1`: Architecture violations found
- `2`: Configuration error

## Migration from Hardcoded Checks

### Old (Hardcoded in Shell)

```bash
# hooks/pre-commit.sh
DOMAIN_FILES=$(echo "$STAGED_FILES" | grep "domain/.*\.py$")
BAD_IMPORTS=$(git diff --cached -- $DOMAIN_FILES | grep -E "^\+.*from hw_checker\.(infrastructure|presentation)")
```

### New (Configurable Python)

```bash
# hooks/pre-commit.sh
python scripts/check_architecture.py --staged
```

**Benefits:**
- ✅ Portable across projects
- ✅ Configurable via TOML
- ✅ Supports any architecture pattern
- ✅ Clear violation messages
- ✅ No hardcoded project names

## Advanced Topics

### Cross-Project Dependencies

For monorepos with multiple projects:

```toml
[architecture.layers.project_a_domain]
path_regex = "(^|/)projects/a/domain/"
module_regex = "(^|\\.)projects\\.a\\.domain(\\.|$)"

[architecture.layers.project_b_domain]
path_regex = "(^|/)projects/b/domain/"
module_regex = "(^|\\.)projects\\.b\\.domain(\\.|$)"

forbid_violations = [
    # Project A domain cannot import from Project B infrastructure
    "project_a_domain -> project_b_infrastructure",
]
```

### Conditional Rules

Enable/disable based on environment:

```toml
[architecture]
enabled = true  # Set to false to disable temporarily

# Or use environment variable:
# export SDP_ARCHITECTURE_CHECKS=false
```

## Contributing

To add support for new architecture patterns:

1. Create example config in `docs/examples/quality-gate-{pattern}.toml`
2. Add integration test in `tests/integration/test_architecture_checking.py`
3. Update this guide with pattern description

## See Also

- [Quality Gates Documentation](quality-gates.md)
- [Configuration Reference](../README.md#configuration)
- [Example Configurations](docs/examples/)
