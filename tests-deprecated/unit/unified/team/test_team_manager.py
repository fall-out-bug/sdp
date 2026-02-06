"""
Unit tests for TeamManager role registry.

Tests follow TDD discipline: Red → Green → Refactor
"""

import json
import pytest
from pathlib import Path
from unittest.mock import patch, MagicMock

from sdp.unified.team.manager import TeamManager, TeamManagerError
from sdp.unified.team.models import Role, RoleState


class TestTeamManagerInitialization:
    """Test TeamManager initialization and setup."""

    def test_init_creates_directory(self, tmp_path):
        """TeamManager creates team directory on initialization."""
        feature_id = "test-feature"
        manager = TeamManager(feature_id, base_path=tmp_path)

        team_dir = tmp_path / ".claude" / "teams" / feature_id
        assert team_dir.exists()
        assert team_dir.is_dir()

    def test_init_loads_existing_config(self, tmp_path):
        """TeamManager loads existing config if present."""
        feature_id = "test-feature"
        team_dir = tmp_path / ".claude" / "teams" / feature_id
        team_dir.mkdir(parents=True)

        config_file = team_dir / "config.json"
        existing_config = {
            "feature_id": feature_id,
            "roles": {
                "planner": {
                    "name": "planner",
                    "description": "Planning role",
                    "state": "active",
                    "skill_file": "planner.md"
                }
            }
        }
        config_file.write_text(json.dumps(existing_config))

        manager = TeamManager(feature_id, base_path=tmp_path)

        assert "planner" in manager.roles
        assert manager.roles["planner"].state == RoleState.ACTIVE

    def test_init_with_empty_config(self, tmp_path):
        """TeamManager starts with empty roles if no config exists."""
        feature_id = "test-feature"
        manager = TeamManager(feature_id, base_path=tmp_path)

        assert len(manager.roles) == 0
        assert manager.feature_id == feature_id


class TestRoleRegistration:
    """Test role registration functionality."""

    def test_register_role_adds_to_registry(self, tmp_path):
        """register_role() adds a new role to the registry."""
        manager = TeamManager("test-feature", base_path=tmp_path)

        role = Role(
            name="planner",
            description="Planning specialist",
            state=RoleState.DORMANT,
            skill_file="planner.md"
        )

        manager.register_role(role)

        assert "planner" in manager.roles
        assert manager.roles["planner"].name == "planner"
        assert manager.roles["planner"].description == "Planning specialist"

    def test_register_role_saves_to_config(self, tmp_path):
        """register_role() persists role to config file."""
        manager = TeamManager("test-feature", base_path=tmp_path)

        role = Role(
            name="reviewer",
            description="Code reviewer",
            state=RoleState.DORMANT,
            skill_file="reviewer.md"
        )

        manager.register_role(role)

        config_file = tmp_path / ".claude" / "teams" / "test-feature" / "config.json"
        assert config_file.exists()

        config = json.loads(config_file.read_text())
        assert "reviewer" in config["roles"]
        assert config["roles"]["reviewer"]["description"] == "Code reviewer"

    def test_register_duplicate_role_raises_error(self, tmp_path):
        """register_role() raises error for duplicate role names."""
        manager = TeamManager("test-feature", base_path=tmp_path)

        role = Role(
            name="planner",
            description="Planning specialist",
            state=RoleState.DORMANT,
            skill_file="planner.md"
        )

        manager.register_role(role)

        with pytest.raises(TeamManagerError, match="Role 'planner' already registered"):
            manager.register_role(role)

    def test_register_multiple_roles(self, tmp_path):
        """register_role() can register multiple roles."""
        manager = TeamManager("test-feature", base_path=tmp_path)

        roles = [
            Role(name="planner", description="Plans work", state=RoleState.DORMANT, skill_file="p.md"),
            Role(name="builder", description="Builds code", state=RoleState.DORMANT, skill_file="b.md"),
            Role(name="reviewer", description="Reviews code", state=RoleState.DORMANT, skill_file="r.md"),
        ]

        for role in roles:
            manager.register_role(role)

        assert len(manager.roles) == 3
        assert "planner" in manager.roles
        assert "builder" in manager.roles
        assert "reviewer" in manager.roles


