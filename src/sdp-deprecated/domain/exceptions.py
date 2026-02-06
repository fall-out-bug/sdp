"""Domain-level exceptions.

Pure domain exceptions with no external dependencies.
All domain errors inherit from DomainError.
"""


class DomainError(Exception):
    """Base exception for domain errors."""
    pass


class ValidationError(DomainError):
    """Validation constraint violated."""
    pass


class WorkstreamNotFoundError(DomainError):
    """Workstream does not exist."""

    def __init__(self, ws_id: str) -> None:
        super().__init__(f"Workstream not found: {ws_id}")
        self.ws_id = ws_id


class DependencyCycleError(DomainError):
    """Circular dependency detected."""

    def __init__(self, cycle: list[str]) -> None:
        cycle_str = " -> ".join(cycle)
        super().__init__(f"Circular dependency detected: {cycle_str}")
        self.cycle = cycle


class MissingDependencyError(DomainError):
    """Required dependency is missing."""

    def __init__(self, ws_id: str, missing_dep: str) -> None:
        super().__init__(
            f"Workstream {ws_id} depends on missing workstream {missing_dep}"
        )
        self.ws_id = ws_id
        self.missing_dep = missing_dep
