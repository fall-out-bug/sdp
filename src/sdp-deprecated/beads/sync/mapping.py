"""
ID mapping between SDP workstreams and Beads tasks.

Maintains bidirectional mapping: PP-FFF-SS ↔ bd-XXXX
"""

import json
from datetime import datetime
from pathlib import Path
from typing import Dict


class BeadsSyncError(Exception):
    """Exception raised during sync operations."""

    pass


def resolve_ws_id_to_beads_id(ws_id: str, mapping_file: Path | None = None) -> str | None:
    """Resolve SDP workstream ID (PP-FFF-SS) to Beads task ID.

    Beads uses hash-based IDs (e.g., sdp-4qq); guard activate accepts ws_id (00-020-03).
    This function looks up the mapping in .beads-sdp-mapping.jsonl.

    Args:
        ws_id: Workstream ID (PP-FFF-SS or beads_id)
        mapping_file: Path to mapping file (default: .beads-sdp-mapping.jsonl)

    Returns:
        Beads task ID if found, else None (caller may use ws_id as-is for beads_id format)
    """
    path = mapping_file or Path.cwd() / ".beads-sdp-mapping.jsonl"
    if not path.exists():
        return None
    try:
        with open(path, "r") as f:
            for line in f:
                if not line.strip():
                    continue
                entry = json.loads(line)
                if entry.get("sdp_id") == ws_id:
                    beads_id = entry.get("beads_id")
                    return beads_id if isinstance(beads_id, str) else None
    except (json.JSONDecodeError, IOError):
        pass
    return None


class MappingManager:
    """Manages bidirectional ID mapping between SDP and Beads."""

    def __init__(self, mapping_file: Path):
        """Initialize mapping manager.

        Args:
            mapping_file: Path to ID mapping table (JSONL format)
        """
        self.mapping_file = mapping_file
        self._mapping: Dict[str, str] = {}  # sdp_id → beads_id
        self._reverse_mapping: Dict[str, str] = {}  # beads_id → sdp_id

    def load(self) -> None:
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

    def save(self) -> None:
        """Save ID mapping to file (overwrite with current state)."""
        try:
            with open(self.mapping_file, "w") as f:
                for sdp_id, beads_id in self._mapping.items():
                    entry = {
                        "sdp_id": sdp_id,
                        "beads_id": beads_id,
                        "updated_at": datetime.utcnow().isoformat(),
                    }
                    f.write(json.dumps(entry) + "\n")

        except IOError as e:
            raise BeadsSyncError(f"Failed to save mapping file: {e}") from e

    def get_beads_id(self, sdp_id: str) -> str | None:
        """Get Beads ID for given SDP workstream ID."""
        return self._mapping.get(sdp_id)

    def get_sdp_id(self, beads_id: str) -> str | None:
        """Get SDP workstream ID for given Beads ID."""
        return self._reverse_mapping.get(beads_id)

    def add_mapping(self, sdp_id: str, beads_id: str) -> None:
        """Add a new ID mapping."""
        self._mapping[sdp_id] = beads_id
        self._reverse_mapping[beads_id] = sdp_id
