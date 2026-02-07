"""Extended tests for traceability service edge cases and error paths."""

from pathlib import Path
from unittest.mock import Mock, patch
import pytest

from sdp.beads.models import BeadsTask, BeadsStatus, BeadsPriority, BeadsTaskCreate
from sdp.beads.mock import MockBeadsClient
from sdp.traceability.service import TraceabilityService


@pytest.fixture
def mock_client():
    """Create mock Beads client."""
    return MockBeadsClient()


@pytest.fixture
def service(mock_client):
    """Create traceability service."""
    return TraceabilityService(mock_client)


class TestGetWSTask:
    """Tests for _get_ws_task method."""

    def test_get_ws_task_with_exception(self, service):
        """Test _get_ws_task handles list_tasks exception."""
        # Mock client to raise exception
        service._client.list_tasks = Mock(side_effect=ConnectionError("API down"))

        result = service._get_ws_task("00-032-01")

        assert result is None

    def test_get_ws_task_not_found(self, service, mock_client):
        """Test _get_ws_task returns None when task not found."""
        # Create task with different external_ref
        mock_client.create_task(
            BeadsTaskCreate(
                title="Other WS",
                description="Test",
                priority=BeadsPriority.MEDIUM,
                external_ref="00-999-01",
            )
        )

        result = service._get_ws_task("00-032-01")

        assert result is None


class TestMarkdownFallback:
    """Tests for markdown fallback paths."""

    def test_get_ws_file_path_not_exists(self, service, tmp_path, monkeypatch):
        """Test _get_ws_file_path returns None when directories don't exist."""
        monkeypatch.chdir(tmp_path)
        # Don't create docs/workstreams directories

        result = service._get_ws_file_path("00-032-01")

        assert result is None

    def test_get_ws_file_path_no_matching_files(self, service, tmp_path, monkeypatch):
        """Test _get_ws_file_path returns None when no matching files found."""
        monkeypatch.chdir(tmp_path)
        # Create directories but no matching files
        ws_dir = tmp_path / "docs" / "workstreams" / "backlog"
        ws_dir.mkdir(parents=True)

        result = service._get_ws_file_path("00-032-01")

        assert result is None

    def test_get_ws_content_from_markdown_not_found(self, service):
        """Test _get_ws_content_from_markdown returns None when file not found."""
        with patch.object(service, '_get_ws_file_path', return_value=None):
            result = service._get_ws_content_from_markdown("00-032-01")

        assert result is None

    def test_check_traceability_markdown_fallback_not_found(self, service):
        """Test check_traceability raises ValueError when WS not in Beads or markdown."""
        with patch.object(service, '_get_ws_task', return_value=None):
            with patch.object(service, '_get_ws_content_from_markdown', return_value=None):
                with pytest.raises(ValueError, match="WS not found: 00-032-01"):
                    service.check_traceability("00-032-01")


class TestGetTraceabilityFromMarkdown:
    """Tests for _get_traceability_from_markdown method."""

    def test_get_traceability_no_frontmatter(self, service):
        """Test returns empty list when no frontmatter."""
        content = "# WS Title\n\nSome content"

        result = service._get_traceability_from_markdown(content)

        assert result == []

    def test_get_traceability_invalid_yaml(self, service):
        """Test returns empty list when YAML is invalid."""
        content = "---\ninvalid: yaml: content:\n---\n\nBody"

        result = service._get_traceability_from_markdown(content)

        assert result == []

    def test_get_traceability_not_dict(self, service):
        """Test returns empty list when frontmatter is not a dict."""
        content = "---\n- list\n- items\n---\n\nBody"

        result = service._get_traceability_from_markdown(content)

        assert result == []

    def test_get_traceability_valid(self, service):
        """Test extracts traceability from valid frontmatter."""
        content = """---
ws_id: "00-032-01"
traceability:
  - ac_id: "AC1"
    test_file: "tests/test_foo.py"
    test_name: "test_bar"
    status: "mapped"
---

Body content
"""

        result = service._get_traceability_from_markdown(content)

        assert len(result) == 1
        assert result[0]["ac_id"] == "AC1"
        assert result[0]["test_file"] == "tests/test_foo.py"


