"""Time estimate checker - validates against forbidden time-based estimates.

Enforces SDP rule: No time-based estimates in workstreams, skills, or templates.
"""

import re
from dataclasses import dataclass
from pathlib import Path


@dataclass
class Violation:
    """Time estimate violation."""

    file: Path
    line_number: int
    line_content: str
    pattern: str
    message: str


# Regex patterns for forbidden time estimates
TIME_PATTERNS: list[str] = [
    r"estimated_duration[:\s]*[\"']?[\w\s-]+[\"']?",  # estimated_duration: "2-3 hours"
    r"estimated_loc[:\s]*\d+",  # estimated_loc: 450
    r"\d+[-–]\d+\s+(hours?|days?|weeks?|months?)\b",  # 2-3 hours, 1-2 days
    r"\b\d+\s+(hours?|days?|weeks?|months?)\b",  # 2 hours, 3 days
    r"\b\d+h\b(?!\))",  # 2h (but not in format strings like :03d})
    r"takes?\s+\d+",  # takes 2, take 3
    r"requires?\s+\d+",  # requires 2, require 3
]

# Compiled patterns for performance
_COMPILED_PATTERNS = [(re.compile(p, re.IGNORECASE), p) for p in TIME_PATTERNS]

# Allowed exceptions (in comments or documentation about rules)
ALLOWED_CONTEXTS = [
    "forbidden",
    "don't use",
    "do not use",
    "remove",
    "violation",
    "anti-pattern",
    "example of bad",
    "bad:",
    "❌",
    ":03d",  # Format strings like ws.sequence:03d
    "ws-{",  # Workstream ID formatting
]


def _is_allowed_context(line: str) -> bool:
    """Check if line is documenting violations (not an actual violation)."""
    line_lower = line.lower()
    return any(ctx in line_lower for ctx in ALLOWED_CONTEXTS)


def check_file(path: Path) -> list[Violation]:
    """Check file for time estimate violations.

    Args:
        path: File to check

    Returns:
        List of violations with line numbers
    """
    violations: list[Violation] = []

    if not path.exists():
        return violations

    try:
        content = path.read_text(encoding="utf-8")
    except (UnicodeDecodeError, PermissionError):
        return violations

    for line_num, line in enumerate(content.splitlines(), 1):
        # Skip lines documenting violations
        if _is_allowed_context(line):
            continue

        # Check against all patterns
        for pattern_re, pattern_str in _COMPILED_PATTERNS:
            if pattern_re.search(line):
                violations.append(
                    Violation(
                        file=path,
                        line_number=line_num,
                        line_content=line.strip(),
                        pattern=pattern_str,
                        message=f"Forbidden time estimate found (pattern: {pattern_str})",
                    )
                )
                break  # Only report first match per line

    return violations


def check_directory(path: Path, glob: str = "**/*.md") -> list[Violation]:
    """Check all matching files in directory.

    Args:
        path: Directory to search
        glob: Glob pattern for files (default: **/*.md)

    Returns:
        List of all violations found
    """
    all_violations: list[Violation] = []

    if not path.exists() or not path.is_dir():
        return all_violations

    for file_path in path.glob(glob):
        if file_path.is_file():
            violations = check_file(file_path)
            all_violations.extend(violations)

    return all_violations


def format_violations(violations: list[Violation]) -> str:
    """Format violations for display.

    Args:
        violations: List of violations

    Returns:
        Formatted string for console output
    """
    if not violations:
        return "✅ No time estimate violations found"

    output = [f"❌ Found {len(violations)} time estimate violation(s):\n"]

    # Group by file
    by_file: dict[Path, list[Violation]] = {}
    for v in violations:
        by_file.setdefault(v.file, []).append(v)

    for file_path, file_violations in sorted(by_file.items()):
        output.append(f"\n{file_path}:")
        for v in file_violations:
            output.append(f"  Line {v.line_number}: {v.line_content}")
            output.append(f"    → {v.message}")

    return "\n".join(output)


def main() -> int:
    """CLI entry point."""
    import sys

    if len(sys.argv) < 2:
        print("Usage: python -m sdp.validators.time_estimate_checker <path>")
        print("  <path>: File or directory to check")
        return 1

    target = Path(sys.argv[1])

    if not target.exists():
        print(f"❌ Path not found: {target}")
        return 1

    violations: list[Violation]
    if target.is_file():
        violations = check_file(target)
    else:
        violations = check_directory(target)

    print(format_violations(violations))

    return 1 if violations else 0


if __name__ == "__main__":
    exit(main())
