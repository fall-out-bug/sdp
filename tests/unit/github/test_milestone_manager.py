"""Tests for GitHub milestone management."""

from unittest.mock import Mock

import pytest
from github import GithubException, Milestone, Repository

from sdp.github.milestone_manager import MilestoneManager


@pytest.fixture
def mock_repo() -> Mock:
    """Create mock GitHub repository."""
    return Mock(spec=Repository.Repository)


@pytest.fixture
def milestone_manager(mock_repo: Mock) -> MilestoneManager:
    """Create MilestoneManager instance."""
    return MilestoneManager(mock_repo)


@pytest.fixture
def mock_milestone() -> Mock:
    """Create mock milestone."""
    milestone = Mock(spec=Milestone.Milestone)
    milestone.number = 1
    milestone.title = "Feature F001: Test Feature"
    milestone.description = "Test description"
    milestone.state = "open"
    milestone.open_issues = 5
    milestone.closed_issues = 3
    return milestone


class TestMilestoneManagerInit:
    """Test MilestoneManager initialization."""

    def test_init_stores_repo(self, mock_repo: Mock) -> None:
        """Verify MilestoneManager stores repository."""
        manager = MilestoneManager(mock_repo)

        assert manager._repo is mock_repo


class TestMilestoneManagerGetOrCreateMilestone:
    """Test get_or_create_milestone method."""

    def test_get_or_create_milestone_returns_existing(
        self, milestone_manager: MilestoneManager, mock_repo: Mock, mock_milestone: Mock
    ) -> None:
        """Verify get_or_create_milestone returns existing milestone."""
        mock_repo.get_milestones.return_value = [mock_milestone]

        result = milestone_manager.get_or_create_milestone(
            feature_id="F001",
            title="Test Feature",
        )

        assert result is mock_milestone
        mock_repo.create_milestone.assert_not_called()

    def test_get_or_create_milestone_creates_new(
        self, milestone_manager: MilestoneManager, mock_repo: Mock, mock_milestone: Mock
    ) -> None:
        """Verify get_or_create_milestone creates new milestone when not exists."""
        mock_repo.get_milestones.return_value = []
        mock_repo.create_milestone.return_value = mock_milestone

        result = milestone_manager.get_or_create_milestone(
            feature_id="F002",
            title="New Feature",
            description="New description",
        )

        assert result is mock_milestone
        mock_repo.create_milestone.assert_called_once()
        call_kwargs = mock_repo.create_milestone.call_args[1]
        assert call_kwargs["title"] == "Feature F002: New Feature"
        assert call_kwargs["description"] == "New description"

    def test_get_or_create_milestone_with_due_date(
        self, milestone_manager: MilestoneManager, mock_repo: Mock, mock_milestone: Mock
    ) -> None:
        """Verify get_or_create_milestone includes due date when provided."""
        mock_repo.get_milestones.return_value = []
        mock_repo.create_milestone.return_value = mock_milestone

        result = milestone_manager.get_or_create_milestone(
            feature_id="F003",
            title="Feature with Due Date",
            due_date="2024-12-31T00:00:00Z",
        )

        assert result is mock_milestone
        call_kwargs = mock_repo.create_milestone.call_args[1]
        assert call_kwargs["due_on"] == "2024-12-31T00:00:00Z"

    def test_get_or_create_milestone_handles_empty_description(
        self, milestone_manager: MilestoneManager, mock_repo: Mock, mock_milestone: Mock
    ) -> None:
        """Verify get_or_create_milestone handles empty description."""
        mock_repo.get_milestones.return_value = []
        mock_repo.create_milestone.return_value = mock_milestone

        result = milestone_manager.get_or_create_milestone(
            feature_id="F004",
            title="Feature",
            description="",
        )

        assert result is mock_milestone
        call_kwargs = mock_repo.create_milestone.call_args[1]
        assert call_kwargs["description"] == ""


