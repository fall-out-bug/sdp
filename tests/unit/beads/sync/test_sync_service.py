"""Unit tests for beads/sync/sync_service.py

Tests for BeadsSyncService to increase coverage from 60% to 80%+.
"""

from pathlib import Path

import pytest

from sdp.beads.mock import MockBeadsClient
from sdp.beads.models import BeadsStatus, BeadsTaskCreate
from sdp.beads.sync import BeadsSyncService


class TestBeadsSyncServiceInit:
    """Test BeadsSyncService initialization."""

    def test_init_loads_mapping(self, tmp_path: Path) -> None:
        """Initialize and load existing mapping."""
        mapping_file = tmp_path / "mapping.jsonl"
        mapping_file.write_text('{"sdp_id": "00-001-01", "beads_id": "bd-abc"}\n')
        
        client = MockBeadsClient()
        sync = BeadsSyncService(client, mapping_file=mapping_file)
        
        assert sync.mapping_manager.get_beads_id("00-001-01") == "bd-abc"

    def test_init_default_mapping_file(self, tmp_path: Path, monkeypatch) -> None:
        """Use default mapping file when none specified."""
        monkeypatch.chdir(tmp_path)
        
        client = MockBeadsClient()
        sync = BeadsSyncService(client)
        
        # Should use .beads-sdp-mapping.jsonl in current directory
        assert sync.mapping_manager.mapping_file == tmp_path / ".beads-sdp-mapping.jsonl"


class TestPersistMapping:
    """Test persist_mapping method."""

    def test_persist_saves_mappings(self, tmp_path: Path) -> None:
        """persist_mapping saves current mappings to file."""
        mapping_file = tmp_path / "mapping.jsonl"
        client = MockBeadsClient()
        sync = BeadsSyncService(client, mapping_file=mapping_file)
        
        sync.mapping_manager.add_mapping("00-001-01", "bd-abc")
        sync.persist_mapping()
        
        assert mapping_file.exists()
        content = mapping_file.read_text()
        assert "00-001-01" in content
        assert "bd-abc" in content


class TestSyncWorkstreamToBeads:
    """Test sync_workstream_to_beads method."""

    def test_missing_ws_id_returns_error(self, tmp_path: Path) -> None:
        """Return error when ws_id missing."""
        client = MockBeadsClient()
        sync = BeadsSyncService(client, mapping_file=tmp_path / "mapping.jsonl")
        
        ws_data = {"title": "Test"}  # No ws_id
        result = sync.sync_workstream_to_beads(tmp_path / "test.md", ws_data)
        
        assert result.success is False
        assert "Missing ws_id" in result.error

    def test_creates_new_task(self, tmp_path: Path) -> None:
        """Create new task when not yet synced."""
        client = MockBeadsClient()
        sync = BeadsSyncService(client, mapping_file=tmp_path / "mapping.jsonl")
        
        ws_data = {
            "ws_id": "00-001-01",
            "title": "Test WS",
            "goal": "Test goal",
            "feature": "F001",
        }
        
        result = sync.sync_workstream_to_beads(tmp_path / "00-001-01.md", ws_data)
        
        assert result.success is True
        assert result.beads_id is not None
        assert "Created new Beads task" in result.message

    def test_updates_existing_task(self, tmp_path: Path) -> None:
        """Update existing task when already synced."""
        client = MockBeadsClient()
        mapping_file = tmp_path / "mapping.jsonl"
        sync = BeadsSyncService(client, mapping_file=mapping_file)
        
        # Create task and mapping
        task = client.create_task(BeadsTaskCreate(title="Test", description="", priority=1))
        sync.mapping_manager.add_mapping("00-001-01", task.id)
        
        ws_data = {"ws_id": "00-001-01", "status": "completed"}
        result = sync.sync_workstream_to_beads(tmp_path / "00-001-01.md", ws_data)
        
        assert result.success is True
        assert result.beads_id == task.id
        assert "Updated existing Beads task" in result.message

    def test_update_task_handles_exception(self, tmp_path: Path, monkeypatch) -> None:
        """Handle exception during task update."""
        client = MockBeadsClient()
        sync = BeadsSyncService(client, mapping_file=tmp_path / "mapping.jsonl")
        
        # Create task and mapping
        task = client.create_task(BeadsTaskCreate(title="Test", description="", priority=1))
        sync.mapping_manager.add_mapping("00-001-01", task.id)
        
        # Mock update to raise exception
        def mock_update(*args, **kwargs):
            raise Exception("Network error")
        monkeypatch.setattr(client, "update_task_status", mock_update)
        
        ws_data = {"ws_id": "00-001-01", "status": "completed"}
        result = sync.sync_workstream_to_beads(tmp_path / "00-001-01.md", ws_data)
        
        assert result.success is False
        assert "Failed to update" in result.error

    def test_create_task_handles_exception(self, tmp_path: Path, monkeypatch) -> None:
        """Handle exception during task creation."""
        client = MockBeadsClient()
        sync = BeadsSyncService(client, mapping_file=tmp_path / "mapping.jsonl")
        
        # Mock create to raise exception
        def mock_create(*args, **kwargs):
            raise Exception("API error")
        monkeypatch.setattr(client, "create_task", mock_create)
        
        ws_data = {"ws_id": "00-001-01", "title": "Test", "goal": "Goal"}
        result = sync.sync_workstream_to_beads(tmp_path / "00-001-01.md", ws_data)
        
        assert result.success is False
        assert "Failed to create" in result.error


