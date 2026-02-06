"""PRD (Product Requirements Document) commands.

Provides CLI commands for validating PRDs and detecting project types.
"""

import sys
from pathlib import Path

import click


@click.group()
def prd() -> None:
    """PRD (Product Requirements Document) operations."""
    pass


@prd.command("validate")
@click.argument(
    "prd_file",
    type=click.Path(exists=True, path_type=Path),
)
@click.option(
    "--exit-code-on-error",
    is_flag=True,
    help="Exit with code 1 if validation errors found",
)
def prd_validate(prd_file: Path, exit_code_on_error: bool) -> None:
    """Validate a PRD document against section limits.

    Args:
        prd_file: Path to PRD file (PROJECT_MAP.md)
        exit_code_on_error: Exit with code 1 if errors found
    """
    from sdp.prd.validator import (
        format_validation_issues,
        has_critical_issues,
        validate_prd_file,
    )

    issues = validate_prd_file(prd_file)

    if issues:
        click.echo(format_validation_issues(issues))
        if has_critical_issues(issues):
            if exit_code_on_error:
                sys.exit(1)
    else:
        click.echo("✅ PRD validation passed")


@prd.command("detect-type")
@click.argument(
    "project_path",
    type=click.Path(exists=True, path_type=Path),
)
def prd_detect_type(project_path: Path) -> None:
    """Detect project type from file structure.

    Args:
        project_path: Path to project root
    """
    from sdp.prd.detector import detect_project_type

    project_type = detect_project_type(project_path)
    click.echo(f"Detected project type: {project_type.value}")
