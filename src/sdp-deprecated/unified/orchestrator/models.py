"""
Data models for orchestrator agent.
"""

from dataclasses import dataclass
from datetime import datetime

from sdp.unified.checkpoint.schema import CheckpointStatus


@dataclass
class ExecutionResult:
    """Result of feature execution."""

    feature_id: str
    agent_id: str
    status: CheckpointStatus
    completed_workstreams: list[str]
    failed_workstreams: list[str]
    duration_seconds: float
    started_at: datetime = datetime.now()
    completed_at: datetime = datetime.now()

    @property
    def is_success(self) -> bool:
        """Check if execution was successful."""
        return (
            self.status == CheckpointStatus.COMPLETED
            and len(self.failed_workstreams) == 0
        )
