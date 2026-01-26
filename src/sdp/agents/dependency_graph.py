"""Dependency graph for workstream execution ordering."""

from __future__ import annotations

import logging
from collections import deque
from dataclasses import dataclass, field
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from pathlib import Path

logger = logging.getLogger(__name__)


@dataclass
class Node:
    """A node in the dependency graph."""

    ws_id: str
    dependencies: set[str] = field(default_factory=set)
    dependents: set[str] = field(default_factory=set)  # Things that depend on this


class DependencyGraph:
    """Manages workstream dependencies and execution order.

    Uses topological sorting to determine valid execution order.
    Detects circular dependencies.
    """

    def __init__(self, ws_dir: str | Path = "docs/workstreams") -> None:
        """Initialize dependency graph.

        Args:
            ws_dir: Path to workstreams directory
        """
        from pathlib import Path

        self._ws_dir = Path(ws_dir)
        self._nodes: dict[str, Node] = {}
        self._built = False

    def build(self, ws_ids: list[str]) -> None:
        """Build dependency graph from workstream files.

        Args:
            ws_ids: List of workstream IDs to include in graph
        """
        self._nodes = {}

        # Create nodes
        for ws_id in ws_ids:
            self._nodes[ws_id] = Node(ws_id=ws_id)

        # Load dependencies from files
        for ws_id in ws_ids:
            deps = self._load_dependencies(ws_id)
            for dep in deps:
                if dep in self._nodes:
                    # Add edge: ws_id depends on dep
                    self._nodes[ws_id].dependencies.add(dep)
                    self._nodes[dep].dependents.add(ws_id)

        self._built = True

    def execution_order(self) -> list[str]:
        """Return workstreams in valid execution order (topological sort).

        Returns:
            List of ws_id in order they should be executed

        Raises:
            ValueError: If circular dependency detected
        """
        if not self._built:
            return []

        # Kahn's algorithm for topological sort
        in_degree: dict[str, int] = {
            ws_id: len(node.dependencies) for ws_id, node in self._nodes.items()
        }
        queue = deque([ws_id for ws_id, degree in in_degree.items() if degree == 0])
        result: list[str] = []

        while queue:
            ws_id = queue.popleft()
            result.append(ws_id)

            # Reduce in-degree for dependents
            for dependent in self._nodes[ws_id].dependents:
                in_degree[dependent] -= 1
                if in_degree[dependent] == 0:
                    queue.append(dependent)

        if len(result) != len(self._nodes):
            # Circular dependency detected
            cycle = self._find_cycle()
            raise ValueError(f"Circular dependency detected: {' -> '.join(cycle)}")

        return result

    def get_ready(self, completed: set[str]) -> list[str]:
        """Get workstreams ready to execute (dependencies satisfied).

        Args:
            completed: Set of already completed ws_ids

        Returns:
            List of ws_ids that can be executed now
        """
        if not self._built:
            return []

        ready: list[str] = []
        for ws_id, node in self._nodes.items():
            if ws_id in completed:
                continue
            if node.dependencies.issubset(completed):
                ready.append(ws_id)

        return sorted(ready)

    def get_dependencies(self, ws_id: str) -> set[str]:
        """Get dependencies for a workstream.

        Args:
            ws_id: Workstream ID

        Returns:
            Set of ws_ids this workstream depends on
        """
        return self._nodes.get(ws_id, Node(ws_id=ws_id)).dependencies.copy()

    def get_dependents(self, ws_id: str) -> set[str]:
        """Get workstreams that depend on this one.

        Args:
            ws_id: Workstream ID

        Returns:
            Set of ws_ids that depend on this workstream
        """
        return self._nodes.get(ws_id, Node(ws_id=ws_id)).dependents.copy()

    def has_circular_dependency(self) -> bool:
        """Check if graph contains circular dependencies.

        Returns:
            True if circular dependency exists
        """
        try:
            self.execution_order()
            return False
        except ValueError:
            return True

    def _load_dependencies(self, ws_id: str) -> set[str]:
        """Load dependencies from workstream file.

        Args:
            ws_id: Workstream ID

        Returns:
            Set of dependency ws_ids
        """
        from pathlib import Path

        deps: set[str] = set()

        # Search for file in any status directory
        for status in ["backlog", "in_progress", "completed"]:
            path = self._ws_dir / status / f"{ws_id}.md"
            if path.exists():
                deps = self._parse_dependencies(path)
                break

        return deps

    def _parse_dependencies(self, path: Path) -> set[str]:
        """Parse dependencies from workstream frontmatter.

        Args:
            path: Path to workstream file

        Returns:
            Set of dependency ws_ids
        """
        deps: set[str] = set()

        try:
            content = path.read_text()
            if not content.startswith("---"):
                return deps

            # Extract frontmatter
            parts = content.split("---", 2)
            if len(parts) < 3:
                return deps

            frontmatter = parts[1]

            # Parse dependencies
            import re

            for match in re.finditer(r"depends_on:\s*\[(.*?)\]", frontmatter, re.MULTILINE):
                dep_list = match.group(1)
                for dep in re.finditer(r'(["\']?)([\w-]+)\1', dep_list):
                    dep_ws = dep.group(2).strip()
                    if dep_ws and dep_ws != '""':
                        deps.add(dep_ws)

        except Exception as e:
            logger.warning(f"Failed to parse dependencies from {path}: {e}")

        return deps

    def _find_cycle(self) -> list[str]:
        """Find a circular dependency cycle for error reporting.

        Returns:
            List of ws_ids forming a cycle
        """
        # Use DFS to find cycle
        visited: set[str] = set()
        rec_stack: set[str] = set()
        path: list[str] = []
        cycle: list[str] = []

        def dfs(ws_id: str) -> bool:
            nonlocal cycle
            visited.add(ws_id)
            rec_stack.add(ws_id)
            path.append(ws_id)

            for dep in self._nodes.get(ws_id, Node(ws_id=ws_id)).dependencies:
                if dep not in visited and dep in self._nodes:
                    if dfs(dep):
                        return True
                elif dep in rec_stack:
                    # Found cycle
                    idx = path.index(dep)
                    cycle = path[idx:] + [dep]
                    return True

            path.pop()
            rec_stack.remove(ws_id)
            return False

        for ws_id in self._nodes:
            if ws_id not in visited:
                if dfs(ws_id):
                    break

        return cycle

    def visualize(self) -> str:
        """Generate a simple text visualization of the graph.

        Returns:
            String representation of the graph
        """
        if not self._built:
            return "Graph not built"

        lines: list[str] = []
        lines.append("Dependency Graph:")
        lines.append("")

        for ws_id in sorted(self._nodes.keys()):
            node = self._nodes[ws_id]
            deps_str = ", ".join(sorted(node.dependencies)) if node.dependencies else "none"
            lines.append(f"  {ws_id} <- [{deps_str}]")

        return "\n".join(lines)
