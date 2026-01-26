"""Tests for GitHub FieldsClient."""

import json
import pytest
from unittest.mock import MagicMock, patch

from sdp.github.fields_client import FieldsClient
from sdp.github.projects_client import ProjectsClient


@pytest.fixture
def mock_projects_client() -> MagicMock:
    """Create mock ProjectsClient."""
    client = MagicMock(spec=ProjectsClient)
    client._query = MagicMock()
    return client


@pytest.fixture
def fields_client(mock_projects_client: MagicMock) -> FieldsClient:
    """Create FieldsClient with mock."""
    return FieldsClient(mock_projects_client)


def test_fields_client_initialization(fields_client: FieldsClient) -> None:
    """Test FieldsClient initialization."""
    assert fields_client._client is not None


def test_list_fields(fields_client: FieldsClient, mock_projects_client: MagicMock) -> None:
    """Test listing project fields."""
    mock_projects_client._query.return_value = {
        "node": {
            "fields": {
                "nodes": [
                    {
                        "id": "field1",
                        "name": "Status",
                        "dataType": "SINGLE_SELECT",
                        "options": [
                            {"id": "opt1", "name": "Backlog"},
                            {"id": "opt2", "name": "In Progress"},
                        ],
                    },
                    {
                        "id": "field2",
                        "name": "Size",
                        "dataType": "TEXT",
                    },
                ]
            }
        }
    }

    fields = fields_client.list_fields("proj123")

    assert len(fields) == 2
    assert fields[0]["name"] == "Status"
    assert fields[0]["dataType"] == "SINGLE_SELECT"
    assert len(fields[0]["options"]) == 2


def test_get_field_by_name(fields_client: FieldsClient, mock_projects_client: MagicMock) -> None:
    """Test finding field by name."""
    mock_projects_client._query.return_value = {
        "node": {
            "fields": {
                "nodes": [
                    {"id": "field1", "name": "Status", "dataType": "TEXT"},
                    {"id": "field2", "name": "Size", "dataType": "TEXT"},
                ]
            }
        }
    }

    field = fields_client.get_field_by_name("proj123", "Status")

    assert field is not None
    assert field["id"] == "field1"

    field = fields_client.get_field_by_name("proj123", "Nonexistent")
    assert field is None


def test_get_option_id_by_name(fields_client: FieldsClient) -> None:
    """Test getting option ID by name."""
    field = {
        "id": "field1",
        "name": "Status",
        "dataType": "SINGLE_SELECT",
        "options": [
            {"id": "opt1", "name": "Backlog"},
            {"id": "opt2", "name": "In Progress"},
        ],
    }

    option_id = fields_client.get_option_id_by_name(field, "backlog")  # Case insensitive

    assert option_id == "opt1"

    option_id = fields_client.get_option_id_by_name(field, "Nonexistent")
    assert option_id is None


def test_create_text_field(fields_client: FieldsClient, mock_projects_client: MagicMock) -> None:
    """Test creating a text field."""
    mock_projects_client._query.return_value = {
        "createProjectV2Field": {
            "projectV2Field": {
                "id": "field3",
                "name": "Custom",
                "dataType": "TEXT",
            }
        }
    }

    field = fields_client.create_text_field("proj123", "Custom")

    assert field["name"] == "Custom"
    assert field["dataType"] == "TEXT"


def test_create_field_option(fields_client: FieldsClient, mock_projects_client: MagicMock) -> None:
    """Test creating a field option."""
    mock_projects_client._query.return_value = {
        "createProjectV2SingleSelectOption": {
            "projectV2SingleSelectOption": {
                "id": "opt3",
                "name": "Review",
            }
        }
    }

    option = fields_client.create_field_option("field1", "Review")

    assert option["name"] == "Review"
    assert option["id"] == "opt3"


def test_update_item_field_single_select(fields_client: FieldsClient, mock_projects_client: MagicMock) -> None:
    """Test updating item single-select field."""
    field = {
        "id": "field1",
        "name": "Status",
        "dataType": "SINGLE_SELECT",
        "options": [
            {"id": "opt1", "name": "Backlog"},
        ],
    }

    fields_client.update_item_field("proj123", "item1", field, "Backlog")

    mock_projects_client.update_item_single_select.assert_called_once()


def test_update_item_field_text(fields_client: FieldsClient, mock_projects_client: MagicMock) -> None:
    """Test updating item text field."""
    field = {
        "id": "field2",
        "name": "Feature",
        "dataType": "TEXT",
    }

    fields_client.update_item_field("proj123", "item1", field, "F001")

    mock_projects_client.update_item_field.assert_called_once()


def test_ensure_options_for_field(fields_client: FieldsClient, mock_projects_client: MagicMock) -> None:
    """Test ensuring options exist for a field."""
    field = {
        "id": "field1",
        "name": "Status",
        "dataType": "SINGLE_SELECT",
        "options": [
            {"id": "opt1", "name": "Backlog"},
        ],
    }

    mock_projects_client._query.return_value = {
        "createProjectV2SingleSelectOption": {
            "projectV2SingleSelectOption": {
                "id": "opt2",
                "name": "In Progress",
            }
        }
    }

    name_to_id = fields_client.ensure_options_for_field(field, ["Backlog", "In Progress"])

    assert "Backlog" in name_to_id
    assert "In Progress" in name_to_id
    assert name_to_id["Backlog"] == "opt1"
    assert name_to_id["In Progress"] == "opt2"
