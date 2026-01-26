"""Tests tab for dashboard."""

from textual.widgets import Static

from ..state import DashboardState


class TestsTab(Static):
    """Tab showing test results."""

    def __init__(self, state_bus, **kwargs) -> None:
        """Initialize tests tab.

        Args:
            state_bus: StateBus for receiving updates
            **kwargs: Additional arguments for Static
        """
        super().__init__("Test Results\n\nLoading...", **kwargs)
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
        results = state.test_results

        if not results:
            content = "No tests run yet"
        elif results.status == "no_tests":
            content = "No tests found in project"
        else:
            status_emoji = {
                "passed": "✅",
                "failed": "❌",
                "error": "⚠️",
                "no_tests": "⚪",
            }

            coverage_str = f"{results.coverage:.1f}%" if results.coverage else "N/A"

            content = f"""{status_emoji} Test Results

Status:    {results.status.upper()}
Passed:    {results.passed}
Failed:    {results.failed}
Coverage:  {coverage_str}
"""

            if results.failed_tests:
                content += "\n\nFailed Tests:\n"
                for test in results.failed_tests[:10]:
                    content += f"  - {test}\n"
                if len(results.failed_tests) > 10:
                    content += f"  ... and {len(results.failed_tests) - 10} more\n"

            if results.error_message:
                content += f"\n\nError: {results.error_message}\n"

        self.update(content)
