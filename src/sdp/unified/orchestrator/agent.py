"""
OrchestratorAgent for autonomous feature execution.

Implements core logic for @oneshot workflow with checkpoint management.
"""

import logging
from datetime import datetime

from sdp.unified.checkpoint.repository import CheckpointRepository
from sdp.unified.checkpoint.schema import Checkpoint, CheckpointStatus
from sdp.unified.orchestrator.agent_extension import AgentCheckpointExtension
from sdp.unified.orchestrator.dispatcher import WorkstreamDispatcher
from sdp.unified.orchestrator.errors import ExecutionError
from sdp.unified.orchestrator.models import ExecutionResult
from sdp.unified.orchestrator.monitor import ProgressMonitor

logger = logging.getLogger(__name__)


class OrchestratorAgent:
    """Orchestrates autonomous feature execution with checkpoint management."""

    def __init__(self, repo: CheckpointRepository) -> None:
        """Initialize orchestrator.

        Args:
            repo: CheckpointRepository for state management

        Raises:
            TypeError: If repo is not CheckpointRepository
        """
        if not isinstance(repo, CheckpointRepository):
            raise TypeError("repo must be CheckpointRepository")
        self.repo = repo
        self.dispatcher = WorkstreamDispatcher()
        self.monitor = ProgressMonitor(repo)
        self.checkpoint_ext = AgentCheckpointExtension(repo)

    def execute_feature(
        self,
        feature_id: str,
        workstreams: list[str],
        agent_id: str,
    ) -> ExecutionResult:
        """Execute feature with workstream orchestration.

        Args:
            feature_id: Feature identifier (e.g., "F01")
            workstreams: List of workstream IDs in execution order
            agent_id: Agent identifier

        Returns:
            ExecutionResult with execution details

        Raises:
            ExecutionError: If execution fails
        """
        started_at = datetime.now()
        logger.info(f"Starting execution of feature: {feature_id} (agent: {agent_id})")

        try:
            # Check for existing checkpoint to resume
            existing_checkpoint = self.repo.load_latest_checkpoint(feature_id)

            previously_completed: list[str] = []
            checkpoint_id = 0
            start_index = 0

            if existing_checkpoint:
                logger.info(f"Resuming from checkpoint: {existing_checkpoint.current_ws}")
                checkpoint_id = (
                    existing_checkpoint.id if hasattr(existing_checkpoint, "id") else 0
                )
                previously_completed = (
                    list(existing_checkpoint.completed_ws)
                    if existing_checkpoint.completed_ws
                    else []
                )

                # Find current position
                if previously_completed:
                    last_completed = previously_completed[-1]
                    if last_completed in workstreams:
                        start_index = workstreams.index(last_completed) + 1
            else:
                # Create new checkpoint
                checkpoint = Checkpoint(
                    feature=feature_id,
                    agent_id=agent_id,
                    status=CheckpointStatus.IN_PROGRESS,
                    completed_ws=[],
                    execution_order=workstreams,
                    started_at=started_at,
                )
                checkpoint_id = self.repo.save_checkpoint(checkpoint)

            # Dispatch workstreams
            newly_completed_ws = self.dispatch_workstreams(
                workstreams=workstreams,
                checkpoint_id=checkpoint_id,
                start_index=start_index,
            )

            # Combine previous and new completions
            all_completed_ws = previously_completed + newly_completed_ws

            # Mark as completed
            completed_at = datetime.now()
            duration = (completed_at - started_at).total_seconds()

            self.repo.update_checkpoint_status(
                checkpoint_id=checkpoint_id,
                new_status=CheckpointStatus.COMPLETED,
                completed_ws=all_completed_ws,
            )

            logger.info(f"Feature execution completed: {feature_id}")

            return ExecutionResult(
                feature_id=feature_id,
                agent_id=agent_id,
                status=CheckpointStatus.COMPLETED,
                completed_workstreams=all_completed_ws,
                failed_workstreams=[],
                duration_seconds=duration,
                started_at=started_at,
                completed_at=completed_at,
            )

        except Exception as e:
            logger.error(f"Failed to execute feature {feature_id}: {e}")
            raise ExecutionError(f"Failed to execute feature: {e}") from e

    def dispatch_workstreams(
        self,
        workstreams: list[str],
        checkpoint_id: int,
        start_index: int,
    ) -> list[str]:
        """Dispatch workstreams for execution.

        Args:
            workstreams: List of workstream IDs
            checkpoint_id: Checkpoint ID for updates
            start_index: Index to start from (for resume)

        Returns:
            List of completed workstream IDs

        Raises:
            ExecutionError: If dispatch fails
        """
        completed_ws: list[str] = []

        try:
            for i in range(start_index, len(workstreams)):
                ws_id = workstreams[i]
                logger.info(f"Dispatching workstream: {ws_id}")

                # Dispatch single workstream
                self._dispatch_single_workstream(ws_id)
                completed_ws.append(ws_id)

                # Update checkpoint progress
                self.repo.update_checkpoint_status(
                    checkpoint_id=checkpoint_id,
                    new_status=CheckpointStatus.IN_PROGRESS,
                    completed_ws=completed_ws,
                )

            return completed_ws

        except Exception as e:
            logger.error(f"Workstream dispatch failed: {e}")
            raise ExecutionError(f"Workstream dispatch failed: {e}") from e

    def _dispatch_single_workstream(self, ws_id: str) -> None:
        """Dispatch a single workstream.

        Args:
            ws_id: Workstream ID to dispatch

        Note:
            This delegates to WorkstreamDispatcher, which will integrate
            with the Task tool in WS-012.
        """
        self.dispatcher.dispatch(ws_id)

    def monitor_progress(self, feature_id: str) -> dict[str, object] | None:
        """Monitor progress of feature execution.

        Args:
            feature_id: Feature identifier

        Returns:
            Progress metrics dict or None if no checkpoint found

        Note:
            Delegates to ProgressMonitor for actual implementation.
        """
        return self.monitor.get_progress(feature_id)
