"""Tests for sdp.hooks.pre_commit."""

from pathlib import Path
from unittest.mock import MagicMock, patch

from sdp.hooks.pre_commit_checks import (
    check_branch,
    check_python_bare_except,
    check_tech_debt,
    check_time_estimates,
    check_ws_format,
    repo_root,
    run_script,
    staged_files,
)
from sdp.hooks.common import CheckResult


def test_check_branch_rejects_main() -> None:
    """check_branch rejects main branch."""
    with patch("subprocess.run") as m:
        m.return_value = MagicMock(stdout="main\n", returncode=0)
        result = check_branch()
    assert not result.passed
    assert "main" in result.message


def test_check_branch_accepts_feature_branch() -> None:
    """check_branch accepts feature branch."""
    with patch("subprocess.run") as m:
        m.return_value = MagicMock(stdout="feature/add-auth\n", returncode=0)
        result = check_branch()
    assert result.passed
    assert "feature/add-auth" in result.message


def test_check_time_estimates_no_ws_files() -> None:
    """check_time_estimates passes when no WS files staged."""
    result = check_time_estimates(["src/foo.py", "README.md"])
    assert result.passed
    assert "No WS files" in result.message


def test_check_time_estimates_allows_relative_sizing() -> None:
    """check_time_estimates allows SMALL/MEDIUM/LARGE."""
    with patch("subprocess.run") as m:
        m.return_value = MagicMock(
            stdout="+size: MEDIUM\n+scope: SMALL",
            returncode=0,
        )
        result = check_time_estimates(["docs/workstreams/backlog/00-001-01-foo.md"])
    # The check looks for day/hour/week patterns - MEDIUM/SMALL don't match
    assert result.passed


def test_check_tech_debt_no_code_files() -> None:
    """check_tech_debt passes when no code files staged."""
    result = check_tech_debt(["hooks/pre-commit.sh"])
    assert result.passed


def test_check_tech_debt_rejects_tech_debt_marker() -> None:
    """check_tech_debt rejects tech debt markers."""
    with patch("subprocess.run") as m:
        m.return_value = MagicMock(
            stdout="+  # fix tech debt later\n",
            returncode=0,
        )
        result = check_tech_debt(["src/foo.py"])
    assert not result.passed
    assert "Tech debt" in result.message


def test_check_tech_debt_allows_no_tech_debt_rule() -> None:
    """check_tech_debt allows 'No Tech Debt' rule in docs."""
    with patch("subprocess.run") as m:
        m.return_value = MagicMock(
            stdout="+  ## No Tech Debt allowed\n",
            returncode=0,
        )
        result = check_tech_debt(["docs/rules.md"])
    # The exclude pattern should match "no.?tech.?debt"
    assert result.passed


def test_check_python_bare_except_no_py_files() -> None:
    """check_python_bare_except passes when no Python files."""
    result = check_python_bare_except(["README.md"])
    assert result.passed


def test_check_python_bare_except_rejects_bare_except() -> None:
    """check_python_bare_except rejects bare except."""
    with patch("subprocess.run") as m:
        m.return_value = MagicMock(
            stdout="+    except:\n+        pass\n",
            returncode=0,
        )
        result = check_python_bare_except(["src/foo.py"])
    assert not result.passed


def test_check_ws_format_no_new_ws() -> None:
    """check_ws_format passes when no new WS files."""
    result = check_ws_format(["src/foo.py"], Path("."))
    assert result.passed


def test_check_ws_format_valid_ws() -> None:
    """check_ws_format accepts valid WS with Goal and AC."""
    with patch("subprocess.run") as m:
        m.return_value = MagicMock(
            stdout="### 🎯 Goal\n\n**Acceptance Criteria:**\n- [ ] AC1",
            returncode=0,
        )
        result = check_ws_format(
            ["docs/workstreams/backlog/00-001-01-foo.md"],
            Path("/tmp"),
        )
    assert result.passed


