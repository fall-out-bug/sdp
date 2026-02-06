"""
SDP Doctor - Health check command.

Performs diagnostic checks on SDP installation and configuration.
"""

import json
import sys
from typing import List

import click

from sdp.health_checks import HealthCheck, HealthCheckResult, get_health_checks


@click.command()
@click.option(
    "--format",
    type=click.Choice(["table", "json"]),
    default="table",
    help="Output format",
)
def doctor(format: str) -> None:
    """Run SDP health checks.

    Performs diagnostic checks on:
    - Python version (>= 3.10)
    - Poetry installation
    - Git hooks configuration
    - Beads CLI (optional)
    - GitHub CLI (optional)
    - Telegram configuration (optional)

    Exit codes:
    - 0: All critical checks passed
    - 1: One or more critical checks failed

    Example:
        sdp doctor
        sdp doctor --format json
    """
    checks: List[HealthCheck] = get_health_checks()
    results = []

    for check in checks:
        try:
            result = check.run()
            results.append(result)
        except Exception as e:
            # If check crashes, treat as failure
            results.append(
                HealthCheckResult(
                    name=check.name,
                    passed=False,
                    message=f"Check failed with error: {e}",
                    remediation="Review error details above",
                )
            )

    # Determine if all critical checks passed
    critical_results = [r for r in results if any(c.name == r.name and c.critical for c in checks)]
    all_passed = all(r.passed for r in critical_results)

    if format == "json":
        # JSON output (machine-readable)
        output = {
            "all_passed": all_passed,
            "checks": [
                {
                    "name": r.name,
                    "passed": r.passed,
                    "message": r.message,
                    "remediation": r.remediation,
                }
                for r in results
            ],
        }
        click.echo(json.dumps(output, indent=2))
    else:
        # Table output (human-readable)
        click.echo("SDP Health Check")
        click.echo("=" * 60)
        click.echo()

        for result in results:
            status = "✓" if result.passed else "✗"
            click.echo(f"{status} {result.name}: {result.message}")

            if result.remediation:
                click.echo(f"  → {result.remediation}")

            click.echo()

        click.echo("=" * 60)

        if all_passed:
            click.echo("✅ All critical checks passed!")
        else:
            failed = sum(1 for r in critical_results if not r.passed)
            click.echo(f"❌ {failed} critical check(s) failed")

    # Exit with appropriate code
    sys.exit(0 if all_passed else 1)
