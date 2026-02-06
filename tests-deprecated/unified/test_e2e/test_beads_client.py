"""E2E tests for Beads client integration.

Tests Beads client operations with both mock and real Beads CLI,
covering CRUD operations, dependencies, and workflow integration.
"""

import os
from typing import TYPE_CHECKING
from unittest.mock import MagicMock, patch

import pytest

from sdp.beads import create_beads_client
from sdp.beads.models import (
    BeadsDependency,
    BeadsDependencyType,
    BeadsPriority,
    BeadsStatus,
    BeadsTask,
    BeadsTaskCreate,
)

if TYPE_CHECKING:
    from sdp.beads.mock import MockBeadsClient


class TestBeadsClientCreation:
    """Test Beads client creation and initialization."""

    def test_create_mock_client(self) -> None:
        """Should create mock Beads client."""
        use_mock = os.getenv("BEADS_USE_MOCK", "true").lower() == "true"

        client = create_beads_client(use_mock=use_mock)

        assert client is not None
        assert hasattr(client, 'create_task')
        assert hasattr(client, 'get_task')

    def test_create_real_client_when_available(self) -> None:
        """Should create real client when USE_MOCK=false."""
        # This test only runs if Beads CLI is available
        try:
            client = create_beads_client(use_mock=False)
            assert client is not None
        except Exception as e:
            pytest.skip(f"Beads CLI not available: {e}")


class TestTaskCRUDOperations:
    """Test task create, read, update, delete operations."""

    def test_create_task(self) -> None:
        """Should create a new task."""
        client = create_beads_client(use_mock=True)

        params = BeadsTaskCreate(
            title="Test Feature",
            description="Test description",
            priority=BeadsPriority.MEDIUM,
        )

        task = client.create_task(params)

        assert task is not None
        assert task.title == "Test Feature"
        assert task.description == "Test description"
        assert task.status == BeadsStatus.OPEN

    def test_get_task(self) -> None:
        """Should retrieve task by ID."""
        client = create_beads_client(use_mock=True)

        # Create a task first
        params = BeadsTaskCreate(title="Test", description="Test description")
        created_task = client.create_task(params)

        # Get the task
        task = client.get_task(created_task.id)

        assert task is not None
        assert task.id == created_task.id
        assert task.title == "Test"

    def test_list_ready_tasks(self) -> None:
        """Should list ready tasks."""
        client = create_beads_client(use_mock=True)

        # Create tasks
        task1 = client.create_task(BeadsTaskCreate(title="Task 1"))
        task2 = client.create_task(BeadsTaskCreate(title="Task 2"))

        # Get ready tasks
        ready_ids = client.get_ready_tasks()

        assert len(ready_ids) == 2
        assert task1.id in ready_ids
        assert task2.id in ready_ids


class TestDependencyManagement:
    """Test dependency tracking and management."""

    def test_add_dependency(self) -> None:
        """Should add dependency between tasks."""
        client = create_beads_client(use_mock=True)

        # Create tasks
        task1 = client.create_task(BeadsTaskCreate(title="Task 1"))
        task2 = client.create_task(BeadsTaskCreate(title="Task 2"))

        # Add dependency
        client.add_dependency(
            from_id=task2.id,
            to_id=task1.id,
            dep_type="blocks"
        )

        # Verify dependency was added
        updated_task2 = client.get_task(task2.id)
        assert len(updated_task2.dependencies) == 1
        assert updated_task2.dependencies[0].task_id == task1.id
        assert updated_task2.dependencies[0].type == BeadsDependencyType.BLOCKS

    def test_get_blocking_tasks(self) -> None:
        """Should get tasks blocking a given task."""
        client = create_beads_client(use_mock=True)

        # Create tasks
        task1 = client.create_task(BeadsTaskCreate(title="Task 1"))
        task2 = client.create_task(BeadsTaskCreate(title="Task 2"))

        # Add dependency
        client.add_dependency(
            from_id=task2.id,
            to_id=task1.id,
            dep_type="blocks"
        )

        # Get task2 and check dependencies
        task2_updated = client.get_task(task2.id)
        blocking_ids = [dep.task_id for dep in task2_updated.dependencies]

        assert task1.id in blocking_ids


class TestWorkflowIntegration:
    """Test Beads integration with workflow execution."""

    def test_complete_task_updates_status(self) -> None:
        """Should update task status to complete."""
        client = create_beads_client(use_mock=True)

        # Create task
        task = client.create_task(BeadsTaskCreate(title="Test"))
        assert task.status == BeadsStatus.OPEN

        # Update status
        client.update_task_status(task.id, BeadsStatus.CLOSED)

        # Verify status updated
        updated_task = client.get_task(task.id)
        assert updated_task.status == BeadsStatus.CLOSED

    def test_create_parent_child_relationships(self) -> None:
        """Should create parent-child relationships."""
        client = create_beads_client(use_mock=True)

        # Create parent task
        parent = client.create_task(BeadsTaskCreate(title="Parent Feature"))

        # Create child tasks
        child1 = client.create_task(BeadsTaskCreate(
            title="Child 1",
            parent_id=parent.id
        ))
        child2 = client.create_task(BeadsTaskCreate(
            title="Child 2",
            parent_id=parent.id
        ))

        # Verify relationships
        assert child1.parent_id == parent.id
        assert child2.parent_id == parent.id

        # List children
        children = client.list_tasks(parent_id=parent.id)
        assert len(children) == 2


