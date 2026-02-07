"""Status collector - gathers project status from various sources."""

import json
from pathlib import Path

from sdp.cli.status.models import (
    BeadsStatus,
    GuardStatus,
    ProjectStatus,
    WorkstreamSummary,
)
from sdp.cli.status.parser import parse_ws_file


class StatusCollector:
    """Collects project status from various sources."""

    def __init__(self, root: Path):
        """Initialize collector.

        Args:
            root: Project root directory
        """
        self.root = root
        self.ws_dir = root / "docs" / "workstreams"

    def collect(self) -> ProjectStatus:
        """Collect complete project status.

        Returns:
            ProjectStatus with all information gathered
        """
        return ProjectStatus(
            in_progress=self._collect_in_progress(),
            blocked=self._collect_blocked(),
            ready=self._collect_ready(),
            guard=self._collect_guard_status(),
            beads=self._collect_beads_status(),
            next_actions=self._suggest_actions(),
        )

    def _collect_in_progress(self) -> list[WorkstreamSummary]:
        """Find workstreams in progress.

        Returns:
            List of in-progress workstreams
        """
        ws_list = []
        in_progress_dir = self.ws_dir / "in_progress"
        if in_progress_dir.exists():
            for ws_file in in_progress_dir.glob("*.md"):
                ws = parse_ws_file(ws_file)
                if ws:
                    ws_list.append(ws)
        return ws_list

    def _collect_blocked(self) -> list[WorkstreamSummary]:
        """Find blocked workstreams.

        Returns:
            List of blocked workstreams with dependencies
        """
        ws_list = []
        backlog_dir = self.ws_dir / "backlog"
        if backlog_dir.exists():
            for ws_file in backlog_dir.glob("*.md"):
                ws = parse_ws_file(ws_file)
                if ws and ws.status == "BLOCKED":
                    ws_list.append(ws)
        return ws_list

    def _collect_ready(self) -> list[WorkstreamSummary]:
        """Find ready workstreams (no blockers).

        Returns:
            List of ready-to-start workstreams
        """
        ws_list = []
        backlog_dir = self.ws_dir / "backlog"
        if backlog_dir.exists():
            for ws_file in backlog_dir.glob("*.md"):
                ws = parse_ws_file(ws_file)
                if ws and ws.status == "READY":
                    ws_list.append(ws)
        return ws_list

    def _collect_guard_status(self) -> GuardStatus:
        """Get guard status.

        Returns:
            GuardStatus with current guard state
        """
        state_file = self.root / ".sdp" / "state.json"
        if state_file.exists():
            try:
                with open(state_file, "r") as f:
                    data = json.load(f)

                return GuardStatus(
                    active=True,
                    workstream_id=data.get("active_workstream"),
                    allowed_files=data.get("allowed_files", []),
                )
            except (json.JSONDecodeError, KeyError):
                pass

        return GuardStatus(active=False)

    def _collect_beads_status(self) -> BeadsStatus:
        """Get Beads integration status.

        Returns:
            BeadsStatus with Beads availability and sync info
        """
        beads_dir = self.root / ".beads"
        if not beads_dir.exists():
            return BeadsStatus(available=False, synced=False)

        # Check metadata file for sync status
        metadata_file = beads_dir / "metadata.json"
        if metadata_file.exists():
            try:
                with open(metadata_file, "r") as f:
                    data = json.load(f)

                return BeadsStatus(
                    available=True,
                    synced=True,
                    last_sync=data.get("last_sync"),
                )
            except (json.JSONDecodeError, KeyError):
                pass

        return BeadsStatus(available=True, synced=False)

    def _suggest_actions(self) -> list[str]:
        """Suggest next actions based on state.

        Returns:
            List of suggested action strings
        """
        actions = []

        # Collect current state
        in_progress = self._collect_in_progress()
        ready = self._collect_ready()
        blocked = self._collect_blocked()

        if in_progress:
            ws = in_progress[0]
            actions.append(f"Complete {ws.id} to continue progress")

            # Check if completing this would unblock others
            unblocked_count = sum(
                1 for b in blocked if ws.id in b.blockers
            )
            if unblocked_count > 0:
                actions.append(
                    f"Completing {ws.id} will unblock {unblocked_count} workstream(s)"
                )
        elif ready:
            ws = ready[0]
            actions.append(f"Start {ws.id}: {ws.title}")
        else:
            actions.append("No workstreams available - create new work with @design")

        # Beads sync suggestion
        beads = self._collect_beads_status()
        if beads.available and not beads.synced:
            actions.append("Run 'bd sync' to synchronize Beads")

        return actions
