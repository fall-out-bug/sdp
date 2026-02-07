"""Tests for GitHub issue synchronization."""

from unittest.mock import Mock, patch

import pytest
from github import GithubException

from sdp.github.client import GitHubClient
from sdp.github.config import GitHubConfig
from sdp.github.issue_sync import IssueSync
from sdp.github.label_manager import LabelManager
from sdp.github.ws_parser import WSMetadata


@pytest.fixture
def mock_config() -> GitHubConfig:
    """Create mock GitHub config."""
    return GitHubConfig(token="test_token", repo="owner/repo")


@pytest.fixture
def sample_ws_metadata() -> WSMetadata:
    """Create sample WSMetadata for tests."""
    return WSMetadata(
        ws_id="00-001-01",
        feature="F001",
        title="Test Workstream",
        goal="Test goal",
        acceptance_criteria=["AC1: Test"],
        dependencies=[],
        size="SMALL",
        status="backlog",
    )


@pytest.fixture
def mock_client(mock_config: GitHubConfig) -> GitHubClient:
    """Create mock GitHub client."""
    with patch("sdp.github.client.Github"):
        return GitHubClient(mock_config)


@pytest.fixture
def mock_label_manager(mock_client: GitHubClient) -> LabelManager:
    """Create mock label manager."""
    return LabelManager(mock_client)


@pytest.fixture
def issue_sync(mock_client: GitHubClient, mock_label_manager: LabelManager) -> IssueSync:
    """Create IssueSync instance."""
    return IssueSync(mock_client, mock_label_manager)


class TestIssueSyncInit:
    """Test IssueSync initialization."""

    def test_init_stores_client_and_label_manager(
        self, mock_client: GitHubClient, mock_label_manager: LabelManager
    ) -> None:
        """Verify IssueSync stores client and label manager."""
        sync = IssueSync(mock_client, mock_label_manager)

        assert sync._client is mock_client
        assert sync._label_manager is mock_label_manager


class TestIssueSyncSyncWS:
    """Test sync_ws method."""

    def test_sync_ws_creates_issue_successfully(
        self,
        issue_sync: IssueSync,
        sample_ws_metadata: WSMetadata,
        mock_client: GitHubClient,
    ) -> None:
        """Verify sync_ws creates issue with correct parameters."""
        mock_issue = Mock()
        mock_issue.number = 123
        mock_client.create_issue = Mock(return_value=mock_issue)

        with patch.object(issue_sync._label_manager, "ensure_labels"), \
             patch.object(issue_sync._label_manager, "derive_labels", return_value=["workstream"]), \
             patch("sdp.github.issue_sync.format_issue_body", return_value="Formatted body"):
            result = issue_sync.sync_ws(sample_ws_metadata, "docs/workstreams/backlog/00-001-01.md")

        assert result == 123
        mock_client.create_issue.assert_called_once()
        call_kwargs = mock_client.create_issue.call_args[1]
        assert call_kwargs["title"] == "00-001-01: Test Workstream (F001)"
        assert call_kwargs["body"] == "Formatted body"
        assert call_kwargs["labels"] == ["workstream"]

    def test_sync_ws_with_milestone(
        self,
        issue_sync: IssueSync,
        sample_ws_metadata: WSMetadata,
        mock_client: GitHubClient,
    ) -> None:
        """Verify sync_ws includes milestone when provided."""
        mock_issue = Mock()
        mock_issue.number = 456
        mock_client.create_issue = Mock(return_value=mock_issue)

        with patch.object(issue_sync._label_manager, "ensure_labels"), \
             patch.object(issue_sync._label_manager, "derive_labels", return_value=["workstream"]), \
             patch("sdp.github.issue_sync.format_issue_body", return_value="Body"):
            result = issue_sync.sync_ws(
                sample_ws_metadata,
                "docs/workstreams/backlog/00-001-01.md",
                milestone_number=5,
            )

        assert result == 456
        call_kwargs = mock_client.create_issue.call_args[1]
        assert call_kwargs["milestone"] == 5

    def test_sync_ws_ensures_labels_before_creating(
        self,
        issue_sync: IssueSync,
        sample_ws_metadata: WSMetadata,
        mock_client: GitHubClient,
    ) -> None:
        """Verify ensure_labels is called before issue creation."""
        mock_issue = Mock()
        mock_issue.number = 789
        mock_client.create_issue = Mock(return_value=mock_issue)

        ensure_labels_mock = Mock()
        issue_sync._label_manager.ensure_labels = ensure_labels_mock

        with patch.object(issue_sync._label_manager, "derive_labels", return_value=["workstream"]), \
             patch("sdp.github.issue_sync.format_issue_body", return_value="Body"):
            issue_sync.sync_ws(sample_ws_metadata, "docs/workstreams/backlog/00-001-01.md")

        ensure_labels_mock.assert_called_once_with(sample_ws_metadata)

    def test_sync_ws_uses_derived_labels(
        self,
        issue_sync: IssueSync,
        sample_ws_metadata: WSMetadata,
        mock_client: GitHubClient,
    ) -> None:
        """Verify derived labels are used in issue creation."""
        mock_issue = Mock()
        mock_issue.number = 999
        mock_client.create_issue = Mock(return_value=mock_issue)

        derived_labels = ["workstream", "feature/F001", "size/SMALL", "status/backlog"]

        with patch.object(issue_sync._label_manager, "ensure_labels"), \
             patch.object(issue_sync._label_manager, "derive_labels", return_value=derived_labels), \
             patch("sdp.github.issue_sync.format_issue_body", return_value="Body"):
            issue_sync.sync_ws(sample_ws_metadata, "docs/workstreams/backlog/00-001-01.md")

        call_kwargs = mock_client.create_issue.call_args[1]
        assert call_kwargs["labels"] == derived_labels