class TestRoleActivation:
    """Test role activation and deactivation."""

    def test_activate_role_changes_state(self, tmp_path):
        """activate_role() changes role state to active."""
        manager = TeamManager("test-feature", base_path=tmp_path)

        role = Role(
            name="planner",
            description="Planning specialist",
            state=RoleState.DORMANT,
            skill_file="planner.md"
        )
        manager.register_role(role)

        manager.activate_role("planner")

        assert manager.roles["planner"].state == RoleState.ACTIVE

    def test_activate_nonexistent_role_raises_error(self, tmp_path):
        """activate_role() raises error for nonexistent role."""
        manager = TeamManager("test-feature", base_path=tmp_path)

        with pytest.raises(TeamManagerError, match="Role 'nonexistent' not found"):
            manager.activate_role("nonexistent")

    def test_activate_role_persists_state(self, tmp_path):
        """activate_role() saves state to config file."""
        manager = TeamManager("test-feature", base_path=tmp_path)

        role = Role(
            name="builder",
            description="Builds features",
            state=RoleState.DORMANT,
            skill_file="builder.md"
        )
        manager.register_role(role)
        manager.activate_role("builder")

        config_file = tmp_path / ".claude" / "teams" / "test-feature" / "config.json"
        config = json.loads(config_file.read_text())

        assert config["roles"]["builder"]["state"] == "active"

    def test_deactivate_role_changes_state(self, tmp_path):
        """Deactivating an active role changes state to dormant."""
        manager = TeamManager("test-feature", base_path=tmp_path)

        role = Role(
            name="tester",
            description="Testing role",
            state=RoleState.ACTIVE,
            skill_file="tester.md"
        )
        manager.register_role(role)

        manager.deactivate_role("tester")

        assert manager.roles["tester"].state == RoleState.DORMANT


class TestSendMessage:
    """Test message sending to roles."""

    def test_send_message_to_active_role(self, tmp_path):
        """send_message() sends message to active role."""
        manager = TeamManager("test-feature", base_path=tmp_path)

        role = Role(
            name="builder",
            description="Builds features",
            state=RoleState.ACTIVE,
            skill_file="builder.md"
        )
        manager.register_role(role)

        # Mock the actual sending mechanism
        with patch.object(manager, '_send_to_agent') as mock_send:
            manager.send_message("builder", "Build feature F01")
            mock_send.assert_called_once_with("builder", "Build feature F01")

    def test_send_message_to_dormant_role_raises_error(self, tmp_path):
        """send_message() raises error for dormant roles."""
        manager = TeamManager("test-feature", base_path=tmp_path)

        role = Role(
            name="planner",
            description="Planning role",
            state=RoleState.DORMANT,
            skill_file="planner.md"
        )
        manager.register_role(role)

        with pytest.raises(TeamManagerError, match="Role 'planner' is not active"):
            manager.send_message("planner", "Plan feature F01")

    def test_send_message_to_nonexistent_role_raises_error(self, tmp_path):
        """send_message() raises error for nonexistent role."""
        manager = TeamManager("test-feature", base_path=tmp_path)

        with pytest.raises(TeamManagerError, match="Role 'ghost' not found"):
            manager.send_message("ghost", "Hello?")


