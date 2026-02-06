"""Unit tests for traceability models."""

import pytest

from sdp.traceability.models import (
    ACTestMapping,
    MappingStatus,
    TraceabilityReport,
)


class TestACTestMapping:
    """Tests for ACTestMapping model."""

    def test_to_dict(self) -> None:
        """AC1: Serialization to dict."""
        mapping = ACTestMapping(
            ac_id="AC1",
            ac_description="User can login",
            test_file="tests/test_auth.py",
            test_name="test_login",
            status=MappingStatus.MAPPED,
        )

        d = mapping.to_dict()

        assert d["ac_id"] == "AC1"
        assert d["ac_description"] == "User can login"
        assert d["test_file"] == "tests/test_auth.py"
        assert d["test_name"] == "test_login"
        assert d["status"] == "mapped"
        assert d["confidence"] == 1.0

    def test_from_dict(self) -> None:
        """AC1: Deserialization from dict."""
        data = {
            "ac_id": "AC1",
            "ac_description": "User can login",
            "test_file": "tests/test_auth.py",
            "test_name": "test_login",
            "status": "mapped",
        }

        mapping = ACTestMapping.from_dict(data)

        assert mapping.ac_id == "AC1"
        assert mapping.ac_description == "User can login"
        assert mapping.test_file == "tests/test_auth.py"
        assert mapping.test_name == "test_login"
        assert mapping.status == MappingStatus.MAPPED
        assert mapping.confidence == 1.0

    def test_from_dict_with_confidence(self) -> None:
        """AC1: Deserialization with confidence score."""
        data = {
            "ac_id": "AC2",
            "ac_description": "User can logout",
            "test_file": None,
            "test_name": None,
            "status": "missing",
            "confidence": 0.8,
        }

        mapping = ACTestMapping.from_dict(data)

        assert mapping.ac_id == "AC2"
        assert mapping.confidence == 0.8


class TestTraceabilityReport:
    """Tests for TraceabilityReport model."""

    def test_coverage_pct_complete(self) -> None:
        """AC2: Coverage percentage when all mapped."""
        report = TraceabilityReport(
            ws_id="00-032-01",
            mappings=[
                ACTestMapping("AC1", "Test 1", "test.py", "test_1", MappingStatus.MAPPED),
                ACTestMapping("AC2", "Test 2", "test.py", "test_2", MappingStatus.MAPPED),
            ],
        )

        assert report.coverage_pct == 100.0
        assert report.is_complete is True

    def test_coverage_pct_partial(self) -> None:
        """AC2: Coverage percentage with partial mapping."""
        report = TraceabilityReport(
            ws_id="00-032-01",
            mappings=[
                ACTestMapping("AC1", "Test 1", "test.py", "test_1", MappingStatus.MAPPED),
                ACTestMapping("AC2", "Test 2", "test.py", "test_2", MappingStatus.MAPPED),
                ACTestMapping("AC3", "Test 3", None, None, MappingStatus.MISSING),
            ],
        )

        assert report.coverage_pct == pytest.approx(66.67, rel=0.01)
        assert report.is_complete is False

    def test_is_complete_true(self) -> None:
        """AC2: Report is complete when all ACs mapped."""
        complete = TraceabilityReport(
            ws_id="00-032-01",
            mappings=[
                ACTestMapping("AC1", "Test", "test.py", "test_1", MappingStatus.MAPPED),
            ],
        )

        assert complete.is_complete is True
        assert complete.missing_acs == 0

    def test_is_complete_false(self) -> None:
        """AC2: Report is incomplete with missing ACs."""
        incomplete = TraceabilityReport(
            ws_id="00-032-01",
            mappings=[
                ACTestMapping("AC1", "Test", None, None, MappingStatus.MISSING),
            ],
        )

        assert incomplete.is_complete is False
        assert incomplete.missing_acs == 1

    def test_empty_report(self) -> None:
        """AC2: Empty report is considered complete."""
        report = TraceabilityReport(ws_id="00-032-01", mappings=[])

        assert report.total_acs == 0
        assert report.coverage_pct == 100.0
        assert report.is_complete is True

    def test_to_dict(self) -> None:
        """AC2: Serialization to dict."""
        report = TraceabilityReport(
            ws_id="00-032-01",
            mappings=[
                ACTestMapping("AC1", "Test 1", "test.py", "test_1", MappingStatus.MAPPED),
                ACTestMapping("AC2", "Test 2", None, None, MappingStatus.MISSING),
            ],
        )

        d = report.to_dict()

        assert d["ws_id"] == "00-032-01"
        assert d["total_acs"] == 2
        assert d["mapped_acs"] == 1
        assert d["missing_acs"] == 1
        assert d["coverage_pct"] == 50.0
        assert len(d["mappings"]) == 2

    def test_to_markdown_table(self) -> None:
        """AC2: Markdown table generation."""
        report = TraceabilityReport(
            ws_id="00-032-01",
            mappings=[
                ACTestMapping("AC1", "User can login", "test.py", "test_login", MappingStatus.MAPPED),
                ACTestMapping("AC2", "User can logout", None, None, MappingStatus.MISSING),
            ],
        )

        table = report.to_markdown_table()

        assert "| AC | Description | Test | Status |" in table
        assert "AC1" in table
        assert "AC2" in table
        assert "✅" in table
        assert "❌" in table
        assert "`test_login`" in table

    def test_long_description_truncation(self) -> None:
        """AC2: Long descriptions are truncated in table."""
        report = TraceabilityReport(
            ws_id="00-032-01",
            mappings=[
                ACTestMapping(
                    "AC1",
                    "This is a very long description that should be truncated",
                    "test.py",
                    "test_1",
                    MappingStatus.MAPPED,
                ),
            ],
        )

        table = report.to_markdown_table()

        # Should contain truncated version with ellipsis
        assert "..." in table
        # Full description should not be in table
        assert "This is a very long description that should be truncated" not in table
