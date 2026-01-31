"""Beads CLI health check."""

import shutil
import subprocess
from dataclasses import dataclass
from typing import Optional

from .base import HealthCheck, HealthCheckResult


@dataclass
class BeadsHealthData:
    """Beads installation status data."""

    installed: bool
    version: Optional[str]
    go_installed: bool
    path: Optional[str]


class BeadsHealthCheck(HealthCheck):
    """Check Beads CLI availability and installation."""

    def __init__(self) -> None:
        """Initialize Beads health check."""
        super().__init__(name="beads", critical=False)

    def run(self) -> HealthCheckResult:
        """Check if Beads CLI is installed and accessible.

        Returns:
            HealthCheckResult with installation status and remediation
        """
        result = self._check_beads()

        if result.installed:
            return HealthCheckResult(
                name=self.name,
                passed=True,
                message=f"Beads CLI v{result.version} at {result.path}",
            )

        # Build remediation message
        remediation_steps = []
        if not result.go_installed:
            remediation_steps.append(
                "1. Install Go: brew install go (macOS) or apt install golang-go (Linux)"
            )
        remediation_steps.append(
            "2. Install Beads: go install github.com/steveyegge/beads/cmd/bd@latest"
        )
        remediation_steps.append("3. Add to PATH: export PATH=$PATH:$(go env GOPATH)/bin")
        remediation_steps.append("4. Verify: bd --version")
        remediation_steps.append("")
        remediation_steps.append("See docs/setup/beads-installation.md for detailed instructions")

        return HealthCheckResult(
            name=self.name,
            passed=False,
            message="Beads CLI not found",
            remediation="\n".join(remediation_steps),
        )

    def _check_beads(self) -> BeadsHealthData:
        """Check Beads installation details.

        Returns:
            BeadsHealthData with installation status
        """
        # Check Go installation
        go_installed = shutil.which("go") is not None

        # Check bd CLI
        bd_path = shutil.which("bd")
        if not bd_path:
            return BeadsHealthData(
                installed=False,
                version=None,
                go_installed=go_installed,
                path=None,
            )

        # Get version
        version = self._get_version(bd_path)

        return BeadsHealthData(
            installed=True,
            version=version,
            go_installed=go_installed,
            path=bd_path,
        )

    def _get_version(self, bd_path: str) -> str:
        """Get Beads CLI version.

        Args:
            bd_path: Path to bd executable

        Returns:
            Version string or "unknown"
        """
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
                    if part[0].isdigit():
                        return part
            return output or "unknown"
        except (subprocess.TimeoutExpired, subprocess.SubprocessError, OSError):
            return "unknown"
