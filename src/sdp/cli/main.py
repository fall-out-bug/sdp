"""Main CLI entry point for SDP package."""

import sys
from pathlib import Path

import click

from sdp import __version__
from sdp.errors import format_error_for_terminal

# Import Beads commands (optional - may not be available in all builds)
try:
    from sdp.cli.beads import beads
    _beads_available = True
except ImportError:
    _beads_available = False

# Import doctor command
try:
    from sdp.doctor import doctor
    _doctor_available = True
except ImportError:
    _doctor_available = False


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
    - Beads integration (if installed)
    """
    pass


# Add Beads commands if available
if _beads_available:
    main.add_command(beads)

# Add doctor command if available
if _doctor_available:
    main.add_command(doctor)


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
        click.echo(format_error_for_terminal(e), err=True)
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


@main.group()
def tier() -> None:
    """Tier management commands (metrics, promotion, demotion)."""
    pass


@tier.command("metrics")
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
        click.echo(f"Successful: {metrics.successful_attempts}")
        click.echo(f"Success Rate: {metrics.success_rate:.1%}")
        click.echo(f"Consecutive Failures: {metrics.consecutive_failures}")
        click.echo(f"Last Updated: {metrics.last_updated.strftime('%Y-%m-%d %H:%M:%S')}")
    else:
        # Show all workstreams
        all_metrics = store._metrics
        if not all_metrics:
            click.echo("No metrics found")
            sys.exit(0)

        click.echo(f"=== Tier Metrics ({len(all_metrics)} workstreams) ===")
        click.echo()

        for ws_id, metrics in sorted(all_metrics.items()):
            click.echo(f"{ws_id}:")
            click.echo(f"  Tier: {metrics.current_tier}")
            click.echo(
                f"  Success: {metrics.successful_attempts}/{metrics.total_attempts} "
                f"({metrics.success_rate:.1%})"
            )
            click.echo(f"  Consecutive Failures: {metrics.consecutive_failures}")


@tier.command("promote-check")
@click.argument("ws_id", default="")
@click.option(
    "--storage",
    type=click.Path(path_type=Path),
    default=Path(".sdp/tier_metrics.json"),
    help="Path to metrics storage file",
)
@click.option(
    "--dry-run",
    is_flag=True,
    help="Check without updating workstream files",
)
def tier_promote_check(ws_id: str, storage: Path, dry_run: bool) -> None:
    """Check workstream(s) for tier promotion/demotion eligibility.

    Args:
        ws_id: Workstream ID (empty for all workstreams)
        storage: Path to metrics storage file
        dry_run: Check without updating files
    """
    from sdp.core.tier_metrics import TierMetricsStore

    store = TierMetricsStore(storage)

    if ws_id:
        ws_ids = [ws_id]
    else:
        ws_ids = list(store._metrics.keys())

    if not ws_ids:
        click.echo("No metrics found")
        sys.exit(0)

    click.echo(f"=== Promotion/Demotion Check ({len(ws_ids)} workstreams) ===")
    click.echo()

    changes = []
    for ws_id in ws_ids:
        new_tier = store.check_promotion_eligible(ws_id)
        if new_tier:
            metrics = store.get_metrics(ws_id)
            if metrics is None:
                continue
            changes.append((ws_id, metrics.current_tier, new_tier))
            click.echo(
                f"⚠ {ws_id}: {metrics.current_tier} → {new_tier} "
                f"({metrics.successful_attempts}/{metrics.total_attempts} attempts, "
                f"{metrics.success_rate:.1%} success)"
            )

    if not changes:
        click.echo("No tier changes needed")
    elif not dry_run:
        click.echo()
        click.echo(f"Found {len(changes)} tier changes")
        click.echo("Note: Automatic file updates not yet implemented")
        click.echo("Use --dry-run to preview changes")


@main.group()
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
            f"⚠️ ALERT: High escalation rate ({escalation_rate:.1%} > {alert_threshold:.1%})"
        )
        click.echo("  Consider reviewing workstream quality or tier assignments")


@main.group()
def prd() -> None:
    """PRD (Product Requirements Document) operations."""
    pass


@prd.command("validate")
@click.argument(
    "prd_file",
    type=click.Path(exists=True, path_type=Path),
)
@click.option(
    "--exit-code-on-error",
    is_flag=True,
    help="Exit with code 1 if validation errors found",
)
def prd_validate(prd_file: Path, exit_code_on_error: bool) -> None:
    """Validate a PRD document against section limits.

    Args:
        prd_file: Path to PRD file (PROJECT_MAP.md)
        exit_code_on_error: Exit with code 1 if errors found
    """
    from sdp.prd.validator import (
        format_validation_issues,
        has_critical_issues,
        validate_prd_file,
    )

    issues = validate_prd_file(prd_file)

    if issues:
        click.echo(format_validation_issues(issues))
        if has_critical_issues(issues):
            if exit_code_on_error:
                sys.exit(1)
    else:
        click.echo("✅ PRD validation passed")


@prd.command("detect-type")
@click.argument(
    "project_path",
    type=click.Path(exists=True, path_type=Path),
)
def prd_detect_type(project_path: Path) -> None:
    """Detect project type from file structure.

    Args:
        project_path: Path to project root
    """
    from sdp.prd.detector import detect_project_type

    project_type = detect_project_type(project_path)
    click.echo(f"Detected project type: {project_type.value}")


# Register extension commands
from sdp.cli_extension import extension
from sdp.cli_init import init

main.add_command(extension)
main.add_command(init)
main.add_command(prd)


if __name__ == "__main__":
    main()
