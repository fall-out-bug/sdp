"""GitHub integration for /design command."""

from pathlib import Path

from sdp.github.client import GitHubClient
from sdp.github.config import GitHubConfig
from sdp.github.milestone_manager import MilestoneManager
from sdp.github.sync_service import SyncResult, SyncService


class DesignGitHubIntegration:
    """Integrate GitHub issue creation into /design command workflow.

    When /design generates workstream files, this class:
    - Creates milestone for feature
    - Creates GitHub issues for each WS
    - Updates WS frontmatter with github_issue numbers
    - Assigns issues to milestone

    Attributes:
        _client: GitHub API client
    """

    def __init__(self, client: GitHubClient) -> None:
        """Initialize design integration.

        Args:
            client: GitHub API client
        """
        self._client = client

    @classmethod
    def from_env(cls) -> "DesignGitHubIntegration":
        """Create integration from environment config.

        Returns:
            DesignGitHubIntegration instance

        Raises:
            ValueError: If GitHub config not found
        """
        config = GitHubConfig.from_env()
        client = GitHubClient(config)
        return cls(client)

    def create_issues_for_feature(
        self,
        feature_id: str,
        feature_name: str,
        feature_spec_path: str,
        ws_files: list[Path],
    ) -> dict[str, object]:
        """Create GitHub milestone and issues for feature workstreams.

        Creates milestone with title "Feature {feature_id}: {feature_name}"
        and creates GitHub issue for each workstream file. Assigns all issues
        to the milestone. Updates WS frontmatter with github_issue numbers.

        Args:
            feature_id: Feature ID (e.g., "F160")
            feature_name: Human-readable feature name
            feature_spec_path: Path to feature specification file
            ws_files: List of Path objects to WS markdown files

        Returns:
            Dictionary with:
            - "milestone": GitHub Milestone object
            - "results": list[SyncResult] for each WS

        Raises:
            ValueError: If GitHub not configured
        """
        repo = self._client.get_repo()

        # Create milestone
        milestone_description = f"Issues for {feature_id}: {feature_name}"

        milestone = MilestoneManager(repo).get_or_create_milestone(
            feature_id=feature_id,
            title=feature_name,
            description=milestone_description,
        )

        # Create issues for each WS
        sync_service = SyncService(self._client)
        results: list[SyncResult] = []

        for ws_file in ws_files:
            try:
                result = sync_service.sync_workstream(ws_file)
                results.append(result)
            except Exception as e:
                ws_id = ws_file.stem
                results.append(
                    SyncResult(
                        ws_id=ws_id,
                        action="failed",
                        error=str(e),
                    )
                )

        return {
            "milestone": milestone,
            "results": results,
        }
