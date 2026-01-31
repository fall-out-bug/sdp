"""Unit tests for traceability service."""

import pytest

from sdp.beads.models import BeadsTask, BeadsStatus
from sdp.beads.mock import MockBeadsClient
from sdp.traceability.models import MappingStatus
from sdp.traceability.service import TraceabilityService


@pytest.fixture
def mock_client():
    """Create mock Beads client with test data."""
    client = MockBeadsClient()
    return client


@pytest.fixture
def ws_task_with_acs(mock_client):
    """Create a workstream task with ACs."""
    from sdp.beads.models import BeadsTaskCreate, BeadsPriority

    task = mock_client.create_task(
        BeadsTaskCreate(
            title="Test WS",
            description="""
## Goal

Test workstream

**Acceptance Criteria:**
- [ ] AC1: User can login
- [ ] AC2: User can logout
- [ ] AC3: Session is maintained
            """,
            priority=BeadsPriority.MEDIUM,
            external_ref="00-032-01",
        )
    )
    return task


class TestTraceabilityService:
    """Tests for TraceabilityService."""

    def test_check_traceability_extracts_acs(self, mock_client, ws_task_with_acs):
        """AC1: Check extracts ACs from description."""
        service = TraceabilityService(mock_client)

        report = service.check_traceability("00-032-01")

        assert report.ws_id == "00-032-01"
        assert report.total_acs == 3
        assert report.mapped_acs == 0
        assert report.missing_acs == 3

    def test_check_traceability_ws_not_found(self, mock_client):
        """AC1: Check raises ValueError for missing WS."""
        service = TraceabilityService(mock_client)

        with pytest.raises(ValueError, match="WS not found"):
            service.check_traceability("99-999-99")

    def test_add_mapping_creates_new(self, mock_client, ws_task_with_acs):
        """AC2: Add mapping creates new entry."""
        service = TraceabilityService(mock_client)

        service.add_mapping(
            "00-032-01", "AC1", "tests/test_auth.py", "test_user_login"
        )

        # Check mapping was stored
        report = service.check_traceability("00-032-01")
        ac1 = next(m for m in report.mappings if m.ac_id == "AC1")

        assert ac1.test_file == "tests/test_auth.py"
        assert ac1.test_name == "test_user_login"
        assert ac1.status == MappingStatus.MAPPED

    def test_add_mapping_updates_existing(self, mock_client, ws_task_with_acs):
        """AC2: Add mapping updates existing entry."""
        service = TraceabilityService(mock_client)

        # Add first mapping
        service.add_mapping("00-032-01", "AC1", "tests/old.py", "test_old")

        # Update with new mapping
        service.add_mapping("00-032-01", "AC1", "tests/new.py", "test_new")

        # Check only one mapping exists
        report = service.check_traceability("00-032-01")
        ac1_mappings = [m for m in report.mappings if m.ac_id == "AC1"]

        assert len(ac1_mappings) == 1
        assert ac1_mappings[0].test_file == "tests/new.py"
        assert ac1_mappings[0].test_name == "test_new"

    def test_extract_acs_checkbox_format(self, mock_client):
        """Extract ACs from checkbox format."""
        service = TraceabilityService(mock_client)

        description = """
        - [ ] AC1: User can login
        - [x] AC2: User can logout
        """

        acs = service._extract_acs(description)

        assert len(acs) == 2
        assert ("AC1", "User can login") in acs
        assert ("AC2", "User can logout") in acs

    def test_extract_acs_plain_format(self, mock_client):
        """Extract ACs from plain format."""
        service = TraceabilityService(mock_client)

        description = """
        AC1: First criterion
        AC2: Second criterion
        """

        acs = service._extract_acs(description)

        assert len(acs) == 2
        assert ("AC1", "First criterion") in acs
        assert ("AC2", "Second criterion") in acs

    def test_extract_acs_case_insensitive(self, mock_client):
        """Extract ACs case-insensitively."""
        service = TraceabilityService(mock_client)

        description = "ac1: lowercase"

        acs = service._extract_acs(description)

        assert len(acs) == 1
        assert ("AC1", "lowercase") in acs
