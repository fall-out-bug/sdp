"""Base classes for health checks."""

from dataclasses import dataclass
from typing import Optional


@dataclass
class HealthCheckResult:
    """Result of a health check."""

    name: str
    passed: bool
    message: str
    remediation: Optional[str] = None


class HealthCheck:
    """Base class for health checks."""

    def __init__(self, name: str, critical: bool = True) -> None:
        """Initialize health check.

        Args:
            name: Check name
            critical: Whether check is critical (failure causes exit code 1)
        """
        self.name = name
        self.critical = critical

    def run(self) -> HealthCheckResult:
        """Run the health check.

        Returns:
            HealthCheckResult with check outcome
        """
        raise NotImplementedError
