"""Guard state management."""

import json
from dataclasses import asdict, dataclass
from pathlib import Path


@dataclass
class GuardState:
    """Guard state for persistence."""

    active_ws: str | None = None
    activated_at: str | None = None
    scope_files: list[str] | None = None


class StateManager:
    """Manages guard state persistence."""

    STATE_FILE = Path(".sdp/state.json")

    @classmethod
    def load(cls) -> GuardState:
        """Load state from file.

        Returns:
            GuardState instance
        """
        if not cls.STATE_FILE.exists():
            return GuardState()

        with open(cls.STATE_FILE) as f:
            data = json.load(f)
        return GuardState(**data)

    @classmethod
    def save(cls, state: GuardState) -> None:
        """Save state to file.

        Args:
            state: State to persist
        """
        cls.STATE_FILE.parent.mkdir(exist_ok=True)
        with open(cls.STATE_FILE, "w") as f:
            json.dump(asdict(state), f, indent=2)

    @classmethod
    def clear(cls) -> None:
        """Clear state."""
        cls.save(GuardState())