class TestBuildDescription:
    """Test _build_description helper method."""

    def test_build_description_with_goal(self, tmp_path: Path) -> None:
        """Build description with goal."""
        client = MockBeadsClient()
        sync = BeadsSyncService(client, mapping_file=tmp_path / "mapping.jsonl")
        
        ws_data = {"ws_id": "00-001-01", "title": "Test", "goal": "Test goal"}
        result = sync.sync_workstream_to_beads(tmp_path / "00-001-01.md", ws_data)
        
        task = client.get_task(result.beads_id)
        assert "Test goal" in task.description

    def test_build_description_with_context(self, tmp_path: Path) -> None:
        """Build description with context."""
        client = MockBeadsClient()
        sync = BeadsSyncService(client, mapping_file=tmp_path / "mapping.jsonl")
        
        ws_data = {
            "ws_id": "00-001-01",
            "title": "Test",
            "goal": "Goal",
            "context": "Important context",
        }
        result = sync.sync_workstream_to_beads(tmp_path / "00-001-01.md", ws_data)
        
        task = client.get_task(result.beads_id)
        assert "Important context" in task.description

    def test_build_description_with_acceptance_criteria(self, tmp_path: Path) -> None:
        """Build description with acceptance criteria."""
        client = MockBeadsClient()
        sync = BeadsSyncService(client, mapping_file=tmp_path / "mapping.jsonl")
        
        ws_data = {
            "ws_id": "00-001-01",
            "title": "Test",
            "goal": "Goal",
            "acceptance_criteria": [
                {"checked": True, "description": "Criterion 1"},
                {"checked": False, "text": "Criterion 2"},
            ],
        }
        result = sync.sync_workstream_to_beads(tmp_path / "00-001-01.md", ws_data)
        
        task = client.get_task(result.beads_id)
        assert "Criterion 1" in task.description
        assert "Criterion 2" in task.description
        assert "✓" in task.description
        assert "☐" in task.description


