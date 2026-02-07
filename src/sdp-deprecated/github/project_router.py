"""Route workstreams to appropriate GitHub projects."""

import os
import re
from pathlib import Path


class ProjectRouter:
    """Determine which project a workstream belongs to."""

    # Path patterns → project names
    PATH_PATTERNS: dict[str, str] = {
        "tools/hw_checker": "MSU AI Masters",  # All hw_checker WS go to main project
        "courses/mlsd": "mlsd",
        "courses/bdde": "bdde",
        "sdp/": "MSU AI Masters",
    }

    # SDP-related WS IDs → "MSU AI Masters" project
    # These WS relate to SDP protocol itself or general hw_checker infrastructure
    SDP_WS_PREFIXES = ["190", "191", "192", "193", "410", "411", "412", "413", "414", "415"]

    DEFAULT_PROJECT = "MSU AI Masters"  # Default for hw_checker and SDP
    MAIN_PROJECT = "MSU AI Masters"

    @classmethod
    def get_project_for_ws(cls, ws_file: Path) -> str:
        """Determine project name for workstream.

        Args:
            ws_file: Path to WS file

        Returns:
            Project name (e.g., "MSU AI Masters", "mlsd", "bdde")

        """
        # Check environment override
        override = os.getenv("GITHUB_PROJECT")
        if override:
            return override

        # Check if WS is SDP-related by WS ID
        ws_id = cls._extract_ws_id(ws_file)
        if ws_id:
            ws_prefix = ws_id.split("-")[1] if "-" in ws_id else None
            if ws_prefix and ws_prefix in cls.SDP_WS_PREFIXES:
                return cls.MAIN_PROJECT

        # Check file path
        ws_path = str(ws_file.resolve())

        for pattern, project_name in cls.PATH_PATTERNS.items():
            if pattern in ws_path:
                return project_name

        # Default
        return cls.DEFAULT_PROJECT

    @classmethod
    def _extract_ws_id(cls, ws_file: Path) -> str | None:
        """Extract WS ID from file name or frontmatter.

        Args:
            ws_file: Path to WS file

        Returns:
            WS ID (e.g., "WS-410-01") or None
        """
        # Try from filename first (faster)
        filename = ws_file.name
        match = re.search(r"WS-(\d+)-(\d+)", filename)
        if match:
            return f"WS-{match.group(1)}-{match.group(2)}"

        # Try from frontmatter
        try:
            content = ws_file.read_text(encoding="utf-8")
            match = re.search(r"ws_id:\s*(WS-\d+-\d+)", content)
            if match:
                return match.group(1)
        except (OSError, ValueError):
            pass

        return None

    @classmethod
    def get_project_from_frontmatter(cls, ws_file: Path) -> str | None:
        """Extract project from WS frontmatter (if specified).

        Args:
            ws_file: Path to WS file

        Returns:
            Project name or None

        """
        try:
            content = ws_file.read_text(encoding="utf-8")
            match = re.search(r"github_project:\s*(\w+)", content)
            if match:
                return match.group(1)
        except (OSError, ValueError):
            # File doesn't exist or encoding issue
            pass
        return None

    @classmethod
    def get_all_projects(cls) -> list[str]:
        """Get list of all configured projects.

        Returns:
            List of project names

        """
        projects = set(cls.PATH_PATTERNS.values())
        projects.add(cls.MAIN_PROJECT)
        return sorted(list(projects))
