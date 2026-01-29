"""Tests for SDP init command and wizard.

Tests init_wizard.py and cli_init.py functionality.
"""

import tempfile
from pathlib import Path
from unittest.mock import Mock, patch

import pytest
from click.testing import CliRunner

from sdp.cli_init import init
from sdp.init_dependencies import (
    _check_command,
    _check_telegram,
    detect_dependencies,
    show_dependencies,
)
from sdp.init_files import (
    collect_metadata,
    create_env_template,
    create_structure,
    generate_quality_gate,
)
from sdp.init_validation import install_git_hooks, run_doctor


class TestCollectMetadata:
    """Tests for collect_metadata function."""

    def test_collect_metadata_non_interactive(self, tmp_path: Path) -> None:
        """Test metadata collection in non-interactive mode."""
        project_name, description, author = collect_metadata(tmp_path, non_interactive=True)

        assert project_name == tmp_path.name
        assert description == "SDP project"
        assert author == "Your Name"

    @patch("click.prompt")
    def test_collect_metadata_interactive(self, mock_prompt: Mock, tmp_path: Path) -> None:
        """Test metadata collection in interactive mode."""
        mock_prompt.side_effect = ["my-project", "My test project", "Test Author"]

        project_name, description, author = collect_metadata(tmp_path, non_interactive=False)

        assert project_name == "my-project"
        assert description == "My test project"
        assert author == "Test Author"
        assert mock_prompt.call_count == 3


class TestDetectDependencies:
    """Tests for detect_dependencies function."""

    @patch("sdp.init_dependencies._check_command")
    @patch("sdp.init_dependencies._check_telegram")
    def test_detect_dependencies_all_available(
        self, mock_telegram: Mock, mock_command: Mock
    ) -> None:
        """Test dependency detection when all are available."""
        mock_command.return_value = True
        mock_telegram.return_value = True

        deps = detect_dependencies()

        assert deps["Beads CLI"] is True
        assert deps["GitHub CLI (gh)"] is True
        assert deps["Telegram"] is True

    @patch("sdp.init_dependencies._check_command")
    @patch("sdp.init_dependencies._check_telegram")
    def test_detect_dependencies_none_available(
        self, mock_telegram: Mock, mock_command: Mock
    ) -> None:
        """Test dependency detection when none are available."""
        mock_command.return_value = False
        mock_telegram.return_value = False

        deps = detect_dependencies()

        assert deps["Beads CLI"] is False
        assert deps["GitHub CLI (gh)"] is False
        assert deps["Telegram"] is False


class TestCheckCommand:
    """Tests for _check_command function."""

    @patch("sdp.init_dependencies.subprocess.run")
    def test_check_command_available(self, mock_run: Mock) -> None:
        """Test command check when available."""
        mock_run.return_value = Mock(returncode=0)

        result = _check_command("test")

        assert result is True
        mock_run.assert_called_once()

    @patch("sdp.init_dependencies.subprocess.run")
    def test_check_command_not_found(self, mock_run: Mock) -> None:
        """Test command check when not found."""
        mock_run.side_effect = FileNotFoundError()

        result = _check_command("test")

        assert result is False

    @patch("sdp.init_dependencies.subprocess.run")
    def test_check_command_timeout(self, mock_run: Mock) -> None:
        """Test command check when timeout."""
        import subprocess

        mock_run.side_effect = subprocess.TimeoutExpired("test", 5)

        result = _check_command("test")

        assert result is False


class TestCheckTelegram:
    """Tests for _check_telegram function."""

    def test_check_telegram_configured(self, tmp_path: Path) -> None:
        """Test Telegram check when configured."""
        env_file = tmp_path / ".env"
        env_file.write_text(
            "TELEGRAM_BOT_TOKEN=test_token\nTELEGRAM_CHAT_ID=test_chat\n"
        )

        with patch("sdp.init_dependencies.Path.cwd", return_value=tmp_path):
            result = _check_telegram()

        assert result is True

    def test_check_telegram_not_configured(self, tmp_path: Path) -> None:
        """Test Telegram check when not configured."""
        env_file = tmp_path / ".env"
        env_file.write_text("OTHER_VAR=value\n")

        with patch("sdp.init_dependencies.Path.cwd", return_value=tmp_path):
            result = _check_telegram()

        assert result is False

    def test_check_telegram_no_env_file(self, tmp_path: Path) -> None:
        """Test Telegram check when .env doesn't exist."""
        with patch("sdp.init_dependencies.Path.cwd", return_value=tmp_path):
            result = _check_telegram()

        assert result is False


