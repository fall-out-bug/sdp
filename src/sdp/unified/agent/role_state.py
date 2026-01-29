"""RoleStateManager for tracking active/dormant agent roles.

Provides state tracking for agent roles with activation/deactivation
and querying capabilities.
"""

import logging
from enum import Enum
from typing import Optional

logger = logging.getLogger(__name__)


class RoleState(Enum):
    """Role activation state."""

    ACTIVE = "active"
    DORMANT = "dormant"


class RoleStateManager:
    """Manager for agent role states.

    Tracks active and dormant roles with activation/deactivation
    functionality and state querying.
    """

    def __init__(self) -> None:
        """Initialize role state manager with empty state tracking."""
        self._role_states: dict[str, RoleState] = {}

    def activate_role(self, role_name: str) -> RoleState:
        """Activate role.

        Args:
            role_name: Name of role to activate

        Returns:
            RoleState.ACTIVE
        """
        self._role_states[role_name] = RoleState.ACTIVE
        logger.info(f"Activated role: {role_name}")
        return RoleState.ACTIVE

    def deactivate_role(self, role_name: str) -> RoleState:
        """Deactivate role.

        Args:
            role_name: Name of role to deactivate

        Returns:
            RoleState.DORMANT
        """
        self._role_states[role_name] = RoleState.DORMANT
        logger.info(f"Deactivated role: {role_name}")
        return RoleState.DORMANT

    def get_state(self, role_name: str) -> Optional[RoleState]:
        """Get role state.

        Args:
            role_name: Name of role to query

        Returns:
            RoleState if tracked, None otherwise
        """
        return self._role_states.get(role_name)

    def is_active(self, role_name: str) -> bool:
        """Check if role is active.

        Args:
            role_name: Name of role to check

        Returns:
            True if role is active, False otherwise
        """
        state = self._role_states.get(role_name)
        return state == RoleState.ACTIVE

    def is_dormant(self, role_name: str) -> bool:
        """Check if role is dormant.

        Args:
            role_name: Name of role to check

        Returns:
            True if role is dormant, False otherwise
        """
        state = self._role_states.get(role_name)
        return state == RoleState.DORMANT

    def list_active(self) -> list[str]:
        """List all active roles.

        Returns:
            List of active role names
        """
        return [
            name for name, state in self._role_states.items()
            if state == RoleState.ACTIVE
        ]

    def list_dormant(self) -> list[str]:
        """List all dormant roles.

        Returns:
            List of dormant role names
        """
        return [
            name for name, state in self._role_states.items()
            if state == RoleState.DORMANT
        ]

    def list_all(self) -> list[str]:
        """List all tracked roles.

        Returns:
            List of all tracked role names
        """
        return list(self._role_states.keys())

    def clear_all(self) -> None:
        """Clear all role states."""
        self._role_states.clear()
        logger.info("Cleared all role states")
