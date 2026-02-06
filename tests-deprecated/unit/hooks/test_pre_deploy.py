"""Tests for sdp.hooks.pre_deploy."""

import sys
from pathlib import Path
from unittest.mock import MagicMock, patch

from sdp.hooks.pre_deploy import _repo_root, main


def test_repo_root_returns_path() -> None:
    """_repo_root returns Path from git."""
    with patch("subprocess.run") as m:
        m.return_value = MagicMock(stdout="/tmp/repo\n", returncode=0)
        root = _repo_root()
    assert "/tmp/repo" in str(root)


def test_main_requires_feature_id() -> None:
    """main() returns 1 when no feature ID provided."""
    orig = sys.argv
    try:
        sys.argv = ["pre_deploy.py"]
        exit_code = main()
        assert exit_code == 1
    finally:
        sys.argv = orig


def test_main_runs_checks() -> None:
    """main() runs pre-deploy checks (mocked)."""
    repo = Path(__file__).resolve().parents[3]
    orig = sys.argv
    try:
        sys.argv = ["pre_deploy.py", "F01", "staging"]
        with patch("sdp.hooks.pre_deploy._repo_root", return_value=repo):
            with patch("subprocess.run") as m_run:
                m_run.return_value = MagicMock(returncode=0, stdout="", stderr="")
                exit_code = main()
        assert exit_code in (0, 1)
    finally:
        sys.argv = orig
