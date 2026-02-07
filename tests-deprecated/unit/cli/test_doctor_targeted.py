"""Targeted tests for uncovered lines in cli/doctor.py."""

import subprocess
from pathlib import Path
from unittest.mock import MagicMock, patch

import pytest

from sdp.cli.doctor import check_environment, check_project_structure, check_workstreams


def test_check_environment_old_python() -> None:
    """Test check_environment with Python < 3.10 (lines 25-27)."""
    import sys
    from collections import namedtuple

    # Create a version_info-like tuple
    VersionInfo = namedtuple("VersionInfo", ["major", "minor", "micro"])
    old_version = VersionInfo(3, 9, 0)

    with patch.object(sys, "version_info", old_version):
        result = check_environment()
        assert not result  # Should fail


def test_check_environment_poetry_fails() -> None:
    """Test check_environment with Poetry command failure (lines 40-41)."""
    with patch("subprocess.run") as mock_run:
        # First call (poetry): returncode != 0
        # Second call (git): success
        mock_run.side_effect = [
            MagicMock(returncode=1, stdout="", stderr=""),
            MagicMock(returncode=0, stdout="git version 2.0", stderr=""),
        ]
        result = check_environment()
        assert not result


def test_check_environment_poetry_exception() -> None:
    """Test check_environment with Poetry exception (lines 42-44)."""
    with patch("subprocess.run") as mock_run:
        # First call (poetry): exception
        # Second call (git): success
        mock_run.side_effect = [
            FileNotFoundError(),
            MagicMock(returncode=0, stdout="git version 2.0", stderr=""),
        ]
        result = check_environment()
        assert not result


def test_check_environment_git_fails() -> None:
    """Test check_environment with Git command failure (lines 58-61)."""
    with patch("subprocess.run") as mock_run:
        # First call (poetry): success
        # Second call (git): returncode != 0
        mock_run.side_effect = [
            MagicMock(returncode=0, stdout="Poetry version 1.0", stderr=""),
            MagicMock(returncode=1, stdout="", stderr=""),
        ]
        result = check_environment()
        assert not result


def test_check_environment_git_exception() -> None:
    """Test check_environment with Git exception (lines 59-61)."""
    with patch("subprocess.run") as mock_run:
        # First call (poetry): success
        # Second call (git): exception
        mock_run.side_effect = [
            MagicMock(returncode=0, stdout="Poetry version 1.0", stderr=""),
            FileNotFoundError(),
        ]
        result = check_environment()
        assert not result


def test_check_workstreams_generic_exception(tmp_path: Path) -> None:
    """Test check_workstreams handles generic exceptions (lines 136-139)."""
    from sdp.core import parse_workstream

    ws_dir = tmp_path / "docs" / "workstreams"
    ws_dir.mkdir(parents=True)

    # Create invalid WS file that will cause generic error
    ws_file = ws_dir / "WS-001-01.md"
    ws_file.write_text("invalid content")

    with patch("sdp.core.parse_workstream") as mock_parse:
        mock_parse.side_effect = RuntimeError("Generic error")
        result = check_workstreams(tmp_path)
        assert not result
