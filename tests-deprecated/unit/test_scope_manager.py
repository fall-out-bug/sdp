"""Tests for Beads scope manager."""

import pytest

from sdp.beads.mock import MockBeadsClient
from sdp.beads.models import BeadsTaskCreate
from sdp.beads.scope_manager import ScopeManager


class TestScopeManager:
    """Test suite for ScopeManager."""

    def test_get_scope_empty_by_default(self) -> None:
        """AC1: New tasks have empty scope (unrestricted)."""
        # Arrange
        client = MockBeadsClient()
        task = client.create_task(BeadsTaskCreate(title="Test Task"))
        manager = ScopeManager(client)

        # Act
        scope = manager.get_scope(task.id)

        # Assert
        assert scope == []

    def test_get_scope_raises_on_nonexistent_task(self) -> None:
        """get_scope raises ValueError for nonexistent task."""
        # Arrange
        client = MockBeadsClient()
        manager = ScopeManager(client)

        # Act & Assert
        with pytest.raises(ValueError) as exc_info:
            manager.get_scope("nonexistent")

        assert "not found" in str(exc_info.value).lower()

    def test_set_scope(self) -> None:
        """AC2: set_scope updates task metadata."""
        # Arrange
        client = MockBeadsClient()
        task = client.create_task(BeadsTaskCreate(title="Test Task"))
        manager = ScopeManager(client)

        files = ["src/file1.py", "src/file2.py"]

        # Act
        manager.set_scope(task.id, files)

        # Assert
        updated_task = client.get_task(task.id)
        assert updated_task is not None
        assert updated_task.sdp_metadata.get("scope_files") == files

    def test_add_file(self) -> None:
        """AC3: add_file adds to scope."""
        # Arrange
        client = MockBeadsClient()
        task = client.create_task(BeadsTaskCreate(title="Test Task"))
        manager = ScopeManager(client)

        # Act
        manager.add_file(task.id, "src/file1.py")
        manager.add_file(task.id, "src/file2.py")

        # Assert
        scope = manager.get_scope(task.id)
        assert "src/file1.py" in scope
        assert "src/file2.py" in scope
        assert len(scope) == 2

    def test_add_file_idempotent(self) -> None:
        """Adding same file twice doesn't duplicate."""
        # Arrange
        client = MockBeadsClient()
        task = client.create_task(BeadsTaskCreate(title="Test Task"))
        manager = ScopeManager(client)

        # Act
        manager.add_file(task.id, "src/file1.py")
        manager.add_file(task.id, "src/file1.py")  # Duplicate

        # Assert
        scope = manager.get_scope(task.id)
        assert scope.count("src/file1.py") == 1

    def test_remove_file(self) -> None:
        """AC4: remove_file removes from scope."""
        # Arrange
        client = MockBeadsClient()
        task = client.create_task(BeadsTaskCreate(title="Test Task"))
        manager = ScopeManager(client)

        manager.set_scope(task.id, ["src/file1.py", "src/file2.py"])

        # Act
        manager.remove_file(task.id, "src/file1.py")

        # Assert
        scope = manager.get_scope(task.id)
        assert "src/file1.py" not in scope
        assert "src/file2.py" in scope

    def test_remove_file_nonexistent_safe(self) -> None:
        """Removing nonexistent file is safe (no-op)."""
        # Arrange
        client = MockBeadsClient()
        task = client.create_task(BeadsTaskCreate(title="Test Task"))
        manager = ScopeManager(client)

        manager.set_scope(task.id, ["src/file1.py"])

        # Act
        manager.remove_file(task.id, "nonexistent.py")  # Should not raise

        # Assert
        scope = manager.get_scope(task.id)
        assert scope == ["src/file1.py"]

    def test_is_in_scope_with_empty_scope(self) -> None:
        """AC5: Empty scope means all files allowed."""
        # Arrange
        client = MockBeadsClient()
        task = client.create_task(BeadsTaskCreate(title="Test Task"))
        manager = ScopeManager(client)

        # Act & Assert - empty scope = unrestricted
        assert manager.is_in_scope(task.id, "any/file.py") is True
        assert manager.is_in_scope(task.id, "another/file.py") is True

    def test_is_in_scope_with_restricted_scope(self) -> None:
        """AC5: Non-empty scope restricts to listed files."""
        # Arrange
        client = MockBeadsClient()
        task = client.create_task(BeadsTaskCreate(title="Test Task"))
        manager = ScopeManager(client)

        manager.set_scope(task.id, ["src/file1.py", "src/file2.py"])

        # Act & Assert
        assert manager.is_in_scope(task.id, "src/file1.py") is True
        assert manager.is_in_scope(task.id, "src/file2.py") is True
        assert manager.is_in_scope(task.id, "other/file.py") is False

    def test_clear_scope(self) -> None:
        """clear_scope makes scope unrestricted."""
        # Arrange
        client = MockBeadsClient()
        task = client.create_task(BeadsTaskCreate(title="Test Task"))
        manager = ScopeManager(client)

        manager.set_scope(task.id, ["src/file1.py"])

        # Act
        manager.clear_scope(task.id)

        # Assert
        scope = manager.get_scope(task.id)
        assert scope == []
        assert manager.is_in_scope(task.id, "any/file.py") is True
