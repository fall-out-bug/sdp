"""Command-line interface for quality gate validation."""

import argparse
import subprocess
import sys
from pathlib import Path

from scripts.quality.checker import QualityGateChecker


def parse_staged_files() -> list[Path]:
    """Parse staged files from git diff --cached.

    Returns:
        List of staged Python files in src/
    """
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
    """Run quality gate validation.

    Returns:
        Exit code (0 = success, 1 = failure)
    """
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
    repo_root = Path(__file__).parent.parent.parent

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
