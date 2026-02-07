"""
Data models for @idea interview process.

Defines core data structures for two-round interview with progressive disclosure.
"""

from dataclasses import dataclass, field
from enum import Enum
from typing import Dict, List


class InterviewRound(Enum):
    """Interview round type."""

    CRITICAL = "critical"  # Round 1: Required
    DEEP_DIVE = "deep_dive"  # Round 2: Optional

    @property
    def is_required(self) -> bool:
        """Whether this round is required."""
        return self == InterviewRound.CRITICAL

    @property
    def is_optional(self) -> bool:
        """Whether this round is optional."""
        return self == InterviewRound.DEEP_DIVE


@dataclass
class AmbiguityResult:
    """Result of ambiguity detection."""

    has_ambiguity: bool
    ambiguous_fields: List[str] = field(default_factory=list)
    confidence: float = 0.0
    details: str | None = None


@dataclass
class InterviewResult:
    """Result of idea interview process."""

    success: bool
    rounds_conducted: int
    feature_id: str
    answers: Dict[str, str]
    suggests_deep_dive: bool = False
    deep_dive_conducted: bool = False
    ambiguous_fields: List[str] = field(default_factory=list)
    confidence: float = 1.0
    duration_minutes: int = 0
    skipped_round_2: bool = False
