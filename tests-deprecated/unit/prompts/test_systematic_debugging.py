"""Tests for systematic debugging prompt structure."""

from pathlib import Path

import pytest


class TestSystematicDebuggingPrompt:
    """Tests for systematic debugging prompt structure."""

    @pytest.fixture
    def prompt_path(self) -> Path:
        """Path to the systematic debugging prompt file."""
        return Path("docs/reference/systematic-debugging.md")

    @pytest.fixture
    def prompt_content(self, prompt_path: Path) -> str:
        """Content of the systematic debugging prompt file."""
        return prompt_path.read_text()

    def test_prompt_file_exists(self, prompt_path: Path) -> None:
        """Prompt file exists at expected location."""
        assert prompt_path.exists()

    def test_four_phases_present(self, prompt_content: str) -> None:
        """All 4 phases documented."""
        assert "Phase 1: Evidence Collection" in prompt_content
        assert "Phase 2: Pattern Analysis" in prompt_content
        assert "Phase 3: Hypothesis Testing" in prompt_content
        assert "Phase 4: Implementation" in prompt_content

    def test_evidence_collection_checklist(self, prompt_content: str) -> None:
        """Phase 1 has complete checklist."""
        assert "Error Messages" in prompt_content
        assert "Reproduce the Issue" in prompt_content
        assert "Recent Changes" in prompt_content
        assert "Environment State" in prompt_content

    def test_failsafe_rule(self, prompt_content: str) -> None:
        """3 strikes rule is documented."""
        assert "Failsafe Rule: 3 Strikes" in prompt_content
        assert "After 3 failed fix attempts" in prompt_content
        assert "STOP" in prompt_content
        assert "architecture" in prompt_content.lower()

    def test_root_cause_tracing(self, prompt_content: str) -> None:
        """Root-cause tracing technique is documented."""
        assert "Root-Cause Tracing" in prompt_content
