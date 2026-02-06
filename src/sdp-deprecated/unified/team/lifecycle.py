"""
Team lifecycle operations for TeamManager.

Provides high-level functions for creating, deleting, and retrieving
team configurations.
"""

import logging
import shutil
from pathlib import Path
from typing import Optional

from sdp.unified.team.errors import TeamManagerError
from sdp.unified.team.manager import TeamManager
from sdp.unified.team.models import Role

logger = logging.getLogger(__name__)


def create_team(
    feature_id: str,
    initial_roles: list[Role],
    base_path: Optional[Path] = None,
) -> TeamManager:
    """Create a new team or load existing team configuration.

    Creates a TeamManager instance with the specified feature_id and
    initial roles. If a team configuration already exists for the
    feature_id, loads the existing configuration instead.

    Args:
        feature_id: Unique feature identifier
        initial_roles: List of initial roles to register
        base_path: Base path for team configs (defaults to home directory)

    Returns:
        TeamManager instance for the feature

    Raises:
        TeamManagerError: If feature_id is empty or invalid
    """
    if not feature_id or not feature_id.strip():
        raise TeamManagerError("feature_id cannot be empty")

    base_path = base_path or Path.home()
    team_dir = base_path / ".claude" / "teams" / feature_id

    # Check if team already exists
    config_file = team_dir / "config.json"
    team_exists = config_file.exists()

    # Create or load team
    manager = TeamManager(feature_id, base_path=base_path)

    # Only register initial roles if team is new
    if not team_exists:
        for role in initial_roles:
            manager.register_role(role)
        logger.info(f"Created new team for feature '{feature_id}' with {len(initial_roles)} roles")
    else:
        logger.info(f"Loaded existing team for feature '{feature_id}'")

    return manager


def delete_team(feature_id: str, base_path: Optional[Path] = None) -> None:
    """Delete team directory and configuration.

    Removes the team configuration directory and all its contents.
    If the team doesn't exist, the function returns silently.

    Args:
        feature_id: Feature identifier of team to delete
        base_path: Base path for team configs (defaults to home directory)
    """
    base_path = base_path or Path.home()
    team_dir = base_path / ".claude" / "teams" / feature_id

    if team_dir.exists():
        shutil.rmtree(team_dir)
        logger.info(f"Deleted team for feature '{feature_id}'")
    else:
        logger.debug(f"Team for feature '{feature_id}' does not exist, skipping deletion")


def get_team(feature_id: str, base_path: Optional[Path] = None) -> Optional[TeamManager]:
    """Get existing team or return None if not found.

    Attempts to load an existing team configuration. If the team
    configuration doesn't exist, returns None instead of raising
    an error.

    Args:
        feature_id: Feature identifier to load
        base_path: Base path for team configs (defaults to home directory)

    Returns:
        TeamManager instance if team exists, None otherwise
    """
    base_path = base_path or Path.home()
    team_dir = base_path / ".claude" / "teams" / feature_id
    config_file = team_dir / "config.json"

    if not config_file.exists():
        return None

    try:
        manager = TeamManager(feature_id, base_path=base_path)
        logger.info(f"Retrieved existing team for feature '{feature_id}'")
        return manager
    except Exception as e:
        logger.warning(f"Failed to load team for feature '{feature_id}': {e}")
        return None