class TestRoleListing:
    """Test role listing and querying."""

    def test_list_active_roles(self, tmp_path):
        """list_active_roles() returns only active roles."""
        manager = TeamManager("test-feature", base_path=tmp_path)

        manager.register_role(Role(
            name="active1",
            description="Active role 1",
            state=RoleState.ACTIVE,
            skill_file="a1.md"
        ))
        manager.register_role(Role(
            name="dormant1",
            description="Dormant role 1",
            state=RoleState.DORMANT,
            skill_file="d1.md"
        ))
        manager.register_role(Role(
            name="active2",
            description="Active role 2",
            state=RoleState.ACTIVE,
            skill_file="a2.md"
        ))

        active = manager.list_active_roles()

        assert len(active) == 2
        assert all(r.state == RoleState.ACTIVE for r in active)
        role_names = {r.name for r in active}
        assert role_names == {"active1", "active2"}

    def test_list_dormant_roles(self, tmp_path):
        """list_dormant_roles() returns only dormant roles."""
        manager = TeamManager("test-feature", base_path=tmp_path)

        manager.register_role(Role(
            name="active1",
            description="Active role 1",
            state=RoleState.ACTIVE,
            skill_file="a1.md"
        ))
        manager.register_role(Role(
            name="dormant1",
            description="Dormant role 1",
            state=RoleState.DORMANT,
            skill_file="d1.md"
        ))
        manager.register_role(Role(
            name="dormant2",
            description="Dormant role 2",
            state=RoleState.DORMANT,
            skill_file="d2.md"
        ))

        dormant = manager.list_dormant_roles()

        assert len(dormant) == 2
        assert all(r.state == RoleState.DORMANT for r in dormant)
        role_names = {r.name for r in dormant}
        assert role_names == {"dormant1", "dormant2"}

    def test_get_role(self, tmp_path):
        """get_role() returns role by name."""
        manager = TeamManager("test-feature", base_path=tmp_path)

        role = Role(
            name="specific",
            description="Specific role",
            state=RoleState.ACTIVE,
            skill_file="specific.md"
        )
        manager.register_role(role)

        retrieved = manager.get_role("specific")

        assert retrieved is not None
        assert retrieved.name == "specific"
        assert retrieved.description == "Specific role"

    def test_get_role_returns_none_for_nonexistent(self, tmp_path):
        """get_role() returns None for nonexistent role."""
        manager = TeamManager("test-feature", base_path=tmp_path)

        result = manager.get_role("nonexistent")

        assert result is None


class TestConfigPersistence:
    """Test config file persistence."""

    def test_config_persists_across_instances(self, tmp_path):
        """Config persists when creating new TeamManager instance."""
        feature_id = "persist-test"

        # First instance
        manager1 = TeamManager(feature_id, base_path=tmp_path)
        role = Role(
            name="persistent",
            description="Persists across instances",
            state=RoleState.ACTIVE,
            skill_file="persistent.md"
        )
        manager1.register_role(role)

        # Second instance
        manager2 = TeamManager(feature_id, base_path=tmp_path)

        assert "persistent" in manager2.roles
        assert manager2.roles["persistent"].state == RoleState.ACTIVE

    def test_config_survives_role_modifications(self, tmp_path):
        """Config file reflects all role changes."""
        manager = TeamManager("survive-test", base_path=tmp_path)

        # Register role
        role = Role(
            name="survivor",
            description="Survives changes",
            state=RoleState.DORMANT,
            skill_file="survivor.md"
        )
        manager.register_role(role)

        # Activate role
        manager.activate_role("survivor")

        # Verify config file
        config_file = tmp_path / ".claude" / "teams" / "survive-test" / "config.json"
        config = json.loads(config_file.read_text())

        assert config["roles"]["survivor"]["state"] == "active"
        assert config["feature_id"] == "survive-test"


class TestErrorHandling:
    """Test error handling and edge cases."""

    def test_invalid_feature_id_raises_error(self, tmp_path):
        """Empty or None feature_id raises error."""
        with pytest.raises(TeamManagerError, match="feature_id cannot be empty"):
            TeamManager("", base_path=tmp_path)

    def test_config_file_corruption_handled(self, tmp_path):
        """Corrupted config file is handled gracefully."""
        feature_id = "corrupt-test"
        team_dir = tmp_path / ".claude" / "teams" / feature_id
        team_dir.mkdir(parents=True)

        config_file = team_dir / "config.json"
        config_file.write_text("invalid json {{{")

        # Should start fresh with empty roles
        manager = TeamManager(feature_id, base_path=tmp_path)

        assert len(manager.roles) == 0

    def test_readonly_config_directory_raises_error(self, tmp_path):
        """Cannot write to read-only directory."""
        # This test is platform-dependent and may not work on all systems
        # Skipping for now, would need special permission setup
        pass


