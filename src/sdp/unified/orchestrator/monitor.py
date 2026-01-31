"""
Progress monitoring for feature execution.
"""

import logging
from typing import Any

from sdp.unified.checkpoint.repository import CheckpointRepository

logger = logging.getLogger(__name__)


class ProgressMonitor:
    """Monitors execution progress for features."""

    def __init__(self, repo: CheckpointRepository) -> None:
        """Initialize progress monitor.

        Args:
            repo: CheckpointRepository for loading checkpoint data
        """
        self.repo = repo

    def get_progress(self, feature_id: str) -> dict[str, Any] | None:
        """Get progress metrics for feature execution.

        Args:
            feature_id: Feature identifier

        Returns:
            Progress metrics dict or None if no checkpoint found

            Metrics include:
            - feature_id: Feature identifier
            - total_workstreams: Total number of workstreams
            - completed_workstreams: Number of completed workstreams
            - current_workstream: Currently executing workstream
            - status: Checkpoint status
            - percentage: Completion percentage (0-100)
            - started_at: Start time (ISO format)
            - completed_at: Completion time (ISO format) if completed
        """
        checkpoint = self.repo.load_checkpoint(feature_id)

        if checkpoint is None:
            return None

        total = len(checkpoint.execution_order)
        completed = len(checkpoint.completed_ws)

        return {
            "feature_id": feature_id,
            "total_workstreams": total,
            "completed_workstreams": completed,
            "current_workstream": checkpoint.current_ws,
            "status": checkpoint.status.value,
            "percentage": (completed / total * 100) if total > 0 else 0,
            "started_at": checkpoint.started_at.isoformat(),
            "completed_at": (
                checkpoint.completed_at.isoformat() if checkpoint.completed_at else None
            ),
        }
