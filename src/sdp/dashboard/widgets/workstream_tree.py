"""Tree widget showing workstreams by status."""

from typing import TYPE_CHECKING

from textual.widgets import Tree

from ..state import DashboardState, WorkstreamState

if TYPE_CHECKING:
    from ..state import StateBus


class WorkstreamTree(Tree):
    """Tree widget showing workstreams by status."""

    def __init__(self, state_bus: "StateBus", **kwargs: object) -> None:
        """Initialize tree with state bus.

        Args:
            state_bus: StateBus for receiving updates
            **kwargs: Additional arguments for Tree widget
        """
        super().__init__(**kwargs)
        self._state_bus = state_bus
        self._state_bus.subscribe(self._on_state_update)
        self.root.set_label("Workstreams")

    def _on_state_update(self, state: DashboardState) -> None:
        """Update tree when state changes.

        Args:
            state: New dashboard state
        """
        self.clear()

        # Group by status
        by_status: dict[str, list[WorkstreamState]] = {}
        for ws in state.workstreams.values():
            if ws.status not in by_status:
                by_status[ws.status] = []
            by_status[ws.status].append(ws)

        # Build tree
        for status, wss in sorted(by_status.items()):
            status_label = f"{status.replace('-', ' ').title()} ({len(wss)})"
            status_node = self.root.add(status_label)
            for ws in sorted(wss, key=lambda w: w.ws_id):
                label = f"{ws.ws_id}: {ws.title}"
                if ws.assignee:
                    label += f" [{ws.assignee}]"
                status_node.add_leaf(label)

    def clear(self) -> None:
        """Clear all tree nodes."""
        self.root.remove_children()
