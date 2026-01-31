"""Tests for tutorial practice file.

This file demonstrates the TDD cycle: Red → Green → Refactor
Follow along with the tutorial at docs/beginner/00-quick-start.md
"""

import pytest
from tutorial_practice import greet, greet_multiple, greet_with_time


class TestGreet:
    """Test suite for greet function."""

    def test_greet_friendly_default(self) -> None:
        """Test default friendly greeting."""
        result = greet("Alice")
        assert result == "Hello, Alice! 👋"

    def test_greet_formal(self) -> None:
        """Test formal greeting style."""
        result = greet("Bob", style="formal")
        assert result == "Good day, Mr. Bob."

    def test_greet_casual(self) -> None:
        """Test casual greeting style."""
        result = greet("Charlie", style="casual")
        assert result == "Yo Charlie!"

    def test_greet_empty_name(self) -> None:
        """Test greeting with empty name."""
        result = greet("", style="friendly")
        assert result == "Hello, ! 👋"


class TestGreetMultiple:
    """Test suite for greet_multiple function."""

    def test_greet_multiple_default(self) -> None:
        """Test greeting multiple people."""
        result = greet_multiple(["Alice", "Bob", "Charlie"])
        assert result == [
            "Hello, Alice! 👋",
            "Hello, Bob! 👋",
            "Hello, Charlie! 👋"
        ]

    def test_greet_multiple_formal(self) -> None:
        """Test greeting multiple people formally."""
        result = greet_multiple(["Alice", "Bob"], style="formal")
        assert result == [
            "Good day, Mr. Alice.",
            "Good day, Mr. Bob."
        ]

    def test_greet_multiple_empty_list(self) -> None:
        """Test greeting empty list."""
        result = greet_multiple([])
        assert result == []


class TestGreetWithTime:
    """Test suite for greet_with_time function."""

    def test_morning_greeting(self) -> None:
        """Test morning greeting (5-11)."""
        assert greet_with_time("Alice", 9) == "Good morning, Alice!"
        assert greet_with_time("Bob", 5) == "Good morning, Bob!"
        assert greet_with_time("Charlie", 11) == "Good morning, Charlie!"

    def test_afternoon_greeting(self) -> None:
        """Test afternoon greeting (12-16)."""
        assert greet_with_time("Alice", 12) == "Good afternoon, Alice!"
        assert greet_with_time("Bob", 14) == "Good afternoon, Bob!"
        assert greet_with_time("Charlie", 16) == "Good afternoon, Charlie!"

    def test_evening_greeting(self) -> None:
        """Test evening greeting (17-4)."""
        assert greet_with_time("Alice", 17) == "Good evening, Alice!"
        assert greet_with_time("Bob", 20) == "Good evening, Bob!"
        assert greet_with_time("Charlie", 0) == "Good evening, Charlie!"

    def test_midnight_boundary(self) -> None:
        """Test boundary conditions at midnight."""
        assert greet_with_time("Alice", 4) == "Good evening, Alice!"
        assert greet_with_time("Bob", 5) == "Good morning, Bob!"


class TestComplexity:
    """Tests to demonstrate complexity checking.

    These tests show functions with different complexity levels.
    Run: radon cc tutorial-tests.py -a
    """

    def test_simple_function(self) -> None:
        """This function has CC = 1 (very simple)."""
        assert 1 + 1 == 2

    def test_medium_function(self) -> None:
        """This function has CC = 3 (acceptable)."""
        x = 5
        if x > 0:
            result = x * 2
        else:
            result = x
        assert result == 10

    def test_complex_function_too_many_branches(self) -> None:
        """This function has CC > 10 (violates quality gate).

        DON'T write code like this! This demonstrates what to avoid.
        """
        x = 5

        if x > 0:
            if x > 2:
                if x > 4:
                    if x > 6:
                        if x > 8:
                            result = "very high"
                        else:
                            result = "high"
                    else:
                        result = "medium"
                else:
                    result = "low"
            else:
                result = "very low"
        else:
            result = "negative"

        assert result == "high"

    # ✅ BETTER: Refactor into smaller functions
    def test_refactored_complex_function(self) -> None:
        """This version has CC < 10 (acceptable)."""
        x = 5
        result = self._categorize_value(x)
        assert result == "high"

    def _categorize_value(self, x: int) -> str:
        """Helper function to categorize value (CC < 10)."""
        if x > 8:
            return "very high"
        elif x > 6:
            return "high"
        elif x > 4:
            return "medium"
        elif x > 2:
            return "low"
        elif x > 0:
            return "very low"
        else:
            return "negative"
