"""GitHub integration for /deploy command."""

import subprocess

from sdp.github.client import GitHubClient
from sdp.github.config import GitHubConfig
from sdp.github.milestone_manager import MilestoneManager


class DeployGitHubIntegration:
    """Integrate GitHub PR creation into /deploy workflow."""

    def __init__(self, github_client: GitHubClient) -> None:
        """Initialize deploy integration.

        Args:
            github_client: GitHub API client

        """
        self._client = github_client
        self._repo = github_client.get_repo()

    def create_pr_with_issues(
        self,
        base_branch: str,
        head_branch: str,
        feature_id: str,  # noqa: ARG002
        feature_name: str,
    ) -> dict[str, str]:
        """Create PR with auto-generated description linking issues.

        Args:
            base_branch: Base branch (usually "main")
            head_branch: Feature branch
            feature_id: Feature ID (e.g., "F60") for future milestone linking
            feature_name: Feature name (e.g., "LMS Integration")

        Returns:
            Dict with PR info (pr_url, description)

        Raises:
            subprocess.CalledProcessError: If git or gh commands fail

        """
        # Generate PR description
        description = self._generate_pr_description(base_branch, feature_name)

        # Create PR via gh CLI
        pr_url = self._create_pr_via_gh(head_branch, base_branch, feature_name, description)

        return {
            "pr_url": pr_url,
            "description": description,
        }

    def _generate_pr_description(self, base_branch: str, feature_name: str) -> str:
        """Generate PR description using script.

        Args:
            base_branch: Base branch
            feature_name: Feature name

        Returns:
            PR description markdown

        Raises:
            subprocess.CalledProcessError: If script fails

        """
        result = subprocess.run(  # noqa: S603
            [  # noqa: S607
                "python3",
                "tools/hw_checker/scripts/generate_pr_description.py",
                base_branch,
                feature_name,
            ],
            capture_output=True,
            text=True,
            check=True,
        )
        return result.stdout

    def _create_pr_via_gh(
        self,
        head_branch: str,
        base_branch: str,
        title: str,
        body: str,
    ) -> str:
        """Create PR using gh CLI.

        Args:
            head_branch: Feature branch
            base_branch: Base branch
            title: PR title
            body: PR description

        Returns:
            PR URL

        Raises:
            subprocess.CalledProcessError: If gh command fails

        """
        result = subprocess.run(  # noqa: S603
            [  # noqa: S607
                "gh",
                "pr",
                "create",
                "--base",
                base_branch,
                "--head",
                head_branch,
                "--title",
                title,
                "--body",
                body,
            ],
            capture_output=True,
            text=True,
            check=True,
        )
        # gh outputs PR URL
        return result.stdout.strip()

    def close_milestone_if_complete(self, feature_id: str) -> bool:
        """Close milestone if all issues are closed.

        Args:
            feature_id: Feature ID (e.g., "F60")

        Returns:
            True if milestone was closed, False otherwise

        """
        manager = MilestoneManager(self._repo)

        # Find milestone
        milestones = self._repo.get_milestones(state="open")
        milestone = None
        for m in milestones:
            if f"Feature {feature_id}:" in m.title:
                milestone = m
                break

        if not milestone:
            return False

        # Check if all issues closed
        if milestone.open_issues == 0:
            manager.close_milestone(milestone)
            return True

        return False

    @classmethod
    def from_env(cls) -> "DeployGitHubIntegration":
        """Create integration from environment config.

        Returns:
            DeployGitHubIntegration instance

        """
        config = GitHubConfig.from_env()
        client = GitHubClient(config)
        return cls(client)
