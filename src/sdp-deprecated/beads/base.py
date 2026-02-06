"""Abstract base class for Beads client implementations."""

from abc import ABC, abstractmethod
from typing import List, Optional

from .models import BeadsStatus, BeadsTask, BeadsTaskCreate


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

    @abstractmethod
    def update_metadata(self, task_id: str, metadata: dict[str, object]) -> None:
        """Update task metadata.

        Args:
            task_id: Beads task ID
            metadata: Metadata dictionary to merge/update

        Raises:
            BeadsClientError: If task not found or update fails
        """
        pass
