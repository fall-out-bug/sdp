"""Tests for Beads client implementations."""

from pathlib import Path
from click.testing import CliRunner
import pytest

from sdp.beads.client import (
    BeadsClientError,
    MockBeadsClient,
    create_beads_client,
)
from sdp.beads.models import (
    BeadsStatus,
    BeadsTask,
    BeadsTaskCreate,
    BeadsDependency,
)


class TestMockBeadsClient:
    """Test mock Beads client implementation."""

    def test_init(self) -> None:
        """Test mock client initialization."""
        client = MockBeadsClient()
        assert client._tasks == {}
        assert client._id_counter == 0

    def test_create_task(self) -> None:
        """Test creating a task."""
        client = MockBeadsClient()
        params = BeadsTaskCreate(
            title="Test task",
            description="Test description",
            priority=0,
        )

        task = client.create_task(params)

        assert task.id.startswith("bd-")
        assert task.title == "Test task"
        assert task.description == "Test description"
        assert task.status == BeadsStatus.OPEN
        assert task.priority == 0

    def test_get_task(self) -> None:
        """Test getting a task by ID."""
        client = MockBeadsClient()
        params = BeadsTaskCreate(
            title="Test task",
            description="Test description",
            priority=0,
        )

        created = client.create_task(params)
        retrieved = client.get_task(created.id)

        assert retrieved is not None
        assert retrieved.id == created.id
        assert retrieved.title == "Test task"

    def test_get_task_not_found(self) -> None:
        """Test getting non-existent task."""
        client = MockBeadsClient()
        task = client.get_task("bd-9999")
        assert task is None

    def test_update_task_status(self) -> None:
        """Test updating task status."""
        client = MockBeadsClient()
        params = BeadsTaskCreate(
            title="Test task",
            description="Test description",
            priority=0,
        )

        task = client.create_task(params)
        client.update_task_status(task.id, BeadsStatus.IN_PROGRESS)

        updated = client.get_task(task.id)
        assert updated.status == BeadsStatus.IN_PROGRESS

    def test_update_task_status_not_found(self) -> None:
        """Test updating status of non-existent task."""
        client = MockBeadsClient()
        with pytest.raises(BeadsClientError):
            client.update_task_status("bd-9999", BeadsStatus.IN_PROGRESS)

    def test_get_ready_tasks_no_dependencies(self) -> None:
        """Test getting ready tasks with no dependencies."""
        client = MockBeadsClient()
        params = BeadsTaskCreate(
            title="Test task",
            description="Test description",
            priority=0,
        )

        task = client.create_task(params)
        ready = client.get_ready_tasks()

        assert task.id in ready

    def test_get_ready_tasks_with_blocking_dependency(self) -> None:
        """Test getting ready tasks with blocking dependencies."""
        client = MockBeadsClient()

        # Create blocking task
        blocker = BeadsTaskCreate(
            title="Blocker",
            description="Blocking task",
            priority=0,
        )
        blocker_task = client.create_task(blocker)

        # Create dependent task
        from sdp.beads.models import BeadsDependencyType
        dependent = BeadsTaskCreate(
            title="Dependent",
            description="Dependent task",
            priority=0,
            dependencies=[BeadsDependency(task_id=blocker_task.id, type=BeadsDependencyType.BLOCKS)],
        )
        dependent_task = client.create_task(dependent)

        # Only blocker should be ready
        ready = client.get_ready_tasks()
        assert blocker_task.id in ready
        assert dependent_task.id not in ready

    def test_add_dependency(self) -> None:
        """Test adding a dependency."""
        client = MockBeadsClient()

        task1 = client.create_task(BeadsTaskCreate(
            title="Task 1",
            description="First task",
            priority=0,
        ))

        task2 = client.create_task(BeadsTaskCreate(
            title="Task 2",
            description="Second task",
            priority=0,
        ))

        client.add_dependency(task2.id, task1.id, "blocks")

        updated = client.get_task(task2.id)
        assert len(updated.dependencies) == 1
        assert updated.dependencies[0].task_id == task1.id
        # The type field is a string, not an enum
        assert updated.dependencies[0].type == "blocks"

    def test_add_dependency_duplicate(self) -> None:
        """Test adding duplicate dependency."""
        client = MockBeadsClient()

        task1 = client.create_task(BeadsTaskCreate(
            title="Task 1",
            description="First task",
            priority=0,
        ))

        task2 = client.create_task(BeadsTaskCreate(
            title="Task 2",
            description="Second task",
            priority=0,
        ))

        client.add_dependency(task2.id, task1.id, "blocks")
        client.add_dependency(task2.id, task1.id, "blocks")  # Duplicate

        updated = client.get_task(task2.id)
        assert len(updated.dependencies) == 1  # Should not duplicate

    def test_add_dependency_task_not_found(self) -> None:
        """Test adding dependency with non-existent task."""
        client = MockBeadsClient()

        task = client.create_task(BeadsTaskCreate(
            title="Task",
            description="Test task",
            priority=0,
        ))

        with pytest.raises(BeadsClientError):
            client.add_dependency(task.id, "bd-9999", "blocks")

    def test_list_tasks_no_filter(self) -> None:
        """Test listing all tasks."""
        client = MockBeadsClient()

        client.create_task(BeadsTaskCreate(
            title="Task 1",
            description="First task",
            priority=0,
        ))

        client.create_task(BeadsTaskCreate(
            title="Task 2",
            description="Second task",
            priority=0,
        ))

        tasks = client.list_tasks()
        assert len(tasks) == 2

    def test_list_tasks_with_status_filter(self) -> None:
        """Test listing tasks with status filter."""
        client = MockBeadsClient()

        task1 = client.create_task(BeadsTaskCreate(
            title="Task 1",
            description="First task",
            priority=0,
        ))

        task2 = client.create_task(BeadsTaskCreate(
            title="Task 2",
            description="Second task",
            priority=0,
        ))

        client.update_task_status(task1.id, BeadsStatus.IN_PROGRESS)

        open_tasks = client.list_tasks(status=BeadsStatus.OPEN)
        assert len(open_tasks) == 1
        assert open_tasks[0].id == task2.id

    def test_list_tasks_with_parent_filter(self) -> None:
        """Test listing tasks with parent filter."""
        client = MockBeadsClient()

        parent = client.create_task(BeadsTaskCreate(
            title="Parent",
            description="Parent task",
            priority=0,
        ))

        child = client.create_task(BeadsTaskCreate(
            title="Child",
            description="Child task",
            priority=0,
            parent_id=parent.id,
        ))

        client.create_task(BeadsTaskCreate(
            title="Other",
            description="Other task",
            priority=0,
        ))

        children = client.list_tasks(parent_id=parent.id)
        assert len(children) == 1
        assert children[0].id == child.id


class TestCreateBeadsClient:
    """Test Beads client factory function."""

    def test_create_mock_client(self) -> None:
        """Test creating mock client."""
        client = create_beads_client(use_mock=True)
        assert isinstance(client, MockBeadsClient)

    def test_create_cli_client_not_available(self) -> None:
        """Test creating CLI client when Beads not available."""
        # This test is environment-dependent
        # Skip if Beads CLI happens to be installed
        import subprocess
        try:
            subprocess.run(
                ["bd", "--version"],
                capture_output=True,
                check=True,
            )
            # Beads is installed, skip this test
            pytest.skip("Beads CLI is installed")
        except (subprocess.CalledProcessError, FileNotFoundError):
            # Beads not installed, should raise error
            with pytest.raises(BeadsClientError):
                create_beads_client(use_mock=False)
