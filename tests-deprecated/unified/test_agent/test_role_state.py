"""Tests for RoleStateManager module.

Tests role state tracking, activation/deactivation, and state
transitions for managing active and dormant agent roles.
"""

import pytest
from enum import Enum

from sdp.unified.agent.role_state import RoleStateManager, RoleState


class TestRoleStateEnum:
    """Test RoleState enum."""

    def test_has_active_state(self):
        """Should have ACTIVE state."""
        assert RoleState.ACTIVE is not None
        assert RoleState.ACTIVE.value == "active"

    def test_has_dormant_state(self):
        """Should have DORMANT state."""
        assert RoleState.DORMANT is not None
        assert RoleState.DORMANT.value == "dormant"


class TestRoleStateManagerInit:
    """Test RoleStateManager initialization."""

    def test_creates_manager(self):
        """Should initialize manager."""
        manager = RoleStateManager()

        assert manager is not None
        assert hasattr(manager, 'activate_role')
        assert hasattr(manager, 'deactivate_role')
        assert hasattr(manager, 'get_state')

    def test_initializes_empty_state(self):
        """Should initialize with empty role state tracking."""
        manager = RoleStateManager()

        assert manager.get_state("planner") is None


class TestRoleActivation:
    """Test role activation functionality."""

    def test_activates_role(self):
        """Should activate role and return ACTIVE state."""
        manager = RoleStateManager()

        state = manager.activate_role("planner")

        assert state == RoleState.ACTIVE
        assert manager.get_state("planner") == RoleState.ACTIVE

    def test_prevents_duplicate_active_roles(self):
        """Should prevent same role from being activated twice."""
        manager = RoleStateManager()

        manager.activate_role("planner")
        state = manager.activate_role("planner")

        # Should still be ACTIVE, not error
        assert state == RoleState.ACTIVE
        assert manager.get_state("planner") == RoleState.ACTIVE

    def test_activates_multiple_roles(self):
        """Should activate multiple different roles."""
        manager = RoleStateManager()

        manager.activate_role("planner")
        manager.activate_role("builder")
        manager.activate_role("tester")

        assert manager.get_state("planner") == RoleState.ACTIVE
        assert manager.get_state("builder") == RoleState.ACTIVE
        assert manager.get_state("tester") == RoleState.ACTIVE


class TestRoleDeactivation:
    """Test role deactivation functionality."""

    def test_deactivates_active_role(self):
        """Should deactivate active role."""
        manager = RoleStateManager()
        manager.activate_role("planner")

        state = manager.deactivate_role("planner")

        assert state == RoleState.DORMANT
        assert manager.get_state("planner") == RoleState.DORMANT

    def test_deactivating_dormant_role_stays_dormant(self):
        """Should keep dormant role dormant."""
        manager = RoleStateManager()

        state = manager.deactivate_role("planner")

        assert state == RoleState.DORMANT
        assert manager.get_state("planner") == RoleState.DORMANT

    def test_activating_after_deactivation(self):
        """Should allow reactivation after deactivation."""
        manager = RoleStateManager()

        manager.activate_role("planner")
        assert manager.get_state("planner") == RoleState.ACTIVE

        manager.deactivate_role("planner")
        assert manager.get_state("planner") == RoleState.DORMANT

        manager.activate_role("planner")
        assert manager.get_state("planner") == RoleState.ACTIVE


