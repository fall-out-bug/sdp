"""Quality gate violation model."""

from dataclasses import dataclass


@dataclass
class QualityGateViolation:
    """Represents a single quality gate violation."""

    category: str
    file_path: str
    line_number: int | None
    message: str
    severity: str = "error"

    def __str__(self) -> str:
        """Format violation for display."""
        loc = f"{self.file_path}:{self.line_number}" if self.line_number else self.file_path
        return f"[{self.severity.upper()}] {loc} - {self.category}: {self.message}"
