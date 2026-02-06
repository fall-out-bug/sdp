"""Tests for dependency security documentation (AC6: 00-025-01)."""

from pathlib import Path


def test_development_md_has_dependency_security_section() -> None:
    """AC6: Documentation updated (dependency security policy)."""
    dev_md = Path("docs/internals/development.md")
    assert dev_md.exists(), "docs/internals/development.md must exist"
    content = dev_md.read_text()
    assert "Dependency Management & Security" in content
    assert "pip-audit" in content
    assert "Dependabot" in content
