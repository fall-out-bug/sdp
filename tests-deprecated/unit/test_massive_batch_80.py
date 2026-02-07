"""Massive batch of tiny tests to cover remaining lines."""

from pathlib import Path


def test_prd_parser_python_line_number_calculation() -> None:
    """Test line number calculation in PRD parser."""
    from sdp.prd.parser_python import parse_python_annotations

    test_file = Path("/tmp/test_lines.py")
    test_file.write_text("""# Comment line 1
# Comment line 2
# Comment line 3

@prd_flow("test")
def handler():
    pass
""")

    result = parse_python_annotations(test_file)
    assert len(result) == 1
    assert result[0].line_number > 1  # Should be calculated

    test_file.unlink()


def test_beads_models_status_enum_values() -> None:
    """Test all BeadsStatus enum values."""
    from sdp.beads.models import BeadsStatus

    assert BeadsStatus.OPEN.value == "open"
    assert BeadsStatus.IN_PROGRESS.value == "in_progress"
    assert BeadsStatus.BLOCKED.value == "blocked"
    assert BeadsStatus.DEFERRED.value == "deferred"
    assert BeadsStatus.CLOSED.value == "closed"
    assert BeadsStatus.TOMBSTONE.value == "tombstone"
    assert BeadsStatus.PINNED.value == "pinned"
    assert BeadsStatus.HOOKED.value == "hooked"


def test_beads_models_priority_enum_values() -> None:
    """Test all BeadsPriority enum values."""
    from sdp.beads.models import BeadsPriority

    assert BeadsPriority.CRITICAL.value == 0
    assert BeadsPriority.HIGH.value == 1
    assert BeadsPriority.MEDIUM.value == 2
    assert BeadsPriority.LOW.value == 3
    assert BeadsPriority.BACKLOG.value == 4


def test_beads_models_dependency_type_enum_values() -> None:
    """Test all BeadsDependencyType enum values."""
    from sdp.beads.models import BeadsDependencyType

    assert BeadsDependencyType.BLOCKS.value == "blocks"
    assert BeadsDependencyType.PARENT_CHILD.value == "parent-child"
    assert BeadsDependencyType.RELATED.value == "related"
    assert BeadsDependencyType.DISCOVERED_FROM.value == "discovered-from"


def test_core_model_parser_multiple_models() -> None:
    """Test parsing multiple models from table."""
    from sdp.core.model.parser import parse_models_table

    section = """
| Provider | Model | Context | Tool Use | Notes |
|----------|-------|---------|----------|-------|
| OpenAI | gpt-4 | 128K | ✅ | Model 1 |
| Anthropic | claude-3 | 200K | ✅ | Model 2 |
| Google | gemini | 1M+ | ✅ | Model 3 |
"""
    result = parse_models_table(section)
    assert len(result) == 3
    assert result[0].provider == "OpenAI"
    assert result[1].provider == "Anthropic"
    assert result[2].provider == "Google"


def test_traceability_models_to_dict() -> None:
    """Test TraceabilityReport.to_dict."""
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
                test_file="tests/test.py",
                test_name="test_foo",
                status=MappingStatus.MAPPED,
            )
        ],
    )
    data = report.to_dict()
    assert data["ws_id"] == "00-001-01"
    assert data["total_acs"] == 1
    assert data["mapped_acs"] == 1
    assert len(data["mappings"]) == 1


def test_prd_annotations_flow_step_source_file() -> None:
    """Test FlowStep with source_file tracking."""
    from sdp.prd.annotations import FlowStep

    source = Path("/tmp/test.py")
    step = FlowStep(
        flow_name="test",
        step_number=1,
        description="Test",
        source_file=source,
        line_number=10,
    )
    assert step.source_file == source


def test_validators_capability_tier_t0_t1_interface_checks() -> None:
    """Test T0/T1 tier enum values."""
    from sdp.validators.capability_tier_models import CapabilityTier

    assert CapabilityTier.T0.value == "T0"
    assert CapabilityTier.T1.value == "T1"
    assert CapabilityTier.T2.value == "T2"
    assert CapabilityTier.T3.value == "T3"


def test_beads_models_task_with_parent_id() -> None:
    """Test BeadsTask with parent_id."""
    from sdp.beads.models import BeadsTask

    task = BeadsTask(
        id="bd-child",
        title="Child Task",
        parent_id="bd-parent",
    )
    data = task.to_dict()
    assert data["parent_id"] == "bd-parent"


