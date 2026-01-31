"""Individual verification checks for workstream completion."""

import subprocess
from pathlib import Path
from typing import Any

from sdp.validators.ws_completion.models import CheckResult


def verify_output_files(ws_data: dict[str, Any]) -> list[CheckResult]:
    """Check all output files in scope exist.

    Args:
        ws_data: Parsed workstream data

    Returns:
        List of check results for each file
    """
    checks: list[CheckResult] = []
    scope_files = ws_data.get("scope_files", [])

    for file_path_str in scope_files:
        file_path = Path(file_path_str)
        exists = file_path.exists()

        checks.append(
            CheckResult(
                name=f"File: {file_path}",
                passed=exists,
                message=str(file_path) if exists else f"Missing: {file_path}",
                evidence=str(file_path.absolute()) if exists else None,
            )
        )

    return checks


def verify_commands(ws_data: dict[str, Any]) -> list[CheckResult]:
    """Run verification commands and check exit codes.

    Args:
        ws_data: Parsed workstream data

    Returns:
        List of check results for each command
    """
    checks: list[CheckResult] = []
    verification_cmds = ws_data.get("verification_commands", [])

    for cmd in verification_cmds:
        try:
            result = subprocess.run(
                cmd,
                shell=True,
                capture_output=True,
                text=True,
                timeout=60,
            )

            passed = result.returncode == 0
            evidence = result.stdout if passed else result.stderr

            checks.append(
                CheckResult(
                    name=f"Command: {cmd[:50]}...",
                    passed=passed,
                    message=f"Exit code: {result.returncode}",
                    evidence=evidence[:500] if evidence else None,
                )
            )
        except subprocess.TimeoutExpired:
            checks.append(
                CheckResult(
                    name=f"Command: {cmd[:50]}...",
                    passed=False,
                    message="Command timed out (60s)",
                    evidence=None,
                )
            )
        except Exception as e:
            checks.append(
                CheckResult(
                    name=f"Command: {cmd[:50]}...",
                    passed=False,
                    message=f"Error: {str(e)}",
                    evidence=None,
                )
            )

    return checks


def verify_coverage(ws_data: dict[str, Any]) -> CheckResult | None:
    """Check test coverage meets threshold.

    Args:
        ws_data: Parsed workstream data

    Returns:
        Check result or None if no coverage requirement
    """
    # Extract module from scope_files
    scope_files = ws_data.get("scope_files", [])
    if not scope_files:
        return None

    # Find Python files in scope
    py_files = [f for f in scope_files if f.endswith(".py") and not f.startswith("tests/")]
    if not py_files:
        return None

    # Run coverage check (basic implementation)
    try:
        result = subprocess.run(
            ["pytest", "--cov=src", "--cov-report=term", "-v"],
            capture_output=True,
            text=True,
            timeout=120,
        )

        # Parse coverage from output
        coverage = _parse_coverage_from_output(result.stdout)

        passed = coverage is not None and coverage >= 80.0

        return CheckResult(
            name="Test Coverage",
            passed=passed,
            message=f"Coverage: {coverage:.1f}%" if coverage else "Coverage not found",
            evidence=result.stdout[:500],
        )
    except Exception as e:
        return CheckResult(
            name="Test Coverage",
            passed=False,
            message=f"Error running coverage: {str(e)}",
            evidence=None,
        )


def _parse_coverage_from_output(output: str) -> float | None:
    """Parse coverage percentage from pytest output.

    Args:
        output: pytest stdout

    Returns:
        Coverage percentage or None
    """
    import re

    # Look for "TOTAL ... XX%"
    match = re.search(r"TOTAL\s+\d+\s+\d+\s+(\d+)%", output)
    if match:
        return float(match.group(1))

    return None


def extract_coverage_value(check: CheckResult | None) -> float | None:
    """Extract coverage value from check result.

    Args:
        check: Coverage check result

    Returns:
        Coverage percentage or None
    """
    if not check or not check.message:
        return None

    import re

    match = re.search(r"Coverage: ([\d.]+)%", check.message)
    if match:
        return float(match.group(1))

    return None
