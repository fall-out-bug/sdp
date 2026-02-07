"""Approval gate manager for @oneshot workflow."""

import logging
from typing import Optional

from sdp.unified.checkpoint.repository import CheckpointRepository
from sdp.unified.gates.errors import GateManagerError
from sdp.unified.gates.integration import SkipFlagIntegration
from sdp.unified.gates.models import ApprovalGate, ApprovalStatus, GateType
from sdp.unified.gates.operations import GateOperations
from sdp.unified.gates.parser import SkipFlagParser
from sdp.unified.gates.storage import GateStorage

logger = logging.getLogger(__name__)


class ApprovalGateManager:
    """Manager for approval gates in @oneshot workflow."""

    def __init__(
        self,
        checkpoint_repo: CheckpointRepository,
        skip_parser: Optional[SkipFlagParser] = None,
    ) -> None:
        """Initialize approval gate manager.

        Args:
            checkpoint_repo: Checkpoint repository for persistence
            skip_parser: Optional skip flag parser for auto-skipping gates
        """
        self.checkpoint_repo = checkpoint_repo
        self.skip_parser = skip_parser
        self._integration = SkipFlagIntegration(checkpoint_repo, skip_parser)

    def approve(
        self,
        feature: str,
        gate_type: GateType,
        approved_by: str,
        comments: str | None = None,
    ) -> None:
        """Approve an approval gate. Skips if already skipped."""
        if not self._integration.should_approve_skipped_gate(feature, gate_type):
            return

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
        """Reject an approval gate."""
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
        """Skip an approval gate."""
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
        """Check if a gate is skipped."""
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
        """Get the status of a gate."""
        checkpoint = self.checkpoint_repo.load_checkpoint(feature)
        if checkpoint is None:
            raise GateManagerError(f"No checkpoint found for feature: {feature}")

        gates = GateStorage.extract_gates_from_checkpoint(checkpoint)
        gate = GateStorage.find_gate(gates, gate_type)

        if gate is None:
            return ApprovalStatus.PENDING

        return gate.status

    def get_all_gates(self, feature: str) -> list[ApprovalGate]:
        """Get all gates for a feature."""
        checkpoint = self.checkpoint_repo.load_checkpoint(feature)
        if checkpoint is None:
            raise GateManagerError(f"No checkpoint found for feature: {feature}")

        return GateStorage.extract_gates_from_checkpoint(checkpoint)

    def request_approval(
        self, feature: str, gate_type: GateType, requestor: str
    ) -> None:
        """Request approval for a gate, auto-skipping if flag is set."""
        self._integration.request_approval(feature, gate_type, requestor, self.skip)

    def auto_skip_gates(self, feature: str) -> None:
        """Automatically skip all gates based on command-line flags."""
        self._integration.auto_skip_gates(feature, self.skip)
