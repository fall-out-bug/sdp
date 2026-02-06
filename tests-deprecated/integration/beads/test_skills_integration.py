"""
BEADS-001 Skill Integration Tests.

Tests for @build, @review, @idea, @design + Beads integration.
AC5: Mock tests work without bd installed.
"""

from pathlib import Path

import pytest

from sdp.beads import (
    BeadsSyncService,
    WorkstreamExecutor,
    create_beads_client,
)
from sdp.beads.models import BeadsStatus, BeadsTaskCreate


class TestBuildBeadsIntegration:
    """AC1: Tests for @build + Beads (bd update/close, status lifecycle)."""

    def test_workstream_executor_status_lifecycle(self) -> None:
        """WorkstreamExecutor: OPEN → IN_PROGRESS → CLOSED on success."""
        client = create_beads_client(use_mock=True)
        task = client.create_task(
            BeadsTaskCreate(title="Test WS", description="", priority=1)
        )
        executor = WorkstreamExecutor(client)

        result = executor.execute(task.id, mock_tdd_success=True)

        assert result.success is True
        updated = client.get_task(task.id)
        assert updated is not None
        assert updated.status == BeadsStatus.CLOSED

    def test_workstream_executor_blocks_on_failure(self) -> None:
        """WorkstreamExecutor: OPEN → IN_PROGRESS → BLOCKED on failure."""
        client = create_beads_client(use_mock=True)
        task = client.create_task(
            BeadsTaskCreate(title="Test WS", description="", priority=1)
        )
        executor = WorkstreamExecutor(client)

        result = executor.execute(task.id, mock_tdd_success=False)

        assert result.success is False
        updated = client.get_task(task.id)
        assert updated is not None
        assert updated.status == BeadsStatus.BLOCKED


class TestReviewBeadsIntegration:
    """AC2: Tests for @review + Beads (bd list --parent, ws_id resolution)."""

    def test_list_tasks_by_parent(self) -> None:
        """Review workflow: list sub-tasks under feature via parent_id."""
        client = create_beads_client(use_mock=True)
        parent = client.create_task(
            BeadsTaskCreate(title="Feature", description="", priority=1)
        )
        child1 = client.create_task(
            BeadsTaskCreate(
                title="WS1", description="", priority=1, parent_id=parent.id
            )
        )
        child2 = client.create_task(
            BeadsTaskCreate(
                title="WS2", description="", priority=1, parent_id=parent.id
            )
        )

        children = client.list_tasks(parent_id=parent.id)

        assert len(children) == 2
        ids = {c.id for c in children}
        assert child1.id in ids
        assert child2.id in ids


class TestIdeaBeadsIntegration:
    """AC3: Tests for @idea + Beads (create_task)."""

    def test_create_task_returns_beads_id(self) -> None:
        """@idea creates Beads task, returns hash-based ID."""
        client = create_beads_client(use_mock=True)
        task = client.create_task(
            BeadsTaskCreate(
                title="Add user auth",
                description="Idea from @idea skill",
                priority=1,
            )
        )

        assert task.id.startswith("bd-")
        assert task.title == "Add user auth"
        assert task.status == BeadsStatus.OPEN


class TestDesignBeadsIntegration:
    """AC4: Tests for @design + Beads (migrate)."""

    def test_sync_workstream_to_beads_creates_task(self, tmp_path: Path) -> None:
        """@design migrate: sync_workstream_to_beads creates Beads task."""
        client = create_beads_client(use_mock=True)
        mapping_file = tmp_path / "mapping.jsonl"
        sync = BeadsSyncService(client, mapping_file=mapping_file)

        ws_data = {
            "ws_id": "00-099-01",
            "title": "Test workstream",
            "goal": "Test goal",
            "feature": "F099",
            "status": "backlog",
            "size": "SMALL",
            "acceptance_criteria": [],
        }

        result = sync.sync_workstream_to_beads(tmp_path / "00-099-01.md", ws_data)

        assert result.success is True
        assert result.beads_id is not None
        assert result.beads_id.startswith("bd-")


class TestMockWithoutBd:
    """AC5: Mock tests work without bd installed."""

    def test_mock_client_works_without_bd(self) -> None:
        """create_beads_client(use_mock=True) works when bd not installed."""
        client = create_beads_client(use_mock=True)
        task = client.create_task(
            BeadsTaskCreate(title="Test", description="", priority=1)
        )
        assert task.id is not None
        assert task.status == BeadsStatus.OPEN
