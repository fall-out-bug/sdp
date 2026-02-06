"""File operations for supersede validation."""

from pathlib import Path


def find_ws_file(ws_id: str, ws_dir: Path) -> Path | None:
    """Find WS file by ID.

    Args:
        ws_id: Workstream ID
        ws_dir: Base workstream directory

    Returns:
        Path to WS file or None
    """
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


def parse_frontmatter(ws_path: Path) -> dict[str, str]:
    """Parse frontmatter from WS file.

    Args:
        ws_path: Path to WS file

    Returns:
        Dict of frontmatter fields
    """
    content = ws_path.read_text(encoding="utf-8")
    data: dict[str, str] = {}

    in_frontmatter = False
    for line in content.splitlines():
        if line.strip() == "---":
            if not in_frontmatter:
                in_frontmatter = True
            else:
                break
            continue

        if in_frontmatter and ":" in line:
            key, value = line.split(":", 1)
            data[key.strip()] = value.strip()

    return data


def update_frontmatter(ws_path: Path, updates: dict[str, str]) -> None:
    """Update frontmatter fields.

    Args:
        ws_path: Path to WS file
        updates: Dict of fields to update

    Raises:
        ValueError: If frontmatter not found
    """
    content = ws_path.read_text(encoding="utf-8")
    lines = content.splitlines()

    # Find frontmatter bounds
    fm_start = -1
    fm_end = -1
    for i, line in enumerate(lines):
        if line.strip() == "---":
            if fm_start == -1:
                fm_start = i
            else:
                fm_end = i
                break

    if fm_start == -1 or fm_end == -1:
        raise ValueError("Frontmatter not found")

    # Update frontmatter
    new_fm_lines: list[str] = []
    updated_keys: set[str] = set()

    for i in range(fm_start + 1, fm_end):
        line = lines[i]
        if ":" in line:
            key, value = line.split(":", 1)
            key = key.strip()

            if key in updates:
                new_fm_lines.append(f"{key}: {updates[key]}")
                updated_keys.add(key)
            else:
                new_fm_lines.append(line)
        else:
            new_fm_lines.append(line)

    # Add new fields
    for key, value in updates.items():
        if key not in updated_keys:
            new_fm_lines.append(f"{key}: {value}")

    # Reconstruct file
    new_lines = (
        lines[:fm_start] + ["---"] + new_fm_lines + ["---"] + lines[fm_end + 1 :]
    )

    ws_path.write_text("\n".join(new_lines) + "\n", encoding="utf-8")
