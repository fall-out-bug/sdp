"""Bidirectional synchronization for GitHub Project custom fields."""

from __future__ import annotations

import logging
from dataclasses import dataclass
from pathlib import Path
from typing import TYPE_CHECKING, Any

if TYPE_CHECKING:
    from .fields_client import FieldsClient
    from .fields_config import FieldsConfig

from .fields_client import FieldsClient
from .fields_config import FieldsConfig, load_config

logger = logging.getLogger(__name__)


@dataclass
class SyncResult:
    """Result of a synchronization operation."""

    updated_fields: int = 0
    created_fields: int = 0
    created_options: int = 0
    errors: list[str] | None = None

    def __post_init__(self) -> None:
        if self.errors is None:
            self.errors = []

    @property
    def success(self) -> bool:
        """Check if sync was successful."""
        return len(self.errors) == 0


class FieldsSync:
    """Manages bidirectional sync between workstream files and GitHub Project fields.

    Features:
    - Sync workstream frontmatter to GitHub Project custom fields
    - Sync GitHub changes back to workstream files
    - Auto-create missing fields and options
    - Configurable field mappings
    """

    def __init__(
        self,
        projects_client: Any,  # ProjectsClient
        config: FieldsConfig | None = None,
        ws_dir: str | Path = "docs/workstreams",
    ) -> None:
        """Initialize fields sync.

        Args:
            projects_client: ProjectsClient instance
            config: FieldsConfig (uses default if None)
            ws_dir: Path to workstreams directory
        """
        self._client = projects_client
        self._config = config or load_config()
        self._ws_dir = Path(ws_dir)
        self._fields_client = FieldsClient(projects_client)

    def sync_ws_to_github(
        self,
        project_id: str,
        item_id: str,
        ws_data: dict[str, Any],
    ) -> SyncResult:
        """Sync workstream data to GitHub Project item.

        Args:
            project_id: Project global ID
            item_id: Project item global ID
            ws_data: Workstream frontmatter data

        Returns:
            SyncResult with sync statistics
        """
        result = SyncResult()

        # Get all fields in project
        fields = self._fields_client.list_fields(project_id)
        field_by_name: dict[str, dict[str, Any]] = {
            f["name"]: f for f in fields
        }

        # Sync each mapped field
        for mapping in self._config.get_all_mappings():
            ws_value = ws_data.get(mapping.ws_field)

            if ws_value is None:
                continue

            # Get or create field
            field = field_by_name.get(mapping.gh_field_name)

            if not field:
                field = self._fields_client.get_or_create_field(
                    project_id,
                    mapping,
                )
                result.created_fields += 1
                field_by_name[field["name"]] = field

            # Convert value for single_select
            if mapping.gh_field_type == "single_select":
                if isinstance(ws_value, str):
                    ws_value = ws_value.replace("_", " ").replace("-", " ").title()
                elif hasattr(ws_value, "value"):
                    ws_value = str(ws_value.value).replace("_", " ").replace("-", " ").title()

                # Ensure option exists
                if field.get("dataType") == "SINGLE_SELECT":
                    option_names = [ws_value] if ws_value else []
                    name_to_id = self._fields_client.ensure_options_for_field(
                        field,
                        option_names,
                    )
                    result.created_options += len(
                        [n for n in option_names if n not in name_to_id]
                    )

            # Update the field value
            try:
                self._fields_client.update_item_field(
                    project_id,
                    item_id,
                    field,
                    str(ws_value),
                )
                result.updated_fields += 1
            except Exception as e:
                result.errors.append(
                    f"Failed to update {mapping.gh_field_name}: {e}"
                )
                logger.exception(f"Failed to update field {mapping.gh_field_name}")

        return result

    def sync_github_to_ws(
        self,
        project_id: str,
        item_id: str,
        ws_file: Path,
    ) -> SyncResult:
        """Sync GitHub Project item data to workstream file.

        Args:
            project_id: Project global ID
            item_id: Project item global ID
            ws_file: Path to workstream file

        Returns:
            SyncResult with sync statistics
        """
        from ..core import parse_workstream

        result = SyncResult()

        # Get item data from GitHub
        item_data = self._get_item_data(project_id, item_id)
        if not item_data:
            result.errors.append("Could not fetch item data from GitHub")
            return result

        # Parse workstream file
        try:
            ws = parse_workstream(ws_file)
        except Exception as e:
            result.errors.append(f"Failed to parse workstream: {e}")
            return result

        # Get fields and mappings
        fields = self._fields_client.list_fields(project_id)
        field_by_name: dict[str, dict[str, Any]] = {
            f["name"]: f for f in fields
        }

        # Update frontmatter from GitHub fields
        updates: dict[str, Any] = {}
        for mapping in self._config.get_all_mappings():
            field = field_by_name.get(mapping.gh_field_name)
            if not field:
                continue

            field_values = item_data.get("fieldValues", [])
            for fv in field_values:
                if fv.get("field", {}).get("id") == field["id"]:
                    # Extract value based on field type
                    value = self._extract_field_value(field, fv)
                    if value is not None:
                        updates[mapping.ws_field] = value

        if updates:
            self._update_ws_frontmatter(ws_file, updates)
            result.updated_fields = len(updates)

        return result

    def ensure_project_fields(self, project_id: str) -> SyncResult:
        """Ensure all configured fields exist in the project.

        Creates missing fields and options.

        Args:
            project_id: Project global ID

        Returns:
            SyncResult with sync statistics
        """
        result = SyncResult()

        for mapping in self._config.get_all_mappings():
            try:
                field = self._fields_client.get_field_by_name(
                    project_id,
                    mapping.gh_field_name,
                )

                if not field:
                    field = self._fields_client.get_or_create_field(
                        project_id,
                        mapping,
                    )
                    result.created_fields += 1

                # Ensure options for single_select fields
                if mapping.gh_field_type == "single_select" and mapping.options:
                    if field.get("dataType") == "SINGLE_SELECT":
                        option_names = list(mapping.options.keys())
                        created = self._fields_client.ensure_options_for_field(
                            field,
                            option_names,
                        )
                        result.created_options += len(created)

            except Exception as e:
                result.errors.append(
                    f"Failed to ensure field {mapping.gh_field_name}: {e}"
                )
                logger.exception(f"Failed to ensure field {mapping.gh_field_name}")

        return result

    def _get_item_data(
        self,
        project_id: str,
        item_id: str,
    ) -> dict[str, Any] | None:
        """Get project item data including field values.

        Args:
            project_id: Project global ID
            item_id: Item global ID

        Returns:
            Item data dict or None
        """
        query = """
        query($projectId: ID!, $itemId: ID!) {
          node(id: $projectId) {
            ... on ProjectV2 {
              item(id: $itemId) {
                id
                fieldValues(first: 20) {
                  nodes {
                    ... on ProjectV2ItemFieldSingleSelectValue {
                      field {
                        id
                        name
                        dataType
                      }
                      optionId
                      option {
                        id
                        name
                      }
                    }
                    ... on ProjectV2ItemFieldTextValue {
                      field {
                        id
                        name
                        dataType
                      }
                      text
                    }
                    ... on ProjectV2ItemFieldNumberValue {
                      field {
                        id
                        name
                        dataType
                      }
                      number
                    }
                  }
                }
              }
            }
          }
        }
        """

        try:
            result = self._client._query(
                query,
                {"projectId": project_id, "itemId": item_id},
            )

            item = result.get("node", {}).get("item")
            if item:
                # Normalize response
                field_values = []
                for fv in item.get("fieldValues", {}).get("nodes", []):
                    field_values.append(fv)

                return {"id": item["id"], "fieldValues": field_values}

        except Exception as e:
            logger.error(f"Failed to get item data: {e}")

        return None

    def _extract_field_value(
        self,
        field: dict[str, Any],
        field_value: dict[str, Any],
    ) -> str | None:
        """Extract value from field value data.

        Args:
            field: Field dict
            field_value: Field value data from GitHub

        Returns:
            Extracted value or None
        """
        data_type = field.get("dataType", "")

        if data_type == "SINGLE_SELECT":
            option = field_value.get("option")
            return option["name"] if option else None

        elif data_type == "TEXT":
            return field_value.get("text")

        elif data_type == "NUMBER":
            return field_value.get("number")

        return None

    def _update_ws_frontmatter(
        self,
        ws_file: Path,
        updates: dict[str, Any],
    ) -> None:
        """Update workstream frontmatter with new values.

        Args:
            ws_file: Path to workstream file
            updates: Dict of field names to new values
        """
        content = ws_file.read_text()

        if not content.startswith("---"):
            logger.warning(f"Workstream file has no frontmatter: {ws_file}")
            return

        # Find frontmatter section
        parts = content.split("---", 2)
        if len(parts) < 3:
            logger.warning(f"Invalid frontmatter in: {ws_file}")
            return

        frontmatter = parts[1]
        body = parts[2]

        # Update frontmatter values
        import re
        import yaml

        try:
            data = yaml.safe_load(frontmatter) or {}
            data.update(updates)

            # Write back
            new_frontmatter = yaml.dump(data, default_flow_style=False)
            new_content = f"---\n{new_frontmatter}---{body}"

            ws_file.write_text(new_content)

        except Exception as e:
            logger.error(f"Failed to update frontmatter: {e}")
