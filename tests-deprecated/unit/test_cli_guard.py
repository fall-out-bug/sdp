"""Unit tests for guard CLI commands."""

import json
import tempfile
from pathlib import Path
from unittest.mock import Mock, patch

import pytest
from click.testing import CliRunner

from sdp.cli.main import main
from sdp.guard.state import StateManager


@pytest.fixture
def runner() -> CliRunner:
    """Create CLI test runner."""
    return CliRunner()


@pytest.fixture
def temp_state_file(tmp_path: Path) -> Path:
    """Create temporary state file."""
    state_file = tmp_path / ".sdp" / "state.json"
    return state_file


@pytest.fixture(autouse=True)
def mock_state_file(temp_state_file: Path, monkeypatch: pytest.MonkeyPatch) -> None:
    """Mock StateManager.STATE_FILE for all tests."""
    monkeypatch.setattr(StateManager, "STATE_FILE", temp_state_file)


class TestGuardCLI:
    """Test guard CLI commands."""

    def test_check_without_active_ws_fails(self, runner: CliRunner) -> None:
        """AC1: Check without active WS returns exit 1."""
        result = runner.invoke(main, ["guard", "check", "any.py"])
        assert result.exit_code == 1
        assert "No active WS" in result.output

    @patch("sdp.cli.guard.create_beads_client")
    def test_activate_saves_state(
        self, mock_create_client: Mock, runner: CliRunner, temp_state_file: Path
    ) -> None:
        """AC2: Activate saves state to file."""
        # Mock Beads client
        mock_client = Mock()
        mock_task = Mock()
        mock_task.sdp_metadata = {"scope_files": ["src/test.py"]}
        mock_client.get_task.return_value = mock_task
        mock_create_client.return_value = mock_client

        result = runner.invoke(main, ["guard", "activate", "00-032-01"])
        assert result.exit_code == 0
        assert "Activated" in result.output

        # Check state file
        state = StateManager.load()
        assert state.active_ws == "00-032-01"
        assert state.scope_files == ["src/test.py"]

    @patch("sdp.cli.guard.create_beads_client")
    def test_status_shows_active_ws(
        self, mock_create_client: Mock, runner: CliRunner
    ) -> None:
        """AC3: Status shows current WS."""
        # Mock Beads client
        mock_client = Mock()
        mock_task = Mock()
        mock_task.sdp_metadata = {}
        mock_client.get_task.return_value = mock_task
        mock_create_client.return_value = mock_client

        # Activate first
        runner.invoke(main, ["guard", "activate", "00-032-01"])

        # Check status
        result = runner.invoke(main, ["guard", "status"])
        assert result.exit_code == 0
        assert "00-032-01" in result.output

    @patch("sdp.cli.guard.create_beads_client")
    def test_deactivate_clears_state(
        self, mock_create_client: Mock, runner: CliRunner
    ) -> None:
        """AC4: Deactivate clears active WS."""
        # Mock Beads client
        mock_client = Mock()
        mock_task = Mock()
        mock_task.sdp_metadata = {}
        mock_client.get_task.return_value = mock_task
        mock_create_client.return_value = mock_client

        # Activate first
        runner.invoke(main, ["guard", "activate", "00-032-01"])

        # Deactivate
        result = runner.invoke(main, ["guard", "deactivate"])
        assert result.exit_code == 0

        state = StateManager.load()
        assert state.active_ws is None

    @patch("sdp.cli.guard.create_beads_client")
    def test_check_allowed_file_succeeds(
        self, mock_create_client: Mock, runner: CliRunner
    ) -> None:
        """AC1: Check with file in scope returns exit 0."""
        # Mock Beads client
        mock_client = Mock()
        mock_task = Mock()
        mock_task.sdp_metadata = {"scope_files": ["src/test.py"]}
        mock_client.get_task.return_value = mock_task
        mock_create_client.return_value = mock_client

        # Activate
        runner.invoke(main, ["guard", "activate", "00-032-01"])

        # Check allowed file
        result = runner.invoke(main, ["guard", "check", "src/test.py"])
        assert result.exit_code == 0
        assert "allowed" in result.output

    @patch("sdp.cli.guard.create_beads_client")
    def test_check_forbidden_file_fails(
        self, mock_create_client: Mock, runner: CliRunner
    ) -> None:
        """AC1: Check with file outside scope returns exit 1."""
        # Mock Beads client
        mock_client = Mock()
        mock_task = Mock()
        mock_task.sdp_metadata = {"scope_files": ["src/test.py"]}
        mock_client.get_task.return_value = mock_task
        mock_create_client.return_value = mock_client

        # Activate
        runner.invoke(main, ["guard", "activate", "00-032-01"])

        # Check forbidden file
        result = runner.invoke(main, ["guard", "check", "src/forbidden.py"])
        assert result.exit_code == 1
        assert "not in WS scope" in result.output

    def test_status_without_active_ws(self, runner: CliRunner) -> None:
        """Status without active WS shows appropriate message."""
        result = runner.invoke(main, ["guard", "status"])
        assert result.exit_code == 0
        assert "No active workstream" in result.output
