"""Tests for capability_tier_checks_interface validator."""

import pytest

from sdp.validators.capability_tier_checks_interface import (
    _check_placeholders_present,
    _check_signatures_complete,
    _check_tests_complete,
)


class TestCheckSignaturesComplete:
    """Tests for _check_signatures_complete."""

    def test_complete_signatures_passes(self) -> None:
        """Complete function signatures pass."""
        code = '''
def foo(x: int) -> str:
    """Docstring."""
    raise NotImplementedError
'''
        result = _check_signatures_complete(code)
        assert result.passed

    def test_pass_body_fails(self) -> None:
        """Function with only pass fails."""
        code = "def foo():\n    pass"
        result = _check_signatures_complete(code)
        assert not result.passed
        assert "incomplete" in result.message.lower()

    def test_ellipsis_body_fails(self) -> None:
        """Function with only ... fails."""
        code = "def bar():\n    ..."
        result = _check_signatures_complete(code)
        assert not result.passed


class TestCheckTestsComplete:
    """Tests for _check_tests_complete."""

    def test_concrete_tests_passes(self) -> None:
        """Tests with implementation pass."""
        code = "def test_foo():\n    assert 1 == 1"
        result = _check_tests_complete(code)
        assert result.passed

    def test_todo_fails(self) -> None:
        """Tests with TODO fail."""
        code = "def test_foo():\n    # TODO implement"
        result = _check_tests_complete(code)
        assert not result.passed
        assert "TODO" in str(result.details or [])

    def test_skip_decorator_fails(self) -> None:
        """Tests with skip decorator fail."""
        code = "@pytest.mark.skip\ndef test_foo(): pass"
        result = _check_tests_complete(code)
        assert not result.passed

    def test_empty_test_body_fails(self) -> None:
        """Tests with pass body fail."""
        code = "def test_foo():\n    pass"
        result = _check_tests_complete(code)
        assert not result.passed


class TestCheckPlaceholdersPresent:
    """Tests for _check_placeholders_present."""

    def test_not_implemented_passes(self) -> None:
        """Interface with NotImplementedError passes."""
        code = "def foo():\n    raise NotImplementedError"
        result = _check_placeholders_present(code)
        assert result.passed

    def test_no_placeholder_fails(self) -> None:
        """Interface without NotImplementedError fails."""
        code = "def foo():\n    return 1"
        result = _check_placeholders_present(code)
        assert not result.passed
        assert "NotImplementedError" in result.message
