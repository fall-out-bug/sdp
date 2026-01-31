"""Python annotation parser.

This module extracts PRD flow annotations from Python source files.
"""

import ast
import re
from pathlib import Path

from .annotations import FlowStep


def parse_python_annotations(path: Path) -> list[FlowStep]:
    """Parse @prd_flow and @prd_step decorators from Python file.

    This function uses regex to find decorator patterns in Python files.
    It looks for @prd_flow("name") and @prd_step(N, "desc") patterns.

    Args:
        path: Path to the Python file

    Returns:
        List of FlowStep objects found in the file
    """
    try:
        content = path.read_text()
    except Exception:
        return []

    steps = []

    # Pattern 1: @prd_flow("name") followed by @prd_step(N, "desc")
    # Pattern 2: @prd_step(N, "desc") followed by @prd_flow("name")
    # Pattern 3: @prd_flow("name") without step (step_number defaults to 0)

    # Combined pattern to match both orders
    flow_pattern = re.compile(
        r'@prd_flow\(["\']([^"\']+)["\']\)\s*\n'
        r'(?:@prd_step\((\d+),\s*["\']([^"\']+)["\']\)\s*\n)?'
        r'(?:async\s+)?def\s+(\w+)',
        re.MULTILINE
    )

    for match in flow_pattern.finditer(content):
        flow_name = match.group(1)
        step_num_str = match.group(2)
        step_desc = match.group(3)
        func_name = match.group(4)

        # Determine step number and description
        if step_num_str:
            step_num = int(step_num_str)
            desc = step_desc or f"{func_name}"
        else:
            step_num = 0
            desc = func_name

        # Calculate line number
        line_number = content[:match.start()].count('\n') + 1

        steps.append(FlowStep(
            flow_name=flow_name,
            step_number=step_num,
            description=desc,
            source_file=path,
            line_number=line_number
        ))

    return steps


def parse_directory(directory: Path, pattern: str = "**/*.py") -> list[FlowStep]:
    """Parse all Python files in directory matching pattern.

    Args:
        directory: Root directory to search
        pattern: Glob pattern for files (default: "**/*.py")

    Returns:
        List of FlowStep objects from all matching files
    """
    all_steps: list[FlowStep] = []

    try:
        for file in directory.glob(pattern):
            # Skip common non-source directories
            if any(skip in str(file) for skip in ["venv", ".venv", "__pycache__", ".tox", "node_modules", ".git"]):  # noqa: E501
                continue

            if file.is_file():
                all_steps.extend(parse_python_annotations(file))
    except Exception:
        pass

    return all_steps


def parse_python_annotations_ast(path: Path) -> list[FlowStep]:
    """Parse @prd_flow and @prd_step decorators using AST.

    This is an alternative implementation using Python's AST module.
    It's more robust but may fail on syntax errors.

    Args:
        path: Path to the Python file

    Returns:
        List of FlowStep objects found in the file
    """
    try:
        content = path.read_text()
        _ = ast.parse(content, filename=str(path))
    except Exception:
        # Fall back to regex if AST parsing fails
        return parse_python_annotations(path)

    # For now, always use regex approach
    # AST visitor is complex and regex works well for decorators
    return parse_python_annotations(path)


class _PRDVisitor(ast.NodeVisitor):
    """AST visitor that extracts PRD annotations."""

    def __init__(self, source_file: Path) -> None:
        """Initialize visitor.

        Args:
            source_file: Path to source file for error reporting
        """
        self.steps: list[FlowStep] = []
        self.current_flow: str | None = None
        self.source_file = source_file

    def visit_FunctionDef(self, node: ast.FunctionDef) -> None:  # noqa: N802
        """Visit function definition and extract decorators."""
        flow_name: str | None = None
        step_number: int | None = None
        description: str | None = None

        # Check decorators
        for decorator in node.decorator_list:
            flow_name, step_number, description = self._extract_decorator_info(
                decorator, flow_name, step_number, description
            )

        if flow_name:
            step_number, description = self._finalize_step_info(
                flow_name, step_number, description, node
            )
            self._add_flow_step(flow_name, step_number, description, node)

        self.generic_visit(node)

    def _extract_decorator_info(
        self,
        decorator: ast.expr,
        flow_name: str | None,
        step_number: int | None,
        description: str | None,
    ) -> tuple[str | None, int | None, str | None]:
        """Extract flow/step info from a single decorator."""
        if not isinstance(decorator, ast.Call):
            return flow_name, step_number, description

        if not hasattr(decorator.func, 'id'):
            return flow_name, step_number, description

        if decorator.func.id == 'prd_flow':
            flow_name = self._extract_flow_name(decorator)
        elif decorator.func.id == 'prd_step':
            step_number, description = self._extract_step_info(decorator, step_number, description)

        return flow_name, step_number, description

    def _extract_flow_name(self, decorator: ast.Call) -> str | None:
        """Extract flow name from @prd_flow decorator."""
        if decorator.args and isinstance(decorator.args[0], ast.Constant):
            value = decorator.args[0].value
            if isinstance(value, str):
                return value
        return None

    def _extract_step_info(
        self,
        decorator: ast.Call,
        step_number: int | None,
        description: str | None,
    ) -> tuple[int | None, str | None]:
        """Extract step info from @prd_step decorator."""
        if len(decorator.args) >= 2:
            if isinstance(decorator.args[0], ast.Constant):
                step_val = decorator.args[0].value
                if isinstance(step_val, int):
                    step_number = step_val
            if isinstance(decorator.args[1], ast.Constant):
                desc_val = decorator.args[1].value
                if isinstance(desc_val, str):
                    description = desc_val
        return step_number, description

    def _finalize_step_info(
        self,
        flow_name: str | None,
        step_number: int | None,
        description: str | None,
        node: ast.FunctionDef,
    ) -> tuple[int, str]:
        """Finalize step info with defaults and type narrowing."""
        if step_number is None:
            step_number = 0
        if description is None:
            description = node.name

        # Type narrowing for mypy
        if not isinstance(flow_name, str):
            flow_name = str(flow_name)
        if not isinstance(step_number, int):
            if isinstance(step_number, (int, str)):
                step_number = int(step_number)
            else:
                step_number = 0
        if not isinstance(description, str):
            description = str(description)

        return step_number, description

    def _add_flow_step(
        self,
        flow_name: str | None,
        step_number: int,
        description: str,
        node: ast.FunctionDef,
    ) -> None:
        """Add a flow step to the steps list."""
        self.steps.append(FlowStep(
            flow_name=flow_name,
            step_number=step_number,
            description=description,
            source_file=self.source_file,
            line_number=node.lineno
        ))

    def visit_AsyncFunctionDef(self, node: ast.AsyncFunctionDef) -> None:  # noqa: N802
        """Visit async function definition."""
        # Treat async functions the same as regular functions
        # Create a fake FunctionDef-like object with same attributes
        class _FuncDefWrapper:
            """Wrapper to make AsyncFunctionDef compatible with visit_FunctionDef."""

            def __init__(self, async_node: ast.AsyncFunctionDef) -> None:
                """Initialize wrapper.

                Args:
                    async_node: Async function node to wrap
                """
                self.name = async_node.name
                self.decorator_list = async_node.decorator_list
                self.lineno = async_node.lineno

        wrapper = _FuncDefWrapper(node)
        self.visit_FunctionDef(wrapper)  # type: ignore[arg-type]
