"""Main CLI entry point for SDP package."""

import os
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

    from sdp.validators import validate_workstream_tier

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


@main.group()
def daemon() -> None:
    """Daemon management commands."""
    pass


@daemon.command("start")
@click.option("--watch", is_flag=True, help="Enable file watch mode")
@click.option(
    "--pid-file",
    default=".sdp/daemon.pid",
    help="PID file path",
)
def daemon_start(watch: bool, pid_file: str) -> None:
    """Start SDP daemon process."""
    import sys

    from sdp.daemon.daemon import Daemon, DaemonConfig
    from sdp.daemon.pid_manager import PIDManager

    manager = PIDManager(pid_file)
    if manager.is_running():
        click.echo(f"Daemon already running (PID: {manager.read()})")
        return

    # Double-fork method for Unix-like systems
    try:
        pid = _fork_daemon()
    except OSError as e:
        click.echo(f"Failed to fork daemon: {e}", err=True)
        sys.exit(1)

    if pid == 0:
        # Child process (daemon)
        from sdp.daemon.daemon import Daemon, DaemonConfig

        config = DaemonConfig(watch_enabled=watch, pid_file=pid_file)
        daemon = Daemon(config)
        daemon.run()
        sys.exit(0)
    else:
        # Parent process
        manager.write(pid)
        click.echo(f"Daemon started (PID: {pid})")


@daemon.command("stop")
@click.option(
    "--pid-file",
    default=".sdp/daemon.pid",
    help="PID file path",
)
def daemon_stop(pid_file: str) -> None:
    """Stop SDP daemon process."""
    import signal

    from sdp.daemon.pid_manager import PIDManager

    manager = PIDManager(pid_file)
    if not manager.is_running():
        click.echo("Daemon not running")
        return

    pid = manager.read()
    try:
        os.kill(pid, signal.SIGTERM)
        manager.remove()
        click.echo(f"Daemon stopped (PID: {pid})")
    except OSError as e:
        click.echo(f"Failed to stop daemon: {e}", err=True)
        sys.exit(1)


@daemon.command("status")
@click.option(
    "--pid-file",
    default=".sdp/daemon.pid",
    help="PID file path",
)
def daemon_status(pid_file: str) -> None:
    """Show daemon status."""
    from sdp.daemon.pid_manager import PIDManager

    manager = PIDManager(pid_file)
    if manager.is_running():
        pid = manager.read()
        click.echo(f"Daemon running (PID: {pid})")
    else:
        click.echo("Daemon not running")


def _fork_daemon() -> int:
    """Double-fork to create a daemon process.

    Returns:
        PID of the daemon process (in parent), 0 in child
    """
    import os

    # First fork
    pid = os.fork()
    if pid > 0:
        # Parent process
        return pid

    # Child process continues
    os.setsid()

    # Second fork
    pid = os.fork()
    if pid > 0:
        # First child exits
        os._exit(0)

    # Daemon process
    return 0


@main.group()
def queue() -> None:
    """Task queue management commands."""
    pass


@queue.command("enqueue")
@click.argument("ws_id")
@click.option(
    "--priority",
    type=click.Choice(["blocked", "backlog", "normal", "active", "urgent"], case_sensitive=False),
    default="normal",
    help="Task priority",
)
def queue_enqueue(ws_id: str, priority: str) -> None:
    """Add workstream to task queue."""
    from sdp.queue import Priority, Task, TaskQueue

    q = TaskQueue()
    task = Task(ws_id=ws_id, priority=Priority.from_string(priority))
    q.enqueue(task)
    click.echo(f"Enqueued: {ws_id} (priority: {task.priority.name})")


@queue.command("dequeue")
@click.option("--wait", type=float, default=0, help="Seconds to wait for task")
def queue_dequeue(wait: float) -> None:
    """Remove and return highest priority task."""
    from sdp.queue import TaskQueue

    q = TaskQueue()
    task = q.dequeue(timeout=wait if wait > 0 else None)
    if task:
        click.echo(f"Dequeued: {task.ws_id} (priority: {task.priority.name})")
    else:
        click.echo("Queue empty")


