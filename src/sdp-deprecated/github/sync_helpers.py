"""
Helper utilities for GitHub sync service.

Provides board sync coordination and milestone management.
"""

from pathlib import Path
from typing import Optional

from github import Issue

from sdp.github.client import GitHubClient
from sdp.github.milestone_manager import MilestoneManager
from sdp.github.project_board_sync import ProjectBoardSync
from sdp.github.project_router import ProjectRouter
from sdp.github.projects_client import ProjectsClient
from sdp.github.status_mapper import StatusMapper


class SyncHelpers:
    """Helper methods for sync service."""

    @staticmethod
    def update_issue_status(issue: Issue.Issue, ws_status: str) -> None:
        """Update GitHub issue status to match WS.

        Updates issue labels and state to match WS status mapping.

        Args:
            issue: PyGithub Issue instance
            ws_status: WS status (backlog, active, completed, blocked)
        """
        # Get expected label and state from WS status
        expected_label = StatusMapper.ws_to_github_label(ws_status)
        expected_state = StatusMapper.ws_to_github_state(ws_status)

        # Get current labels
        current_labels = [label.name for label in issue.labels]

        # Remove old status label, keep others
        new_labels = [
            label for label in current_labels if not label.startswith("status/")
        ]
        # Add new status label
        new_labels.append(expected_label)

        # Update issue with new labels and state
        issue.edit(
            labels=new_labels,
            state=expected_state,
        )

    @staticmethod
    def sync_to_board(
        board_sync: ProjectBoardSync | None,
        issue: Issue.Issue,
        ws_status: str,
    ) -> None:
        """Sync issue to project board (non-blocking).

        Args:
            board_sync: Board sync instance (can be None)
            issue: PyGithub Issue instance
            ws_status: WS status for column placement
        """
        if board_sync is None:
            return

        try:
            board_sync.update_issue_status(issue, ws_status)
        except Exception as e:
            # Non-critical: log but don't fail sync
            print(f"Warning: Board sync failed for #{issue.number}: {e}")

    @staticmethod
    def resolve_project_name(
        ws_file: Path, project_override: str | None
    ) -> str:
        """Resolve project name for a workstream file."""
        if project_override and project_override != "auto":
            return project_override
        frontmatter_project = ProjectRouter.get_project_from_frontmatter(ws_file)
        if frontmatter_project:
            return frontmatter_project
        return ProjectRouter.get_project_for_ws(ws_file)

    @staticmethod
    def get_board_sync(
        client: GitHubClient,
        project_name: str,
        cache: dict[str, ProjectBoardSync | None],
    ) -> ProjectBoardSync | None:
        """Get or create board sync for project (non-blocking)."""
        if project_name in cache:
            return cache[project_name]

        try:
            config = client._config
            owner = config.org or config.repo.split("/")[0]
            projects_client = ProjectsClient(config.token, owner)
            board_sync = ProjectBoardSync(projects_client, project_name)
        except Exception as e:
            print(f"Warning: Project board sync disabled: {e}")
            board_sync = None

        cache[project_name] = board_sync
        return board_sync

    @staticmethod
    def get_feature_milestone(
        client: GitHubClient,
        feature_id: str,
        milestone_manager: Optional[MilestoneManager],
    ) -> tuple[object, Optional[MilestoneManager]]:
        """Get or create milestone for feature.

        Returns:
            Tuple of (milestone, milestone_manager)
        """
        if milestone_manager is None:
            repo = client.get_repo()
            milestone_manager = MilestoneManager(repo)
        milestone = milestone_manager.get_or_create_feature_milestone(feature_id)
        return milestone, milestone_manager

    @staticmethod
    def ensure_issue_milestone(
        issue: Issue.Issue, milestone_number: int | None, milestone: object
    ) -> None:
        """Ensure issue milestone matches expected milestone."""
        if milestone_number is None:
            return
        current = getattr(issue, "milestone", None)
        if current is None or getattr(current, "number", None) != milestone_number:
            issue.edit(milestone=milestone)  # type: ignore[arg-type]
