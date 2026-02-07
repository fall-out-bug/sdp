"""Intent schema validator for AI-human communication."""

import json
from pathlib import Path
from typing import Any, cast


class ValidationError(Exception):
    """Raised when intent validation fails."""

    def __init__(self, errors: list[str]):
        self.errors = errors
        super().__init__("\n".join(errors))


class IntentValidator:
    """Validates intent against JSON schema."""

    def __init__(self, schema_path: str | None = None):
        if schema_path is None:
            schema_path = "docs/schema/intent.schema.json"
        self._schema_path = Path(schema_path)
        self._schema = self._load_schema()

    def validate(self, intent: dict[str, Any]) -> None:
        """Validate intent against schema.

        Args:
            intent: Intent dictionary to validate

        Raises:
            ValidationError: If intent is invalid
        """
        import jsonschema

        try:
            jsonschema.validate(instance=intent, schema=self._schema)
        except jsonschema.ValidationError as e:
            raise ValidationError([self._format_error(e)])

    def validate_file(self, intent_path: str) -> dict[str, Any]:
        """Validate intent file.

        Args:
            intent_path: Path to intent JSON file

        Returns:
            Parsed and validated intent dict

        Raises:
            ValidationError: If file is invalid
        """
        path = Path(intent_path)
        if not path.exists():
            raise ValidationError([f"File not found: {intent_path}"])

        intent: dict[str, Any] = cast(dict[str, Any], json.loads(path.read_text()))
        self.validate(intent)
        return intent

    def _load_schema(self) -> dict[str, object]:
        """Load JSON schema from file or embedded."""
        if not self._schema_path.exists():
            # Fallback to embedded schema
            return self._embedded_schema()

        schema: dict[str, object] = json.loads(self._schema_path.read_text())
        return schema

    def _embedded_schema(self) -> dict[str, object]:
        """Embedded schema for when file doesn't exist."""
        return {
            "$schema": "http://json-schema.org/draft-07/schema#",
            "type": "object",
            "required": ["problem", "users", "success_criteria"],
            "properties": {
                "problem": {"type": "string", "minLength": 50},
                "users": {"type": "array", "minItems": 1},
                "success_criteria": {"type": "array", "minItems": 1}
            }
        }

    def _format_error(self, error: object) -> str:
        """Format jsonschema error for readability."""
        # jsonschema ValidationError has .path and .message attributes
        if hasattr(error, "path"):
            path = " -> ".join(str(p) for p in error.path) if error.path else "root"
        else:
            path = "root"
        message = getattr(error, "message", str(error))
        return f"{path}: {message}"