class TestMilestoneManagerFindMilestoneByTitle:
    """Test _find_milestone_by_title method."""

    def test_find_milestone_by_title_returns_match(
        self, milestone_manager: MilestoneManager, mock_repo: Mock, mock_milestone: Mock
    ) -> None:
        """Verify _find_milestone_by_title returns matching milestone."""
        other_milestone = Mock(spec=Milestone.Milestone)
        other_milestone.title = "Feature F002: Other Feature"
        mock_repo.get_milestones.return_value = [mock_milestone, other_milestone]

        result = milestone_manager._find_milestone_by_title("Feature F001: Test Feature")

        assert result is mock_milestone

    def test_find_milestone_by_title_returns_none_when_not_found(
        self, milestone_manager: MilestoneManager, mock_repo: Mock
    ) -> None:
        """Verify _find_milestone_by_title returns None when not found."""
        mock_milestone = Mock(spec=Milestone.Milestone)
        mock_milestone.title = "Feature F001: Test Feature"
        mock_repo.get_milestones.return_value = [mock_milestone]

        result = milestone_manager._find_milestone_by_title("Feature F999: Not Found")

        assert result is None

    def test_find_milestone_by_title_handles_empty_list(
        self, milestone_manager: MilestoneManager, mock_repo: Mock
    ) -> None:
        """Verify _find_milestone_by_title handles empty milestone list."""
        mock_repo.get_milestones.return_value = []

        result = milestone_manager._find_milestone_by_title("Feature F001: Test")

        assert result is None


class TestMilestoneManagerFindMilestoneForFeature:
    """Test find_milestone_for_feature method."""

    def test_find_milestone_for_feature_returns_match(
        self, milestone_manager: MilestoneManager, mock_repo: Mock, mock_milestone: Mock
    ) -> None:
        """Verify find_milestone_for_feature returns matching milestone."""
        mock_repo.get_milestones.return_value = [mock_milestone]

        result = milestone_manager.find_milestone_for_feature("F001")

        assert result is mock_milestone

    def test_find_milestone_for_feature_handles_prefix_match(
        self, milestone_manager: MilestoneManager, mock_repo: Mock
    ) -> None:
        """Verify find_milestone_for_feature handles prefix matching."""
        milestone = Mock(spec=Milestone.Milestone)
        milestone.title = "Feature F001: Some Title"
        mock_repo.get_milestones.return_value = [milestone]

        result = milestone_manager.find_milestone_for_feature("F001")

        assert result is milestone

    def test_find_milestone_for_feature_handles_exact_match(
        self, milestone_manager: MilestoneManager, mock_repo: Mock
    ) -> None:
        """Verify find_milestone_for_feature handles exact match."""
        milestone = Mock(spec=Milestone.Milestone)
        milestone.title = "Feature F001"
        mock_repo.get_milestones.return_value = [milestone]

        result = milestone_manager.find_milestone_for_feature("F001")

        assert result is milestone

    def test_find_milestone_for_feature_returns_none_when_not_found(
        self, milestone_manager: MilestoneManager, mock_repo: Mock
    ) -> None:
        """Verify find_milestone_for_feature returns None when not found."""
        milestone = Mock(spec=Milestone.Milestone)
        milestone.title = "Feature F001: Test"
        mock_repo.get_milestones.return_value = [milestone]

        result = milestone_manager.find_milestone_for_feature("F999")

        assert result is None


