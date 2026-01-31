"""Tests for skills_build.py - Workstream execution with TDD cycle."""

import pytest

from sdp.beads.client import MockBeadsClient, BeadsClientError
from sdp.beads.skills_build import WorkstreamExecutor, ExecutionResult
from sdp.beads.models import BeadsStatus, BeadsTaskCreate


class TestExecutionResult:
    """Test ExecutionResult dataclass."""

    def test_create_success_result(self) -> None:
        """Test creating successful execution result."""
        result = ExecutionResult(
            success=True,
            task_id="bd-0001",
            newly_ready=["bd-0002", "bd-0003"],
        )
        assert result.success is True
        assert result.task_id == "bd-0001"
        assert result.error is None
        assert result.newly_ready == ["bd-0002", "bd-0003"]

    def test_create_failure_result(self) -> None:
        """Test creating failed execution result."""
        result = ExecutionResult(
            success=False,
            task_id="bd-0001",
            error="TDD cycle failed",
        )
        assert result.success is False
        assert result.task_id == "bd-0001"
        assert result.error == "TDD cycle failed"
        assert result.newly_ready == []


class TestWorkstreamExecutor:
    """Test WorkstreamExecutor class."""

    def test_init(self) -> None:
        """Test executor initialization."""
        client = MockBeadsClient()
        executor = WorkstreamExecutor(client)
        assert executor.client is client

    def test_execute_success_flow(self) -> None:
        """Test successful execution flow: OPEN → IN_PROGRESS → CLOSED."""
        client = MockBeadsClient()
        executor = WorkstreamExecutor(client)

        from sdp.beads.models import BeadsPriority
        task = client.create_task(
            BeadsTaskCreate(title="Test WS", description="", priority=BeadsPriority.HIGH)
        )
        assert task.status == BeadsStatus.OPEN

        result = executor.execute(task.id, mock_tdd_success=True)

        assert result.success is True
        assert result.task_id == task.id
        assert result.error is None

        updated = client.get_task(task.id)
        assert updated is not None
        assert updated.status == BeadsStatus.CLOSED

    def test_execute_failure_flow(self) -> None:
        """Test failed execution flow: OPEN → IN_PROGRESS → BLOCKED."""
        client = MockBeadsClient()
        executor = WorkstreamExecutor(client)

        from sdp.beads.models import BeadsPriority
        task = client.create_task(
            BeadsTaskCreate(title="Test WS", description="", priority=BeadsPriority.HIGH)
        )
        assert task.status == BeadsStatus.OPEN

        result = executor.execute(task.id, mock_tdd_success=False)

        assert result.success is False
        assert result.task_id == task.id
        assert result.error == "TDD cycle failed"

        updated = client.get_task(task.id)
        assert updated is not None
        assert updated.status == BeadsStatus.BLOCKED

    def test_execute_exception_handling(self) -> None:
        """Test exception handling during execution."""
        client = MockBeadsClient()
        executor = WorkstreamExecutor(client)

        from sdp.beads.models import BeadsPriority
        task = client.create_task(
            BeadsTaskCreate(title="Test WS", description="", priority=BeadsPriority.HIGH)
        )

        # Mock execute_tdd_cycle to raise exception
        original_tdd = executor.execute_tdd_cycle

        def failing_tdd(task_id: str, mock_tdd_success: bool) -> bool:
            raise BeadsClientError("Network error")

        executor.execute_tdd_cycle = failing_tdd  # type: ignore[assignment]

        result = executor.execute(task.id, mock_tdd_success=True)

        assert result.success is False
        assert result.task_id == task.id
        assert result.error is not None
        assert "Network error" in result.error

        # Should be marked as BLOCKED due to exception
        updated = client.get_task(task.id)
        assert updated is not None
        assert updated.status == BeadsStatus.BLOCKED

    def test_execute_tdd_cycle_success(self) -> None:
        """Test TDD cycle returns True on success."""
        client = MockBeadsClient()
        executor = WorkstreamExecutor(client)

        result = executor.execute_tdd_cycle("bd-0001", mock_tdd_success=True)
        assert result is True

    def test_execute_tdd_cycle_failure(self) -> None:
        """Test TDD cycle returns False on failure."""
        client = MockBeadsClient()
        executor = WorkstreamExecutor(client)

        result = executor.execute_tdd_cycle("bd-0001", mock_tdd_success=False)
        assert result is False

    def test_run_tdd_phase(self) -> None:
        """Test _run_tdd_phase helper method."""
        client = MockBeadsClient()
        executor = WorkstreamExecutor(client)

        called = []

        def phase_fn(task_id: str) -> None:
            called.append(task_id)

        executor._run_tdd_phase("bd-0001", phase_fn)

        assert len(called) == 1
        assert called[0] == "bd-0001"

    def test_execute_newly_ready_tasks(self) -> None:
        """Test that execute returns newly ready tasks."""
        client = MockBeadsClient()
        executor = WorkstreamExecutor(client)

        # Create blocker task
        from sdp.beads.models import BeadsPriority, BeadsDependency, BeadsDependencyType
        blocker = client.create_task(
            BeadsTaskCreate(title="Blocker", description="", priority=BeadsPriority.HIGH)
        )

        # Create dependent task
        dependent = client.create_task(
            BeadsTaskCreate(
                title="Dependent",
                description="",
                priority=BeadsPriority.HIGH,
                dependencies=[
                    BeadsDependency(
                        task_id=blocker.id, type=BeadsDependencyType.BLOCKS
                    )
                ],
            )
        )

        # Initially, only blocker is ready
        ready_before = set(client.get_ready_tasks())
        assert blocker.id in ready_before
        assert dependent.id not in ready_before

        # Execute blocker
        result = executor.execute(blocker.id, mock_tdd_success=True)

        assert result.success is True
        # After blocker completes, dependent should be ready
        ready_after = set(client.get_ready_tasks())
        assert dependent.id in ready_after

    def test_execute_status_transitions(self) -> None:
        """Test status transitions during execution."""
        client = MockBeadsClient()
        executor = WorkstreamExecutor(client)

        from sdp.beads.models import BeadsPriority
        task = client.create_task(
            BeadsTaskCreate(title="Test WS", description="", priority=BeadsPriority.HIGH)
        )

        # Track status changes
        statuses = []

        original_update = client.update_task_status

        def track_update(task_id: str, status: BeadsStatus) -> None:
            statuses.append(status)
            original_update(task_id, status)

        client.update_task_status = track_update  # type: ignore[assignment]

        executor.execute(task.id, mock_tdd_success=True)

        # Should transition: OPEN → IN_PROGRESS → CLOSED
        assert BeadsStatus.IN_PROGRESS in statuses
        assert BeadsStatus.CLOSED in statuses

    def test_execute_exception_during_tdd(self) -> None:
        """Test exception during TDD cycle execution."""
        client = MockBeadsClient()
        executor = WorkstreamExecutor(client)

        from sdp.beads.models import BeadsPriority
        task = client.create_task(
            BeadsTaskCreate(title="Test WS", description="", priority=BeadsPriority.HIGH)
        )

        # Mock execute_tdd_cycle to raise exception
        original_tdd = executor.execute_tdd_cycle

        def failing_tdd(task_id: str, mock_tdd_success: bool) -> bool:
            raise RuntimeError("Test execution failed")

        executor.execute_tdd_cycle = failing_tdd  # type: ignore[assignment]

        result = executor.execute(task.id, mock_tdd_success=True)

        assert result.success is False
        assert result.error is not None
        assert "Test execution failed" in result.error

        # Should be marked as BLOCKED
        updated = client.get_task(task.id)
        assert updated is not None
        assert updated.status == BeadsStatus.BLOCKED
