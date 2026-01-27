"""Queue state persistence."""

import logging
import json
from dataclasses import asdict
from datetime import datetime
from pathlib import Path
from typing import Any

from queue import PriorityQueue


logger = logging.getLogger(__name__)


class QueueState:
    """Manages queue state persistence."""

    @staticmethod
    def save(queue: PriorityQueue[object], state_file: str | Path) -> None:
        """Save queue state to file.

        Args:
            queue: PriorityQueue to save
            state_file: Path to state file
        """
        state_path = Path(state_file)
        state_path.parent.mkdir(parents=True, exist_ok=True)

        # Extract tasks from queue (not straightforward with PriorityQueue)
        # We'll store the items that haven't been processed yet
        items = []
        temp_list = []

        try:
            while not queue.empty():
                item = queue.get_nowait()
                items.append(item)
                temp_list.append(item)

            # Put items back
            for item in temp_list:
                queue.put_nowait(item)
        except Exception as e:
            logger.debug("Error extracting queue items: %s", e)

        # Convert to serializable format
        serializable = []
        for item in items:
            if hasattr(item, "__dict__"):
                data = {}
                for k, v in item.__dict__.items():
                    if isinstance(v, datetime):
                        data[k] = v.isoformat()
                    else:
                        data[k] = v
                serializable.append(data)

        state_path.write_text(json.dumps({"tasks": serializable}, indent=2))

    @staticmethod
    def load(state_file: str | Path) -> list[dict[str, Any]]:
        """Load queue state from file.

        Args:
            state_file: Path to state file

        Returns:
            List of task dicts
        """
        state_path = Path(state_file)
        if not state_path.exists():
            return []

        try:
            data = json.loads(state_path.read_text())
            return data.get("tasks", [])
        except Exception as e:
            logger.warning("Error loading queue state from %s: %s", state_path, e)
            return []

    @staticmethod
    def clear(state_file: str | Path) -> None:
        """Clear queue state file.

        Args:
            state_file: Path to state file
        """
        state_path = Path(state_file)
        if state_path.exists():
            state_path.unlink()
