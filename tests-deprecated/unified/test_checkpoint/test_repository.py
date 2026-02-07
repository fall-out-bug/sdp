"""
Tests for CheckpointRepository.

Repository pattern for checkpoint management with error handling and logging.
"""

import pytest
from pathlib import Path
from datetime import datetime, timezone

from sdp.unified.checkpoint.schema import Checkpoint, CheckpointStatus
from sdp.unified.checkpoint.repository import CheckpointRepository, RepositoryError


class TestCheckpointRepositoryInit:
    """Test repository initialization."""

    def test_creates_repository_with_database_path(self, tmp_path):
        """Should initialize with database path."""
        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))

        assert repo.db_path == db_path

    def test_initializes_database_on_creation(self, tmp_path):
        """Should initialize database schema."""
        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))

        repo.initialize()

        assert db_path.exists()


class TestSaveCheckpoint:
    """Test save_checkpoint method."""

    def test_saves_checkpoint_to_database(self, tmp_path):
        """Should save checkpoint and return checkpoint ID."""
        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))
        repo.initialize()

        checkpoint = Checkpoint(
            feature="sdp-118",
            agent_id="agent-001",
            status=CheckpointStatus.IN_PROGRESS,
            completed_ws=[],
            execution_order=["ws-001", "ws-002"],
            started_at=datetime.now(timezone.utc),
        )

        saved_id = repo.save_checkpoint(checkpoint)

        assert saved_id > 0

    def test_logs_on_successful_save(self, tmp_path, caplog):
        """Should log info when checkpoint saved."""
        import logging

        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))
        repo.initialize()

        checkpoint = Checkpoint(
            feature="sdp-118",
            agent_id="agent-001",
            status=CheckpointStatus.IN_PROGRESS,
            completed_ws=[],
            execution_order=[],
            started_at=datetime.now(timezone.utc),
        )

        with caplog.at_level(logging.INFO):
            repo.save_checkpoint(checkpoint)

        assert "Checkpoint saved" in caplog.text
        assert "sdp-118" in caplog.text


class TestLoadCheckpoint:
    """Test load_checkpoint methods."""

    def test_load_checkpoint_by_feature(self, tmp_path):
        """Should load latest checkpoint for feature."""
        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))
        repo.initialize()

        checkpoint = Checkpoint(
            feature="sdp-118",
            agent_id="agent-001",
            status=CheckpointStatus.IN_PROGRESS,
            completed_ws=["ws-001"],
            execution_order=["ws-001", "ws-002"],
            started_at=datetime.now(timezone.utc),
        )
        repo.save_checkpoint(checkpoint)

        loaded = repo.load_checkpoint("sdp-118")

        assert loaded is not None
        assert loaded.feature == "sdp-118"
        assert loaded.completed_ws == ["ws-001"]


class TestUpdateCheckpointStatus:
    """Test update_checkpoint_status method."""

    def test_updates_checkpoint_status(self, tmp_path):
        """Should update checkpoint status."""
        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))
        repo.initialize()

        checkpoint = Checkpoint(
            feature="sdp-118",
            agent_id="agent-001",
            status=CheckpointStatus.IN_PROGRESS,
            completed_ws=[],
            execution_order=[],
            started_at=datetime.now(timezone.utc),
        )
        checkpoint_id = repo.save_checkpoint(checkpoint)

        repo.update_checkpoint_status(
            checkpoint_id,
            new_status=CheckpointStatus.COMPLETED,
            completed_ws=["ws-001"],
        )

        loaded = repo.load_checkpoint("sdp-118")
        assert loaded.status == CheckpointStatus.COMPLETED


class TestListActiveCheckpoints:
    """Test list_active_checkpoints method."""

    def test_returns_active_checkpoints(self, tmp_path):
        """Should return in_progress and failed checkpoints."""
        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))
        repo.initialize()

        checkpoint1 = Checkpoint(
            feature="sdp-118",
            agent_id="agent-001",
            status=CheckpointStatus.IN_PROGRESS,
            completed_ws=[],
            execution_order=[],
            started_at=datetime.now(timezone.utc),
        )
        repo.save_checkpoint(checkpoint1)

        active = repo.list_active_checkpoints()
        assert len(active) == 1
        assert active[0].status == CheckpointStatus.IN_PROGRESS


