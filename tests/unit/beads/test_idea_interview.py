"""
Tests for @idea two-round interview (F014.02).
"""

import pytest
from unittest.mock import Mock, patch

from sdp.beads.idea_interview import (
    InterviewRound,
    CriticalQuestions,
    AmbiguityDetector,
    IdeaInterviewer,
)


class TestInterviewRound:
    """Test InterviewRound enum."""

    def test_round_1_is_critical(self):
        """Round 1 should be critical questions only."""
        round_type = InterviewRound.CRITICAL
        assert round_type.is_required is True
        assert round_type.is_optional is False

    def test_round_2_is_deep_dive(self):
        """Round 2 should be optional deep dive."""
        round_type = InterviewRound.DEEP_DIVE
        assert round_type.is_required is False
        assert round_type.is_optional is True


class TestCriticalQuestions:
    """Test critical questions definition."""

    def test_critical_questions_count(self):
        """Should have 3-5 critical questions."""
        questions = CriticalQuestions.get_round_1_questions()
        assert len(questions) >= 3
        assert len(questions) <= 5

    def test_critical_questions_have_required_fields(self):
        """Each critical question should have required fields."""
        questions = CriticalQuestions.get_round_1_questions()
        
        for q in questions:
            assert "question" in q
            assert "header" in q
            assert "options" in q
            assert len(q["options"]) >= 2


class TestAmbiguityDetector:
    """Test ambiguity detection logic."""

    def test_no_ambiguity_in_clear_answers(self):
        """Should detect no ambiguity in clear answers."""
        detector = AmbiguityDetector()
        
        answers = {
            "mission": "Solve user authentication problem",
            "users": "End users",
            "tech_approach": "PostgreSQL + bcrypt",
        }
        
        ambiguity = detector.detect_ambiguity(answers)
        assert ambiguity.has_ambiguity is False

    def test_ambiguity_in_vague_answers(self):
        """Should detect ambiguity in vague answers."""
        detector = AmbiguityDetector()
        
        answers = {
            "mission": "Something with users",
            "users": "Not sure yet",
            "tech_approach": "Maybe a database",
        }
        
        ambiguity = detector.detect_ambiguity(answers)
        assert ambiguity.has_ambiguity is True
        assert len(ambiguity.ambiguous_fields) > 0

    def test_ambiguity_in_conflicting_answers(self):
        """Should detect ambiguity in conflicting answers."""
        detector = AmbiguityDetector()
        
        answers = {
            "mission": "Add authentication",
            "users": "No users",  # Conflicts with having auth
            "tech_approach": "JWT tokens",
        }
        
        ambiguity = detector.detect_ambiguity(answers)
        assert ambiguity.has_ambiguity is True


class TestIdeaInterviewer:
    """Test idea interviewer with two-round logic."""

    def test_round_1_only_if_no_ambiguity(self):
        """Should only do Round 1 if no ambiguity detected."""
        from sdp.beads.idea_interview import AmbiguityResult

        client = Mock()
        detector = AmbiguityDetector()
        interviewer = IdeaInterviewer(client, detector)

        clear_answers = {
            "mission": "Solve X",
            "users": "Developers",
            "tech_approach": "Python",
        }

        ambiguity_result = AmbiguityResult(
            has_ambiguity=False,
            ambiguous_fields=[],
            confidence=1.0,
        )

        with patch.object(detector, 'detect_ambiguity', return_value=ambiguity_result):
            result = interviewer.conduct_interview(
                "feature-1",
                answers=clear_answers,
                deep_dive_requested=False,
            )

        assert result.rounds_conducted == 1
        assert result.skipped_round_2 is True

    def test_round_2_if_ambiguity_detected(self):
        """Should automatically suggest Round 2 if ambiguity detected."""
        from sdp.beads.idea_interview import AmbiguityResult

        client = Mock()
        detector = AmbiguityDetector()
        interviewer = IdeaInterviewer(client, detector)

        vague_answers = {
            "mission": "Something",
            "users": "Not sure",
            "tech_approach": "Maybe",
        }

        ambiguity_result = AmbiguityResult(
            has_ambiguity=True,
            ambiguous_fields=["mission", "users"],
            confidence=0.4,  # Below 0.5 threshold
        )

        with patch.object(detector, 'detect_ambiguity', return_value=ambiguity_result):
            result = interviewer.conduct_interview(
                "feature-1",
                answers=vague_answers,
                deep_dive_requested=False,
            )

        assert result.rounds_conducted == 2  # Auto-conducted because confidence < 0.5
        assert result.deep_dive_conducted is True
        assert len(result.ambiguous_fields) > 0

    def test_round_2_if_explicitly_requested(self):
        """Should do Round 2 if user explicitly requests."""
        client = Mock()
        detector = AmbiguityDetector()
        interviewer = IdeaInterviewer(client, detector)
        
        answers = {
            "mission": "Solve X",
            "users": "Developers",
        }
        
        result = interviewer.conduct_interview(
            "feature-1",
            answers=answers,
            deep_dive_requested=True,  # Explicit request
        )
        
        assert result.rounds_conducted == 2
        assert result.deep_dive_conducted is True
