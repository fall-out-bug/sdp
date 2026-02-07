"""Metadata collection for SDP init wizard."""

from pathlib import Path


def collect_metadata(
    target_dir: Path, non_interactive: bool
) -> tuple[str, str, str]:
    """Collect project metadata from user.

    Args:
        target_dir: Target directory for project
        non_interactive: Use defaults without prompting

    Returns:
        Tuple of (project_name, description, author)
    """
    import click

    project_name = target_dir.name
    description = "SDP project"
    author = "Your Name"

    if non_interactive:
        return project_name, description, author

    # Interactive prompts
    click.echo(click.style("Step 1: Project Metadata", fg="cyan", bold=True))

    project_name = click.prompt(
        "Project name",
        default=project_name,
        type=str,
    )

    description = click.prompt(
        "Description",
        default=description,
        type=str,
    )

    author = click.prompt(
        "Author",
        default=author,
        type=str,
    )

    return project_name, description, author
