"""Performance checks for quality gates.

Checks for:
- Excessive nesting depth (>5)
"""

import ast
from pathlib import Path

from scripts.quality.models import Violation


class PerformanceChecker:
    """Performance anti-pattern checker."""

    def __init__(self, max_nesting_depth: int = 5) -> None:
        """Initialize performance checker.

        Args:
            max_nesting_depth: Maximum allowed nesting depth
        """
        self.max_nesting_depth = max_nesting_depth
        self._violations: list[Violation] = []

    def check(self, path: Path, tree: ast.AST) -> list[Violation]:
        """Run performance checks.

        Args:
            path: File path
            tree: AST tree

        Returns:
            List of violations
        """
        self._violations.clear()

        if self.max_nesting_depth:
            self._check_nesting_depth(path, tree)

        return self._violations

    def _check_nesting_depth(self, path: Path, tree: ast.AST) -> None:
        """Check for excessive nesting depth."""
        for node in ast.walk(tree):
            if isinstance(node, (ast.FunctionDef, ast.AsyncFunctionDef)):
                depth = self._calculate_nesting_depth(node)
                if depth > self.max_nesting_depth:
                    self._violations.append(
                        Violation(
                            "performance",
                            str(path),
                            node.lineno,
                            f"Function '{node.name}' has nesting depth {depth} "
                            f"(max: {self.max_nesting_depth})",
                            "warning",
                        )
                    )

    def _calculate_nesting_depth(self, node: ast.AST) -> int:
        """Calculate maximum nesting depth in a function.

        Uses recursive approach with return value instead of closure.
        """
        return self._depth_at(node, 0)

    def _depth_at(self, child_node: ast.AST, current_depth: int) -> int:
        """Calculate nesting depth at a specific node.

        Args:
            child_node: AST node to check
            current_depth: Current nesting depth

        Returns:
            Maximum depth found at this node
        """
        child_depth = current_depth

        for grandchild in ast.iter_child_nodes(child_node):
            if isinstance(
                grandchild,
                (ast.If, ast.While, ast.For, ast.AsyncFor, ast.With, ast.Try),
            ):
                child_depth = max(child_depth, self._depth_at(grandchild, current_depth + 1))

        return child_depth
