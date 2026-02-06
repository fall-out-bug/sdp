"""
Tests for OrchestratorAgent core logic (WS-003).

Tests follow TDD: Red -> Green -> Refactor.
"""

import pytest
from datetime import datetime
from unittest.mock import Mock, MagicMock, patch
from pathlib import Path

from sdp.unified.checkpoint.schema import Checkpoint, CheckpointStatus
from sdp.unified.checkpoint.repository import CheckpointRepository, RepositoryError
from sdp.unified.orchestrator.agent import OrchestratorAgent
from sdp.unified.orchestrator.models import ExecutionResult
from sdp.unified.orchestrator.errors import ExecutionError


class TestOrchestratorAgentInit:
    """Test OrchestratorAgent initialization."""

    def test_initialization_with_repository(self):
        """Should initialize with CheckpointRepository."""
        mock_repo = Mock(spec=CheckpointRepository)
        agent = OrchestratorAgent(mock_repo)

        assert agent.repo == mock_repo

    def test_initialization_fails_without_repository(self):
        """Should fail if repository is not CheckpointRepository."""
        with pytest.raises(TypeError):
            OrchestratorAgent(None)  # type: ignore


class TestExecuteFeature:
    """Test execute_feature method."""

    def test_execute_feature_creates_checkpoint(self):
        """Should create initial checkpoint when starting feature."""
        mock_repo = Mock(spec=CheckpointRepository)
        agent = OrchestratorAgent(mock_repo)

        # Mock repository to return checkpoint ID and no existing checkpoint
        mock_repo.save_checkpoint.return_value = 1
        mock_repo.load_latest_checkpoint.return_value = None

        # Mock workstream list
        workstreams = ["WS-001", "WS-002", "WS-003"]

        result = agent.execute_feature(
            feature_id="F01",
            workstreams=workstreams,
            agent_id="agent-123",
        )

        # Verify checkpoint was created
        assert mock_repo.save_checkpoint.called
        checkpoint_arg = mock_repo.save_checkpoint.call_args[0][0]
        assert checkpoint_arg.feature == "F01"
        assert checkpoint_arg.agent_id == "agent-123"
        assert checkpoint_arg.status == CheckpointStatus.IN_PROGRESS

    def test_execute_feature_returns_execution_result(self):
        """Should return ExecutionResult with feature details."""
        mock_repo = Mock(spec=CheckpointRepository)
        agent = OrchestratorAgent(mock_repo)

        mock_repo.save_checkpoint.return_value = 1
        mock_repo.load_latest_checkpoint.return_value = None  # No existing checkpoint

        workstreams = ["WS-001", "WS-002"]
        result = agent.execute_feature(
            feature_id="F01",
            workstreams=workstreams,
            agent_id="agent-123",
        )

        assert isinstance(result, ExecutionResult)
        assert result.feature_id == "F01"
        assert result.agent_id == "agent-123"
        assert result.status == CheckpointStatus.COMPLETED

    def test_execute_feature_handles_empty_workstream_list(self):
        """Should handle empty workstream list gracefully."""
        mock_repo = Mock(spec=CheckpointRepository)
        agent = OrchestratorAgent(mock_repo)

        mock_repo.save_checkpoint.return_value = 1
        mock_repo.load_latest_checkpoint.return_value = None

        result = agent.execute_feature(
            feature_id="F01",
            workstreams=[],
            agent_id="agent-123",
        )

        assert result.status == CheckpointStatus.COMPLETED
        assert len(result.completed_workstreams) == 0

    def test_execute_feature_resumes_from_checkpoint(self):
        """Should resume execution from existing checkpoint."""
        mock_repo = Mock(spec=CheckpointRepository)
        agent = OrchestratorAgent(mock_repo)

        # Create existing checkpoint
        existing_checkpoint = Checkpoint(
            feature="F01",
            agent_id="agent-123",
            status=CheckpointStatus.IN_PROGRESS,
            completed_ws=["WS-001"],
            execution_order=["WS-001", "WS-002", "WS-003"],
            started_at=datetime.now(),
            current_ws="WS-002",
        )

        mock_repo.load_latest_checkpoint.return_value = existing_checkpoint
        mock_repo.save_checkpoint.return_value = 2

        workstreams = ["WS-001", "WS-002", "WS-003"]
        result = agent.execute_feature(
            feature_id="F01",
            workstreams=workstreams,
            agent_id="agent-123",
        )

        # Should resume from WS-002 (skip WS-001)
        assert result.completed_workstreams == ["WS-001", "WS-002", "WS-003"]

    def test_execute_feature_handles_repository_error(self):
        """Should handle RepositoryError gracefully."""
        mock_repo = Mock(spec=CheckpointRepository)
        agent = OrchestratorAgent(mock_repo)

        mock_repo.save_checkpoint.side_effect = RepositoryError("Database error")

        workstreams = ["WS-001"]
        with pytest.raises(ExecutionError, match="Failed to execute feature"):
            agent.execute_feature(
                feature_id="F01",
                workstreams=workstreams,
                agent_id="agent-123",
            )


