"""Feature decomposition and dependency management.

This module provides Feature abstraction that contains multiple Workstreams
with dependency graph calculation and validation.
"""

from collections import defaultdict, deque
from dataclasses import dataclass, field
from pathlib import Path
from typing import Optional

from sdp.core.workstream import Workstream, WorkstreamParseError, parse_workstream


class CircularDependencyError(Exception):
    """Raised when circular dependencies detected in workstreams."""

    pass


class MissingDependencyError(Exception):
    """Raised when a workstream references a missing dependency."""

    pass


@dataclass
class Feature:
    """Feature containing multiple workstreams with dependency management.

    Attributes:
        feature_id: Feature identifier (FXXX)
        workstreams: List of workstreams in this feature
        dependency_graph: Adjacency list representation of dependencies
        execution_order: Topologically sorted execution order
    """

    feature_id: str
    workstreams: list[Workstream] = field(default_factory=list)
    dependency_graph: dict[str, list[str]] = field(default_factory=dict)
    execution_order: list[str] = field(default_factory=list)

    def __post_init__(self) -> None:
        """Build dependency graph and execution order after initialization."""
        self._build_dependency_graph()
        self._validate_dependencies()
        self._calculate_execution_order()

    def _build_dependency_graph(self) -> None:
        """Build adjacency list representation of workstream dependencies."""
        ws_by_id: dict[str, Workstream] = {ws.ws_id: ws for ws in self.workstreams}
        self.dependency_graph = defaultdict(list)

        for ws in self.workstreams:
            for dep_id in ws.dependencies:
                if dep_id not in ws_by_id:
                    raise MissingDependencyError(
                        f"Workstream {ws.ws_id} depends on {dep_id}, but it's not in this feature"
                    )
                self.dependency_graph[ws.ws_id].append(dep_id)

    def _validate_dependencies(self) -> None:
        """Validate no circular dependencies exist."""
        ws_ids = {ws.ws_id for ws in self.workstreams}
        visited: set[str] = set()
        rec_stack: set[str] = set()

        def has_cycle(ws_id: str) -> bool:
            """Check for cycles using DFS."""
            visited.add(ws_id)
            rec_stack.add(ws_id)

            for dep_id in self.dependency_graph.get(ws_id, []):
                if dep_id not in ws_ids:
                    continue  # External dependency, skip
                if dep_id not in visited:
                    if has_cycle(dep_id):
                        return True
                elif dep_id in rec_stack:
                    return True  # Back edge found, cycle detected

            rec_stack.remove(ws_id)
            return False

        for ws_id in ws_ids:
            if ws_id not in visited:
                if has_cycle(ws_id):
                    raise CircularDependencyError(
                        f"Circular dependency detected involving workstream {ws_id}"
                    )  # noqa: E501

    def _build_reverse_graph(self, ws_ids: set[str]) -> dict[str, list[str]]:
        """Build reverse dependency graph.

        Args:
            ws_ids: Set of workstream IDs

        Returns:
            Reverse graph mapping dependency to dependents
        """
        reverse_graph: dict[str, list[str]] = defaultdict(list)
        for ws_id in ws_ids:
            for dep_id in self.dependency_graph.get(ws_id, []):
                if dep_id in ws_ids:
                    reverse_graph[dep_id].append(ws_id)
        return reverse_graph

    def _calculate_in_degrees(self, ws_ids: set[str]) -> dict[str, int]:
        """Calculate in-degree for each workstream.

        Args:
            ws_ids: Set of workstream IDs

        Returns:
            Dictionary mapping ws_id to number of dependencies
        """
        return {ws_id: len(self.dependency_graph.get(ws_id, [])) for ws_id in ws_ids}

    def _calculate_execution_order(self) -> None:
        """Calculate topological sort for execution order."""
        ws_ids = {ws.ws_id for ws in self.workstreams}
        reverse_graph = self._build_reverse_graph(ws_ids)
        in_degree = self._calculate_in_degrees(ws_ids)

        # Kahn's algorithm for topological sort
        queue: deque[str] = deque(ws_id for ws_id in ws_ids if in_degree[ws_id] == 0)
        result: list[str] = []

        while queue:
            ws_id = queue.popleft()
            result.append(ws_id)

            for dependent_id in reverse_graph.get(ws_id, []):
                in_degree[dependent_id] -= 1
                if in_degree[dependent_id] == 0:
                    queue.append(dependent_id)

        # Check if all workstreams were processed
        if len(result) != len(ws_ids):
            remaining = ws_ids - set(result)
            raise CircularDependencyError(
                f"Could not determine execution order. Remaining: {remaining}"
            )

        self.execution_order = result

    def get_workstream(self, ws_id: str) -> Optional[Workstream]:
        """Get workstream by ID.

        Args:
            ws_id: Workstream identifier

        Returns:
            Workstream instance or None if not found
        """
        for ws in self.workstreams:
            if ws.ws_id == ws_id:
                return ws
        return None

    def get_dependencies(self, ws_id: str) -> list[str]:
        """Get direct dependencies for a workstream.

        Args:
            ws_id: Workstream identifier

        Returns:
            List of dependency WS IDs
        """
        return self.dependency_graph.get(ws_id, []).copy()


def load_feature_from_directory(
    feature_id: str, directory: Path, pattern: str = "WS-*.md"
) -> Feature:
    """Load feature from directory containing workstream files.

    Args:
        feature_id: Feature identifier (FXXX)
        directory: Directory containing workstream markdown files
        pattern: Glob pattern for workstream files (default: "WS-*.md")

    Returns:
        Feature instance with loaded workstreams

    Raises:
        WorkstreamParseError: If any workstream file fails to parse
        ValueError: If no workstreams found or feature_id mismatch
    """
    ws_files = sorted(directory.glob(pattern))
    if not ws_files:
        raise ValueError(f"No workstream files found in {directory} matching {pattern}")

    workstreams: list[Workstream] = []
    for ws_file in ws_files:
        try:
            ws = parse_workstream(ws_file)
            if ws.feature != feature_id:
                raise ValueError(
                    f"Workstream {ws.ws_id} has feature {ws.feature}, expected {feature_id}"
                )
            workstreams.append(ws)
        except WorkstreamParseError as e:
            raise WorkstreamParseError(f"Failed to parse {ws_file}: {e}") from e

    return Feature(feature_id=feature_id, workstreams=workstreams)


def load_feature_from_spec(feature_id: str, spec_file: Path) -> Feature:
    """Load feature from spec file (future: parse feature.md and find WS).

    For now, this is a placeholder that loads from the workstreams directory
    based on feature_id.

    Args:
        feature_id: Feature identifier (FXXX)
        spec_file: Path to feature spec file (e.g., feature.md)

    Returns:
        Feature instance

    Raises:
        NotImplementedError: This is a placeholder for future implementation
    """
    # Future: Parse feature.md to find workstream directory
    # For now, infer from spec_file location
    spec_dir = spec_file.parent
    workstreams_dir = spec_dir.parent.parent / "workstreams" / "backlog"

    if not workstreams_dir.exists():
        raise ValueError(f"Workstreams directory not found: {workstreams_dir}")

    return load_feature_from_directory(feature_id, workstreams_dir)
