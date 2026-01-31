"""BugReportFlow for bug report creation and /issue skill integration.

Creates bug reports with automatic severity classification and routes
to /hotfix, /bugfix, or /issue skills based on severity.
"""

import logging
import uuid
from dataclasses import dataclass, field
from enum import Enum
from typing import Optional

logger = logging.getLogger(__name__)


class BugSeverity(Enum):
    """Bug severity levels."""

    P0 = "P0"
    P1 = "P1"
    P2 = "P2"
    P3 = "P3"

    @property
    def label(self) -> str:
        """Get skill label for severity."""
        if self == BugSeverity.P0:
            return "hotfix"
        elif self in [BugSeverity.P1, BugSeverity.P2]:
            return "bugfix"
        else:
            return "backlog"


class BugStatus(Enum):
    """Bug status values."""

    OPEN = "open"
    IN_PROGRESS = "in_progress"
    RESOLVED = "resolved"


@dataclass
class BugReport:
    """Bug report with metadata.

    Attributes:
        id: Unique bug report ID
        title: Bug title
        description: Bug description
        severity: Bug severity level
        status: Bug status
        workstream_id: Optional associated workstream ID
    """
    id: str = field(default_factory=lambda: str(uuid.uuid4()))
    title: str = ""
    description: str = ""
    severity: BugSeverity = BugSeverity.P3
    status: BugStatus = BugStatus.OPEN
    workstream_id: Optional[str] = None


class BugReportFlow:
    """Flow for creating and routing bug reports with auto-classification."""

    # Keywords for severity classification
    P0_KEYWORDS = ["outage", "crash", "critical", "production", "down"]
    P1_KEYWORDS = [
        "security", "data loss", "corruption", "injection",
        "login", "auth", "fail", "cannot", "unable"
    ]
    P2_KEYWORDS = ["ui", "glitch", "slow", "performance", "error"]
    P3_KEYWORDS = ["typo", "documentation", "cosmetic", "minor"]

    def __init__(self) -> None:
        """Initialize bug report flow."""
        self._reports: list[BugReport] = []

    def create_report(
        self,
        title: str,
        description: str,
        severity: Optional[BugSeverity] = None,
        workstream_id: Optional[str] = None,
    ) -> BugReport:
        """Create bug report with auto-classification if severity not provided."""
        if severity is None:
            severity = self.classify_severity(title, description)

        report = BugReport(
            title=title,
            description=description,
            severity=severity,
            workstream_id=workstream_id,
        )

        self._reports.append(report)
        logger.info(f"Created bug report: {title} ({severity.value})")

        return report

    def classify_severity(self, title: str, description: str) -> BugSeverity:
        """Classify bug severity based on title and description keywords."""
        text = f"{title} {description}".lower()

        if any(keyword in text for keyword in self.P0_KEYWORDS):
            return BugSeverity.P0
        if any(keyword in text for keyword in self.P1_KEYWORDS):
            return BugSeverity.P1
        if any(keyword in text for keyword in self.P2_KEYWORDS):
            return BugSeverity.P2

        return BugSeverity.P3

    def get_reports(self) -> list[BugReport]:
        """Get all bug reports."""
        return self._reports.copy()

    def get_reports_by_severity(self, severity: BugSeverity) -> list[BugReport]:
        """Get reports by severity."""
        return [r for r in self._reports if r.severity == severity]

    def get_reports_by_workstream(self, workstream_id: str) -> list[BugReport]:
        """Get reports by workstream."""
        return [r for r in self._reports if r.workstream_id == workstream_id]

    def get_reports_by_status(self, status: BugStatus) -> list[BugReport]:
        """Get reports by status."""
        return [r for r in self._reports if r.status == status]

    def get_blocking_bugs(self) -> list[BugReport]:
        """Get blocking bugs (P0 and P1)."""
        return [
            r for r in self._reports
            if r.severity in [BugSeverity.P0, BugSeverity.P1]
        ]

    def get_blocking_workstreams(self) -> list[str]:
        """Get workstreams blocked by bugs."""
        workstreams = set()
        for report in self.get_blocking_bugs():
            if report.workstream_id:
                workstreams.add(report.workstream_id)
        return list(workstreams)

    def update_status(self, report_id: str, status: BugStatus) -> None:
        """Update report status."""
        for report in self._reports:
            if report.id == report_id:
                report.status = status
                logger.info(f"Updated bug {report_id} to {status.value}")
                return

    def update_severity(self, report_id: str, severity: BugSeverity) -> None:
        """Update report severity."""
        for report in self._reports:
            if report.id == report_id:
                report.severity = severity
                logger.info(f"Updated bug {report_id} to {severity.value}")
                return

    def route_to_issue_skill(self, report: BugReport) -> None:
        """Route bug report to /hotfix (P0), /bugfix (P1/P2), or /issue (P3)."""
        # Determine skill based on severity
        if report.severity == BugSeverity.P0:
            skill_name = "hotfix"
        elif report.severity in [BugSeverity.P1, BugSeverity.P2]:
            skill_name = "bugfix"
        else:
            skill_name = "issue"

        # In production, would invoke Skill tool here
        logger.info(f"Routing bug {report.id} to /{skill_name}")

        # Placeholder for skill invocation
        # from ... import Skill
        # Skill(skill_name, args=report.title)
