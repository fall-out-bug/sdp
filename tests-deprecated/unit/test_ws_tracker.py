"""Unit tests for WorkstreamTracker."""

import json
import tempfile
from pathlib import Path
from unittest.mock import Mock

import pytest

from sdp.beads.models import BeadsStatus
from sdp.guard.tracker import WorkstreamInProgressError, WorkstreamTracker


class TestWorkstreamTracker:
    """Test WorkstreamTracker."""

    def test_activate_updates_beads_status(self, tmp_path: Path) -> None:
        """AC1: Activate sets IN_PROGRESS in Beads."""
        client = Mock()
        mock_task = Mock()
        mock_task.sdp_metadata = {}
        client.get_task.return_value = mock_task

        state_file = tmp_path / "state.json"
        tracker = WorkstreamTracker(client, state_file)
        tracker.activate("00-032-01")

        client.update_task_status.assert_called_with("00-032-01", BeadsStatus.IN_PROGRESS)

    def test_complete_sets_closed_status(self, tmp_path: Path) -> None:
        """AC2: Complete sets CLOSED in Beads."""
        client = Mock()
        mock_task = Mock()
        mock_task.sdp_metadata = {}
        client.get_task.return_value = mock_task

        state_file = tmp_path / "state.json"
        tracker = WorkstreamTracker(client, state_file)
        tracker.activate("00-032-01")
        tracker.complete("00-032-01")

        client.update_task_status.assert_called_with("00-032-01", BeadsStatus.CLOSED)

    def test_cannot_activate_second_ws(self, tmp_path: Path) -> None:
        """AC3: Second activation raises error."""
        client = Mock()
        mock_task = Mock()
        mock_task.sdp_metadata = {}
        client.get_task.return_value = mock_task

        state_file = tmp_path / "state.json"
        tracker = WorkstreamTracker(client, state_file)
        tracker.activate("00-032-01")

        with pytest.raises(WorkstreamInProgressError):
            tracker.activate("00-032-02")

    def test_get_active_returns_none_when_no_state(self, tmp_path: Path) -> None:
        """get_active() returns None when no state file exists."""
        client = Mock()
        state_file = tmp_path / "state.json"
        tracker = WorkstreamTracker(client, state_file)

        assert tracker.get_active() is None

    def test_get_active_returns_ws_id(self, tmp_path: Path) -> None:
        """AC4: get_active() returns current WS ID."""
        client = Mock()
        mock_task = Mock()
        mock_task.sdp_metadata = {}
        client.get_task.return_value = mock_task

        state_file = tmp_path / "state.json"
        tracker = WorkstreamTracker(client, state_file)
        tracker.activate("00-032-01")

        assert tracker.get_active() == "00-032-01"

    def test_abort_returns_to_open_status(self, tmp_path: Path) -> None:
        """abort() sets status back to OPEN."""
        client = Mock()
        mock_task = Mock()
        mock_task.sdp_metadata = {}
        client.get_task.return_value = mock_task

        state_file = tmp_path / "state.json"
        tracker = WorkstreamTracker(client, state_file)
        tracker.activate("00-032-01")
        tracker.abort("00-032-01")

        client.update_task_status.assert_called_with("00-032-01", BeadsStatus.OPEN)

    def test_complete_wrong_ws_raises_error(self, tmp_path: Path) -> None:
        """complete() raises error if wrong WS ID."""
        client = Mock()
        mock_task = Mock()
        mock_task.sdp_metadata = {}
        client.get_task.return_value = mock_task

        state_file = tmp_path / "state.json"
        tracker = WorkstreamTracker(client, state_file)
        tracker.activate("00-032-01")

        with pytest.raises(ValueError, match="not active"):
            tracker.complete("00-032-02")

    def test_abort_wrong_ws_raises_error(self, tmp_path: Path) -> None:
        """abort() raises error if wrong WS ID."""
        client = Mock()
        mock_task = Mock()
        mock_task.sdp_metadata = {}
        client.get_task.return_value = mock_task

        state_file = tmp_path / "state.json"
        tracker = WorkstreamTracker(client, state_file)
        tracker.activate("00-032-01")

        with pytest.raises(ValueError, match="not active"):
            tracker.abort("00-032-02")

    def test_activate_saves_scope_files(self, tmp_path: Path) -> None:
        """activate() saves scope_files from metadata."""
        client = Mock()
        mock_task = Mock()
        mock_task.sdp_metadata = {"scope_files": ["src/test.py", "src/other.py"]}
        client.get_task.return_value = mock_task

        state_file = tmp_path / "state.json"
        tracker = WorkstreamTracker(client, state_file)
        tracker.activate("00-032-01")

        # Verify state file
        with open(state_file) as f:
            state = json.load(f)
        assert state["scope_files"] == ["src/test.py", "src/other.py"]
