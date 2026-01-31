"""Tests for ws_template_checker validator."""

from pathlib import Path

import pytest

from sdp.validators.ws_template_checker import (
    ValidationResult,
    format_result,
    validate_ws_structure,
)


class TestValidateWsStructure:
    """Tests for validate_ws_structure function."""

    def test_file_not_found_returns_violation(self) -> None:
        """Non-existent file returns violation."""
        result = validate_ws_structure(Path("/nonexistent/ws.md"))
        assert not result.passed
        assert "not found" in result.violations[0].lower()

    def test_valid_short_ws_passes(self, tmp_path: Path) -> None:
        """Short WS with required sections passes."""
        ws = tmp_path / "00-032-01.md"
        ws.write_text("""---
ws_id: 00-032-01
---

## Goal
Target

## Acceptance Criteria
- AC1: Something

## Contract
```python
raise NotImplementedError
```

## Scope
Small

## Verification
```bash
pytest tests/
```
""")
        result = validate_ws_structure(ws)
        assert result.passed
        assert result.ws_id == "00-032-01"

    def test_missing_section_fails(self, tmp_path: Path) -> None:
        """Missing required section fails."""
        ws = tmp_path / "bad.md"
        ws.write_text("---\nws_id: bad\n---\n\n## Goal\nOnly goal")
        result = validate_ws_structure(ws)
        assert not result.passed
        assert any("Missing" in v for v in result.violations)

    def test_too_long_ws_fails(self, tmp_path: Path) -> None:
        """WS exceeding MAX_WS_LINES fails."""
        ws = tmp_path / "long.md"
        content = "---\nws_id: long\n---\n\n## Goal\n## Acceptance Criteria\n## Contract\n## Scope\n## Verification\n"
        content += "\n".join([f"Line {i}" for i in range(200)])
        ws.write_text(content)
        result = validate_ws_structure(ws)
        assert not result.passed
        assert any("too long" in v.lower() for v in result.violations)


class TestFormatResult:
    """Tests for format_result function."""

    def test_passed_format(self) -> None:
        """Passed result shows PASSED."""
        result = ValidationResult(ws_id="x", passed=True, violations=[], warnings=[])
        output = format_result(result)
        assert "PASSED" in output
        assert "x" in output

    def test_failed_format(self) -> None:
        """Failed result shows violations."""
        result = ValidationResult(
            ws_id="x", passed=False, violations=["Missing section"], warnings=[]
        )
        output = format_result(result)
        assert "FAILED" in output
        assert "Missing section" in output