def test_beads_models_task_with_external_ref() -> None:
    """Test BeadsTask with external_ref."""
    from sdp.beads.models import BeadsTask

    task = BeadsTask(
        id="bd-1234",
        title="Test",
        external_ref="00-001-01",
    )
    data = task.to_dict()
    assert data["external_ref"] == "00-001-01"


def test_beads_models_task_with_issue_type() -> None:
    """Test BeadsTask with issue_type."""
    from sdp.beads.models import BeadsTask

    task = BeadsTask(
        id="bd-1234",
        title="Test",
        issue_type="feature",
    )
    data = task.to_dict()
    assert data["issue_type"] == "feature"


def test_core_model_parser_empty_section() -> None:
    """Test parse_models_table with section but no matching table."""
    from sdp.core.model.parser import parse_models_table

    section = """
# Some content

No table here.
"""
    result = parse_models_table(section)
    assert result == []


def test_validators_ws_completion_parser_multiple_files() -> None:
    """Test parse_frontmatter_scope with multiple files."""
    from sdp.validators.ws_completion.parser import parse_frontmatter_scope

    content = """---
scope_files:
  - src/file1.py
  - src/file2.py
  - src/file3.py
  - tests/test_file1.py
---
"""
    result = parse_frontmatter_scope(content)
    assert len(result) == 4
    assert "src/file1.py" in result
    assert "tests/test_file1.py" in result


def test_validators_ws_completion_parser_multiple_commands() -> None:
    """Test parse_verification_commands with multiple commands."""
    from sdp.validators.ws_completion.parser import parse_verification_commands

    content = """
### Verification

```bash
pytest tests/unit/
pytest tests/integration/
mypy src/
black --check src/
ruff check src/
```
"""
    result = parse_verification_commands(content)
    assert len(result) == 5
    assert "pytest tests/unit/" in result
    assert "ruff check src/" in result


def test_traceability_models_mapping_status_enum() -> None:
    """Test MappingStatus enum values."""
    from sdp.traceability.models import MappingStatus

    assert MappingStatus.MAPPED.value == "mapped"
    assert MappingStatus.MISSING.value == "missing"
    assert MappingStatus.FAILED.value == "failed"


def test_prd_parser_python_directory_with_valid_files() -> None:
    """Test parse_directory with actual Python files."""
    from sdp.prd.parser_python import parse_directory
    from pathlib import Path

    test_dir = Path("/tmp/test_prd_dir")
    test_dir.mkdir(exist_ok=True)

    test_file = test_dir / "module.py"
    test_file.write_text("""
@prd_flow("test")
def handler():
    pass
""")

    result = parse_directory(test_dir)
    assert len(result) >= 1

    test_file.unlink()
    test_dir.rmdir()


def test_beads_models_task_from_dict_with_all_fields() -> None:
    """Test BeadsTask.from_dict with comprehensive data."""
    from sdp.beads.models import BeadsTask

    data = {
        "id": "bd-1234",
        "title": "Full Task",
        "description": "Full description",
        "status": "in_progress",
        "priority": 1,
        "issue_type": "feature",
        "parent_id": "bd-parent",
        "dependencies": [{"task_id": "bd-dep", "type": "blocks"}],
        "external_ref": "00-001-01",
        "created_at": "2024-01-01T12:00:00",
        "updated_at": "2024-01-01T13:00:00",
        "metadata": {"sdp": {"ws_id": "00-001-01"}},
    }
    task = BeadsTask.from_dict(data)
    assert task.id == "bd-1234"
    assert task.title == "Full Task"
    assert task.description == "Full description"
    assert len(task.dependencies) == 1
    assert task.sdp_metadata == {"ws_id": "00-001-01"}


def test_core_model_parser_context_window_1m_plus() -> None:
    """Test parse_context_window with 1M+ format."""
    from sdp.core.model.parser import parse_context_window

    assert parse_context_window("1M+") == 1_000_000
    assert parse_context_window("2M+") == 2_000_000


def test_core_model_parser_context_window_fractional() -> None:
    """Test parse_context_window with fractional values."""
    from sdp.core.model.parser import parse_context_window

    assert parse_context_window("0.5M") == 500_000
    assert parse_context_window("1.5M") == 1_500_000
