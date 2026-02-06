"""Workstream tracking with Beads synchronization."""

import json
from datetime import datetime, timezone
from pathlib import Path

from sdp.beads import BeadsClient
from sdp.beads.models import BeadsStatus
from sdp.errors import ErrorCategory, SDPError


class WorkstreamInProgressError(SDPError):
    """Another WS is already in progress."""

    def __init__(self, current_ws: str, requested_ws: str) -> None:
        """Initialize error.

        Args:
            current_ws: Currently active workstream
            requested_ws: Requested workstream to activate
        """
        super().__init__(
            category=ErrorCategory.VALIDATION,
            message=f"WS {current_ws} already in progress",
            remediation=(
                f"1. Complete current WS: sdp guard complete {current_ws}\n"
                f"2. Or abort: sdp guard abort {current_ws}\n"
                f"3. Then activate: sdp guard activate {requested_ws}"
            ),
            context={"current_ws": current_ws, "requested_ws": requested_ws},
        )


class WorkstreamTracker:
    """Track active workstream with Beads sync."""

    def __init__(
        self, client: BeadsClient, state_file: Path = Path(".sdp/state.json")
    ) -> None:
        """Initialize tracker.

        Args:
            client: Beads client instance
            state_file: Path to state file
        """
        self._client = client
        self._state_file = state_file

    def get_active(self) -> str | None:
        """Get currently active WS ID.

        Returns:
            Active workstream ID or None
        """
        if not self._state_file.exists():
            return None
        with open(self._state_file) as f:
            state = json.load(f)
        active = state.get("active_ws")
        return active if isinstance(active, str) else None

    def activate(self, ws_id: str) -> None:
        """Activate WS and update Beads status.

        Args:
            ws_id: Workstream ID to activate

        Raises:
            WorkstreamInProgressError: If another WS is already active
        """
        # Check no other WS is active
        current = self.get_active()
        if current and current != ws_id:
            raise WorkstreamInProgressError(current, ws_id)

        # Update Beads status
        self._client.update_task_status(ws_id, BeadsStatus.IN_PROGRESS)

        # Get scope from WS metadata
        ws = self._client.get_task(ws_id)
        scope = ws.sdp_metadata.get("scope_files", []) if ws else []

        # Save local state
        self._save_state(
            {
                "active_ws": ws_id,
                "started_at": datetime.now(timezone.utc).isoformat(),
                "scope_files": scope,
            }
        )

    def complete(self, ws_id: str) -> None:
        """Mark WS as complete.

        Args:
            ws_id: Workstream ID to complete

        Raises:
            ValueError: If WS is not active
        """
        current = self.get_active()
        if current != ws_id:
            raise ValueError(f"WS {ws_id} is not active (active: {current})")

        self._client.update_task_status(ws_id, BeadsStatus.CLOSED)
        self._clear_state()

    def abort(self, ws_id: str) -> None:
        """Abort WS without completing.

        Args:
            ws_id: Workstream ID to abort

        Raises:
            ValueError: If WS is not active
        """
        current = self.get_active()
        if current != ws_id:
            raise ValueError(f"WS {ws_id} is not active")

        # Return to OPEN status
        self._client.update_task_status(ws_id, BeadsStatus.OPEN)
        self._clear_state()

    def _save_state(self, state: dict[str, str | list[str]]) -> None:
        """Save state to file.

        Args:
            state: State dictionary to save
        """
        self._state_file.parent.mkdir(exist_ok=True)
        with open(self._state_file, "w") as f:
            json.dump(state, f, indent=2)

    def _clear_state(self) -> None:
        """Clear state."""
        self._save_state({"active_ws": None})  # type: ignore[dict-item]
