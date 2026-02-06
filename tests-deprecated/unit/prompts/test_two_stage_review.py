"""Tests for two-stage review protocol structure."""

from pathlib import Path

import pytest


class TestTwoStageReviewProtocol:
    """Tests for two-stage review protocol structure."""

    @pytest.fixture
    def protocol_path(self) -> Path:
        """Path to the two-stage review protocol file."""
        return Path("docs/reference/two-stage-review.md")

    @pytest.fixture
    def protocol_content(self, protocol_path: Path) -> str:
        """Content of the two-stage review protocol file."""
        return protocol_path.read_text()

    def test_protocol_file_exists(self, protocol_path: Path) -> None:
        """Protocol file exists at expected location."""
        assert protocol_path.exists()

    def test_stage1_checklist_present(self, protocol_content: str) -> None:
        """Stage 1 checklist includes all 5 checks."""
        assert "Stage 1: Spec Compliance" in protocol_content
        assert "Goal Achievement" in protocol_content
        assert "Specification Alignment" in protocol_content
        assert "AC Coverage" in protocol_content
        assert "No Over-Engineering" in protocol_content
        assert "No Under-Engineering" in protocol_content

    def test_stage2_checklist_present(self, protocol_content: str) -> None:
        """Stage 2 checklist includes all 10 checks."""
        assert "Stage 2: Code Quality" in protocol_content
        assert "Tests & Coverage" in protocol_content
        assert "Regression" in protocol_content
        assert "AI-Readiness" in protocol_content
        assert "Clean Architecture" in protocol_content
        assert "Type Hints" in protocol_content
        assert "Error Handling" in protocol_content
        assert "Security" in protocol_content
        assert "No Tech Debt" in protocol_content
        assert "Documentation" in protocol_content
        assert "Git History" in protocol_content

    def test_review_loop_logic(self, protocol_content: str) -> None:
        """Review loop logic is documented."""
        assert "Review Loop Logic" in protocol_content
        assert "Re-review" in protocol_content

    def test_verdict_rules(self, protocol_content: str) -> None:
        """Verdict rules are clear (APPROVED / CHANGES REQUESTED)."""
        assert "APPROVED" in protocol_content
        assert "CHANGES REQUESTED" in protocol_content
        # No "APPROVED WITH NOTES"
        assert "APPROVED WITH NOTES" not in protocol_content
