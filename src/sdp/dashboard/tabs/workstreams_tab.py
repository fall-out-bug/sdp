"""Workstreams tab for dashboard."""

from textual.widgets import Static

from ..state import DashboardState


class WorkstreamsTab(Static):
    """Tab showing workstreams by status."""

    def __init__(self, state_bus, **kwargs) -> None:
        """Initialize workstreams tab.

        Args:
            state_bus: StateBus for receiving updates
            **kwargs: Additional arguments for Static
        """
        super().__init__("Workstreams\n\nLoading...", **kwargs)
        self._state_bus = state_bus
        self._state_bus.subscribe(self._on_state_update)
        self._update_state(state_bus.state or DashboardState())

    def _on_state_update(self, state: DashboardState) -> None:
        """Update tab when state changes.

        Args:
            state: New dashboard state
        """
        self._update_state(state)

    def _update_state(self, state: DashboardState) -> None:
        """Update display with new state.

        Args:
            state: New dashboard state
        """
        if not state.workstreams:
            self.update("No workstreams found")
            return

        # Group by status
        by_status: dict[str, list] = {}
        for ws in state.workstreams.values():
            if ws.status not in by_status:
                by_status[ws.status] = []
            by_status[ws.status].append(ws)

        content = ""
        for status in ["backlog", "in-progress", "completed", "blocked"]:
            if status not in by_status:
                continue

            wss = by_status[status]
            status_label = status.replace("-", " ").title()
            content += f"\n{status_label} ({len(wss)})\n"

            for ws in sorted(wss, key=lambda w: w.ws_id):
                assignee = f" [{ws.assignee}]" if ws.assignee else ""
                size = f" [{ws.size}]" if ws.size else ""
                content += f"  {ws.ws_id}: {ws.title}{assignee}{size}\n"

        self.update(content.lstrip())