@queue.command("list")
def queue_list() -> None:
    """Show all tasks in queue."""
    from sdp.queue import TaskQueue

    q = TaskQueue()

    if q.is_empty():
        click.echo("Queue empty")
        return

    click.echo(f"Queue size: {q.size()}")
    click.echo()

    tasks = []
    while True:
        task = q.dequeue()
        if task is None:
            break
        tasks.append(task)
        click.echo(f"  {task.ws_id} - {task.priority.name}")

    # Put tasks back
    for task in tasks:
        q.enqueue(task)


@queue.command("clear")
def queue_clear() -> None:
    """Clear all tasks from queue."""
    from sdp.queue import TaskQueue

    q = TaskQueue()
    size = q.size()
    q.clear()
    click.echo(f"Cleared {size} task(s)")


@queue.command("status")
def queue_status() -> None:
    """Show queue status."""
    from sdp.queue import TaskQueue

    q = TaskQueue()
    click.echo(f"Queue size: {q.size()}")

    next_task = q.peek()
    if next_task:
        click.echo(f"Next: {next_task.ws_id} (priority: {next_task.priority.name})")
    else:
        click.echo("Queue empty")


# Test command group
@main.group()
def test() -> None:
    """Test execution commands."""
    pass


@test.command("run")
@click.option("--coverage", is_flag=True, help="Enable coverage report")
@click.option("--pattern", help="Filter tests by pattern")
def test_run(coverage: bool, pattern: str | None) -> None:
    """Run tests once."""
    from sdp.test_watch.runner import WatchTestRunner

    runner = WatchTestRunner(".", coverage=coverage, pattern=pattern)
    results = runner.run()

    status_emoji = {"passed": "✅", "failed": "❌", "error": "⚠️", "no_tests": "⚪"}
    click.echo(f"{status_emoji.get(results.status, '?')} Test Results")
    click.echo(f"Status: {results.status.upper()}")

    if results.status != "no_tests":
        click.echo(f"Passed: {results.passed}")
        click.echo(f"Failed: {results.failed}")
        if results.coverage:
            click.echo(f"Coverage: {results.coverage:.1f}%")

    if results.error_message:
        click.echo(f"Error: {results.error_message}", err=True)

    sys.exit(0 if results.status in ["passed", "no_tests"] else 1)


@test.command("watch")
@click.option("--coverage", is_flag=True, default=True, help="Enable coverage report")
@click.option("--pattern", help="Filter tests by pattern")
@click.option("--debounce", type=float, default=0.5, help="Debounce delay in seconds")
def test_watch(coverage: bool, pattern: str | None, debounce: float) -> None:
    """Watch for file changes and run tests automatically."""
    import time

    from sdp.test_watch.runner import WatchTestRunner
    from sdp.test_watch.watcher import watch_tests

    runner = WatchTestRunner(".", coverage=coverage, pattern=pattern)

    def on_change(changed_file: str) -> None:
        click.echo(f"\n📝 {changed_file} changed")
        click.echo("Running tests...")

        results = runner.run_affected(changed_file)

        status_emoji = {"passed": "✅", "failed": "❌", "error": "⚠️", "no_tests": "⚪"}
        click.echo(f"{status_emoji.get(results.status, '?')} {results.status.upper()}")

        if results.failed_tests:
            click.echo(f"Failed: {results.failed}")

    observer = watch_tests(".", on_change, debounce=debounce)

    click.echo("Watching for file changes... (Ctrl+C to stop)")
    try:
        # Run initial tests
        click.echo("Running initial tests...")
        results = runner.run()
        status_emoji = {"passed": "✅", "failed": "❌", "error": "⚠️", "no_tests": "⚪"}
        click.echo(f"{status_emoji.get(results.status, '?')} Initial: {results.status.upper()}")

        # Keep watcher alive
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        observer.stop()
        observer.join()
        click.echo("\nWatch mode stopped")



