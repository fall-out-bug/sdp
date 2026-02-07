"""Data models for GitHub sync operations."""

from dataclasses import dataclass
from typing import Optional


@dataclass
class SyncResult:
    """Result of syncing single workstream.

    Attributes:
        ws_id: Workstream ID
        action: Action taken ("created", "updated", "skipped", "failed")
        issue_number: GitHub issue number if created
        error: Error message if failed
    """

    ws_id: str
    action: str
    issue_number: Optional[int] = None
    error: Optional[str] = None
