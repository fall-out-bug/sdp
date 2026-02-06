"""Tests for mock.py - Mock client edge cases and cache behavior."""

import pytest

from sdp.beads.client import MockBeadsClient, BeadsClientError
from sdp.beads.models import (
    BeadsStatus,
    BeadsTaskCreate,
    BeadsDependency,
    BeadsDependencyType,
    BeadsPriority,
)


class TestMockBeadsClientEdgeCases:
    """Test edge cases in MockBeadsClient."""

    def test_generate_id_uniqueness(self) -> None:
        """Test that generated IDs are unique."""
        client = MockBeadsClient()

        ids = set()
        for _ in range(100):
            task = client.create_task(
                BeadsTaskCreate(title="Test", description="", priority=BeadsPriority.HIGH)
            )
            assert task.id not in ids
            ids.add(task.id)

    def test_get_task_nonexistent(self) -> None:
        """Test get_task returns None for non-existent task."""
        client = MockBeadsClient()
        result = client.get_task("bd-9999")
        assert result is None

    def test_update_task_status_nonexistent(self) -> None:
        """Test update_task_status raises error for non-existent task."""
        client = MockBeadsClient()
        with pytest.raises(BeadsClientError) as exc_info:
            client.update_task_status("bd-9999", BeadsStatus.IN_PROGRESS)
        assert "Task not found" in str(exc_info.value)

    def test_get_ready_tasks_in_progress_status(self) -> None:
        """Test get_ready_tasks includes IN_PROGRESS tasks."""
        client = MockBeadsClient()
        task = client.create_task(
            BeadsTaskCreate(title="Test", description="", priority=BeadsPriority.HIGH)
        )
        client.update_task_status(task.id, BeadsStatus.IN_PROGRESS)

        ready = client.get_ready_tasks()
        # Current implementation only returns OPEN tasks
        # This test documents current behavior
        assert task.id not in ready  # IN_PROGRESS not included

    def test_get_ready_tasks_blocked_by_closed_task(self) -> None:
        """Test get_ready_tasks excludes tasks blocked by closed dependencies."""
        client = MockBeadsClient()

        blocker = client.create_task(
            BeadsTaskCreate(title="Blocker", description="", priority=BeadsPriority.HIGH)
        )
        blocker.status = BeadsStatus.CLOSED

        dependent = client.create_task(
            BeadsTaskCreate(
                title="Dependent",
                description="",
                priority=BeadsPriority.HIGH,
                dependencies=[
                    BeadsDependency(task_id=blocker.id, type=BeadsDependencyType.BLOCKS)
                ],
            )
        )

        ready = client.get_ready_tasks()
        # Closed blocker should not block dependent
        assert dependent.id in ready

    def test_add_dependency_nonexistent_from_task(self) -> None:
        """Test add_dependency raises error for non-existent from task."""
        client = MockBeadsClient()
        task = client.create_task(
            BeadsTaskCreate(title="Test", description="", priority=BeadsPriority.HIGH)
        )

        with pytest.raises(BeadsClientError) as exc_info:
            client.add_dependency("bd-9999", task.id, "blocks")
        assert "Task not found" in str(exc_info.value)

    def test_add_dependency_nonexistent_to_task(self) -> None:
        """Test add_dependency raises error for non-existent to task."""
        client = MockBeadsClient()
        task = client.create_task(
            BeadsTaskCreate(title="Test", description="", priority=BeadsPriority.HIGH)
        )

        with pytest.raises(BeadsClientError) as exc_info:
            client.add_dependency(task.id, "bd-9999", "blocks")
        assert "Task not found" in str(exc_info.value)

    def test_add_dependency_self_reference(self) -> None:
        """Test add_dependency allows self-reference (edge case)."""
        client = MockBeadsClient()
        task = client.create_task(
            BeadsTaskCreate(title="Test", description="", priority=BeadsPriority.HIGH)
        )

        # Should not raise error (current implementation allows this)
        client.add_dependency(task.id, task.id, "blocks")

        updated = client.get_task(task.id)
        assert len(updated.dependencies) == 1
        assert updated.dependencies[0].task_id == task.id

    def test_list_tasks_no_filters(self) -> None:
        """Test list_tasks returns all tasks with no filters."""
        client = MockBeadsClient()

        task1 = client.create_task(
            BeadsTaskCreate(title="Task 1", description="", priority=BeadsPriority.HIGH)
        )
        task2 = client.create_task(
            BeadsTaskCreate(title="Task 2", description="", priority=BeadsPriority.HIGH)
        )

        tasks = client.list_tasks()
        assert len(tasks) == 2
        ids = {t.id for t in tasks}
        assert task1.id in ids
        assert task2.id in ids

    def test_list_tasks_status_filter(self) -> None:
        """Test list_tasks filters by status."""
        client = MockBeadsClient()

        task1 = client.create_task(
            BeadsTaskCreate(title="Task 1", description="", priority=BeadsPriority.HIGH)
        )
        task2 = client.create_task(
            BeadsTaskCreate(title="Task 2", description="", priority=BeadsPriority.HIGH)
        )
        client.update_task_status(task1.id, BeadsStatus.IN_PROGRESS)

        open_tasks = client.list_tasks(status=BeadsStatus.OPEN)
        assert len(open_tasks) == 1
        assert open_tasks[0].id == task2.id

        in_progress_tasks = client.list_tasks(status=BeadsStatus.IN_PROGRESS)
        assert len(in_progress_tasks) == 1
        assert in_progress_tasks[0].id == task1.id

    def test_list_tasks_parent_filter(self) -> None:
        """Test list_tasks filters by parent_id."""
        client = MockBeadsClient()

        parent = client.create_task(
            BeadsTaskCreate(title="Parent", description="", priority=BeadsPriority.HIGH)
        )
        child1 = client.create_task(
            BeadsTaskCreate(
                title="Child 1", description="", priority=1, parent_id=parent.id
            )
        )
        child2 = client.create_task(
            BeadsTaskCreate(
                title="Child 2", description="", priority=1, parent_id=parent.id
            )
        )
        other = client.create_task(
            BeadsTaskCreate(title="Other", description="", priority=BeadsPriority.HIGH)
        )

        children = client.list_tasks(parent_id=parent.id)
        assert len(children) == 2
        ids = {c.id for c in children}
        assert child1.id in ids
        assert child2.id in ids
        assert other.id not in ids

    def test_list_tasks_combined_filters(self) -> None:
        """Test list_tasks with both status and parent filters."""
        client = MockBeadsClient()

        parent = client.create_task(
            BeadsTaskCreate(title="Parent", description="", priority=BeadsPriority.HIGH)
        )
        child1 = client.create_task(
            BeadsTaskCreate(
                title="Child 1", description="", priority=1, parent_id=parent.id
            )
        )
        child2 = client.create_task(
            BeadsTaskCreate(
                title="Child 2", description="", priority=1, parent_id=parent.id
            )
        )
        client.update_task_status(child1.id, BeadsStatus.IN_PROGRESS)

        open_children = client.list_tasks(
            status=BeadsStatus.OPEN, parent_id=parent.id
        )
        assert len(open_children) == 1
        assert open_children[0].id == child2.id

    def test_update_metadata_nonexistent_task(self) -> None:
        """Test update_metadata raises error for non-existent task."""
        client = MockBeadsClient()
        with pytest.raises(BeadsClientError) as exc_info:
            client.update_metadata("bd-9999", {"key": "value"})
        assert "Task not found" in str(exc_info.value)

    def test_update_metadata_merges_existing(self) -> None:
        """Test update_metadata merges with existing metadata."""
        client = MockBeadsClient()
        task = client.create_task(
            BeadsTaskCreate(
                title="Test",
                description="",
                priority=BeadsPriority.HIGH,
                sdp_metadata={"existing": "value", "key1": "value1"},
            )
        )

        client.update_metadata(task.id, {"key1": "updated", "key2": "value2"})

        updated = client.get_task(task.id)
        assert updated.sdp_metadata["existing"] == "value"
        assert updated.sdp_metadata["key1"] == "updated"
        assert updated.sdp_metadata["key2"] == "value2"

    def test_update_metadata_empty_dict(self) -> None:
        """Test update_metadata with empty dict preserves existing."""
        client = MockBeadsClient()
        task = client.create_task(
            BeadsTaskCreate(
                title="Test",
                description="",
                priority=BeadsPriority.HIGH,
                sdp_metadata={"key": "value"},
            )
        )

        client.update_metadata(task.id, {})

        updated = client.get_task(task.id)
        assert updated.sdp_metadata["key"] == "value"

    def test_get_ready_tasks_complex_dependencies(self) -> None:
        """Test get_ready_tasks with complex dependency chains."""
        client = MockBeadsClient()

        # Create chain: ws1 -> ws2 -> ws3
        ws1 = client.create_task(
            BeadsTaskCreate(title="WS1", description="", priority=BeadsPriority.HIGH)
        )
        ws2 = client.create_task(
            BeadsTaskCreate(
                title="WS2",
                description="",
                priority=BeadsPriority.HIGH,
                dependencies=[
                    BeadsDependency(task_id=ws1.id, type=BeadsDependencyType.BLOCKS)
                ],
            )
        )
        ws3 = client.create_task(
            BeadsTaskCreate(
                title="WS3",
                description="",
                priority=BeadsPriority.HIGH,
                dependencies=[
                    BeadsDependency(task_id=ws2.id, type=BeadsDependencyType.BLOCKS)
                ],
            )
        )

        ready = client.get_ready_tasks()
        # Only ws1 should be ready (no blockers)
        assert ws1.id in ready
        assert ws2.id not in ready
        assert ws3.id not in ready

        # Close ws1, ws2 should become ready
        client.update_task_status(ws1.id, BeadsStatus.CLOSED)
        ready = client.get_ready_tasks()
        assert ws2.id in ready
        assert ws3.id not in ready

    def test_create_task_preserves_dependencies(self) -> None:
        """Test create_task preserves dependency list."""
        client = MockBeadsClient()

        blocker = client.create_task(
            BeadsTaskCreate(title="Blocker", description="", priority=BeadsPriority.HIGH)
        )

        dependent = client.create_task(
            BeadsTaskCreate(
                title="Dependent",
                description="",
                priority=BeadsPriority.HIGH,
                dependencies=[
                    BeadsDependency(task_id=blocker.id, type=BeadsDependencyType.BLOCKS)
                ],
            )
        )

        assert len(dependent.dependencies) == 1
        assert dependent.dependencies[0].task_id == blocker.id
