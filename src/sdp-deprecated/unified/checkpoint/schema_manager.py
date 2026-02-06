"""
Checkpoint database schema management.

Handles schema creation, versioning, and migrations.
"""

import sqlite3


class SchemaManager:
    """Manages database schema and migrations."""

    SCHEMA_VERSION = 1

    def __init__(self, conn: sqlite3.Connection):
        """Initialize schema manager.

        Args:
            conn: SQLite connection
        """
        self.conn = conn

    def ensure_schema(self) -> None:
        """Create tables and indexes if they don't exist."""
        cursor = self.conn.cursor()

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
        self._create_indexes()

        self.conn.commit()

    def _create_indexes(self) -> None:
        """Create database indexes for query performance."""
        cursor = self.conn.cursor()

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

    def get_schema_version(self) -> int:
        """Get current schema version.

        Returns:
            Schema version number
        """
        cursor = self.conn.cursor()

        cursor.execute("SELECT MAX(version) as version FROM schema_version")
        row = cursor.fetchone()

        return row["version"] if row else 0
