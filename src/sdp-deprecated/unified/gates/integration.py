"""Skip flag integration logic for ApprovalGateManager."""

import logging
from typing import Any, Optional

from sdp.unified.checkpoint.repository import CheckpointRepository
from sdp.unified.gates.models import ApprovalStatus, GateType
from sdp.unified.gates.parser import SkipFlagParser

logger = logging.getLogger(__name__)


class SkipFlagIntegration:
    """Handles skip flag integration with approval gates."""

    def __init__(
        self,
        checkpoint_repo: CheckpointRepository,
        skip_parser: Optional[SkipFlagParser] = None,
    ) -> None:
        """Initialize skip flag integration.

        Args:
            checkpoint_repo: Checkpoint repository for persistence
            skip_parser: Optional skip flag parser for auto-skipping gates
        """
        self.checkpoint_repo = checkpoint_repo
        self.skip_parser = skip_parser

    def request_approval(
        self,
        feature: str,
        gate_type: GateType,
        requestor: str,
        skip_method: Any,
    ) -> None:
        """Request approval for a gate, auto-skipping if flag is set.

        This method checks if a skip flag is set for the gate type.
        If the skip flag is present, the gate is automatically skipped.
        Otherwise, this method is a no-op (approval must be done manually).

        Args:
            feature: Feature ID
            gate_type: Type of gate to request approval for
            requestor: ID of the agent/user requesting approval
            skip_method: Method to call for skipping (manager.skip)
        """
        if self.skip_parser is None:
            logger.debug(
                f"No skip parser configured for {feature}, gate {gate_type.value}"
            )
            return

        if self.skip_parser.is_skip_required(gate_type):
            logger.info(
                f"Auto-skipping {gate_type.value} gate for {feature} "
                f"(flag detected, requested by {requestor})"
            )
            skip_method(
                feature=feature,
                gate_type=gate_type,
                reason=f"Auto-skipped via --skip-{gate_type.value} flag",
            )
        else:
            logger.debug(
                f"No skip flag for {gate_type.value} gate, "
                f"manual approval required for {feature}"
            )

    def auto_skip_gates(self, feature: str, skip_method: Any) -> None:
        """Automatically skip all gates based on command-line flags.

        This method checks all gate types and skips those that have
        corresponding skip flags set. Gates without skip flags are unchanged.

        Args:
            feature: Feature ID
            skip_method: Method to call for skipping (manager.skip)
        """
        if self.skip_parser is None:
            logger.debug(f"No skip parser configured for {feature}, skipping none")
            return

        for gate_type in GateType:
            if self.skip_parser.is_skip_required(gate_type):
                logger.info(f"Auto-skipping {gate_type.value} gate for {feature}")
                skip_method(
                    feature=feature,
                    gate_type=gate_type,
                    reason=f"Auto-skipped via --skip-{gate_type.value} flag",
                )

    def should_approve_skipped_gate(self, feature: str, gate_type: GateType) -> bool:
        """Check if an approve operation should proceed for a skipped gate.

        Args:
            feature: Feature ID
            gate_type: Type of gate to check

        Returns:
            False if gate is skipped (approve should be ignored), True otherwise
        """
        from sdp.unified.gates.storage import GateStorage

        checkpoint = self.checkpoint_repo.load_checkpoint(feature)
        if checkpoint is None:
            return True

        gates = GateStorage.extract_gates_from_checkpoint(checkpoint)
        gate = GateStorage.find_gate(gates, gate_type)

        # If gate is skipped, approve should not proceed
        if gate is not None and gate.status == ApprovalStatus.SKIPPED:
            logger.warning(
                f"Cannot approve {gate_type.value} gate for {feature}: "
                f"already skipped"
            )
            return False

        return True
