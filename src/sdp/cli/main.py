"""Main CLI entry point for SDP package."""

from typing import TYPE_CHECKING

import click

from sdp import __version__

if TYPE_CHECKING:
    from click import Command, Group

# Import Beads commands (optional - may not be available in all builds)
_beads_available: bool = False
try:
    from sdp.cli.beads import beads
    _beads_available = True
except ImportError:
    _beads_available = False

# Import doctor command
_doctor_available: bool = False
try:
    from sdp.doctor import doctor
    _doctor_available = True
except ImportError:
    _doctor_available = False

# Import command groups
try:
    from sdp.cli.workstream import workstream
except ImportError:
    # Fallback if workstream module doesn't exist yet
    @click.group()
    def workstream() -> None:
        """Core SDP operations (workstreams, features, project maps)."""
        pass

tier: Group | None = None
try:
    from sdp.cli.tier import tier
except ImportError:
    tier = None

metrics: Group | None = None
metrics_escalations: Command | None = None
try:
    from sdp.cli.metrics import metrics, metrics_escalations
except ImportError:
    metrics = None
    metrics_escalations = None

prd: Group | None = None
prd_detect_type: Command | None = None
prd_validate: Command | None = None
try:
    from sdp.cli.prd import prd, prd_detect_type, prd_validate
except ImportError:
    prd = None
    prd_detect_type = None
    prd_validate = None

# Import extension commands
_extension_available: bool = False
try:
    from sdp.cli_extension import extension  # noqa: E402
    _extension_available = True
except ImportError:
    _extension_available = False

_init_available: bool = False
try:
    from sdp.cli_init import init  # noqa: E402
    _init_available = True
except ImportError:
    _init_available = False


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

# Add workstream commands
main.add_command(workstream)

# Add tier commands (tier group includes metrics and promote-check)
if tier:
    main.add_command(tier)

# Add metrics commands
if metrics:
    main.add_command(metrics)
if metrics_escalations:
    main.add_command(metrics_escalations)

# Add PRD commands
if prd:
    main.add_command(prd)
if prd_detect_type:
    main.add_command(prd_detect_type)
if prd_validate:
    main.add_command(prd_validate)

# Add extension commands
if _extension_available:
    main.add_command(extension)

# Add init command
if _init_available:
    main.add_command(init)


@main.command()
def version() -> None:
    """Show SDP version."""
    click.echo(f"sdp version {__version__}")


if __name__ == "__main__":
    main()
