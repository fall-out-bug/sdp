"""Tests for CLI PRD commands."""

from pathlib import Path
from unittest.mock import patch

import pytest
from click.testing import CliRunner

from sdp.cli.prd import prd, prd_detect_type, prd_validate
from sdp.prd.validator import Severity, ValidationIssue


class TestPrdGroup:
    """Test PRD command group."""

    def test_prd_group_exists(self) -> None:
        """Test that PRD command group exists."""
        runner = CliRunner()
        result = runner.invoke(prd, ["--help"])
        assert result.exit_code == 0
        assert "PRD (Product Requirements Document)" in result.output

    def test_prd_validate_command_exists(self) -> None:
        """Test that PRD validate command exists."""
        runner = CliRunner()
        result = runner.invoke(prd, ["validate", "--help"])
        assert result.exit_code == 0
        assert "Validate a PRD document" in result.output

    def test_prd_detect_type_command_exists(self) -> None:
        """Test that PRD detect-type command exists."""
        runner = CliRunner()
        result = runner.invoke(prd, ["detect-type", "--help"])
        assert result.exit_code == 0
        assert "Detect project type" in result.output


class TestPrdValidateCommand:
    """Test prd_validate command execution."""

    def test_prd_validate_success(self, tmp_path: Path) -> None:
        """Test PRD validation when validation passes."""
        prd_file = tmp_path / "PROJECT_MAP.md"
        prd_file.write_text("""# Project Map

## 1. Назначение
Short purpose description.

## 2. Модель БД
Simple DB model.
""")

        runner = CliRunner()
        # Patch at the source module since it's imported inside the function
        with patch("sdp.prd.validator.validate_prd_file", return_value=[]):
            result = runner.invoke(prd_validate, [str(prd_file)])

            assert result.exit_code == 0
            assert "PRD validation passed" in result.output

    def test_prd_validate_with_warnings(self, tmp_path: Path) -> None:
        """Test PRD validation with warnings."""
        prd_file = tmp_path / "PROJECT_MAP.md"
        prd_file.write_text("Test content")

        issues = [
            ValidationIssue(
                section="Назначение",
                message="Too long",
                severity=Severity.WARNING,
                current=600,
                limit=500,
            ),
        ]

        runner = CliRunner()
        # Patch at the source module since it's imported inside the function
        with patch("sdp.prd.validator.validate_prd_file", return_value=issues):
            result = runner.invoke(prd_validate, [str(prd_file)])

            assert result.exit_code == 0
            assert "Too long" in result.output

    def test_prd_validate_with_errors(self, tmp_path: Path) -> None:
        """Test PRD validation with errors."""
        prd_file = tmp_path / "PROJECT_MAP.md"
        prd_file.write_text("Test content")

        issues = [
            ValidationIssue(
                section="Модель БД",
                message="Line too long",
                severity=Severity.ERROR,
                current=150,
                limit=120,
            ),
        ]

        runner = CliRunner()
        # Patch at the source module since it's imported inside the function
        with patch("sdp.prd.validator.validate_prd_file", return_value=issues):
            result = runner.invoke(prd_validate, [str(prd_file)])

            assert result.exit_code == 0  # Default: no exit code on error

    def test_prd_validate_exit_code_on_error(self, tmp_path: Path) -> None:
        """Test PRD validation with --exit-code-on-error flag."""
        prd_file = tmp_path / "PROJECT_MAP.md"
        prd_file.write_text("Test content")

        issues = [
            ValidationIssue(
                section="Модель БД",
                message="Line too long",
                severity=Severity.ERROR,
                current=150,
                limit=120,
            ),
        ]

        runner = CliRunner()
        # Patch at the source module since it's imported inside the function
        with patch("sdp.prd.validator.validate_prd_file", return_value=issues):
            result = runner.invoke(
                prd_validate,
                [str(prd_file), "--exit-code-on-error"],
            )

            assert result.exit_code == 1

    def test_prd_validate_multiple_issues(self, tmp_path: Path) -> None:
        """Test PRD validation with multiple issues."""
        prd_file = tmp_path / "PROJECT_MAP.md"
        prd_file.write_text("Test content")

        issues = [
            ValidationIssue(
                section="Назначение",
                message="Warning 1",
                severity=Severity.WARNING,
                current=600,
                limit=500,
            ),
            ValidationIssue(
                section="Модель БД",
                message="Error 1",
                severity=Severity.ERROR,
                current=150,
                limit=120,
            ),
            ValidationIssue(
                section="Модель БД",
                message="Error 2",
                severity=Severity.ERROR,
                current=200,
                limit=120,
            ),
        ]

        runner = CliRunner()
        # Patch at the source module since it's imported inside the function
        with patch("sdp.prd.validator.validate_prd_file", return_value=issues):
            result = runner.invoke(prd_validate, [str(prd_file)])

            assert "Warning 1" in result.output
            assert "Error 1" in result.output
            assert "Error 2" in result.output


class TestPrdDetectTypeCommand:
    """Test prd_detect_type command execution."""

    def test_detect_type_service(self, tmp_path: Path) -> None:
        """Test project type detection for service."""
        project_path = tmp_path
        (project_path / "docker-compose.yml").write_text("version: '3'")

        runner = CliRunner()
        with patch("sdp.prd.detector.detect_project_type") as mock_detect:
            from sdp.prd.profiles import ProjectType

            mock_detect.return_value = ProjectType.SERVICE
            result = runner.invoke(prd_detect_type, [str(project_path)])

            assert result.exit_code == 0
            assert "Detected project type: service" in result.output

    def test_detect_type_cli(self, tmp_path: Path) -> None:
        """Test project type detection for CLI."""
        project_path = tmp_path
        (project_path / "cli.py").write_text("import click")

        runner = CliRunner()
        with patch("sdp.prd.detector.detect_project_type") as mock_detect:
            from sdp.prd.profiles import ProjectType

            mock_detect.return_value = ProjectType.CLI
            result = runner.invoke(prd_detect_type, [str(project_path)])

            assert result.exit_code == 0
            assert "Detected project type: cli" in result.output

    def test_detect_type_library(self, tmp_path: Path) -> None:
        """Test project type detection for library."""
        project_path = tmp_path

        runner = CliRunner()
        with patch("sdp.prd.detector.detect_project_type") as mock_detect:
            from sdp.prd.profiles import ProjectType

            mock_detect.return_value = ProjectType.LIBRARY
            result = runner.invoke(prd_detect_type, [str(project_path)])

            assert result.exit_code == 0
            assert "Detected project type: library" in result.output
