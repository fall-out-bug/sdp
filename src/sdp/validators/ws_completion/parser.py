"""Workstream file parsing utilities."""

from pathlib import Path
from typing import Any


def find_ws_file(ws_id: str, ws_dir: Path = Path("docs/workstreams")) -> Path | None:
    """Find WS file by ID.

    Args:
        ws_id: Workstream ID
        ws_dir: Base directory for workstream files

    Returns:
        Path to WS file or None
    """
    # Search in multiple locations
    search_dirs = [
        ws_dir / "backlog",
        ws_dir / "in_progress",
        ws_dir / "completed",
    ]

    for search_dir in search_dirs:
        if not search_dir.exists():
            continue

        for ws_file in search_dir.glob(f"{ws_id}*.md"):
            return ws_file

    return None


def parse_frontmatter_scope(content: str) -> list[str]:
    """Extract scope_files from frontmatter.

    Args:
        content: Full markdown file content

    Returns:
        List of scope file paths
    """
    scope_files: list[str] = []
    in_frontmatter = False
    in_scope_files = False
    for line in content.splitlines():
        if line.strip() == "---":
            if not in_frontmatter:
                in_frontmatter = True
            else:
                break
            continue
        if in_frontmatter:
            if line.startswith("scope_files:"):
                in_scope_files = True
            elif in_scope_files:
                if line.startswith("  - "):
                    scope_files.append(line.strip()[2:].strip())
                elif not line.startswith(" "):
                    in_scope_files = False
    return scope_files


def parse_verification_commands(content: str) -> list[str]:
    """Extract verification commands from ### Verification section.

    Args:
        content: Full markdown file content

    Returns:
        List of shell commands
    """
    commands: list[str] = []
    in_verification = False
    in_code_block = False
    for line in content.splitlines():
        if line.startswith("### Verification"):
            in_verification = True
            continue
        if in_verification:
            if line.strip().startswith("```bash") or line.strip().startswith("```sh"):
                in_code_block = True
                continue
            if line.strip() == "```":
                in_code_block = False
                continue
            if line.startswith("##"):
                break
            if in_code_block and line.strip() and not line.strip().startswith("#"):
                commands.append(line.strip())
    return commands


def parse_ws_file(ws_path: Path) -> dict[str, Any]:
    """Parse WS file for verification data.

    Args:
        ws_path: Path to WS file

    Returns:
        Dict with scope_files, verification_commands, etc.
    """
    content = ws_path.read_text(encoding="utf-8")
    return {
        "scope_files": parse_frontmatter_scope(content),
        "verification_commands": parse_verification_commands(content),
    }
