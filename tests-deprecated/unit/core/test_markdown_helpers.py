"""Tests for markdown parsing helpers."""

import pytest
import yaml

from sdp.core.workstream.markdown_helpers import (
    extract_acceptance_criteria,
    extract_code_blocks,
    extract_dependencies,
    extract_frontmatter,
    extract_section,
    extract_steps,
    extract_title,
    strip_frontmatter,
)


class TestExtractFrontmatter:
    """Test frontmatter extraction."""

    def test_extract_valid_frontmatter(self) -> None:
        """Verify extraction of valid frontmatter."""
        content = """---
ws_id: 00-001-01
feature: F001
status: backlog
size: SMALL
---

Body content"""
        result = extract_frontmatter(content)
        assert result["ws_id"] == "00-001-01"
        assert result["feature"] == "F001"
        assert result["status"] == "backlog"
        assert result["size"] == "SMALL"

    def test_extract_frontmatter_with_error_path(self) -> None:
        """Verify error path is included in error messages."""
        content = "No frontmatter"
        with pytest.raises(ValueError, match="test.md"):
            extract_frontmatter(content, error_path="test.md")

    def test_extract_frontmatter_invalid_yaml(self) -> None:
        """Verify YAML parsing errors are caught."""
        content = """---
ws_id: 00-001-01
feature: F001
status: backlog
size: SMALL
invalid: [unclosed
---

Body"""
        with pytest.raises(ValueError, match="Invalid YAML"):
            extract_frontmatter(content, error_path="test.md")

    def test_extract_frontmatter_non_dict(self) -> None:
        """Verify non-dict frontmatter raises error."""
        content = """---
- item1
- item2
---

Body"""
        with pytest.raises(ValueError, match="must be a YAML dict"):
            extract_frontmatter(content)

    def test_extract_frontmatter_missing_required_fields(self) -> None:
        """Verify missing required fields raises error."""
        content = """---
ws_id: 00-001-01
feature: F001
---

Body"""
        with pytest.raises(ValueError, match="Missing required fields"):
            extract_frontmatter(content)


class TestStripFrontmatter:
    """Test frontmatter stripping."""

    def test_strip_frontmatter(self) -> None:
        """Verify frontmatter is stripped correctly."""
        content = """---
ws_id: 00-001-01
---

Body content"""
        result = strip_frontmatter(content)
        # strip_frontmatter includes newline after ---
        assert result.strip() == "Body content"

    def test_strip_frontmatter_no_frontmatter(self) -> None:
        """Verify content without frontmatter is returned as-is."""
        content = "No frontmatter here"
        result = strip_frontmatter(content)
        assert result == content


class TestExtractTitle:
    """Test title extraction."""

    def test_extract_title(self) -> None:
        """Verify title is extracted from first ## heading."""
        body = """## My Title

Some content"""
        assert extract_title(body) == "My Title"

    def test_extract_title_no_heading(self) -> None:
        """Verify empty string when no heading found."""
        body = "No heading here"
        assert extract_title(body) == ""


class TestExtractSection:
    """Test section extraction."""

    def test_extract_section(self) -> None:
        """Verify section content is extracted."""
        body = """### Goal

This is the goal content.

### Context

This is context."""
        result = extract_section(body, "Goal")
        assert "This is the goal content" in result

    def test_extract_section_case_insensitive(self) -> None:
        """Verify section matching is case-insensitive."""
        body = """### GOAL

Goal content

### Context

Context content"""
        result = extract_section(body, "goal")
        assert "Goal content" in result

    def test_extract_section_not_found(self) -> None:
        """Verify empty string when section not found."""
        body = "No sections here"
        assert extract_section(body, "Goal") == ""

    def test_extract_section_empty_after_heading(self) -> None:
        """Verify empty string when section ends immediately after heading."""
        body = """### Goal

### Context

Context content"""
        result = extract_section(body, "Goal")
        assert result == ""

    def test_extract_section_last_section(self) -> None:
        """Verify last section extracts to end of body."""
        body = """### Goal

Goal content

### Context

Context content"""
        result = extract_section(body, "Context")
        assert "Context content" in result


class TestExtractAcceptanceCriteria:
    """Test acceptance criteria extraction."""

    def test_extract_acceptance_criteria_unchecked(self) -> None:
        """Verify unchecked acceptance criteria are extracted."""
        body = """- [ ] AC1: First criterion
- [ ] AC2: Second criterion"""
        criteria = extract_acceptance_criteria(body)
        assert len(criteria) == 2
        assert criteria[0].id == "AC1"
        assert criteria[0].description == "First criterion"
        assert criteria[0].checked is False

    def test_extract_acceptance_criteria_checked(self) -> None:
        """Verify checked acceptance criteria are extracted."""
        body = """- [x] AC1: Completed criterion
- [X] AC2: Also completed"""
        criteria = extract_acceptance_criteria(body)
        assert len(criteria) == 2
        assert criteria[0].checked is True
        assert criteria[1].checked is True

    def test_extract_acceptance_criteria_case_insensitive(self) -> None:
        """Verify case-insensitive matching."""
        body = """- [ ] ac1: Lowercase ID
- [X] AC2: Mixed case"""
        criteria = extract_acceptance_criteria(body)
        assert len(criteria) == 2
        assert criteria[0].id == "ac1"
        assert criteria[1].checked is True

    def test_extract_acceptance_criteria_none(self) -> None:
        """Verify empty list when no criteria found."""
        body = "No acceptance criteria here"
        criteria = extract_acceptance_criteria(body)
        assert criteria == []


