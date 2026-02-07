"""Tests for capability tier validation data models."""

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

    def test_tier_iteration(self) -> None:
        """Test that we can iterate over all tiers."""
        tiers = list(CapabilityTier)
        assert len(tiers) == 4
        assert CapabilityTier.T0 in tiers
        assert CapabilityTier.T3 in tiers


class TestValidationCheck:
    """Test ValidationCheck dataclass."""

    def test_create_check_passed(self) -> None:
        """Test creating a passed check."""
        check = ValidationCheck(
            name="test_check",
            passed=True,
            message="Check passed",
        )
        assert check.name == "test_check"
        assert check.passed is True
        assert check.message == "Check passed"
        assert check.details == []

    def test_create_check_with_details(self) -> None:
        """Test creating a check with details."""
        check = ValidationCheck(
            name="test_check",
            passed=False,
            message="Check failed",
            details=["Error 1", "Error 2"],
        )
        assert check.passed is False
        assert len(check.details) == 2
        assert "Error 1" in check.details

    def test_create_check_with_list_factory(self) -> None:
        """Test that details list factory works correctly."""
        check1 = ValidationCheck(
            name="check1",
            passed=True,
            message="Message",
        )
        check2 = ValidationCheck(
            name="check2",
            passed=True,
            message="Message",
        )
        # Details should be separate lists
        check1.details.append("Detail")
        assert "Detail" in check1.details
        assert "Detail" not in check2.details


class TestValidationResult:
    """Test ValidationResult dataclass."""

    def test_create_result_passed(self) -> None:
        """Test creating a passed result."""
        check = ValidationCheck(
            name="test",
            passed=True,
            message="OK",
        )
        result = ValidationResult(
            tier=CapabilityTier.T2,
            passed=True,
            checks=[check],
        )
        assert result.tier == CapabilityTier.T2
        assert result.passed is True
        assert len(result.checks) == 1

    def test_create_result_with_errors(self) -> None:
        """Test creating a result with errors (via failed checks)."""
        check = ValidationCheck(
            name="test",
            passed=False,
            message="Failed",
        )
        result = ValidationResult(
            tier=CapabilityTier.T1,
            passed=False,
            checks=[check],
        )
        assert result.passed is False
        assert len(result.checks) == 1

    def test_create_result_with_multiple_checks(self) -> None:
        """Test creating a result with multiple checks."""
        checks = [
            ValidationCheck(
                name=f"check{i}",
                passed=i % 2 == 0,
                message=f"Message {i}",
            )
            for i in range(5)
        ]
        result = ValidationResult(
            tier=CapabilityTier.T3,
            passed=all(check.passed for check in checks),
            checks=checks,
        )
        assert len(result.checks) == 5
        assert result.passed is False  # Not all checks pass
