"""Gate operation implementations."""

import logging
from datetime import datetime

from sdp.unified.checkpoint.repository import CheckpointRepository
from sdp.unified.gates.errors import GateManagerError
from sdp.unified.gates.models import ApprovalGate, ApprovalStatus, GateType
from sdp.unified.gates.storage import GateStorage

logger = logging.getLogger(__name__)


class GateOperations:
    """Operations for approval gates."""

    @staticmethod
    def approve_gate(
        checkpoint_repo: CheckpointRepository,
        feature: str,
        gate_type: GateType,
        approved_by: str,
        comments: str | None = None,
    ) -> None:
        """Approve an approval gate.

        Args:
            checkpoint_repo: Checkpoint repository
            feature: Feature ID
            gate_type: Type of gate to approve
            approved_by: User ID who approved
            comments: Optional approval comments

        Raises:
            GateManagerError: If approval fails
        """
        checkpoint = checkpoint_repo.load_checkpoint(feature)
        if checkpoint is None:
            raise GateManagerError(f"No checkpoint found for feature: {feature}")

        # Get checkpoint ID from database
        checkpoint_id = GateStorage.get_checkpoint_id(checkpoint_repo, feature)

        # Extract existing gates
        gates = GateStorage.extract_gates_from_checkpoint(checkpoint)

        # Find and update the target gate
        gate = GateStorage.find_gate(gates, gate_type)
        if gate is None:
            gate = ApprovalGate(
                gate_type=gate_type,
                status=ApprovalStatus.PENDING,
            )
            gates.append(gate)

        gate.status = ApprovalStatus.APPROVED
        gate.approved_by = approved_by
        gate.approved_at = datetime.now()
        gate.comments = comments

        # Save back to checkpoint
        GateStorage.save_gates_to_checkpoint(checkpoint, gates)

        # Persist to database
        GateStorage.update_checkpoint_in_db(
            checkpoint_repo, checkpoint_id, checkpoint
        )

        logger.info(f"Gate approved: {feature} - {gate_type.value} by {approved_by}")

    @staticmethod
    def reject_gate(
        checkpoint_repo: CheckpointRepository,
        feature: str,
        gate_type: GateType,
        rejected_by: str,
        comments: str | None = None,
    ) -> None:
        """Reject an approval gate.

        Args:
            checkpoint_repo: Checkpoint repository
            feature: Feature ID
            gate_type: Type of gate to reject
            rejected_by: User ID who rejected
            comments: Optional rejection comments

        Raises:
            GateManagerError: If rejection fails
        """
        checkpoint = checkpoint_repo.load_checkpoint(feature)
        if checkpoint is None:
            raise GateManagerError(f"No checkpoint found for feature: {feature}")

        # Get checkpoint ID from database
        checkpoint_id = GateStorage.get_checkpoint_id(checkpoint_repo, feature)

        # Extract existing gates
        gates = GateStorage.extract_gates_from_checkpoint(checkpoint)

        # Find and update the target gate
        gate = GateStorage.find_gate(gates, gate_type)
        if gate is None:
            gate = ApprovalGate(
                gate_type=gate_type,
                status=ApprovalStatus.PENDING,
            )
            gates.append(gate)

        gate.status = ApprovalStatus.REJECTED
        gate.approved_by = rejected_by
        gate.approved_at = datetime.now()
        gate.comments = comments

        # Save back to checkpoint
        GateStorage.save_gates_to_checkpoint(checkpoint, gates)

        # Persist to database
        GateStorage.update_checkpoint_in_db(
            checkpoint_repo, checkpoint_id, checkpoint
        )

        logger.info(f"Gate rejected: {feature} - {gate_type.value} by {rejected_by}")

    @staticmethod
    def skip_gate(
        checkpoint_repo: CheckpointRepository,
        feature: str,
        gate_type: GateType,
        reason: str | None = None,
    ) -> None:
        """Skip an approval gate.

        Args:
            checkpoint_repo: Checkpoint repository
            feature: Feature ID
            gate_type: Type of gate to skip
            reason: Optional skip reason

        Raises:
            GateManagerError: If skip fails
        """
        checkpoint = checkpoint_repo.load_checkpoint(feature)
        if checkpoint is None:
            raise GateManagerError(f"No checkpoint found for feature: {feature}")

        # Get checkpoint ID from database
        checkpoint_id = GateStorage.get_checkpoint_id(checkpoint_repo, feature)

        # Extract existing gates
        gates = GateStorage.extract_gates_from_checkpoint(checkpoint)

        # Find and update the target gate
        gate = GateStorage.find_gate(gates, gate_type)
        if gate is None:
            gate = ApprovalGate(
                gate_type=gate_type,
                status=ApprovalStatus.PENDING,
            )
            gates.append(gate)

        gate.status = ApprovalStatus.SKIPPED
        gate.approved_by = None
        gate.approved_at = datetime.now()
        gate.comments = reason

        # Save back to checkpoint
        GateStorage.save_gates_to_checkpoint(checkpoint, gates)

        # Persist to database
        GateStorage.update_checkpoint_in_db(
            checkpoint_repo, checkpoint_id, checkpoint
        )

        logger.info(f"Gate skipped: {feature} - {gate_type.value}")
