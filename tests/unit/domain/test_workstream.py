"""Tests for workstream domain entities."""

import pytest
from pathlib import Path

from sdp.domain.workstream import (
    AcceptanceCriterion,
    Workstream,
    WorkstreamID,
    WorkstreamSize,
    WorkstreamStatus,
)


class TestWorkstreamID:
    """Tests for WorkstreamID value object."""

    def test_parse_new_format(self) -> None:
        """Parse PP-FFF-SS format."""
        ws_id = WorkstreamID.parse("00-001-01")
        assert ws_id.project_id == 0
        assert ws_id.feature_id == 1
        assert ws_id.sequence == 1

    def test_parse_legacy_format(self) -> None:
        """Parse WS-FFF-SS format (legacy)."""
        ws_id = WorkstreamID.parse("WS-500-01")
        assert ws_id.project_id == 0  # Assumes SDP
        assert ws_id.feature_id == 500
        assert ws_id.sequence == 1

    def test_str_representation(self) -> None:
        """String representation uses PP-FFF-SS format."""
        ws_id = WorkstreamID(project_id=2, feature_id=42, sequence=7)
        assert str(ws_id) == "02-042-07"

    def test_parse_invalid_format(self) -> None:
        """Invalid format raises ValueError."""
        with pytest.raises(ValueError, match="Invalid WS ID format"):
            WorkstreamID.parse("invalid")

    def test_is_sdp(self) -> None:
        """Test project type checks."""
        assert WorkstreamID(0, 1, 1).is_sdp
        assert not WorkstreamID(2, 1, 1).is_sdp

    def test_is_hw_checker(self) -> None:
        """Test hw_checker project check."""
        assert WorkstreamID(2, 1, 1).is_hw_checker
        assert not WorkstreamID(0, 1, 1).is_hw_checker

    def test_validate_project_id_success(self) -> None:
        """Valid project ID passes validation."""
        ws_id = WorkstreamID(0, 1, 1)
        ws_id.validate_project_id()  # Should not raise

    def test_validate_project_id_failure(self) -> None:
        """Invalid project ID raises ValueError."""
        ws_id = WorkstreamID(99, 1, 1)
        with pytest.raises(ValueError, match="Invalid project_id: 99"):
            ws_id.validate_project_id()

    def test_immutability(self) -> None:
        """WorkstreamID is immutable (frozen dataclass)."""
        ws_id = WorkstreamID(0, 1, 1)
        with pytest.raises(AttributeError):
            ws_id.project_id = 2  # type: ignore


class TestAcceptanceCriterion:
    """Tests for AcceptanceCriterion."""

    def test_creation(self) -> None:
        """Create acceptance criterion."""
        ac = AcceptanceCriterion(id="AC1", description="Must do X")
        assert ac.id == "AC1"
        assert ac.description == "Must do X"
        assert not ac.checked

    def test_checked_flag(self) -> None:
        """Checked flag defaults to False."""
        ac = AcceptanceCriterion(id="AC1", description="Test", checked=True)
        assert ac.checked


class TestWorkstream:
    """Tests for Workstream entity."""

    def test_minimal_workstream(self) -> None:
        """Create workstream with minimal fields."""
        ws = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
        )
        assert ws.ws_id == "00-001-01"
        assert ws.feature == "F001"
        assert ws.status == WorkstreamStatus.BACKLOG
        assert ws.size == WorkstreamSize.SMALL

    def test_full_workstream(self) -> None:
        """Create workstream with all fields."""
        ws = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.ACTIVE,
            size=WorkstreamSize.MEDIUM,
            github_issue=123,
            assignee="dev",
            title="Test WS",
            goal="Test goal",
            acceptance_criteria=[
                AcceptanceCriterion("AC1", "Criterion 1"),
            ],
            context="Test context",
            dependencies=["00-001-00"],
            steps=["Step 1", "Step 2"],
            code_blocks=["code"],
            file_path=Path("/tmp/ws.md"),
        )
        assert ws.title == "Test WS"
        assert ws.goal == "Test goal"
        assert len(ws.acceptance_criteria) == 1
        assert ws.dependencies == ["00-001-00"]
        assert ws.file_path == Path("/tmp/ws.md")

    def test_default_values(self) -> None:
        """Check default values for optional fields."""
        ws = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
        )
        assert ws.github_issue is None
        assert ws.assignee is None
        assert ws.title == ""
        assert ws.goal == ""
        assert ws.acceptance_criteria == []
        assert ws.context == ""
        assert ws.dependencies == []
        assert ws.steps == []
        assert ws.code_blocks == []
        assert ws.file_path is None
