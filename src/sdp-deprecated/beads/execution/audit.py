"""
Audit logging for execution tracking.

Tracks execution history for compliance and debugging.
"""

import json
from dataclasses import dataclass
from datetime import datetime, timezone
from typing import Any, List, Optional

from .models import ExecutionMode


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
        import getpass
        import os

        # Ensure audit directory exists
        os.makedirs(
            os.path.dirname(self.audit_file) if os.path.dirname(self.audit_file) else ".",
            exist_ok=True,
        )

        # Determine user
        if user is None:
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