class TestMilestoneManagerGetOrCreateFeatureMilestone:
    """Test get_or_create_feature_milestone method."""

    def test_get_or_create_feature_milestone_returns_existing(
        self, milestone_manager: MilestoneManager, mock_repo: Mock, mock_milestone: Mock
    ) -> None:
        """Verify get_or_create_feature_milestone returns existing milestone."""
        mock_repo.get_milestones.return_value = [mock_milestone]

        result = milestone_manager.get_or_create_feature_milestone("F001")

        assert result is mock_milestone
        mock_repo.create_milestone.assert_not_called()

    def test_get_or_create_feature_milestone_creates_without_title(
        self, milestone_manager: MilestoneManager, mock_repo: Mock, mock_milestone: Mock
    ) -> None:
        """Verify get_or_create_feature_milestone creates milestone without title."""
        mock_repo.get_milestones.return_value = []
        mock_repo.create_milestone.return_value = mock_milestone

        result = milestone_manager.get_or_create_feature_milestone("F005")

        assert result is mock_milestone
        call_kwargs = mock_repo.create_milestone.call_args[1]
        assert call_kwargs["title"] == "Feature F005"

    def test_get_or_create_feature_milestone_creates_with_title(
        self, milestone_manager: MilestoneManager, mock_repo: Mock, mock_milestone: Mock
    ) -> None:
        """Verify get_or_create_feature_milestone creates milestone with title."""
        mock_repo.get_milestones.return_value = []
        mock_repo.create_milestone.return_value = mock_milestone

        result = milestone_manager.get_or_create_feature_milestone(
            "F006",
            title="Custom Title",
        )

        assert result is mock_milestone
        call_kwargs = mock_repo.create_milestone.call_args[1]
        assert call_kwargs["title"] == "Feature F006: Custom Title"

    def test_get_or_create_feature_milestone_with_due_date(
        self, milestone_manager: MilestoneManager, mock_repo: Mock, mock_milestone: Mock
    ) -> None:
        """Verify get_or_create_feature_milestone includes due date."""
        mock_repo.get_milestones.return_value = []
        mock_repo.create_milestone.return_value = mock_milestone

        result = milestone_manager.get_or_create_feature_milestone(
            "F007",
            due_date="2025-01-01T00:00:00Z",
        )

        assert result is mock_milestone
        call_kwargs = mock_repo.create_milestone.call_args[1]
        assert call_kwargs["due_on"] == "2025-01-01T00:00:00Z"


class TestMilestoneManagerUpdateMilestone:
    """Test update_milestone method."""

    def test_update_milestone_updates_title(
        self, milestone_manager: MilestoneManager, mock_milestone: Mock
    ) -> None:
        """Verify update_milestone updates title."""
        milestone_manager.update_milestone(mock_milestone, title="New Title")

        mock_milestone.edit.assert_called_once()
        call_kwargs = mock_milestone.edit.call_args[1]
        assert call_kwargs["title"] == "New Title"
        assert call_kwargs["state"] == mock_milestone.state

    def test_update_milestone_updates_description(
        self, milestone_manager: MilestoneManager, mock_milestone: Mock
    ) -> None:
        """Verify update_milestone updates description."""
        milestone_manager.update_milestone(mock_milestone, description="New Description")

        call_kwargs = mock_milestone.edit.call_args[1]
        assert call_kwargs["description"] == "New Description"

    def test_update_milestone_updates_state(
        self, milestone_manager: MilestoneManager, mock_milestone: Mock
    ) -> None:
        """Verify update_milestone updates state."""
        milestone_manager.update_milestone(mock_milestone, state="closed")

        call_kwargs = mock_milestone.edit.call_args[1]
        assert call_kwargs["state"] == "closed"

    def test_update_milestone_updates_multiple_fields(
        self, milestone_manager: MilestoneManager, mock_milestone: Mock
    ) -> None:
        """Verify update_milestone updates multiple fields."""
        milestone_manager.update_milestone(
            mock_milestone,
            title="Updated Title",
            description="Updated Description",
            state="closed",
        )

        call_kwargs = mock_milestone.edit.call_args[1]
        assert call_kwargs["title"] == "Updated Title"
        assert call_kwargs["description"] == "Updated Description"
        assert call_kwargs["state"] == "closed"

    def test_update_milestone_preserves_state_when_not_specified(
        self, milestone_manager: MilestoneManager, mock_milestone: Mock
    ) -> None:
        """Verify update_milestone preserves state when not specified."""
        milestone_manager.update_milestone(mock_milestone, title="New Title")

        call_kwargs = mock_milestone.edit.call_args[1]
        assert call_kwargs["state"] == mock_milestone.state


