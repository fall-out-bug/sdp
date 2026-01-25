"""PRD document scaffolding.

This module generates PRD document templates based on project profiles.
"""

import json
from datetime import datetime
from pathlib import Path

from .profiles import PRDProfile, PRDSection, ProjectType, get_profile
from .detector import detect_project_type


FRONTmatter_TEMPLATE = """---
project_type: {project_type}
prd_version: "2.0"
last_updated: {last_updated}
diagrams_hash: {diagrams_hash}
---

"""


def generate_prd_template(
    project_path: Path,
    project_name: str,
    project_type: ProjectType | None = None,
) -> str:
    """Generate a PRD template for a project.

    Args:
        project_path: Path to the project
        project_name: Name of the project
        project_type: Project type (auto-detected if None)

    Returns:
        Complete PRD document as string
    """
    # Auto-detect project type if not specified
    if project_type is None:
        project_type = detect_project_type(project_path)

    # Get profile
    profile = get_profile(project_type)

    # Generate frontmatter
    frontmatter = FRONTmatter_TEMPLATE.format(
        project_type=project_type.value,
        last_updated=datetime.now().strftime("%Y-%m-%d"),
        diagrams_hash="",  # Empty initially, set after diagram generation
    )

    # Generate document header
    header = f"# PROJECT_MAP: {project_name}\n\n"

    # Generate sections
    sections_content = []
    for idx, section in enumerate(profile.sections, start=1):
        section_text = f"{section.template}\n"
        sections_content.append(section_text)

    # Combine all parts
    prd_content = frontmatter + header + "\n".join(sections_content)

    return prd_content


def create_prd_file(
    project_path: Path,
    project_name: str,
    output_path: Path | None = None,
    project_type: ProjectType | None = None,
) -> Path:
    """Create a PRD file for a project.

    Args:
        project_path: Path to the project
        project_name: Name of the project
        output_path: Where to save the PRD file (default: project_path/docs/PROJECT_MAP.md)
        project_type: Project type (auto-detected if None)

    Returns:
        Path to the created PRD file
    """
    # Determine output path
    if output_path is None:
        output_path = project_path / "docs" / "PROJECT_MAP.md"

    # Ensure parent directory exists
    output_path.parent.mkdir(parents=True, exist_ok=True)

    # Generate content
    content = generate_prd_template(
        project_path=project_path,
        project_name=project_name,
        project_type=project_type,
    )

    # Check if file exists
    if output_path.exists():
        # For updates, preserve manual edits
        content = _merge_prd_content(output_path, content)

    # Write file
    output_path.write_text(content)

    return output_path


def _merge_prd_content(existing_path: Path, new_content: str) -> str:
    """Merge existing PRD content with new template.

    This preserves manual edits while updating structure.

    Args:
        existing_path: Path to existing PRD file
        new_content: New template content

    Returns:
        Merged content
    """
    existing_content = existing_path.read_text()

    # For now, just return existing content
    # In a full implementation, this would:
    # 1. Parse existing sections
    # 2. Update frontmatter
    # 3. Add new sections
    # 4. Preserve manual edits in existing sections

    return existing_content


def update_prd_frontmatter(
    prd_path: Path,
    diagrams_hash: str | None = None,
) -> None:
    """Update the frontmatter of an existing PRD file.

    Args:
        prd_path: Path to the PRD file
        diagrams_hash: New diagrams hash (if None, leaves unchanged)
    """
    if not prd_path.exists():
        return

    content = prd_path.read_text()
    lines = content.split("\n")

    # Find and update frontmatter
    in_frontmatter = False
    updated_lines = []

    for line in lines:
        if line.strip() == "---":
            if not in_frontmatter:
                in_frontmatter = True
            else:
                in_frontmatter = False
            updated_lines.append(line)
        elif in_frontmatter:
            if line.startswith("last_updated:"):
                updated_lines.append(f"last_updated: {datetime.now().strftime('%Y-%m-%d')}")
            elif diagrams_hash and line.startswith("diagrams_hash:"):
                updated_lines.append(f"diagrams_hash: {diagrams_hash}")
            elif not diagrams_hash and line.startswith("diagrams_hash:") and diagrams_hash is not None:
                updated_lines.append(f"diagrams_hash: {diagrams_hash}")
            else:
                updated_lines.append(line)
        else:
            updated_lines.append(line)

    prd_path.write_text("\n".join(updated_lines))
