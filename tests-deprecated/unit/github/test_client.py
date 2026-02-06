"""Tests for GitHub API client wrapper."""

from unittest.mock import Mock, patch

import pytest

from sdp.github.client import GitHubClient
from sdp.github.config import GitHubConfig
from sdp.github.exceptions import (
    GitHubSyncError,
    RateLimitError,
    AuthenticationError,
)


@pytest.fixture
def mock_config() -> GitHubConfig:
    """Create mock GitHub config."""
    return GitHubConfig(token="test_token", repo="owner/repo")


class TestGitHubClientInit:
    """Test GitHub client initialization."""

    @patch("sdp.github.client.Github")
    def test_init_creates_github_instance(self, mock_github: Mock, mock_config: GitHubConfig) -> None:
        """Verify client creates PyGithub instance with token."""
        GitHubClient(mock_config)

        mock_github.assert_called_once_with("test_token")

    @patch("sdp.github.client.Github")
    def test_get_repo_caches_repo(self, mock_github: Mock, mock_config: GitHubConfig) -> None:
        """Verify get_repo returns cached repository."""
        mock_repo = Mock()
        mock_github.return_value.get_repo.return_value = mock_repo

        client = GitHubClient(mock_config)
        repo1 = client.get_repo()
        repo2 = client.get_repo()

        assert repo1 is repo2
        mock_github.return_value.get_repo.assert_called_once_with("owner/repo")


class TestGitHubClientCreateIssue:
    """Test GitHub client create_issue."""

    @patch("sdp.github.client.Github")
    def test_create_issue_returns_issue(self, mock_github: Mock, mock_config: GitHubConfig) -> None:
        """Verify create_issue returns Issue object."""
        mock_issue = Mock()
        mock_issue.number = 123
        mock_repo = Mock()
        mock_repo.create_issue.return_value = mock_issue
        mock_github.return_value.get_repo.return_value = mock_repo

        client = GitHubClient(mock_config)
        result = client.create_issue(
            title="Test",
            body="Body",
            labels=["workstream"],
        )

        assert result is mock_issue
        mock_repo.create_issue.assert_called_once_with(
            title="Test",
            body="Body",
            labels=["workstream"],
        )


class TestGitHubClientGetIssue:
    """Test GitHub client get_issue."""

    @patch("sdp.github.client.Github")
    def test_get_issue_returns_issue(self, mock_github: Mock, mock_config: GitHubConfig) -> None:
        """Verify get_issue returns Issue by number."""
        mock_issue = Mock()
        mock_repo = Mock()
        mock_repo.get_issue.return_value = mock_issue
        mock_github.return_value.get_repo.return_value = mock_repo

        client = GitHubClient(mock_config)
        result = client.get_issue(42)

        assert result is mock_issue
        mock_repo.get_issue.assert_called_once_with(42)
