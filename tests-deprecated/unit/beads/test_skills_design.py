"""Tests for skills_design.py - Feature decomposition into workstreams."""

import pytest

from sdp.beads.client import MockBeadsClient
from sdp.beads.skills_design import FeatureDecomposer, WorkstreamSpec
from sdp.beads.models import BeadsDependencyType


class TestWorkstreamSpec:
    """Test WorkstreamSpec dataclass."""

    def test_create_with_defaults(self) -> None:
        """Test creating WorkstreamSpec with defaults."""
        spec = WorkstreamSpec(title="Test WS")
        assert spec.title == "Test WS"
        assert spec.sequence == 1
        assert spec.size == "MEDIUM"
        assert spec.dependencies == []

    def test_create_with_all_fields(self) -> None:
        """Test creating WorkstreamSpec with all fields."""
        spec = WorkstreamSpec(
            title="Custom WS",
            sequence=5,
            size="LARGE",
            dependencies=["bd-0001", "bd-0002"],
        )
        assert spec.title == "Custom WS"
        assert spec.sequence == 5
        assert spec.size == "LARGE"
        assert spec.dependencies == ["bd-0001", "bd-0002"]


class TestFeatureDecomposer:
    """Test FeatureDecomposer class."""

    def test_init(self) -> None:
        """Test decomposer initialization."""
        client = MockBeadsClient()
        decomposer = FeatureDecomposer(client)
        assert decomposer.client is client

    def test_decompose_with_default_workstreams(self) -> None:
        """Test decompose with default workstreams."""
        client = MockBeadsClient()
        decomposer = FeatureDecomposer(client)

        # Create parent feature
        from sdp.beads.models import BeadsTaskCreate, BeadsPriority
        feature = client.create_task(
            BeadsTaskCreate(title="Feature", description="", priority=BeadsPriority.HIGH)
        )

        ws_ids = decomposer.decompose(feature.id)

        # Should create 3 default workstreams
        assert len(ws_ids) == 3

        # Verify sequential dependencies
        ws1 = client.get_task(ws_ids[0])
        ws2 = client.get_task(ws_ids[1])
        ws3 = client.get_task(ws_ids[2])

        assert ws1 is not None
        assert ws2 is not None
        assert ws3 is not None

        # ws1 has no dependencies
        assert len(ws1.dependencies) == 0

        # ws2 depends on ws1
        assert len(ws2.dependencies) == 1
        assert ws2.dependencies[0].task_id == ws_ids[0]
        assert ws2.dependencies[0].type == BeadsDependencyType.BLOCKS

        # ws3 depends on ws2
        assert len(ws3.dependencies) == 1
        assert ws3.dependencies[0].task_id == ws_ids[1]
        assert ws3.dependencies[0].type == BeadsDependencyType.BLOCKS

    def test_decompose_with_custom_workstreams(self) -> None:
        """Test decompose with custom workstream specs."""
        client = MockBeadsClient()
        decomposer = FeatureDecomposer(client)

        # Create parent feature
        from sdp.beads.models import BeadsTaskCreate, BeadsPriority
        feature = client.create_task(
            BeadsTaskCreate(title="Feature", description="", priority=BeadsPriority.HIGH)
        )

        custom_ws = [
            WorkstreamSpec(title="Setup", sequence=1, size="SMALL"),
            WorkstreamSpec(title="Implementation", sequence=2, size="LARGE"),
            WorkstreamSpec(title="Testing", sequence=3, size="MEDIUM"),
            WorkstreamSpec(title="Documentation", sequence=4, size="SMALL"),
        ]

        ws_ids = decomposer.decompose(feature.id, workstreams=custom_ws)

        assert len(ws_ids) == 4

        # Verify titles
        ws1 = client.get_task(ws_ids[0])
        ws2 = client.get_task(ws_ids[1])
        ws3 = client.get_task(ws_ids[2])
        ws4 = client.get_task(ws_ids[3])

        assert ws1.title == "Setup"
        assert ws2.title == "Implementation"
        assert ws3.title == "Testing"
        assert ws4.title == "Documentation"

        # Verify sequential dependencies
        assert len(ws1.dependencies) == 0
        assert len(ws2.dependencies) == 1
        assert ws2.dependencies[0].task_id == ws_ids[0]
        assert len(ws3.dependencies) == 1
        assert ws3.dependencies[0].task_id == ws_ids[1]
        assert len(ws4.dependencies) == 1
        assert ws4.dependencies[0].task_id == ws_ids[2]

    def test_decompose_single_workstream(self) -> None:
        """Test decompose with single workstream."""
        client = MockBeadsClient()
        decomposer = FeatureDecomposer(client)

        from sdp.beads.models import BeadsTaskCreate, BeadsPriority
        feature = client.create_task(
            BeadsTaskCreate(title="Feature", description="", priority=BeadsPriority.HIGH)
        )

        single_ws = [WorkstreamSpec(title="Single WS", sequence=1)]

        ws_ids = decomposer.decompose(feature.id, workstreams=single_ws)

        assert len(ws_ids) == 1
        ws = client.get_task(ws_ids[0])
        assert ws is not None
        assert ws.title == "Single WS"
        assert len(ws.dependencies) == 0  # No previous workstream

    def test_decompose_with_existing_dependencies(self) -> None:
        """Test decompose with workstreams that have existing dependencies."""
        client = MockBeadsClient()
        decomposer = FeatureDecomposer(client)

        from sdp.beads.models import BeadsTaskCreate, BeadsPriority
        feature = client.create_task(
            BeadsTaskCreate(title="Feature", description="", priority=BeadsPriority.HIGH)
        )

        # Create external blocker
        blocker = client.create_task(
            BeadsTaskCreate(title="Blocker", description="", priority=BeadsPriority.HIGH)
        )

        custom_ws = [
            WorkstreamSpec(
                title="WS1",
                sequence=1,
                dependencies=[blocker.id],  # External dependency
            ),
            WorkstreamSpec(title="WS2", sequence=2),
        ]

        ws_ids = decomposer.decompose(feature.id, workstreams=custom_ws)

        # Note: The current implementation doesn't preserve WorkstreamSpec.dependencies
        # This test documents current behavior - WS1 should only depend on blocker
        # if the implementation is updated to support it
        ws1 = client.get_task(ws_ids[0])
        ws2 = client.get_task(ws_ids[1])

        # Current implementation: WS1 has no dependencies (WorkstreamSpec.dependencies ignored)
        # WS2 depends on WS1 (sequential dependency)
        assert len(ws1.dependencies) == 0
        assert len(ws2.dependencies) == 1
        assert ws2.dependencies[0].task_id == ws_ids[0]

    def test_decompose_workstream_metadata(self) -> None:
        """Test that workstreams include SDP metadata."""
        client = MockBeadsClient()
        decomposer = FeatureDecomposer(client)

        from sdp.beads.models import BeadsTaskCreate, BeadsPriority
        feature = client.create_task(
            BeadsTaskCreate(title="Feature", description="", priority=BeadsPriority.HIGH)
        )

        ws_ids = decomposer.decompose(feature.id)

        ws1 = client.get_task(ws_ids[0])
        assert ws1 is not None
        assert "sequence" in ws1.sdp_metadata
        assert "size" in ws1.sdp_metadata
        assert ws1.sdp_metadata["sequence"] == 1
        assert ws1.sdp_metadata["size"] == "MEDIUM"

    def test_decompose_parent_id_set(self) -> None:
        """Test that workstreams have correct parent_id."""
        client = MockBeadsClient()
        decomposer = FeatureDecomposer(client)

        from sdp.beads.models import BeadsTaskCreate, BeadsPriority
        feature = client.create_task(
            BeadsTaskCreate(title="Feature", description="", priority=BeadsPriority.HIGH)
        )

        ws_ids = decomposer.decompose(feature.id)

        for ws_id in ws_ids:
            ws = client.get_task(ws_id)
            assert ws is not None
            assert ws.parent_id == feature.id
