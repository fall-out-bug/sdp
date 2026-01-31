"""Tests for CLI doctor commands."""

from pathlib import Path
from click.testing import CliRunner
import pytest


class TestDoctorFunctions:
    """Test doctor utility functions."""

    def test_check_environment(self) -> None:
        """Test environment check function."""
        from sdp.cli.doctor import check_environment

        # Should return boolean
        result = check_environment()
        assert isinstance(result, bool)

    def test_check_project_structure(self, tmp_path: Path) -> None:
        """Test project structure check."""
        from sdp.cli.doctor import check_project_structure

        # Create minimal project structure
        (tmp_path / "docs").mkdir()
        (tmp_path / "src").mkdir()
        (tmp_path / "tests").mkdir()
        (tmp_path / "README.md").write_text("# Test")
        (tmp_path / "pyproject.toml").write_text("[tool.poetry]\nname = 'test'")

        result = check_project_structure(tmp_path)
        assert result is True

    def test_check_project_structure_missing_dirs(self, tmp_path: Path) -> None:
        """Test project structure check with missing directories."""
        from sdp.cli.doctor import check_project_structure

        # Only create one required directory
        (tmp_path / "src").mkdir()
        (tmp_path / "README.md").write_text("# Test")

        result = check_project_structure(tmp_path)
        assert result is False

    def test_check_workstreams_no_directory(self, tmp_path: Path) -> None:
        """Test workstream check with no workstreams directory."""
        from sdp.cli.doctor import check_workstreams

        result = check_workstreams(tmp_path)
        assert result is True

    def test_check_workstreams_empty_directory(self, tmp_path: Path) -> None:
        """Test workstream check with empty workstreams directory."""
        from sdp.cli.doctor import check_workstreams

        docs_dir = tmp_path / "docs" / "workstreams"
        docs_dir.mkdir(parents=True)

        result = check_workstreams(tmp_path)
        assert result is True
