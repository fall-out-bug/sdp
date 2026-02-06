"""Table parsing utilities for project map parsing."""

import re

from sdp.core.project_map_types import Constraint, Decision, TechStackItem


def parse_decisions_table(content: str) -> list[Decision]:
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


def parse_constraints(content: str) -> list[Constraint]:
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


def parse_tech_stack_table(content: str) -> list[TechStackItem]:
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
