"""Status output formatters."""

import json
from dataclasses import asdict

from sdp.cli.status.models import BeadsStatus, GuardStatus, ProjectStatus, WorkstreamSummary


def format_status_human(status: ProjectStatus, verbose: bool = False) -> str:
    """Format status for human reading.

    Args:
        status: Project status to format
        verbose: Include additional details

    Returns:
        Formatted string output
    """
    lines = []

    # Header
    lines.append("=" * 50)
    lines.append("       SDP Project Status")
    lines.append("=" * 50)
    lines.append("")

    # Sections
    _add_in_progress_section(lines, status.in_progress, verbose)
    _add_ready_section(lines, status.ready, verbose)
    _add_blocked_section(lines, status.blocked)
    _add_guard_section(lines, status.guard, verbose)
    _add_beads_section(lines, status.beads, verbose)
    _add_next_actions_section(lines, status.next_actions)

    return "\n".join(lines)


def _add_in_progress_section(
    lines: list[str], in_progress: list[WorkstreamSummary], verbose: bool
) -> None:
    """Add in-progress section to output."""
    if in_progress:
        lines.append("⏳ In Progress")
        for ws in in_progress:
            lines.append(f"  • {ws.id}: {ws.title}")
            if verbose:
                lines.append(f"    Status: {ws.status}, Scope: {ws.scope}")
    else:
        lines.append("No workstreams in progress")
    lines.append("")


def _add_ready_section(
    lines: list[str], ready: list[WorkstreamSummary], verbose: bool
) -> None:
    """Add ready section to output."""
    if ready:
        lines.append("✅ Ready to Start")
        for ws in ready:
            lines.append(f"  • {ws.id}: {ws.title}")
            if verbose:
                lines.append(f"    Scope: {ws.scope}")
    lines.append("")


def _add_blocked_section(lines: list[str], blocked: list[WorkstreamSummary]) -> None:
    """Add blocked section to output."""
    if blocked:
        lines.append("🚫 Blocked")
        for ws in blocked:
            blockers = ", ".join(ws.blockers) if ws.blockers else "unknown"
            lines.append(f"  • {ws.id}: {ws.title} (by: {blockers})")
    lines.append("")


def _add_guard_section(lines: list[str], guard: GuardStatus, verbose: bool) -> None:
    """Add guard status section to output."""
    lines.append("🛡️  Guard")
    if guard.active:
        lines.append(f"  Active: {guard.workstream_id}")
        if verbose and guard.allowed_files:
            lines.append(f"  Allowed files: {len(guard.allowed_files)}")
    else:
        lines.append("  Inactive")
    lines.append("")


def _add_beads_section(lines: list[str], beads: BeadsStatus, verbose: bool) -> None:
    """Add Beads status section to output."""
    if beads.available:
        lines.append("📿 Beads")
        sync_status = "✅ Synced" if beads.synced else "⚠️  Needs sync"
        lines.append(f"  {sync_status}")
        if beads.ready_tasks:
            lines.append(f"  Ready: {len(beads.ready_tasks)} tasks")
        if verbose and beads.last_sync:
            lines.append(f"  Last sync: {beads.last_sync}")
        lines.append("")


def _add_next_actions_section(lines: list[str], next_actions: list[str]) -> None:
    """Add next actions section to output."""
    if next_actions:
        lines.append("💡 Suggested Actions")
        for action in next_actions:
            lines.append(f"  → {action}")
        lines.append("")


def format_status_json(status: ProjectStatus) -> str:
    """Format status as JSON.

    Args:
        status: Project status to format

    Returns:
        JSON string
    """
    return json.dumps(asdict(status), indent=2)
