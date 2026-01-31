"""Pure domain entities for features.

Feature aggregates represent collections of workstreams with dependency management.
"""

from collections import defaultdict, deque
from dataclasses import dataclass, field
from typing import Optional

from sdp.domain.exceptions import (
    DependencyCycleError,
    MissingDependencyError,
)
from sdp.domain.workstream import Workstream


@dataclass
class Feature:
    """Feature containing multiple workstreams with dependency management.

    This is a domain aggregate root that manages workstream dependencies
    and calculates execution order.

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
                        ws_id=ws.ws_id,
                        missing_dep=dep_id,
                    )
                self.dependency_graph[ws.ws_id].append(dep_id)

    def _validate_dependencies(self) -> None:
        """Validate no circular dependencies exist."""
        ws_ids = {ws.ws_id for ws in self.workstreams}
        visited: set[str] = set()
        rec_stack: set[str] = set()
        path: list[str] = []

        def extract_cycle(ws_id: str) -> list[str] | None:
            """Extract cycle using DFS, returns cycle path if found."""
            visited.add(ws_id)
            rec_stack.add(ws_id)
            path.append(ws_id)

            for dep_id in self.dependency_graph.get(ws_id, []):
                if dep_id not in ws_ids:
                    continue
                if dep_id not in visited:
                    cycle = extract_cycle(dep_id)
                    if cycle is not None:
                        return cycle
                elif dep_id in rec_stack:
                    idx = path.index(dep_id)
                    return path[idx:] + [dep_id]

            path.pop()
            rec_stack.remove(ws_id)
            return None

        for ws_id in ws_ids:
            if ws_id not in visited:
                cycle = extract_cycle(ws_id)
                if cycle is not None:
                    raise DependencyCycleError(
                        cycle=cycle[:-1] if cycle[-1] == cycle[0] else cycle,
                    )

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
            raise DependencyCycleError(
                cycle=list(remaining),
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
