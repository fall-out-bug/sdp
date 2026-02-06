"""Section extraction utilities for project map parsing."""

import re


def extract_project_name(content: str) -> str | None:
    """Extract project name from title line.

    Args:
        content: Markdown content

    Returns:
        Project name or None if not found
    """
    # Try multiple formats:
    # 1. "# PROJECT_MAP: SDP" (SDP format)
    # 2. "# Project Map: <name>" (standard format)
    # 3. "# PROJECT: <name>" (alternative format)
    patterns = [
        r"^# PROJECT_MAP:\s*(.+)$",
        r"^# Project Map:\s*(.+)$",
        r"^# PROJECT:\s*(.+)$",
    ]
    for pattern in patterns:
        match = re.search(pattern, content, re.MULTILINE)
        if match:
            return match.group(1).strip()
    return None


def extract_section(content: str, section_name: str) -> str:
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
