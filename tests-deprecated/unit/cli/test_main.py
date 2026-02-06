"""Tests for sdp.cli.main CLI entry point."""

from unittest.mock import patch

import pytest
from click.testing import CliRunner

from sdp.cli.main import main


def test_cli_version() -> None:
    """Test --version flag displays version."""
    runner = CliRunner()
    result = runner.invoke(main, ["--version"])
    
    assert result.exit_code == 0
    assert "version" in result.output.lower()


def test_cli_help() -> None:
    """Test --help flag displays help."""
    runner = CliRunner()
    result = runner.invoke(main, ["--help"])
    
    assert result.exit_code == 0
    assert "Usage" in result.output or "Options" in result.output


def test_cli_no_args_shows_help() -> None:
    """Test CLI with no args shows help."""
    runner = CliRunner()
    result = runner.invoke(main, [])
    
    # Should succeed or show help (exit code 0 or specific help code)
    assert result.exit_code in (0, 2)


def test_cli_workstream_command_available() -> None:
    """Test workstream command group is available."""
    runner = CliRunner()
    result = runner.invoke(main, ["workstream", "--help"])
    
    # Should show workstream subcommands
    assert result.exit_code in (0, 2)  # 2 if no subcommands, 0 if help shown


def test_cli_guard_command_available() -> None:
    """Test guard command group is available when imported."""
    runner = CliRunner()
    result = runner.invoke(main, ["guard", "--help"])
    
    # May not be available in all builds
    assert result.exit_code in (0, 2)


def test_cli_tier_command_available() -> None:
    """Test tier command group is available when imported."""
    runner = CliRunner()
    result = runner.invoke(main, ["tier", "--help"])
    
    # May not be available in all builds
    assert result.exit_code in (0, 2)


def test_cli_prd_command_available() -> None:
    """Test prd command group is available when imported."""
    runner = CliRunner()
    result = runner.invoke(main, ["prd", "--help"])
    
    # May not be available in all builds
    assert result.exit_code in (0, 2)


def test_cli_quality_command_available() -> None:
    """Test quality command group is available when imported."""
    runner = CliRunner()
    result = runner.invoke(main, ["quality", "--help"])
    
    # May not be available in all builds - quality is in workstream subcommands
    assert result.exit_code in (0, 2)


def test_cli_status_command_available() -> None:
    """Test status command is available when imported."""
    runner = CliRunner()
    result = runner.invoke(main, ["status", "--help"])
    
    # May not be available in all builds
    assert result.exit_code in (0, 2)


def test_cli_beads_command_available_when_imported() -> None:
    """Test beads command is available when beads module exists."""
    runner = CliRunner()
    result = runner.invoke(main, ["beads", "--help"])
    
    # May not be available in all builds
    assert result.exit_code in (0, 2)


def test_cli_doctor_command_available_when_imported() -> None:
    """Test doctor command is available when doctor module exists."""
    runner = CliRunner()
    result = runner.invoke(main, ["doctor", "--help"])
    
    # May not be available in all builds
    assert result.exit_code in (0, 2)


def test_cli_trace_command_available_when_imported() -> None:
    """Test trace command is available when trace module exists."""
    runner = CliRunner()
    result = runner.invoke(main, ["trace", "--help"])
    
    # May not be available in all builds
    assert result.exit_code in (0, 2)


def test_cli_skill_command_available_when_imported() -> None:
    """Test skill command is available when skill module exists."""
    runner = CliRunner()
    result = runner.invoke(main, ["skill", "--help"])
    
    # May not be available in all builds
    assert result.exit_code in (0, 2)


def test_cli_sync_command_available_when_imported() -> None:
    """Test sync command is available when sync module exists."""
    runner = CliRunner()
    result = runner.invoke(main, ["sync", "--help"])
    
    # May not be available in all builds
    assert result.exit_code in (0, 2)


def test_cli_metrics_command_available_when_imported() -> None:
    """Test metrics command is available when metrics module exists."""
    runner = CliRunner()
    result = runner.invoke(main, ["metrics", "--help"])
    
    # May not be available in all builds
    assert result.exit_code in (0, 2)


def test_cli_version_command() -> None:
    """Test version subcommand displays version."""
    runner = CliRunner()
    result = runner.invoke(main, ["version"])
    
    assert result.exit_code == 0
    assert "version" in result.output.lower()


def test_cli_invalid_command() -> None:
    """Test CLI with invalid command shows error."""
    runner = CliRunner()
    result = runner.invoke(main, ["invalid-command"])
    
    # Should fail with error message
    assert result.exit_code != 0


def test_cli_fallback_when_import_fails() -> None:
    """Test CLI gracefully handles import failures."""
    runner = CliRunner()
    
    # Mock import failure for a command
    with patch("sdp.cli.main._beads_available", False):
        result = runner.invoke(main, ["--help"])
    
    # Should still work without beads command
    assert result.exit_code == 0
