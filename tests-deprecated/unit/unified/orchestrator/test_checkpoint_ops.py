"""
Tests for CheckpointOperations error handling and edge cases.

Tests error paths in save_checkpoint and resume_from_checkpoint.
"""

import pytest
import json
from pathlib import Path
from datetime import datetime
from unittest.mock import MagicMock, patch

from sdp.unified.checkpoint.repository import CheckpointRepository
from sdp.unified.checkpoint.schema import CheckpointStatus
from sdp.unified.orchestrator.checkpoint_ops import CheckpointOperations
from sdp.unified.orchestrator.checkpoint import CheckpointFileManager
from sdp.unified.gates.manager import ApprovalGateManager
from sdp.unified.team.manager import TeamManager


@pytest.fixture
def temp_checkpoint_dir(tmp_path: Path) -> Path:
    """Create temporary directory for checkpoint files."""
    checkpoint_dir = tmp_path / ".oneshot"
    checkpoint_dir.mkdir(parents=True, exist_ok=True)
    return checkpoint_dir


@pytest.fixture
def mock_repo() -> MagicMock:
    """Create mock CheckpointRepository."""
    repo = MagicMock(spec=CheckpointRepository)
    return repo


@pytest.fixture
def checkpoint_file_manager(temp_checkpoint_dir: Path) -> CheckpointFileManager:
    """Create CheckpointFileManager instance."""
    return CheckpointFileManager(str(temp_checkpoint_dir))


@pytest.fixture
def checkpoint_ops(mock_repo: MagicMock, checkpoint_file_manager: CheckpointFileManager) -> CheckpointOperations:
    """Create CheckpointOperations instance."""
    return CheckpointOperations(checkpoint_file_manager, mock_repo)


class TestCheckpointOperationsSaveCheckpoint:
    """Test save_checkpoint error handling."""

    def test_save_checkpoint_handles_gate_manager_exception(
        self, checkpoint_ops: CheckpointOperations, mock_repo: MagicMock
    ):
        """Should handle exception when loading gate state."""
        gate_manager = MagicMock(spec=ApprovalGateManager)
        checkpoint_ops.set_gate_manager(gate_manager)

        # Mock load_checkpoint to raise exception
        mock_repo.load_checkpoint.side_effect = Exception("Load failed")

        # Should not raise, just log warning
        checkpoint_ops.save_checkpoint(
            feature_id="F01",
            agent_id="agent-123",
            workstreams=["WS-001"],
            completed_workstreams=[],
            current_workstream="WS-001",
            status=CheckpointStatus.IN_PROGRESS,
        )

        # Checkpoint should still be saved
        checkpoint_file = checkpoint_ops.checkpoint_manager.base_path / "F01-checkpoint.json"
        assert checkpoint_file.exists()

    def test_save_checkpoint_handles_gate_manager_without_metrics(
        self, checkpoint_ops: CheckpointOperations, mock_repo: MagicMock
    ):
        """Should handle checkpoint without metrics."""
        gate_manager = MagicMock(spec=ApprovalGateManager)
        checkpoint_ops.set_gate_manager(gate_manager)

        # Mock checkpoint without metrics
        mock_checkpoint = MagicMock()
        mock_checkpoint.metrics = None
        mock_repo.load_checkpoint.return_value = mock_checkpoint

        checkpoint_ops.save_checkpoint(
            feature_id="F01",
            agent_id="agent-123",
            workstreams=["WS-001"],
            completed_workstreams=[],
            current_workstream="WS-001",
            status=CheckpointStatus.IN_PROGRESS,
        )

        # Should not include gates in checkpoint data
        checkpoint_file = checkpoint_ops.checkpoint_manager.base_path / "F01-checkpoint.json"
        with open(checkpoint_file) as f:
            data = json.load(f)

        assert "gates" not in data

    def test_save_checkpoint_handles_team_manager_exception(
        self, checkpoint_ops: CheckpointOperations, temp_checkpoint_dir: Path
    ):
        """Should handle exception when loading team state."""
        team_manager = MagicMock(spec=TeamManager)
        # Mock roles as a dict with values() method
        mock_roles = MagicMock()
        mock_roles.values.side_effect = Exception("Team load failed")
        team_manager.roles = mock_roles
        checkpoint_ops.set_team_manager(team_manager)

        # Should not raise, just log warning
        checkpoint_ops.save_checkpoint(
            feature_id="F01",
            agent_id="agent-123",
            workstreams=["WS-001"],
            completed_workstreams=[],
            current_workstream="WS-001",
            status=CheckpointStatus.IN_PROGRESS,
        )

        # Checkpoint should still be saved
        checkpoint_file = checkpoint_ops.checkpoint_manager.base_path / "F01-checkpoint.json"
        assert checkpoint_file.exists()

    def test_save_checkpoint_includes_gates_when_available(
        self, checkpoint_ops: CheckpointOperations, mock_repo: MagicMock
    ):
        """Should include gates when checkpoint has metrics."""
        gate_manager = MagicMock(spec=ApprovalGateManager)
        checkpoint_ops.set_gate_manager(gate_manager)

        # Mock checkpoint with gates in metrics
        mock_checkpoint = MagicMock()
        mock_checkpoint.metrics = {"gates": {"requirements": {"status": "approved"}}}
        mock_repo.load_checkpoint.return_value = mock_checkpoint

        checkpoint_ops.save_checkpoint(
            feature_id="F01",
            agent_id="agent-123",
            workstreams=["WS-001"],
            completed_workstreams=[],
            current_workstream="WS-001",
            status=CheckpointStatus.IN_PROGRESS,
        )

        checkpoint_file = checkpoint_ops.checkpoint_manager.base_path / "F01-checkpoint.json"
        with open(checkpoint_file) as f:
            data = json.load(f)

        assert "gates" in data
        assert data["gates"]["requirements"]["status"] == "approved"


