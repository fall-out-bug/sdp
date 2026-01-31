"""Tests for GitHub label management."""

from unittest.mock import Mock, patch

import pytest
from github import GithubException

from sdp.github.client import GitHubClient
from sdp.github.config import GitHubConfig
from sdp.github.label_manager import LabelConfig, LabelManager
from sdp.github.ws_parser import WSMetadata


@pytest.fixture
def mock_config() -> GitHubConfig:
    """Create mock GitHub config."""
    return GitHubConfig(token="test_token", repo="owner/repo")


@pytest.fixture
def mock_client(mock_config: GitHubConfig) -> GitHubClient:
    """Create mock GitHub client."""
    with patch("sdp.github.client.Github"):
        return GitHubClient(mock_config)


@pytest.fixture
def label_manager(mock_client: GitHubClient) -> LabelManager:
    """Create LabelManager instance."""
    return LabelManager(mock_client)


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


class TestLabelManagerInit:
    """Test LabelManager initialization."""

    def test_init_stores_client(self, mock_client: GitHubClient) -> None:
        """Verify LabelManager stores client."""
        manager = LabelManager(mock_client)

        assert manager._client is mock_client


class TestLabelManagerDeriveLabels:
    """Test derive_labels method."""

    def test_derive_labels_returns_base_labels(
        self, label_manager: LabelManager, sample_ws_metadata: WSMetadata
    ) -> None:
        """Verify derive_labels returns base labels."""
        labels = label_manager.derive_labels(sample_ws_metadata)

        assert "workstream" in labels
        assert f"feature/{sample_ws_metadata.feature}" in labels
        assert f"size/{sample_ws_metadata.size}" in labels
        assert f"status/{sample_ws_metadata.status}" in labels

    def test_derive_labels_with_different_feature(
        self, label_manager: LabelManager
    ) -> None:
        """Verify derive_labels handles different features."""
        ws = WSMetadata(
            ws_id="00-002-01",
            feature="F150",
            title="Test",
            goal="Goal",
            acceptance_criteria=[],
            dependencies=[],
            size="MEDIUM",
            status="active",
        )

        labels = label_manager.derive_labels(ws)

        assert "feature/F150" in labels
        assert "size/MEDIUM" in labels
        assert "status/active" in labels

    def test_derive_labels_with_large_size(
        self, label_manager: LabelManager
    ) -> None:
        """Verify derive_labels handles LARGE size."""
        ws = WSMetadata(
            ws_id="00-003-01",
            feature="F001",
            title="Test",
            goal="Goal",
            acceptance_criteria=[],
            dependencies=[],
            size="LARGE",
            status="completed",
        )

        labels = label_manager.derive_labels(ws)

        assert "size/LARGE" in labels
        assert "status/completed" in labels


class TestLabelManagerCreateLabel:
    """Test create_label method."""

    def test_create_label_creates_new_label(
        self, label_manager: LabelManager, mock_client: GitHubClient
    ) -> None:
        """Verify create_label creates label when not exists."""
        mock_repo = Mock()
        mock_repo.create_label = Mock()
        mock_client.get_repo = Mock(return_value=mock_repo)

        config = LabelConfig(name="test-label", color="ff0000", description="Test")

        label_manager.create_label(config)

        mock_repo.create_label.assert_called_once_with(
            name="test-label",
            color="ff0000",
            description="Test",
        )

    def test_create_label_ignores_already_exists_error(
        self, label_manager: LabelManager, mock_client: GitHubClient
    ) -> None:
        """Verify create_label ignores 422 (already exists) errors."""
        mock_repo = Mock()
        already_exists_error = GithubException(422, {"message": "already exists"}, headers={})
        mock_repo.create_label = Mock(side_effect=already_exists_error)
        mock_client.get_repo = Mock(return_value=mock_repo)

        config = LabelConfig(name="existing-label", color="00ff00", description="Exists")

        # Should not raise
        label_manager.create_label(config)

        mock_repo.create_label.assert_called_once()

    def test_create_label_raises_other_errors(
        self, label_manager: LabelManager, mock_client: GitHubClient
    ) -> None:
        """Verify create_label raises non-422 errors."""
        mock_repo = Mock()
        forbidden_error = GithubException(403, {"message": "Forbidden"}, headers={})
        mock_repo.create_label = Mock(side_effect=forbidden_error)
        mock_client.get_repo = Mock(return_value=mock_repo)

        config = LabelConfig(name="test-label", color="0000ff", description="Test")

        with pytest.raises(GithubException) as exc_info:
            label_manager.create_label(config)

        assert exc_info.value.status == 403

    def test_create_label_handles_rate_limit(
        self, label_manager: LabelManager, mock_client: GitHubClient
    ) -> None:
        """Verify create_label handles rate limit errors."""
        mock_repo = Mock()
        rate_limit_error = GithubException(
            403,
            {"message": "API rate limit exceeded"},
            headers={"X-RateLimit-Remaining": "0"},
        )
        mock_repo.create_label = Mock(side_effect=rate_limit_error)
        mock_client.get_repo = Mock(return_value=mock_repo)

        config = LabelConfig(name="test-label", color="ffffff", description="Test")

        with pytest.raises(GithubException) as exc_info:
            label_manager.create_label(config)

        assert exc_info.value.status == 403

    def test_create_label_handles_network_error(
        self, label_manager: LabelManager, mock_client: GitHubClient
    ) -> None:
        """Verify create_label handles network errors."""
        mock_repo = Mock()
        mock_repo.create_label = Mock(side_effect=TimeoutError("Connection timeout"))
        mock_client.get_repo = Mock(return_value=mock_repo)

        config = LabelConfig(name="test-label", color="aaaaaa", description="Test")

        with pytest.raises(TimeoutError):
            label_manager.create_label(config)


