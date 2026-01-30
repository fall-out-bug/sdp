#!/usr/bin/env python3
"""Quality gate validation script for git hooks.

Runs security, documentation, and performance checks on Python files.
"""

import argparse
import ast
import re
import subprocess
import sys
from collections import namedtuple
from pathlib import Path
from typing import cast

Violation = namedtuple("Violation", ["category", "file_path", "line_no", "message", "severity"])


class QualityGateChecker:
    """Minimal quality gate checker for git hooks."""

    def __init__(self, repo_root: Path) -> None:
        """Initialize checker."""
        self.repo_root = repo_root
        self.config_file = repo_root / "quality-gate.toml"
        self._violations: list[Violation] = []

        # Load config (simple TOML parsing)
        self.security_enabled = True
        self.forbid_hardcoded_secrets = True
        self.forbid_eval_usage = True
        self.require_module_docstrings = True
        self.max_nesting_depth = 5

    def validate_file(self, file_path: Path) -> list[Violation]:
        """Validate a single Python file."""
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
        """Run enabled quality checks."""
        # Security checks
        if self.security_enabled:
            self._check_security(path, source_code)

        # Documentation checks
        if self.require_module_docstrings:
            self._check_documentation(path, tree)

        # Performance checks
        if self.max_nesting_depth:
            self._check_performance(path, tree)

    def _check_security(self, path: Path, source_code: str) -> None:
        """Check for security issues."""
        if self.forbid_hardcoded_secrets:
            secret_patterns = [
                r'(?:password|passwd|pwd)\s*=\s*["\']([^"\']{8,})["\']',
                r'(?:api_key|apikey|api-key)\s*=\s*["\']([^"\']{8,})["\']',
                r'(?:secret|secret_key|secret-key)\s*=\s*["\']([^"\']{8,})["\']',
                r'(?:token|auth_token|auth-token)\s*=\s*["\']([^"\']{8,})["\']',
                r'(?:private_key|private-key|privatekey)\s*=\s*["\']([^"\']{8,})["\']',
            ]

            for pattern in secret_patterns:
                matches = re.finditer(pattern, source_code, re.IGNORECASE)
                for match in matches:
                    # Exclude obvious test/example values
                    value = match.group(1)
                    if not re.search(r'^(test|example|mock|dummy|xxx|xxx+|\*+)$', value, re.IGNORECASE):
                        line_num = source_code[: match.start()].count("\n") + 1
                        self._violations.append(
                            Violation(
                                "security",
                                str(path),
                                line_num,
                                f"Possible hardcoded secret: {match.group(1)[:10]}...",
                                "error",
                            )
                        )

        if self.forbid_eval_usage:
            if "eval(" in source_code:
                line_num = source_code.index("eval(")
                line_num = source_code[:line_num].count("\n") + 1
                self._violations.append(
                    Violation(
                        "security",
                        str(path),
                        line_num,
                        "Use of eval() detected (security risk)",
                        "error",
                    )
                )

    def _check_documentation(self, path: Path, tree: ast.AST) -> None:
        """Check documentation requirements."""
        has_docstring = ast.get_docstring(cast(ast.Module, tree)) is not None
        if not has_docstring:
            self._violations.append(
                Violation("documentation", str(path), 1, "Module missing docstring", "warning")
            )

    def _check_performance(self, path: Path, tree: ast.AST) -> None:
        """Check performance anti-patterns."""
        if not self.max_nesting_depth:
            return

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
        """Calculate maximum nesting depth in a function."""
        max_depth = 0

        def _depth_at(child_node: ast.AST, current_depth: int) -> None:
            nonlocal max_depth
            max_depth = max(max_depth, current_depth)
            for grandchild in ast.iter_child_nodes(child_node):
                if isinstance(
                    grandchild,
                    (ast.If, ast.While, ast.For, ast.AsyncFor, ast.With, ast.Try),
                ):
                    _depth_at(grandchild, current_depth + 1)

        _depth_at(node, 0)
        return max_depth


def parse_staged_files() -> list[Path]:
    """Parse staged files from git diff --cached."""
    result = subprocess.run(
        ["git", "diff", "--cached", "--name-only", "--diff-filter=ACM"],
        capture_output=True,
        text=True,
        check=False,
    )

    if result.returncode != 0:
        return []

    files = []
    for line in result.stdout.strip().split("\n"):
        if line:
            file_path = Path(line)
            # Only validate Python files in src/ directory
            if file_path.suffix == ".py" and "src/" in str(file_path):
                files.append(file_path)

    return files


def main() -> int:
    """Run quality gate validation."""
    parser = argparse.ArgumentParser(
        description="Validate Python files against quality gates"
    )
    parser.add_argument(
        "--staged",
        action="store_true",
        help="Check only staged files (git diff --cached)",
    )
    parser.add_argument(
        "files",
        nargs="*",
        type=Path,
        help="Specific files to validate (optional)",
    )

    args = parser.parse_args()

    # Get repository root
    repo_root = Path(__file__).parent.parent

    # Determine which files to validate
    if args.staged:
        files = parse_staged_files()
    elif args.files:
        files = [f for f in args.files if f.suffix == ".py"]
    else:
        print("❌ No files specified. Use --staged or provide file paths.")
        return 1

    if not files:
        print("  No Python files to validate")
        return 0

    # Initialize checker
    checker = QualityGateChecker(repo_root)

    # Validate all files
    violations = []
    for file_path in files:
        violations.extend(checker.validate_file(file_path))

    # Group violations by severity
    errors = [v for v in violations if v.severity == "error"]
    warnings = [v for v in violations if v.severity == "warning"]

    # Report results
    if violations:
        print(f"\n{'='*60}")
        print("Quality Gate Validation Report")
        print(f"{'='*60}")
        print(f"Files checked: {len(files)}")
        print(f"Total violations: {len(violations)}")
        print(f"  Errors: {len(errors)}")
        print(f"  Warnings: {len(warnings)}")

        # Group by category
        by_category: dict[str, list] = {}
        for v in violations:
            if v.category not in by_category:
                by_category[v.category] = []
            by_category[v.category].append(v)

        if by_category:
            print("\nViolations by category:")
            for category in sorted(by_category.keys()):
                print(f"  {category}: {len(by_category[category])}")

        # Print detailed violations
        print(f"\n{'='*60}")
        print("Detailed violations:")
        print(f"{'='*60}")
        for v in violations:
            print(f"{v.file_path}:{v.line_no}: [{v.category}] {v.message}")
        print(f"{'='*60}\n")

        # Exit with error if any errors found
        if errors:
            print("❌ Quality gate validation FAILED (errors found)")
            return 1
        elif warnings:
            print("⚠️ Quality gate validation passed with warnings")
            return 0
    else:
        print("✓ Quality gate validation passed")

    return 0


if __name__ == "__main__":
    sys.exit(main())
