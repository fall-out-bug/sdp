"""Project map data types for SDP projects.

This module provides core dataclasses for project map representation.
"""

from dataclasses import dataclass, field
from pathlib import Path
from typing import Optional


@dataclass
class Decision:
    """Single architectural decision record.

    Attributes:
        area: Decision area (e.g., "Architecture", "Storage")
        decision: Decision description
        adr: ADR identifier (e.g., "ADR-001")
        date: Decision date (YYYY-MM-DD)
    """

    area: str
    decision: str
    adr: str
    date: str


@dataclass
class Constraint:
    """Single project constraint.

    Attributes:
        category: Constraint category (e.g., "AI-Readiness", "Clean Architecture")
        description: Constraint description text
    """

    category: str
    description: str


@dataclass
class TechStackItem:
    """Single tech stack entry.

    Attributes:
        layer: Layer name (e.g., "Language", "API")
        technology: Technology name (e.g., "Python 3.11+", "FastAPI")
    """

    layer: str
    technology: str


@dataclass
class ProjectMap:
    """Parsed project map specification.

    Attributes:
        project_name: Project name extracted from title
        decisions: List of architectural decisions
        constraints: List of active constraints
        current_state: Current state section content (raw markdown)
        patterns: Patterns & conventions section content (raw markdown)
        tech_stack: List of tech stack items
        file_path: Path to source PROJECT_MAP.md file
    """

    project_name: str
    decisions: list[Decision] = field(default_factory=list)
    constraints: list[Constraint] = field(default_factory=list)
    current_state: str = ""
    patterns: str = ""
    tech_stack: list[TechStackItem] = field(default_factory=list)
    file_path: Optional[Path] = None


class ProjectMapParseError(Exception):
    """Error parsing project map file."""

    pass


def get_decision(
    project_map: ProjectMap,
    *,
    area: Optional[str] = None,
    adr: Optional[str] = None,
) -> Optional[Decision]:
    """Query decision by area or ADR.

    Args:
        project_map: ProjectMap instance
        area: Decision area to search for
        adr: ADR identifier to search for

    Returns:
        Decision if found, None otherwise

    Raises:
        ValueError: If neither area nor adr is provided
    """
    if area is None and adr is None:
        raise ValueError("Must provide either 'area' or 'adr' parameter")

    for decision in project_map.decisions:
        if area is not None and decision.area == area:
            return decision
        if adr is not None and decision.adr == adr:
            return decision

    return None


def get_constraint(
    project_map: ProjectMap,
    *,
    category: Optional[str] = None,
    keyword: Optional[str] = None,
) -> list[Constraint]:
    """Query constraints by category or keyword.

    Args:
        project_map: ProjectMap instance
        category: Constraint category to filter by
        keyword: Keyword to search in constraint descriptions

    Returns:
        List of matching constraints (empty if none found)

    Raises:
        ValueError: If neither category nor keyword is provided
    """
    if category is None and keyword is None:
        raise ValueError("Must provide either 'category' or 'keyword' parameter")

    results: list[Constraint] = []

    for constraint in project_map.constraints:
        if category is not None and constraint.category == category:
            results.append(constraint)
        elif keyword is not None and keyword.lower() in constraint.description.lower():
            results.append(constraint)

    return results
