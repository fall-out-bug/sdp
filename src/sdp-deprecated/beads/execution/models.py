"""
Data models for execution modes and results.

Defines core types for @oneshot execution configuration and results.
"""

from dataclasses import dataclass, field
from enum import Enum
from typing import List, Optional


class ExecutionMode(Enum):
    """Execution mode for @oneshot."""

    STANDARD = "standard"
    AUTO_APPROVE = "auto_approve"
    SANDBOX = "sandbox"
    DRY_RUN = "dry_run"

    @property
    def requires_pr(self) -> bool:
        """Whether this mode requires PR approval."""
        return self == ExecutionMode.STANDARD

    @property
    def allows_production(self) -> bool:
        """Whether this mode allows production deployment."""
        return self in (ExecutionMode.STANDARD, ExecutionMode.AUTO_APPROVE)

    @property
    def is_sandbox(self) -> bool:
        """Whether this is a sandbox-only deployment."""
        return self == ExecutionMode.SANDBOX

    @property
    def is_preview(self) -> bool:
        """Whether this is a preview-only mode."""
        return self == ExecutionMode.DRY_RUN


@dataclass
class OneshotResult:
    """Result of oneshot feature execution."""

    success: bool
    feature_id: str
    total_executed: int = 0
    error: Optional[str] = None
    failed_tasks: List[str] = field(default_factory=list)
    mode: ExecutionMode = ExecutionMode.STANDARD
    deployment_target: str = "production"
    pr_created: bool = False
    preview_only: bool = False
    tasks_preview: List[str] = field(default_factory=list)
