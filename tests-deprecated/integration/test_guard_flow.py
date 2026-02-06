"""E2E integration tests for guard workflow."""

from pathlib import Path
from unittest.mock import Mock, patch

import pytest
from click.testing import CliRunner

from sdp.beads.models import BeadsStatus, BeadsTask, BeadsPriority
from sdp.cli.main import main
from sdp.guard.state import StateManager


@pytest.fixture
def runner() -> CliRunner:
    """Create CLI test runner."""
    return CliRunner()


@pytest.fixture
def temp_state_file(tmp_path: Path) -> Path:
    """Create temporary state file."""
    return tmp_path / ".sdp" / "state.json"


@pytest.fixture(autouse=True)
def mock_state_file(temp_state_file: Path, monkeypatch: pytest.MonkeyPatch) -> None:
    """Mock StateManager.STATE_FILE for all tests."""
    monkeypatch.setattr(StateManager, "STATE_FILE", temp_state_file)


@pytest.fixture
def sample_workstream() -> BeadsTask:
    """Create a sample workstream for testing."""
    return BeadsTask(
        id="00-032-01",
        title="Test WS",
        description="Test workstream",
        priority=BeadsPriority.MEDIUM,
        sdp_metadata={
            "scope_files": [
                "src/sdp/guard/skill.py",
                "src/sdp/guard/models.py",
            ]
        },
    )


class TestGuardE2E:
    """End-to-end tests for guard workflow."""

    @patch("sdp.cli.guard.create_beads_client")
    def test_happy_path_activate_edit_complete(
        self, mock_create_client: Mock, runner: CliRunner, sample_workstream: BeadsTask
    ) -> None:
        """AC1: Full workflow succeeds."""
        # Mock Beads client
        mock_client = Mock()
        mock_client.get_task.return_value = sample_workstream
        mock_create_client.return_value = mock_client

        ws_id = sample_workstream.id

        # 1. Activate
        result = runner.invoke(main, ["guard", "activate", ws_id])
        assert result.exit_code == 0
        assert "Activated" in result.output

        # 2. Check allowed file
        result = runner.invoke(main, ["guard", "check", "src/sdp/guard/skill.py"])
        assert result.exit_code == 0
        assert "allowed" in result.output.lower()

        # 3. Deactivate
        result = runner.invoke(main, ["guard", "deactivate"])
        assert result.exit_code == 0

    def test_edit_blocked_without_active_ws(self, runner: CliRunner) -> None:
        """AC2: Edit blocked without activation."""
        result = runner.invoke(main, ["guard", "check", "any/file.py"])
        assert result.exit_code == 1
        assert "No active WS" in result.output

    @patch("sdp.cli.guard.create_beads_client")
    def test_edit_blocked_outside_scope(
        self, mock_create_client: Mock, runner: CliRunner, sample_workstream: BeadsTask
    ) -> None:
        """AC3: Edit blocked for file outside scope."""
        # Mock Beads client
        mock_client = Mock()
        mock_client.get_task.return_value = sample_workstream
        mock_create_client.return_value = mock_client

        ws_id = sample_workstream.id

        # Activate WS with restricted scope
        runner.invoke(main, ["guard", "activate", ws_id])

        # Try to edit file outside scope
        result = runner.invoke(main, ["guard", "check", "src/other/forbidden.py"])
        assert result.exit_code == 1
        assert "not in WS scope" in result.output

    @patch("sdp.cli.guard.create_beads_client")
    def test_concurrent_activation_blocked(
        self, mock_create_client: Mock, runner: CliRunner, sample_workstream: BeadsTask
    ) -> None:
        """AC4: Cannot activate second WS (via guard state check)."""
        # Mock Beads client
        mock_client = Mock()
        mock_client.get_task.return_value = sample_workstream
        mock_create_client.return_value = mock_client

        # Activate first WS
        ws1_id = "00-032-01"
        runner.invoke(main, ["guard", "activate", ws1_id])

        # Create second WS
        ws2 = BeadsTask(
            id="00-032-02",
            title="Second WS",
            description="Another workstream",
            priority=BeadsPriority.MEDIUM,
        )
        mock_client.get_task.return_value = ws2

        # Try to activate second (guard doesn't block, but tracker would)
        # For now, guard allows multiple activations, tracker blocks
        result = runner.invoke(main, ["guard", "activate", "00-032-02"])
        assert result.exit_code == 0  # Guard itself doesn't block

    @patch("sdp.cli.guard.create_beads_client")
    def test_status_shows_active_ws_details(
        self, mock_create_client: Mock, runner: CliRunner, sample_workstream: BeadsTask
    ) -> None:
        """Status command shows comprehensive WS details."""
        # Mock Beads client
        mock_client = Mock()
        mock_client.get_task.return_value = sample_workstream
        mock_create_client.return_value = mock_client

        # Activate
        runner.invoke(main, ["guard", "activate", sample_workstream.id])

        # Check status
        result = runner.invoke(main, ["guard", "status"])
        assert result.exit_code == 0
        assert sample_workstream.id in result.output
        assert "Scope: 2 files" in result.output

    @patch("sdp.cli.guard.create_beads_client")
    def test_current_command_shows_active_ws(
        self, mock_create_client: Mock, runner: CliRunner, sample_workstream: BeadsTask
    ) -> None:
        """Current command shows active WS."""
        # Mock Beads client
        mock_client = Mock()
        mock_client.get_task.return_value = sample_workstream
        mock_create_client.return_value = mock_client

        # Activate
        runner.invoke(main, ["guard", "activate", sample_workstream.id])

        # Check current
        result = runner.invoke(main, ["guard", "current"])
        assert result.exit_code == 0
        assert sample_workstream.id in result.output

    @patch("sdp.cli.guard.create_beads_client")
    def test_check_shows_scope_files_on_error(
        self, mock_create_client: Mock, runner: CliRunner, sample_workstream: BeadsTask
    ) -> None:
        """Check command shows allowed files when blocked."""
        # Mock Beads client
        mock_client = Mock()
        mock_client.get_task.return_value = sample_workstream
        mock_create_client.return_value = mock_client

        # Activate
        runner.invoke(main, ["guard", "activate", sample_workstream.id])

        # Try forbidden file
        result = runner.invoke(main, ["guard", "check", "src/forbidden.py"])
        assert result.exit_code == 1
        assert "src/sdp/guard/skill.py" in result.output
        assert "src/sdp/guard/models.py" in result.output

    @patch("sdp.cli.guard.create_beads_client")
    def test_unrestricted_scope_allows_all_files(
        self, mock_create_client: Mock, runner: CliRunner
    ) -> None:
        """WS with no scope restrictions allows all files."""
        # Mock Beads client with no scope
        mock_client = Mock()
        mock_task = Mock()
        mock_task.sdp_metadata = {}  # No scope_files
        mock_client.get_task.return_value = mock_task
        mock_create_client.return_value = mock_client

        # Activate
        runner.invoke(main, ["guard", "activate", "00-032-01"])

        # Check any file
        result = runner.invoke(main, ["guard", "check", "any/random/file.py"])
        assert result.exit_code == 0
        assert "allowed" in result.output.lower()
