"""Tests for GitHub sync service."""

from pathlib import Path
from unittest.mock import Mock, patch

import pytest

from sdp.github.client import GitHubClient
from sdp.github.config import GitHubConfig
from sdp.github.sync_service import SyncService, SyncResult
from sdp.github.ws_parser import WSMetadata


@pytest.fixture
def mock_config() -> GitHubConfig:
    """Create mock GitHub config."""
    return GitHubConfig(token="test_token", repo="owner/repo")


@pytest.fixture
def sample_ws_metadata() -> WSMetadata:
    """Create sample WSMetadata for sync tests."""
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


class TestSyncServiceInit:
    """Test SyncService initialization."""

    def test_init_accepts_client(self, mock_config: GitHubConfig) -> None:
        """Verify SyncService accepts GitHubClient."""
        with patch("sdp.github.client.Github"):
            client = GitHubClient(mock_config)
            service = SyncService(client=client)

        assert service._client is client


class TestSyncServiceSyncWorkstream:
    """Test sync_workstream method."""

    def test_sync_returns_result_on_parse_error(self, mock_config: GitHubConfig) -> None:
        """Verify SyncResult returned when WS file invalid."""
        with patch("sdp.github.client.Github"):
            client = GitHubClient(mock_config)
            service = SyncService(client=client)

        invalid_file = Path("/nonexistent/00-001-01.md")
        result = service.sync_workstream(invalid_file)

        assert isinstance(result, SyncResult)
        assert result.action == "failed"
        assert result.error is not None

    @patch("sdp.github.sync_service.FrontmatterUpdater")
    @patch("sdp.github.sync_service.parse_ws_file")
    def test_sync_create_path_mocked(
        self,
        mock_parse: Mock,
        mock_frontmatter: Mock,
        mock_config: GitHubConfig,
        sample_ws_metadata: WSMetadata,
        tmp_path: Path,
    ) -> None:
        """Verify sync creates issue when no existing issue (mocked)."""
        mock_parse.return_value = sample_ws_metadata
        mock_frontmatter.get_github_issue.return_value = None

        mock_issue = Mock()
        mock_issue.number = 123
        mock_issue.labels = []
        mock_issue.state = "open"

        with patch("sdp.github.client.Github"):
            client = GitHubClient(mock_config)
            client.get_issue = Mock(return_value=mock_issue)

            with patch("sdp.github.sync_helpers.SyncHelpers.resolve_project_name", return_value="SDP"), \
                 patch("sdp.github.sync_helpers.SyncHelpers.get_feature_milestone", return_value=(None, None)), \
                 patch("sdp.github.sync_helpers.SyncHelpers.sync_to_board"), \
                 patch("sdp.github.sync_helpers.SyncHelpers.update_issue_status"), \
                 patch("sdp.github.sync_helpers.SyncHelpers.ensure_issue_milestone"), \
                 patch("sdp.github.sync_helpers.SyncHelpers.get_board_sync", return_value=None):
                service = SyncService(client=client)
                service._issue_sync.sync_ws = Mock(return_value=123)

                ws_file = tmp_path / "00-001-01.md"
                ws_file.write_text("---\nws_id: 00-001-01\n---")

                result = service.sync_workstream(ws_file)

        assert result.action == "created"
        assert result.issue_number == 123
        assert result.ws_id == "00-001-01"

    @patch("sdp.github.sync_service.FrontmatterUpdater")
    @patch("sdp.github.sync_service.parse_ws_file")
    def test_sync_update_path_mocked(
        self,
        mock_parse: Mock,
        mock_frontmatter: Mock,
        mock_config: GitHubConfig,
        sample_ws_metadata: WSMetadata,
        tmp_path: Path,
    ) -> None:
        """Verify sync updates when issue exists (mocked)."""
        mock_parse.return_value = sample_ws_metadata
        mock_frontmatter.get_github_issue.return_value = 123

        mock_issue = Mock()
        mock_issue.labels = []
        mock_issue.state = "open"

        with patch("sdp.github.client.Github"):
            client = GitHubClient(mock_config)
            client.get_issue = Mock(return_value=mock_issue)

            with patch("sdp.github.sync_helpers.SyncHelpers.resolve_project_name", return_value="SDP"), \
                 patch("sdp.github.sync_helpers.SyncHelpers.get_feature_milestone", return_value=(None, None)), \
                 patch("sdp.github.sync_helpers.SyncHelpers.sync_to_board"), \
                 patch("sdp.github.sync_helpers.SyncHelpers.update_issue_status"), \
                 patch("sdp.github.sync_helpers.SyncHelpers.ensure_issue_milestone"), \
                 patch("sdp.github.sync_helpers.SyncHelpers.get_board_sync", return_value=None):
                service = SyncService(client=client)

                ws_file = tmp_path / "00-001-01.md"
                ws_file.write_text("---\nws_id: 00-001-01\n---")

                result = service.sync_workstream(ws_file)

        assert result.action == "updated"
        assert result.issue_number == 123
