"""GitHub sync orchestration service."""

from pathlib import Path
from typing import Optional

from sdp.github.client import GitHubClient
from sdp.github.frontmatter_updater import FrontmatterUpdater
from sdp.github.issue_sync import IssueSync
from sdp.github.label_manager import LabelManager
from sdp.github.milestone_manager import MilestoneManager
from sdp.github.project_board_sync import ProjectBoardSync
from sdp.github.status_mapper import StatusMapper
from sdp.github.sync_helpers import SyncHelpers
from sdp.github.sync_models import SyncResult
from sdp.github.ws_parser import parse_ws_file


class SyncService:
    """Orchestrate workstream → GitHub sync operations.

    Coordinates:
    - WS file parsing
    - Issue creation
    - Label management
    - Project board sync
    - Result tracking
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
            project_name = SyncHelpers.resolve_project_name(ws_file, self._project_override)
            milestone, self._milestone_manager = SyncHelpers.get_feature_milestone(
                self._client, ws.feature, self._milestone_manager
            )
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
                SyncHelpers.update_issue_status(issue, ws.status)
                SyncHelpers.ensure_issue_milestone(issue, milestone_number, milestone)

                # Update project board status
                board_sync = SyncHelpers.get_board_sync(
                    self._client, project_name, self._board_sync_by_project
                )
                SyncHelpers.sync_to_board(board_sync, issue, ws.status)

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
                board_sync = SyncHelpers.get_board_sync(
                    self._client, project_name, self._board_sync_by_project
                )
                SyncHelpers.sync_to_board(board_sync, issue, ws.status)

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
