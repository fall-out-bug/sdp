"""Tests for checkpoint save/resume functionality in OrchestratorAgent.

This module tests the checkpoint save/resume logic for the @oneshot workflow,
including file-based checkpoint storage, agent ID verification, and integration
with ApprovalGateManager and TeamManager.
"""

import json
import pytest
from datetime import datetime
from pathlib import Path
from unittest.mock import MagicMock, Mock

from sdp.unified.checkpoint.repository import CheckpointRepository
from sdp.unified.checkpoint.schema import Checkpoint, CheckpointStatus
from sdp.unified.gates.manager import ApprovalGateManager
from sdp.unified.gates.parser import SkipFlagParser
from sdp.unified.orchestrator.agent import OrchestratorAgent
from sdp.unified.orchestrator.agent_extension import AgentCheckpointExtension
from sdp.unified.orchestrator.checkpoint import CheckpointFileManager
from sdp.unified.team.manager import TeamManager
from sdp.unified.team.models import Role, RoleState


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
    repo.initialize = MagicMock()
    return repo


@pytest.fixture
def checkpoint_file_manager(temp_checkpoint_dir: Path) -> CheckpointFileManager:
    """Create CheckpointFileManager instance."""
    return CheckpointFileManager(str(temp_checkpoint_dir))


@pytest.fixture
def checkpoint_extension(mock_repo: MagicMock, temp_checkpoint_dir: Path) -> AgentCheckpointExtension:
    """Create AgentCheckpointExtension with checkpoint file manager."""
    from sdp.unified.orchestrator.agent_extension import AgentCheckpointExtension

    extension = AgentCheckpointExtension(mock_repo)
    checkpoint_manager = CheckpointFileManager(str(temp_checkpoint_dir))
    extension.set_checkpoint_manager(checkpoint_manager)
    return extension


class TestCheckpointFileManager:
    """Tests for CheckpointFileManager class."""

    def test_init_creates_directory(self, tmp_path: Path) -> None:
        """Test that initialization creates checkpoint directory."""
        checkpoint_dir = tmp_path / ".oneshot"
        manager = CheckpointFileManager(str(checkpoint_dir))

        assert checkpoint_dir.exists()
        assert manager.base_path == checkpoint_dir

    def test_save_checkpoint_creates_file(
        self, checkpoint_file_manager: CheckpointFileManager, temp_checkpoint_dir: Path
    ) -> None:
        """Test saving checkpoint creates JSON file."""
        checkpoint_data = {
            "feature_id": "F01",
            "agent_id": "agent-123",
            "workstreams": ["WS-001", "WS-002", "WS-003"],
            "completed_workstreams": ["WS-001"],
            "current_workstream": "WS-002",
            "status": "in_progress",
            "started_at": "2026-01-28T10:00:00",
        }

        checkpoint_file_manager.save_checkpoint("F01", checkpoint_data)

        checkpoint_file = temp_checkpoint_dir / "F01-checkpoint.json"
        assert checkpoint_file.exists()

        with open(checkpoint_file) as f:
            loaded_data = json.load(f)

        assert loaded_data == checkpoint_data

    def test_save_checkpoint_overwrites_existing(
        self, checkpoint_file_manager: CheckpointFileManager, temp_checkpoint_dir: Path
    ) -> None:
        """Test that saving checkpoint overwrites existing file."""
        initial_data = {"feature_id": "F01", "agent_id": "agent-123", "status": "in_progress"}
        updated_data = {"feature_id": "F01", "agent_id": "agent-123", "status": "completed"}

        checkpoint_file_manager.save_checkpoint("F01", initial_data)
        checkpoint_file_manager.save_checkpoint("F01", updated_data)

        checkpoint_file = temp_checkpoint_dir / "F01-checkpoint.json"
        with open(checkpoint_file) as f:
            loaded_data = json.load(f)

        assert loaded_data["status"] == "completed"

    def test_load_checkpoint_returns_data(
        self, checkpoint_file_manager: CheckpointFileManager, temp_checkpoint_dir: Path
    ) -> None:
        """Test loading checkpoint returns correct data."""
        checkpoint_data = {
            "feature_id": "F01",
            "agent_id": "agent-123",
            "workstreams": ["WS-001", "WS-002"],
            "completed_workstreams": ["WS-001"],
            "current_workstream": "WS-002",
            "status": "in_progress",
            "started_at": "2026-01-28T10:00:00",
        }

        checkpoint_file = temp_checkpoint_dir / "F01-checkpoint.json"
        with open(checkpoint_file, "w") as f:
            json.dump(checkpoint_data, f)

        loaded_data = checkpoint_file_manager.load_checkpoint("F01")

        assert loaded_data == checkpoint_data

    def test_load_checkpoint_returns_none_if_not_exists(
        self, checkpoint_file_manager: CheckpointFileManager
    ) -> None:
        """Test loading nonexistent checkpoint returns None."""
        loaded_data = checkpoint_file_manager.load_checkpoint("NONEXISTENT")
        assert loaded_data is None

    def test_delete_checkpoint_removes_file(
        self, checkpoint_file_manager: CheckpointFileManager, temp_checkpoint_dir: Path
    ) -> None:
        """Test deleting checkpoint removes file."""
        checkpoint_data = {"feature_id": "F01", "status": "in_progress"}
        checkpoint_file_manager.save_checkpoint("F01", checkpoint_data)

        checkpoint_file_manager.delete_checkpoint("F01")

        checkpoint_file = temp_checkpoint_dir / "F01-checkpoint.json"
        assert not checkpoint_file.exists()

    def test_delete_checkpoint_handles_nonexistent(
        self, checkpoint_file_manager: CheckpointFileManager
    ) -> None:
        """Test deleting nonexistent checkpoint doesn't raise error."""
        # Should not raise exception
        checkpoint_file_manager.delete_checkpoint("NONEXISTENT")

    def test_checkpoint_exists(self, checkpoint_file_manager: CheckpointFileManager) -> None:
        """Test checking if checkpoint exists."""
        assert not checkpoint_file_manager.checkpoint_exists("F01")

        checkpoint_data = {"feature_id": "F01", "status": "in_progress"}
        checkpoint_file_manager.save_checkpoint("F01", checkpoint_data)

        assert checkpoint_file_manager.checkpoint_exists("F01")