class TestIssueSyncErrorHandling:
    """Test error handling in sync_ws."""

    def test_sync_ws_raises_on_create_issue_failure(
        self,
        issue_sync: IssueSync,
        sample_ws_metadata: WSMetadata,
        mock_client: GitHubClient,
    ) -> None:
        """Verify sync_ws raises exception when create_issue fails."""
        github_error = GithubException(500, {"message": "Internal server error"}, headers={})
        mock_client.create_issue = Mock(side_effect=github_error)

        with patch.object(issue_sync._label_manager, "ensure_labels"), \
             patch.object(issue_sync._label_manager, "derive_labels", return_value=["workstream"]), \
             patch("sdp.github.issue_sync.format_issue_body", return_value="Body"):
            with pytest.raises(GithubException) as exc_info:
                issue_sync.sync_ws(sample_ws_metadata, "docs/workstreams/backlog/00-001-01.md")

            assert exc_info.value.status == 500

    def test_sync_ws_raises_on_label_ensure_failure(
        self,
        issue_sync: IssueSync,
        sample_ws_metadata: WSMetadata,
    ) -> None:
        """Verify sync_ws raises exception when ensure_labels fails."""
        github_error = GithubException(403, {"message": "Forbidden"}, headers={})
        issue_sync._label_manager.ensure_labels = Mock(side_effect=github_error)

        with pytest.raises(GithubException) as exc_info:
            issue_sync.sync_ws(sample_ws_metadata, "docs/workstreams/backlog/00-001-01.md")

        assert exc_info.value.status == 403

    def test_sync_ws_handles_network_timeout(
        self,
        issue_sync: IssueSync,
        sample_ws_metadata: WSMetadata,
        mock_client: GitHubClient,
    ) -> None:
        """Verify sync_ws handles network timeout errors."""
        timeout_error = TimeoutError("Connection timeout")
        mock_client.create_issue = Mock(side_effect=timeout_error)

        with patch.object(issue_sync._label_manager, "ensure_labels"), \
             patch.object(issue_sync._label_manager, "derive_labels", return_value=["workstream"]), \
             patch("sdp.github.issue_sync.format_issue_body", return_value="Body"):
            with pytest.raises(TimeoutError):
                issue_sync.sync_ws(sample_ws_metadata, "docs/workstreams/backlog/00-001-01.md")

    def test_sync_ws_handles_rate_limit_error(
        self,
        issue_sync: IssueSync,
        sample_ws_metadata: WSMetadata,
        mock_client: GitHubClient,
    ) -> None:
        """Verify sync_ws handles rate limit errors."""
        rate_limit_error = GithubException(
            403,
            {"message": "API rate limit exceeded"},
            headers={"X-RateLimit-Remaining": "0"},
        )
        mock_client.create_issue = Mock(side_effect=rate_limit_error)

        with patch.object(issue_sync._label_manager, "ensure_labels"), \
             patch.object(issue_sync._label_manager, "derive_labels", return_value=["workstream"]), \
             patch("sdp.github.issue_sync.format_issue_body", return_value="Body"):
            with pytest.raises(GithubException) as exc_info:
                issue_sync.sync_ws(sample_ws_metadata, "docs/workstreams/backlog/00-001-01.md")

            assert exc_info.value.status == 403

    def test_sync_ws_handles_authentication_error(
        self,
        issue_sync: IssueSync,
        sample_ws_metadata: WSMetadata,
        mock_client: GitHubClient,
    ) -> None:
        """Verify sync_ws handles authentication errors."""
        auth_error = GithubException(401, {"message": "Bad credentials"}, headers={})
        mock_client.create_issue = Mock(side_effect=auth_error)

        with patch.object(issue_sync._label_manager, "ensure_labels"), \
             patch.object(issue_sync._label_manager, "derive_labels", return_value=["workstream"]), \
             patch("sdp.github.issue_sync.format_issue_body", return_value="Body"):
            with pytest.raises(GithubException) as exc_info:
                issue_sync.sync_ws(sample_ws_metadata, "docs/workstreams/backlog/00-001-01.md")

            assert exc_info.value.status == 401

    def test_sync_ws_handles_not_found_error(
        self,
        issue_sync: IssueSync,
        sample_ws_metadata: WSMetadata,
        mock_client: GitHubClient,
    ) -> None:
        """Verify sync_ws handles not found errors."""
        not_found_error = GithubException(404, {"message": "Not Found"}, headers={})
        mock_client.create_issue = Mock(side_effect=not_found_error)

        with patch.object(issue_sync._label_manager, "ensure_labels"), \
             patch.object(issue_sync._label_manager, "derive_labels", return_value=["workstream"]), \
             patch("sdp.github.issue_sync.format_issue_body", return_value="Body"):
            with pytest.raises(GithubException) as exc_info:
                issue_sync.sync_ws(sample_ws_metadata, "docs/workstreams/backlog/00-001-01.md")

            assert exc_info.value.status == 404