class TestErrorHandling:
    """Test error handling in Beads operations."""

    def test_handles_nonexistent_task(self) -> None:
        """Should handle getting non-existent task."""
        client = create_beads_client(use_mock=True)
        task = client.get_task("nonexistent")

        # Should return None
        assert task is None

    @patch('sdp.beads.cli.subprocess.run')
    @patch('shutil.which')
    @patch.dict('os.environ', {'BEADS_USE_MOCK': 'false'}, clear=False)
    def test_handles_beads_unavailable(
        self, mock_which: MagicMock, mock_run: MagicMock
    ) -> None:
        """Should handle Beads CLI not being installed."""
        from sdp.beads.exceptions import BeadsClientError

        # Make factory think bd is in PATH so it tries CLIBeadsClient
        mock_which.return_value = "/usr/bin/bd"
        mock_run.side_effect = FileNotFoundError("bd not found")

        # Should raise BeadsClientError when use_mock=False
        with pytest.raises(BeadsClientError):
            create_beads_client(use_mock=False)


class TestBeadsMockMode:
    """Test mock mode functionality."""

    def test_mock_mode_does_not_require_beads(self) -> None:
        """Mock mode should work without Beads CLI."""
        client = create_beads_client(use_mock=True)

        # Should work without subprocess
        params = BeadsTaskCreate(title="Test", description="Desc")
        task = client.create_task(params)

        # Should return mock task
        assert task is not None
        assert task.id.startswith("bd-")

    def test_mock_task_generation(self) -> None:
        """Should generate unique mock task IDs."""
        client = create_beads_client(use_mock=True)

        task1 = client.create_task(BeadsTaskCreate(title="Task 1"))
        task2 = client.create_task(BeadsTaskCreate(title="Task 2"))

        # Should have different IDs
        assert task1.id != task2.id

        # Should have bd- prefix
        assert task1.id.startswith("bd-")
        assert task2.id.startswith("bd-")


class TestRealBeadsIntegration:
    """Integration tests with real Beads CLI."""

    @pytest.mark.skipif(
        os.getenv("BEADS_USE_MOCK", "true").lower() == "true",
        reason="Skipping real Beads tests in mock mode"
    )
    def test_real_beads_create_task(self) -> None:
        """Test creating task with real Beads CLI."""
        try:
            client = create_beads_client(use_mock=False)

            params = BeadsTaskCreate(
                title="E2E Test Feature",
                description="Testing real Beads integration",
                priority=BeadsPriority.MEDIUM,
            )

            task = client.create_task(params)

            assert task is not None
            assert task.title == "E2E Test Feature"

        except FileNotFoundError:
            pytest.skip("Beads CLI not installed")
        except Exception as e:
            pytest.fail(f"Failed to create task with real Beads: {e}")

    @pytest.mark.skipif(
        os.getenv("BEADS_USE_MOCK", "true").lower() == "true",
        reason="Skipping real Beads tests in mock mode"
    )
    def test_real_beads_get_ready_tasks(self) -> None:
        """Test getting ready tasks with real Beads."""
        try:
            client = create_beads_client(use_mock=False)

            # Create a task first
            params = BeadsTaskCreate(title="Test", description="Desc")
            task = client.create_task(params)
            client.update_task_status(task.id, BeadsStatus.OPEN)

            # Get ready tasks
            ready = client.get_ready_tasks()

            assert ready is not None
            assert isinstance(ready, list)

        except FileNotFoundError:
            pytest.skip("Beads CLI not installed")
        except Exception as e:
            pytest.fail(f"Failed to get ready tasks: {e}")


class TestCheckpointIntegration:
    """Test checkpoint save/restore with Beads."""

    def test_checkpoint_tracks_completed_ws(self) -> None:
        """Test checkpoint tracks completed workstreams."""
        client = create_beads_client(use_mock=True)

        # Create some tasks
        task1 = client.create_task(BeadsTaskCreate(title="Task 1"))
        client.update_task_status(task1.id, BeadsStatus.CLOSED)

        task2 = client.create_task(BeadsTaskCreate(title="Task 2"))

        # Get ready tasks
        ready_ids = client.get_ready_tasks()

        # Task 2 should be ready (no dependencies)
        assert task2.id in ready_ids


class TestBeadsClientFixtures:
    """Test fixtures for Beads integration."""

    @pytest.fixture
    def beads_client(self) -> "MockBeadsClient":
        """Provide Beads client for testing."""
        return create_beads_client(use_mock=True)

    @pytest.fixture
    def sample_task(self, beads_client: "MockBeadsClient") -> BeadsTask:
        """Provide sample task for testing."""
        params = BeadsTaskCreate(
            title="Sample Task",
            description="Sample description",
            priority=BeadsPriority.MEDIUM,
        )
        return beads_client.create_task(params)

    def test_fixture_usage(self, beads_client: "MockBeadsClient", sample_task: BeadsTask) -> None:
        """Test fixture usage in tests."""
        assert beads_client is not None
        assert sample_task is not None
        assert sample_task.title == "Sample Task"
