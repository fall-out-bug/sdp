"""Tests for Orchestrator."""

import asyncio
import json
import pytest
from pathlib import Path

from sdp.agents.orchestrator import (
    Orchestrator,
    OrchestratorState,
    ExecutionResult,
)


@pytest.fixture
def sample_workstreams(tmp_path: Path) -> Path:
    """Create sample workstream files."""
    backlog = tmp_path / "backlog"
    backlog.mkdir(parents=True)

    (backlog / "00-001-01.md").write_text("""---
ws_id: 00-001-01
status: backlog
title: First
feature: F001
---
# First
""")

    (backlog / "00-001-02.md").write_text("""---
ws_id: 00-001-02
status: backlog
title: Second
feature: F001
depends_on: ["00-001-01"]
---
# Second
""")

    return tmp_path


@pytest.mark.asyncio
async def test_orchestrator_initialization(sample_workstreams: Path) -> None:
    """Test orchestrator initializes correctly."""
    orch = Orchestrator(max_agents=2, ws_dir=sample_workstreams)

    assert orch._max_agents == 2
    assert orch._ws_dir == sample_workstreams
    assert orch._state is None


@pytest.mark.asyncio
async def test_orchestrator_get_feature_workstreams(sample_workstreams: Path) -> None:
    """Test getting workstreams for a feature."""
    orch = Orchestrator(ws_dir=sample_workstreams)

    ws_ids = orch._get_feature_workstreams("F001")

    assert len(ws_ids) == 2
    assert "00-001-01" in ws_ids
    assert "00-001-02" in ws_ids


@pytest.mark.asyncio
async def test_orchestrator_get_feature_workstreams_empty(tmp_path: Path) -> None:
    """Test getting workstreams for non-existent feature."""
    orch = Orchestrator(ws_dir=tmp_path)

    ws_ids = orch._get_feature_workstreams("F999")

    assert ws_ids == []


@pytest.mark.asyncio
async def test_orchestrator_get_stats(sample_workstreams: Path) -> None:
    """Test getting pool statistics."""
    orch = Orchestrator(ws_dir=sample_workstreams)

    stats = orch.get_stats()

    assert stats.total_agents == 3  # default
    assert stats.busy_agents == 0
    assert stats.available_agents == 3


@pytest.mark.asyncio
async def test_orchestrator_save_and_load_state(sample_workstreams: Path, tmp_path: Path) -> None:
    """Test state persistence."""
    from datetime import datetime

    state_file = tmp_path / "orch_state.json"
    orch = Orchestrator(ws_dir=sample_workstreams, state_file=state_file)

    state = OrchestratorState(
        feature_id="F001",
        started_at=datetime(2024, 1, 1, 0, 0, 0),
        status="running",
        results={"00-001-01": True},
        current_ws="00-001-02",
    )

    orch._state = state
    orch._save_state()

    # Load with new instance
    orch2 = Orchestrator(ws_dir=sample_workstreams, state_file=state_file)
    loaded = orch2.load_state()

    assert loaded is not None
    assert loaded.feature_id == "F001"
    assert loaded.status == "running"
    assert loaded.results == {"00-001-01": True}
    assert loaded.current_ws == "00-001-02"


@pytest.mark.asyncio
async def test_orchestrator_load_state_no_file(tmp_path: Path) -> None:
    """Test loading when no state file exists."""
    state_file = tmp_path / "nonexistent.json"
    orch = Orchestrator(state_file=state_file)

    state = orch.load_state()

    assert state is None


@pytest.mark.asyncio
async def test_orchestrator_on_progress(sample_workstreams: Path) -> None:
    """Test progress callback registration."""
    orch = Orchestrator(ws_dir=sample_workstreams)

    received = []

    def callback(ws_id: str, success: bool, error: str | None) -> None:
        received.append((ws_id, success, error))

    orch.on_progress(callback)

    # Trigger callbacks
    for cb in orch._callbacks:
        cb("00-001-01", True, None)

    assert received == [("00-001-01", True, None)]


@pytest.mark.asyncio
async def test_orchestrator_enqueue_feature(sample_workstreams: Path) -> None:
    """Test enqueuing feature workstreams."""
    from sdp.queue.task_queue import Priority

    orch = Orchestrator(ws_dir=sample_workstreams)

    task_ids = await orch.enqueue_feature("F001", Priority.URGENT)

    assert len(task_ids) == 2

    # Check queue size
    assert orch._queue.size() == 2


@pytest.mark.asyncio
async def test_orchestrator_run_feature_empty(tmp_path: Path) -> None:
    """Test running feature with no workstreams."""
    orch = Orchestrator(ws_dir=tmp_path)

    result = await orch.run_feature("F999")

    assert result.feature_id == "F999"
    assert result.success is True  # Empty is "success"
    assert len(result.completed) == 0


@pytest.mark.asyncio
async def test_orchestrator_run_feature_with_circular_deps(tmp_path: Path) -> None:
    """Test running feature with circular dependencies."""
    backlog = tmp_path / "backlog"
    backlog.mkdir(parents=True)

    # Create circular dependency
    (backlog / "00-001-01.md").write_text("""---
ws_id: 00-001-01
depends_on: ["00-001-02"]
feature: F001
---
# A
""")
    (backlog / "00-001-02.md").write_text("""---
ws_id: 00-001-02
depends_on: ["00-001-01"]
feature: F001
---
# B
""")

    orch = Orchestrator(ws_dir=tmp_path)

    with pytest.raises(ValueError, match="Circular dependency"):
        await orch.run_feature("F001")


@pytest.mark.asyncio
async def test_orchestrator_run_feature_ordered(sample_workstreams: Path) -> None:
    """Test running feature in ordered mode."""
    orch = Orchestrator(max_agents=1, ws_dir=sample_workstreams)

    result = await orch.run_feature_ordered("F001")

    assert result.feature_id == "F001"
    # Should have both workstreams in results
    assert len(result.completed) + len(result.failed) == 2


@pytest.mark.asyncio
async def test_orchestrator_group_by_level(sample_workstreams: Path) -> None:
    """Test grouping workstreams by dependency level."""
    orch = Orchestrator(ws_dir=sample_workstreams)

    # Build graph
    ws_ids = ["00-001-01", "00-001-02"]
    orch._graph.build(ws_ids)

    levels = orch._group_by_level([*ws_ids])

    # 001 is level 0, 002 is level 1
    assert levels[0][0] == 0
    assert "00-001-01" in levels[0][1]

    assert levels[1][0] == 1
    assert "00-001-02" in levels[1][1]


def test_execution_result() -> None:
    """Test ExecutionResult dataclass."""
    result = ExecutionResult(
        feature_id="F001",
        success=True,
        completed=["ws-001", "ws-002"],
        failed=[],
        duration_seconds=10.5,
        errors={},
    )

    assert result.feature_id == "F001"
    assert result.success is True
    assert result.duration_seconds == 10.5


def test_orchestrator_state_defaults() -> None:
    """Test OrchestratorState default values."""
    from datetime import datetime

    state = OrchestratorState(
        feature_id="F001",
        started_at=datetime.now(),
    )

    assert state.status == "running"
    assert state.results == {}
    assert state.current_ws is None
    assert state.errors == {}
