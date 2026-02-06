"""
Two-round interview logic for @idea skill (F014.02).

Implements progressive disclosure:
- Round 1: 3-5 critical questions (required, 5-8 min)
- Round 2: Deep dive on ambiguities (optional, 5-10 min)
"""

from .ambiguity_detector import AmbiguityDetector
from .interviewer import IdeaInterviewer
from .models import AmbiguityResult, InterviewResult, InterviewRound
from .questions import CriticalQuestions

__all__ = [
    "InterviewRound",
    "AmbiguityResult",
    "InterviewResult",
    "AmbiguityDetector",
    "CriticalQuestions",
    "IdeaInterviewer",
]
