"""Enhanced sync service with conflict detection."""

import logging
from pathlib import Path

from .conflict_resolver import Conflict, ConflictResolver
from .sync_service import SyncService

logger = logging.getLogger(__name__)


class EnhancedSyncService:
    """Enhanced sync with conflict detection."""

    def __init__(self, base_service: SyncService) -> None:
        """Initialize enhanced sync service.

        Args:
            base_service: Base SyncService instance
        """
        self._base = base_service
        self._resolver = ConflictResolver()

    def sync_with_conflict_detection(
        self, ws_id: str, dry_run: bool = False
    ) -> dict[str, bool | None]:
        """Sync workstream with conflict detection.

        Args:
            ws_id: Workstream ID to sync
            dry_run: If True, preview changes without applying

        Returns:
            Dict with sync results
        """
        result = {
            "has_conflict": False,
            "resolved": False,
            "synced": False,
        }

        # Read states
        ws_state = self._read_ws_state(ws_id)
        gh_state = self._read_gh_state(ws_id)

        # Detect conflicts
        conflict = self._resolver.detect(ws_state, gh_state)
        if conflict:
            result["has_conflict"] = True
            resolved_value = self._resolver.resolve(conflict)

            if dry_run:
                logger.info(
                    f"[DRY RUN] Would resolve conflict: {conflict.field} "
                    f"{conflict.gh_value} -> {resolved_value}"
                )
            else:
                logger.info(
                    f"Resolving conflict: {conflict.field} "
                    f"{conflict.gh_value} -> {resolved_value}"
                )
                self._update_gh_state(ws_id, {conflict.field: resolved_value})

            result["resolved"] = True
            result["synced"] = not dry_run
        else:
            result["synced"] = True

        return result

    def sync_backlog(
        self, ws_dir: str = "docs/workstreams", dry_run: bool = False
    ) -> list[dict]:
        """Sync all backlog workstreams with conflict detection.

        Args:
            ws_dir: Workstreams directory
            dry_run: If True, preview changes without applying

        Returns:
            List of sync results for each workstream
        """
        results = []
        backlog_dir = Path(ws_dir) / "backlog"

        if not backlog_dir.exists():
            return results

        for ws_file in backlog_dir.glob("*.md"):
            ws_id = ws_file.stem
            try:
                result = self.sync_with_conflict_detection(ws_id, dry_run)
                result["ws_id"] = ws_id
                results.append(result)
            except Exception as e:
                logger.warning(f"Failed to sync {ws_id}: {e}")
                results.append({"ws_id": ws_id, "error": str(e)})

        return results

    def _read_ws_state(self, ws_id: str) -> dict:
        """Read workstream state from file.

        Args:
            ws_id: Workstream ID

        Returns:
            Dict with workstream state
        """
        from .ws_parser import WorkstreamParser

        # Try to find the workstream file
        for status in ["backlog", "in_progress", "completed"]:
            ws_path = Path(self._base.ws_dir) / status / f"{ws_id}.md"
            if ws_path.exists():
                parser = WorkstreamParser()
                ws = parser.parse(ws_path)
                return {"status": ws.status.value}

        return {}

    def _read_gh_state(self, ws_id: str) -> dict:
        """Read state from GitHub.

        Args:
            ws_id: Workstream ID

        Returns:
            Dict with GitHub state
        """
        # This would call the actual GitHub API
        # For now, return empty dict
        return {}

    def _update_gh_state(self, ws_id: str, updates: dict) -> None:
        """Update GitHub state.

        Args:
            ws_id: Workstream ID
            updates: Dict of fields to update
        """
        # This would call the actual GitHub API
        logger.info(f"Would update GitHub for {ws_id}: {updates}")
