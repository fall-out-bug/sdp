"""Bidirectional synchronization between Beads and local state.

Provides sync service for keeping Beads task status and local guard state
in sync with conflict detection and resolution.
"""

from dataclasses import dataclass
from enum import Enum
from pathlib import Path
from typing import List, Optional

from sdp.beads.base import BeadsClient
from sdp.beads.models import BeadsStatus


class SyncSource(str, Enum):
    """Source of truth for sync resolution."""

    LOCAL = "local"
    BEADS = "beads"


@dataclass
class SyncConflict:
    """Detected sync conflict."""

    ws_id: str
    local_status: Optional[str]
    beads_status: str
    field: str  # 'existence', 'status', etc.


@dataclass
class SyncResult:
    """Result of sync operation."""

    synced: bool
    conflicts: List[SyncConflict]
    changes: List[str]


class BeadsSyncService:
    """Synchronize local state with Beads.

    Handles bidirectional sync between:
    - Local guard state (active workstream tracking)
    - Beads task status

    Detects conflicts and provides resolution strategies.
    """

    def __init__(self, client: BeadsClient):
        """Initialize sync service.

        Args:
            client: Beads client for task operations
        """
        self._client = client

    def check_sync(self, active_ws: Optional[str] = None) -> SyncResult:
        """Check if local state matches Beads.

        Args:
            active_ws: Active workstream ID from local state

        Returns:
            SyncResult with detected conflicts
        """
        conflicts: List[SyncConflict] = []

        if not active_ws:
            # No active workstream locally - in sync
            return SyncResult(synced=True, conflicts=[], changes=[])

        # Check if task exists in Beads
        beads_task = self._client.get_task(active_ws)
        if not beads_task:
            conflicts.append(
                SyncConflict(
                    ws_id=active_ws,
                    local_status="active",
                    beads_status="not_found",
                    field="existence",
                )
            )
            return SyncResult(synced=False, conflicts=conflicts, changes=[])

        # Check status match
        if beads_task.status != BeadsStatus.IN_PROGRESS:
            conflicts.append(
                SyncConflict(
                    ws_id=active_ws,
                    local_status="active",
                    beads_status=beads_task.status.value,
                    field="status",
                )
            )

        return SyncResult(
            synced=len(conflicts) == 0,
            conflicts=conflicts,
            changes=[],
        )

    def sync(
        self,
        active_ws: Optional[str] = None,
        source: SyncSource = SyncSource.BEADS,
    ) -> SyncResult:
        """Sync state from specified source.

        Args:
            active_ws: Active workstream ID from local state
            source: Which side is source of truth

        Returns:
            SyncResult with applied changes
        """
        check = self.check_sync(active_ws)

        if check.synced:
            return check

        changes: List[str] = []

        for conflict in check.conflicts:
            if source == SyncSource.BEADS:
                # Update local from Beads
                if conflict.beads_status == "not_found":
                    self._clear_local_state()
                    changes.append(
                        f"Cleared local (WS {conflict.ws_id} not in Beads)"
                    )
                elif conflict.beads_status != "in_progress":
                    self._clear_local_state()
                    changes.append(
                        f"Cleared local (Beads status: {conflict.beads_status})"
                    )
            else:
                # Update Beads from local
                if conflict.field == "status":
                    self._client.update_task_status(
                        conflict.ws_id,
                        BeadsStatus.IN_PROGRESS,
                    )
                    changes.append(f"Updated Beads {conflict.ws_id} to IN_PROGRESS")

        return SyncResult(synced=True, conflicts=[], changes=changes)

    def _clear_local_state(self) -> None:
        """Clear local guard state file.

        Removes the .guard_state file that tracks active workstream.
        """
        state_file = Path.cwd() / ".guard_state"
        if state_file.exists():
            state_file.unlink()

    def _get_local_active_ws(self) -> Optional[str]:
        """Get active workstream from local state.

        Returns:
            Active workstream ID or None
        """
        state_file = Path.cwd() / ".guard_state"
        if not state_file.exists():
            return None

        try:
            content = state_file.read_text().strip()
            if content:
                # Parse simple format: active_ws=<id>
                if "=" in content:
                    return content.split("=", 1)[1]
                return content
        except (OSError, ValueError):
            return None

        return None
