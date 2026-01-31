"""Tests for cli.py - CLI Beads client error handling and edge cases."""

import json
import subprocess
from pathlib import Path
from unittest.mock import Mock, patch

import pytest

from sdp.beads.cli import CLIBeadsClient
from sdp.beads.exceptions import BeadsClientError
from sdp.beads.models import BeadsStatus, BeadsTaskCreate, BeadsPriority


class TestCLIBeadsClientErrorHandling:
    """Test error handling in CLI client."""

    def test_init_beads_not_installed(self) -> None:
        """Test initialization fails when Beads CLI not found."""
        # Skip if bd is actually installed (environment-dependent test)
        import subprocess
        try:
            subprocess.run(["bd", "--version"], capture_output=True, check=True, timeout=1)
            pytest.skip("Beads CLI is installed, cannot test failure case")
        except (subprocess.CalledProcessError, FileNotFoundError, subprocess.TimeoutExpired):
            # Beads not installed - test should work
            with patch("sdp.beads.cli.subprocess.run") as mock_run:
                mock_run.side_effect = FileNotFoundError("bd: command not found")
                with pytest.raises(BeadsClientError) as exc_info:
                    CLIBeadsClient()
                assert "Beads CLI not found" in str(exc_info.value)

    def test_init_beads_version_check_fails(self) -> None:
        """Test initialization fails when version check fails."""
        with patch("subprocess.run") as mock_run:
            mock_run.side_effect = subprocess.CalledProcessError(1, "bd")
            with pytest.raises(BeadsClientError) as exc_info:
                CLIBeadsClient()
            assert "Beads CLI not found" in str(exc_info.value)

    def test_create_task_json_decode_error(self) -> None:
        """Test create_task handles JSON decode errors."""
        client = CLIBeadsClient.__new__(CLIBeadsClient)
        client.project_dir = Path.cwd()

        with patch("sdp.beads.cli.subprocess.run") as mock_run:
            mock_result = Mock()
            mock_result.stdout = "invalid json"
            mock_result.returncode = 0
            mock_run.return_value = mock_result

            with pytest.raises((BeadsClientError, json.JSONDecodeError)):
                client.create_task(
                    BeadsTaskCreate(title="Test", description="", priority=BeadsPriority.HIGH)
                )

    def test_get_task_not_found(self) -> None:
        """Test get_task returns None for non-existent task."""
        client = CLIBeadsClient.__new__(CLIBeadsClient)
        client.project_dir = Path.cwd()

        with patch("subprocess.run") as mock_run:
            mock_run.side_effect = subprocess.CalledProcessError(1, "bd")
            result = client.get_task("bd-9999")
            assert result is None

    def test_get_task_json_decode_error(self) -> None:
        """Test get_task handles JSON decode errors."""
        client = CLIBeadsClient.__new__(CLIBeadsClient)
        client.project_dir = Path.cwd()

        with patch("subprocess.run") as mock_run:
            mock_result = Mock()
            mock_result.stdout = "invalid json"
            mock_run.return_value = mock_result

            result = client.get_task("bd-0001")
            assert result is None

    def test_get_task_empty_array(self) -> None:
        """Test get_task handles empty array response."""
        client = CLIBeadsClient.__new__(CLIBeadsClient)
        client.project_dir = Path.cwd()

        with patch("sdp.beads.cli.subprocess.run") as mock_run:
            mock_result = Mock()
            mock_result.stdout = json.dumps([])
            mock_result.returncode = 0
            mock_run.return_value = mock_result

            # Empty array will cause from_dict to fail, which is caught
            result = client.get_task("bd-0001")
            # Should return None due to exception handling
            assert result is None

    def test_get_task_dict_response(self) -> None:
        """Test get_task handles dict response (non-array)."""
        client = CLIBeadsClient.__new__(CLIBeadsClient)
        client.project_dir = Path.cwd()

        with patch("subprocess.run") as mock_run:
            mock_result = Mock()
            mock_result.stdout = json.dumps({"id": "bd-0001", "title": "Test"})
            mock_run.return_value = mock_result

            result = client.get_task("bd-0001")
            assert result is not None
            assert result.id == "bd-0001"

    def test_update_task_status_command_fails(self) -> None:
        """Test update_task_status handles command failure."""
        client = CLIBeadsClient.__new__(CLIBeadsClient)
        client.project_dir = Path.cwd()

        with patch("subprocess.run") as mock_run:
            mock_run.side_effect = subprocess.CalledProcessError(1, "bd", stderr="Task not found")
            with pytest.raises(BeadsClientError) as exc_info:
                client.update_task_status("bd-9999", BeadsStatus.IN_PROGRESS)
            assert "Command failed" in str(exc_info.value)

    def test_get_ready_tasks_dict_response(self) -> None:
        """Test get_ready_tasks handles dict response with ready_tasks key."""
        client = CLIBeadsClient.__new__(CLIBeadsClient)
        client.project_dir = Path.cwd()

        with patch("subprocess.run") as mock_run:
            mock_result = Mock()
            mock_result.stdout = json.dumps({"ready_tasks": ["bd-0001", "bd-0002"]})
            mock_run.return_value = mock_result

            ready = client.get_ready_tasks()
            assert len(ready) == 2
            assert "bd-0001" in ready
            assert "bd-0002" in ready

    def test_list_tasks_dict_response(self) -> None:
        """Test list_tasks handles dict response with tasks key."""
        client = CLIBeadsClient.__new__(CLIBeadsClient)
        client.project_dir = Path.cwd()

        with patch("subprocess.run") as mock_run:
            mock_result = Mock()
            mock_result.stdout = json.dumps({
                "tasks": [
                    {"id": "bd-0001", "title": "Task 1"},
                    {"id": "bd-0002", "title": "Task 2"},
                ]
            })
            mock_run.return_value = mock_result

            tasks = client.list_tasks()
            assert len(tasks) == 2

    def test_create_task_with_all_options(self) -> None:
        """Test create_task builds command with all optional parameters."""
        client = CLIBeadsClient.__new__(CLIBeadsClient)
        client.project_dir = Path.cwd()

        with patch("subprocess.run") as mock_run:
            mock_result = Mock()
            mock_result.stdout = json.dumps({"id": "bd-0001", "title": "Test"})
            mock_run.return_value = mock_result

            from sdp.beads.models import BeadsDependency, BeadsDependencyType
            params = BeadsTaskCreate(
                title="Test Task",
                description="Test description",
                priority=BeadsPriority.HIGH,
                parent_id="bd-parent",
                external_ref="ref-123",
                dependencies=[
                    BeadsDependency(task_id="bd-dep1", type=BeadsDependencyType.BLOCKS)
                ],
            )

            client.create_task(params)

            # Verify command was called
            assert mock_run.called
            call_args = mock_run.call_args[0][0]
            assert "bd" in call_args
            assert "create" in call_args
            assert "Test Task" in call_args
            assert "--description" in call_args
            assert "--priority" in call_args
            assert "--parent" in call_args
            assert "--external-ref" in call_args
            assert "--deps" in call_args

    def test_create_task_long_title_truncated(self) -> None:
        """Test create_task truncates titles over 500 characters."""
        client = CLIBeadsClient.__new__(CLIBeadsClient)
        client.project_dir = Path.cwd()

        with patch("subprocess.run") as mock_run:
            mock_result = Mock()
            mock_result.stdout = json.dumps({"id": "bd-0001", "title": "Test"})
            mock_run.return_value = mock_result

            long_title = "x" * 600
            params = BeadsTaskCreate(title=long_title, description="", priority=BeadsPriority.HIGH)

            client.create_task(params)

            call_args = mock_run.call_args[0][0]
            title_arg = call_args[call_args.index("create") + 1]
            assert len(title_arg) == 500

    def test_create_task_empty_title(self) -> None:
        """Test create_task handles empty title."""
        client = CLIBeadsClient.__new__(CLIBeadsClient)
        client.project_dir = Path.cwd()

        with patch("subprocess.run") as mock_run:
            mock_result = Mock()
            mock_result.stdout = json.dumps({"id": "bd-0001", "title": "Untitled"})
            mock_run.return_value = mock_result

            params = BeadsTaskCreate(title="", description="", priority=BeadsPriority.HIGH)
            result = client.create_task(params)

            assert result is not None

    def test_run_command_json_decode_error_handling(self) -> None:
        """Test _run_command handles JSON decode errors."""
        client = CLIBeadsClient.__new__(CLIBeadsClient)
        client.project_dir = Path.cwd()

        with patch("subprocess.run") as mock_run:
            # Simulate JSON decode error during command execution
            mock_run.side_effect = json.JSONDecodeError("Invalid JSON", "", 0)
            with pytest.raises(BeadsClientError) as exc_info:
                client._run_command(["bd", "create", "test"], capture_output=True)
            assert "Invalid JSON response" in str(exc_info.value)
