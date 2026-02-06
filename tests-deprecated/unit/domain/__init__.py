"""Tests for domain package initialization."""

from sdp.domain import (
    AcceptanceCriterion,
    DependencyCycleError,
    DomainError,
    Feature,
    MissingDependencyError,
    ValidationError,
    Workstream,
    WorkstreamID,
    WorkstreamSize,
    WorkstreamStatus,
    WorkstreamNotFoundError,
)


def test_all_exports_available() -> None:
    """Verify all public exports are accessible."""
    # Workstream types
    assert Workstream is not None
    assert WorkstreamID is not None
    assert WorkstreamStatus is not None
    assert WorkstreamSize is not None
    assert AcceptanceCriterion is not None

    # Feature types
    assert Feature is not None

    # Exceptions
    assert DomainError is not None
    assert ValidationError is not None
    assert WorkstreamNotFoundError is not None
    assert DependencyCycleError is not None
    assert MissingDependencyError is not None
