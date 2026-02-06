"""Guard data models."""

from dataclasses import dataclass


@dataclass
class GuardResult:
    """Result of pre-edit guard check."""

    allowed: bool
    ws_id: str | None
    reason: str
    scope_files: list[str]
