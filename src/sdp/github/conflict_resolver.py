"""Conflict detection and resolution for GitHub sync."""

from dataclasses import dataclass


@dataclass(frozen=True)
class Conflict:
    """Represents a conflict between WS and GitHub state."""

    field: str
    ws_value: object
    gh_value: object


class ConflictResolver:
    """Detects and resolves conflicts between WS and GitHub state."""

    # Status mapping between WS and GitHub
    _STATUS_MAP: dict[str, set[str]] = {
        "backlog": {"backlog", "Backlog"},
        "in-progress": {"in-progress", "in_progress", "In Progress", "In progress"},
        "completed": {"completed", "Completed", "Done", "done"},
        "blocked": {"blocked", "Blocked"},
    }

    def detect(self, ws_state: dict, gh_state: dict) -> Conflict | None:
        """Detect if there's a conflict between states.

        Args:
            ws_state: Workstream state from file
            gh_state: State from GitHub

        Returns:
            Conflict if detected, None otherwise
        """
        ws_status = ws_state.get("status")
        gh_status = gh_state.get("status")

        if not ws_status or not gh_status:
            return None

        # Normalize GitHub status to match WS format
        normalized_gh = self._normalize_status(gh_status)

        if ws_status != normalized_gh:
            return Conflict(
                field="status",
                ws_value=ws_status,
                gh_value=gh_status,
            )

        return None

    def resolve(self, conflict: Conflict) -> object:
        """Resolve conflict: WS file always wins (source of truth).

        Args:
            conflict: Conflict to resolve

        Returns:
            Resolved value (WS value)
        """
        return conflict.ws_value

    def _normalize_status(self, gh_status: str) -> str:
        """Normalize GitHub status to WS format.

        Args:
            gh_status: Status from GitHub

        Returns:
            Normalized status string
        """
        gh_lower = gh_status.lower().replace(" ", "_").replace("-", "_")

        for ws_status, variants in self._STATUS_MAP.items():
            if gh_lower in {v.lower().replace(" ", "_").replace("-", "_") for v in variants}:
                return ws_status

        return gh_lower
