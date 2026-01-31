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

from concurrent.futures import ThreadPoolExecutor, as_completed
from typing import TYPE_CHECKING, Any, List, Optional

from .client import BeadsClient
from .execution_mode import (
    AuditLogger,
    DestructiveOperationDetector,
    ExecutionMode,
    OneshotResult,
)
from .skills_build import WorkstreamExecutor

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
        self.destructive_detector = DestructiveOperationDetector()

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
            return self._execute_dry_run(feature_id)

        # Determine deployment target
        deployment_target = "sandbox" if mode == ExecutionMode.SANDBOX else "production"

        # Check for destructive operations if not dry-run
        if mode in (ExecutionMode.AUTO_APPROVE, ExecutionMode.SANDBOX):
            if not self._check_destructive_operations_confirmation(feature_id):
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

    def _check_destructive_operations_confirmation(self, feature_id: str) -> bool:
        """Check if user confirms destructive operations.

        Args:
            feature_id: Parent feature task ID

        Returns:
            True if user confirms or no destructive operations, False otherwise

        Implementation:
        1. Get all subtasks for the feature
        2. Check task titles/descriptions for destructive keywords
        3. Prompt user for confirmation if found
        4. Return False if user declines
        """
        try:
            # Get all subtasks for this feature
            all_tasks = self.client.list_tasks(parent_id=feature_id)

            if not all_tasks:
                # No tasks found, auto-confirm
                return True

            # Check each task for destructive patterns
            destructive_tasks = []

            for task in all_tasks:
                # Combine title and description for checking
                text_to_check = f"{task.title} {task.description or ''}"

                # Check against destructive patterns
                for category, patterns in DestructiveOperationDetector.DESTRUCTIVE_PATTERNS.items():
                    for pattern in patterns:
                        if pattern.lower() in text_to_check.lower():
                            destructive_tasks.append({
                                'task_id': task.id,
                                'title': task.title,
                                'operation_type': category,
                                'pattern': pattern
                            })
                            break

            if not destructive_tasks:
                # No destructive operations found, auto-confirm
                return True

            # Destructive operations detected - need user confirmation
            operation_summary = self._build_destructive_operations_summary(destructive_tasks)

            # Prompt user for confirmation
            return self._console_prompt_confirmation(operation_summary)

        except Exception:
            # If detection fails, fail-open (allow execution)
            # This prevents blocking all executions due to bugs in detection logic
            return True

    def _build_destructive_operations_summary(
        self, destructive_tasks: list[dict[str, Any]]
    ) -> str:
        """Build a summary of destructive operations for user confirmation.

        Args:
            destructive_tasks: List of dicts with operation details

        Returns:
            Formatted summary string
        """
        lines = [
            "⚠️  DESTRUCTIVE OPERATIONS DETECTED",
            "",
            f"Found {len(destructive_tasks)} task(s) with destructive operations:",
            ""
        ]

        for task_info in destructive_tasks:
            lines.append(f"  • {task_info['task_id']}: {task_info['title'][:60]}")
            lines.append(f"    Operation: {task_info['operation_type']}")
            lines.append(f"    Pattern matched: {task_info['pattern']}")
            lines.append("")

        lines.extend([
            "These operations may:",
            "  - Delete files or data",
            "  - Modify database schema",
            "  - Cause irreversible changes",
            "",
            "Do you want to proceed?"
        ])

        return "\n".join(lines)

    def _console_prompt_confirmation(self, summary: str) -> bool:
        """Prompt user via console for confirmation.

        Args:
            summary: Summary message to show user

        Returns:
            True if user confirms, False otherwise
        """
        print("\n" + "=" * 70)
        print(summary)
        print("=" * 70 + "\n")

        while True:
            response = input("Type 'yes' to proceed, 'no' to cancel: ").strip().lower()

            if response in ('yes', 'y'):
                print("✓ User confirmed - proceeding with destructive operations\n")
                return True
            elif response in ('no', 'n'):
                print("✗ User declined - cancelling execution\n")
                return False
            else:
                print("Please type 'yes' or 'no'\n")

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
