"""Individual health check implementations."""

import shutil
import sys
from pathlib import Path
from typing import List

from .base import HealthCheck, HealthCheckResult


class PythonVersionCheck(HealthCheck):
    """Check Python version meets minimum requirements."""

    MIN_VERSION = (3, 10)

    def __init__(self) -> None:
        super().__init__("Python Version", critical=True)

    def run(self) -> HealthCheckResult:
        """Check Python version."""
        current_version = sys.version_info[:2]

        if current_version >= self.MIN_VERSION:
            message = (
                f"Python {current_version[0]}.{current_version[1]} "
                f"(>= {self.MIN_VERSION[0]}.{self.MIN_VERSION[1]})"
            )
            return HealthCheckResult(
                name=self.name,
                passed=True,
                message=message,
            )

        message = (
            f"Python {current_version[0]}.{current_version[1]} "
            f"(< {self.MIN_VERSION[0]}.{self.MIN_VERSION[1]})"
        )
        return HealthCheckResult(
            name=self.name,
            passed=False,
            message=message,
            remediation=f"Upgrade to Python {self.MIN_VERSION[0]}.{self.MIN_VERSION[1]}+",
        )


class PoetryCheck(HealthCheck):
    """Check Poetry is installed."""

    def __init__(self) -> None:
        super().__init__("Poetry", critical=True)

    def run(self) -> HealthCheckResult:
        """Check Poetry installation."""
        poetry_path = shutil.which("poetry")

        if poetry_path:
            return HealthCheckResult(
                name=self.name,
                passed=True,
                message=f"Poetry installed at {poetry_path}",
            )

        return HealthCheckResult(
            name=self.name,
            passed=False,
            message="Poetry not found",
            remediation="Install Poetry: curl -sSL https://install.python-poetry.org | python3 -",
        )


class GitHooksCheck(HealthCheck):
    """Check git hooks are configured."""

    def __init__(self) -> None:
        super().__init__("Git Hooks", critical=True)

    def run(self) -> HealthCheckResult:
        """Check git hooks directory exists."""
        cwd = Path.cwd()
        hooks_dir = cwd / ".git" / "hooks"

        if hooks_dir.exists():
            return HealthCheckResult(
                name=self.name,
                passed=True,
                message=f"Git hooks directory found at {hooks_dir}",
            )

        return HealthCheckResult(
            name=self.name,
            passed=False,
            message="Git hooks directory not found",
            remediation="Initialize git repository: git init",
        )


class BeadsCLICheck(HealthCheck):
    """Check Beads CLI is installed (optional)."""

    def __init__(self) -> None:
        super().__init__("Beads CLI", critical=False)

    def run(self) -> HealthCheckResult:
        """Check Beads CLI installation with version and Go detection."""
        import subprocess

        # Check if Go is installed
        go_path = shutil.which("go")
        go_installed = go_path is not None

        # Check if bd CLI is installed
        beads_path = shutil.which("bd")

        if beads_path:
            # Get version
            version = self._get_version(beads_path)
            return HealthCheckResult(
                name=self.name,
                passed=True,
                message=f"Beads CLI installed (v{version}) at {beads_path}",
            )

        # Build remediation message
        remediation_steps = []
        if not go_installed:
            remediation_steps.append("1. Install Go: brew install go (macOS) or apt install golang-go (Linux)")
        remediation_steps.append("2. Install Beads: go install github.com/steveyegge/beads/cmd/bd@latest")
        remediation_steps.append("3. Add to PATH: export PATH=$PATH:$(go env GOPATH)/bin")
        remediation_steps.append("4. Verify: bd --version")
        remediation_steps.append("")
        remediation_steps.append("See docs/setup/beads-installation.md for detailed instructions")

        return HealthCheckResult(
            name=self.name,
            passed=True,  # Optional, so still True
            message="Beads CLI not installed (optional)",
            remediation="\n".join(remediation_steps),
        )

    def _get_version(self, bd_path: str) -> str:
        """Get Beads CLI version.

        Args:
            bd_path: Path to bd executable

        Returns:
            Version string or "unknown"
        """
        import subprocess

        try:
            result = subprocess.run(
                [bd_path, "--version"],
                capture_output=True,
                text=True,
                timeout=5,
            )
            # Parse version from output like "bd version 0.1.0"
            output = result.stdout.strip()
            if output:
                # Try to extract version number
                parts = output.split()
                for part in parts:
                    if part and part[0].isdigit():
                        return part
            return output or "unknown"
        except (subprocess.TimeoutExpired, subprocess.SubprocessError, OSError):
            return "unknown"


class GitHubCLICheck(HealthCheck):
    """Check GitHub CLI is installed (optional)."""

    def __init__(self) -> None:
        super().__init__("GitHub CLI", critical=False)

    def run(self) -> HealthCheckResult:
        """Check GitHub CLI installation."""
        gh_path = shutil.which("gh")

        if gh_path:
            return HealthCheckResult(
                name=self.name,
                passed=True,
                message=f"GitHub CLI installed at {gh_path}",
            )

        return HealthCheckResult(
            name=self.name,
            passed=True,  # Optional, so still True
            message="GitHub CLI not installed (optional)",
            remediation="Install GitHub CLI: https://cli.github.com/",
        )


class TelegramConfigCheck(HealthCheck):
    """Check Telegram configuration (optional)."""

    def __init__(self) -> None:
        super().__init__("Telegram Config", critical=False)

    def run(self) -> HealthCheckResult:
        """Check Telegram configuration."""
        cwd = Path.cwd()
        env_file = cwd / ".env"

        if not env_file.exists():
            return HealthCheckResult(
                name=self.name,
                passed=True,  # Optional, so still True
                message="Telegram not configured (.env not found, optional)",
                remediation="Create .env with TELEGRAM_TOKEN and TELEGRAM_CHAT_ID",
            )

        env_content = env_file.read_text()

        if "TELEGRAM_TOKEN" in env_content and "TELEGRAM_CHAT_ID" in env_content:
            return HealthCheckResult(
                name=self.name,
                passed=True,
                message="Telegram configured in .env",
            )

        return HealthCheckResult(
            name=self.name,
            passed=True,  # Optional, so still True
            message="Telegram not fully configured (optional)",
            remediation="Add TELEGRAM_TOKEN and TELEGRAM_CHAT_ID to .env",
        )


def get_health_checks() -> List[HealthCheck]:
    """Get list of all health checks.

    Returns:
        List of HealthCheck instances
    """
    return [
        PythonVersionCheck(),
        PoetryCheck(),
        GitHooksCheck(),
        BeadsCLICheck(),
        GitHubCLICheck(),
        TelegramConfigCheck(),
    ]
