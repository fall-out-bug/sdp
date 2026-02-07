"""Tests for sdp status command."""

import json
from pathlib import Path

import pytest

from sdp.cli.status.collector import StatusCollector
from sdp.cli.status.formatter import format_status_json, format_status_human
from sdp.cli.status.models import (
    BeadsStatus,
    GuardStatus,
    ProjectStatus,
    WorkstreamSummary,
)


class TestStatusCollector:
    """Tests for StatusCollector."""

    def test_empty_project(self, tmp_path: Path) -> None:
        """Test status collection on empty project."""
        collector = StatusCollector(tmp_path)
        status = collector.collect()

        assert status.in_progress == []
        assert status.ready == []
        assert status.blocked == []
        assert status.guard.active is False
        assert status.beads.available is False

    def test_with_in_progress_workstreams(self, tmp_path: Path) -> None:
        """Test status collection with workstreams in progress."""
        # Setup workstream files
        ws_dir = tmp_path / "docs" / "workstreams" / "in_progress"
        ws_dir.mkdir(parents=True)
        (ws_dir / "00-034-01-test.md").write_text(
            "---\nws_id: 00-034-01\nstatus: in_progress\n---\n# Test WS\n\nTest workstream."
        )

        collector = StatusCollector(tmp_path)
        status = collector.collect()

        assert len(status.in_progress) == 1
        assert status.in_progress[0].id == "00-034-01"
        assert status.in_progress[0].status == "in_progress"

    def test_with_ready_workstreams(self, tmp_path: Path) -> None:
        """Test status collection with ready workstreams (no blockers)."""
        ws_dir = tmp_path / "docs" / "workstreams" / "backlog"
        ws_dir.mkdir(parents=True)
        (ws_dir / "00-034-02-test.md").write_text(
            "---\nws_id: 00-034-02\nstatus: READY\n---\n# Ready WS"
        )

        collector = StatusCollector(tmp_path)
        status = collector.collect()

        assert len(status.ready) >= 1
        ready_ids = [ws.id for ws in status.ready]
        assert "00-034-02" in ready_ids

    def test_with_blocked_workstreams(self, tmp_path: Path) -> None:
        """Test status collection with blocked workstreams."""
        ws_dir = tmp_path / "docs" / "workstreams" / "backlog"
        ws_dir.mkdir(parents=True)
        (ws_dir / "00-034-03-test.md").write_text(
            "---\nws_id: 00-034-03\nstatus: BLOCKED\ndepends_on:\n  - 00-034-01\n---\n# Blocked WS"
        )

        collector = StatusCollector(tmp_path)
        status = collector.collect()

        assert len(status.blocked) >= 1
        blocked = next((ws for ws in status.blocked if ws.id == "00-034-03"), None)
        assert blocked is not None
        assert "00-034-01" in blocked.blockers

    def test_guard_status_active(self, tmp_path: Path) -> None:
        """Test guard status detection when active."""
        sdp_dir = tmp_path / ".sdp"
        sdp_dir.mkdir()
        (sdp_dir / "state.json").write_text(
            json.dumps({
                "active_workstream": "00-034-01",
                "allowed_files": ["src/sdp/cli/status/models.py"],
            })
        )

        collector = StatusCollector(tmp_path)
        status = collector.collect()

        assert status.guard.active is True
        assert status.guard.workstream_id == "00-034-01"
        assert len(status.guard.allowed_files) > 0

    def test_guard_status_inactive(self, tmp_path: Path) -> None:
        """Test guard status when no state file exists."""
        collector = StatusCollector(tmp_path)
        status = collector.collect()

        assert status.guard.active is False
        assert status.guard.workstream_id is None

    def test_beads_not_available(self, tmp_path: Path) -> None:
        """Test Beads status when .beads directory doesn't exist."""
        collector = StatusCollector(tmp_path)
        status = collector.collect()

        assert status.beads.available is False
        assert status.beads.synced is False

    def test_beads_available(self, tmp_path: Path) -> None:
        """Test Beads status when available."""
        beads_dir = tmp_path / ".beads"
        beads_dir.mkdir()
        (beads_dir / "metadata.json").write_text(
            json.dumps({
                "last_sync": "2026-01-31T10:00:00Z",
            })
        )

        collector = StatusCollector(tmp_path)
        status = collector.collect()

        assert status.beads.available is True

    def test_next_actions_suggestions(self, tmp_path: Path) -> None:
        """Test that next actions are suggested."""
        collector = StatusCollector(tmp_path)
        status = collector.collect()

        # Should always have some suggestions
        assert isinstance(status.next_actions, list)


