"""Project type auto-detection.

This module detects the project type (service, library, cli)
based on the project's file structure and configuration.
"""

from pathlib import Path

from .profiles import ProjectType


def detect_project_type(project_path: Path) -> ProjectType:
    """Auto-detect project type from file structure.

    Detection strategy:
    1. If docker-compose.yml exists → SERVICE
    2. If cli.py with Click/Typer exists → CLI
    3. Default → LIBRARY

    Args:
        project_path: Path to the project root

    Returns:
        Detected project type
    """
    # Check for docker-compose.yml + API endpoints → service
    docker_compose_paths = [
        project_path / "docker-compose.yml",
        project_path / "docker-compose.yaml",
        project_path / "compose.yml",
        project_path / "compose.yaml",
    ]

    for docker_path in docker_compose_paths:
        if docker_path.exists():
            # Additional check: look for API indicators
            # (FastAPI, Flask, Django, etc.)
            if _has_api_framework(project_path):
                return ProjectType.SERVICE
            # Docker compose is strong signal for service
            return ProjectType.SERVICE

    # Check for cli.py with Click/Typer → cli
    cli_files = list(project_path.glob("**/cli.py")) + list(project_path.glob("**/main.py"))

    for cli_file in cli_files:
        if _has_cli_framework(cli_file):
            return ProjectType.CLI

    # Check for pyproject.toml with entry points → cli
    pyproject = project_path / "pyproject.toml"
    if pyproject.exists():
        content = pyproject.read_text()
        if "[project.scripts]" in content or "[tool.poetry.scripts]" in content:
            return ProjectType.CLI

    # Default: library
    return ProjectType.LIBRARY


def _has_api_framework(project_path: Path) -> bool:
    """Check if project uses an API framework.

    Args:
        project_path: Path to the project root

    Returns:
        True if FastAPI, Flask, Django, etc. detected
    """
    # Check requirements files
    for req_file in ["requirements.txt", "pyproject.toml", "setup.py", "Pipfile"]:
        req_path = project_path / req_file
        if req_path.exists():
            content = req_path.read_text().lower()
            if any(fw in content for fw in ["fastapi", "flask", "django", "tornado", "aiohttp"]):
                return True

    # Check source files
    for src_file in project_path.rglob("*.py"):
        # Skip common non-source directories
        if any(skip in str(src_file) for skip in ["venv", ".venv", "__pycache__", ".tox", "node_modules"]):
            continue
        try:
            content = src_file.read_text()
            # Check for imports
            if any(fw in content for fw in ["from fastapi", "import fastapi", "from flask", "import flask", "from django", "import django"]):
                return True
        except Exception:
            continue

    return False


def _has_cli_framework(cli_file: Path) -> bool:
    """Check if file uses a CLI framework.

    Args:
        cli_file: Path to potential CLI file

    Returns:
        True if Click, Typer, or argparse detected
    """
    try:
        content = cli_file.read_text()
        content_lower = content.lower()

        # Check for common CLI frameworks
        if any(fw in content_lower for fw in ["import click", "from click", "import typer", "from typer", "argparse"]):
            return True

        # Check for common CLI patterns
        if any(pattern in content for pattern in ["@click.command", "@app.command", "ArgumentParser", "click.option"]):
            return True

    except Exception:
        pass

    return False