class TestMapDependencies:
    """Test _map_dependencies helper method."""

    def test_map_dependencies_existing(self, tmp_path: Path) -> None:
        """Map dependencies that are already synced."""
        client = MockBeadsClient()
        sync = BeadsSyncService(client, mapping_file=tmp_path / "mapping.jsonl")
        
        # Create dependency task
        dep_task = client.create_task(BeadsTaskCreate(title="Dep", description="", priority=1))
        sync.mapping_manager.add_mapping("00-001-00", dep_task.id)
        
        ws_data = {
            "ws_id": "00-001-01",
            "title": "Test",
            "goal": "Goal",
            "dependencies": ["00-001-00"],
        }
        result = sync.sync_workstream_to_beads(tmp_path / "00-001-01.md", ws_data)
        
        task = client.get_task(result.beads_id)
        assert len(task.dependencies) == 1
        assert task.dependencies[0].task_id == dep_task.id

    def test_map_dependencies_skips_unsynced(self, tmp_path: Path) -> None:
        """Skip dependencies not yet synced."""
        client = MockBeadsClient()
        sync = BeadsSyncService(client, mapping_file=tmp_path / "mapping.jsonl")
        
        ws_data = {
            "ws_id": "00-001-01",
            "title": "Test",
            "goal": "Goal",
            "dependencies": ["00-001-00"],  # Not synced
        }
        result = sync.sync_workstream_to_beads(tmp_path / "00-001-01.md", ws_data)
        
        task = client.get_task(result.beads_id)
        assert len(task.dependencies) == 0


class TestTitleTruncation:
    """Test title truncation to 500 chars."""

    def test_truncates_long_title(self, tmp_path: Path) -> None:
        """Truncate title exceeding 500 chars."""
        client = MockBeadsClient()
        sync = BeadsSyncService(client, mapping_file=tmp_path / "mapping.jsonl")
        
        long_title = "A" * 600
        ws_data = {"ws_id": "00-001-01", "title": long_title, "goal": "Goal"}
        result = sync.sync_workstream_to_beads(tmp_path / "00-001-01.md", ws_data)
        
        task = client.get_task(result.beads_id)
        assert len(task.title) == 500
        assert task.title.endswith("...")

    def test_preserves_short_title(self, tmp_path: Path) -> None:
        """Preserve title under 500 chars."""
        client = MockBeadsClient()
        sync = BeadsSyncService(client, mapping_file=tmp_path / "mapping.jsonl")
        
        ws_data = {"ws_id": "00-001-01", "title": "Short title", "goal": "Goal"}
        result = sync.sync_workstream_to_beads(tmp_path / "00-001-01.md", ws_data)
        
        task = client.get_task(result.beads_id)
        assert "00-001-01: Short title" in task.title
        assert not task.title.endswith("...")


class TestSyncBeadsToWorkstream:
    """Test sync_beads_to_workstream method."""

    def test_task_not_found(self, tmp_path: Path) -> None:
        """Return error when Beads task not found."""
        client = MockBeadsClient()
        sync = BeadsSyncService(client, mapping_file=tmp_path / "mapping.jsonl")
        
        result = sync.sync_beads_to_workstream("bd-nonexistent", tmp_path / "test.md")
        
        assert result.success is False
        assert "not found" in result.error

    def test_sync_status_from_beads(self, tmp_path: Path) -> None:
        """Sync status from Beads task."""
        client = MockBeadsClient()
        sync = BeadsSyncService(client, mapping_file=tmp_path / "mapping.jsonl")
        
        # Create task with specific status
        task = client.create_task(BeadsTaskCreate(title="Test", description="", priority=1))
        client.update_task_status(task.id, BeadsStatus.IN_PROGRESS)
        
        result = sync.sync_beads_to_workstream(task.id, tmp_path / "test.md")
        
        assert result.success is True
        assert "active" in result.message  # Maps to SDP 'active'

    def test_sync_handles_exception(self, tmp_path: Path, monkeypatch) -> None:
        """Handle exception during sync."""
        client = MockBeadsClient()
        sync = BeadsSyncService(client, mapping_file=tmp_path / "mapping.jsonl")
        
        # Mock get_task to raise exception
        def mock_get(*args, **kwargs):
            raise Exception("Connection error")
        monkeypatch.setattr(client, "get_task", mock_get)
        
        result = sync.sync_beads_to_workstream("bd-test", tmp_path / "test.md")
        
        assert result.success is False
        assert "Failed to sync" in result.error
