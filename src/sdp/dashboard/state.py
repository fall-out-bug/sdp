"""Dashboard state management."""

from dataclasses import dataclass, field
from datetime import datetime
from typing import Callable, Optional


@dataclass
class WorkstreamState:
    """State of a single workstream."""

    ws_id: str
    status: str  # backlog, in-progress, completed, blocked
    title: str
    feature: str
    assignee: Optional[str] = None
    size: Optional[str] = None
    started: Optional[datetime] = None
    completed: Optional[datetime] = None


@dataclass
class TestResults:
    """Results of test run."""

    status: str  # passed, failed, error, no_tests
    total: int = 0
    passed: int = 0
    failed: int = 0
    skipped: int = 0
    coverage: Optional[float] = None
    failed_tests: list[str] = field(default_factory=list)
    error_message: Optional[str] = None


@dataclass
class AgentEvent:
    """Event from agent execution."""

    timestamp: datetime
    event_type: str  # started, completed, error
    ws_id: Optional[str] = None
    message: str = ""


@dataclass
class DashboardState:
    """Complete dashboard state."""

    workstreams: dict[str, WorkstreamState] = field(default_factory=dict)
    test_results: Optional[TestResults] = None
    agent_activity: list[AgentEvent] = field(default_factory=list)
    last_update: Optional[datetime] = None


class StateBus:
    """Pub/sub for state updates."""

    def __init__(self) -> None:
        self._subscribers: list[Callable[[DashboardState], None]] = []
        self.state: Optional[DashboardState] = None

    def subscribe(self, callback: Callable[[DashboardState], None]) -> None:
        """Subscribe to state updates.

        Args:
            callback: Function to call when state updates
        """
        self._subscribers.append(callback)

    def publish(self, state: DashboardState) -> None:
        """Publish state update to all subscribers.

        Args:
            state: New dashboard state
        """
        self.state = state
        for callback in self._subscribers:
            callback(state)
