"""
Multi-agent oneshot execution using Beads ready detection.

Implements @oneshot skill logic: execute all workstreams for a feature
using multiple agents in parallel with Beads dependency tracking.
"""

from dataclasses import dataclass, field
from typing import List, Optional
from concurrent.futures import ThreadPoolExecutor, as_completed
from .client import BeadsClient, BeadsStatus
from .skills_build import WorkstreamExecutor


@dataclass
class OneshotResult:
    """Result of oneshot feature execution."""

    success: bool
    feature_id: str
    total_executed: int = 0
    error: Optional[str] = None
    failed_tasks: List[str] = field(default_factory=list)


class MultiAgentExecutor:
    """Execute feature workstreams with multi-agent coordination.

    Uses Beads `get_ready_tasks()` to discover executable workstreams
    and executes them in parallel using ThreadPoolExecutor.
    """

    def __init__(self, client: BeadsClient, num_agents: int = 3):
        """Initialize multi-agent executor.

        Args:
            client: BeadsClient instance (mock or real)
            num_agents: Maximum number of parallel agents
        """
        self.client = client
        self.num_agents = num_agents
        self.build_executor = WorkstreamExecutor(client)

    def execute_feature(
        self,
        feature_id: str,
        mock_success: bool = True,
    ) -> OneshotResult:
        """Execute all workstreams for a feature.

        Uses get_ready_tasks() to discover executable workstreams,
        executes them in parallel, and repeats until none remain.

        Args:
            feature_id: Parent feature task ID
            mock_success: Mock success for testing

        Returns:
            OneshotResult with execution summary

        Example:
            executor = MultiAgentExecutor(client, num_agents=3)
            result = executor.execute_feature("bd-0001")

            # Executes ready tasks in parallel:
            # Round 1: bd-0001.1, bd-0001.4 (parallel)
            # Round 2: bd-0001.2 (after bd-0001.1 completes)
            # Round 3: bd-0001.3 (after bd-0001.2 completes)
            # Round 4: bd-0001.5 (after bd-0001.4 completes)
        """
        total_executed = 0
        failed_tasks = []

        try:
            with ThreadPoolExecutor(max_workers=self.num_agents) as executor:
                while True:
                    # Get ready tasks
                    ready_tasks = self.client.get_ready_tasks()

                    # Filter to only this feature's sub-tasks
                    feature_tasks = self._filter_feature_tasks(ready_tasks, feature_id)

                    if not feature_tasks:
                        # No more workstreams for this feature
                        break

                    # Execute ready tasks in parallel
                    futures = {
                        executor.submit(
                            self._execute_single, task_id, mock_success
                        ): task_id
                        for task_id in feature_tasks
                    }

                    # Wait for completion and collect results
                    for future in as_completed(futures):
                        task_id = futures[future]
                        try:
                            success = future.result()
                            total_executed += 1

                            if not success:
                                failed_tasks.append(task_id)

                        except Exception as e:
                            failed_tasks.append(task_id)
                            total_executed += 1

            # Check if any tasks failed
            if failed_tasks:
                return OneshotResult(
                    success=False,
                    feature_id=feature_id,
                    total_executed=total_executed,
                    error=f"{len(failed_tasks)} tasks failed: {failed_tasks}",
                    failed_tasks=failed_tasks,
                )
            else:
                return OneshotResult(
                    success=True,
                    feature_id=feature_id,
                    total_executed=total_executed,
                )

        except Exception as e:
            return OneshotResult(
                success=False,
                feature_id=feature_id,
                total_executed=total_executed,
                error=str(e),
            )

    def _filter_feature_tasks(self, task_ids: List[str], feature_id: str) -> List[str]:
        """Filter tasks to only include sub-tasks of this feature.

        Args:
            task_ids: List of ready task IDs
            feature_id: Parent feature task ID

        Returns:
            Filtered list of task IDs
        """
        feature_tasks = []

        for task_id in task_ids:
            task = self.client.get_task(task_id)
            if task and task.parent_id == feature_id:
                feature_tasks.append(task_id)

        return feature_tasks

    def _execute_single(self, task_id: str, mock_success: bool) -> bool:
        """Execute a single workstream.

        Args:
            task_id: Beads task ID
            mock_success: Mock success for testing

        Returns:
            True if successful, False otherwise
        """
        result = self.build_executor.execute(task_id, mock_tdd_success=mock_success)
        return result.success