class TestOrchestratorAgentSaveCheckpoint:
    """Tests for AgentCheckpointExtension.save_checkpoint() method."""

    def test_save_checkpoint_creates_file(
        self, checkpoint_extension: AgentCheckpointExtension, temp_checkpoint_dir: Path
    ) -> None:
        """Test save_checkpoint() creates checkpoint file."""
        feature_id = "F01"
        agent_id = "agent-123"
        workstreams = ["WS-001", "WS-002", "WS-003"]
        completed_ws = ["WS-001"]
        current_ws = "WS-002"
        status = CheckpointStatus.IN_PROGRESS

        checkpoint_extension.save_checkpoint(
            feature_id=feature_id,
            agent_id=agent_id,
            workstreams=workstreams,
            completed_workstreams=completed_ws,
            current_workstream=current_ws,
            status=status,
        )

        checkpoint_file = temp_checkpoint_dir / "F01-checkpoint.json"
        assert checkpoint_file.exists()

        with open(checkpoint_file) as f:
            data = json.load(f)

        assert data["feature_id"] == feature_id
        assert data["agent_id"] == agent_id
        assert data["workstreams"] == workstreams
        assert data["completed_workstreams"] == completed_ws
        assert data["current_workstream"] == current_ws
        assert data["status"] == status.value

    def test_save_checkpoint_includes_metadata(
        self, checkpoint_extension: AgentCheckpointExtension, temp_checkpoint_dir: Path
    ) -> None:
        """Test save_checkpoint() includes timestamps and metadata."""
        checkpoint_extension.save_checkpoint(
            feature_id="F01",
            agent_id="agent-123",
            workstreams=["WS-001"],
            completed_workstreams=[],
            current_workstream="WS-001",
            status=CheckpointStatus.IN_PROGRESS,
        )

        checkpoint_file = temp_checkpoint_dir / "F01-checkpoint.json"
        with open(checkpoint_file) as f:
            data = json.load(f)

        assert "started_at" in data
        assert "last_updated" in data
        assert isinstance(data["started_at"], str)
        assert isinstance(data["last_updated"], str)

    def test_save_checkpoint_includes_gates(
        self, checkpoint_extension: AgentCheckpointExtension, temp_checkpoint_dir: Path
    ) -> None:
        """Test save_checkpoint() includes approval gate status."""
        gate_manager = ApprovalGateManager(checkpoint_extension.repo)
        checkpoint_extension.set_gate_manager(gate_manager)

        # Mock checkpoint repo to return gates
        mock_checkpoint = Checkpoint(
            feature="F01",
            agent_id="agent-123",
            status=CheckpointStatus.IN_PROGRESS,
            completed_ws=[],
            execution_order=["WS-001"],
            started_at=datetime.now(),
            metrics={"gates": {"requirements": {"status": "approved"}}},
        )
        checkpoint_extension.repo.load_checkpoint.return_value = mock_checkpoint

        checkpoint_extension.save_checkpoint(
            feature_id="F01",
            agent_id="agent-123",
            workstreams=["WS-001"],
            completed_workstreams=[],
            current_workstream="WS-001",
            status=CheckpointStatus.IN_PROGRESS,
        )

        checkpoint_file = temp_checkpoint_dir / "F01-checkpoint.json"
        with open(checkpoint_file) as f:
            data = json.load(f)

        assert "gates" in data

    def test_save_checkpoint_includes_team_state(
        self, checkpoint_extension: AgentCheckpointExtension, temp_checkpoint_dir: Path
    ) -> None:
        """Test save_checkpoint() includes team configuration."""
        # Create team manager
        team_manager = TeamManager("F01", base_path=temp_checkpoint_dir)
        team_manager.register_role(
            Role(
                name="orchestrator",
                description="Orchestration agent",
                state=RoleState.ACTIVE,
                skill_file="/path/to/skill.md",
            )
        )

        checkpoint_extension.set_team_manager(team_manager)

        checkpoint_extension.save_checkpoint(
            feature_id="F01",
            agent_id="agent-123",
            workstreams=["WS-001"],
            completed_workstreams=[],
            current_workstream="WS-001",
            status=CheckpointStatus.IN_PROGRESS,
        )

        checkpoint_file = temp_checkpoint_dir / "F01-checkpoint.json"
        with open(checkpoint_file) as f:
            data = json.load(f)

        assert "team" in data
        assert "roles" in data["team"]
        assert len(data["team"]["roles"]) == 1


