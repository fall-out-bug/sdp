"""GitHub sync CLI commands."""

from pathlib import Path

import click

from sdp.github.client import GitHubClient
from sdp.github.config import GitHubConfig
from sdp.github.sync_service import SyncResult, SyncService


@click.group()
def github() -> None:
    """GitHub integration commands."""
    pass


@github.command("sync-all")
@click.option(
    "--ws-dir",
    default="tools/hw_checker/docs/workstreams",
    help="Workstreams directory",
    type=str,
)
@click.option(
    "--dry-run",
    is_flag=True,
    help="Show what would be synced without making changes",
)
@click.option(
    "--project",
    default="auto",
    help="GitHub Project name or 'auto' for routing",
)
def sync_all(ws_dir: str, dry_run: bool, project: str) -> None:
    """Sync all workstreams to GitHub.

    Syncs WS files from backlog/, active/, and completed/ directories.
    Creates GitHub issues and adds them to project board.
    """
    ws_path = Path(ws_dir)

    if not ws_path.exists():
        click.echo(f"Error: Directory not found: {ws_path}")
        raise SystemExit(1)

    if dry_run:
        click.echo("=== DRY RUN MODE ===")
        click.echo(f"Would sync from: {ws_path}")
        # Count WS files
        backlog_dir = ws_path / "backlog"
        backlog = list(backlog_dir.glob("WS-*.md")) if backlog_dir.exists() else []
        active_dir = ws_path / "active"
        active = list(active_dir.glob("WS-*.md")) if active_dir.exists() else []
        completed_dir = ws_path / "completed"
        completed = (
            list(completed_dir.rglob("WS-*.md")) if completed_dir.exists() else []
        )
        click.echo(f"  Backlog: {len(backlog)} files")
        click.echo(f"  Active: {len(active)} files")
        click.echo(f"  Completed: {len(completed)} files")
        click.echo("No changes made.")
        return

    # Initialize client and service
    config = GitHubConfig.from_env()
    client = GitHubClient(config)
    service = SyncService(client, project_name=project)

    # Sync backlog
    click.echo("Syncing backlog...")
    backlog_dir = ws_path / "backlog"
    if backlog_dir.exists():
        results = service.sync_all(backlog_dir)
        _print_results(results)

    # Sync active
    click.echo("Syncing active...")
    active_dir = ws_path / "active"
    if active_dir.exists():
        results = service.sync_all(active_dir)
        _print_results(results)

    # Sync completed (includes nested date folders)
    click.echo("Syncing completed...")
    completed_dir = ws_path / "completed"
    if completed_dir.exists():
        if list(completed_dir.glob("WS-*.md")):
            results = service.sync_all(completed_dir)
            _print_results(results)
        for subdir in sorted(
            path for path in completed_dir.iterdir() if path.is_dir()
        ):
            results = service.sync_all(subdir)
            _print_results(results)

    click.echo("Done!")


@github.command("sync-ws")
@click.argument("ws_file", type=click.Path(exists=True))
@click.option(
    "--project",
    default="auto",
    help="GitHub Project name or 'auto' for routing",
)
def sync_ws(ws_file: str, project: str) -> None:
    """Sync single workstream file to GitHub."""
    ws_path = Path(ws_file)

    config = GitHubConfig.from_env()
    client = GitHubClient(config)
    service = SyncService(client, project_name=project)

    result = service.sync_workstream(ws_path)

    if result.action == "failed":
        click.echo(f"Error: {result.error}")
        raise SystemExit(1)

    click.echo(f"{result.ws_id}: {result.action} (#{result.issue_number})")


def _print_results(results: list[SyncResult]) -> None:
    """Print sync results."""
    for result in results:
        if result.action == "failed":
            click.echo(f"  {result.ws_id}: FAILED - {result.error}")
        else:
            click.echo(f"  {result.ws_id}: {result.action} (#{result.issue_number})")


if __name__ == "__main__":
    github()
