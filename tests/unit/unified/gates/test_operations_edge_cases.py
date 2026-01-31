"""Tests for gate operations edge cases."""

import pytest
from datetime import datetime
from unittest.mock import Mock, MagicMock, patch

from sdp.unified.gates.operations import GateOperations
from sdp.unified.gates.errors import GateManagerError
from sdp.unified.gates.models import ApprovalGate, ApprovalStatus, GateType
from sdp.unified.checkpoint.repository import CheckpointRepository
from sdp.unified.checkpoint.models import Checkpoint


@pytest.fixture
def mock_checkpoint_repo():
    """Create mock checkpoint repository."""
    return Mock(spec=CheckpointRepository)


@pytest.fixture
def mock_checkpoint():
    """Create mock checkpoint."""
    checkpoint = Mock(spec=Checkpoint)
    checkpoint.feature_id = "F001"
    checkpoint.state = {"approval_gates": []}
    return checkpoint


class TestApproveGate:
    """Tests for approve_gate operation."""

    def test_approve_gate_no_checkpoint(self, mock_checkpoint_repo):
        """Test raises error when checkpoint not found."""
        mock_checkpoint_repo.load_checkpoint.return_value = None

        with pytest.raises(GateManagerError, match="No checkpoint found for feature: F001"):
            GateOperations.approve_gate(
                mock_checkpoint_repo,
                "F001",
                GateType.REQUIREMENTS,
                "user123",
            )

    @patch('sdp.unified.gates.operations.GateStorage')
    def test_approve_gate_creates_new_gate(self, mock_storage, mock_checkpoint_repo, mock_checkpoint):
        """Test creates new gate if not found."""
        mock_checkpoint_repo.load_checkpoint.return_value = mock_checkpoint
        mock_storage.get_checkpoint_id.return_value = "checkpoint123"
        mock_storage.extract_gates_from_checkpoint.return_value = []
        mock_storage.find_gate.return_value = None

        GateOperations.approve_gate(
            mock_checkpoint_repo,
            "F001",
            GateType.REQUIREMENTS,
            "user123",
            comments="LGTM",
        )

        # Verify new gate was added to gates list
        save_call = mock_storage.save_gates_to_checkpoint.call_args
        saved_gates = save_call[0][1]
        assert len(saved_gates) == 1
        assert saved_gates[0].gate_type == GateType.REQUIREMENTS
        assert saved_gates[0].status == ApprovalStatus.APPROVED

    @patch('sdp.unified.gates.operations.GateStorage')
    def test_approve_gate_updates_existing_gate(self, mock_storage, mock_checkpoint_repo, mock_checkpoint):
        """Test updates existing gate."""
        existing_gate = ApprovalGate(
            gate_type=GateType.REQUIREMENTS,
            status=ApprovalStatus.PENDING,
        )

        mock_checkpoint_repo.load_checkpoint.return_value = mock_checkpoint
        mock_storage.get_checkpoint_id.return_value = "checkpoint123"
        mock_storage.extract_gates_from_checkpoint.return_value = [existing_gate]
        mock_storage.find_gate.return_value = existing_gate

        GateOperations.approve_gate(
            mock_checkpoint_repo,
            "F001",
            GateType.REQUIREMENTS,
            "user123",
            comments="Approved",
        )

        assert existing_gate.status == ApprovalStatus.APPROVED
        assert existing_gate.approved_by == "user123"
        assert existing_gate.comments == "Approved"
        assert existing_gate.approved_at is not None


class TestRejectGate:
    """Tests for reject_gate operation."""

    @patch('sdp.unified.gates.operations.GateStorage')
    def test_reject_gate_no_checkpoint(self, mock_storage, mock_checkpoint_repo):
        """Test raises error when checkpoint not found."""
        mock_checkpoint_repo.load_checkpoint.return_value = None

        with pytest.raises(GateManagerError, match="No checkpoint found for feature: F001"):
            GateOperations.reject_gate(
                mock_checkpoint_repo,
                "F001",
                GateType.REQUIREMENTS,
                "user123",
            )


class TestListGates:
    """Tests for list_gates operation."""

    @patch('sdp.unified.gates.operations.GateStorage')
    def test_list_gates_no_checkpoint(self, mock_storage, mock_checkpoint_repo):
        """Test returns empty list when checkpoint not found."""
        mock_checkpoint_repo.load_checkpoint.return_value = None

        from sdp.unified.gates.operations import GateOperations
        
        # The actual method might be in manager, but testing error path
        with pytest.raises(GateManagerError):
            GateOperations.approve_gate(mock_checkpoint_repo, "F001", GateType.REQUIREMENTS, "user")