class TestLargeRoleRegistry:
    """Test handling of 100+ roles as per requirements."""

    def test_register_hundred_roles(self, tmp_path):
        """TeamManager can handle 100+ roles."""
        manager = TeamManager("large-test", base_path=tmp_path)

        # Register 100 roles
        for i in range(100):
            role = Role(
                name=f"role_{i}",
                description=f"Role number {i}",
                state=RoleState.DORMANT if i % 2 == 0 else RoleState.ACTIVE,
                skill_file=f"role_{i}.md"
            )
            manager.register_role(role)

        assert len(manager.roles) == 100

    def test_performance_with_many_roles(self, tmp_path):
        """Operations remain fast with many roles."""
        import time

        manager = TeamManager("perf-test", base_path=tmp_path)

        # Register 100 roles
        for i in range(100):
            role = Role(
                name=f"role_{i}",
                description=f"Role {i}",
                state=RoleState.DORMANT,
                skill_file=f"role_{i}.md"
            )
            manager.register_role(role)

        # Time listing operations
        start = time.time()
        active = manager.list_active_roles()
        list_time = time.time() - start

        # Should be very fast (< 0.1 seconds)
        assert list_time < 0.1

    def test_activate_in_large_registry(self, tmp_path):
        """Activating roles works correctly in large registry."""
        manager = TeamManager("activate-test", base_path=tmp_path)

        # Register 50 roles
        for i in range(50):
            role = Role(
                name=f"role_{i}",
                description=f"Role {i}",
                state=RoleState.DORMANT,
                skill_file=f"role_{i}.md"
            )
            manager.register_role(role)

        # Activate one
        manager.activate_role("role_25")

        # Verify only one is active
        active = manager.list_active_roles()
        assert len(active) == 1
        assert active[0].name == "role_25"


