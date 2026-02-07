"""Main parser for PROJECT_MAP.md files."""

from pathlib import Path

from sdp.core.project_map.extractors import extract_project_name, extract_section
from sdp.core.project_map.table_parsers import (
    parse_constraints,
    parse_decisions_table,
    parse_tech_stack_table,
)
from sdp.core.project_map_types import ProjectMap, ProjectMapParseError


def parse_project_map(file_path: Path) -> ProjectMap:
    """Parse PROJECT_MAP.md file.

    Args:
        file_path: Path to PROJECT_MAP.md file

    Returns:
        Parsed ProjectMap instance

    Raises:
        FileNotFoundError: If file doesn't exist
        ProjectMapParseError: If file format is invalid
    """
    if not file_path.exists():
        raise FileNotFoundError(f"Project map file not found: {file_path}")

    content = file_path.read_text(encoding="utf-8")

    # Extract project name from title
    project_name = extract_project_name(content)
    if not project_name:
        raise ProjectMapParseError("Could not extract project name from title")

    # Parse decisions table
    decisions = parse_decisions_table(content)

    # Parse constraints
    constraints = parse_constraints(content)

    # Extract current state section
    current_state = extract_section(content, "Current State")

    # Extract patterns section
    patterns = extract_section(content, "Patterns & Conventions")

    # Parse tech stack table
    tech_stack = parse_tech_stack_table(content)

    return ProjectMap(
        project_name=project_name,
        decisions=decisions,
        constraints=constraints,
        current_state=current_state,
        patterns=patterns,
        tech_stack=tech_stack,
        file_path=file_path,
    )
