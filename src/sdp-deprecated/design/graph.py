"""Dependency graph for workstream execution planning."""

from dataclasses import dataclass, field


@dataclass
class WorkstreamNode:
    """A node in the workstream dependency graph."""

    ws_id: str
    depends_on: list[str] = field(default_factory=list)
    oneshot_ready: bool = True
    estimated_loc: int = 0
    estimated_duration: str = ""


class DependencyGraph:
    """Manages workstream dependencies and provides execution order."""

    def __init__(self) -> None:
        self._nodes: dict[str, WorkstreamNode] = {}

    def add(self, node: WorkstreamNode) -> None:
        """Add a workstream node to the graph.

        Args:
            node: WorkstreamNode to add
        """
        self._nodes[node.ws_id] = node

    def get(self, ws_id: str) -> WorkstreamNode | None:
        """Get a workstream node by ID.

        Args:
            ws_id: Workstream ID

        Returns:
            WorkstreamNode if found, None otherwise
        """
        return self._nodes.get(ws_id)

    def topological_sort(self) -> list[str]:
        """Return workstreams in dependency order (Kahn's algorithm).

        Returns:
            List of workstream IDs in execution order

        Raises:
            ValueError: If graph contains a cycle
        """
        # Calculate in-degree for each node
        in_degree: dict[str, int] = {ws_id: 0 for ws_id in self._nodes}
        for ws_id, node in self._nodes.items():
            for dep in node.depends_on:
                if dep not in in_degree:
                    raise ValueError(f"Dependency {dep} not found in graph")
                in_degree[ws_id] += 1

        # Start with nodes that have no dependencies
        queue: list[str] = [
            ws_id for ws_id, degree in in_degree.items() if degree == 0
        ]
        result: list[str] = []

        while queue:
            # Sort for deterministic order
            queue.sort()
            current = queue.pop(0)
            result.append(current)

            # Reduce in-degree for dependent nodes
            for ws_id, node in self._nodes.items():
                if current in node.depends_on:
                    in_degree[ws_id] -= 1
                    if in_degree[ws_id] == 0:
                        queue.append(ws_id)

        if len(result) != len(self._nodes):
            # Cycle detected
            remaining = [ws_id for ws_id in in_degree if in_degree[ws_id] > 0]
            raise ValueError(f"Cycle detected in dependencies: {remaining}")

        return result

    def get_ready_workstreams(self, completed: list[str]) -> list[str]:
        """Get workstreams that are ready to execute.

        Args:
            completed: List of completed workstream IDs

        Returns:
            List of workstream IDs whose dependencies are all satisfied
        """
        ready: list[str] = []
        completed_set = set(completed)

        for ws_id, node in self._nodes.items():
            if ws_id in completed_set:
                continue
            if all(dep in completed_set for dep in node.depends_on):
                ready.append(ws_id)

        return sorted(ready)

    def to_mermaid(self) -> str:
        """Generate Mermaid graph visualization.

        Returns:
            Mermaid graph string
        """
        lines = ["graph TD"]
        for ws_id, node in self._nodes.items():
            for dep in node.depends_on:
                lines.append(f"  {dep} --> {ws_id}")
        return "\n".join(lines)
