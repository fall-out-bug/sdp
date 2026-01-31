"""Shared utilities for Git hooks."""

from dataclasses import dataclass
from pathlib import Path
from typing import cast

import tomllib


def find_project_root(start_dir: Path | None = None) -> Path:
    """Find project root by looking for markers.

    Markers (in order of precedence):
    1. .sdp-root file (explicit root marker)
    2. docs/workstreams/ directory
    3. .git directory + pyproject.toml with [tool.sdp] section

    Args:
        start_dir: Directory to start search from (default: cwd)

    Returns:
        Path to project root directory

    Raises:
        RuntimeError: If project root cannot be found
    """
    start = start_dir or Path.cwd()
    if not start.is_absolute():
        start = start.resolve()

    for path in [start, *list(start.parents)]:
        if (path / ".sdp-root").exists():
            return path
        if (path / "docs" / "workstreams").exists():
            return path
        if (path / ".git").exists() and (path / "pyproject.toml").exists():
            try:
                config = tomllib.loads((path / "pyproject.toml").read_text())
                if "tool" in config and "sdp" in config.get("tool", {}):
                    return path
            except (tomllib.TOMLDecodeError, OSError):
                pass  # TOML parse or read failed, try next marker

    raise RuntimeError(
        "SDP project root not found. "
        "Initialize with: sdp init or create .sdp-root marker"
    )


def find_workstream_dir(project_root: Path) -> Path:
    """Find workstream directory with fallbacks.

    Search order:
    1. quality-gate.toml [workstreams.dir] config
    2. SDP_WORKSTREAM_DIR environment variable
    3. docs/workstreams/ (default)
    4. workstreams/ (legacy fallback)

    Args:
        project_root: Path to project root

    Returns:
        Path to workstream directory

    Raises:
        RuntimeError: If no workstream directory found
    """
    import os

    config_file = project_root / "quality-gate.toml"
    if config_file.exists():
        try:
            config = tomllib.loads(config_file.read_text())
            ws_config = config.get("workstreams", {})
            if "dir" in ws_config:
                ws_dir = project_root / cast(str, ws_config["dir"])
                if ws_dir.exists():
                    return ws_dir
        except (tomllib.TOMLDecodeError, OSError):
            pass  # Config parse failed, try env or defaults

    env_ws_dir = os.getenv("SDP_WORKSTREAM_DIR")
    if env_ws_dir:
        ws_dir = Path(env_ws_dir)
        if ws_dir.is_absolute() and ws_dir.exists():
            return ws_dir
        abs_ws = (project_root / env_ws_dir).resolve()
        if abs_ws.exists():
            return abs_ws

    default_ws = project_root / "docs" / "workstreams"
    if default_ws.exists():
        return default_ws

    legacy_ws = project_root / "workstreams"
    if legacy_ws.exists():
        return legacy_ws

    raise RuntimeError(
        f"Workstream directory not found in {project_root}. "
        "Create docs/workstreams/ or configure in quality-gate.toml"
    )


@dataclass
class CheckResult:
    """Result of a quality check."""

    passed: bool
    message: str
    violations: list[tuple[Path, int | None, str]]  # (file, line, issue)

    def format_terminal(self) -> str:
        """Format result for terminal output."""
        if self.passed:
            return f"✓ {self.message}"
        output = [f"❌ {self.message}"]
        for file_path, line, issue in self.violations:
            line_str = str(line) if line is not None else "?"
            output.append(f"  {file_path}:{line_str} - {issue}")
        return "\n".join(output)
