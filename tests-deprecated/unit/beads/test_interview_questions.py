"""Tests for interview/questions.py - Question generation for interview rounds."""

import pytest

from sdp.beads.interview.questions import CriticalQuestions


class TestCriticalQuestions:
    """Test CriticalQuestions class."""

    def test_get_round_1_questions(self) -> None:
        """Test get_round_1_questions returns structured questions."""
        questions = CriticalQuestions.get_round_1_questions()

        assert len(questions) == 4

        # Verify structure
        for q in questions:
            assert "question" in q
            assert "header" in q
            assert "options" in q
            assert "multiSelect" in q
            assert isinstance(q["options"], list)
            assert len(q["options"]) > 0

    def test_round_1_question_structure(self) -> None:
        """Test Round 1 questions have correct structure."""
        questions = CriticalQuestions.get_round_1_questions()

        # Problem question
        problem_q = questions[0]
        assert "problem" in problem_q["question"].lower() or "problem" in problem_q["header"].lower()
        assert len(problem_q["options"]) == 4
        assert problem_q["multiSelect"] is False

        # Users question
        users_q = questions[1]
        assert "user" in users_q["question"].lower() or "user" in users_q["header"].lower()
        assert users_q["multiSelect"] is True

        # Technical approach question
        tech_q = questions[2]
        assert "technical" in tech_q["question"].lower() or "technical" in tech_q["header"].lower()
        assert tech_q["multiSelect"] is False

        # Risk question
        risk_q = questions[3]
        assert "risk" in risk_q["question"].lower() or "risk" in risk_q["header"].lower()
        assert risk_q["multiSelect"] is False

    def test_get_round_2_questions_empty_ambiguous_fields(self) -> None:
        """Test get_round_2_questions with empty ambiguous fields."""
        questions = CriticalQuestions.get_round_2_questions([])
        assert questions == []

    def test_get_round_2_questions_single_field(self) -> None:
        """Test get_round_2_questions with single ambiguous field."""
        questions = CriticalQuestions.get_round_2_questions(["problem"])

        assert len(questions) == 1
        assert "problem" in questions[0]["question"].lower() or "problem" in questions[0]["header"].lower()

    def test_get_round_2_questions_multiple_fields(self) -> None:
        """Test get_round_2_questions with multiple ambiguous fields."""
        questions = CriticalQuestions.get_round_2_questions(["problem", "users", "tech_approach"])

        assert len(questions) == 3

        # Verify all questions are structured correctly
        for q in questions:
            assert "question" in q
            assert "header" in q
            assert "options" in q
            assert "multiSelect" in q

    def test_get_round_2_questions_unknown_field(self) -> None:
        """Test get_round_2_questions with unknown ambiguous field."""
        questions = CriticalQuestions.get_round_2_questions(["unknown_field"])

        # Unknown fields should not generate questions
        assert len(questions) == 0

    def test_get_round_2_questions_problem_field(self) -> None:
        """Test get_round_2_questions for problem field."""
        questions = CriticalQuestions.get_round_2_questions(["problem"])

        assert len(questions) == 1
        q = questions[0]
        assert "problem" in q["header"].lower() or "problem" in q["question"].lower()
        assert len(q["options"]) == 2
        assert q["multiSelect"] is False

    def test_get_round_2_questions_users_field(self) -> None:
        """Test get_round_2_questions for users field."""
        questions = CriticalQuestions.get_round_2_questions(["users"])

        assert len(questions) == 1
        q = questions[0]
        assert "user" in q["header"].lower() or "user" in q["question"].lower()
        assert len(q["options"]) == 3
        assert q["multiSelect"] is True

    def test_get_round_2_questions_tech_approach_field(self) -> None:
        """Test get_round_2_questions for tech_approach field."""
        questions = CriticalQuestions.get_round_2_questions(["tech_approach"])

        assert len(questions) == 1
        q = questions[0]
        assert "technical" in q["header"].lower() or "technical" in q["question"].lower()
        assert len(q["options"]) == 3
        assert q["multiSelect"] is False

    def test_get_round_2_questions_duplicate_fields(self) -> None:
        """Test get_round_2_questions handles duplicate fields."""
        questions = CriticalQuestions.get_round_2_questions(["problem", "problem", "users"])

        # Should deduplicate or handle gracefully
        # Current implementation may add duplicates - test documents behavior
        assert len(questions) >= 2
        problem_count = sum(1 for q in questions if "problem" in q["header"].lower())
        users_count = sum(1 for q in questions if "user" in q["header"].lower())
        assert problem_count >= 1
        assert users_count >= 1

    def test_round_2_question_options_structure(self) -> None:
        """Test Round 2 question options have correct structure."""
        questions = CriticalQuestions.get_round_2_questions(["problem"])

        assert len(questions) > 0
        q = questions[0]

        for option in q["options"]:
            assert "label" in option
            assert "description" in option
            assert isinstance(option["label"], str)
            assert isinstance(option["description"], str)
