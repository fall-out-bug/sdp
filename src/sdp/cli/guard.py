"""CLI commands for guard operations."""

import sys
from datetime import datetime, timezone

import click

from sdp.beads import create_beads_client
from sdp.beads.sync import resolve_ws_id_to_beads_id
from sdp.guard.skill import GuardSkill
from sdp.guard.state import GuardState, StateManager


@click.group()
def guard() -> None:
    """Pre-edit guard commands."""
    pass


@guard.command("activate")
@click.argument("ws_id")
def activate(ws_id: str) -> None:
    """Activate workstream for editing.

    Args:
        ws_id: Workstream ID (PP-FFF-SS) or Beads task ID
    """
    client = create_beads_client()
    guard_skill = GuardSkill(client)

    # Resolve ws_id (00-020-03) to beads_id (sdp-4qq) if needed
    task_id = resolve_ws_id_to_beads_id(ws_id) or ws_id

    try:
        guard_skill.activate(task_id)
    except ValueError as e:
        click.echo(f"❌ {e}")
        sys.exit(1)

    ws = client.get_task(task_id)
    scope = ws.sdp_metadata.get("scope_files", []) if ws else []

    state = GuardState(
        active_ws=ws_id,
        activated_at=datetime.now(timezone.utc).isoformat(),
        scope_files=scope,
    )
    StateManager.save(state)

    click.echo(f"✅ Activated WS: {ws_id}")
    if scope:
        click.echo(f"   Scope: {len(scope)} files")
    else:
        click.echo("   Scope: unrestricted")


@guard.command("check")
@click.argument("file_path", type=click.Path())
def check_file(file_path: str) -> None:
    """Check if file edit is allowed.

    Args:
        file_path: Path to file to check
    """
    state = StateManager.load()

    if not state.active_ws:
        click.echo("❌ No active WS. Run: sdp guard activate <ws_id>")
        sys.exit(1)

    # Resolve ws_id (00-020-03) to beads_id (sdp-4qq) for get_task
    task_id = resolve_ws_id_to_beads_id(state.active_ws) or state.active_ws

    client = create_beads_client()
    guard_skill = GuardSkill(client)
    guard_skill._active_ws = task_id

    result = guard_skill.check_edit(file_path)

    if not result.allowed:
        click.echo(f"❌ {result.reason}")
        if result.scope_files:
            click.echo("   Allowed files:")
            for f in result.scope_files[:10]:
                click.echo(f"     - {f}")
        sys.exit(1)

    click.echo(f"✅ Edit allowed: {file_path}")


@guard.command("status")
def status() -> None:
    """Show current guard status."""
    state = StateManager.load()

    if not state.active_ws:
        click.echo("No active workstream")
        return

    click.echo(f"Active WS: {state.active_ws}")
    click.echo(f"Activated: {state.activated_at}")
    if state.scope_files:
        click.echo(f"Scope: {len(state.scope_files)} files")
    else:
        click.echo("Scope: unrestricted")


@guard.command("deactivate")
def deactivate() -> None:
    """Deactivate current workstream."""
    StateManager.clear()
    click.echo("✅ Guard deactivated")


@guard.command("current")
def current() -> None:
    """Show currently active workstream (AC4)."""
    state = StateManager.load()

    if not state.active_ws:
        click.echo("No active workstream")
        return

    click.echo(f"Active WS: {state.active_ws}")
