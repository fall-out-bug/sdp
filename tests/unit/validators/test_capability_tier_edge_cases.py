"""Extended tests for capability tier validation edge cases."""

from pathlib import Path
import pytest

from sdp.validators.capability_tier import validate_workstream_tier, _extract_body
from sdp.validators.capability_tier_models import CapabilityTier
from sdp.core.workstream import WorkstreamParseError


class TestValidateWorkstreamTier:
    """Tests for validate_workstream_tier function."""

    def test_validate_tier_lowercase_raises(self):
        """Test raises ValueError for lowercase tier."""
        with pytest.raises(ValueError, match="Invalid tier: t0. Must be uppercase"):
            validate_workstream_tier(Path("dummy.md"), "t0")

    def test_validate_tier_mixed_case_raises(self):
        """Test raises ValueError for mixed case tier."""
        with pytest.raises(ValueError, match="Invalid tier: T0x. Must be uppercase"):
            validate_workstream_tier(Path("dummy.md"), "T0x")

    def test_validate_tier_invalid_value_raises(self):
        """Test raises ValueError for invalid tier value."""
        with pytest.raises(ValueError, match="Invalid tier: T9. Must be one of T0, T1, T2, T3"):
            validate_workstream_tier(Path("dummy.md"), "T9")

    def test_validate_tier_file_not_found(self, tmp_path):
        """Test raises WorkstreamParseError when file doesn't exist."""
        ws_path = tmp_path / "nonexistent.md"

        with pytest.raises(WorkstreamParseError, match="File not found"):
            validate_workstream_tier(ws_path, "T0")

    def test_validate_tier_parse_error(self, tmp_path):
        """Test returns failed result when parsing fails."""
        ws_path = tmp_path / "bad_ws.md"
        ws_path.write_text("Not valid frontmatter")

        result = validate_workstream_tier(ws_path, "T0")

        assert result.passed is False
        assert any("Failed to parse workstream" in check.message for check in result.checks)

    def test_validate_t0_tier(self, tmp_path):
        """Test validates T0 tier workstream."""
        ws_path = tmp_path / "ws.md"
        ws_path.write_text("""---
ws_id: "00-001-01"
title: "Test WS"
---

# Test WS

## Goal
Test goal

## Acceptance Criteria
- [ ] AC1: Test criterion
""")

        result = validate_workstream_tier(ws_path, "T0")

        assert result.tier == CapabilityTier.T0

    def test_validate_t1_tier(self, tmp_path):
        """Test validates T1 tier workstream."""
        ws_path = tmp_path / "ws.md"
        ws_path.write_text("""---
ws_id: "00-001-01"
title: "Test WS"
---

# Test WS

## Goal
Test goal

## Acceptance Criteria
- [ ] AC1: Test criterion
""")

        result = validate_workstream_tier(ws_path, "T1")

        assert result.tier == CapabilityTier.T1

    def test_validate_t2_tier(self, tmp_path):
        """Test validates T2 tier workstream."""
        ws_path = tmp_path / "ws.md"
        ws_path.write_text("""---
ws_id: "00-001-01"
title: "Test WS"
---

# Test WS

## Goal
Test goal

## Acceptance Criteria
- [ ] AC1: Test criterion

## Technical Specification
Details here
""")

        result = validate_workstream_tier(ws_path, "T2")

        assert result.tier == CapabilityTier.T2

    def test_validate_t3_tier(self, tmp_path):
        """Test validates T3 tier workstream."""
        ws_path = tmp_path / "ws.md"
        ws_path.write_text("""---
ws_id: "00-001-01"
title: "Test WS"
---

# Test WS

## Goal
Test goal

## Acceptance Criteria
- [ ] AC1: Test criterion

## Technical Specification
Details here
""")

        result = validate_workstream_tier(ws_path, "T3")

        assert result.tier == CapabilityTier.T3


class TestExtractBody:
    """Tests for _extract_body function."""

    def test_extract_body_with_frontmatter(self):
        """Test extracts body from content with frontmatter."""
        content = """---
title: "Test"
---

Body content here
"""

        result = _extract_body(content)

        assert result == "\nBody content here\n"

    def test_extract_body_without_frontmatter(self):
        """Test returns full content when no frontmatter."""
        content = "Just body content"

        result = _extract_body(content)

        assert result == "Just body content"

    def test_extract_body_multiline_frontmatter(self):
        """Test extracts body from content with multiline frontmatter."""
        content = """---
title: "Test"
ws_id: "00-001-01"
tags:
  - test
  - demo
---

Body starts here
Line 2
"""

        result = _extract_body(content)

        assert "Body starts here" in result
        assert "Line 2" in result
        assert "title:" not in result
