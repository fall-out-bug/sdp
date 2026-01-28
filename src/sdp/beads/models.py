"""
Beads task data models.

Defines the data structures for interacting with Beads issues/tasks.
Compatible with Beads JSONL format and CLI output.
"""

from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum
from typing import Optional, List
import json


class BeadsStatus(str, Enum):
    """Beads task status values."""

    OPEN = "open"
    IN_PROGRESS = "in_progress"
    BLOCKED = "blocked"
    DEFERRED = "deferred"
    CLOSED = "closed"
    TOMBSTONE = "tombstone"
    PINNED = "pinned"
    HOOKED = "hooked"


class BeadsPriority(int, Enum):
    """Beads priority levels (0=critical, 4=backlog)."""

    CRITICAL = 0
    HIGH = 1
    MEDIUM = 2
    LOW = 3
    BACKLOG = 4


class BeadsDependencyType(str, Enum):
    """Beads dependency relationship types."""

    BLOCKS = "blocks"
    PARENT_CHILD = "parent-child"
    RELATED = "related"
    DISCOVERED_FROM = "discovered-from"


@dataclass
class BeadsDependency:
    """A dependency relationship between tasks."""

    task_id: str
    type: BeadsDependencyType


@dataclass
class BeadsTask:
    """A Beads task/issue.

    Compatible with Beads JSONL format and CLI output.
    """

    # Core fields
    id: str  # Hash-based ID like "bd-a3f8"
    title: str
    description: Optional[str] = None
    status: BeadsStatus = BeadsStatus.OPEN
    priority: BeadsPriority = BeadsPriority.MEDIUM

    # Relationships
    parent_id: Optional[str] = None  # For sub-tasks
    dependencies: List[BeadsDependency] = field(default_factory=list)

    # Metadata
    external_ref: Optional[str] = None  # Reference to external system (e.g., SDP WS ID)
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None

    # SDP-specific fields (stored in metadata)
    sdp_metadata: dict = field(default_factory=dict)

    def to_dict(self) -> dict:
        """Convert to dictionary for JSON serialization."""
        return {
            "id": self.id,
            "title": self.title,
            "description": self.description,
            "status": self.status.value,
            "priority": self.priority.value,
            "parent_id": self.parent_id,
            "dependencies": [
                {"task_id": d.task_id, "type": d.type.value} for d in self.dependencies
            ],
            "external_ref": self.external_ref,
            "created_at": self.created_at.isoformat() if self.created_at else None,
            "updated_at": self.updated_at.isoformat() if self.updated_at else None,
            "metadata": {"sdp": self.sdp_metadata} if self.sdp_metadata else None,
        }

    @classmethod
    def from_dict(cls, data: dict) -> "BeadsTask":
        """Create from dictionary (JSON deserialization)."""
        dependencies = [
            BeadsDependency(
                task_id=d["task_id"], type=BeadsDependencyType(d["type"])
            )
            for d in data.get("dependencies", [])
        ]

        sdp_metadata = data.get("metadata", {}).get("sdp", {})

        return cls(
            id=data["id"],
            title=data["title"],
            description=data.get("description"),
            status=BeadsStatus(data.get("status", "open")),
            priority=BeadsPriority(data.get("priority", 2)),
            parent_id=data.get("parent_id"),
            dependencies=dependencies,
            external_ref=data.get("external_ref"),
            created_at=(
                datetime.fromisoformat(data["created_at"]) if data.get("created_at") else None
            ),
            updated_at=(
                datetime.fromisoformat(data["updated_at"]) if data.get("updated_at") else None
            ),
            sdp_metadata=sdp_metadata,
        )

    def to_jsonl(self) -> str:
        """Convert to JSONL string (one line)."""
        return json.dumps(self.to_dict())

    @classmethod
    def from_jsonl(cls, line: str) -> "BeadsTask":
        """Create from JSONL line."""
        return cls.from_dict(json.loads(line))

    def is_ready(self) -> bool:
        """Check if task is ready to work on (no open blockers)."""
        # Task is ready if:
        # - Status is OPEN or IN_PROGRESS
        # - No blocking dependencies are open
        if self.status not in (BeadsStatus.OPEN, BeadsStatus.IN_PROGRESS):
            return False

        # Check dependencies (should be checked externally in real implementation)
        return True


@dataclass
class BeadsTaskCreate:
    """Parameters for creating a new Beads task."""

    title: str
    description: Optional[str] = None
    priority: BeadsPriority = BeadsPriority.MEDIUM
    parent_id: Optional[str] = None
    dependencies: List[BeadsDependency] = field(default_factory=list)
    external_ref: Optional[str] = None
    sdp_metadata: dict = field(default_factory=dict)

    def to_dict(self) -> dict:
        """Convert to dictionary for CLI/API call."""
        return {
            "title": self.title,
            "description": self.description,
            "priority": self.priority.value,
            "parent_id": self.parent_id,
            "dependencies": [
                {"task_id": d.task_id, "type": d.type.value} for d in self.dependencies
            ],
            "external_ref": self.external_ref,
            "metadata": {"sdp": self.sdp_metadata} if self.sdp_metadata else None,
        }


@dataclass
class BeadsSyncResult:
    """Result of a sync operation."""

    success: bool
    task_id: str
    beads_id: Optional[str] = None
    message: Optional[str] = None
    error: Optional[str] = None
