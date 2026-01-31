"""Chain tracing utilities for supersede validation."""

from pathlib import Path

from sdp.validators.supersede.models import SupersedeChain
from sdp.validators.supersede.parser import find_ws_file, parse_frontmatter


def trace_chain(ws_id: str, ws_dir: Path) -> SupersedeChain:
    """Trace supersede chain to final WS.

    Args:
        ws_id: Starting workstream ID
        ws_dir: Base workstream directory

    Returns:
        SupersedeChain with full path
    """
    visited: set[str] = set()
    chain: list[str] = [ws_id]
    current = ws_id

    while True:
        if current in visited:
            # Cycle detected
            return SupersedeChain(
                original_ws=ws_id, replacements=chain, has_cycle=True, final_ws=None
            )

        visited.add(current)

        # Find current WS
        ws_path = find_ws_file(current, ws_dir)
        if not ws_path:
            # Dead end
            return SupersedeChain(
                original_ws=ws_id, replacements=chain, has_cycle=False, final_ws=None
            )

        # Check if superseded
        data = parse_frontmatter(ws_path)
        if data.get("status") != "superseded":
            # Found final WS
            return SupersedeChain(
                original_ws=ws_id, replacements=chain, has_cycle=False, final_ws=current
            )

        # Continue to replacement
        next_ws = data.get("superseded_by")
        if not next_ws:
            # Orphan
            return SupersedeChain(
                original_ws=ws_id, replacements=chain, has_cycle=False, final_ws=None
            )

        chain.append(next_ws)
        current = next_ws


def find_orphans(ws_dir: Path) -> list[str]:
    """Find superseded WS without valid replacement.

    Args:
        ws_dir: Directory to search

    Returns:
        List of orphaned WS IDs
    """
    orphans: list[str] = []

    for ws_file in ws_dir.rglob("*.md"):
        data = parse_frontmatter(ws_file)

        if data.get("status") == "superseded":
            superseded_by = data.get("superseded_by")

            if not superseded_by:
                # No replacement specified
                orphans.append(data.get("ws_id", ws_file.stem))
            else:
                # Check if replacement exists
                replacement_path = find_ws_file(superseded_by, ws_dir)
                if not replacement_path:
                    orphans.append(data.get("ws_id", ws_file.stem))

    return orphans
