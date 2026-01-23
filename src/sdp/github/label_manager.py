"""GitHub label management."""

from dataclasses import dataclass

from github import GithubException

from sdp.github.client import GitHubClient
from sdp.github.ws_parser import WSMetadata


@dataclass
class LabelConfig:
    """GitHub label configuration.

    Attributes:
        name: Label name
        color: Hex color (without #)
        description: Label description
    """

    name: str
    color: str
    description: str


class LabelManager:
    """Manage GitHub labels for workstreams.

    Attributes:
        _client: GitHub API client
    """

    def __init__(self, client: GitHubClient) -> None:
        """Initialize label manager.

        Args:
            client: GitHubClient instance
        """
        self._client = client

    def derive_labels(self, ws: WSMetadata) -> list[str]:
        """Derive labels from WS metadata.

        Args:
            ws: Parsed WS metadata

        Returns:
            List of label names to apply
        """
        labels = [
            "workstream",  # Base label for all WS
            f"feature/{ws.feature}",  # Feature label
            f"size/{ws.size}",  # Size label
            f"status/{ws.status}",  # Status label
        ]
        return labels

    def create_label(self, config: LabelConfig) -> None:
        """Create label if not exists.

        Args:
            config: LabelConfig with name, color, description

        Raises:
            GithubException: If API error (except "already exists")
        """
        repo = self._client.get_repo()
        try:
            repo.create_label(
                name=config.name,
                color=config.color,
                description=config.description,
            )
        except GithubException as e:
            # Ignore "already exists" errors (422 status)
            if e.status == 422:
                return
            raise

    def ensure_labels(self, ws: WSMetadata) -> None:
        """Ensure all derived labels exist.

        Creates labels if they don't exist. Silently ignores
        "already exists" errors.

        Args:
            ws: Parsed WS metadata
        """
        labels = self.derive_labels(ws)
        label_configs = self._get_label_configs(labels)

        for config in label_configs:
            self.create_label(config)

    def _get_label_configs(self, label_names: list[str]) -> list[LabelConfig]:
        """Get label configurations for names.

        Args:
            label_names: List of label names

        Returns:
            List of LabelConfig objects
        """
        configs = []
        for name in label_names:
            config = self._config_for_label(name)
            if config:
                configs.append(config)
        return configs

    def _config_for_label(self, name: str) -> LabelConfig | None:
        """Get label config for name.

        Args:
            name: Label name

        Returns:
            LabelConfig or None if unknown label
        """
        label_configs: dict[str, LabelConfig] = {
            "workstream": LabelConfig(
                name="workstream",
                color="0366d6",
                description="Workstream task",
            ),
            "feature/F100": LabelConfig(
                name="feature/F100",
                color="7057ff",
                description="Feature F100",
            ),
            "feature/F150": LabelConfig(
                name="feature/F150",
                color="7057ff",
                description="Feature F150",
            ),
            "feature/F160": LabelConfig(
                name="feature/F160",
                color="7057ff",
                description="Feature F160",
            ),
            "size/SMALL": LabelConfig(
                name="size/SMALL",
                color="d4c5f9",
                description="Small workstream (< 200 LOC)",
            ),
            "size/MEDIUM": LabelConfig(
                name="size/MEDIUM",
                color="9d8de9",
                description="Medium workstream (200-400 LOC)",
            ),
            "size/LARGE": LabelConfig(
                name="size/LARGE",
                color="6b5297",
                description="Large workstream (> 400 LOC)",
            ),
            "status/backlog": LabelConfig(
                name="status/backlog",
                color="cccccc",
                description="Workstream in backlog",
            ),
            "status/active": LabelConfig(
                name="status/active",
                color="ededed",
                description="Workstream in progress",
            ),
            "status/completed": LabelConfig(
                name="status/completed",
                color="6f42c1",
                description="Workstream completed",
            ),
        }

        # Check for exact match first
        if name in label_configs:
            return label_configs[name]

        # Check for feature/* pattern
        if name.startswith("feature/"):
            return LabelConfig(
                name=name,
                color="7057ff",
                description=name,
            )

        # Check for size/* pattern
        if name.startswith("size/"):
            return LabelConfig(
                name=name,
                color="d4c5f9",
                description=name,
            )

        # Check for status/* pattern
        if name.startswith("status/"):
            return LabelConfig(
                name=name,
                color="cccccc",
                description=name,
            )

        return None
