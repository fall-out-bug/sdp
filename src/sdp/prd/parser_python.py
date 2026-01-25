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
    all_steps = []

    try:
        for file in directory.glob(pattern):
            # Skip common non-source directories
            if any(skip in str(file) for skip in ["venv", ".venv", "__pycache__", ".tox", "node_modules", ".git"]):
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
        tree = ast.parse(content, filename=str(path))
    except Exception:
        # Fall back to regex if AST parsing fails
        return parse_python_annotations(path)

    steps = []

    class PRDVisitor(ast.NodeVisitor):
        """AST visitor that extracts PRD annotations."""

        def visit_FunctionDef(self, node: ast.FunctionDef) -> None:
            """Visit function definition and extract decorators."""
            flow_name = None
            step_number = None
            description = None

            # Check decorators
            for decorator in node.decorator_list:
                # @prd_flow("name")
                if isinstance(decorator, ast.Call):
                    if hasattr(decorator.func, 'id') and decorator.func.id == 'prd_flow':
                        if decorator.args and isinstance(decorator.args[0], ast.Constant):
                            flow_name = decorator.args[0].value

                    # @prd_step(N, "desc")
                    elif hasattr(decorator.func, 'id') and decorator.func.id == 'prd_step':
                        if len(decorator.args) >= 2:
                            if isinstance(decorator.args[0], ast.Constant):
                                step_number = decorator.args[0].value
                            if isinstance(decorator.args[1], ast.Constant):
                                description = decorator.args[1].value

            if flow_name:
                if step_number is None:
                    step_number = 0
                if description is None:
                    description = node.name

                steps.append(FlowStep(
                    flow_name=flow_name,
                    step_number=step_number,
                    description=description,
                    source_file=path,
                    line_number=node.lineno
                ))

            self.generic_visit(node)

        def visit_AsyncFunctionDef(self, node: ast.AsyncFunctionDef) -> None:
            """Visit async function definition."""
            # Treat async functions the same as regular functions
            self.visit_FunctionDef(node)

    visitor = PRDVisitor()
    visitor.visit(tree)

    return steps
