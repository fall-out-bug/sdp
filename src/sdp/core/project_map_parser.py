"""Project map parsing for SDP projects.

This module provides parsing functionality for PROJECT_MAP.md files.
"""

import re
from pathlib import Path

from sdp.core.project_map_types import (
    Constraint,
    Decision,
    ProjectMap,
    ProjectMapParseError,
    TechStackItem,
)


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
    project_name = _extract_project_name(content)
    if not project_name:
        raise ProjectMapParseError("Could not extract project name from title")

    # Parse decisions table
    decisions = _parse_decisions_table(content)

    # Parse constraints
    constraints = _parse_constraints(content)

    # Extract current state section
    current_state = _extract_section(content, "Current State")

    # Extract patterns section
    patterns = _extract_section(content, "Patterns & Conventions")

    # Parse tech stack table
    tech_stack = _parse_tech_stack_table(content)

    return ProjectMap(
        project_name=project_name,
        decisions=decisions,
        constraints=constraints,
        current_state=current_state,
        patterns=patterns,
        tech_stack=tech_stack,
        file_path=file_path,
    )


# =============================================================================
# PRIVATE HELPER FUNCTIONS
# =============================================================================


def _extract_project_name(content: str) -> str | None:
    """Extract project name from title line.

    Args:
        content: Markdown content

    Returns:
        Project name or None if not found
    """
    match = re.search(r"^# Project Map:\s*(.+)$", content, re.MULTILINE)
    if match:
        return match.group(1).strip()
    return None


def _parse_decisions_table(content: str) -> list[Decision]:
    """Parse decisions from Quick Reference table.

    Args:
        content: Markdown content

    Returns:
        List of Decision objects
    """
    decisions: list[Decision] = []

    # Find the decisions table section
    table_start = content.find("## Ключевые решения")
    if table_start == -1:
        return decisions

    # Find the table (starts after header row)
    # Note: There may be blank lines between section header and table
    table_match = re.search(
        r"\|\s*Область\s*\|[^\n]*\n\|[-\s|]+\|\s*\n(.*?)(?=\n\s*\n|\n\s*##|\Z)",
        content[table_start:],
        re.DOTALL,
    )

    if not table_match:
        return decisions

    table_rows = table_match.group(1).strip()

    # Parse each row
    for line in table_rows.split("\n"):
        line = line.strip()
        if not line.startswith("|"):
            continue

        # Split by | and remove empty strings
        parts = [p.strip() for p in line.split("|") if p.strip()]

        if len(parts) >= 4:
            area = parts[0].replace("**", "").strip()
            decision = parts[1].strip()
            adr_raw = parts[2].strip()
            date = parts[3].strip()

            # Extract ADR from markdown link if present: [ADR-001](path) -> ADR-001
            adr_match = re.search(r"\[(ADR-\d+)\]", adr_raw)
            adr = adr_match.group(1) if adr_match else adr_raw

            # Skip empty rows (check for placeholder text)
            if area and decision and adr and date and not area.startswith("_"):
                decisions.append(Decision(area=area, decision=decision, adr=adr, date=date))

    return decisions


def _parse_constraints(content: str) -> list[Constraint]:
    """Parse constraints from Active Constraints section.

    Args:
        content: Markdown content

    Returns:
        List of Constraint objects
    """
    constraints: list[Constraint] = []

    # Find Active Constraints section
    section_start = content.find("## Active Constraints")
    if section_start == -1:
        return constraints

    # Extract section content until next ## or ---
    section_match = re.search(
        r"## Active Constraints\s*\n(.*?)(?=\n\s*##|\n\s*---|\Z)",
        content,
        re.DOTALL,
    )

    if not section_match:
        return constraints

    section_content = section_match.group(1)

    # Parse subsections (### Category)
    current_category = ""
    for line in section_content.split("\n"):
        line = line.strip()

        # Check for category header
        category_match = re.match(r"^### (.+)$", line)
        if category_match:
            current_category = category_match.group(1).strip()
            continue

        # Check for constraint item (starts with -)
        if line.startswith("-") and current_category:
            # Remove markdown list marker and extract description
            description = re.sub(r"^-\s*", "", line).strip()
            if description:
                constraints.append(Constraint(category=current_category, description=description))

    return constraints


def _extract_section(content: str, section_name: str) -> str:
    """Extract section content by name.

    Args:
        content: Markdown content
        section_name: Section name to extract

    Returns:
        Section content (raw markdown) or empty string
    """
    # Try different header patterns (## and ###)
    # Allow optional trailing content after section name (e.g., "Current State (Live)")
    escaped_name = re.escape(section_name)

    # Find section header (## or ###)
    for header_level in ["##", "###"]:
        # Pattern to match header with section name and optional trailing content
        header_pattern = rf"^{re.escape(header_level)}\s+{escaped_name}(?:[^\n]*)?$"
        header_match = re.search(header_pattern, content, re.MULTILINE)

        if not header_match:
            continue

        # Find the start of content (after the header line and its newline)
        # header_match.end() is at the end of the matched line
        # Skip any newline character to get to the actual content
        content_start = header_match.end()
        # Skip the newline if present
        if content_start < len(content) and content[content_start] == "\n":
            content_start += 1

        # Find the next section header at the same or higher level
        # For ##, stop at next ##
        # For ###, stop at next ### or ##
        if header_level == "##":
            next_section_pattern = r"^\s*##\s+"
        else:  # ###
            next_section_pattern = r"^\s*(?:###|##)\s+"

        # Find the next section header
        next_match = re.search(next_section_pattern, content[content_start:], re.MULTILINE)

        if next_match:
            # Extract content up to the next section (excluding the newline before it)
            content_end = content_start + next_match.start()
            section_content = content[content_start:content_end]
        else:
            # No next section, extract to end of file
            section_content = content[content_start:]

        return section_content.strip()

    return ""


def _parse_tech_stack_table(content: str) -> list[TechStackItem]:
    """Parse tech stack from table.

    Args:
        content: Markdown content

    Returns:
        List of TechStackItem objects
    """
    items: list[TechStackItem] = []

    # Find Tech Stack section
    section_start = content.find("## Tech Stack")
    if section_start == -1:
        return items

    # Find the table
    table_match = re.search(
        r"\|\s*Layer\s*\|[^\n]*\n\|[-\s|]+\|\s*\n(.*?)(?=\n\s*\n|\n\s*##|\Z)",
        content[section_start:],
        re.DOTALL,
    )

    if not table_match:
        return items

    table_rows = table_match.group(1).strip()

    # Parse each row
    for line in table_rows.split("\n"):
        line = line.strip()
        if not line.startswith("|"):
            continue

        # Split by | and remove empty strings
        parts = [p.strip() for p in line.split("|") if p.strip()]

        if len(parts) >= 2:
            layer = parts[0].strip()
            technology = parts[1].strip()

            # Skip placeholder rows
            if layer and technology and not layer.startswith("_"):
                items.append(TechStackItem(layer=layer, technology=technology))

    return items
