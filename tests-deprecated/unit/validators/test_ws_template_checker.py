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

    def test_format_with_warnings(self) -> None:
        """Format result includes warnings."""
        result = ValidationResult(
            ws_id="x",
            passed=True,
            violations=[],
            warnings=["Large code block without NotImplementedError"],
        )
        output = format_result(result)
        assert "Warnings:" in output
        assert "Large code block" in output

    def test_extract_ws_id_unknown(self, tmp_path: Path) -> None:
        """Test ws_id extraction returns 'unknown' when not found."""
        ws = tmp_path / "no_id.md"
        ws.write_text("---\n---\n\n## Goal\nTest")
        result = validate_ws_structure(ws)
        assert result.ws_id == "unknown"

    def test_code_block_too_large_violation(self, tmp_path: Path) -> None:
        """Test code block exceeding MAX_CODE_BLOCK_LINES."""
        ws = tmp_path / "large.md"
        large_code = "\n".join([f"line {i}" for i in range(35)])
        ws.write_text(
            f"""---
ws_id: large
---

## Goal
Test

## Acceptance Criteria
- AC1

## Contract
```python
{large_code}
```

## Scope
Small

## Verification
```bash
pytest
```
"""
        )
        result = validate_ws_structure(ws)
        assert not result.passed
        assert any("too large" in v.lower() for v in result.violations)

    def test_large_code_block_without_notimplemented_warning(self, tmp_path: Path) -> None:
        """Test large code block without NotImplementedError generates warning."""
        ws = tmp_path / "warning.md"
        large_code = "\n".join([f"line {i}" for i in range(15)])
        ws.write_text(
            f"""---
ws_id: warning
---

## Goal
Test

## Acceptance Criteria
- AC1

## Contract
```python
{large_code}
```

## Scope
Small

## Verification
```bash
pytest
```
"""
        )
        result = validate_ws_structure(ws)
        assert any("without" in w.lower() and "notimplementederror" in w.lower() for w in result.warnings)

    def test_code_block_with_notimplemented_no_warning(self, tmp_path: Path) -> None:
        """Test large code block with NotImplementedError doesn't generate warning."""
        ws = tmp_path / "no_warning.md"
        large_code = "\n".join([f"line {i}" for i in range(15)])
        ws.write_text(
            f"""---
ws_id: no_warning
---

## Goal
Test

## Acceptance Criteria
- AC1

## Contract
```python
{large_code}
raise NotImplementedError
```

## Scope
Small

## Verification
```bash
pytest
```
"""
        )
        result = validate_ws_structure(ws)
        # Should not have warning about missing NotImplementedError
        assert not any("without" in w.lower() and "notimplementederror" in w.lower() for w in result.warnings)

    def test_multiple_code_blocks(self, tmp_path: Path) -> None:
        """Test validation with multiple code blocks."""
        ws = tmp_path / "multi.md"
        ws.write_text(
            """---
ws_id: multi
---

## Goal
Test

## Acceptance Criteria
- AC1

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

## Another Section
```python
def helper():
    raise NotImplementedError
```
"""
        )
        result = validate_ws_structure(ws)
        # Should handle multiple code blocks
        assert result.ws_id == "multi"

    def test_cli_main_no_args(self) -> None:
        """Test CLI main with no args returns 1."""
        import sys
        from unittest.mock import patch

        with patch.object(sys, "argv", ["ws_template_checker"]):
            from sdp.validators.ws_template_checker import main

            assert main() == 1

    def test_cli_main_with_file(self, tmp_path: Path) -> None:
        """Test CLI main with valid file."""
        import sys
        from unittest.mock import patch

        ws = tmp_path / "valid.md"
        ws.write_text(
            """---
ws_id: valid
---

## Goal
Test

## Acceptance Criteria
- AC1

## Contract
```python
raise NotImplementedError
```

## Scope
Small

## Verification
```bash
pytest
```
"""
        )
        with patch.object(sys, "argv", ["ws_template_checker", str(ws)]):
            from sdp.validators.ws_template_checker import main

            assert main() == 0

    def test_cli_main_with_invalid_file(self, tmp_path: Path) -> None:
        """Test CLI main with invalid file returns 1."""
        import sys
        from unittest.mock import patch

        ws = tmp_path / "invalid.md"
        ws.write_text("---\nws_id: invalid\n---\n\n## Goal\nOnly goal")
        with patch.object(sys, "argv", ["ws_template_checker", str(ws)]):
            from sdp.validators.ws_template_checker import main

            assert main() == 1

    def test_cli_entry_point(self, tmp_path: Path) -> None:
        """Test CLI entry point via __main__."""
        import subprocess
        import sys

        ws = tmp_path / "cli_test.md"
        ws.write_text(
            """---
ws_id: cli_test
---

## Goal
Test

## Acceptance Criteria
- AC1

## Contract
```python
raise NotImplementedError
```

## Scope
Small

## Verification
```bash
pytest
```
"""
        )
        result = subprocess.run(
            [sys.executable, "-m", "sdp.validators.ws_template_checker", str(ws)],
            capture_output=True,
            text=True,
        )
        assert result.returncode == 0
