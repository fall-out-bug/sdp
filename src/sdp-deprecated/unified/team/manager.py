"""
TeamManager role registry for managing 100+ agent roles.

Provides role registration, activation/deactivation, and message routing
for coordinated multi-agent workflows.
"""

import logging
from pathlib import Path
from typing import Optional

from sdp.unified.team.errors import TeamManagerError
from sdp.unified.team.models import Role, RoleState
from sdp.unified.team.persistence import TeamConfigStore

logger = logging.getLogger(__name__)


class TeamManager:
    """Manager for agent role registry and lifecycle.

    Maintains a registry of agent roles with active/dormant states,
    persists configuration to ~/.claude/teams/{feature_id}/config.json,
    and handles message routing to active agents.

    Attributes:
        feature_id: Unique feature identifier
        base_path: Base path for team configurations
        roles: Dictionary of registered roles {name: Role}
    """

    def __init__(self, feature_id: str, base_path: Optional[Path] = None) -> None:
        """Initialize TeamManager for a feature.

        Args:
            feature_id: Unique feature identifier
            base_path: Base path for team configs (defaults to home directory)

        Raises:
            TeamManagerError: If feature_id is empty
        """
        if not feature_id or not feature_id.strip():
            raise TeamManagerError("feature_id cannot be empty")

        self.feature_id = feature_id
        self.base_path = base_path or Path.home()
        self.roles: dict[str, Role] = {}

        # Create team directory structure
        self._team_dir = self.base_path / ".claude" / "teams" / feature_id
        self._team_dir.mkdir(parents=True, exist_ok=True)

        self._config_file = self._team_dir / "config.json"
        self._config_store = TeamConfigStore(self._config_file)

        # Load existing configuration if present
        self._load_config()

        logger.info(f"TeamManager initialized for feature '{feature_id}' at {self._team_dir}")

    def register_role(self, role: Role) -> None:
        """Register a new agent role.

        Args:
            role: Role to register

        Raises:
            TeamManagerError: If role already registered
        """
        if role.name in self.roles:
            raise TeamManagerError(f"Role '{role.name}' already registered")

        self.roles[role.name] = role
        self._config_store.save(self.feature_id, self.roles)

        logger.info(f"Role '{role.name}' registered in state {role.state.value}")

    def activate_role(self, role_name: str) -> None:
        """Activate a dormant role.

        Args:
            role_name: Name of role to activate

        Raises:
            TeamManagerError: If role not found
        """
        if role_name not in self.roles:
            raise TeamManagerError(f"Role '{role_name}' not found")

        self.roles[role_name].state = RoleState.ACTIVE
        self._config_store.save(self.feature_id, self.roles)

        logger.info(f"Role '{role_name}' activated")

    def deactivate_role(self, role_name: str) -> None:
        """Deactivate an active role.

        Args:
            role_name: Name of role to deactivate

        Raises:
            TeamManagerError: If role not found
        """
        if role_name not in self.roles:
            raise TeamManagerError(f"Role '{role_name}' not found")

        self.roles[role_name].state = RoleState.DORMANT
        self._config_store.save(self.feature_id, self.roles)

        logger.info(f"Role '{role_name}' deactivated")

    def send_message(self, role_name: str, message: str) -> None:
        """Send message to an active role.

        Args:
            role_name: Name of role to send message to
            message: Message content

        Raises:
            TeamManagerError: If role not found or not active
        """
        if role_name not in self.roles:
            raise TeamManagerError(f"Role '{role_name}' not found")

        role = self.roles[role_name]

        if role.state != RoleState.ACTIVE:
            raise TeamManagerError(f"Role '{role_name}' is not active")

        # Delegate to actual agent communication
        self._send_to_agent(role_name, message)

    def list_active_roles(self) -> list[Role]:
        """List all active roles.

        Returns:
            List of active Role objects
        """
        return [role for role in self.roles.values() if role.state == RoleState.ACTIVE]

    def list_dormant_roles(self) -> list[Role]:
        """List all dormant roles.

        Returns:
            List of dormant Role objects
        """
        return [role for role in self.roles.values() if role.state == RoleState.DORMANT]

    def get_role(self, role_name: str) -> Optional[Role]:
        """Get role by name.

        Args:
            role_name: Name of role to retrieve

        Returns:
            Role object or None if not found
        """
        return self.roles.get(role_name)

    def _send_to_agent(self, role_name: str, message: str) -> None:
        """Send message to actual agent via communication channel.

        This is a placeholder for the actual agent communication mechanism
        that will be implemented in WS-013 (SendMessage router).

        Args:
            role_name: Target role name
            message: Message content
        """
        logger.debug(f"Sending message to '{role_name}': {message}")
        # Actual implementation will use Task tool or agent communication
        # This is a stub for now

    def _load_config(self) -> None:
        """Load configuration from config file."""
        self.roles = self._config_store.load()
