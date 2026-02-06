"""Tests for health_checks module.

Tests health check functionality for SDP doctor command.
"""

from pathlib import Path
from unittest.mock import Mock, patch

import pytest
from sdp.health_checks import (
    BeadsCLICheck,
    GitHubCLICheck,
    GitHooksCheck,
    HealthCheck,
    HealthCheckResult,
    PoetryCheck,
    PythonVersionCheck,
    TelegramConfigCheck,
    get_health_checks,
)


class TestHealthCheckResult:
    """Tests for HealthCheckResult class."""

    def test_health_check_result_init(self) -> None:
        """Test HealthCheckResult initialization."""
        result = HealthCheckResult(
            name="Test Check",
            passed=True,
            message="Test passed",
            remediation=None,
        )

        assert result.name == "Test Check"
        assert result.passed is True
        assert result.message == "Test passed"
        assert result.remediation is None

    def test_health_check_result_with_remediation(self) -> None:
        """Test HealthCheckResult with remediation."""
        result = HealthCheckResult(
            name="Test Check",
            passed=False,
            message="Test failed",
            remediation="Fix it",
        )

        assert result.remediation == "Fix it"


class TestPythonVersionCheck:
    """Tests for PythonVersionCheck class."""

    @patch("sdp.health_checks.checks.sys")
    def test_python_version_check_310(self, mock_sys: Mock) -> None:
        """Test Python version check with 3.10."""
        mock_sys.version_info = (3, 10, 0)

        check = PythonVersionCheck()
        result = check.run()

        assert result.passed is True
        assert "3.10" in result.message

    @patch("sdp.health_checks.checks.sys")
    def test_python_version_check_311(self, mock_sys: Mock) -> None:
        """Test Python version check with 3.11."""
        mock_sys.version_info = (3, 11, 5)

        check = PythonVersionCheck()
        result = check.run()

        assert result.passed is True
        assert "3.11" in result.message

    @patch("sdp.health_checks.checks.sys")
    def test_python_version_check_39(self, mock_sys: Mock) -> None:
        """Test Python version check with 3.9 (fails)."""
        mock_sys.version_info = (3, 9, 0)

        check = PythonVersionCheck()
        result = check.run()

        assert result.passed is False
        assert "3.9" in result.message
        assert result.remediation is not None

    def test_python_version_check_critical(self) -> None:
        """Test PythonVersionCheck is critical."""
        check = PythonVersionCheck()
        assert check.critical is True


class TestPoetryCheck:
    """Tests for PoetryCheck class."""

    @patch("sdp.health_checks.checks.shutil.which")
    def test_poetry_check_installed(self, mock_which: Mock) -> None:
        """Test Poetry check when installed."""
        mock_which.return_value = "/usr/bin/poetry"

        check = PoetryCheck()
        result = check.run()

        assert result.passed is True
        assert "installed" in result.message.lower()

    @patch("sdp.health_checks.checks.shutil.which")
    def test_poetry_check_not_installed(self, mock_which: Mock) -> None:
        """Test Poetry check when not installed."""
        mock_which.return_value = None

        check = PoetryCheck()
        result = check.run()

        assert result.passed is False
        assert "not found" in result.message.lower()
        assert result.remediation is not None

    def test_poetry_check_critical(self) -> None:
        """Test PoetryCheck is critical."""
        check = PoetryCheck()
        assert check.critical is True


class TestGitHooksCheck:
    """Tests for GitHooksCheck class."""

    def test_git_hooks_check_installed(self, tmp_path: Path) -> None:
        """Test git hooks check when installed."""
        git_dir = tmp_path / ".git"
        git_dir.mkdir()
        hooks_dir = git_dir / "hooks"
        hooks_dir.mkdir()

        with patch("sdp.health_checks.checks.Path.cwd", return_value=tmp_path):
            check = GitHooksCheck()
            result = check.run()

        assert result.passed is True
        assert "found" in result.message.lower()

    def test_git_hooks_check_not_installed(self, tmp_path: Path) -> None:
        """Test git hooks check when not installed."""
        with patch("sdp.health_checks.checks.Path.cwd", return_value=tmp_path):
            check = GitHooksCheck()
            result = check.run()

        assert result.passed is False
        assert "not found" in result.message.lower()
        assert result.remediation is not None

    def test_git_hooks_check_critical(self) -> None:
        """Test GitHooksCheck is critical."""
        check = GitHooksCheck()
        assert check.critical is True


