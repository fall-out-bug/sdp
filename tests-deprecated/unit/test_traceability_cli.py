"""Unit tests for traceability CLI."""

import json
import tempfile
from pathlib import Path

import pytest
from click.testing import CliRunner

from sdp.beads.models import BeadsTaskCreate, BeadsPriority
from sdp.beads.mock import MockBeadsClient
from sdp.cli.trace import trace


@pytest.fixture
def mock_client(monkeypatch):
    """Mock the Beads client."""
    client = MockBeadsClient()

    def mock_create_beads_client():
        return client

    monkeypatch.setattr("sdp.cli.trace.create_beads_client", mock_create_beads_client)
    return client


@pytest.fixture
def ws_with_acs(mock_client):
    """Create a workstream with ACs."""
    task = mock_client.create_task(
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
    return task


class TestTraceCLI:
    """Tests for trace CLI commands."""

    def test_check_shows_table(self, mock_client, ws_with_acs):
        """AC1: check shows mapping table."""
        runner = CliRunner()
        result = runner.invoke(trace, ["check", "00-032-01"])

        assert "Traceability Report" in result.output
        assert "AC" in result.output
        assert "Status" in result.output

    def test_check_exits_1_if_incomplete(self, mock_client, ws_with_acs):
        """AC3: exit code 1 if unmapped ACs."""
        runner = CliRunner()
        result = runner.invoke(trace, ["check", "00-032-01"])

        assert result.exit_code == 1
        assert "INCOMPLETE" in result.output

    def test_check_exits_0_if_complete(self, mock_client, ws_with_acs):
        """AC3: exit code 0 if all ACs mapped."""
        # Add mappings for all ACs
        from sdp.traceability.service import TraceabilityService

        service = TraceabilityService(mock_client)
        service.add_mapping("00-032-01", "AC1", "test.py", "test_1")
        service.add_mapping("00-032-01", "AC2", "test.py", "test_2")

        runner = CliRunner()
        result = runner.invoke(trace, ["check", "00-032-01"])

        assert result.exit_code == 0
        assert "COMPLETE" in result.output

    def test_json_output(self, mock_client, ws_with_acs):
        """AC4: --json flag returns JSON."""
        runner = CliRunner()
        result = runner.invoke(trace, ["check", "00-032-01", "--json"])

        data = json.loads(result.output)
        assert "ws_id" in data
        assert "mappings" in data
        assert data["ws_id"] == "00-032-01"

    def test_check_ws_not_found(self, mock_client):
        """AC1: error if WS not found."""
        runner = CliRunner()
        result = runner.invoke(trace, ["check", "99-999-99"])

        assert result.exit_code == 1
        assert "❌" in result.output
        assert "not found" in result.output.lower()

    def test_add_mapping(self, mock_client, ws_with_acs):
        """AC2: add command creates mapping."""
        runner = CliRunner()
        result = runner.invoke(
            trace,
            [
                "add",
                "00-032-01",
                "--ac",
                "AC1",
                "--test",
                "test_login",
                "--file",
                "tests/test_auth.py",
            ],
        )

        assert result.exit_code == 0
        assert "✅" in result.output
        assert "AC1" in result.output
        assert "test_login" in result.output

    def test_add_mapping_without_file(self, mock_client, ws_with_acs):
        """AC2: add command works without file."""
        runner = CliRunner()
        result = runner.invoke(
            trace, ["add", "00-032-01", "--ac", "AC1", "--test", "test_login"]
        )

        assert result.exit_code == 0

    def test_auto_not_implemented(self, mock_client):
        """AC5: auto command detects mappings."""
        from sdp.beads.models import BeadsTaskCreate, BeadsPriority
        import tempfile

        # Create WS with ACs
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

        # Create temporary test directory
        with tempfile.TemporaryDirectory() as tmpdir:
            test_file = Path(tmpdir) / "test_auth.py"
            test_file.write_text(
                '''
def test_ac1_login():
    """Tests AC1: User can login."""
    pass
'''
            )

            runner = CliRunner()
            result = runner.invoke(trace, ["auto", "00-032-01", "--test-dir", tmpdir])

            assert result.exit_code == 0
            assert "Detected" in result.output
            assert "AC1" in result.output
            assert "test_ac1_login" in result.output
