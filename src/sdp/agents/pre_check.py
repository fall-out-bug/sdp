"""Pre-execution validation for workstreams."""

import logging
from pathlib import Path

import yaml

from .errors import PreCheckError

logger = logging.getLogger(__name__)


class PreExecutionChecker:
    """Validates workstream before execution."""

    def __init__(self, ws_dir: str = "docs/workstreams") -> None:
        """Initialize checker.

        Args:
            ws_dir: Workstreams directory
        """
        self._ws_dir = Path(ws_dir)

    def check(self, ws_id: str) -> list[str]:
        """Run all pre-execution checks.

        Args:
            ws_id: Workstream ID

        Returns:
            List of error messages (empty if all checks pass)
        """
        errors = []

        # Check 1: WS file exists and valid YAML
        ws_file = self._find_ws_file(ws_id)
        if not ws_file:
            errors.append(f"Workstream file not found: {ws_id}")
            return errors  # Can't continue without file

        try:
            self._parse_yaml(ws_file)
        except Exception as e:
            errors.append(f"Invalid YAML: {e}")

        # Check 2: Dependencies satisfied
        deps = self._get_dependencies(ws_file)
        for dep in deps:
            if not self._is_completed(dep):
                errors.append(f"Dependency not completed: {dep}")

        # Check 3: No circular dependencies
        if self._has_circular_dependency(ws_id):
            errors.append("Circular dependency detected")

        # Check 4: Size ≤ MEDIUM
        if self._get_size(ws_file) == "LARGE":
            errors.append("Workstream too large (LARGE), please split")

        return errors

    def can_execute(self, ws_id: str) -> bool:
        """Check if workstream can be executed.

        Args:
            ws_id: Workstream ID

        Returns:
            True if all checks pass
        """
        return len(self.check(ws_id)) == 0

    def _find_ws_file(self, ws_id: str) -> Path | None:
        """Find workstream file.

        Args:
            ws_id: Workstream ID

        Returns:
            Path to workstream file or None
        """
        for status_dir in ["backlog", "in_progress", "completed", "blocked"]:
            path = self._ws_dir / status_dir / f"{ws_id}.md"
            if path.exists():
                return path
        return None

    def _parse_yaml(self, ws_file: Path) -> dict:
        """Parse YAML frontmatter from workstream file.

        Args:
            ws_file: Path to workstream file

        Returns:
            Parsed YAML frontmatter dict

        Raises:
            yaml.YAMLError: If YAML is invalid
        """
        content = ws_file.read_text()

        if content.startswith("---"):
            parts = content.split("---", 2)
            if len(parts) >= 2:
                return yaml.safe_load(parts[1]) or {}
        return {}

    def _get_dependencies(self, ws_file: Path) -> list[str]:
        """Get dependency list from workstream.

        Args:
            ws_file: Path to workstream file

        Returns:
            List of dependency WS IDs
        """
        frontmatter = self._parse_yaml(ws_file)
        return frontmatter.get("dependencies", [])

    def _is_completed(self, dep_ws_id: str) -> bool:
        """Check if dependency workstream is completed.

        Args:
            dep_ws_id: Dependency workstream ID

        Returns:
            True if dependency is marked completed
        """
        dep_file = self._find_ws_file(dep_ws_id)
        if not dep_file:
            return False  # Dependency doesn't exist

        try:
            frontmatter = self._parse_yaml(dep_file)
            status = frontmatter.get("status", "")
            return status.replace("-", "_") == "completed"
        except Exception:
            return False

    def _has_circular_dependency(self, ws_id: str) -> bool:
        """Check for circular dependencies using DFS.

        Args:
            ws_id: Workstream ID to check

        Returns:
            True if circular dependency detected
        """
        visited = set()

        def check(current_id: str) -> bool:
            if current_id in visited:
                return True  # Cycle detected

            visited.add(current_id)
            current_file = self._find_ws_file(current_id)
            if not current_file:
                return False

            deps = self._get_dependencies(current_file)
            for dep in deps:
                if check(dep):
                    return True

            visited.remove(current_id)
            return False

        return check(ws_id)

    def _get_size(self, ws_file: Path) -> str:
        """Get workstream size from frontmatter.

        Args:
            ws_file: Path to workstream file

        Returns:
            Size value (TINY, SMALL, MEDIUM, LARGE, AI)
        """
        frontmatter = self._parse_yaml(ws_file)
        return frontmatter.get("size", "MEDIUM").upper()
