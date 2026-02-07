"""Tests for RoleLoader module.

Tests role loading from filesystem, caching, and role lookup
functionality for agent prompt management.
"""

import pytest
from pathlib import Path
from unittest.mock import Mock, patch, MagicMock

from sdp.unified.agent.role_loader import RoleLoader, Role


class TestRoleDataclass:
    """Test Role dataclass."""

    def test_create_role_with_required_fields(self):
        """Should create role with required fields."""
        role = Role(
            name="planner",
            prompt="You are a planning agent.",
        )

        assert role.name == "planner"
        assert role.prompt == "You are a planning agent."

    def test_create_role_with_optional_fields(self):
        """Should create role with optional fields."""
        role = Role(
            name="builder",
            prompt="Build features",
            description="Feature implementation agent",
            capabilities=["build", "test", "refactor"],
        )

        assert role.description == "Feature implementation agent"
        assert role.capabilities == ["build", "test", "refactor"]

    def test_role_defaults(self):
        """Should have correct default values."""
        role = Role(
            name="test",
            prompt="Test",
        )

        assert role.description is None
        assert role.capabilities == []
        assert role.metadata == {}


class TestRoleLoaderInit:
    """Test RoleLoader initialization."""

    def test_creates_loader(self):
        """Should initialize loader."""
        loader = RoleLoader()

        assert loader is not None
        assert hasattr(loader, 'load_role')
        assert hasattr(loader, 'get_role')

    def test_initializes_with_default_agents_dir(self):
        """Should initialize with default .claude/agents directory."""
        loader = RoleLoader()

        assert loader._agents_dir == Path(".claude/agents")

    def test_initializes_with_custom_agents_dir(self):
        """Should initialize with custom agents directory."""
        custom_dir = Path("/custom/agents")
        loader = RoleLoader(agents_dir=custom_dir)

        assert loader._agents_dir == custom_dir

    def test_initializes_empty_cache(self):
        """Should initialize with empty role cache."""
        loader = RoleLoader()

        assert loader._role_cache == {}


class TestRoleLoading:
    """Test role loading from filesystem."""

    @patch('sdp.unified.agent.role_loader.Path.exists')
    @patch('sdp.unified.agent.role_loader.Path.is_file')
    def test_loads_role_from_file(self, mock_is_file, mock_exists, tmp_path):
        """Should load role from markdown file."""
        mock_exists.return_value = True
        mock_is_file.return_value = True

        # Create actual file in tmp_path
        role_file = tmp_path / "planner.md"
        role_file.write_text("""# Planner

You are a planning agent.

**Capabilities:**
- Plan features
- Decompose tasks
""")

        # Use temporary directory for loading
        loader = RoleLoader(agents_dir=tmp_path)
        role = loader.load_role("planner")

        assert role is not None
        assert role.name == "planner"
        assert "planning agent" in role.prompt.lower()

    @patch('sdp.unified.agent.role_loader.Path.exists')
    def test_returns_none_for_missing_role(self, mock_exists):
        """Should return None for non-existent role."""
        mock_exists.return_value = False

        loader = RoleLoader()
        role = loader.load_role("missing")

        assert role is None

    def test_caches_loaded_role(self, tmp_path):
        """Should cache role after loading."""
        # Create role file
        role_file = tmp_path / "builder.md"
        role_file.write_text("# Builder\nBuild features")

        loader = RoleLoader(agents_dir=tmp_path)
        role1 = loader.load_role("builder")
        role2 = loader.load_role("builder")

        assert role1 is role2  # Same cached instance
        assert "builder" in loader._role_cache


class TestRoleRetrieval:
    """Test role retrieval from cache."""

    def test_gets_cached_role(self):
        """Should return role from cache if available."""
        loader = RoleLoader()

        # Pre-populate cache
        role = Role(name="test", prompt="Test")
        loader._role_cache["test"] = role

        retrieved = loader.get_role("test")

        assert retrieved is role

    def test_returns_none_for_uncached_role(self):
        """Should return None if role not in cache."""
        loader = RoleLoader()

        retrieved = loader.get_role("uncached")

        assert retrieved is None


class TestRoleListing:
    """Test listing available roles."""

    def test_lists_cached_roles(self):
        """Should list all cached roles."""
        loader = RoleLoader()

        # Pre-populate cache
        loader._role_cache["planner"] = Mock()
        loader._role_cache["builder"] = Mock()
        loader._role_cache["tester"] = Mock()

        roles = loader.list_roles()

        assert set(roles) == {"planner", "builder", "tester"}

    def test_returns_empty_list_when_no_roles(self):
        """Should return empty list when cache is empty."""
        loader = RoleLoader()

        roles = loader.list_roles()

        assert roles == []


class TestRoleParsing:
    """Test role file parsing."""

    def test_parses_role_name_from_heading(self, tmp_path):
        """Should parse role name from markdown heading."""
        role_file = tmp_path / "super-planner.md"
        role_file.write_text("""# Super Planner

You are a super planner.
""")

        loader = RoleLoader(agents_dir=tmp_path)
        role = loader.load_role("super-planner")

        assert role.name == "super-planner"
        assert "super planner" in role.prompt.lower()

    def test_extracts_capabilities_from_section(self, tmp_path):
        """Should extract capabilities from markdown section."""
        role_file = tmp_path / "builder.md"
        role_file.write_text("""# Builder

Build features.

**Capabilities:**
- Implement features
- Write tests
- Refactor code
""")

        loader = RoleLoader(agents_dir=tmp_path)
        role = loader.load_role("builder")

        assert len(role.capabilities) == 3
        assert "Implement features" in role.capabilities

    def test_handles_missing_capabilities(self, tmp_path):
        """Should handle role without capabilities section."""
        role_file = tmp_path / "tester.md"
        role_file.write_text("""# Tester

Test everything.
""")

        loader = RoleLoader(agents_dir=tmp_path)
        role = loader.load_role("tester")

        assert role.capabilities == []


class TestCacheInvalidation:
    """Test cache invalidation."""

    def test_clear_cache(self):
        """Should clear role cache."""
        loader = RoleLoader()

        # Pre-populate cache
        loader._role_cache["test"] = Mock()

        loader.clear_cache()

        assert loader._role_cache == {}

    def test_reload_role_updates_cache(self, tmp_path):
        """Should reload role and update cache."""
        role_file = tmp_path / "test.md"
        role_file.write_text("# Test\nOriginal")

        loader = RoleLoader(agents_dir=tmp_path)

        # First load
        role1 = loader.load_role("test")

        # Modify file and reload
        role_file.write_text("# Test\nUpdated")
        role2 = loader.load_role("test")

        # Should return cached version, not reload
        assert role1 is role2


class TestErrorHandling:
    """Test error handling."""

    def test_handles_file_read_error(self, tmp_path, monkeypatch):
        """Should handle file read errors gracefully."""
        # Create file that will raise IOError on read
        role_file = tmp_path / "error.md"

        def mock_read_text(self):
            raise IOError("Permission denied")

        # Create file
        role_file.write_text("test")
        # Monkeypatch read_text to raise error
        monkeypatch.setattr(role_file.__class__, "read_text", mock_read_text)

        loader = RoleLoader(agents_dir=tmp_path)
        role = loader.load_role("error")

        assert role is None

    @patch('sdp.unified.agent.role_loader.Path.exists')
    def test_handles_directory_not_found(self, mock_exists):
        """Should handle missing agents directory."""
        mock_exists.return_value = False

        loader = RoleLoader(agents_dir=Path("/nonexistent"))
        role = loader.load_role("test")

        assert role is None
