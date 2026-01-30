"""Beads client factory and public API.

This module provides the factory function for creating Beads client instances
and re-exports the public API.
"""

from pathlib import Path
from typing import Optional

from .base import BeadsClient
from .cli import CLIBeadsClient
from .exceptions import BeadsClientError, BeadsNotInstalledError
from .mock import MockBeadsClient

__all__ = [
    "BeadsClient",
    "BeadsClientError",
    "BeadsNotInstalledError",
    "MockBeadsClient",
    "CLIBeadsClient",
    "create_beads_client",
]


def create_beads_client(
    use_mock: bool = False, project_dir: Optional[Path] = None
) -> BeadsClient:
    """Factory function to create appropriate Beads client.

    Default behavior (as of v0.6.0):
    - If bd CLI installed → Returns CLIBeadsClient (real Beads)
    - If bd not installed → Returns MockBeadsClient with warning
    - Override with BEADS_USE_MOCK=true for tests

    Args:
        use_mock: Force mock client (for testing)
        project_dir: Project directory (for CLI client)

    Returns:
        BeadsClient instance

    Example:
        # Auto-detect (real if bd installed)
        client = create_beads_client()

        # Force mock for testing
        client = create_beads_client(use_mock=True)

        # Use real Beads in specific directory
        client = create_beads_client(project_dir=Path("/my/project"))

    The factory determines the client type as follows:
    1. If use_mock=True, returns MockBeadsClient
    2. If BEADS_USE_MOCK env var is "true", returns MockBeadsClient
    3. If bd CLI installed, returns CLIBeadsClient
    4. Otherwise, returns MockBeadsClient with warning
    """
    import os
    import shutil
    import warnings

    # Check if mock is explicitly requested
    if use_mock:
        return MockBeadsClient()

    # Check environment variable
    env_mock = os.getenv("BEADS_USE_MOCK")
    if env_mock is not None:
        if env_mock.lower() == "true":
            return MockBeadsClient()

    # Auto-detect: use real if bd CLI is available
    bd_available = shutil.which("bd") is not None

    if not bd_available:
        warnings.warn(
            "Beads CLI (bd) not found. Using mock client. "
            "Install with: go install github.com/steveyegge/beads/cmd/bd@latest"
        )
        return MockBeadsClient()

    # Use real Beads CLI client
    return CLIBeadsClient(project_dir)
