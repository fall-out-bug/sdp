"""Main SDP Dashboard application."""

from datetime import datetime

from textual.app import App, ComposeResult
from textual.binding import Binding
from textual.widgets import Tabs, TabPane

from .state import DashboardState, StateBus
from .sources.agent_reader import AgentReader
from .sources.test_runner import TestRunner
from .sources.workstream_reader import WorkstreamReader
from .tabs.tests_tab import TestsTab
from .tabs.activity_tab import ActivityTab
from .tabs.workstreams_tab import WorkstreamsTab


class DashboardApp(App):
    """Main SDP Dashboard application."""

    TITLE = "SDP Dashboard"
    BINDINGS = [
        Binding("q", "quit", "Quit"),
        Binding("1", "switch_tab(1)", "Workstreams"),
        Binding("2", "switch_tab(2)", "Tests"),
        Binding("3", "switch_tab(3)", "Activity"),
        Binding("w", "switch_tab(1)", "Workstreams"),
        Binding("t", "switch_tab(2)", "Tests"),
        Binding("a", "switch_tab(3)", "Activity"),
        Binding("r", "refresh", "Refresh"),
    ]

    def __init__(self) -> None:
        """Initialize dashboard."""
        super().__init__()
        self._state_bus = StateBus()
        self._ws_reader = WorkstreamReader("docs/workstreams")
        self._test_runner = TestRunner(".")
        self._agent_reader = AgentReader(".")
        self._polling_active = True

    def compose(self) -> ComposeResult:
        """Compose dashboard UI."""
        yield Tabs("Workstreams", "Tests", "Activity")
        yield TabPane(
            WorkstreamsTab(self._state_bus, id="ws"),
            TestsTab(self._state_bus, id="tests"),
            ActivityTab(self._state_bus, id="activity"),
        )

    def on_mount(self) -> None:
        """Start background updates."""
        self._start_ws_polling()
        self._start_test_watching()

    def _start_ws_polling(self) -> None:
        """Poll workstream state every 1s."""
        async def poll() -> None:
            import asyncio

            while self._polling_active:
                try:
                    state = self._ws_reader.read()
                    state.last_update = datetime.now()
                    self._state_bus.publish(state)
                except Exception as e:
                    # Don't crash on polling errors
                    pass
                await asyncio.sleep(1)

        self.run_worker(poll())

    def _start_test_watching(self) -> None:
        """Start test file watcher (simplified)."""
        # For now, just run tests once on mount
        # Full implementation would use the test_watch module
        try:
            results = self._test_runner.run()
            current_state = self._state_bus.state or DashboardState()
            current_state.test_results = results
            self._state_bus.publish(current_state)
        except Exception:
            pass

    def action_switch_tab(self, tab_index: int) -> None:
        """Switch to tab by index.

        Args:
            tab_index: Tab index (1-based)
        """
        tabs = self.query_one(Tabs)
        if tabs:
            tabs.active = tab_index - 1

    def action_refresh(self) -> None:
        """Force refresh all data."""
        try:
            state = self._ws_reader.read()
            state.last_update = datetime.now()

            # Add test results
            if self._state_bus.state and self._state_bus.state.test_results:
                state.test_results = self._state_bus.state.test_results

            self._state_bus.publish(state)
        except Exception:
            pass

    def on_unmount(self) -> None:
        """Cleanup when dashboard closes."""
        self._polling_active = False
