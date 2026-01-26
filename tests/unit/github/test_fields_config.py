"""Tests for GitHub Project Fields configuration."""

import json
import pytest
from pathlib import Path

from sdp.github.fields_config import (
    FieldMapping,
    FieldsConfig,
    load_config,
    save_default_config,
)


def test_field_mapping_defaults() -> None:
    """Test FieldMapping default values."""
    mapping = FieldMapping(
        ws_field="status",
        gh_field_name="Status",
        gh_field_type="single_select",
    )

    assert mapping.ws_field == "status"
    assert mapping.gh_field_name == "Status"
    assert mapping.gh_field_type == "single_select"
    assert mapping.options is None


def test_field_mapping_with_options() -> None:
    """Test FieldMapping with options."""
    mapping = FieldMapping(
        ws_field="size",
        gh_field_name="Size",
        gh_field_type="single_select",
        options={"small": "opt1", "medium": "opt2"},
    )

    assert mapping.get_option_id("small") == "opt1"
    assert mapping.get_option_id("medium") == "opt2"
    assert mapping.get_option_id("large") is None


def test_fields_config_defaults() -> None:
    """Test FieldsConfig default values."""
    config = FieldsConfig()

    assert config.project_name == "SDP"
    assert config.status_field.ws_field == "status"
    assert config.size_field.ws_field == "size"
    assert config.feature_field.ws_field == "feature"
    assert config.assignee_field.ws_field == "assignee"


def test_fields_config_get_all_mappings() -> None:
    """Test getting all field mappings."""
    config = FieldsConfig()
    mappings = config.get_all_mappings()

    assert len(mappings) == 4  # status, size, feature, assignee


def test_fields_config_get_mapping_for_ws_field() -> None:
    """Test finding mapping by workstream field."""
    config = FieldsConfig()

    mapping = config.get_mapping_for_ws_field("status")
    assert mapping is not None
    assert mapping.gh_field_name == "Status"

    mapping = config.get_mapping_for_ws_field("nonexistent")
    assert mapping is None


def test_fields_config_to_file(tmp_path: Path) -> None:
    """Test saving configuration to file."""
    config_file = tmp_path / "test_config.toml"
    config = FieldsConfig(project_name="TestProject")

    config.to_file(config_file)

    assert config_file.exists()

    content = config_file.read_text()
    assert "TestProject" in content
    assert "status_field" in content
    assert "size_field" in content


def test_fields_config_from_file(tmp_path: Path) -> None:
    """Test loading configuration from file."""
    config_file = tmp_path / "test_config.toml"
    config_file.write_text("""
project_name = "TestProject"

[status_field]
ws_field = "status"
gh_field_name = "Status"
gh_field_type = "single_select"

[size_field]
ws_field = "size"
gh_field_name = "Size"
gh_field_type = "single_select"
options = { small = "opt1", medium = "opt2" }
""")

    config = FieldsConfig.from_file(config_file)

    assert config.project_name == "TestProject"
    assert config.status_field.gh_field_name == "Status"
    assert config.size_field.options == {"small": "opt1", "medium": "opt2"}


def test_fields_config_from_missing_file(tmp_path: Path) -> None:
    """Test loading from missing file returns defaults."""
    config = FieldsConfig.from_file(tmp_path / "nonexistent.toml")

    assert config.project_name == "SDP"  # Default value


def test_load_config() -> None:
    """Test load_config function."""
    config = load_config()

    # Should return default config when file doesn't exist
    assert isinstance(config, FieldsConfig)


def test_save_default_config(tmp_path: Path) -> None:
    """Test save_default_config function."""
    config_file = tmp_path / "github_fields.toml"

    save_default_config(config_file)

    assert config_file.exists()


def test_fields_config_custom_fields(tmp_path: Path) -> None:
    """Test custom fields in configuration."""
    config_file = tmp_path / "custom.toml"
    # Use array of tables syntax
    config_file.write_text("""
project_name = "CustomProject"

[[custom_fields]]
ws_field = "priority"
gh_field_name = "Priority"
gh_field_type = "single_select"
""")

    config = FieldsConfig.from_file(config_file)

    assert len(config.custom_fields) == 1
    assert config.custom_fields[0].ws_field == "priority"