class TestDispatchWorkstreams:
    """Test dispatch_workstreams method."""

    def test_dispatch_workstreams_returns_all_completed(self):
        """Should dispatch all workstreams and return completed list."""
        mock_repo = Mock(spec=CheckpointRepository)
        agent = OrchestratorAgent(mock_repo)

        workstreams = ["WS-001", "WS-002", "WS-003"]
        checkpoint_id = 1

        completed = agent.dispatch_workstreams(
            workstreams=workstreams,
            checkpoint_id=checkpoint_id,
            start_index=0,
        )

        assert completed == workstreams

    def test_dispatch_workstreams_starts_from_index(self):
        """Should start dispatching from given index."""
        mock_repo = Mock(spec=CheckpointRepository)
        agent = OrchestratorAgent(mock_repo)

        workstreams = ["WS-001", "WS-002", "WS-003"]
        checkpoint_id = 1

        completed = agent.dispatch_workstreams(
            workstreams=workstreams,
            checkpoint_id=checkpoint_id,
            start_index=1,  # Start from WS-002
        )

        # Should only execute WS-002 and WS-003
        assert completed == ["WS-002", "WS-003"]

    def test_dispatch_workstreams_updates_checkpoint(self):
        """Should update checkpoint after each workstream."""
        mock_repo = Mock(spec=CheckpointRepository)
        agent = OrchestratorAgent(mock_repo)

        workstreams = ["WS-001", "WS-002"]
        checkpoint_id = 1

        agent.dispatch_workstreams(
            workstreams=workstreams,
            checkpoint_id=checkpoint_id,
            start_index=0,
        )

        # Verify checkpoint was updated twice
        assert mock_repo.update_checkpoint_status.call_count == 2

    def test_dispatch_workstreams_handles_dispatch_error(self):
        """Should handle errors during workstream dispatch."""
        mock_repo = Mock(spec=CheckpointRepository)
        agent = OrchestratorAgent(mock_repo)

        workstreams = ["WS-001", "WS-002"]
        checkpoint_id = 1

        # Mock dispatch failure
        with patch.object(agent, '_dispatch_single_workstream', side_effect=Exception("Dispatch failed")):
            with pytest.raises(ExecutionError, match="Workstream dispatch failed"):
                agent.dispatch_workstreams(
                    workstreams=workstreams,
                    checkpoint_id=checkpoint_id,
                    start_index=0,
                )


class TestMonitorProgress:
    """Test monitor_progress method."""

    def test_monitor_progress_returns_progress_metrics(self):
        """Should return progress metrics for feature."""
        mock_repo = Mock(spec=CheckpointRepository)
        agent = OrchestratorAgent(mock_repo)

        checkpoint = Checkpoint(
            feature="F01",
            agent_id="agent-123",
            status=CheckpointStatus.IN_PROGRESS,
            completed_ws=["WS-001", "WS-002"],
            execution_order=["WS-001", "WS-002", "WS-003", "WS-004"],
            started_at=datetime.now(),
            current_ws="WS-003",
        )

        mock_repo.load_checkpoint.return_value = checkpoint

        progress = agent.monitor_progress(feature_id="F01")

        assert progress["total_workstreams"] == 4
        assert progress["completed_workstreams"] == 2
        assert progress["current_workstream"] == "WS-003"
        assert progress["percentage"] == 50.0

    def test_monitor_progress_handles_no_checkpoint(self):
        """Should return None if no checkpoint exists."""
        mock_repo = Mock(spec=CheckpointRepository)
        agent = OrchestratorAgent(mock_repo)

        mock_repo.load_checkpoint.return_value = None

        progress = agent.monitor_progress(feature_id="F01")

        assert progress is None

    def test_monitor_progress_handles_completed_feature(self):
        """Should show 100% for completed feature."""
        mock_repo = Mock(spec=CheckpointRepository)
        agent = OrchestratorAgent(mock_repo)

        checkpoint = Checkpoint(
            feature="F01",
            agent_id="agent-123",
            status=CheckpointStatus.COMPLETED,
            completed_ws=["WS-001", "WS-002", "WS-003"],
            execution_order=["WS-001", "WS-002", "WS-003"],
            started_at=datetime.now(),
            completed_at=datetime.now(),
        )

        mock_repo.load_checkpoint.return_value = checkpoint

        progress = agent.monitor_progress(feature_id="F01")

        assert progress["percentage"] == 100.0
        assert progress["status"] == CheckpointStatus.COMPLETED.value


class TestExecutionResult:
    """Test ExecutionResult dataclass."""

    def test_execution_result_creation(self):
        """Should create ExecutionResult with all fields."""
        result = ExecutionResult(
            feature_id="F01",
            agent_id="agent-123",
            status=CheckpointStatus.COMPLETED,
            completed_workstreams=["WS-001", "WS-002"],
            failed_workstreams=[],
            duration_seconds=3600,
        )

        assert result.feature_id == "F01"
        assert result.agent_id == "agent-123"
        assert result.status == CheckpointStatus.COMPLETED
        assert result.duration_seconds == 3600

    def test_execution_result_has_success_property(self):
        """Should provide success property based on status."""
        success_result = ExecutionResult(
            feature_id="F01",
            agent_id="agent-123",
            status=CheckpointStatus.COMPLETED,
            completed_workstreams=["WS-001"],
            failed_workstreams=[],
            duration_seconds=100,
        )

        failed_result = ExecutionResult(
            feature_id="F01",
            agent_id="agent-123",
            status=CheckpointStatus.FAILED,
            completed_workstreams=["WS-001"],
            failed_workstreams=["WS-002"],
            duration_seconds=100,
        )

        assert success_result.is_success is True
        assert failed_result.is_success is False