class TestOrchestratorAgentResumeCheckpoint:
    """Tests for OrchestratorAgent.resume_from_checkpoint() method."""

    def test_resume_from_checkpoint_returns_data(
        self, checkpoint_extension: AgentCheckpointExtension, temp_checkpoint_dir: Path
    ) -> None:
        """Test resume_from_checkpoint() returns checkpoint data."""
        # Save checkpoint first
        checkpoint_data = {
            "feature_id": "F01",
            "agent_id": "agent-123",
            "workstreams": ["WS-001", "WS-002", "WS-003"],
            "completed_workstreams": ["WS-001"],
            "current_workstream": "WS-002",
            "status": "in_progress",
            "started_at": "2026-01-28T10:00:00",
            "last_updated": "2026-01-28T11:00:00",
        }

        checkpoint_file = temp_checkpoint_dir / "F01-checkpoint.json"
        with open(checkpoint_file, "w") as f:
            json.dump(checkpoint_data, f)

        # Resume checkpoint
        resumed_data = checkpoint_extension.resume_from_checkpoint(
            feature_id="F01", agent_id="agent-123"
        )

        assert resumed_data is not None
        assert resumed_data["feature_id"] == "F01"
        assert resumed_data["agent_id"] == "agent-123"
        assert resumed_data["completed_workstreams"] == ["WS-001"]
        assert resumed_data["current_workstream"] == "WS-002"

    def test_resume_from_checkpoint_verifies_agent_id(
        self, checkpoint_extension: AgentCheckpointExtension, temp_checkpoint_dir: Path
    ) -> None:
        """Test resume_from_checkpoint() verifies agent ID."""
        checkpoint_data = {
            "feature_id": "F01",
            "agent_id": "agent-123",
            "workstreams": ["WS-001"],
            "completed_workstreams": [],
            "current_workstream": "WS-001",
            "status": "in_progress",
            "started_at": "2026-01-28T10:00:00",
        }

        checkpoint_file = temp_checkpoint_dir / "F01-checkpoint.json"
        with open(checkpoint_file, "w") as f:
            json.dump(checkpoint_data, f)

        # Try to resume with different agent ID
        resumed_data = checkpoint_extension.resume_from_checkpoint(
            feature_id="F01", agent_id="different-agent-456"
        )

        assert resumed_data is None

    def test_resume_from_checkpoint_returns_none_if_not_exists(
        self, checkpoint_extension: AgentCheckpointExtension
    ) -> None:
        """Test resume_from_checkpoint() returns None if checkpoint doesn't exist."""
        resumed_data = checkpoint_extension.resume_from_checkpoint(
            feature_id="NONEXISTENT", agent_id="agent-123"
        )

        assert resumed_data is None

    def test_resume_from_checkpoint_restores_gates(
        self, checkpoint_extension: AgentCheckpointExtension, temp_checkpoint_dir: Path
    ) -> None:
        """Test resume_from_checkpoint() restores approval gate state."""
        checkpoint_data = {
            "feature_id": "F01",
            "agent_id": "agent-123",
            "workstreams": ["WS-001"],
            "completed_workstreams": [],
            "current_workstream": "WS-001",
            "status": "in_progress",
            "started_at": "2026-01-28T10:00:00",
            "gates": {
                "requirements": {"status": "approved", "approved_by": "user", "comments": "LGTM"}
            },
        }

        checkpoint_file = temp_checkpoint_dir / "F01-checkpoint.json"
        with open(checkpoint_file, "w") as f:
            json.dump(checkpoint_data, f)

        # Create gate manager
        gate_manager = ApprovalGateManager(checkpoint_extension.repo)
        checkpoint_extension.set_gate_manager(gate_manager)

        # Mock checkpoint repo
        mock_checkpoint = Checkpoint(
            feature="F01",
            agent_id="agent-123",
            status=CheckpointStatus.IN_PROGRESS,
            completed_ws=[],
            execution_order=["WS-001"],
            started_at=datetime.now(),
            metrics={"gates": checkpoint_data["gates"]},
        )
        checkpoint_extension.repo.load_checkpoint.return_value = mock_checkpoint

        resumed_data = checkpoint_extension.resume_from_checkpoint(
            feature_id="F01", agent_id="agent-123"
        )

        assert resumed_data is not None
        assert "gates" in resumed_data

    def test_resume_from_checkpoint_restores_team_state(
        self, checkpoint_extension: AgentCheckpointExtension, temp_checkpoint_dir: Path
    ) -> None:
        """Test resume_from_checkpoint() restores team configuration."""
        checkpoint_data = {
            "feature_id": "F01",
            "agent_id": "agent-123",
            "workstreams": ["WS-001"],
            "completed_workstreams": [],
            "current_workstream": "WS-001",
            "status": "in_progress",
            "started_at": "2026-01-28T10:00:00",
            "team": {
                "roles": [
                    {
                        "name": "orchestrator",
                        "state": "active",
                        "capabilities": ["coordinate"],
                    }
                ]
            },
        }

        checkpoint_file = temp_checkpoint_dir / "F01-checkpoint.json"
        with open(checkpoint_file, "w") as f:
            json.dump(checkpoint_data, f)

        resumed_data = checkpoint_extension.resume_from_checkpoint(
            feature_id="F01", agent_id="agent-123"
        )

        assert resumed_data is not None
        assert "team" in resumed_data
        assert len(resumed_data["team"]["roles"]) == 1