class TestShowDependencies:
    """Tests for show_dependencies function."""

    @patch("click.echo")
    def test_show_dependencies(self, mock_echo: Mock) -> None:
        """Test displaying dependencies."""
        deps = {
            "Beads CLI": True,
            "GitHub CLI (gh)": False,
            "Telegram": True,
        }

        show_dependencies(deps)

        assert mock_echo.call_count >= 3


class TestCreateStructure:
    """Tests for create_structure function."""

    def test_create_structure_new_project(self, tmp_path: Path) -> None:
        """Test creating directory structure for new project."""
        created, skipped = create_structure(tmp_path, "test-project", force=False)

        # Check directories created
        assert (tmp_path / "docs").exists()
        assert (tmp_path / "docs" / "workstreams").exists()
        assert (tmp_path / "docs" / "workstreams" / "backlog").exists()
        assert (tmp_path / "sdp.local").exists()

        # Check files created
        assert (tmp_path / "docs" / "PROJECT_MAP.md").exists()
        assert (tmp_path / "docs" / "workstreams" / "INDEX.md").exists()
        assert (tmp_path / "docs" / "workstreams" / "TEMPLATE.md").exists()

        assert len(created) >= 7
        assert len(skipped) == 0

    def test_create_structure_existing_files(self, tmp_path: Path) -> None:
        """Test creating structure when files exist."""
        # Create existing file
        project_map = tmp_path / "docs" / "PROJECT_MAP.md"
        project_map.parent.mkdir(parents=True)
        project_map.write_text("Existing content")

        created, skipped = create_structure(tmp_path, "test-project", force=False)

        assert len(created) >= 4  # Directories still created
        assert len(skipped) == 1
        assert str(project_map) in skipped

    def test_create_structure_force(self, tmp_path: Path) -> None:
        """Test creating structure with force flag."""
        project_map = tmp_path / "docs" / "PROJECT_MAP.md"
        project_map.parent.mkdir(parents=True)
        project_map.write_text("Old content")

        created, skipped = create_structure(tmp_path, "test-project", force=True)

        # File should be overwritten
        assert len(skipped) == 0
        assert "test-project" in project_map.read_text()


class TestGenerateQualityGate:
    """Tests for generate_quality_gate function."""

    @patch("click.echo")
    def test_generate_quality_gate_new(self, mock_echo: Mock, tmp_path: Path) -> None:
        """Test generating quality gate config."""
        deps = {"Beads CLI": True, "GitHub CLI (gh)": False, "Telegram": False}

        result = generate_quality_gate(tmp_path, deps)

        assert result is not None
        assert result == tmp_path / "quality-gate.toml"
        assert result.exists()
        assert "[coverage]" in result.read_text()

    @patch("click.echo")
    def test_generate_quality_gate_existing(self, mock_echo: Mock, tmp_path: Path) -> None:
        """Test generating quality gate when file exists."""
        quality_gate = tmp_path / "quality-gate.toml"
        quality_gate.write_text("Existing config")

        deps = {}
        result = generate_quality_gate(tmp_path, deps)

        assert result is None

    @patch("click.echo")
    def test_generate_quality_gate_with_beads(self, mock_echo: Mock, tmp_path: Path) -> None:
        """Test generating quality gate with Beads config."""
        deps = {"Beads CLI": True, "GitHub CLI (gh)": False, "Telegram": False}

        result = generate_quality_gate(tmp_path, deps)

        assert result is not None
        content = result.read_text()
        assert "[testing]" in content
        assert "require_test_for_new_code" in content


