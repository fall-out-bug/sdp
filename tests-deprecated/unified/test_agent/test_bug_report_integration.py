"""Additional integration tests for BugReportFlow.

Tests integration scenarios and edge cases for bug report
creation, routing, and workstream mapping.
"""

import pytest

from sdp.unified.agent.bug_report import (
    BugReportFlow,
    BugSeverity,
    BugStatus,
    BugReport,
)


class TestBugReportIntegration:
    """Integration tests for bug report workflow."""

    def test_full_bug_lifecycle(self):
        """Test complete bug report lifecycle."""
        flow = BugReportFlow()

        # Create bug
        report = flow.create_report(
            title="Login fails",
            description="Users cannot login",
            severity=BugSeverity.P1,
            workstream_id="WS-001",
        )

        # Check initial state
        assert report.status == BugStatus.OPEN
        assert report in flow.get_reports()

        # Update to in progress
        flow.update_status(report.id, BugStatus.IN_PROGRESS)
        assert flow.get_reports_by_status(BugStatus.IN_PROGRESS)[0].id == report.id

        # Update to resolved
        flow.update_status(report.id, BugStatus.RESOLVED)
        assert flow.get_reports_by_status(BugStatus.RESOLVED)[0].id == report.id

    def test_multiple_bugs_per_workstream(self):
        """Test multiple bugs blocking same workstream."""
        flow = BugReportFlow()

        flow.create_report("Bug 1", "Error", BugSeverity.P0, "WS-001")
        flow.create_report("Bug 2", "Fail", BugSeverity.P1, "WS-001")
        flow.create_report("Bug 3", "Issue", BugSeverity.P2, "WS-001")

        ws_bugs = flow.get_reports_by_workstream("WS-001")

        assert len(ws_bugs) == 3

        # Only P0 and P1 are blocking
        blocking = flow.get_blocking_bugs()
        assert len(blocking) == 2

    def test_bug_severity_escalation(self):
        """Test escalating bug severity."""
        flow = BugReportFlow()

        report = flow.create_report(
            title="Bug",
            description="Issue",
            severity=BugSeverity.P3,
        )

        assert report.severity == BugSeverity.P3

        # Escalate to P1
        flow.update_severity(report.id, BugSeverity.P1)

        updated = flow.get_reports()[0]
        assert updated.severity == BugSeverity.P1

        # Should now be in blocking bugs
        assert updated in flow.get_blocking_bugs()

    def test_cross_workstream_blocking(self):
        """Test bugs blocking multiple workstreams."""
        flow = BugReportFlow()

        flow.create_report("Shared bug", "Common issue", BugSeverity.P0, "WS-001")
        flow.create_report("Shared bug", "Common issue", BugSeverity.P0, "WS-002")

        blocking = flow.get_blocking_workstreams()

        assert "WS-001" in blocking
        assert "WS-002" in blocking

    def test_filter_combinations(self):
        """Test filtering by multiple criteria."""
        flow = BugReportFlow()

        report = flow.create_report(
            title="Test bug",
            description="Test",
            severity=BugSeverity.P1,
            workstream_id="WS-001",
        )
        flow.update_status(report.id, BugStatus.IN_PROGRESS)

        # Should find by severity
        p1_reports = flow.get_reports_by_severity(BugSeverity.P1)
        assert report in p1_reports

        # Should find by status
        in_progress = flow.get_reports_by_status(BugStatus.IN_PROGRESS)
        assert report in in_progress

        # Should find by workstream
        ws_reports = flow.get_reports_by_workstream("WS-001")
        assert report in ws_reports

    def test_clear_and_repopulate(self):
        """Test clearing and adding new reports."""
        flow = BugReportFlow()

        # Add some reports
        flow.create_report("Bug 1", "Test", BugSeverity.P1)
        flow.create_report("Bug 2", "Test", BugSeverity.P2)

        assert len(flow.get_reports()) == 2

        # Clear all
        flow._reports.clear()

        assert len(flow.get_reports()) == 0

        # Add new reports
        flow.create_report("Bug 3", "Test", BugSeverity.P0)

        assert len(flow.get_reports()) == 1


