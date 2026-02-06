"""Checkpoint storage helpers for approval gates."""

import logging
from datetime import datetime

from sdp.unified.checkpoint.repository import CheckpointRepository
from sdp.unified.checkpoint.schema import Checkpoint, CheckpointDatabase
from sdp.unified.gates.errors import GateManagerError
from sdp.unified.gates.models import ApprovalGate, ApprovalStatus, GateType

logger = logging.getLogger(__name__)


class GateStorage:
    """Storage helpers for approval gates in checkpoints."""

    # Metric key for storing gates in checkpoint
    GATES_METRIC_KEY = "approval_gates"

    @staticmethod
    def extract_gates_from_checkpoint(checkpoint: Checkpoint) -> list[ApprovalGate]:
        """Extract gates from checkpoint metrics.

        Args:
            checkpoint: Checkpoint to extract from

        Returns:
            List of approval gates
        """
        gates_data_raw = checkpoint.metrics.get(
            GateStorage.GATES_METRIC_KEY, []
        )

        if not gates_data_raw:
            # Return default gates if none exist
            return [
                ApprovalGate(gate_type=gt, status=ApprovalStatus.PENDING)
                for gt in GateType
            ]

        # Deserialize gates from JSON
        gates = []
        gates_data = (
            list(gates_data_raw) if isinstance(gates_data_raw, list) else []
        )
        for gate_dict in gates_data:
            gate = ApprovalGate(
                gate_type=GateType(gate_dict["gate_type"]),
                status=ApprovalStatus(gate_dict["status"]),
                approved_by=gate_dict.get("approved_by"),
                approved_at=(
                    datetime.fromisoformat(gate_dict["approved_at"])
                    if gate_dict.get("approved_at")
                    else None
                ),
                comments=gate_dict.get("comments"),
            )
            gates.append(gate)

        # Ensure all gate types exist
        existing_types = {gate.gate_type for gate in gates}
        for gt in GateType:
            if gt not in existing_types:
                gates.append(
                    ApprovalGate(gate_type=gt, status=ApprovalStatus.PENDING)
                )

        return gates

    @staticmethod
    def find_gate(
        gates: list[ApprovalGate], gate_type: GateType
    ) -> ApprovalGate | None:
        """Find a gate by type.

        Args:
            gates: List of gates to search
            gate_type: Type of gate to find

        Returns:
            Gate if found, None otherwise
        """
        for gate in gates:
            if gate.gate_type == gate_type:
                return gate
        return None

    @staticmethod
    def save_gates_to_checkpoint(
        checkpoint: Checkpoint, gates: list[ApprovalGate]
    ) -> None:
        """Save gates to checkpoint metrics.

        Args:
            checkpoint: Checkpoint to update
            gates: Gates to save
        """
        gates_data = []
        for gate in gates:
            gate_dict = {
                "gate_type": gate.gate_type.value,
                "status": gate.status.value,
                "approved_by": gate.approved_by,
                "approved_at": (
                    gate.approved_at.isoformat() if gate.approved_at else None
                ),
                "comments": gate.comments,
            }
            gates_data.append(gate_dict)

        checkpoint.metrics[GateStorage.GATES_METRIC_KEY] = gates_data

    @staticmethod
    def get_checkpoint_id(
        checkpoint_repo: CheckpointRepository, feature: str
    ) -> int:
        """Get checkpoint ID from database.

        Args:
            checkpoint_repo: Checkpoint repository
            feature: Feature ID

        Returns:
            Checkpoint ID

        Raises:
            GateManagerError: If checkpoint not found
        """
        db = CheckpointDatabase(str(checkpoint_repo.db_path))

        # Query the checkpoint table to get the ID
        conn = db._get_connection()
        cursor = conn.cursor()
        cursor.execute(
            """
            SELECT id FROM checkpoints
            WHERE feature = ?
            ORDER BY id DESC
            LIMIT 1
        """,
            (feature,),
        )
        row = cursor.fetchone()

        if row is None:
            raise GateManagerError(f"Checkpoint ID not found for feature: {feature}")

        checkpoint_id = row[0]
        return int(checkpoint_id)

    @staticmethod
    def update_checkpoint_in_db(
        checkpoint_repo: CheckpointRepository, checkpoint_id: int, checkpoint: Checkpoint
    ) -> None:
        """Update checkpoint in database.

        Args:
            checkpoint_repo: Checkpoint repository
            checkpoint_id: Checkpoint ID
            checkpoint: Updated checkpoint data
        """
        db = CheckpointDatabase(str(checkpoint_repo.db_path))
        db.update_checkpoint(checkpoint_id, checkpoint)
