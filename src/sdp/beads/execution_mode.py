"""
Execution modes and related classes for @oneshot workflow efficiency feature.

Defines:
- ExecutionMode enum for different execution modes
- DestructiveOperationDetector for detecting dangerous operations
- AuditLogger for tracking --auto-approve executions
- OneshotResult dataclass for execution results
"""

import json
from dataclasses import dataclass, field
from datetime import datetime, timezone
from enum import Enum
from typing import Any, List, Optional


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
class DestructiveOperations:
    """Result of destructive operations check."""

    has_destructive_operations: bool
    operation_types: List[str]
    files_affected: List[str]
    details: Optional[str] = None


class DestructiveOperationDetector:
    """Detect destructive operations that require user confirmation."""

    # Patterns that indicate destructive operations
    DESTRUCTIVE_PATTERNS = {
        "database_migration": ["migration", "migrate", "schema", "alembic"],
        "file_deletion": ["delete", "remove", "rm"],
        "data_loss": ["drop", "truncate", "wipe"],
    }

    def check_operations(
        self,
        files_to_create: List[str],
        files_to_modify: List[str],
        files_to_delete: List[str],
    ) -> DestructiveOperations:
        """Check if operations are destructive.

        Args:
            files_to_create: List of file paths to be created
            files_to_modify: List of file paths to be modified
            files_to_delete: List of file paths to be deleted

        Returns:
            DestructiveOperations with check result
        """
        operation_types = []
        files_affected = []

        # Check file deletions (always destructive)
        if files_to_delete:
            operation_types.append("file_deletion")
            files_affected.extend(files_to_delete)

        # Check for database migrations in created files
        for file_path in files_to_create:
            if self._is_database_migration(file_path):
                operation_types.append("database_migration")
                files_affected.append(file_path)

        # Check for destructive patterns in modified files
        for file_path in files_to_modify:
            if self._has_destructive_changes(file_path):
                operation_types.append("data_loss")
                files_affected.append(file_path)

        return DestructiveOperations(
            has_destructive_operations=len(operation_types) > 0,
            operation_types=list(set(operation_types)),  # Deduplicate
            files_affected=list(set(files_affected)),  # Deduplicate
            details=f"Found {len(operation_types)} destructive operation types",
        )

    def _is_database_migration(self, file_path: str) -> bool:
        """Check if file is a database migration."""
        path_lower = file_path.lower()
        return any(
            pattern in path_lower
            for patterns in self.DESTRUCTIVE_PATTERNS.values()
            for pattern in patterns
            if pattern in path_lower
        )

    def _has_destructive_changes(self, file_path: str) -> bool:
        """Check if file modification involves destructive changes."""
        # Simplified check - in real implementation would read file content
        path_lower = file_path.lower()
        return any(
            pattern in path_lower
            for pattern in self.DESTRUCTIVE_PATTERNS.get("data_loss", [])
        )


@dataclass
class AuditLogEntry:
    """Single audit log entry."""

    timestamp: str
    user: str
    feature: str
    mode: str
    workstreams_executed: int
    result: str
    deployment_target: Optional[str] = None

    def to_dict(self) -> dict[str, Any]:
        """Convert to dictionary for JSON serialization."""
        return {
            "timestamp": self.timestamp,
            "user": self.user,
            "feature": self.feature,
            "mode": self.mode,
            "workstreams_executed": self.workstreams_executed,
            "result": self.result,
            "deployment_target": self.deployment_target,
        }


class AuditLogger:
    """Audit logger for tracking --auto-approve executions."""

    def __init__(self, audit_file: str = ".sdp/audit.log"):
        """Initialize audit logger.

        Args:
            audit_file: Path to audit log file
        """
        self.audit_file = audit_file

    def log_execution(
        self,
        feature_id: str,
        mode: ExecutionMode,
        workstreams_executed: int,
        result: str,
        user: Optional[str] = None,
        deployment_target: Optional[str] = None,
    ) -> None:
        """Log an execution to audit file.

        Args:
            feature_id: Feature task ID
            mode: Execution mode used
            workstreams_executed: Number of workstreams executed
            result: Execution result ("success" or "failure")
            user: User who initiated execution (optional)
            deployment_target: Deployment target ("production" or "sandbox")
        """
        import os

        # Ensure audit directory exists
        os.makedirs(os.path.dirname(self.audit_file) if os.path.dirname(self.audit_file) else ".", exist_ok=True)  # noqa: E501

        # Determine user
        if user is None:
            import getpass
            user = getpass.getuser()

        # Create log entry
        entry = AuditLogEntry(
            timestamp=datetime.now(timezone.utc).isoformat(),
            user=user,
            feature=feature_id,
            mode=mode.value,
            workstreams_executed=workstreams_executed,
            result=result,
            deployment_target=deployment_target,
        )

        # Append to audit log
        with open(self.audit_file, "a") as f:
            f.write(json.dumps(entry.to_dict()) + "\n")

    def read_recent(self, count: int = 10) -> List[dict[str, Any]]:
        """Read recent audit log entries.

        Args:
            count: Number of recent entries to read

        Returns:
            List of audit log entries (most recent last)
        """
        try:
            with open(self.audit_file, "r") as f:
                lines = f.readlines()

            # Get last N lines
            recent_lines = lines[-count:] if len(lines) >= count else lines

            return [json.loads(line.strip()) for line in recent_lines if line.strip()]

        except FileNotFoundError:
            return []


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
