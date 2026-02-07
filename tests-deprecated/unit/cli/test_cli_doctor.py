"""Tests for CLI doctor commands."""

from pathlib import Path
from unittest.mock import MagicMock, patch
import subprocess

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

    def test_check_environment_python_version_ok(self) -> None:
        """Test environment check with good Python version."""
        from sdp.cli.doctor import check_environment
        
        # Current Python should be 3.10+
        result = check_environment()
        # Should pass Python check (at minimum)
        assert isinstance(result, bool)

    def test_check_environment_poetry_not_found(self) -> None:
        """Test environment check when Poetry is not available."""
        from sdp.cli.doctor import check_environment
        
        with patch("subprocess.run") as mock_run:
            # Mock Poetry check to fail
            mock_run.side_effect = [
                MagicMock(returncode=1, stdout=""),  # Poetry fails
                MagicMock(returncode=0, stdout="git version 2.0"),  # Git succeeds
            ]
            
            result = check_environment()
            assert result is False

    def test_check_environment_git_not_found(self) -> None:
        """Test environment check when Git is not available."""
        from sdp.cli.doctor import check_environment
        
        with patch("subprocess.run") as mock_run:
            # Mock Git check to fail
            mock_run.side_effect = [
                MagicMock(returncode=0, stdout="Poetry version 1.0"),  # Poetry succeeds
                MagicMock(returncode=1, stdout=""),  # Git fails
            ]
            
            result = check_environment()
            assert result is False

    def test_check_environment_exception_handling(self) -> None:
        """Test environment check handles exceptions."""
        from sdp.cli.doctor import check_environment
        
        with patch("subprocess.run") as mock_run:
            # Mock exception during Poetry check
            mock_run.side_effect = [
                Exception("Command not found"),  # Poetry raises
                MagicMock(returncode=0, stdout="git version 2.0"),  # Git succeeds
            ]
            
            result = check_environment()
            assert result is False

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

    def test_check_project_structure_missing_files(self, tmp_path: Path) -> None:
        """Test project structure check with missing files."""
        from sdp.cli.doctor import check_project_structure

        # Create all directories but missing pyproject.toml
        (tmp_path / "docs").mkdir()
        (tmp_path / "src").mkdir()
        (tmp_path / "tests").mkdir()
        (tmp_path / "README.md").write_text("# Test")

        result = check_project_structure(tmp_path)
        assert result is False

    def test_check_project_structure_all_missing(self, tmp_path: Path) -> None:
        """Test project structure check with everything missing."""
        from sdp.cli.doctor import check_project_structure

        # Empty directory
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

    def test_check_workstreams_valid_files(self, tmp_path: Path) -> None:
        """Test workstream check with valid workstream files."""
        from sdp.cli.doctor import check_workstreams

        ws_dir = tmp_path / "docs" / "workstreams"
        ws_dir.mkdir(parents=True)
        
        # Create valid workstream file
        ws_file = ws_dir / "WS-001-01-test.md"
        ws_file.write_text("""---
ws_id: WS-001-01
title: Test Workstream
feature: F001
status: backlog
size: SMALL
---

# Test Workstream

## Acceptance Criteria
- [ ] Test criterion
""")
        
        result = check_workstreams(tmp_path)
        assert result is True

    def test_check_workstreams_invalid_file(self, tmp_path: Path) -> None:
        """Test workstream check with invalid workstream file."""
        from sdp.cli.doctor import check_workstreams

        ws_dir = tmp_path / "docs" / "workstreams"
        ws_dir.mkdir(parents=True)
        
        # Create invalid workstream file
        ws_file = ws_dir / "WS-001-01-test.md"
        ws_file.write_text("Invalid content")
        
        result = check_workstreams(tmp_path)
        assert result is False

    def test_check_workstreams_parse_error(self, tmp_path: Path) -> None:
        """Test workstream check handles parse errors."""
        from sdp.cli.doctor import check_workstreams

        ws_dir = tmp_path / "docs" / "workstreams"
        ws_dir.mkdir(parents=True)
        
        # Create workstream file with missing required fields
        ws_file = ws_dir / "WS-001-01-test.md"
        ws_file.write_text("""---
title: Test Workstream
---

# Test
""")
        
        result = check_workstreams(tmp_path)
        assert result is False

    def test_check_workstreams_generic_exception(self, tmp_path: Path) -> None:
        """Test workstream check handles generic exceptions."""
        from sdp.cli.doctor import check_workstreams

        ws_dir = tmp_path / "docs" / "workstreams"
        ws_dir.mkdir(parents=True)
        
        # Create a file that will cause issues
        ws_file = ws_dir / "WS-001-01-test.md"
        ws_file.write_text("---\n" * 1000)  # Malformed frontmatter
        
        result = check_workstreams(tmp_path)
        # Should handle gracefully
        assert isinstance(result, bool)

    def test_check_workstreams_multiple_files(self, tmp_path: Path) -> None:
        """Test workstream check with multiple workstream files."""
        from sdp.cli.doctor import check_workstreams

        ws_dir = tmp_path / "docs" / "workstreams"
        ws_dir.mkdir(parents=True)
        
        # Create multiple valid workstream files
        for i in range(1, 4):
            ws_file = ws_dir / f"WS-001-0{i}-test.md"
            ws_file.write_text(f"""---
ws_id: WS-001-0{i}
title: Test Workstream {i}
feature: F001
status: backlog
size: SMALL
---

# Test Workstream {i}
""")
        
        result = check_workstreams(tmp_path)
        assert result is True

    def test_check_workstreams_mixed_valid_invalid(self, tmp_path: Path) -> None:
        """Test workstream check with mix of valid and invalid files."""
        from sdp.cli.doctor import check_workstreams

        ws_dir = tmp_path / "docs" / "workstreams"
        ws_dir.mkdir(parents=True)
        
        # Create one valid file
        valid_file = ws_dir / "WS-001-01-test.md"
        valid_file.write_text("""---
ws_id: WS-001-01
title: Test Workstream
feature: F001
status: backlog
size: SMALL
---

# Test Workstream
""")
        
        # Create one invalid file
        invalid_file = ws_dir / "WS-001-02-test.md"
        invalid_file.write_text("Invalid content")
        
        result = check_workstreams(tmp_path)
        # Should fail due to invalid file
        assert result is False
