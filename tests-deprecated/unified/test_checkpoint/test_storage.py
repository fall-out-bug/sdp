"""
Tests for CheckpointDatabase storage operations.

Tests error handling, edge cases, and context manager behavior.
"""

import pytest
import sqlite3
from datetime import datetime, timezone
from unittest.mock import MagicMock, patch

from sdp.unified.checkpoint.storage import CheckpointDatabase
from sdp.unified.checkpoint.models import Checkpoint, CheckpointStatus


@pytest.fixture
def temp_db_path(tmp_path):
    """Create temporary database path."""
    return str(tmp_path / "test_storage.db")


@pytest.fixture
def checkpoint_db(temp_db_path):
    """Create initialized checkpoint database."""
    db = CheckpointDatabase(temp_db_path)
    db.initialize()
    yield db
    db.close()


@pytest.fixture
def sample_checkpoint():
    """Create sample checkpoint for testing."""
    return Checkpoint(
        feature="F01",
        agent_id="agent-001",
        status=CheckpointStatus.IN_PROGRESS,
        completed_ws=[],
        execution_order=["WS-001", "WS-002"],
        started_at=datetime.now(timezone.utc),
    )


class TestCheckpointDatabaseInit:
    """Test database initialization."""

    def test_creates_database_file(self, temp_db_path):
        """Should create database file on initialization."""
        db = CheckpointDatabase(temp_db_path)
        db.initialize()
        db.close()

        from pathlib import Path
        assert Path(temp_db_path).exists()

    def test_get_schema_version_creates_schema_manager_if_none(self, checkpoint_db):
        """Should create schema manager if None when getting version."""
        # Reset schema manager
        checkpoint_db._schema_manager = None

        version = checkpoint_db.get_schema_version()

        assert isinstance(version, int)
        assert checkpoint_db._schema_manager is not None


class TestCreateCheckpoint:
    """Test checkpoint creation."""

    def test_create_checkpoint_returns_id(self, checkpoint_db, sample_checkpoint):
        """Should return checkpoint ID after creation."""
        checkpoint_id = checkpoint_db.create_checkpoint(sample_checkpoint)

        assert checkpoint_id > 0

    def test_create_checkpoint_raises_on_no_lastrowid(self, checkpoint_db, sample_checkpoint):
        """Should raise RuntimeError if lastrowid is None."""
        # Mock cursor to return None for lastrowid
        mock_cursor = MagicMock()
        mock_cursor.lastrowid = None
        mock_cursor.execute = MagicMock()
        
        mock_conn = MagicMock()
        mock_conn.cursor.return_value = mock_cursor
        
        # Replace _get_connection to return mock
        checkpoint_db._get_connection = MagicMock(return_value=mock_conn)

        with pytest.raises(RuntimeError) as exc_info:
            checkpoint_db.create_checkpoint(sample_checkpoint)

        assert "Failed to get checkpoint ID" in str(exc_info.value)


class TestContextManager:
    """Test context manager behavior."""

    def test_context_manager_enters_and_exits(self, temp_db_path):
        """Should work as context manager."""
        with CheckpointDatabase(temp_db_path) as db:
            db.initialize()
            assert db._conn is not None

        # Connection closed after context
        assert db._conn is None

    def test_context_manager_exit_closes_connection(self, temp_db_path):
        """Should close connection on exit."""
        db = CheckpointDatabase(temp_db_path)
        db.initialize()
        assert db._conn is not None

        # Call __exit__ directly
        db.__exit__(None, None, None)

        assert db._conn is None

    def test_context_manager_exit_handles_exceptions(self, temp_db_path):
        """Should close connection even if exception occurs."""
        db = CheckpointDatabase(temp_db_path)
        db.initialize()

        # Simulate exception during context
        try:
            db.__exit__(ValueError, ValueError("test"), None)
        except ValueError:
            pass

        # Connection should still be closed
        assert db._conn is None
