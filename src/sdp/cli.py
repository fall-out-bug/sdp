"""Main CLI entry point for SDP package."""

import sys
from pathlib import Path

import click

from sdp import __version__


@click.group()
@click.version_option(version=__version__, prog_name="sdp")
def main() -> None:
    """SDP (Spec-Driven Protocol) - Workstream automation tools.

    This CLI provides commands for:
    - Workstream parsing and validation
    - Feature decomposition
    - Project map management
    - GitHub integration
    - Extension management
    """
    pass


@main.command()
def version() -> None:
    """Show SDP version."""
    click.echo(f"sdp version {__version__}")


@main.group()
def core() -> None:
    """Core SDP operations (workstreams, features, project maps)."""
    pass


@core.command("parse-ws")
@click.argument("ws_file", type=click.Path(exists=True, path_type=Path))
def parse_workstream(ws_file: Path) -> None:
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
        click.echo(f"Error parsing workstream: {e}", err=True)
        sys.exit(1)


@core.command("parse-project-map")
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
        click.echo(f"Error parsing project map: {e}", err=True)
        sys.exit(1)


@core.command("validate-tier")
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
    import json

    from sdp.validators import ValidationResult, validate_workstream_tier

    try:
        result = validate_workstream_tier(ws_file, tier)
    except ValueError as e:
        click.echo(f"Error: {e}", err=True)
        sys.exit(1)
    except Exception as e:
        click.echo(f"Unexpected error: {e}", err=True)
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


# Register extension commands
from sdp.cli_extension import extension
from sdp.cli_init import init

main.add_command(extension)
main.add_command(init)


if __name__ == "__main__":
    main()
