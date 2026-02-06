"""Targeted tests for uncovered lines in beads/models.py."""

from datetime import datetime

import pytest

from sdp.beads.models import (
    BeadsDependency,
    BeadsDependencyType,
    BeadsPriority,
    BeadsStatus,
    BeadsTask,
    BeadsTaskCreate,
)


def test_beads_task_to_dict_with_dates() -> None:
    """Test BeadsTask.to_dict with dates (lines 96-97)."""
    now = datetime(2024, 1, 1, 12, 0, 0)
    task = BeadsTask(
        id="bd-1234",
        title="Test",
        created_at=now,
        updated_at=now,
    )
    data = task.to_dict()
    assert data["created_at"] == "2024-01-01T12:00:00"
    assert data["updated_at"] == "2024-01-01T12:00:00"


def test_beads_task_to_dict_with_sdp_metadata() -> None:
    """Test BeadsTask.to_dict with sdp_metadata (line 98)."""
    task = BeadsTask(
        id="bd-1234",
        title="Test",
        sdp_metadata={"ws_id": "00-001-01"},
    )
    data = task.to_dict()
    assert data["metadata"] == {"sdp": {"ws_id": "00-001-01"}}


def test_beads_task_to_dict_without_sdp_metadata() -> None:
    """Test BeadsTask.to_dict without sdp_metadata (line 98)."""
    task = BeadsTask(
        id="bd-1234",
        title="Test",
        sdp_metadata={},
    )
    data = task.to_dict()
    assert data["metadata"] is None


def test_beads_task_from_dict_invalid_dependency() -> None:
    """Test BeadsTask.from_dict skips invalid dependencies (lines 106-107)."""
    data = {
        "id": "bd-1234",
        "title": "Test",
        "dependencies": [
            {"task_id": "bd-5678", "type": "blocks"},  # Valid
            "invalid_string",  # Invalid - not a dict
            {"no_task_id": "foo"},  # Invalid - missing task_id
        ],
    }
    task = BeadsTask.from_dict(data)
    # Should only include valid dependency
    assert len(task.dependencies) == 1
    assert task.dependencies[0].task_id == "bd-5678"


def test_beads_task_from_dict_dependency_with_id_field() -> None:
    """Test BeadsTask.from_dict handles 'id' field for backward compat (line 108)."""
    data = {
        "id": "bd-1234",
        "title": "Test",
        "dependencies": [
            {"id": "bd-5678", "type": "blocks"},  # Old format: 'id' instead of 'task_id'
        ],
    }
    task = BeadsTask.from_dict(data)
    assert len(task.dependencies) == 1
    assert task.dependencies[0].task_id == "bd-5678"


def test_beads_task_from_dict_dependency_missing_type() -> None:
    """Test BeadsTask.from_dict skips dependency without type (line 110)."""
    data = {
        "id": "bd-1234",
        "title": "Test",
        "dependencies": [
            {"task_id": "bd-5678"},  # Missing type
        ],
    }
    task = BeadsTask.from_dict(data)
    assert len(task.dependencies) == 0


def test_beads_task_from_dict_with_dates() -> None:
    """Test BeadsTask.from_dict with dates (lines 130-135)."""
    data = {
        "id": "bd-1234",
        "title": "Test",
        "created_at": "2024-01-01T12:00:00",
        "updated_at": "2024-01-01T13:00:00",
    }
    task = BeadsTask.from_dict(data)
    assert task.created_at == datetime(2024, 1, 1, 12, 0, 0)
    assert task.updated_at == datetime(2024, 1, 1, 13, 0, 0)


def test_beads_task_from_dict_without_dates() -> None:
    """Test BeadsTask.from_dict without dates (lines 130-135)."""
    data = {
        "id": "bd-1234",
        "title": "Test",
    }
    task = BeadsTask.from_dict(data)
    assert task.created_at is None
    assert task.updated_at is None


def test_beads_task_is_ready_blocked_status() -> None:
    """Test BeadsTask.is_ready returns False for blocked status (lines 153-154)."""
    task = BeadsTask(
        id="bd-1234",
        title="Test",
        status=BeadsStatus.BLOCKED,
    )
    assert not task.is_ready()


def test_beads_task_is_ready_closed_status() -> None:
    """Test BeadsTask.is_ready returns False for closed status (lines 153-154)."""
    task = BeadsTask(
        id="bd-1234",
        title="Test",
        status=BeadsStatus.CLOSED,
    )
    assert not task.is_ready()


def test_beads_task_is_ready_open_status() -> None:
    """Test BeadsTask.is_ready returns True for open status (line 157)."""
    task = BeadsTask(
        id="bd-1234",
        title="Test",
        status=BeadsStatus.OPEN,
    )
    assert task.is_ready()


def test_beads_task_create_to_dict() -> None:
    """Test BeadsTaskCreate.to_dict (lines 172-184)."""
    create = BeadsTaskCreate(
        title="Test",
        description="Test desc",
        priority=BeadsPriority.HIGH,
        parent_id="bd-parent",
        dependencies=[
            BeadsDependency(task_id="bd-dep", type=BeadsDependencyType.BLOCKS)
        ],
        external_ref="00-001-01",
        sdp_metadata={"foo": "bar"},
    )
    data = create.to_dict()
    assert data["title"] == "Test"
    assert data["description"] == "Test desc"
    assert data["priority"] == 1  # HIGH
    assert data["parent_id"] == "bd-parent"
    assert len(data["dependencies"]) == 1
    assert data["external_ref"] == "00-001-01"
    assert data["metadata"] == {"sdp": {"foo": "bar"}}


def test_beads_task_create_to_dict_no_metadata() -> None:
    """Test BeadsTaskCreate.to_dict without metadata (line 183)."""
    create = BeadsTaskCreate(
        title="Test",
        sdp_metadata={},
    )
    data = create.to_dict()
    assert data["metadata"] is None