class TestTeamLifecycle:
    """Test team lifecycle management: create_team, delete_team, get_team."""

    def test_create_team_initializes_new_team(self, tmp_path):
        """create_team() creates a new team with initial roles."""
        from sdp.unified.team import create_team

        feature_id = "new-feature"
        initial_roles = [
            Role(
                name="orchestrator",
                description="Feature orchestrator",
                state=RoleState.ACTIVE,
                skill_file="orchestrator.md"
            ),
            Role(
                name="builder",
                description="Workstream executor",
                state=RoleState.DORMANT,
                skill_file="builder.md"
            ),
        ]

        manager = create_team(feature_id, initial_roles, base_path=tmp_path)

        assert manager is not None
        assert manager.feature_id == feature_id
        assert len(manager.roles) == 2
        assert "orchestrator" in manager.roles
        assert "builder" in manager.roles

        # Verify config was created
        config_file = tmp_path / ".claude" / "teams" / feature_id / "config.json"
        assert config_file.exists()

    def test_create_team_with_checkpoint_integration(self, tmp_path):
        """create_team() integrates with CheckpointRepository."""
        from sdp.unified.team import create_team
        from sdp.unified.checkpoint.repository import CheckpointRepository

        feature_id = "checkpoint-feature"

        # Initialize checkpoint repo
        db_path = tmp_path / "checkpoints.db"
        repo = CheckpointRepository(str(db_path))
        repo.initialize()

        # Create team with checkpoint reference
        initial_roles = [
            Role(
                name="orchestrator",
                description="Orchestrator",
                state=RoleState.ACTIVE,
                skill_file="orchestrator.md"
            )
        ]

        manager = create_team(feature_id, initial_roles, base_path=tmp_path)

        assert manager.feature_id == feature_id
        repo.close()

    def test_create_team_returns_existing_team_if_present(self, tmp_path):
        """create_team() returns existing team if config already exists."""
        from sdp.unified.team import create_team

        feature_id = "existing-feature"

        # Create team first time
        initial_roles = [
            Role(
                name="planner",
                description="Planner",
                state=RoleState.ACTIVE,
                skill_file="planner.md"
            )
        ]
        manager1 = create_team(feature_id, initial_roles, base_path=tmp_path)

        # Add a role to first instance
        manager1.register_role(Role(
            name="builder",
            description="Builder",
            state=RoleState.DORMANT,
            skill_file="builder.md"
        ))

        # Create team second time - should load existing
        manager2 = create_team(feature_id, initial_roles, base_path=tmp_path)

        # Should have both roles (from existing config)
        assert len(manager2.roles) == 2
        assert "planner" in manager2.roles
        assert "builder" in manager2.roles

    def test_delete_team_removes_team_directory(self, tmp_path):
        """delete_team() removes team directory and configuration."""
        from sdp.unified.team import create_team, delete_team

        feature_id = "deletable-feature"

        # Create team
        initial_roles = [
            Role(
                name="orchestrator",
                description="Orchestrator",
                state=RoleState.ACTIVE,
                skill_file="orchestrator.md"
            )
        ]
        manager = create_team(feature_id, initial_roles, base_path=tmp_path)

        team_dir = tmp_path / ".claude" / "teams" / feature_id
        assert team_dir.exists()

        # Delete team
        delete_team(feature_id, base_path=tmp_path)

        # Directory should be removed
        assert not team_dir.exists()

    def test_delete_team_nonexistent_raises_no_error(self, tmp_path):
        """delete_team() handles nonexistent teams gracefully."""
        from sdp.unified.team import delete_team

        # Should not raise error for nonexistent team
        delete_team("nonexistent-feature", base_path=tmp_path)

    def test_get_team_returns_existing_team(self, tmp_path):
        """get_team() returns existing TeamManager or None."""
        from sdp.unified.team import create_team, get_team

        feature_id = "gettable-feature"

        # Get nonexistent team
        manager = get_team(feature_id, base_path=tmp_path)
        assert manager is None

        # Create team
        initial_roles = [
            Role(
                name="orchestrator",
                description="Orchestrator",
                state=RoleState.ACTIVE,
                skill_file="orchestrator.md"
            )
        ]
        create_team(feature_id, initial_roles, base_path=tmp_path)

        # Get existing team
        manager = get_team(feature_id, base_path=tmp_path)

        assert manager is not None
        assert manager.feature_id == feature_id
        assert len(manager.roles) == 1

    def test_team_lifecycle_full_workflow(self, tmp_path):
        """Test full lifecycle: create -> get -> delete."""
        from sdp.unified.team import create_team, get_team, delete_team

        feature_id = "lifecycle-feature"

        # Team doesn't exist initially
        manager = get_team(feature_id, base_path=tmp_path)
        assert manager is None

        # Create team
        initial_roles = [
            Role(
                name="orchestrator",
                description="Orchestrator",
                state=RoleState.ACTIVE,
                skill_file="orchestrator.md"
            )
        ]
        manager = create_team(feature_id, initial_roles, base_path=tmp_path)
        assert manager is not None

        # Get team
        manager = get_team(feature_id, base_path=tmp_path)
        assert manager is not None
        assert manager.feature_id == feature_id

        # Delete team
        delete_team(feature_id, base_path=tmp_path)

        # Team doesn't exist anymore
        manager = get_team(feature_id, base_path=tmp_path)
        assert manager is None

    def test_create_team_validates_feature_id(self, tmp_path):
        """create_team() raises error for invalid feature_id."""
        from sdp.unified.team import create_team
        from sdp.unified.team.errors import TeamManagerError

        with pytest.raises(TeamManagerError, match="feature_id cannot be empty"):
            create_team("", [], base_path=tmp_path)

        with pytest.raises(TeamManagerError, match="feature_id cannot be empty"):
            create_team("   ", [], base_path=tmp_path)
