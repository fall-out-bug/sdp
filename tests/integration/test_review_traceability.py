"""Integration tests for review with traceability."""

import pytest
from click.testing import CliRunner

from sdp.beads.models import BeadsTaskCreate, BeadsPriority
from sdp.beads.mock import MockBeadsClient
from sdp.cli.trace import trace
from sdp.traceability.service import TraceabilityService


@pytest.fixture
def mock_client(monkeypatch):
    """Mock the Beads client."""
    client = MockBeadsClient()

    def mock_create_beads_client():
        return client

    monkeypatch.setattr("sdp.cli.trace.create_beads_client", mock_create_beads_client)
    return client


class TestReviewWithTraceability:
    """Tests for review integration with traceability."""

    def test_review_fails_on_missing_traceability(self, mock_client):
        """Review should fail if traceability incomplete."""
        # Create WS with unmapped ACs
        mock_client.create_task(
            BeadsTaskCreate(
                title="Test WS",
                description="""
**Acceptance Criteria:**
- [ ] AC1: User can login
- [ ] AC2: User can logout
                """,
                priority=BeadsPriority.MEDIUM,
                external_ref="00-032-01",
            )
        )

        # Run traceability check (part of review flow)
        runner = CliRunner()
        result = runner.invoke(trace, ["check", "00-032-01"])

        assert result.exit_code == 1
        assert "INCOMPLETE" in result.output

    def test_review_passes_with_complete_traceability(self, mock_client):
        """Review should pass if traceability complete."""
        # Create WS with ACs
        mock_client.create_task(
            BeadsTaskCreate(
                title="Test WS",
                description="""
**Acceptance Criteria:**
- [ ] AC1: User can login
- [ ] AC2: User can logout
                """,
                priority=BeadsPriority.MEDIUM,
                external_ref="00-032-01",
            )
        )

        # Add mappings
        service = TraceabilityService(mock_client)
        service.add_mapping("00-032-01", "AC1", "tests/test_auth.py", "test_login")
        service.add_mapping("00-032-01", "AC2", "tests/test_auth.py", "test_logout")

        # Run traceability check
        runner = CliRunner()
        result = runner.invoke(trace, ["check", "00-032-01"])

        assert result.exit_code == 0
        assert "COMPLETE" in result.output
        assert "100%" in result.output

    def test_review_shows_traceability_table(self, mock_client):
        """Review should display traceability table."""
        # Create WS
        mock_client.create_task(
            BeadsTaskCreate(
                title="Test WS",
                description="""
**Acceptance Criteria:**
- [ ] AC1: User can login
                """,
                priority=BeadsPriority.MEDIUM,
                external_ref="00-032-01",
            )
        )

        # Add mapping
        service = TraceabilityService(mock_client)
        service.add_mapping("00-032-01", "AC1", "tests/test.py", "test_login")

        # Check output format
        runner = CliRunner()
        result = runner.invoke(trace, ["check", "00-032-01"])

        assert "Traceability Report" in result.output
        assert "AC1" in result.output
        assert "test_login" in result.output
        assert "✅" in result.output

    def test_review_with_auto_detection(self, mock_client):
        """Review can use auto-detection for missing mappings."""
        import tempfile
        from pathlib import Path

        # Create WS
        mock_client.create_task(
            BeadsTaskCreate(
                title="Test WS",
                description="""
**Acceptance Criteria:**
- [ ] AC1: User can login
                """,
                priority=BeadsPriority.MEDIUM,
                external_ref="00-032-01",
            )
        )

        # Create test file with AC reference
        with tempfile.TemporaryDirectory() as tmpdir:
            test_file = Path(tmpdir) / "test_auth.py"
            test_file.write_text(
                '''
def test_ac1_login():
    """Tests AC1: User can login."""
    pass
'''
            )

            # Run auto-detection
            runner = CliRunner()
            result = runner.invoke(
                trace, ["auto", "00-032-01", "--test-dir", tmpdir, "--apply"]
            )

            assert result.exit_code == 0
            assert "Detected" in result.output
            assert "Applied" in result.output

            # Verify traceability is now complete
            result = runner.invoke(trace, ["check", "00-032-01"])
            assert result.exit_code == 0
            assert "COMPLETE" in result.output
