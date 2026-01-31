"""Dependency detection for SDP init wizard."""

import subprocess
from pathlib import Path


def detect_dependencies() -> dict[str, bool]:
    """Detect optional dependencies.

    Returns:
        Dict mapping dependency name to availability
    """
    deps = {
        "Beads CLI": _check_command("beads"),
        "GitHub CLI (gh)": _check_command("gh"),
        "Telegram": _check_telegram(),
    }
    return deps


def _check_command(command: str) -> bool:
    """Check if a command is available.

    Args:
        command: Command to check

    Returns:
        True if command is available
    """
    try:
        result = subprocess.run(
            [command, "--version"],
            capture_output=True,
            timeout=5,
        )
        return result.returncode == 0
    except (FileNotFoundError, subprocess.TimeoutExpired):
        return False


def _check_telegram() -> bool:
    """Check if Telegram is configured.

    Returns:
        True if .env has Telegram configuration
    """
    env_file = Path.cwd() / ".env"
    if not env_file.exists():
        return False

    content = env_file.read_text()
    return "TELEGRAM_BOT_TOKEN" in content and "TELEGRAM_CHAT_ID" in content


def show_dependencies(deps: dict[str, bool]) -> None:
    """Show detected dependencies.

    Args:
        deps: Dict of dependency names to availability
    """
    import click

    click.echo()
    click.echo(click.style("Step 2: Dependency Detection", fg="cyan", bold=True))

    for name, available in deps.items():
        if available:
            status = click.style("✓ Detected", fg="green")
        else:
            status = click.style("⊘ Not found", fg="yellow")
        click.echo(f"  {name}: {status}")
