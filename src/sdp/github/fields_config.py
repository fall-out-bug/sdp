"""Configuration for GitHub Project custom fields mapping."""

from __future__ import annotations

import logging
from dataclasses import dataclass, field
from pathlib import Path
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from typing import Any

try:
    import tomllib
except ImportError:
    import tomli as tomllib

try:
    import tomli_w
except ImportError:
    tomli_w = None

logger = logging.getLogger(__name__)


@dataclass
class FieldMapping:
    """Mapping between workstream frontmatter and GitHub Project field."""

    ws_field: str  # Workstream frontmatter field name
    gh_field_name: str  # GitHub Project field display name
    gh_field_type: str  # single_select, text, number, date
    options: dict[str, str] | None = None  # For single_select: ws_value -> option_id mapping

    def get_option_id(self, ws_value: str) -> str | None:
        """Get GitHub option ID for workstream value.

        Args:
            ws_value: Workstream field value

        Returns:
            GitHub option ID or None if not found
        """
        if self.options is None:
            return None
        return self.options.get(ws_value)


@dataclass
class FieldsConfig:
    """Configuration for GitHub Project custom fields synchronization."""

    project_name: str = "SDP"
    # Standard field mappings
    status_field: FieldMapping = field(
        default_factory=lambda: FieldMapping(
            ws_field="status",
            gh_field_name="Status",
            gh_field_type="single_select",
            options={},
        )
    )
    size_field: FieldMapping = field(
        default_factory=lambda: FieldMapping(
            ws_field="size",
            gh_field_name="Size",
            gh_field_type="single_select",
            options={},
        )
    )
    feature_field: FieldMapping = field(
        default_factory=lambda: FieldMapping(
            ws_field="feature",
            gh_field_name="Feature",
            gh_field_type="text",
        )
    )
    assignee_field: FieldMapping = field(
        default_factory=lambda: FieldMapping(
            ws_field="assignee",
            gh_field_name="Assignee",
            gh_field_type="text",
        )
    )
    custom_fields: list[FieldMapping] = field(default_factory=list)

    @classmethod
    def from_file(cls, path: str | Path) -> FieldsConfig:
        """Load configuration from TOML file.

        Args:
            path: Path to config file

        Returns:
            FieldsConfig instance
        """
        path = Path(path)

        if not path.exists():
            logger.info(f"Config file not found: {path}, using defaults")
            return cls()

        try:
            data: dict[str, Any] = tomllib.loads(path.read_text())

            config = cls(project_name=data.get("project_name", "SDP"))

            # Load status field
            if "status_field" in data:
                config.status_field = cls._load_field_mapping(data["status_field"])

            # Load size field
            if "size_field" in data:
                config.size_field = cls._load_field_mapping(data["size_field"])

            # Load feature field
            if "feature_field" in data:
                config.feature_field = cls._load_field_mapping(data["feature_field"])

            # Load assignee field
            if "assignee_field" in data:
                config.assignee_field = cls._load_field_mapping(data["assignee_field"])

            # Load custom fields
            if "custom_fields" in data:
                config.custom_fields = [
                    cls._load_field_mapping(f) for f in data["custom_fields"]
                ]

            return config

        except Exception as e:
            logger.warning(f"Failed to load config from {path}: {e}, using defaults")
            return cls()

    @classmethod
    def _load_field_mapping(cls, data: dict[str, Any]) -> FieldMapping:
        """Load FieldMapping from dict.

        Args:
            data: Field data dict

        Returns:
            FieldMapping instance
        """
        return FieldMapping(
            ws_field=data.get("ws_field", ""),
            gh_field_name=data.get("gh_field_name", ""),
            gh_field_type=data.get("gh_field_type", "text"),
            options=data.get("options"),
        )

    def to_file(self, path: str | Path) -> None:
        """Save configuration to TOML file.

        Args:
            path: Path to save config
        """
        path = Path(path)
        path.parent.mkdir(parents=True, exist_ok=True)

        data: dict[str, Any] = {
            "project_name": self.project_name,
            "status_field": self._field_to_dict(self.status_field),
            "size_field": self._field_to_dict(self.size_field),
            "feature_field": self._field_to_dict(self.feature_field),
            "assignee_field": self._field_to_dict(self.assignee_field),
        }

        if self.custom_fields:
            data["custom_fields"] = [
                self._field_to_dict(f) for f in self.custom_fields
            ]

        # Serialize to TOML
        if tomli_w is not None:
            content = tomli_w.dumps(data)
        else:
            content = self._simple_toml_dumps(data)

        path.write_text(content)

    def _field_to_dict(self, field: FieldMapping) -> dict[str, Any]:
        """Convert FieldMapping to dict.

        Args:
            field: FieldMapping to convert

        Returns:
            Dict representation
        """
        result: dict[str, Any] = {
            "ws_field": field.ws_field,
            "gh_field_name": field.gh_field_name,
            "gh_field_type": field.gh_field_type,
        }
        if field.options:
            result["options"] = field.options
        return result

    def _simple_toml_dumps(self, data: dict[str, Any]) -> str:
        """Simple TOML serialization fallback.

        Args:
            data: Dict to serialize

        Returns:
            TOML formatted string
        """
        lines: list[str] = []

        for key, value in data.items():
            if isinstance(value, str):
                lines.append(f'{key} = "{value}"')
            elif isinstance(value, (int, float, bool)):
                lines.append(f"{key} = {value}")
            elif isinstance(value, dict):
                lines.append(f"\n[{key}]")
                for k, v in value.items():
                    if isinstance(v, str):
                        lines.append(f'{k} = "{v}"')
                    elif isinstance(v, (int, float, bool)):
                        lines.append(f"{k} = {v}")
                    elif isinstance(v, dict):
                        lines.append(f"\n[{key}.{k}]")
                        for ok, ov in v.items():
                            if isinstance(ov, str):
                                lines.append(f'"{ok}" = "{ov}"')
                            else:
                                lines.append(f'"{ok}" = {ov}')
            elif isinstance(value, list):
                for item in value:
                    if isinstance(item, dict):
                        lines.append(f"\n[[{key}]]")
                        for k, v in item.items():
                            if isinstance(v, str):
                                lines.append(f'{k} = "{v}"')
                            elif isinstance(v, dict):
                                # Serialize options dict
                                opt_parts = []
                                for ok, ov in v.items():
                                    opt_parts.append(f'"{ok}" = "{ov}"')
                                lines.append(f'{k} = {{ {", ".join(opt_parts)} }}')

        return "\n".join(lines)

    def get_all_mappings(self) -> list[FieldMapping]:
        """Get all field mappings including custom ones.

        Returns:
            List of all FieldMapping instances
        """
        return [
            self.status_field,
            self.size_field,
            self.feature_field,
            self.assignee_field,
            *self.custom_fields,
        ]

    def get_mapping_for_ws_field(self, ws_field: str) -> FieldMapping | None:
        """Find field mapping for workstream field.

        Args:
            ws_field: Workstream frontmatter field name

        Returns:
            FieldMapping or None if not found
        """
        for mapping in self.get_all_mappings():
            if mapping.ws_field == ws_field:
                return mapping
        return None


def load_config(config_path: str | Path = ".sdp/github_fields.toml") -> FieldsConfig:
    """Load fields configuration.

    Args:
        config_path: Path to config file

    Returns:
        FieldsConfig instance
    """
    return FieldsConfig.from_file(config_path)


def save_default_config(config_path: str | Path = ".sdp/github_fields.toml") -> None:
    """Save default configuration template.

    Args:
        config_path: Path to save config
    """
    config = FieldsConfig()
    config.to_file(config_path)
    logger.info(f"Saved default config to {config_path}")
