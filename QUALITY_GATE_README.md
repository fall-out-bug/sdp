# Quality Gate System for SDP

**Version:** 1.0.0
**Status:** Production Ready

## Overview

The Quality Gate System provides automated code quality validation for SDP projects. It enforces coding standards through a declarative `quality-gate.toml` configuration file.

## Quick Start

### 1. Installation

The quality gate system is included with SDP. No additional installation needed.

### 2. Configuration

Create a `quality-gate.toml` file in your project root:

```toml
[coverage]
minimum = 80
fail_under = 80

[complexity]
max_cc = 10

[file_size]
max_lines = 200

[type_hints]
require_return_types = true
require_param_types = true

[error_handling]
forbid_bare_except = true
forbid_pass_with_except = true

[architecture]
enabled = true
enforce_layer_separation = true
```

### 3. Usage

#### Command Line

```bash
# Validate a single file
python -m sdp.quality validate path/to/file.py

# Validate a directory
python -m sdp.quality validate path/to/directory/

# Generate report
python -m sdp.quality validate --report src/
```

#### Python API

```python
from sdp.quality import QualityGateValidator

# Initialize validator
validator = QualityGateValidator()

# Validate a file
violations = validator.validate_file("src/module.py")

# Validate a directory
violations = validator.validate_directory("src/")

# Print report
validator.print_report()

# Get summary
summary = validator.get_summary()
print(f"Total violations: {summary['total']}")
```

## Features

### Quality Checks

| Check | Description | Configurable |
|-------|-------------|--------------|
| **Coverage** | Test coverage percentage | ✅ |
| **Complexity** | Cyclomatic complexity | ✅ |
| **File Size** | Lines of code, imports, functions | ✅ |
| **Type Hints** | Type annotation requirements | ✅ |
| **Error Handling** | Exception handling patterns | ✅ |
| **Architecture** | Layer separation rules | ✅ |
| **Documentation** | Docstring coverage | ✅ |
| **Testing** | Test quality patterns | ✅ |
| **Naming** | PEP 8 conventions | ✅ |
| **Security** | Hardcoded secrets, eval usage | ✅ |
| **Performance** | Nesting depth, SQL queries | ✅ |

### Severity Levels

- **Error**: Must be fixed before merge
- **Warning**: Should be fixed when possible

## Configuration Reference

See [docs/quality-gate-schema.md](docs/quality-gate-schema.md) for complete configuration reference.

## Integration

### Pre-commit Hook

Add to `.git/hooks/pre-commit`:

```bash
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

### GitHub Actions

```yaml
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
      - name: Run quality gate
        run: |
          python -m sdp.quality validate --report src/
```

### VS Code

Add to `.vscode/settings.json`:

```json
{
    "python.linting.enabled": true,
    "python.linting.sdpQualityEnabled": true
}
```

## Examples

See [examples/quality_gate_example.py](examples/quality_gate_example.py) for usage examples.

## Testing

Run the test suite:

```bash
poetry run pytest tests/test_quality_gate.py -v
```

## Architecture

```
src/sdp/quality/
├── __init__.py           # Public API
├── models.py             # Dataclass models
├── config.py             # TOML parser
└── validator.py          # Validation logic
```

## Contributing

When adding new quality checks:

1. Add configuration to `models.py`
2. Add parser logic to `config.py`
3. Add validation logic to `validator.py`
4. Add tests to `tests/test_quality_gate.py`
5. Update documentation

## License

Same as SDP project.

## Support

For issues and questions, see the main SDP documentation.
