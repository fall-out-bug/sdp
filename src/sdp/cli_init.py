"""SDP init command - Initialize SDP in a project.

Creates standard directory structure, templates, and configuration.
"""

from datetime import datetime
from pathlib import Path

import click

from sdp.templates.init_templates import (
    INDEX_TEMPLATE,
    PROJECT_MAP_TEMPLATE,
    WS_TEMPLATE,
)


@click.command()
@click.option(
    "--force",
    is_flag=True,
    help="Overwrite existing files",
)
@click.option(
    "--path",
    type=click.Path(file_okay=False, dir_okay=True, path_type=Path),
    default=None,
    help="Target directory (defaults to current directory)",
)
def init(force: bool, path: Path | None) -> None:
    """Initialize SDP in current project.

    Creates standard directory structure:
    - docs/workstreams/ (INDEX.md, TEMPLATE.md, backlog/)
    - docs/PROJECT_MAP.md
    - sdp.local/ (for project-local extensions)

    Example:
        $ sdp init
        ✓ Created docs/workstreams/INDEX.md
        ✓ Created docs/workstreams/TEMPLATE.md
        ✓ Created docs/PROJECT_MAP.md
        ✓ Created sdp.local/

        $ sdp init --path /tmp/my-project
        ✓ Created /tmp/my-project/docs/workstreams/INDEX.md
        ...

        SDP initialized! Next steps:
        1. Edit docs/PROJECT_MAP.md with your project info
        2. Run: sdp extension list
        3. Start: /idea "your first feature"
    """
    # Determine target directory
    target_dir = path.resolve() if path else Path.cwd()

    # Create target directory if it doesn't exist
    if not target_dir.exists():
        target_dir.mkdir(parents=True)

    # Get project name from target directory
    project_name = target_dir.name
    current_date = datetime.now().strftime("%Y-%m-%d")

    # Define paths relative to target directory
    docs_dir = target_dir / "docs"
    workstreams_dir = docs_dir / "workstreams"
    backlog_dir = workstreams_dir / "backlog"
    project_map_file = docs_dir / "PROJECT_MAP.md"
    index_file = workstreams_dir / "INDEX.md"
    template_file = workstreams_dir / "TEMPLATE.md"
    sdp_local_dir = target_dir / "sdp.local"
    
    # Track created files
    created_files = []
    skipped_files = []
    
    # Create directories
    for directory in [docs_dir, workstreams_dir, backlog_dir, sdp_local_dir]:
        if not directory.exists():
            directory.mkdir(parents=True)
            created_files.append(str(directory) + "/")
    
    # Create PROJECT_MAP.md
    if project_map_file.exists() and not force:
        skipped_files.append(str(project_map_file))
    else:
        content = PROJECT_MAP_TEMPLATE.format(
            project_name=project_name,
            date=current_date,
        )
        project_map_file.write_text(content)
        created_files.append(str(project_map_file))
    
    # Create INDEX.md
    if index_file.exists() and not force:
        skipped_files.append(str(index_file))
    else:
        index_file.write_text(INDEX_TEMPLATE)
        created_files.append(str(index_file))
    
    # Create TEMPLATE.md
    if template_file.exists() and not force:
        skipped_files.append(str(template_file))
    else:
        template_file.write_text(WS_TEMPLATE)
        created_files.append(str(template_file))
    
    # Display results
    if created_files:
        for file in created_files:
            click.echo(f"✓ Created {file}")
    
    if skipped_files:
        click.echo("")
        click.echo("Skipped (already exists, use --force to overwrite):")
        for file in skipped_files:
            click.echo(f"  - {file}")
    
    # Success message
    click.echo("")
    click.echo("SDP initialized! Next steps:")
    click.echo(f"  1. Edit {project_map_file} with your project info")
    click.echo("  2. Run: sdp extension list")
    click.echo("  3. Start: /idea \"your first feature\"")
