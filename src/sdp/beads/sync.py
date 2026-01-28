"""
Bidirectional sync between SDP workstreams and Beads tasks.

Handles conversion between:
- SDP workstream format (markdown with YAML frontmatter)
- Beads task format (JSONL with hash-based IDs)

Mapping:
- SDP workstream ID (PP-FFF-SS) ↔ Beads task ID (bd-XXXX) via mapping table
- SDP status (backlog|active|completed|blocked) ↔ Beads status (open|in_progress|closed|blocked)
- SDP size (SMALL|MEDIUM|LARGE) ↔ Beads priority (2|1|0)
- SDP dependencies (list of WS IDs) ↔ Beads dependencies (type: "blocks")
"""

import hashlib
import json
from datetime import datetime
from pathlib import Path
from typing import Optional, Dict

from .client import BeadsClient
from .models import (
    BeadsTask,
    BeadsTaskCreate,
    BeadsStatus,
    BeadsPriority,
    BeadsDependency,
    BeadsDependencyType,
    BeadsSyncResult,
)

# SDP workstream status values
SDP_STATUS_BACKLOG = "backlog"
SDP_STATUS_ACTIVE = "active"
SDP_STATUS_COMPLETED = "completed"
SDP_STATUS_BLOCKED = "blocked"

# Beads status values (for reference)
# open, in_progress, blocked, deferred, closed, tombstone, pinned, hooked


class BeadsSyncError(Exception):
    """Exception raised during sync operations."""

    pass


