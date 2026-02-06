"""Pre-edit guard for WS scope enforcement."""

from sdp.beads import BeadsClient
from sdp.guard.models import GuardResult


class GuardSkill:
    """Pre-edit guard for WS scope enforcement."""

    def __init__(self, beads_client: BeadsClient) -> None:
        """Initialize guard with Beads client.

        Args:
            beads_client: Beads client instance
        """
        self._client = beads_client
        self._active_ws: str | None = None

    def activate(self, task_id: str) -> None:
        """Set active workstream.

        Args:
            task_id: Beads task ID (e.g., sdp-4qq) or workstream ID (PP-FFF-SS)

        Raises:
            ValueError: If task not found
        """
        ws = self._client.get_task(task_id)
        if not ws:
            raise ValueError(f"WS not found: {task_id}")
        self._active_ws = task_id

    def check_edit(self, file_path: str) -> GuardResult:
        """Check if file edit is allowed.

        Args:
            file_path: Path to file being edited

        Returns:
            GuardResult with allowed status and reason
        """
        if not self._active_ws:
            return GuardResult(
                allowed=False,
                ws_id=None,
                reason="No active WS. Run @build <ws_id> first.",
                scope_files=[],
            )

        ws = self._client.get_task(self._active_ws)
        scope = ws.sdp_metadata.get("scope_files", []) if ws else []

        if not scope:
            # No scope defined = all files allowed
            return GuardResult(
                allowed=True,
                ws_id=self._active_ws,
                reason="No scope restrictions",
                scope_files=[],
            )

        if file_path not in scope:
            return GuardResult(
                allowed=False,
                ws_id=self._active_ws,
                reason=f"File {file_path} not in WS scope",
                scope_files=scope,
            )

        return GuardResult(
            allowed=True, ws_id=self._active_ws, reason="File in scope", scope_files=scope
        )
