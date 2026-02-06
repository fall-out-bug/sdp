"""Tests for BugReportFlow module.

Tests bug report creation, severity classification, and integration
with /issue skill for bug routing and tracking.
"""

import pytest
from dataclasses import dataclass
from enum import Enum
from unittest.mock import Mock, patch

from sdp.unified.agent.bug_report import (
    BugReportFlow,
    BugSeverity,
    BugStatus,
    BugReport,
)


class TestBugSeverityEnum:
    """Test BugSeverity enum."""

    def test_has_p0_severity(self):
        """Should have P0 severity for hotfix."""
        assert BugSeverity.P0 is not None
        assert BugSeverity.P0.value == "P0"
        assert BugSeverity.P0.label == "hotfix"

    def test_has_p1_severity(self):
        """Should have P1 severity for bugfix."""
        assert BugSeverity.P1 is not None
        assert BugSeverity.P1.value == "P1"

    def test_has_p2_severity(self):
        """Should have P2 severity for bugfix."""
        assert BugSeverity.P2 is not None
        assert BugSeverity.P2.value == "P2"

    def test_has_p3_severity(self):
        """Should have P3 severity for backlog."""
        assert BugSeverity.P3 is not None
        assert BugSeverity.P3.value == "P3"
        assert BugSeverity.P3.label == "backlog"


class TestBugStatusEnum:
    """Test BugStatus enum."""

    def test_has_open_status(self):
        """Should have OPEN status."""
        assert BugStatus.OPEN is not None
        assert BugStatus.OPEN.value == "open"

    def test_has_in_progress_status(self):
        """Should have IN_PROGRESS status."""
        assert BugStatus.IN_PROGRESS is not None

    def test_has_resolved_status(self):
        """Should have RESOLVED status."""
        assert BugStatus.RESOLVED is not None


class TestBugReportDataclass:
    """Test BugReport dataclass."""

    def test_create_report_with_required_fields(self):
        """Should create report with required fields."""
        report = BugReport(
            title="Test failure",
            description="Test failed unexpectedly",
            severity=BugSeverity.P1,
        )

        assert report.title == "Test failure"
        assert report.severity == BugSeverity.P1

    def test_create_report_with_optional_fields(self):
        """Should create report with optional fields."""
        report = BugReport(
            title="Bug",
            description="Bug description",
            severity=BugSeverity.P2,
            workstream_id="WS-001",
            status=BugStatus.IN_PROGRESS,
        )

        assert report.workstream_id == "WS-001"
        assert report.status == BugStatus.IN_PROGRESS

    def test_report_defaults(self):
        """Should have correct default values."""
        report = BugReport(
            title="Bug",
            description="Desc",
            severity=BugSeverity.P3,
        )

        assert report.workstream_id is None
        assert report.status == BugStatus.OPEN


class TestBugReportFlowInit:
    """Test BugReportFlow initialization."""

    def test_creates_flow(self):
        """Should initialize flow."""
        flow = BugReportFlow()

        assert flow is not None
        assert hasattr(flow, 'create_report')
        assert hasattr(flow, 'classify_severity')

    def test_initializes_empty_reports(self):
        """Should initialize with empty report tracking."""
        flow = BugReportFlow()

        assert flow.get_reports() == []


class TestSeverityClassification:
    """Test severity classification."""

    def test_classifies_p0_for_critical(self):
        """Should classify as P0 for critical bugs."""
        flow = BugReportFlow()

        severity = flow.classify_severity(
            "Production outage",
            "API is down for all users",
        )

        assert severity == BugSeverity.P0

    def test_classifies_p1_for_high_priority(self):
        """Should classify as P1 for high priority bugs."""
        flow = BugReportFlow()

        severity = flow.classify_severity(
            "Login fails",
            "Users cannot login to application",
        )

        assert severity == BugSeverity.P1

    def test_classifies_p2_for_medium_priority(self):
        """Should classify as P2 for medium priority bugs."""
        flow = BugReportFlow()

        severity = flow.classify_severity(
            "UI glitch",
            "Button alignment is off",
        )

        assert severity == BugSeverity.P2

    def test_classifies_p3_for_low_priority(self):
        """Should classify as P3 for low priority bugs."""
        flow = BugReportFlow()

        severity = flow.classify_severity(
            "Typo in docs",
            "Minor typo in documentation",
        )

        assert severity == BugSeverity.P3

    def test_classifies_based_on_keywords(self):
        """Should classify based on description keywords."""
        flow = BugReportFlow()

        # Critical keywords
        assert flow.classify_severity("Outage", "Production is down") == BugSeverity.P0
        assert flow.classify_severity("Crash", "Application crashes") == BugSeverity.P0

        # High priority keywords
        assert flow.classify_severity("Security", "SQL injection vulnerability") == BugSeverity.P1
        assert flow.classify_severity("Data loss", "User data deleted") == BugSeverity.P1


