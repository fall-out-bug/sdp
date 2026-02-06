"""Data models for status command."""

from dataclasses import dataclass, field
from typing import Optional


@dataclass
class WorkstreamSummary:
    """Summary of a workstream for status display."""

    id: str
    title: str
    status: str
    scope: str
    blockers: list[str] = field(default_factory=list)


@dataclass
class GuardStatus:
    """Guard state summary."""

    active: bool
    workstream_id: Optional[str] = None
    allowed_files: list[str] = field(default_factory=list)


@dataclass
class BeadsStatus:
    """Beads integration status."""

    available: bool
    synced: bool
    ready_tasks: list[str] = field(default_factory=list)
    last_sync: Optional[str] = None


@dataclass
class ProjectStatus:
    """Complete project status."""

    # Workstreams
    in_progress: list[WorkstreamSummary]
    blocked: list[WorkstreamSummary]
    ready: list[WorkstreamSummary]

    # Integrations
    guard: GuardStatus
    beads: BeadsStatus

    # Suggestions
    next_actions: list[str]
