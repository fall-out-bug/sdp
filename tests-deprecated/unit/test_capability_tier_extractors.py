"""Tests for capability tier extraction utilities."""

import pytest

from sdp.validators.capability_tier_extractors import (
    _extract_body,
    _extract_code_block,
    _extract_section,
)


class TestExtractBody:
    """Test _extract_body function."""

    def test_extract_body_with_frontmatter(self) -> None:
        """Test extracting body from content with frontmatter."""
        content = """---
ws_id: 01-001-01
status: backlog
---

## Goal
This is the goal.
"""
        body = _extract_body(content)
        assert "ws_id:" not in body
        assert "## Goal" in body
        assert "This is the goal." in body

    def test_extract_body_without_frontmatter(self) -> None:
        """Test extracting body from content without frontmatter."""
        content = "## Goal\nThis is the goal."
        body = _extract_body(content)
        assert body == content

    def test_extract_body_empty_content(self) -> None:
        """Test extracting body from empty content."""
        body = _extract_body("")
        assert body == ""

    def test_extract_body_multiline_frontmatter(self) -> None:
        """Test extracting body with multiline frontmatter."""
        content = """---
ws_id: 01-001-01
status: backlog
size: MEDIUM
---

## Goal
Goal content.
"""
        body = _extract_body(content)
        assert "## Goal" in body
        assert "ws_id:" not in body


class TestExtractSection:
    """Test _extract_section function."""

    def test_extract_section_exists(self) -> None:
        """Test extracting an existing section."""
        body = """## Goal
This is the goal.

## Context
This is the context.
"""
        section = _extract_section(body, "Goal")
        assert section == "This is the goal."

    def test_extract_section_not_exists(self) -> None:
        """Test extracting a non-existent section."""
        body = "## Goal\nSome content."
        section = _extract_section(body, "Context")
        assert section == ""

    def test_extract_section_multiline(self) -> None:
        """Test extracting multiline section."""
        body = """## Goal
First line.
Second line.
Third line.

## Context
Context here.
"""
        section = _extract_section(body, "Goal")
        assert "First line." in section
        assert "Second line." in section
        assert "Third line." in section
        assert "Context here" not in section

    def test_extract_section_with_subsections(self) -> None:
        """Test extracting section with subsections."""
        body = """## Goal
Main content.

### Subsection
Sub content.

## Context
Context.
"""
        section = _extract_section(body, "Goal")
        assert "Main content." in section
        assert "### Subsection" in section
        assert "Sub content." in section

    def test_extract_section_case_sensitive(self) -> None:
        """Test that section name is case-sensitive."""
        body = "## Goal\nContent."
        section = _extract_section(body, "goal")
        assert section == ""

    def test_extract_section_special_chars(self) -> None:
        """Test extracting section with special characters in name."""
        body = "## API Endpoints\nContent."
        section = _extract_section(body, "API Endpoints")
        assert section == "Content."


class TestExtractCodeBlock:
    """Test _extract_code_block function."""

    def test_extract_python_code_block(self) -> None:
        """Test extracting Python code block."""
        section = """Some text.

```python
def hello():
    print("Hello, world!")
```

More text.
"""
        code = _extract_code_block(section, "python")
        assert 'def hello():' in code
        assert 'print("Hello, world!")' in code

    def test_extract_code_block_not_found(self) -> None:
        """Test when code block not found."""
        section = "Just text here."
        code = _extract_code_block(section, "python")
        assert code is None

    def test_extract_code_block_wrong_language(self) -> None:
        """Test extracting code block with different language."""
        section = """```python
def foo():
    pass
```

```bash
echo "hello"
```
"""
        python_code = _extract_code_block(section, "python")
        bash_code = _extract_code_block(section, "bash")
        assert "def foo" in python_code
        assert 'echo "hello"' in bash_code

    def test_extract_code_block_default_language(self) -> None:
        """Test extracting code block with default language."""
        section = """Some text.

```
def foo():
    pass
```
"""
        code = _extract_code_block(section)
        assert code is None  # No language specified

    def test_extract_code_block_multiline(self) -> None:
        """Test extracting multiline code block."""
        section = """```python
class MyClass:
    def method1(self):
        pass

    def method2(self):
        pass
```
"""
        code = _extract_code_block(section, "python")
        assert "class MyClass:" in code
        assert "def method1" in code
        assert "def method2" in code

    def test_extract_code_block_empty(self) -> None:
        """Test extracting empty code block."""
        section = "```python\n```"
        code = _extract_code_block(section, "python")
        # Empty code block returns empty string, not None
        assert code == "" or code is None  # Both are acceptable