# Task commands (separate from queue)
@main.group()
def task() -> None:
    """Task execution and management commands."""
    pass


@task.command("enqueue")
@click.argument("ws_id")
@click.option(
    "--priority",
    type=click.Choice(["blocked", "backlog", "normal", "active", "urgent"], case_sensitive=False),
    default="normal",
    help="Task priority",
)
def task_enqueue(ws_id: str, priority: str) -> None:
    """Add workstream to task queue."""
    from sdp.queue import Priority, Task, TaskQueue

    q = TaskQueue()
    task = Task(ws_id=ws_id, priority=Priority.from_string(priority))
    q.enqueue(task)
    click.echo(f"Enqueued: {ws_id} (priority: {task.priority.name})")


@task.command("execute")
@click.argument("ws_id")
@click.option("--timeout", type=int, default=3600, help="Timeout in seconds")
@click.option("--dry-run", is_flag=True, help="Show what would be done without executing")
def task_execute(ws_id: str, timeout: int, dry_run: bool) -> None:
    """Execute workstream immediately."""
    from sdp.agents import AgentExecutor

    if dry_run:
        click.echo(f"Would execute: {ws_id}")
        click.echo(f"  Timeout: {timeout}s")
        return

    executor = AgentExecutor(timeout=timeout)
    click.echo(f"Executing {ws_id}...")

    try:
        success = executor.execute(ws_id, timeout=timeout)
        if success:
            click.echo(f"✅ {ws_id} completed successfully")
        else:
            click.echo(f"❌ {ws_id} failed")
            sys.exit(1)
    except Exception as e:
        click.echo(f"❌ Error: {e}", err=True)
        sys.exit(1)


@task.command("list")
def task_list() -> None:
    """Show all pending/running/completed tasks."""
    from sdp.queue import TaskQueue

    q = TaskQueue()

    if q.is_empty():
        click.echo("No tasks in queue")
        return

    click.echo(f"Queue size: {q.size()}")
    task = q.peek()
    if task:
        click.echo(f"Next: {task.ws_id} (priority: {task.priority.name})")


@task.command("cancel")
@click.argument("task_id", type=int)
def task_cancel(task_id: int) -> None:
    """Cancel pending task (by queue position)."""
    from sdp.queue import TaskQueue

    q = TaskQueue()

    if q.is_empty():
        click.echo("Queue is empty")
        return

    if task_id == 1:
        # Cancel first task
        task = q.dequeue()
        if task:
            click.echo(f"Cancelled: {task.ws_id}")
    else:
        click.echo(f"Cancelling specific positions not yet implemented")


# Workspace commands
@main.group()
def ws() -> None:
    """Workstream state management commands."""
    pass


@ws.command("move")
@click.argument("ws_id")
@click.option("--to", "to_status", required=True, type=click.Choice(["backlog", "in-progress", "completed", "blocked"]), help="Target status")
@click.option("--no-index", is_flag=True, help="Skip updating INDEX.md")
def ws_move(ws_id: str, to_status: str, no_index: bool) -> None:
    """Move workstream to different status directory."""
    from sdp.workspace.mover import WorkstreamMover
    from sdp.workspace.validator import MoveValidationError

    mover = WorkstreamMover()

    try:
        new_path = mover.move(ws_id, to_status, update_index=not no_index)
        click.echo(f"Moved {ws_id} to {to_status}")
        click.echo(f"  {new_path}")
    except MoveValidationError as e:
        click.echo(f"Error: {e}", err=True)
        sys.exit(1)


