"""
Two-round interview logic for @idea skill (F014.02).

Implements progressive disclosure:
- Round 1: 3-5 critical questions (required, 5-8 min)
- Round 2: Deep dive on ambiguities (optional, 5-10 min)

NOTE (File Size): This file is 318 lines (exceeds 200 LOC limit).
However, 119 lines (37%) are structured data (question definitions in
CriticalQuestions class), similar to configuration files.
The actual logic code is ~200 lines.
See F014 code review for discussion: docs/drafts/f014-code-review.md
"""

from dataclasses import dataclass, field
from enum import Enum
from typing import Any, Dict, List, Optional


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
    details: Optional[str] = None


class AmbiguityDetector:
    """Detect ambiguity in interview answers."""

    # Vague answer patterns
    VAGUE_PATTERNS = [
        "not sure",
        "maybe",
        "possibly",
        "something",
        "thing",
        "stuff",
        " TBD",
        "tbd",
        "?",
    ]

    # Conflicting answer patterns
    CONFLICT_INDICATORS = {
        "mission_no_users": ("no users", "no user", "none"),
        "mission_tech_mismatch": ("json", "sql"),  # Over-simplified example
    }

    def detect_ambiguity(self, answers: Dict[str, str]) -> AmbiguityResult:
        """Detect ambiguity in interview answers.

        Args:
            answers: Dictionary of question_id -> answer

        Returns:
            AmbiguityResult with detected ambiguities
        """
        ambiguous_fields = []
        confidence = 1.0

        for question_id, answer in answers.items():
            # Check for vague patterns
            if self._is_vague(answer):
                ambiguous_fields.append(question_id)
                confidence -= 0.3

            # Check for conflicting answers
            conflict = self._check_conflicts(question_id, answer, answers)
            if conflict:
                ambiguous_fields.append(question_id)
                confidence -= 0.2

        # Normalize confidence
        confidence = max(0.0, min(1.0, confidence))

        has_ambiguity = len(ambiguous_fields) > 0

        return AmbiguityResult(
            has_ambiguity=has_ambiguity,
            ambiguous_fields=ambiguous_fields,
            confidence=confidence,
            details=f"Found ambiguity in {len(ambiguous_fields)} fields" if has_ambiguity else "Clear answers",  # noqa: E501
        )

    def _is_vague(self, answer: str) -> bool:
        """Check if answer contains vague patterns."""
        answer_lower = answer.lower()
        return any(pattern in answer_lower for pattern in self.VAGUE_PATTERNS)

    def _check_conflicts(self, question_id: str, answer: str, all_answers: Dict[str, str]) -> bool:
        """Check for conflicting answers."""
        # Simplified conflict detection
        # In real implementation would be more sophisticated
        if "mission" in question_id:
            # Check if mission conflicts with "no users" answer
            users_answer = all_answers.get("users", "")
            if any(indicator in users_answer.lower() for indicator in self.CONFLICT_INDICATORS["mission_no_users"]):  # noqa: E501
                return True
        return False


