"""Log widget showing agent activity."""

from datetime import datetime

from textual.widgets import Log

from ..state import AgentEvent, DashboardState


class ActivityLog(Log):
    """Log widget showing agent activity."""

    def __init__(self, state_bus: "StateBus", **kwargs: object) -> None:
        """Initialize log with state bus.

        Args:
            state_bus: StateBus for receiving updates
            **kwargs: Additional arguments for Log widget
        """
        super().__init__(**kwargs)
        self._state_bus = state_bus
        self._state_bus.subscribe(self._on_state_update)
        self._seen_events: set[int] = set()

    def _on_state_update(self, state: DashboardState) -> None:
        """Add new events to log.

        Args:
            state: New dashboard state
        """
        for event in state.agent_activity:
            # Create a simple hash to track seen events
            event_key = hash((event.timestamp, event.event_type, event.ws_id or "", event.message))
            if event_key not in self._seen_events:
                self._seen_events.add(event_key)

                # Format timestamp
                ts_str = event.timestamp.strftime("%H:%M:%S") if isinstance(event.timestamp, datetime) else str(event.timestamp)

                # Format event
                ws_part = f"[{event.ws_id}] " if event.ws_id else ""
                self.write_line(f"{ts_str} {ws_part}{event.event_type}: {event.message}")
