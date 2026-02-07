"""Data models for workstream completion verification."""

from dataclasses import dataclass


@dataclass
class CheckResult:
    """Result of a single verification check."""

    name: str
    passed: bool
    message: str
    evidence: str | None  # Command output or file path


@dataclass
class VerificationResult:
    """Result of full WS verification."""

    ws_id: str
    passed: bool
    checks: list[CheckResult]
    coverage_actual: float | None
    missing_files: list[str]
    failed_commands: list[str]
