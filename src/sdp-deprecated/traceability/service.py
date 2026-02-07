"""Traceability service for AC→Test mapping."""

import re
from pathlib import Path

import yaml

from sdp.beads.base import BeadsClient
from sdp.beads.models import BeadsTask
from sdp.traceability.models import (
    ACTestMapping,
    MappingStatus,
    TraceabilityReport,
)


class TraceabilityService:
    """Check and manage AC→Test traceability."""

    def __init__(self, client: BeadsClient):
        """Initialize service with Beads client.

        Args:
            client: Beads client instance
        """
        self._client = client

    def check_traceability(self, ws_id: str) -> TraceabilityReport:
        """Check traceability for workstream.

        Args:
            ws_id: Workstream ID (e.g., "00-032-01")

        Returns:
            Traceability report

        Raises:
            ValueError: If workstream not found
        """
        # Get WS from Beads or markdown fallback
        task = self._get_ws_task(ws_id)
        if task:
            acs = self._extract_acs(task.description or "")
            stored_mappings = task.sdp_metadata.get("traceability", [])
        else:
            # Markdown fallback: read from docs/workstreams/
            content = self._get_ws_content_from_markdown(ws_id)
            if not content:
                raise ValueError(f"WS not found: {ws_id}")
            acs = self._extract_acs(content)
            stored_mappings = self._get_traceability_from_markdown(content)

        # Build report
        mappings = []
        for ac_id, ac_desc in acs:
            # Find stored mapping
            stored = next(
                (m for m in stored_mappings if m["ac_id"] == ac_id),
                None,
            )

            if stored:
                mappings.append(ACTestMapping.from_dict(stored))
            else:
                mappings.append(
                    ACTestMapping(
                        ac_id=ac_id,
                        ac_description=ac_desc,
                        test_file=None,
                        test_name=None,
                        status=MappingStatus.MISSING,
                    )
                )

        return TraceabilityReport(ws_id=ws_id, mappings=mappings)

    def add_mapping(
        self, ws_id: str, ac_id: str, test_file: str, test_name: str
    ) -> None:
        """Add AC→Test mapping.

        Persists to Beads task metadata if WS in Beads, else to markdown frontmatter.

        Args:
            ws_id: Workstream ID
            ac_id: Acceptance criterion ID (e.g., "AC1")
            test_file: Test file path
            test_name: Test function name

        Raises:
            ValueError: If workstream not found
        """
        task = self._get_ws_task(ws_id)
        if task:
            self._add_mapping_to_beads(task, ws_id, ac_id, test_file, test_name)
        else:
            self._add_mapping_to_markdown(ws_id, ac_id, test_file, test_name)

    def _add_mapping_to_beads(
        self,
        task: object,
        ws_id: str,
        ac_id: str,
        test_file: str,
        test_name: str,
    ) -> None:
        """Add mapping to Beads task metadata."""
        metadata = getattr(task, "sdp_metadata", {}).copy()
        mappings = list(metadata.get("traceability", []))
        acs = self._extract_acs(getattr(task, "description", "") or "")

        existing = next((m for m in mappings if m.get("ac_id") == ac_id), None)
        if existing:
            existing["test_file"] = test_file
            existing["test_name"] = test_name
            existing["status"] = "mapped"
        else:
            ac_desc = next((d for aid, d in acs if aid == ac_id), "")
            mappings.append(
                {
                    "ac_id": ac_id,
                    "ac_description": ac_desc,
                    "test_file": test_file,
                    "test_name": test_name,
                    "status": "mapped",
                    "confidence": 1.0,
                }
            )
        metadata["traceability"] = mappings
        setattr(task, "sdp_metadata", metadata)

    def _add_mapping_to_markdown(
        self, ws_id: str, ac_id: str, test_file: str, test_name: str
    ) -> None:
        """Add mapping to markdown frontmatter (fallback when Beads has no task)."""
        ws_path = self._get_ws_file_path(ws_id)
        if not ws_path:
            raise ValueError(f"WS not found: {ws_id}")

        content = ws_path.read_text(encoding="utf-8")
        acs = self._extract_acs(content)
        ac_desc = next((d for aid, d in acs if aid == ac_id), "")

        match = re.match(r"^---\n(.*?)\n---", content, re.DOTALL)
        if not match:
            raise ValueError(f"No frontmatter in {ws_path}")

        fm_text = match.group(1)
        fm = yaml.safe_load(fm_text) or {}
        mappings = list(fm.get("traceability", []))

        existing = next((m for m in mappings if m.get("ac_id") == ac_id), None)
        if existing:
            existing["test_file"] = test_file
            existing["test_name"] = test_name
            existing["status"] = "mapped"
        else:
            mappings.append(
                {
                    "ac_id": ac_id,
                    "ac_description": ac_desc,
                    "test_file": test_file,
                    "test_name": test_name,
                    "status": "mapped",
                    "confidence": 1.0,
                }
            )
        fm["traceability"] = mappings

        new_fm = yaml.dump(fm, default_flow_style=False, allow_unicode=True)
        body = re.sub(r"^---\n.*?\n---", "", content, count=1, flags=re.DOTALL).lstrip()
        ws_path.write_text(f"---\n{new_fm}---\n\n{body}", encoding="utf-8")

    def _get_ws_task(self, ws_id: str) -> BeadsTask | None:
        """Get Beads task for workstream.

        Args:
            ws_id: Workstream ID

        Returns:
            Beads task or None if not found

        Note:
            Currently uses external_ref to find task.
            In future, may use dedicated WS index.
        """
        try:
            tasks = self._client.list_tasks()
        except Exception:
            return None
        for task in tasks:
            if task.external_ref == ws_id:
                return task
        return None

    def _get_ws_file_path(self, ws_id: str) -> Path | None:
        """Get path to WS markdown file.

        Searches docs/workstreams/backlog and docs/workstreams/completed.
        """
        cwd = Path.cwd()
        for subdir in ("backlog", "completed"):
            ws_dir = cwd / "docs" / "workstreams" / subdir
            if not ws_dir.exists():
                continue
            for f in ws_dir.glob(f"{ws_id}-*.md"):
                return f
        return None

    def _get_ws_content_from_markdown(self, ws_id: str) -> str | None:
        """Get WS content from markdown file (fallback when Beads has no task).

        Searches docs/workstreams/backlog and docs/workstreams/completed.
        """
        path = self._get_ws_file_path(ws_id)
        return path.read_text(encoding="utf-8") if path else None

    def _get_traceability_from_markdown(self, content: str) -> list[dict[str, object]]:
        """Extract traceability mappings from markdown frontmatter."""
        match = re.match(r"^---\n(.*?)\n---", content, re.DOTALL)
        if not match:
            return []
        try:
            fm = yaml.safe_load(match.group(1))
            return list(fm.get("traceability", [])) if isinstance(fm, dict) else []
        except yaml.YAMLError:
            return []

    def _extract_acs(self, description: str) -> list[tuple[str, str]]:
        """Extract ACs from WS description.

        Looks for patterns like:
        - [ ] AC1: Description
        - AC1: Description

        Args:
            description: Task description text

        Returns:
            List of (ac_id, ac_description) tuples
        """
        pattern = r"(?:- \[[ x]\] )?(AC\d+):\s*(.+?)(?:\n|$)"
        matches = re.findall(pattern, description, re.IGNORECASE)
        return [(m[0].upper(), m[1].strip()) for m in matches]
