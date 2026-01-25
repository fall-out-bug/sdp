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
    in_frontmatter = False
    frontmatter_seen = False
    updated_lines = []
    frontmatter_end_idx = -1

    # First, find where frontmatter ends
    for idx, line in enumerate(lines):
        if line.strip() == "---":
            if not frontmatter_seen:
                frontmatter_seen = True
                in_frontmatter = True
            else:
                in_frontmatter = False
                frontmatter_end_idx = idx
                break

    # If no frontmatter, return content unchanged
    if frontmatter_end_idx == -1:
        return content

    # Process frontmatter
    in_frontmatter = False
    for idx, line in enumerate(lines):
        if line.strip() == "---":
            if not in_frontmatter:
                in_frontmatter = True
            else:
                in_frontmatter = False
                # Add any remaining updates before closing ---
                for key, value in updates.items():
                    updated_lines.append(f"{key}: {value}")
                updated_lines.append(line)
                continue

        if in_frontmatter:
            if ":" in line:
                key = line.split(":", 1)[0].strip()
                if key in updates:
                    updated_lines.append(f"{key}: {updates[key]}")
                    del updates[key]
                    continue

        updated_lines.append(line)

    return "\n".join(updated_lines)
