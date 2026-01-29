"""Tier metrics command.

Provides CLI command for viewing tier metrics.
"""

import sys
from pathlib import Path

import click


@click.command("metrics")
@click.argument("ws_id", default="")
@click.option(
    "--storage",
    type=click.Path(path_type=Path),
    default=Path(".sdp/tier_metrics.json"),
    help="Path to metrics storage file",
)
def tier_metrics(ws_id: str, storage: Path) -> None:
    """Show tier metrics for workstream(s).

    Args:
        ws_id: Workstream ID (empty for all workstreams)
        storage: Path to metrics storage file
    """
    from sdp.core.tier_metrics import TierMetricsStore

    store = TierMetricsStore(storage)

    if ws_id:
        # Show specific workstream
        metrics = store.get_metrics(ws_id)
        if not metrics:
            click.echo(f"No metrics found for {ws_id}")
            sys.exit(1)

        click.echo(f"=== Tier Metrics: {ws_id} ===")
        click.echo(f"Current Tier: {metrics.current_tier}")
        click.echo(f"Total Attempts: {metrics.total_attempts}")

        if metrics.history:
            click.echo("\nPromotion History:")
            for entry in metrics.history[-5:]:  # Show last 5
                click.echo(
                    f"  {entry.timestamp}: {entry.from_tier} → {entry.to_tier} "
                    f"({entry.reason})"
                )
    else:
        # Show all workstreams
        all_metrics = store.list_all()
        click.echo(f"=== Tier Metrics ({len(all_metrics)} workstreams) ===")
        for ws_id, metrics in all_metrics:
            click.echo(f"\n{ws_id}: {metrics.current_tier}")
            click.echo(f"  Attempts: {metrics.total_attempts}")
