"""Targeted tests for uncovered lines in validators/ws_completion/parser.py."""

from pathlib import Path

from sdp.validators.ws_completion.parser import (
    parse_frontmatter_scope,
    parse_verification_commands,
)


def test_parse_frontmatter_scope_not_in_frontmatter() -> None:
    """Test parse_frontmatter_scope with scope_files outside frontmatter (lines 59-60)."""
    content = """---
title: Test
---

scope_files:
  - file1.py
  - file2.py
"""
    result = parse_frontmatter_scope(content)
    # Should not parse scope_files outside frontmatter
    assert result == []


def test_parse_verification_commands_bash_block() -> None:
    """Test parse_verification_commands with bash code block (lines 81-82)."""
    content = """
### Verification

```bash
pytest tests/test_foo.py
mypy src/
```
"""
    result = parse_verification_commands(content)
    assert result == ["pytest tests/test_foo.py", "mypy src/"]


def test_parse_verification_commands_sh_block() -> None:
    """Test parse_verification_commands with sh code block (lines 81-82)."""
    content = """
### Verification

```sh
pytest tests/test_foo.py
mypy src/
```
"""
    result = parse_verification_commands(content)
    assert result == ["pytest tests/test_foo.py", "mypy src/"]


def test_parse_verification_commands_with_comments() -> None:
    """Test parse_verification_commands skips comments (line 89)."""
    content = """
### Verification

```bash
# Run tests
pytest tests/test_foo.py
# Type check
mypy src/
```
"""
    result = parse_verification_commands(content)
    # Should exclude comments
    assert result == ["pytest tests/test_foo.py", "mypy src/"]


def test_parse_verification_commands_stops_at_next_section() -> None:
    """Test parse_verification_commands stops at next section (lines 87-88)."""
    content = """
### Verification

```bash
pytest tests/test_foo.py
```

## Next Section

```bash
should_not_be_parsed
```
"""
    result = parse_verification_commands(content)
    assert result == ["pytest tests/test_foo.py"]
    assert "should_not_be_parsed" not in result


def test_parse_verification_commands_empty_lines() -> None:
    """Test parse_verification_commands skips empty lines (line 89)."""
    content = """
### Verification

```bash
pytest tests/test_foo.py

mypy src/
```
"""
    result = parse_verification_commands(content)
    assert result == ["pytest tests/test_foo.py", "mypy src/"]