class TestIssueSyncEdgeCases:
    """Test edge cases in sync_ws."""

    def test_sync_ws_with_empty_labels(
        self,
        issue_sync: IssueSync,
        sample_ws_metadata: WSMetadata,
        mock_client: GitHubClient,
    ) -> None:
        """Verify sync_ws handles empty label list."""
        mock_issue = Mock()
        mock_issue.number = 111
        mock_client.create_issue = Mock(return_value=mock_issue)

        with patch.object(issue_sync._label_manager, "ensure_labels"), \
             patch.object(issue_sync._label_manager, "derive_labels", return_value=[]), \
             patch("sdp.github.issue_sync.format_issue_body", return_value="Body"):
            result = issue_sync.sync_ws(sample_ws_metadata, "docs/workstreams/backlog/00-001-01.md")

        assert result == 111
        call_kwargs = mock_client.create_issue.call_args[1]
        assert call_kwargs["labels"] == []

    def test_sync_ws_with_long_title(
        self,
        issue_sync: IssueSync,
        mock_client: GitHubClient,
    ) -> None:
        """Verify sync_ws handles long titles correctly."""
        long_title = "A" * 200
        ws_metadata = WSMetadata(
            ws_id="00-001-01",
            feature="F001",
            title=long_title,
            goal="Test goal",
            acceptance_criteria=[],
            dependencies=[],
            size="SMALL",
            status="backlog",
        )

        mock_issue = Mock()
        mock_issue.number = 222
        mock_client.create_issue = Mock(return_value=mock_issue)

        with patch.object(issue_sync._label_manager, "ensure_labels"), \
             patch.object(issue_sync._label_manager, "derive_labels", return_value=["workstream"]), \
             patch("sdp.github.issue_sync.format_issue_body", return_value="Body"):
            result = issue_sync.sync_ws(ws_metadata, "docs/workstreams/backlog/00-001-01.md")

        assert result == 222
        call_kwargs = mock_client.create_issue.call_args[1]
        expected_title = f"00-001-01: {long_title} (F001)"
        assert call_kwargs["title"] == expected_title

    def test_sync_ws_with_special_characters_in_title(
        self,
        issue_sync: IssueSync,
        mock_client: GitHubClient,
    ) -> None:
        """Verify sync_ws handles special characters in title."""
        special_title = "Test & < > \" ' Title"
        ws_metadata = WSMetadata(
            ws_id="00-001-01",
            feature="F001",
            title=special_title,
            goal="Test goal",
            acceptance_criteria=[],
            dependencies=[],
            size="SMALL",
            status="backlog",
        )

        mock_issue = Mock()
        mock_issue.number = 333
        mock_client.create_issue = Mock(return_value=mock_issue)

        with patch.object(issue_sync._label_manager, "ensure_labels"), \
             patch.object(issue_sync._label_manager, "derive_labels", return_value=["workstream"]), \
             patch("sdp.github.issue_sync.format_issue_body", return_value="Body"):
            result = issue_sync.sync_ws(ws_metadata, "docs/workstreams/backlog/00-001-01.md")

        assert result == 333
        call_kwargs = mock_client.create_issue.call_args[1]
        assert special_title in call_kwargs["title"]
