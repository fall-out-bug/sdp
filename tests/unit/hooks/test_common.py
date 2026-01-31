"""Tests for sdp.hooks.common."""

from pathlib import Path

import pytest

from sdp.hooks.common import CheckResult, find_project_root, find_workstream_dir


def test_find_project_root_via_sdp_root_marker(tmp_path: Path) -> None:
    """find_project_root returns path with .sdp-root marker."""
    (tmp_path / ".sdp-root").touch()
    assert find_project_root(tmp_path) == tmp_path
    assert find_project_root(tmp_path / "sub" / "dir") == tmp_path


def test_find_project_root_via_docs_workstreams(tmp_path: Path) -> None:
    """find_project_root returns path with docs/workstreams directory."""
    (tmp_path / "docs" / "workstreams").mkdir(parents=True)
    assert find_project_root(tmp_path) == tmp_path
    assert find_project_root(tmp_path / "sub") == tmp_path


def test_find_project_root_via_git_and_pyproject(tmp_path: Path) -> None:
    """find_project_root returns path with .git and [tool.sdp] in pyproject."""
    (tmp_path / ".git").mkdir()  # .git is typically a directory
    (tmp_path / "pyproject.toml").write_text(
        '[tool.sdp]\nversion = "0.1"\n', encoding="utf-8"
    )
    assert find_project_root(tmp_path) == tmp_path


def test_find_project_root_skips_invalid_toml(tmp_path: Path) -> None:
    """find_project_root skips path when pyproject.toml has invalid TOML."""
    (tmp_path / ".git").mkdir()
    (tmp_path / "pyproject.toml").write_text("invalid toml [", encoding="utf-8")
    with pytest.raises(RuntimeError, match="SDP project root not found"):
        find_project_root(tmp_path)


def test_find_project_root_raises_when_not_found(tmp_path: Path) -> None:
    """find_project_root raises RuntimeError when no markers found."""
    with pytest.raises(RuntimeError, match="SDP project root not found"):
        find_project_root(tmp_path)


def test_find_workstream_dir_from_quality_gate(tmp_path: Path) -> None:
    """find_workstream_dir uses quality-gate.toml [workstreams.dir]."""
    (tmp_path / "custom_ws").mkdir()
    (tmp_path / "quality-gate.toml").write_text(
        '[workstreams]\ndir = "custom_ws"', encoding="utf-8"
    )
    assert find_workstream_dir(tmp_path) == tmp_path / "custom_ws"


def test_find_workstream_dir_falls_through_on_invalid_toml(tmp_path: Path) -> None:
    """find_workstream_dir falls through to default when quality-gate.toml invalid."""
    (tmp_path / "docs" / "workstreams").mkdir(parents=True)
    (tmp_path / "quality-gate.toml").write_text("invalid toml [", encoding="utf-8")
    assert find_workstream_dir(tmp_path) == tmp_path / "docs" / "workstreams"


def test_find_workstream_dir_from_env(tmp_path: Path, monkeypatch: pytest.MonkeyPatch) -> None:
    """find_workstream_dir uses SDP_WORKSTREAM_DIR when set."""
    (tmp_path / "docs" / "workstreams").mkdir(parents=True)
    (tmp_path / "env_ws").mkdir()
    monkeypatch.setenv("SDP_WORKSTREAM_DIR", str(tmp_path / "env_ws"))
    assert find_workstream_dir(tmp_path) == tmp_path / "env_ws"


def test_find_workstream_dir_default_docs_workstreams(tmp_path: Path) -> None:
    """find_workstream_dir defaults to docs/workstreams."""
    (tmp_path / "docs" / "workstreams").mkdir(parents=True)
    assert find_workstream_dir(tmp_path) == tmp_path / "docs" / "workstreams"


def test_find_workstream_dir_legacy_fallback(tmp_path: Path) -> None:
    """find_workstream_dir falls back to workstreams/ when docs/workstreams missing."""
    (tmp_path / "workstreams").mkdir()
    assert find_workstream_dir(tmp_path) == tmp_path / "workstreams"


def test_find_workstream_dir_raises_when_not_found(tmp_path: Path) -> None:
    """find_workstream_dir raises when no workstream dir exists."""
    with pytest.raises(RuntimeError, match="Workstream directory not found"):
        find_workstream_dir(tmp_path)


def test_check_result_passed_format() -> None:
    """CheckResult formats passed result correctly."""
    result = CheckResult(passed=True, message="All good", violations=[])
    assert "✓" in result.format_terminal()
    assert "All good" in result.format_terminal()


def test_check_result_failed_format() -> None:
    """CheckResult formats failed result with violations."""
    result = CheckResult(
        passed=False,
        message="Time estimates found",
        violations=[
            (Path("foo.py"), 10, "2 hours"),
            (Path("bar.md"), None, "soon"),
        ],
    )
    output = result.format_terminal()
    assert "❌" in output
    assert "Time estimates found" in output
    assert "foo.py:10" in output
    assert "2 hours" in output
    assert "bar.md:?" in output
    assert "soon" in output
