"""Data models for supersede validation."""

from dataclasses import dataclass


@dataclass
class SupersedeChain:
    """Supersede chain representation."""

    original_ws: str
    replacements: list[str]
    has_cycle: bool
    final_ws: str | None  # None if cycle


@dataclass
class SupersedeResult:
    """Result of supersede operation."""

    success: bool
    old_ws: str
    new_ws: str
    error: str | None


@dataclass
class ValidationReport:
    """Report of all supersede validations."""

    total_superseded: int
    orphans: list[str]
    cycles: list[SupersedeChain]
    valid_chains: list[SupersedeChain]
