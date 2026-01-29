"""Beads client factory and public API.

This module provides the factory function for creating Beads client instances
and re-exports the public API.
"""

from pathlib import Path
from typing import Optional

from .base import BeadsClient
from .cli import CLIBeadsClient
from .exceptions import BeadsClientError
from .mock import MockBeadsClient

__all__ = [
    "BeadsClient",
    "BeadsClientError",
    "MockBeadsClient",
    "CLIBeadsClient",
    "create_beads_client",
]


def create_beads_client(
    use_mock: bool = False, project_dir: Optional[Path] = None
) -> BeadsClient:
    """Factory function to create appropriate Beads client.

    Args:
        use_mock: Force mock client (for testing)
        project_dir: Project directory (for CLI client)

    Returns:
        BeadsClient instance

    Example:
        # Use mock for testing
        client = create_beads_client(use_mock=True)

        # Use real Beads (must be installed)
        client = create_beads_client(project_dir=Path("/my/project"))

    The factory determines the client type as follows:
    1. If use_mock=True, returns MockBeadsClient
    2. If BEADS_USE_MOCK env var is "true", returns MockBeadsClient
    3. Otherwise, returns CLIBeadsClient (requires Beads installation)
    """
    # Check environment variable
    import os

    if use_mock or os.getenv("BEADS_USE_MOCK", "false").lower() == "true":
        return MockBeadsClient()

    return CLIBeadsClient(project_dir)
