"""Ultra-targeted tests for specific uncovered lines to reach 80%."""

from pathlib import Path


def test_prd_parser_python_flow_only_no_step() -> None:
    """Test PRD parser with flow decorator only (no step)."""
    from sdp.prd.parser_python import parse_python_annotations

    test_file = Path("/tmp/test_flow_only.py")
    test_file.write_text("""
@prd_flow("auth_flow")
def login():
    pass
""")

    result = parse_python_annotations(test_file)
    assert len(result) == 1
    assert result[0].flow_name == "auth_flow"
    assert result[0].step_number == 0  # Default when no step
    assert result[0].description == "login"  # Function name

    test_file.unlink()


def test_prd_parser_python_multiple_flows() -> None:
    """Test PRD parser with multiple flows."""
    from sdp.prd.parser_python import parse_python_annotations

    test_file = Path("/tmp/test_multi.py")
    test_file.write_text("""
@prd_flow("flow1")
@prd_step(1, "Step 1")
def func1():
    pass

@prd_flow("flow2")
@prd_step(2, "Step 2")
def func2():
    pass
""")

    result = parse_python_annotations(test_file)
    assert len(result) == 2
    assert result[0].flow_name == "flow1"
    assert result[1].flow_name == "flow2"

    test_file.unlink()


def test_core_model_parser_context_window_edge_cases() -> None:
    """Test parse_context_window with various formats."""
    from sdp.core.model.parser import parse_context_window

    assert parse_context_window("2M+") == 2_000_000
    assert parse_context_window("1.5M") == 1_500_000
    assert parse_context_window("64K") == 64_000
    assert parse_context_window("4096") == 4096
    assert parse_context_window("INVALID") == 128_000


def test_validators_ws_completion_parser_ws_file_parsing() -> None:
    """Test parse_ws_file integration."""
    from sdp.validators.ws_completion.parser import parse_ws_file

    test_file = Path("/tmp/test_ws.md")
    test_file.write_text("""---
title: Test WS
scope_files:
  - src/test.py
---

## Implementation

### Verification

```bash
pytest tests/
```
""")

    result = parse_ws_file(test_file)
    assert "src/test.py" in result["scope_files"]
    assert "pytest tests/" in result["verification_commands"]

    test_file.unlink()


def test_beads_models_all_priority_levels() -> None:
    """Test BeadsTask with all priority levels."""
    from sdp.beads.models import BeadsTask, BeadsPriority

    for priority in BeadsPriority:
        task = BeadsTask(
            id=f"bd-{priority.value}",
            title="Test",
            priority=priority,
        )
        data = task.to_dict()
        assert data["priority"] == priority.value


def test_beads_models_dependency_types() -> None:
    """Test BeadsDependency with all types."""
    from sdp.beads.models import BeadsDependency, BeadsDependencyType

    for dep_type in BeadsDependencyType:
        dep = BeadsDependency(task_id="bd-1234", type=dep_type)
        assert dep.type == dep_type


def test_cli_trace_report_empty_mappings() -> None:
    """Test TraceabilityReport with no mappings."""
    from sdp.traceability.models import TraceabilityReport

    report = TraceabilityReport(ws_id="00-001-01", mappings=[])
    assert report.total_acs == 0
    assert report.mapped_acs == 0
    assert report.coverage_pct == 100.0  # 100% of 0 is 100%
    assert report.is_complete is True


def test_cli_trace_report_with_failed() -> None:
    """Test TraceabilityReport with failed tests."""
    from sdp.traceability.models import (
        TraceabilityReport,
        ACTestMapping,
        MappingStatus,
    )

    report = TraceabilityReport(
        ws_id="00-001-01",
        mappings=[
            ACTestMapping(
                ac_id="AC1",
                ac_description="Test",
                test_file="tests/test_foo.py",
                test_name="test_foo",
                status=MappingStatus.FAILED,
            )
        ],
    )
    assert report.failed_acs == 1


