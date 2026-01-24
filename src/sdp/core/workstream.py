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
class WorkstreamID:
    """Parsed workstream ID in PP-FFF-SS format.

    Format: PP-FFF-SS where:
    - PP = Project ID (00-99), e.g., 00=SDP, 02=hw_checker, 03=mlsd, 04=bdde, 05=meta
    - FFF = Feature ID (000-999)
    - SS = Workstream sequence (00-99)
    """

    project_id: int  # 00-99
    feature_id: int  # 000-999
    sequence: int  # 00-99

    def __str__(self) -> str:
        return f"{self.project_id:02d}-{self.feature_id:03d}-{self.sequence:02d}"

    @classmethod
    def parse(cls, ws_id: str) -> "WorkstreamID":
        """Parse WS ID string like '00-500-01' or 'WS-500-01' (legacy).

        Args:
            ws_id: Workstream ID string

        Returns:
            WorkstreamID instance

        Raises:
            ValueError: If format is invalid
        """
        # Support both formats: PP-FFF-SS (new) and WS-FFF-SS (legacy)
        # Legacy format assumes project_id 00 (SDP)
        pattern_legacy = r"^WS-(\d{3})-(\d{2})$"
        match_legacy = re.match(pattern_legacy, ws_id)
        if match_legacy:
            feature_id, sequence = match_legacy.groups()
            return cls(
                project_id=0,  # SDP
                feature_id=int(feature_id),
                sequence=int(sequence)
            )

        pattern = r"^(\d{2})-(\d{3})-(\d{2})$"
        match = re.match(pattern, ws_id)
        if not match:
            raise ValueError(
                f"Invalid WS ID format: {ws_id}. "
                f"Expected PP-FFF-SS (e.g., 00-500-01) or WS-FFF-SS (legacy)"
            )
        project_id, feature_id, sequence = match.groups()
        return cls(
            project_id=int(project_id),
            feature_id=int(feature_id),
            sequence=int(sequence)
        )

    @property
    def is_sdp(self) -> bool:
        """Check if this is an SDP Protocol workstream (Project 00)."""
        return self.project_id == 0

    @property
    def is_hw_checker(self) -> bool:
        """Check if this is a hw_checker workstream (Project 02)."""
        return self.project_id == 2

    @property
    def is_mlsd(self) -> bool:
        """Check if this is an MLSD course workstream (Project 03)."""
        return self.project_id == 3

    @property
    def is_bdde(self) -> bool:
        """Check if this is a BDDE course workstream (Project 04)."""
        return self.project_id == 4

    @property
    def is_meta_repo(self) -> bool:
        """Check if this is a meta-repo workstream (Project 05)."""
        return self.project_id == 5

    def validate_project_id(self, valid_ids: set[int] | None = None) -> None:
        """Validate project ID against known registry.

        Args:
            valid_ids: Set of valid project IDs. Defaults to {0, 2, 3, 4, 5}

        Raises:
            ValueError: If project_id is not in valid_ids
        """
        if valid_ids is None:
            valid_ids = {0, 2, 3, 4, 5}  # SDP, hw_checker, mlsd, bdde, meta

        if self.project_id not in valid_ids:
            raise ValueError(
                f"Invalid project_id: {self.project_id:02d}. "
                f"Valid IDs: {', '.join(f'{i:02d}' for i in sorted(valid_ids))}"
            )


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
    steps: list[str] = field(default_factory=list)
    code_blocks: list[str] = field(default_factory=list)
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
    steps = _extract_steps(body)
    code_blocks = _extract_code_blocks(body)

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
        steps=steps,
        code_blocks=code_blocks,
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


def _extract_steps(body: str) -> list[str]:
    """Extract numbered steps from Steps section."""
    steps_section = _extract_section(body, "Steps")
    steps: list[str] = []

    # Match patterns like "1. Step description" or "#### 1. Step"
    # Look for lines starting with a number followed by a dot
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


def _extract_code_blocks(body: str) -> list[str]:
    """Extract fenced code blocks from body."""
    # Match ```language\ncode```
    pattern = r"```[\w]*\n(.+?)```"
    return re.findall(pattern, body, re.DOTALL)
