"""GitHub issue synchronization from WS files."""

from sdp.github.client import GitHubClient
from sdp.github.issue_formatter import format_issue_body
from sdp.github.label_manager import LabelManager
from sdp.github.ws_parser import WSMetadata


class IssueSync:
    """Synchronize workstream files to GitHub issues.

    Coordinates:
    - WS parsing
    - Label management
    - Issue creation
    - Body formatting

    Attributes:
        _client: GitHub API client
        _label_manager: Label manager
    """

    def __init__(self, client: GitHubClient, label_manager: LabelManager) -> None:
        """Initialize issue sync.

        Args:
            client: GitHubClient instance
            label_manager: LabelManager instance
        """
        self._client = client
        self._label_manager = label_manager

    def sync_ws(
        self,
        ws: WSMetadata,
        ws_file_path: str,
        milestone_number: int | None = None,
    ) -> int:
        """Sync workstream to GitHub issue.

        Creates GitHub issue from WS metadata with:
        - Title from WS ID and feature
        - Body formatted with goal, AC, dependencies
        - Labels derived from WS metadata

        Args:
            ws: Parsed WS metadata
            ws_file_path: Relative path to WS file (for links)

        Returns:
            Created issue number

        Raises:
            Exception: If GitHub API call fails
        """
        # Ensure all labels exist first
        self._label_manager.ensure_labels(ws)

        # Derive labels for this WS
        labels = self._label_manager.derive_labels(ws)

        # Format issue body
        body = format_issue_body(ws, ws_file_path)

        # Create title: "WS-100-01: Title (F100)"
        title = f"{ws.ws_id}: {ws.title} ({ws.feature})"

        # Create issue
        issue = self._client.create_issue(
            title=title,
            body=body,
            labels=labels,
            milestone=milestone_number,
        )

        return issue.number
