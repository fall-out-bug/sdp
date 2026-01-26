"""Activity tab for dashboard."""

from datetime import datetime

from textual.widgets import Static

from ..state import AgentEvent, DashboardState


class ActivityTab(Static):
    """Tab showing agent activity."""

    def __init__(self, state_bus, **kwargs) -> None:
        """Initialize activity tab.

        Args:
            state_bus: StateBus for receiving updates
            **kwargs: Additional arguments for Static
        """
        super().__init__("Activity Log\n\nLoading...", **kwargs)
        self._state_bus = state_bus
        self._state_bus.subscribe(self._on_state_update)
        self._update_state(state_bus.state or DashboardState())
        self._seen_events = set()

    def _on_state_update(self, state: DashboardState) -> None:
        """Add new events to log.

        Args:
            state: New dashboard state
        """
        # Only show new events
        for event in state.agent_activity:
            event_key = (
                event.timestamp.isoformat()
                if isinstance(event.timestamp, datetime)
                else str(event.timestamp),
                event.event_type,
                event.ws_id or "",
                event.message,
            )
            if event_key not in self._seen_events:
                self._seen_events.add(event_key)
                self._add_event(event)

    def _update_state(self, state: DashboardState) -> None:
        """Update display with current state.

        Args:
            state: New dashboard state
        """
        if not state.agent_activity:
            self.update("No activity yet")
            return

        content = "Activity Log\n\n"
        for event in state.agent_activity[-20:]:  # Show last 20
            ts_str = (
                event.timestamp.strftime("%H:%M:%S")
                if isinstance(event.timestamp, datetime)
                else str(event.timestamp)
            )
            ws_part = f"[{event.ws_id}] " if event.ws_id else ""
            content += f"{ts_str} {ws_part}{event.event_type}: {event.message}\n"

        self.update(content)

    def _add_event(self, event: AgentEvent) -> None:
        """Add event to display.

        Args:
            event: Agent event to add
        """
        current = self.render_plain(self._get_rich_content() or "")
        ts_str = (
            event.timestamp.strftime("%H:%M:%S")
            if isinstance(event.timestamp, datetime)
            else str(event.timestamp)
        )
        ws_part = f"[{event.ws_id}] " if event.ws_id else ""
        new_line = f"{ts_str} {ws_part}{event.event_type}: {event.message}\n"

        self.update(current + new_line)
