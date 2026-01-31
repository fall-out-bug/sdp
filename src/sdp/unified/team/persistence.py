"""
Configuration persistence for TeamManager.

Handles loading and saving team configuration to JSON files.
"""

import json
import logging
from pathlib import Path

from sdp.unified.team.models import Role

logger = logging.getLogger(__name__)


class TeamConfigStore:
    """Handles configuration persistence for team roles.

    Manages loading and saving role configurations to/from
    ~/.claude/teams/{feature_id}/config.json files.
    """

    def __init__(self, config_file: Path) -> None:
        """Initialize config store.

        Args:
            config_file: Path to config.json file
        """
        self._config_file = config_file

    def load(self) -> dict[str, Role]:
        """Load roles from configuration file.

        Returns:
            Dictionary of role_name -> Role

        Handles corrupted files gracefully by returning empty dict.
        """
        if not self._config_file.exists():
            logger.debug("No existing config file found")
            return {}

        try:
            data = json.loads(self._config_file.read_text())
            roles: dict[str, Role] = {}

            # Load roles from config
            for role_data in data.get("roles", {}).values():
                role = Role.from_dict(role_data)
                roles[role.name] = role

            logger.debug(f"Loaded {len(roles)} roles from config")
            return roles

        except (json.JSONDecodeError, KeyError, ValueError) as e:
            logger.warning(f"Failed to load config: {e}. Starting with empty roles.")
            return {}

    def save(self, feature_id: str, roles: dict[str, Role]) -> None:
        """Save roles to configuration file.

        Args:
            feature_id: Feature identifier
            roles: Dictionary of role_name -> Role
        """
        config = {
            "feature_id": feature_id,
            "roles": {name: role.to_dict() for name, role in roles.items()},
        }

        self._config_file.write_text(json.dumps(config, indent=2))
        logger.debug(f"Saved config with {len(roles)} roles")
