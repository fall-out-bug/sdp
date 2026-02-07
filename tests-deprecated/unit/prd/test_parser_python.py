"""Tests for Python PRD annotation parser."""

from pathlib import Path

import pytest

from sdp.prd.annotations import FlowStep
from sdp.prd.parser_python import (
    parse_directory,
    parse_python_annotations,
    parse_python_annotations_ast,
)


class TestParsePythonAnnotations:
    """Test regex-based Python parser."""

    def test_parse_empty_file(self, tmp_path: Path) -> None:
        """Test parsing empty file."""
        test_file = tmp_path / "test.py"
        test_file.write_text("")

        steps = parse_python_annotations(test_file)
        assert steps == []

    def test_parse_file_without_decorators(self, tmp_path: Path) -> None:
        """Test parsing file without PRD decorators."""
        test_file = tmp_path / "test.py"
        test_file.write_text("""
def ordinary_function():
    pass

class OrdinaryClass:
    pass
""")

        steps = parse_python_annotations(test_file)
        assert steps == []

    def test_parse_flow_without_step(self, tmp_path: Path) -> None:
        """Test parsing @prd_flow without @prd_step."""
        test_file = tmp_path / "test.py"
        test_file.write_text("""
@prd_flow("test_flow")
def my_function():
    pass
""")

        steps = parse_python_annotations(test_file)
        assert len(steps) == 1
        assert steps[0].flow_name == "test_flow"
        assert steps[0].step_number == 0
        assert steps[0].description == "my_function"
        assert steps[0].source_file == test_file
        assert steps[0].line_number == 2

    def test_parse_flow_with_step(self, tmp_path: Path) -> None:
        """Test parsing @prd_flow with @prd_step."""
        test_file = tmp_path / "test.py"
        test_file.write_text("""
@prd_flow("test_flow")
@prd_step(1, "First step")
def my_function():
    pass
""")

        steps = parse_python_annotations(test_file)
        assert len(steps) == 1
        assert steps[0].flow_name == "test_flow"
        assert steps[0].step_number == 1
        assert steps[0].description == "First step"

    def test_parse_multiple_flows(self, tmp_path: Path) -> None:
        """Test parsing multiple flows in one file."""
        test_file = tmp_path / "test.py"
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

        steps = parse_python_annotations(test_file)
        assert len(steps) == 2
        assert steps[0].flow_name == "flow1"
        assert steps[1].flow_name == "flow2"

    def test_parse_async_function(self, tmp_path: Path) -> None:
        """Test parsing async function with decorators."""
        test_file = tmp_path / "test.py"
        test_file.write_text("""
@prd_flow("async_flow")
@prd_step(1, "Async step")
async def my_async_function():
    pass
""")

        steps = parse_python_annotations(test_file)
        assert len(steps) == 1
        assert steps[0].flow_name == "async_flow"
        assert steps[0].description == "Async step"

    def test_parse_handles_missing_file(self) -> None:
        """Test handling of missing file."""
        steps = parse_python_annotations(Path("/nonexistent/file.py"))
        assert steps == []

    def test_line_number_calculation(self, tmp_path: Path) -> None:
        """Test correct line number calculation."""
        test_file = tmp_path / "test.py"
        test_file.write_text("""
# Line 1
# Line 2
# Line 3
@prd_flow("test")
def func():
    pass
""")

        steps = parse_python_annotations(test_file)
        assert len(steps) == 1
        # @prd_flow is on line 5 (after the initial blank line)
        assert steps[0].line_number == 5


class TestParseDirectory:
    """Test directory parsing."""

    def test_parse_empty_directory(self, tmp_path: Path) -> None:
        """Test parsing empty directory."""
        steps = parse_directory(tmp_path)
        assert steps == []

    def test_parse_directory_with_files(self, tmp_path: Path) -> None:
        """Test parsing directory with multiple files."""
        (tmp_path / "test1.py").write_text("""
@prd_flow("flow1")
def func1():
    pass
""")

        (tmp_path / "test2.py").write_text("""
@prd_flow("flow2")
def func2():
    pass
""")

        steps = parse_directory(tmp_path)
        assert len(steps) == 2
        flow_names = {s.flow_name for s in steps}
        assert flow_names == {"flow1", "flow2"}

    def test_parse_directory_skips_venv(self, tmp_path: Path) -> None:
        """Test that venv directories are skipped."""
        venv_dir = tmp_path / "venv"
        venv_dir.mkdir()
        (venv_dir / "test.py").write_text("""
@prd_flow("should_be_skipped")
def func():
    pass
""")

        (tmp_path / "test.py").write_text("""
@prd_flow("should_be_found")
def func():
    pass
""")

        steps = parse_directory(tmp_path)
        assert len(steps) == 1
        assert steps[0].flow_name == "should_be_found"

    def test_parse_directory_with_pattern(self, tmp_path: Path) -> None:
        """Test parsing with custom pattern."""
        (tmp_path / "test.py").write_text("""
@prd_flow("flow1")
def func1():
    pass
""")

        (tmp_path / "test.txt").write_text("Not a Python file")

        steps = parse_directory(tmp_path, pattern="*.py")
        assert len(steps) == 1


class TestParsePythonAnnotationsAST:
    """Test AST-based parser."""

    def test_ast_falls_back_to_regex(self, tmp_path: Path) -> None:
        """Test that AST parser falls back to regex on syntax error."""
        test_file = tmp_path / "test.py"
        # Write valid Python file
        test_file.write_text("""
@prd_flow("test_flow")
def my_function():
    pass
""")

        steps = parse_python_annotations_ast(test_file)
        assert len(steps) == 1
        assert steps[0].flow_name == "test_flow"

    def test_ast_handles_missing_file(self) -> None:
        """Test handling of missing file."""
        steps = parse_python_annotations_ast(Path("/nonexistent/file.py"))
        assert steps == []
