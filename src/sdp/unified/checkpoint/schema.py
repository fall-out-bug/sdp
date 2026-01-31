"""
Checkpoint database schema and management.

Implements SQLite-based checkpoint storage for @oneshot execution.
"""

import sqlite3
from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum
from pathlib import Path
from typing import Optional


class CheckpointStatus(Enum):
    """Status of a checkpoint execution."""

    IN_PROGRESS = "in_progress"
    COMPLETED = "completed"
    FAILED = "failed"
    CANCELLED = "cancelled"


@dataclass
class Checkpoint:
    """Checkpoint data model."""

    feature: str
    agent_id: str
    status: CheckpointStatus
    completed_ws: list[str]
    execution_order: list[str]
    started_at: datetime
    current_ws: Optional[str] = None
    completed_at: Optional[datetime] = None
    failed_tasks: list[str] = field(default_factory=list)
    error: Optional[str] = None
    metrics: dict[str, object] = field(default_factory=dict)


class CheckpointDatabase:
    """SQLite database for checkpoint storage."""

    SCHEMA_VERSION = 1

    def __init__(self, db_path: str):
        """Initialize database connection.

        Args:
            db_path: Path to SQLite database file
        """
        self.db_path = Path(db_path)
        self._conn: Optional[sqlite3.Connection] = None

    def initialize(self) -> None:
        """Create database schema if not exists."""
        self._ensure_schema()

    def _get_connection(self) -> sqlite3.Connection:
        """Get or create database connection."""
        if self._conn is None:
            self._conn = sqlite3.connect(str(self.db_path))
            self._conn.row_factory = sqlite3.Row
        return self._conn

    def _ensure_schema(self) -> None:
        """Create tables and indexes if they don't exist."""
        conn = self._get_connection()
        cursor = conn.cursor()

        # Create checkpoints table
        cursor.execute(
            """
            CREATE TABLE IF NOT EXISTS checkpoints (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                feature TEXT NOT NULL,
                agent_id TEXT NOT NULL,
                status TEXT NOT NULL,
                completed_ws TEXT NOT NULL DEFAULT '[]',
                execution_order TEXT NOT NULL DEFAULT '[]',
                started_at TEXT NOT NULL,
                current_ws TEXT,
                completed_at TEXT,
                failed_tasks TEXT NOT NULL DEFAULT '[]',
                error TEXT,
                metrics TEXT NOT NULL DEFAULT '{}',
                created_at TEXT DEFAULT CURRENT_TIMESTAMP
            )
        """
        )

        # Create schema_version table
        cursor.execute(
            """
            CREATE TABLE IF NOT EXISTS schema_version (
                version INTEGER PRIMARY KEY,
                applied_at TEXT DEFAULT CURRENT_TIMESTAMP
            )
        """
        )

        # Insert schema version if not exists
        cursor.execute(
            """
            INSERT OR IGNORE INTO schema_version (version)
            VALUES (?)
        """,
            (self.SCHEMA_VERSION,),
        )

        # Create indexes for performance
        cursor.execute(
            """
            CREATE INDEX IF NOT EXISTS idx_checkpoints_feature
            ON checkpoints(feature)
        """
        )

        cursor.execute(
            """
            CREATE INDEX IF NOT EXISTS idx_checkpoints_status
            ON checkpoints(status)
        """
        )

        cursor.execute(
            """
            CREATE INDEX IF NOT EXISTS idx_checkpoints_feature_status
            ON checkpoints(feature, status)
        """
        )

        conn.commit()

    def get_schema_version(self) -> int:
        """Get current schema version."""
        conn = self._get_connection()
        cursor = conn.cursor()

        cursor.execute("SELECT MAX(version) as version FROM schema_version")
        row = cursor.fetchone()

        return row["version"] if row else 0

    def create_checkpoint(self, checkpoint: Checkpoint) -> int:
        """Create a new checkpoint.

        Args:
            checkpoint: Checkpoint to save

        Returns:
            ID of created checkpoint
        """
        import json

        conn = self._get_connection()
        cursor = conn.cursor()

        cursor.execute(
            """
            INSERT INTO checkpoints (
                feature, agent_id, status, completed_ws, execution_order,
                started_at, current_ws, completed_at, failed_tasks, error, metrics
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        """,
            (
                checkpoint.feature,
                checkpoint.agent_id,
                checkpoint.status.value,
                json.dumps(checkpoint.completed_ws),
                json.dumps(checkpoint.execution_order),
                checkpoint.started_at.isoformat(),
                checkpoint.current_ws,
                checkpoint.completed_at.isoformat() if checkpoint.completed_at else None,
                json.dumps(checkpoint.failed_tasks),
                checkpoint.error,
                json.dumps(checkpoint.metrics),
            ),
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

        return self._row_to_checkpoint(row)

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

        return self._row_to_checkpoint(row)

    def update_checkpoint(self, checkpoint_id: int, checkpoint: Checkpoint) -> None:
        """Update existing checkpoint.

        Args:
            checkpoint_id: Checkpoint ID
            checkpoint: Updated checkpoint data
        """
        import json

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
            (
                checkpoint.status.value,
                json.dumps(checkpoint.completed_ws),
                json.dumps(checkpoint.execution_order),
                checkpoint.current_ws,
                checkpoint.completed_at.isoformat() if checkpoint.completed_at else None,
                json.dumps(checkpoint.failed_tasks),
                checkpoint.error,
                json.dumps(checkpoint.metrics),
                checkpoint_id,
            ),
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

        return [self._row_to_checkpoint(row) for row in cursor.fetchall()]

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

        return [self._row_to_checkpoint(row) for row in cursor.fetchall()]

    def _row_to_checkpoint(self, row: sqlite3.Row) -> Checkpoint:
        """Convert database row to Checkpoint object.

        Args:
            row: Database row

        Returns:
            Checkpoint object
        """
        import json

        return Checkpoint(
            feature=row["feature"],
            agent_id=row["agent_id"],
            status=CheckpointStatus(row["status"]),
            completed_ws=json.loads(row["completed_ws"]),
            execution_order=json.loads(row["execution_order"]),
            started_at=datetime.fromisoformat(row["started_at"]),
            current_ws=row["current_ws"],
            completed_at=(
                datetime.fromisoformat(row["completed_at"]) if row["completed_at"] else None
            ),
            failed_tasks=json.loads(row["failed_tasks"]),
            error=row["error"],
            metrics=json.loads(row["metrics"]),
        )

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
