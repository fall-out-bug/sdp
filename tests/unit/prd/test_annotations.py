"""Tests for sdp.prd.annotations data structures."""

from pathlib import Path

from sdp.prd.annotations import Flow, FlowStep


def test_flow_add_step() -> None:
    """Test Flow.add_step() adds step to steps list."""
    flow = Flow(name="TestFlow", steps=[])
    step = FlowStep(
        flow_name="TestFlow",
        step_number=1,
        description="Test step",
        source_file=Path("/test/file.py"),
        line_number=10,
    )
    
    flow.add_step(step)
    
    assert len(flow.steps) == 1
    assert flow.steps[0] == step


def test_flow_get_sorted_steps() -> None:
    """Test Flow.get_sorted_steps() returns steps sorted by step_number."""
    flow = Flow(name="TestFlow", steps=[])
    step3 = FlowStep(
        flow_name="TestFlow",
        step_number=3,
        description="Third step",
        source_file=Path("/test/file.py"),
        line_number=30,
    )
    step1 = FlowStep(
        flow_name="TestFlow",
        step_number=1,
        description="First step",
        source_file=Path("/test/file.py"),
        line_number=10,
    )
    step2 = FlowStep(
        flow_name="TestFlow",
        step_number=2,
        description="Second step",
        source_file=Path("/test/file.py"),
        line_number=20,
    )
    
    flow.steps = [step3, step1, step2]
    sorted_steps = flow.get_sorted_steps()
    
    assert len(sorted_steps) == 3
    assert sorted_steps[0].step_number == 1
    assert sorted_steps[1].step_number == 2
    assert sorted_steps[2].step_number == 3


def test_flow_len() -> None:
    """Test Flow.__len__() returns correct step count."""
    flow = Flow(name="TestFlow", steps=[])
    assert len(flow) == 0
    
    flow.steps.append(
        FlowStep(
            flow_name="TestFlow",
            step_number=1,
            description="Step 1",
            source_file=Path("/test/file.py"),
            line_number=10,
        )
    )
    assert len(flow) == 1
    
    flow.steps.append(
        FlowStep(
            flow_name="TestFlow",
            step_number=2,
            description="Step 2",
            source_file=Path("/test/file.py"),
            line_number=20,
        )
    )
    assert len(flow) == 2
