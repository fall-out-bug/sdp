"""Metrics and monitoring commands.

Provides CLI commands for tracking escalation metrics and build statistics.
"""

from pathlib import Path

import click


@click.group()
def metrics() -> None:
    """Metrics and monitoring commands."""
    pass


@metrics.command("escalations")
@click.option(
    "--tier",
    type=click.Choice(["T2", "T3"]),
    help="Filter by capability tier",
)
@click.option(
    "--days",
    type=int,
    default=7,
    help="Time window in days (default: 7)",
)
@click.option(
    "--top",
    type=int,
    default=10,
    help="Show top N escalating workstreams",
)
@click.option(
    "--storage",
    type=click.Path(path_type=Path),
    default=Path(".sdp/escalation_metrics.json"),
    help="Path to escalation metrics storage",
)
@click.option(
    "--total-builds",
    type=int,
    default=20,
    help="Total builds in period for rate calculation",
)
def metrics_escalations(tier: str, days: int, top: int, storage: Path, total_builds: int) -> None:
    """Show escalation metrics and analysis.

    Args:
        tier: Filter by capability tier
        days: Time window in days
        top: Number of top workstreams to show
        storage: Path to metrics storage file
        total_builds: Total builds for rate calculation
    """
    from sdp.core.escalation_metrics import EscalationMetricsStore

    store = EscalationMetricsStore(storage)

    click.echo(f"=== Escalation Metrics (last {days} days) ===")
    click.echo()

    # Total escalations
    escalation_count = store.get_escalation_count(tier=tier, days=days)
    click.echo(f"Total Escalations: {escalation_count}")

    # Escalation rate
    escalation_rate = store.get_escalation_rate(tier=tier, days=days, total_builds=total_builds)
    click.echo(f"Escalation Rate: {escalation_rate:.1%} ({escalation_count}/{total_builds} builds)")

    # Average attempts
    avg_attempts = store.get_average_attempts(tier=tier, days=days)
    if avg_attempts > 0:
        click.echo(f"Avg Attempts Before Escalation: {avg_attempts:.1f}")

    # Top escalating workstreams
    top_ws = store.get_top_escalating_ws(limit=top, days=days)
    if top_ws:
        click.echo()
        click.echo(f"Top {len(top_ws)} Escalating Workstreams:")
        for ws_id, count in top_ws:
            click.echo(f"  {ws_id}: {count} escalations")

    # Alert if high escalation rate
    alert_threshold = 0.20  # 20%
    if escalation_rate > alert_threshold:
        click.echo()
        click.echo(
            f"⚠️  ALERT: High escalation rate ({escalation_rate:.1%} > {alert_threshold:.1%})"
        )
        click.echo("  Consider reviewing workstream quality or tier assignments")
