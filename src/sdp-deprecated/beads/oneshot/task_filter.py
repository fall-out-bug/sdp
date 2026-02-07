"""
Task filtering utilities for @oneshot execution.

Handles filtering and execution of individual workstream tasks.
"""

from typing import TYPE_CHECKING, List

if TYPE_CHECKING:
    from ..client import BeadsClient
    from ..skills_build import WorkstreamExecutor


def filter_feature_tasks(
    client: "BeadsClient", task_ids: List[str], feature_id: str
) -> List[str]:
    """Filter tasks to only include sub-tasks of this feature.

    Args:
        client: BeadsClient instance
        task_ids: List of ready task IDs
        feature_id: Parent feature task ID

    Returns:
        Filtered list of task IDs
    """
    feature_tasks = []

    for task_id in task_ids:
        task = client.get_task(task_id)
        if task and task.parent_id == feature_id:
            feature_tasks.append(task_id)

    return feature_tasks


def execute_single_task(
    build_executor: "WorkstreamExecutor", task_id: str, mock_success: bool
) -> bool:
    """Execute a single workstream.

    Args:
        build_executor: WorkstreamExecutor instance
        task_id: Beads task ID
        mock_success: Mock success for testing

    Returns:
        True if successful, False otherwise
    """
    result = build_executor.execute(task_id, mock_tdd_success=mock_success)
    return result.success
