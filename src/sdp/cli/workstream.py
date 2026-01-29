"""Workstream parsing commands.

Provides CLI commands for parsing workstreams and project maps.
"""

import sys
from pathlib import Path
from typing import TYPE_CHECKING

import click

from sdp.errors import format_error_for_terminal

if TYPE_CHECKING:
    from click import Command

# Import validate_tier if available
validate_tier: Command | None = None
try:
    from sdp.cli.tier import validate_tier
except ImportError:
    validate_tier = None


@click.group()
def workstream() -> None:
    """Core SDP operations (workstreams, features, project maps)."""
    pass


@workstream.command("parse")
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


@workstream.command("parse-project-map")
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


# Add validate_tier command if available (as 'validate' for compatibility)
if validate_tier:
    workstream.add_command(validate_tier, name='validate')
