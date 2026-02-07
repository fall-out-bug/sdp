"""Checkpoint save/resume extension for OrchestratorAgent.

This module provides extension methods for file-based checkpoint
persistence with approval gate and team configuration integration.
"""

import logging
from typing import Any, Optional

from sdp.unified.checkpoint.repository import CheckpointRepository
from sdp.unified.checkpoint.schema import CheckpointStatus
from sdp.unified.gates.manager import ApprovalGateManager
from sdp.unified.orchestrator.checkpoint import CheckpointFileManager
from sdp.unified.orchestrator.checkpoint_ops import CheckpointOperations
from sdp.unified.team.manager import TeamManager

logger = logging.getLogger(__name__)


class AgentCheckpointExtension:
    """Extension for OrchestratorAgent to add checkpoint save/resume.

    Provides methods for file-based checkpoint persistence with
    integration for approval gates and team configuration.
    """

    def __init__(self, repo: CheckpointRepository) -> None:
        """Initialize checkpoint extension.

        Args:
            repo: Checkpoint repository for database state
        """
        self.repo = repo
        self.checkpoint_manager: Optional[CheckpointFileManager] = None
        self.checkpoint_ops: Optional[CheckpointOperations] = None
        self.gate_manager: Optional[ApprovalGateManager] = None
        self.team_manager: Optional[TeamManager] = None

    def set_checkpoint_manager(
        self, checkpoint_manager: CheckpointFileManager
    ) -> None:
        """Set checkpoint file manager.

        Args:
            checkpoint_manager: Checkpoint file manager instance
        """
        self.checkpoint_manager = checkpoint_manager
        self.checkpoint_ops = CheckpointOperations(checkpoint_manager, self.repo)

        # Configure checkpoint_ops with gate and team managers if available
        if self.gate_manager:
            self.checkpoint_ops.set_gate_manager(self.gate_manager)
        if self.team_manager:
            self.checkpoint_ops.set_team_manager(self.team_manager)

    def set_gate_manager(self, gate_manager: ApprovalGateManager) -> None:
        """Set approval gate manager.

        Args:
            gate_manager: Approval gate manager instance
        """
        self.gate_manager = gate_manager
        if self.checkpoint_ops:
            self.checkpoint_ops.set_gate_manager(gate_manager)

    def set_team_manager(self, team_manager: TeamManager) -> None:
        """Set team manager.

        Args:
            team_manager: Team manager instance
        """
        self.team_manager = team_manager
        if self.checkpoint_ops:
            self.checkpoint_ops.set_team_manager(team_manager)

    def save_checkpoint(
        self,
        feature_id: str,
        agent_id: str,
        workstreams: list[str],
        completed_workstreams: list[str],
        current_workstream: str,
        status: CheckpointStatus,
    ) -> None:
        """Save checkpoint to file for resumption.

        Args:
            feature_id: Feature identifier
            agent_id: Agent identifier
            workstreams: List of workstream IDs in execution order
            completed_workstreams: List of completed workstream IDs
            current_workstream: Current workstream ID
            status: Checkpoint status

        Raises:
            ValueError: If checkpoint_ops not initialized
        """
        if self.checkpoint_ops is None:
            raise ValueError("checkpoint_ops not initialized")

        self.checkpoint_ops.save_checkpoint(
            feature_id=feature_id,
            agent_id=agent_id,
            workstreams=workstreams,
            completed_workstreams=completed_workstreams,
            current_workstream=current_workstream,
            status=status,
        )

    def resume_from_checkpoint(
        self, feature_id: str, agent_id: str
    ) -> Optional[dict[str, Any]]:
        """Resume execution from checkpoint file.

        Args:
            feature_id: Feature identifier
            agent_id: Agent identifier for verification

        Returns:
            Checkpoint data dictionary or None if not found/invalid

        Raises:
            ValueError: If checkpoint_ops not initialized

        Note:
            Verifies agent ID matches checkpoint before resuming.
        """
        if self.checkpoint_ops is None:
            raise ValueError("checkpoint_ops not initialized")

        return self.checkpoint_ops.resume_from_checkpoint(feature_id, agent_id)