class TestLabelManagerEnsureLabels:
    """Test ensure_labels method."""

    def test_ensure_labels_creates_all_labels(
        self, label_manager: LabelManager, sample_ws_metadata: WSMetadata, mock_client: GitHubClient
    ) -> None:
        """Verify ensure_labels creates all derived labels."""
        mock_repo = Mock()
        mock_repo.create_label = Mock()
        mock_client.get_repo = Mock(return_value=mock_repo)

        label_manager.ensure_labels(sample_ws_metadata)

        # Should create 4 labels: workstream, feature/F001, size/SMALL, status/backlog
        assert mock_repo.create_label.call_count == 4

    def test_ensure_labels_handles_existing_labels(
        self, label_manager: LabelManager, sample_ws_metadata: WSMetadata, mock_client: GitHubClient
    ) -> None:
        """Verify ensure_labels handles already existing labels."""
        mock_repo = Mock()
        already_exists_error = GithubException(422, {"message": "already exists"}, headers={})
        mock_repo.create_label = Mock(side_effect=already_exists_error)
        mock_client.get_repo = Mock(return_value=mock_repo)

        # Should not raise
        label_manager.ensure_labels(sample_ws_metadata)

        assert mock_repo.create_label.call_count == 4

    def test_ensure_labels_handles_partial_failure(
        self, label_manager: LabelManager, sample_ws_metadata: WSMetadata, mock_client: GitHubClient
    ) -> None:
        """Verify ensure_labels raises on non-422 errors."""
        mock_repo = Mock()
        call_count = 0

        def side_effect(*args: object, **kwargs: object) -> None:
            nonlocal call_count
            call_count += 1
            if call_count == 2:  # Fail on second label
                raise GithubException(403, {"message": "Forbidden"}, headers={})
            raise GithubException(422, {"message": "already exists"}, headers={})

        mock_repo.create_label = Mock(side_effect=side_effect)
        mock_client.get_repo = Mock(return_value=mock_repo)

        with pytest.raises(GithubException) as exc_info:
            label_manager.ensure_labels(sample_ws_metadata)

        assert exc_info.value.status == 403


class TestLabelManagerConfigForLabel:
    """Test _config_for_label method."""

    def test_config_for_label_returns_exact_match(
        self, label_manager: LabelManager
    ) -> None:
        """Verify _config_for_label returns config for exact match."""
        config = label_manager._config_for_label("workstream")

        assert config is not None
        assert config.name == "workstream"
        assert config.color == "0366d6"
        assert config.description == "Workstream task"

    def test_config_for_label_returns_feature_pattern(
        self, label_manager: LabelManager
    ) -> None:
        """Verify _config_for_label returns config for feature pattern."""
        config = label_manager._config_for_label("feature/F200")

        assert config is not None
        assert config.name == "feature/F200"
        assert config.color == "7057ff"

    def test_config_for_label_returns_size_pattern(
        self, label_manager: LabelManager
    ) -> None:
        """Verify _config_for_label returns config for size pattern."""
        config = label_manager._config_for_label("size/EXTRA_LARGE")

        assert config is not None
        assert config.name == "size/EXTRA_LARGE"
        assert config.color == "d4c5f9"

    def test_config_for_label_returns_status_pattern(
        self, label_manager: LabelManager
    ) -> None:
        """Verify _config_for_label returns config for status pattern."""
        config = label_manager._config_for_label("status/blocked")

        assert config is not None
        assert config.name == "status/blocked"
        assert config.color == "cccccc"

    def test_config_for_label_returns_none_for_unknown(
        self, label_manager: LabelManager
    ) -> None:
        """Verify _config_for_label returns None for unknown label."""
        config = label_manager._config_for_label("unknown-label")

        assert config is None

    def test_config_for_label_handles_predefined_features(
        self, label_manager: LabelManager
    ) -> None:
        """Verify _config_for_label handles predefined feature labels."""
        config = label_manager._config_for_label("feature/F100")

        assert config is not None
        assert config.name == "feature/F100"
        assert config.description == "Feature F100"

    def test_config_for_label_handles_predefined_sizes(
        self, label_manager: LabelManager
    ) -> None:
        """Verify _config_for_label handles predefined size labels."""
        for size in ["SMALL", "MEDIUM", "LARGE"]:
            config = label_manager._config_for_label(f"size/{size}")

            assert config is not None
            assert config.name == f"size/{size}"

    def test_config_for_label_handles_predefined_statuses(
        self, label_manager: LabelManager
    ) -> None:
        """Verify _config_for_label handles predefined status labels."""
        for status in ["backlog", "active", "completed"]:
            config = label_manager._config_for_label(f"status/{status}")

            assert config is not None
            assert config.name == f"status/{status}"


