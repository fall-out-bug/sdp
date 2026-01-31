"""Tests for approval gates."""

import pytest
from datetime import datetime

from sdp.unified.gates.models import ApprovalGate, ApprovalStatus, GateType
from sdp.unified.gates.manager import ApprovalGateManager
from sdp.unified.gates.storage import GateStorage
from sdp.unified.checkpoint.repository import CheckpointRepository
from sdp.unified.checkpoint.schema import Checkpoint, CheckpointStatus


@pytest.fixture
def temp_db_path(tmp_path):
    """Create temporary database path."""
    return str(tmp_path / "test.db")


@pytest.fixture
def checkpoint_repo(temp_db_path):
    """Create checkpoint repository for testing."""
    repo = CheckpointRepository(temp_db_path)
    repo.initialize()
    yield repo
    repo.close()


@pytest.fixture
def sample_checkpoint():
    """Create sample checkpoint for testing."""
    return Checkpoint(
        feature="F01",
        agent_id="agent-001",
        status=CheckpointStatus.IN_PROGRESS,
        completed_ws=[],
        execution_order=["WS-001", "WS-002"],
        started_at=datetime.now(),
        current_ws="WS-001",
    )


class TestApprovalGate:
    """Tests for ApprovalGate data model."""

    def test_create_gate_with_required_fields(self):
        """Test creating a gate with required fields."""
        gate = ApprovalGate(
            gate_type=GateType.REQUIREMENTS,
            status=ApprovalStatus.PENDING,
        )

        assert gate.gate_type == GateType.REQUIREMENTS
        assert gate.status == ApprovalStatus.PENDING
        assert gate.approved_by is None
        assert gate.approved_at is None
        assert gate.comments is None

    def test_create_gate_with_all_fields(self):
        """Test creating a gate with all fields."""
        now = datetime.now()
        gate = ApprovalGate(
            gate_type=GateType.ARCHITECTURE,
            status=ApprovalStatus.APPROVED,
            approved_by="user123",
            approved_at=now,
            comments="Looks good",
        )

        assert gate.gate_type == GateType.ARCHITECTURE
        assert gate.status == ApprovalStatus.APPROVED
        assert gate.approved_by == "user123"
        assert gate.approved_at == now
        assert gate.comments == "Looks good"

    def test_gate_type_enum_values(self):
        """Test that GateType enum has correct values."""
        assert GateType.REQUIREMENTS.value == "requirements"
        assert GateType.ARCHITECTURE.value == "architecture"
        assert GateType.UAT.value == "uat"

    def test_approval_status_enum_values(self):
        """Test that ApprovalStatus enum has correct values."""
        assert ApprovalStatus.PENDING.value == "pending"
        assert ApprovalStatus.APPROVED.value == "approved"
        assert ApprovalStatus.REJECTED.value == "rejected"
        assert ApprovalStatus.SKIPPED.value == "skipped"


