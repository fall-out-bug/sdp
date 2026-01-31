"""Workstream template validator - ensures WS follow simplified template."""

import re
from dataclasses import dataclass
from pathlib import Path


@dataclass
class ValidationResult:
    """Result of WS template validation."""

    ws_id: str
    passed: bool
    violations: list[str]
    warnings: list[str]


# Configuration
MAX_WS_LINES: int = 150
MAX_CODE_BLOCK_LINES: int = 30
REQUIRED_SECTIONS: list[str] = ["Goal", "Acceptance Criteria", "Contract", "Scope", "Verification"]


def _extract_ws_id(lines: list[str]) -> str:
    """Extract ws_id from frontmatter."""
    for line in lines[:20]:
        if line.startswith("ws_id:"):
            return line.split(":", 1)[1].strip()
    return "unknown"


def _check_required_sections(lines: list[str]) -> list[str]:
    """Check all required sections present."""
    violations: list[str] = []
    for section in REQUIRED_SECTIONS:
        pattern = re.compile(rf"^##+ {section}", re.IGNORECASE)
        if not any(pattern.match(line) for line in lines):
            violations.append(f"Missing required section: {section}")
    return violations


def _check_code_blocks(lines: list[str]) -> tuple[list[str], list[str]]:
    """Check code block sizes. Returns (violations, warnings)."""
    violations: list[str] = []
    warnings: list[str] = []
    in_code_block = False
    code_block_lines = 0
    code_block_start = 0
    has_not_implemented = False

    for i, line in enumerate(lines, 1):
        if line.strip().startswith("```"):
            if in_code_block:
                if code_block_lines > MAX_CODE_BLOCK_LINES:
                    msg = (
                        f"Code block too large at line {code_block_start}: "
                        f"{code_block_lines} lines (max {MAX_CODE_BLOCK_LINES})"
                    )
                    violations.append(msg)
                if code_block_lines > 10 and not has_not_implemented:
                    msg = (
                        f"Large code block at line {code_block_start} without "
                        "'raise NotImplementedError' - may contain full implementation"
                    )
                    warnings.append(msg)
                in_code_block = False
                code_block_lines = 0
                has_not_implemented = False
            else:
                in_code_block = True
                code_block_start = i
        elif in_code_block:
            code_block_lines += 1
            if "NotImplementedError" in line or "raise NotImplementedError" in line:
                has_not_implemented = True

    return violations, warnings


def validate_ws_structure(ws_path: Path) -> ValidationResult:
    """Validate WS follows simplified template.

    Checks:
    - Total lines ≤ MAX_WS_LINES
    - Code blocks ≤ MAX_CODE_BLOCK_LINES
    - All required sections present
    - No full implementations (detect by line count + no raise NotImplementedError)

    Args:
        ws_path: Path to workstream file

    Returns:
        ValidationResult with pass/fail and violations
    """
    if not ws_path.exists():
        return ValidationResult(
            ws_id=ws_path.stem, passed=False, violations=[f"File not found: {ws_path}"], warnings=[]
        )

    content = ws_path.read_text(encoding="utf-8")
    lines = content.splitlines()
    ws_id = _extract_ws_id(lines)
    violations: list[str] = []
    warnings: list[str] = []

    total_lines = len(lines)
    if total_lines > MAX_WS_LINES:
        violations.append(f"Workstream too long: {total_lines} lines (max {MAX_WS_LINES})")
    violations.extend(_check_required_sections(lines))
    cb_violations, cb_warnings = _check_code_blocks(lines)
    violations.extend(cb_violations)
    warnings.extend(cb_warnings)

    passed = len(violations) == 0
    return ValidationResult(ws_id=ws_id, passed=passed, violations=violations, warnings=warnings)


def format_result(result: ValidationResult) -> str:
    """Format validation result for display.

    Args:
        result: Validation result

    Returns:
        Formatted string
    """
    lines = [f"Workstream: {result.ws_id}"]

    if result.passed:
        lines.append("✅ PASSED")
    else:
        lines.append("❌ FAILED")

    if result.violations:
        lines.append("\nViolations:")
        for v in result.violations:
            lines.append(f"  - {v}")

    if result.warnings:
        lines.append("\nWarnings:")
        for w in result.warnings:
            lines.append(f"  - {w}")

    return "\n".join(lines)


def main() -> int:
    """CLI entry point."""
    import sys

    if len(sys.argv) < 2:
        print("Usage: python -m sdp.validators.ws_template_checker <ws_file>")
        return 1

    ws_path = Path(sys.argv[1])
    result = validate_ws_structure(ws_path)
    print(format_result(result))

    return 0 if result.passed else 1


if __name__ == "__main__":
    exit(main())
