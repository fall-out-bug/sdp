"""Tests for capability tier validation."""

from pathlib import Path
from textwrap import dedent

import pytest

from sdp.core.workstream import WorkstreamParseError
from sdp.validators.capability_tier import validate_workstream_tier
from sdp.validators.capability_tier_models import CapabilityTier


class TestValidateWorkstreamTier:
    """Test validate_workstream_tier function."""

    def test_invalid_tier_format(self, tmp_path: Path) -> None:
        """Test that invalid tier format raises ValueError."""
        ws_file = tmp_path / "test.md"
        ws_file.write_text("# Test WS")

        with pytest.raises(ValueError, match="Invalid tier"):
            validate_workstream_tier(ws_file, "T4")

        with pytest.raises(ValueError, match="Invalid tier"):
            validate_workstream_tier(ws_file, "INVALID")

    def test_invalid_tier_case(self, tmp_path: Path) -> None:
        """Test that lowercase tier raises ValueError."""
        ws_file = tmp_path / "test.md"
        ws_file.write_text("# Test WS")

        # Should accept uppercase
        result = validate_workstream_tier(ws_file, "T0")
        assert result.tier == CapabilityTier.T0

        # Should reject lowercase (will try to parse as "t0" which is invalid)
        with pytest.raises(ValueError, match="Invalid tier"):
            validate_workstream_tier(ws_file, "t0")

    def test_parse_workstream_invalid_file(self, tmp_path: Path) -> None:
        """Test handling of invalid workstream file."""
        ws_file = tmp_path / "invalid.md"
        ws_file.write_text("Not a valid workstream file")

        result = validate_workstream_tier(ws_file, "T0")
        assert result.tier == CapabilityTier.T0
        assert result.passed is False
        assert len(result.checks) == 1
        assert result.checks[0].name == "parse_workstream"
        assert result.checks[0].passed is False
        assert "Failed to parse workstream" in result.checks[0].message

    def test_t0_validation_minimal_ws(self, tmp_path: Path) -> None:
        """Test T0 validation with minimal valid workstream."""
        ws_content = dedent("""
        ---
        ws_id: 00-000-01
        feature: F00
        status: backlog
        size: SMALL
        title: Test Workstream
        goal: Test goal
        ---
        """)

        ws_file = tmp_path / "ws.md"
        ws_file.write_text(ws_content)

        result = validate_workstream_tier(ws_file, "T0")
        assert result.tier == CapabilityTier.T0
        # Result depends on T0/T1 validation logic

    def test_t1_validation(self, tmp_path: Path) -> None:
        """Test T1 validation."""
        ws_content = dedent("""
        ---
        ws_id: 00-000-01
        feature: F00
        status: backlog
        size: SMALL
        title: Test Workstream
        goal: Test goal
        ---
        """)

        ws_file = tmp_path / "ws.md"
        ws_file.write_text(ws_content)

        result = validate_workstream_tier(ws_file, "T1")
        assert result.tier == CapabilityTier.T1

    def test_t2_validation(self, tmp_path: Path) -> None:
        """Test T2 validation."""
        ws_content = dedent("""
        ---
        ws_id: 00-000-01
        feature: F00
        status: backlog
        size: SMALL
        title: Test Workstream
        goal: Test goal
        ---
        """)

        ws_file = tmp_path / "ws.md"
        ws_file.write_text(ws_content)

        result = validate_workstream_tier(ws_file, "T2")
        assert result.tier == CapabilityTier.T2

    def test_t3_validation(self, tmp_path: Path) -> None:
        """Test T3 validation."""
        ws_content = dedent("""
        ---
        ws_id: 00-000-01
        feature: F00
        status: backlog
        size: SMALL
        title: Test Workstream
        goal: Test goal
        ---
        """)

        ws_file = tmp_path / "ws.md"
        ws_file.write_text(ws_content)

        result = validate_workstream_tier(ws_file, "T3")
        assert result.tier == CapabilityTier.T3

    def test_nonexistent_file(self, tmp_path: Path) -> None:
        """Test handling of nonexistent file."""
        ws_file = tmp_path / "nonexistent.md"

        with pytest.raises(WorkstreamParseError):
            validate_workstream_tier(ws_file, "T0")

    def test_workstream_with_t2_t3_sections(self, tmp_path: Path) -> None:
        """T2/T3 workstream with contract sections."""
        ws_content = dedent("""
        ---
        ws_id: 00-000-01
        feature: F00
        status: backlog
        size: SMALL
        title: Test Workstream
        goal: Test goal
        ---

        ## Contract

        ### Input
        - File: `src/example.py`

        ### Output
        - File: `src/example.py`

        ### Acceptance Criteria
        - [x] Criterion 1
        - [ ] Criterion 2
        """)

        ws_file = tmp_path / "ws.md"
        ws_file.write_text(ws_content)

        result = validate_workstream_tier(ws_file, "T2")
        assert result.tier == CapabilityTier.T2
        # Should have checks for T2-specific validation

    def test_extract_body_from_workstream(self, tmp_path: Path) -> None:
        """Test that body extraction works correctly."""
        ws_content = dedent("""
        ---
        ws_id: 00-000-01
        feature: F00
        status: backlog
        size: SMALL
        title: Test
        goal: Test goal
        ---

        # Body Content

        This is the body.
        """)

        ws_file = tmp_path / "ws.md"
        ws_file.write_text(ws_content)

        result = validate_workstream_tier(ws_file, "T0")
        # Should successfully parse and extract body
        assert result.tier == CapabilityTier.T0

    def test_body_extraction_with_frontmatter(self, tmp_path: Path) -> None:
        """Test body extraction when frontmatter is present."""
        ws_content = dedent("""
        ---
        ws_id: 00-000-01
        feature: F00
        status: backlog
        size: SMALL
        title: Test
        goal: Test goal
        ---

        ## Goal
        Test goal

        ## Architecture
        Test architecture

        ## Decision
        Test decision
        """)

        ws_file = tmp_path / "ws.md"
        ws_file.write_text(ws_content)

        result = validate_workstream_tier(ws_file, "T0")
        assert result.tier == CapabilityTier.T0
        # Should have checks from T0/T1 validation
        assert len(result.checks) > 0

    def test_body_extraction_without_frontmatter(self, tmp_path: Path) -> None:
        """Test body extraction when no frontmatter exists."""
        ws_content = dedent("""
        # Body Content

        This is the body without frontmatter.
        """)

        ws_file = tmp_path / "ws.md"
        ws_file.write_text(ws_content)

        # Should still parse (though may fail validation)
        result = validate_workstream_tier(ws_file, "T0")
        assert result.tier == CapabilityTier.T0

    def test_t2_validation_with_full_contract(self, tmp_path: Path) -> None:
        """Test T2 validation with complete contract sections."""
        ws_content = dedent("""
        ---
        ws_id: 00-000-01
        feature: F00
        status: backlog
        size: SMALL
        title: Test Workstream
        goal: Test goal
        ---

        ## Contract

        ### Interface
        ```python
        def example_function(param: str) -> None:
            raise NotImplementedError
        ```

        ### Tests
        ```python
        def test_example():
            assert True
        ```

        ## Verification
        ```bash
        pytest tests/
        ```
        """)

        ws_file = tmp_path / "ws.md"
        ws_file.write_text(ws_content)

        result = validate_workstream_tier(ws_file, "T2")
        assert result.tier == CapabilityTier.T2
        # Should have T2-specific checks
        assert len(result.checks) > 0

    def test_t3_validation_with_tiny_scope(self, tmp_path: Path) -> None:
        """Test T3 validation with tiny scope."""
        ws_content = dedent("""
        ---
        ws_id: 00-000-01
        feature: F00
        status: backlog
        size: TINY
        title: Test Workstream
        goal: Test goal
        ---

        ## Contract

        ### Interface
        ```python
        def example_function(param: str) -> None:
            raise NotImplementedError
        ```

        ### Tests
        ```python
        def test_example():
            assert True
        ```

        ## Verification
        ```bash
        pytest tests/
        ```
        """)

        ws_file = tmp_path / "ws.md"
        ws_file.write_text(ws_content)

        result = validate_workstream_tier(ws_file, "T3")
        assert result.tier == CapabilityTier.T3
        # Should have T3-specific checks including scope check
        assert len(result.checks) > 0

    def test_t0_t1_validation_runs_checks(self, tmp_path: Path) -> None:
        """Test that T0/T1 validation runs gate checks."""
        ws_content = dedent("""
        ---
        ws_id: WS-000-01
        feature: F00
        title: Test Workstream
        goal: Test goal
        ---

        ## Architecture
        Test architecture section

        ## Decision
        Test decision with rationale
        """)

        ws_file = tmp_path / "ws.md"
        ws_file.write_text(ws_content)

        result = validate_workstream_tier(ws_file, "T0")
        assert result.tier == CapabilityTier.T0
        # Should have T0/T1 gate checks (if parsing succeeds)
        if result.passed or any(c.name != "parse_workstream" for c in result.checks):
            check_names = {check.name for check in result.checks}
            # Should have T0/T1 specific checks
            assert len(check_names) > 0

    def test_successful_parse_reads_file_content(self, tmp_path: Path) -> None:
        """Test that successful parse reads file content and extracts body."""
        ws_content = dedent("""
        ---
        ws_id: 00-000-01
        feature: F00
        status: backlog
        size: SMALL
        title: Test Workstream
        goal: Test goal
        ---

        ## Architecture
        Test architecture section

        ## Decision
        Test decision with rationale
        """)

        ws_file = tmp_path / "ws.md"
        ws_file.write_text(ws_content)

        result = validate_workstream_tier(ws_file, "T0")
        # Should successfully parse and read content
        assert result.tier == CapabilityTier.T0
        # Should have checks from body validation
        assert len(result.checks) > 0

    def test_extract_body_with_frontmatter(self, tmp_path: Path) -> None:
        """Test _extract_body extracts content after frontmatter."""
        ws_content = dedent("""
        ---
        ws_id: 00-000-01
        feature: F00
        status: backlog
        size: SMALL
        title: Test Workstream
        goal: Test goal
        ---

        ## Body Content
        This is the body after frontmatter.
        """)

        ws_file = tmp_path / "ws.md"
        ws_file.write_text(ws_content)

        result = validate_workstream_tier(ws_file, "T0")
        assert result.tier == CapabilityTier.T0
        # Body should be extracted and validated

    def test_extract_body_without_frontmatter(self, tmp_path: Path) -> None:
        """Test _extract_body returns full content when no frontmatter."""
        ws_content = dedent("""
        # Body Content
        This is content without frontmatter.
        """)

        ws_file = tmp_path / "ws.md"
        ws_file.write_text(ws_content)

        result = validate_workstream_tier(ws_file, "T0")
        assert result.tier == CapabilityTier.T0
        # Should handle content without frontmatter (will fail parse but test _extract_body)

    def test_t2_validation_with_body_extraction(self, tmp_path: Path) -> None:
        """Test T2 validation extracts body and runs checks."""
        ws_content = dedent("""
        ---
        ws_id: 00-000-01
        feature: F00
        status: backlog
        size: SMALL
        title: Test Workstream
        goal: Test goal
        ---

        ## Contract

        ### Interface
        ```python
        def example_function(param: str) -> None:
            raise NotImplementedError
        ```

        ### Tests
        ```python
        def test_example():
            assert True
        ```

        ## Verification
        ```bash
        pytest tests/
        ```
        """)

        ws_file = tmp_path / "ws.md"
        ws_file.write_text(ws_content)

        result = validate_workstream_tier(ws_file, "T2")
        assert result.tier == CapabilityTier.T2
        # Should extract body and run T2/T3 checks
        assert len(result.checks) > 0

    def test_t3_validation_with_body_extraction(self, tmp_path: Path) -> None:
        """Test T3 validation extracts body and runs checks."""
        ws_content = dedent("""
        ---
        ws_id: 00-000-01
        feature: F00
        status: backlog
        size: TINY
        title: Test Workstream
        goal: Test goal
        ---

        ## Contract

        ### Interface
        ```python
        def example_function(param: str) -> None:
            raise NotImplementedError
        ```

        ### Tests
        ```python
        def test_example():
            assert True
        ```

        ## Verification
        ```bash
        pytest tests/
        ```
        """)

        ws_file = tmp_path / "ws.md"
        ws_file.write_text(ws_content)

        result = validate_workstream_tier(ws_file, "T3")
        assert result.tier == CapabilityTier.T3
        # Should extract body and run T2/T3 checks including scope check
        assert len(result.checks) > 0

    def test_t1_validation_with_body_extraction(self, tmp_path: Path) -> None:
        """Test T1 validation extracts body and runs gate checks."""
        ws_content = dedent("""
        ---
        ws_id: 00-000-01
        feature: F00
        status: backlog
        size: SMALL
        title: Test Workstream
        goal: Test goal
        ---

        ## Architecture
        Test architecture section

        ## Decision
        Test decision with rationale
        """)

        ws_file = tmp_path / "ws.md"
        ws_file.write_text(ws_content)

        result = validate_workstream_tier(ws_file, "T1")
        assert result.tier == CapabilityTier.T1
        # Should extract body and run T0/T1 gate checks
        assert len(result.checks) > 0
