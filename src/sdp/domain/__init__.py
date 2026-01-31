"""Domain layer - pure business entities with no external dependencies.

This package contains the core domain model for SDP:
- Workstreams: atomic units of work
- Features: aggregates of workstreams
- Value objects: immutable identifiers
- Domain exceptions: business rule violations

The domain layer has ZERO dependencies on other SDP layers (core, beads, unified).
All other layers can depend on domain.

Clean Architecture dependency flow:
    domain ← core ← (beads, unified, cli)
"""

from sdp.domain.exceptions import (
    DependencyCycleError,
    DomainError,
    MissingDependencyError,
    ValidationError,
    WorkstreamNotFoundError,
)
from sdp.domain.feature import Feature
from sdp.domain.workstream import (
    AcceptanceCriterion,
    Workstream,
    WorkstreamID,
    WorkstreamSize,
    WorkstreamStatus,
)

__all__ = [
    # Workstream entities
    "Workstream",
    "WorkstreamID",
    "WorkstreamStatus",
    "WorkstreamSize",
    "AcceptanceCriterion",
    # Feature entities
    "Feature",
    # Exceptions
    "DomainError",
    "ValidationError",
    "WorkstreamNotFoundError",
    "DependencyCycleError",
    "MissingDependencyError",
]
