"""GitHub milestone management."""

from typing import Any, Optional

from github import Milestone, Repository


class MilestoneManager:
    """Manage GitHub milestones for features.

    Attributes:
        _repo: GitHub repository instance
    """

    def __init__(self, repo: Repository.Repository) -> None:
        """Initialize milestone manager.

        Args:
            repo: PyGithub Repository instance
        """
        self._repo = repo

    def get_or_create_milestone(
        self,
        feature_id: str,
        title: str,
        description: str = "",
        due_date: Optional[str] = None,
    ) -> Milestone.Milestone:
        """Get or create milestone for feature.

        Checks if milestone with formatted title exists. If not, creates it.

        Args:
            feature_id: Feature identifier (e.g., "F100")
            title: Milestone title
            description: Milestone description
            due_date: Optional due date

        Returns:
            GitHub Milestone object
        """
        formatted_title = f"Feature {feature_id}: {title}"
        milestone = self._find_milestone_by_title(formatted_title)

        if milestone:
            return milestone

        kwargs: dict[str, Any] = {
            "title": formatted_title,
            "description": description,
        }
        if due_date:
            kwargs["due_on"] = due_date

        return self._repo.create_milestone(**kwargs)

    def _find_milestone_by_title(
        self, title: str
    ) -> Optional[Milestone.Milestone]:
        """Find milestone by exact title match.

        Args:
            title: Milestone title to find

        Returns:
            Milestone if found, None otherwise
        """
        milestones = self._repo.get_milestones(state="open")
        for milestone in milestones:
            if milestone.title == title:
                return milestone
        return None

    def find_milestone_for_feature(
        self, feature_id: str
    ) -> Optional[Milestone.Milestone]:
        """Find milestone matching a feature ID.

        Args:
            feature_id: Feature identifier (e.g., "F160")

        Returns:
            Milestone if found, None otherwise
        """
        prefix = f"Feature {feature_id}:"
        exact = f"Feature {feature_id}"
        milestones = self._repo.get_milestones(state="open")
        for milestone in milestones:
            if milestone.title.startswith(prefix) or milestone.title == exact:
                return milestone
        return None

    def get_or_create_feature_milestone(
        self,
        feature_id: str,
        title: Optional[str] = None,
        description: str = "",
        due_date: Optional[str] = None,
    ) -> Milestone.Milestone:
        """Get or create milestone for feature with optional title.

        Args:
            feature_id: Feature identifier (e.g., "F160")
            title: Optional feature title for milestone naming
            description: Milestone description
            due_date: Optional due date

        Returns:
            GitHub Milestone object
        """
        milestone = self.find_milestone_for_feature(feature_id)
        if milestone:
            return milestone

        milestone_title = f"Feature {feature_id}"
        if title:
            milestone_title = f"{milestone_title}: {title}"

        kwargs: dict[str, Any] = {
            "title": milestone_title,
            "description": description,
        }
        if due_date:
            kwargs["due_on"] = due_date

        return self._repo.create_milestone(**kwargs)

    def update_milestone(
        self,
        milestone: Milestone.Milestone,
        title: Optional[str] = None,
        description: Optional[str] = None,
        state: Optional[str] = None,
    ) -> None:
        """Update milestone properties.

        Args:
            milestone: Milestone object to update
            title: New title (optional)
            description: New description (optional)
            state: New state - "open" or "closed" (optional)
        """
        kwargs: dict[str, Any] = {}
        if title is not None:
            kwargs["title"] = title
        if description is not None:
            kwargs["description"] = description
        if state is not None:
            kwargs["state"] = state
        else:
            # Preserve current state if not specified
            kwargs["state"] = milestone.state

        milestone.edit(**kwargs)

    def close_milestone(self, milestone: Milestone.Milestone) -> None:
        """Close a milestone.

        Args:
            milestone: Milestone object to close
        """
        milestone.edit(state="closed")  # type: ignore[call-arg]

    def get_milestone_progress(
        self, milestone: Milestone.Milestone
    ) -> dict[str, int]:
        """Calculate milestone progress.

        Args:
            milestone: Milestone object

        Returns:
            Dictionary with open_issues, closed_issues, progress_percent
        """
        open_count = milestone.open_issues
        closed_count = milestone.closed_issues
        total = open_count + closed_count

        progress_percent = 0
        if total > 0:
            progress_percent = int((closed_count / total) * 100)

        return {
            "open_issues": open_count,
            "closed_issues": closed_count,
            "progress_percent": progress_percent,
        }
