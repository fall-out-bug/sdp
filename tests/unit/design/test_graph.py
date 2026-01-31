"""Tests for dependency graph."""

import pytest

from sdp.design.graph import DependencyGraph, WorkstreamNode


def test_topological_sort_simple():
    """Test simple dependency ordering."""
    graph = DependencyGraph()
    graph.add(WorkstreamNode("00-001-01", depends_on=[]))
    graph.add(WorkstreamNode("00-001-02", depends_on=["00-001-01"]))
    graph.add(WorkstreamNode("00-001-03", depends_on=["00-001-02"]))

    order = graph.topological_sort()

    assert order.index("00-001-01") < order.index("00-001-02")
    assert order.index("00-001-02") < order.index("00-001-03")


def test_topological_sort_branching():
    """Test branching dependencies."""
    graph = DependencyGraph()
    graph.add(WorkstreamNode("00-001-01", depends_on=[]))
    graph.add(WorkstreamNode("00-001-02", depends_on=["00-001-01"]))
    graph.add(WorkstreamNode("00-001-03", depends_on=["00-001-01"]))
    graph.add(WorkstreamNode("00-001-04", depends_on=["00-001-02", "00-001-03"]))

    order = graph.topological_sort()

    assert order[0] == "00-001-01"
    assert order.index("00-001-01") < order.index("00-001-02")
    assert order.index("00-001-01") < order.index("00-001-03")
    assert order.index("00-001-02") < order.index("00-001-04")
    assert order.index("00-001-03") < order.index("00-001-04")


def test_topological_sort_cycle_detection():
    """Test that cycles are detected."""
    graph = DependencyGraph()
    graph.add(WorkstreamNode("00-001-01", depends_on=["00-001-02"]))
    graph.add(WorkstreamNode("00-001-02", depends_on=["00-001-01"]))

    with pytest.raises(ValueError, match="Cycle"):
        graph.topological_sort()


def test_get_ready_workstreams():
    """Test getting ready workstreams."""
    graph = DependencyGraph()
    graph.add(WorkstreamNode("00-001-01", depends_on=[]))
    graph.add(WorkstreamNode("00-001-02", depends_on=["00-001-01"]))
    graph.add(WorkstreamNode("00-001-03", depends_on=["00-001-01"]))

    # Initially only 00-001-01 is ready
    assert graph.get_ready_workstreams([]) == ["00-001-01"]

    # After 00-001-01, both 02 and 03 are ready
    ready = graph.get_ready_workstreams(["00-001-01"])
    assert set(ready) == {"00-001-02", "00-001-03"}


def test_to_mermaid():
    """Test Mermaid graph generation."""
    graph = DependencyGraph()
    graph.add(WorkstreamNode("00-001-01", depends_on=[]))
    graph.add(WorkstreamNode("00-001-02", depends_on=["00-001-01"]))

    mermaid = graph.to_mermaid()

    assert "graph TD" in mermaid
    assert "00-001-01 --> 00-001-02" in mermaid


def test_get_node():
    """Test retrieving individual nodes."""
    graph = DependencyGraph()
    node = WorkstreamNode("00-001-01", estimated_loc=100)
    graph.add(node)

    retrieved = graph.get("00-001-01")
    assert retrieved is not None
    assert retrieved.ws_id == "00-001-01"
    assert retrieved.estimated_loc == 100


def test_get_missing_node():
    """Test retrieving non-existent node."""
    graph = DependencyGraph()
    assert graph.get("nonexistent") is None