class TestBeadsCLICheck:
    """Tests for BeadsCLICheck class."""

    @patch("sdp.health_checks.checks.shutil.which")
    def test_beads_check_installed(self, mock_which: Mock) -> None:
        """Test Beads CLI check when installed."""
        mock_which.return_value = "/usr/bin/bd"

        check = BeadsCLICheck()
        result = check.run()

        assert result.passed is True
        assert "installed" in result.message.lower()

    @patch("sdp.health_checks.checks.shutil.which")
    def test_beads_check_not_installed(self, mock_which: Mock) -> None:
        """Test Beads CLI check when not installed."""
        mock_which.return_value = None

        check = BeadsCLICheck()
        result = check.run()

        # Optional checks still return True
        assert result.passed is True
        assert "optional" in result.message.lower()

    def test_beads_check_not_critical(self) -> None:
        """Test BeadsCLICheck is not critical."""
        check = BeadsCLICheck()
        assert check.critical is False


class TestGitHubCLICheck:
    """Tests for GitHubCLICheck class."""

    @patch("sdp.health_checks.checks.shutil.which")
    def test_github_check_installed(self, mock_which: Mock) -> None:
        """Test GitHub CLI check when installed."""
        mock_which.return_value = "/usr/bin/gh"

        check = GitHubCLICheck()
        result = check.run()

        assert result.passed is True
        assert "installed" in result.message.lower()

    @patch("sdp.health_checks.checks.shutil.which")
    def test_github_check_not_installed(self, mock_which: Mock) -> None:
        """Test GitHub CLI check when not installed."""
        mock_which.return_value = None

        check = GitHubCLICheck()
        result = check.run()

        # Optional checks still return True
        assert result.passed is True
        assert "optional" in result.message.lower()

    def test_github_check_not_critical(self) -> None:
        """Test GitHubCLICheck is not critical."""
        check = GitHubCLICheck()
        assert check.critical is False


class TestTelegramConfigCheck:
    """Tests for TelegramConfigCheck class."""

    def test_telegram_check_configured(self, tmp_path: Path) -> None:
        """Test Telegram check when configured."""
        env_file = tmp_path / ".env"
        env_file.write_text(
            "TELEGRAM_TOKEN=test_token\nTELEGRAM_CHAT_ID=test_chat\n"
        )

        with patch("sdp.health_checks.checks.Path.cwd", return_value=tmp_path):
            check = TelegramConfigCheck()
            result = check.run()

        assert result.passed is True
        assert "configured" in result.message.lower()

    def test_telegram_check_partial_config(self, tmp_path: Path) -> None:
        """Test Telegram check with partial config."""
        env_file = tmp_path / ".env"
        env_file.write_text("TELEGRAM_TOKEN=test_token\n")

        with patch("sdp.health_checks.checks.Path.cwd", return_value=tmp_path):
            check = TelegramConfigCheck()
            result = check.run()

        # Optional checks still return True
        assert result.passed is True
        assert "optional" in result.message.lower()

    def test_telegram_check_no_env_file(self, tmp_path: Path) -> None:
        """Test Telegram check with no .env file."""
        with patch("sdp.health_checks.checks.Path.cwd", return_value=tmp_path):
            check = TelegramConfigCheck()
            result = check.run()

        # Optional checks still return True
        assert result.passed is True
        assert "optional" in result.message.lower()

    def test_telegram_check_not_critical(self) -> None:
        """Test TelegramConfigCheck is not critical."""
        check = TelegramConfigCheck()
        assert check.critical is False


class TestGetHealthChecks:
    """Tests for get_health_checks function."""

    def test_get_health_checks_returns_list(self) -> None:
        """Test get_health_checks returns list."""
        checks = get_health_checks()

        assert isinstance(checks, list)
        assert len(checks) == 6

    def test_get_health_checks_all_critical_flags(self) -> None:
        """Test get_health_checks has correct critical flags."""
        checks = get_health_checks()

        critical_checks = [c for c in checks if c.critical]
        optional_checks = [c for c in checks if not c.critical]

        # Python, Poetry, Git Hooks should be critical
        assert len(critical_checks) == 3
        critical_names = [c.name for c in critical_checks]
        assert "Python Version" in critical_names
        assert "Poetry" in critical_names
        assert "Git Hooks" in critical_names

        # Beads, GitHub CLI, Telegram should be optional
        assert len(optional_checks) == 3
        optional_names = [c.name for c in optional_checks]
        assert "Beads CLI" in optional_names
        assert "GitHub CLI" in optional_names
        assert "Telegram Config" in optional_names

    def test_get_health_checks_runnable(self) -> None:
        """Test all health checks can run."""
        checks = get_health_checks()

        for check in checks:
            result = check.run()
            assert isinstance(result, HealthCheckResult)
            assert isinstance(result.passed, bool)
            assert isinstance(result.name, str)
            assert isinstance(result.message, str)

    def test_get_health_checks_correct_types(self) -> None:
        """Test all health checks are correct types."""
        checks = get_health_checks()

        assert isinstance(checks[0], PythonVersionCheck)
        assert isinstance(checks[1], PoetryCheck)
        assert isinstance(checks[2], GitHooksCheck)
        assert isinstance(checks[3], BeadsCLICheck)
        assert isinstance(checks[4], GitHubCLICheck)
        assert isinstance(checks[5], TelegramConfigCheck)
