"""PRD document parser.

This module parses PRD documents into sections for validation and processing.
"""

import re


def parse_prd_sections(content: str) -> dict[str, str]:
    """Parse PRD content into sections.

    Args:
        content: PRD document content

    Returns:
        Dictionary mapping section names to their content
    """
    sections: dict[str, str] = {}
    current_section: str | None = None
    current_content: list[str] = []
    frontmatter_skipped = False
    frontmarker_count = 0
    has_frontmatter = content.strip().startswith("---")

    for line in content.split("\n"):
        # Skip frontmatter
        if has_frontmatter and not frontmatter_skipped:
            if line.strip() == "---":
                frontmarker_count += 1
                if frontmarker_count == 2:
                    frontmatter_skipped = True
                continue
            continue

        # Match section headers: ## 1. Section Name or ## Section Name
        if match := re.match(r'^##\s+(\d+\.\s+)?(.+)$', line):
            if current_section is not None:
                sections[current_section] = "\n".join(current_content).strip()
            current_section = match.group(2).strip()
            current_content = []
        elif current_section is not None:
            current_content.append(line)

    # Don't forget the last section
    if current_section is not None:
        sections[current_section] = "\n".join(current_content).strip()

    return sections


def get_frontmatter(content: str) -> dict[str, str]:
    """Extract frontmatter from PRD content.

    Args:
        content: PRD document content

    Returns:
        Dictionary of frontmatter key-value pairs
    """
    frontmatter = {}
    in_frontmatter = False

    for line in content.split("\n"):
        if line.strip() == "---":
            if not in_frontmatter:
                in_frontmatter = True
            else:
                break
            continue

        if in_frontmatter:
            if ":" in line:
                key, value = line.split(":", 1)
                frontmatter[key.strip()] = value.strip()

    return frontmatter


def update_frontmatter(content: str, updates: dict[str, str]) -> str:
    """Update frontmatter fields in PRD content.

    Args:
        content: PRD document content
        updates: Dictionary of fields to update

    Returns:
        Updated content
    """
    lines = content.split("\n")
    frontmatter_end_idx = _find_frontmatter_end(lines)

    # If no frontmatter, return content unchanged
    if frontmatter_end_idx == -1:
        return content

    return _process_frontmatter_updates(lines, frontmatter_end_idx, updates)


def _find_frontmatter_end(lines: list[str]) -> int:
    """Find the line index where frontmatter ends.

    Returns:
        Index of closing "---", or -1 if no frontmatter found
    """
    frontmatter_seen = False

    for idx, line in enumerate(lines):
        if line.strip() == "---":
            if not frontmatter_seen:
                frontmatter_seen = True
            else:
                return idx

    return -1


def _process_frontmatter_updates(
    lines: list[str],
    frontmatter_end_idx: int,
    updates: dict[str, str],
) -> str:
    """Process frontmatter and apply updates.

    Args:
        lines: Original content lines
        frontmatter_end_idx: Index of closing "---"
        updates: Dictionary of fields to update (will be mutated)

    Returns:
        Updated content
    """
    updated_lines: list[str] = []
    in_frontmatter = False

    for idx, line in enumerate(lines):
        if line.strip() == "---":
            in_frontmatter = _handle_frontmatter_marker(
                in_frontmatter, updated_lines, line, updates
            )
            continue

        if in_frontmatter:
            line = _update_frontmatter_field(line, updates)

        updated_lines.append(line)

    return "\n".join(updated_lines)


def _handle_frontmatter_marker(
    in_frontmatter: bool,
    updated_lines: list[str],
    line: str,
    updates: dict[str, str],
) -> bool:
    """Handle frontmatter marker ("---") line.

    Args:
        in_frontmatter: Whether we're currently in frontmatter
        updated_lines: List to append updated lines to
        line: Current line being processed
        updates: Remaining updates to apply

    Returns:
        New in_frontmatter state
    """
    if not in_frontmatter:
        return True
    else:
        # Add any remaining updates before closing ---
        for key, value in updates.items():
            updated_lines.append(f"{key}: {value}")
        updated_lines.append(line)
        return False


def _update_frontmatter_field(line: str, updates: dict[str, str]) -> str:
    """Update a frontmatter field if it's in the updates dict.

    Args:
        line: Current line being processed
        updates: Dictionary of fields to update (will be mutated)

    Returns:
        Updated line (or original if no update)
    """
    if ":" in line:
        key = line.split(":", 1)[0].strip()
        if key in updates:
            updated_line = f"{key}: {updates[key]}"
            del updates[key]
            return updated_line
    return line
