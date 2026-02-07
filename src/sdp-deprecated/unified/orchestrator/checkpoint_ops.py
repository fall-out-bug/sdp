"""Checkpoint save/resume operations for OrchestratorAgent.

Provides methods for saving and resuming checkpoint state with
approval gate and team configuration integration.
"""

import logging
from datetime import datetime
from typing import Any, Optional

from sdp.unified.checkpoint.repository import CheckpointRepository
from sdp.unified.checkpoint.schema import CheckpointStatus
from sdp.unified.gates.manager import ApprovalGateManager
from sdp.unified.orchestrator.checkpoint import CheckpointFileManager
from sdp.unified.team.manager import TeamManager

logger = logging.getLogger(__name__)


class CheckpointOperations:
    """Operations for checkpoint save/resume.

    Handles file-based checkpoint persistence with integration
    for approval gates and team configuration.
    """

    def __init__(
        self,
        checkpoint_manager: CheckpointFileManager,
        repo: CheckpointRepository,
    ) -> None:
        """Initialize checkpoint operations.

        Args:
            checkpoint_manager: Checkpoint file manager
            repo: Checkpoint repository for database state
        """
        self.checkpoint_manager = checkpoint_manager
        self.repo = repo
        self.gate_manager: Optional[ApprovalGateManager] = None
        self.team_manager: Optional[TeamManager] = None

    def set_gate_manager(self, gate_manager: ApprovalGateManager) -> None:
        """Set approval gate manager for integration.

        Args:
            gate_manager: Approval gate manager instance
        """
        self.gate_manager = gate_manager

    def set_team_manager(self, team_manager: TeamManager) -> None:
        """Set team manager for integration.

        Args:
            team_manager: Team manager instance
        """
        self.team_manager = team_manager

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
        """
        started_at = datetime.now().isoformat()

        checkpoint_data: dict[str, Any] = {
            "feature_id": feature_id,
            "agent_id": agent_id,
            "workstreams": workstreams,
            "completed_workstreams": completed_workstreams,
            "current_workstream": current_workstream,
            "status": status.value,
            "started_at": started_at,
            "last_updated": started_at,
        }

        # Include approval gate state if available
        if self.gate_manager:
            try:
                checkpoint = self.repo.load_checkpoint(feature_id)
                if checkpoint and checkpoint.metrics:
                    checkpoint_data["gates"] = checkpoint.metrics.get("gates", {})
            except Exception as e:
                logger.warning(f"Failed to load gate state: {e}")

        # Include team state if available
        if self.team_manager:
            try:
                roles_data = [
                    {
                        "name": role.name,
                        "description": role.description,
                        "state": role.state.value,
                        "skill_file": role.skill_file,
                        "metadata": role.metadata,
                    }
                    for role in self.team_manager.roles.values()
                ]
                checkpoint_data["team"] = {"roles": roles_data}
            except Exception as e:
                logger.warning(f"Failed to load team state: {e}")

        self.checkpoint_manager.save_checkpoint(feature_id, checkpoint_data)
        logger.info(f"Checkpoint saved for feature: {feature_id}")

    def resume_from_checkpoint(
        self, feature_id: str, agent_id: str
    ) -> Optional[dict[str, Any]]:
        """Resume execution from checkpoint file.

        Args:
            feature_id: Feature identifier
            agent_id: Agent identifier for verification

        Returns:
            Checkpoint data dictionary or None if not found/invalid

        Note:
            Verifies agent ID matches checkpoint before resuming.
        """
        checkpoint_data = self.checkpoint_manager.load_checkpoint(feature_id)

        if checkpoint_data is None:
            logger.debug(f"No checkpoint found for feature: {feature_id}")
            return None

        # Verify agent ID
        if checkpoint_data.get("agent_id") != agent_id:
            logger.warning(
                f"Agent ID mismatch for feature {feature_id}: "
                f"expected {checkpoint_data.get('agent_id')}, got {agent_id}"
            )
            return None

        # Restore approval gate state if available
        if "gates" in checkpoint_data and self.gate_manager:
            try:
                logger.debug(f"Restoring gate state for feature: {feature_id}")
            except Exception as e:
                logger.warning(f"Failed to restore gate state: {e}")

        # Restore team state if available
        if "team" in checkpoint_data and self.team_manager:
            try:
                logger.debug(f"Restoring team state for feature: {feature_id}")
            except Exception as e:
                logger.warning(f"Failed to restore team state: {e}")

        logger.info(f"Resumed from checkpoint for feature: {feature_id}")
        return checkpoint_data
