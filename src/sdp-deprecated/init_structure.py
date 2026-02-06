"""Directory structure creation for SDP init wizard."""

from datetime import datetime
from pathlib import Path

from sdp.templates.init_templates import (
    INDEX_TEMPLATE,
    PROJECT_MAP_TEMPLATE,
    WS_TEMPLATE,
)


def create_structure(
    target_dir: Path, project_name: str, force: bool
) -> tuple[list[str], list[str]]:
    """Create directory structure.

    Args:
        target_dir: Target directory
        project_name: Name of project
        force: Overwrite existing files

    Returns:
        Tuple of (created_files, skipped_files)
    """
    current_date = datetime.now().strftime("%Y-%m-%d")

    # Define paths
    docs_dir = target_dir / "docs"
    workstreams_dir = docs_dir / "workstreams"
    backlog_dir = workstreams_dir / "backlog"
    project_map_file = docs_dir / "PROJECT_MAP.md"
    index_file = workstreams_dir / "INDEX.md"
    template_file = workstreams_dir / "TEMPLATE.md"
    sdp_local_dir = target_dir / "sdp.local"

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

    return created_files, skipped_files
