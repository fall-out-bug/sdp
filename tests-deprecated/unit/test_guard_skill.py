"""Unit tests for GuardSkill."""

import pytest
from unittest.mock import Mock

from sdp.guard.skill import GuardSkill
from sdp.guard.models import GuardResult


class TestGuardSkill:
    """Test GuardSkill pre-edit checks."""

    def test_no_active_ws_blocks_edit(self) -> None:
        """AC3: No active WS returns allowed=False."""
        client = Mock()
        guard = GuardSkill(client)

        result = guard.check_edit("any/file.py")

        assert result.allowed is False
        assert "No active WS" in result.reason
        assert result.ws_id is None
        assert result.scope_files == []

    def test_file_not_in_scope_blocks_edit(self) -> None:
        """AC4: File outside scope returns allowed=False."""
        client = Mock()
        mock_task = Mock()
        mock_task.sdp_metadata = {"scope_files": ["src/allowed.py"]}
        client.get_task.return_value = mock_task

        guard = GuardSkill(client)
        guard._active_ws = "00-032-01"

        result = guard.check_edit("src/forbidden.py")

        assert result.allowed is False
        assert "not in WS scope" in result.reason
        assert "src/allowed.py" in result.scope_files
        assert result.ws_id == "00-032-01"

    def test_file_in_scope_allows_edit(self) -> None:
        """File in scope returns allowed=True."""
        client = Mock()
        mock_task = Mock()
        mock_task.sdp_metadata = {"scope_files": ["src/allowed.py"]}
        client.get_task.return_value = mock_task

        guard = GuardSkill(client)
        guard._active_ws = "00-032-01"

        result = guard.check_edit("src/allowed.py")

        assert result.allowed is True
        assert result.ws_id == "00-032-01"

    def test_no_scope_allows_all(self) -> None:
        """No scope_files defined allows all edits."""
        client = Mock()
        mock_task = Mock()
        mock_task.sdp_metadata = {}
        client.get_task.return_value = mock_task

        guard = GuardSkill(client)
        guard._active_ws = "00-032-01"

        result = guard.check_edit("any/file.py")

        assert result.allowed is True
        assert result.ws_id == "00-032-01"
        assert "No scope restrictions" in result.reason

    def test_activate_sets_active_ws(self) -> None:
        """AC2: activate() sets the active workstream."""
        client = Mock()
        mock_task = Mock()
        mock_task.id = "00-032-01"
        client.get_task.return_value = mock_task

        guard = GuardSkill(client)
        guard.activate("00-032-01")

        assert guard._active_ws == "00-032-01"
        client.get_task.assert_called_once_with("00-032-01")

    def test_activate_nonexistent_ws_raises_error(self) -> None:
        """activate() raises ValueError for non-existent WS."""
        client = Mock()
        client.get_task.return_value = None

        guard = GuardSkill(client)

        with pytest.raises(ValueError, match="WS not found"):
            guard.activate("invalid-ws")
