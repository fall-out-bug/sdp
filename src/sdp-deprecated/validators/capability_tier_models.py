"""Capability tier validation data models."""


from dataclasses import dataclass, field
from enum import Enum


class CapabilityTier(str, Enum):
    """SDP capability tier levels."""

    T0 = "T0"  # Architect (Contract writer)
    T1 = "T1"  # Integrator (Complex build)
    T2 = "T2"  # Implementer (Contract-driven build)
    T3 = "T3"  # Autocomplete (Micro build)


@dataclass
class ValidationCheck:
    """Single validation check result."""

    name: str
    passed: bool
    message: str = ""
    details: list[str] = field(default_factory=list)


@dataclass
class ValidationResult:
    """Result of tier validation."""

    tier: CapabilityTier
    passed: bool
    checks: list[ValidationCheck] = field(default_factory=list)

    def add_check(self, check: ValidationCheck) -> None:
        """Add a validation check result."""
        self.checks.append(check)
        if not check.passed:
            self.passed = False