class TestMilestoneManagerCloseMilestone:
    """Test close_milestone method."""

    def test_close_milestone_closes_milestone(
        self, milestone_manager: MilestoneManager, mock_milestone: Mock
    ) -> None:
        """Verify close_milestone closes milestone."""
        milestone_manager.close_milestone(mock_milestone)

        mock_milestone.edit.assert_called_once_with(state="closed")


class TestMilestoneManagerGetMilestoneProgress:
    """Test get_milestone_progress method."""

    def test_get_milestone_progress_calculates_correctly(
        self, milestone_manager: MilestoneManager, mock_milestone: Mock
    ) -> None:
        """Verify get_milestone_progress calculates progress correctly."""
        mock_milestone.open_issues = 5
        mock_milestone.closed_issues = 3

        progress = milestone_manager.get_milestone_progress(mock_milestone)

        assert progress["open_issues"] == 5
        assert progress["closed_issues"] == 3
        assert progress["progress_percent"] == 37  # 3/8 * 100

    def test_get_milestone_progress_handles_zero_total(
        self, milestone_manager: MilestoneManager, mock_milestone: Mock
    ) -> None:
        """Verify get_milestone_progress handles zero total issues."""
        mock_milestone.open_issues = 0
        mock_milestone.closed_issues = 0

        progress = milestone_manager.get_milestone_progress(mock_milestone)

        assert progress["open_issues"] == 0
        assert progress["closed_issues"] == 0
        assert progress["progress_percent"] == 0

    def test_get_milestone_progress_handles_all_closed(
        self, milestone_manager: MilestoneManager, mock_milestone: Mock
    ) -> None:
        """Verify get_milestone_progress handles all issues closed."""
        mock_milestone.open_issues = 0
        mock_milestone.closed_issues = 10

        progress = milestone_manager.get_milestone_progress(mock_milestone)

        assert progress["progress_percent"] == 100

    def test_get_milestone_progress_handles_all_open(
        self, milestone_manager: MilestoneManager, mock_milestone: Mock
    ) -> None:
        """Verify get_milestone_progress handles all issues open."""
        mock_milestone.open_issues = 10
        mock_milestone.closed_issues = 0

        progress = milestone_manager.get_milestone_progress(mock_milestone)

        assert progress["progress_percent"] == 0


