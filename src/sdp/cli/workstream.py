"""Workstream parsing commands.

Provides CLI commands for parsing workstreams and project maps.
"""

import sys
from pathlib import Path
from typing import TYPE_CHECKING

import click

from sdp.errors import format_error_for_terminal

if TYPE_CHECKING:
    from click import Command

# Import validate_tier if available
validate_tier: Command | None = None
try:
    from sdp.cli.tier import validate_tier
except ImportError:
    validate_tier = None


@click.group()
def workstream() -> None:
    """Core SDP operations (workstreams, features, project maps)."""
    pass


@workstream.command("parse")
@click.argument("ws_file", type=click.Path(exists=True, path_type=Path))
def parse_workstream_cmd(ws_file: Path) -> None:
    """Parse a workstream markdown file.

    Args:
        ws_file: Path to workstream markdown file
    """
    from sdp.core import WorkstreamParseError, parse_workstream

    try:
        ws = parse_workstream(ws_file)
        click.echo(f"✓ Parsed {ws.ws_id}: {ws.title}")
        click.echo(f"  Feature: {ws.feature}")
        click.echo(f"  Status: {ws.status.value}")
        click.echo(f"  Size: {ws.size.value}")
        if ws.acceptance_criteria:
            click.echo(f"  Acceptance Criteria: {len(ws.acceptance_criteria)}")
    except WorkstreamParseError as e:
        click.echo(format_error_for_terminal(e), err=True)
        sys.exit(1)


@workstream.command("parse-project-map")
@click.argument("project_map_file", type=click.Path(exists=True, path_type=Path))
def parse_project_map(project_map_file: Path) -> None:
    """Parse a PROJECT_MAP.md file.

    Args:
        project_map_file: Path to PROJECT_MAP.md file
    """
    from sdp.core import ProjectMapParseError, parse_project_map

    try:
        pm = parse_project_map(project_map_file)
        click.echo(f"✓ Parsed project map: {pm.project_name}")
        click.echo(f"  Decisions: {len(pm.decisions)}")
        click.echo(f"  Constraints: {len(pm.constraints)}")
        if pm.tech_stack:
            click.echo(f"  Tech Stack Items: {len(pm.tech_stack)}")
    except ProjectMapParseError as e:
        click.echo(format_error_for_terminal(e), err=True)
        sys.exit(1)


@workstream.group("scope")
def scope() -> None:
    """Manage workstream file scope."""
    pass


@scope.command("show")
@click.argument("ws_id")
def show_scope(ws_id: str) -> None:
    """Show scope files for workstream.

    Args:
        ws_id: Workstream ID (e.g., bd-0001 or 00-032-15)
    """
    from sdp.beads import create_beads_client
    from sdp.beads.scope_manager import ScopeManager

    manager = ScopeManager(create_beads_client())

    try:
        scope_files = manager.get_scope(ws_id)
    except ValueError as e:
        click.echo(f"❌ {e}", err=True)
        sys.exit(1)

    if not scope_files:
        click.echo(f"Scope for {ws_id}: unrestricted (all files allowed)")
        return

    click.echo(f"Scope for {ws_id}: {len(scope_files)} files")
    for file_path in scope_files:
        click.echo(f"  - {file_path}")


@scope.command("add")
@click.argument("ws_id")
@click.argument("file_path")
def add_to_scope(ws_id: str, file_path: str) -> None:
    """Add file to workstream scope.

    Args:
        ws_id: Workstream ID
        file_path: File path to add
    """
    from sdp.beads import create_beads_client
    from sdp.beads.scope_manager import ScopeManager

    manager = ScopeManager(create_beads_client())

    try:
        manager.add_file(ws_id, file_path)
        click.echo(f"✅ Added {file_path} to {ws_id} scope")
    except ValueError as e:
        click.echo(f"❌ {e}", err=True)
        sys.exit(1)


