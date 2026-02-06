"""Tests for feature domain entities."""

import pytest

from sdp.domain.exceptions import DependencyCycleError, MissingDependencyError
from sdp.domain.feature import Feature
from sdp.domain.workstream import Workstream, WorkstreamSize, WorkstreamStatus


class TestFeature:
    """Tests for Feature aggregate."""

    def test_empty_feature(self) -> None:
        """Create feature with no workstreams."""
        feature = Feature(feature_id="F001")
        assert feature.feature_id == "F001"
        assert feature.workstreams == []
        assert feature.execution_order == []

    def test_single_workstream(self) -> None:
        """Feature with single workstream."""
        ws = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
        )
        feature = Feature(feature_id="F001", workstreams=[ws])
        assert len(feature.workstreams) == 1
        assert feature.execution_order == ["00-001-01"]

    def test_linear_dependencies(self) -> None:
        """Feature with linear dependency chain."""
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
        feature = Feature(feature_id="F001", workstreams=[ws3, ws1, ws2])
        # Should order: ws1 -> ws2 -> ws3
        assert feature.execution_order == ["00-001-01", "00-001-02", "00-001-03"]

    def test_parallel_workstreams(self) -> None:
        """Feature with parallel workstreams (no dependencies)."""
        ws1 = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
        )
        ws2 = Workstream(
            ws_id="00-001-02",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
        )
        feature = Feature(feature_id="F001", workstreams=[ws1, ws2])
        # Both can run in parallel - order not guaranteed but both present
        assert set(feature.execution_order) == {"00-001-01", "00-001-02"}

    def test_circular_dependency_detection(self) -> None:
        """Detect circular dependencies."""
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
        with pytest.raises(DependencyCycleError) as exc_info:
            Feature(feature_id="F001", workstreams=[ws1, ws2])
        assert len(exc_info.value.cycle) == 2

    def test_missing_dependency_error(self) -> None:
        """Raise error if dependency not found."""
        ws = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            dependencies=["00-001-99"],  # Does not exist
        )
        with pytest.raises(MissingDependencyError) as exc_info:
            Feature(feature_id="F001", workstreams=[ws])
        assert exc_info.value.ws_id == "00-001-01"
        assert exc_info.value.missing_dep == "00-001-99"

    def test_get_workstream(self) -> None:
        """Retrieve workstream by ID."""
        ws = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
        )
        feature = Feature(feature_id="F001", workstreams=[ws])
        found = feature.get_workstream("00-001-01")
        assert found is not None
        assert found.ws_id == "00-001-01"

    def test_get_workstream_not_found(self) -> None:
        """Return None if workstream not found."""
        feature = Feature(feature_id="F001")
        found = feature.get_workstream("00-001-99")
        assert found is None

    def test_get_dependencies(self) -> None:
        """Get direct dependencies for a workstream."""
        ws1 = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
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

    def test_complex_dependency_graph(self) -> None:
        """Feature with complex DAG."""
        # Diamond dependency:
        #   01
        #  /  \
        # 02  03
        #  \  /
        #   04
        ws1 = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
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
        # ws1 must be first, ws4 must be last
        order = feature.execution_order
        assert order[0] == "00-001-01"
        assert order[-1] == "00-001-04"
        assert "00-001-02" in order
        assert "00-001-03" in order
