"""Tier validation and management commands.

Provides CLI commands for validating workstreams against capability tiers
and managing tier metrics.
"""

import json
import sys
from pathlib import Path

import click

from sdp.errors import format_error_for_terminal


@click.command()
@click.argument("ws_file", type=click.Path(exists=True, path_type=Path))
@click.option(
    "--tier",
    type=click.Choice(["T0", "T1", "T2", "T3"], case_sensitive=False),
    required=True,
    help="Capability tier to validate against (T0, T1, T2, T3)",
)
@click.option(
    "--json",
    "output_json",
    is_flag=True,
    help="Output results as JSON (machine-readable)",
)
def validate_tier(ws_file: Path, tier: str, output_json: bool) -> None:
    """Validate workstream against capability tier.

    Validates a workstream markdown file against the specified capability tier
    (T0-T3) according to Contract-Driven WS v2.0 specification.

    Args:
        ws_file: Path to workstream markdown file
        tier: Capability tier (T0, T1, T2, T3)
        output_json: Output results as JSON
    """
    from sdp.validators import validate_workstream_tier

    try:
        result = validate_workstream_tier(ws_file, tier)
    except ValueError as e:
        click.echo(format_error_for_terminal(e), err=True)
        sys.exit(1)
    except Exception as e:
        click.echo(format_error_for_terminal(e), err=True)
        sys.exit(1)

    if output_json:
        # Machine-readable JSON output
        output = {
            "tier": result.tier.value,
            "passed": result.passed,
            "checks": [
                {
                    "name": check.name,
                    "passed": check.passed,
                    "message": check.message,
                    "details": check.details,
                }
                for check in result.checks
            ],
        }
        click.echo(json.dumps(output, indent=2))
        sys.exit(0 if result.passed else 1)
    else:
        # Human-readable output
        click.echo(f"=== Capability Tier Validation ({result.tier.value}) ===")
        click.echo(f"Workstream: {ws_file}")
        click.echo()

        for check in result.checks:
            status = "✓" if check.passed else "✗"
            click.echo(f"{status} {check.name}: {check.message}")
            if check.details:
                for detail in check.details:
                    click.echo(f"    - {detail}")

        click.echo()
        if result.passed:
            click.echo(f"Result: {result.tier.value}-READY ✓")
            sys.exit(0)
        else:
            click.echo(f"Result: {result.tier.value}-READY ✗")
            failed_count = sum(1 for check in result.checks if not check.passed)
            click.echo(f"Failed checks: {failed_count}/{len(result.checks)}")
            sys.exit(1)


@click.group()
def tier() -> None:
    """Tier management commands (validate, metrics, promotion)."""
    pass


# Add validate_tier as subcommand
tier.add_command(validate_tier, name="validate")
