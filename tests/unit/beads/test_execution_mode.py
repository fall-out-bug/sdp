"""
Tests for ExecutionMode and workflow efficiency features.

Tests follow AAA pattern (Arrange-Act-Assert) and use descriptive names.
"""

import pytest
from datetime import datetime, timezone
from pathlib import Path
from unittest.mock import Mock, patch, call, MagicMock
import json

from sdp.beads.execution_mode import (
    ExecutionMode,
    AuditLogger,
    DestructiveOperationDetector,
    DestructiveOperations,
)
from sdp.beads.skills_oneshot import MultiAgentExecutor
from sdp.beads.models import BeadsStatus


class TestExecutionMode:
    """Test ExecutionMode enum and its properties."""

    def test_standard_mode_requires_pr(self):
        """Standard mode should require PR approval."""
        mode = ExecutionMode.STANDARD
        assert mode.requires_pr is True
        assert mode.allows_production is True
        assert mode.is_sandbox is False

    def test_auto_approve_mode_skips_pr(self):
        """Auto-approve mode should skip PR and allow production."""
        mode = ExecutionMode.AUTO_APPROVE
        assert mode.requires_pr is False
        assert mode.allows_production is True
        assert mode.is_sandbox is False

    def test_sandbox_mode_skips_pr(self):
        """Sandbox mode should skip PR and not allow production."""
        mode = ExecutionMode.SANDBOX
        assert mode.requires_pr is False
        assert mode.allows_production is False
        assert mode.is_sandbox is True

    def test_dry_run_mode_is_preview(self):
        """Dry-run mode should be preview-only."""
        mode = ExecutionMode.DRY_RUN
        assert mode.requires_pr is False
        assert mode.allows_production is False
        assert mode.is_sandbox is False
        assert mode.is_preview is True


class TestAuditLogger:
    """Test audit logging for --auto-approve executions."""

    def test_log_auto_approve_execution(self, tmp_path):
        """Should log auto-approve execution to audit file."""
        logger = AuditLogger(audit_file=str(tmp_path / "audit.log"))

        logger.log_execution(
            feature_id="bd-0001",
            mode=ExecutionMode.AUTO_APPROVE,
            workstreams_executed=4,
            result="success",
            user="developer@example.com",
        )

        # Verify log file created and contains correct data
        assert (tmp_path / "audit.log").exists()

        with open(tmp_path / "audit.log") as f:
            logs = [json.loads(line) for line in f]

        assert len(logs) == 1
        log = logs[0]
        assert log["feature"] == "bd-0001"
        assert log["mode"] == "auto_approve"
        assert log["workstreams_executed"] == 4
        assert log["result"] == "success"
        assert log["user"] == "developer@example.com"
        assert "timestamp" in log

    def test_log_multiple_executions(self, tmp_path):
        """Should append multiple executions to audit log."""
        logger = AuditLogger(audit_file=str(tmp_path / "audit.log"))

        logger.log_execution(
            feature_id="bd-0001",
            mode=ExecutionMode.AUTO_APPROVE,
            workstreams_executed=4,
            result="success",
        )

        logger.log_execution(
            feature_id="bd-0002",
            mode=ExecutionMode.SANDBOX,
            workstreams_executed=3,
            result="success",
        )

        with open(tmp_path / "audit.log") as f:
            logs = [json.loads(line) for line in f]

        assert len(logs) == 2
        assert logs[0]["feature"] == "bd-0001"
        assert logs[1]["feature"] == "bd-0002"

    def test_read_recent_logs(self, tmp_path):
        """Should read last N log entries."""
        logger = AuditLogger(audit_file=str(tmp_path / "audit.log"))

        # Create 5 log entries
        for i in range(5):
            logger.log_execution(
                feature_id=f"bd-{i:04d}",
                mode=ExecutionMode.AUTO_APPROVE,
                workstreams_executed=i + 1,
                result="success",
            )

        # Read last 3
        recent = logger.read_recent(count=3)

        assert len(recent) == 3
        assert recent[0]["feature"] == "bd-0002"  # Third from last
        assert recent[1]["feature"] == "bd-0003"
        assert recent[2]["feature"] == "bd-0004"  # Most recent


class TestDestructiveOperationDetector:
    """Test detection of destructive operations."""

    def test_detect_database_migration(self):
        """Should detect database migration files."""
        detector = DestructiveOperationDetector()

        files_to_create = [
            "src/sdp/migrations/002_add_users_table.py",
        ]

        destructive = detector.check_operations(
            files_to_create=files_to_create,
            files_to_modify=[],
            files_to_delete=[],
        )

        assert destructive.has_destructive_operations is True
        assert "database_migration" in destructive.operation_types
        assert len(destructive.files_affected) > 0

    def test_detect_file_deletion(self):
        """Should detect file deletion operations."""
        detector = DestructiveOperationDetector()

        files_to_delete = [
            "src/sdp/old_module.py",
            "docs/deprecated.md",
        ]

        destructive = detector.check_operations(
            files_to_create=[],
            files_to_modify=[],
            files_to_delete=files_to_delete,
        )

        assert destructive.has_destructive_operations is True
        assert "file_deletion" in destructive.operation_types
        assert len(destructive.files_affected) == 2

    def test_no_destructive_operations(self):
        """Should return False when no destructive operations."""
        detector = DestructiveOperationDetector()

        destructive = detector.check_operations(
            files_to_create=["src/sdp/new_module.py"],
            files_to_modify=["src/sdp/existing.py"],
            files_to_delete=[],
        )

        assert destructive.has_destructive_operations is False
        assert len(destructive.operation_types) == 0


