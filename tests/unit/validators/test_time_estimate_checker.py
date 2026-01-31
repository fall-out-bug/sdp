"""Tests for time_estimate_checker validator."""

import sys
from pathlib import Path
from unittest.mock import patch

import pytest

from sdp.validators.time_estimate_checker import (
    Violation,
    check_directory,
    check_file,
    format_violations,
)


class TestCheckFile:
    """Tests for check_file function."""

    def test_empty_file_returns_no_violations(self, tmp_path: Path) -> None:
        """Empty file has no violations."""
        f = tmp_path / "empty.md"
        f.write_text("")
        assert check_file(f) == []

    def test_file_not_found_returns_empty(self) -> None:
        """Non-existent file returns empty list."""
        assert check_file(Path("/nonexistent/file.md")) == []

    def test_detects_estimated_duration(self, tmp_path: Path) -> None:
        """Detects estimated_duration pattern."""
        f = tmp_path / "ws.md"
        f.write_text("estimated_duration: \"2-3 hours\"")
        violations = check_file(f)
        assert len(violations) == 1
        assert violations[0].pattern
        assert "estimated_duration" in violations[0].message

    def test_detects_estimated_loc(self, tmp_path: Path) -> None:
        """Detects estimated_loc pattern."""
        f = tmp_path / "ws.md"
        f.write_text("estimated_loc: 450")
        violations = check_file(f)
        assert len(violations) >= 1

    def test_detects_hours_pattern(self, tmp_path: Path) -> None:
        """Detects 2-3 hours pattern."""
        f = tmp_path / "ws.md"
        f.write_text("Takes 2-3 hours to complete")
        violations = check_file(f)
        assert len(violations) >= 1

    def test_allowed_context_skipped(self, tmp_path: Path) -> None:
        """Lines documenting violations are skipped."""
        f = tmp_path / "ws.md"
        f.write_text("❌ estimated_duration is forbidden")
        assert check_file(f) == []

    def test_format_violations_empty(self) -> None:
        """format_violations with empty list."""
        assert "No time estimate violations" in format_violations([])

    def test_format_violations_with_violations(self, tmp_path: Path) -> None:
        """format_violations formats violations."""
        f = tmp_path / "ws.md"
        f.write_text("estimated_duration: 2h")
        violations = check_file(f)
        result = format_violations(violations)
        assert "violation" in result.lower()
        assert str(f) in result

    def test_read_error_returns_empty(self, tmp_path: Path) -> None:
        """UnicodeDecodeError returns empty violations."""
        f = tmp_path / "binary.bin"
        f.write_bytes(b"\xff\xfe\x00")
        # Rename to .md to trigger check
        md = tmp_path / "bad.md"
        f.rename(md)
        violations = check_file(md)
        assert violations == []


class TestCheckDirectory:
    """Tests for check_directory function."""

    def test_nonexistent_dir_returns_empty(self) -> None:
        """Non-existent directory returns empty."""
        assert check_directory(Path("/nonexistent")) == []

    def test_not_dir_returns_empty(self, tmp_path: Path) -> None:
        """File path returns empty."""
        f = tmp_path / "file.md"
        f.write_text("")
        assert check_directory(f) == []

    def test_finds_violations_in_md_files(self, tmp_path: Path) -> None:
        """Scans .md files in directory."""
        (tmp_path / "a.md").write_text("estimated_duration: 1h")
        (tmp_path / "b.md").write_text("clean content")
        violations = check_directory(tmp_path)
        assert len(violations) >= 1

    def test_custom_glob_pattern(self, tmp_path: Path) -> None:
        """check_directory respects glob parameter."""
        (tmp_path / "a.md").write_text("estimated_duration: 1h")
        violations = check_directory(tmp_path, glob="*.txt")
        assert violations == []


class TestMain:
    """Tests for main CLI."""

    def test_main_no_args_returns_1(self) -> None:
        """main with no args returns 1."""
        with patch.object(sys, "argv", ["time_estimate_checker"]):
            from sdp.validators.time_estimate_checker import main

            assert main() == 1

    def test_main_path_not_found_returns_1(self) -> None:
        """main with nonexistent path returns 1."""
        with patch.object(sys, "argv", ["tc", "/nonexistent/path"]):
            from sdp.validators.time_estimate_checker import main

            assert main() == 1

    def test_main_file_no_violations_returns_0(self, tmp_path: Path) -> None:
        """main with clean file returns 0."""
        f = tmp_path / "clean.md"
        f.write_text("Clean content only")
        with patch.object(sys, "argv", ["tc", str(f)]):
            from sdp.validators.time_estimate_checker import main

            assert main() == 0
