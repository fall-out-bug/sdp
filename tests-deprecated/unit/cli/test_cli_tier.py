"""Tests for CLI tier commands."""

import json
from pathlib import Path
from unittest.mock import patch

import pytest
from click.testing import CliRunner

from sdp.cli.tier import tier, validate_tier
from sdp.validators.capability_tier_models import (
    CapabilityTier,
    ValidationCheck,
    ValidationResult,
)


class TestTierGroup:
    """Test tier command group."""

    def test_tier_group_exists(self) -> None:
        """Test that tier command group exists."""
        runner = CliRunner()
        result = runner.invoke(tier, ["--help"])
        assert result.exit_code == 0
        assert "Tier management" in result.output

    def test_tier_validate_command_exists(self) -> None:
        """Test that tier validate command exists."""
        runner = CliRunner()
        result = runner.invoke(tier, ["validate", "--help"])
        assert result.exit_code == 0
        assert "Validate workstream" in result.output


class TestValidateTierCommand:
    """Test validate_tier command execution."""

    def test_validate_tier_t0_success(self, tmp_path: Path) -> None:
        """Test tier validation with T0 tier that passes."""
        ws_file = tmp_path / "test_ws.md"
        ws_file.write_text("""---
ws_id: 00-001-01
tier: T0
---

# Test Workstream

## Goal
Test goal
""")

        mock_result = ValidationResult(
            tier=CapabilityTier.T0,
            passed=True,
            checks=[
                ValidationCheck(
                    name="Test Check",
                    passed=True,
                    message="Check passed",
                    details=[],
                ),
            ],
        )

        runner = CliRunner()
        with patch("sdp.validators.validate_workstream_tier", return_value=mock_result):
            result = runner.invoke(
                validate_tier,
                [str(ws_file), "--tier", "T0"],
            )

            assert result.exit_code == 0
            assert "T0-READY ✓" in result.output

    def test_validate_tier_t1_failure(self, tmp_path: Path) -> None:
        """Test tier validation with T1 tier that fails."""
        ws_file = tmp_path / "test_ws.md"
        ws_file.write_text("""---
ws_id: 00-001-01
tier: T1
---

# Test Workstream
""")

        mock_result = ValidationResult(
            tier=CapabilityTier.T1,
            passed=False,
            checks=[
                ValidationCheck(
                    name="Test Check",
                    passed=False,
                    message="Check failed",
                    details=["Detail 1", "Detail 2"],
                ),
            ],
        )

        runner = CliRunner()
        with patch("sdp.validators.validate_workstream_tier", return_value=mock_result):
            result = runner.invoke(
                validate_tier,
                [str(ws_file), "--tier", "T1"],
            )

            assert result.exit_code == 1
            assert "T1-READY ✗" in result.output
            assert "Failed checks: 1/1" in result.output

    def test_validate_tier_json_output(self, tmp_path: Path) -> None:
        """Test tier validation with JSON output."""
        ws_file = tmp_path / "test_ws.md"
        ws_file.write_text("""---
ws_id: 00-001-01
tier: T2
---

# Test Workstream
""")

        mock_result = ValidationResult(
            tier=CapabilityTier.T2,
            passed=True,
            checks=[
                ValidationCheck(
                    name="Test Check",
                    passed=True,
                    message="Check passed",
                    details=["Detail"],
                ),
            ],
        )

        runner = CliRunner()
        with patch("sdp.validators.validate_workstream_tier", return_value=mock_result):
            result = runner.invoke(
                validate_tier,
                [str(ws_file), "--tier", "T2", "--json"],
            )

            assert result.exit_code == 0
            output_json = json.loads(result.output)
            assert output_json["tier"] == "T2"
            assert output_json["passed"] is True
            assert len(output_json["checks"]) == 1

    def test_validate_tier_json_output_failure(self, tmp_path: Path) -> None:
        """Test tier validation JSON output when validation fails."""
        ws_file = tmp_path / "test_ws.md"
        ws_file.write_text("""---
ws_id: 00-001-01
tier: T3
---

# Test Workstream
""")

        mock_result = ValidationResult(
            tier=CapabilityTier.T3,
            passed=False,
            checks=[
                ValidationCheck(
                    name="Test Check",
                    passed=False,
                    message="Check failed",
                    details=[],
                ),
            ],
        )

        runner = CliRunner()
        with patch("sdp.validators.validate_workstream_tier", return_value=mock_result):
            result = runner.invoke(
                validate_tier,
                [str(ws_file), "--tier", "T3", "--json"],
            )

            assert result.exit_code == 1
            output_json = json.loads(result.output)
            assert output_json["tier"] == "T3"
            assert output_json["passed"] is False

    def test_validate_tier_case_insensitive(self, tmp_path: Path) -> None:
        """Test tier validation accepts case-insensitive tier names."""
        ws_file = tmp_path / "test_ws.md"
        ws_file.write_text("""---
ws_id: 00-001-01
tier: T0
---

# Test Workstream
""")

        mock_result = ValidationResult(
            tier=CapabilityTier.T0,
            passed=True,
            checks=[],
        )

        runner = CliRunner()
        with patch("sdp.validators.validate_workstream_tier", return_value=mock_result):
            # Test lowercase
            result = runner.invoke(
                validate_tier,
                [str(ws_file), "--tier", "t0"],
            )
            assert result.exit_code == 0

            # Test mixed case
            result = runner.invoke(
                validate_tier,
                [str(ws_file), "--tier", "T0"],
            )
            assert result.exit_code == 0

    def test_validate_tier_value_error(self, tmp_path: Path) -> None:
        """Test tier validation handles ValueError."""
        ws_file = tmp_path / "test_ws.md"
        ws_file.write_text("Invalid content")

        runner = CliRunner()
        with patch(
            "sdp.validators.validate_workstream_tier",
            side_effect=ValueError("Invalid tier"),
        ):
            result = runner.invoke(
                validate_tier,
                [str(ws_file), "--tier", "T0"],
            )

            assert result.exit_code == 1

    def test_validate_tier_generic_exception(self, tmp_path: Path) -> None:
        """Test tier validation handles generic exceptions."""
        ws_file = tmp_path / "test_ws.md"
        ws_file.write_text("Test content")

        runner = CliRunner()
        with patch(
            "sdp.validators.validate_workstream_tier",
            side_effect=Exception("Unexpected error"),
        ):
            result = runner.invoke(
                validate_tier,
                [str(ws_file), "--tier", "T0"],
            )

            assert result.exit_code == 1

    def test_validate_tier_multiple_checks(self, tmp_path: Path) -> None:
        """Test tier validation with multiple checks."""
        ws_file = tmp_path / "test_ws.md"
        ws_file.write_text("""---
ws_id: 00-001-01
tier: T2
---

# Test Workstream
""")

        mock_result = ValidationResult(
            tier=CapabilityTier.T2,
            passed=False,
            checks=[
                ValidationCheck(
                    name="Check 1",
                    passed=True,
                    message="Passed",
                    details=[],
                ),
                ValidationCheck(
                    name="Check 2",
                    passed=False,
                    message="Failed",
                    details=["Detail"],
                ),
                ValidationCheck(
                    name="Check 3",
                    passed=False,
                    message="Also failed",
                    details=[],
                ),
            ],
        )

        runner = CliRunner()
        with patch("sdp.validators.validate_workstream_tier", return_value=mock_result):
            result = runner.invoke(
                validate_tier,
                [str(ws_file), "--tier", "T2"],
            )

            assert result.exit_code == 1
            assert "Failed checks: 2/3" in result.output
            assert "Check 1" in result.output
            assert "Check 2" in result.output
            assert "Check 3" in result.output