class TestBugReportCreation:
    """Test bug report creation."""

    def test_creates_report(self):
        """Should create bug report."""
        flow = BugReportFlow()

        report = flow.create_report(
            title="Test failure",
            description="Test failed",
            severity=BugSeverity.P1,
        )

        assert report.title == "Test failure"
        assert report.severity == BugSeverity.P1
        assert len(flow.get_reports()) == 1

    def test_assigns_workstream_to_report(self):
        """Should assign workstream to report."""
        flow = BugReportFlow()

        report = flow.create_report(
            title="Bug",
            description="Desc",
            severity=BugSeverity.P2,
            workstream_id="WS-123",
        )

        assert report.workstream_id == "WS-123"

    def test_auto_classifies_severity(self):
        """Should auto-classify severity if not provided."""
        flow = BugReportFlow()

        report = flow.create_report(
            title="Production outage",
            description="API is down",
        )

        assert report.severity == BugSeverity.P0


class TestBugReportRetrieval:
    """Test bug report retrieval and querying."""

    def test_gets_all_reports(self):
        """Should get all reports."""
        flow = BugReportFlow()

        flow.create_report("Bug 1", "Desc 1", BugSeverity.P1)
        flow.create_report("Bug 2", "Desc 2", BugSeverity.P2)

        reports = flow.get_reports()

        assert len(reports) == 2

    def test_filters_reports_by_severity(self):
        """Should filter reports by severity."""
        flow = BugReportFlow()

        flow.create_report("P0 bug", "Critical", BugSeverity.P0)
        flow.create_report("P1 bug", "High", BugSeverity.P1)
        flow.create_report("P2 bug", "Medium", BugSeverity.P2)

        p0_reports = flow.get_reports_by_severity(BugSeverity.P0)
        p1_reports = flow.get_reports_by_severity(BugSeverity.P1)

        assert len(p0_reports) == 1
        assert len(p1_reports) == 1
        assert p0_reports[0].title == "P0 bug"

    def test_filters_reports_by_workstream(self):
        """Should filter reports by workstream."""
        flow = BugReportFlow()

        flow.create_report("Bug 1", "Desc", BugSeverity.P1, "WS-001")
        flow.create_report("Bug 2", "Desc", BugSeverity.P2, "WS-002")
        flow.create_report("Bug 3", "Desc", BugSeverity.P1, "WS-001")

        ws_reports = flow.get_reports_by_workstream("WS-001")

        assert len(ws_reports) == 2

    def test_gets_blocking_bugs(self):
        """Should get blocking bugs (P0, P1)."""
        flow = BugReportFlow()

        flow.create_report("P0 bug", "Critical", BugSeverity.P0)
        flow.create_report("P1 bug", "High", BugSeverity.P1)
        flow.create_report("P2 bug", "Medium", BugSeverity.P2)

        blocking = flow.get_blocking_bugs()

        assert len(blocking) == 2

    def test_filters_reports_by_status(self):
        """Should filter reports by status."""
        flow = BugReportFlow()

        report1 = flow.create_report("Bug 1", "Desc", BugSeverity.P1)
        report2 = flow.create_report("Bug 2", "Desc", BugSeverity.P2)

        # Manually update status
        report2.status = BugStatus.IN_PROGRESS

        open_reports = flow.get_reports_by_status(BugStatus.OPEN)

        assert len(open_reports) == 1
        assert open_reports[0].title == "Bug 1"


class TestBugReportUpdate:
    """Test bug report updates."""

    def test_updates_report_status(self):
        """Should update report status."""
        flow = BugReportFlow()

        report = flow.create_report("Bug", "Desc", BugSeverity.P1)

        flow.update_status(report.id, BugStatus.IN_PROGRESS)

        updated = flow.get_reports()[0]
        assert updated.status == BugStatus.IN_PROGRESS

    def test_updates_report_severity(self):
        """Should update report severity."""
        flow = BugReportFlow()

        report = flow.create_report("Bug", "Desc", BugSeverity.P2)

        flow.update_severity(report.id, BugSeverity.P0)

        updated = flow.get_reports()[0]
        assert updated.severity == BugSeverity.P0


class TestBugWorkstreamMapping:
    """Test bug-to-workstream mapping."""

    def test_maps_bug_to_workstream(self):
        """Should map bug to blocking workstream."""
        flow = BugReportFlow()

        report = flow.create_report(
            title="Bug in WS-001",
            description="Failed",
            severity=BugSeverity.P1,
            workstream_id="WS-001",
        )

        blocking = flow.get_blocking_workstreams()

        assert "WS-001" in blocking

    def test_tracks_multiple_blocking_bugs(self):
        """Should track multiple blocking bugs per workstream."""
        flow = BugReportFlow()

        flow.create_report("Bug 1", "High", BugSeverity.P1, "WS-001")
        flow.create_report("Bug 2", "Critical", BugSeverity.P0, "WS-001")
        flow.create_report("Bug 3", "Medium", BugSeverity.P2, "WS-002")

        blocking = flow.get_blocking_workstreams()

        # Only WS-001 has blocking bugs (P0 and P1)
        # WS-002 has P2 which is not blocking
        assert "WS-001" in blocking
        assert "WS-002" not in blocking
