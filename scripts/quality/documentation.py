"""Documentation checks for quality gates.

Checks for:
- Missing module docstrings
"""

import ast
from pathlib import Path
from typing import cast

from scripts.quality.models import Violation


class DocumentationChecker:
    """Documentation requirements checker."""

    def __init__(self, require_module_docstrings: bool = True) -> None:
        """Initialize documentation checker.

        Args:
            require_module_docstrings: Require module-level docstrings
        """
        self.require_module_docstrings = require_module_docstrings
        self._violations: list[Violation] = []

    def check(self, path: Path, tree: ast.AST) -> list[Violation]:
        """Run documentation checks.

        Args:
            path: File path
            tree: AST tree

        Returns:
            List of violations
        """
        self._violations.clear()

        if self.require_module_docstrings:
            self._check_module_docstring(path, tree)

        return self._violations

    def _check_module_docstring(self, path: Path, tree: ast.AST) -> None:
        """Check for module docstring."""
        has_docstring = ast.get_docstring(cast(ast.Module, tree)) is not None
        if not has_docstring:
            self._violations.append(
                Violation("documentation", str(path), 1, "Module missing docstring", "warning")
            )
