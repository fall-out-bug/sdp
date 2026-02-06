"""
Data models for TeamManager role registry.

Defines Role and RoleState for managing agent roles.
"""

from dataclasses import dataclass, field
from enum import Enum


class RoleState(Enum):
    """Role activation state."""

    DORMANT = "dormant"
    ACTIVE = "active"


@dataclass
class Role:
    """Agent role definition.

    Attributes:
        name: Unique role identifier (e.g., "planner", "builder")
        description: Human-readable role description
        state: Current activation state (active/dormant)
        skill_file: Path to skill definition file
        metadata: Optional additional role metadata
    """

    name: str
    description: str
    state: RoleState
    skill_file: str
    metadata: dict[str, str] = field(default_factory=dict)

    def to_dict(self) -> dict[str, str | dict[str, str]]:
        """Convert role to dictionary for JSON serialization."""
        return {
            "name": self.name,
            "description": self.description,
            "state": self.state.value,
            "skill_file": self.skill_file,
            "metadata": self.metadata,
        }

    @classmethod
    def from_dict(cls, data: dict[str, str | dict[str, str]]) -> "Role":
        """Create Role from dictionary (JSON deserialization)."""
        metadata_value = data.get("metadata", {})
        metadata = dict(metadata_value) if isinstance(metadata_value, dict) else {}

        return cls(
            name=str(data["name"]),
            description=str(data["description"]),
            state=RoleState(data["state"]),
            skill_file=str(data["skill_file"]),
            metadata=metadata,
        )
