"""
Two-round interviewer for @idea skill.

Conducts progressive disclosure interview with automatic ambiguity detection.
"""

from typing import Any, Dict, Optional

from .ambiguity_detector import AmbiguityDetector
from .models import InterviewResult


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
