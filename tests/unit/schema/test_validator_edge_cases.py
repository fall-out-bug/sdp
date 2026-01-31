"""Tests for intent schema validator."""

from pathlib import Path
import json
import pytest
from unittest.mock import Mock, patch

from sdp.schema.validator import IntentValidator, ValidationError


class TestValidatorInit:
    """Tests for validator initialization."""

    def test_init_with_custom_schema_path(self, tmp_path):
        """Test initializes with custom schema path."""
        schema_file = tmp_path / "custom.schema.json"
        schema_file.write_text(json.dumps({
            "$schema": "http://json-schema.org/draft-07/schema#",
            "type": "object",
            "required": ["title"],
            "properties": {"title": {"type": "string"}}
        }))

        validator = IntentValidator(schema_path=str(schema_file))

        assert validator._schema_path == schema_file
        assert validator._schema["type"] == "object"

    def test_init_with_default_path(self, tmp_path, monkeypatch):
        """Test initializes with default schema path."""
        monkeypatch.chdir(tmp_path)
        schema_dir = tmp_path / "docs" / "schema"
        schema_dir.mkdir(parents=True)
        schema_file = schema_dir / "intent.schema.json"
        schema_file.write_text(json.dumps({
            "$schema": "http://json-schema.org/draft-07/schema#",
            "type": "object"
        }))

        validator = IntentValidator()

        assert validator._schema_path == Path("docs/schema/intent.schema.json")

    def test_init_fallback_to_embedded(self, tmp_path, monkeypatch):
        """Test falls back to embedded schema when file doesn't exist."""
        monkeypatch.chdir(tmp_path)

        validator = IntentValidator()

        assert validator._schema["type"] == "object"
        assert "problem" in validator._schema["required"]


class TestValidate:
    """Tests for validate method."""

    def test_validate_success(self, tmp_path):
        """Test validates correct intent."""
        schema_file = tmp_path / "test.schema.json"
        schema_file.write_text(json.dumps({
            "$schema": "http://json-schema.org/draft-07/schema#",
            "type": "object",
            "required": ["title"],
            "properties": {"title": {"type": "string"}}
        }))

        validator = IntentValidator(schema_path=str(schema_file))
        intent = {"title": "Test Intent"}

        # Should not raise
        validator.validate(intent)

    def test_validate_failure(self, tmp_path):
        """Test raises ValidationError for invalid intent."""
        schema_file = tmp_path / "test.schema.json"
        schema_file.write_text(json.dumps({
            "$schema": "http://json-schema.org/draft-07/schema#",
            "type": "object",
            "required": ["title"],
            "properties": {"title": {"type": "string"}}
        }))

        validator = IntentValidator(schema_path=str(schema_file))
        intent = {"title": 123}  # Wrong type

        with pytest.raises(ValidationError) as exc_info:
            validator.validate(intent)

        assert len(exc_info.value.errors) == 1


class TestValidateFile:
    """Tests for validate_file method."""

    def test_validate_file_not_found(self, tmp_path):
        """Test raises ValidationError when file doesn't exist."""
        validator = IntentValidator()

        with pytest.raises(ValidationError, match="File not found"):
            validator.validate_file(str(tmp_path / "nonexistent.json"))

    def test_validate_file_success(self, tmp_path):
        """Test validates file successfully."""
        schema_file = tmp_path / "test.schema.json"
        schema_file.write_text(json.dumps({
            "$schema": "http://json-schema.org/draft-07/schema#",
            "type": "object",
            "required": ["title"],
            "properties": {"title": {"type": "string"}}
        }))

        intent_file = tmp_path / "intent.json"
        intent_file.write_text(json.dumps({"title": "Test"}))

        validator = IntentValidator(schema_path=str(schema_file))
        result = validator.validate_file(str(intent_file))

        assert result["title"] == "Test"

    def test_validate_file_invalid_content(self, tmp_path):
        """Test raises ValidationError for invalid content."""
        schema_file = tmp_path / "test.schema.json"
        schema_file.write_text(json.dumps({
            "$schema": "http://json-schema.org/draft-07/schema#",
            "type": "object",
            "required": ["title"],
            "properties": {"title": {"type": "string"}}
        }))

        intent_file = tmp_path / "intent.json"
        intent_file.write_text(json.dumps({"title": 123}))

        validator = IntentValidator(schema_path=str(schema_file))

        with pytest.raises(ValidationError):
            validator.validate_file(str(intent_file))


class TestLoadSchema:
    """Tests for _load_schema method."""

    def test_load_schema_from_file(self, tmp_path):
        """Test loads schema from file."""
        schema_file = tmp_path / "test.schema.json"
        schema_file.write_text(json.dumps({
            "$schema": "http://json-schema.org/draft-07/schema#",
            "type": "object",
            "properties": {"custom": {"type": "boolean"}}
        }))

        validator = IntentValidator(schema_path=str(schema_file))

        assert "custom" in validator._schema["properties"]

    def test_load_schema_fallback_embedded(self, tmp_path):
        """Test falls back to embedded schema."""
        validator = IntentValidator(schema_path=str(tmp_path / "nonexistent.json"))

        # Should use embedded schema
        assert validator._schema["type"] == "object"
        assert "problem" in validator._schema["required"]


class TestEmbeddedSchema:
    """Tests for _embedded_schema method."""

    def test_embedded_schema_structure(self):
        """Test embedded schema has required structure."""
        validator = IntentValidator(schema_path="/nonexistent/path.json")

        schema = validator._schema

        assert schema["$schema"] == "http://json-schema.org/draft-07/schema#"
        assert schema["type"] == "object"
        assert "problem" in schema["required"]
        assert "users" in schema["required"]
        assert "success_criteria" in schema["required"]

    def test_embedded_schema_validates(self):
        """Test embedded schema validates correct intent."""
        validator = IntentValidator(schema_path="/nonexistent/path.json")

        intent = {
            "problem": "A" * 50,  # At least 50 chars
            "users": ["user1"],
            "success_criteria": ["criterion1"]
        }

        # Should not raise
        validator.validate(intent)


class TestFormatError:
    """Tests for _format_error method."""

    def test_format_error_with_path(self):
        """Test formats error with path."""
        validator = IntentValidator(schema_path="/nonexistent/path.json")

        # Create mock error with path
        error = Mock()
        error.path = ["field", "nested"]
        error.message = "Invalid value"

        result = validator._format_error(error)

        assert "field -> nested" in result
        assert "Invalid value" in result

    def test_format_error_without_path(self):
        """Test formats error without path."""
        validator = IntentValidator(schema_path="/nonexistent/path.json")

        # Create mock error without path
        error = Mock()
        error.path = []
        error.message = "Invalid root"

        result = validator._format_error(error)

        assert "root" in result
        assert "Invalid root" in result

    def test_format_error_no_path_attribute(self):
        """Test formats error when no path attribute."""
        validator = IntentValidator(schema_path="/nonexistent/path.json")

        # Create mock error without path attribute
        error = Mock(spec=[])
        del error.path  # Remove path attribute

        result = validator._format_error(error)

        assert "root" in result


class TestValidationErrorException:
    """Tests for ValidationError exception."""

    def test_validation_error_init(self):
        """Test ValidationError initialization."""
        errors = ["Error 1", "Error 2"]

        exc = ValidationError(errors)

        assert exc.errors == errors
        assert "Error 1" in str(exc)
        assert "Error 2" in str(exc)

    def test_validation_error_single(self):
        """Test ValidationError with single error."""
        errors = ["Single error"]

        exc = ValidationError(errors)

        assert exc.errors == errors
        assert str(exc) == "Single error"
