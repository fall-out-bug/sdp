"""PRD annotation decorators for Python code.

This module provides decorators that can be used to annotate Python code
with PRD flow information for automatic diagram generation.
"""

from collections.abc import Callable
from typing import TypeVar

F = TypeVar('F', bound=Callable[..., object])


def prd_flow(flow_name: str) -> Callable[[F], F]:
    """Mark function as part of a PRD flow.

    This decorator marks a function as belonging to a specific PRD flow.
    It should be used together with @prd_step to provide complete information.

    Args:
        flow_name: Name of the flow this function belongs to

    Returns:
        Decorator function

    Example:
        ```python
        @prd_flow("submission-processing")
        @prd_step(1, "Receive submission from queue")
        async def process_submission(self, job: Job) -> RunResult:
            ...
        ```
    """
    def decorator(func: F) -> F:
        func._prd_flow = flow_name  # type: ignore
        return func
    return decorator


def prd_step(step_number: int, description: str) -> Callable[[F], F]:
    """Mark function as a step in the PRD flow.

    This decorator marks a function as a specific step in a PRD flow
    and should be used together with @prd_flow.

    Args:
        step_number: Sequential number of this step in the flow
        description: Human-readable description of what this step does

    Returns:
        Decorator function

    Example:
        ```python
        @prd_flow("submission-processing")
        @prd_step(1, "Receive submission from queue")
        async def process_submission(self, job: Job) -> RunResult:
            ...
        ```
    """
    def decorator(func: F) -> F:
        func._prd_step = step_number  # type: ignore
        func._prd_step_desc = description  # type: ignore
        return func
    return decorator


def get_flow_info(func: Callable[..., object]) -> tuple[str | None, int | None, str | None]:
    """Extract PRD flow information from a function.

    Args:
        func: The function to extract info from

    Returns:
        Tuple of (flow_name, step_number, description) or (None, None, None)
    """
    flow_name = getattr(func, '_prd_flow', None)
    step_number = getattr(func, '_prd_step', None)
    description = getattr(func, '_prd_step_desc', None)
    return flow_name, step_number, description
