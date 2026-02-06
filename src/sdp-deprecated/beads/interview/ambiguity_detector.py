"""
Ambiguity detection for interview answers.

Detects vague or conflicting answers that require follow-up questions.
"""

from typing import Dict

from .models import AmbiguityResult


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
            details=(
                f"Found ambiguity in {len(ambiguous_fields)} fields"
                if has_ambiguity
                else "Clear answers"
            ),
        )

    def _is_vague(self, answer: str) -> bool:
        """Check if answer contains vague patterns."""
        answer_lower = answer.lower()
        return any(pattern in answer_lower for pattern in self.VAGUE_PATTERNS)

    def _check_conflicts(
        self, question_id: str, answer: str, all_answers: Dict[str, str]
    ) -> bool:
        """Check for conflicting answers."""
        # Simplified conflict detection
        # In real implementation would be more sophisticated
        if "mission" in question_id:
            # Check if mission conflicts with "no users" answer
            users_answer = all_answers.get("users", "")
            if any(
                indicator in users_answer.lower()
                for indicator in self.CONFLICT_INDICATORS["mission_no_users"]
            ):
                return True
        return False
