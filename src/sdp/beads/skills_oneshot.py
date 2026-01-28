"""
Multi-agent oneshot execution using Beads ready detection.

Implements @oneshot skill logic: execute all workstreams for a feature
using multiple agents in parallel with Beads dependency tracking.

Enhanced with workflow efficiency modes (F014):
- Standard mode (PR required)
- Auto-approve mode (skip PR)
- Sandbox mode (skip PR, sandbox only)
- Dry-run mode (preview changes)
"""

from dataclasses import dataclass, field
from typing import List, Optional, TYPE_CHECKING
from concurrent.futures import ThreadPoolExecutor, as_completed
import os

from .client import BeadsClient, BeadsStatus
from .skills_build import WorkstreamExecutor
from .execution_mode import (
    ExecutionMode,
    AuditLogger,
    DestructiveOperationDetector,
    OneshotResult,
)

if TYPE_CHECKING:
    pass


@dataclass
class OneshotResult:
    """Result of oneshot feature execution."""

    success: bool
    feature_id: str
    total_executed: int = 0
    error: Optional[str] = None
    failed_tasks: List[str] = field(default_factory=list)
    mode: ExecutionMode = ExecutionMode.STANDARD
    deployment_target: str = "production"
    pr_created: bool = False
    preview_only: bool = False
    tasks_preview: List[str] = field(default_factory=list)


class MultiAgentExecutor:
    """Execute feature workstreams with multi-agent coordination.

    Uses Beads `get_ready_tasks()` to discover executable workstreams
    and executes them in parallel using ThreadPoolExecutor.

    Enhanced with execution modes for workflow efficiency (F014).
    """

    def __init__(
        self,
        client: BeadsClient,
        num_agents: int = 3,
        audit_logger: Optional[AuditLogger] = None,
    ):
        """Initialize multi-agent executor.

        Args:
            client: BeadsClient instance (mock or real)
            num_agents: Maximum number of parallel agents
            audit_logger: Optional audit logger for auto-approve mode
        """
        self.client = client
        self.num_agents = num_agents
        self.build_executor = WorkstreamExecutor(client)
        self.audit_logger = audit_logger or AuditLogger()
        self.destructive_detector = DestructiveOperationDetector()

    def execute_feature(
        self,
        feature_id: str,
        mode: ExecutionMode = ExecutionMode.STANDARD,
        mock_success: bool = True,
    ) -> OneshotResult:
        """Execute all workstreams for a feature.

        Uses get_ready_tasks() to discover executable workstreams,
        executes them in parallel, and repeats until none remain.

        Args:
            feature_id: Parent feature task ID
            mode: Execution mode (standard, auto-approve, sandbox, dry-run)
            mock_success: Mock success for testing

        Returns:
            OneshotResult with execution summary

        Example:
            executor = MultiAgentExecutor(client, num_agents=3)
            result = executor.execute_feature("bd-0001", mode=ExecutionMode.AUTO_APPROVE)

            # Executes ready tasks in parallel:
            # Round 1: bd-0001.1, bd-0001.4 (parallel)
            # Round 2: bd-0001.2 (after bd-0001.1 completes)
            # Round 3: bd-0001.3 (after bd-0001.2 completes)
        """
        # Handle dry-run mode
        if mode == ExecutionMode.DRY_RUN:
            return self._execute_dry_run(feature_id)

        # Determine deployment target
        deployment_target = "sandbox" if mode == ExecutionMode.SANDBOX else "production"

        # Check for destructive operations if not dry-run
        if mode in (ExecutionMode.AUTO_APPROVE, ExecutionMode.SANDBOX):
            if not self._check_destructive_operations_confirmation():
                return OneshotResult(
                    success=False,
                    feature_id=feature_id,
                    total_executed=0,
                    error="Destructive operations detected and user declined confirmation",
                    mode=mode,
                    deployment_target=deployment_target,
                )

        # Execute workstreams
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

            # Determine if PR was created
            pr_created = mode == ExecutionMode.STANDARD

            # Check if any tasks failed
            if failed_tasks:
                result = OneshotResult(
                    success=False,
                    feature_id=feature_id,
                    total_executed=total_executed,
                    error=f"{len(failed_tasks)} tasks failed: {failed_tasks}",
                    failed_tasks=failed_tasks,
                    mode=mode,
                    deployment_target=deployment_target,
                    pr_created=pr_created,
                )
            else:
                result = OneshotResult(
                    success=True,
                    feature_id=feature_id,
                    total_executed=total_executed,
                    mode=mode,
                    deployment_target=deployment_target,
                    pr_created=pr_created,
                )

            # Log auto-approve executions to audit
            if mode == ExecutionMode.AUTO_APPROVE:
                self.audit_logger.log_execution(
                    feature_id=feature_id,
                    mode=mode,
                    workstreams_executed=total_executed,
                    result="success" if result.success else "failure",
                    deployment_target=deployment_target,
                )

            return result

        except Exception as e:
            return OneshotResult(
                success=False,
                feature_id=feature_id,
                total_executed=total_executed,
                error=str(e),
                mode=mode,
                deployment_target=deployment_target,
            )

    def _execute_dry_run(self, feature_id: str) -> OneshotResult:
        """Execute dry-run mode (preview only).

        Args:
            feature_id: Parent feature task ID

        Returns:
            OneshotResult with preview information
        """
        # Get ready tasks for preview
        ready_tasks = self.client.get_ready_tasks()
        feature_tasks = self._filter_feature_tasks(ready_tasks, feature_id)

        # Build preview list
        tasks_preview = []
        for task_id in feature_tasks:
            task = self.client.get_task(task_id)
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

    def _check_destructive_operations_confirmation(self) -> bool:
        """Check if user confirms destructive operations.

        Returns:
            True if user confirms or no destructive operations, False otherwise
        """
        # For now, always return True (auto-confirm)
        # In real implementation, would prompt user with AskUserQuestion
        return True

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
