"""Interface and test validation checks for workstream contracts.

Provides validation functions for interface signatures and test completeness.
"""

import re

from sdp.validators.capability_tier_models import ValidationCheck


def _check_signatures_complete(interface_code: str) -> ValidationCheck:
    """Check interface code has complete function signatures (no pass)."""
    # Look for functions with only "pass" or "..." as body
    incomplete_pattern = r"def\s+\w+\([^)]*\):\s*\n\s+(pass|\.\.\.)"

    if re.search(incomplete_pattern, interface_code):
        return ValidationCheck(
            name="signatures_complete",
            passed=False,
            message="Interface has incomplete function signatures (pass/...)",
            details=["All functions must have complete type hints and docstrings"],
        )

    return ValidationCheck(
        name="signatures_complete",
        passed=True,
        message="All function signatures are complete",
    )


def _check_tests_complete(tests_code: str) -> ValidationCheck:
    """Check tests are concrete (no TODO/FIXME/skip)."""
    issues = []

    # Check for TODO/FIXME
    if re.search(r"(TODO|FIXME)", tests_code, re.IGNORECASE):
        issues.append("Tests contain TODO/FIXME markers")

    # Check for skip decorators
    if re.search(r"@pytest\.mark\.skip|@skip", tests_code):
        issues.append("Tests have skip decorators")

    # Check for empty test bodies
    empty_test = re.search(r"def test_\w+\([^)]*\):\s*(pass|#\s*TODO)", tests_code)
    if empty_test:
        issues.append("Tests have empty bodies (pass/TODO)")

    if issues:
        return ValidationCheck(
            name="tests_complete",
            passed=False,
            message="Tests are incomplete",
            details=issues,
        )

    return ValidationCheck(
        name="tests_complete",
        passed=True,
        message="All tests are concrete and complete",
    )


def _check_placeholders_present(interface_code: str) -> ValidationCheck:
    """Check interface uses NotImplementedError placeholders."""
    # Count functions
    func_count = len(re.findall(r"def\s+\w+", interface_code))

    # Count NotImplementedError
    not_impl_count = len(re.findall(r"raise NotImplementedError", interface_code))

    if func_count > 0 and not_impl_count == 0:
        return ValidationCheck(
            name="placeholders_present",
            passed=False,
            message="Interface functions should have NotImplementedError placeholders",
            details=[f"Found {func_count} functions, 0 with NotImplementedError"],
        )

    return ValidationCheck(
        name="placeholders_present",
        passed=True,
        message=f"Interface uses NotImplementedError placeholders ({not_impl_count}/{func_count})",
    )
