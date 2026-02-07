"""Mock Beads client implementation for testing/development."""

from typing import List, Optional

from .base import BeadsClient
from .exceptions import BeadsClientError
from .models import BeadsDependency, BeadsDependencyType, BeadsStatus, BeadsTask, BeadsTaskCreate


class MockBeadsClient(BeadsClient):
    """In-memory mock implementation for testing/development.

    Simulates Beads behavior without requiring actual Beads installation.
    Useful for:
    - Unit tests
    - Development without Go/Beads
    - CI/CD pipelines
    """

    def __init__(self) -> None:
        """Initialize mock client with empty task store."""
        self._tasks: dict[str, BeadsTask] = {}
        self._id_counter = 0

    def _generate_id(self) -> str:  # noqa: ANN202
        """Generate a mock Beads-style ID.

        In real Beads, this would be a content-addressed hash.
        For mocking, we use a simple counter.
        """
        self._id_counter += 1
        # Simulate hash format: bd-XXXX
        return f"bd-{self._id_counter:04x}"

    def create_task(self, params: BeadsTaskCreate) -> BeadsTask:
        """Create a new task (mock)."""
        task_id = self._generate_id()

        task = BeadsTask(
            id=task_id,
            title=params.title,
            description=params.description,
            status=BeadsStatus.OPEN,
            priority=params.priority,
            parent_id=params.parent_id,
            dependencies=params.dependencies.copy(),
            external_ref=params.external_ref,
            sdp_metadata=params.sdp_metadata.copy(),
        )

        self._tasks[task_id] = task
        return task

    def get_task(self, task_id: str) -> Optional[BeadsTask]:
        """Get a task by ID (mock)."""
        return self._tasks.get(task_id)

    def update_task_status(self, task_id: str, status: BeadsStatus) -> None:
        """Update task status (mock)."""
        task = self._tasks.get(task_id)
        if not task:
            raise BeadsClientError(f"Task not found: {task_id}")

        task.status = status

    def get_ready_tasks(self) -> List[str]:
        """Get ready task IDs (mock).

        Simple implementation: returns OPEN tasks with no blockers.
        """
        ready = []

        for task_id, task in self._tasks.items():
            if task.status != BeadsStatus.OPEN:
                continue

            # Check if any blocking dependencies are still open
            blocked = False
            for dep in task.dependencies:
                # dep.type is a string, not an enum
                if dep.type == "blocks":
                    dep_task = self._tasks.get(dep.task_id)
                    if dep_task and dep_task.status == BeadsStatus.OPEN:
                        blocked = True
                        break

            if not blocked:
                ready.append(task_id)

        return ready

    def add_dependency(
        self, from_id: str, to_id: str, dep_type: str = "blocks"
    ) -> None:
        """Add dependency (mock)."""
        from_task = self._tasks.get(from_id)
        to_task = self._tasks.get(to_id)

        if not from_task:
            raise BeadsClientError(f"Task not found: {from_id}")
        if not to_task:
            raise BeadsClientError(f"Task not found: {to_id}")

        # Check if dependency already exists
        for dep in from_task.dependencies:
            if dep.task_id == to_id:
                return  # Already exists

        from_task.dependencies.append(
            BeadsDependency(task_id=to_id, type=BeadsDependencyType(dep_type))
        )

    def list_tasks(
        self,
        status: Optional[BeadsStatus] = None,
        parent_id: Optional[str] = None,
    ) -> List[BeadsTask]:
        """List tasks with filters (mock)."""
        tasks = list(self._tasks.values())

        if status:
            tasks = [t for t in tasks if t.status == status]

        if parent_id:
            tasks = [t for t in tasks if t.parent_id == parent_id]

        return tasks

    def update_metadata(self, task_id: str, metadata: dict[str, object]) -> None:
        """Update task metadata (mock)."""
        task = self._tasks.get(task_id)
        if not task:
            raise BeadsClientError(f"Task not found: {task_id}")

        # Merge metadata
        task.sdp_metadata.update(metadata)