class TestOrchestratorAgentCheckpointIntegration:
    """Integration tests for checkpoint save/resume with approval gates and team."""

    def test_auto_skip_gates_on_resume(
        self, checkpoint_extension: AgentCheckpointExtension, temp_checkpoint_dir: Path
    ) -> None:
        """Test that skipped gates are auto-skipped on resume."""
        # Create gate manager with skip parser
        skip_parser = SkipFlagParser(["--skip-requirements"])
        gate_manager = ApprovalGateManager(
            checkpoint_extension.repo, skip_parser=skip_parser
        )
        checkpoint_extension.set_gate_manager(gate_manager)

        # Save checkpoint with skipped gate
        checkpoint_data = {
            "feature_id": "F01",
            "agent_id": "agent-123",
            "workstreams": ["WS-001"],
            "completed_workstreams": [],
            "current_workstream": "WS-001",
            "status": "in_progress",
            "started_at": "2026-01-28T10:00:00",
            "gates": {"requirements": {"status": "skipped", "reason": "--skip-requirements"}},
        }

        checkpoint_file = temp_checkpoint_dir / "F01-checkpoint.json"
        with open(checkpoint_file, "w") as f:
            json.dump(checkpoint_data, f)

        # Mock checkpoint repo
        mock_checkpoint = Checkpoint(
            feature="F01",
            agent_id="agent-123",
            status=CheckpointStatus.IN_PROGRESS,
            completed_ws=[],
            execution_order=["WS-001"],
            started_at=datetime.now(),
            metrics={"gates": checkpoint_data["gates"]},
        )
        checkpoint_extension.repo.load_checkpoint.return_value = mock_checkpoint

        # Resume and verify gate is still skipped
        resumed_data = checkpoint_extension.resume_from_checkpoint(
            feature_id="F01", agent_id="agent-123"
        )

        assert resumed_data is not None
        assert resumed_data["gates"]["requirements"]["status"] == "skipped"

    def test_save_and_resume_cycle(
        self, checkpoint_extension: AgentCheckpointExtension, temp_checkpoint_dir: Path
    ) -> None:
        """Test complete save and resume cycle."""
        # Save checkpoint
        checkpoint_extension.save_checkpoint(
            feature_id="F01",
            agent_id="agent-123",
            workstreams=["WS-001", "WS-002", "WS-003"],
            completed_workstreams=["WS-001"],
            current_workstream="WS-002",
            status=CheckpointStatus.IN_PROGRESS,
        )

        # Resume checkpoint
        resumed_data = checkpoint_extension.resume_from_checkpoint(
            feature_id="F01", agent_id="agent-123"
        )

        assert resumed_data is not None
        assert resumed_data["feature_id"] == "F01"
        assert resumed_data["completed_workstreams"] == ["WS-001"]
        assert resumed_data["current_workstream"] == "WS-002"
        assert resumed_data["status"] == "in_progress"

    def test_progress_tracking_across_executions(
        self, checkpoint_extension: AgentCheckpointExtension, temp_checkpoint_dir: Path
    ) -> None:
        """Test progress tracking is maintained across multiple executions."""
        # First execution: complete WS-001
        checkpoint_extension.save_checkpoint(
            feature_id="F01",
            agent_id="agent-123",
            workstreams=["WS-001", "WS-002", "WS-003"],
            completed_workstreams=["WS-001"],
            current_workstream="WS-002",
            status=CheckpointStatus.IN_PROGRESS,
        )

        # Simulate resume and complete WS-002
        resumed_data = checkpoint_extension.resume_from_checkpoint(
            feature_id="F01", agent_id="agent-123"
        )

        assert resumed_data is not None
        assert len(resumed_data["completed_workstreams"]) == 1

        # Update checkpoint with WS-002 completed
        checkpoint_extension.save_checkpoint(
            feature_id="F01",
            agent_id="agent-123",
            workstreams=["WS-001", "WS-002", "WS-003"],
            completed_workstreams=["WS-001", "WS-002"],
            current_workstream="WS-003",
            status=CheckpointStatus.IN_PROGRESS,
        )

        # Resume again and verify progress
        final_data = checkpoint_extension.resume_from_checkpoint(
            feature_id="F01", agent_id="agent-123"
        )

        assert final_data is not None
        assert len(final_data["completed_workstreams"]) == 2
        assert final_data["current_workstream"] == "WS-003"