class TestStatusFormatters:
    """Tests for status formatters."""

    def test_json_formatter_valid_output(self) -> None:
        """Test JSON formatter produces valid JSON."""
        status = ProjectStatus(
            in_progress=[],
            blocked=[],
            ready=[],
            guard=GuardStatus(active=False),
            beads=BeadsStatus(available=False, synced=False),
            next_actions=["Run @build to start"],
        )

        output = format_status_json(status)
        parsed = json.loads(output)

        assert "in_progress" in parsed
        assert "blocked" in parsed
        assert "ready" in parsed
        assert "guard" in parsed
        assert "beads" in parsed
        assert "next_actions" in parsed
        assert parsed["next_actions"] == ["Run @build to start"]

    def test_json_formatter_with_workstreams(self) -> None:
        """Test JSON formatter with workstreams."""
        status = ProjectStatus(
            in_progress=[
                WorkstreamSummary(
                    id="00-034-01",
                    title="Test WS",
                    status="in_progress",
                    scope="MEDIUM",
                )
            ],
            blocked=[],
            ready=[],
            guard=GuardStatus(active=True, workstream_id="00-034-01"),
            beads=BeadsStatus(available=False, synced=False),
            next_actions=[],
        )

        output = format_status_json(status)
        parsed = json.loads(output)

        assert len(parsed["in_progress"]) == 1
        assert parsed["in_progress"][0]["id"] == "00-034-01"
        assert parsed["guard"]["active"] is True

    def test_human_formatter_produces_output(self) -> None:
        """Test human formatter produces readable output."""
        status = ProjectStatus(
            in_progress=[],
            blocked=[],
            ready=[],
            guard=GuardStatus(active=False),
            beads=BeadsStatus(available=False, synced=False),
            next_actions=["Run @build to start"],
        )

        output = format_status_human(status)

        # Check for key sections
        assert "SDP Project Status" in output
        assert "Guard" in output or "guard" in output.lower()
        assert "Suggested Actions" in output or "next" in output.lower()

    def test_human_formatter_with_workstreams(self) -> None:
        """Test human formatter with workstreams."""
        status = ProjectStatus(
            in_progress=[
                WorkstreamSummary(
                    id="00-034-01",
                    title="Test WS",
                    status="in_progress",
                    scope="MEDIUM",
                )
            ],
            blocked=[
                WorkstreamSummary(
                    id="00-034-02",
                    title="Blocked WS",
                    status="BLOCKED",
                    scope="MEDIUM",
                    blockers=["00-034-01"],
                )
            ],
            ready=[
                WorkstreamSummary(
                    id="00-034-03",
                    title="Ready WS",
                    status="READY",
                    scope="SMALL",
                )
            ],
            guard=GuardStatus(active=True, workstream_id="00-034-01"),
            beads=BeadsStatus(available=True, synced=True, ready_tasks=["bd-0001"]),
            next_actions=["Complete 00-034-01"],
        )

        output = format_status_human(status)

        # Check workstreams are mentioned
        assert "00-034-01" in output
        assert "00-034-02" in output
        assert "00-034-03" in output

    def test_human_formatter_verbose(self) -> None:
        """Test human formatter in verbose mode."""
        status = ProjectStatus(
            in_progress=[],
            blocked=[],
            ready=[],
            guard=GuardStatus(active=False),
            beads=BeadsStatus(available=False, synced=False),
            next_actions=[],
        )

        output = format_status_human(status, verbose=True)
        assert isinstance(output, str)
        assert len(output) > 0


class TestStatusModels:
    """Tests for status data models."""

    def test_workstream_summary_creation(self) -> None:
        """Test WorkstreamSummary dataclass."""
        ws = WorkstreamSummary(
            id="00-034-01",
            title="Test",
            status="READY",
            scope="MEDIUM",
        )
        assert ws.id == "00-034-01"
        assert ws.blockers == []

    def test_workstream_summary_with_blockers(self) -> None:
        """Test WorkstreamSummary with blockers."""
        ws = WorkstreamSummary(
            id="00-034-02",
            title="Test",
            status="BLOCKED",
            scope="MEDIUM",
            blockers=["00-034-01"],
        )
        assert len(ws.blockers) == 1

    def test_guard_status_creation(self) -> None:
        """Test GuardStatus dataclass."""
        guard = GuardStatus(active=True, workstream_id="00-034-01")
        assert guard.active is True
        assert guard.workstream_id == "00-034-01"
        assert guard.allowed_files == []

    def test_beads_status_creation(self) -> None:
        """Test BeadsStatus dataclass."""
        beads = BeadsStatus(available=True, synced=True)
        assert beads.available is True
        assert beads.synced is True
        assert beads.ready_tasks == []

    def test_project_status_creation(self) -> None:
        """Test ProjectStatus dataclass."""
        status = ProjectStatus(
            in_progress=[],
            blocked=[],
            ready=[],
            guard=GuardStatus(active=False),
            beads=BeadsStatus(available=False, synced=False),
            next_actions=[],
        )
        assert isinstance(status.in_progress, list)
        assert isinstance(status.guard, GuardStatus)
        assert isinstance(status.beads, BeadsStatus)