class TestMultiAgentExecutorWithModes:
    """Test MultiAgentExecutor with execution modes."""

    def test_execute_with_auto_approve_mode(self):
        """Should execute without PR in auto-approve mode."""
        client = Mock()
        # Return ready tasks for 2 rounds, then empty
        client.get_ready_tasks.side_effect = [
            ["bd-0001.1", "bd-0001.2"],  # Round 1
            [],  # Round 2 (no more tasks)
        ]

        # Setup mock tasks
        mock_task_1 = Mock()
        mock_task_1.id = "bd-0001.1"
        mock_task_1.parent_id = "bd-0001"
        mock_task_1.title = "First task"

        mock_task_2 = Mock()
        mock_task_2.id = "bd-0001.2"
        mock_task_2.parent_id = "bd-0001"
        mock_task_2.title = "Second task"

        client.get_task.side_effect = [
            mock_task_1,
            mock_task_2,
        ]

        executor = MultiAgentExecutor(client, num_agents=2)

        # Mock the build executor to return success
        with patch.object(executor.build_executor, 'execute') as mock_execute:
            mock_execute.return_value = Mock(success=True)

            result = executor.execute_feature(
                "bd-0001",
                mode=ExecutionMode.AUTO_APPROVE,
                mock_success=True,
            )

        assert result.success is True
        assert result.total_executed == 2
        assert result.mode == ExecutionMode.AUTO_APPROVE
        assert result.pr_created is False

    def test_execute_with_sandbox_mode(self):
        """Should execute in sandbox mode without PR."""
        client = Mock()
        client.get_ready_tasks.side_effect = [
            ["bd-0001.1"],
            [],
        ]

        mock_task = Mock()
        mock_task.id = "bd-0001.1"
        mock_task.parent_id = "bd-0001"

        client.get_task.return_value = mock_task

        executor = MultiAgentExecutor(client, num_agents=1)

        with patch.object(executor.build_executor, 'execute') as mock_execute:
            mock_execute.return_value = Mock(success=True)

            result = executor.execute_feature(
                "bd-0001",
                mode=ExecutionMode.SANDBOX,
                mock_success=True,
            )

        assert result.success is True
        assert result.total_executed == 1
        assert result.deployment_target == "sandbox"
        assert result.mode == ExecutionMode.SANDBOX

    def test_execute_with_dry_run_mode(self):
        """Should preview changes without executing in dry-run mode."""
        client = Mock()
        client.get_ready_tasks.return_value = ["bd-0001.1", "bd-0001.2"]

        mock_task_1 = Mock()
        mock_task_1.id = "bd-0001.1"
        mock_task_1.parent_id = "bd-0001"
        mock_task_1.title = "Domain entities"

        mock_task_2 = Mock()
        mock_task_2.id = "bd-0001.2"
        mock_task_2.parent_id = "bd-0001"
        mock_task_2.title = "Repository layer"

        # Use return_value instead of side_effect since we need to call get_task multiple times
        # Create a mapping by id
        def get_task_side_effect(task_id):
            if task_id == "bd-0001.1":
                return mock_task_1
            elif task_id == "bd-0001.2":
                return mock_task_2
            return None

        client.get_task.side_effect = get_task_side_effect

        executor = MultiAgentExecutor(client, num_agents=2)
        result = executor.execute_feature(
            "bd-0001",
            mode=ExecutionMode.DRY_RUN,
            mock_success=True,
        )

        assert result.success is True
        assert result.total_executed == 0  # No actual execution
        assert result.preview_only is True
        assert len(result.tasks_preview) == 2
        assert "bd-0001.1: Domain entities" in result.tasks_preview
        assert "bd-0001.2: Repository layer" in result.tasks_preview

    def test_audit_logging_for_auto_approve(self, tmp_path):
        """Should log auto-approve executions to audit file."""
        client = Mock()
        client.get_ready_tasks.side_effect = [
            ["bd-0001.1"],
            [],
        ]

        mock_task = Mock()
        mock_task.id = "bd-0001.1"
        mock_task.parent_id = "bd-0001"

        client.get_task.return_value = mock_task

        audit_logger = AuditLogger(audit_file=str(tmp_path / "audit.log"))
        executor = MultiAgentExecutor(client, num_agents=1, audit_logger=audit_logger)

        with patch.object(executor.build_executor, 'execute') as mock_execute:
            mock_execute.return_value = Mock(success=True)

            result = executor.execute_feature(
                "bd-0001",
                mode=ExecutionMode.AUTO_APPROVE,
                mock_success=True,
            )

        assert result.success is True

        # Verify audit log was written
        with open(tmp_path / "audit.log") as f:
            logs = [json.loads(line) for line in f]

        assert len(logs) == 1
        assert logs[0]["mode"] == "auto_approve"
        assert logs[0]["feature"] == "bd-0001"
        assert logs[0]["workstreams_executed"] == 1
