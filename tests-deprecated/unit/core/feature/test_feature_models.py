"""Tests for sdp.core.feature.models."""

import pytest

from sdp.core.feature.errors import CircularDependencyError, MissingDependencyError
from sdp.core.feature.models import Feature
from sdp.domain.workstream import Workstream, WorkstreamSize, WorkstreamStatus


def test_feature_single_workstream() -> None:
    """Test feature with single workstream and no dependencies."""
    ws = Workstream(
        ws_id="00-001-01",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=[],
    )
    feature = Feature(feature_id="F001", workstreams=[ws])
    
    assert feature.feature_id == "F001"
    assert len(feature.workstreams) == 1
    assert feature.execution_order == ["00-001-01"]


def test_feature_linear_dependencies() -> None:
    """Test feature with linear dependency chain."""
    ws1 = Workstream(
        ws_id="00-001-01",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=[],
    )
    ws2 = Workstream(
        ws_id="00-001-02",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=["00-001-01"],
    )
    ws3 = Workstream(
        ws_id="00-001-03",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=["00-001-02"],
    )
    feature = Feature(feature_id="F001", workstreams=[ws1, ws2, ws3])
    
    assert feature.execution_order == ["00-001-01", "00-001-02", "00-001-03"]


def test_feature_parallel_workstreams() -> None:
    """Test feature with independent parallel workstreams."""
    ws1 = Workstream(
        ws_id="00-001-01",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=[],
    )
    ws2 = Workstream(
        ws_id="00-001-02",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=[],
    )
    feature = Feature(feature_id="F001", workstreams=[ws1, ws2])
    
    # Both can execute in parallel (order doesn't matter)
    assert len(feature.execution_order) == 2
    assert set(feature.execution_order) == {"00-001-01", "00-001-02"}


def test_feature_diamond_dependencies() -> None:
    """Test feature with diamond-shaped dependency graph."""
    ws1 = Workstream(
        ws_id="00-001-01",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=[],
    )
    ws2 = Workstream(
        ws_id="00-001-02",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=["00-001-01"],
    )
    ws3 = Workstream(
        ws_id="00-001-03",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=["00-001-01"],
    )
    ws4 = Workstream(
        ws_id="00-001-04",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=["00-001-02", "00-001-03"],
    )
    feature = Feature(feature_id="F001", workstreams=[ws1, ws2, ws3, ws4])
    
    # ws1 must come first, ws4 must come last
    assert feature.execution_order[0] == "00-001-01"
    assert feature.execution_order[-1] == "00-001-04"
    # ws2 and ws3 can be in any order in the middle
    assert set(feature.execution_order[1:3]) == {"00-001-02", "00-001-03"}


def test_feature_get_workstream() -> None:
    """Test getting workstream by ID."""
    ws1 = Workstream(
        ws_id="00-001-01",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=[],
    )
    ws2 = Workstream(
        ws_id="00-001-02",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=[],
    )
    feature = Feature(feature_id="F001", workstreams=[ws1, ws2])
    
    result = feature.get_workstream("00-001-01")
    assert result == ws1
    
    result = feature.get_workstream("00-001-02")
    assert result == ws2
    
    result = feature.get_workstream("00-001-99")
    assert result is None


def test_feature_get_dependencies() -> None:
    """Test getting dependencies for a workstream."""
    ws1 = Workstream(
        ws_id="00-001-01",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=[],
    )
    ws2 = Workstream(
        ws_id="00-001-02",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=["00-001-01"],
    )
    feature = Feature(feature_id="F001", workstreams=[ws1, ws2])
    
    deps = feature.get_dependencies("00-001-02")
    assert deps == ["00-001-01"]
    
    deps = feature.get_dependencies("00-001-01")
    assert deps == []
    
    # Non-existent workstream returns empty list
    deps = feature.get_dependencies("00-001-99")
    assert deps == []


def test_feature_missing_dependency_raises() -> None:
    """Test that missing dependency raises MissingDependencyError."""
    ws = Workstream(
        ws_id="00-001-01",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=["00-001-99"],  # Missing dependency
    )
    
    with pytest.raises(MissingDependencyError) as exc_info:
        Feature(feature_id="F001", workstreams=[ws])
    
    assert "00-001-01" in str(exc_info.value)
    assert "00-001-99" in str(exc_info.value)


def test_feature_circular_dependency_raises() -> None:
    """Test that circular dependency raises CircularDependencyError."""
    ws1 = Workstream(
        ws_id="00-001-01",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=["00-001-02"],
    )
    ws2 = Workstream(
        ws_id="00-001-02",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=["00-001-01"],
    )
    
    with pytest.raises(CircularDependencyError) as exc_info:
        Feature(feature_id="F001", workstreams=[ws1, ws2])
    
    # Should mention the cycle
    assert "00-001-01" in str(exc_info.value)
    assert "00-001-02" in str(exc_info.value)


def test_feature_self_dependency_raises() -> None:
    """Test that self-referencing dependency raises CircularDependencyError."""
    ws = Workstream(
        ws_id="00-001-01",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=["00-001-01"],  # Self-reference
    )
    
    with pytest.raises(CircularDependencyError):
        Feature(feature_id="F001", workstreams=[ws])


def test_feature_complex_circular_dependency_raises() -> None:
    """Test complex circular dependency (A→B→C→A) raises error."""
    ws1 = Workstream(
        ws_id="00-001-01",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=["00-001-02"],
    )
    ws2 = Workstream(
        ws_id="00-001-02",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=["00-001-03"],
    )
    ws3 = Workstream(
        ws_id="00-001-03",
        feature="F001",
        status=WorkstreamStatus.BACKLOG,
        size=WorkstreamSize.SMALL,
        dependencies=["00-001-01"],
    )
    
    with pytest.raises(CircularDependencyError):
        Feature(feature_id="F001", workstreams=[ws1, ws2, ws3])
