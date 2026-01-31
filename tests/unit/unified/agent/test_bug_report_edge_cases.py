"""Tests for bug report flow edge cases."""

import pytest
from sdp.unified.agent.bug_report import (
    BugReportFlow,
    BugReport,
    BugSeverity,
    BugStatus,
)


class TestBugSeverity:
    """Tests for BugSeverity enum."""

    def test_p0_label(self):
        """Test P0 returns hotfix label."""
        assert BugSeverity.P0.label == "hotfix"

    def test_p1_label(self):
        """Test P1 returns bugfix label."""
        assert BugSeverity.P1.label == "bugfix"

    def test_p2_label(self):
        """Test P2 returns bugfix label."""
        assert BugSeverity.P2.label == "bugfix"

    def test_p3_label(self):
        """Test P3 returns backlog label."""
        assert BugSeverity.P3.label == "backlog"


class TestBugReport:
    """Tests for BugReport dataclass."""

    def test_default_initialization(self):
        """Test default bug report initialization."""
        report = BugReport()

        assert report.id != ""
        assert report.title == ""
        assert report.description == ""
        assert report.severity == BugSeverity.P3
        assert report.status == BugStatus.OPEN
        assert report.workstream_id is None

    def test_custom_initialization(self):
        """Test custom bug report initialization."""
        report = BugReport(
            title="Test Bug",
            description="Test description",
            severity=BugSeverity.P1,
            status=BugStatus.IN_PROGRESS,
            workstream_id="00-001-01",
        )

        assert report.title == "Test Bug"
        assert report.description == "Test description"
        assert report.severity == BugSeverity.P1
        assert report.status == BugStatus.IN_PROGRESS
        assert report.workstream_id == "00-001-01"


class TestBugReportFlow:
    """Tests for BugReportFlow."""

    def test_classify_severity_p0_outage(self):
        """Test classifies outage as P0."""
        flow = BugReportFlow()

        report = flow.create_report(
            title="Production outage",
            description="System is down",
        )

        assert report.severity == BugSeverity.P0

    def test_classify_severity_p0_critical(self):
        """Test classifies critical as P0."""
        flow = BugReportFlow()

        report = flow.create_report(
            title="Critical bug",
            description="Crash on startup",
        )

        assert report.severity == BugSeverity.P0

    def test_classify_severity_p1_security(self):
        """Test classifies security as P1."""
        flow = BugReportFlow()

        report = flow.create_report(
            title="Security vulnerability",
            description="SQL injection possible",
        )

        assert report.severity == BugSeverity.P1

    def test_classify_severity_p1_data_loss(self):
        """Test classifies data loss as P1."""
        flow = BugReportFlow()

        report = flow.create_report(
            title="Data loss bug",
            description="User data not saved",
        )

        assert report.severity == BugSeverity.P1

    def test_classify_severity_p2_ui(self):
        """Test classifies UI issue as P2."""
        flow = BugReportFlow()

        report = flow.create_report(
            title="UI glitch",
            description="Button is misaligned",
        )

        assert report.severity == BugSeverity.P2

    def test_classify_severity_p2_performance(self):
        """Test classifies performance issue as P2."""
        flow = BugReportFlow()

        report = flow.create_report(
            title="Slow loading",
            description="Page takes too long",
        )

        assert report.severity == BugSeverity.P2

    def test_classify_severity_p3_typo(self):
        """Test classifies typo as P3."""
        flow = BugReportFlow()

        report = flow.create_report(
            title="Typo in readme",
            description="Cosmetic typo in documentation",
        )

        assert report.severity == BugSeverity.P3

    def test_classify_severity_p3_cosmetic(self):
        """Test classifies cosmetic issue as P3."""
        flow = BugReportFlow()

        report = flow.create_report(
            title="Cosmetic issue",
            description="Color slightly off",
        )

        assert report.severity == BugSeverity.P3

    def test_classify_severity_default_p3(self):
        """Test defaults to P3 when no keywords match."""
        flow = BugReportFlow()

        report = flow.create_report(
            title="General issue",
            description="Something happened",
        )

        assert report.severity == BugSeverity.P3

    def test_create_report_stores_in_list(self):
        """Test created reports are stored."""
        flow = BugReportFlow()

        report1 = flow.create_report(title="Bug 1", description="Desc 1")
        report2 = flow.create_report(title="Bug 2", description="Desc 2")

        assert len(flow._reports) == 2
        assert flow._reports[0].id == report1.id
        assert flow._reports[1].id == report2.id

    def test_get_reports_returns_all(self):
        """Test retrieves all reports."""
        flow = BugReportFlow()
        report1 = flow.create_report(title="Bug 1", description="Desc 1")
        report2 = flow.create_report(title="Bug 2", description="Desc 2")

        reports = flow.get_reports()

        assert len(reports) == 2
        assert reports[0].id == report1.id
        assert reports[1].id == report2.id

    def test_get_reports_empty(self):
        """Test returns empty list when no reports."""
        flow = BugReportFlow()

        reports = flow.get_reports()

        assert reports == []

    def test_case_insensitive_classification(self):
        """Test keyword matching is case insensitive."""
        flow = BugReportFlow()

        report = flow.create_report(
            title="PRODUCTION OUTAGE",
            description="CRITICAL system DOWN",
        )

        assert report.severity == BugSeverity.P0
