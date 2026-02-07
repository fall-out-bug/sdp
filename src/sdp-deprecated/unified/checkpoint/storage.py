"""
Checkpoint database storage operations.

Implements SQLite-based checkpoint CRUD operations.
"""

import sqlite3
from pathlib import Path
from typing import Optional

from .models import Checkpoint, CheckpointStatus
from .schema_manager import SchemaManager
from .serialization import (
    checkpoint_to_insert_params,
    checkpoint_to_update_params,
    row_to_checkpoint,
)


class CheckpointDatabase:
    """SQLite database for checkpoint storage."""

    def __init__(self, db_path: str):
        """Initialize database connection.

        Args:
            db_path: Path to SQLite database file
        """
        self.db_path = Path(db_path)
        self._conn: Optional[sqlite3.Connection] = None
        self._schema_manager: Optional[SchemaManager] = None

    def initialize(self) -> None:
        """Create database schema if not exists."""
        conn = self._get_connection()
        self._schema_manager = SchemaManager(conn)
        self._schema_manager.ensure_schema()

    def _get_connection(self) -> sqlite3.Connection:
        """Get or create database connection."""
        if self._conn is None:
            self._conn = sqlite3.connect(str(self.db_path))
            self._conn.row_factory = sqlite3.Row
        return self._conn

    def get_schema_version(self) -> int:
        """Get current schema version."""
        if self._schema_manager is None:
            conn = self._get_connection()
            self._schema_manager = SchemaManager(conn)
        return self._schema_manager.get_schema_version()

    def create_checkpoint(self, checkpoint: Checkpoint) -> int:
        """Create a new checkpoint.

        Args:
            checkpoint: Checkpoint to save

        Returns:
            ID of created checkpoint
        """
        conn = self._get_connection()
        cursor = conn.cursor()

        cursor.execute(
            """
            INSERT INTO checkpoints (
                feature, agent_id, status, completed_ws, execution_order,
                started_at, current_ws, completed_at, failed_tasks, error, metrics
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        """,
            checkpoint_to_insert_params(checkpoint),
        )

        conn.commit()

        lastrowid = cursor.lastrowid
        if lastrowid is None:
            raise RuntimeError("Failed to get checkpoint ID after insert")
        return lastrowid

    def get_checkpoint(self, checkpoint_id: int) -> Optional[Checkpoint]:
        """Get checkpoint by ID.

        Args:
            checkpoint_id: Checkpoint ID

        Returns:
            Checkpoint or None if not found
        """

        conn = self._get_connection()
        cursor = conn.cursor()

        cursor.execute("SELECT * FROM checkpoints WHERE id = ?", (checkpoint_id,))
        row = cursor.fetchone()

        if not row:
            return None

        return row_to_checkpoint(row)

    def get_checkpoint_by_feature(self, feature: str) -> Optional[Checkpoint]:
        """Get latest checkpoint for a feature.

        Args:
            feature: Feature ID

        Returns:
            Latest checkpoint or None
        """
        conn = self._get_connection()
        cursor = conn.cursor()

        cursor.execute(
            """
            SELECT * FROM checkpoints
            WHERE feature = ?
            ORDER BY id DESC
            LIMIT 1
        """,
            (feature,),
        )

        row = cursor.fetchone()

        if not row:
            return None

        return row_to_checkpoint(row)

    def update_checkpoint(self, checkpoint_id: int, checkpoint: Checkpoint) -> None:
        """Update existing checkpoint.

        Args:
            checkpoint_id: Checkpoint ID
            checkpoint: Updated checkpoint data
        """
        conn = self._get_connection()
        cursor = conn.cursor()

        cursor.execute(
            """
            UPDATE checkpoints SET
                status = ?,
                completed_ws = ?,
                execution_order = ?,
                current_ws = ?,
                completed_at = ?,
                failed_tasks = ?,
                error = ?,
                metrics = ?
            WHERE id = ?
        """,
            checkpoint_to_update_params(checkpoint, checkpoint_id),
        )

        conn.commit()

    def delete_checkpoint(self, checkpoint_id: int) -> None:
        """Delete checkpoint by ID.

        Args:
            checkpoint_id: Checkpoint ID
        """
        conn = self._get_connection()
        cursor = conn.cursor()

        cursor.execute("DELETE FROM checkpoints WHERE id = ?", (checkpoint_id,))

        conn.commit()

    def list_checkpoints_by_status(self, status: CheckpointStatus) -> list[Checkpoint]:
        """List all checkpoints with given status.

        Args:
            status: Status to filter by

        Returns:
            List of checkpoints
        """
        conn = self._get_connection()
        cursor = conn.cursor()

        cursor.execute(
            "SELECT * FROM checkpoints WHERE status = ? ORDER BY created_at DESC",
            (status.value,),
        )

        return [row_to_checkpoint(row) for row in cursor.fetchall()]

    def get_active_checkpoints(self) -> list[Checkpoint]:
        """Get checkpoints that need attention (in_progress or failed).

        Returns:
            List of active checkpoints
        """
        conn = self._get_connection()
        cursor = conn.cursor()

        cursor.execute(
            """
            SELECT * FROM checkpoints
            WHERE status IN ('in_progress', 'failed')
            ORDER BY created_at DESC
        """
        )

        return [row_to_checkpoint(row) for row in cursor.fetchall()]

    def close(self) -> None:
        """Close database connection."""
        if self._conn:
            self._conn.close()
            self._conn = None

    def __enter__(self) -> "CheckpointDatabase":
        """Context manager entry."""
        return self

    def __exit__(
        self,
        exc_type: type[BaseException] | None,
        exc_val: BaseException | None,
        exc_tb: object,
    ) -> None:
        """Context manager exit."""
        self.close()
