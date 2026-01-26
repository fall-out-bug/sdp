"""File watcher for triggering test runs."""

import logging
import time
from pathlib import Path
from typing import Callable

from watchdog.events import FileSystemEventHandler, FileModifiedEvent
from watchdog.observers import Observer

logger = logging.getLogger(__name__)


class TestWatcher(FileSystemEventHandler):
    """Watches for file changes and triggers test runs."""

    def __init__(self, on_change: Callable[[str], None], debounce: float = 0.5) -> None:
        """Initialize watcher.

        Args:
            on_change: Callback function when a Python file changes
            debounce: Seconds to wait before triggering callback (debounce)
        """
        self._on_change = on_change
        self._debounce = debounce
        self._last_change = 0
        self._last_file = ""

    def on_modified(self, event: FileModifiedEvent) -> None:
        """Handle file modification event.

        Args:
            event: File modified event
        """
        if event.is_directory:
            return

        src_path = event.src_path
        if not src_path.endswith(".py"):
            return

        now = time.time()
        # Debounce: only trigger if enough time has passed or different file
        if (now - self._last_change > self._debounce) or (src_path != self._last_file):
            self._last_change = now
            self._last_file = src_path
            logger.info(f"File changed: {src_path}")
            self._on_change(src_path)


def watch_tests(
    project_dir: str | Path,
    callback: Callable[[str], None],
    debounce: float = 0.5,
) -> Observer:
    """Start watching for test file changes.

    Args:
        project_dir: Project directory to watch
        callback: Function to call when a Python file changes
        debounce: Debounce delay in seconds

    Returns:
        Observer instance for later stopping
    """
    observer = Observer()
    handler = TestWatcher(callback, debounce=debounce)

    project_path = Path(project_dir)

    # Watch src and tests directories
    for subdir in ["src", "tests"]:
        watch_path = project_path / subdir
        if watch_path.exists():
            observer.schedule(handler, str(watch_path), recursive=True)

    observer.start()
    logger.info(f"Started watching {project_dir} for file changes")
    return observer