class TestMilestoneManagerErrorHandling:
    """Test error handling."""

    def test_get_or_create_milestone_handles_api_error(
        self, milestone_manager: MilestoneManager, mock_repo: Mock
    ) -> None:
        """Verify get_or_create_milestone handles API errors."""
        api_error = GithubException(500, {"message": "Internal server error"}, headers={})
        mock_repo.get_milestones.return_value = []
        mock_repo.create_milestone = Mock(side_effect=api_error)

        with pytest.raises(GithubException) as exc_info:
            milestone_manager.get_or_create_milestone("F001", "Test")

        assert exc_info.value.status == 500

    def test_get_or_create_milestone_handles_rate_limit(
        self, milestone_manager: MilestoneManager, mock_repo: Mock
    ) -> None:
        """Verify get_or_create_milestone handles rate limit errors."""
        rate_limit_error = GithubException(
            403,
            {"message": "API rate limit exceeded"},
            headers={"X-RateLimit-Remaining": "0"},
        )
        mock_repo.get_milestones = Mock(side_effect=rate_limit_error)

        with pytest.raises(GithubException) as exc_info:
            milestone_manager.get_or_create_milestone("F001", "Test")

        assert exc_info.value.status == 403

    def test_get_or_create_milestone_handles_authentication_error(
        self, milestone_manager: MilestoneManager, mock_repo: Mock
    ) -> None:
        """Verify get_or_create_milestone handles authentication errors."""
        auth_error = GithubException(401, {"message": "Bad credentials"}, headers={})
        mock_repo.get_milestones = Mock(side_effect=auth_error)

        with pytest.raises(GithubException) as exc_info:
            milestone_manager.get_or_create_milestone("F001", "Test")

        assert exc_info.value.status == 401

    def test_update_milestone_handles_api_error(
        self, milestone_manager: MilestoneManager, mock_milestone: Mock
    ) -> None:
        """Verify update_milestone handles API errors."""
        api_error = GithubException(500, {"message": "Internal server error"}, headers={})
        mock_milestone.edit = Mock(side_effect=api_error)

        with pytest.raises(GithubException) as exc_info:
            milestone_manager.update_milestone(mock_milestone, title="New Title")

        assert exc_info.value.status == 500

    def test_close_milestone_handles_api_error(
        self, milestone_manager: MilestoneManager, mock_milestone: Mock
    ) -> None:
        """Verify close_milestone handles API errors."""
        api_error = GithubException(500, {"message": "Internal server error"}, headers={})
        mock_milestone.edit = Mock(side_effect=api_error)

        with pytest.raises(GithubException) as exc_info:
            milestone_manager.close_milestone(mock_milestone)

        assert exc_info.value.status == 500

    def test_find_milestone_for_feature_handles_api_error(
        self, milestone_manager: MilestoneManager, mock_repo: Mock
    ) -> None:
        """Verify find_milestone_for_feature handles API errors."""
        api_error = GithubException(500, {"message": "Internal server error"}, headers={})
        mock_repo.get_milestones = Mock(side_effect=api_error)

        with pytest.raises(GithubException) as exc_info:
            milestone_manager.find_milestone_for_feature("F001")

        assert exc_info.value.status == 500


class TestMilestoneManagerEdgeCases:
    """Test edge cases."""

    def test_get_or_create_milestone_handles_special_characters(
        self, milestone_manager: MilestoneManager, mock_repo: Mock, mock_milestone: Mock
    ) -> None:
        """Verify get_or_create_milestone handles special characters in title."""
        mock_repo.get_milestones.return_value = []
        mock_repo.create_milestone.return_value = mock_milestone

        result = milestone_manager.get_or_create_milestone(
            "F001",
            title="Feature & < > \" ' Title",
        )

        assert result is mock_milestone
        call_kwargs = mock_repo.create_milestone.call_args[1]
        assert "Feature & < > \" ' Title" in call_kwargs["title"]

    def test_find_milestone_for_feature_handles_case_sensitivity(
        self, milestone_manager: MilestoneManager, mock_repo: Mock
    ) -> None:
        """Verify find_milestone_for_feature handles case sensitivity."""
        milestone = Mock(spec=Milestone.Milestone)
        milestone.title = "Feature F001: Test"
        mock_repo.get_milestones.return_value = [milestone]

        # Should match case-insensitively based on prefix
        result = milestone_manager.find_milestone_for_feature("F001")

        assert result is milestone

    def test_get_milestone_progress_handles_large_numbers(
        self, milestone_manager: MilestoneManager, mock_milestone: Mock
    ) -> None:
        """Verify get_milestone_progress handles large numbers."""
        mock_milestone.open_issues = 1000
        mock_milestone.closed_issues = 500

        progress = milestone_manager.get_milestone_progress(mock_milestone)

        assert progress["progress_percent"] == 33  # 500/1500 * 100

    def test_update_milestone_with_none_values(
        self, milestone_manager: MilestoneManager, mock_milestone: Mock
    ) -> None:
        """Verify update_milestone handles None values correctly."""
        milestone_manager.update_milestone(mock_milestone, title=None, description=None)

        call_kwargs = mock_milestone.edit.call_args[1]
        # Should only include state, not None values
        assert "title" not in call_kwargs or call_kwargs.get("title") is None
        assert "description" not in call_kwargs or call_kwargs.get("description") is None
        assert "state" in call_kwargs
