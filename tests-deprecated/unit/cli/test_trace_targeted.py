"""Targeted tests for uncovered lines in cli/trace.py."""

import json
from pathlib import Path
from unittest.mock import MagicMock, patch

import pytest
from click.testing import CliRunner

from sdp.cli.trace import check_traceability, trace
from sdp.traceability.models import ACTestMapping, MappingStatus, TraceabilityReport


@pytest.fixture
def runner() -> CliRunner:
    """Create CLI runner."""
    return CliRunner()


def test_check_traceability_value_error(runner: CliRunner) -> None:
    """Test check_traceability handles ValueError (lines 36-38)."""
    with patch("sdp.cli.trace.TraceabilityService") as mock_service:
        mock_service.return_value.check_traceability.side_effect = ValueError(
            "Test error"
        )
        result = runner.invoke(trace, ["check", "00-001-01"])
        assert result.exit_code == 1
        assert "❌ Test error" in result.output


def test_check_traceability_incomplete_report(runner: CliRunner) -> None:
    """Test check_traceability exits 1 for incomplete report (lines 58-59)."""
    incomplete_report = TraceabilityReport(
        ws_id="00-001-01",
        mappings=[
            ACTestMapping(
                ac_id="AC1",
                ac_description="Test",
                test_file=None,
                test_name=None,
                status=MappingStatus.MISSING,
            )
        ],
    )

    with patch("sdp.cli.trace.TraceabilityService") as mock_service:
        mock_service.return_value.check_traceability.return_value = incomplete_report
        result = runner.invoke(trace, ["check", "00-001-01"])
        assert result.exit_code == 1
        assert "INCOMPLETE" in result.output


def test_add_mapping_value_error(runner: CliRunner) -> None:
    """Test add_mapping handles ValueError (lines 84-86)."""
    with patch("sdp.cli.trace.TraceabilityService") as mock_service:
        mock_service.return_value.add_mapping.side_effect = ValueError("Test error")
        result = runner.invoke(
            trace, ["add", "00-001-01", "--ac", "AC1", "--test", "test_foo"]
        )
        assert result.exit_code == 1
        assert "❌ Test error" in result.output


def test_auto_detect_value_error(runner: CliRunner) -> None:
    """Test auto_detect handles ValueError during check (lines 121-123)."""
    with patch("sdp.cli.trace.TraceabilityService") as mock_service:
        mock_service.return_value.check_traceability.side_effect = ValueError(
            "Test error"
        )
        result = runner.invoke(trace, ["auto", "00-001-01"])
        assert result.exit_code == 1
        assert "❌ Test error" in result.output


def test_auto_detect_test_dir_not_found(runner: CliRunner) -> None:
    """Test auto_detect with non-existent test directory (lines 130-131)."""
    report = TraceabilityReport(
        ws_id="00-001-01",
        mappings=[],
    )

    with patch("sdp.cli.trace.TraceabilityService") as mock_service:
        mock_service.return_value.check_traceability.return_value = report
        result = runner.invoke(
            trace, ["auto", "00-001-01", "--test-dir", "/nonexistent"]
        )
        assert result.exit_code == 1
        assert "not found" in result.output


def test_auto_detect_no_mappings_detected(runner: CliRunner) -> None:
    """Test auto_detect with no mappings detected (lines 136-137)."""
    from sdp.traceability.detector import ACDetector

    report = TraceabilityReport(
        ws_id="00-001-01",
        mappings=[
            ACTestMapping(
                ac_id="AC1",
                ac_description="Test",
                test_file=None,
                test_name=None,
                status=MappingStatus.MISSING,
            )
        ],
    )

    with patch("sdp.cli.trace.TraceabilityService") as mock_service:
        mock_service.return_value.check_traceability.return_value = report

        with patch("sdp.traceability.detector.ACDetector") as mock_detector:
            mock_detector.return_value.detect_all.return_value = []

            with runner.isolated_filesystem():
                Path("tests").mkdir()
                result = runner.invoke(trace, ["auto", "00-001-01"])
                assert "No mappings detected" in result.output
