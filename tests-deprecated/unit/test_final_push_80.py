"""Final push tests for 80% coverage - targeting high-value modules."""

from pathlib import Path
from datetime import datetime

from sdp.beads.models import BeadsTask, BeadsStatus, BeadsPriority


def test_beads_task_to_jsonl() -> None:
    """Test BeadsTask.to_jsonl serialization."""
    task = BeadsTask(
        id="bd-1234",
        title="Test",
        description="Test description",
        status=BeadsStatus.OPEN,
        priority=BeadsPriority.HIGH,
    )
    jsonl = task.to_jsonl()
    assert "bd-1234" in jsonl
    assert "Test" in jsonl


def test_beads_task_from_jsonl() -> None:
    """Test BeadsTask.from_jsonl deserialization."""
    jsonl = '{"id": "bd-1234", "title": "Test", "status": "open", "priority": 1}'
    task = BeadsTask.from_jsonl(jsonl)
    assert task.id == "bd-1234"
    assert task.title == "Test"
    assert task.status == BeadsStatus.OPEN
    assert task.priority == BeadsPriority.HIGH


def test_beads_task_in_progress_is_ready() -> None:
    """Test BeadsTask.is_ready for IN_PROGRESS status."""
    task = BeadsTask(
        id="bd-1234",
        title="Test",
        status=BeadsStatus.IN_PROGRESS,
    )
    assert task.is_ready()


def test_core_model_parser_tool_use_variations() -> None:
    """Test model parser handles different tool_use formats."""
    from sdp.core.model.parser import parse_models_table

    section = """
| Provider | Model | Context | Tool Use | Notes |
|----------|-------|---------|----------|-------|
| OpenAI | gpt-4 | 128K | Full | Great |
"""
    result = parse_models_table(section)
    assert len(result) == 1
    assert result[0].tool_use is True  # "Full" should be detected


def test_prd_parser_python_step_without_description() -> None:
    """Test PRD parser with step and proper description."""
    from sdp.prd.parser_python import parse_python_annotations

    test_file = Path("/tmp/test_prd.py")
    test_file.write_text("""
@prd_flow("test_flow")
@prd_step(1, "Do something")
def handler():
    pass
""")

    result = parse_python_annotations(test_file)
    assert len(result) == 1
    assert result[0].step_number == 1
    assert result[0].description == "Do something"

    test_file.unlink()


def test_validators_ws_completion_parser_verification_no_section() -> None:
    """Test verification parser with no ### Verification section."""
    from sdp.validators.ws_completion.parser import parse_verification_commands

    content = """
# Some content

No verification section here.
"""
    result = parse_verification_commands(content)
    assert result == []


def test_validators_ws_completion_parser_scope_ends_on_non_indent() -> None:
    """Test scope parser stops on non-indented line."""
    from sdp.validators.ws_completion.parser import parse_frontmatter_scope

    content = """---
scope_files:
  - src/foo.py
feature: F01
---
"""
    result = parse_frontmatter_scope(content)
    assert len(result) == 1
    assert "feature: F01" not in result


def test_cli_trace_markdown_table_no_test_name() -> None:
    """Test trace markdown table with missing test name."""
    from sdp.traceability.models import TraceabilityReport, ACTestMapping, MappingStatus

    report = TraceabilityReport(
        ws_id="00-001-01",
        mappings=[
            ACTestMapping(
                ac_id="AC1",
                ac_description="Short",
                test_file=None,
                test_name=None,
                status=MappingStatus.MISSING,
            )
        ],
    )

    table = report.to_markdown_table()
    assert "AC1" in table
    assert "-" in table  # Should show "-" for missing test


def test_cli_doctor_check_environment_success() -> None:
    """Test check_environment with all tools available."""
    import sys
    from unittest.mock import MagicMock, patch
    from collections import namedtuple

    VersionInfo = namedtuple("VersionInfo", ["major", "minor", "micro"])
    good_version = VersionInfo(3, 10, 0)

    with patch.object(sys, "version_info", good_version), patch(
        "subprocess.run"
    ) as mock_run:
        mock_run.side_effect = [
            MagicMock(returncode=0, stdout="Poetry version 1.0", stderr=""),
            MagicMock(returncode=0, stdout="git version 2.0", stderr=""),
        ]

        from sdp.cli.doctor import check_environment

        result = check_environment()
        assert result is True


def test_cli_doctor_check_project_structure_all_present(tmp_path: Path) -> None:
    """Test check_project_structure with all required items."""
    from sdp.cli.doctor import check_project_structure

    # Create required dirs
    (tmp_path / "docs").mkdir()
    (tmp_path / "src").mkdir()
    (tmp_path / "tests").mkdir()

    # Create required files
    (tmp_path / "README.md").write_text("# Test")
    (tmp_path / "pyproject.toml").write_text("[tool.poetry]")

    result = check_project_structure(tmp_path)
    assert result is True


def test_beads_models_pinned_status() -> None:
    """Test beads models with PINNED status."""
    task = BeadsTask(
        id="bd-1234",
        title="Test",
        status=BeadsStatus.PINNED,
    )

    assert not task.is_ready()


def test_beads_models_hooked_status() -> None:
    """Test beads models with HOOKED status."""
    task = BeadsTask(
        id="bd-1234",
        title="Test",
        status=BeadsStatus.HOOKED,
    )

    assert not task.is_ready()
