"""Tests for sdp.hooks.pre_push."""

from pathlib import Path
from unittest.mock import MagicMock, patch

from sdp.hooks.pre_push import _repo_root, _files_to_push


def test_repo_root_returns_path() -> None:
    """_repo_root returns Path from git."""
    with patch("subprocess.run") as m:
        m.return_value = MagicMock(stdout="/tmp/repo\n", returncode=0)
        root = _repo_root()
    assert "/tmp/repo" in str(root)


def test_files_to_push_handles_no_upstream() -> None:
    """_files_to_push falls back to HEAD~1 HEAD when no upstream."""
    with patch("subprocess.run") as m:
        m.side_effect = [
            MagicMock(returncode=1, stdout=""),  # HEAD @{u} fails
            MagicMock(stdout="src/foo.py\n", returncode=0),
        ]
        files = _files_to_push()
    assert "src/foo.py" in files


def test_main_skips_when_no_py_files() -> None:
    """main() returns 0 when no Python files to push."""
    from sdp.hooks.pre_push import main

    with patch("sdp.hooks.pre_push._repo_root") as m_root:
        m_root.return_value = Path("/tmp/repo")
        with patch("sdp.hooks.pre_push._files_to_push") as m_files:
            m_files.return_value = ["README.md", "docs/foo.md"]
            exit_code = main()
    assert exit_code == 0


def test_main_coverage_check_branch() -> None:
    """main() runs coverage check when .coverage exists."""
    from sdp.hooks.pre_push import main

    repo = Path(__file__).resolve().parents[3]
    (repo / ".coverage").touch()
    try:
        with patch("sdp.hooks.pre_push._repo_root") as m_root:
            m_root.return_value = repo
            with patch("sdp.hooks.pre_push._files_to_push") as m_files:
                m_files.return_value = ["src/foo.py"]
                with patch("subprocess.run") as m_run:
                    m_run.return_value = MagicMock(
                        returncode=0, stdout="TOTAL 100 85 85%"
                    )
                    exit_code = main()
        assert exit_code == 0
    finally:
        (repo / ".coverage").unlink(missing_ok=True)


def test_main_runs_tests_when_py_files() -> None:
    """main() runs pytest when Python files to push."""
    from sdp.hooks.pre_push import main

    repo = Path(__file__).resolve().parents[3]
    with patch("sdp.hooks.pre_push._repo_root") as m_root:
        m_root.return_value = repo
        with patch("sdp.hooks.pre_push._files_to_push") as m_files:
            m_files.return_value = ["src/foo.py"]
            with patch("subprocess.run") as m_run:
                m_run.return_value = MagicMock(returncode=0, stdout="")
                exit_code = main()
    assert exit_code == 0