class TestAddMappingToMarkdown:
    """Tests for _add_mapping_to_markdown method."""

    def test_add_mapping_to_markdown_ws_not_found(self, service):
        """Test raises ValueError when WS file not found."""
        with patch.object(service, '_get_ws_file_path', return_value=None):
            with pytest.raises(ValueError, match="WS not found: 00-032-01"):
                service._add_mapping_to_markdown("00-032-01", "AC1", "test.py", "test_func")

    def test_add_mapping_to_markdown_no_frontmatter(self, service, tmp_path):
        """Test raises ValueError when markdown has no frontmatter."""
        ws_file = tmp_path / "00-032-01-test.md"
        ws_file.write_text("# WS Title\n\nNo frontmatter")

        with patch.object(service, '_get_ws_file_path', return_value=ws_file):
            with patch.object(service, '_extract_acs', return_value=[("AC1", "Test AC")]):
                with pytest.raises(ValueError, match="No frontmatter"):
                    service._add_mapping_to_markdown("00-032-01", "AC1", "test.py", "test_func")

    def test_add_mapping_to_markdown_new_mapping(self, service, tmp_path):
        """Test adds new mapping to markdown frontmatter."""
        ws_file = tmp_path / "00-032-01-test.md"
        ws_file.write_text("""---
ws_id: "00-032-01"
---

# WS Title
- [ ] AC1: Test criterion
""")

        with patch.object(service, '_get_ws_file_path', return_value=ws_file):
            with patch.object(service, '_extract_acs', return_value=[("AC1", "Test criterion")]):
                service._add_mapping_to_markdown("00-032-01", "AC1", "test.py", "test_func")

        content = ws_file.read_text()
        assert "traceability:" in content
        assert "ac_id: AC1" in content
        assert "test_file: test.py" in content
        assert "test_name: test_func" in content

    def test_add_mapping_to_markdown_update_existing(self, service, tmp_path):
        """Test updates existing mapping in markdown frontmatter."""
        ws_file = tmp_path / "00-032-01-test.md"
        ws_file.write_text("""---
ws_id: "00-032-01"
traceability:
  - ac_id: "AC1"
    test_file: "old.py"
    test_name: "old_test"
    status: "mapped"
---

# WS Title
- [ ] AC1: Test criterion
""")

        with patch.object(service, '_get_ws_file_path', return_value=ws_file):
            with patch.object(service, '_extract_acs', return_value=[("AC1", "Test criterion")]):
                service._add_mapping_to_markdown("00-032-01", "AC1", "new.py", "new_test")

        content = ws_file.read_text()
        assert "test_file: new.py" in content
        assert "test_name: new_test" in content
        assert "old.py" not in content


class TestAddMappingRouting:
    """Tests for add_mapping routing logic."""

    def test_add_mapping_routes_to_beads(self, service):
        """Test add_mapping routes to Beads when task exists."""
        mock_task = Mock()
        mock_task.description = "- [ ] AC1: Test"
        mock_task.sdp_metadata = {}

        with patch.object(service, '_get_ws_task', return_value=mock_task):
            with patch.object(service, '_add_mapping_to_beads') as mock_add_beads:
                service.add_mapping("00-032-01", "AC1", "test.py", "test_func")

                mock_add_beads.assert_called_once()

    def test_add_mapping_routes_to_markdown(self, service, tmp_path):
        """Test add_mapping routes to markdown when task not in Beads."""
        ws_file = tmp_path / "00-032-01-test.md"
        ws_file.write_text("""---
ws_id: "00-032-01"
---

# WS Title
- [ ] AC1: Test
""")

        with patch.object(service, '_get_ws_task', return_value=None):
            with patch.object(service, '_get_ws_file_path', return_value=ws_file):
                with patch.object(service, '_extract_acs', return_value=[("AC1", "Test")]):
                    service.add_mapping("00-032-01", "AC1", "test.py", "test_func")

        # Verify file was updated
        content = ws_file.read_text()
        assert "traceability:" in content


class TestCheckTraceabilityMarkdownPath:
    """Tests for check_traceability using markdown path."""

    def test_check_traceability_from_markdown(self, service, tmp_path):
        """Test check_traceability reads from markdown when not in Beads."""
        ws_file = tmp_path / "00-032-01-test.md"
        ws_file.write_text("""---
ws_id: "00-032-01"
traceability:
  - ac_id: "AC1"
    ac_description: "Test AC"
    test_file: "test.py"
    test_name: "test_func"
    status: "mapped"
    confidence: 1.0
---

# WS Title
- [ ] AC1: Test AC
- [ ] AC2: Another AC
""")

        with patch.object(service, '_get_ws_task', return_value=None):
            with patch.object(service, '_get_ws_content_from_markdown', return_value=ws_file.read_text()):
                report = service.check_traceability("00-032-01")

        assert report.ws_id == "00-032-01"
        assert len(report.mappings) == 2
        assert report.mapped_acs == 1
        assert report.missing_acs == 1