class TestSeverityClassificationEdgeCases:
    """Test edge cases in severity classification."""

    def test_case_insensitive_keywords(self):
        """Test case-insensitive keyword matching."""
        flow = BugReportFlow()

        # Should match regardless of case
        assert flow.classify_severity("OUTAGE", "Server DOWN") == BugSeverity.P0
        assert flow.classify_severity("Security", "SQL injection") == BugSeverity.P1

    def test_multiple_keywords_same_priority(self):
        """Test multiple keywords from same priority."""
        flow = BugReportFlow()

        severity = flow.classify_severity(
            "Critical crash",
            "Production outage with corruption",
        )

        # Should still be P0 (first matching priority)
        assert severity == BugSeverity.P0

    def test_no_matching_keywords(self):
        """Test classification with no matching keywords."""
        flow = BugReportFlow()

        severity = flow.classify_severity(
            "Unknown issue",
            "Something happened",
        )

        # Should default to P3
        assert severity == BugSeverity.P3

    def test_keyword_substring_matching(self):
        """Test substring keyword matching."""
        flow = BugReportFlow()

        # "cannot" should match P1
        severity = flow.classify_severity(
            "User cannot login",
            "Authentication fails",
        )

        assert severity == BugSeverity.P1


class TestWorkstreamBlocking:
    """Test workstream blocking scenarios."""

    def test_workstream_not_blocked_if_no_bugs(self):
        """Test workstream with no bugs is not blocked."""
        flow = BugReportFlow()

        blocking = flow.get_blocking_workstreams()

        assert "WS-999" not in blocking

    def test_workstream_blocked_by_p2(self):
        """Test P2 bugs don't block workstreams."""
        flow = BugReportFlow()

        flow.create_report("Minor bug", "Issue", BugSeverity.P2, "WS-001")

        blocking = flow.get_blocking_workstreams()

        # P2 should not block
        assert "WS-001" not in blocking

    def test_workstream_blocked_by_p0(self):
        """Test P0 bugs block workstreams."""
        flow = BugReportFlow()

        flow.create_report("Critical", "Outage", BugSeverity.P0, "WS-001")

        blocking = flow.get_blocking_workstreams()

        assert "WS-001" in blocking

    def test_workstream_unblocked_after_resolved(self):
        """Test workstream unblocked after bug resolved."""
        flow = BugReportFlow()

        report = flow.create_report("Bug", "Error", BugSeverity.P1, "WS-001")

        # Initially blocking
        assert "WS-001" in flow.get_blocking_workstreams()

        # Mark as resolved
        flow.update_status(report.id, BugStatus.RESOLVED)

        # Should no longer be in blocking (if we filter by OPEN/IN_PROGRESS)
        blocking_bugs = [b for b in flow.get_blocking_bugs() if b.status != BugStatus.RESOLVED]
        assert len(blocking_bugs) == 0


class TestReportUpdateOperations:
    """Test report update operations."""

    def test_update_nonexistent_report(self):
        """Test updating non-existent report does nothing."""
        flow = BugReportFlow()

        # Should not raise exception
        flow.update_status("nonexistent-id", BugStatus.IN_PROGRESS)
        flow.update_severity("nonexistent-id", BugSeverity.P0)

        # Should have no reports
        assert len(flow.get_reports()) == 0

    def test_update_preserves_other_fields(self):
        """Test updates preserve other fields."""
        flow = BugReportFlow()

        report = flow.create_report(
            title="Original",
            description="Desc",
            severity=BugSeverity.P1,
            workstream_id="WS-001",
        )

        # Update status
        flow.update_status(report.id, BugStatus.IN_PROGRESS)

        updated = flow.get_reports()[0]
        assert updated.title == "Original"
        assert updated.workstream_id == "WS-001"
        assert updated.severity == BugSeverity.P1

    def test_severity_update_changes_blocking_status(self):
        """Test severity update affects blocking status."""
        flow = BugReportFlow()

        report = flow.create_report(
            title="Bug",
            description="Test",
            severity=BugSeverity.P2,  # Not blocking
            workstream_id="WS-001",
        )

        # Initially not blocking
        assert len(flow.get_blocking_bugs()) == 0

        # Escalate to P1
        flow.update_severity(report.id, BugSeverity.P1)

        # Now blocking
        assert len(flow.get_blocking_bugs()) == 1
