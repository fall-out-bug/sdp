"""Beads sync CLI commands.

Provides commands for synchronizing local state with Beads.
"""


import click

from sdp.beads import create_beads_client
from sdp.beads.sync_service import BeadsSyncService, SyncSource


@click.group()
def sync() -> None:
    """Beads synchronization commands."""
    pass


@sync.command("check")
@click.option(
    "--ws-id",
    default=None,
    help="Active workstream ID (reads from .guard_state if not provided)",
)
def check_sync(ws_id: str | None) -> None:
    """Check sync status between local and Beads.

    Args:
        ws_id: Optional active workstream ID
    """
    service = BeadsSyncService(create_beads_client())

    # Get active WS from local state if not provided
    if ws_id is None:
        ws_id = service._get_local_active_ws()

    result = service.check_sync(ws_id)

    if result.synced:
        click.echo("✅ Local and Beads are in sync")
        return

    click.echo("❌ Sync conflicts detected:")
    for conflict in result.conflicts:
        click.echo(
            f"  - {conflict.ws_id}: local={conflict.local_status}, "
            f"beads={conflict.beads_status} (field: {conflict.field})"
        )


@sync.command("run")
@click.option(
    "--source",
    type=click.Choice(["beads", "local"], case_sensitive=False),
    default="beads",
    help="Source of truth for conflict resolution",
)
@click.option(
    "--ws-id",
    default=None,
    help="Active workstream ID (reads from .guard_state if not provided)",
)
def run_sync(source: str, ws_id: str | None) -> None:
    """Sync local state with Beads.

    Args:
        source: Which side is source of truth ('beads' or 'local')
        ws_id: Optional active workstream ID
    """
    sync_source = SyncSource(source.lower())
    service = BeadsSyncService(create_beads_client())

    # Get active WS from local state if not provided
    if ws_id is None:
        ws_id = service._get_local_active_ws()

    result = service.sync(ws_id, sync_source)

    if not result.changes:
        click.echo("✅ Already in sync")
        return

    click.echo(f"✅ Synced from {source}:")
    for change in result.changes:
        click.echo(f"  - {change}")
