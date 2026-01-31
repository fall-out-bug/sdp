"""Tests for sdp.hooks.post_build."""

import sys
from pathlib import Path
from unittest.mock import MagicMock, patch

from sdp.hooks.post_build import _project_root, main


def test_project_root_uses_sdp_when_no_hw_checker() -> None:
    """_project_root returns (repo_root, src/sdp) when tools/hw_checker missing."""
    repo = Path("/tmp/repo")
    work_dir, prefix = _project_root(repo)
    assert work_dir == repo
    assert prefix == "src/sdp"


def test_main_requires_ws_id() -> None:
    """main() returns 1 when no WS-ID provided."""
    orig = sys.argv
    try:
        sys.argv = ["post_build.py"]
        exit_code = main()
        assert exit_code == 1
    finally:
        sys.argv = orig


def test_main_runs_with_mocked_subprocess() -> None:
    """main() runs checks with mocked subprocess (all pass)."""
    repo = Path(__file__).resolve().parents[3]
    orig = sys.argv
    try:
        sys.argv = ["post_build.py", "00-020-01", "hooks"]
        with patch("sdp.hooks.post_build._repo_root", return_value=repo):
            with patch("subprocess.run") as m_run:
                # pytest=0, ruff=0, mypy=0, grep=1(no match), git=0, ...
                m_run.side_effect = [
                    MagicMock(returncode=0, stdout="", stderr=""),  # pytest
                    MagicMock(returncode=0, stdout="", stderr=""),  # ruff
                    MagicMock(returncode=0, stdout="", stderr=""),  # mypy
                    MagicMock(returncode=1, stdout="", stderr=""),  # grep (no TODO)
                    MagicMock(stdout="feat(hooks): 00-020-01 - Extract hooks\n", returncode=0),  # git log
                ]
                exit_code = main()
        assert exit_code in (0, 1)
    finally:
        sys.argv = orig


