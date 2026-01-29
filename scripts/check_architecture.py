#!/usr/bin/env python3
"""Check clean architecture layer boundaries.

Usage:
    python scripts/check_architecture.py <file_paths...>

This script replaces hardcoded architecture checks in hooks/pre-commit.sh.
It uses configurable layer patterns from quality-gate.toml.

Examples:
    # Check specific files
    python scripts/check_architecture.py src/domain/entities/user.py

    # Check from pre-commit hook
    python scripts/check_architecture.py --staged

Exit codes:
    0: All checks passed
    1: Architecture violations found
    2: Configuration error
"""

import argparse
import ast
import sys
from pathlib import Path

from sdp.quality.architecture import ArchitectureChecker
from sdp.quality.config import QualityGateConfigLoader
from sdp.quality.validator_models import QualityGateViolation


def check_file(
    file_path: Path,
    config_loader: QualityGateConfigLoader,
) -> list[QualityGateViolation]:
    """Check a single file for architecture violations.

    Args:
        file_path: Path to Python file.
        config_loader: Quality gate configuration loader.

    Returns:
        List of violations (empty if none).
    """
    violations: list[QualityGateViolation] = []

    try:
        content = file_path.read_text(encoding="utf-8")
        tree = ast.parse(content, filename=str(file_path))

        checker = ArchitectureChecker(
            config=config_loader.config.architecture,
            violations=violations,
        )
        checker.check_architecture(file_path, tree)

    except (SyntaxError, OSError) as e:
        violations.append(
            QualityGateViolation(
                category="architecture",
                file_path=str(file_path),
                line_number=None,
                message=f"Error parsing file: {e}",
                severity="error",
            )
        )

    return violations


def main() -> int:
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description="Check clean architecture layer boundaries",
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument(
        "files",
        nargs="*",
        help="Python files to check",
    )
    parser.add_argument(
        "--staged",
        action="store_true",
        help="Check only staged Python files",
    )
    parser.add_argument(
        "--config",
        type=Path,
        default=None,
        help="Path to quality-gate.toml file",
    )

    args = parser.parse_args()

    # Load configuration
    try:
        config_loader = QualityGateConfigLoader(args.config)
    except Exception as e:
        print(f"❌ Configuration error: {e}", file=sys.stderr)
        return 2

    if not config_loader.config.architecture.enabled:
        print("⚠️  Architecture checks disabled in configuration")
        return 0

    # Determine files to check
    if args.staged:
        import subprocess

        try:
            result = subprocess.run(
                ["git", "diff", "--cached", "--name-only", "--diff-filter=ACM"],
                capture_output=True,
                text=True,
                check=True,
            )
            file_paths = [
                Path(f) for f in result.stdout.splitlines() if f.endswith(".py")
            ]
        except (subprocess.CalledProcessError, FileNotFoundError):
            print("⚠️  Not in git repository or git not available")
            file_paths = []
    else:
        file_paths = [Path(f) for f in args.files]

    if not file_paths:
        print("No Python files to check")
        return 0

    # Check all files
    all_violations: list[QualityGateViolation] = []
    for file_path in file_paths:
        if not file_path.exists():
            continue
        violations = check_file(file_path, config_loader)
        all_violations.extend(violations)

    # Report results
    if all_violations:
        print(f"❌ Architecture violations found ({len(all_violations)})")
        print()

        for violation in all_violations:
            print(f"  {violation.file_path}:{violation.line_number or '?'}")
            print(f"    {violation.message}")
            print()

        return 1

    print(f"✓ Architecture checks passed ({len(file_paths)} files)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
