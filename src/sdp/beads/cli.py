"""Real Beads client using CLI subprocess calls."""

import json
import subprocess
from pathlib import Path
from typing import List, Optional

from .base import BeadsClient
from .exceptions import BeadsClientError
from .models import BeadsStatus, BeadsTask, BeadsTaskCreate


class CLIBeadsClient(BeadsClient):
    """Real Beads client using CLI subprocess calls.

    Requires:
    - Go 1.24+ installed
    - Beads installed: `go install github.com/steveyegge/beads/cmd/bd@latest`
    - Beads initialized: `bd init` in project directory
    """

    def __init__(self, project_dir: Optional[Path] = None):
        """Initialize CLI client.

        Args:
            project_dir: Project directory (defaults to current dir)
        """
        self.project_dir = project_dir or Path.cwd()

        # Verify Beads is available
        try:
            subprocess.run(
                ["bd", "--version"],
                capture_output=True,
                cwd=self.project_dir,
                check=True,
            )
        except (subprocess.CalledProcessError, FileNotFoundError) as e:
            raise BeadsClientError(
                "Beads CLI not found. Install with: "
                "go install github.com/steveyegge/beads/cmd/bd@latest"
            ) from e

    def create_task(self, params: BeadsTaskCreate) -> BeadsTask:
        """Create task via Beads CLI.

        Beads expects: bd create [title] --description=... --priority=...
        The --json flag is for OUTPUT format, not input.
        """
        # Title as positional arg (Beads limit: 500 chars)
        title = (params.title or "Untitled")[:500]

        cmd = ["bd", "create", title, "--json"]

        if params.description:
            cmd.extend(["--description", params.description])
        if params.priority is not None:
            cmd.extend(["--priority", str(params.priority.value)])
        if params.parent_id:
            cmd.extend(["--parent", params.parent_id])
        if params.external_ref:
            cmd.extend(["--external-ref", params.external_ref])
        if params.dependencies:
            deps_str = ",".join(f"blocks:{d.task_id}" for d in params.dependencies)
            cmd.extend(["--deps", deps_str])

        result = self._run_command(cmd, capture_output=True)
        return BeadsTask.from_dict(json.loads(result.stdout))

    def get_task(self, task_id: str) -> Optional[BeadsTask]:
        """Get task via Beads CLI.

        Example:
            bd show --json bd-a3f8
        """
        cmd = ["bd", "show", "--json", task_id]

        try:
            result = self._run_command(cmd, capture_output=True)
            data = json.loads(result.stdout)
            # bd show --json returns array with one element
            if isinstance(data, list) and len(data) > 0:
                return BeadsTask.from_dict(data[0])
            else:
                return BeadsTask.from_dict(data)

        except (BeadsClientError, json.JSONDecodeError):
            # Task not found or invalid response
            return None

    def update_task_status(self, task_id: str, status: BeadsStatus) -> None:
        """Update task status via Beads CLI.

        Example:
            bd update bd-a3f8 --status in_progress
        """
        cmd = ["bd", "update", task_id, "--status", status.value]
        self._run_command(cmd)

    def get_ready_tasks(self) -> List[str]:
        """Get ready tasks via Beads CLI.

        Example:
            bd ready --json
        """
        cmd = ["bd", "ready", "--json"]
        result = self._run_command(cmd, capture_output=True)

        data = json.loads(result.stdout)
        # bd ready --json returns array directly
        if isinstance(data, list):
            return [str(item) for item in data]
        else:
            return [str(item) for item in data.get("ready_tasks", [])]

    def add_dependency(
        self, from_id: str, to_id: str, dep_type: str = "blocks"
    ) -> None:
        """Add dependency via Beads CLI.

        Example:
            bd dep add bd-a3f8.1 bd-a3f8 --type blocks
        """
        cmd = ["bd", "dep", "add", from_id, to_id, "--type", dep_type]
        self._run_command(cmd)

    def list_tasks(
        self,
        status: Optional[BeadsStatus] = None,
        parent_id: Optional[str] = None,
    ) -> List[BeadsTask]:
        """List tasks via Beads CLI.

        Example:
            bd list --status open --json
        """
        cmd = ["bd", "list", "--json"]

        if status:
            cmd.extend(["--status", status.value])

        if parent_id:
            cmd.extend(["--parent", parent_id])

        result = self._run_command(cmd, capture_output=True)
        data = json.loads(result.stdout)
        # bd list --json returns array directly
        if isinstance(data, list):
            return [BeadsTask.from_dict(t) for t in data]
        else:
            return [BeadsTask.from_dict(t) for t in data.get("tasks", [])]

    def update_metadata(self, task_id: str, metadata: dict) -> None:
        """Update task metadata via Beads CLI.

        Example:
            bd update bd-a3f8 --metadata '{"sdp": {...}}'
        """
        cmd = ["bd", "update", task_id, "--metadata", json.dumps(metadata)]
        self._run_command(cmd)

    def _run_command(
        self, cmd: List[str], capture_output: bool = False
    ) -> subprocess.CompletedProcess:
        """Run a Beads CLI command.

        Args:
            cmd: Command and arguments
            capture_output: Whether to capture stdout/stderr

        Returns:
            Completed process result

        Raises:
            BeadsClientError: If command fails
        """
        try:
            result = subprocess.run(
                cmd,
                capture_output=True,
                text=capture_output,
                cwd=self.project_dir,
                check=True,
            )
            return result

        except subprocess.CalledProcessError as e:
            error_msg = e.stderr if capture_output else str(e)
            raise BeadsClientError(f"Command failed: {error_msg}") from e

        except json.JSONDecodeError as e:
            error_msg = e.stderr if capture_output else str(e)
            raise BeadsClientError(f"Invalid JSON response: {error_msg}") from e
