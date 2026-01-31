"""Tests for sdp.cli.status.command."""

from pathlib import Path
from unittest.mock import MagicMock, patch

from click.testing import CliRunner

from sdp.cli.status.command import status
from sdp.cli.status.models import BeadsStatus, GuardStatus, ProjectStatus


def test_status_command_default_output() -> None:
    """Test status command with default human-readable output."""
    runner = CliRunner()
    mock_status = ProjectStatus(
        in_progress=[],
        blocked=[],
        ready=[],
        guard=GuardStatus(active=False),
        beads=BeadsStatus(available=False, synced=False),
        next_actions=[],
    )
    
    with patch("sdp.cli.status.command.StatusCollector") as mock_collector:
        mock_collector.return_value.collect.return_value = mock_status
        result = runner.invoke(status, [])
    
    assert result.exit_code == 0
    # Should use human format by default


def test_status_command_json_output() -> None:
    """Test status command with JSON output."""
    runner = CliRunner()
    mock_status = ProjectStatus(
        in_progress=[],
        blocked=[],
        ready=[],
        guard=GuardStatus(active=False),
        beads=BeadsStatus(available=False, synced=False),
        next_actions=[],
    )
    
    with patch("sdp.cli.status.command.StatusCollector") as mock_collector:
        mock_collector.return_value.collect.return_value = mock_status
        with patch("sdp.cli.status.command.format_status_json") as mock_json:
            mock_json.return_value = '{"status": "ok"}'
            result = runner.invoke(status, ["--json"])
    
    assert result.exit_code == 0
    assert '{"status": "ok"}' in result.output


def test_status_command_verbose_output() -> None:
    """Test status command with verbose flag."""
    runner = CliRunner()
    mock_status = ProjectStatus(
        in_progress=[],
        blocked=[],
        ready=[],
        guard=GuardStatus(active=False),
        beads=BeadsStatus(available=False, synced=False),
        next_actions=[],
    )
    
    with patch("sdp.cli.status.command.StatusCollector") as mock_collector:
        mock_collector.return_value.collect.return_value = mock_status
        with patch("sdp.cli.status.command.format_status_human") as mock_human:
            mock_human.return_value = "Verbose output"
            result = runner.invoke(status, ["--verbose"])
            # Verify verbose=True was passed
            mock_human.assert_called_once()
            assert mock_human.call_args[1]["verbose"] is True
    
    assert result.exit_code == 0


def test_status_command_short_verbose_flag() -> None:
    """Test status command with -v short flag."""
    runner = CliRunner()
    mock_status = ProjectStatus(
        in_progress=[],
        blocked=[],
        ready=[],
        guard=GuardStatus(active=False),
        beads=BeadsStatus(available=False, synced=False),
        next_actions=[],
    )
    
    with patch("sdp.cli.status.command.StatusCollector") as mock_collector:
        mock_collector.return_value.collect.return_value = mock_status
        with patch("sdp.cli.status.command.format_status_human") as mock_human:
            mock_human.return_value = "Verbose output"
            result = runner.invoke(status, ["-v"])
            # Verify verbose=True was passed
            mock_human.assert_called_once()
            assert mock_human.call_args[1]["verbose"] is True
    
    assert result.exit_code == 0


def test_status_command_json_and_verbose() -> None:
    """Test status command with both --json and --verbose flags."""
    runner = CliRunner()
    mock_status = ProjectStatus(
        in_progress=[],
        blocked=[],
        ready=[],
        guard=GuardStatus(active=False),
        beads=BeadsStatus(available=False, synced=False),
        next_actions=[],
    )
    
    with patch("sdp.cli.status.command.StatusCollector") as mock_collector:
        mock_collector.return_value.collect.return_value = mock_status
        with patch("sdp.cli.status.command.format_status_json") as mock_json:
            mock_json.return_value = '{"status": "ok"}'
            result = runner.invoke(status, ["--json", "--verbose"])
    
    # JSON takes precedence
    assert result.exit_code == 0
    assert '{"status": "ok"}' in result.output
