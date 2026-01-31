"""Integration tests: hooks work on SDP repository itself (dogfooding)."""

from pathlib import Path

import pytest

from sdp.hooks.common import find_project_root, find_workstream_dir


def test_finds_sdp_project_root() -> None:
    """Verify detector finds SDP repo root correctly."""
    sdp_root = find_project_root(Path(__file__).resolve().parents[3])
    assert (sdp_root / "docs" / "workstreams").exists()
    assert (sdp_root / "pyproject.toml").exists()


def test_find_workstream_dir_on_sdp() -> None:
    """Verify find_workstream_dir returns docs/workstreams on SDP repo."""
    sdp_root = find_project_root(Path(__file__).resolve().parents[3])
    ws_dir = find_workstream_dir(sdp_root)
    assert ws_dir == sdp_root / "docs" / "workstreams"
    assert ws_dir.exists()


@pytest.mark.integration
def test_pre_commit_runs_on_sdp() -> None:
    """Verify pre-commit hook runs on SDP codebase without external deps."""
    from sdp.hooks.pre_commit import main

    exit_code = main()
    assert exit_code == 0