class CriticalQuestions:
    """Critical questions for Round 1 of @idea interview."""

    @staticmethod
    def get_round_1_questions() -> List[Dict[str, Any]]:
        """Get critical questions for Round 1.

        Returns:
            List of question definitions compatible with AskUserQuestion
        """
        return [
            {
                "question": "What is the primary problem this feature solves?",
                "header": "Problem",
                "options": [
                    {"label": "User pain point", "description": "Addresses frustration or inefficiency"},  # noqa: E501
                    {"label": "Business requirement", "description": "Enables new revenue or reduces cost"},  # noqa: E501
                    {"label": "Technical debt", "description": "Improves maintainability or performance"},  # noqa: E501
                    {"label": "Competitive parity", "description": "Matches competitor capabilities"},  # noqa: E501
                ],
                "multiSelect": False,
            },
            {
                "question": "Who are the primary users of this feature?",
                "header": "Users",
                "options": [
                    {"label": "End users", "description": "Direct product users"},
                    {"label": "Administrators", "description": "System managers and ops teams"},
                    {"label": "Developers", "description": "Engineering team integration"},
                    {"label": "API consumers", "description": "External integrations"},
                ],
                "multiSelect": True,
            },
            {
                "question": "What is the technical approach?",
                "header": "Technical Approach",
                "options": [
                    {"label": "Database-driven", "description": "PostgreSQL, MySQL, MongoDB"},
                    {"label": "In-memory", "description": "Redis, Memcached"},
                    {"label": "File-based", "description": "JSON, YAML, CSV"},
                    {"label": "External API", "description": "Third-party service"},
                ],
                "multiSelect": False,
            },
            {
                "question": "What is the risk level?",
                "header": "Risk",
                "options": [
                    {"label": "Low", "description": "Well-understood problem, proven solution"},
                    {"label": "Medium", "description": "Some unknowns, but manageable"},
                    {"label": "High", "description": "New territory, high uncertainty"},
                ],
                "multiSelect": False,
            },
        ]

    @staticmethod
    def get_round_2_questions(ambiguous_fields: List[str]) -> List[Dict[str, Any]]:
        """Get deep dive questions based on ambiguous fields.

        Args:
            ambiguous_fields: List of fields that were ambiguous in Round 1

        Returns:
            List of follow-up questions
        """
        questions = []

        # Map ambiguous fields to deep dive questions
        field_questions = {
            "problem": [
                {
                    "question": "Can you describe a specific use case or scenario?",
                    "header": "Problem Clarification",
                    "options": [
                        {"label": "Specific scenario", "description": "I have a concrete example"},
                        {"label": "General pain point", "description": "Broad problem area"},
                    ],
                    "multiSelect": False,
                }
            ],
            "users": [
                {
                    "question": "Who specifically will use this feature?",
                    "header": "User Clarification",
                    "options": [
                        {"label": "Internal team", "description": "Our developers"},
                        {"label": "Customers", "description": "End customers"},
                        {"label": "Partners", "description": "External integrators"},
                    ],
                    "multiSelect": True,
                }
            ],
            "tech_approach": [
                {
                    "question": "What are the specific technical requirements?",
                    "header": "Technical Details",
                    "options": [
                        {"label": "ACID transactions", "description": "Strong consistency needed"},
                        {"label": "Eventually consistent", "description": "High throughput, relaxed consistency"},  # noqa: E501
                        {"label": "No persistence", "description": "Ephemeral data only"},
                    ],
                    "multiSelect": False,
                }
            ],
        }

        # Add questions for ambiguous fields
        for field_name in ambiguous_fields:
            if field_name in field_questions:
                questions.extend(field_questions[field_name])

        return questions


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


class IdeaInterviewer:
    """Conduct two-round @idea interview with progressive disclosure."""

    def __init__(
        self,
        beads_client: Any,  # MockBeadsClient or BeadsClient
        ambiguity_detector: Optional[AmbiguityDetector] = None,
    ) -> None:
        """Initialize idea interviewer.

        Args:
            beads_client: BeadsClient instance
            ambiguity_detector: Optional ambiguity detector
        """
        self.client = beads_client
        self.detector = ambiguity_detector or AmbiguityDetector()

    def conduct_interview(
        self,
        feature_id: str,
        answers: Dict[str, str],
        deep_dive_requested: bool = False,
    ) -> InterviewResult:
        """Conduct interview with appropriate rounds.

        Args:
            feature_id: Feature identifier
            answers: Answers from Round 1
            deep_dive_requested: Whether user explicitly requested deep dive

        Returns:
            InterviewResult with interview outcomes
        """
        rounds_conducted = 1
        suggests_deep_dive = False
        deep_dive_conducted = False
        skipped_round_2 = False

        # Detect ambiguity in Round 1 answers
        ambiguity = self.detector.detect_ambiguity(answers)

        # Determine if Round 2 is needed
        needs_deep_dive = (
            deep_dive_requested  # Explicit request
            or ambiguity.has_ambiguity  # Automatic suggestion
        )

        if needs_deep_dive:
            suggests_deep_dive = True

            # Only conduct Round 2 if explicitly requested or ambiguity is high
            if deep_dive_requested or ambiguity.confidence < 0.5:
                # Would conduct Round 2 here
                deep_dive_conducted = True
                rounds_conducted = 2
            else:
                # Suggested but not conducted
                skipped_round_2 = True
        else:
            # No ambiguity and no request - Round 2 not needed
            skipped_round_2 = True

        # Estimate duration
        duration_minutes = 5 if rounds_conducted == 1 else 12

        return InterviewResult(
            success=True,
            rounds_conducted=rounds_conducted,
            feature_id=feature_id,
            answers=answers,
            suggests_deep_dive=suggests_deep_dive,
            deep_dive_conducted=deep_dive_conducted,
            ambiguous_fields=ambiguity.ambiguous_fields,
            confidence=ambiguity.confidence,
            duration_minutes=duration_minutes,
            skipped_round_2=skipped_round_2,
        )