@scope.command("remove")
@click.argument("ws_id")
@click.argument("file_path")
def remove_from_scope(ws_id: str, file_path: str) -> None:
    """Remove file from workstream scope.

    Args:
        ws_id: Workstream ID
        file_path: File path to remove
    """
    from sdp.beads import create_beads_client
    from sdp.beads.scope_manager import ScopeManager

    manager = ScopeManager(create_beads_client())

    try:
        manager.remove_file(ws_id, file_path)
        click.echo(f"✅ Removed {file_path} from {ws_id} scope")
    except ValueError as e:
        click.echo(f"❌ {e}", err=True)
        sys.exit(1)


@scope.command("clear")
@click.argument("ws_id")
def clear_scope(ws_id: str) -> None:
    """Clear scope (make unrestricted).

    Args:
        ws_id: Workstream ID
    """
    from sdp.beads import create_beads_client
    from sdp.beads.scope_manager import ScopeManager

    manager = ScopeManager(create_beads_client())

    try:
        manager.clear_scope(ws_id)
        click.echo(f"✅ Cleared scope for {ws_id} (now unrestricted)")
    except ValueError as e:
        click.echo(f"❌ {e}", err=True)
        sys.exit(1)


@workstream.command("verify")
@click.argument("ws_id")
def verify_completion(ws_id: str) -> None:
    """Verify workstream completion with evidence.

    Checks:
    - All scope_files output exist
    - All Verification commands pass
    - Test coverage meets threshold

    Args:
        ws_id: Workstream ID (e.g., 00-032-26)
    """
    from sdp.validators.ws_completion import WSCompletionVerifier

    verifier = WSCompletionVerifier()
    result = verifier.verify(ws_id)

    # Print results
    if result.passed:
        click.echo(f"✅ Workstream {ws_id} verification PASSED")
    else:
        click.echo(f"❌ Workstream {ws_id} verification FAILED")

    click.echo(f"\nChecks run: {len(result.checks)}")

    for check in result.checks:
        status = "✅" if check.passed else "❌"
        click.echo(f"  {status} {check.name}: {check.message}")

    if result.coverage_actual is not None:
        click.echo(f"\nCoverage: {result.coverage_actual:.1f}%")

    if result.missing_files:
        click.echo(f"\nMissing files ({len(result.missing_files)}):")
        for f in result.missing_files:
            click.echo(f"  - {f}")

    if result.failed_commands:
        click.echo(f"\nFailed commands ({len(result.failed_commands)}):")
        for cmd in result.failed_commands:
            click.echo(f"  - {cmd}")

    sys.exit(0 if result.passed else 1)


@workstream.command("supersede")
@click.argument("old_ws")
@click.option("--replacement", required=True, help="Replacement workstream ID")
def supersede_ws(old_ws: str, replacement: str) -> None:
    """Mark workstream as superseded by another.

    Args:
        old_ws: Workstream to supersede
        replacement: Replacement workstream ID
    """
    from sdp.validators.supersede_checker import SupersedeValidator

    validator = SupersedeValidator()
    result = validator.supersede(old_ws, replacement)

    if result.success:
        click.echo(f"✅ Marked {old_ws} as superseded by {replacement}")
    else:
        click.echo(f"❌ Failed to supersede {old_ws}: {result.error}")
        sys.exit(1)


@workstream.command("orphans")
def find_orphans() -> None:
    """Find superseded workstreams without valid replacement."""
    from sdp.validators.supersede_checker import SupersedeValidator

    validator = SupersedeValidator()
    orphans = validator.find_orphans()

    if not orphans:
        click.echo("✅ No orphaned superseded workstreams found")
        return

    click.echo(f"❌ Found {len(orphans)} orphaned superseded workstream(s):")
    for ws_id in orphans:
        click.echo(f"  - {ws_id}")

    click.echo("\n💡 Fix with: sdp ws supersede <ws_id> --replacement <new_ws_id>")
    sys.exit(1)


# Add validate_tier command if available (as 'validate' for compatibility)
if validate_tier:
    workstream.add_command(validate_tier, name='validate')
