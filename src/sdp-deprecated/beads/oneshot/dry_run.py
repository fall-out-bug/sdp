"""
Dry-run mode for @oneshot execution.

Provides preview-only mode without actual execution.
"""

from typing import TYPE_CHECKING

from ..execution_mode import ExecutionMode, OneshotResult

if TYPE_CHECKING:
    from ..client import BeadsClient


def execute_dry_run(client: "BeadsClient", feature_id: str, filter_fn) -> OneshotResult:
    """Execute dry-run mode (preview only).

    Args:
        client: BeadsClient instance
        feature_id: Parent feature task ID
        filter_fn: Task filtering function

    Returns:
        OneshotResult with preview information
    """
    # Get ready tasks for preview
    ready_tasks = client.get_ready_tasks()
    feature_tasks = filter_fn(client, ready_tasks, feature_id)

    # Build preview list
    tasks_preview = []
    for task_id in feature_tasks:
        task = client.get_task(task_id)
        if task:
            tasks_preview.append(f"{task_id}: {task.title}")

    return OneshotResult(
        success=True,
        feature_id=feature_id,
        total_executed=0,  # No actual execution
        mode=ExecutionMode.DRY_RUN,
        deployment_target="none",
        pr_created=False,
        preview_only=True,
        tasks_preview=tasks_preview,
    )
