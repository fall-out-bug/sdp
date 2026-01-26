"""Workstream file reader for dashboard."""

import logging
from pathlib import Path

from ..state import DashboardState, WorkstreamState

logger = logging.getLogger(__name__)


class WorkstreamReader:
    """Reads workstream files and builds dashboard state."""

    def __init__(self, base_dir: str | Path) -> None:
        """Initialize reader with base directory.

        Args:
            base_dir: Base directory containing workstream subdirectories
        """
        self._base = Path(base_dir)

    def read(self) -> DashboardState:
        """Read all workstreams and build state.

        Returns:
            DashboardState with all discovered workstreams
        """
        workstreams: dict[str, WorkstreamState] = {}

        for status_dir in ["backlog", "in_progress", "completed"]:
            ws_dir = self._base / status_dir
            if not ws_dir.exists():
                continue

            for ws_file in ws_dir.glob("*.md"):
                try:
                    ws_state = self._parse_ws_file(ws_file, status_dir)
                    workstreams[ws_state.ws_id] = ws_state
                except Exception as e:
                    logger.warning(f"Failed to parse {ws_file}: {e}")

        return DashboardState(workstreams=workstreams)

    def _parse_ws_file(self, path: Path, status_dir: str) -> WorkstreamState:
        """Parse single workstream file.

        Args:
            path: Path to workstream markdown file
            status_dir: Status directory name (backlog, in_progress, completed)

        Returns:
            WorkstreamState for the file

        Raises:
            ValueError: If YAML frontmatter is invalid
        """
        import yaml

        content = path.read_text()

        # Extract YAML frontmatter
        if content.startswith("---"):
            parts = content.split("---", 2)
            if len(parts) >= 3:
                try:
                    frontmatter = yaml.safe_load(parts[1]) or {}
                except yaml.YAMLError as e:
                    raise ValueError(f"Invalid YAML in {path}: {e}")
            else:
                frontmatter = {}
        else:
            frontmatter = {}

        return WorkstreamState(
            ws_id=frontmatter.get("ws_id", path.stem),
            status=status_dir.replace("_", "-"),  # in_progress -> in-progress
            title=frontmatter.get("title", path.stem),
            feature=frontmatter.get("feature", "Unknown"),
            assignee=frontmatter.get("assignee"),
            size=frontmatter.get("size"),
        )
