"""Markdown parsing helpers for workstream files."""

import re
from typing import Any, Optional

import yaml

from sdp.domain.workstream import AcceptanceCriterion


def extract_frontmatter(content: str, error_path: Optional[str] = None) -> dict[str, Any]:
    """Extract YAML frontmatter from markdown.

    Args:
        content: Markdown file content
        error_path: Optional path for error messages

    Returns:
        Parsed frontmatter as dict

    Raises:
        ValueError: If frontmatter is missing or invalid
    """
    match = re.match(r"^---\n(.*?)\n---", content, re.DOTALL)
    if not match:
        raise ValueError(
            f"No frontmatter found (must start with ---) in {error_path or 'file'}"
        )

    frontmatter_text = match.group(1)
    try:
        data: Any = yaml.safe_load(frontmatter_text)
    except yaml.YAMLError as e:
        raise ValueError(
            f"Invalid YAML in frontmatter: {e} in {error_path or 'file'}"
        ) from e

    if not isinstance(data, dict):
        raise ValueError(f"Frontmatter must be a YAML dict in {error_path or 'file'}")

    required_fields = {"ws_id", "feature", "status", "size"}
    missing = required_fields - set(data.keys())
    if missing:
        raise ValueError(
            f"Missing required fields: {missing} in {error_path or 'file'}"
        )

    return data


def strip_frontmatter(content: str) -> str:
    """Remove frontmatter, return body only."""
    match = re.match(r"^---\n.*?\n---\n(.*)", content, re.DOTALL)
    return match.group(1) if match else content


def extract_title(body: str) -> str:
    """Extract title from first ## heading."""
    match = re.search(r"^## (.+)$", body, re.MULTILINE)
    return match.group(1) if match else ""


def extract_section(body: str, section_name: str) -> str:
    """Extract content of a ### section by name (case-insensitive).

    Args:
        body: Markdown body content
        section_name: Section heading to find

    Returns:
        Section content without heading
    """
    heading_pattern = rf"^### .*{section_name}.*$"
    heading_match = re.search(heading_pattern, body, re.MULTILINE | re.IGNORECASE)

    if not heading_match:
        return ""

    start_pos = heading_match.end() + 1
    if start_pos >= len(body):
        return ""

    remaining = body[start_pos:]
    next_heading = re.search(r"^###", remaining, re.MULTILINE)

    end_pos = next_heading.start() if next_heading else len(remaining)
    content = remaining[:end_pos]
    content = re.sub(r"\n---\s*$", "", content)
    return content.strip()


def extract_acceptance_criteria(body: str) -> list[AcceptanceCriterion]:
    """Extract acceptance criteria from body.

    Args:
        body: Markdown body content

    Returns:
        List of parsed AcceptanceCriterion objects
    """
    criteria: list[AcceptanceCriterion] = []
    pattern = r"- \[([ x])\] (AC\d+): (.+)"
    for match in re.finditer(pattern, body, re.IGNORECASE):
        checked_char = match.group(1)
        ac_id = match.group(2)
        description = match.group(3)
        checked = checked_char.lower() == "x"
        criteria.append(
            AcceptanceCriterion(id=ac_id, description=description, checked=checked)
        )
    return criteria


def extract_dependencies(body: str) -> list[str]:
    """Extract WS dependencies from Dependencies section.

    Args:
        body: Markdown body content

    Returns:
        List of workstream IDs (PP-FFF-SS or WS-FFF-SS format)
    """
    dep_section = extract_section(body, "Dependencies")
    if not dep_section or dep_section.lower() in ("none", ""):
        return []
    # Support PP-FFF-SS (00-032-18) and legacy WS-FFF-SS formats
    pattern = r"(?:\d{2}-\d{3}-\d{2}|WS-\d+-\d+)"
    return re.findall(pattern, dep_section)


def extract_steps(body: str) -> list[str]:
    """Extract numbered steps from Steps section.

    Args:
        body: Markdown body content

    Returns:
        List of step descriptions
    """
    steps_section = extract_section(body, "Steps")
    steps: list[str] = []

    # Match patterns like "1. Step description" or "#### 1. Step"
    for line in steps_section.split("\n"):
        line = line.strip()
        # Skip empty lines and headings
        if not line or line.startswith("#"):
            continue
        # Match: "1. Step description" or "#### 1. Step"
        match = re.match(r"^(?:####\s*)?(\d+)\.\s+(.+)", line)
        if match:
            steps.append(match.group(2).strip())

    return steps


def extract_code_blocks(body: str) -> list[str]:
    """Extract fenced code blocks from body.

    Args:
        body: Markdown body content

    Returns:
        List of code block contents
    """
    # Match ```language\ncode```
    pattern = r"```[\w]*\n(.+?)```"
    return re.findall(pattern, body, re.DOTALL)
