"""Parse workstream markdown files."""

import re
from dataclasses import dataclass
from pathlib import Path


@dataclass
class WSMetadata:
    """Parsed WS metadata from markdown file.

    Attributes:
        ws_id: Workstream ID (e.g., WS-100-01)
        feature: Feature ID (e.g., F100)
        title: WS title
        goal: Goal description
        acceptance_criteria: List of AC strings
        dependencies: List of dependent WS IDs
        size: Size (SMALL, MEDIUM, LARGE)
        status: Status (backlog, active, completed)
    """

    ws_id: str
    feature: str
    title: str
    goal: str
    acceptance_criteria: list[str]
    dependencies: list[str]
    size: str
    status: str


def parse_ws_file(file_path: Path) -> WSMetadata:
    """Parse WS markdown file.

    Args:
        file_path: Path to WS-XXX-YY.md file

    Returns:
        WSMetadata with parsed data

    Raises:
        ValueError: If file format invalid or required fields missing
    """
    content = file_path.read_text(encoding="utf-8")

    # Parse frontmatter (YAML between ---)
    frontmatter_match = re.search(
        r"^---\n(.*?)\n---", content, re.MULTILINE | re.DOTALL
    )
    if not frontmatter_match:
        raise ValueError(f"No frontmatter found in {file_path}")

    frontmatter = frontmatter_match.group(1)

    # Extract fields from frontmatter
    ws_id = _extract_field(frontmatter, "ws_id")
    feature = _extract_field(frontmatter, "feature")
    size = _extract_field(frontmatter, "size")
    status = _extract_field(frontmatter, "status")

    # Extract title from ## WS-XXX-YY: Title
    title_match = re.search(r"## WS-\d{3}-\d{2}: (.+)", content)
    title = title_match.group(1) if title_match else "Untitled"

    # Extract Goal section
    goal_match = re.search(
        r"### 🎯 Цель \(Goal\).*?\*\*Что должно РАБОТАТЬ.*?:\*\*\n(.+?)(?=\n\*\*|---)",
        content,
        re.DOTALL,
    )
    goal = goal_match.group(1).strip() if goal_match else ""

    # Extract Acceptance Criteria
    ac_matches = re.findall(r"- \[ \] (AC\d+: .+)", content)
    acceptance_criteria = ac_matches

    # Extract Dependencies
    deps_match = re.search(
        r"### Зависимость\n\n(.+?)(?=\n###|\nWS-|\Z)", content, re.DOTALL
    )
    dependencies = []
    if deps_match:
        deps_text = deps_match.group(1).strip()
        if deps_text and "Независимый" not in deps_text:
            # Extract WS-XXX-YY patterns
            dependencies = re.findall(r"WS-\d{3}-\d{2}", deps_text)

    return WSMetadata(
        ws_id=ws_id,
        feature=feature,
        title=title,
        goal=goal,
        acceptance_criteria=acceptance_criteria,
        dependencies=dependencies,
        size=size,
        status=status,
    )


def _extract_field(frontmatter: str, field: str) -> str:
    """Extract field from YAML frontmatter.

    Args:
        frontmatter: YAML frontmatter text
        field: Field name to extract

    Returns:
        Field value

    Raises:
        ValueError: If field not found
    """
    match = re.search(rf"^{field}:\s*(.+)$", frontmatter, re.MULTILINE)
    if not match:
        raise ValueError(f"Field '{field}' not found in frontmatter")
    return match.group(1).strip()
