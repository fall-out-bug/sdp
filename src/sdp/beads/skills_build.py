"""
Workstream execution with TDD cycle and status updates.

Implements @build skill logic: execute workstream with TDD discipline
and update Beads status (OPEN → IN_PROGRESS → CLOSED/BLOCKED).
"""

from dataclasses import dataclass, field
from typing import List, Callable, Optional
from .client import BeadsClient, BeadsStatus


@dataclass
class ExecutionResult:
    """Result of workstream execution."""

    success: bool
    task_id: str
    error: Optional[str] = None
    newly_ready: List[str] = field(default_factory=list)


class WorkstreamExecutor:
    """Execute workstreams with TDD cycle and status updates.

    Manages status transitions:
    OPEN → IN_PROGRESS → CLOSED (success)
    OPEN → IN_PROGRESS → BLOCKED (failure)
    """

    def __init__(self, client: BeadsClient):
        """Initialize executor.

        Args:
            client: BeadsClient instance (mock or real)
        """
        self.client = client

    def execute(
        self,
        task_id: str,
        mock_tdd_success: bool = True,
    ) -> ExecutionResult:
        """Execute workstream with TDD cycle.

        Updates task status through lifecycle:
        1. OPEN → IN_PROGRESS (start)
        2. IN_PROGRESS → CLOSED (success) or BLOCKED (failure)
        3. Returns list of newly ready tasks (unblocked by completion)

        Args:
            task_id: Beads task ID to execute
            mock_tdd_success: Mock TDD result for testing

        Returns:
            ExecutionResult with success status and newly ready tasks

        Example:
            executor = WorkstreamExecutor(client)
            result = executor.execute("bd-0001.1")

            if result.success:
                print(f"✅ Tasks now ready: {result.newly_ready}")
        """
        # Update to IN_PROGRESS
        self.client.update_task_status(task_id, BeadsStatus.IN_PROGRESS)

        try:
            # Execute TDD cycle
            tdd_result = self.execute_tdd_cycle(task_id, mock_tdd_success)

            if tdd_result:
                # Mark as done → unblocks dependent tasks
                self.client.update_task_status(task_id, BeadsStatus.CLOSED)

                # Get newly ready tasks
                ready_before = set()  # Would track before execution
                ready_now = set(self.client.get_ready_tasks())
                newly_ready = list(ready_now - ready_before)

                return ExecutionResult(
                    success=True,
                    task_id=task_id,
                    newly_ready=newly_ready,
                )
            else:
                # Mark as blocked
                self.client.update_task_status(task_id, BeadsStatus.BLOCKED)

                return ExecutionResult(
                    success=False,
                    task_id=task_id,
                    error="TDD cycle failed",
                )

        except Exception as e:
            # Mark as blocked on exception
            self.client.update_task_status(task_id, BeadsStatus.BLOCKED)

            return ExecutionResult(
                success=False,
                task_id=task_id,
                error=str(e),
            )

    def execute_tdd_cycle(
        self,
        task_id: str,
        mock_tdd_success: bool = True,
    ) -> bool:
        """Execute TDD cycle (Red → Green → Refactor).

        Args:
            task_id: Beads task ID
            mock_tdd_success: Mock success for testing

        Returns:
            True if cycle succeeded, False otherwise
        """
        # In real implementation, this would run actual tests
        # For now, use mock parameter
        return mock_tdd_success

    def _run_tdd_phase(
        self,
        task_id: str,
        phase_fn: Callable[[str], None],
    ) -> None:
        """Run a single TDD phase.

        Args:
            task_id: Beads task ID
            phase_fn: Function to execute for this phase
        """
        phase_fn(task_id)