class TestErrorHandling:
    """Test error handling."""

    def test_raises_on_save_without_init(self, tmp_path):
        """Should raise RepositoryError when saving without init."""
        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))

        checkpoint = Checkpoint(
            feature="sdp-118",
            agent_id="agent-001",
            status=CheckpointStatus.IN_PROGRESS,
            completed_ws=[],
            execution_order=[],
            started_at=datetime.now(timezone.utc),
        )

        with pytest.raises(RepositoryError) as exc_info:
            repo.save_checkpoint(checkpoint)

        assert "not initialized" in str(exc_info.value)

    def test_raises_on_load_without_init(self, tmp_path):
        """Should raise RepositoryError when loading without init."""
        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))

        with pytest.raises(RepositoryError) as exc_info:
            repo.load_checkpoint("sdp-118")

        assert "not initialized" in str(exc_info.value)

    def test_raises_on_update_without_init(self, tmp_path):
        """Should raise RepositoryError when updating without init."""
        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))

        with pytest.raises(RepositoryError) as exc_info:
            repo.update_checkpoint_status(1, CheckpointStatus.COMPLETED, [])

        assert "not initialized" in str(exc_info.value)

    def test_raises_on_list_without_init(self, tmp_path):
        """Should raise RepositoryError when listing without init."""
        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))

        with pytest.raises(RepositoryError) as exc_info:
            repo.list_active_checkpoints()

        assert "not initialized" in str(exc_info.value)

    def test_raises_on_update_nonexistent_checkpoint(self, tmp_path):
        """Should raise RepositoryError when updating non-existent checkpoint."""
        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))
        repo.initialize()

        with pytest.raises(RepositoryError) as exc_info:
            repo.update_checkpoint_status(999, CheckpointStatus.COMPLETED, [])

        assert "not found" in str(exc_info.value)

    def test_load_latest_returns_none_for_completed_checkpoints(self, tmp_path):
        """Should return None if checkpoint is completed."""
        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))
        repo.initialize()

        checkpoint = Checkpoint(
            feature="sdp-118",
            agent_id="agent-001",
            status=CheckpointStatus.COMPLETED,
            completed_ws=[],
            execution_order=[],
            started_at=datetime.now(timezone.utc),
        )
        repo.save_checkpoint(checkpoint)

        latest = repo.load_latest_checkpoint("sdp-118")
        assert latest is None

    def test_close_closes_database(self, tmp_path):
        """Should close database connection."""
        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))
        repo.initialize()

        repo.close()

        # Should not raise error
        repo.close()

    def test_context_manager_initializes_and_closes(self, tmp_path):
        """Should work as context manager."""
        db_path = tmp_path / "test_checkpoints.db"

        with CheckpointRepository(str(db_path)) as repo:
            repo.initialize()
            assert repo._db is not None

        # Database closed after context
        assert repo._db is None

    def test_initialize_raises_on_database_error(self, tmp_path, monkeypatch):
        """Should raise RepositoryError when database initialization fails."""
        import logging
        from unittest.mock import MagicMock, patch

        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))

        # Mock CheckpointDatabase to raise exception
        with patch("sdp.unified.checkpoint.repository.CheckpointDatabase") as mock_db:
            mock_instance = MagicMock()
            mock_instance.initialize.side_effect = Exception("Database error")
            mock_db.return_value = mock_instance

            with pytest.raises(RepositoryError) as exc_info:
                repo.initialize()

            assert "Failed to initialize repository" in str(exc_info.value)

    def test_save_checkpoint_raises_on_database_error(self, tmp_path, monkeypatch):
        """Should raise RepositoryError when save fails."""
        import logging
        from unittest.mock import MagicMock, patch

        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))
        repo.initialize()

        checkpoint = Checkpoint(
            feature="sdp-118",
            agent_id="agent-001",
            status=CheckpointStatus.IN_PROGRESS,
            completed_ws=[],
            execution_order=[],
            started_at=datetime.now(timezone.utc),
        )

        # Mock database to raise exception
        repo._db.create_checkpoint = MagicMock(side_effect=Exception("Save failed"))

        with pytest.raises(RepositoryError) as exc_info:
            repo.save_checkpoint(checkpoint)

        assert "Failed to save checkpoint" in str(exc_info.value)

    def test_load_checkpoint_raises_on_database_error(self, tmp_path, monkeypatch):
        """Should raise RepositoryError when load fails."""
        from unittest.mock import MagicMock

        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))
        repo.initialize()

        # Mock database to raise exception
        repo._db.get_checkpoint_by_feature = MagicMock(side_effect=Exception("Load failed"))

        with pytest.raises(RepositoryError) as exc_info:
            repo.load_checkpoint("sdp-118")

        assert "Failed to load checkpoint" in str(exc_info.value)

    def test_load_latest_checkpoint_raises_on_database_error(self, tmp_path):
        """Should raise RepositoryError when load_latest_checkpoint fails."""
        from unittest.mock import MagicMock

        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))
        repo.initialize()

        # Mock database to raise exception
        repo._db.get_checkpoint_by_feature = MagicMock(side_effect=Exception("Load failed"))

        with pytest.raises(RepositoryError) as exc_info:
            repo.load_latest_checkpoint("sdp-118")

        assert "Failed to load latest checkpoint" in str(exc_info.value)

    def test_load_latest_checkpoint_returns_none_for_none_checkpoint(self, tmp_path):
        """Should return None when checkpoint doesn't exist."""
        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))
        repo.initialize()

        latest = repo.load_latest_checkpoint("nonexistent")
        assert latest is None

    def test_load_latest_checkpoint_logs_when_found(self, tmp_path, caplog):
        """Should log info when active checkpoint found."""
        import logging

        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))
        repo.initialize()

        checkpoint = Checkpoint(
            feature="sdp-118",
            agent_id="agent-001",
            status=CheckpointStatus.IN_PROGRESS,
            completed_ws=[],
            execution_order=[],
            started_at=datetime.now(timezone.utc),
        )
        repo.save_checkpoint(checkpoint)

        with caplog.at_level(logging.INFO):
            latest = repo.load_latest_checkpoint("sdp-118")

        assert latest is not None
        assert "Latest active checkpoint found" in caplog.text

    def test_update_checkpoint_status_raises_on_database_error(self, tmp_path):
        """Should raise RepositoryError when update fails."""
        from unittest.mock import MagicMock

        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))
        repo.initialize()

        checkpoint = Checkpoint(
            feature="sdp-118",
            agent_id="agent-001",
            status=CheckpointStatus.IN_PROGRESS,
            completed_ws=[],
            execution_order=[],
            started_at=datetime.now(timezone.utc),
        )
        checkpoint_id = repo.save_checkpoint(checkpoint)

        # Mock database to raise exception during update
        repo._db.update_checkpoint = MagicMock(side_effect=Exception("Update failed"))

        with pytest.raises(RepositoryError) as exc_info:
            repo.update_checkpoint_status(
                checkpoint_id, CheckpointStatus.COMPLETED, ["ws-001"]
            )

        assert "Failed to update checkpoint" in str(exc_info.value)

    def test_list_active_checkpoints_raises_on_database_error(self, tmp_path):
        """Should raise RepositoryError when list fails."""
        from unittest.mock import MagicMock

        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))
        repo.initialize()

        # Mock database to raise exception
        repo._db.get_active_checkpoints = MagicMock(side_effect=Exception("List failed"))

        with pytest.raises(RepositoryError) as exc_info:
            repo.list_active_checkpoints()

        assert "Failed to list active checkpoints" in str(exc_info.value)

    def test_load_latest_checkpoint_raises_when_not_initialized(self, tmp_path):
        """Should raise RepositoryError when repository not initialized."""
        db_path = tmp_path / "test_checkpoints.db"
        repo = CheckpointRepository(str(db_path))

        with pytest.raises(RepositoryError) as exc_info:
            repo.load_latest_checkpoint("F01")

        assert "not initialized" in str(exc_info.value)
