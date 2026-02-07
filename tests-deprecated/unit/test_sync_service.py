"""Tests for Beads sync service."""

from pathlib import Path
from unittest.mock import Mock, patch

import pytest

from sdp.beads.mock import MockBeadsClient
from sdp.beads.models import BeadsStatus, BeadsTaskCreate
from sdp.beads.sync_service import BeadsSyncService, SyncSource


class TestBeadsSyncService:
    """Test suite for BeadsSyncService."""

    def test_check_sync_no_active_ws(self) -> None:
        """AC1: No active WS means in sync."""
        # Arrange
        client = MockBeadsClient()
        service = BeadsSyncService(client)

        # Act
        result = service.check_sync(active_ws=None)

        # Assert
        assert result.synced is True
        assert len(result.conflicts) == 0

    def test_check_sync_task_not_found(self) -> None:
        """AC3: Detects when local has active WS but Beads doesn't."""
        # Arrange
        client = MockBeadsClient()
        service = BeadsSyncService(client)

        # Act
        result = service.check_sync(active_ws="nonexistent")

        # Assert
        assert result.synced is False
        assert len(result.conflicts) == 1
        assert result.conflicts[0].ws_id == "nonexistent"
        assert result.conflicts[0].beads_status == "not_found"

    def test_check_sync_status_mismatch(self) -> None:
        """AC3: Detects status mismatch between local and Beads."""
        # Arrange
        client = MockBeadsClient()
        task = client.create_task(BeadsTaskCreate(title="Test"))
        service = BeadsSyncService(client)

        # Task is OPEN but local thinks it's active (IN_PROGRESS)
        # Act
        result = service.check_sync(active_ws=task.id)

        # Assert
        assert result.synced is False
        assert len(result.conflicts) == 1
        assert result.conflicts[0].field == "status"
        assert result.conflicts[0].beads_status == "open"

    def test_check_sync_matching_status(self) -> None:
        """AC1: No conflicts when statuses match."""
        # Arrange
        client = MockBeadsClient()
        task = client.create_task(BeadsTaskCreate(title="Test"))
        client.update_task_status(task.id, BeadsStatus.IN_PROGRESS)
        service = BeadsSyncService(client)

        # Act
        result = service.check_sync(active_ws=task.id)

        # Assert
        assert result.synced is True
        assert len(result.conflicts) == 0

    @patch.object(BeadsSyncService, "_clear_local_state")
    def test_sync_clears_local_when_task_not_found(
        self, mock_clear: Mock
    ) -> None:
        """AC2: Sync clears local state when Beads task not found."""
        # Arrange
        client = MockBeadsClient()
        service = BeadsSyncService(client)

        # Act
        result = service.sync(active_ws="nonexistent", source=SyncSource.BEADS)

        # Assert
        mock_clear.assert_called_once()
        assert len(result.changes) == 1
        assert "Cleared local" in result.changes[0]

    @patch.object(BeadsSyncService, "_clear_local_state")
    def test_sync_clears_local_when_beads_not_in_progress(
        self, mock_clear: Mock
    ) -> None:
        """AC2: Sync clears local when Beads status != IN_PROGRESS."""
        # Arrange
        client = MockBeadsClient()
        task = client.create_task(BeadsTaskCreate(title="Test"))
        # Task is OPEN, not IN_PROGRESS
        service = BeadsSyncService(client)

        # Act
        result = service.sync(active_ws=task.id, source=SyncSource.BEADS)

        # Assert
        mock_clear.assert_called_once()
        assert len(result.changes) == 1
        assert "open" in result.changes[0].lower()

    def test_sync_updates_beads_from_local(self) -> None:
        """AC2: Sync updates Beads status when local is source."""
        # Arrange
        client = MockBeadsClient()
        task = client.create_task(BeadsTaskCreate(title="Test"))
        # Task is OPEN
        service = BeadsSyncService(client)

        # Act
        result = service.sync(active_ws=task.id, source=SyncSource.LOCAL)

        # Assert
        assert len(result.changes) == 1
        assert "IN_PROGRESS" in result.changes[0]

        # Verify Beads was updated
        updated_task = client.get_task(task.id)
        assert updated_task is not None
        assert updated_task.status == BeadsStatus.IN_PROGRESS

    def test_sync_no_changes_when_already_synced(self) -> None:
        """Sync returns no changes when already in sync."""
        # Arrange
        client = MockBeadsClient()
        task = client.create_task(BeadsTaskCreate(title="Test"))
        client.update_task_status(task.id, BeadsStatus.IN_PROGRESS)
        service = BeadsSyncService(client)

        # Act
        result = service.sync(active_ws=task.id, source=SyncSource.BEADS)

        # Assert
        assert result.synced is True
        assert len(result.changes) == 0

    def test_clear_local_state_removes_file(self, tmp_path: Path) -> None:
        """_clear_local_state removes .guard_state file."""
        # Arrange
        state_file = tmp_path / ".guard_state"
        state_file.write_text("active_ws=bd-0001")

        client = MockBeadsClient()
        service = BeadsSyncService(client)

        # Patch Path.cwd() to return tmp_path
        with patch("pathlib.Path.cwd", return_value=tmp_path):
            # Act
            service._clear_local_state()

        # Assert
        assert not state_file.exists()

    def test_get_local_active_ws_from_file(self, tmp_path: Path) -> None:
        """_get_local_active_ws reads from .guard_state."""
        # Arrange
        state_file = tmp_path / ".guard_state"
        state_file.write_text("active_ws=bd-0001")

        client = MockBeadsClient()
        service = BeadsSyncService(client)

        # Patch Path.cwd() to return tmp_path
        with patch("pathlib.Path.cwd", return_value=tmp_path):
            # Act
            ws_id = service._get_local_active_ws()

        # Assert
        assert ws_id == "bd-0001"

    def test_get_local_active_ws_no_file(self, tmp_path: Path) -> None:
        """_get_local_active_ws returns None when file doesn't exist."""
        # Arrange
        client = MockBeadsClient()
        service = BeadsSyncService(client)

        # Patch Path.cwd() to return tmp_path
        with patch("pathlib.Path.cwd", return_value=tmp_path):
            # Act
            ws_id = service._get_local_active_ws()

        # Assert
        assert ws_id is None
