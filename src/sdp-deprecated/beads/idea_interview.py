"""
Two-round interview logic for @idea skill (F014.02).

DEPRECATED: This module is split into smaller modules for better maintainability.
Import from sdp.beads.interview instead:
- models: InterviewRound, AmbiguityResult, InterviewResult
- ambiguity_detector: AmbiguityDetector
- questions: CriticalQuestions
- interviewer: IdeaInterviewer

This module remains for backward compatibility.

Implements progressive disclosure:
- Round 1: 3-5 critical questions (required, 5-8 min)
- Round 2: Deep dive on ambiguities (optional, 5-10 min)
"""

# Re-export all public APIs for backward compatibility
from .interview import (  # noqa: F401
    AmbiguityDetector,
    AmbiguityResult,
    CriticalQuestions,
    IdeaInterviewer,
    InterviewResult,
    InterviewRound,
)

__all__ = [
    "InterviewRound",
    "AmbiguityResult",
    "AmbiguityDetector",
    "CriticalQuestions",
    "InterviewResult",
    "IdeaInterviewer",
]
