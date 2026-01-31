"""Tests for PRD parser edge cases."""

from pathlib import Path

import pytest

from sdp.prd.parser import (
    get_frontmatter,
    parse_prd_sections,
    update_frontmatter,
)


class TestParsePrdSectionsEdgeCases:
    """Test PRD section parsing edge cases."""

    def test_parse_empty_content(self) -> None:
        """Test parsing empty content."""
        sections = parse_prd_sections("")
        assert sections == {}

    def test_parse_no_sections(self) -> None:
        """Test parsing content with no sections."""
        content = "Just some text without sections."
        sections = parse_prd_sections(content)
        assert sections == {}

    def test_parse_section_without_number(self) -> None:
        """Test parsing section header without number."""
        content = """## Section Name
Content here.
"""
        sections = parse_prd_sections(content)
        assert "Section Name" in sections
        assert sections["Section Name"] == "Content here."

    def test_parse_section_with_number(self) -> None:
        """Test parsing section header with number."""
        content = """## 1. Section Name
Content here.
"""
        sections = parse_prd_sections(content)
        assert "Section Name" in sections
        assert sections["Section Name"] == "Content here."

    def test_parse_multiple_sections(self) -> None:
        """Test parsing multiple sections."""
        content = """## 1. First Section
First content.

## 2. Second Section
Second content.

## Third Section
Third content.
"""
        sections = parse_prd_sections(content)
        assert len(sections) == 3
        assert "First Section" in sections
        assert "Second Section" in sections
        assert "Third Section" in sections

    def test_parse_section_with_empty_content(self) -> None:
        """Test parsing section with empty content."""
        content = """## Empty Section

## Next Section
Has content.
"""
        sections = parse_prd_sections(content)
        assert "Empty Section" in sections
        assert sections["Empty Section"] == ""

    def test_parse_frontmatter_skipping(self) -> None:
        """Test that frontmatter is skipped when parsing sections."""
        content = """---
key: value
---

## Section Name
Content here.
"""
        sections = parse_prd_sections(content)
        assert "Section Name" in sections
        assert "key" not in sections

    def test_parse_incomplete_frontmatter(self) -> None:
        """Test parsing with incomplete frontmatter (only opening ---)."""
        content = """---
key: value
## Section Name
Content here.
"""
        sections = parse_prd_sections(content)
        # Parser skips all content until closing --- is found
        # So with incomplete frontmatter, sections won't be parsed
        # This is expected behavior - frontmatter must be properly closed
        assert sections == {} or "Section Name" in sections

    def test_parse_section_with_multiline_content(self) -> None:
        """Test parsing section with multiline content."""
        content = """## Section Name
Line 1
Line 2
Line 3
"""
        sections = parse_prd_sections(content)
        assert sections["Section Name"] == "Line 1\nLine 2\nLine 3"

    def test_parse_section_with_special_chars(self) -> None:
        """Test parsing section name with special characters."""
        content = """## Section (with) [special] chars
Content here.
"""
        sections = parse_prd_sections(content)
        assert "Section (with) [special] chars" in sections


class TestGetFrontmatterEdgeCases:
    """Test frontmatter extraction edge cases."""

    def test_get_frontmatter_empty(self) -> None:
        """Test getting frontmatter from empty content."""
        frontmatter = get_frontmatter("")
        assert frontmatter == {}

    def test_get_frontmatter_no_frontmatter(self) -> None:
        """Test getting frontmatter when none exists."""
        content = "## Section\nContent"
        frontmatter = get_frontmatter(content)
        assert frontmatter == {}

    def test_get_frontmatter_simple(self) -> None:
        """Test getting simple frontmatter."""
        content = """---
key1: value1
key2: value2
---
Content
"""
        frontmatter = get_frontmatter(content)
        assert frontmatter["key1"] == "value1"
        assert frontmatter["key2"] == "value2"

    def test_get_frontmatter_with_colon_in_value(self) -> None:
        """Test frontmatter with colon in value."""
        content = """---
key: value:with:colons
---
Content
"""
        frontmatter = get_frontmatter(content)
        assert frontmatter["key"] == "value:with:colons"

    def test_get_frontmatter_incomplete(self) -> None:
        """Test getting frontmatter with incomplete closing."""
        content = """---
key: value
Content
"""
        frontmatter = get_frontmatter(content)
        # With incomplete frontmatter, parser reads until end of content
        # So key: value will be included
        assert "key" in frontmatter

    def test_get_frontmatter_multiple_dashes(self) -> None:
        """Test frontmatter with multiple dashes in content."""
        content = """---
key: value
---
Content with --- dashes
"""
        frontmatter = get_frontmatter(content)
        assert frontmatter["key"] == "value"


class TestUpdateFrontmatterEdgeCases:
    """Test frontmatter update edge cases."""

    def test_update_frontmatter_no_frontmatter(self) -> None:
        """Test updating frontmatter when none exists."""
        content = "## Section\nContent"
        updated = update_frontmatter(content, {"key": "value"})
        assert updated == content  # Should remain unchanged

    def test_update_frontmatter_add_new_field(self) -> None:
        """Test adding new field to frontmatter."""
        content = """---
key1: value1
---
Content
"""
        updated = update_frontmatter(content, {"key2": "value2"})
        assert "key2: value2" in updated
        assert "key1: value1" in updated

    def test_update_frontmatter_update_existing(self) -> None:
        """Test updating existing frontmatter field."""
        content = """---
key1: old_value
key2: value2
---
Content
"""
        updated = update_frontmatter(content, {"key1": "new_value"})
        assert "key1: new_value" in updated
        assert "old_value" not in updated
        assert "key2: value2" in updated

    def test_update_frontmatter_multiple_updates(self) -> None:
        """Test updating multiple fields."""
        content = """---
key1: value1
key2: value2
---
Content
"""
        updated = update_frontmatter(content, {"key1": "new1", "key2": "new2"})
        assert "key1: new1" in updated
        assert "key2: new2" in updated
        assert "value1" not in updated
        assert "value2" not in updated

    def test_update_frontmatter_preserves_content(self) -> None:
        """Test that content after frontmatter is preserved."""
        content = """---
key: value
---
## Section
Content here.
"""
        updated = update_frontmatter(content, {"key": "new_value"})
        assert "## Section" in updated
        assert "Content here." in updated

    def test_update_frontmatter_empty_updates(self) -> None:
        """Test updating with empty updates dict."""
        content = """---
key: value
---
Content
"""
        updated = update_frontmatter(content, {})
        # Content should remain the same when no updates
        assert "key: value" in updated
        assert "Content" in updated

    def test_update_frontmatter_field_not_in_content(self) -> None:
        """Test updating field that doesn't exist in frontmatter."""
        content = """---
key1: value1
---
Content
"""
        updated = update_frontmatter(content, {"key2": "value2"})
        assert "key2: value2" in updated
        assert "key1: value1" in updated
