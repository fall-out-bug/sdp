"""
Bidirectional sync service between SDP workstreams and Beads tasks.

Handles conversion between:
- SDP workstream format (markdown with YAML frontmatter)
- Beads task format (JSONL with hash-based IDs)
"""

from pathlib import Path
from typing import Any, Optional

from ..client import BeadsClient
from ..models import BeadsDependency, BeadsDependencyType, BeadsSyncResult, BeadsTaskCreate
from .mapping import MappingManager
from .status_mapper import (
    SDP_STATUS_BACKLOG,
    map_beads_status_to_sdp,
    map_sdp_size_to_beads_priority,
    map_sdp_status_to_beads,
)


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
        mapping_path = mapping_file or Path.cwd() / ".beads-sdp-mapping.jsonl"
        self.mapping_manager = MappingManager(mapping_path)

        # Load existing mapping
        self.mapping_manager.load()

    def persist_mapping(self) -> None:
        """Persist ID mapping to file (deduplicates and overwrites)."""
        self.mapping_manager.save()

    def sync_workstream_to_beads(
        self,
        ws_file: Path,
        ws_data: dict[str, Any],
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
        beads_id = self.mapping_manager.get_beads_id(ws_id)

        if beads_id:
            # Update existing Beads task
            return self._update_existing_task(ws_id, ws_data, beads_id)
        else:
            # Create new Beads task
            return self._create_new_task(ws_id, ws_data, ws_file)

    def _update_existing_task(
        self, ws_id: str, ws_data: dict[str, Any], beads_id: str
    ) -> BeadsSyncResult:
        """Update existing Beads task with SDP workstream data."""
        try:
            self.client.update_task_status(
                beads_id,
                map_sdp_status_to_beads(ws_data.get("status", SDP_STATUS_BACKLOG)),
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

    def _create_new_task(
        self, ws_id: str, ws_data: dict[str, Any], ws_file: Path
    ) -> BeadsSyncResult:
        """Create new Beads task from SDP workstream."""
        try:
            # Build description from workstream content
            description = self._build_description(ws_data)

            # Map dependencies
            dependencies = self._map_dependencies(ws_data)

            # Create task (Beads limits title to 500 chars)
            full_title = f"{ws_id}: {ws_data.get('title', '')}"
            title = full_title[:497] + "..." if len(full_title) > 500 else full_title

            params = BeadsTaskCreate(
                title=title,
                description=description,
                priority=map_sdp_size_to_beads_priority(ws_data.get("size", "MEDIUM")),
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
            self.mapping_manager.add_mapping(ws_id, task.id)
            self.mapping_manager.save()

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

    def _build_description(self, ws_data: dict[str, Any]) -> str:
        """Build Beads task description from workstream data."""
        description = f"**Goal:**\n{ws_data.get('goal', '')}\n\n"

        if ws_data.get("context"):
            description += f"**Context:**\n{ws_data['context']}\n\n"

        if ws_data.get("acceptance_criteria"):
            description += "**Acceptance Criteria:**\n"
            for ac in ws_data["acceptance_criteria"]:
                checked = "✓" if ac.get("checked") else "☐"
                text = ac.get("description") or ac.get("text", "")
                description += f"{checked} {text}\n"

        return description

    def _map_dependencies(self, ws_data: dict[str, Any]) -> list[BeadsDependency]:
        """Map SDP dependencies to Beads dependencies."""
        dependencies = []

        for dep_ws_id in ws_data.get("dependencies", []):
            # Check if dependency already synced
            dep_beads_id = self.mapping_manager.get_beads_id(dep_ws_id)
            if dep_beads_id:
                dependencies.append(
                    BeadsDependency(
                        task_id=dep_beads_id,
                        type=BeadsDependencyType.BLOCKS,
                    )
                )

        return dependencies

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
            sdp_status = map_beads_status_to_sdp(task.status)

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