class BeadsSyncService:
    """Bidirectional sync service between SDP and Beads.

    Responsibilities:
    - Convert SDP workstreams → Beads tasks
    - Convert Beads tasks → SDP workstreams (status updates only)
    - Maintain ID mapping table (PP-FFF-SS ↔ bd-XXXX)
    - Handle conflict resolution (SDP wins for content, Beads wins for execution metadata)
    """

    def __init__(
        self,
        client: BeadsClient,
        mapping_file: Optional[Path] = None,
    ):
        """Initialize sync service.

        Args:
            client: BeadsClient instance (mock or real)
            mapping_file: Path to ID mapping table (defaults to .beads-sdp-mapping.jsonl)
        """
        self.client = client
        self.mapping_file = mapping_file or Path.cwd() / ".beads-sdp-mapping.jsonl"
        self._mapping: Dict[str, str] = {}  # sdp_id → beads_id
        self._reverse_mapping: Dict[str, str] = {}  # beads_id → sdp_id

        # Load existing mapping
        self._load_mapping()

    def _load_mapping(self) -> None:
        """Load ID mapping from file."""
        if not self.mapping_file.exists():
            return

        try:
            with open(self.mapping_file, "r") as f:
                for line in f:
                    if not line.strip():
                        continue

                    entry = json.loads(line)
                    sdp_id = entry.get("sdp_id")
                    beads_id = entry.get("beads_id")

                    if sdp_id and beads_id:
                        self._mapping[sdp_id] = beads_id
                        self._reverse_mapping[beads_id] = sdp_id

        except (json.JSONDecodeError, IOError) as e:
            raise BeadsSyncError(f"Failed to load mapping file: {e}") from e

    def _save_mapping(self) -> None:
        """Save ID mapping to file (append-only)."""
        try:
            with open(self.mapping_file, "a") as f:
                # Write entire mapping (inefficient but simple)
                for sdp_id, beads_id in self._mapping.items():
                    entry = {
                        "sdp_id": sdp_id,
                        "beads_id": beads_id,
                        "updated_at": datetime.utcnow().isoformat(),
                    }
                    f.write(json.dumps(entry) + "\n")

        except IOError as e:
            raise BeadsSyncError(f"Failed to save mapping file: {e}") from e

    def _map_sdp_status_to_beads(self, sdp_status: str) -> BeadsStatus:
        """Map SDP status to Beads status."""
        mapping = {
            SDP_STATUS_BACKLOG: BeadsStatus.OPEN,
            SDP_STATUS_ACTIVE: BeadsStatus.IN_PROGRESS,
            SDP_STATUS_COMPLETED: BeadsStatus.CLOSED,
            SDP_STATUS_BLOCKED: BeadsStatus.BLOCKED,
        }
        return mapping.get(sdp_status, BeadsStatus.OPEN)

    def _map_beads_status_to_sdp(self, beads_status: BeadsStatus) -> str:
        """Map Beads status to SDP status."""
        mapping = {
            BeadsStatus.OPEN: SDP_STATUS_BACKLOG,
            BeadsStatus.IN_PROGRESS: SDP_STATUS_ACTIVE,
            BeadsStatus.CLOSED: SDP_STATUS_COMPLETED,
            BeadsStatus.BLOCKED: SDP_STATUS_BLOCKED,
        }
        return mapping.get(beads_status, SDP_STATUS_BACKLOG)

    def _map_sdp_size_to_beads_priority(self, size: str) -> BeadsPriority:
        """Map SDP size to Beads priority."""
        # Larger workstream = higher priority (lower number = more critical)
        mapping = {
            "SMALL": BeadsPriority(2),
            "MEDIUM": BeadsPriority(1),
            "LARGE": BeadsPriority(0),
        }
        return mapping.get(size.upper(), BeadsPriority(2))

    def sync_workstream_to_beads(
        self,
        ws_file: Path,
        ws_data: dict,
    ) -> BeadsSyncResult:
        """Sync SDP workstream → Beads task (create or update).

        Args:
            ws_file: Path to workstream markdown file
            ws_data: Parsed workstream data (frontmatter + body)

        Returns:
            SyncResult with beads_id and status

        Example:
            sync = BeadsSyncService(client)
            result = sync.sync_workstream_to_beads(
                Path("docs/workstreams/backlog/00-001-01.md"),
                {"ws_id": "00-001-01", "title": "Domain entities", ...}
            )
            print(f"Created: {result.beads_id}")
        """
        ws_id = ws_data.get("ws_id")
        if not ws_id:
            return BeadsSyncResult(
                success=False,
                task_id=ws_file.name,
                error="Missing ws_id in workstream data",
            )

        # Check if already synced
        beads_id = self._mapping.get(ws_id)

        if beads_id:
            # Update existing Beads task
            try:
                self.client.update_task_status(
                    beads_id,
                    self._map_sdp_status_to_beads(ws_data.get("status", SDP_STATUS_BACKLOG)),
                )

                return BeadsSyncResult(
                    success=True,
                    task_id=ws_id,
                    beads_id=beads_id,
                    message="Updated existing Beads task",
                )

            except Exception as e:
                return BeadsSyncResult(
                    success=False,
                    task_id=ws_id,
                    beads_id=beads_id,
                    error=f"Failed to update: {e}",
                )

        else:
            # Create new Beads task
            try:
                # Build description from workstream content
                description = f"**Goal:**\n{ws_data.get('goal', '')}\n\n"
                if ws_data.get("context"):
                    description += f"**Context:**\n{ws_data['context']}\n\n"
                if ws_data.get("acceptance_criteria"):
                    description += f"**Acceptance Criteria:**\n"
                    for ac in ws_data["acceptance_criteria"]:
                        checked = "✓" if ac.get("checked") else "☐"
                        description += f"{checked} {ac.get('text', '')}\n"

                # Map dependencies
                dependencies = []
                for dep_ws_id in ws_data.get("dependencies", []):
                    # Check if dependency already synced
                    dep_beads_id = self._mapping.get(dep_ws_id)
                    if dep_beads_id:
                        dependencies.append(
                            BeadsDependency(
                                task_id=dep_beads_id,
                                type=BeadsDependencyType.BLOCKS,
                            )
                        )

                # Create task
                params = BeadsTaskCreate(
                    title=f"{ws_id}: {ws_data.get('title', '')}",
                    description=description,
                    priority=self._map_sdp_size_to_beads_priority(ws_data.get("size", "MEDIUM")),
                    dependencies=dependencies,
                    external_ref=f"PP-FFF-SS:{ws_id}",
                    sdp_metadata={
                        "ws_id": ws_id,
                        "feature": ws_data.get("feature"),
                        "file_path": str(ws_file),
                    },
                )

                task = self.client.create_task(params)

                # Store mapping
                self._mapping[ws_id] = task.id
                self._reverse_mapping[task.id] = ws_id
                self._save_mapping()

                return BeadsSyncResult(
                    success=True,
                    task_id=ws_id,
                    beads_id=task.id,
                    message="Created new Beads task",
                )

            except Exception as e:
                return BeadsSyncResult(
                    success=False,
                    task_id=ws_id,
                    error=f"Failed to create: {e}",
                )

    def sync_beads_to_workstream(
        self,
        beads_id: str,
        ws_file: Path,
    ) -> BeadsSyncResult:
        """Sync Beads task → SDP workstream (status updates only).

        Args:
            beads_id: Beads task ID
            ws_file: Path to workstream markdown file

        Returns:
            SyncResult with updated status

        Note:
            SDP is authoritative for content (title, goal, acceptance criteria).
            This method only updates status from Beads.
        """
        try:
            # Get Beads task
            task = self.client.get_task(beads_id)
            if not task:
                return BeadsSyncResult(
                    success=False,
                    task_id=beads_id,
                    error=f"Beads task not found: {beads_id}",
                )

            # Map Beads status → SDP status
            sdp_status = self._map_beads_status_to_sdp(task.status)

            # Update workstream frontmatter
            # (Implementation would read yaml, update status field, write back)
            # For now, just return success with mapped status
            return BeadsSyncResult(
                success=True,
                task_id=beads_id,
                message=f"Status would be updated to: {sdp_status}",
            )

        except Exception as e:
            return BeadsSyncResult(
                success=False,
                task_id=beads_id,
                error=f"Failed to sync: {e}",
            )
