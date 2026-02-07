"""Workstream file parser."""

from pathlib import Path
from typing import Any, Optional

from sdp.cli.status.models import WorkstreamSummary


def parse_ws_file(ws_file: Path) -> Optional[WorkstreamSummary]:
    """Parse workstream markdown file.

    Args:
        ws_file: Path to workstream file

    Returns:
        WorkstreamSummary if parsed successfully, None otherwise
    """
    try:
        content = ws_file.read_text()
        frontmatter, body = _extract_frontmatter(content)
        if frontmatter is None or body is None:
            return None

        metadata = _parse_metadata(frontmatter)
        blockers = _extract_blockers(frontmatter)
        title = _extract_title(body)

        ws_id = metadata.get("ws_id", "")
        if not ws_id:
            return None

        return WorkstreamSummary(
            id=ws_id,
            title=title,
            status=metadata.get("status", "UNKNOWN"),
            scope=metadata.get("complexity", "UNKNOWN"),
            blockers=blockers,
        )
    except Exception:
        return None


def _extract_frontmatter(content: str) -> tuple[Optional[str], Optional[str]]:
    """Extract frontmatter and body from markdown content."""
    if not content.startswith("---"):
        return None, None

    parts = content.split("---", 2)
    if len(parts) < 3:
        return None, None

    return parts[1].strip(), parts[2].strip()


def _parse_metadata(frontmatter: str) -> dict[str, Any]:
    """Parse YAML-like frontmatter into metadata dict."""
    metadata: dict[str, Any] = {}
    for line in frontmatter.split("\n"):
        if ":" in line:
            key, value = line.split(":", 1)
            key = key.strip()
            value = value.strip().strip('"')
            metadata[key] = value
    return metadata


def _extract_blockers(frontmatter: str) -> list[str]:
    """Extract dependency blockers from frontmatter."""
    blockers = []
    in_depends = False
    for line in frontmatter.split("\n"):
        if "depends_on:" in line:
            in_depends = True
            continue
        if in_depends:
            if line.strip().startswith("-"):
                blocker = line.strip().lstrip("- ").strip()
                if blocker:
                    blockers.append(blocker)
            elif line.strip() and not line.startswith(" "):
                break
    return blockers


def _extract_title(body: str) -> str:
    """Extract title from first heading in body."""
    for line in body.split("\n"):
        if line.startswith("# "):
            return line[2:].strip()
    return "Untitled"
