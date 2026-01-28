"""
Beads client interface and implementations.

Provides both mock and real implementations for interacting with Beads.
"""

import json
import subprocess
from abc import ABC, abstractmethod
from pathlib import Path
from typing import List, Optional

from .models import (
    BeadsTask,
    BeadsTaskCreate,
    BeadsStatus,
    BeadsDependency,
    BeadsSyncResult,
)


class BeadsClientError(Exception):
    """Base exception for Beads client errors."""

    pass


class BeadsClient(ABC):
    """Abstract client for interacting with Beads.

    Implementations can be:
    - MockBeadsClient: In-memory mock for testing/development
    - CLIBeadsClient: Real Beads via subprocess CLI calls
    - APIBeadsClient: Direct API calls (future)
    """

    @abstractmethod
    def create_task(self, params: BeadsTaskCreate) -> BeadsTask:
        """Create a new task.

        Args:
            params: Task creation parameters

        Returns:
            Created task with generated ID

        Raises:
            BeadsClientError: If creation fails
        """
        pass

    @abstractmethod
    def get_task(self, task_id: str) -> Optional[BeadsTask]:
        """Get a task by ID.

        Args:
            task_id: Beads task ID (e.g., "bd-a3f8")

        Returns:
            Task if found, None otherwise
        """
        pass

    @abstractmethod
    def update_task_status(self, task_id: str, status: BeadsStatus) -> None:
        """Update task status.

        Args:
            task_id: Beads task ID
            status: New status

        Raises:
            BeadsClientError: If task not found or update fails
        """
        pass

    @abstractmethod
    def get_ready_tasks(self) -> List[str]:
        """Get IDs of tasks ready to work on.

        Ready tasks are:
        - Status is OPEN or IN_PROGRESS
        - No blocking dependencies are open

        Returns:
            List of task IDs
        """
        pass

    @abstractmethod
    def add_dependency(
        self, from_id: str, to_id: str, dep_type: str = "blocks"
    ) -> None:
        """Add a dependency relationship.

        Args:
            from_id: Task that depends on to_id
            to_id: Task being depended on
            dep_type: Dependency type (blocks, parent-child, etc.)

        Raises:
            BeadsClientError: If tasks not found or dependency invalid
        """
        pass

    @abstractmethod
    def list_tasks(
        self,
        status: Optional[BeadsStatus] = None,
        parent_id: Optional[str] = None,
    ) -> List[BeadsTask]:
        """List tasks with optional filters.

        Args:
            status: Filter by status
            parent_id: Filter by parent ID (for sub-tasks)

        Returns:
            List of matching tasks
        """
        pass


class MockBeadsClient(BeadsClient):
    """In-memory mock implementation for testing/development.

    Simulates Beads behavior without requiring actual Beads installation.
    Useful for:
    - Unit tests
    - Development without Go/Beads
    - CI/CD pipelines
    """

    def __init__(self):
        """Initialize mock client with empty task store."""
        self._tasks: dict[str, BeadsTask] = {}
        self._id_counter = 0

    def _generate_id(self) -> str:
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
                if dep.type.value == "blocks":
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
            BeadsDependency(task_id=to_id, type=dep_type)
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


