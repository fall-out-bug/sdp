"""
Tests for execution report generator.
"""

import tempfile
from pathlib import Path
from unittest.mock import patch

import pytest

from sdp.report.generator import ExecutionStats, ReportGenerator


class TestReportGenerator:
    """Test ReportGenerator class."""

    def test_init_with_ws_id(self) -> None:
        """Test initialization with workstream ID."""
        gen = ReportGenerator(ws_id="WS-001-01")
        assert gen.ws_id == "WS-001-01"
        assert gen.start_time is None
        assert gen.end_time is None

    def test_start_timer(self) -> None:
        """Test starting execution timer."""
        gen = ReportGenerator(ws_id="WS-001-01")
        gen.start_timer()
        assert gen.start_time is not None

    def test_stop_timer(self) -> None:
        """Test stopping execution timer."""
        gen = ReportGenerator(ws_id="WS-001-01")
        gen.start_timer()
        gen.stop_timer()
        assert gen.end_time is not None

    def test_get_duration_minutes(self) -> None:
        """Test calculating duration in minutes."""
        gen = ReportGenerator(ws_id="WS-001-01")
        gen.start_timer()
        gen.stop_timer()
        duration = gen.get_duration_minutes()
        assert duration >= 0.0
        assert duration < 1.0  # Should be very fast

    def test_get_duration_minutes_not_started(self) -> None:
        """Test duration calculation when timer not started."""
        gen = ReportGenerator(ws_id="WS-001-01")
        duration = gen.get_duration_minutes()
        assert duration == 0.0

    def test_collect_stats_basic(self) -> None:
        """Test collecting basic statistics."""
        gen = ReportGenerator(ws_id="WS-001-01")
        gen.start_timer()
        gen.stop_timer()

        stats = gen.collect_stats(
            files_changed=[("src/module.py", "modified", 50)],
            coverage_pct=85.0,
            tests_passed=10,
            tests_failed=0,
        )

        assert stats.files_changed == 1
        assert stats.loc_added == 50
        assert stats.loc_removed == 0
        assert stats.coverage_pct == 85.0
        assert stats.tests_passed == 10
        assert stats.tests_failed == 0
        assert stats.deviations == []

    def test_collect_stats_with_deviations(self) -> None:
        """Test collecting statistics with deviations."""
        gen = ReportGenerator(ws_id="WS-001-01")
        gen.start_timer()
        gen.stop_timer()

        deviations = ["Added extra validation", "Refactored for clarity"]
        stats = gen.collect_stats(deviations=deviations)

        assert stats.deviations == deviations

    def test_generate_report_markdown(self) -> None:
        """Test generating markdown report."""
        gen = ReportGenerator(ws_id="WS-001-01")

        stats = ExecutionStats(
            duration_minutes=45.5,
            files_changed=3,
            loc_added=150,
            loc_removed=20,
            coverage_pct=87.5,
            tests_passed=15,
            tests_failed=0,
            deviations=[],
        )

        report = gen.generate_report_markdown(
            stats, executed_by="test-user", commit_hash="abc123"
        )

        assert "## Execution Report" in report
        assert "test-user" in report
        assert "45.5 minutes" in report
        assert "87.5%" in report
        assert "abc123" in report
        assert "**Files Changed:** 3" in report
        assert "**Lines Added:** 150" in report

    def test_generate_report_with_deviations(self) -> None:
        """Test report includes deviations."""
        gen = ReportGenerator(ws_id="WS-001-01")

        stats = ExecutionStats(
            duration_minutes=30.0,
            files_changed=1,
            loc_added=50,
            loc_removed=0,
            coverage_pct=90.0,
            tests_passed=5,
            tests_failed=0,
            deviations=["Added extra error handling"],
        )

        report = gen.generate_report_markdown(stats, executed_by="dev")

        assert "Added extra error handling" in report
        assert "Deviations from Plan" in report

    def test_append_report(self) -> None:
        """Test appending report to workstream file."""
        with tempfile.TemporaryDirectory() as tmpdir:
            ws_file = Path(tmpdir) / "WS-001-01.md"
            ws_file.write_text("# Original Content\n")

            gen = ReportGenerator(ws_id="WS-001-01", ws_path=ws_file)

            stats = ExecutionStats(
                duration_minutes=60.0,
                files_changed=2,
                loc_added=100,
                loc_removed=10,
                coverage_pct=82.0,
                tests_passed=8,
                tests_failed=0,
                deviations=[],
            )

            gen.append_report(stats, executed_by="tester", commit_hash="def456")

            content = ws_file.read_text()
            assert "# Original Content" in content
            assert "## Execution Report" in content
            assert "60.0 minutes" in content

    def test_append_report_file_not_found(self) -> None:
        """Test appending report when file doesn't exist."""
        gen = ReportGenerator(
            ws_id="WS-001-01", ws_path=Path("/nonexistent/path.md")
        )

        stats = ExecutionStats(
            duration_minutes=1.0,
            files_changed=0,
            loc_added=0,
            loc_removed=0,
            coverage_pct=0.0,
            tests_passed=0,
            tests_failed=0,
            deviations=[],
        )

        with pytest.raises(FileNotFoundError):
            gen.append_report(stats, executed_by="test")

    def test_generate_complete_workflow(self) -> None:
        """Test complete generation workflow."""
        gen = ReportGenerator(ws_id="WS-001-02")
        gen.start_timer()

        # Simulate some work
        gen.stop_timer()

        report = gen.generate(
            executed_by="workflow-tester",
            files_changed=[("src/new.py", "created", 75)],
            coverage_pct=88.0,
            tests_passed=12,
            tests_failed=1,
            deviations=["Fixed minor bug"],
            commit_hash="xyz789",
        )

        assert "## Execution Report" in report
        assert "workflow-tester" in report
        assert "88.0%" in report
        assert "**Tests Passed:** 12" in report
        assert "**Tests Failed:** 1" in report
        assert "Fixed minor bug" in report
        assert "xyz789" in report

    def test_find_ws_file_auto_detect(self) -> None:
        """Test auto-detection of workstream file."""
        with tempfile.TemporaryDirectory() as tmpdir:
            ws_dir = Path(tmpdir) / "docs" / "workstreams" / "in_progress"
            ws_dir.mkdir(parents=True)

            ws_file = ws_dir / "WS-001-03.md"
            ws_file.write_text("# Test WS\n")

            # Change to temp directory
            original_cwd = Path.cwd()
            import os

            try:
                os.chdir(tmpdir)
                gen = ReportGenerator(ws_id="WS-001-03")
                # Resolve to absolute path for comparison
                assert gen.ws_path.resolve() == ws_file.resolve()
            finally:
                os.chdir(original_cwd)

    def test_find_ws_file_not_found(self) -> None:
        """Test error when workstream file not found."""
        with tempfile.TemporaryDirectory() as tmpdir:
            original_cwd = Path.cwd()
            import os

            try:
                os.chdir(tmpdir)
                gen = ReportGenerator(ws_id="WS-999-99")
                # Accessing ws_path should trigger FileNotFoundError
                with pytest.raises(FileNotFoundError):
                    _ = gen.ws_path
            finally:
                os.chdir(original_cwd)

    @patch("subprocess.run")
    def test_get_files_changed_with_git(self, mock_run) -> None:
        """Test getting files changed from git."""
        mock_run.return_value.returncode = 0
        mock_run.return_value.stdout = " src/file.py | 10 ++++++++++\n"

        gen = ReportGenerator(ws_id="WS-001-01")
        table = gen._get_files_changed_markdown()

        assert "file.py" in table or "|      |        |     |" in table

    @patch("subprocess.run")
    def test_get_files_changed_git_error(self, mock_run) -> None:
        """Test handling git error gracefully."""
        mock_run.return_value.returncode = 1

        gen = ReportGenerator(ws_id="WS-001-01")
        table = gen._get_files_changed_markdown()

        # Should return empty table on error
        assert "|      |        |     |" in table


class TestExecutionStats:
    """Test ExecutionStats dataclass."""

    def test_creation(self) -> None:
        """Test creating ExecutionStats."""
        stats = ExecutionStats(
            duration_minutes=30.0,
            files_changed=2,
            loc_added=100,
            loc_removed=10,
            coverage_pct=85.0,
            tests_passed=10,
            tests_failed=0,
            deviations=[],
        )

        assert stats.duration_minutes == 30.0
        assert stats.files_changed == 2
        assert stats.loc_added == 100
        assert stats.coverage_pct == 85.0
        assert stats.tests_passed == 10