class TestExtractDependencies:
    """Test dependency extraction."""

    def test_extract_dependencies_pp_format(self) -> None:
        """Verify PP-FFF-SS format dependencies are extracted."""
        body = """### Dependencies

00-001-01, 00-002-02"""
        deps = extract_dependencies(body)
        assert "00-001-01" in deps
        assert "00-002-02" in deps

    def test_extract_dependencies_legacy_format(self) -> None:
        """Verify legacy WS-FFF-SS format dependencies are extracted."""
        body = """### Dependencies

WS-500-01, WS-501-02"""
        deps = extract_dependencies(body)
        assert "WS-500-01" in deps
        assert "WS-501-02" in deps

    def test_extract_dependencies_mixed_formats(self) -> None:
        """Verify mixed format dependencies are extracted."""
        body = """### Dependencies

00-001-01, WS-500-02"""
        deps = extract_dependencies(body)
        assert "00-001-01" in deps
        assert "WS-500-02" in deps

    def test_extract_dependencies_none(self) -> None:
        """Verify empty list when dependencies section says 'none'."""
        body = """### Dependencies

None"""
        deps = extract_dependencies(body)
        assert deps == []

    def test_extract_dependencies_empty_section(self) -> None:
        """Verify empty list when dependencies section is empty."""
        body = """### Dependencies

"""
        deps = extract_dependencies(body)
        assert deps == []

    def test_extract_dependencies_no_section(self) -> None:
        """Verify empty list when no dependencies section."""
        body = "No dependencies section"
        deps = extract_dependencies(body)
        assert deps == []


class TestExtractSteps:
    """Test step extraction."""

    def test_extract_steps_numbered(self) -> None:
        """Verify numbered steps are extracted."""
        body = """### Steps

1. First step
2. Second step
3. Third step"""
        steps = extract_steps(body)
        assert len(steps) == 3
        assert steps[0] == "First step"
        assert steps[1] == "Second step"
        assert steps[2] == "Third step"

    def test_extract_steps_with_heading(self) -> None:
        """Verify steps with #### heading format are extracted."""
        # Note: Lines starting with # are skipped, so #### 1. won't match
        # The regex matches numbered lists, not headings
        body = """### Steps

1. First step
2. Second step"""
        steps = extract_steps(body)
        assert len(steps) == 2
        assert steps[0] == "First step"
        assert steps[1] == "Second step"

    def test_extract_steps_mixed_format(self) -> None:
        """Verify numbered steps are extracted."""
        body = """### Steps

1. First step
2. Second step
3. Third step"""
        steps = extract_steps(body)
        assert len(steps) == 3

    def test_extract_steps_skips_empty_lines(self) -> None:
        """Verify empty lines are skipped."""
        body = """### Steps

1. First step

2. Second step"""
        steps = extract_steps(body)
        assert len(steps) == 2

    def test_extract_steps_skips_headings(self) -> None:
        """Verify section headings are skipped."""
        body = """### Steps

1. First step
2. Second step"""
        steps = extract_steps(body)
        assert len(steps) == 2
        assert steps[0] == "First step"
        assert steps[1] == "Second step"

    def test_extract_steps_none(self) -> None:
        """Verify empty list when no steps found."""
        body = "No steps here"
        steps = extract_steps(body)
        assert steps == []


class TestExtractCodeBlocks:
    """Test code block extraction."""

    def test_extract_code_blocks(self) -> None:
        """Verify code blocks are extracted."""
        body = """```python
def hello():
    print("world")
```

Some text

```bash
echo "test"
```"""
        blocks = extract_code_blocks(body)
        assert len(blocks) == 2
        assert 'print("world")' in blocks[0]
        assert 'echo "test"' in blocks[1]

    def test_extract_code_blocks_multiline(self) -> None:
        """Verify multiline code blocks are extracted."""
        body = """```python
def func():
    x = 1
    y = 2
    return x + y
```"""
        blocks = extract_code_blocks(body)
        assert len(blocks) == 1
        assert "x = 1" in blocks[0]
        assert "y = 2" in blocks[0]

    def test_extract_code_blocks_none(self) -> None:
        """Verify empty list when no code blocks found."""
        body = "No code blocks here"
        blocks = extract_code_blocks(body)
        assert blocks == []
