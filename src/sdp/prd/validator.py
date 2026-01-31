"""PRD section validation.

This module validates PRD document sections against size and format limits.
"""

import re
from dataclasses import dataclass
from enum import Enum
from pathlib import Path

from .parser import parse_prd_sections
from .profiles import PROFILES, ProjectType


class Severity(Enum):
    """Validation issue severity."""
    WARNING = "warning"
    ERROR = "error"


@dataclass
class ValidationIssue:
    """A validation issue found in a PRD section.

    Attributes:
        section: Name of the section with the issue
        message: Human-readable description of the issue
        severity: Severity level (WARNING or ERROR)
        current: Current value (e.g., character count)
        limit: Limit that was exceeded
    """
    section: str
    message: str
    severity: Severity
    current: int
    limit: int


def validate_prd(
    content: str,
    project_type: ProjectType = ProjectType.SERVICE
) -> list[ValidationIssue]:
    """Validate PRD content against line limits.

    Args:
        content: PRD document content
        project_type: Type of project (affects which sections to validate)

    Returns:
        List of validation issues
    """
    issues = []

    # Parse sections
    sections = parse_prd_sections(content)

    # Get profile for this project type
    profile = PROFILES[project_type]

    # Validate each section
    for section_def in profile.sections:
        section_name = section_def.name
        section_content = sections.get(section_name, "")

        # Check character limit if defined
        if section_def.max_chars is not None:
            char_count = len(section_content.strip())
            if char_count > section_def.max_chars:
                issues.append(ValidationIssue(
                    section=section_name,
                    message=f"Превышен лимит символов: {char_count}/{section_def.max_chars}",
                    severity=Severity.WARNING,
                    current=char_count,
                    limit=section_def.max_chars
                ))

    # Special validation for "Модель БД" format
    if "Модель БД" in sections:
        issues.extend(_validate_db_section(sections["Модель БД"]))

    # Special validation for "Назначение" length
    if "Назначение" in sections:
        issues.extend(_validate_purpose_section(sections["Назначение"]))

    return issues


def _validate_db_section(db_section: str) -> list[ValidationIssue]:
    """Validate database model section format.

    Each field should be on a single line (max 120 chars per line).

    Args:
        db_section: Content of the database section

    Returns:
        List of validation issues
    """
    issues = []
    max_line_length = 120

    for line_num, line in enumerate(db_section.strip().split("\n"), start=1):
        # Skip headers and empty lines
        if line.startswith("#") or not line.strip():
            continue

        # Each field should be on single line
        if len(line) > max_line_length:
            issues.append(ValidationIssue(
                section="Модель БД",
                message=f"Поле превышает {max_line_length} символов (строка {line_num}): {line[:50]}...",  # noqa: E501
                severity=Severity.ERROR,
                current=len(line),
                limit=max_line_length
            ))

    return issues


def _validate_purpose_section(purpose_section: str) -> list[ValidationIssue]:
    """Validate purpose section length.

    Args:
        purpose_section: Content of the purpose section

    Returns:
        List of validation issues
    """
    issues = []
    max_chars = 500
    text = purpose_section.strip()

    # Remove markdown formatting from character count
    text_clean = re.sub(r'[#*`\-\[\]]', '', text)
    char_count = len(text_clean)

    if char_count > max_chars:
        issues.append(ValidationIssue(
            section="Назначение",
            message=f"Превышен лимит: {char_count}/{max_chars} символов",
            severity=Severity.WARNING,
            current=char_count,
            limit=max_chars
        ))

    return issues


def validate_prd_file(prd_path: Path) -> list[ValidationIssue]:
    """Validate a PRD file.

    Args:
        prd_path: Path to the PRD file

    Returns:
        List of validation issues
    """
    try:
        content = prd_path.read_text()

        # Detect project type from frontmatter
        project_type = _detect_project_type_from_frontmatter(content)

        return validate_prd(content, project_type)
    except Exception:
        return []


def _detect_project_type_from_frontmatter(content: str) -> ProjectType:
    """Detect project type from PRD frontmatter.

    Args:
        content: PRD document content

    Returns:
        Detected project type (defaults to SERVICE)
    """
    match = re.search(r'project_type:\s*(\w+)', content)
    if match:
        type_str = match.group(1)
        try:
            return ProjectType(type_str)
        except ValueError:
            pass
    return ProjectType.SERVICE


def format_validation_issues(issues: list[ValidationIssue]) -> str:
    """Format validation issues for display.

    Args:
        issues: List of validation issues

    Returns:
        Formatted string for display
    """
    if not issues:
        return "✅ PRD validation passed"

    lines = ["❌ PRD validation failed:\n"]

    for issue in issues:
        symbol = "⚠️" if issue.severity == Severity.WARNING else "❌"
        lines.append(
            f"  {symbol} [{issue.section}] {issue.message}\n"
            f"     Current: {issue.current}, Limit: {issue.limit}"
        )

    return "\n".join(lines)


def has_critical_issues(issues: list[ValidationIssue]) -> bool:
    """Check if there are any ERROR severity issues.

    Args:
        issues: List of validation issues

    Returns:
        True if any ERROR issues found
    """
    return any(issue.severity == Severity.ERROR for issue in issues)
