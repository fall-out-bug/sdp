"""Tier promotion check command.

Provides CLI command for checking tier promotion readiness.
"""

import sys
from pathlib import Path

import click

from sdp.errors import format_error_for_terminal


@click.command("promote-check")
@click.argument("ws_id")
@click.option(
    "--storage",
    type=click.Path(path_type=Path),
    default=Path(".sdp/tier_metrics.json"),
    help="Path to metrics storage file",
)
@click.option(
    "--dry-run",
    is_flag=True,
    help="Check without recording metrics",
)
def tier_promote_check(ws_id: str, storage: Path, dry_run: bool) -> None:
    """Check if workstream is ready for tier promotion.

    Args:
        ws_id: Workstream ID
        storage: Path to metrics storage file
        dry_run: Check without recording metrics
    """
    from sdp.core.tier_metrics import TierMetricsStore
    from sdp.validators.capability_tier import validate_workstream_tier

    store = TierMetricsStore(storage)

    # Get current metrics
    metrics = store.get_metrics(ws_id)
    if not metrics:
        click.echo(f"No metrics found for {ws_id}")
        sys.exit(1)

    current_tier = metrics.current_tier
    click.echo(f"Current Tier: {current_tier}")

    # Determine next tier
    tier_order = ["T0", "T1", "T2", "T3"]
    current_idx = tier_order.index(current_tier.value)

    if current_idx >= len(tier_order) - 1:
        click.echo("Already at highest tier (T3)")
        sys.exit(0)

    next_tier = tier_order[current_idx + 1]
    click.echo(f"Checking promotion to: {next_tier}")

    # Validate against next tier
    # Find workstream file from metrics
    ws_file = metrics.metadata.get("file_path")
    if not ws_file:
        click.echo("Error: Workstream file path not in metrics")
        sys.exit(1)

    try:
        result = validate_workstream_tier(Path(ws_file), next_tier)
    except Exception as e:
        click.echo(format_error_for_terminal(e), err=True)
        sys.exit(1)

    # Show results
    click.echo()
    click.echo(f"Validation Result: {result.passed}")
    click.echo(f"Checks Passed: {sum(1 for c in result.checks if c.passed)}/{len(result.checks)}")

    failed_checks = [c for c in result.checks if not c.passed]
    if failed_checks:
        click.echo("\nFailed Checks:")
        for check in failed_checks:
            click.echo(f"  ✗ {check.name}: {check.message}")
            if check.details:
                for detail in check.details:
                    click.echo(f"      - {detail}")

    # Record metrics if not dry run
    if result.passed and not dry_run:
        store.record_promotion(ws_id, current_tier.value, next_tier, "Validation passed")
        click.echo(f"\n✓ Promoted to {next_tier}!")
        sys.exit(0)
    elif result.passed and dry_run:
        click.echo(f"\n✓ Would be promoted to {next_tier} (dry-run)")
        sys.exit(0)
    else:
        click.echo(f"\n✗ Not ready for {next_tier}")
        sys.exit(1)
