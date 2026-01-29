# Error Patterns in SDP

This document describes common error patterns and how to use the SDP error framework effectively.

## Table of Contents

- [Overview](#overview)
- [When to Raise Errors](#when-to-raise-errors)
- [Error Categories](#error-categories)
- [Predefined Error Types](#predefined-error-types)
- [Custom Errors](#custom-errors)
- [Best Practices](#best-practices)
- [Examples](#examples)

---

## Overview

The SDP error framework provides structured error handling with:

1. **Consistent error format** - All errors have category, message, remediation, and optional context
2. **Actionable guidance** - Each error includes specific remediation steps
3. **Documentation links** - Errors link to relevant documentation
4. **Rich context** - Errors include relevant debugging information

### Benefits

- **Better UX**: Users get clear, actionable error messages
- **Faster debugging**: Context information helps identify root causes
- **Consistent API**: All errors follow the same structure
- **Documentation-driven**: Errors link to troubleshooting guides

---

## When to Raise Errors

### Validation Errors

Raise when user input or configuration is invalid:

```python
from sdp.errors import WorkstreamValidationError

def validate_workstream(ws_id: str, ws_file: Path) -> None:
    errors = []

    if not ws_file.exists():
        errors.append(f"File not found: {ws_file}")

    if not has_goal_section(ws_file):
        errors.append("Missing Goal section")

    if errors:
        raise WorkstreamValidationError(
            ws_id=ws_id,
            errors=errors,
            file_path=str(ws_file),
        )
```

### Build Errors

Raise when build validation fails:

```python
from sdp.errors import BuildValidationError

def check_pre_build(ws_id: str) -> None:
    if not has_goal_section(ws_id):
        raise BuildValidationError(
            ws_id=ws_id,
            stage="pre-build",
            check_name="Goal section",
            details="Goal section is missing from workstream file",
        )
```

### Coverage Errors

Raise when test coverage is insufficient:

```python
from sdp.errors import CoverageTooLowError

def check_coverage(module: str, required_pct: float) -> None:
    actual_pct = get_coverage(module)

    if actual_pct < required_pct:
        missing = get_uncovered_files(module)
        raise CoverageTooLowError(
            coverage_pct=actual_pct,
            required_pct=required_pct,
            module=module,
            missing_files=missing,
        )
```

### Test Failures

Raise when tests fail:

```python
from sdp.errors import TestFailureError

def run_tests(test_command: str) -> None:
    result = subprocess.run(test_command, capture_output=True)

    if result.returncode != 0:
        failed = parse_failed_tests(result.stdout)
        raise TestFailureError(
            test_command=test_command,
            failed_tests=failed,
            total_tests=result.total_tests,
            passed_tests=result.passed_tests,
        )
```

---

## Error Categories

Choose the appropriate category for your error:

| Category | Use When... | Example |
|----------|-------------|---------|
| **validation** | Input/data validation fails | Invalid WS format, missing sections |
| **build** | Build checks fail | Pre/post-build validation |
| **test** | Test execution fails | Unit tests fail, low coverage |
| **configuration** | Config is invalid | Bad TOML, missing keys |
| **dependency** | Dependencies not met | WS dependency not completed |
| **hook** | Git/build hooks fail | Pre-commit, post-build hooks |
| **artifact** | Output validation fails | Code quality, documentation |
| **beads** | Beads operations fail | Task not found, server down |
| **coverage** | Coverage thresholds not met | Coverage below 80% |

---

## Predefined Error Types

### SDPError (Base)

Use for custom error types:

```python
from sdp.errors import SDPError, ErrorCategory

raise SDPError(
    category=ErrorCategory.VALIDATION,
    message="Custom validation failed",
    remediation="Check your input and try again",
    docs_url="https://docs.sdp.dev/custom",
    context={"input": input_value},
)
```

### BeadsNotFoundError

Beads task not found:

```python
from sdp.errors import BeadsNotFoundError

raise BeadsNotFoundError(
    task_id="TASK-001",
    search_paths=["/path1", "/path2"],
)
```

### CoverageTooLowError

Test coverage below threshold:

```python
from sdp.errors import CoverageTooLowError

raise CoverageTooLowError(
    coverage_pct=65.5,
    required_pct=80.0,
    module="sdp.core",
    missing_files=["parser.py", "validator.py"],
)
```

### QualityGateViolationError

Quality gate check failed:

```python
from sdp.errors import QualityGateViolationError

raise QualityGateViolationError(
    gate_name="file_size",
    violations=["module.py: 250 lines (max: 200)"],
    severity="warning",
)
```

### WorkstreamValidationError

Workstream validation failed:

```python
from sdp.errors import WorkstreamValidationError

raise WorkstreamValidationError(
    ws_id="WS-001-01",
    errors=["Missing Goal", "No Acceptance Criteria"],
    file_path="docs/workstreams/backlog/WS-001-01.md",
)
```

### ConfigurationError

Configuration validation failed:

```python
from sdp.errors import ConfigurationError

raise ConfigurationError(
    config_file="quality-gate.toml",
    errors=["Invalid TOML syntax"],
    missing_keys=["max_lines", "min_coverage"],
)
```

### DependencyNotFoundError

Workstream dependency not found:

```python
from sdp.errors import DependencyNotFoundError

raise DependencyNotFoundError(
    dependency="WS-001-01",
    ws_id="WS-001-02",
    available_ws=["WS-001-01", "WS-001-03"],
)
```

### HookExecutionError

Git/build hook failed:

```python
from sdp.errors import HookExecutionError

raise HookExecutionError(
    hook_name="pre-commit",
    stage="pre-commit",
    output="Time estimates found",
    exit_code=1,
)
```

### TestFailureError

Test execution failed:

```python
from sdp.errors import TestFailureError

raise TestFailureError(
    test_command="pytest",
    failed_tests=["test_foo", "test_bar"],
    total_tests=10,
    passed_tests=8,
)
```

### BuildValidationError

Build validation failed:

```python
from sdp.errors import BuildValidationError

raise BuildValidationError(
    ws_id="WS-001-01",
    stage="pre-build",
    check_name="Goal section",
    details="Goal section is missing",
)
```

### ArtifactValidationError

Artifact validation failed:

```python
from sdp.errors import ArtifactValidationError

raise ArtifactValidationError(
    artifact_type="code",
    artifact_path="src/sdp/module.py",
    errors=["File too large", "Missing type hints"],
)
```

---

## Custom Errors

Create custom error types by inheriting from `SDPError`:

```python
from sdp.errors import SDPError, ErrorCategory

class MyCustomError(SDPError):
    """Custom error for my specific use case."""

    def __init__(
        self,
        resource_id: str,
        issue: str,
    ) -> None:
        """Initialize custom error.

        Args:
            resource_id: ID of the resource that failed
            issue: Description of the issue
        """
        super().__init__(
            category=ErrorCategory.VALIDATION,
            message=f"Resource '{resource_id}' failed: {issue}",
            remediation=(
                "1. Check resource configuration\n"
                "2. Verify resource exists\n"
                "3. Review resource logs"
            ),
            docs_url="https://docs.sdp.dev/custom#resource",
            context={
                "resource_id": resource_id,
                "issue": issue,
            },
        )
```

---

## Best Practices

### DO ✅

1. **Use specific error types** - Choose the most specific predefined error
2. **Provide clear remediation** - Give actionable steps to fix the error
3. **Include relevant context** - Add debugging information to context dict
4. **Link to documentation** - Add docs_url for complex errors
5. **Format for terminal** - Use `format_error_for_terminal()` for display

```python
from sdp.errors import CoverageTooLowError, format_error_for_terminal

try:
    check_coverage(module)
except CoverageTooLowError as e:
    print(format_error_for_terminal(e))
```

6. **Format for JSON** - Use `format_error_for_json()` for APIs/CLI tools

```python
from sdp.errors import format_error_for_json
import json

try:
    validate_workstream(ws_id)
except SDPError as e:
    error_json = format_error_for_json(e)
    print(json.dumps(error_json, indent=2))
```

### DON'T ❌

1. **Don't raise bare exceptions** - Always use specific SDP error types
2. **Don't omit remediation** - Always provide how to fix the error
3. **Don't include sensitive data** - Don't add passwords/tokens to context
4. **Don't create vague errors** - Be specific about what failed
5. **Don't duplicate existing errors** - Use predefined types when possible

---

## Examples

### CLI Error Handling

```python
import click
from sdp.errors import format_error_for_terminal

@click.command()
@click.argument("ws_file")
def parse_ws(ws_file: str) -> None:
    """Parse a workstream file."""
    try:
        ws = parse_workstream(Path(ws_file))
        click.echo(f"Parsed: {ws.ws_id}")
    except SDPError as e:
        click.echo(format_error_for_terminal(e), err=True)
        raise click.Abort()
```

### API Error Handling

```python
from fastapi import HTTPException
from sdp.errors import format_error_for_json

@app.post("/validate")
async def validate_workstream(ws_id: str):
    """Validate a workstream."""
    try:
        validate(ws_id)
        return {"status": "ok"}
    except SDPError as e:
        error_json = format_error_for_json(e)
        raise HTTPException(
            status_code=400,
            detail=error_json,
        )
```

### Script Error Handling

```python
from sdp.errors import WorkstreamValidationError

def main():
    """Main script entry point."""
    try:
        validate_all_workstreams()
    except WorkstreamValidationError as e:
        logger.error(str(e))
        sys.exit(1)
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        sys.exit(2)
```

### Testing Error Handling

```python
import pytest
from sdp.errors import CoverageTooLowError

def test_coverage_check():
    """Test coverage check raises error."""
    with pytest.raises(CoverageTooLowError) as exc_info:
        check_coverage("module", 80.0)

    error = exc_info.value
    assert error.context["actual_coverage"] == "65.5%"
    assert "65.5%" in str(error)
```

---

## Migration Guide

### Before (Pattern 1: Bare exceptions)

```python
# ❌ Old way
if not ws_file.exists():
    raise FileNotFoundError(f"WS file not found: {ws_file}")
```

### After (Pattern 1: Structured errors)

```python
# ✅ New way
from sdp.errors import WorkstreamValidationError

if not ws_file.exists():
    raise WorkstreamValidationError(
        ws_id=ws_id,
        errors=[f"WS file not found: {ws_file}"],
        file_path=str(ws_file),
    )
```

### Before (Pattern 2: Generic errors)

```python
# ❌ Old way
if coverage < 80:
    raise Exception(f"Coverage too low: {coverage}%")
```

### After (Pattern 2: Specific errors)

```python
# ✅ New way
from sdp.errors import CoverageTooLowError

if coverage < 80:
    raise CoverageTooLowError(
        coverage_pct=coverage,
        required_pct=80.0,
        module=module_name,
    )
```

---

## Resources

- [Troubleshooting Guide](./troubleshooting.md)
- [Quality Gates](./quality-gates.md)
- [Workstreams](./workstreams.md)
- [API Reference](../api/errors.md)

---

**Last Updated:** 2025-01-29
**SDP Version:** 0.3.0