def test_main_passes_when_no_staged_files() -> None:
    """main() returns 0 when no staged files."""
    from sdp.hooks.pre_commit import main

    with patch("sdp.hooks.pre_commit.repo_root") as m_root:
        m_root.return_value = Path("/tmp/repo")
        with patch("sdp.hooks.pre_commit.staged_files") as m_files:
            m_files.return_value = []
            exit_code = main()
    assert exit_code == 0


def test_main_fails_on_branch_check() -> None:
    """main() returns 1 when branch check fails."""
    from sdp.hooks.pre_commit import main

    with patch("sdp.hooks.pre_commit.repo_root") as m_root:
        m_root.return_value = Path("/tmp/repo")
        with patch("sdp.hooks.pre_commit.staged_files") as m_files:
            m_files.return_value = ["src/foo.py"]
            with patch("sdp.hooks.pre_commit.check_branch") as m_branch:
                m_branch.return_value = MagicMock(
                    passed=False,
                    message="main",
                    format_terminal=lambda: "❌ main",
                    violations=[],
                )
                exit_code = main()
    assert exit_code == 1


def test_repo_root_returns_path() -> None:
    """repo_root returns Path from git."""
    with patch("subprocess.run") as m:
        m.return_value = MagicMock(stdout="/tmp/repo\n", returncode=0)
        root = repo_root()
    assert "/tmp/repo" in str(root)


def test_staged_files_returns_list() -> None:
    """staged_files returns list of staged file paths."""
    with patch("subprocess.run") as m:
        m.return_value = MagicMock(stdout="src/foo.py\nREADME.md\n", returncode=0)
        files = staged_files()
    assert "src/foo.py" in files
    assert "README.md" in files


def test_run_script_returns_success_when_exists() -> None:
    """run_script returns (True, output) when script succeeds."""
    repo = Path(__file__).resolve().parents[3]
    script = "scripts/check_workstreams_layout.py"
    with patch("subprocess.run") as m:
        m.return_value = MagicMock(returncode=0, stdout="ok", stderr="")
        ok, out = run_script(repo, script, [])
    assert ok is True
    assert "ok" in out


def test_run_script_skips_when_not_exists() -> None:
    """run_script returns (True, msg) when script not found."""
    ok, out = run_script(Path("/tmp"), "nonexistent/script.py", [])
    assert ok is True
    assert "not found" in out or "skipped" in out


def test_main_invoked_as_module() -> None:
    """main() runs when invoked as python -m sdp.hooks.pre_commit (no staged files)."""
    import subprocess
    import sys
    from pathlib import Path
    repo = Path(__file__).resolve().parents[3]
    proc = subprocess.run(
        [sys.executable, "-m", "sdp.hooks.pre_commit"],
        cwd=repo,
        capture_output=True,
        text=True,
        timeout=5,
    )
    assert proc.returncode == 0
    assert "Pre-commit checks" in proc.stdout
    assert "No staged files" in proc.stdout or "staged" in proc.stdout.lower()


def test_main_passes_all_checks() -> None:
    """main() returns 0 when all checks pass (mocked)."""
    from sdp.hooks.pre_commit import main

    passed = MagicMock(passed=True, format_terminal=lambda: "✓ ok", violations=[])
    with patch("sdp.hooks.pre_commit.repo_root") as m_root:
        m_root.return_value = Path("/tmp/repo")
        with patch("sdp.hooks.pre_commit.staged_files") as m_files:
            m_files.return_value = ["README.md"]  # No workstreams, no src py
            with patch("sdp.hooks.pre_commit.check_branch", return_value=passed):
                with patch("sdp.hooks.pre_commit.check_time_estimates", return_value=passed):
                    with patch("sdp.hooks.pre_commit.check_tech_debt", return_value=passed):
                        with patch("sdp.hooks.pre_commit.check_python_bare_except", return_value=passed):
                            with patch("sdp.hooks.pre_commit.check_ws_format", return_value=passed):
                                exit_code = main()
    assert exit_code == 0
