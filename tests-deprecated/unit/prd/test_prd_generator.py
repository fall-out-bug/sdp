"""Tests for PRD diagram generator."""

from pathlib import Path

import pytest

from sdp.prd.annotations import Flow, FlowStep
from sdp.prd.generator import generate_diagrams, generate_flow_from_steps


class TestGenerateFlowFromSteps:
    """Test generate_flow_from_steps function."""

    def test_generate_flow_single_step(self, tmp_path: Path) -> None:
        """Test generating flow with single step."""
        step = FlowStep(
            flow_name="test_flow",
            step_number=1,
            description="First step",
            source_file=tmp_path / "test.py",
            line_number=10,
        )

        flow = generate_flow_from_steps("test_flow", [step])

        assert flow.name == "test_flow"
        assert len(flow.steps) == 1
        assert flow.steps[0].step_number == 1

    def test_generate_flow_multiple_steps(self, tmp_path: Path) -> None:
        """Test generating flow with multiple steps."""
        steps = [
            FlowStep(
                flow_name="test_flow",
                step_number=1,
                description="Step 1",
                source_file=tmp_path / "test.py",
                line_number=10,
            ),
            FlowStep(
                flow_name="test_flow",
                step_number=2,
                description="Step 2",
                source_file=tmp_path / "test.py",
                line_number=20,
            ),
            FlowStep(
                flow_name="test_flow",
                step_number=3,
                description="Step 3",
                source_file=tmp_path / "test.py",
                line_number=30,
            ),
        ]

        flow = generate_flow_from_steps("test_flow", steps)

        assert flow.name == "test_flow"
        assert len(flow.steps) == 3

    def test_generate_flow_empty_steps(self) -> None:
        """Test generating flow with empty steps list."""
        flow = generate_flow_from_steps("empty_flow", [])

        assert flow.name == "empty_flow"
        assert len(flow.steps) == 0

    def test_generate_flow_with_participant(self, tmp_path: Path) -> None:
        """Test generating flow with participant."""
        step = FlowStep(
            flow_name="test_flow",
            step_number=1,
            description="Step with participant",
            source_file=tmp_path / "test.py",
            line_number=10,
            participant="User",
        )

        flow = generate_flow_from_steps("test_flow", [step])

        assert flow.steps[0].participant == "User"


class TestGenerateDiagrams:
    """Test generate_diagrams function."""

    def test_generate_diagrams_single_flow(self, tmp_path: Path) -> None:
        """Test generating diagrams for single flow."""
        output_dir = tmp_path / "diagrams"
        step = FlowStep(
            flow_name="test_flow",
            step_number=1,
            description="Test step",
            source_file=tmp_path / "test.py",
            line_number=10,
        )
        flow = Flow(name="test_flow", steps=[step])

        # Default project_type is "service", so includes component/deployment diagrams
        created_files = generate_diagrams([flow], output_dir)

        assert len(created_files) >= 2  # At least Mermaid + PlantUML sequence
        assert output_dir.exists()
        assert any("sequence-test_flow.mmd" in str(f) for f in created_files)
        assert any("sequence-test_flow.puml" in str(f) for f in created_files)

    def test_generate_diagrams_multiple_flows(self, tmp_path: Path) -> None:
        """Test generating diagrams for multiple flows."""
        output_dir = tmp_path / "diagrams"
        flow1 = Flow(
            name="flow1",
            steps=[
                FlowStep(
                    flow_name="flow1",
                    step_number=1,
                    description="Step 1",
                    source_file=tmp_path / "test1.py",
                    line_number=10,
                ),
            ],
        )
        flow2 = Flow(
            name="flow2",
            steps=[
                FlowStep(
                    flow_name="flow2",
                    step_number=1,
                    description="Step 1",
                    source_file=tmp_path / "test2.py",
                    line_number=10,
                ),
            ],
        )

        # Default project_type is "service", so includes component/deployment diagrams
        created_files = generate_diagrams([flow1, flow2], output_dir)

        assert len(created_files) >= 4  # At least 2 flows × 2 formats
        assert any("sequence-flow1.mmd" in str(f) for f in created_files)
        assert any("sequence-flow2.mmd" in str(f) for f in created_files)

    def test_generate_diagrams_service_type(self, tmp_path: Path) -> None:
        """Test generating diagrams for service project type."""
        output_dir = tmp_path / "diagrams"
        flow = Flow(
            name="test_flow",
            steps=[
                FlowStep(
                    flow_name="test_flow",
                    step_number=1,
                    description="Test step",
                    source_file=tmp_path / "test.py",
                    line_number=10,
                ),
            ],
        )

        created_files = generate_diagrams([flow], output_dir, project_type="service")

        # Should include component and deployment diagrams for service
        assert len(created_files) >= 4
        assert any("component-overview.mmd" in str(f) for f in created_files)
        assert any("deployment-production.puml" in str(f) for f in created_files)

    def test_generate_diagrams_library_type(self, tmp_path: Path) -> None:
        """Test generating diagrams for library project type."""
        output_dir = tmp_path / "diagrams"
        flow = Flow(
            name="test_flow",
            steps=[
                FlowStep(
                    flow_name="test_flow",
                    step_number=1,
                    description="Test step",
                    source_file=tmp_path / "test.py",
                    line_number=10,
                ),
            ],
        )

        created_files = generate_diagrams([flow], output_dir, project_type="library")

        # Library type should not include component/deployment diagrams
        assert len(created_files) == 2  # Only sequence diagrams

    def test_generate_diagrams_creates_directory(self, tmp_path: Path) -> None:
        """Test that generate_diagrams creates output directory."""
        output_dir = tmp_path / "nonexistent" / "diagrams"
        flow = Flow(
            name="test_flow",
            steps=[
                FlowStep(
                    flow_name="test_flow",
                    step_number=1,
                    description="Test step",
                    source_file=tmp_path / "test.py",
                    line_number=10,
                ),
            ],
        )

        assert not output_dir.exists()
        generate_diagrams([flow], output_dir)
        assert output_dir.exists()

    def test_generate_diagrams_empty_flow(self, tmp_path: Path) -> None:
        """Test generating diagrams for flow with no steps."""
        output_dir = tmp_path / "diagrams"
        flow = Flow(name="empty_flow", steps=[])

        # Default project_type is "service", so includes component/deployment diagrams
        created_files = generate_diagrams([flow], output_dir)

        # Should still create diagram files even for empty flow
        assert len(created_files) >= 2
        assert any("sequence-empty_flow.mmd" in str(f) for f in created_files)

    def test_generate_diagrams_flow_with_participants(self, tmp_path: Path) -> None:
        """Test generating diagrams for flow with participants."""
        output_dir = tmp_path / "diagrams"
        flow = Flow(
            name="test_flow",
            steps=[
                FlowStep(
                    flow_name="test_flow",
                    step_number=1,
                    description="Step 1",
                    source_file=tmp_path / "test.py",
                    line_number=10,
                    participant="User",
                ),
                FlowStep(
                    flow_name="test_flow",
                    step_number=2,
                    description="Step 2",
                    source_file=tmp_path / "api.py",
                    line_number=20,
                    participant="API",
                ),
            ],
        )

        created_files = generate_diagrams([flow], output_dir)

        # Check that diagrams were created
        assert len(created_files) >= 2

        # Check mermaid file content includes participants
        mmd_file = next(f for f in created_files if "sequence-test_flow.mmd" in str(f))
        content = mmd_file.read_text()
        assert "participant" in content.lower()
