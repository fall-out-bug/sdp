"""GitHub API client wrapper.

Provides typed interface to PyGithub for hw_checker operations.
"""

from typing import Any, Optional

from github import Github, Issue, Repository

from sdp.github.config import GitHubConfig


class GitHubClient:
    """Wrapper around PyGithub for hw_checker operations.

    Attributes:
        _config: GitHub configuration
        _github: PyGithub Github instance
        _repo: Cached repository object
    """

    def __init__(self, config: GitHubConfig) -> None:
        """Initialize GitHub client.

        Args:
            config: GitHub configuration

        Raises:
            ValueError: If token is invalid
        """
        self._config = config
        self._github = Github(config.token)
        self._repo: Optional[Repository.Repository] = None

    def get_repo(self) -> Repository.Repository:
        """Get repository object.

        Returns:
            PyGithub Repository instance

        Raises:
            GithubException: If repo not found or access denied
        """
        if self._repo is None:
            self._repo = self._github.get_repo(self._config.repo)
        return self._repo

    def get_issue(self, issue_number: int) -> Issue.Issue:
        """Get issue by number.

        Args:
            issue_number: GitHub issue number

        Returns:
            PyGithub Issue instance

        Raises:
            GithubException: If issue not found
        """
        repo = self.get_repo()
        return repo.get_issue(issue_number)

    def get_issue_global_id(self, issue_number: int) -> str:
        """Get issue global ID (for Projects API).

        Args:
            issue_number: Issue number

        Returns:
            Global ID (node_id)
        """
        issue = self.get_issue(issue_number)
        return issue.node_id

    def create_issue(
        self,
        title: str,
        body: str,
        labels: Optional[list[str]] = None,
        milestone: Optional[int] = None,
    ) -> Issue.Issue:
        """Create new issue.

        Args:
            title: Issue title
            body: Issue body (markdown)
            labels: List of label names
            milestone: Milestone number (optional)

        Returns:
            Created Issue instance

        Raises:
            GithubException: If creation failed
        """
        repo = self.get_repo()
        kwargs: dict[str, Any] = {
            "title": title,
            "body": body,
            "labels": labels or [],
        }
        if milestone:
            kwargs["milestone"] = repo.get_milestone(milestone)
        return repo.create_issue(**kwargs)

    def close(self) -> None:
        """Close GitHub connection."""
        self._github.close()
