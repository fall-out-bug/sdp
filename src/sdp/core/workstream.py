"""Workstream parsing and validation for SDP markdown files."""

import re
from dataclasses import dataclass, field
from enum import Enum
from pathlib import Path
from typing import Any, Optional

import yaml  # type: ignore[import-untyped]


class WorkstreamStatus(Enum):
    """Workstream lifecycle status."""

    BACKLOG = "backlog"
    ACTIVE = "active"
    COMPLETED = "completed"
    BLOCKED = "blocked"


class WorkstreamSize(Enum):
    """Workstream scope size."""

    SMALL = "SMALL"
    MEDIUM = "MEDIUM"
    LARGE = "LARGE"


@dataclass
class AcceptanceCriterion:
    """Single acceptance criterion."""

    id: str
    description: str
    checked: bool = False


@dataclass
class Workstream:
    """Parsed workstream specification."""

    ws_id: str
    feature: str
    status: WorkstreamStatus
    size: WorkstreamSize
    github_issue: Optional[int] = None
    assignee: Optional[str] = None
    title: str = ""
    goal: str = ""
    acceptance_criteria: list[AcceptanceCriterion] = field(default_factory=list)
    context: str = ""
    dependencies: list[str] = field(default_factory=list)
    file_path: Optional[Path] = None


class WorkstreamParseError(Exception):
    """Error parsing workstream file."""

    pass


def parse_workstream(file_path: Path) -> Workstream:
    """Parse workstream markdown file.

    Args:
        file_path: Path to WS markdown file

    Returns:
        Parsed Workstream instance

    Raises:
        WorkstreamParseError: If file has no frontmatter or required fields missing
    """
    content = file_path.read_text(encoding="utf-8")
    frontmatter = _extract_frontmatter(content)

    ws_id: str = str(frontmatter["ws_id"])
    feature: str = str(frontmatter["feature"])
    status_str: str = str(frontmatter["status"])
    size_str: str = str(frontmatter["size"])

    try:
        status = WorkstreamStatus(status_str)
    except ValueError as e:
        raise WorkstreamParseError(f"Invalid status: {status_str}") from e

    try:
        size = WorkstreamSize(size_str)
    except ValueError as e:
        raise WorkstreamParseError(f"Invalid size: {size_str}") from e

    github_issue_val = frontmatter.get("github_issue")
    github_issue: Optional[int] = None
    if github_issue_val is not None:
        github_issue = int(github_issue_val)

    assignee_val = frontmatter.get("assignee")
    assignee: Optional[str] = None
    if assignee_val is not None:
        assignee = str(assignee_val)

    body = _strip_frontmatter(content)
    title = _extract_title(body)
    goal = _extract_section(body, "Goal")
    context = _extract_section(body, "Context")
    criteria = _extract_acceptance_criteria(body)
    deps = _extract_dependencies(body)

    return Workstream(
        ws_id=ws_id,
        feature=feature,
        status=status,
        size=size,
        github_issue=github_issue,
        assignee=assignee,
        title=title,
        goal=goal,
        acceptance_criteria=criteria,
        context=context,
        dependencies=deps,
        file_path=file_path,
    )


def _extract_frontmatter(content: str) -> dict[str, Any]:
    """Extract YAML frontmatter from markdown."""
    match = re.match(r"^---\n(.*?)\n---", content, re.DOTALL)
    if not match:
        raise WorkstreamParseError("No frontmatter found (must start with ---)")

    frontmatter_text = match.group(1)
    try:
        data: Any = yaml.safe_load(frontmatter_text)
    except yaml.YAMLError as e:
        raise WorkstreamParseError(f"Invalid YAML in frontmatter: {e}")

    if not isinstance(data, dict):
        raise WorkstreamParseError("Frontmatter must be a YAML dict")

    required_fields = {"ws_id", "feature", "status", "size"}
    missing = required_fields - set(data.keys())
    if missing:
        raise WorkstreamParseError(f"Missing required fields: {missing}")

    return data


def _strip_frontmatter(content: str) -> str:
    """Remove frontmatter, return body only."""
    match = re.match(r"^---\n.*?\n---\n(.*)", content, re.DOTALL)
    return match.group(1) if match else content


def _extract_title(body: str) -> str:
    """Extract title from first ## heading."""
    match = re.search(r"^## (.+)$", body, re.MULTILINE)
    return match.group(1) if match else ""


def _extract_section(body: str, section_name: str) -> str:
    """Extract content of a ### section by name (case-insensitive)."""
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


def _extract_acceptance_criteria(body: str) -> list[AcceptanceCriterion]:
    """Extract acceptance criteria from body."""
    criteria: list[AcceptanceCriterion] = []
    pattern = r"- \[([ x])\] (AC\d+): (.+)"
    for match in re.finditer(pattern, body, re.IGNORECASE):
        checked_char = match.group(1)
        ac_id = match.group(2)
        description = match.group(3)
        checked = checked_char.lower() == "x"
        criteria.append(AcceptanceCriterion(id=ac_id, description=description, checked=checked))
    return criteria


def _extract_dependencies(body: str) -> list[str]:
    """Extract WS dependencies from Dependencies section."""
    dep_section = _extract_section(body, "Dependencies")
    if not dep_section or dep_section.lower() in ("none", ""):
        return []
    pattern = r"WS-\d+-\d+"
    return re.findall(pattern, dep_section)