@ws.command("start")
@click.argument("ws_id")
def ws_start(ws_id: str) -> None:
    """Start workstream (backlog -> in-progress)."""
    from sdp.workspace.mover import WorkstreamMover
    from sdp.workspace.validator import MoveValidationError

    mover = WorkstreamMover()

    try:
        new_path = mover.start(ws_id)
        click.echo(f"Started {ws_id}")
        click.echo(f"  {new_path}")
    except MoveValidationError as e:
        click.echo(f"Error: {e}", err=True)
        sys.exit(1)


@ws.command("complete")
@click.argument("ws_id")
@click.option("--no-validate", is_flag=True, help="Skip acceptance criteria validation")
def ws_complete(ws_id: str, no_validate: bool) -> None:
    """Complete workstream (in-progress -> completed)."""
    from sdp.workspace.mover import WorkstreamMover
    from sdp.workspace.validator import MoveValidationError

    mover = WorkstreamMover()

    try:
        new_path = mover.complete(ws_id, validate=not no_validate)
        click.echo(f"Completed {ws_id}")
        click.echo(f"  {new_path}")
    except MoveValidationError as e:
        click.echo(f"Error: {e}", err=True)
        sys.exit(1)


@ws.command("block")
@click.argument("ws_id")
def ws_block(ws_id: str) -> None:
    """Block workstream."""
    from sdp.workspace.mover import WorkstreamMover
    from sdp.workspace.validator import MoveValidationError

    mover = WorkstreamMover()

    try:
        new_path = mover.block(ws_id)
        click.echo(f"Blocked {ws_id}")
        click.echo(f"  {new_path}")
    except MoveValidationError as e:
        click.echo(f"Error: {e}", err=True)
        sys.exit(1)


@ws.command("list")
@click.option("--status", type=click.Choice(["backlog", "in-progress", "completed", "blocked"]), help="Filter by status")
def ws_list(status: str | None) -> None:
    """List workstreams in status directory."""
    from sdp.workspace.mover import WorkstreamMover

    mover = WorkstreamMover()

    if status:
        workstreams = mover.list_in_status(status)
        click.echo(f"{status} ({len(workstreams)})")
        for ws in sorted(workstreams, key=lambda p: p.name):
            click.echo(f"  {ws.name}")
    else:
        total = 0
        for s in ["backlog", "in_progress", "completed", "blocked"]:
            workstreams = mover.list_in_status(s)
            if workstreams:
                click.echo(f"\n{s.replace('_', '-')} ({len(workstreams)})")
                for ws in sorted(workstreams, key=lambda p: p.name):
                    click.echo(f"  {ws.name}")
                total += len(workstreams)

        click.echo(f"\nTotal: {total}")


# Register extension commands
from sdp.cli_extension import extension
from sdp.cli_init import init
from sdp.status.command import status

main.add_command(extension)
main.add_command(init)
main.add_command(prd)
main.add_command(daemon)
main.add_command(queue)
main.add_command(status)
main.add_command(test)
main.add_command(task)
main.add_command(ws)
main.add_command(orchestrator)


@main.command()
def dashboard() -> None:
    """Launch the SDP Dashboard TUI application."""
    from sdp.dashboard.dashboard_app import DashboardApp

    try:
        app = DashboardApp()
        app.run()
    except Exception as e:
        click.echo(f"Error starting dashboard: {e}", err=True)
        sys.exit(1)


# Orchestrator commands
@main.group()
def orchestrator() -> None:
    """Multi-agent orchestration commands."""
    pass


