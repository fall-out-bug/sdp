"""Task tool wrapper for agent spawning.

Provides wrapper interface for Claude Code's Task tool to enable
mocking and testing of agent spawning functionality.
"""

from typing import Optional


class TaskResult:
    """Result from Task tool invocation."""

    def __init__(
        self,
        task_id: str,
        output_file: Optional[str] = None,
    ) -> None:
        """Initialize Task result.

        Args:
            task_id: Task/agent ID
            output_file: Output file path for background tasks
        """
        self.task_id = task_id
        self.output_file = output_file


def Task(  # noqa: N802
    description: str,
    prompt: str,
    subagent_type: Optional[str] = None,
    model: Optional[str] = None,
    run_in_background: bool = False,
) -> TaskResult:
    """Invoke Task tool to spawn agent.

    This is a wrapper function that will be replaced with the actual
    Task tool invocation in production. In tests, this function is
    mocked to return TaskResult objects.

    Args:
        description: Short description for the task
        prompt: Detailed prompt for the agent
        subagent_type: Optional subagent type
        model: Optional model to use
        run_in_background: Whether to run in background

    Returns:
        TaskResult with task_id and optional output_file
    """
    # This is a placeholder - in production this calls the actual Task tool
    # For now, return a mock result
    import time
    import uuid

    task_id = f"agent-{int(time.time())}-{uuid.uuid4().hex[:8]}"

    if run_in_background:
        output_file = f"/tmp/agent_{task_id}.log"
        return TaskResult(task_id=task_id, output_file=output_file)

    return TaskResult(task_id=task_id)