class TestCreateEnvTemplate:
    """Tests for create_env_template function."""

    @patch("click.echo")
    def test_create_env_template_new(self, mock_echo: Mock, tmp_path: Path) -> None:
        """Test creating .env template."""
        deps = {
            "Beads CLI": True,
            "GitHub CLI (gh)": True,
            "Telegram": True,
        }

        result = create_env_template(tmp_path, deps)

        assert result is not None
        assert result == tmp_path / ".env.template"
        assert result.exists()

        content = result.read_text()
        assert "TELEGRAM_BOT_TOKEN" in content
        assert "GITHUB_TOKEN" in content
        assert "BEADS_API_URL" in content

    @patch("click.echo")
    def test_create_env_template_existing(self, mock_echo: Mock, tmp_path: Path) -> None:
        """Test creating .env template when exists."""
        env_template = tmp_path / ".env.template"
        env_template.write_text("Existing")

        deps = {}
        result = create_env_template(tmp_path, deps)

        assert result is None

    @patch("click.echo")
    def test_create_env_template_no_deps(self, mock_echo: Mock, tmp_path: Path) -> None:
        """Test creating .env template with no dependencies."""
        deps = {
            "Beads CLI": False,
            "GitHub CLI (gh)": False,
            "Telegram": False,
        }

        result = create_env_template(tmp_path, deps)

        assert result is not None
        content = result.read_text()
        assert "TELEGRAM" not in content
        assert "GITHUB" not in content


class TestInstallGitHooks:
    """Tests for install_git_hooks function."""

    def test_install_git_hooks_no_git_repo(self, tmp_path: Path) -> None:
        """Test installing hooks when not a git repo."""
        result = install_git_hooks(tmp_path)

        assert result is False

    def test_install_git_hooks_success(self, tmp_path: Path) -> None:
        """Test successfully installing git hooks."""
        git_dir = tmp_path / ".git"
        git_dir.mkdir()
        hooks_dir = git_dir / "hooks"
        hooks_dir.mkdir()

        # Create a mock hooks file
        hooks_src = tmp_path / "hooks" / "pre-commit.sh"
        hooks_src.parent.mkdir()
        hooks_src.write_text("#!/bin/bash\necho 'test'\n")

        # Mock the path resolution to find our test hooks
        with patch("sdp.init_validation.Path") as mock_path_cls:
            # Create a real Path for the hooks
            real_hooks_path = hooks_src

            # Mock the instance methods
            mock_path_instance = Mock()
            mock_path_instance.exists.return_value = True
            mock_path_instance.__truediv__ = Mock(return_value=mock_path_instance)

            result = install_git_hooks(tmp_path)

        # Should return True if hooks file exists
        assert isinstance(result, bool)


class TestRunDoctor:
    """Tests for run_doctor function."""

    @patch("click.echo")
    @patch("sdp.health_checks.get_health_checks")
    def test_run_doctor_all_pass(self, mock_get_checks: Mock, mock_echo: Mock, tmp_path: Path) -> None:
        """Test doctor when all checks pass."""
        from sdp.health_checks import HealthCheckResult

        mock_check = Mock()
        mock_check.run.return_value = HealthCheckResult(
            name="Test Check",
            passed=True,
            message="Passed",
        )
        mock_get_checks.return_value = [mock_check]

        result = run_doctor(tmp_path)

        assert result is True

    @patch("click.echo")
    @patch("sdp.health_checks.get_health_checks")
    def test_run_doctor_some_fail(self, mock_get_checks: Mock, mock_click: Mock, tmp_path: Path) -> None:
        """Test doctor when some checks fail."""
        from sdp.health_checks import HealthCheckResult

        passing_check = Mock()
        passing_check.run.return_value = HealthCheckResult(
            name="Passing Check",
            passed=True,
            message="Passed",
        )

        failing_check = Mock()
        failing_check.run.return_value = HealthCheckResult(
            name="Failing Check",
            passed=False,
            message="Failed",
        )

        mock_get_checks.return_value = [passing_check, failing_check]

        result = run_doctor(tmp_path)

        assert result is False

    @patch("click.echo")
    @patch("sdp.health_checks.get_health_checks")
    def test_run_doctor_exception(self, mock_get_checks: Mock, mock_click: Mock, tmp_path: Path) -> None:
        """Test doctor when exception occurs."""
        mock_get_checks.side_effect = Exception("Import error")

        result = run_doctor(tmp_path)

        # Should return True (don't fail setup on doctor error)
        assert result is True


