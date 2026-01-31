"""Tests for sdp.hooks.common."""

import pytest
from pathlib import Path

from sdp.hooks.common import CheckResult


def test_check_result_passed_format() -> None:
    """CheckResult formats passed result correctly."""
    result = CheckResult(passed=True, message="All good", violations=[])
    assert "✓" in result.format_terminal()
    assert "All good" in result.format_terminal()


def test_check_result_failed_format() -> None:
    """CheckResult formats failed result with violations."""
    result = CheckResult(
        passed=False,
        message="Time estimates found",
        violations=[
            (Path("foo.py"), 10, "2 hours"),
            (Path("bar.md"), None, "soon"),
        ],
    )
    output = result.format_terminal()
    assert "❌" in output
    assert "Time estimates found" in output
    assert "foo.py:10" in output
    assert "2 hours" in output
    assert "bar.md:?" in output
    assert "soon" in output
