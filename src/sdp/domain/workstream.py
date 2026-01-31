"""Pure domain entities for workstreams.

This module contains pure business logic entities with no external dependencies.
These can be imported by any layer (core, beads, unified, etc).
"""

import re
from dataclasses import dataclass, field
from enum import Enum
from pathlib import Path
from typing import Optional


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


@dataclass(frozen=True)
class WorkstreamID:
    """Parsed workstream ID in PP-FFF-SS format.

    Format: PP-FFF-SS where:
    - PP = Project ID (00-99), e.g., 00=SDP, 02=hw_checker, 03=mlsd, 04=bdde, 05=meta
    - FFF = Feature ID (000-999)
    - SS = Workstream sequence (00-99)

    This is a value object - immutable after creation.
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
    """Core workstream entity with business logic.

    This is the central domain entity representing a unit of work.
    Contains no I/O logic or external dependencies.
    """

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
