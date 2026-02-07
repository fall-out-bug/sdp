"""
Beads migration commands.

Convert existing SDP markdown workstreams to Beads tasks.
"""

from dataclasses import asdict
from pathlib import Path

import click

from ..beads.client import create_beads_client
from ..beads.sync import BeadsSyncService


@click.group()
def beads() -> None:
    """Beads integration commands."""
    pass


@beads.command()
@click.argument(
    "workstreams_dir",
    type=click.Path(exists=True),
)
@click.option(
    "--real",
    "use_real",
    is_flag=True,
    default=False,
    help="Use real Beads CLI (default: mock)",
)
@click.option(
    "--use-mock",
    "use_mock",
    is_flag=True,
    default=None,
    help="Use mock Beads client (deprecated: use --real for real Beads)",
)
def migrate(workstreams_dir: Path, use_real: bool, use_mock: bool | None) -> None:
    """Migrate markdown workstreams to Beads tasks.

    Reads all markdown workstream files and converts them to Beads tasks.

    Example:
        sdp beads migrate docs/workstreams/backlog/

    This creates:
    - Beads tasks for each workstream
    - ID mapping in .beads-sdp-mapping.jsonl
    - Preserves dependencies between workstreams
    """

    click.echo(f"🔄 Migrating workstreams from {workstreams_dir}")

    # Initialize client (--real = real Beads, else mock)
    force_mock = use_mock if use_mock is not None else (not use_real)
    client = create_beads_client(use_mock=force_mock)
    sync = BeadsSyncService(client)

    # Find all workstream markdown files (exclude feature overviews, epics)
    all_files = list(Path(workstreams_dir).rglob("*.md"))
    skip_patterns = ("00-032-00-", "BEADS-001-")  # Feature overview, Epic (no ws_id)
    ws_files = [
        f for f in all_files
        if not any(f.name.startswith(p) for p in skip_patterns)
    ]

    if not ws_files:
        click.echo("⚠️  No workstream files found")
        return

    click.echo(f"Found {len(ws_files)} workstream files")

    # Migrate each workstream
    success = 0
    failed = 0

    for ws_file in ws_files:
        click.echo(f"\n📄 Processing {ws_file.name}")

        # Parse workstream file
        try:
            from ..core.workstream import parse_workstream

            ws = parse_workstream(ws_file)

            # Convert to Beads
            result = sync.sync_workstream_to_beads(ws_file, asdict(ws))

            if result.success:
                click.echo(f"  ✅ {ws.ws_id} → {result.beads_id}")
                success += 1
            else:
                click.echo(f"  ❌ {ws.ws_id}: {result.error}")
                failed += 1

        except Exception as e:
            click.echo(f"  ❌ {ws_file.name}: {e}")
            failed += 1

    # Persist deduplicated mapping (fixes legacy append-duplicates)
    sync.persist_mapping()

    # Summary
    click.echo(f"\n{'='*60}")
    click.echo("Migration complete!")
    click.echo(f"{'='*60}")
    click.echo(f"Total: {len(ws_files)}")
    click.echo(f"Success: {success} ✅")
    click.echo(f"Failed: {failed} ❌")

    if failed == 0:
        click.echo("\n✅ All workstreams migrated successfully!")
        click.echo("\nNext steps:")
        click.echo("  1. Verify: bd list")
        click.echo("  2. Check ready: bd ready")
        click.echo("  3. Start execution: @build <task_id>")
    else:
        click.echo("\n⚠️  Some workstreams failed to migrate")
        click.echo("   Fix issues and run migration again")


@beads.command()
@click.option(
    "--format",
    type=click.Choice(["table", "json"]),
    default="table",
    help="Output format",
)
def status(format: str) -> None:
    """Show Beads integration status.

    Displays information about Beads setup and migration state.
    """
    import os

    # Check if mock or real
    use_mock = os.getenv("BEADS_USE_MOCK", "true").lower() == "true"

    if format == "json":
        import json

        status_info = {
            "client_type": "mock" if use_mock else "real",
            "mapping_file": ".beads-sdp-mapping.jsonl",
            "workstreams_migrated": _count_migrated_workstreams(),
        }

        click.echo(json.dumps(status_info, indent=2))
    else:
        click.echo("Beads Integration Status")
        click.echo("=" * 40)
        click.echo(f"Client: {'Mock (dev)' if use_mock else 'Real (Beads CLI)'}")
        click.echo("Mapping: .beads-sdp-mapping.jsonl")

        migrated = _count_migrated_workstreams()
        click.echo(f"Migrated: {migrated} workstreams")


def _count_migrated_workstreams() -> int:
    """Count migrated workstreams from mapping file."""
    mapping_file = Path(".beads-sdp-mapping.jsonl")

    if not mapping_file.exists():
        return 0

    count = 0
    try:
        with open(mapping_file, "r") as f:
            for line in f:
                if line.strip():
                    count += 1
    except FileNotFoundError:
        # File doesn't exist yet
        pass

    return count
