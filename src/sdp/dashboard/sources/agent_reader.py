"""Agent activity reader for dashboard."""

import json
import logging
from datetime import datetime
from pathlib import Path

from ..state import AgentEvent

logger = logging.getLogger(__name__)


class AgentReader:
    """Reads agent activity from daemon queue."""

    def __init__(self, project_dir: str | Path) -> None:
        """Initialize reader with project directory.

        Args:
            project_dir: Project root directory
        """
        self._project = Path(project_dir)
        self._queue_file = self._project / ".sdp" / "daemon_queue.json"

    def read(self) -> list[AgentEvent]:
        """Read agent events from queue file.

        Returns:
            List of AgentEvent objects, empty if file doesn't exist
        """
        if not self._queue_file.exists():
            return []

        try:
            data = json.loads(self._queue_file.read_text())
            events = []
            for e in data.get("events", []):
                ts = e.get("timestamp")
                if isinstance(ts, str):
                    timestamp = datetime.fromisoformat(ts)
                elif isinstance(ts, datetime):
                    timestamp = ts
                else:
                    timestamp = datetime.now()

                events.append(
                    AgentEvent(
                        timestamp=timestamp,
                        event_type=e.get("event_type", "unknown"),
                        ws_id=e.get("ws_id"),
                        message=e.get("message", ""),
                    )
                )
            return events
        except Exception as e:
            logger.warning(f"Failed to read agent queue: {e}")
            return []
