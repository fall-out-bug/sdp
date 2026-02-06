"""
Tests for checkpoint database schema.

Test-first approach: Write failing tests, then implement schema.
"""

import pytest
import sqlite3
from pathlib import Path
from datetime import datetime, timezone

from sdp.unified.checkpoint.schema import (
    CheckpointDatabase,
    Checkpoint,
    CheckpointStatus,
)


class TestCheckpointDatabaseInit:
    """Test database initialization and schema creation."""

    def test_creates_database_file_if_not_exists(self, tmp_path):
        """Database should be created if it doesn't exist."""
        db_path = tmp_path / "test_checkpoints.db"

        assert not db_path.exists()

        db = CheckpointDatabase(str(db_path))
        db.initialize()

        assert db_path.exists()

    def test_initializes_with_correct_schema_version(self, tmp_path):
        """Database should track schema version."""
        db_path = tmp_path / "test_checkpoints.db"
        db = CheckpointDatabase(str(db_path))
        db.initialize()

        version = db.get_schema_version()

        assert version == 1  # First schema version

    def test_creates_checkpoints_table(self, tmp_path):
        """Checkpoints table should exist with correct columns."""
        db_path = tmp_path / "test_checkpoints.db"
        db = CheckpointDatabase(str(db_path))
        db.initialize()

        conn = sqlite3.connect(str(db_path))
        cursor = conn.cursor()

        # Get table schema
        cursor.execute(
            "SELECT sql FROM sqlite_master WHERE type='table' AND name='checkpoints'"
        )
        table_sql = cursor.fetchone()[0]

        # Verify key columns exist
        assert "feature" in table_sql
        assert "agent_id" in table_sql
        assert "status" in table_sql
        assert "completed_ws" in table_sql
        assert "execution_order" in table_sql
        assert "started_at" in table_sql
        assert "metrics" in table_sql

        conn.close()

    def test_creates_indexes_for_performance(self, tmp_path):
        """Indexes should exist on frequently queried columns."""
        db_path = tmp_path / "test_checkpoints.db"
        db = CheckpointDatabase(str(db_path))
        db.initialize()

        conn = sqlite3.connect(str(db_path))
        cursor = conn.cursor()

        # Get indexes
        cursor.execute(
            "SELECT name FROM sqlite_master WHERE type='index' AND tbl_name='checkpoints'"
        )
        indexes = [row[0] for row in cursor.fetchall()]

        # Should have index on feature (for lookups)
        assert any("feature" in idx for idx in indexes)

        # Should have index on status (for filtering)
        assert any("status" in idx for idx in indexes)

        conn.close()


