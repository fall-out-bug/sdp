"""
Checkpoint repository with error handling and logging.

Provides high-level interface for checkpoint management operations.
"""

import logging
from pathlib import Path
from typing import Optional

from sdp.unified.checkpoint.schema import Checkpoint, CheckpointDatabase, CheckpointStatus

logger = logging.getLogger(__name__)


class RepositoryError(Exception):
    """Repository operation failed."""

    pass


class CheckpointRepository:
    """Repository for checkpoint management with error handling."""

    def __init__(self, db_path: str) -> None:
        """Initialize repository.

        Args:
            db_path: Path to SQLite database file
        """
        self.db_path = Path(db_path)
        self._db: Optional[CheckpointDatabase] = None

    def initialize(self) -> None:
        """Initialize database schema."""
        try:
            self._db = CheckpointDatabase(str(self.db_path))
            self._db.initialize()
            logger.info(f"CheckpointRepository initialized with database: {self.db_path}")
        except Exception as e:
            logger.error(f"Failed to initialize repository: {e}")
            raise RepositoryError(f"Failed to initialize repository: {e}") from e

    def save_checkpoint(self, checkpoint: Checkpoint) -> int:
        """Save checkpoint to database.

        Args:
            checkpoint: Checkpoint to save

        Returns:
            Checkpoint ID

        Raises:
            RepositoryError: If save operation fails
        """
        if self._db is None:
            raise RepositoryError("Repository not initialized")

        try:
            checkpoint_id = self._db.create_checkpoint(checkpoint)
            logger.info(
                f"Checkpoint saved: {checkpoint.feature} "
                f"(ID: {checkpoint_id}, agent: {checkpoint.agent_id})"
            )
            return checkpoint_id
        except Exception as e:
            logger.error(f"Failed to save checkpoint for {checkpoint.feature}: {e}")
            raise RepositoryError(f"Failed to save checkpoint: {e}") from e

    def load_checkpoint(self, feature: str) -> Optional[Checkpoint]:
        """Load latest checkpoint for feature.

        Args:
            feature: Feature ID

        Returns:
            Checkpoint or None if not found
        """
        if self._db is None:
            raise RepositoryError("Repository not initialized")

        try:
            logger.info(f"Loading checkpoint for feature: {feature}")
            checkpoint = self._db.get_checkpoint_by_feature(feature)
            return checkpoint
        except Exception as e:
            logger.error(f"Failed to load checkpoint for {feature}: {e}")
            raise RepositoryError(f"Failed to load checkpoint: {e}") from e

    def load_latest_checkpoint(self, feature: str) -> Optional[Checkpoint]:
        """Load latest in_progress or failed checkpoint for resume.

        Args:
            feature: Feature ID

        Returns:
            Latest active checkpoint or None
        """
        if self._db is None:
            raise RepositoryError("Repository not initialized")

        try:
            checkpoint = self._db.get_checkpoint_by_feature(feature)
            if checkpoint and checkpoint.status in (
                CheckpointStatus.IN_PROGRESS,
                CheckpointStatus.FAILED,
            ):
                logger.info(f"Latest active checkpoint found for {feature}")
                return checkpoint
            return None
        except Exception as e:
            logger.error(f"Failed to load latest checkpoint for {feature}: {e}")
            raise RepositoryError(f"Failed to load latest checkpoint: {e}") from e

    def update_checkpoint_status(
        self,
        checkpoint_id: int,
        new_status: CheckpointStatus,
        completed_ws: list[str],
    ) -> None:
        """Update checkpoint status and completed workstreams.

        Args:
            checkpoint_id: Checkpoint ID
            new_status: New status
            completed_ws: List of completed workstream IDs
        """
        if self._db is None:
            raise RepositoryError("Repository not initialized")

        try:
            # Load existing checkpoint
            checkpoint = self._db.get_checkpoint(checkpoint_id)
            if checkpoint is None:
                raise RepositoryError(f"Checkpoint {checkpoint_id} not found")

            # Update fields
            checkpoint.status = new_status
            checkpoint.completed_ws = completed_ws

            # Save back
            self._db.update_checkpoint(checkpoint_id, checkpoint)
            logger.info(
                f"Checkpoint updated: {checkpoint.feature} "
                f"(status: {new_status.value}, completed: {len(completed_ws)} WS)"
            )
        except RepositoryError:
            raise
        except Exception as e:
            logger.error(f"Failed to update checkpoint {checkpoint_id}: {e}")
            raise RepositoryError(f"Failed to update checkpoint: {e}") from e

    def list_active_checkpoints(self) -> list[Checkpoint]:
        """List all in_progress or failed checkpoints.

        Returns:
            List of active checkpoints
        """
        if self._db is None:
            raise RepositoryError("Repository not initialized")

        try:
            return self._db.get_active_checkpoints()
        except Exception as e:
            logger.error(f"Failed to list active checkpoints: {e}")
            raise RepositoryError(f"Failed to list active checkpoints: {e}") from e

    def close(self) -> None:
        """Close database connection."""
        if self._db:
            self._db.close()
            self._db = None

    def __enter__(self) -> "CheckpointRepository":
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
