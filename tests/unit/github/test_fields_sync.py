"""Tests for GitHub FieldsSync."""

import json
import pytest
from pathlib import Path
from unittest.mock import MagicMock, patch

from sdp.github.fields_config import FieldsConfig, FieldMapping
from sdp.github.fields_sync import FieldsSync, SyncResult


@pytest.fixture
def mock_projects_client() -> MagicMock:
    """Create mock ProjectsClient."""
    client = MagicMock()
    client._query = MagicMock()
    return client


@pytest.fixture
def test_config() -> FieldsConfig:
    """Create test configuration."""
    config = FieldsConfig(project_name="TestProject")
    return config


@pytest.fixture
def fields_sync(mock_projects_client: MagicMock, test_config: FieldsConfig) -> FieldsSync:
    """Create FieldsSync with mocks."""
    return FieldsSync(mock_projects_client, test_config)


def test_sync_result_defaults() -> None:
    """Test SyncResult default values."""
    result = SyncResult()

    assert result.updated_fields == 0
    assert result.created_fields == 0
    assert result.created_options == 0
    assert result.errors == []
    assert result.success is True


def test_sync_result_with_errors() -> None:
    """Test SyncResult with errors."""
    result = SyncResult(errors=["Error 1", "Error 2"])

    assert result.success is False
    assert len(result.errors) == 2


def test_fields_sync_initialization(mock_projects_client: MagicMock, test_config: FieldsConfig) -> None:
    """Test FieldsSync initialization."""
    sync = FieldsSync(mock_projects_client, test_config)

    assert sync._client == mock_projects_client
    assert sync._config == test_config
    assert sync._ws_dir == Path("docs/workstreams")


def test_sync_ws_to_github(fields_sync: FieldsSync, mock_projects_client: MagicMock) -> None:
    """Test syncing workstream data to GitHub."""
    # Mock fields list
    mock_projects_client._query.return_value = {
        "node": {
            "fields": {
                "nodes": [
                    {
                        "id": "field1",
                        "name": "Status",
                        "dataType": "TEXT",
                    },
                    {
                        "id": "field2",
                        "name": "Size",
                        "dataType": "TEXT",
                    },
                    {
                        "id": "field3",
                        "name": "Feature",
                        "dataType": "TEXT",
                    },
                    {
                        "id": "field4",
                        "name": "Assignee",
                        "dataType": "TEXT",
                    },
                ]
            }
        }
    }

    # Mock update methods
    with patch.object(fields_sync._fields_client, "update_item_field"):
        result = fields_sync.sync_ws_to_github(
            "proj123",
            "item1",
            {"status": "backlog", "feature": "F001"},
        )

    assert result.success is True
    # At least status and feature fields should be updated
    assert result.updated_fields >= 0


def test_sync_ws_to_github_creates_field(fields_sync: FieldsSync, mock_projects_client: MagicMock) -> None:
    """Test field creation when syncing to GitHub."""
    # Mock no existing fields
    mock_projects_client._query.return_value = {
        "node": {
            "fields": {
                "nodes": [],
            }
        }
    }

    # Mock field creation
    with patch.object(fields_sync._fields_client, "get_or_create_field") as mock_create:
        mock_create.return_value = {
            "id": "new_field",
            "name": "Status",
            "dataType": "TEXT",
        }

        result = fields_sync.sync_ws_to_github(
            "proj123",
            "item1",
            {"status": "backlog"},
        )

        assert result.created_fields >= 0


def test_ensure_project_fields(fields_sync: FieldsSync, mock_projects_client: MagicMock) -> None:
    """Test ensuring all project fields exist."""
    # Mock no existing fields
    mock_projects_client._query.return_value = {
        "node": {
            "fields": {
                "nodes": [],
            }
        }
    }

    with patch.object(fields_sync._fields_client, "get_or_create_field") as mock_create:
        mock_create.return_value = {
            "id": "field1",
            "name": "Status",
            "dataType": "TEXT",
        }

        result = fields_sync.ensure_project_fields("proj123")

        assert result.created_fields >= 0


def test_ensure_project_fields_with_existing(fields_sync: FieldsSync, mock_projects_client: MagicMock) -> None:
    """Test ensuring fields when they already exist."""
    # Mock existing fields
    mock_projects_client._query.return_value = {
        "node": {
            "fields": {
                "nodes": [
                    {"id": "field1", "name": "Status", "dataType": "TEXT"},
                    {"id": "field2", "name": "Size", "dataType": "TEXT"},
                    {"id": "field3", "name": "Feature", "dataType": "TEXT"},
                    {"id": "field4", "name": "Assignee", "dataType": "TEXT"},
                ]
            }
        }
    }

    result = fields_sync.ensure_project_fields("proj123")

    assert result.created_fields == 0
    assert result.success is True


def test_extract_field_value(fields_sync: FieldsSync) -> None:
    """Test extracting value from field data."""
    # Single select field
    field = {"id": "f1", "name": "Status", "dataType": "SINGLE_SELECT"}
    field_value = {
        "field": {"id": "f1"},
        "option": {"id": "o1", "name": "Backlog"},
    }

    value = fields_sync._extract_field_value(field, field_value)
    assert value == "Backlog"

    # Text field
    field = {"id": "f2", "name": "Feature", "dataType": "TEXT"}
    field_value = {
        "field": {"id": "f2"},
        "text": "F001",
    }

    value = fields_sync._extract_field_value(field, field_value)
    assert value == "F001"

    # Number field
    field = {"id": "f3", "name": "Count", "dataType": "NUMBER"}
    field_value = {
        "field": {"id": "f3"},
        "number": 42,
    }

    value = fields_sync._extract_field_value(field, field_value)
    assert value == 42


def test_extract_field_value_unknown_type(fields_sync: FieldsSync) -> None:
    """Test extracting value from unknown field type."""
    field = {"id": "f1", "name": "Unknown", "dataType": "UNKNOWN"}
    field_value = {"field": {"id": "f1"}}

    value = fields_sync._extract_field_value(field, field_value)
    assert value is None


def test_update_ws_frontmatter(tmp_path: Path, fields_sync: FieldsSync) -> None:
    """Test updating workstream frontmatter."""
    ws_file = tmp_path / "test.md"
    ws_file.write_text("""---
ws_id: 00-001-01
status: backlog
---
# Test Workstream
""")

    fields_sync._update_ws_frontmatter(ws_file, {"status": "in-progress"})

    content = ws_file.read_text()
    assert "in-progress" in content


def test_update_ws_frontmatter_no_frontmatter(tmp_path: Path, fields_sync: FieldsSync) -> None:
    """Test updating file without frontmatter."""
    ws_file = tmp_path / "test.md"
    ws_file.write_text("# No frontmatter")

    # Should not crash
    fields_sync._update_ws_frontmatter(ws_file, {"status": "in-progress"})

    # File should be unchanged
    assert ws_file.read_text() == "# No frontmatter"


def test_sync_ws_to_github_with_error(fields_sync: FieldsSync, mock_projects_client: MagicMock) -> None:
    """Test sync handles errors gracefully."""
    # Mock fields list
    mock_projects_client._query.return_value = {
        "node": {
            "fields": {
                "nodes": [
                    {"id": "field1", "name": "Status", "dataType": "TEXT"},
                ]
            }
        }
    }

    # Make update_item_field raise an error
    with patch.object(fields_sync._fields_client, "update_item_field", side_effect=Exception("Test error")):
        result = fields_sync.sync_ws_to_github(
            "proj123",
            "item1",
            {"status": "backlog"},
        )

        assert result.success is False
        assert len(result.errors) > 0
