"""Tests for real Beads CLI client (CLIBeadsClient)."""

import json
import subprocess
from pathlib import Path
from unittest.mock import Mock, patch

import pytest
from sdp.beads.cli import CLIBeadsClient
from sdp.beads.exceptions import BeadsClientError
from sdp.beads.models import BeadsPriority, BeadsStatus, BeadsTaskCreate


class TestCLIBeadsClient:
    """Test suite for CLIBeadsClient."""

    @patch("subprocess.run")
    def test_initialization_checks_bd_available(self, mock_run: Mock) -> None:
        """Client initialization verifies bd CLI is available."""
        # Arrange
        mock_run.return_value = Mock(returncode=0, stdout="bd version 0.1.0", stderr="")

        # Act
        client = CLIBeadsClient(project_dir=Path("/tmp/test"))

        # Assert
        mock_run.assert_called_once()
        call_args = mock_run.call_args[0][0]
        assert "bd" in call_args
        assert "--version" in call_args

    @patch("subprocess.run")
    def test_initialization_raises_when_bd_not_found(self, mock_run: Mock) -> None:
        """AC2: Initialization raises BeadsClientError when bd not installed."""
        # Arrange
        mock_run.side_effect = FileNotFoundError("bd not found")

        # Act & Assert
        with pytest.raises(BeadsClientError) as exc_info:
            CLIBeadsClient(project_dir=Path("/tmp/test"))

        assert "not found" in str(exc_info.value).lower()
        assert "go install" in str(exc_info.value)

    @patch("subprocess.run")
    def test_create_task_calls_bd_create(self, mock_run: Mock) -> None:
        """AC2: create_task method calls bd create with JSON params."""
        # Arrange - initialization check
        mock_run.return_value = Mock(returncode=0, stdout="bd version 0.1.0", stderr="")
        client = CLIBeadsClient(project_dir=Path("/tmp/test"))

        # Arrange - create task
        task_data = {
            "id": "bd-0001",
            "title": "Test Task",
            "description": "Test description",
            "status": "open",
            "priority": 2,
        }
        mock_run.return_value = Mock(
            returncode=0,
            stdout=json.dumps(task_data),
            stderr="",
        )

        params = BeadsTaskCreate(
            title="Test Task",
            description="Test description",
            priority=BeadsPriority.MEDIUM,
        )

        # Act
        task = client.create_task(params)

        # Assert
        assert task.id == "bd-0001"
        assert task.title == "Test Task"

        # Check subprocess call
        create_call = mock_run.call_args_list[-1]
        call_cmd = create_call[0][0]
        assert "bd" in call_cmd
        assert "create" in call_cmd
        assert "--json" in call_cmd

    @patch("subprocess.run")
    def test_get_task_returns_task(self, mock_run: Mock) -> None:
        """AC2: get_task method calls bd show and parses response."""
        # Arrange - initialization
        mock_run.return_value = Mock(returncode=0, stdout="bd version 0.1.0", stderr="")
        client = CLIBeadsClient(project_dir=Path("/tmp/test"))

        # Arrange - get task
        task_data = [
            {
                "id": "bd-0001",
                "title": "Test Task",
                "status": "open",
                "priority": 2,
            }
        ]
        mock_run.return_value = Mock(
            returncode=0,
            stdout=json.dumps(task_data),
            stderr="",
        )

        # Act
        task = client.get_task("bd-0001")

        # Assert
        assert task is not None
        assert task.id == "bd-0001"

        # Check subprocess call
        get_call = mock_run.call_args_list[-1]
        call_cmd = get_call[0][0]
        assert "bd" in call_cmd
        assert "show" in call_cmd
        assert "bd-0001" in call_cmd

    @patch("subprocess.run")
    def test_get_task_returns_none_when_not_found(self, mock_run: Mock) -> None:
        """get_task returns None when task not found."""
        # Arrange - initialization
        mock_run.return_value = Mock(returncode=0, stdout="bd version 0.1.0", stderr="")
        client = CLIBeadsClient(project_dir=Path("/tmp/test"))

        # Arrange - task not found
        mock_run.side_effect = subprocess.CalledProcessError(
            returncode=1,
            cmd=["bd", "show", "nonexistent"],
            stderr="Task not found",
        )

        # Act
        task = client.get_task("nonexistent")

        # Assert
        assert task is None

    @patch("subprocess.run")
    def test_update_task_status(self, mock_run: Mock) -> None:
        """AC2: update_task_status calls bd update with status."""
        # Arrange - initialization
        mock_run.return_value = Mock(returncode=0, stdout="bd version 0.1.0", stderr="")
        client = CLIBeadsClient(project_dir=Path("/tmp/test"))

        # Arrange - update status
        mock_run.return_value = Mock(returncode=0, stdout="", stderr="")

        # Act
        client.update_task_status("bd-0001", BeadsStatus.IN_PROGRESS)

        # Assert
        update_call = mock_run.call_args_list[-1]
        call_cmd = update_call[0][0]
        assert "bd" in call_cmd
        assert "update" in call_cmd
        assert "bd-0001" in call_cmd
        assert "--status" in call_cmd
        assert "in_progress" in call_cmd

    @patch("subprocess.run")
    def test_update_metadata(self, mock_run: Mock) -> None:
        """update_metadata calls bd update with metadata JSON."""
        # Arrange - initialization
        mock_run.return_value = Mock(returncode=0, stdout="bd version 0.1.0", stderr="")
        client = CLIBeadsClient(project_dir=Path("/tmp/test"))

        # Arrange - update metadata
        mock_run.return_value = Mock(returncode=0, stdout="", stderr="")

        metadata = {"scope_files": ["src/file1.py", "src/file2.py"]}

        # Act
        client.update_metadata("bd-0001", metadata)

        # Assert
        update_call = mock_run.call_args_list[-1]
        call_cmd = update_call[0][0]
        assert "bd" in call_cmd
        assert "update" in call_cmd
        assert "bd-0001" in call_cmd
        assert "--metadata" in call_cmd
        # Check metadata JSON is in command
        assert any(json.dumps(metadata) in str(arg) for arg in call_cmd)

    @patch("subprocess.run")
    def test_cli_error_handling(self, mock_run: Mock) -> None:
        """AC3: CLI errors are caught and raised as BeadsClientError."""
        # Arrange - initialization
        mock_run.return_value = Mock(returncode=0, stdout="bd version 0.1.0", stderr="")
        client = CLIBeadsClient(project_dir=Path("/tmp/test"))

        # Arrange - command fails
        mock_run.side_effect = subprocess.CalledProcessError(
            returncode=1,
            cmd=["bd", "update", "bd-0001"],
            stderr="Update failed: invalid status",
        )

        # Act & Assert
        with pytest.raises(BeadsClientError) as exc_info:
            client.update_task_status("bd-0001", BeadsStatus.IN_PROGRESS)

        assert "failed" in str(exc_info.value).lower()

    @patch("subprocess.run")
    def test_json_parsing_error_handling(self, mock_run: Mock) -> None:
        """AC4: Invalid JSON responses are handled (return None for get_task)."""
        # Arrange - initialization
        mock_run.return_value = Mock(returncode=0, stdout="bd version 0.1.0", stderr="")
        client = CLIBeadsClient(project_dir=Path("/tmp/test"))

        # Arrange - invalid JSON
        mock_run.return_value = Mock(
            returncode=0,
            stdout="This is not valid JSON",
            stderr="",
        )

        # Act - get_task returns None on JSON error
        task = client.get_task("bd-0001")

        # Assert
        assert task is None

    @patch("subprocess.run")
    def test_list_tasks_with_filters(self, mock_run: Mock) -> None:
        """list_tasks calls bd list with optional filters."""
        # Arrange - initialization
        mock_run.return_value = Mock(returncode=0, stdout="bd version 0.1.0", stderr="")
        client = CLIBeadsClient(project_dir=Path("/tmp/test"))

        # Arrange - list tasks
        tasks_data = [
            {"id": "bd-0001", "title": "Task 1", "status": "open", "priority": 2},
            {"id": "bd-0002", "title": "Task 2", "status": "open", "priority": 2},
        ]
        mock_run.return_value = Mock(
            returncode=0,
            stdout=json.dumps(tasks_data),
            stderr="",
        )

        # Act
        tasks = client.list_tasks(status=BeadsStatus.OPEN, parent_id="bd-0000")

        # Assert
        assert len(tasks) == 2
        assert tasks[0].id == "bd-0001"

        # Check subprocess call includes filters
        list_call = mock_run.call_args_list[-1]
        call_cmd = list_call[0][0]
        assert "bd" in call_cmd
        assert "list" in call_cmd
        assert "--status" in call_cmd
        assert "open" in call_cmd
        assert "--parent" in call_cmd
        assert "bd-0000" in call_cmd

    @patch("subprocess.run")
    def test_get_ready_tasks(self, mock_run: Mock) -> None:
        """get_ready_tasks calls bd ready."""
        # Arrange - initialization
        mock_run.return_value = Mock(returncode=0, stdout="bd version 0.1.0", stderr="")
        client = CLIBeadsClient(project_dir=Path("/tmp/test"))

        # Arrange - ready tasks
        ready_data = ["bd-0001", "bd-0002"]
        mock_run.return_value = Mock(
            returncode=0,
            stdout=json.dumps(ready_data),
            stderr="",
        )

        # Act
        ready_ids = client.get_ready_tasks()

        # Assert
        assert len(ready_ids) == 2
        assert "bd-0001" in ready_ids
        assert "bd-0002" in ready_ids

    @patch("subprocess.run")
    def test_add_dependency(self, mock_run: Mock) -> None:
        """add_dependency calls bd dep add."""
        # Arrange - initialization
        mock_run.return_value = Mock(returncode=0, stdout="bd version 0.1.0", stderr="")
        client = CLIBeadsClient(project_dir=Path("/tmp/test"))

        # Arrange - add dependency
        mock_run.return_value = Mock(returncode=0, stdout="", stderr="")

        # Act
        client.add_dependency("bd-0001", "bd-0002", dep_type="blocks")

        # Assert
        dep_call = mock_run.call_args_list[-1]
        call_cmd = dep_call[0][0]
        assert "bd" in call_cmd
        assert "dep" in call_cmd
        assert "add" in call_cmd
        assert "bd-0001" in call_cmd
        assert "bd-0002" in call_cmd
        assert "--type" in call_cmd
        assert "blocks" in call_cmd