class TestCheckpointCRUD:
    """Test create, read, update operations."""

    def test_create_checkpoint_saves_all_fields(self, tmp_path):
        """Should save checkpoint with all required fields."""
        db_path = tmp_path / "test_checkpoints.db"
        db = CheckpointDatabase(str(db_path))
        db.initialize()

        checkpoint = Checkpoint(
            feature="sdp-118",
            agent_id="agent-20260128-120000",
            status=CheckpointStatus.IN_PROGRESS,
            completed_ws=[],
            execution_order=["ws-001", "ws-002"],
            started_at=datetime.now(timezone.utc),
        )

        saved_id = db.create_checkpoint(checkpoint)

        # Retrieve and verify
        retrieved = db.get_checkpoint(saved_id)

        assert retrieved.feature == "sdp-118"
        assert retrieved.agent_id == "agent-20260128-120000"
        assert retrieved.status == CheckpointStatus.IN_PROGRESS
        assert retrieved.completed_ws == []
        assert retrieved.execution_order == ["ws-001", "ws-002"]

    def test_get_checkpoint_by_feature_returns_latest(self, tmp_path):
        """Should return the most recent checkpoint for a feature."""
        import time

        db_path = tmp_path / "test_checkpoints.db"
        db = CheckpointDatabase(str(db_path))
        db.initialize()

        # Create first checkpoint
        checkpoint1 = Checkpoint(
            feature="sdp-118",
            agent_id="agent-1",
            status=CheckpointStatus.IN_PROGRESS,
            completed_ws=[],
            execution_order=["ws-001"],
            started_at=datetime.now(timezone.utc),
        )
        db.create_checkpoint(checkpoint1)

        # Small delay to ensure different timestamps
        time.sleep(0.01)

        # Create second checkpoint (later)
        checkpoint2 = Checkpoint(
            feature="sdp-118",
            agent_id="agent-2",
            status=CheckpointStatus.IN_PROGRESS,
            completed_ws=["ws-001"],
            execution_order=["ws-001", "ws-002"],
            started_at=datetime.now(timezone.utc),
        )
        db.create_checkpoint(checkpoint2)

        # Get latest by feature
        latest = db.get_checkpoint_by_feature("sdp-118")

        assert latest.agent_id == "agent-2"
        assert len(latest.completed_ws) == 1

    def test_update_checkpoint_modifies_status_and_completed_ws(self, tmp_path):
        """Should update checkpoint fields."""
        db_path = tmp_path / "test_checkpoints.db"
        db = CheckpointDatabase(str(db_path))
        db.initialize()

        # Create checkpoint
        checkpoint = Checkpoint(
            feature="sdp-118",
            agent_id="agent-1",
            status=CheckpointStatus.IN_PROGRESS,
            completed_ws=[],
            execution_order=["ws-001", "ws-002"],
            started_at=datetime.now(timezone.utc),
        )
        saved_id = db.create_checkpoint(checkpoint)

        # Update
        checkpoint.status = CheckpointStatus.COMPLETED
        checkpoint.completed_ws = ["ws-001", "ws-002"]
        db.update_checkpoint(saved_id, checkpoint)

        # Verify update
        retrieved = db.get_checkpoint(saved_id)

        assert retrieved.status == CheckpointStatus.COMPLETED
        assert retrieved.completed_ws == ["ws-001", "ws-002"]

    def test_delete_checkpoint_removes_from_database(self, tmp_path):
        """Should delete checkpoint by ID."""
        db_path = tmp_path / "test_checkpoints.db"
        db = CheckpointDatabase(str(db_path))
        db.initialize()

        checkpoint = Checkpoint(
            feature="sdp-118",
            agent_id="agent-1",
            status=CheckpointStatus.IN_PROGRESS,
            completed_ws=[],
            execution_order=["ws-001"],
            started_at=datetime.now(timezone.utc),
        )
        saved_id = db.create_checkpoint(checkpoint)

        # Delete
        db.delete_checkpoint(saved_id)

        # Verify deleted
        retrieved = db.get_checkpoint(saved_id)

        assert retrieved is None


class TestCheckpointQueries:
    """Test common query patterns."""

    def test_list_checkpoints_by_status(self, tmp_path):
        """Should filter checkpoints by status."""
        db_path = tmp_path / "test_checkpoints.db"
        db = CheckpointDatabase(str(db_path))
        db.initialize()

        # Create checkpoints with different statuses
        for i, status in enumerate(
            [CheckpointStatus.IN_PROGRESS, CheckpointStatus.COMPLETED, CheckpointStatus.FAILED]
        ):
            checkpoint = Checkpoint(
                feature=f"sdp-{i}",
                agent_id=f"agent-{i}",
                status=status,
                completed_ws=[],
                execution_order=[],
                started_at=datetime.now(timezone.utc),
            )
            db.create_checkpoint(checkpoint)

        # Query in_progress
        in_progress = db.list_checkpoints_by_status(CheckpointStatus.IN_PROGRESS)

        assert len(in_progress) == 1
        assert in_progress[0].status == CheckpointStatus.IN_PROGRESS

    def test_get_active_checkpoints_returns_in_progress_or_failed(self, tmp_path):
        """Should return checkpoints that need attention."""
        db_path = tmp_path / "test_checkpoints.db"
        db = CheckpointDatabase(str(db_path))
        db.initialize()

        # Create checkpoints
        for i, status in enumerate(
            [
                CheckpointStatus.IN_PROGRESS,
                CheckpointStatus.COMPLETED,
                CheckpointStatus.FAILED,
            ]
        ):
            checkpoint = Checkpoint(
                feature=f"sdp-{i}",
                agent_id=f"agent-{i}",
                status=status,
                completed_ws=[],
                execution_order=[],
                started_at=datetime.now(timezone.utc),
            )
            db.create_checkpoint(checkpoint)

        # Get active (in_progress + failed)
        active = db.get_active_checkpoints()

        assert len(active) == 2  # in_progress + failed
