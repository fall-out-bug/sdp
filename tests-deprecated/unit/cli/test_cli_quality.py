"""Tests for CLI quality commands."""

import subprocess
from pathlib import Path
from unittest.mock import MagicMock, patch

import pytest
from click.testing import CliRunner

from sdp.cli.quality import quality, quality_check


class TestQualityGroup:
    """Test quality command group."""

    def test_quality_group_exists(self) -> None:
        """Test that quality command group exists."""
        runner = CliRunner()
        result = runner.invoke(quality, ["--help"])
        assert result.exit_code == 0
        assert "Quality assurance" in result.output

    def test_quality_check_command_exists(self) -> None:
        """Test that quality check command exists."""
        runner = CliRunner()
        result = runner.invoke(quality, ["check", "--help"])
        assert result.exit_code == 0
        assert "Run quality gates" in result.output


class TestQualityCheckCommand:
    """Test quality_check command execution."""

    def test_quality_check_with_file(self, tmp_path: Path) -> None:
        """Test quality check with a file path."""
        test_file = tmp_path / "test.py"
        test_file.write_text("def test(): pass\n")

        runner = CliRunner()
        with patch("subprocess.run") as mock_run:
            # Mock all subprocess calls to succeed
            mock_run.return_value = MagicMock(returncode=0, stdout="", stderr="")

            result = runner.invoke(
                quality_check,
                [str(test_file), "--threshold", "80", "--max-cc", "10", "--max-loc", "200"],
            )

            # Should call pytest, mypy, ruff, radon, and file size check
            assert mock_run.call_count == 5

    def test_quality_check_with_directory(self, tmp_path: Path) -> None:
        """Test quality check with a directory path."""
        test_dir = tmp_path / "src"
        test_dir.mkdir()
        (test_dir / "test.py").write_text("def test(): pass\n")

        runner = CliRunner()
        with patch("subprocess.run") as mock_run:
            mock_run.return_value = MagicMock(returncode=0, stdout="", stderr="")

            result = runner.invoke(quality_check, [str(test_dir)])

            assert mock_run.call_count == 5

    def test_quality_check_custom_threshold(self, tmp_path: Path) -> None:
        """Test quality check with custom coverage threshold."""
        test_file = tmp_path / "test.py"
        test_file.write_text("def test(): pass\n")

        runner = CliRunner()
        with patch("subprocess.run") as mock_run:
            mock_run.return_value = MagicMock(returncode=0, stdout="", stderr="")

            result = runner.invoke(
                quality_check,
                [str(test_file), "--threshold", "90"],
            )

            # Check that threshold was passed to pytest
            pytest_calls = [
                call for call in mock_run.call_args_list
                if "pytest" in str(call)
            ]
            assert len(pytest_calls) > 0
            assert "--cov-fail-under=90" in str(pytest_calls[0])

    def test_quality_check_custom_max_cc(self, tmp_path: Path) -> None:
        """Test quality check with custom max cyclomatic complexity."""
        test_file = tmp_path / "test.py"
        test_file.write_text("def test(): pass\n")

        runner = CliRunner()
        with patch("subprocess.run") as mock_run:
            mock_run.return_value = MagicMock(
                returncode=0,
                stdout="test.py:1:test - A (5)",
                stderr="",
            )

            result = runner.invoke(
                quality_check,
                [str(test_file), "--max-cc", "15"],
            )

            assert mock_run.call_count == 5

    def test_quality_check_custom_max_loc(self, tmp_path: Path) -> None:
        """Test quality check with custom max lines of code."""
        test_file = tmp_path / "test.py"
        test_file.write_text("def test(): pass\n")

        runner = CliRunner()
        with patch("subprocess.run") as mock_run:
            mock_run.return_value = MagicMock(returncode=0, stdout="", stderr="")

            result = runner.invoke(
                quality_check,
                [str(test_file), "--max-loc", "300"],
            )

            assert mock_run.call_count == 5

    def test_quality_check_coverage_failure(self, tmp_path: Path) -> None:
        """Test quality check when coverage check fails."""
        test_file = tmp_path / "test.py"
        test_file.write_text("def test(): pass\n")

        runner = CliRunner()
        with patch("subprocess.run") as mock_run:
            # First call (pytest) fails, others succeed
            mock_run.side_effect = [
                MagicMock(returncode=1, stdout="", stderr="Coverage too low"),
                MagicMock(returncode=0, stdout="", stderr=""),
                MagicMock(returncode=0, stdout="", stderr=""),
                MagicMock(returncode=0, stdout="", stderr=""),
                MagicMock(returncode=0, stdout="", stderr=""),
            ]

            result = runner.invoke(quality_check, [str(test_file)])

            assert result.exit_code == 1
            assert "Coverage <" in result.output

    def test_quality_check_type_check_failure(self, tmp_path: Path) -> None:
        """Test quality check when type checking fails."""
        test_file = tmp_path / "test.py"
        test_file.write_text("def test(): pass\n")

        runner = CliRunner()
        with patch("subprocess.run") as mock_run:
            mock_run.side_effect = [
                MagicMock(returncode=0, stdout="", stderr=""),
                MagicMock(returncode=1, stdout="", stderr="Type error"),
                MagicMock(returncode=0, stdout="", stderr=""),
                MagicMock(returncode=0, stdout="", stderr=""),
                MagicMock(returncode=0, stdout="", stderr=""),
            ]

            result = runner.invoke(quality_check, [str(test_file)])

            assert result.exit_code == 1
            assert "Type checking failed" in result.output

    def test_quality_check_linting_failure(self, tmp_path: Path) -> None:
        """Test quality check when linting fails."""
        test_file = tmp_path / "test.py"
        test_file.write_text("def test(): pass\n")

        runner = CliRunner()
        with patch("subprocess.run") as mock_run:
            mock_run.side_effect = [
                MagicMock(returncode=0, stdout="", stderr=""),
                MagicMock(returncode=0, stdout="", stderr=""),
                MagicMock(returncode=1, stdout="", stderr="Lint error"),
                MagicMock(returncode=0, stdout="", stderr=""),
                MagicMock(returncode=0, stdout="", stderr=""),
            ]

            result = runner.invoke(quality_check, [str(test_file)])

            assert result.exit_code == 1
            assert "Linting failed" in result.output

    def test_quality_check_complexity_failure(self, tmp_path: Path) -> None:
        """Test quality check when complexity check fails."""
        test_file = tmp_path / "test.py"
        test_file.write_text("def test(): pass\n")

        runner = CliRunner()
        with patch("subprocess.run") as mock_run:
            mock_run.side_effect = [
                MagicMock(returncode=0, stdout="", stderr=""),
                MagicMock(returncode=0, stdout="", stderr=""),
                MagicMock(returncode=0, stdout="", stderr=""),
                MagicMock(
                    returncode=0,
                    stdout="test.py:1:complex_function - A (15)",
                    stderr="",
                ),
                MagicMock(returncode=0, stdout="", stderr=""),
            ]

            result = runner.invoke(quality_check, [str(test_file)])

            assert result.exit_code == 1
            assert "exceed CC" in result.output

    def test_quality_check_file_size_failure(self, tmp_path: Path) -> None:
        """Test quality check when file size check fails."""
        test_file = tmp_path / "test.py"
        test_file.write_text("def test(): pass\n")

        runner = CliRunner()
        with patch("subprocess.run") as mock_run:
            mock_run.side_effect = [
                MagicMock(returncode=0, stdout="", stderr=""),
                MagicMock(returncode=0, stdout="", stderr=""),
                MagicMock(returncode=0, stdout="", stderr=""),
                MagicMock(returncode=0, stdout="", stderr=""),
                MagicMock(returncode=1, stdout="", stderr="File too large"),
            ]

            result = runner.invoke(quality_check, [str(test_file)])

            assert result.exit_code == 1
            assert "exceed" in result.output

    def test_quality_check_exception_handling(self, tmp_path: Path) -> None:
        """Test quality check handles exceptions gracefully."""
        test_file = tmp_path / "test.py"
        test_file.write_text("def test(): pass\n")

        runner = CliRunner()
        with patch("subprocess.run") as mock_run:
            mock_run.side_effect = Exception("Subprocess error")

            result = runner.invoke(quality_check, [str(test_file)])

            assert result.exit_code == 1
            assert "failed" in result.output.lower()

    def test_quality_check_all_passed(self, tmp_path: Path) -> None:
        """Test quality check when all gates pass."""
        test_file = tmp_path / "test.py"
        test_file.write_text("def test(): pass\n")

        runner = CliRunner()
        with patch("subprocess.run") as mock_run:
            mock_run.return_value = MagicMock(returncode=0, stdout="", stderr="")

            result = runner.invoke(quality_check, [str(test_file)])

            assert result.exit_code == 0
            assert "All quality gates passed" in result.output
