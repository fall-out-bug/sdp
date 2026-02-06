"""Tests for skill validator."""
from pathlib import Path

import pytest
from click.testing import CliRunner

from sdp.cli.skill import skill

runner = CliRunner()


def test_validate_skill_valid(tmp_path: Path) -> None:
    """Test validation passes for valid skill."""
    skill_file = tmp_path / "SKILL.md"
    skill_file.write_text("""---
name: test
description: Test skill
tools: Read, Write
---

# @test - Test Skill

This is a test skill for validation.

## Quick Reference

| Step | Action | Gate |
|------|--------|------|
| 1 | Do something | Something done |

## Workflow

### Step 1: Do something

Do the thing.

## See Also

- [Reference](../docs/reference/test.md)
""")

    result = runner.invoke(skill, ["validate", str(skill_file)])
    assert result.exit_code == 0
    assert "✅" in result.output


def test_validate_skill_too_long(tmp_path: Path) -> None:
    """Test validation warns for long skills."""
    skill_file = tmp_path / "SKILL.md"
    content = """---
name: test
---

# Test Skill

## Quick Reference

| Step | Action | Gate |
|------|--------|------|
| 1 | Do something | Something done |

## Workflow

### Step 1

"""
    content += "\n".join([f"Line {i}" for i in range(120)])
    content += "\n\n## See Also\n\n- [Reference](../docs/reference/test.md)"
    skill_file.write_text(content)

    result = runner.invoke(skill, ["validate", str(skill_file)])
    assert result.exit_code == 0
    assert "⚠️" in result.output
    assert "Consider shortening" in result.output


def test_validate_skill_missing_sections(tmp_path: Path) -> None:
    """Test validation fails for missing sections."""
    skill_file = tmp_path / "SKILL.md"
    skill_file.write_text("""---
name: test
---

# Test

No required sections.
""")

    result = runner.invoke(skill, ["validate", str(skill_file)])
    assert result.exit_code == 1
    assert "❌" in result.output
    assert "Missing section" in result.output


def test_validate_skill_not_found() -> None:
    """Test validation fails for non-existent file."""
    result = runner.invoke(skill, ["validate", "/nonexistent/file.md"])
    assert result.exit_code != 0