class TestApprovalGateManager:
    """Tests for ApprovalGateManager."""

    def test_create_manager_with_checkpoint_repo(self, checkpoint_repo):
        """Test creating manager with checkpoint repository."""
        manager = ApprovalGateManager(checkpoint_repo)

        assert manager.checkpoint_repo == checkpoint_repo

    def test_approve_requirements_gate(
        self, checkpoint_repo, sample_checkpoint
    ):
        """Test approving the requirements gate."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        manager = ApprovalGateManager(checkpoint_repo)

        manager.approve(
            feature="F01",
            gate_type=GateType.REQUIREMENTS,
            approved_by="user123",
            comments="Requirements approved",
        )

        # Reload checkpoint and verify
        updated_checkpoint = checkpoint_repo.load_checkpoint("F01")
        assert updated_checkpoint is not None

        gates = GateStorage.extract_gates_from_checkpoint(
            updated_checkpoint
        )
        requirements_gate = next(
            (g for g in gates if g.gate_type == GateType.REQUIREMENTS),
            None,
        )

        assert requirements_gate is not None
        assert requirements_gate.status == ApprovalStatus.APPROVED
        assert requirements_gate.approved_by == "user123"
        assert requirements_gate.comments == "Requirements approved"
        assert requirements_gate.approved_at is not None

    def test_approve_architecture_gate(
        self, checkpoint_repo, sample_checkpoint
    ):
        """Test approving the architecture gate."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        manager = ApprovalGateManager(checkpoint_repo)

        manager.approve(
            feature="F01",
            gate_type=GateType.ARCHITECTURE,
            approved_by="architect",
            comments="Design approved",
        )

        updated_checkpoint = checkpoint_repo.load_checkpoint("F01")
        assert updated_checkpoint is not None

        gates = GateStorage.extract_gates_from_checkpoint(
            updated_checkpoint
        )
        architecture_gate = next(
            (g for g in gates if g.gate_type == GateType.ARCHITECTURE),
            None,
        )

        assert architecture_gate is not None
        assert architecture_gate.status == ApprovalStatus.APPROVED
        assert architecture_gate.approved_by == "architect"

    def test_approve_uat_gate(self, checkpoint_repo, sample_checkpoint):
        """Test approving the UAT gate."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        manager = ApprovalGateManager(checkpoint_repo)

        manager.approve(
            feature="F01",
            gate_type=GateType.UAT,
            approved_by="tester",
            comments="UAT passed",
        )

        updated_checkpoint = checkpoint_repo.load_checkpoint("F01")
        gates = GateStorage.extract_gates_from_checkpoint(
            updated_checkpoint
        )
        uat_gate = next(
            (g for g in gates if g.gate_type == GateType.UAT),
            None,
        )

        assert uat_gate is not None
        assert uat_gate.status == ApprovalStatus.APPROVED
        assert uat_gate.approved_by == "tester"

    def test_reject_requirements_gate(
        self, checkpoint_repo, sample_checkpoint
    ):
        """Test rejecting the requirements gate."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        manager = ApprovalGateManager(checkpoint_repo)

        manager.reject(
            feature="F01",
            gate_type=GateType.REQUIREMENTS,
            rejected_by="user123",
            comments="Incomplete requirements",
        )

        updated_checkpoint = checkpoint_repo.load_checkpoint("F01")
        gates = GateStorage.extract_gates_from_checkpoint(
            updated_checkpoint
        )
        requirements_gate = next(
            (g for g in gates if g.gate_type == GateType.REQUIREMENTS),
            None,
        )

        assert requirements_gate is not None
        assert requirements_gate.status == ApprovalStatus.REJECTED
        assert requirements_gate.approved_by == "user123"
        assert requirements_gate.comments == "Incomplete requirements"

    def test_reject_architecture_gate(
        self, checkpoint_repo, sample_checkpoint
    ):
        """Test rejecting the architecture gate."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        manager = ApprovalGateManager(checkpoint_repo)

        manager.reject(
            feature="F01",
            gate_type=GateType.ARCHITECTURE,
            rejected_by="architect",
            comments="Design needs revision",
        )

        updated_checkpoint = checkpoint_repo.load_checkpoint("F01")
        gates = GateStorage.extract_gates_from_checkpoint(
            updated_checkpoint
        )
        architecture_gate = next(
            (g for g in gates if g.gate_type == GateType.ARCHITECTURE),
            None,
        )

        assert architecture_gate is not None
        assert architecture_gate.status == ApprovalStatus.REJECTED
        assert architecture_gate.comments == "Design needs revision"

    def test_is_skipped_returns_false_when_not_skipped(
        self, checkpoint_repo, sample_checkpoint
    ):
        """Test that is_skipped returns False when gate is not skipped."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        manager = ApprovalGateManager(checkpoint_repo)

        is_skipped = manager.is_skipped(
            feature="F01",
            gate_type=GateType.REQUIREMENTS,
        )

        assert is_skipped is False

    def test_is_skipped_returns_true_when_skipped(
        self, checkpoint_repo, sample_checkpoint
    ):
        """Test that is_skipped returns True when gate is skipped."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        manager = ApprovalGateManager(checkpoint_repo)

        # Skip the gate
        manager.skip(
            feature="F01",
            gate_type=GateType.REQUIREMENTS,
            reason="Skipping for now",
        )

        is_skipped = manager.is_skipped(
            feature="F01",
            gate_type=GateType.REQUIREMENTS,
        )

        assert is_skipped is True

    def test_get_gate_status(
        self, checkpoint_repo, sample_checkpoint
    ):
        """Test getting gate status."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        manager = ApprovalGateManager(checkpoint_repo)

        # Initially pending
        status = manager.get_gate_status(
            feature="F01",
            gate_type=GateType.REQUIREMENTS,
        )
        assert status == ApprovalStatus.PENDING

        # After approval
        manager.approve(
            feature="F01",
            gate_type=GateType.REQUIREMENTS,
            approved_by="user123",
        )

        status = manager.get_gate_status(
            feature="F01",
            gate_type=GateType.REQUIREMENTS,
        )
        assert status == ApprovalStatus.APPROVED

    def test_get_all_gates(
        self, checkpoint_repo, sample_checkpoint
    ):
        """Test getting all gates for a feature."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        manager = ApprovalGateManager(checkpoint_repo)

        gates = manager.get_all_gates(feature="F01")

        assert len(gates) == 3
        gate_types = {gate.gate_type for gate in gates}
        assert gate_types == {
            GateType.REQUIREMENTS,
            GateType.ARCHITECTURE,
            GateType.UAT,
        }

    def test_approve_nonexistent_feature_raises_error(
        self, checkpoint_repo
    ):
        """Test that approving a nonexistent feature raises an error."""
        manager = ApprovalGateManager(checkpoint_repo)

        with pytest.raises(Exception):  # Will be specific exception later
            manager.approve(
                feature="NONEXISTENT",
                gate_type=GateType.REQUIREMENTS,
                approved_by="user123",
            )

    def test_multiple_gates_can_be_approved(
        self, checkpoint_repo, sample_checkpoint
    ):
        """Test that multiple gates can be approved independently."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        manager = ApprovalGateManager(checkpoint_repo)

        # Approve requirements
        manager.approve(
            feature="F01",
            gate_type=GateType.REQUIREMENTS,
            approved_by="product_manager",
        )

        # Approve architecture
        manager.approve(
            feature="F01",
            gate_type=GateType.ARCHITECTURE,
            approved_by="architect",
        )

        # Approve UAT
        manager.approve(
            feature="F01",
            gate_type=GateType.UAT,
            approved_by="tester",
        )

        gates = manager.get_all_gates(feature="F01")

        for gate in gates:
            assert gate.status == ApprovalStatus.APPROVED

    def test_gate_decisions_persist_across_checkpoint_loads(
        self, checkpoint_repo, sample_checkpoint
    ):
        """Test that gate decisions persist across checkpoint loads."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        manager = ApprovalGateManager(checkpoint_repo)

        # Approve gate
        manager.approve(
            feature="F01",
            gate_type=GateType.REQUIREMENTS,
            approved_by="user123",
            comments="First approval",
        )

        # Create new manager instance (simulating restart)
        new_manager = ApprovalGateManager(checkpoint_repo)

        status = new_manager.get_gate_status(
            feature="F01",
            gate_type=GateType.REQUIREMENTS,
        )

        assert status == ApprovalStatus.APPROVED
