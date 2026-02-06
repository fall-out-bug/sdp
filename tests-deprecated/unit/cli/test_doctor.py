"""Tests for sdp doctor command."""

import sys
from pathlib import Path
from unittest.mock import MagicMock, patch

import pytest
from click.testing import CliRunner

from sdp.doctor import doctor
from sdp.health_checks import (
    HealthCheck,
    HealthCheckResult,
    PythonVersionCheck,
    PoetryCheck,
    GitHooksCheck,
    BeadsCLICheck,
    GitHubCLICheck,
    TelegramConfigCheck,
)


class TestHealthCheckResult:
    """Test HealthCheckResult model."""

    def test_success_result(self):
        """Test successful health check result."""
        result = HealthCheckResult(
            name="test_check",
            passed=True,
            message="Check passed",
            remediation=None,
        )
        assert result.name == "test_check"
        assert result.passed is True
        assert result.message == "Check passed"
        assert result.remediation is None

    def test_failure_result(self):
        """Test failed health check result."""
        result = HealthCheckResult(
            name="test_check",
            passed=False,
            message="Check failed",
            remediation="Run this command",
        )
        assert result.name == "test_check"
        assert result.passed is False
        assert result.message == "Check failed"
        assert result.remediation == "Run this command"


class TestPythonVersionCheck:
    """Test Python version check."""

    def test_python_version_pass(self):
        """Test Python version check passes."""
        check = PythonVersionCheck()

        with patch("sys.version_info", (3, 12, 0)):
            result = check.run()

        assert result.passed is True
        assert "3.12" in result.message
        assert result.remediation is None

    def test_python_version_fail(self):
        """Test Python version check fails."""
        check = PythonVersionCheck()

        with patch("sys.version_info", (3, 9, 0)):
            result = check.run()

        assert result.passed is False
        assert "3.9" in result.message
        assert result.remediation is not None

    def test_python_version_exact_minimum(self):
        """Test Python version at minimum threshold."""
        check = PythonVersionCheck()

        with patch("sys.version_info", (3, 10, 0)):
            result = check.run()

        assert result.passed is True


class TestPoetryCheck:
    """Test Poetry check."""

    def test_poetry_installed(self):
        """Test Poetry is installed."""
        check = PoetryCheck()

        with patch("shutil.which", return_value="/usr/bin/poetry"):
            result = check.run()

        assert result.passed is True
        assert "installed" in result.message.lower()
        assert result.remediation is None

    def test_poetry_not_installed(self):
        """Test Poetry is not installed."""
        check = PoetryCheck()

        with patch("shutil.which", return_value=None):
            result = check.run()

        assert result.passed is False
        assert "not found" in result.message.lower()
        assert result.remediation is not None
        assert "curl" in result.remediation.lower()


class TestGitHooksCheck:
    """Test git hooks check."""

    def test_hooks_exist(self, tmp_path: Path):
        """Test git hooks directory exists."""
        check = GitHooksCheck()

        # Create .git/hooks directory
        hooks_dir = tmp_path / ".git" / "hooks"
        hooks_dir.mkdir(parents=True)

        with patch("pathlib.Path.cwd", return_value=tmp_path):
            result = check.run()

        assert result.passed is True
        assert "found" in result.message.lower()
        assert result.remediation is None

    def test_hooks_missing(self, tmp_path: Path):
        """Test git hooks directory missing."""
        check = GitHooksCheck()

        with patch("pathlib.Path.cwd", return_value=tmp_path):
            result = check.run()

        assert result.passed is False
        assert "not found" in result.message.lower()
        assert result.remediation is not None


class TestBeadsCLICheck:
    """Test Beads CLI check (optional)."""

    def test_beads_installed(self):
        """Test Beads CLI is installed."""
        check = BeadsCLICheck()

        with patch("shutil.which", return_value="/usr/bin/bd"):
            result = check.run()

        assert result.passed is True
        assert result.remediation is None

    def test_beads_not_installed(self):
        """Test Beads CLI is not installed (optional check)."""
        check = BeadsCLICheck()

        with patch("shutil.which", return_value=None):
            result = check.run()

        # Should not fail because it's optional
        assert result.passed is True
        assert "optional" in result.message.lower() or "not installed" in result.message.lower()
        assert result.remediation is not None


