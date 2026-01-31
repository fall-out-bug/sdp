"""Integration tests for RoleLoader module.

Tests role loading edge cases, caching behavior, and
error scenarios.
"""

import pytest
from pathlib import Path
from unittest.mock import patch

from sdp.unified.agent.role_loader import RoleLoader, Role


class TestRoleLoaderEdgeCases:
    """Test edge cases and error handling."""

    def test_load_role_with_empty_file(self, tmp_path):
        """Should handle empty role file."""
        role_file = tmp_path / "empty.md"
        role_file.write_text("")

        loader = RoleLoader(agents_dir=tmp_path)
        role = loader.load_role("empty")

        # Should still create a role
        assert role is not None
        assert role.name == "empty"

    def test_load_role_with_no_heading(self, tmp_path):
        """Should handle role file without heading."""
        role_file = tmp_path / "no-heading.md"
        role_file.write_text("Just some text without a heading.")

        loader = RoleLoader(agents_dir=tmp_path)
        role = loader.load_role("no-heading")

        assert role is not None
        assert role.description is None

    def test_load_role_with_multiline_description(self, tmp_path):
        """Should handle multiline description."""
        role_file = tmp_path / "multiline.md"
        role_file.write_text("""# Multi-line Role

This is a long description
that spans multiple lines.
""")

        loader = RoleLoader(agents_dir=tmp_path)
        role = loader.load_role("multiline")

        assert role is not None
        # Description is from heading
        assert role.description == "Multi-line Role"

    def test_cache_invalidation_on_file_change(self, tmp_path):
        """Should respect cache (no auto-reload)."""
        role_file = tmp_path / "cached.md"
        role_file.write_text("# Cached\nOriginal content")

        loader = RoleLoader(agents_dir=tmp_path)
        role1 = loader.load_role("cached")

        # Modify file
        role_file.write_text("# Cached\nUpdated content")

        # Load again - should return cached version
        role2 = loader.load_role("cached")

        assert role1 is role2  # Same cached instance

    def test_cache_after_clear(self, tmp_path):
        """Should reload after cache clear."""
        role_file = tmp_path / "cache-clear.md"
        role_file.write_text("# Test\nOriginal")

        loader = RoleLoader(agents_dir=tmp_path)
        role1 = loader.load_role("cache-clear")

        # Clear cache
        loader.clear_cache()

        # Modify file
        role_file.write_text("# Test\nUpdated")

        # Load again - should get new content
        role2 = loader.load_role("cache-clear")

        assert role1 is not role2  # Different instances
        assert "Updated" in role2.prompt

    def test_get_nonexistent_role_from_cache(self):
        """Should return None for uncached role."""
        loader = RoleLoader()

        role = loader.get_role("nonexistent")

        assert role is None


class TestRoleParsingEdgeCases:
    """Test edge cases in role parsing."""

    def test_capability_with_special_chars(self, tmp_path):
        """Should handle capabilities with special characters."""
        role_file = tmp_path / "special.md"
        role_file.write_text("""# Special

**Capabilities:**
- Feature with "quotes"
- Feature with <brackets>
- Feature with (parens)
""")

        loader = RoleLoader(agents_dir=tmp_path)
        role = loader.load_role("special")

        assert len(role.capabilities) == 3

    def test_capability_list_with_empty_lines(self, tmp_path):
        """Should handle capability list with empty lines."""
        role_file = tmp_path / "empty-caps.md"
        role_file.write_text("""# Test

**Capabilities:**
- Valid capability

- Another valid
""")

        loader = RoleLoader(agents_dir=tmp_path)
        role = loader.load_role("empty-caps")

        # Should only get non-empty lines
        assert len(role.capabilities) == 2

    def test_capabilities_section_with_different_format(self, tmp_path):
        """Should handle different capability formats."""
        role_file = tmp_path / "format.md"
        role_file.write_text("""# Test

**Capabilities:**
- Capability one
- Capability two
- Capability three
""")

        loader = RoleLoader(agents_dir=tmp_path)
        role = loader.load_role("format")

        assert len(role.capabilities) == 3

    def test_case_insensitive_capabilities_section(self, tmp_path):
        """Should find capabilities section case-insensitively."""
        role_file = tmp_path / "case.md"
        role_file.write_text("""# Test

**capabilities:**
- Lowercase section name
""")

        loader = RoleLoader(agents_dir=tmp_path)
        role = loader.load_role("case")

        # Should still find the section
        assert len(role.capabilities) == 1


class TestRoleDataclass:
    """Test Role dataclass edge cases."""

    def test_role_with_empty_metadata(self):
        """Should create role with empty metadata."""
        role = Role(
            name="test",
            prompt="Test",
            metadata={},
        )

        assert role.metadata == {}

    def test_role_with_custom_metadata(self):
        """Should create role with custom metadata."""
        role = Role(
            name="test",
            prompt="Test",
            metadata={"key": "value", "number": 123},
        )

        assert role.metadata["key"] == "value"
        assert role.metadata["number"] == 123


class TestRoleListing:
    """Test role listing edge cases."""

    def test_list_roles_after_multiple_loads(self, tmp_path):
        """Should list all cached roles."""
        # Create role files first
        for role_name in ["role1", "role2", "role3"]:
            role_file = tmp_path / f"{role_name}.md"
            role_file.write_text(f"# {role_name.title()}\nTest")

        loader = RoleLoader(agents_dir=tmp_path)

        loader.load_role("role1")
        loader.load_role("role2")
        loader.load_role("role3")

        roles = loader.list_roles()

        assert set(roles) == {"role1", "role2", "role3"}

    def test_list_roles_isolated_between_loaders(self, tmp_path):
        """Should isolate roles between different loaders."""
        # Create role files
        for role_name in ["role1", "role2"]:
            role_file = tmp_path / f"{role_name}.md"
            role_file.write_text(f"# {role_name.title()}\nTest")

        loader1 = RoleLoader(agents_dir=tmp_path)
        loader2 = RoleLoader(agents_dir=tmp_path)

        loader1.load_role("role1")
        loader2.load_role("role2")

        assert loader1.list_roles() == ["role1"]
        assert loader2.list_roles() == ["role2"]
