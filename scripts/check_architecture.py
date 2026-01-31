#!/usr/bin/env python3
"""Check Clean Architecture dependency rules.

Validates that domain layer has no external dependencies and that
other layers follow the dependency flow: domain ← core ← (beads, unified, cli).
"""

import ast
import sys
from pathlib import Path
from typing import Any

# Forbidden import rules for each layer
FORBIDDEN_IMPORTS = {
    # Domain layer must have ZERO dependencies on other SDP layers
    "src/sdp/domain/": [
        "sdp.core",
        "sdp.beads",
        "sdp.unified",
        "sdp.cli",
        "sdp.prd",
        "sdp.github",
        "sdp.hooks",
        "sdp.validators",
        "sdp.quality",
        "sdp.tdd",
        "sdp.adapters",
    ],
    # Beads should not import from core (use domain instead)
    "src/sdp/beads/": [
        "sdp.core.workstream.models",  # Use domain.workstream
        "sdp.core.feature.models",  # Use domain.feature
    ],
    # Unified should not import from core models (use domain)
    "src/sdp/unified/": [
        "sdp.core.workstream.models",
        "sdp.core.feature.models",
    ],
}


class ImportVisitor(ast.NodeVisitor):
    """AST visitor to collect imports."""

    def __init__(self) -> None:
        self.imports: list[tuple[str, int]] = []

    def visit_Import(self, node: ast.Import) -> None:
        """Visit import statement."""
        for alias in node.names:
            self.imports.append((alias.name, node.lineno))

    def visit_ImportFrom(self, node: ast.ImportFrom) -> None:
        """Visit from-import statement."""
        if node.module:
            self.imports.append((node.module, node.lineno))


def check_imports(file_path: Path, forbidden: list[str]) -> list[str]:
    """Check file for forbidden imports.

    Args:
        file_path: Python file to check
        forbidden: List of forbidden import prefixes

    Returns:
        List of violation messages
    """
    try:
        content = file_path.read_text()
        tree = ast.parse(content, str(file_path))
    except SyntaxError as e:
        return [f"{file_path}:0: Syntax error: {e}"]

    visitor = ImportVisitor()
    visitor.visit(tree)

    violations = []
    for module, lineno in visitor.imports:
        for forbidden_prefix in forbidden:
            if module == forbidden_prefix or module.startswith(f"{forbidden_prefix}."):
                violations.append(
                    f"{file_path}:{lineno}: Forbidden import '{module}' "
                    f"(violates '{forbidden_prefix}' rule)"
                )

    return violations


def find_python_files(directory: Path) -> list[Path]:
    """Find all Python files in directory.

    Args:
        directory: Directory to search

    Returns:
        List of .py files
    """
    return list(directory.rglob("*.py"))


def main() -> int:
    """Run architecture checks.

    Returns:
        Exit code: 0 if all checks pass, 1 if violations found
    """
    workspace = Path(__file__).parent.parent
    src_dir = workspace / "src" / "sdp"

    if not src_dir.exists():
        print(f"Error: Source directory not found: {src_dir}")
        return 1

    violations_found = False

    # Check each layer
    for layer_path, forbidden_imports in FORBIDDEN_IMPORTS.items():
        layer_dir = workspace / layer_path
        if not layer_dir.exists():
            continue

        print(f"Checking {layer_path}...")
        python_files = find_python_files(layer_dir)

        for py_file in python_files:
            violations = check_imports(py_file, forbidden_imports)
            if violations:
                violations_found = True
                for violation in violations:
                    print(f"  ❌ {violation}")

    if violations_found:
        print("\n❌ Architecture violations found!")
        print("\nClean Architecture rules:")
        print("  • domain/ must not import from any other SDP layer")
        print("  • beads/ should use domain/, not core/workstream/models")
        print("  • unified/ should use domain/, not core/feature/models")
        print("\nSee docs/concepts/clean-architecture/ for details.")
        return 1

    print("\n✅ All architecture checks passed!")
    return 0


if __name__ == "__main__":
    sys.exit(main())
