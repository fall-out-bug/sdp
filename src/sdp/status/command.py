"""Status command implementation."""

import sys
import time

import click

from sdp.dashboard.sources.workstream_reader import WorkstreamReader


@click.command()
@click.option(
    "--filter",
    "status_filter",
    type=click.Choice(["backlog", "in-progress", "completed"]),
    help="Filter by status",
)
@click.option("--feature", help="Filter by feature ID")
@click.option("--watch", is_flag=True, help="Auto-refresh every 2s")
def status(status_filter: str | None, feature: str | None, watch: bool) -> None:
    """Show workstream status.

    Displays all workstreams grouped by status, with optional filtering
    by status or feature ID.
    """
    reader = WorkstreamReader("docs/workstreams")

    while True:
        state = reader.read()
        _print_status(state, status_filter, feature)

        if not watch:
            break

        time.sleep(2)


def _print_status(
    state, status_filter: str | None, feature: str | None
) -> None:
    """Print status table.

    Args:
        state: DashboardState from WorkstreamReader
        status_filter: Optional status filter
        feature: Optional feature filter
    """
    # Group by status
    by_status: dict[str, list] = {}
    for ws in state.workstreams.values():
        if status_filter and ws.status != status_filter:
            continue
        if feature and ws.feature != feature:
            continue

        if ws.status not in by_status:
            by_status[ws.status] = []
        by_status[ws.status].append(ws)

    if not by_status:
        click.echo("No workstreams found")
        return

    # Print tables by status
    status_order = ["backlog", "in-progress", "completed"]
    for s in status_order:
        if s not in by_status:
            continue

        wss = by_status[s]
        click.echo(f"\n{s.replace('-', ' ').title()} ({len(wss)})")

        for ws in sorted(wss, key=lambda w: w.ws_id):
            assignee = f" [{ws.assignee}]" if ws.assignee else ""
            size = f" [{ws.size}]" if ws.size else ""
            click.echo(f"  {ws.ws_id}: {ws.title}{assignee}{size}")

    click.echo(f"\nTotal: {sum(len(v) for v in by_status.values())} workstream(s)")
