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

from .mapping import BeadsSyncError, MappingManager, resolve_ws_id_to_beads_id
from .status_mapper import (
    map_beads_status_to_sdp,
    map_sdp_size_to_beads_priority,
    map_sdp_status_to_beads,
)
from .sync_service import BeadsSyncService

__all__ = [
    "BeadsSyncError",
    "resolve_ws_id_to_beads_id",
    "MappingManager",
    "map_sdp_status_to_beads",
    "map_beads_status_to_sdp",
    "map_sdp_size_to_beads_priority",
    "BeadsSyncService",
]
