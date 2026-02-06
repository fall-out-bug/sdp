"""
Multi-agent executor for @oneshot execution.

Coordinates parallel execution of feature workstreams with Beads dependency tracking.
"""

from concurrent.futures import ThreadPoolExecutor, as_completed
from typing import TYPE_CHECKING, Optional

from ..client import BeadsClient
from ..execution_mode import AuditLogger, ExecutionMode, OneshotResult
from ..skills_build import WorkstreamExecutor
from .destructive_checker import check_destructive_operations_confirmation
from .dry_run import execute_dry_run
from .task_filter import execute_single_task, filter_feature_tasks

if TYPE_CHECKING:
    pass


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

    def execute_feature(  # noqa: C901
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
            return execute_dry_run(self.client, feature_id, filter_feature_tasks)

        # Determine deployment target
        deployment_target = "sandbox" if mode == ExecutionMode.SANDBOX else "production"

        # Check for destructive operations if not dry-run
        if mode in (ExecutionMode.AUTO_APPROVE, ExecutionMode.SANDBOX):
            if not check_destructive_operations_confirmation(self.client, feature_id):
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
                    feature_tasks = filter_feature_tasks(
                        self.client, ready_tasks, feature_id
                    )

                    if not feature_tasks:
                        # No more workstreams for this feature
                        break

                    # Execute ready tasks in parallel
                    futures = {
                        executor.submit(
                            execute_single_task,
                            self.build_executor,
                            task_id,
                            mock_success,
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

                        except Exception:
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
