"""GitHub Projects custom fields API client."""

from __future__ import annotations

import logging
from typing import TYPE_CHECKING, Any

if TYPE_CHECKING:
    from .fields_config import FieldMapping

from .projects_client import ProjectsClient

logger = logging.getLogger(__name__)


class FieldsClient:
    """Client for managing GitHub Project custom fields.

    Provides methods to:
    - List custom fields in a project
    - Create custom fields
    - Get field options (for single_select fields)
    - Create field options
    - Get field by name
    """

    def __init__(self, projects_client: ProjectsClient) -> None:
        """Initialize fields client.

        Args:
            projects_client: Initialized ProjectsClient instance
        """
        self._client = projects_client

    def list_fields(self, project_id: str) -> list[dict[str, Any]]:
        """List all custom fields in a project.

        Args:
            project_id: Project global ID

        Returns:
            List of field dicts with id, name, type, options
        """
        query = """
        query($projectId: ID!) {
          node(id: $projectId) {
            ... on ProjectV2 {
              fields(first: 50) {
                nodes {
                  id
                  name
                  dataType
                  ... on ProjectV2SingleSelectField {
                    options {
                      id
                      name
                    }
                  }
                }
              }
            }
          }
        }
        """

        result = self._client._query(query, {"projectId": project_id})
        fields = (
            result.get("node", {})
            .get("fields", {})
            .get("nodes", [])
        )

        return fields

    def get_field_by_name(
        self,
        project_id: str,
        field_name: str,
    ) -> dict[str, Any] | None:
        """Find a field by its display name.

        Args:
            project_id: Project global ID
            field_name: Field display name

        Returns:
            Field dict or None if not found
        """
        fields = self.list_fields(project_id)

        for field in fields:
            if field["name"] == field_name:
                return field

        return None

    def create_single_select_field(
        self,
        project_id: str,
        name: str,
        options: list[str],
    ) -> dict[str, Any]:
        """Create a single-select custom field.

        Args:
            project_id: Project global ID
            name: Field display name
            options: List of option names

        Returns:
            Created field dict
        """
        # Create field first
        query = """
        mutation($projectId: ID!, $name: String!) {
          createProjectV2Field(
            input: {
              projectId: $projectId
              dataType: SINGLE_SELECT
              name: $name
            }
          ) {
            projectV2Field {
              id
              name
              dataType
            }
          }
        }
        """

        result = self._client._query(query, {"projectId": project_id, "name": name})
        field = result["createProjectV2Field"]["projectV2Field"]

        # Create options
        for option_name in options:
            self.create_field_option(field["id"], option_name)

        # Return field with options
        return self.get_field_by_name(project_id, name) or field

    def create_text_field(
        self,
        project_id: str,
        name: str,
    ) -> dict[str, Any]:
        """Create a text custom field.

        Args:
            project_id: Project global ID
            name: Field display name

        Returns:
            Created field dict
        """
        query = """
        mutation($projectId: ID!, $name: String!) {
          createProjectV2Field(
            input: {
              projectId: $projectId
              dataType: TEXT
              name: $name
            }
          ) {
            projectV2Field {
              id
              name
              dataType
            }
          }
        }
        """

        result = self._client._query(query, {"projectId": project_id, "name": name})
        return result["createProjectV2Field"]["projectV2Field"]

    def create_field_option(
        self,
        field_id: str,
        option_name: str,
    ) -> dict[str, Any]:
        """Create an option for a single-select field.

        Args:
            field_id: Field global ID
            option_name: Option display name

        Returns:
            Created option dict
        """
        query = """
        mutation($fieldId: ID!, $name: String!) {
          createProjectV2SingleSelectOption(
            input: {
              fieldId: $fieldId
              name: $name
            }
          ) {
            projectV2SingleSelectOption {
              id
              name
            }
          }
        }
        """

        result = self._client._query(query, {"fieldId": field_id, "name": option_name})
        return result["createProjectV2SingleSelectOption"]["projectV2SingleSelectOption"]

    def get_option_id_by_name(
        self,
        field: dict[str, Any],
        option_name: str,
    ) -> str | None:
        """Get option ID by name for a single-select field.

        Args:
            field: Field dict (must have options)
            option_name: Option display name

        Returns:
            Option ID or None if not found
        """
        if "options" not in field:
            return None

        for option in field["options"]:
            # Case-insensitive comparison
            if option["name"].lower() == option_name.lower():
                return option["id"]

        return None

    def get_or_create_field(
        self,
        project_id: str,
        mapping: FieldMapping,
    ) -> dict[str, Any]:
        """Get existing field or create new one.

        Args:
            project_id: Project global ID
            mapping: FieldMapping with configuration

        Returns:
            Field dict
        """
        # Try to find existing field
        field = self.get_field_by_name(project_id, mapping.gh_field_name)

        if field:
            logger.debug(f"Found existing field: {mapping.gh_field_name}")
            return field

        # Create new field based on type
        logger.info(f"Creating field: {mapping.gh_field_name} ({mapping.gh_field_type})")

        if mapping.gh_field_type == "single_select":
            # Create with options from config
            option_names = list(mapping.gh_field_name.values() if isinstance(mapping.gh_field_name, dict) else [])
            return self.create_single_select_field(project_id, mapping.gh_field_name, option_names)
        else:
            return self.create_text_field(project_id, mapping.gh_field_name)

    def ensure_options_for_field(
        self,
        field: dict[str, Any],
        option_names: list[str],
    ) -> dict[str, str]:
        """Ensure all options exist for a single-select field.

        Creates missing options and returns mapping of name -> id.

        Args:
            field: Field dict (must have options)
            option_names: Required option names

        Returns:
            Dict mapping option name to option ID
        """
        name_to_id: dict[str, str] = {}

        # Collect existing options
        existing = set()
        if "options" in field:
            for option in field["options"]:
                name_to_id[option["name"]] = option["id"]
                existing.add(option["name"])

        # Create missing options
        for name in option_names:
            if name not in existing:
                logger.info(f"Creating option: {name}")
                option = self.create_field_option(field["id"], name)
                name_to_id[name] = option["id"]

        return name_to_id

    def update_item_field(
        self,
        project_id: str,
        item_id: str,
        field: dict[str, Any],
        value: str,
    ) -> None:
        """Update a project item's field value.

        Args:
            project_id: Project global ID
            item_id: Item global ID
            field: Field dict
            value: New value
        """
        data_type = field.get("dataType", "")

        if data_type == "SINGLE_SELECT":
            # Get option ID
            option_id = self.get_option_id_by_name(field, value)
            if option_id:
                self._client.update_item_single_select(
                    project_id,
                    item_id,
                    field["id"],
                    option_id,
                )
            else:
                logger.warning(f"Option not found: {value} for field {field['name']}")
        else:
            # Text field
            self._client.update_item_field(
                project_id,
                item_id,
                field["id"],
                value,
            )