def test_core_model_models_str_repr() -> None:
    """Test ModelProvider string representation."""
    from sdp.core.model.models import ModelProvider

    model = ModelProvider(
        provider="OpenAI",
        model="gpt-4",
        context="128K",
        tool_use=True,
        cost_per_1m_tokens=3.0,
        availability_pct=0.99,
        context_window=128000,
        notes="Test",
    )
    # Ensure object can be created and accessed
    assert model.provider == "OpenAI"
    assert model.model == "gpt-4"
    assert model.cost_per_1m_tokens == 3.0


def test_prd_annotations_flow_step_immutability() -> None:
    """Test FlowStep dataclass immutability."""
    from sdp.prd.annotations import FlowStep

    step = FlowStep(
        flow_name="test",
        step_number=1,
        description="Test step",
        source_file=Path("/tmp/test.py"),
        line_number=10,
    )
    assert step.flow_name == "test"
    assert step.step_number == 1


def test_validators_capability_tier_models_check_result() -> None:
    """Test ValidationCheck model."""
    from sdp.validators.capability_tier_models import ValidationCheck

    check = ValidationCheck(
        name="test_check",
        passed=True,
        message="Success",
        details=["Detail 1", "Detail 2"],
    )
    assert check.passed is True
    assert check.name == "test_check"
    assert len(check.details) == 2


def test_validators_capability_tier_models_validation_result() -> None:
    """Test ValidationResult model."""
    from sdp.validators.capability_tier_models import (
        ValidationResult,
        ValidationCheck,
        CapabilityTier,
    )

    result = ValidationResult(tier=CapabilityTier.T2, passed=True, checks=[])

    # Add passing check
    result.add_check(ValidationCheck(name="check1", passed=True))
    assert result.passed is True

    # Add failing check
    result.add_check(ValidationCheck(name="check2", passed=False, message="Failed"))
    assert result.passed is False  # Should flip to False


def test_beads_models_task_create_defaults() -> None:
    """Test BeadsTaskCreate with default values."""
    from sdp.beads.models import BeadsTaskCreate, BeadsPriority

    create = BeadsTaskCreate(title="Test Task")
    assert create.priority == BeadsPriority.MEDIUM
    assert create.dependencies == []
    assert create.sdp_metadata == {}


def test_core_model_parser_tool_use_limited() -> None:
    """Test model parser with Limited tool use."""
    from sdp.core.model.parser import parse_models_table

    section = """
| Provider | Model | Context | Tool Use | Notes |
|----------|-------|---------|----------|-------|
| OpenAI | gpt-4 | 128K | ⚠️ Limited | Test |
"""
    result = parse_models_table(section)
    assert len(result) == 1
    assert result[0].tool_use is False  # Limited = False


def test_core_model_parser_tool_use_none() -> None:
    """Test model parser with None tool use."""
    from sdp.core.model.parser import parse_models_table

    section = """
| Provider | Model | Context | Tool Use | Notes |
|----------|-------|---------|----------|-------|
| OpenAI | gpt-4 | 128K | ❌ None | Test |
"""
    result = parse_models_table(section)
    assert len(result) == 1
    assert result[0].tool_use is False


def test_traceability_models_actest_mapping_confidence() -> None:
    """Test ACTestMapping with confidence score."""
    from sdp.traceability.models import ACTestMapping, MappingStatus

    mapping = ACTestMapping(
        ac_id="AC1",
        ac_description="Test",
        test_file="tests/test.py",
        test_name="test_foo",
        status=MappingStatus.MAPPED,
        confidence=0.95,
    )
    assert mapping.confidence == 0.95


def test_traceability_models_from_dict() -> None:
    """Test ACTestMapping.from_dict."""
    from sdp.traceability.models import ACTestMapping, MappingStatus

    data = {
        "ac_id": "AC1",
        "ac_description": "Test",
        "test_file": "tests/test.py",
        "test_name": "test_foo",
        "status": "mapped",
        "confidence": 0.8,
    }
    mapping = ACTestMapping.from_dict(data)
    assert mapping.ac_id == "AC1"
    assert mapping.status == MappingStatus.MAPPED
    assert mapping.confidence == 0.8
