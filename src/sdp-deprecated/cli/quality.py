"""Quality gate CLI commands."""

import sys
from pathlib import Path

import click


@click.group()
def quality() -> None:
    """Quality assurance and validation commands."""
    pass


@quality.command("check")
@click.argument("target", type=click.Path(exists=True, path_type=Path))
@click.option(
    "--threshold",
    type=float,
    default=80.0,
    help="Coverage threshold percentage (default: 80)",
)
@click.option(
    "--max-cc",
    type=int,
    default=10,
    help="Maximum cyclomatic complexity (default: 10)",
)
@click.option(
    "--max-loc",
    type=int,
    default=200,
    help="Maximum lines of code per file (default: 200)",
)
def quality_check(  # noqa: C901
    target: Path,
    threshold: float,
    max_cc: int,
    max_loc: int,
) -> None:
    """Run quality gates on target file or directory.

    Args:
        target: File or directory to check
        threshold: Coverage threshold percentage
        max_cc: Maximum cyclomatic complexity
        max_loc: Maximum lines of code per file
    """
    import subprocess

    click.echo(f"Running quality gates on {target}...")
    click.echo()

    all_passed = True

    # Coverage check
    click.echo("1. Coverage check...")
    try:
        result = subprocess.run(
            ["pytest", "--cov", str(target), f"--cov-fail-under={threshold}"],
            capture_output=True,
            text=True,
        )
        if result.returncode == 0:
            click.echo(f"   ✓ Coverage ≥{threshold}%")
        else:
            click.echo(f"   ✗ Coverage <{threshold}%")
            all_passed = False
    except Exception as e:
        click.echo(f"   ✗ Coverage check failed: {e}")
        all_passed = False

    # Type checking
    click.echo("2. Type checking (mypy)...")
    try:
        result = subprocess.run(
            ["mypy", "--strict", str(target)],
            capture_output=True,
            text=True,
        )
        if result.returncode == 0:
            click.echo("   ✓ Type checking passed")
        else:
            click.echo("   ✗ Type checking failed")
            all_passed = False
    except Exception as e:
        click.echo(f"   ✗ Type check failed: {e}")
        all_passed = False

    # Linting
    click.echo("3. Linting (ruff)...")
    try:
        result = subprocess.run(
            ["ruff", "check", str(target)],
            capture_output=True,
            text=True,
        )
        if result.returncode == 0:
            click.echo("   ✓ Linting passed")
        else:
            click.echo("   ✗ Linting failed")
            all_passed = False
    except Exception as e:
        click.echo(f"   ✗ Linting failed: {e}")
        all_passed = False

    # Complexity
    click.echo(f"4. Complexity check (max CC: {max_cc})...")
    try:
        result = subprocess.run(
            ["radon", "cc", str(target), "-a", "-nb"],
            capture_output=True,
            text=True,
        )
        # Check if any function exceeds max CC
        high_cc = [
            line for line in result.stdout.split("\n")
            if line.strip() and " - " in line
        ]
        has_violations = any(
            int(line.split("(")[1].split(")")[0]) > max_cc
            for line in high_cc
            if "(" in line and ")" in line
        )
        if not has_violations:
            click.echo(f"   ✓ All functions ≤ CC {max_cc}")
        else:
            click.echo(f"   ✗ Some functions exceed CC {max_cc}")
            all_passed = False
    except Exception as e:
        click.echo(f"   ✗ Complexity check failed: {e}")
        all_passed = False

    # File size
    click.echo(f"5. File size check (max: {max_loc} LOC)...")
    try:
        result = subprocess.run(
            ["python", "scripts/check_file_size.py"],
            capture_output=True,
            text=True,
        )
        if result.returncode == 0:
            click.echo(f"   ✓ All files ≤ {max_loc} LOC")
        else:
            click.echo(f"   ✗ Some files exceed {max_loc} LOC")
            all_passed = False
    except Exception as e:
        click.echo(f"   ✗ File size check failed: {e}")
        all_passed = False

    click.echo()
    if all_passed:
        click.echo("✅ All quality gates passed")
        sys.exit(0)
    else:
        click.echo("❌ Some quality gates failed")
        sys.exit(1)
