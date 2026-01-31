"""Section extraction utilities for PRD validation."""

import re
from typing import Optional


def _extract_body(content: str) -> str:
    """Extract body from markdown content.

    Args:
        content: Full markdown content

    Returns:
        Body content (without frontmatter)
    """
    # Remove frontmatter
    match = re.match(r"^---\n(.*?)\n---", content, re.DOTALL)
    if match:
        return content[match.end():]
    return content


def _extract_section(body: str, section_name: str) -> str:
    """Extract a section from markdown body.

    Args:
        body: Markdown body content
        section_name: Name of section to extract (without ##)

    Returns:
        Section content
    """
    # Look for ## Section Name
    pattern = rf"##\s+{re.escape(section_name)}\s*\n(.*?)(?=\n##\s|\Z)"
    match = re.search(pattern, body, re.DOTALL)
    if match:
        return match.group(1).strip()
    return ""


def _extract_code_block(section: str, language: str = "python") -> Optional[str]:
    """Extract code block from section.

    Args:
        section: Section content
        language: Code block language (default: "python")

    Returns:
        Code block content or None
    """
    # Look for ```language ... ```
    pattern = rf"```{language}\s*\n(.*?)\n```"
    match = re.search(pattern, section, re.DOTALL)
    if match:
        return match.group(1).strip()
    return None
