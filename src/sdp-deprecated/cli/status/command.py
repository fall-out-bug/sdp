"""SDP status command."""

from pathlib import Path

import click

from sdp.cli.status.collector import StatusCollector
from sdp.cli.status.formatter import format_status_human, format_status_json


@click.command()
@click.option(
    "--json",
    "json_output",
    is_flag=True,
    help="Output as JSON",
)
@click.option(
    "-v",
    "--verbose",
    is_flag=True,
    help="Show more details",
)
def status(json_output: bool, verbose: bool) -> None:
    """Show current project status.

    Displays:
    - Workstreams in progress
    - Ready workstreams
    - Blocked workstreams
    - Guard status
    - Beads integration status
    - Suggested next actions
    """
    root = Path.cwd()
    collector = StatusCollector(root)
    project_status = collector.collect()

    if json_output:
        click.echo(format_status_json(project_status))
    else:
        click.echo(format_status_human(project_status, verbose=verbose))
