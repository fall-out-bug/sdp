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
        WS: WS-000-01
        Feature: F00
        Title: Test Workstream
        Goal: Test goal
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
        WS: WS-000-01
        Feature: F00
        Title: Test Workstream
        Goal: Test goal
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
        WS: WS-000-01
        Feature: F00
        Title: Test Workstream
        Goal: Test goal
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
        WS: WS-000-01
        Feature: F00
        Title: Test Workstream
        Goal: Test goal
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
        WS: WS-000-01
        Feature: F00
        Title: Test Workstream
        Goal: Test goal
        ---

        ## Contract

        ### Input
        - File: \`src/example.py\`

        ### Output
        - File: \`src/example.py\`

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
        WS: WS-000-01
        Feature: F00
        Title: Test
        ---

        # Body Content

        This is the body.
        """)

        ws_file = tmp_path / "ws.md"
        ws_file.write_text(ws_content)

        result = validate_workstream_tier(ws_file, "T0")
        # Should successfully parse and extract body
        assert result.tier == CapabilityTier.T0
