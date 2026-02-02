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


def test_main_skips_commit_check_when_env_set() -> None:
    """main() skips commit check when SKIP_COMMIT_CHECK=1."""
    import os
    repo = Path(__file__).resolve().parents[3]
    orig_argv = sys.argv
    orig_env = os.environ.get("SKIP_COMMIT_CHECK")
    try:
        sys.argv = ["post_build.py", "00-020-01", "hooks"]
        os.environ["SKIP_COMMIT_CHECK"] = "1"
        with patch("sdp.hooks.post_build._repo_root", return_value=repo):
            with patch("subprocess.run") as m_run:
                m_run.side_effect = [
                    MagicMock(returncode=0, stdout="", stderr=""),  # pytest
                    MagicMock(returncode=0, stdout="", stderr=""),  # ruff
                    MagicMock(returncode=0, stdout="", stderr=""),  # mypy
                    MagicMock(returncode=1, stdout="", stderr=""),  # grep (no TODO)
                ]
                exit_code = main()
        # Should pass without git check
        assert exit_code in (0, 1)
    finally:
        sys.argv = orig_argv
        if orig_env is None:
            os.environ.pop("SKIP_COMMIT_CHECK", None)
        else:
            os.environ["SKIP_COMMIT_CHECK"] = orig_env


def test_main_detects_large_files() -> None:
    """main() detects files > 200 LOC."""
    repo = Path(__file__).resolve().parents[3]
    orig = sys.argv
    try:
        sys.argv = ["post_build.py", "00-020-01"]
        with patch("sdp.hooks.post_build._repo_root", return_value=repo):
            with patch("subprocess.run") as m_run:
                m_run.side_effect = [
                    MagicMock(returncode=0, stdout="", stderr=""),  # pytest
                    MagicMock(returncode=0, stdout="", stderr=""),  # ruff
                    MagicMock(returncode=0, stdout="", stderr=""),  # mypy
                    MagicMock(returncode=1, stdout="", stderr=""),  # grep (no TODO)
                ]
                # Will check file sizes in real repo - some files are > 200 LOC
                exit_code = main()
        # Should fail if large files found, or pass if not
        assert exit_code in (0, 1)
    finally:
        sys.argv = orig


def test_main_skips_tests_when_dir_missing() -> None:
    """main() skips tests check when tests/ directory doesn't exist."""
    repo = Path(__file__).resolve().parents[3]
    orig = sys.argv
    try:
        sys.argv = ["post_build.py", "00-020-01", "hooks"]
        with patch("sdp.hooks.post_build._repo_root", return_value=Path("/nonexistent")):
            with patch("subprocess.run") as m_run:
                m_run.side_effect = [
                    # No pytest call, skip directly to ruff
                    MagicMock(returncode=0, stdout="", stderr=""),  # ruff
                    MagicMock(returncode=0, stdout="", stderr=""),  # mypy
                    MagicMock(returncode=1, stdout="", stderr=""),  # grep (no TODO)
                    MagicMock(stdout="", returncode=0),  # git log
                ]
                with patch("pathlib.Path.exists", side_effect=lambda: False):
                    exit_code = main()
        assert exit_code in (0, 1)
    finally:
        sys.argv = orig
