"""Tier promotion check command.

Provides CLI command for checking tier promotion readiness.
"""

import sys
from pathlib import Path

import click

from sdp.errors import format_error_for_terminal


def _get_next_tier(current_tier: str) -> str | None:
    """Get next tier in promotion order. Returns None if at T3."""
    tier_order = ["T0", "T1", "T2", "T3"]
    idx = tier_order.index(current_tier)
    if idx >= len(tier_order) - 1:
        return None
    return tier_order[idx + 1]


def _show_validation_result(result) -> None:
    """Display validation result to user."""
    click.echo()
    click.echo(f"Validation Result: {result.passed}")
    passed = sum(1 for c in result.checks if c.passed)
    click.echo(f"Checks Passed: {passed}/{len(result.checks)}")
    failed_checks = [c for c in result.checks if not c.passed]
    if failed_checks:
        click.echo("\nFailed Checks:")
        for check in failed_checks:
            click.echo(f"  ✗ {check.name}: {check.message}")
            if check.details:
                for detail in check.details:
                    click.echo(f"      - {detail}")


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
    metrics = store.get_metrics(ws_id)
    if not metrics:
        click.echo(f"No metrics found for {ws_id}")
        sys.exit(1)

    current_tier = metrics.current_tier
    click.echo(f"Current Tier: {current_tier}")

    next_tier = _get_next_tier(current_tier.value)
    if not next_tier:
        click.echo("Already at highest tier (T3)")
        sys.exit(0)

    click.echo(f"Checking promotion to: {next_tier}")
    ws_file = metrics.metadata.get("file_path")
    if not ws_file:
        click.echo("Error: Workstream file path not in metrics")
        sys.exit(1)

    try:
        result = validate_workstream_tier(Path(ws_file), next_tier)
    except Exception as e:
        click.echo(format_error_for_terminal(e), err=True)
        sys.exit(1)

    _show_validation_result(result)
    if result.passed and not dry_run:
        store.record_promotion(ws_id, current_tier.value, next_tier, "Validation passed")
        click.echo(f"\n✓ Promoted to {next_tier}!")
    elif result.passed and dry_run:
        click.echo(f"\n✓ Would be promoted to {next_tier} (dry-run)")
    else:
        click.echo(f"\n✗ Not ready for {next_tier}")
    sys.exit(0 if result.passed else 1)
