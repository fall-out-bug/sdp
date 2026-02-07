"""
Checkpoint serialization utilities.

Handles conversion between database rows and Checkpoint objects.
"""

import json
import sqlite3
from datetime import datetime

from .models import Checkpoint, CheckpointStatus


def row_to_checkpoint(row: sqlite3.Row) -> Checkpoint:
    """Convert database row to Checkpoint object.

    Args:
        row: Database row

    Returns:
        Checkpoint object
    """
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


def checkpoint_to_insert_params(checkpoint: Checkpoint) -> tuple:
    """Convert Checkpoint to database insert parameters.

    Args:
        checkpoint: Checkpoint to serialize

    Returns:
        Tuple of values for INSERT query
    """
    return (
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
    )


def checkpoint_to_update_params(checkpoint: Checkpoint, checkpoint_id: int) -> tuple:
    """Convert Checkpoint to database update parameters.

    Args:
        checkpoint: Checkpoint to serialize
        checkpoint_id: ID of checkpoint to update

    Returns:
        Tuple of values for UPDATE query
    """
    return (
        checkpoint.status.value,
        json.dumps(checkpoint.completed_ws),
        json.dumps(checkpoint.execution_order),
        checkpoint.current_ws,
        checkpoint.completed_at.isoformat() if checkpoint.completed_at else None,
        json.dumps(checkpoint.failed_tasks),
        checkpoint.error,
        json.dumps(checkpoint.metrics),
        checkpoint_id,
    )