class TestLabelManagerGetLabelConfigs:
    """Test _get_label_configs method."""

    def test_get_label_configs_returns_configs(
        self, label_manager: LabelManager
    ) -> None:
        """Verify _get_label_configs returns configs for known labels."""
        label_names = ["workstream", "feature/F001", "size/SMALL", "status/backlog"]

        configs = label_manager._get_label_configs(label_names)

        assert len(configs) == 4
        assert all(config is not None for config in configs)

    def test_get_label_configs_filters_unknown_labels(
        self, label_manager: LabelManager
    ) -> None:
        """Verify _get_label_configs filters out unknown labels."""
        label_names = ["workstream", "unknown-label", "size/SMALL"]

        configs = label_manager._get_label_configs(label_names)

        assert len(configs) == 2
        assert all(config.name != "unknown-label" for config in configs)

    def test_get_label_configs_handles_empty_list(
        self, label_manager: LabelManager
    ) -> None:
        """Verify _get_label_configs handles empty list."""
        configs = label_manager._get_label_configs([])

        assert len(configs) == 0


class TestLabelManagerErrorHandling:
    """Test error handling."""

    def test_create_label_handles_authentication_error(
        self, label_manager: LabelManager, mock_client: GitHubClient
    ) -> None:
        """Verify create_label handles authentication errors."""
        mock_repo = Mock()
        auth_error = GithubException(401, {"message": "Bad credentials"}, headers={})
        mock_repo.create_label = Mock(side_effect=auth_error)
        mock_client.get_repo = Mock(return_value=mock_repo)

        config = LabelConfig(name="test-label", color="000000", description="Test")

        with pytest.raises(GithubException) as exc_info:
            label_manager.create_label(config)

        assert exc_info.value.status == 401

    def test_create_label_handles_not_found_error(
        self, label_manager: LabelManager, mock_client: GitHubClient
    ) -> None:
        """Verify create_label handles not found errors."""
        mock_repo = Mock()
        not_found_error = GithubException(404, {"message": "Not Found"}, headers={})
        mock_repo.create_label = Mock(side_effect=not_found_error)
        mock_client.get_repo = Mock(return_value=mock_repo)

        config = LabelConfig(name="test-label", color="111111", description="Test")

        with pytest.raises(GithubException) as exc_info:
            label_manager.create_label(config)

        assert exc_info.value.status == 404

    def test_ensure_labels_handles_repo_access_error(
        self, label_manager: LabelManager, sample_ws_metadata: WSMetadata, mock_client: GitHubClient
    ) -> None:
        """Verify ensure_labels handles repo access errors."""
        access_error = GithubException(403, {"message": "Forbidden"}, headers={})
        mock_client.get_repo = Mock(side_effect=access_error)

        with pytest.raises(GithubException) as exc_info:
            label_manager.ensure_labels(sample_ws_metadata)

        assert exc_info.value.status == 403


class TestLabelManagerEdgeCases:
    """Test edge cases."""

    def test_derive_labels_with_empty_metadata(
        self, label_manager: LabelManager
    ) -> None:
        """Verify derive_labels handles minimal metadata."""
        ws = WSMetadata(
            ws_id="00-999-99",
            feature="F999",
            title="",
            goal="",
            acceptance_criteria=[],
            dependencies=[],
            size="SMALL",
            status="backlog",
        )

        labels = label_manager.derive_labels(ws)

        assert len(labels) == 4
        assert "feature/F999" in labels

    def test_create_label_with_special_characters(
        self, label_manager: LabelManager, mock_client: GitHubClient
    ) -> None:
        """Verify create_label handles special characters in name."""
        mock_repo = Mock()
        mock_repo.create_label = Mock()
        mock_client.get_repo = Mock(return_value=mock_repo)

        config = LabelConfig(
            name="label-with-special-chars-&-symbols",
            color="abcdef",
            description="Test & description",
        )

        label_manager.create_label(config)

        mock_repo.create_label.assert_called_once()
        call_kwargs = mock_repo.create_label.call_args[1]
        assert call_kwargs["name"] == "label-with-special-chars-&-symbols"

    def test_ensure_labels_with_custom_feature(
        self, label_manager: LabelManager, mock_client: GitHubClient
    ) -> None:
        """Verify ensure_labels handles custom feature labels."""
        ws = WSMetadata(
            ws_id="00-001-01",
            feature="CUSTOM_FEATURE",
            title="Test",
            goal="Goal",
            acceptance_criteria=[],
            dependencies=[],
            size="SMALL",
            status="backlog",
        )

        mock_repo = Mock()
        mock_repo.create_label = Mock()
        mock_client.get_repo = Mock(return_value=mock_repo)

        label_manager.ensure_labels(ws)

        # Should create label for custom feature
        call_args_list = [call[1]["name"] for call in mock_repo.create_label.call_args_list]
        assert "feature/CUSTOM_FEATURE" in call_args_list
