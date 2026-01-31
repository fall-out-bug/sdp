"""Checkpoint file management for @oneshot workflow.

Provides file-based checkpoint storage for agent resumption, including
approval gate state and team configuration persistence.
"""

import json
import logging
from pathlib import Path
from typing import Any, Optional

logger = logging.getLogger(__name__)


class CheckpointFileManager:
    """Manager for checkpoint file operations.

    Handles saving, loading, and deleting checkpoint files in the
    .oneshot directory. Checkpoint files store execution state,
    approval gate decisions, and team configuration for resumption.

    Attributes:
        base_path: Base directory for checkpoint files
    """

    def __init__(self, base_path: str) -> None:
        """Initialize checkpoint file manager.

        Args:
            base_path: Base directory for checkpoint files
        """
        self.base_path = Path(base_path)
        self._ensure_directory()

    def _ensure_directory(self) -> None:
        """Create checkpoint directory if it doesn't exist."""
        self.base_path.mkdir(parents=True, exist_ok=True)
        logger.debug(f"Checkpoint directory ready: {self.base_path}")

    def _get_checkpoint_file(self, feature_id: str) -> Path:
        """Get checkpoint file path for feature.

        Args:
            feature_id: Feature identifier

        Returns:
            Path to checkpoint file
        """
        return self.base_path / f"{feature_id}-checkpoint.json"

    def save_checkpoint(self, feature_id: str, checkpoint_data: dict[str, Any]) -> None:
        """Save checkpoint data to file.

        Args:
            feature_id: Feature identifier
            checkpoint_data: Checkpoint data dictionary

        Raises:
            OSError: If file write fails
        """
        checkpoint_file = self._get_checkpoint_file(feature_id)

        try:
            with open(checkpoint_file, "w") as f:
                json.dump(checkpoint_data, f, indent=2)

            logger.info(f"Checkpoint saved: {feature_id} -> {checkpoint_file}")
        except OSError as e:
            logger.error(f"Failed to save checkpoint for {feature_id}: {e}")
            raise

    def load_checkpoint(self, feature_id: str) -> Optional[dict[str, Any]]:
        """Load checkpoint data from file.

        Args:
            feature_id: Feature identifier

        Returns:
            Checkpoint data dictionary or None if not found
        """
        checkpoint_file = self._get_checkpoint_file(feature_id)

        if not checkpoint_file.exists():
            logger.debug(f"Checkpoint file not found: {checkpoint_file}")
            return None

        try:
            with open(checkpoint_file) as f:
                data: dict[str, Any] = json.load(f)

            logger.debug(f"Checkpoint loaded: {feature_id}")
            return data
        except (OSError, json.JSONDecodeError) as e:
            logger.error(f"Failed to load checkpoint for {feature_id}: {e}")
            return None

    def delete_checkpoint(self, feature_id: str) -> None:
        """Delete checkpoint file.

        Args:
            feature_id: Feature identifier
        """
        checkpoint_file = self._get_checkpoint_file(feature_id)

        if checkpoint_file.exists():
            try:
                checkpoint_file.unlink()
                logger.info(f"Checkpoint deleted: {feature_id}")
            except OSError as e:
                logger.error(f"Failed to delete checkpoint for {feature_id}: {e}")
        else:
            logger.debug(f"Checkpoint file not found for deletion: {checkpoint_file}")

    def checkpoint_exists(self, feature_id: str) -> bool:
        """Check if checkpoint file exists.

        Args:
            feature_id: Feature identifier

        Returns:
            True if checkpoint file exists, False otherwise
        """
        checkpoint_file = self._get_checkpoint_file(feature_id)
        return checkpoint_file.exists()
