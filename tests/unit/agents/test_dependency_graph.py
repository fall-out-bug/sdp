"""Tests for DependencyGraph."""

import pytest
from pathlib import Path

from sdp.agents.dependency_graph import DependencyGraph, Node


@pytest.fixture
def sample_workstreams(tmp_path: Path) -> Path:
    """Create sample workstream files with dependencies."""
    backlog = tmp_path / "backlog"
    backlog.mkdir(parents=True)

    # WS-001: No dependencies
    (backlog / "00-001-01.md").write_text("""---
ws_id: 00-001-01
status: backlog
title: Foundation
feature: F001
---
# Foundation
No dependencies.
""")

    # WS-002: Depends on WS-001
    (backlog / "00-001-02.md").write_text("""---
ws_id: 00-001-02
status: backlog
title: Depends on 001
feature: F001
depends_on: ["00-001-01"]
---
# Depends on 001
""")

    # WS-003: Depends on WS-002
    (backlog / "00-001-03.md").write_text("""---
ws_id: 00-001-03
status: backlog
title: Depends on 002
feature: F001
depends_on: ["00-001-02"]
---
# Depends on 002
""")

    return tmp_path


def test_dependency_graph_initialization() -> None:
    """Test graph initializes correctly."""
    graph = DependencyGraph()
    assert graph._built is False
    assert len(graph._nodes) == 0


def test_dependency_graph_build(sample_workstreams: Path) -> None:
    """Test building graph from workstream files."""
    graph = DependencyGraph(sample_workstreams)
    graph.build(["00-001-01", "00-001-02", "00-001-03"])

    assert graph._built is True
    assert len(graph._nodes) == 3


def test_dependency_graph_execution_order(sample_workstreams: Path) -> None:
    """Test topological sort returns correct order."""
    graph = DependencyGraph(sample_workstreams)
    graph.build(["00-001-01", "00-001-02", "00-001-03"])

    order = graph.execution_order()

    # 001 must come before 002, which must come before 003
    assert order.index("00-001-01") < order.index("00-001-02")
    assert order.index("00-001-02") < order.index("00-001-03")


def test_dependency_graph_get_ready(sample_workstreams: Path) -> None:
    """Test getting ready workstreams."""
    graph = DependencyGraph(sample_workstreams)
    graph.build(["00-001-01", "00-001-02", "00-001-03"])

    # Initially only WS-001 is ready
    ready = graph.get_ready(set())
    assert ready == ["00-001-01"]

    # After completing WS-001, WS-002 is ready
    ready = graph.get_ready({"00-001-01"})
    assert "00-001-02" in ready

    # After completing WS-002, WS-003 is ready
    ready = graph.get_ready({"00-001-01", "00-001-02"})
    assert "00-001-03" in ready


def test_dependency_graph_circular_detection(tmp_path: Path) -> None:
    """Test circular dependency detection."""
    backlog = tmp_path / "backlog"
    backlog.mkdir(parents=True)

    # Create circular dependency: A -> B -> A
    (backlog / "00-001-01.md").write_text("""---
ws_id: 00-001-01
depends_on: ["00-001-02"]
---
# A
""")
    (backlog / "00-001-02.md").write_text("""---
ws_id: 00-001-02
depends_on: ["00-001-01"]
---
# B
""")

    graph = DependencyGraph(tmp_path)
    graph.build(["00-001-01", "00-001-02"])

    with pytest.raises(ValueError, match="Circular dependency"):
        graph.execution_order()


def test_dependency_graph_has_circular_dependency(tmp_path: Path) -> None:
    """Test has_circular_dependency method."""
    backlog = tmp_path / "backlog"
    backlog.mkdir(parents=True)

    # No circular deps
    (backlog / "00-001-01.md").write_text("""---
ws_id: 00-001-01
---
# A
""")

    graph = DependencyGraph(tmp_path)
    graph.build(["00-001-01"])
    assert graph.has_circular_dependency() is False

    # Add circular dependency
    (backlog / "00-001-02.md").write_text("""---
ws_id: 00-001-02
depends_on: ["00-001-01"]
---
# B
""")
    (backlog / "00-001-01.md").write_text("""---
ws_id: 00-001-01
depends_on: ["00-001-02"]
---
# A
""")

    graph = DependencyGraph(tmp_path)
    graph.build(["00-001-01", "00-001-02"])
    assert graph.has_circular_dependency() is True


def test_dependency_graph_get_dependencies(sample_workstreams: Path) -> None:
    """Test getting dependencies for a workstream."""
    graph = DependencyGraph(sample_workstreams)
    graph.build(["00-001-01", "00-001-02", "00-001-03"])

    deps = graph.get_dependencies("00-001-03")
    assert deps == {"00-001-02"}

    deps = graph.get_dependencies("00-001-01")
    assert deps == set()


def test_dependency_graph_get_dependents(sample_workstreams: Path) -> None:
    """Test getting dependents of a workstream."""
    graph = DependencyGraph(sample_workstreams)
    graph.build(["00-001-01", "00-001-02", "00-001-03"])

    dependents = graph.get_dependents("00-001-01")
    assert "00-001-02" in dependents

    dependents = graph.get_dependents("00-001-03")
    assert len(dependents) == 0


def test_dependency_graph_visualize(sample_workstreams: Path) -> None:
    """Test graph visualization."""
    graph = DependencyGraph(sample_workstreams)
    graph.build(["00-001-01", "00-001-02", "00-001-03"])

    viz = graph.visualize()
    assert "Dependency Graph:" in viz
    assert "00-001-01" in viz
    assert "00-001-02" in viz
    assert "00-001-03" in viz


def test_dependency_graph_empty_build() -> None:
    """Test building with empty list."""
    graph = DependencyGraph()
    graph.build([])

    assert graph._built is True
    assert len(graph._nodes) == 0
    assert graph.execution_order() == []


def test_dependency_graph_not_built() -> None:
    """Test methods when graph not built."""
    graph = DependencyGraph()

    assert graph.execution_order() == []
    assert graph.get_ready(set()) == []
    assert graph.get_dependencies("nonexistent") == set()
    assert graph.get_dependents("nonexistent") == set()
