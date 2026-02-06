"""Doctor command implementations for health checking."""

import sys
from pathlib import Path

import click


def check_environment() -> bool:
    """Check environment setup.

    Returns:
        True if environment is healthy
    """
    import subprocess

    all_ok = True

    # Check Python version
    click.echo("Python version:", nl=False)
    sys_version = sys.version_info
    if sys_version >= (3, 10):
        click.echo(f" {sys_version.major}.{sys_version.minor}.{sys_version.micro} ✓")
    else:
        click.echo(f" {sys_version.major}.{sys_version.minor}.{sys_version.micro} ✗")
        click.echo("  Warning: Python 3.10+ recommended")
        all_ok = False

    # Check Poetry
    click.echo("Poetry:", nl=False)
    try:
        result = subprocess.run(
            ["poetry", "--version"],
            capture_output=True,
            text=True,
        )
        if result.returncode == 0:
            click.echo(f" {result.stdout.strip()} ✓")
        else:
            click.echo(" not found ✗")
            all_ok = False
    except Exception:
        click.echo(" not found ✗")
        all_ok = False

    # Check Git
    click.echo("Git:", nl=False)
    try:
        result = subprocess.run(
            ["git", "--version"],
            capture_output=True,
            text=True,
        )
        if result.returncode == 0:
            click.echo(f" {result.stdout.strip()} ✓")
        else:
            click.echo(" not found ✗")
            all_ok = False
    except Exception:
        click.echo(" not found ✗")
        all_ok = False

    return all_ok


def check_project_structure(project_dir: Path) -> bool:
    """Check project structure.

    Args:
        project_dir: Project directory to check

    Returns:
        True if structure is valid
    """
    all_ok = True

    # Check for required directories
    required_dirs = ["docs", "src", "tests"]
    for dir_name in required_dirs:
        dir_path = project_dir / dir_name
        click.echo(f"{dir_name}/:", nl=False)
        if dir_path.is_dir():
            click.echo(" ✓")
        else:
            click.echo(" ✗")
            all_ok = False

    # Check for required files
    required_files = ["README.md", "pyproject.toml"]
    for file_name in required_files:
        file_path = project_dir / file_name
        click.echo(f"{file_name}:", nl=False)
        if file_path.is_file():
            click.echo(" ✓")
        else:
            click.echo(" ✗")
            all_ok = False

    return all_ok


def check_workstreams(project_dir: Path) -> bool:
    """Check workstream files.

    Args:
        project_dir: Project directory

    Returns:
        True if workstreams are valid
    """
    from sdp.core import WorkstreamParseError, parse_workstream

    ws_dir = project_dir / "docs" / "workstreams"
    if not ws_dir.is_dir():
        click.echo("No workstreams directory found")
        return True

    all_ok = True
    ws_files = list(ws_dir.rglob("WS-*.md"))

    if not ws_files:
        click.echo("No workstream files found")
        return True

    click.echo(f"\nChecking {len(ws_files)} workstream(s)...")

    for ws_file in ws_files:
        click.echo(f"  {ws_file.name}:", nl=False)
        try:
            parse_workstream(ws_file)
            click.echo(" ✓")
        except WorkstreamParseError as e:
            click.echo(" ✗")
            click.echo(f"    Error: {e}", err=True)
            all_ok = False
        except Exception as e:
            click.echo(" ✗")
            click.echo(f"    Error: {e}", err=True)
            all_ok = False

    return all_ok
