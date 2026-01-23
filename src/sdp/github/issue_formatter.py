"""Format GitHub issue body from WS metadata."""

from sdp.github.ws_parser import WSMetadata


def format_issue_body(ws: WSMetadata, ws_file_path: str) -> str:
    """Format issue body from WS metadata.

    Args:
        ws: Parsed WS metadata
        ws_file_path: Relative path to WS file (for link)

    Returns:
        Formatted markdown body for GitHub issue
    """
    lines = []

    # Goal section
    lines.append("## 🎯 Goal")
    lines.append("")
    lines.append(ws.goal)
    lines.append("")

    # Acceptance Criteria section
    lines.append("## Acceptance Criteria")
    lines.append("")
    for ac in ws.acceptance_criteria:
        lines.append(f"- [ ] {ac}")
    lines.append("")

    # Dependencies section (if any)
    if ws.dependencies:
        lines.append("## Dependencies")
        lines.append("")
        for dep in ws.dependencies:
            lines.append(f"- {dep}")
        lines.append("")

    # WS file link
    lines.append("## Workstream File")
    lines.append("")
    lines.append(f"- [View WS]({ws_file_path})")
    lines.append("")

    # Metadata section
    lines.append("## Metadata")
    lines.append("")
    lines.append(f"- **WS ID:** {ws.ws_id}")
    lines.append(f"- **Feature:** {ws.feature}")
    lines.append(f"- **Size:** {ws.size}")
    lines.append(f"- **Status:** {ws.status}")
    lines.append("")

    return "\n".join(lines)
