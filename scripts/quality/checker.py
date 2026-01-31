"""Quality gate validation script for git hooks.

Main checker that orchestrates security, documentation, and performance checks.
"""

import ast
from pathlib import Path
from typing import cast

from scripts.quality.documentation import DocumentationChecker
from scripts.quality.models import Violation
from scripts.quality.performance import PerformanceChecker
from scripts.quality.security import SecurityChecker


class QualityGateChecker:
    """Minimal quality gate checker for git hooks."""

    def __init__(self, repo_root: Path) -> None:
        """Initialize checker.

        Args:
            repo_root: Repository root path
        """
        self.repo_root = repo_root
        self.config_file = repo_root / "quality-gate.toml"
        self._violations: list[Violation] = []

        # Load config (simple TOML parsing)
        self.security_enabled = True
        self.require_module_docstrings = True
        self.max_nesting_depth = 5

        # Initialize specialized checkers
        self.security_checker = SecurityChecker(
            forbid_hardcoded_secrets=True,
            forbid_eval_usage=True,
        )
        self.documentation_checker = DocumentationChecker(
            require_module_docstrings=True,
        )
        self.performance_checker = PerformanceChecker(
            max_nesting_depth=5,
        )

    def validate_file(self, file_path: Path) -> list[Violation]:
        """Validate a single Python file.

        Args:
            file_path: Path to file

        Returns:
            List of violations
        """
        self._violations.clear()

        if not file_path.exists():
            self._violations.append(
                Violation("file_not_found", str(file_path), None, "File not found", "error")
            )
            return self._violations

        if file_path.suffix != ".py":
            return self._violations

        try:
            source_code = file_path.read_text()
            tree = ast.parse(source_code, filename=str(file_path))
            self._run_checks(file_path, source_code, tree)
        except SyntaxError as e:
            self._violations.append(
                Violation(
                    "syntax_error",
                    str(file_path),
                    e.lineno,
                    f"Syntax error: {e.msg}",
                    "error",
                )
            )

        return self._violations

    def _run_checks(self, path: Path, source_code: str, tree: ast.AST) -> None:
        """Run enabled quality checks.

        Args:
            path: File path
            source_code: Source code string
            tree: AST tree
        """
        # Security checks
        if self.security_enabled:
            self._violations.extend(self.security_checker.check(path, source_code, tree))

        # Documentation checks
        if self.require_module_docstrings:
            self._violations.extend(self.documentation_checker.check(path, tree))

        # Performance checks
        if self.max_nesting_depth:
            self._violations.extend(self.performance_checker.check(path, tree))