class TestCheckpointOperationsResumeFromCheckpoint:
    """Test resume_from_checkpoint error handling."""

    def test_resume_from_checkpoint_handles_gate_restore_exception(
        self, checkpoint_ops: CheckpointOperations, temp_checkpoint_dir: Path
    ):
        """Should handle exception when restoring gate state."""
        # Save checkpoint with gates
        checkpoint_data = {
            "feature_id": "F01",
            "agent_id": "agent-123",
            "workstreams": ["WS-001"],
            "completed_workstreams": [],
            "current_workstream": "WS-001",
            "status": "in_progress",
            "gates": {"requirements": {"status": "approved"}},
        }

        checkpoint_file = temp_checkpoint_dir / "F01-checkpoint.json"
        with open(checkpoint_file, "w") as f:
            json.dump(checkpoint_data, f)

        gate_manager = MagicMock(spec=ApprovalGateManager)
        # Mock gate manager to raise exception
        gate_manager.restore_gates = MagicMock(side_effect=Exception("Restore failed"))
        checkpoint_ops.set_gate_manager(gate_manager)

        # Should not raise, just log warning
        resumed = checkpoint_ops.resume_from_checkpoint("F01", "agent-123")

        assert resumed is not None
        assert resumed["feature_id"] == "F01"

    def test_resume_from_checkpoint_handles_team_restore_exception(
        self, checkpoint_ops: CheckpointOperations, temp_checkpoint_dir: Path
    ):
        """Should handle exception when restoring team state."""
        # Save checkpoint with team
        checkpoint_data = {
            "feature_id": "F01",
            "agent_id": "agent-123",
            "workstreams": ["WS-001"],
            "completed_workstreams": [],
            "current_workstream": "WS-001",
            "status": "in_progress",
            "team": {"roles": [{"name": "orchestrator", "state": "active"}]},
        }

        checkpoint_file = temp_checkpoint_dir / "F01-checkpoint.json"
        with open(checkpoint_file, "w") as f:
            json.dump(checkpoint_data, f)

        team_manager = MagicMock(spec=TeamManager)
        # Mock team manager to raise exception
        team_manager.restore_roles = MagicMock(side_effect=Exception("Restore failed"))
        checkpoint_ops.set_team_manager(team_manager)

        # Should not raise, just log warning
        resumed = checkpoint_ops.resume_from_checkpoint("F01", "agent-123")

        assert resumed is not None
        assert resumed["feature_id"] == "F01"

    def test_resume_from_checkpoint_restores_gates_when_available(
        self, checkpoint_ops: CheckpointOperations, temp_checkpoint_dir: Path
    ):
        """Should restore gates when present in checkpoint."""
        checkpoint_data = {
            "feature_id": "F01",
            "agent_id": "agent-123",
            "workstreams": ["WS-001"],
            "completed_workstreams": [],
            "current_workstream": "WS-001",
            "status": "in_progress",
            "gates": {"requirements": {"status": "approved"}},
        }

        checkpoint_file = temp_checkpoint_dir / "F01-checkpoint.json"
        with open(checkpoint_file, "w") as f:
            json.dump(checkpoint_data, f)

        gate_manager = MagicMock(spec=ApprovalGateManager)
        checkpoint_ops.set_gate_manager(gate_manager)

        resumed = checkpoint_ops.resume_from_checkpoint("F01", "agent-123")

        assert resumed is not None
        assert "gates" in resumed

    def test_resume_from_checkpoint_restores_team_when_available(
        self, checkpoint_ops: CheckpointOperations, temp_checkpoint_dir: Path
    ):
        """Should restore team when present in checkpoint."""
        checkpoint_data = {
            "feature_id": "F01",
            "agent_id": "agent-123",
            "workstreams": ["WS-001"],
            "completed_workstreams": [],
            "current_workstream": "WS-001",
            "status": "in_progress",
            "team": {"roles": [{"name": "orchestrator", "state": "active"}]},
        }

        checkpoint_file = temp_checkpoint_dir / "F01-checkpoint.json"
        with open(checkpoint_file, "w") as f:
            json.dump(checkpoint_data, f)

        team_manager = MagicMock(spec=TeamManager)
        checkpoint_ops.set_team_manager(team_manager)

        resumed = checkpoint_ops.resume_from_checkpoint("F01", "agent-123")

        assert resumed is not None
        assert "team" in resumed

    def test_resume_from_checkpoint_handles_logger_exception_in_gate_restore(
        self, checkpoint_ops: CheckpointOperations, temp_checkpoint_dir: Path, monkeypatch
    ):
        """Should handle exception from logger.debug when restoring gates."""
        import logging

        checkpoint_data = {
            "feature_id": "F01",
            "agent_id": "agent-123",
            "workstreams": ["WS-001"],
            "completed_workstreams": [],
            "current_workstream": "WS-001",
            "status": "in_progress",
            "gates": {"requirements": {"status": "approved"}},
        }

        checkpoint_file = temp_checkpoint_dir / "F01-checkpoint.json"
        with open(checkpoint_file, "w") as f:
            json.dump(checkpoint_data, f)

        gate_manager = MagicMock(spec=ApprovalGateManager)
        checkpoint_ops.set_gate_manager(gate_manager)

        # Mock logger.debug to raise exception
        logger = logging.getLogger("sdp.unified.orchestrator.checkpoint_ops")
        original_debug = logger.debug

        def failing_debug(msg):
            if "Restoring gate state" in msg:
                raise Exception("Logger failed")
            return original_debug(msg)

        monkeypatch.setattr(logger, "debug", failing_debug)

        # Should not raise, just log warning
        resumed = checkpoint_ops.resume_from_checkpoint("F01", "agent-123")

        assert resumed is not None

    def test_resume_from_checkpoint_handles_logger_exception_in_team_restore(
        self, checkpoint_ops: CheckpointOperations, temp_checkpoint_dir: Path, monkeypatch
    ):
        """Should handle exception from logger.debug when restoring team."""
        import logging

        checkpoint_data = {
            "feature_id": "F01",
            "agent_id": "agent-123",
            "workstreams": ["WS-001"],
            "completed_workstreams": [],
            "current_workstream": "WS-001",
            "status": "in_progress",
            "team": {"roles": [{"name": "orchestrator", "state": "active"}]},
        }

        checkpoint_file = temp_checkpoint_dir / "F01-checkpoint.json"
        with open(checkpoint_file, "w") as f:
            json.dump(checkpoint_data, f)

        team_manager = MagicMock(spec=TeamManager)
        checkpoint_ops.set_team_manager(team_manager)

        # Mock logger.debug to raise exception
        logger = logging.getLogger("sdp.unified.orchestrator.checkpoint_ops")
        original_debug = logger.debug

        def failing_debug(msg):
            if "Restoring team state" in msg:
                raise Exception("Logger failed")
            return original_debug(msg)

        monkeypatch.setattr(logger, "debug", failing_debug)

        # Should not raise, just log warning
        resumed = checkpoint_ops.resume_from_checkpoint("F01", "agent-123")

        assert resumed is not None