@orchestrator.command("run")
@click.argument("feature_id")
@click.option("--max-agents", default=3, type=int, help="Maximum concurrent agents")
@click.option("--ordered", is_flag=True, help="Run sequentially (for debugging)")
@click.option("--ws-dir", default="docs/workstreams", help="Workstreams directory")
def orchestrator_run(feature_id: str, max_agents: int, ordered: bool, ws_dir: str) -> None:
    """Execute all workstreams for a feature.

    Example: sdp orchestrator run F012
    """
    import asyncio

    from sdp.agents.orchestrator import Orchestrator

    async def run() -> None:
        orch = Orchestrator(max_agents=max_agents, ws_dir=ws_dir)

        # Progress callback
        def on_progress(ws_id: str, success: bool, error: str | None) -> None:
            status = click.style("✓", fg="green") if success else click.style("✗", fg="red")
            click.echo(f"{status} {ws_id}" + (f": {error}" if error else ""))

        orch.on_progress(on_progress)

        click.echo(f"Executing feature {feature_id}...")

        if ordered:
            result = await orch.run_feature_ordered(feature_id)
        else:
            result = await orch.run_feature(feature_id)

        # Summary
        click.echo(f"\n{'='*50}")
        click.echo(f"Feature: {result.feature_id}")
        click.echo(f"Status: {click.style('SUCCESS', fg='green') if result.success else click.style('FAILED', fg='red')}")
        click.echo(f"Duration: {result.duration_seconds:.1f}s")
        click.echo(f"Completed: {len(result.completed)}")
        click.echo(f"Failed: {len(result.failed)}")

        if result.failed:
            click.echo(f"\nFailed workstreams:")
            for ws_id in result.failed:
                error = result.errors.get(ws_id, "Unknown error")
                click.echo(f"  {ws_id}: {error}")

    try:
        asyncio.run(run())
    except ValueError as e:
        click.echo(f"Error: {e}", err=True)
        sys.exit(1)
    except Exception as e:
        click.echo(f"Unexpected error: {e}", err=True)
        sys.exit(1)


@orchestrator.command("enqueue")
@click.argument("feature_id")
@click.option(
    "--priority",
    type=click.Choice(["blocked", "backlog", "normal", "active", "urgent"], case_sensitive=False),
    default="normal",
    help="Task priority",
)
@click.option("--ws-dir", default="docs/workstreams", help="Workstreams directory")
def orchestrator_enqueue(feature_id: str, priority: str, ws_dir: str) -> None:
    """Enqueue all workstreams for a feature in the task queue.

    Example: sdp orchestrator enqueue F012 --priority urgent
    """
    import asyncio

    from sdp.agents.orchestrator import Orchestrator
    from sdp.queue.task_queue import Priority

    async def run() -> None:
        orch = Orchestrator(ws_dir=ws_dir)

        priority_map = {
            "blocked": Priority.BLOCKED,
            "backlog": Priority.BACKLOG,
            "normal": Priority.NORMAL,
            "active": Priority.ACTIVE,
            "urgent": Priority.URGENT,
        }

        task_ids = await orch.enqueue_feature(feature_id, priority_map[priority])
        click.echo(f"Enqueued {len(task_ids)} tasks for feature {feature_id}")

    asyncio.run(run())


@orchestrator.command("status")
@click.option("--ws-dir", default="docs/workstreams", help="Workstreams directory")
def orchestrator_status(ws_dir: str) -> None:
    """Show orchestrator state and agent pool status."""
    import asyncio

    from sdp.agents.orchestrator import Orchestrator

    orch = Orchestrator(ws_dir=ws_dir)

    # Load saved state
    state = orch.load_state()
    if state:
        click.echo(f"Last execution:")
        click.echo(f"  Feature: {state.feature_id}")
        click.echo(f"  Status: {state.status}")
        click.echo(f"  Started: {state.started_at}")
        click.echo(f"  Current: {state.current_ws or 'None'}")
        click.echo(f"  Completed: {len(state.results)}/{len(state.results) + len(state.errors)}")
    else:
        click.echo("No previous execution state found.")

    # Pool stats
    stats = orch.get_stats()
    click.echo(f"\nAgent Pool:")
    click.echo(f"  Total: {stats.total_agents}")
    click.echo(f"  Busy: {stats.busy_agents}")
    click.echo(f"  Available: {stats.available_agents}")


if __name__ == "__main__":
    main()