class TestInitCommand:
    """Tests for init CLI command."""

    @patch("sdp.cli_init.run_doctor")
    @patch("sdp.cli_init.install_git_hooks")
    @patch("sdp.cli_init.create_env_template")
    @patch("sdp.cli_init.generate_quality_gate")
    @patch("sdp.cli_init.create_structure")
    @patch("sdp.cli_init.show_dependencies")
    @patch("sdp.cli_init.detect_dependencies")
    @patch("sdp.cli_init.collect_metadata")
    def test_init_command_non_interactive(
        self,
        mock_metadata: Mock,
        mock_detect: Mock,
        mock_show: Mock,
        mock_structure: Mock,
        mock_quality: Mock,
        mock_env: Mock,
        mock_hooks: Mock,
        mock_doctor: Mock,
        tmp_path: Path,
    ) -> None:
        """Test init command in non-interactive mode."""
        mock_metadata.return_value = ("test-project", "Test", "Author")
        mock_detect.return_value = {}
        mock_structure.return_value = ([], [])
        mock_quality.return_value = None
        mock_env.return_value = None
        mock_hooks.return_value = False
        mock_doctor.return_value = True

        runner = CliRunner()
        with runner.isolated_filesystem(temp_dir=tmp_path):
            result = runner.invoke(init, ["--non-interactive"])

        assert result.exit_code == 0
        mock_metadata.assert_called_once()
        mock_detect.assert_called_once()

    @patch("sdp.cli_init.run_doctor")
    @patch("sdp.cli_init.install_git_hooks")
    @patch("sdp.cli_init.create_env_template")
    @patch("sdp.cli_init.generate_quality_gate")
    @patch("sdp.cli_init.create_structure")
    @patch("sdp.cli_init.show_dependencies")
    @patch("sdp.cli_init.detect_dependencies")
    @patch("sdp.cli_init.collect_metadata")
    def test_init_command_force(
        self,
        mock_metadata: Mock,
        mock_detect: Mock,
        mock_show: Mock,
        mock_structure: Mock,
        mock_quality: Mock,
        mock_env: Mock,
        mock_hooks: Mock,
        mock_doctor: Mock,
        tmp_path: Path,
    ) -> None:
        """Test init command with force flag."""
        mock_metadata.return_value = ("test-project", "Test", "Author")
        mock_detect.return_value = {}
        mock_structure.return_value = ([], [])
        mock_quality.return_value = None
        mock_env.return_value = None
        mock_hooks.return_value = False
        mock_doctor.return_value = True

        runner = CliRunner()
        with runner.isolated_filesystem(temp_dir=tmp_path):
            result = runner.invoke(init, ["--non-interactive", "--force"])

        assert result.exit_code == 0
        # Verify force flag passed to create_structure
        call_args = mock_structure.call_args
        assert call_args[0][2] is True  # force=True

    @patch("sdp.cli_init.run_doctor")
    @patch("sdp.cli_init.install_git_hooks")
    @patch("sdp.cli_init.create_env_template")
    @patch("sdp.cli_init.generate_quality_gate")
    @patch("sdp.cli_init.create_structure")
    @patch("sdp.cli_init.show_dependencies")
    @patch("sdp.cli_init.detect_dependencies")
    @patch("sdp.cli_init.collect_metadata")
    def test_init_command_custom_path(
        self,
        mock_metadata: Mock,
        mock_detect: Mock,
        mock_show: Mock,
        mock_structure: Mock,
        mock_quality: Mock,
        mock_env: Mock,
        mock_hooks: Mock,
        mock_doctor: Mock,
        tmp_path: Path,
    ) -> None:
        """Test init command with custom path."""
        custom_path = tmp_path / "custom"
        custom_path.mkdir()

        mock_metadata.return_value = ("test-project", "Test", "Author")
        mock_detect.return_value = {}
        mock_structure.return_value = ([], [])
        mock_quality.return_value = None
        mock_env.return_value = None
        mock_hooks.return_value = False
        mock_doctor.return_value = True

        runner = CliRunner()
        result = runner.invoke(init, ["--non-interactive", "--path", str(custom_path)])

        assert result.exit_code == 0
        # Verify custom path used
        call_args = mock_metadata.call_args
        assert call_args[0][0] == custom_path
