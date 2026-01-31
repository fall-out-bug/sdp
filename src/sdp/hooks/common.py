"""Shared utilities for Git hooks."""

from dataclasses import dataclass
from pathlib import Path


@dataclass
class CheckResult:
    """Result of a quality check."""

    passed: bool
    message: str
    violations: list[tuple[Path, int | None, str]]  # (file, line, issue)

    def format_terminal(self) -> str:
        """Format result for terminal output."""
        if self.passed:
            return f"✓ {self.message}"
        output = [f"❌ {self.message}"]
        for file_path, line, issue in self.violations:
            line_str = str(line) if line is not None else "?"
            output.append(f"  {file_path}:{line_str} - {issue}")
        return "\n".join(output)
