"""Tests for testing anti-patterns guide."""

from pathlib import Path

import pytest


class TestTestingAntipatterns:
    """Tests for testing anti-patterns guide."""

    @pytest.fixture
    def prompt_path(self) -> Path:
        """Path to the testing anti-patterns prompt file."""
        return Path("docs/reference/testing-antipatterns.md")

    @pytest.fixture
    def prompt_content(self, prompt_path: Path) -> str:
        """Content of the testing anti-patterns prompt file."""
        return prompt_path.read_text()

    def test_prompt_file_exists(self, prompt_path: Path) -> None:
        """Anti-patterns file exists."""
        assert prompt_path.exists()

    def test_seven_antipatterns_documented(self, prompt_content: str) -> None:
        """All 7 anti-patterns are documented."""
        assert "Anti-Pattern 1: Mocking What You're Testing" in prompt_content
        assert "Anti-Pattern 2: Test-Only Code Paths" in prompt_content
        assert "Anti-Pattern 3: Incomplete Mocks" in prompt_content
        assert "Anti-Pattern 4: Testing Implementation Details" in prompt_content
        assert "Anti-Pattern 5: Flaky Tests with Timeouts" in prompt_content
        assert "Anti-Pattern 6: Testing Multiple Things" in prompt_content
        assert "Anti-Pattern 7: Tests Without Assertions" in prompt_content

    def test_each_antipattern_has_examples(self, prompt_content: str) -> None:
        """Each anti-pattern has bad and good examples."""
        # Check for example markers
        bad_count = prompt_content.count("❌ Bad Example")
        good_count = prompt_content.count("✅ Good Example")
        assert bad_count >= 7  # At least one bad example per anti-pattern
        assert good_count >= 7  # At least one good example per anti-pattern

    def test_detection_rules_present(self, prompt_content: str) -> None:
        """Detection rules for lint tools are present."""
        assert "Detection Rules Summary" in prompt_content
        assert "ANTIPATTERN_" in prompt_content  # Lint rule identifiers
