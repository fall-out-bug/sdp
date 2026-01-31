"""Tests for capability tier validation models."""

import pytest

from sdp.validators.capability_tier_models import CapabilityTier, ValidationCheck, ValidationResult


class TestCapabilityTier:
    """Test CapabilityTier enum."""

    def test_tier_values(self) -> None:
        """Test that all tier values are defined."""
        assert CapabilityTier.T0.value == "T0"
        assert CapabilityTier.T1.value == "T1"
        assert CapabilityTier.T2.value == "T2"
        assert CapabilityTier.T3.value == "T3"

    def test_tier_from_string(self) -> None:
        """Test creating tier from string."""
        assert CapabilityTier("T0") == CapabilityTier.T0
        assert CapabilityTier("T1") == CapabilityTier.T1
        assert CapabilityTier("T2") == CapabilityTier.T2
        assert CapabilityTier("T3") == CapabilityTier.T3

    def test_tier_from_string_case_insensitive(self) -> None:
        """Test that tier enum is case-sensitive."""
        with pytest.raises(ValueError):
            CapabilityTier("t0")
        with pytest.raises(ValueError):
            CapabilityTier("t1")

    def test_invalid_tier(self) -> None:
        """Test that invalid tier raises ValueError."""
        with pytest.raises(ValueError):
            CapabilityTier("T4")
        with pytest.raises(ValueError):
            CapabilityTier("INVALID")


class TestValidationCheck:
    """Test ValidationCheck dataclass."""

    def test_create_passed_check(self) -> None:
        """Test creating a passed validation check."""
        check = ValidationCheck(name="test_check", passed=True, message="Check passed")
        assert check.name == "test_check"
        assert check.passed is True
        assert check.message == "Check passed"
        assert check.details == []

    def test_create_failed_check(self) -> None:
        """Test creating a failed validation check."""
        check = ValidationCheck(name="test_check", passed=False, message="Check failed")
        assert check.name == "test_check"
        assert check.passed is False
        assert check.message == "Check failed"

    def test_check_with_details(self) -> None:
        """Test creating a check with details."""
        details = ["Issue 1", "Issue 2", "Issue 3"]
        check = ValidationCheck(
            name="test_check", passed=False, message="Multiple issues", details=details
        )
        assert check.details == details
        assert len(check.details) == 3

    def test_check_empty_message(self) -> None:
        """Test creating a check with empty message (default)."""
        check = ValidationCheck(name="test_check", passed=True)
        assert check.message == ""

    def test_check_empty_details_default(self) -> None:
        """Test that details default to empty list."""
        check = ValidationCheck(name="test_check", passed=True)
        assert check.details == []
        assert isinstance(check.details, list)


class TestValidationResult:
    """Test ValidationResult dataclass."""

    def test_create_passed_result(self) -> None:
        """Test creating a passed validation result."""
        result = ValidationResult(tier=CapabilityTier.T0, passed=True)
        assert result.tier == CapabilityTier.T0
        assert result.passed is True
        assert result.checks == []

    def test_create_failed_result(self) -> None:
        """Test creating a failed validation result."""
        result = ValidationResult(tier=CapabilityTier.T1, passed=False)
        assert result.tier == CapabilityTier.T1
        assert result.passed is False

    def test_add_check_passed(self) -> None:
        """Test adding a passed check doesn't change result status."""
        result = ValidationResult(tier=CapabilityTier.T0, passed=True)
        check = ValidationCheck(name="test_check", passed=True, message="OK")
        result.add_check(check)
        assert len(result.checks) == 1
        assert result.checks[0] == check
        assert result.passed is True

    def test_add_check_failed(self) -> None:
        """Test adding a failed check changes result status to failed."""
        result = ValidationResult(tier=CapabilityTier.T0, passed=True)
        check = ValidationCheck(name="test_check", passed=False, message="Failed")
        result.add_check(check)
        assert len(result.checks) == 1
        assert result.checks[0] == check
        assert result.passed is False

    def test_add_multiple_checks(self) -> None:
        """Test adding multiple checks."""
        result = ValidationResult(tier=CapabilityTier.T2, passed=True)
        check1 = ValidationCheck(name="check1", passed=True)
        check2 = ValidationCheck(name="check2", passed=False, message="Error")
        check3 = ValidationCheck(name="check3", passed=True)
        result.add_check(check1)
        result.add_check(check2)
        result.add_check(check3)
        assert len(result.checks) == 3
        assert result.passed is False  # Failed because check2 failed

    def test_add_check_all_passed(self) -> None:
        """Test that result remains passed when all checks pass."""
        result = ValidationResult(tier=CapabilityTier.T3, passed=True)
        for i in range(5):
            check = ValidationCheck(name=f"check{i}", passed=True)
            result.add_check(check)
        assert len(result.checks) == 5
        assert result.passed is True

    def test_initially_failed_add_passed_check(self) -> None:
        """Test that adding passed check to initially failed result keeps it failed."""
        result = ValidationResult(tier=CapabilityTier.T1, passed=False)
        check = ValidationCheck(name="test_check", passed=True)
        result.add_check(check)
        assert result.passed is False  # Stays failed

    def test_add_check_with_details(self) -> None:
        """Test adding a check with detailed information."""
        result = ValidationResult(tier=CapabilityTier.T2, passed=True)
        check = ValidationCheck(
            name="complex_check",
            passed=False,
            message="Multiple issues found",
            details=["Issue 1", "Issue 2", "Issue 3"],
        )
        result.add_check(check)
        assert len(result.checks) == 1
        assert result.checks[0].details == ["Issue 1", "Issue 2", "Issue 3"]
        assert result.passed is False
