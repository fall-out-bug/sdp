"""Approval gate manager for @oneshot workflow."""

import logging

from sdp.unified.checkpoint.repository import CheckpointRepository
from sdp.unified.gates.errors import GateManagerError
from sdp.unified.gates.models import ApprovalGate, ApprovalStatus, GateType
from sdp.unified.gates.operations import GateOperations
from sdp.unified.gates.storage import GateStorage

logger = logging.getLogger(__name__)


class ApprovalGateManager:
    """Manager for approval gates in @oneshot workflow."""

    def __init__(self, checkpoint_repo: CheckpointRepository) -> None:
        """Initialize approval gate manager.

        Args:
            checkpoint_repo: Checkpoint repository for persistence
        """
        self.checkpoint_repo = checkpoint_repo

    def approve(
        self,
        feature: str,
        gate_type: GateType,
        approved_by: str,
        comments: str | None = None,
    ) -> None:
        """Approve an approval gate.

        Args:
            feature: Feature ID
            gate_type: Type of gate to approve
            approved_by: User ID who approved
            comments: Optional approval comments

        Raises:
            GateManagerError: If approval fails
        """
        try:
            GateOperations.approve_gate(
                self.checkpoint_repo, feature, gate_type, approved_by, comments
            )
        except GateManagerError:
            raise
        except Exception as e:
            logger.error(f"Failed to approve gate {gate_type.value} for {feature}: {e}")
            raise GateManagerError(f"Failed to approve gate: {e}") from e

    def reject(
        self,
        feature: str,
        gate_type: GateType,
        rejected_by: str,
        comments: str | None = None,
    ) -> None:
        """Reject an approval gate.

        Args:
            feature: Feature ID
            gate_type: Type of gate to reject
            rejected_by: User ID who rejected
            comments: Optional rejection comments

        Raises:
            GateManagerError: If rejection fails
        """
        try:
            GateOperations.reject_gate(
                self.checkpoint_repo, feature, gate_type, rejected_by, comments
            )
        except GateManagerError:
            raise
        except Exception as e:
            logger.error(f"Failed to reject gate {gate_type.value} for {feature}: {e}")
            raise GateManagerError(f"Failed to reject gate: {e}") from e

    def skip(
        self,
        feature: str,
        gate_type: GateType,
        reason: str | None = None,
    ) -> None:
        """Skip an approval gate.

        Args:
            feature: Feature ID
            gate_type: Type of gate to skip
            reason: Optional skip reason

        Raises:
            GateManagerError: If skip fails
        """
        try:
            GateOperations.skip_gate(
                self.checkpoint_repo, feature, gate_type, reason
            )
        except GateManagerError:
            raise
        except Exception as e:
            logger.error(f"Failed to skip gate {gate_type.value} for {feature}: {e}")
            raise GateManagerError(f"Failed to skip gate: {e}") from e

    def is_skipped(self, feature: str, gate_type: GateType) -> bool:
        """Check if a gate is skipped.

        Args:
            feature: Feature ID
            gate_type: Type of gate to check

        Returns:
            True if gate is skipped, False otherwise
        """
        try:
            checkpoint = self.checkpoint_repo.load_checkpoint(feature)
            if checkpoint is None:
                return False

            gates = GateStorage.extract_gates_from_checkpoint(checkpoint)
            gate = GateStorage.find_gate(gates, gate_type)

            return gate is not None and gate.status == ApprovalStatus.SKIPPED
        except Exception as e:
            logger.error(f"Failed to check skip status for {feature}: {e}")
            return False

    def get_gate_status(
        self, feature: str, gate_type: GateType
    ) -> ApprovalStatus:
        """Get the status of a gate.

        Args:
            feature: Feature ID
            gate_type: Type of gate to check

        Returns:
            Current status of the gate

        Raises:
            GateManagerError: If checkpoint not found
        """
        checkpoint = self.checkpoint_repo.load_checkpoint(feature)
        if checkpoint is None:
            raise GateManagerError(f"No checkpoint found for feature: {feature}")

        gates = GateStorage.extract_gates_from_checkpoint(checkpoint)
        gate = GateStorage.find_gate(gates, gate_type)

        if gate is None:
            return ApprovalStatus.PENDING

        return gate.status

    def get_all_gates(self, feature: str) -> list[ApprovalGate]:
        """Get all gates for a feature.

        Args:
            feature: Feature ID

        Returns:
            List of all gates (created if not exist)

        Raises:
            GateManagerError: If checkpoint not found
        """
        checkpoint = self.checkpoint_repo.load_checkpoint(feature)
        if checkpoint is None:
            raise GateManagerError(f"No checkpoint found for feature: {feature}")

        return GateStorage.extract_gates_from_checkpoint(checkpoint)
