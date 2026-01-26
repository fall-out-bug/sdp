"""Tests for PreExecutionChecker."""

from pathlib import Path

import pytest

from sdp.agents.pre_check import PreExecutionChecker


def test_pre_check_file_not_found() -> None:
    checker = PreExecutionChecker()
    errors = checker.check("00-012-99")

    assert len(errors) > 0
    assert "not found" in errors[0].lower()


def test_pre_check_missing_dependency(tmp_path: Path) -> None:
    # Create workstream with dependency
    backlog_dir = tmp_path / "docs" / "workstreams" / "backlog"
    backlog_dir.mkdir(parents=True)
    ws_file = backlog_dir / "00-012-01.md"
    ws_file.write_text("""---
dependencies:
  - 00-012-99
status: backlog
---
# Test
""")

    import os
    original_cwd = os.getcwd()
    try:
        os.chdir(tmp_path)
        checker = PreExecutionChecker()
        errors = checker.check("00-012-01")

        assert any("00-012-99" in e and "not completed" in e.lower() for e in errors)
    finally:
        os.chdir(original_cwd)


def test_pre_check_circular_dependency(tmp_path: Path) -> None:
    # Create two workstreams that depend on each other
    backlog_dir = tmp_path / "docs" / "workstreams" / "backlog"
    backlog_dir.mkdir(parents=True)

    ws1 = backlog_dir / "00-012-01.md"
    ws1.write_text("""---
dependencies:
  - 00-012-02
status: backlog
---
# WS1
""")

    ws2 = backlog_dir / "00-012-02.md"
    ws2.write_text("""---
dependencies:
  - 00-012-01
status: backlog
---
# WS2
""")

    import os
    original_cwd = os.getcwd()
    try:
        os.chdir(tmp_path)
        checker = PreExecutionChecker()
        errors = checker.check("00-012-01")

        assert any("circular" in e.lower() for e in errors)
    finally:
        os.chdir(original_cwd)


def test_pre_check_large_size(tmp_path: Path) -> None:
    backlog_dir = tmp_path / "docs" / "workstreams" / "backlog"
    backlog_dir.mkdir(parents=True)
    ws_file = backlog_dir / "00-012-01.md"
    ws_file.write_text("""---
size: LARGE
status: backlog
---
# Large workstream
""")

    import os
    original_cwd = os.getcwd()
    try:
        os.chdir(tmp_path)
        checker = PreExecutionChecker()
        errors = checker.check("00-012-01")

        assert any("too large" in e.lower() for e in errors)
    finally:
        os.chdir(original_cwd)


def test_pre_check_can_execute(tmp_path: Path) -> None:
    backlog_dir = tmp_path / "docs" / "workstreams" / "backlog"
    backlog_dir.mkdir(parents=True)
    ws_file = backlog_dir / "00-012-01.md"
    ws_file.write_text("""---
size: SMALL
status: backlog
---
# Valid workstream
""")

    import os
    original_cwd = os.getcwd()
    try:
        os.chdir(tmp_path)
        checker = PreExecutionChecker()
        assert checker.can_execute("00-012-01") is True
    finally:
        os.chdir(original_cwd)


def test_pre_check_get_size(tmp_path: Path) -> None:
    backlog_dir = tmp_path / "docs" / "workstreams" / "backlog"
    backlog_dir.mkdir(parents=True)
    ws_file = backlog_dir / "00-012-01.md"
    ws_file.write_text("---\nsize: TINY\n---\n# Test")

    import os
    original_cwd = os.getcwd()
    try:
        os.chdir(tmp_path)
        checker = PreExecutionChecker()
        size = checker._get_size(ws_file)
        assert size == "TINY"
    finally:
        os.chdir(original_cwd)