class CLIBeadsClient(BeadsClient):
    """Real Beads client using CLI subprocess calls.

    Requires:
    - Go 1.24+ installed
    - Beads installed: `go install github.com/steveyegge/beads/cmd/bd@latest`
    - Beads initialized: `bd init` in project directory
    """

    def __init__(self, project_dir: Optional[Path] = None):
        """Initialize CLI client.

        Args:
            project_dir: Project directory (defaults to current dir)
        """
        self.project_dir = project_dir or Path.cwd()

        # Verify Beads is available
        try:
            result = subprocess.run(
                ["bd", "--version"],
                capture_output=True,
                cwd=self.project_dir,
                check=True,
            )
        except (subprocess.CalledProcessError, FileNotFoundError) as e:
            raise BeadsClientError(
                "Beads CLI not found. Install with: "
                "go install github.com/steveyegge/beads/cmd/bd@latest"
            ) from e

    def create_task(self, params: BeadsTaskCreate) -> BeadsTask:
        """Create task via Beads CLI.

        Example:
            bd create --json '{"title": "...", "priority": 0}'
        """
        cmd = [
            "bd",
            "create",
            "--json",
            json.dumps(params.to_dict()),
        ]

        try:
            result = subprocess.run(
                cmd,
                capture_output=True,
                text=True,
                cwd=self.project_dir,
                check=True,
            )

            return BeadsTask.from_dict(json.loads(result.stdout))

        except subprocess.CalledProcessError as e:
            raise BeadsClientError(f"Failed to create task: {e.stderr}") from e
        except json.JSONDecodeError as e:
            raise BeadsClientError(f"Invalid JSON response from Beads: {e}") from e

    def get_task(self, task_id: str) -> Optional[BeadsTask]:
        """Get task via Beads CLI.

        Example:
            bd show --json bd-a3f8
        """
        cmd = ["bd", "show", "--json", task_id]

        try:
            result = subprocess.run(
                cmd,
                capture_output=True,
                text=True,
                cwd=self.project_dir,
                check=True,
            )

            return BeadsTask.from_dict(json.loads(result.stdout))

        except subprocess.CalledProcessError:
            # Task not found
            return None
        except json.JSONDecodeError as e:
            raise BeadsClientError(f"Invalid JSON response from Beads: {e}") from e

    def update_task_status(self, task_id: str, status: BeadsStatus) -> None:
        """Update task status via Beads CLI.

        Example:
            bd update bd-a3f8 --status in_progress
        """
        cmd = ["bd", "update", task_id, "--status", status.value]

        try:
            subprocess.run(
                cmd,
                capture_output=True,
                cwd=self.project_dir,
                check=True,
            )

        except subprocess.CalledProcessError as e:
            raise BeadsClientError(f"Failed to update status: {e.stderr}") from e

    def get_ready_tasks(self) -> List[str]:
        """Get ready tasks via Beads CLI.

        Example:
            bd ready --json
        """
        cmd = ["bd", "ready", "--json"]

        try:
            result = subprocess.run(
                cmd,
                capture_output=True,
                text=True,
                cwd=self.project_dir,
                check=True,
            )

            data = json.loads(result.stdout)
            return data.get("ready_tasks", [])

        except subprocess.CalledProcessError as e:
            raise BeadsClientError(f"Failed to get ready tasks: {e.stderr}") from e
        except json.JSONDecodeError as e:
            raise BeadsClientError(f"Invalid JSON response from Beads: {e}") from e

    def add_dependency(
        self, from_id: str, to_id: str, dep_type: str = "blocks"
    ) -> None:
        """Add dependency via Beads CLI.

        Example:
            bd dep add bd-a3f8.1 bd-a3f8 --type blocks
        """
        cmd = ["bd", "dep", "add", from_id, to_id, "--type", dep_type]

        try:
            subprocess.run(
                cmd,
                capture_output=True,
                cwd=self.project_dir,
                check=True,
            )

        except subprocess.CalledProcessError as e:
            raise BeadsClientError(f"Failed to add dependency: {e.stderr}") from e

    def list_tasks(
        self,
        status: Optional[BeadsStatus] = None,
        parent_id: Optional[str] = None,
    ) -> List[BeadsTask]:
        """List tasks via Beads CLI.

        Example:
            bd list --status open --json
        """
        cmd = ["bd", "list", "--json"]

        if status:
            cmd.extend(["--status", status.value])

        if parent_id:
            cmd.extend(["--parent", parent_id])

        try:
            result = subprocess.run(
                cmd,
                capture_output=True,
                text=True,
                cwd=self.project_dir,
                check=True,
            )

            data = json.loads(result.stdout)
            return [BeadsTask.from_dict(t) for t in data.get("tasks", [])]

        except subprocess.CalledProcessError as e:
            raise BeadsClientError(f"Failed to list tasks: {e.stderr}") from e
        except json.JSONDecodeError as e:
            raise BeadsClientError(f"Invalid JSON response from Beads: {e}") from e


def create_beads_client(
    use_mock: bool = False, project_dir: Optional[Path] = None
) -> BeadsClient:
    """Factory function to create appropriate Beads client.

    Args:
        use_mock: Force mock client (for testing)
        project_dir: Project directory (for CLI client)

    Returns:
        BeadsClient instance

    Example:
        # Use mock for testing
        client = create_beads_client(use_mock=True)

        # Use real Beads (must be installed)
        client = create_beads_client()
    """
    if use_mock:
        return MockBeadsClient()
    else:
        return CLIBeadsClient(project_dir)
