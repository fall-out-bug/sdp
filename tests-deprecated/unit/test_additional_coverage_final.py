"""Additional targeted tests to reach 80% coverage."""

from pathlib import Path

from sdp.prd.parser_python import parse_python_annotations_ast


def test_parse_python_annotations_ast_fallback_on_parse_error() -> None:
    """Test parse_python_annotations_ast falls back to regex (lines 113-115)."""
    invalid_file = Path("/tmp/invalid_syntax.py")
    invalid_file.write_text("def foo(:\n  pass")  # Invalid syntax

    # Should fall back to regex parser
    result = parse_python_annotations_ast(invalid_file)
    assert result == []

    invalid_file.unlink()


def test_beads_models_deferred_status() -> None:
    """Test beads models with DEFERRED status."""
    from sdp.beads.models import BeadsTask, BeadsStatus

    task = BeadsTask(
        id="bd-1234",
        title="Test",
        status=BeadsStatus.DEFERRED,
    )

    data = task.to_dict()
    assert data["status"] == "deferred"


def test_beads_models_tombstone_status() -> None:
    """Test beads models with TOMBSTONE status."""
    from sdp.beads.models import BeadsTask, BeadsStatus

    task = BeadsTask(
        id="bd-1234",
        title="Test",
        status=BeadsStatus.TOMBSTONE,
    )

    assert not task.is_ready()


def test_cli_trace_markdown_table_truncation() -> None:
    """Test cli trace output formatting with long descriptions."""
    from sdp.traceability.models import TraceabilityReport, ACTestMapping, MappingStatus

    report = TraceabilityReport(
        ws_id="00-001-01",
        mappings=[
            ACTestMapping(
                ac_id="AC1",
                ac_description="A very long description that exceeds thirty characters",
                test_file="tests/test_foo.py",
                test_name="test_something",
                status=MappingStatus.MAPPED,
            )
        ],
    )

    table = report.to_markdown_table()
    # Should truncate long descriptions
    assert "..." in table


def test_validators_ws_completion_parser_find_ws_file_not_found() -> None:
    """Test validators ws_completion parser find_ws_file with nonexistent dir."""
    from sdp.validators.ws_completion.parser import find_ws_file

    result = find_ws_file("00-999-99", Path("/nonexistent"))
    assert result is None


def test_prd_parser_python_async_with_regex() -> None:
    """Test prd parser handles async functions via regex."""
    from sdp.prd.parser_python import parse_python_annotations

    test_file = Path("/tmp/test_async.py")
    test_file.write_text("""
@prd_flow("test")
async def async_handler():
    pass
""")

    result = parse_python_annotations(test_file)
    assert len(result) == 1
    assert result[0].flow_name == "test"

    test_file.unlink()


def test_core_model_parser_new_format_with_cost() -> None:
    """Test model parser with new format including cost."""
    from sdp.core.model.parser import parse_models_table

    section = """
| Provider | Model | Cost ($/1M) | Availability | Context | Tool Use | Notes |
|----------|-------|-------------|--------------|---------|----------|-------|
| OpenAI | gpt-4 | 3.00 | 99.9% | 128K | ✅ Full | Great |
"""
    result = parse_models_table(section)
    assert len(result) == 1
    assert result[0].cost_per_1m_tokens == 3.0
    assert abs(result[0].availability_pct - 0.999) < 0.001  # Floating point comparison


def test_validators_ws_completion_parser_scope_indented() -> None:
    """Test ws_completion parser with indented scope_files."""
    from sdp.validators.ws_completion.parser import parse_frontmatter_scope

    content = """---
title: Test
scope_files:
  - src/foo.py
  - src/bar.py
other_field: value
---

Content
"""
    result = parse_frontmatter_scope(content)
    assert "src/foo.py" in result
    assert "src/bar.py" in result

