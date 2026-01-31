"""Tests for domain exceptions."""

import pytest

from sdp.domain.exceptions import (
    DependencyCycleError,
    DomainError,
    MissingDependencyError,
    ValidationError,
    WorkstreamNotFoundError,
)


def test_domain_error_hierarchy() -> None:
    """Verify exception hierarchy."""
    assert issubclass(ValidationError, DomainError)
    assert issubclass(WorkstreamNotFoundError, DomainError)
    assert issubclass(DependencyCycleError, DomainError)
    assert issubclass(MissingDependencyError, DomainError)


def test_workstream_not_found_error() -> None:
    """Test WorkstreamNotFoundError stores ws_id."""
    error = WorkstreamNotFoundError("00-001-01")
    assert error.ws_id == "00-001-01"
    assert "00-001-01" in str(error)


def test_dependency_cycle_error() -> None:
    """Test DependencyCycleError formats cycle."""
    cycle = ["00-001-01", "00-001-02", "00-001-03"]
    error = DependencyCycleError(cycle)
    assert error.cycle == cycle
    assert "00-001-01 -> 00-001-02 -> 00-001-03" in str(error)


def test_missing_dependency_error() -> None:
    """Test MissingDependencyError stores details."""
    error = MissingDependencyError("00-001-02", "00-001-01")
    assert error.ws_id == "00-001-02"
    assert error.missing_dep == "00-001-01"
    assert "00-001-02" in str(error)
    assert "00-001-01" in str(error)