class TestStateQuerying:
    """Test state querying functionality."""

    def test_gets_active_role_state(self):
        """Should return ACTIVE for active role."""
        manager = RoleStateManager()
        manager.activate_role("planner")

        state = manager.get_state("planner")

        assert state == RoleState.ACTIVE

    def test_gets_dormant_role_state(self):
        """Should return DORMANT for deactivated role."""
        manager = RoleStateManager()
        manager.deactivate_role("planner")

        state = manager.get_state("planner")

        assert state == RoleState.DORMANT

    def test_returns_none_for_unknown_role(self):
        """Should return None for role with no state."""
        manager = RoleStateManager()

        state = manager.get_state("unknown")

        assert state is None

    def test_is_active_returns_true_for_active(self):
        """Should return True for active role."""
        manager = RoleStateManager()
        manager.activate_role("planner")

        assert manager.is_active("planner") is True

    def test_is_active_returns_false_for_dormant(self):
        """Should return False for dormant role."""
        manager = RoleStateManager()
        manager.deactivate_role("planner")

        assert manager.is_active("planner") is False

    def test_is_active_returns_false_for_unknown(self):
        """Should return False for unknown role."""
        manager = RoleStateManager()

        assert manager.is_active("unknown") is False

    def test_is_dormant_returns_true_for_dormant(self):
        """Should return True for dormant role."""
        manager = RoleStateManager()
        manager.deactivate_role("planner")

        assert manager.is_dormant("planner") is True

    def test_is_dormant_returns_false_for_active(self):
        """Should return False for active role."""
        manager = RoleStateManager()
        manager.activate_role("planner")

        assert manager.is_dormant("planner") is False


class TestRoleListing:
    """Test listing roles by state."""

    def test_lists_active_roles(self):
        """Should list all active roles."""
        manager = RoleStateManager()

        manager.activate_role("planner")
        manager.activate_role("builder")
        manager.deactivate_role("tester")

        active = manager.list_active()

        assert set(active) == {"planner", "builder"}

    def test_lists_dormant_roles(self):
        """Should list all dormant roles."""
        manager = RoleStateManager()

        manager.activate_role("planner")
        manager.deactivate_role("builder")
        manager.deactivate_role("tester")

        dormant = manager.list_dormant()

        assert set(dormant) == {"builder", "tester"}

    def test_lists_all_roles(self):
        """Should list all tracked roles."""
        manager = RoleStateManager()

        manager.activate_role("planner")
        manager.deactivate_role("builder")
        manager.activate_role("tester")

        all_roles = manager.list_all()

        assert set(all_roles) == {"planner", "builder", "tester"}

    def test_returns_empty_list_when_no_roles(self):
        """Should return empty list when no roles tracked."""
        manager = RoleStateManager()

        assert manager.list_active() == []
        assert manager.list_dormant() == []
        assert manager.list_all() == []


class TestStateTracking:
    """Test state tracking and transitions."""

    def test_tracks_state_transitions(self):
        """Should track role state transitions."""
        manager = RoleStateManager()

        # Initial state: None
        assert manager.get_state("planner") is None

        # Activate
        manager.activate_role("planner")
        assert manager.get_state("planner") == RoleState.ACTIVE

        # Deactivate
        manager.deactivate_role("planner")
        assert manager.get_state("planner") == RoleState.DORMANT

        # Reactivate
        manager.activate_role("planner")
        assert manager.get_state("planner") == RoleState.ACTIVE

    def test_clear_all_states(self):
        """Should clear all role states."""
        manager = RoleStateManager()

        manager.activate_role("planner")
        manager.activate_role("builder")
        manager.deactivate_role("tester")

        manager.clear_all()

        assert manager.get_state("planner") is None
        assert manager.get_state("builder") is None
        assert manager.get_state("tester") is None


class TestEdgeCases:
    """Test edge cases and error handling."""

    def test_case_sensitive_role_names(self):
        """Should treat role names as case-sensitive."""
        manager = RoleStateManager()

        manager.activate_role("Planner")
        manager.activate_role("planner")

        # Should be different roles
        assert manager.get_state("Planner") == RoleState.ACTIVE
        assert manager.get_state("planner") == RoleState.ACTIVE
        assert len(manager.list_all()) == 2

    def test_empty_role_name(self):
        """Should handle empty role name gracefully."""
        manager = RoleStateManager()

        state = manager.activate_role("")

        assert state == RoleState.ACTIVE
        assert manager.get_state("") == RoleState.ACTIVE

    def test_special_characters_in_role_name(self):
        """Should handle special characters in role names."""
        manager = RoleStateManager()

        state = manager.activate_role("role-with_special.chars")

        assert state == RoleState.ACTIVE
        assert manager.get_state("role-with_special.chars") == RoleState.ACTIVE