class TestGitHubCLICheck:
    """Test GitHub CLI check (optional)."""

    def test_gh_installed(self):
        """Test GitHub CLI is installed."""
        check = GitHubCLICheck()

        with patch("shutil.which", return_value="/usr/bin/gh"):
            result = check.run()

        assert result.passed is True
        assert result.remediation is None

    def test_gh_not_installed(self):
        """Test GitHub CLI is not installed (optional check)."""
        check = GitHubCLICheck()

        with patch("shutil.which", return_value=None):
            result = check.run()

        # Should not fail because it's optional
        assert result.passed is True
        assert "optional" in result.message.lower() or "not installed" in result.message.lower()
        assert result.remediation is not None


class TestTelegramConfigCheck:
    """Test Telegram config check (optional)."""

    def test_telegram_configured(self, tmp_path: Path):
        """Test Telegram is configured."""
        check = TelegramConfigCheck()

        # Create mock .env file with TELEGRAM_TOKEN
        env_file = tmp_path / ".env"
        env_file.write_text("TELEGRAM_TOKEN=test_token\nTELEGRAM_CHAT_ID=123\n")

        with patch("pathlib.Path.cwd", return_value=tmp_path):
            result = check.run()

        assert result.passed is True
        assert result.remediation is None

    def test_telegram_not_configured(self, tmp_path: Path):
        """Test Telegram is not configured (optional check)."""
        check = TelegramConfigCheck()

        # Create .env without Telegram config
        env_file = tmp_path / ".env"
        env_file.write_text("OTHER_VAR=value\n")

        with patch("pathlib.Path.cwd", return_value=tmp_path):
            result = check.run()

        # Should not fail because it's optional
        assert result.passed is True
        assert "optional" in result.message.lower() or "not configured" in result.message.lower()
        assert result.remediation is not None

    def test_no_env_file(self, tmp_path: Path):
        """Test no .env file exists (optional check)."""
        check = TelegramConfigCheck()

        with patch("pathlib.Path.cwd", return_value=tmp_path):
            result = check.run()

        # Should not fail because it's optional
        assert result.passed is True
        assert result.remediation is not None


class TestDoctorCommand:
    """Test sdp doctor command."""

    def test_doctor_all_passed(self, tmp_path: Path):
        """Test doctor command when all checks pass."""
        runner = CliRunner()

        # Create .git/hooks directory
        hooks_dir = tmp_path / ".git" / "hooks"
        hooks_dir.mkdir(parents=True)

        # Create .env with Telegram config
        env_file = tmp_path / ".env"
        env_file.write_text("TELEGRAM_TOKEN=test_token\nTELEGRAM_CHAT_ID=123\n")

        # Patch Path.cwd to return our temp directory
        with patch("pathlib.Path.cwd", return_value=tmp_path):
            with patch("sys.version_info", (3, 12, 0)):
                with patch("shutil.which", return_value="/usr/bin/tool"):
                    result = runner.invoke(doctor)

        # Debug output
        if result.exit_code != 0:
            print(f"Output: {result.output}")
            print(f"Exception: {result.exception}")

        assert result.exit_code == 0

    def test_doctor_critical_fail(self):
        """Test doctor command fails on critical check failure."""
        runner = CliRunner()

        with patch("sys.version_info", (3, 9, 0)):  # Python version too low
            result = runner.invoke(doctor)

        assert result.exit_code == 1

    def test_doctor_performance(self):
        """Test doctor command completes in under 5 seconds."""
        import time

        runner = CliRunner()

        start_time = time.time()

        with patch("sys.version_info", (3, 12, 0)):
            with patch("shutil.which", return_value="/usr/bin/tool"):
                result = runner.invoke(doctor)

        elapsed = time.time() - start_time

        assert result.exit_code == 0
        assert elapsed < 5.0, "Doctor command took longer than 5 seconds"

    def test_doctor_json_output(self):
        """Test doctor command JSON output."""
        runner = CliRunner()

        with patch("sys.version_info", (3, 12, 0)):
            with patch("shutil.which", return_value="/usr/bin/tool"):
                result = runner.invoke(doctor, ["--format", "json"])

        assert result.exit_code == 0
        # Should be valid JSON
        import json

        data = json.loads(result.output)
        assert "checks" in data
        assert "all_passed" in data
