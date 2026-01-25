"""GitHub sync orchestration service."""

from dataclasses import dataclass
from pathlib import Path
from typing import Optional

from github import Issue

from sdp.github.client import GitHubClient
from sdp.github.frontmatter_updater import FrontmatterUpdater
from sdp.github.issue_sync import IssueSync
from sdp.github.label_manager import LabelManager
from sdp.github.milestone_manager import MilestoneManager
from sdp.github.project_board_sync import ProjectBoardSync
from sdp.github.project_router import ProjectRouter
from sdp.github.projects_client import ProjectsClient
from sdp.github.status_mapper import StatusMapper
from sdp.github.ws_parser import parse_ws_file


@dataclass
class SyncResult:
    """Result of syncing single workstream.

    Attributes:
        ws_id: Workstream ID
        action: Action taken ("created", "updated", "skipped", "failed")
        issue_number: GitHub issue number if created
        error: Error message if failed
    """

    ws_id: str
    action: str
    issue_number: Optional[int] = None
    error: Optional[str] = None


class SyncService:
    """Orchestrate workstream → GitHub sync operations.

    Coordinates:
    - WS file parsing
    - Issue creation
    - Label management
    - Project board sync
    - Result tracking

    Attributes:
        _client: GitHub API client
        _issue_sync: Issue sync coordinator
        _board_sync: Project board sync coordinator (optional)
    """

    def __init__(
        self,
        client: GitHubClient,
        project_name: str | None = None,
    ) -> None:
        """Initialize sync service with project board integration.

        Args:
            client: GitHub API client
            project_name: GitHub Project name for board sync
        """
        self._client = client
        label_manager = LabelManager(client)
        self._issue_sync = IssueSync(client, label_manager)
        self._project_override = project_name
        self._board_sync_by_project: dict[str, ProjectBoardSync | None] = {}
        self._milestone_manager: Optional[MilestoneManager] = None

    def sync_workstream(self, ws_file: Path) -> SyncResult:
        """Sync single workstream to GitHub (bidirectional).

        Creates or updates GitHub issue to match WS file. Updates WS
        frontmatter with GitHub issue number. Detects conflicts where
        WS is source of truth.

        Args:
            ws_file: Path to WS markdown file

        Returns:
            SyncResult with action taken (created, updated, or failed)
        """
        try:
            ws = parse_ws_file(ws_file)
            ws_file_path = str(ws_file)
            project_name = self._resolve_project_name(ws_file)
            milestone = self._get_feature_milestone(ws.feature)
            milestone_number = getattr(milestone, "number", None) if milestone else None

            # Check if issue already exists (from WS frontmatter)
            issue_number = FrontmatterUpdater.get_github_issue(ws_file)

            if issue_number:
                # Update existing issue
                issue = self._client.get_issue(issue_number)

                # Check for conflicts
                current_labels = [label.name for label in issue.labels]
                status_label = next(
                    (label for label in current_labels if label.startswith("status/")),
                    None,
                )

                if status_label and StatusMapper.detect_conflict(
                    ws.status, status_label, issue.state
                ):
                    print(
                        f"⚠️ Conflict detected for {ws.ws_id}: "
                        f"WS={ws.status}, GitHub={status_label}/{issue.state}"
                    )

                # Update issue to match WS (WS is source of truth)
                self._update_issue_status(issue, ws.status)
                self._ensure_issue_milestone(issue, milestone_number)

                # Update project board status
                self._sync_to_board(issue, ws.status, project_name)

                return SyncResult(
                    ws_id=ws.ws_id,
                    action="updated",
                    issue_number=issue_number,
                )
            else:
                # Create new issue
                issue_number = self._issue_sync.sync_ws(
                    ws,
                    ws_file_path,
                    milestone_number=milestone_number,
                )

                # Update WS frontmatter with GitHub issue number
                FrontmatterUpdater.update_github_issue(ws_file, issue_number)

                # Add to project board
                issue = self._client.get_issue(issue_number)
                self._sync_to_board(issue, ws.status, project_name)

                return SyncResult(
                    ws_id=ws.ws_id,
                    action="created",
                    issue_number=issue_number,
                )
        except Exception as e:
            # Extract WS ID from filename if parsing fails
            ws_id = ws_file.stem
            return SyncResult(
                ws_id=ws_id,
                action="failed",
                error=str(e),
            )

    def _update_issue_status(self, issue: Issue.Issue, ws_status: str) -> None:
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

    def _sync_to_board(
        self, issue: Issue.Issue, ws_status: str, project_name: str
    ) -> None:
        """Sync issue to project board (non-blocking).

        Args:
            issue: PyGithub Issue instance
            ws_status: WS status for column placement
            project_name: Project name for routing
        """
        board_sync = self._get_board_sync(project_name)
        if board_sync is None:
            return

        try:
            board_sync.update_issue_status(issue, ws_status)
        except Exception as e:
            # Non-critical: log but don't fail sync
            print(f"Warning: Board sync failed for #{issue.number}: {e}")

    def _resolve_project_name(self, ws_file: Path) -> str:
        """Resolve project name for a workstream file."""
        if self._project_override and self._project_override != "auto":
            return self._project_override
        frontmatter_project = ProjectRouter.get_project_from_frontmatter(ws_file)
        if frontmatter_project:
            return frontmatter_project
        return ProjectRouter.get_project_for_ws(ws_file)

    def _get_board_sync(self, project_name: str) -> ProjectBoardSync | None:
        """Get or create board sync for project (non-blocking)."""
        if project_name in self._board_sync_by_project:
            return self._board_sync_by_project[project_name]

        try:
            config = self._client._config
            owner = config.org or config.repo.split("/")[0]
            projects_client = ProjectsClient(config.token, owner)
            board_sync = ProjectBoardSync(projects_client, project_name)
        except Exception as e:
            print(f"Warning: Project board sync disabled: {e}")
            board_sync = None

        self._board_sync_by_project[project_name] = board_sync
        return board_sync

    def _get_feature_milestone(self, feature_id: str) -> object:
        """Get or create milestone for feature."""
        if self._milestone_manager is None:
            repo = self._client.get_repo()
            self._milestone_manager = MilestoneManager(repo)
        return self._milestone_manager.get_or_create_feature_milestone(feature_id)

    def _ensure_issue_milestone(
        self, issue: Issue.Issue, milestone_number: int | None
    ) -> None:
        """Ensure issue milestone matches expected milestone."""
        if milestone_number is None:
            return
        current = getattr(issue, "milestone", None)
        if current is None or getattr(current, "number", None) != milestone_number:
            issue.edit(milestone=self._get_feature_milestone(""))  # type: ignore[arg-type]

    def sync_feature(self, feature_id: str, ws_dir: Path) -> list[SyncResult]:
        """Sync all workstreams for feature.

        Finds all WS files matching feature pattern and syncs them.

        Args:
            feature_id: Feature ID (e.g., "F60")
            ws_dir: Directory with workstream files

        Returns:
            List of SyncResult for each WS
        """
        results = []

        # Find all WS files for feature
        # F60 → WS-*60*.md, F150 → WS-*150*.md
        feature_num = feature_id[1:]  # Remove 'F' prefix
        pattern = f"WS-*{feature_num}*.md"
        ws_files = sorted(ws_dir.glob(pattern))

        for ws_file in ws_files:
            result = self.sync_workstream(ws_file)
            results.append(result)

        return results

    def sync_all(self, ws_dir: Path) -> list[SyncResult]:
        """Sync all workstreams in directory.

        Finds all WS markdown files and syncs them.

        Args:
            ws_dir: Directory with workstream files

        Returns:
            List of SyncResult for each WS
        """
        results = []

        # Find all WS files
        ws_files = sorted(ws_dir.glob("WS-*.md"))

        for ws_file in ws_files:
            result = self.sync_workstream(ws_file)
            results.append(result)

        return results
