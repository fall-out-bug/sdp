"""
Checkpoint data models.

Defines core data structures for checkpoint storage.
"""

from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum
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
